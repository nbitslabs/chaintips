package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/nbitslabs/chaintips/types"
)

func (d *SqliteBackend) FirstBlock(chainID int) (types.Block, error) {
	row := d.db.QueryRow("SELECT * FROM blocks WHERE chain_id = ? ORDER BY height ASC LIMIT 1", chainID)
	return scanBlock(row)
}

func (d *SqliteBackend) LastBlock(chainID int) (types.Block, error) {
	row := d.db.QueryRow("SELECT * FROM blocks WHERE chain_id = ? ORDER BY height DESC LIMIT 1", chainID)
	return scanBlock(row)
}

func (d *SqliteBackend) BlockCount(chainID int) (int, error) {
	var count int
	err := d.db.QueryRow("SELECT COUNT(*) FROM blocks WHERE chain_id = ?", chainID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count blocks: %w", err)
	}
	return count, nil
}

func (d *SqliteBackend) UpsertBlock(block types.Block, chainID int) error {
	_, err := d.db.Exec("INSERT INTO blocks (height, hash, version, merkleroot, time, mediantime, nonce, bits, difficulty, chainwork, previousblockhash, chain_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) ON CONFLICT (hash, height, chain_id) DO NOTHING", block.Height, block.Hash, block.Version, block.MerkleRoot, block.Time, block.MedianTime, block.Nonce, block.Bits, block.Difficulty, block.ChainWork, block.PreviousBlockHash, chainID)
	return err
}

func scanBlock(row *sql.Row) (types.Block, error) {
	var block types.Block
	err := row.Scan(&block.ID, &block.Height, &block.Hash, &block.Version, &block.MerkleRoot, &block.Time, &block.MedianTime, &block.Nonce, &block.Bits, &block.Difficulty, &block.ChainWork, &block.PreviousBlockHash, &block.ChainID)
	if err != nil {
		return types.Block{}, fmt.Errorf("failed to scan block: %w", err)
	}
	return block, nil
}
