package grpc

import (
	"time"

	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

// ConnectOptions -
type ConnectOptions struct {
	creds            credentials.TransportCredentials
	reconnectTimeout time.Duration
	reconnectionTime time.Duration
	userAgent        string
	wait             bool
}

func newConnectOptions() ConnectOptions {
	return ConnectOptions{
		creds:            insecure.NewCredentials(),
		reconnectTimeout: time.Second * 5,
		reconnectionTime: time.Minute * 30,
		userAgent:        "grpc-client",
		wait:             false,
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

// WithReconnectTimeout -
func WithReconnectTimeout(duration time.Duration) ConnectOption {
	return func(opts *ConnectOptions) {
		opts.reconnectTimeout = duration
	}
}

// WithReconnectionTime -
func WithReconnectionTime(duration time.Duration) ConnectOption {
	return func(opts *ConnectOptions) {
		opts.reconnectionTime = duration
	}
}

// WaitServer -
func WaitServer() ConnectOption {
	return func(opts *ConnectOptions) {
		opts.wait = true
	}
}

// WithUserAgent -
func WithUserAgent(s string) ConnectOption {
	return func(opts *ConnectOptions) {
		opts.userAgent = s
	}
}
