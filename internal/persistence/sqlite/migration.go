package sqlite

import (
	"embed"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var migrations embed.FS

// Migrate apply db migration, then return result.
func Migrate(db *sqlx.DB) error {
	goose.SetBaseFS(migrations)
	if err := goose.SetDialect("sqlite3"); err != nil {
		return err
	}
	if err := goose.Up(db.DB, "migrations"); err != nil {
		return fmt.Errorf("goose Up: %w", err)
	}
	return nil
}
