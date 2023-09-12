package modules

import (
	"context"
	"github.com/dipdup-io/workerpool"
	"github.com/dipdup-net/indexer-sdk/pkg/sync"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var _ Module = &BaseModule{}

type BaseModule struct {
	name    string
	Inputs  *sync.Map[string, *Input]
	Outputs *sync.Map[string, *Output]
	Log     zerolog.Logger
	G       workerpool.Group
}

func (m *BaseModule) Init(name string) {
	m.name = name
	m.Inputs = sync.NewMap[string, *Input]()
	m.Outputs = sync.NewMap[string, *Output]()
	m.Log = log.With().Str("module", name).Logger()
	m.G = workerpool.NewGroup()
}

func (*BaseModule) Name() string {
	return "base_module"
}

func (*BaseModule) Start(_ context.Context) {}

func (*BaseModule) Close() error {
	return nil
}

func (m *BaseModule) Input(name string) (*Input, error) {
	input, ok := m.Inputs.Get(name)
	if !ok {
		return nil, errors.Wrap(ErrUnknownInput, name)
	}
	return input, nil
}

func (m *BaseModule) Output(name string) (*Output, error) {
	output, ok := m.Outputs.Get(name)
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

// func (m *BaseModule) AttachTo(output Module, channelName string) error {
//	outputChannel, err := output.Output(channelName)
//	if err != nil {
//		return err
//	}
//
//	input, e := m.Input(channelName)
//	if e != nil {
//		return e
//	}
//
//	outputChannel.Attach(input)
//	return nil
// }

// func (m *BaseModule) Connect(input Module, channelName string) error {
//	inputChannel, err := input.Input(channelName)
//	if err != nil {
//		return err
//	}
//
//	output, e := m.Output(channelName)
//	if e != nil {
//		return e
//	}
//
//	output.Attach(inputChannel)
//	return nil
// }
