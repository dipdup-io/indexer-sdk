# Contract

Transforms contract metadata (ABI) into JSON Schema representation.

## Supported Types

```go
type Type string

const TypeEvm Type = "evm"
```

## Transformer Interface

```go
type Transformer interface {
    JSONSchema(data []byte) ([]byte, error)
}
```

## Usage

### Direct

```go
import "github.com/dipdup-net/indexer-sdk/pkg/contract"

schema, err := contract.JSONSchema(contract.TypeEvm, abiBytes)
```

### EVM Transformer

```go
evm := contract.NewEVM()
schema, err := evm.JSONSchema(abiBytes)
```

Takes an Ethereum contract ABI (JSON) as input and produces a JSON Schema describing all events and methods with their input/output parameters. Uses `go-ethereum/accounts/abi` for ABI parsing.
