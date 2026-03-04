# Storage: PostgreSQL

PostgreSQL implementation of the [storage](../) interfaces. Built on top of [pgx](https://github.com/jackc/pgx) driver and [bun](https://github.com/uptrace/bun) ORM.

## Components

### Storage

Top-level struct that holds the database connection and provides `Transactable`:

```go
type Storage struct {
    Transactable storage.Transactable
}
```

Create a storage instance with `Create`:

```go
func Create(ctx context.Context, cfg config.Database, init Init) (*Storage, error)
```

`Init` is a callback invoked after connection is established — use it for schema migrations, table creation, etc.:

```go
storage, err := postgres.Create(ctx, cfg, func(ctx context.Context, conn *database.Bun) error {
    _, err := conn.DB().NewCreateTable().Model((*Person)(nil)).IfNotExists().Exec(ctx)
    return err
})
```

Methods:

- `Connection() *database.Bun` — returns the underlying bun connection
- `Close() error` — closes the storage

### Table[M storage.Model]

Generic implementation of `storage.Table[M]` for a single model type:

```go
persons := postgres.NewTable[Person](storage.Connection())

person, err := persons.GetByID(ctx, 42)
err = persons.Save(ctx, Person{Name: "John"})
err = persons.Update(ctx, person)
list, err := persons.List(ctx, 10, 0, storage.SortOrderDesc)
items, err := persons.CursorList(ctx, lastID, 100, storage.SortOrderAsc, storage.ComparatorGt)
lastID, err := persons.LastID(ctx)
```

`DB() *bun.DB` provides direct access to the bun database for building custom queries.

### Transaction

Implements `storage.Transaction`. Obtained via `Transactable.BeginTransaction`:

```go
tx, err := storage.Transactable.BeginTransaction(ctx)
if err != nil {
    return err
}
defer tx.Close(ctx)

if err := tx.Add(ctx, model); err != nil {
    return tx.HandleError(ctx, err)
}
if err := tx.Flush(ctx); err != nil {
    return tx.HandleError(ctx, err)
}
```

Methods:
- `Add` / `Update` — insert or update a single model
- `BulkSave` — insert multiple models
- `Flush` — commit the transaction
- `Rollback` — rollback the transaction
- `Exec` — execute raw SQL
- `Tx() *bun.Tx` — access the underlying bun transaction
- `Pool() *pgx.Conn` — access the underlying pgx connection

### Numeric

Wrapper around `big.Int` for PostgreSQL `numeric` columns:

```go
n := postgres.NewNumeric(bigInt)
n = postgres.FromInt64(42)
n, err = postgres.FromString("123456789012345678901234567890")

result := n.Add(other)
result = n.Sub(other)
result = n.Mul(other)
result = n.Div(other)
result = n.Neg()

i64 := n.ToInt64()
u64 := n.ToUInt64()
```

Implements `driver.Valuer` and `sql.Scanner` for seamless database integration.

### Helper Functions

Query builders for common patterns:

```go
// Offset-based pagination
query = postgres.Pagination(query, limit, offset, order)

// Cursor-based pagination
query = postgres.CursorPagination(query, id, limit, order, comparator)

// WHERE field IN (...)
query = postgres.In(query, "status", []string{"active", "pending"})

// WHERE field = ANY(...)
query = postgres.Any(query, "id", ids)
```

### Bulk Copy

Efficient bulk insert using PostgreSQL `COPY FROM`:

```go
err := postgres.SaveBulkWithCopy(ctx, tx, data, 100)
```

Falls back to regular insert if data size is below threshold.
