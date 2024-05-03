package tracker

import (
	"strconv"
	"time"

	"github.com/nbitslabs/chaintips/bitcoinrpc"
	"github.com/nbitslabs/chaintips/types"
	"github.com/rs/zerolog/log"
)

const ZeroHash = "0000000000000000000000000000000000000000000000000000000000000000"

func (t *Tracker) indexBlocks(chain types.Chain) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	endpoints, err := t.db.GetEnabledEndpoints(chain.ID)
	if err != nil {
		log.Error().Err(err).
			Int("chain_id", chain.ID).
			Str("chain_identifier", chain.Identifier).
			Msg("Failed to get enabled endpoints")
		return
	}

	if len(endpoints) == 0 {
		log.Warn().
			Int("chain_id", chain.ID).
			Msg("No enabled endpoints")
		return
	}

	endpoint := endpoints[0]

	rpcclient := bitcoinrpc.NewRpcClient(endpoint.Protocol, endpoint.Host, endpoint.Port, endpoint.Username, endpoint.Password)

	blockCount, err := t.db.BlockCount(chain.ID)
	if err != nil {
		log.Error().Err(err).
			Int("chain_id", chain.ID).
			Msg("Failed to get block count")
		return
	}

	if blockCount == 0 {
		log.Info().
			Int("chain_id", chain.ID).
			Msg("No blocks in database, starting initial sync")

		go t.startBackfill(chain, rpcclient)
	} else {
		firstBlock, err := t.db.FirstBlock(chain.ID)
		if err != nil {
			log.Error().Err(err).
				Int("chain_id", chain.ID).
				Msg("Failed to get first block")
			return
		}

		if firstBlock.Height != 0 {
			log.Info().
				Int("chain_id", chain.ID).
				Str("chain_identifier", chain.Identifier).
				Str("first_block_hash", firstBlock.Hash).
				Int64("first_block_height", firstBlock.Height).
				Msg("Not yet synced to genesis block, continuing backfill")

			go t.backfillBlocks(chain, rpcclient, firstBlock.PreviousBlockHash, ZeroHash)
		}
	}

	for {
		select {
		case <-ticker.C:
			log.Info().
				Int("chain_id", chain.ID).
				Str("chain_identifier", chain.Identifier).
				Msg("Indexing blocks")

			if blockCount == 0 {
				blockCount, err = t.db.BlockCount(chain.ID)
				if err != nil {
					log.Error().Err(err).
						Int("chain_id", chain.ID).
						Msg("Failed to get block count")
					continue
				}
			}

			if blockCount == 0 {
				continue
			}

			lastBlock, err := t.db.LastBlock(chain.ID)
			if err != nil {
				log.Error().Err(err).
					Int("chain_id", chain.ID).
					Msg("Failed to get last block")
				continue
			}

			bestBlockHash, err := rpcclient.GetBestBlockHash()
			if err != nil {
				log.Error().Err(err).
					Int("chain_id", chain.ID).
					Msg("Failed to get best block hash")
				continue
			}

			if lastBlock.Hash != bestBlockHash {
				log.Info().
					Str("chain_identifier", chain.Identifier).
					Int("chain_id", chain.ID).
					Str("last_block_hash", lastBlock.Hash).
					Str("best_block_hash", bestBlockHash).
					Msg("Populating block segment")
				t.backfillBlocks(chain, rpcclient, bestBlockHash, lastBlock.Hash)
			}
		}
	}

	syncWg.Done()
}

func (t *Tracker) startBackfill(chain types.Chain, rpcclient *bitcoinrpc.RpcClient) {
	bestBlockHash, err := rpcclient.GetBestBlockHash()
	if err != nil {
		log.Error().Err(err).
			Str("chain_identifier", chain.Identifier).
			Int("chain_id", chain.ID).
			Msg("Failed to get best block hash")
		return
	}

	t.backfillBlocks(chain, rpcclient, bestBlockHash, ZeroHash)
}

func (t *Tracker) backfillBlocks(chain types.Chain, rpcclient *bitcoinrpc.RpcClient, bestBlockHash string, target string) {
	log.Info().
		Str("chain_identifier", chain.Identifier).
		Int("chain_id", chain.ID).
		Str("block_hash", bestBlockHash).
		Str("target_block", target).
		Msg("Backfilling blocks")

	lastHeight := int64(-1)

	for (bestBlockHash != target) || (target == ZeroHash && lastHeight == 0) {
		block, err := rpcclient.GetBlockHeader(bestBlockHash)
		if err != nil {
			log.Error().Err(err).
				Str("chain_identifier", chain.Identifier).
				Int("chain_id", chain.ID).
				Str("block_hash", bestBlockHash).
				Msg("Failed to get block")
			return
		}

		if err := t.db.UpsertBlock(types.Block{
			ChainID:           chain.ID,
			Height:            block.Height,
			Hash:              block.Hash,
			Version:           block.VersionHex,
			MerkleRoot:        block.MerkleRoot,
			Time:              block.Time,
			MedianTime:        block.MedianTime,
			Nonce:             int64(block.Nonce),
			Bits:              block.Bits,
			Difficulty:        strconv.FormatFloat(block.Difficulty, 'f', -1, 64),
			ChainWork:         block.ChainWork,
			PreviousBlockHash: block.PreviousBlockHash,
		}, chain.ID); err != nil {
			log.Error().Err(err).
				Str("chain_identifier", chain.Identifier).
				Int("chain_id", chain.ID).
				Str("block_hash", block.Hash).
				Msg("Failed to upsert block")
			return
		}

		log.Debug().
			Int("chain_id", chain.ID).
			Str("chain_identifier", chain.Identifier).
			Str("block_hash", block.Hash).
			Int64("block_height", block.Height).
			Msg("Inserted block")

		bestBlockHash = block.PreviousBlockHash
		lastHeight = block.Height
	}

	log.Info().
		Str("chain_identifier", chain.Identifier).
		Int("chain_id", chain.ID).
		Str("target_block", target).
		Msg("Backfill complete")
}
