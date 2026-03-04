# DipDup Indexer SDK

SDK for building indexers by [DipDup](https://dipdup.io). Provides a set of packages for constructing modular, flow-based indexing pipelines.

## Architecture

The SDK follows a **flow-based programming (FBP)** approach: independent asynchronous modules communicate through typed inputs and outputs, forming a processing pipeline (workflow).

```
┌──────────┐     ┌──────────┐     ┌──────────┐
│  Source  │────>│ Process  │────>│  Sink    │
│ (gRPC,   │     │ (custom  │     │ (storage,│
│  cron)   │     │  logic)  │     │  printer)│
└──────────┘     └──────────┘     └──────────┘
```

## Packages

| Package | Description |
|---------|-------------|
| [`pkg/modules`](pkg/modules/) | Flow-based module system: interfaces, inputs/outputs, workflow orchestration |
| [`pkg/modules/grpc`](pkg/modules/grpc/) | gRPC client and server modules with subscription support |
| [`pkg/modules/cron`](pkg/modules/cron/) | Cron scheduler module |
| [`pkg/modules/zipper`](pkg/modules/zipper/) | Aggregates two input streams by key |
| [`pkg/modules/stopper`](pkg/modules/stopper/) | Graceful shutdown via context cancellation |
| [`pkg/modules/printer`](pkg/modules/printer/) | Debug module that logs received messages |
| [`pkg/storage`](pkg/storage/) | Abstract storage layer with generic interfaces |
| [`pkg/storage/postgres`](pkg/storage/postgres/) | PostgreSQL implementation (pgx + bun) |
| [`pkg/sync`](pkg/sync/) | Thread-safe generic `Map[K, V]` |
| [`pkg/contract`](pkg/contract/) | Contract ABI to JSON Schema transformer (EVM) |
| [`pkg/rlp`](pkg/rlp/) | RLP encoding/decoding for Ethereum logs |
| [`pkg/jsonschema`](pkg/jsonschema/) | JSON Schema types (Draft 2019-09) |

## Quick Start

### Module Workflow

```go
import (
    "github.com/dipdup-net/indexer-sdk/pkg/modules"
    "github.com/dipdup-net/indexer-sdk/pkg/modules/cron"
)

// Create modules
cronModule, _ := cron.NewModule(cfg.Cron)
customModule := NewCustomModule()

// Connect cron output to custom module input
modules.Connect(cronModule, customModule, "every_second", "input")

// Start workflow
ctx, cancel := context.WithCancel(context.Background())
cronModule.Start(ctx)
customModule.Start(ctx)
```

### Storage

```go
import (
    "github.com/dipdup-io/go-lib/config"
    "github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
)

cfg := config.Database{
    Host: "127.0.0.1", Port: 5432,
    User: "user", Password: "password", Database: "mydb",
}

storage, _ := postgres.Create(ctx, cfg, func(ctx context.Context, conn *database.Bun) error {
    _, err := conn.DB().NewCreateTable().Model((*MyModel)(nil)).IfNotExists().Exec(ctx)
    return err
})
defer storage.Close()
```

## Code Generation

The `cmd/dipdup-gen` tool generates boilerplate for new indexer projects from EVM contract ABIs.

```bash
go run ./cmd/dipdup-gen abi --input contract.json --output ./generated
```

## Examples

Working examples are available in the [`examples/`](examples/) directory:

- [`examples/cron`](examples/cron/) — cron scheduler usage
- [`examples/storage`](examples/storage/) — PostgreSQL storage layer
- [`examples/grpc`](examples/grpc/) — gRPC client/server with subscriptions
- [`examples/zipper`](examples/zipper/) — stream aggregation by key

## Installation

```bash
go get github.com/dipdup-net/indexer-sdk
```
