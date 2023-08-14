package postgres

import (
	"context"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/uptrace/bun"
)

var (
	errNilTx = errors.New("nil transaction pointer")
)

// Transaction -
type Transaction struct {
	tx *bun.Tx
}

// Flush -
func (t *Transaction) Flush(ctx context.Context) error {
	if t.tx == nil {
		return errNilTx
	}
	if err := t.tx.Commit(); err != nil {
		return err
	}

	return nil
}

// Add -
func (t *Transaction) Add(ctx context.Context, model any) error {
	if t.tx == nil {
		return errNilTx
	}

	_, err := t.tx.NewInsert().Model(model).Returning("id").Exec(ctx)
	return err
}

// Rollback -
func (t *Transaction) Rollback(ctx context.Context) error {
	if t.tx == nil {
		return errNilTx
	}
	return t.tx.Rollback()
}

// Close -
func (t *Transaction) Close(ctx context.Context) error {
	t.tx = nil
	return nil
}

// Update -
func (t *Transaction) Update(ctx context.Context, model any) error {
	if t.tx == nil {
		return errNilTx
	}

	_, err := t.tx.NewUpdate().Model(model).WherePK().Exec(ctx)
	return err
}

// BulkSave -
func (t *Transaction) BulkSave(ctx context.Context, models []any) error {
	if t.tx == nil {
		return errNilTx
	}

	if len(models) == 0 {
		return nil
	}

	_, err := t.tx.NewInsert().Model(&models).Returning("id").Exec(ctx)
	return err
}

// HandleError -
func (t *Transaction) HandleError(ctx context.Context, err error) error {
	processorErr := errors.Wrap(err, "transaction error")
	if err := t.Rollback(ctx); err != nil {
		return errors.Wrap(processorErr, errors.Wrap(err, "rollback").Error())
	}
	return processorErr
}

// Exec - executes query and returns the number of affected rows
func (t *Transaction) Exec(ctx context.Context, query string, params ...any) (int64, error) {
	if t.tx == nil {
		return 0, errNilTx
	}

	result, err := t.tx.NewRaw(query, params...).Exec(ctx)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// CopyFrom -
func (t *Transaction) CopyFrom(ctx context.Context, tableName string, data []storage.Copiable) error {
	if len(data) == 0 {
		return nil
	}
	if t.tx == nil {
		return errNilTx
	}

	stmt, err := t.tx.PrepareContext(ctx, pq.CopyIn(tableName, data[0].Columns()...))
	if err != nil {
		return err
	}

	for i := range data {
		if _, err := stmt.ExecContext(ctx, data[i].Flat()...); err != nil {
			return err
		}
	}

	if _, err := stmt.ExecContext(ctx); err != nil {
		return err
	}

	return stmt.Close()
}

// Tx - returns bun.Tx pointer for custom queries
func (t *Transaction) Tx() *bun.Tx {
	return t.tx
}
