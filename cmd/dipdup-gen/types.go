package main

import (
	"fmt"
	"strings"

	js "github.com/dipdup-net/indexer-sdk/pkg/jsonschema"
)

const (
	BigIntType  = "*big.Int"
	DecimalType = "*decimal.Decimal"
	StringType  = "string"
)

type field struct {
	RawName  string
	Name     string
	Type     string
	IsNested bool
}

// Tags -
func (f *field) Tags() string {
	var tag strings.Builder

	tag.WriteString(`json:"`)
	tag.WriteString(f.RawName)
	tag.WriteString(`"`)

	switch f.Type {
	case "bool":
		tag.WriteString(` pg:"default:false"`)
	case BigIntType:
		tag.WriteString(` pg:",type:numeric"`)
	default:
		if f.IsNested {
			tag.WriteString(` pg:",type:jsonb"`)
		}
	}

	return tag.String()
}

func (f field) isBigInt() bool {
	return strings.Contains(f.Type, BigIntType)
}

func (f field) isDecimal() bool {
	return strings.Contains(f.Type, DecimalType)
}

type goType struct {
	Name   string
	Type   string
	Fields []field

	HasBigInt  bool
	HasDecimal bool
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
	"address":  "string",
	"bytes":    "[]byte",
}

func generateTypes(name, postfix string, schema *js.JSONSchema, types map[string]goType) goType {
	name = buildName(name, postfix)
	resultType := goType{
		Fields: make([]field, 0),
		Name:   name,
	}

	switch {
	case strings.HasPrefix(schema.Comment, "bytes"):
		count := strings.TrimPrefix(schema.Comment, "bytes")
		resultType.Type = fmt.Sprintf("[%s]byte", count)
	case strings.HasPrefix(schema.Comment, "uint") || strings.HasPrefix(schema.Comment, "int"):
		resultType.Type = BigIntType
		resultType.HasBigInt = true
	case strings.HasPrefix(schema.Comment, "fixed"):
		resultType.Type = DecimalType
		resultType.HasDecimal = true
	case schema.Comment == "address":
		resultType.Type = StringType
	case schema.Type == "object":
		var props map[string]js.JSONSchema
		if schema != nil && schema.ObjectItem != nil {
			props = schema.Properties
		}

		for title, prop := range props {
			f := generateField(title, &prop, types)
			resultType.Fields = append(resultType.Fields, f)
			resultType.HasBigInt = resultType.HasBigInt || f.isBigInt()
			resultType.HasDecimal = resultType.HasDecimal || f.isDecimal()
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
	}
	if typ, ok := abiToGo[prop.Comment]; ok {
		f.Type = typ
		return f
	}

	switch {
	case strings.HasPrefix(prop.Comment, "bytes"):
		count := strings.TrimPrefix(prop.Comment, "bytes")
		f.Type = fmt.Sprintf("[%s]byte", count)
	case strings.HasPrefix(prop.Comment, "uint") || strings.HasPrefix(prop.Comment, "int"):
		f.Type = BigIntType
	case strings.HasPrefix(prop.Comment, "fixed"):
		f.Type = DecimalType
	case prop.Type == "array":
		goTyp := generateArrayItem(f.Name, prop, types)
		f.Type = fmt.Sprintf("[]%s", goTyp.Type)
	case prop.Type == "object":
		goTyp := generateTypes(f.Name, "Type", prop, types)
		f.Type = goTyp.Name
		f.IsNested = true
	default:
		f.Type = "any"
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
				resultType.HasBigInt = resultType.HasBigInt || f.isBigInt()
				resultType.HasDecimal = resultType.HasDecimal || f.isDecimal()
			}

			return resultType
		}
	}

	return goType{
		Fields: make([]field, 0),
		Name:   name,
	}
}
