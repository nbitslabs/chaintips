package storage

import "github.com/nbitslabs/chaintips/types"

type Storage interface {
	GetChains() ([]types.Chain, error)
	GetEnabledEndpoints(chainID int) ([]types.Endpoint, error)
	UpsertChainTip(tip types.ChainTip) error
}
