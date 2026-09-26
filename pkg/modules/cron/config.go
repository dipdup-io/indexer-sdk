package cron

// Config -
type Config struct {
	Jobs map[string]string `validate:"required" yaml:"jobs"`
}
