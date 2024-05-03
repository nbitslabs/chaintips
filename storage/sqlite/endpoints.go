package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/nbitslabs/chaintips/types"
)

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

func (d *SqliteBackend) GetEnabledEndpoint(id int) (types.Endpoint, error) {
	row := d.db.QueryRow("SELECT id, chain_id, host, port, protocol, username, password, enabled FROM endpoints WHERE id = ? AND enabled = 1;", id)
	return scanEndpoint(row)
}

func scanEndpoint(row *sql.Row) (types.Endpoint, error) {
	var endpoint types.Endpoint
	err := row.Scan(&endpoint.ID, &endpoint.ChainID, &endpoint.Host, &endpoint.Port, &endpoint.Protocol, &endpoint.Username, &endpoint.Password, &endpoint.Enabled)
	if err != nil {
		return types.Endpoint{}, fmt.Errorf("failed to scan endpoint: %w", err)
	}
	return endpoint, nil
}
