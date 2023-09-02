package cron

import (
	"testing"

	"github.com/dipdup-net/indexer-sdk/pkg/modules"
	"github.com/stretchr/testify/require"
)

func TestRegisterModule(t *testing.T) {
	err := modules.Register(&Module{})
	require.NoError(t, err)
}
