package jsonschema

const Draft201909 = "http://json-schema.org/draft/2019-09/schema#"

// Type -
type Type struct {
	Type      string      `json:"type"`
	Inputs    *JSONSchema `json:"inputs,omitempty"`
	Outputs   *JSONSchema `json:"outputs,omitempty"`
	Signature string      `json:"signature"`
}

// JSONSchema -
type JSONSchema struct {
	Schema      string `json:"$schema,omitempty"`
	ID          string `json:"$id,omitempty"`
	Comment     string `json:"$comment,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Default     any    `json:"default,omitempty"`
	Examples    []any  `json:"examples,omitempty"`
	Enum        []any  `json:"enum,omitempty"`

	Type        ItemType `json:"type"`
	*StringItem `json:",omitempty"`
	*NumberItem `json:",omitempty"`
	*ObjectItem `json:",omitempty"`
	*ArrayItem  `json:",omitempty"`
}

// ItemType -
type ItemType string

// item types
const (
	ItemTypeString  ItemType = "string"
	ItemTypeNumber  ItemType = "number"
	ItemTypeInteger ItemType = "integer"
	ItemTypeObject  ItemType = "object"
	ItemTypeBoolean ItemType = "boolean"
	ItemTypeArray   ItemType = "array"
	ItemTypeNull    ItemType = "null"
)

// StringItem -
type StringItem struct {
	MinLength int64      `json:"minLength,omitempty"`
	MaxLength int64      `json:"maxLength,omitempty"`
	Pattern   string     `json:"pattern,omitempty"`
	Format    FormatKind `json:"format,omitempty"`
}

// FormatKind -
type FormatKind string

// string format kinds
const (
	FormatKindDateTime FormatKind = "date-time"
	FormatKindDate     FormatKind = "date"
	FormatKindTime     FormatKind = "time"
	FormatKindDuration FormatKind = "duration"

	FormatKindEmail    FormatKind = "email"
	FormatKindIDNEmail FormatKind = "idn-email"

	FormatKindHostname    FormatKind = "hostname"
	FormatKindIDNHostname FormatKind = "idn-hostname"

	FormatKindIPv4 FormatKind = "ipv4"
	FormatKindIPv6 FormatKind = "ipv6"

	FormatKindUUID         FormatKind = "uuid"
	FormatKindURI          FormatKind = "uri"
	FormatKindURIReference FormatKind = "uri-reference"
	FormatKindIRI          FormatKind = "iri"
	FormatKindIRIReference FormatKind = "iri-reference"

	FormatKindURITemplate FormatKind = "uri-template"

	FormatKindRegex FormatKind = "regex"
)

// NumberItem -
type NumberItem struct {
	MultipleOf       int64 `json:"multipleOf,omitempty"`
	Minimum          int64 `json:"minimum,omitempty"`
	ExclusiveMinimum int64 `json:"exclusiveMinimum,omitempty"`
	Maximum          int64 `json:"maximum,omitempty"`
	ExclusiveMaximum int64 `json:"exclusiveMaximum,omitempty"`
}

// ObjectItem -
type ObjectItem struct {
	Properties           map[string]JSONSchema `json:"properties,omitempty"`
	AdditionalProperties bool                  `json:"additionalProperties,omitempty"`
	Required             []string              `json:"required,omitempty"`
	MinProperties        int64                 `json:"minProperties,omitempty"`
	MaxProperties        int64                 `json:"maxProperties,omitempty"`
}

// ArrayItem -
type ArrayItem struct {
	Items       []JSONSchema `json:"items,omitempty"`
	Contains    *JSONSchema  `json:"contains,omitempty"`
	MinContains int64        `json:"minContains,omitempty"`
	MaxContains int64        `json:"maxContains,omitempty"`
	MinItems    int64        `json:"minItems,omitempty"`
	MaxItems    int64        `json:"maxItems,omitempty"`
	UniqueItems bool         `json:"uniqueItems"`
}
