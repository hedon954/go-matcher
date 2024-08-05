package pto

type Glicko2Info struct {
	MMR           float64 `json:"mmr"`
	Star          int     `json:"star"`
	StartMatchSec int64   `json:"start_match_sec"`
	Rank          int     `json:"rank"`
}
