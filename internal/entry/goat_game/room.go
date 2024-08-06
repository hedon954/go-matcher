package goat_game

import (
	"fmt"

	"github.com/hedon954/go-matcher/internal/constant"
	"github.com/hedon954/go-matcher/internal/entry"
	"github.com/hedon954/go-matcher/internal/entry/glicko2"
)

type Room struct {
	*glicko2.RoomBaseGlicko2
}

func CreateRoom(base *entry.RoomBase) (entry.Room, error) {
	room := &Room{}

	if err := room.withMatchStrategy(base); err != nil {
		return nil, err
	}
	return room, nil
}

func (r *Room) withMatchStrategy(base *entry.RoomBase) error {
	switch base.MatchStrategy {
	case constant.MatchStrategyGlicko2:
		r.RoomBaseGlicko2 = glicko2.CreateRoomBase(base)
	default:
		return fmt.Errorf("unsupported match strategy: %d", base.MatchStrategy)
	}
	return nil
}
