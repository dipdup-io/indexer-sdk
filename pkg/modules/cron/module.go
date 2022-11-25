package cron

import (
	"context"

	"github.com/dipdup-net/indexer-sdk/pkg/modules"
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
)

// Module - cron module
type Module struct {
	cron *cron.Cron

	outputs map[string]*modules.Output
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
		outputs: make(map[string]*modules.Output),
	}
	for job, pattern := range cfg.Jobs {
		module.outputs[job] = modules.NewOutput(job)

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

// Output -
func (module *Module) Output(name string) (*modules.Output, error) {
	output, ok := module.outputs[name]
	if !ok {
		return nil, errors.Wrap(modules.ErrUnknownOutput, name)
	}
	return output, nil
}

// Input -
func (module *Module) Input(name string) (*modules.Input, error) {
	return nil, errors.Wrap(modules.ErrUnknownInput, name)
}

// AttachTo -
func (module *Module) AttachTo(name string, input *modules.Input) error {
	output, err := module.Output(name)
	if err != nil {
		return err
	}

	output.Attach(input)
	return nil
}

func (module *Module) notify(job, pattern string) {
	module.outputs[job].Push(struct{}{})
}

func newHandler(module *Module, job, pattern string) func() {
	return func() {
		module.notify(job, pattern)
	}
}
