package main

import (
	"context"
	"sync"

	"github.com/dipdup-net/indexer-sdk/examples/grpc/pb"
	"github.com/dipdup-net/indexer-sdk/pkg/modules"
	"github.com/dipdup-net/indexer-sdk/pkg/modules/grpc"
	generalPB "github.com/dipdup-net/indexer-sdk/pkg/modules/grpc/pb"
	"github.com/pkg/errors"
)

// Client -
type Client struct {
	*grpc.Client

	output *modules.Output

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
}

// SubscribeOnMetadata -
func (client *Client) SubscribeOnTime(ctx context.Context) (uint64, error) {
	stream, err := client.client.SubscribeOnTime(ctx, new(pb.Request))
	if err != nil {
		return 0, err
	}

	return grpc.Subscribe[*pb.Response](
		stream,
		client.handleTime,
		client.wg,
	)
}

func (client *Client) handleTime(ctx context.Context, data *pb.Response, id uint64) error {
	client.output.Push(data)
	return nil
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
