package main

import "github.com/dipdup-net/indexer-sdk/pkg/modules/cron"

// Config -
type Config struct {
	Cron *cron.Config `yaml:"cron" validate:"required"`
}

// Substitute -
func (c *Config) Substitute() error {
	return nil
}
