# Modules

Flow-based programming (FBP) framework. Modules are independent asynchronous units with N inputs and M outputs that can be composed into workflows.

## Module Interface

Every module implements:

```go
type Module interface {
    io.Closer

    Name() string
    Start(ctx context.Context)
    Input(name string) (*Input, error)
    MustInput(name string) *Input
    Output(name string) (*Output, error)
    MustOutput(name string) *Output
    AttachTo(output Module, outputName, inputName string) error
}
```

- `Start` — starts asynchronous processing
- `Close` — gracefully stops the module (from `io.Closer`)
- `Name` — returns the module name used in workflow construction
- `Input` / `MustInput` — returns input by name (or panics)
- `Output` / `MustOutput` — returns output by name (or panics)
- `AttachTo` — connects this module's input to another module's output

## BaseModule

`BaseModule` provides a default implementation to reduce boilerplate:

```go
type BaseModule struct {
    Log zerolog.Logger   // structured logger
    G   workerpool.Group // goroutine group
}
```

Usage:

```go
type MyModule struct {
    modules.BaseModule
}

func NewMyModule() *MyModule {
    m := &MyModule{
        BaseModule: modules.New("my_module"),
    }
    m.CreateInput("input")
    m.CreateOutput("output")
    return m
}

func (m *MyModule) Start(ctx context.Context) {
    m.G.GoCtx(ctx, m.listen)
}

func (m *MyModule) listen(ctx context.Context) {
    input := m.MustInput("input")
    output := m.MustOutput("output")
    for {
        select {
        case <-ctx.Done():
            return
        case msg := <-input.Listen():
            // process msg
            output.Push(result)
        }
    }
}
```

Helper methods:
- `CreateInput(name string)` — creates an input with default capacity (1024)
- `CreateInputWithCapacity(name string, cap int)` — creates an input with custom capacity
- `CreateOutput(name string)` — creates an output

## Inputs and Outputs

### Input

Channel-based message receiver:

```go
input := modules.NewInput("name")           // default capacity 1024
input := modules.NewInputWithCapacity("name", 100)

input.Push(msg)           // send message
<-input.Listen()          // receive message
input.Close()             // close channel
```

### Output

Fan-out broadcaster to connected inputs:

```go
output := modules.NewOutput("name")

output.Attach(input)      // connect an input
output.Push(msg)           // broadcast to all connected inputs
inputs := output.ConnectedInputs()
```

### Connecting Modules

```go
// Helper function
modules.Connect(srcModule, dstModule, "output_name", "input_name")

// Or directly
dstModule.AttachTo(srcModule, "output_name", "input_name")
```

## Workflow

Orchestrates multiple modules:

```go
wf := modules.NewWorkflow(module1, module2, module3)

// Add more modules
wf.Add(module4)
wf.AddWithName(module5, "custom_name")

// Connect modules by name
wf.Connect("cron", "every_second", "processor", "input")

// Start all modules
wf.Start(ctx)

// Retrieve module by name
mod, err := wf.Get("cron")
```

## Built-in Modules

| Module | Description | Docs |
|--------|-------------|------|
| [gRPC](grpc/) | Client/server with subscriptions and metrics | [README](grpc/) |
| [Cron](cron/) | Cron scheduler with named jobs | [README](cron/) |
| [Zipper](zipper/) | Aggregates two streams by key | [README](zipper/) |
| [Stopper](stopper/) | Cancels context on signal | see below |
| [Printer](printer/) | Logs received messages | see below |

### Stopper

Cancels context when a signal is received on its input:

```go
ctx, cancel := context.WithCancel(context.Background())
stop := stopper.NewModule(cancel)
// Input name: stopper.InputName ("signal")
```

### Printer

Debug module that logs all messages received on its input:

```go
p := printer.NewPrinter()
// Input name: printer.InputName ("input")
```

## Example

```go
cronModule, _ := cron.NewModule(cfg.Cron)
customModule := NewCustomModule()

modules.Connect(cronModule, customModule, "every_second", "every_second")
modules.Connect(cronModule, customModule, "every_five_second", "every_five_second")

ctx, cancel := context.WithCancel(context.Background())
customModule.Start(ctx)
cronModule.Start(ctx)

// Wait for shutdown signal...

cancel()
customModule.Close()
cronModule.Close()
```

Full example: [`examples/cron/`](/examples/cron/)
