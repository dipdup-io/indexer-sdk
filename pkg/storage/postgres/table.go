package postgres

import (
	"context"
	"errors"
	"reflect"

	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/go-pg/pg/v10"
)

// Table - Postgres realization of Table interface
type Table[M storage.Model] struct {
	db *database.PgGo
}

// NewTable - creates Table structure
func NewTable[M storage.Model](db *database.PgGo) *Table[M] {
	return &Table[M]{db}
}

// Save - inserts row to table and returns id.
func (s *Table[M]) Save(ctx context.Context, m M) error {
	_, err := s.db.DB().ModelContext(ctx, m).Returning("id").Insert()
	return err
}

// Update - updates table row by primary key.
func (s *Table[M]) Update(ctx context.Context, m M) error {
	_, err := s.db.DB().ModelContext(ctx, m).WherePK().Update()
	return err
}

// List - returns array of rows
func (s *Table[M]) List(ctx context.Context, limit, offset uint64, order storage.SortOrder) ([]M, error) {
	var models []M
	query := s.db.DB().ModelContext(ctx, &models)
	query = Pagination(query, limit, offset, order)

	err := query.Select(&models)
	return models, err
}

// GetByID - returns row by id
func (s *Table[M]) GetByID(ctx context.Context, id uint64) (m M, err error) {
	typ := reflect.TypeOf(m)
	if typ.Kind() == reflect.Ptr {
		value := reflect.New(typ.Elem())
		val := value.Interface()
		err = s.db.DB().ModelContext(ctx, val).Where("id = ?", id).First()
		return val.(M), err
	} else {
		err = s.db.DB().ModelContext(ctx, &m).Where("id = ?", id).First()
	}
	return
}

// IsNoRows - checks errors is pg.ErrNoRows
func (s *Table[M]) IsNoRows(err error) bool {
	return errors.Is(err, pg.ErrNoRows)
}

// DB - returns Postgres connection
func (s *Table[M]) DB() *pg.DB {
	return s.db.DB()
}

// CursorList - returns array of rows by cursor pagination
func (s *Table[M]) CursorList(ctx context.Context, id, limit uint64, order storage.SortOrder, cmp storage.Comparator) ([]M, error) {
	var models []M
	query := s.db.DB().ModelContext(ctx, &models)
	query = CursorPagination(query, id, limit, order, cmp)

	err := query.Select(&models)
	return models, err
}
