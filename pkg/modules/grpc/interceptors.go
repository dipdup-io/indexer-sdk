package grpc

import (
	"context"
	"errors"
	"strconv"
	"sync/atomic"

	"github.com/dipdup-net/indexer-sdk/pkg/messages"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type subscriptionKey struct{}

// headers
const (
	SubscriptionIdHeader = "subscription_id"
)

// GetSubscriptionID - returns current subscriptions id from context. If subscription id is not found then raise an error.
func GetSubscriptionID(ctx context.Context) (uint64, error) {
	if id := ctx.Value(subscriptionKey{}); id != nil {
		return id.(uint64), nil
	}
	return 0, errors.New("there aren't any subscriptions")
}

// ReceiveSubscriptionID -
func ReceiveSubscriptionID(stream grpc.ClientStream) (messages.SubscriptionID, error) {
	md, err := stream.Header()
	if err != nil {
		return 0, err
	}

	values := md.Get(SubscriptionIdHeader)
	if len(values) == 0 {
		return 0, errors.New("can't get subscription id")
	}

	return strconv.ParseUint(values[0], 10, 64)
}

var subscriptionsCounter = new(atomic.Uint64)

func subscriptionStreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		newSubscriptionID := subscriptionsCounter.Add(1)
		ctx := context.WithValue(stream.Context(), subscriptionKey{}, newSubscriptionID)
		wrapper := newStreamContextWrapper(stream)
		wrapper.SetContext(ctx)

		if err := grpc.SendHeader(ctx, metadata.Pairs(SubscriptionIdHeader, strconv.FormatUint(newSubscriptionID, 10))); err != nil {
			return err
		}

		return handler(srv, wrapper)
	}
}

type streamContextWrapper interface {
	grpc.ServerStream
	SetContext(context.Context)
}

type wrapper struct {
	grpc.ServerStream
	//nolint
	ctx context.Context
}

// Context -
func (w *wrapper) Context() context.Context {
	return w.ctx
}

// SetContext -
func (w *wrapper) SetContext(ctx context.Context) {
	w.ctx = ctx
}

func newStreamContextWrapper(inner grpc.ServerStream) streamContextWrapper {
	ctx := inner.Context()
	return &wrapper{
		inner,
		ctx,
	}
}
