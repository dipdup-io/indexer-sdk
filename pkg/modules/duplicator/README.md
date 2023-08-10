# Duplicator

Duplicator module repeats signal from one of inputs to all outputs

## Usage

Usage of duplicator module is described by the [example](/examples/duplicator/).

To import module in your code write following line:

```go
import "github.com/dipdup-net/indexer-sdk/pkg/modules/duplicator"
```

Duplicator module implements interface `Module`. So you can use it like any other module. For example:

```go
// create duplicator module with 2 inputs and 2 outputs. Any signal from one of 2 output will be passed to 2 outputs.
duplicatorModule, err := duplicator.NewModule(2, 2)
if err != nil {
    log.Panic(err)
}
// start duplicator module
duplicatorModule.Start(ctx)

// your code is here

// close duplicator module
if err := duplicatorModule.Close(); err != nil {
    log.Panic(err)
}
```

## Input

Module reply to its outputs all signals from inputs. You should connect it with next code:

```go
// with helper function. You should pass a input number to function GetInputName for receiving input name

if err := modules.Connect(customModule, dup, "your_output_of_module", duplicator.GetInputName(0)); err != nil {
    log.Panic(err)
}

```

## Output

You should connect it with next code:

```go
// with helper function. You should pass a output number to function GetOutputName for receiving output name

if err := modules.Connect(dup, customModule, duplicator.GetOutputName(0), "your_output_of_module"); err != nil {
    log.Panic(err)
}

```
