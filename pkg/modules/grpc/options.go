package grpc

import (
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

// ConnectOptions -
type ConnectOptions struct {
	creds credentials.TransportCredentials
}

func newConnectOptions() ConnectOptions {
	return ConnectOptions{
		creds: insecure.NewCredentials(),
	}
}

// ConnectOption -
type ConnectOption func(opts *ConnectOptions)

// WithTlsFromCert -
func WithTlsFromCert(domain string) ConnectOption {
	return func(opts *ConnectOptions) {
		if domain != "" {
			opts.creds = credentials.NewClientTLSFromCert(nil, domain)
		}
	}
}
