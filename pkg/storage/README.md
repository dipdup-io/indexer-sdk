# Storage package

Abstract layer of data storage is described in the package. Also package contains Postgres realization of the abstract layer.

## Abstract layer

Basic interfaces are described in the package. It describes data and communication with storage. Descriptions of interfaces can be found below.

### Model

Basic data interface. Its declaration is:

```go
type Model interface {
	TableName() string
}
```

It declares one method `TableName` which returns name of collection or table where it will be stored.

### Table

Interface describes the default way how developer can communicate with data storage. Its declaration is:

```go
type Table[M Model] interface {
	GetByID(ctx context.Context, id uint64) (M, error)
	Save(ctx context.Context, m M) error
	Update(ctx context.Context, m M) error
	List(ctx context.Context, limit, offset uint64, order SortOrder) ([]M, error)

	DB() *pg.DB
	IsNoRows(err error) bool
}
```

Interface is generic and allows communication with only one `Model`. But you can add needed methods to your storage implemention over that interface.

### Transactable

Interface allows to begin atomic transaction operation. Its declaration is:

```go
type Transactable interface {
	BeginTransaction(ctx context.Context) (Transaction, error)
}
``` 

### Transaction

Atomic transaction operation interface. Its declaration is:

```go
type Transaction interface {
	Flush(ctx context.Context) error
	Add(ctx context.Context, model any) error
	Update(ctx context.Context, model any) error
	Rollback(ctx context.Context) error
	BulkSave(ctx context.Context, models []any) error
	Close(ctx context.Context) error
	HandleError(ctx context.Context, err error) error
}
```

## Implementations

Now only one implementation is avalable. It's Postgres. It can be found [here](/pkg/storage/postgres/).

## Usage

Example of usage can be found [here](/examples/storage/)