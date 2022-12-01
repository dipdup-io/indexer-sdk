package contract

import (
	"bytes"
	"encoding/hex"
	stdJSON "encoding/json"
	"fmt"
	"strings"

	js "github.com/dipdup-net/indexer-sdk/pkg/jsonschema"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/pkg/errors"
)

// EVM -
type EVM struct{}

// NewEVM -
func NewEVM() *EVM {
	return new(EVM)
}

// JSONSchema -
func (evm *EVM) JSONSchema(data []byte) ([]byte, error) {
	var contractABI abi.ABI
	if err := json.Unmarshal(data, &contractABI); err != nil {
		return nil, err
	}

	schema, err := evm.createEntrypointsSchema(contractABI)
	if err != nil {
		return nil, err
	}

	jsData, err := json.Marshal(schema)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	err = stdJSON.Compact(&buf, jsData)
	return buf.Bytes(), err
}

func (vm *EVM) createEntrypointsSchema(contractABI abi.ABI) (map[string]js.Type, error) {
	result := make(map[string]js.Type)

	for _, event := range contractABI.Events {
		var (
			typ = js.Type{
				Type: "event",
				Inputs: &js.JSONSchema{
					Schema: js.Draft201909,
					Type:   js.ItemTypeObject,
				},
				Signature: strings.TrimPrefix(event.ID.Hex(), "0x"),
			}
		)

		inputsBody, err := getBodyByArgs(event.Inputs)
		if err != nil {
			return nil, err
		}
		if inputsBody != nil {
			typ.Inputs.ObjectItem = inputsBody
		}

		result[event.Name] = typ
	}

	for _, method := range contractABI.Methods {
		var (
			typ = js.Type{
				Type: "method",
				Inputs: &js.JSONSchema{
					Schema: js.Draft201909,
					Type:   js.ItemTypeObject,
				},
				Outputs: &js.JSONSchema{
					Schema: js.Draft201909,
					Type:   js.ItemTypeObject,
				},
				Signature: hex.EncodeToString(method.ID),
			}
		)

		inputsBody, err := getBodyByArgs(method.Inputs)
		if err != nil {
			return nil, err
		}
		if inputsBody != nil {
			typ.Inputs.ObjectItem = inputsBody
		}

		outputsBody, err := getBodyByArgs(method.Outputs)
		if err != nil {
			return nil, err
		}
		if outputsBody != nil {
			typ.Outputs.ObjectItem = outputsBody
		}

		result[method.Name] = typ
	}

	return result, nil
}

func getBodyByArgs(args abi.Arguments) (*js.ObjectItem, error) {
	if len(args) == 0 {
		return nil, nil
	}
	body := &js.ObjectItem{
		Properties: make(map[string]js.JSONSchema),
		Required:   []string{},
	}
	for idx, arg := range args {
		argSchema, err := createSchemaItem(arg.Name, idx, &arg.Type)
		if err != nil {
			return nil, err
		}
		argSchema.Indexed = arg.Indexed
		body.Properties[argSchema.Title] = argSchema
		body.Required = append(body.Required, argSchema.Title)
	}

	return body, nil
}

func createSchemaItem(name string, idx int, typ *abi.Type) (js.JSONSchema, error) {
	if name == "" {
		name = fmt.Sprintf("%s_%02d", typ.String(), idx)
	}
	switch typ.T {
	case abi.AddressTy, abi.StringTy:
		return js.JSONSchema{
			Type:         js.ItemTypeString,
			Title:        name,
			InternalType: typ.String(),
		}, nil

	case abi.ArrayTy, abi.SliceTy:
		schema := js.JSONSchema{
			Type:  js.ItemTypeArray,
			Title: name,
		}

		elemName := fmt.Sprintf("%s_elem", name)
		elem, err := createSchemaItem(elemName, idx, typ.Elem)
		if err != nil {
			return elem, err
		}

		schema.ArrayItem = &js.ArrayItem{
			Items: []js.JSONSchema{
				elem,
			},
		}

		return schema, nil

	case abi.TupleTy:
		schema := js.JSONSchema{
			Type:  js.ItemTypeObject,
			Title: name,
			ObjectItem: &js.ObjectItem{
				Properties: make(map[string]js.JSONSchema),
				Required:   make([]string, 0),
			},
		}

		for compIdx, component := range typ.TupleElems {
			elem, err := createSchemaItem(typ.TupleRawNames[compIdx], compIdx, component)
			if err != nil {
				return elem, err
			}
			schema.ObjectItem.Properties[typ.TupleRawNames[compIdx]] = elem
			schema.ObjectItem.Required = append(schema.ObjectItem.Required, typ.TupleRawNames[compIdx])
		}

		return schema, nil
	case abi.BoolTy:
		return js.JSONSchema{
			Type:         js.ItemTypeBoolean,
			Title:        name,
			InternalType: typ.String(),
		}, nil
	case abi.BytesTy, abi.FixedBytesTy, abi.FunctionTy:
		return js.JSONSchema{
			Type:         js.ItemTypeString,
			Title:        name,
			InternalType: typ.String(),
		}, nil
	case abi.IntTy, abi.UintTy, abi.FixedPointTy:
		return js.JSONSchema{
			Type:         js.ItemTypeNumber,
			Title:        name,
			InternalType: typ.String(),
		}, nil
	default:
		return js.JSONSchema{}, errors.Errorf("unknown argument type: %d", typ.T)
	}
}
