package goat_game

import (
	"encoding/gob"

	"github.com/hedon954/go-matcher/internal/entry"
	"github.com/hedon954/go-matcher/internal/entry/glicko2"
)

func init() {
	gob.Register(&Room{})
}

type Room struct {
	*glicko2.RoomBaseGlicko2
}

func (r *Room) withMatchStrategy(base *entry.RoomBase, mgr *entry.TeamMgr) {
	r.RoomBaseGlicko2 = glicko2.CreateRoomBase(base, mgr)
}

func (r *Room) Encode() ([]byte, error) {
	return entry.Encode(r)
}

func (r *Room) Decode(data []byte) error {
	return entry.Decode(data, r)
}
