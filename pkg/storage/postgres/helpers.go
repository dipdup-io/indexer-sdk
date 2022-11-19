package postgres

import (
	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

// Pagination - adds limit, offset and sort to query. Query being like this:
//
//	query.Limit(limit).Offset(offset).Order("id ?", order)
func Pagination(query *orm.Query, limit, offset uint64, order storage.SortOrder) *orm.Query {
	if limit == 0 {
		limit = 10
	}

	query.Limit(int(limit)).Offset(int(offset))

	switch order {
	case storage.SortOrderAsc:
		query.Order("id asc")
	case storage.SortOrderDesc:
		query.Order("id desc")
	default:
		query.Order("id asc")
	}
	return query
}

// CursorPagination - adds limit, id where clause and sort to query. Query being like this:
//
//	query.Where("id > ?", id).Limit(limit).Order("id ?", order)
func CursorPagination(query *orm.Query, id, limit uint64, order storage.SortOrder, cmp storage.Comparator) *orm.Query {
	if id > 0 {
		query.Where("id ? ?", pg.Safe(cmp.String()), id)
	}

	if limit == 0 {
		limit = 10
	}

	query.Limit(int(limit))

	switch order {
	case storage.SortOrderAsc:
		query.Order("id asc")
	case storage.SortOrderDesc:
		query.Order("id desc")
	default:
		query.Order("id asc")
	}
	return query
}

// In - adds IN clause to query:
//
//	WHERE field IN (1,2,3)
//
// If length of array equals 0 condition skips.
func In[M any](query *orm.Query, field string, arr []M) *orm.Query {
	if len(arr) == 0 {
		return query
	}

	query.Where("? IN (?)", pg.Ident(field), pg.In(arr))

	return query
}
