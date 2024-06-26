package sqlite

import (
	"database/sql"
	"embed"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/nbitslabs/chaintips/types"
	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog/log"
)

//go:embed migrations/*.sql
var embeddedMigrations embed.FS

var logger = log.With().Str("module", "sqlite").Logger()

type SqliteBackend struct {
	db *sql.DB
}

func NewSqliteBackend() (*SqliteBackend, error) {
	path := os.Getenv("DB_PATH")
	if path == "" {
		return nil, fmt.Errorf("DB_PATH environment variable must be set")
	}

	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	log.Info().
		Str("path", path).
		Msg("Database opened")

	backend := &SqliteBackend{db: db}
	if err := backend.Migrate(); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return backend, nil
}

func (d *SqliteBackend) Close() error {
	return d.db.Close()
}

func (d *SqliteBackend) Migrate() error {
	goose.SetBaseFS(embeddedMigrations)
	if err := goose.SetDialect("sqlite3"); err != nil {
		return fmt.Errorf("failed to set goose dialect: %w", err)
	}

	if err := goose.Up(d.db, "migrations"); err != nil {
		return fmt.Errorf("failed to run goose up: %w", err)
	}
	return nil
}

func (d *SqliteBackend) GetChains() ([]types.Chain, error) {
	rows, err := d.db.Query("SELECT id, identifier, title, icon FROM chains;")
	if err != nil {
		return nil, fmt.Errorf("failed to query chains: %w", err)
	}

	defer rows.Close()

	var chains []types.Chain

	for rows.Next() {
		var chain types.Chain
		if err := rows.Scan(&chain.ID, &chain.Identifier, &chain.Title, &chain.Icon); err != nil {
			return nil, fmt.Errorf("failed to scan chain: %w", err)
		}
		chains = append(chains, chain)
	}

	return chains, nil
}
