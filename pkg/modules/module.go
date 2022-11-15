package modules

import (
	"context"
	"io"
)

// Module is the interface which modules has to implement.
type Module interface {
	io.Closer

	Start(ctx context.Context)
}
