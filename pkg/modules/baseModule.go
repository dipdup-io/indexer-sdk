package modules

import (
	"context"
	"github.com/dipdup-net/indexer-sdk/pkg/sync"
	"github.com/pkg/errors"
)

var _ Module = &BaseModule{}

type BaseModule struct {
	inputs  sync.Map[string, *Input]
	outputs sync.Map[string, *Output]
}

func (*BaseModule) Name() string {
	return "base_module"
}

func (*BaseModule) Start(ctx context.Context) {
}

func (*BaseModule) Close() error {
	// TODO: close all inputs
	return nil
}

func (m *BaseModule) Input(name string) (*Input, error) {
	input, ok := m.inputs.Get(name)
	if !ok {
		return nil, errors.Wrap(ErrUnknownInput, name)
	}
	return input, nil
}

func (m *BaseModule) Output(name string) (*Output, error) {
	output, ok := m.outputs.Get(name)
	if !ok {
		return nil, errors.Wrap(ErrUnknownOutput, name)
	}
	return output, nil
}

func (m *BaseModule) AttachTo(outputName string, input *Input) error {
	output, err := m.Output(outputName)
	if err != nil {
		return err
	}

	output.Attach(input)
	return nil
}
