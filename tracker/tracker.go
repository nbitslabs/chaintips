package tracker

import (
	"log"

	"github.com/nbitslabs/chaintips/bitcoinrpc"
	"github.com/nbitslabs/chaintips/storage"
	"github.com/nbitslabs/chaintips/types"
)

type Tracker struct {
	db storage.Storage
}

func NewTracker(db storage.Storage) *Tracker {
	return &Tracker{db: db}
}

func (t *Tracker) Run() {
	chains, err := t.db.GetChains()
	if err != nil {
		log.Println("Failed to get chains:", err)
		return
	}

	for _, chain := range chains {
		endpoints, err := t.db.GetEnabledEndpoints(chain.ID)
		if err != nil {
			log.Println("Failed to get enabled endpoints for chain", chain.Identifier, ":", err)
			continue
		}

		for _, endpoint := range endpoints {
			log.Printf("Checking endpoint %+v\n", endpoint)

			rpcclient := bitcoinrpc.NewRpcClient(endpoint.Protocol, endpoint.Host, endpoint.Port, endpoint.Username, endpoint.Password)
			tips, err := rpcclient.GetChainTips()
			if err != nil {
				log.Println("Failed to get chaintips from", endpoint.Host, ":", err)
				continue
			}

			for _, tip := range tips {
				dbTip := types.ChainTip{
					ChainID:    chain.ID,
					EndpointID: endpoint.ID,
					Height:     tip.Height,
					Hash:       tip.Hash,
					Branchlen:  tip.Branchlen,
					Status:     tip.Status,
				}
				if err := t.db.UpsertChainTip(dbTip); err != nil {
					log.Printf("Failed to upsert chaintip %+v: %v\n", tip, err)
				} else {
					log.Printf("Upserted chaintip %+v successfully\n", tip)
				}
			}
		}
	}
}
