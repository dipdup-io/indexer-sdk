package grpc

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/time/rate"
)

type simpleLimiter struct {
	limiter *rate.Limiter
}

func newSimpleLimiter(rps int) *simpleLimiter {
	return &simpleLimiter{limiter: rate.NewLimiter(rate.Every(time.Second/time.Duration(rps)), rps)}
}

// Limit tells wether or not to allow a certain request
// depending on the underlying limiter decision.
func (l *simpleLimiter) Limit(_ context.Context) error {
	if !l.limiter.Allow() {
		return errors.New("rate limit exceeded")
	}

	return nil
}
