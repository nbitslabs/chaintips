package tracker

import (
	"time"

	"github.com/nbitslabs/chaintips/bitcoinrpc"
	"github.com/nbitslabs/chaintips/types"
	"github.com/rs/zerolog/log"
)

func (t *Tracker) trackTips(chain types.Chain) {
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			endpoints, err := t.db.GetEnabledEndpoints(chain.ID)
			if err != nil {
				log.Error().Err(err).
					Int("chain_id", chain.ID).
					Str("chain_identifier", chain.Identifier).
					Msg("Failed to get enabled endpoints")
				continue
			}

			firstBlock, err := t.db.FirstBlock(chain.ID)
			if err != nil {
				log.Error().Err(err).
					Msg("Failed to get first block")
				continue
			}

			for _, endpoint := range endpoints {
				log.Info().
					Str("endpoint", endpoint.Host).
					Msg("Checking endpoint")

				rpcclient := bitcoinrpc.NewRpcClient(endpoint.Protocol, endpoint.Host, endpoint.Port, endpoint.Username, endpoint.Password)
				tips, err := rpcclient.GetChainTips()
				if err != nil {
					log.Error().Err(err).
						Str("endpoint", endpoint.Host).
						Msg("Failed to get chaintips")
					continue
				}

				for _, tip := range tips {
					// We don't want to track the current chaintip, only forks
					if tip.Status == "active" {
						continue
					}

					// Don't track forks from before our synced progress
					// This can cause the backfilling task to create gaps by
					// unexpectedly inserting block heights that are disjointed
					// from the current chain
					if tip.Height <= firstBlock.Height {
						log.Debug().
							Int("chain_id", chain.ID).
							Str("chain_identifier", chain.Identifier).
							Int64("height", tip.Height).
							Int64("first_block_height", firstBlock.Height).
							Msg("Skipping chaintip, height is less than first block")
						continue
					}

					dbTip := types.ChainTip{
						ChainID:    chain.ID,
						EndpointID: endpoint.ID,
						Height:     tip.Height,
						Hash:       tip.Hash,
						Branchlen:  tip.Branchlen,
						Status:     tip.Status,
					}

					if err := t.db.UpsertChainTip(dbTip); err != nil {
						log.Error().Err(err).
							Msg("Failed to upsert chaintip")
					} else {
						log.Info().
							Str("chaintip", dbTip.Hash).
							Int64("height", dbTip.Height).
							Msg("Upserted chaintip successfully")
					}
				}
			}
		}
	}
	syncWg.Done()
}
