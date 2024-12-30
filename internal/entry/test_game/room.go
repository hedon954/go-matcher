package test_game

import (
	"encoding/gob"

	"github.com/hedon954/go-matcher/internal/entry"
)

func init() {
	gob.Register(&Room{})
}

type Room struct {
	*entry.RoomBase
}

func CreateRoom(base *entry.RoomBase) entry.Room {
	room := &Room{
		RoomBase: base,
	}
	return room
}

func (r *Room) Encode() ([]byte, error) {
	return entry.Encode(r)
}

func (r *Room) Decode(data []byte) error {
	return entry.Decode(data, r)
}
