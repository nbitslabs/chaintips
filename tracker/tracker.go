package tracker

import (
	"github.com/nbitslabs/chaintips/bitcoinrpc"
	"github.com/nbitslabs/chaintips/storage"
	"github.com/nbitslabs/chaintips/types"
	"github.com/rs/zerolog/log"
)

var logger = log.With().Str("module", "tracker").Logger()

type Tracker struct {
	db storage.Storage
}

func NewTracker(db storage.Storage) *Tracker {
	return &Tracker{db: db}
}

func (t *Tracker) Run() {
	chains, err := t.db.GetChains()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get chains")
		return
	}

	for _, chain := range chains {
		endpoints, err := t.db.GetEnabledEndpoints(chain.ID)
		if err != nil {
			log.Error().Err(err).
				Int("chain_id", chain.ID).
				Str("chain_identifier", chain.Identifier).
				Msg("Failed to get enabled endpoints")
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
				if tip.Status == "active" {
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
