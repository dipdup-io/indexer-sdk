package storage

import (
	"context"

	"github.com/uptrace/bun"
)

// SortOrder - asc or desc
type SortOrder string

// sort orders
const (
	SortOrderAsc  SortOrder = "asc"
	SortOrderDesc SortOrder = "desc"
)

// Comparator - enum for cursor pagination
type Comparator uint64

const (
	ComparatorEq Comparator = iota
	ComparatorNeq
	ComparatorLt
	ComparatorLte
	ComparatorGt
	ComparatorGte
)

// String -
func (c Comparator) String() string {
	switch c {
	case ComparatorEq:
		return "="
	case ComparatorGt:
		return ">"
	case ComparatorGte:
		return ">="
	case ComparatorLt:
		return "<"
	case ComparatorLte:
		return "<="
	case ComparatorNeq:
		return "!="
	default:
		return ""
	}
}

// Table - interface to communication with one data type (like table, collection, etc)
type Table[M Model] interface {
	GetByID(ctx context.Context, id uint64) (M, error)
	Save(ctx context.Context, m M) error
	Update(ctx context.Context, m M) error
	List(ctx context.Context, limit, offset uint64, order SortOrder) ([]M, error)
	CursorList(ctx context.Context, id, limit uint64, order SortOrder, cmp Comparator) ([]M, error)
	LastID(ctx context.Context) (uint64, error)

	IsNoRows(err error) bool
}

// Transactable - interface which allows to begin atomic transaction operation
type Transactable interface {
	BeginTransaction(ctx context.Context) (Transaction, error)
}

// Transaction - atomic transaction operation interface
type Transaction interface {
	Flush(ctx context.Context) error
	Add(ctx context.Context, model any) error
	Update(ctx context.Context, model any) error
	Rollback(ctx context.Context) error
	BulkSave(ctx context.Context, models []any) error
	Close(ctx context.Context) error
	HandleError(ctx context.Context, err error) error
	Exec(ctx context.Context, query string, params ...any) (int64, error)
	CopyFrom(ctx context.Context, tableName string, data []Copiable) error
	Tx() *bun.Tx
}

// Model - general data type interface
type Model interface {
	TableName() string
}

// Copiable -
type Copiable interface {
	Columns() []string
	Flat() []any
}
