package modules

import (
	"context"

	"github.com/pkg/errors"
)

// errors
var (
	ErrUnknownModule = errors.New("unknown module")
)

// Workflow -
type Workflow struct {
	modules map[string]Module
}

// NewWorkflow -
func NewWorkflow(modules ...Module) *Workflow {
	names := make(map[string]Module)
	for i := range modules {
		names[modules[i].Name()] = modules[i]
	}
	return &Workflow{
		modules: names,
	}
}

// Add - adds module to workflow
func (wf *Workflow) Add(module Module) error {
	return wf.AddWithName(module, module.Name())
}

// AddWithName - adds module to workflow with custom name
func (wf *Workflow) AddWithName(module Module, name string) error {
	if _, ok := wf.modules[name]; ok {
		return errors.Errorf("module with name '%s' is already in the workflow", module.Name())
	}

	wf.modules[name] = module
	return nil
}

// Get - gets module from the workflow by name
func (wf *Workflow) Get(name string) (Module, error) {
	module, ok := wf.modules[name]
	if !ok {
		return nil, errors.Wrap(ErrUnknownModule, name)
	}
	return module, nil
}

// // Connect - connect destination nodule input to source module output
// func (wf *Workflow) Connect(srcModule, srcOutput, destModule, destInput string) error {
// 	src, ok := wf.modules[srcModule]
// 	if !ok {
// 		return errors.Wrap(ErrUnknownModule, srcModule)
// 	}
// 	output, err := src.Output(srcOutput)
// 	if err != nil {
// 		return err
// 	}

// 	dest, ok := wf.modules[destModule]
// 	if !ok {
// 		return errors.Wrap(ErrUnknownModule, destModule)
// 	}
// 	input, err := dest.Input(destInput)
// 	if err != nil {
// 		return err
// 	}

// 	output.Attach(input)
// 	return nil
// }

// Start - starts workflow
func (wf *Workflow) Start(ctx context.Context) {
	// TODO: check on cyclic dependencies.
	// TODO: detect starting order by connections: leafs are first, root is last.
	for _, module := range wf.modules {
		module.Start(ctx)
	}
}
