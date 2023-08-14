package postgres

import (
	"context"
	"time"

	"github.com/dipdup-net/go-lib/config"
	"github.com/dipdup-net/go-lib/database"
)

// Init - type of initialization function which called after creating connection to database. For example, can be used for indexes creation.
type Init func(ctx context.Context, conn *database.Bun) error

// Create - creates storage connection entity
func Create(ctx context.Context, cfg config.Database, init Init) (*Storage, error) {
	conn := database.NewBun()

	connectCtx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	if err := conn.Connect(connectCtx, cfg); err != nil {
		return nil, err
	}

	database.Wait(ctx, conn, time.Second*5)

	conn.DB().AddQueryHook(&logQueryHook{})

	if init != nil {
		if err := init(ctx, conn); err != nil {
			return nil, err
		}
	}

	return &Storage{
		Transactable: NewTransactable(conn),
		db:           conn,
	}, nil
}
