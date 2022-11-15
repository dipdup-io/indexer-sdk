package grpc

import (
	"context"
	"net"
	"sync"
	"time"

	"github.com/dipdup-net/indexer-sdk/pkg/modules/grpc/pb"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/stats"

	"github.com/dipdup-net/indexer-sdk/pkg/messages"
	gogrpc "google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

// Server - basic server structure which implemented module interface and handle stats endpoints.
type Server struct {
	*messages.Subscriber
	pb.UnimplementedHelloServiceServer

	bind string

	users    map[string]struct{}
	useresMx sync.RWMutex

	server *gogrpc.Server

	wg *sync.WaitGroup
}

// NewServer - constructor of server struture
func NewServer(cfg *ServerConfig) (*Server, error) {
	if cfg == nil {
		return nil, errors.New("configuration structure of gRPC server is nil")
	}
	subscriber, err := messages.NewSubscriber()
	if err != nil {
		return nil, err
	}
	module := &Server{
		bind:       cfg.Bind,
		Subscriber: subscriber,
		users:      make(map[string]struct{}),
		wg:         new(sync.WaitGroup),
	}
	module.server = gogrpc.NewServer(
		gogrpc.StatsHandler(module),
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

////////////////////////////////////////////////
//////////////    HANDLERS    //////////////////
////////////////////////////////////////////////

// Hello - implementation server handshake endpoint
func (module *Server) Hello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	id := ctx.Value(clientID)
	if id == nil {
		return nil, errors.New("unknown client")
	}

	return &pb.HelloResponse{
		Id: id.(string),
	}, nil
}

////////////////////////////////////////////////
////////////////    STATS    ///////////////////
////////////////////////////////////////////////

// TagRPC -
func (module *Server) TagRPC(ctx context.Context, info *stats.RPCTagInfo) context.Context {
	return ctx
}

// HandleRPC -
func (module *Server) HandleRPC(ctx context.Context, s stats.RPCStats) {}

// TagConn -
func (module *Server) TagConn(ctx context.Context, info *stats.ConnTagInfo) context.Context {
	id, err := randomString(32)
	if err != nil {
		log.Err(err).Msg("invalid random string")
	}
	return context.WithValue(ctx, clientID, id)
}

// HandleConn -
func (module *Server) HandleConn(ctx context.Context, s stats.ConnStats) {
	id := ctx.Value(clientID).(string)

	switch s.(type) {
	case *stats.ConnEnd:
		module.useresMx.Lock()
		{
			delete(module.users, id)
		}
		module.useresMx.Unlock()
	case *stats.ConnBegin:
		module.useresMx.Lock()
		{
			module.users[id] = struct{}{}
		}
		module.useresMx.Unlock()
	}
}

// Server - returns current grpc.Server to register handlers
func (module *Server) Server() *gogrpc.Server {
	return module.server
}
