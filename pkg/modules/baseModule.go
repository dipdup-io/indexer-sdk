package modules

import (
	"context"
	"github.com/dipdup-net/indexer-sdk/pkg/sync"
	"github.com/pkg/errors"
)

var _ Module = &BaseModule{}

type BaseModule struct {
	Inputs  *sync.Map[string, *Input]
	Outputs *sync.Map[string, *Output]
}

func (*BaseModule) Name() string {
	return "base_module"
}

func (*BaseModule) Start(ctx context.Context) {}

func (m *BaseModule) Close() error { // TODO-DISCUSS
	if m.Inputs == nil {
		return nil
	}

	return m.Inputs.Range(func(name string, input *Input) (error, bool) { // TODO-DISCUSS
		return input.Close(), false
	})
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

//func (m *BaseModule) AttachTo(output Module, channelName string) error {
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
//}

//func (m *BaseModule) Connect(input Module, channelName string) error {
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
//}
