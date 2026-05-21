// Package db handles database initialization and migrations.
package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite" // registers the sqlite driver with database/sql
)

// Open opens (or creates) the SQLite database at path, configures connection
// pragmas, runs any pending migrations, and returns a ready-to-use *sql.DB.
func Open(path string) (*sql.DB, error) {
	log.Println("Opening database...")

	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("opening database %q: %w", path, err)
	}

	// These pragmas are not persisted across connections in all SQLite builds,
	// so we set them on every open.
	pragmas := []string{
		"PRAGMA journal_mode = WAL;",      // better concurrent read performance
		"PRAGMA foreign_keys = ON;",       // enforce FK constraints (off by default)
		"PRAGMA busy_timeout = 5000;",     // wait up to 5s instead of failing on lock
		"PRAGMA synchronous = NORMAL;",    // safe with WAL, faster than FULL
		"PRAGMA cache_size = -20000;",     // 20 MB page cache
		"PRAGMA temp_store = MEMORY;",     // keep temp tables in memory
		"PRAGMA optimize;",                // update query planner statistics from last session
	}
	for _, p := range pragmas {
		if _, err := db.Exec(p); err != nil {
			db.Close()
			return nil, fmt.Errorf("setting pragma %q: %w", p, err)
		}
	}

	log.Println("Database opened successfully.")

	if err := migrate(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("running migrations: %w", err)
	}

	return db, nil
}

// IsSeeded returns true if the database has already been seeded.
func IsSeeded(db *sql.DB) (bool, error) {
	var value string
	err := db.QueryRow(`SELECT value FROM config WHERE key = 'seeded'`).Scan(&value)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("checking seed status: %w", err)
	}
	return value == "true", nil
}
 
// MarkSeeded writes a 'seeded' key to the config table so SeedAll is never
// run again on this database.
func MarkSeeded(db *sql.DB) error {
	_, err := db.Exec(`
		INSERT INTO config (key, value) VALUES ('seeded', 'true')
		ON CONFLICT (key) DO UPDATE SET value = 'true', updated_at = unixepoch()`)
	if err != nil {
		return fmt.Errorf("marking database as seeded: %w", err)
	}
	return nil
}
