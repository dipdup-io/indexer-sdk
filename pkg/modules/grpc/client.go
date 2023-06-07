package grpc

import (
	"context"
	"io"
	"sync"
	"time"

	"github.com/dipdup-net/indexer-sdk/pkg/modules/grpc/pb"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	gogrpc "google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/keepalive"
)

// Client - the structure which is responsible for connection to server
type Client struct {
	conn *gogrpc.ClientConn

	serverAddress string
	reconnect     chan struct{}

	wg *sync.WaitGroup
}

// NewClient - constructor of client structure
func NewClient(server string) *Client {
	return &Client{
		serverAddress: server,
		reconnect:     make(chan struct{}, 1),
		wg:            new(sync.WaitGroup),
	}
}

// Name -
func (client *Client) Name() string {
	return "grpc_client"
}

// Connect - connects to server
func (client *Client) Connect(ctx context.Context, opts ...ConnectOption) error {
	dialOpts := []gogrpc.DialOption{
		gogrpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:    10 * time.Second,
			Timeout: 10 * time.Second,
		}),
	}
	connectOpts := newConnectOptions()
	for i := range opts {
		opts[i](&connectOpts)
	}
	dialOpts = append(dialOpts, gogrpc.WithTransportCredentials(connectOpts.creds))
	dialOpts = append(dialOpts, gogrpc.WithConnectParams(
		gogrpc.ConnectParams{
			MinConnectTimeout: connectOpts.reconnectTimeout,
			Backoff: backoff.Config{
				BaseDelay:  1.0 * time.Second,
				Multiplier: 1.6,
				Jitter:     0.2,
				MaxDelay:   connectOpts.reconnectionTime,
			},
		},
	))
	dialOpts = append(dialOpts, gogrpc.WithUserAgent(connectOpts.userAgent))

	if connectOpts.wait {
		dialOpts = append(dialOpts, gogrpc.WithBlock())
	}

	conn, err := gogrpc.Dial(
		client.serverAddress,
		dialOpts...,
	)
	if err != nil {
		return errors.Wrap(err, "dial connection")
	}
	client.conn = conn

	return nil
}

// Start - starts authentication client module
func (client *Client) Start(ctx context.Context) {
	client.wg.Add(1)
	go client.checkConnectionState(ctx)
}

// Reconnect - returns channel with reconnection events
func (client *Client) Reconnect() <-chan struct{} {
	return client.reconnect
}

// Close - closes authentication client module
func (client *Client) Close() error {
	client.wg.Wait()

	if err := client.conn.Close(); err != nil {
		return err
	}
	close(client.reconnect)
	return nil
}

func (client *Client) checkConnectionState(ctx context.Context) {
	defer client.wg.Done()

	ticker := time.NewTicker(time.Millisecond * 100)
	defer ticker.Stop()

	status := connectivity.Ready
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			state := client.conn.GetState()
			switch state {
			case connectivity.Connecting:
			case connectivity.Idle:
				client.conn.Connect()
			case connectivity.Ready:
				if status != state {
					client.reconnect <- struct{}{}
				}
			case connectivity.Shutdown:
			case connectivity.TransientFailure:
			}
			status = state
		}
	}
}

// Connection - receives connection entity
func (client *Client) Connection() *grpc.ClientConn {
	return client.conn
}

// Stream -
type Stream[T any] struct {
	stream grpc.ClientStream
	data   chan *T

	wg *sync.WaitGroup
}

// NewStream - creates new stream
func NewStream[T any](stream grpc.ClientStream) *Stream[T] {
	return &Stream[T]{
		stream: stream,
		data:   make(chan *T, 1024),

		wg: new(sync.WaitGroup),
	}
}

// Subscribe - generic function to subscribe on service stream
func (s *Stream[T]) Subscribe(ctx context.Context) (uint64, error) {
	var msg pb.SubscribeResponse
	if err := s.stream.RecvMsg(&msg); err != nil {
		return 0, err
	}

	s.wg.Add(1)
	go s.listen(ctx, msg.Id)

	return msg.Id, nil
}

// Listen - channel with received messages
func (s *Stream[T]) Listen() <-chan *T {
	return s.data
}

func (s *Stream[T]) listen(ctx context.Context, id uint64) {
	defer s.wg.Done()

	for {
		select {
		case <-ctx.Done():
			log.Info().Msg("stop listening")
			return
		default:
			var msg T
			err := s.stream.RecvMsg(&msg)
			switch {
			case err == io.EOF:
				log.Info().Msg("connection to gRPC was closed")
				return
			case err != nil:
				log.Err(err).Msg("receiving subscription error")
				return
			default:
				s.data <- &msg
			}
		}
	}
}

// Close - closes stream
func (s *Stream[T]) Close() error {
	s.wg.Wait()

	close(s.data)
	return nil
}

// Unsubscribe -
func (s *Stream[T]) Unsubscribe(ctx context.Context, id uint64) error {
	return s.stream.SendMsg(&pb.UnsubscribeRequest{Id: id})
}

// Context -
func (s *Stream[T]) Context() context.Context {
	return s.stream.Context()
}
