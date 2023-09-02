package main

import (
	"context"
	"sync"

	"github.com/dipdup-net/indexer-sdk/pkg/modules/cron"
	"github.com/rs/zerolog/log"
)

// CustomModule -
type CustomModule struct {
	Messages chan cron.Message

	wg *sync.WaitGroup
}

// NewCustomModule -
func NewCustomModule() *CustomModule {
	return &CustomModule{
		Messages: make(chan cron.Message, 16),
		wg:       new(sync.WaitGroup),
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
		case msg := <-m.Messages:
			log.Info().Str("job", msg.Job).Msg("arrived from cron module")
		}
	}
}

// Close -
func (m *CustomModule) Close() error {
	m.wg.Wait()

	close(m.Messages)
	return nil
}

// Name -
func (*CustomModule) Name() string {
	return "custom"
}
