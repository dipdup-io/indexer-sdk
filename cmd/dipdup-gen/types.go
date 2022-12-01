package main

import (
	"fmt"
	"strings"

	js "github.com/dipdup-net/indexer-sdk/pkg/jsonschema"
)

const (
	BoolType    = "bool"
	DecimalType = "*decimal.Decimal"
	BigIntType  = "*big.Int"
	NumericType = "*postgres.Numeric"
	AddressType = "common.Address"
	StringType  = "string"
)

type field struct {
	RawName    string
	Name       string
	Type       string
	UnpackType string
	Index      int
	IsNested   bool
}

// Tags -
func (f *field) Tags() string {
	var tag strings.Builder

	tag.WriteString(`json:"`)
	tag.WriteString(f.RawName)
	tag.WriteString(`"`)

	switch f.Type {
	case BoolType:
		tag.WriteString(` pg:"default:false"`)
	case NumericType:
		tag.WriteString(` pg:",type:numeric"`)
	case DecimalType:
		tag.WriteString(` pg:",type:decimal(10,10)"`)
	default:
		if f.IsNested {
			tag.WriteString(` pg:",type:jsonb"`)
		}
	}

	return tag.String()
}

func (f field) isBigInt() bool {
	return f.UnpackType == BigIntType
}

func (f field) isDecimal() bool {
	return f.Type == DecimalType
}

func (f field) isAddress() bool {
	return f.UnpackType == AddressType
}

type goType struct {
	Name   string
	Type   string
	Fields []field

	HasDecimal bool
	HasBigInt  bool
	HasAddress bool
}

var abiToGo = map[string]string{
	"uint8":    "uint8",
	"uint16":   "uint16",
	"uint32":   "uint32",
	"uint64":   "uint64",
	"int8":     "int8",
	"int16":    "int16",
	"int32":    "int32",
	"int64":    "int64",
	"function": "[24]byte",
	"bool":     "bool",
	"string":   "string",
	"bytes":    "[]byte",
}

func generateTypes(name, postfix string, schema *js.JSONSchema, types map[string]goType) goType {
	name = buildName(name, postfix)
	resultType := goType{
		Fields: make([]field, 0),
		Name:   name,
	}

	switch {
	case strings.HasPrefix(schema.InternalType, "bytes"):
		count := strings.TrimPrefix(schema.InternalType, "bytes")
		resultType.Type = fmt.Sprintf("[%s]byte", count)
	case strings.HasPrefix(schema.InternalType, "uint") || strings.HasPrefix(schema.InternalType, "int"):
		resultType.Type = BigIntType
		resultType.HasBigInt = true
	case strings.HasPrefix(schema.InternalType, "fixed"):
		resultType.Type = DecimalType
		resultType.HasDecimal = true
	case schema.InternalType == "address":
		resultType.Type = StringType
	case schema.Type == "object":
		var props map[string]js.JSONSchema
		if schema != nil && schema.ObjectItem != nil {
			props = schema.Properties
		}

		for title, prop := range props {
			f := generateField(title, &prop, types)
			resultType.Fields = append(resultType.Fields, f)
			resultType.HasDecimal = resultType.HasDecimal || f.isDecimal()
			resultType.HasBigInt = resultType.HasBigInt || f.isBigInt()
			resultType.HasAddress = resultType.HasAddress || f.isAddress()
		}
	default:
		resultType.Type = "any"
	}

	if types != nil {
		types[name] = resultType
	}
	return resultType
}

func generateField(title string, prop *js.JSONSchema, types map[string]goType) field {
	f := field{
		RawName: title,
		Name:    buildName(title, ""),
		Index:   prop.Index,
	}
	if typ, ok := abiToGo[prop.InternalType]; ok {
		f.Type = typ
		f.UnpackType = f.Type
		return f
	}

	switch {
	case prop.InternalType == "address":
		f.Type = "string"
		f.UnpackType = AddressType
	case strings.HasPrefix(prop.InternalType, "bytes"):
		count := strings.TrimPrefix(prop.InternalType, "bytes")
		f.Type = fmt.Sprintf("[%s]byte", count)
		f.UnpackType = f.Type
	case strings.HasPrefix(prop.InternalType, "uint") || strings.HasPrefix(prop.InternalType, "int"):
		f.Type = NumericType
		f.UnpackType = BigIntType
	case strings.HasPrefix(prop.InternalType, "fixed"):
		f.Type = DecimalType
		f.UnpackType = f.Type
	case prop.Type == "array":
		goTyp := generateArrayItem(f.Name, prop, types)
		f.Type = fmt.Sprintf("[]%s", goTyp.Type)
		f.UnpackType = f.Type
	case prop.Type == "object":
		goTyp := generateTypes(f.Name, "Type", prop, types)
		f.Type = goTyp.Name
		f.UnpackType = f.Type
		f.IsNested = true
	default:
		f.Type = "any"
		f.UnpackType = f.Type
	}
	return f
}

func generateArrayItem(name string, prop *js.JSONSchema, types map[string]goType) goType {
	name = buildName(name, "Item")

	if prop.ArrayItem != nil {
		switch len(prop.Items) {
		case 0:
			return goType{
				Fields: make([]field, 0),
				Name:   name,
			}
		case 1:
			return generateTypes(name, "", &prop.Items[0], nil)
		default:
			resultType := goType{
				Fields: make([]field, 0),
				Name:   name,
			}

			for i := range prop.Items {
				f := generateField(prop.Items[i].Title, &prop.Items[i], types)
				resultType.Fields = append(resultType.Fields, f)
				resultType.HasDecimal = resultType.HasDecimal || f.isDecimal()
				resultType.HasBigInt = resultType.HasBigInt || f.isBigInt()
				resultType.HasAddress = resultType.HasAddress || f.isAddress()
			}

			return resultType
		}
	}

	return goType{
		Fields: make([]field, 0),
		Name:   name,
	}
}
