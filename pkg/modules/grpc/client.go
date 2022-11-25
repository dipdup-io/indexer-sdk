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
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

// Client - the structure which is responsible for connection to server
type Client struct {
	conn *gogrpc.ClientConn

	serverAddress string
}

// NewClient - constructor of client structure
func NewClient(server string) *Client {
	return &Client{
		serverAddress: server,
	}
}

// Name -
func (client *Client) Name() string {
	return "grpc_client"
}

// Connect - connects to server
func (client *Client) Connect(ctx context.Context) error {
	conn, err := gogrpc.Dial(
		client.serverAddress,
		gogrpc.WithTransportCredentials(insecure.NewCredentials()),
		gogrpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:    20 * time.Second,
			Timeout: 10 * time.Second,
		}),
	)
	if err != nil {
		return errors.Wrap(err, "dial connection")
	}
	client.conn = conn

	return nil
}

// WaitConnect - trying to connect to server
func (client *Client) WaitConnect(ctx context.Context) error {
	conn, err := gogrpc.Dial(
		client.serverAddress,
		gogrpc.WithBlock(),
		gogrpc.WithTransportCredentials(insecure.NewCredentials()),
		gogrpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:    20 * time.Second,
			Timeout: 10 * time.Second,
		}),
	)
	if err != nil {
		return errors.Wrap(err, "dial connection")
	}
	client.conn = conn

	return nil
}

// Start - starts authentication client module
func (client *Client) Start(ctx context.Context) {}

// Close - closes authentication client module
func (client *Client) Close() error {
	if err := client.conn.Close(); err != nil {
		return err
	}
	return nil
}

// Connection - receives connection entity
func (client *Client) Connection() *grpc.ClientConn {
	return client.conn
}

// ClientStream -
type ClientStream[T any] interface {
	Recv() (T, error)
	grpc.ClientStream
}

// Subscribe - generic function to subscribe on events from server
func Subscribe[T any](
	stream ClientStream[T],
	handler SubscriptionHandler[T],
	wg *sync.WaitGroup,
) (uint64, error) {
	var msg pb.SubscribeResponse
	if err := stream.RecvMsg(&msg); err != nil {
		return 0, err
	}

	wg.Add(1)
	go listen(stream, msg.Id, handler, wg)

	return msg.Id, nil
}

// SubscriptionHandler is handled on subscription message
type SubscriptionHandler[T any] func(ctx context.Context, data T, subscriptionID uint64) error

func listen[T any](stream ClientStream[T], subscriptionID uint64, handler SubscriptionHandler[T], wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-stream.Context().Done():
			return
		default:
			data, err := stream.Recv()
			if err == io.EOF {
				continue
			}
			if err != nil {
				log.Err(err).Msg("receiving subscription error")
				continue
			}

			if handler != nil {
				if err := handler(stream.Context(), data, subscriptionID); err != nil {
					log.Err(err).Msg("subscription handler error")
				}
			}
		}
	}
}
