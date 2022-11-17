package grpc

import (
	"context"
	"io"
	"net"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	"github.com/dipdup-net/indexer-sdk/pkg/messages"
	"github.com/dipdup-net/indexer-sdk/pkg/modules/grpc/pb"
	"google.golang.org/grpc"
	gogrpc "google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

// Server - basic server structure which implemented module interface and handle stats endpoints.
type Server struct {
	*messages.Subscriber

	bind string

	server *gogrpc.Server

	wg *sync.WaitGroup
}

// NewServer - constructor of server struture
func NewServer(cfg *ServerConfig) (*Server, error) {
	if cfg == nil {
		return nil, errors.New("configuration structure of gRPC server is nil")
	}
	subscriber := messages.NewSubscriber()
	module := &Server{
		bind:       cfg.Bind,
		Subscriber: subscriber,
		wg:         new(sync.WaitGroup),
	}
	module.server = gogrpc.NewServer(
		gogrpc.StreamInterceptor(subscriptionStreamServerInterceptor()),
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
	if err := module.Subscriber.Close(); err != nil {
		return err
	}
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

// DefaultSubscribeOn - default subscribe server handler
func DefaultSubscribeOn[T any, P any](stream ServerStream[P], subscriptions *Subscriptions[T, P], subscription Subscription[T, P]) error {
	ctx := stream.Context()

	subscriptionID, err := GetSubscriptionID(ctx)
	if err != nil {
		return err
	}
	subscriptions.Add(subscriptionID, subscription)

	var end bool
	for !end {
		select {
		case <-ctx.Done():
			end = true
		case msg := <-subscription.Listen():
			if err := stream.Send(msg); err != nil {
				if err == io.EOF {
					end = true
				} else {
					log.Err(err).Msg("sending message error")
				}
			}
		}
	}

	return subscriptions.Remove(subscriptionID)
}

// DefaultUnsubscribe - default unsubscribe server handler
func DefaultUnsubscribe[T any, P any](ctx context.Context, subscriptions *Subscriptions[T, P]) (*pb.UnsubscribeResponse, error) {
	subscriptionID, err := GetSubscriptionID(ctx)
	if err != nil {
		return nil, err
	}
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
