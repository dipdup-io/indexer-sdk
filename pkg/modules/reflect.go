package modules

import (
	"reflect"

	"github.com/pkg/errors"
)

type ports map[string]reflect.Value

type modulePorts struct {
	inputs  ports
	outputs ports
}

func newModulePorts() modulePorts {
	return modulePorts{
		inputs:  make(ports),
		outputs: make(ports),
	}
}

func getModulePorts(module any) (modulePorts, error) {
	result := newModulePorts()

	value := reflect.ValueOf(module)
	typ := reflect.TypeOf(module)

	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		value = value.Elem()
	}

	if typ.Kind() != reflect.Struct {
		return result, errors.New("module should be struct")
	}

	for i := 0; i < typ.NumField(); i++ {
		field := value.Field(i)
		fieldType := typ.Field(i)

		switch field.Kind() {
		case reflect.Chan:
			switch fieldType.Type.ChanDir() {
			case reflect.RecvDir, reflect.BothDir:
				result.inputs[fieldType.Name] = field
			}
		case reflect.Pointer:
			elemType := fieldType.Type.Elem()
			if elemType.Kind() != reflect.Struct {
				continue
			}
			if fieldType.Name != "Output" {
				continue
			}
			result.outputs[fieldType.Name] = field
		default:
			continue
		}

	}

	return result, nil
}
