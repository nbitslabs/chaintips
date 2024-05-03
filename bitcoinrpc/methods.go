package bitcoinrpc

import (
	"encoding/json"
	"fmt"

	"github.com/nbitslabs/chaintips/bitcoinrpc/types"
)

func (rpc *RpcClient) GetChainTips() ([]types.ChainTip, error) {
	result, err := rpc.Do("getchaintips", nil)
	if err != nil {
		return nil, err
	}

	var tips []types.ChainTip
	if err := json.Unmarshal(result, &tips); err != nil {
		return nil, fmt.Errorf("failed to unmarshal getchaintips response: %v", err)
	}

	return tips, nil
}

func (rpc *RpcClient) GetBestBlockHash() (string, error) {
	result, err := rpc.Do("getbestblockhash", nil)
	if err != nil {
		return "", err
	}

	var hash string
	if err := json.Unmarshal(result, &hash); err != nil {
		return "", fmt.Errorf("failed to unmarshal getbestblockhash response: %v", err)
	}

	return hash, nil
}

func (rpc *RpcClient) GetBlockHeader(hash string) (types.BlockHeader, error) {
	result, err := rpc.Do("getblockheader", []interface{}{hash})
	if err != nil {
		return types.BlockHeader{}, err
	}

	var header types.BlockHeader
	if err := json.Unmarshal(result, &header); err != nil {
		return types.BlockHeader{}, fmt.Errorf("failed to unmarshal getblockheader response: %v", err)
	}

	return header, nil
}
