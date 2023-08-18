package postgres

import (
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage"
)

// Storage - default storage structure containing Transactable interface and connection to database
type Storage struct {
	Transactable storage.Transactable

	db *database.Bun
}

// Close - closes storage
func (s *Storage) Close() error {
	return s.db.Close()
}

// Connection - returns connection structure
func (s *Storage) Connection() *database.Bun {
	return s.db
}
