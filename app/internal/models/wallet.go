package models

type Wallet struct {
	ID     string         `json:"id"`
	Stocks map[string]int `json:"stocks"`
}