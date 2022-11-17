package main

import (
	"time"

	"github.com/dipdup-net/indexer-sdk/examples/grpc/pb"
	generalPB "github.com/dipdup-net/indexer-sdk/pkg/modules/grpc/pb"
)

// TimeSubscription -
type TimeSubscription struct {
	data chan *pb.Response
}

// NewTimeSubscription -
func NewTimeSubscription() *TimeSubscription {
	return &TimeSubscription{
		data: make(chan *pb.Response, 1024),
	}
}

// Filter -
func (ts *TimeSubscription) Filter(data time.Time) bool {
	// here you can realize filters
	// data is a structure which can be used for filtering rules
	// for example, it can be a model of another module on which server subscirbed
	// it returns true if message has to send.
	return true
}

// Send -
func (ts *TimeSubscription) Send(data *pb.Response) {
	ts.data <- data
}

// Close -
func (ts *TimeSubscription) Close() error {
	close(ts.data)
	return nil
}

// Listen -
func (ts *TimeSubscription) Listen() <-chan *pb.Response {
	return ts.data
}

// Response -
func Response(id uint64, t time.Time) *pb.Response {
	return &pb.Response{
		Response: &generalPB.SubscribeResponse{
			Id: id,
		},
		Time: t.String(),
	}
}
