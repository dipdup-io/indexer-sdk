package main

import (
	"context"

	"github.com/dipdup-io/workerpool"
	"github.com/dipdup-net/indexer-sdk/examples/grpc/pb"
	"github.com/dipdup-net/indexer-sdk/pkg/modules"
	"github.com/dipdup-net/indexer-sdk/pkg/modules/grpc"
	generalPB "github.com/dipdup-net/indexer-sdk/pkg/modules/grpc/pb"
	"github.com/rs/zerolog/log"
)

// Client -
type Client struct {
	*grpc.Client

	Output *modules.Output[string]

	stream *grpc.Stream[pb.Response]
	client pb.TimeServiceClient
	g      workerpool.Group
}

// NewClient -
func NewClient(server string) *Client {
	return &Client{
		Client: grpc.NewClient(server),
		Output: modules.NewOutput[string](),
		g:      workerpool.NewGroup(),
	}
}

// Start -
func (client *Client) Start(ctx context.Context) {
	client.client = pb.NewTimeServiceClient(client.Connection())

	client.g.GoCtx(ctx, client.reconnect)
}

// SubscribeOnTime -
func (client *Client) SubscribeOnTime(ctx context.Context) (uint64, error) {
	stream, err := client.client.SubscribeOnTime(ctx, new(pb.Request))
	if err != nil {
		return 0, err
	}
	client.stream = grpc.NewStream[pb.Response](stream)

	client.g.GoCtx(ctx, client.handleTime)
	return client.stream.Subscribe(ctx)
}

func (client *Client) reconnect(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-client.Reconnect():
			if client.stream != nil {
				if err := client.stream.Close(); err != nil {
					log.Err(err).Msg("closing stream")
					continue
				}
			}

			if _, err := client.SubscribeOnTime(ctx); err != nil {
				log.Err(err).Msg("subscription error")
			}
		}
	}
}

func (client *Client) handleTime(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return

		case msg := <-client.stream.Listen():
			client.Output.Push(msg.Time)
		}
	}
}

// UnsubscribeFromTime -
func (client *Client) UnsubscribeFromTime(ctx context.Context, id uint64) error {
	if _, err := client.client.UnsubscribeFromTime(ctx, &generalPB.UnsubscribeRequest{
		Id: id,
	}); err != nil {
		return err
	}
	return nil
}

// Close -
func (client *Client) Close() error {
	client.g.Wait()

	if client.stream != nil {
		if err := client.stream.Close(); err != nil {
			return err
		}
	}
	if err := client.Client.Close(); err != nil {
		return err
	}
	return nil
}
