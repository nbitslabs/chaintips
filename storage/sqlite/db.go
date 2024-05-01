package sqlite

import (
	"database/sql"
	"embed"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/nbitslabs/chaintips/types"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embeddedMigrations embed.FS

type SqliteBackend struct {
	db *sql.DB
}

func NewSqliteBackend() (*SqliteBackend, error) {
	db, err := sql.Open("sqlite3", os.Getenv("DB_PATH"))
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

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

func (d *SqliteBackend) GetEnabledEndpoints(chainID int) ([]types.Endpoint, error) {
	rows, err := d.db.Query("SELECT id, chain_id, host, port, protocol, username, password, enabled FROM endpoints WHERE chain_id = ? AND enabled = 1;", chainID)
	if err != nil {
		return nil, fmt.Errorf("failed to query endpoints: %w", err)
	}

	defer rows.Close()

	var endpoints []types.Endpoint

	for rows.Next() {
		var endpoint types.Endpoint
		if err := rows.Scan(&endpoint.ID, &endpoint.ChainID, &endpoint.Host, &endpoint.Port, &endpoint.Protocol, &endpoint.Username, &endpoint.Password, &endpoint.Enabled); err != nil {
			return nil, fmt.Errorf("failed to scan endpoint: %w", err)
		}
		endpoints = append(endpoints, endpoint)
	}

	return endpoints, nil
}
