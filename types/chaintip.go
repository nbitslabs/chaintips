package types

type ChainTip struct {
	ID         int    `json:"-"`
	ChainID    int    `json:"chain_id"`
	EndpointID int    `json:"endpoint_id"`
	Height     int64  `json:"height"`
	Hash       string `json:"hash"`
	Branchlen  int    `json:"branchlen"`
	Status     string `json:"status"`
	InsertedAt string `json:"inserted_at"`
}
