package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/nbitslabs/chaintips/types"
)

func (d *SqliteBackend) FirstBlock(chainID int) (types.Block, error) {
	row := d.db.QueryRow("SELECT id, height, hash, version, merkleroot, time, mediantime, nonce, bits, difficulty, chainwork, previousblockhash, chain_id FROM blocks WHERE chain_id = ? ORDER BY height ASC LIMIT 1", chainID)
	return scanBlock(row)
}

func (d *SqliteBackend) LastBlock(chainID int) (types.Block, error) {
	row := d.db.QueryRow("SELECT id, height, hash, version, merkleroot, time, mediantime, nonce, bits, difficulty, chainwork, previousblockhash, chain_id FROM blocks WHERE chain_id = ? ORDER BY height DESC LIMIT 1", chainID)
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

func (d *SqliteBackend) BlockAtHeight(height int64, chainID int) (types.Block, error) {
	row := d.db.QueryRow("SELECT id, height, hash, version, merkleroot, time, mediantime, nonce, bits, difficulty, chainwork, previousblockhash, chain_id FROM blocks WHERE height = ? AND chain_id = ?", height, chainID)
	return scanBlock(row)
}

func (d *SqliteBackend) UpsertBlock(block types.Block, chainID int) error {
	_, err := d.db.Exec("INSERT INTO blocks (height, hash, version, merkleroot, time, mediantime, nonce, bits, difficulty, chainwork, previousblockhash, chain_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) ON CONFLICT (hash, height, chain_id) DO NOTHING", block.Height, block.Hash, block.Version, block.MerkleRoot, block.Time, block.MedianTime, block.Nonce, block.Bits, block.Difficulty, block.ChainWork, block.PreviousBlockHash, chainID)
	return err
}

func (d *SqliteBackend) GetNotableBlocks(chainID int) ([]types.Block, error) {
	var blocks []types.Block
	firstBlock, err := d.FirstBlock(chainID)
	if err != nil {
		return nil, fmt.Errorf("failed to get first block: %w", err)
	}

	blocks = append(blocks, firstBlock)

	query := `WITH RECURSIVE
					    height_range AS (
					        SELECT
					            chain_id,
					            endpoint_id,
					            height - branchlen AS start_height,
					            height AS end_height,
					            hash AS tip_hash
					        FROM chaintips WHERE chain_id = ?
					        
					        UNION ALL
					        
					        SELECT
					            hr.chain_id,
					            hr.endpoint_id,
					            hr.start_height + 1,
					            hr.end_height,
					            hr.tip_hash
					        FROM height_range hr
					        WHERE hr.start_height + 1 <= hr.end_height
					    ),
					    fork_blocks AS (
					        SELECT
					            b.*,
					            hr.tip_hash,
					            CASE
					                WHEN b.hash = hr.tip_hash THEN 'true'
					                ELSE 'false'
					            END AS fork
					        FROM blocks b
					        JOIN height_range hr ON b.height BETWEEN hr.start_height AND hr.end_height
					    )
					SELECT
							fb.id,
					    fb.height,
					    fb.hash,
							fb.version,
							fb.merkleroot,
							fb.time,
							fb.mediantime,
							fb.nonce,
							fb.bits,
							fb.difficulty,
							fb.chainwork,
							fb.previousblockhash,
							fb.chain_id,
					    MAX(fb.fork) AS fork -- Aggregate to ensure we consider any true values as definitive for fork status
					FROM fork_blocks fb
					WHERE chain_id = ?
					GROUP BY fb.height, fb.hash
					ORDER BY fb.height;`

	rows, err := d.db.Query(query, chainID, chainID)
	if err != nil {
		return nil, fmt.Errorf("failed to query notable blocks: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var block types.Block
		if err := rows.Scan(&block.ID, &block.Height, &block.Hash, &block.Version, &block.MerkleRoot, &block.Time, &block.MedianTime, &block.Nonce, &block.Bits, &block.Difficulty, &block.ChainWork, &block.PreviousBlockHash, &block.ChainID, &block.Fork); err != nil {
			return nil, fmt.Errorf("failed to scan block: %w", err)
		}
		blocks = append(blocks, block)
	}

	lastBlock, err := d.LastBlock(chainID)
	if err != nil {
		return nil, fmt.Errorf("failed to get last block: %w", err)
	}

	blocks = append(blocks, lastBlock)

	return blocks, nil
}

func scanBlock(row *sql.Row) (types.Block, error) {
	var block types.Block
	err := row.Scan(&block.ID, &block.Height, &block.Hash, &block.Version, &block.MerkleRoot, &block.Time, &block.MedianTime, &block.Nonce, &block.Bits, &block.Difficulty, &block.ChainWork, &block.PreviousBlockHash, &block.ChainID)
	if err != nil {
		return types.Block{}, fmt.Errorf("failed to scan block: %w", err)
	}
	return block, nil
}
