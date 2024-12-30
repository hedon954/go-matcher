package test_game

import (
	"encoding/gob"

	"github.com/hedon954/go-matcher/internal/entry"
)

func init() {
	gob.Register(&Player{})
}

type Player struct {
	*entry.PlayerBase
}

func CreatePlayer(base *entry.PlayerBase) entry.Player {
	return &Player{
		PlayerBase: base,
	}
}

func (p *Player) Encode() ([]byte, error) {
	return entry.Encode(p)
}

func (p *Player) Decode(data []byte) error {
	return entry.Decode(data, p)
}
