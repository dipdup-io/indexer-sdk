package main

import (
	"context"
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
			log.Info().Str("job", msg.SubscriptionID().(string)).Msg("arrived from cron module")
		}
	}
}

// Close -
func (m *CustomModule) Close() error {
	m.wg.Wait()
	return nil
}
