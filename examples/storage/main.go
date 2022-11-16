package main

import (
	"context"
	"log"

	"github.com/dipdup-net/go-lib/config"
	"github.com/dipdup-net/indexer-sdk/examples/storage/storage"
	"github.com/dipdup-net/indexer-sdk/examples/storage/storage/postgres"
)

func main() {
	cfg := config.Database{
		Host:     "127.0.0.1",
		Port:     5432,
		User:     "user",
		Password: "password",
		Database: "database",
	}
	ctx, cancel := context.WithCancel(context.Background())
	strg, err := postgres.Create(ctx, cfg)
	if err != nil {
		cancel()
		log.Panic(err)
	}

	if err := strg.Persons.Save(ctx, storage.Person{
		Name:  "John",
		Phone: "+1234567890",
	}); err != nil {
		cancel()
		log.Panic(err)
	}

	if err := strg.Persons.Save(ctx, storage.Person{
		Name:  "Mike",
		Phone: "+098765432",
	}); err != nil {
		cancel()
		log.Panic(err)
	}

	contact, err := strg.Persons.GetByName(ctx, "Mike")
	if err != nil {
		cancel()
		log.Panic(err)
	}
	log.Printf("%s: %s", contact.Name, contact.Phone)

	cancel()

	if err := strg.Close(); err != nil {
		log.Panic(err)
	}
}
