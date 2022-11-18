package cron

import (
	"context"

	"github.com/dipdup-net/indexer-sdk/pkg/messages"
	"github.com/robfig/cron/v3"
)

// Module - cron module
type Module struct {
	*messages.Publisher

	cron *cron.Cron
}

// NewModule - creates cron module
func NewModule(cfg *Config) (*Module, error) {
	module := &Module{
		Publisher: messages.NewPublisher(),
		cron: cron.New(
			cron.WithParser(cron.NewParser(
				cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor,
			)),
			// cron.WithLogger(cron.VerbosePrintfLogger(&log.Logger)),
		),
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
	module.Publisher.Notify(messages.NewMessage(job, struct{}{}))
}

type handler func()

func newHandler(module *Module, job, pattern string) handler {
	return func() {
		module.notify(job, pattern)
	}
}
