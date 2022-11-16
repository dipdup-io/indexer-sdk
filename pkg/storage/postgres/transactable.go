package postgres

import (
	"context"

	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage"
)

// Transactable - realization of Transactable interface for Postgres
type Transactable struct {
	db *database.PgGo
}

// NewTransactable - creates Transactable structure
func NewTransactable(db *database.PgGo) *Transactable {
	return &Transactable{db}
}

// BeginTransaction - opens atomic transaction to communication with Postgres
func (t *Transactable) BeginTransaction(ctx context.Context) (storage.Transaction, error) {
	tx, err := t.db.DB().BeginContext(ctx)
	if err != nil {
		return nil, err
	}
	return &Transaction{tx}, nil
}
