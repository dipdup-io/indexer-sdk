package modules

import (
	"errors"
)

// errors
var (
	ErrUnknownInput  = errors.New("unknown input")
	ErrUnknownOutput = errors.New("unknown output")
)

// Input -
type Input struct {
	data chan any
	name string
}

// NewInput -
func NewInput(name string) *Input {
	return &Input{
		data: make(chan any, 1024),
		name: name,
	}
}

// Close -
func (input *Input) Close() error {
	close(input.data)
	return nil
}

// Push -
func (input *Input) Push(msg any) {
	input.data <- msg
}

// Listen -
func (input *Input) Listen() <-chan any {
	return input.data
}

// Name -
func (input *Input) Name() string {
	return input.name
}
