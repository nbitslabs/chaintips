package tracker

import (
	"fmt"

	"github.com/nbitslabs/chaintips/bitcoinrpc"
	"github.com/nbitslabs/chaintips/storage"
)

type Tracker struct {
	db storage.Storage
}

func NewTracker(db storage.Storage) *Tracker {
	return &Tracker{db: db}
}

func (t *Tracker) Run() {
	// Get all chains
	chains, err := t.db.GetChains()
	if err != nil {
		fmt.Println("Failed to get chains:", err)
	}

	// For each chain, get enabled endpoints
	for _, chain := range chains {
		endpoints, err := t.db.GetEnabledEndpoints(chain.ID)
		if err != nil {
			fmt.Println("Failed to get enabled endpoints:", err)
		}

		// For each endpoint, check if it's up
		for _, endpoint := range endpoints {
			fmt.Printf("Checking endpoint %+v\n", endpoint)

			rpcclient := bitcoinrpc.NewRpcClient(endpoint.Protocol, endpoint.Host, endpoint.Port, endpoint.Username, endpoint.Password)
			tips, err := rpcclient.GetChainTips()
			if err != nil {
				fmt.Println("Failed to get chaintips:", err)
				break
			}
			for _, tip := range tips {
				fmt.Printf("Chaintip: %+v\n", tip)
			}
		}
	}
}
