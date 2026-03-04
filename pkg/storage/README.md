# Storage

Abstract storage layer for data persistence. Defines generic interfaces for CRUD operations, transactions, and cursor-based pagination. Includes a PostgreSQL implementation.

## Interfaces

### Model

Every entity stored in the database must implement `Model`:

```go
type Model interface {
    TableName() string
}
```

Example:

```go
type Person struct {
    bun.BaseModel `bun:"table:persons"`
    ID    uint64 `bun:"id,pk,autoincrement"`
    Name  string `bun:"name"`
    Phone string `bun:"phone"`
}

func (Person) TableName() string { return "persons" }
```

### Table[M Model]

Generic interface for typed access to a single table:

```go
type Table[M Model] interface {
    GetByID(ctx context.Context, id uint64) (M, error)
    Save(ctx context.Context, m M) error
    Update(ctx context.Context, m M) error
    List(ctx context.Context, limit, offset uint64, order SortOrder) ([]M, error)
    CursorList(ctx context.Context, id, limit uint64, order SortOrder, cmp Comparator) ([]M, error)
    LastID(ctx context.Context) (uint64, error)
    IsNoRows(err error) bool
}
```

### Transactable

Opens atomic transactions:

```go
type Transactable interface {
    BeginTransaction(ctx context.Context) (Transaction, error)
}
```

### Transaction

Atomic operations within a single transaction:

```go
type Transaction interface {
    Flush(ctx context.Context) error
    Add(ctx context.Context, model any) error
    Update(ctx context.Context, model any) error
    Rollback(ctx context.Context) error
    BulkSave(ctx context.Context, models []any) error
    Close(ctx context.Context) error
    HandleError(ctx context.Context, err error) error
    Exec(ctx context.Context, query string, params ...any) (int64, error)
    Tx() *bun.Tx
    Pool() *pgx.Conn
}
```

### Copiable

For bulk inserts using PostgreSQL `COPY FROM`:

```go
type Copiable interface {
    Flat() ([]any, error)
    Columns() []string
    TableName() string
}
```

## Types

### SortOrder

```go
type SortOrder string

const (
    SortOrderAsc  SortOrder = "asc"
    SortOrderDesc SortOrder = "desc"
)
```

### Comparator

Used for cursor-based pagination:

```go
type Comparator uint64

const (
    ComparatorEq  Comparator = iota // =
    ComparatorNeq                   // !=
    ComparatorLt                    // <
    ComparatorLte                   // <=
    ComparatorGt                    // >
    ComparatorGte                   // >=
)
```

## PostgreSQL Implementation

See [`postgres/`](postgres/) for the PostgreSQL implementation details.

## Usage

```go
import (
    "github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
)

// Create storage with initialization function
storage, err := postgres.Create(ctx, cfg, func(ctx context.Context, conn *database.Bun) error {
    _, err := conn.DB().NewCreateTable().Model((*Person)(nil)).IfNotExists().Exec(ctx)
    return err
})
defer storage.Close()

// Use typed tables
persons := postgres.NewTable[Person](storage.Connection())
err = persons.Save(ctx, Person{Name: "John", Phone: "+1234567890"})

// Cursor pagination
items, err := persons.CursorList(ctx, lastID, 100, storage.SortOrderAsc, storage.ComparatorGt)

// Transactions
tx, err := storage.Transactable.BeginTransaction(ctx)
_ = tx.Add(ctx, person1)
_ = tx.Add(ctx, person2)
_ = tx.Flush(ctx)
_ = tx.Close(ctx)
```

Full example: [`examples/storage/`](/examples/storage/)
