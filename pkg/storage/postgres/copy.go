package postgres

import (
	"context"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/jackc/pgx/v5"
)

// SaveBulkWithCopy -
func SaveBulkWithCopy[T storage.Copiable](ctx context.Context, tx storage.Transaction, data []T, copyThreashold int) error {
	switch {
	case len(data) == 0:
		return nil
	case len(data) < copyThreashold:
		_, err := tx.Tx().NewInsert().Model(&data).Exec(ctx)
		return err
	default:
		_, err := tx.Pool().CopyFrom(
			ctx,
			pgx.Identifier{data[0].TableName()},
			data[0].Columns(),
			pgx.CopyFromSlice(len(data), func(i int) ([]any, error) {
				return data[i].Flat()
			},
			))
		return err
	}
}
