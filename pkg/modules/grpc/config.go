package grpc

// ServerConfig - config for server
type ServerConfig struct {
	Bind    string `yaml:"bind" validate:"required,hostname_port"`
	Log     bool   `yaml:"log" validate:"omitempty"`
	Metrics bool   `yaml:"metrics" validate:"omitempty"`
	RPS     int    `yaml:"rps" validate:"omitempty,min=1"`
}
