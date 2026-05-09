package migrations

import (
	"context"
	"database/sql"

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

func RollbackLast(ctx context.Context, db *sql.DB, migrationsDir string) error {
	if err := goose.DownContext(ctx, db, migrationsDir); err != nil {
		return err
	}
	return nil
}
