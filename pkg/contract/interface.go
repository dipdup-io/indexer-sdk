package contract

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// Transformer -
type Transformer interface {
	JSONSchema(data []byte) ([]byte, error)
}

// Type -
type Type string

// types
const (
	TypeEvm Type = "evm"
)

// JSONSchema -
func JSONSchema(typ Type, data []byte) ([]byte, error) {
	var t Transformer
	switch typ {
	case TypeEvm:
		t = NewEVM()
	default:
		return nil, errors.Errorf("unknown metadata transformer type: %s", typ)
	}

	return t.JSONSchema(data)
}
