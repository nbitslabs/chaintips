package types

type Chain struct {
	ID         int    `json:"id"`
	Identifier string `json:"identifier"`
	Title      string `json:"title"`
	Icon       string `json:"icon"`
}
