package sqlite

import (
	"fmt"

	"github.com/nbitslabs/chaintips/types"
)

func (d *SqliteBackend) UpsertChainTip(tip types.ChainTip) error {
	sqlStmt := `INSERT INTO chaintips (chain_id, endpoint_id, height, hash, branchlen, status) VALUES (?, ?, ?, ?, ?, ?)
				ON CONFLICT (height, hash, endpoint_id, chain_id) DO UPDATE SET branchlen = excluded.branchlen, status = excluded.status;`
	_, err := d.db.Exec(sqlStmt, tip.ChainID, tip.EndpointID, tip.Height, tip.Hash, tip.Branchlen, tip.Status)
	if err != nil {
		return fmt.Errorf("failed to upsert chaintip: %w", err)
	}
	return nil
}

func (d *SqliteBackend) GetUnlinkedChainTips(chainID int) ([]types.ChainTip, error) {
	rows, err := d.db.Query("SELECT id, chain_id, endpoint_id, height, hash, branchlen, status FROM chaintips WHERE chain_id = ? AND hash NOT IN (SELECT hash FROM blocks);", chainID)
	if err != nil {
		return nil, fmt.Errorf("failed to query chaintips: %w", err)
	}

	defer rows.Close()

	var tips []types.ChainTip

	for rows.Next() {
		var tip types.ChainTip
		if err := rows.Scan(&tip.ID, &tip.ChainID, &tip.EndpointID, &tip.Height, &tip.Hash, &tip.Branchlen, &tip.Status); err != nil {
			return nil, fmt.Errorf("failed to scan chaintip: %w", err)
		}
		tips = append(tips, tip)
	}

	return tips, nil
}
