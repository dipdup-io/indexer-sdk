package printer

import (
	"context"
	"fmt"
	"github.com/dipdup-net/indexer-sdk/pkg/modules"
)

// Printer - the structure which is responsible for print received messages
type Printer struct {
	modules.BaseModule
	cancelListener context.CancelFunc
}

const InputName = "input"

var _ modules.Module = (*Printer)(nil)

// NewPrinter - constructor of printer structure
func NewPrinter() Printer {
	p := Printer{BaseModule: modules.New("printer")}
	p.CreateInput(InputName)
	return p
}

// Close - gracefully stops module
func (printer *Printer) Close() error {
	if printer.cancelListener != nil {
		printer.cancelListener()
	}
	printer.G.Wait()
	return nil
}

// Start - starts module
func (printer *Printer) Start(ctx context.Context) {
	listenerCtx, c := context.WithCancel(ctx)
	printer.cancelListener = c
	printer.G.GoCtx(listenerCtx, printer.listen)
}

func (printer *Printer) listen(ctx context.Context) {
	input, err := printer.Input(InputName)
	if err != nil {
		printer.Log.Panic().Msg("while getting default input channel")
	}

	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-input.Listen():
			if !ok {
				return
			}
			printer.Log.Info().Str("obj_type", fmt.Sprintf("%T", msg)).Msgf("%##v", msg)
		}
	}
}
