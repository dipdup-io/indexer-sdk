package main

import (
	"context"
	"sync"
	"time"

	"github.com/dipdup-net/indexer-sdk/pkg/modules"
	"github.com/pkg/errors"
)

// CustomModule -
type CustomModule struct {
	additional int
	startValue int
	name       string
	output     *modules.Output
	wg         *sync.WaitGroup
}

// NewCustomModule -
func NewCustomModule(startValue, additional int, name string) *CustomModule {
	return &CustomModule{
		additional: additional,
		startValue: startValue,
		name:       name,
		output:     modules.NewOutput("output"),
		wg:         new(sync.WaitGroup),
	}
}

// Input -
func (m *CustomModule) Input(name string) (*modules.Input, error) {
	return nil, errors.Wrap(modules.ErrUnknownInput, name)
}

// Output -
func (m *CustomModule) Output(name string) (*modules.Output, error) {
	if name != "output" {
		return nil, errors.Wrap(modules.ErrUnknownOutput, name)
	}
	return m.output, nil
}

// AttachTo -
func (m *CustomModule) AttachTo(name string, input *modules.Input) error {
	output, err := m.Output(name)
	if err != nil {
		return err
	}
	output.Attach(input)
	return nil
}

// Close -
func (m *CustomModule) Close() error {
	m.wg.Wait()
	return nil
}

// Name -
func (m *CustomModule) Name() string {
	return m.name
}

// Start -
func (m *CustomModule) Start(ctx context.Context) {
	m.wg.Add(1)
	go m.generate(ctx)
}

func (m *CustomModule) generate(ctx context.Context) {
	defer m.wg.Done()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	data := ZipData{m.startValue, m.name}

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			m.output.Push(data)
			data.key += m.additional
		}
	}
}
