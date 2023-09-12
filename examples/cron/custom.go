package main

import (
	"context"
	"sync"

	"github.com/dipdup-net/indexer-sdk/pkg/modules"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

var _ modules.Module = &CustomModule{}

// CustomModule -
type CustomModule struct {
	everySecond     *modules.Input
	everyFiveSecond *modules.Input

	wg *sync.WaitGroup
}

// NewCustomModule -
func NewCustomModule() *CustomModule {
	return &CustomModule{
		everySecond:     modules.NewInput("every_second"),
		everyFiveSecond: modules.NewInput("every_five_second"),
		wg:              new(sync.WaitGroup),
	}
}

// Start -
func (m *CustomModule) Start(ctx context.Context) {
	m.wg.Add(1)
	go m.listen(ctx)
}

// Input -
func (m *CustomModule) Input(name string) (*modules.Input, error) {
	switch name {
	case "every_second":
		return m.everySecond, nil
	case "every_five_second":
		return m.everyFiveSecond, nil
	default:
		return nil, errors.Wrap(modules.ErrUnknownInput, name)
	}
}

// Output -
func (m *CustomModule) Output(name string) (*modules.Output, error) {
	return nil, errors.Wrap(modules.ErrUnknownOutput, name)
}

// AttachTo -
func (m *CustomModule) AttachTo(outputM modules.Module, outputName, inputName string) error {
	output, err := outputM.Output(outputName)
	if err != nil {
		return nil
	}

	switch inputName {
	case "every_second":
		output.Attach(m.everySecond)
		return nil
	case "every_five_second":
		output.Attach(m.everyFiveSecond)
		return nil
	default:
		return errors.Wrap(modules.ErrUnknownInput, inputName)
	}
}

func (m *CustomModule) listen(ctx context.Context) {
	defer m.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case <-m.everySecond.Listen():
			log.Info().Msg("arrived from cron module")
		case <-m.everyFiveSecond.Listen():
			log.Info().Msg("arrived from cron module")
		}
	}
}

// Close -
func (m *CustomModule) Close() error {
	m.wg.Wait()

	if err := m.everyFiveSecond.Close(); err != nil {
		return err
	}
	if err := m.everySecond.Close(); err != nil {
		return err
	}

	return nil
}

// Name -
func (*CustomModule) Name() string {
	return "custom"
}
