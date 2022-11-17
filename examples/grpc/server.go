package main

import (
	"context"
	"sync"
	"time"

	"github.com/dipdup-net/indexer-sdk/examples/grpc/pb"
	"github.com/dipdup-net/indexer-sdk/pkg/modules/grpc"
	generalPB "github.com/dipdup-net/indexer-sdk/pkg/modules/grpc/pb"
)

// Server -
type Server struct {
	*grpc.Server
	pb.UnimplementedTimeServiceServer

	subscriptions *grpc.Subscriptions[time.Time, *pb.Response]

	wg *sync.WaitGroup
}

// NewServer -
func NewServer(cfg *grpc.ServerConfig) (*Server, error) {
	server, err := grpc.NewServer(cfg)
	if err != nil {
		return nil, err
	}

	return &Server{
		Server:        server,
		subscriptions: grpc.NewSubscriptions[time.Time, *pb.Response](),
		wg:            new(sync.WaitGroup),
	}, nil
}

// Start -
func (server *Server) Start(ctx context.Context) {
	pb.RegisterTimeServiceServer(server.Server.Server(), server)

	server.Server.Start(ctx)

	server.wg.Add(1)
	go server.listen(ctx)
}

func (server *Server) listen(ctx context.Context) {
	defer server.wg.Done()

	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case t := <-ticker.C:
			server.subscriptions.NotifyAll(t, Response)
		}
	}
}

// SubscribeOnTime -
func (server *Server) SubscribeOnTime(req *pb.Request, stream pb.TimeService_SubscribeOnTimeServer) error {
	return grpc.DefaultSubscribeOn[time.Time, *pb.Response](
		stream,
		server.subscriptions,
		NewTimeSubscription(),
	)
}

// UnsubscribeFromTime -
func (server *Server) UnsubscribeFromTime(ctx context.Context, req *generalPB.UnsubscribeRequest) (*generalPB.UnsubscribeResponse, error) {
	return grpc.DefaultUnsubscribe(ctx, server.subscriptions, req.Id)
}
