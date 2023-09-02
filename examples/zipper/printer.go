package main

import (
	"context"
	"fmt"

	"github.com/dipdup-io/workerpool"
	"github.com/dipdup-net/indexer-sdk/pkg/modules/zipper"
	"github.com/rs/zerolog/log"
)

// predefined constants
const (
	ModuleName = "printer"
	InputName  = "Input"
)

// Printer - the structure which is responsible for print received messages
type Printer struct {
	Input chan *zipper.Result[int]

	g workerpool.Group
}

// NewPrinter - constructor of printer structure
func NewPrinter() *Printer {
	return &Printer{
		Input: make(chan *zipper.Result[int], 16),
		g:     workerpool.NewGroup(),
	}
}

// Name -
func (printer *Printer) Name() string {
	return ModuleName
}

// Close - gracefully stops module
func (printer *Printer) Close() error {
	printer.g.Wait()

	close(printer.Input)
	return nil
}

// Start - starts module
func (printer *Printer) Start(ctx context.Context) {
	printer.g.GoCtx(ctx, printer.listen)
}

func (printer *Printer) listen(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-printer.Input:
			if !ok {
				return
			}
			log.Info().Str("obj_type", fmt.Sprintf("%T", msg)).Msgf("%##v", msg)
		}
	}
}
