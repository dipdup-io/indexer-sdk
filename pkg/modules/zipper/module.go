package zipper

import (
	"context"
	"sync"

	"github.com/dipdup-net/indexer-sdk/pkg/modules"
	"github.com/pkg/errors"
)

// Module - zip module
type Module[Key comparable] struct {
	firstInput  *modules.Input
	secondInput *modules.Input

	output *modules.Output

	firstStream  map[Key]Zippable[Key]
	secondStream map[Key]Zippable[Key]

	zipFunc ZipFunction[Key]

	mx *sync.RWMutex
	wg *sync.WaitGroup
}

// NewModule - creates zip module
func NewModule[Key comparable]() *Module[Key] {
	return &Module[Key]{
		firstInput:   modules.NewInput(FirstInputName),
		secondInput:  modules.NewInput(SecondInputName),
		output:       modules.NewOutput(OutputName),
		firstStream:  make(map[Key]Zippable[Key]),
		secondStream: make(map[Key]Zippable[Key]),
		zipFunc:      defaultZip[Key],
		mx:           new(sync.RWMutex),
		wg:           new(sync.WaitGroup),
	}
}

// NewModuleWithFunc - creates zip module with custom zip function
func NewModuleWithFunc[Key comparable](f ZipFunction[Key]) (*Module[Key], error) {
	if f == nil {
		return nil, ErrNilZipFunc
	}
	return &Module[Key]{
		firstInput:  modules.NewInput(FirstInputName),
		secondInput: modules.NewInput(SecondInputName),
		output:      modules.NewOutput(OutputName),
		zipFunc:     f,
		mx:          new(sync.RWMutex),
		wg:          new(sync.WaitGroup),
	}, nil
}

// Name - returns module name
func (*Module[Key]) Name() string {
	return ModuleName
}

// Input - returns input by name
func (m *Module[Key]) Input(name string) (*modules.Input, error) {
	switch name {
	case FirstInputName:
		return m.firstInput, nil
	case SecondInputName:
		return m.secondInput, nil
	default:
		return nil, errors.Wrap(modules.ErrUnknownInput, name)
	}
}

// Output - returns output by name
func (m *Module[Key]) Output(name string) (*modules.Output, error) {
	if name != OutputName {
		return nil, errors.Wrap(modules.ErrUnknownOutput, name)
	}
	return m.output, nil
}

// AttachTo - attach input to output with name
func (m *Module[Key]) AttachTo(name string, input *modules.Input) error {
	output, err := m.Output(name)
	if err != nil {
		return err
	}
	output.Attach(input)
	return nil
}

// Close - gracefully stops module
func (m *Module[Key]) Close() error {
	m.wg.Wait()

	if err := m.firstInput.Close(); err != nil {
		return err
	}
	if err := m.secondInput.Close(); err != nil {
		return err
	}

	return nil
}

// Start - starts module
func (m *Module[Key]) Start(ctx context.Context) {
	m.wg.Add(1)
	go m.listen(ctx)
}

func (m *Module[Key]) listen(ctx context.Context) {
	defer m.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-m.firstInput.Listen():
			if !ok {
				return
			}
			value, ok := msg.(Zippable[Key])
			if !ok {
				continue
			}
			m.zip(value, m.firstStream, m.secondStream)
		case msg, ok := <-m.secondInput.Listen():
			if !ok {
				return
			}
			value, ok := msg.(Zippable[Key])
			if !ok {
				continue
			}
			m.zip(value, m.secondStream, m.firstStream)
		}
	}
}

func (m *Module[Key]) zip(value Zippable[Key], first, second map[Key]Zippable[Key]) {
	if data, ok := second[value.Key()]; !ok {
		first[value.Key()] = value
	} else {
		if result := m.zipFunc(value, data); result != nil {
			m.output.Push(result)
			delete(second, value.Key())
		}
	}
}
