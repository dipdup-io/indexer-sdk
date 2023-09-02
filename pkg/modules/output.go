package modules

import (
	"reflect"
	"sync"

	"github.com/pkg/errors"
)

// Output -
type Output[T any] struct {
	connectedInputs []chan T

	mx sync.RWMutex
}

// NewOutput -
func NewOutput[T any]() *Output[T] {
	return &Output[T]{
		connectedInputs: make([]chan T, 0),
	}
}

// ConnectedInputs -
func (output *Output[T]) ConnectedInputs() []chan T {
	return output.connectedInputs
}

// Push -
func (output *Output[T]) Push(msg T) {
	output.mx.RLock()
	{
		for i := range output.connectedInputs {
			output.connectedInputs[i] <- msg
		}
	}
	output.mx.RUnlock()
}

// Attach -
func (output *Output[T]) Attach(input chan T) {
	output.mx.Lock()
	{
		output.connectedInputs = append(output.connectedInputs, input)
	}
	output.mx.Unlock()
}

// Connect -
func Connect(outputModule, inputModule, outputName, inputName string) error {
	outPorts, ok := globalRegistry.modules[outputModule]
	if !ok {
		return errors.Errorf("unregistered output module: %s", outputModule)
	}
	inPorts, ok := globalRegistry.modules[inputModule]
	if !ok {
		return errors.Errorf("unregistered input module: %s", inputModule)
	}
	out, ok := outPorts.outputs[outputName]
	if !ok {
		return errors.Errorf("unregistered output %s in module %s", outputName, outputModule)
	}
	in, ok := inPorts.inputs[inputName]
	if !ok {
		return errors.Errorf("unregistered input %s in module %s", inputName, inputModule)
	}

	attach := out.MethodByName("Attach")
	attach.Call([]reflect.Value{in})
	return nil
}
