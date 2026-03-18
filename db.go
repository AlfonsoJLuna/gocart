package main

import (
	"fmt"
	"log"

	"gocart/models"
	"gocart/seeds"

	"go.etcd.io/bbolt"
)

const schemaVersion = 1

func dbInit(path string) (*bbolt.DB, error) {
	log.Println("Opening database...")

	db, err := bbolt.Open(path, 0600, nil)
	if err != nil {
		return nil, fmt.Errorf("opening database %q: %w", path, err)
	}

	log.Println("Database opened successfully.")

	meta, err := models.MetaRead(db)

	if err != nil || !meta.Initialized {
		log.Println("Database is empty. Initializing database...")

		if err := seeds.SeedAll(db); err != nil {
			db.Close()
			return nil, fmt.Errorf("seeding database: %w", err)
		}

		log.Println("Writing database metadata...")
		if err := models.MetaWrite(db, &models.Meta{
			SchemaVersion: schemaVersion,
			Initialized:   true,
		}); err != nil {
			db.Close()
			return nil, fmt.Errorf("writing meta: %w", err)
		}

		log.Println("Database initialized successfully.")
	}

	// To do: if schemaVersion > meta.SchemaVersion, apply migrations

	return db, nil
}
