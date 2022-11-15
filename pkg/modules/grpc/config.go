package grpc

// ServerConfig - config for server
type ServerConfig struct {
	Bind string `yaml:"bind" validate:"required,hostname_port"`
}
