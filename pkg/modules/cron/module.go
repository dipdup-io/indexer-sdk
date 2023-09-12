package cron

import (
	"context"

	"github.com/dipdup-net/indexer-sdk/pkg/modules"
	"github.com/robfig/cron/v3"
)

// Module - cron module
type Module struct {
	cron *cron.Cron
	modules.BaseModule
}

var _ modules.Module = (*Module)(nil)

// NewModule - creates cron module
func NewModule(cfg *Config) (*Module, error) {
	module := &Module{
		BaseModule: modules.New("cron"),
		cron: cron.New(
			cron.WithParser(cron.NewParser(
				cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor,
			)),
			// cron.WithLogger(cron.VerbosePrintfLogger(&log.Logger)),
		),
	}

	for job, pattern := range cfg.Jobs {
		module.CreateOutput(job)

		if _, err := module.cron.AddFunc(
			pattern,
			newHandler(module, job, pattern),
		); err != nil {
			return nil, err
		}
	}

	return module, nil
}

// Start - starts cron scheduler
func (module *Module) Start(ctx context.Context) {
	module.cron.Start()
}

// Close - closes cron scheduler
func (module *Module) Close() error {
	module.cron.Stop()
	return nil
}

func (module *Module) notify(job, pattern string) {
	output, err := module.Output(job)
	if err != nil {
		module.Log.Panic().Msg("while getting output for notification")
	}
	output.Push(struct{}{})
}

func newHandler(module *Module, job, pattern string) func() {
	return func() {
		module.notify(job, pattern)
	}
}
