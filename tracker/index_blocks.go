package tracker

import (
	"strconv"
	"time"

	"github.com/nbitslabs/chaintips/bitcoinrpc"
	"github.com/nbitslabs/chaintips/types"
	"github.com/rs/zerolog/log"
)

func (t *Tracker) indexBlocks(chain types.Chain) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			log.Info().
				Int("chain_id", chain.ID).
				Str("chain_identifier", chain.Identifier).
				Msg("Indexing blocks")

			endpoints, err := t.db.GetEnabledEndpoints(chain.ID)
			if err != nil {
				log.Error().Err(err).
					Int("chain_id", chain.ID).
					Str("chain_identifier", chain.Identifier).
					Msg("Failed to get enabled endpoints")
				continue
			}

			if len(endpoints) == 0 {
				log.Warn().
					Int("chain_id", chain.ID).
					Msg("No enabled endpoints")
				continue
			}

			endpoint := endpoints[0]

			rpcclient := bitcoinrpc.NewRpcClient(endpoint.Protocol, endpoint.Host, endpoint.Port, endpoint.Username, endpoint.Password)

			blockCount, err := t.db.BlockCount(chain.ID)
			if err != nil {
				log.Error().Err(err).
					Int("chain_id", chain.ID).
					Msg("Failed to get block count")
				continue
			}

			if blockCount == 0 {
				log.Info().
					Int("chain_id", chain.ID).
					Msg("No blocks in database, starting initial sync")

				t.startBackfill(chain, rpcclient)
			} else {
				firstBlock, err := t.db.FirstBlock(chain.ID)
				if err != nil {
					log.Error().Err(err).
						Int("chain_id", chain.ID).
						Msg("Failed to get first block")
					continue
				}

				if firstBlock.Height != 0 {
					log.Info().
						Int("chain_id", chain.ID).
						Str("chain_identifier", chain.Identifier).
						Str("first_block_hash", firstBlock.Hash).
						Int64("first_block_height", firstBlock.Height).
						Msg("Not yet synced to genesis block, continuing backfill")

					backFillLimit := 20
					if firstBlock.Height < 20 {
						backFillLimit = int(firstBlock.Height)
					}

					t.backfillBlocks(chain, rpcclient, firstBlock.PreviousBlockHash, backFillLimit)
				}
			}
		}
	}
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

	t.backfillBlocks(chain, rpcclient, bestBlockHash, 20)
}

func (t *Tracker) backfillBlocks(chain types.Chain, rpcclient *bitcoinrpc.RpcClient, bestBlockHash string, count int) {
	log.Info().
		Str("chain_identifier", chain.Identifier).
		Int("chain_id", chain.ID).
		Str("block_hash", bestBlockHash).
		Int("count", count).
		Msg("Backfilling blocks")

	for i := 0; i < count; i++ {
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
	}
}
