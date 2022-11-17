package main

import (
	"context"
	"errors"
	"sync"

	"github.com/dipdup-net/indexer-sdk/examples/grpc/pb"
	"github.com/dipdup-net/indexer-sdk/pkg/messages"
	"github.com/dipdup-net/indexer-sdk/pkg/modules/grpc"
	generalPB "github.com/dipdup-net/indexer-sdk/pkg/modules/grpc/pb"
	"github.com/rs/zerolog/log"
)

// Client -
type Client struct {
	*grpc.Client
	client pb.TimeServiceClient
	wg     *sync.WaitGroup
}

// NewClient -
func NewClient(server string) *Client {
	return &Client{
		Client: grpc.NewClient(server),
		wg:     new(sync.WaitGroup),
	}
}

// Start -
func (client *Client) Start(ctx context.Context) {
	client.client = pb.NewTimeServiceClient(client.Connection())
}

// SubscribeOnMetadata -
func (client *Client) SubscribeOnTime(ctx context.Context, s *messages.Subscriber) (messages.SubscriptionID, error) {
	stream, err := client.client.SubscribeOnTime(ctx, new(pb.Request))
	if err != nil {
		return 0, err
	}

	return grpc.Subscribe[*pb.Response](
		client.Publisher(),
		s,
		stream,
		client.handleTime,
		client.wg,
	)
}

func (client *Client) handleTime(ctx context.Context, data *pb.Response, id messages.SubscriptionID) error {
	log.Info().Str("time", data.Time).Msg("now")
	client.Publisher().Notify(messages.NewMessage(id, data))
	return nil
}

// UnsubscribeFromTime -
func (client *Client) UnsubscribeFromTime(ctx context.Context, s *messages.Subscriber, id messages.SubscriptionID) error {
	subscriptionID, ok := id.(uint64)
	if !ok {
		return errors.New("invalid subscription id")
	}

	if _, err := client.client.UnsubscribeFromTime(ctx, &generalPB.UnsubscribeRequest{
		Id: subscriptionID,
	}); err != nil {
		return err
	}

	client.Publisher().Unsubscribe(s, id)
	return nil
}
