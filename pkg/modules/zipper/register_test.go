package zipper

import (
	"testing"

	"github.com/dipdup-net/indexer-sdk/pkg/modules"
	"github.com/stretchr/testify/require"
)

func TestRegisterModule(t *testing.T) {
	err := modules.Register(NewModule[int]())
	require.NoError(t, err)
}
