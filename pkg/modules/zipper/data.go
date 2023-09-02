package zipper

import "github.com/pkg/errors"

// predefined names
const (
	FirstInputName  = "First"
	SecondInputName = "Second"
	OutputName      = "Output"
	ModuleName      = "zipper"
)

// Zippable - interface of data which can be zipped
type Zippable[Type comparable] interface {
	Key() Type
}

// Result - data structure is result of zip operation. It has two entity
type Result[Type comparable] struct {
	Key    Type
	First  any
	Second any
}

// ZipFunction - function to check can we zip two entity and if we can zip it.
// If function returns nil then it can't zip 2 entities
type ZipFunction[Type comparable] func(x Zippable[Type], y Zippable[Type]) *Result[Type]

// errors
var (
	ErrNilZipFunc = errors.New("nil zip function")
)
