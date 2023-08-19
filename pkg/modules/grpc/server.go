package grpc

import (
	"context"
	"io"
	"net"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"

	"github.com/dipdup-net/indexer-sdk/pkg/modules/grpc/pb"
	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/ratelimit"
	"google.golang.org/grpc"
	gogrpc "google.golang.org/grpc"
	_ "google.golang.org/grpc/encoding/gzip" // Install the gzip compressor
	"google.golang.org/grpc/keepalive"
)

// Server - basic server structure which implemented module interface and handle stats endpoints.
type Server struct {
	bind string

	server        *gogrpc.Server
	metricsServer *http.Server
	srvMetrics    *grpcprom.ServerMetrics

	wg *sync.WaitGroup
}

// NewServer - constructor of server struture
func NewServer(cfg *ServerConfig) (*Server, error) {
	if cfg == nil {
		return nil, errors.New("configuration structure of gRPC server is nil")
	}
	module := &Server{
		bind: cfg.Bind,
		wg:   new(sync.WaitGroup),
	}

	opts := []gogrpc.ServerOption{
		gogrpc.KeepaliveParams(
			keepalive.ServerParameters{
				Time:    20 * time.Second,
				Timeout: 10 * time.Second,
			},
		),
		gogrpc.KeepaliveEnforcementPolicy(
			keepalive.EnforcementPolicy{
				MinTime:             10 * time.Second,
				PermitWithoutStream: true,
			},
		),
	}

	streamInterceptors := make([]gogrpc.StreamServerInterceptor, 0)
	unaryInterceptors := make([]gogrpc.UnaryServerInterceptor, 0)

	if cfg.Log {
		streamInterceptors = append(streamInterceptors, logging.StreamServerInterceptor(logCalls()))
		unaryInterceptors = append(unaryInterceptors, logging.UnaryServerInterceptor(logCalls()))

	}

	if cfg.Metrics {
		module.metricsServer, module.srvMetrics = newMetricsServer("127.0.0.1:6789")
		streamInterceptors = append(streamInterceptors, module.srvMetrics.StreamServerInterceptor())
		unaryInterceptors = append(unaryInterceptors, module.srvMetrics.UnaryServerInterceptor())
	}

	if cfg.RPS > 0 {
		limiter := newSimpleLimiter(cfg.RPS)
		streamInterceptors = append(streamInterceptors, ratelimit.StreamServerInterceptor(limiter))
		unaryInterceptors = append(unaryInterceptors, ratelimit.UnaryServerInterceptor(limiter))
	}

	if len(streamInterceptors) > 0 {
		opts = append(opts,
			gogrpc.ChainStreamInterceptor(streamInterceptors...),
		)
	}
	if len(unaryInterceptors) > 0 {
		opts = append(opts,
			gogrpc.ChainUnaryInterceptor(unaryInterceptors...),
		)
	}

	module.server = gogrpc.NewServer(opts...)
	return module, nil
}

func newMetricsServer(httpAddr string) (*http.Server, *grpcprom.ServerMetrics) {
	httpSrv := &http.Server{Addr: httpAddr}
	m := http.NewServeMux()

	srvMetrics := grpcprom.NewServerMetrics(
		grpcprom.WithServerHandlingTimeHistogram(
			grpcprom.WithHistogramBuckets([]float64{0.001, 0.01, 0.1, 0.3, 0.6, 1, 3, 6, 9, 20, 30, 60, 90, 120}),
		),
	)
	reg := prometheus.NewRegistry()
	reg.MustRegister(srvMetrics)
	m.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	httpSrv.Handler = m
	return httpSrv, srvMetrics
}

// Name -
func (*Server) Name() string {
	return "grpc_server"
}

// Start - starts server module
func (module *Server) Start(ctx context.Context) {
	module.wg.Add(1)
	go module.grpc(ctx)

	module.wg.Add(1)
	go module.startMetricsServer(ctx)
}

func (module *Server) grpc(ctx context.Context) {
	defer module.wg.Done()

	log.Info().Str("bind", module.bind).Msg("running grpc...")

	listener, err := net.Listen("tcp", module.bind)
	if err != nil {
		log.Err(err).Msg("net.Listen")
		return
	}
	if err := module.server.Serve(listener); err != nil {
		log.Err(err).Msg("grpcServer.Serve")
	}
}

func (module *Server) startMetricsServer(ctx context.Context) {
	defer module.wg.Done()

	if module.metricsServer == nil {
		return
	}

	if err := module.metricsServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Err(err).Msg("failed to serve metrics")
		return
	}
	log.Info().Msg("metrics server shutdown")
}

// Close - closes server module
func (module *Server) Close() error {
	if module.metricsServer != nil {
		if err := module.metricsServer.Shutdown(context.Background()); err != nil {
			return err
		}
	}
	module.server.GracefulStop()

	module.wg.Wait()
	return nil
}

// Server - returns current grpc.Server to register handlers
func (module *Server) Server() *gogrpc.Server {
	return module.server
}

// ServerStream -
type ServerStream[T any] interface {
	Send(T) error
	grpc.ServerStream
}

var subscriptionsCounter = new(atomic.Uint64)

// DefaultSubscribeOn - default subscribe server handler
func DefaultSubscribeOn[T any, P any](stream ServerStream[P], subscriptions *Subscriptions[T, P], subscription Subscription[T, P], handler func(id uint64) error, onEndOfSync func(id uint64) error) error {
	subscriptionID := subscriptionsCounter.Add(1)

	if err := stream.SendMsg(&pb.SubscribeResponse{
		Id: subscriptionID,
	}); err != nil {
		return err
	}

	if handler != nil {
		if err := handler(subscriptionID); err != nil {
			return errors.Wrap(err, "synchronization")
		}
	}

	subscriptions.Add(subscriptionID, subscription)

	if onEndOfSync != nil {
		if err := onEndOfSync(subscriptionID); err != nil {
			return errors.Wrap(err, "end of sync handler")
		}
	}

loop:
	for {
		select {
		case <-stream.Context().Done():
			break loop
		case msg, ok := <-subscription.Listen():
			if !ok {
				break loop
			}
			if err := stream.Send(msg); err != nil {
				if err == io.EOF {
					break loop
				} else {
					log.Err(err).Msg("sending message error")
				}
			}
		}
	}

	return subscriptions.Remove(subscriptionID)
}

// DefaultUnsubscribe - default unsubscribe server handler
func DefaultUnsubscribe[T any, P any](ctx context.Context, subscriptions *Subscriptions[T, P], subscriptionID uint64) (*pb.UnsubscribeResponse, error) {
	if err := subscriptions.Remove(subscriptionID); err != nil {
		return nil, err
	}

	return &pb.UnsubscribeResponse{
		Id: subscriptionID,
		Response: &pb.Message{
			Message: SuccessMessage,
		},
	}, nil
}
