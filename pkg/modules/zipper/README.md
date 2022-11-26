# Zipper

It's a module receiving two input streams and aggregate them by key.

## Usage

Usage of module is described by the [example](/examples/zipper/).

To import module in your code write following line:

```go
import "github.com/dipdup-net/indexer-sdk/pkg/modules/zippper"
```

Zipper module implements interface `Module`. So you can use it like any other module. For example:

```go
// create zip module
zip := zipper.NewModule[int]()

// start zip module
zip.Start(ctx)

// your code is here

// close zip module
if err := zip.Close(); err != nil {
    log.Panic(err)
}
```

Zipper is the generic structure. Type parameter `Key` is `comparable` constraint

```go
type Module[Key comparable] struct {
	firstInput  *modules.Input
	secondInput *modules.Input

	output *modules.Output

	firstStream  map[Key]Zippable[Key]
	secondStream map[Key]Zippable[Key]

	zipFunc ZipFunction[Key]

	mx *sync.RWMutex
	wg *sync.WaitGroup
}
```

Module has 2 constructors:

```go
// creates zip module with default ZipFunction
func NewModule[Key comparable]() *Module[Key]

// creates zip module with custom ZipFunction
func NewModuleWithFunc[Key comparable](f ZipFunction[Key]) (*Module[Key], error)
```

`ZipFunction` is the function which can be used to redeclare key comparing rules. The function has to return nil if it can't zip structures. It has following declaration:

```go
type ZipFunction[Type comparable] func(x Zippable[Type], y Zippable[Type]) *Result[Type]
```

Default zip function declares like that:

```go
func defaultZip[Type comparable](x Zippable[Type], y Zippable[Type]) *Result[Type] {
	if x.Key() != y.Key() {
		return nil
	}
	return &Result[Type]{
		First:  x,
		Second: y,
		Key:    x.Key(),
	}
}
```

## Input

Inputs receives types realizing `Zippable` interface.

```go
type Zippable[Type comparable] interface {
	Key() Type
}
```

`Zippable` requires realization of one method `Key` which returns key. Received key will be used for zipping data streams. For example:

```go
zip := zipper.NewModule[int]()

type ZipData struct {
	key   int
	value string
}

func (z ZipData) Key() int {
	return z.key
}
```

In the example zip module created with `Key` type `int`. That's why to use the module you should realize method `Key` on data structure which returns `int`. 

Package contains declared constants of inputs name:

```go
FirstInputName  = "first"
SecondInputName = "second"
```

## Output

When module zipped data stream it sends `Result` structure to output.

```go
type Result[Type comparable] struct {
	Key    Type
	First  any
	Second any
}
```

It contains key which was used to zip. Also it contains data of first and second stream types.

Package contains declared constant of output name:

```go
OutputName = "output"
```