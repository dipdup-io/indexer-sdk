package grpc

import (
	"context"
	"io"
	"sync"
	"time"

	"github.com/dipdup-net/indexer-sdk/pkg/messages"
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
	publisher *messages.Publisher
	conn      *gogrpc.ClientConn

	serverAddress string
}

// NewClient - constructor of client structure
func NewClient(server string) *Client {
	return &Client{
		publisher:     messages.NewPublisher(),
		serverAddress: server,
	}
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

// Subscribe - subscribes on events of client module
func (client *Client) Subscribe(s *messages.Subscriber, id messages.SubscriptionID) {
	client.publisher.Subscribe(s, id)
}

// Unsubscribe - unsubscribes from events of client module
func (client *Client) Unsubscribe(s *messages.Subscriber, id messages.SubscriptionID) {
	client.publisher.Unsubscribe(s, id)
}

// Connection - receives connection entity
func (client *Client) Connection() *grpc.ClientConn {
	return client.conn
}

// Publisher - receives publisher entity
func (client *Client) Publisher() *messages.Publisher {
	return client.publisher
}

// ClientStream -
type ClientStream[T any] interface {
	Recv() (T, error)
	grpc.ClientStream
}

// Subscribe - generic function to subscribe on events from server
func Subscribe[T any](
	p *messages.Publisher,
	s *messages.Subscriber,
	stream ClientStream[T],
	handler SubscriptionHandler[T],
	wg *sync.WaitGroup,
) (messages.SubscriptionID, error) {
	var msg pb.SubscribeResponse
	if err := stream.RecvMsg(&msg); err != nil {
		return 0, err
	}

	p.Subscribe(s, msg.Id)

	wg.Add(1)
	go listen(stream, msg.Id, handler, wg)

	return msg.Id, nil
}

// SubscriptionHandler is handled on subscription message
type SubscriptionHandler[T any] func(ctx context.Context, data T, id messages.SubscriptionID) error

func listen[T any](stream ClientStream[T], id messages.SubscriptionID, handler SubscriptionHandler[T], wg *sync.WaitGroup) {
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
				if err := handler(stream.Context(), data, id); err != nil {
					log.Err(err).Msg("subscription handler error")
				}
			}
		}
	}
}
