package modules

import (
	"context"
	"github.com/pkg/errors"
	"io"
)

// errors
var (
	ErrUnknownInput  = errors.New("unknown input")
	ErrUnknownOutput = errors.New("unknown output")
)

// Module is the interface which modules have to implement.
type Module interface {
	io.Closer

	Name() string
	Start(ctx context.Context)
}
