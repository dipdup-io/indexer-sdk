package printer

import (
	"context"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestPrinter_ExpectedPrintedValue(t *testing.T) {
	p := NewPrinter()
	logs := &logSink{}
	logger := zerolog.New(logs)
	p.Log = logger

	input, err := p.Input(InputName)
	assert.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	p.Start(ctx)

	input.Push("hello")

	time.Sleep(time.Millisecond * 100) // TODO-UGLY
	err = p.Close()
	assert.NoError(t, err)
	assert.Contains(t, logs.Index(0), "hello")
}

type logSink struct {
	logs []string
}

func (l *logSink) Write(p []byte) (n int, err error) {
	l.logs = append(l.logs, string(p))
	return len(p), nil
}

func (l *logSink) Index(i int) string {
	return l.logs[i]
}
