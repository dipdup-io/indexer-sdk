package zipper

import (
	"context"
	"sync"

	"github.com/dipdup-io/workerpool"
	"github.com/dipdup-net/indexer-sdk/pkg/modules"
)

// Module - zip module
type Module[Key comparable] struct {
	First  chan Zippable[Key]
	Second chan Zippable[Key]

	Output *modules.Output[*Result[Key]]

	firstStream  map[Key]Zippable[Key]
	secondStream map[Key]Zippable[Key]

	zipFunc ZipFunction[Key]

	mx *sync.RWMutex
	g  workerpool.Group
}

// NewModule - creates zip module
func NewModule[Key comparable]() *Module[Key] {
	return &Module[Key]{
		First:        make(chan Zippable[Key], 1024),
		Second:       make(chan Zippable[Key], 1024),
		Output:       modules.NewOutput[*Result[Key]](),
		firstStream:  make(map[Key]Zippable[Key]),
		secondStream: make(map[Key]Zippable[Key]),
		zipFunc:      defaultZip[Key],
		mx:           new(sync.RWMutex),
		g:            workerpool.NewGroup(),
	}
}

// NewModuleWithFunc - creates zip module with custom zip function
func NewModuleWithFunc[Key comparable](f ZipFunction[Key]) (*Module[Key], error) {
	if f == nil {
		return nil, ErrNilZipFunc
	}
	return &Module[Key]{
		First:   make(chan Zippable[Key], 1024),
		Second:  make(chan Zippable[Key], 1024),
		Output:  modules.NewOutput[*Result[Key]](),
		zipFunc: f,
		mx:      new(sync.RWMutex),
		g:       workerpool.NewGroup(),
	}, nil
}

// Name - returns module name
func (*Module[Key]) Name() string {
	return ModuleName
}

// Close - gracefully stops module
func (m *Module[Key]) Close() error {
	m.g.Wait()

	close(m.First)
	close(m.Second)

	return nil
}

// Start - starts module
func (m *Module[Key]) Start(ctx context.Context) {
	m.g.GoCtx(ctx, m.listen)
}

func (m *Module[Key]) listen(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-m.First:
			if !ok {
				return
			}
			m.zip(msg, m.firstStream, m.secondStream)
		case msg, ok := <-m.Second:
			if !ok {
				return
			}
			m.zip(msg, m.secondStream, m.firstStream)
		}
	}
}

func (m *Module[Key]) zip(value Zippable[Key], first, second map[Key]Zippable[Key]) {
	if data, ok := second[value.Key()]; !ok {
		first[value.Key()] = value
	} else {
		if result := m.zipFunc(value, data); result != nil {
			m.Output.Push(result)
			delete(second, value.Key())
		}
	}
}
