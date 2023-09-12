package modules

import (
	"github.com/dipdup-net/indexer-sdk/pkg/sync"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBaseModule_ExistingInput(t *testing.T) {
	bm := New("module")
	existingChannelName := "input-channel"
	bm.CreateInput(existingChannelName)

	// Act
	input, err := bm.Input(existingChannelName)
	assert.NoError(t, err)
	assert.Equal(t, existingChannelName, input.Name())
}

func TestBaseModule_NonExistingInput(t *testing.T) {
	bm := New("module")
	nonExistingChannelName := "non-existing-input-channel"

	// Act
	input, err := bm.Input(nonExistingChannelName)
	assert.ErrorIs(t, err, ErrUnknownInput)
	assert.Errorf(t, err, "%s: %s", ErrUnknownInput.Error(), nonExistingChannelName)
	assert.Nil(t, input)
}

func TestBaseModule_ExistingOutput(t *testing.T) {
	bm := New("module")
	existingChannelName := "output-channel"
	bm.CreateOutput(existingChannelName)

	// Act
	output, err := bm.Output(existingChannelName)
	assert.NoError(t, err)
	assert.Equal(t, existingChannelName, output.Name())
}

func TestBaseModule_NonExistingOutput(t *testing.T) {
	bm := New("module")
	nonExistingChannelName := "non-existing-output-channel"

	// Act
	output, err := bm.Output(nonExistingChannelName)
	assert.ErrorIs(t, err, ErrUnknownOutput)
	assert.Errorf(t, err, "%s: %s", ErrUnknownOutput.Error(), nonExistingChannelName)
	assert.Nil(t, output)
}

func TestBaseModule_AttachToOnExistingChannel(t *testing.T) {
	bmSrc := &BaseModule{outputs: sync.NewMap[string, *Output]()}
	bmDst := &BaseModule{inputs: sync.NewMap[string, *Input]()}
	inputName := "data-in"
	outputName := "data-out"

	bmSrc.CreateOutput(outputName)
	bmDst.CreateInput(inputName)

	input, err := bmDst.Input(inputName)
	assert.NoError(t, err)

	err = bmDst.AttachTo(bmSrc, outputName, inputName)
	assert.NoError(t, err)

	output, err := bmSrc.Output(outputName)
	assert.NoError(t, err)

	output.Push("hello")

	msg := <-input.Listen()
	assert.Equal(t, "hello", msg)

	err = bmSrc.Close()
	assert.NoError(t, err)

	err = bmDst.Close()
	assert.NoError(t, err)
}

func TestBaseModule_ReturnsCorrectName(t *testing.T) {
	bm := New("module")

	name := bm.Name()
	assert.Equal(t, "module", name)
}
