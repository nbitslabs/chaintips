package types

type ChainTip struct {
	Height    int64  `json:"height"`
	Hash      string `json:"hash"`
	Branchlen int    `json:"branchlen"`
	Status    string `json:"status"`
}
