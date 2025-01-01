package database

import (
	"database/sql"
	"embed"
	"fmt"
	"log"

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
		fmt.Println("Something went wrong with getting the database driver:", err)
		return err
	}

	migrations, err := iofs.New(migrationsFS, "migrations")

	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithInstance("iofs", migrations, "sqlite3", driver)

	if err != nil {
		fmt.Println("Something went wrong with initializing the migrator:", err)
		return err
	}

	err = m.Up()

	if err != nil && err != migrate.ErrNoChange {
		fmt.Println("Something went wrong with migrating:", err)
	}

	return nil
}
