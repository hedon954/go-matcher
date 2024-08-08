package goat_game

import (
	"github.com/hedon954/go-matcher/internal/entry"
	"github.com/hedon954/go-matcher/internal/entry/glicko2"
)

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
