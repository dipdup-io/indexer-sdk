package stopper

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStopperCallsStop(t *testing.T) {
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
