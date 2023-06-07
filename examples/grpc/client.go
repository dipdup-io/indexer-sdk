package main

import (
	"context"
	"sync"

	"github.com/dipdup-net/indexer-sdk/examples/grpc/pb"
	"github.com/dipdup-net/indexer-sdk/pkg/modules"
	"github.com/dipdup-net/indexer-sdk/pkg/modules/grpc"
	generalPB "github.com/dipdup-net/indexer-sdk/pkg/modules/grpc/pb"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// Client -
type Client struct {
	*grpc.Client

	output *modules.Output
	stream *grpc.Stream[*pb.Response]

	client pb.TimeServiceClient
	wg     *sync.WaitGroup
}

// NewClient -
func NewClient(server string) *Client {
	return &Client{
		Client: grpc.NewClient(server),
		output: modules.NewOutput("time"),
		wg:     new(sync.WaitGroup),
	}
}

// Start -
func (client *Client) Start(ctx context.Context) {
	client.client = pb.NewTimeServiceClient(client.Connection())

	client.wg.Add(1)
	go client.reconnect(ctx)
}

// SubscribeOnTime -
func (client *Client) SubscribeOnTime(ctx context.Context) (uint64, error) {
	stream, err := client.client.SubscribeOnTime(ctx, new(pb.Request))
	if err != nil {
		return 0, err
	}
	client.stream = grpc.NewStream[*pb.Response](stream)

	client.wg.Add(1)
	go client.handleTime(ctx)

	return client.stream.Subscribe(ctx)
}

func (client *Client) reconnect(ctx context.Context) {
	defer client.wg.Done()

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
	defer client.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-client.stream.Listen():
			client.output.Push(msg)
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

// Input -
func (client *Client) Input(name string) (*modules.Input, error) {
	return nil, errors.Wrap(modules.ErrUnknownInput, name)
}

// Output -
func (client *Client) Output(name string) (*modules.Output, error) {
	if name != "time" {
		return nil, errors.Wrap(modules.ErrUnknownOutput, name)
	}
	return client.output, nil
}

// AttachTo -
func (client *Client) AttachTo(name string, input *modules.Input) error {
	output, err := client.Output(name)
	if err != nil {
		return err
	}
	output.Attach(input)
	return nil
}

// Close -
func (client *Client) Close() error {
	client.wg.Wait()

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
