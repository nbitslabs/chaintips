package types

type Block struct {
	ID                int    `json:"id"`
	Height            int64  `json:"height"`
	Hash              string `json:"hash"`
	Version           string `json:"version"`
	MerkleRoot        string `json:"merkle_root"`
	Time              int64  `json:"time"`
	MedianTime        int64  `json:"median_time"`
	Nonce             int64  `json:"nonce"`
	Bits              string `json:"bits"`
	Difficulty        string `json:"difficulty"`
	ChainWork         string `json:"chain_work"`
	PreviousBlockHash string `json:"previous_block_hash"`
	ChainID           int    `json:"chain_id"`
}
