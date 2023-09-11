package modules

import (
	"github.com/dipdup-net/indexer-sdk/pkg/sync"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBaseModule_ExistingInput(t *testing.T) {
	bm := &BaseModule{
		Inputs: sync.NewMap[string, *Input](),
	}
	existingChannelName := "input-channel"
	bm.Inputs.Set(existingChannelName, NewInput(existingChannelName))

	// Act
	input, err := bm.Input(existingChannelName)
	assert.NoError(t, err)
	assert.Equal(t, existingChannelName, input.Name())
}

func TestBaseModule_NonExistingInput(t *testing.T) {
	bm := &BaseModule{
		Inputs: sync.NewMap[string, *Input](),
	}
	nonExistingChannelName := "non-existing-input-channel"

	// Act
	input, err := bm.Input(nonExistingChannelName)
	assert.ErrorIs(t, err, ErrUnknownInput)
	assert.Errorf(t, err, "%s: %s", ErrUnknownInput.Error(), nonExistingChannelName)
	assert.Nil(t, input)
}

func TestBaseModule_ExistingOutput(t *testing.T) {
	bm := &BaseModule{
		Outputs: sync.NewMap[string, *Output](),
	}
	existingChannelName := "output-channel"
	bm.Outputs.Set(existingChannelName, NewOutput(existingChannelName))

	// Act
	output, err := bm.Output(existingChannelName)
	assert.NoError(t, err)
	assert.Equal(t, existingChannelName, output.Name())
}

func TestBaseModule_NonExistingOutput(t *testing.T) {
	bm := &BaseModule{
		Outputs: sync.NewMap[string, *Output](),
	}
	nonExistingChannelName := "non-existing-output-channel"

	// Act
	output, err := bm.Output(nonExistingChannelName)
	assert.ErrorIs(t, err, ErrUnknownOutput)
	assert.Errorf(t, err, "%s: %s", ErrUnknownOutput.Error(), nonExistingChannelName)
	assert.Nil(t, output)
}
