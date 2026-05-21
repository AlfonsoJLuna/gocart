package db

import (
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"sort"
	"strings"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// migrate applies any pending migrations to the database.
//
// Migration files are numbered SQL scripts in the migrations/ directory
// (e.g. 001_initial_schema.sql). They are applied in filename order, one per
// version increment. The current schema version is tracked via SQLite's
// built-in PRAGMA user_version.
//
// Each migration runs inside a transaction together with the user_version
// increment, so a failure leaves the database completely unchanged.
func migrate(db *sql.DB) error {
	var version int
	if err := db.QueryRow("PRAGMA user_version;").Scan(&version); err != nil {
		return fmt.Errorf("reading user_version: %w", err)
	}

	entries, err := fs.ReadDir(migrationsFS, "migrations")
	if err != nil {
		return fmt.Errorf("reading migrations dir: %w", err)
	}

	var files []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".sql") {
			files = append(files, e.Name())
		}
	}
	sort.Strings(files)

	// files[i] corresponds to schema version i+1, so we skip everything
	// already covered by the current user_version.
	pending := files[version:]
	if len(pending) == 0 {
		log.Printf("Database schema is up to date (version %d).", version)
		return nil
	}

	log.Printf("Database schema is at version %d, %d migration(s) to apply.", version, len(pending))

	for i, name := range pending {
		nextVersion := version + i + 1
		log.Printf("Applying migration %d: %s", nextVersion, name)

		data, err := migrationsFS.ReadFile("migrations/" + name)
		if err != nil {
			return fmt.Errorf("reading migration %q: %w", name, err)
		}

		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("beginning transaction for %q: %w", name, err)
		}

		if _, err := tx.Exec(string(data)); err != nil {
			tx.Rollback()
			return fmt.Errorf("executing migration %q: %w", name, err)
		}

		// user_version is set inside the transaction so it rolls back together
		// with the migration SQL if anything goes wrong.
		if _, err := tx.Exec(fmt.Sprintf("PRAGMA user_version = %d;", nextVersion)); err != nil {
			tx.Rollback()
			return fmt.Errorf("updating user_version to %d: %w", nextVersion, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("committing migration %q: %w", name, err)
		}

		log.Printf("Migration %d applied successfully.", nextVersion)
	}

	return nil
}
