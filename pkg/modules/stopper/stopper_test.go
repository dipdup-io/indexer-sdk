package stopper

import (
	"context"
	"github.com/dipdup-net/indexer-sdk/pkg/modules"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStopper_CallsStop(t *testing.T) {
	stopWasCalled := false

	var stopFunc context.CancelFunc = func() {
		stopWasCalled = true
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stopperModule := NewModule(stopFunc)
	stopperModule.Start(ctx)

	stopperInput, err := stopperModule.Input(InputName)
	assert.NoError(t, err)

	// Act: send stop signal to stopper
	stopperInput.Push(struct{}{})

	err = stopperModule.Close()
	assert.NoError(t, err)

	assert.True(t, stopWasCalled, "stop was never called")
}

func TestStopper_CallsStopFromAnotherModule(t *testing.T) {
	stopWasCalled := false

	var stopFunc context.CancelFunc = func() {
		stopWasCalled = true
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stopperModule := NewModule(stopFunc)
	workerModule := &modules.BaseModule{}
	workerModule.Init("worker")
	workerModule.CreateOutput("stop")

	err := stopperModule.AttachTo(workerModule, "stop", InputName)
	assert.NoError(t, err)

	stopperModule.Start(ctx)

	workerOutput, err := workerModule.Output("stop")
	assert.NoError(t, err)

	// Act: send stop signal to stopper
	workerOutput.Push(struct{}{})

	err = stopperModule.Close()
	assert.NoError(t, err)

	assert.True(t, stopWasCalled, "stop was never called")
}
