# Modules

 Flow-based programming is the base of SDK. It is a set of asynchronous modules which can be combinated in workflow. Each module implement a function that has N inputs and M outputs and executes asynchronously (N, M >= 0). In our implementation inputs are untyped. So each module has to check type of received message. 

 ## Interface

 Module has to implement the interface `Module`:

 ```go
type Module interface {
	io.Closer

	Name() string
	Start(ctx context.Context)
	Input(name string) (*Input, error)
	Output(name string) (*Output, error)
	AttachTo(outputName string, input *Input) error
}
 ``` 

Interface contains following methods:

 * `Start` - starts asynchronous waiting of messages in inputs and initialize module state.
 * `Close` - gracefully stops all module activities (inherited from `io.Closer` interface).
 * `Name` - returns name of module which will be used in workflow construction.
 * `Input` - returns input by its name.
 * `Output` - returns output by its name.
 * `AttachTo` - connects output with name to passed input of another module. 

## Inputs and outputs

All communication between modules is implemented via inputs/outputs. Input is the structure contains channel with `any` as data. It also has name which will be its identity in module's scope.

```go
type Input struct {
	data chan any
	name string
}
```
It has following methods:

* `Close` - closes channel of input
* `Push` - sends message to channel
* `Listen` - waits new message
* `Name` - returns input name

Output is the set of inputs which connected to it. When module send message to output it iterates over all connected inputs and pushes message to them. Output also has name which identifies it.

```go
type Output struct {
	connectedInputs []*Input
	name            string

	mx sync.RWMutex
}
```

It has following methods:

* `ConnectedInputs` - returns all connected inputs
* `Push` - pushes message to all connected inputs
* `Attach` - adds input to connected inputs array
* `Name` - returns output name

SDK has helper function `Connect`:

```go
func Connect(outputModule, inputModule Module, outputName, inputName string) error 
```

The function receives `outputModule` and `inputModule`: modules which will be connected. Also it receives input and output names in that modules which will be connected.

## Workflow

Modules can be united in workflow. Workflow is set of module which connected in certain seqeunce. To create `Workflow` you can call function `NewWorkflow`:

```go
func NewWorkflow(modules ...Module) *Workflow
```

`Workflow` has following functions:

* `Add` - adds module with name returning from its `Name` method
* `AddWithName` - adds module with name passed to it
* `Get` - returns module by name which module was created with
* `Connect` - connects modules with certain names by names of its input and output
* `Start` - starts all modules in workflow

## Implemented modules

SDK has some modules which can be used during workflow creation.

### gRPC

gRPC module where realized default client and server. Detailed docs can be found [here](/pkg/modules/grpc/).

### Cron

Cron module implements cron scheduler. Detailed docs can be found [here](/pkg/modules/cron/).