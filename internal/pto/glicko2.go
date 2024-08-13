package pto

type Glicko2Info struct {
	MMR  float64 `json:"mmr"`
	Star int64   `json:"star"`
	Rank int64   `json:"rank"`
}
