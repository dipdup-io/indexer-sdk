package postgres

import (
	"context"

	"github.com/go-pg/pg/v10"
	"github.com/pkg/errors"
)

var (
	errNilTx = errors.New("nil transaction pointer")
)

// Transaction -
type Transaction struct {
	tx *pg.Tx
}

// Flush -
func (t *Transaction) Flush(ctx context.Context) error {
	if t.tx == nil {
		return errNilTx
	}
	if err := t.tx.CommitContext(ctx); err != nil {
		return err
	}

	return nil
}

// Add -
func (t *Transaction) Add(ctx context.Context, model any) error {
	if t.tx == nil {
		return errNilTx
	}

	_, err := t.tx.ModelContext(ctx, model).Returning("id").Insert()
	return err
}

// Rollback -
func (t *Transaction) Rollback(ctx context.Context) error {
	if t.tx == nil {
		return errNilTx
	}
	return t.tx.RollbackContext(ctx)
}

// Close -
func (t *Transaction) Close(ctx context.Context) error {
	if t.tx == nil {
		return errNilTx
	}

	if err := t.tx.CloseContext(ctx); err != nil {
		return err
	}

	t.tx = nil
	return nil
}

// Update -
func (t *Transaction) Update(ctx context.Context, model any) error {
	if t.tx == nil {
		return errNilTx
	}

	_, err := t.tx.ModelContext(ctx, model).WherePK().Update()
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

	_, err := t.tx.ModelContext(ctx, &models).Returning("id").Insert()
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
