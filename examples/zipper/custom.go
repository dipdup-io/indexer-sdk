package main

import (
	"context"
	"time"

	"github.com/dipdup-io/workerpool"
	"github.com/dipdup-net/indexer-sdk/pkg/modules"
	"github.com/dipdup-net/indexer-sdk/pkg/modules/zipper"
)

// CustomModule -
type CustomModule struct {
	Output *modules.Output[zipper.Zippable[int]]

	additional int
	startValue int
	name       string
	g          workerpool.Group
}

// NewCustomModule -
func NewCustomModule(startValue, additional int, name string) *CustomModule {
	return &CustomModule{
		Output: modules.NewOutput[zipper.Zippable[int]](),

		additional: additional,
		startValue: startValue,
		name:       name,
		g:          workerpool.NewGroup(),
	}
}

// Close -
func (m *CustomModule) Close() error {
	m.g.Wait()
	return nil
}

// Name -
func (m *CustomModule) Name() string {
	return m.name
}

// Start -
func (m *CustomModule) Start(ctx context.Context) {
	m.g.GoCtx(ctx, m.generate)
}

func (m *CustomModule) generate(ctx context.Context) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	data := ZipData{m.startValue, m.name}

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			m.Output.Push(data)
			data.key += m.additional
		}
	}
}
