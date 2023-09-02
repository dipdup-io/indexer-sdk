package main

import (
	"context"
	"time"

	"github.com/dipdup-io/workerpool"
	"github.com/dipdup-net/indexer-sdk/examples/grpc/pb"
	"github.com/dipdup-net/indexer-sdk/pkg/modules/grpc"
	generalPB "github.com/dipdup-net/indexer-sdk/pkg/modules/grpc/pb"
)

// Server -
type Server struct {
	*grpc.Server
	pb.UnimplementedTimeServiceServer

	subscriptions *grpc.Subscriptions[time.Time, *pb.Response]

	g workerpool.Group
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
		g:             workerpool.NewGroup(),
	}, nil
}

// Start -
func (server *Server) Start(ctx context.Context) {
	pb.RegisterTimeServiceServer(server.Server.Server(), server)

	server.Server.Start(ctx)

	server.g.GoCtx(ctx, server.listen)
}

func (server *Server) listen(ctx context.Context) {
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
		nil,
		nil,
	)
}

// UnsubscribeFromTime -
func (server *Server) UnsubscribeFromTime(ctx context.Context, req *generalPB.UnsubscribeRequest) (*generalPB.UnsubscribeResponse, error) {
	return grpc.DefaultUnsubscribe(ctx, server.subscriptions, req.Id)
}
