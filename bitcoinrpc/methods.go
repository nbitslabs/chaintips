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
