package database

import (
	"database/sql"
	"embed"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/*sql
var migrationsFS embed.FS

func MigrateDatabase(db *sql.DB) error {
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{
		MigrationsTable: "migrations",
	})

	if err != nil {
		return fmt.Errorf("something went wrong with getting the database driver: %s", err)
	}

	migrations, err := iofs.New(migrationsFS, "migrations")

	if err != nil {
		return fmt.Errorf("something went wrong with getting the migrations: %s", err)
	}

	m, err := migrate.NewWithInstance("iofs", migrations, "sqlite3", driver)

	if err != nil {
		return fmt.Errorf("something went wrong with initializing the migrator: %s", err)
	}

	err = m.Up()

	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("something went wrong with migrating: %s", err)
	}

	return nil
}
