package database

import (
	"database/sql"
	"testing"
)

func TestMigrateDatabase(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")

	if err != nil {
		t.Error("Expected nil, got", err)
	}

	err = MigrateDatabase(db)

	if err != nil {
		t.Error("Expected nil, got", err)
	}
}

func TestMigrateDatabaseError(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")

	if err != nil {
		t.Error("Expected nil, got", err)
	}

	db.Close()

	err = MigrateDatabase(db)

	if err == nil {
		t.Error("Expected error, got nil")
	}
}
