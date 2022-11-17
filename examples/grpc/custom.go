package main

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/dipdup-net/indexer-sdk/pkg/messages"
	"github.com/rs/zerolog/log"
)

// CustomModule -
type CustomModule struct {
	*messages.Subscriber

	wg *sync.WaitGroup
}

// NewCustomModule -
func NewCustomModule() *CustomModule {
	return &CustomModule{
		Subscriber: messages.NewSubscriber(),
		wg:         new(sync.WaitGroup),
	}
}

// Start -
func (m *CustomModule) Start(ctx context.Context) {
	m.wg.Add(1)
	go m.listen(ctx)
}

func (m *CustomModule) listen(ctx context.Context) {
	defer m.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-m.Listen():
			b, _ := json.Marshal(msg.Data())
			log.Info().Str("msg", string(b)).Msg("arrived from grpc module")
		}
	}
}

// Close -
func (m *CustomModule) Close() error {
	m.wg.Wait()
	return nil
}
