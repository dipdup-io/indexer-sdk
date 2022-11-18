package cron

// Config -
type Config struct {
	Jobs map[string]string `yaml:"jobs" validate:"required"`
}
