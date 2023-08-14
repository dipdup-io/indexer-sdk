package postgres

import (
	"context"
	"database/sql"

	"github.com/dipdup-net/go-lib/config"
	"github.com/dipdup-net/go-lib/database"
	models "github.com/dipdup-net/indexer-sdk/examples/storage/storage"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
	"github.com/uptrace/bun"
)

// Storage -
type Storage struct {
	*postgres.Storage

	Persons models.IPerson
}

// Create -
func Create(ctx context.Context, cfg config.Database) (*Storage, error) {
	strg, err := postgres.Create(ctx, cfg, initDatabase)
	if err != nil {
		return nil, err
	}

	return &Storage{
		Storage: strg,
		Persons: NewPerson(strg.Connection()),
	}, nil
}

func initDatabase(ctx context.Context, conn *database.Bun) error {
	// here you can create schemas, user grants or indexes

	return conn.DB().RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		_, err := tx.ExecContext(ctx, `CREATE INDEX IF NOT EXISTS persons_name ON name (name)`)
		return err
	})
}
