package cron

import (
	"context"

	"github.com/dipdup-net/indexer-sdk/pkg/modules"
	"github.com/robfig/cron/v3"
)

type Message struct {
	Job string
}

// Module - cron module
type Module struct {
	cron *cron.Cron

	Output *modules.Output[Message]
}

// NewModule - creates cron module
func NewModule(cfg *Config) (*Module, error) {
	module := &Module{
		cron: cron.New(
			cron.WithParser(cron.NewParser(
				cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor,
			)),
			// cron.WithLogger(cron.VerbosePrintfLogger(&log.Logger)),
		),
		Output: modules.NewOutput[Message](),
	}
	for job, pattern := range cfg.Jobs {
		if _, err := module.cron.AddFunc(
			pattern,
			newHandler(module, job, pattern),
		); err != nil {
			return nil, err
		}
	}

	return module, nil
}

// Name -
func (*Module) Name() string {
	return "cron"
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
	module.Output.Push(Message{
		Job: job,
	})
}

func newHandler(module *Module, job, pattern string) func() {
	return func() {
		module.notify(job, pattern)
	}
}
