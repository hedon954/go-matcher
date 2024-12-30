package test_game

import (
	"encoding/gob"

	"github.com/hedon954/go-matcher/internal/entry"
)

func init() {
	gob.Register(&Group{})
}

type Group struct {
	*entry.GroupBase
}

func CreateGroup(base *entry.GroupBase) entry.Group {
	return &Group{
		GroupBase: base,
	}
}

func (g *Group) Json() string {
	return entry.Json(g)
}

func (g *Group) Encode() ([]byte, error) {
	return entry.Encode(g)
}

func (g *Group) Decode(data []byte) error {
	return entry.Decode(data, g)
}
