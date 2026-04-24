package migrations

import (
	"database/sql"
	"log"

	"github.com/pressly/goose/v3"
)

func RunMigrations(db *sql.DB, migrationsDir string) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	if err := goose.Up(db, migrationsDir); err != nil {
		return err
	}
	return nil
}

func RollbackLast(db *sql.DB, migrationsDir string) {
	if err := goose.Down(db, migrationsDir); err != nil {
		log.Fatalf("goose down failed: %v", err)
	}
}
