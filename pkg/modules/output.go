package modules

import (
	"sync"
)

// Output -
type Output struct {
	connectedInputs []*Input
	name            string

	mx sync.RWMutex
}

// NewOutput -
func NewOutput(name string) *Output {
	return &Output{
		connectedInputs: make([]*Input, 0),
		name:            name,
	}
}

// ConnectedInputs -
func (output *Output) ConnectedInputs() []*Input {
	return output.connectedInputs
}

// Push -
func (output *Output) Push(msg any) {
	output.mx.RLock()
	{
		for i := range output.connectedInputs {
			output.connectedInputs[i].Push(msg)
		}
	}
	output.mx.RUnlock()
}

// Attach -
func (output *Output) Attach(input *Input) {
	output.mx.Lock()
	{
		output.connectedInputs = append(output.connectedInputs, input)
	}
	output.mx.Unlock()
}

// Name -
func (output *Output) Name() string {
	return output.name
}

// Connect -
func Connect(outputModule, inputModule Module, outputName, inputName string) error {
	return inputModule.AttachTo(outputModule, outputName, inputName)
}
