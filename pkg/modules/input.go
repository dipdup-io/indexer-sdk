package modules

import (
	"errors"
	"log"
	"sync"
)

// errors
var (
	ErrUnknownInput  = errors.New("unknown input")
	ErrUnknownOutput = errors.New("unknown output")
)

// Input -
type Input struct {
	data chan any
	name string
}

// NewInput -
func NewInput(name string) *Input {
	return &Input{
		data: make(chan any, 1024),
		name: name,
	}
}

// Close -
func (input *Input) Close() error {
	close(input.data)
	return nil
}

// Push -
func (input *Input) Push(msg any) {
	input.data <- msg
}

// Listen -
func (input *Input) Listen() <-chan any {
	return input.data
}

// Name -
func (input *Input) Name() string {
	return input.name
}

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
	everySecond, err := inputModule.Input(inputName)
	if err != nil {
		log.Panic(err)
	}
	return outputModule.AttachTo(outputName, everySecond)
}
