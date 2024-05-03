package storage

import "github.com/nbitslabs/chaintips/types"

type Storage interface {
	GetChains() ([]types.Chain, error)
	GetEnabledEndpoints(chainID int) ([]types.Endpoint, error)
	GetEnabledEndpoint(id int) (types.Endpoint, error)
	UpsertChainTip(tip types.ChainTip) error
	FirstBlock(chainID int) (types.Block, error)
	LastBlock(chainID int) (types.Block, error)
	BlockCount(chainID int) (int, error)
	UpsertBlock(block types.Block, chainID int) error
	BlockAtHeight(height int, chainID int) (types.Block, error)
}
