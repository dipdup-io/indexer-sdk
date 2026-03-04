# Zipper

Module that receives two input streams and aggregates them by key. When matching keys arrive from both streams, the module emits a combined result.

## Usage

```go
import "github.com/dipdup-net/indexer-sdk/pkg/modules/zipper"

zip := zipper.NewModule[int]()
zip.Start(ctx)
defer zip.Close()
```

Zipper is a generic module parameterized by key type (`comparable`).

## Constructors

```go
// Default zip function (exact key match)
func NewModule[Key comparable]() *Module[Key]

// Custom zip function
func NewModuleWithFunc[Key comparable](f ZipFunction[Key]) (*Module[Key], error)
```

## Inputs

Two named inputs, both accept types implementing `Zippable[Key]`:

```go
const (
    FirstInputName  = "first"
    SecondInputName = "second"
)
```

### Zippable Interface

```go
type Zippable[Type comparable] interface {
    Key() Type
}
```

Example:

```go
type OrderEvent struct {
    orderID int
    data    string
}

func (o OrderEvent) Key() int { return o.orderID }
```

## Output

When both streams provide data with matching keys, emits `Result`:

```go
const OutputName = "output"

type Result[Type comparable] struct {
    Key    Type
    First  any
    Second any
}
```

## Custom Zip Function

Override the default key-matching logic:

```go
type ZipFunction[Type comparable] func(x Zippable[Type], y Zippable[Type]) *Result[Type]
```

Default implementation:

```go
func defaultZip[Type comparable](x Zippable[Type], y Zippable[Type]) *Result[Type] {
    if x.Key() != y.Key() {
        return nil
    }
    return &Result[Type]{First: x, Second: y, Key: x.Key()}
}
```

Full example: [`examples/zipper/`](/examples/zipper/)
