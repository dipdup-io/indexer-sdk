package printer

import (
	"context"
	"fmt"
	"sync"

	"github.com/dipdup-net/indexer-sdk/pkg/modules"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// input name
const (
	InputName = "input"
)

// Printer - the structure which is responsible for print received messages
type Printer struct {
	input *modules.Input

	wg *sync.WaitGroup
}

// NewPrinter - constructor of printer structure
func NewPrinter() *Printer {
	return &Printer{
		input: modules.NewInput(InputName),

		wg: new(sync.WaitGroup),
	}
}

// Name -
func (*Printer) Name() string {
	return "printer"
}

// Input - returns input by name
func (printer *Printer) Input(name string) (*modules.Input, error) {
	switch name {
	case InputName:
		return printer.input, nil
	default:
		return nil, errors.Wrap(modules.ErrUnknownInput, name)
	}
}

// Output - returns output by name
func (printer *Printer) Output(name string) (*modules.Output, error) {
	return nil, errors.Wrap(modules.ErrUnknownOutput, name)
}

// AttachTo - attach input to output with name
func (printer *Printer) AttachTo(name string, input *modules.Input) error {
	output, err := printer.Output(name)
	if err != nil {
		return err
	}
	output.Attach(input)
	return nil
}

// Close - gracefully stops module
func (printer *Printer) Close() error {
	printer.wg.Wait()

	if err := printer.input.Close(); err != nil {
		return err
	}

	return nil
}

// Start - starts module
func (printer *Printer) Start(ctx context.Context) {
	printer.wg.Add(1)
	go printer.listen(ctx)
}

func (printer *Printer) listen(ctx context.Context) {
	defer printer.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-printer.input.Listen():
			if !ok {
				return
			}
			log.Info().Str("obj_type", fmt.Sprintf("%T", msg)).Msgf("%##v", msg)
		}
	}
}
