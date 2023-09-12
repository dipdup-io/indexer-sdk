package modules

import (
	"context"
	"github.com/dipdup-io/workerpool"
	"github.com/dipdup-net/indexer-sdk/pkg/sync"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var _ Module = (*BaseModule)(nil)

type BaseModule struct {
	name    string
	inputs  *sync.Map[string, *Input]
	outputs *sync.Map[string, *Output]
	Log     zerolog.Logger
	G       workerpool.Group
}

func New(name string) BaseModule {
	m := BaseModule{
		name:    name,
		inputs:  sync.NewMap[string, *Input](),
		outputs: sync.NewMap[string, *Output](),
		Log:     log.With().Str("module", name).Logger(),
		G:       workerpool.NewGroup(),
	}

	return m
}

func (m *BaseModule) Name() string {
	return m.name
}

func (*BaseModule) Start(_ context.Context) {}

func (*BaseModule) Close() error {
	return nil
}

func (m *BaseModule) Input(name string) (*Input, error) {
	input, ok := m.inputs.Get(name)
	if !ok {
		return nil, errors.Wrap(ErrUnknownInput, name)
	}
	return input, nil
}

func (m *BaseModule) CreateInput(name string) {
	m.inputs.Set(name, NewInput(name))
}

func (m *BaseModule) Output(name string) (*Output, error) {
	output, ok := m.outputs.Get(name)
	if !ok {
		return nil, errors.Wrap(ErrUnknownOutput, name)
	}
	return output, nil
}

func (m *BaseModule) CreateOutput(name string) {
	m.outputs.Set(name, NewOutput(name))
}

func (m *BaseModule) AttachTo(outputModule Module, outputName, inputName string) error {
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
