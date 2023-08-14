package postgres

import (
	"context"
	"database/sql"
	"errors"
	"reflect"

	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/uptrace/bun"
)

// Table - Postgres realization of Table interface
type Table[M storage.Model] struct {
	db *database.Bun
}

// NewTable - creates Table structure
func NewTable[M storage.Model](db *database.Bun) *Table[M] {
	return &Table[M]{db}
}

// Save - inserts row to table and returns id.
func (s *Table[M]) Save(ctx context.Context, m M) error {
	_, err := s.db.DB().NewInsert().Model(m).Returning("id").Exec(ctx)
	return err
}

// Update - updates table row by primary key.
func (s *Table[M]) Update(ctx context.Context, m M) error {
	_, err := s.db.DB().NewUpdate().Model(m).WherePK().Exec(ctx)
	return err
}

// List - returns array of rows
func (s *Table[M]) List(ctx context.Context, limit, offset uint64, order storage.SortOrder) ([]M, error) {
	var models []M
	query := s.db.DB().NewSelect().Model(&models)
	query = Pagination(query, limit, offset, order)

	err := query.Scan(ctx)
	return models, err
}

// GetByID - returns row by id
func (s *Table[M]) GetByID(ctx context.Context, id uint64) (m M, err error) {
	typ := reflect.TypeOf(m)
	if typ.Kind() == reflect.Ptr {
		value := reflect.New(typ.Elem())
		val := value.Interface()
		err = s.db.DB().NewSelect().Model(val).Where("id = ?", id).Scan(ctx)
		return val.(M), err
	} else {
		err = s.db.DB().NewSelect().Model(&m).Where("id = ?", id).Scan(ctx)
	}
	return
}

// IsNoRows - checks errors is pg.ErrNoRows
func (s *Table[M]) IsNoRows(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}

// DB - returns Postgres connection
func (s *Table[M]) DB() *bun.DB {
	return s.db.DB()
}

// CursorList - returns array of rows by cursor pagination
func (s *Table[M]) CursorList(ctx context.Context, id, limit uint64, order storage.SortOrder, cmp storage.Comparator) ([]M, error) {
	var models []M
	query := s.db.DB().NewSelect().Model(&models)
	query = CursorPagination(query, id, limit, order, cmp)

	err := query.Scan(ctx)
	return models, err
}

// LastID - returns last used id
func (s *Table[M]) LastID(ctx context.Context) (uint64, error) {
	var (
		m  M
		id uint64
	)
	err := s.DB().NewSelect().Model(m).ColumnExpr("max(id)").Scan(ctx, &id)
	return id, err
}
