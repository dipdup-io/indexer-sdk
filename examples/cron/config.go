package main

import "github.com/dipdup-net/indexer-sdk/pkg/modules/cron"

// Config -
type Config struct {
	Cron *cron.Config `validate:"required" yaml:"cron"`
}

// Substitute -
func (c *Config) Substitute() error {
	return nil
}
