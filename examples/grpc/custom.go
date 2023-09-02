package main

import (
	"context"
	"encoding/json"

	"github.com/dipdup-io/workerpool"
	"github.com/rs/zerolog/log"
)

// CustomModule -
type CustomModule struct {
	Input chan string

	g workerpool.Group
}

// NewCustomModule -
func NewCustomModule() *CustomModule {
	return &CustomModule{
		Input: make(chan string),
		g:     workerpool.NewGroup(),
	}
}

// Start -
func (m *CustomModule) Start(ctx context.Context) {
	m.g.GoCtx(ctx, m.listen)
}

func (m *CustomModule) listen(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-m.Input:
			b, _ := json.Marshal(msg)
			log.Info().Str("msg", string(b)).Msg("arrived from grpc module")
		}
	}
}

// Close -
func (m *CustomModule) Close() error {
	m.g.Wait()
	close(m.Input)
	return nil
}

// Name -
func (*CustomModule) Name() string {
	return "custom"
}
