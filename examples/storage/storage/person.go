package storage

import (
	"context"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
)

// IPerson -
type IPerson interface {
	storage.Table[Person]

	// add here custom queries
	GetByName(ctx context.Context, name string) (Person, error)
}

// Person -
type Person struct {
	// nolint
	tableName struct{} `pg:"persons"`

	ID    uint64
	Name  string
	Phone string
}

// TableName -
func (Person) TableName() string {
	return "persons"
}
