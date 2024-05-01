package types

type Endpoint struct {
	ID       int    `json:"id"`
	ChainID  int    `json:"chain_id"`
	IP       string `json:"ip"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Enabled  bool   `json:"enabled"`
}
