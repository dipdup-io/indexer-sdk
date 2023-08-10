package duplicator

import (
	"context"
	"reflect"
	"sync"

	"github.com/dipdup-net/indexer-sdk/pkg/modules"
	"github.com/pkg/errors"
)

// names
const (
	InputName  = "input"
	OutputName = "output"
)

// Duplicator - the structure which is responsible for duplicate signal from one of inputs to all outputs
type Duplicator struct {
	inputs  map[string]*modules.Input
	outputs map[string]*modules.Output

	wg *sync.WaitGroup
}

// NewDuplicator - constructor of Duplicator structure
func NewDuplicator(inputsCount, outputsCount int) *Duplicator {
	d := &Duplicator{
		inputs:  make(map[string]*modules.Input),
		outputs: make(map[string]*modules.Output),

		wg: new(sync.WaitGroup),
	}

	for i := 0; i < inputsCount; i++ {
		name := GetInputName(i)
		d.inputs[name] = modules.NewInput(name)
	}

	for i := 0; i < outputsCount; i++ {
		name := GetOutputName(i)
		d.outputs[name] = modules.NewOutput(name)
	}

	return d
}

// Name -
func (duplicator *Duplicator) Name() string {
	return "duplicator"
}

// Input - returns input by name
func (duplicator *Duplicator) Input(name string) (*modules.Input, error) {
	input, ok := duplicator.inputs[name]
	if !ok {
		return nil, errors.Wrap(modules.ErrUnknownInput, name)
	}
	return input, nil
}

// Output - returns output by name
func (duplicator *Duplicator) Output(name string) (*modules.Output, error) {
	output, ok := duplicator.outputs[name]
	if !ok {
		return nil, errors.Wrap(modules.ErrUnknownOutput, name)
	}
	return output, nil
}

// AttachTo - attach input to output with name
func (duplicator *Duplicator) AttachTo(name string, input *modules.Input) error {
	output, err := duplicator.Output(name)
	if err != nil {
		return err
	}
	output.Attach(input)
	return nil
}

// Close - gracefully stops module
func (duplicator *Duplicator) Close() error {
	duplicator.wg.Wait()

	for _, i := range duplicator.inputs {
		if err := i.Close(); err != nil {
			return err
		}
	}

	return nil
}

// Start - starts module
func (duplicator *Duplicator) Start(ctx context.Context) {
	duplicator.wg.Add(1)
	go duplicator.listen(ctx)
}

func (duplicator *Duplicator) listen(ctx context.Context) {
	defer duplicator.wg.Done()

	cases := []reflect.SelectCase{
		{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(ctx.Done()),
		},
	}
	for _, input := range duplicator.inputs {
		cases = append(cases,
			reflect.SelectCase{
				Dir:  reflect.SelectRecv,
				Chan: reflect.ValueOf(input.Listen()),
			},
		)
	}

	for {
		_, value, ok := reflect.Select(cases)
		if !ok {
			return
		}
		for _, output := range duplicator.outputs {
			output.Push(value.Interface())
		}
	}
}
