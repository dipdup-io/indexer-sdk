package postgres

import (
	"context"

	"github.com/dipdup-io/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pkg/errors"
)

// Transactable - realization of Transactable interface for Postgres
type Transactable struct {
	db *database.Bun
}

// NewTransactable - creates Transactable structure
func NewTransactable(db *database.Bun) *Transactable {
	return &Transactable{db}
}

// BeginTransaction - opens atomic transaction to communication with Postgres
func (t *Transactable) BeginTransaction(ctx context.Context) (storage.Transaction, error) {
	bunConn, err := t.db.DB().Conn(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "tx connection")
	}

	var pgxConn *pgx.Conn
	if err := bunConn.Raw(func(c any) error {
		pgxConn = c.(*stdlib.Conn).Conn()
		return nil
	}); err != nil {
		return nil, errors.Wrap(err, "raw")
	}

	tx, err := bunConn.BeginTx(ctx, nil)
	if err != nil {
		return nil, errors.Wrap(err, "begin tx")
	}

	return &Transaction{
		conn:    bunConn,
		tx:      &tx,
		pgxConn: pgxConn,
	}, nil
}
