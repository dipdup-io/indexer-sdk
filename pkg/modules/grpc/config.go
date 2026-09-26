package grpc

// ServerConfig - config for server
type ServerConfig struct {
	Bind    string `validate:"required,hostname_port" yaml:"bind"`
	Log     bool   `validate:"omitempty"              yaml:"log"`
	Metrics bool   `validate:"omitempty"              yaml:"metrics"`
	RPS     int    `validate:"omitempty,min=1"        yaml:"rps"`
}
