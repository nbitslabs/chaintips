package tracker

import (
	"time"

	"github.com/nbitslabs/chaintips/bitcoinrpc"
	"github.com/nbitslabs/chaintips/types"
	"github.com/rs/zerolog/log"
)

func (t *Tracker) linkChainTips(chain types.Chain) {
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			unlinkedTips, err := t.db.GetUnlinkedChainTips(chain.ID)
			if err != nil {
				log.Error().Err(err).
					Int("chain_id", chain.ID).
					Str("chain_identifier", chain.Identifier).
					Msg("Failed to get unlinked chaintips")
				continue
			}

			for _, tip := range unlinkedTips {
				log.Info().
					Str("tip_hash", tip.Hash).
					Str("chain_identifier", chain.Identifier).
					Int("chain_id", chain.ID).
					Int64("height", tip.Height).
					Int("branch_length", tip.Branchlen).
					Msg("Linking chaintip")

				endpoint, err := t.db.GetEnabledEndpoint(tip.EndpointID)
				if err != nil {
					log.Error().Err(err).
						Int("chain_id", chain.ID).
						Str("chain_identifier", chain.Identifier).
						Int("endpoint_id", tip.EndpointID).
						Msg("Failed to get endpoint")
					continue
				}

				commonAncestor, err := t.db.BlockAtHeight(tip.Height-int64(tip.Branchlen), chain.ID)
				if err != nil {
					log.Error().Err(err).
						Str("chain_identifier", chain.Identifier).
						Int("chain_id", chain.ID).
						Int64("height", tip.Height).
						Int("branch_length", tip.Branchlen).
						Msg("Failed to get common ancestor")
					continue
				}

				log.Info().
					Str("chain_identifier", chain.Identifier).
					Int("chain_id", chain.ID).
					Str("tip_hash", tip.Hash).
					Int64("height", tip.Height).
					Int("branch_length", tip.Branchlen).
					Str("common_ancestor_hash", commonAncestor.Hash).
					Int64("common_ancestor_height", commonAncestor.Height).
					Msg("Linking chaintip")

				rpcclient := bitcoinrpc.NewRpcClient(endpoint.Protocol, endpoint.Host, endpoint.Port, endpoint.Username, endpoint.Password)

				t.backfillBlocks(chain, rpcclient, tip.Hash, commonAncestor.Hash)
			}
		}
	}
}
