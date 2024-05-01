package types

type Endpoint struct {
	ID       int    `json:"id"`
	ChainID  int    `json:"chain_id"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
	Username string `json:"username"`
	Password string `json:"password"`
	Enabled  bool   `json:"enabled"`
}
