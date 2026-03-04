# RLP

RLP (Recursive Length Prefix) encoding and decoding for Ethereum log data. Uses `go-ethereum/rlp` under the hood.

## Types

```go
type Log struct {
    Topics [][32]byte
    Data   []byte
}
```

## Functions

### Encode

```go
import (
    "github.com/dipdup-net/indexer-sdk/pkg/rlp"
    "github.com/ethereum/go-ethereum/core/types"
)

encoded, err := rlp.EncodeLogData(ethereumLog)
```

Converts an `ethereum/types.Log` into RLP-encoded bytes.

### Decode

```go
log, err := rlp.DecodeLogData(encoded)
```

Decodes RLP bytes back into `rlp.Log` with topics and data.
