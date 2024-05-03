package types

type Block struct {
	ID      int    `json:"id"`
	Height  int64  `json:"height"`
	Hash    string `json:"hash"`
	ChainID int    `json:"chain_id"`
}
