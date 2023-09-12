package main

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/dipdup-net/indexer-sdk/pkg/modules"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// CustomModule -
type CustomModule struct {
	input *modules.Input

	wg *sync.WaitGroup
}

// NewCustomModule -
func NewCustomModule() *CustomModule {
	return &CustomModule{
		input: modules.NewInput("input"),
		wg:    new(sync.WaitGroup),
	}
}

// Start -
func (m *CustomModule) Start(ctx context.Context) {
	m.wg.Add(1)
	go m.listen(ctx)
}

// Input -
func (m *CustomModule) Input(name string) (*modules.Input, error) {
	if name != "input" {
		return nil, errors.Wrap(modules.ErrUnknownInput, name)
	}
	return m.input, nil
}

// Output -
func (m *CustomModule) Output(name string) (*modules.Output, error) {
	return nil, errors.Wrap(modules.ErrUnknownOutput, name)
}

// AttachTo -
func (m *CustomModule) AttachTo(outputModule modules.Module, outputName, inputName string) error {
	outputChannel, err := outputModule.Output(outputName)
	if err != nil {
		return err
	}

	input, err := m.Input(inputName)
	if err != nil {
		return err
	}

	outputChannel.Attach(input)
	return nil
}

func (m *CustomModule) listen(ctx context.Context) {
	defer m.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-m.input.Listen():
			b, _ := json.Marshal(msg)
			log.Info().Str("msg", string(b)).Msg("arrived from grpc module")
		}
	}
}

// Close -
func (m *CustomModule) Close() error {
	m.wg.Wait()
	return m.input.Close()
}

// Name -
func (*CustomModule) Name() string {
	return "custom"
}
