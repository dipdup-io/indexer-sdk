# JSON Schema

Go types for representing [JSON Schema Draft 2019-09](http://json-schema.org/draft/2019-09/schema#). Used internally by the [`contract`](../contract/) package to generate schemas from contract ABIs.

## Types

### Type

Top-level contract type schema:

```go
type Type struct {
    Type      string
    Inputs    *JSONSchema
    Outputs   *JSONSchema
    Signature string
}
```

### JSONSchema

Core schema representation supporting all Draft 2019-09 features:

```go
type JSONSchema struct {
    Schema       string        // $schema URI
    ID           string        // $id
    Title        string
    Description  string
    Type         ItemType      // string, number, integer, object, boolean, array, null
    InternalType string        // original type name (e.g., Solidity type)
    Index        int

    // Composition
    OneOf []*JSONSchema
    AllOf []*JSONSchema
    AnyOf []*JSONSchema
    Not   *JSONSchema

    // Type-specific constraints
    StringItem                // MinLength, MaxLength, Pattern, Format
    NumberItem                // Minimum, Maximum, MultipleOf, etc.
    ObjectItem                // Properties, Required, AdditionalProperties
    ArrayItem                 // Items, MinItems, MaxItems, UniqueItems
}
```

### ItemType

```go
const (
    ItemTypeString  ItemType = "string"
    ItemTypeNumber  ItemType = "number"
    ItemTypeInteger ItemType = "integer"
    ItemTypeObject  ItemType = "object"
    ItemTypeBoolean ItemType = "boolean"
    ItemTypeArray   ItemType = "array"
    ItemTypeNull    ItemType = "null"
)
```

### FormatKind

String format validation:

```go
const (
    FormatKindDateTime    FormatKind = "date-time"
    FormatKindDate        FormatKind = "date"
    FormatKindEmail       FormatKind = "email"
    FormatKindIPv4        FormatKind = "ipv4"
    FormatKindIPv6        FormatKind = "ipv6"
    FormatKindUUID        FormatKind = "uuid"
    FormatKindURI         FormatKind = "uri"
    FormatKindRegex       FormatKind = "regex"
    // ... and more
)
```
