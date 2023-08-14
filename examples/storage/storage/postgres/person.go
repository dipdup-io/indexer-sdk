package postgres

import (
	"context"

	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/examples/storage/storage"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
)

// Person -
type Person struct {
	*postgres.Table[storage.Person]
}

// NewPerson -
func NewPerson(db *database.Bun) *Person {
	return &Person{
		Table: postgres.NewTable[storage.Person](db),
	}
}

// GetByName-
func (p *Person) GetByName(ctx context.Context, name string) (storage.Person, error) {
	var person storage.Person
	err := p.DB().NewSelect().Model(&person).Where("name = ?", name).Scan(ctx)
	return person, err
}
