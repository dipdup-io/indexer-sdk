package stopper

import (
	"context"
	"github.com/dipdup-net/indexer-sdk/pkg/modules"
	"github.com/rs/zerolog/log"
)

// Module - cancels context of all application if get signal.
//
//	                |----------------|
//	                |                |
//	-- struct{} ->  |     MODULE     |
//	                |                |
//	                |----------------|
type Module struct {
	modules.BaseModule
	stop context.CancelFunc
}

var _ modules.Module = &Module{}

const (
	InputName = "signal"
)

func NewModule(cancelFunc context.CancelFunc) Module {
	m := Module{
		BaseModule: modules.New("stopper"),
		stop:       cancelFunc,
	}
	m.CreateInput(InputName)

	return m
}

// Start -
func (s *Module) Start(ctx context.Context) {
	s.G.GoCtx(ctx, s.listen)
}

func (s *Module) listen(ctx context.Context) {
	input, err := s.Input(InputName)
	if err != nil {
		s.Log.Panic().Msg("while getting default input channel in listen")
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-input.Listen():
			log.Info().Msg("stop signal received")
			if s.stop != nil {
				log.Info().Msg("cancelling context...")
				s.stop()
				return
			}
		}
	}
}

// Close -
func (s *Module) Close() error {
	s.G.Wait()
	return nil
}
