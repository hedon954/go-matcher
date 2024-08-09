package entry

import (
	"sync"

	"github.com/hedon954/go-matcher/internal/constant"
	"github.com/hedon954/go-matcher/internal/pto"
	"github.com/hedon954/go-matcher/pkg/rand"
)

type Room interface {
	Base() *RoomBase
	ID() int64
	NeedAI() bool
	GetMatchInfo() *pto.MatchInfo
}

type RoomBase struct {
	sync.RWMutex
	id        int64
	teams     []Team
	TeamLimit int

	GameMode       constant.GameMode
	MatchStrategy  constant.MatchStrategy
	ModeVersion    int64
	FinishMatchSec int64

	GameServerInfo pto.GameServerInfo
}

func NewRoomBase(id int64, teamLimit int, t Team) *RoomBase {
	r := &RoomBase{
		id:            id,
		TeamLimit:     teamLimit,
		teams:         make([]Team, 0),
		GameMode:      t.Base().GameMode,
		MatchStrategy: t.Base().MatchStrategy,
		ModeVersion:   t.Base().ModeVersion,
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

func (r *RoomBase) NeedAI() bool {
	return false
}

func (r *RoomBase) ShuffleTeamOrder() {
	ids := rand.PermFrom1(len(r.teams))
	for i := 0; i < len(r.teams); i++ {
		team := r.teams[i]
		team.Base().TeamID = ids[i]
	}
}
