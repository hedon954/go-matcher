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

func CreateRoom(base *entry.RoomBase) entry.Room {
	room := &Room{}

	room.withMatchStrategy(base)
	return room
}

func (r *Room) withMatchStrategy(base *entry.RoomBase) {
	r.RoomBaseGlicko2 = glicko2.CreateRoomBase(base)
}

func (r *Room) Encode() ([]byte, error) {
	return entry.Encode(r)
}

func (r *Room) Decode(data []byte) error {
	return entry.Decode(data, r)
}
