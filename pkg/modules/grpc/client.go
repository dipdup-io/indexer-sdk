package grpc

import (
	"context"
	"time"

	"github.com/dipdup-net/indexer-sdk/pkg/messages"
	"github.com/dipdup-net/indexer-sdk/pkg/modules/grpc/pb"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	gogrpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

// AuthClient - the structure which is responsible for connection to server and handshake
type AuthClient struct {
	publisher *messages.Publisher
	conn      *gogrpc.ClientConn

	authClient pb.HelloServiceClient

	SubscriptionID string
	serverAddress  string
}

// NewAuthClient - constructor of authentication client
func NewAuthClient(server string) *AuthClient {
	return &AuthClient{
		publisher:     messages.NewPublisher(),
		serverAddress: server,
	}
}

// Connect - connects to server
func (client *AuthClient) Connect(ctx context.Context) error {
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
	client.authClient = pb.NewHelloServiceClient(conn)

	hello, err := client.authClient.Hello(ctx, new(pb.HelloRequest))
	if err != nil {
		return errors.Wrap(err, "error after hello request")
	}
	client.SubscriptionID = hello.Id
	return nil
}

// Start - starts authentication client module
func (client *AuthClient) Start(ctx context.Context) {}

// Close - closes authentication client module
func (client *AuthClient) Close() error {
	if err := client.conn.Close(); err != nil {
		return err
	}
	return nil
}

// Subscribe - subscribes on events of authentication client module
func (client *AuthClient) Subscribe(s *messages.Subscriber, topic messages.Topic) {
	client.publisher.Subscribe(s, topic)
}

// Unsubscribe - unsubscribes from events of authentication client module
func (client *AuthClient) Unsubscribe(s *messages.Subscriber, topic messages.Topic) {
	client.publisher.Unsubscribe(s, topic)
}

// Connection - receives connection entity
func (client *AuthClient) Connection() *grpc.ClientConn {
	return client.conn
}

// Publisher - receives publisher entity
func (client *AuthClient) Publisher() *messages.Publisher {
	return client.publisher
}
