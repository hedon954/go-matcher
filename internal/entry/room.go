package entry

import (
	"sync"

	"github.com/hedon954/go-matcher/internal/constant"
	"github.com/hedon954/go-matcher/internal/pto"
)

type Room interface {
	Base() *RoomBase
	ID() int64
	GetMatchInfo() *pto.MatchInfo
}

type RoomBase struct {
	sync.RWMutex
	id            int64
	teams         []Team
	GameMode      constant.GameMode
	MatchStrategy constant.MatchStrategy
}

func NewRoomBase(id int64, t Team) *RoomBase {
	r := &RoomBase{
		id:            id,
		teams:         make([]Team, 0),
		GameMode:      t.Base().GameMode,
		MatchStrategy: t.Base().MatchStrategy,
	}
	return r
}

func (r *RoomBase) Base() *RoomBase {
	return r
}

func (r *RoomBase) ID() int64 {
	return r.id
}

func (r *RoomBase) GetTeams() []Team {
	return r.teams
}

func (r *RoomBase) AddTeam(t Team) {
	r.teams = append(r.teams, t)
}

func (r *RoomBase) RemoveTeam(id int64) {
	for i, t := range r.teams {
		if t.ID() == id {
			r.teams = append(r.teams[:i], r.teams[i+1:]...)
			break
		}
	}
}

func (r *RoomBase) GetMatchInfo() *pto.MatchInfo {
	// TODO
	return nil
}
