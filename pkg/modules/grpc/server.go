package grpc

import (
	"context"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	"github.com/dipdup-net/indexer-sdk/pkg/modules/grpc/pb"
	"google.golang.org/grpc"
	gogrpc "google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

// Server - basic server structure which implemented module interface and handle stats endpoints.
type Server struct {
	bind string

	server *gogrpc.Server

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
	module.server = gogrpc.NewServer(
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
		))
	return module, nil
}

// Name -
func (*Server) Name() string {
	return "grpc_server"
}

// Start - starts server module
func (module *Server) Start(ctx context.Context) {
	module.wg.Add(1)
	go module.grpc(ctx)
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

// Close - closes server module
func (module *Server) Close() error {
	module.server.Stop()
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
func DefaultSubscribeOn[T any, P any](stream ServerStream[P], subscriptions *Subscriptions[T, P], subscription Subscription[T, P]) error {
	ctx := stream.Context()

	subscriptionID := subscriptionsCounter.Add(1)
	subscriptions.Add(subscriptionID, subscription)

	if err := stream.SendMsg(&pb.SubscribeResponse{
		Id: subscriptionID,
	}); err != nil {
		return err
	}

loop:
	for {
		select {
		case <-ctx.Done():
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
