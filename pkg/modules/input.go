package modules

const (
	defaultCapacity = 1024
)

// Input -
type Input struct {
	data chan any
	name string
}

// NewInput -
func NewInput(name string) *Input {
	return &Input{
		data: make(chan any, defaultCapacity),
		name: name,
	}
}

// NewInputWithCapacity -
func NewInputWithCapacity(name string, cap int) *Input {
	if cap < 0 {
		cap = defaultCapacity
	}
	return &Input{
		data: make(chan any, cap),
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
