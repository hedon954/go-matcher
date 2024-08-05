package goat_game

import (
	"fmt"

	"github.com/hedon954/go-matcher/internal/constant"
	"github.com/hedon954/go-matcher/internal/entry"
	"github.com/hedon954/go-matcher/internal/entry/goat_game/glicko2"
)

type Room struct {
	*entry.RoomBase
}

func CreateRoom(base *entry.RoomBase, t entry.Team) (entry.Room, error) {
	room := &Room{base}
	switch t.Base().MatchStrategy() {
	case constant.MatchStrategyGlicko2:
		return glicko2.CreateRoom(room, t.(*glicko2.Team)), nil
	default:
		return nil, fmt.Errorf("unsupported match strategy: %d", t.Base().MatchStrategy())
	}
}
