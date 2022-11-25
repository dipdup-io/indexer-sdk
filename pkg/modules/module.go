package modules

import (
	"context"
	"io"
)

// Module is the interface which modules has to implement.
type Module interface {
	io.Closer

	Name() string

	Start(ctx context.Context)

	Input(name string) (*Input, error)
	Output(name string) (*Output, error)
	AttachTo(outputName string, input *Input) error
}
