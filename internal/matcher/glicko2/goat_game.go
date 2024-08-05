package glicko2

import (
	"fmt"

	"github.com/hedon954/go-matcher/internal/constant"
	"github.com/hedon954/go-matcher/internal/entry"
	"github.com/hedon954/go-matcher/pkg/algorithm/glicko2"
)

func (m *Matcher) resgiterGoatGame() {
	m.AddMode(constant.GameModeGoatGame, &Funcs{
		ArgsFunc:          m.goatGameArgs,
		NewTeamFunc:       m.newGoatGameTeam,
		NewRoomFunc:       m.newGoatGameRoom,
		NewRoomWithAIFunc: m.newGoatGameRoomWithAI,
	})
}

func (m *Matcher) goatGameArgs() *glicko2.QueueArgs {
	return &glicko2.QueueArgs{
		MatchTimeoutSec:              50,
		TeamPlayerLimit:              1,
		RoomTeamLimit:                3,
		NewerWithNewer:               false,
		UnfriendlyTeamMMRVarianceMin: 0,
		MaliciousTeamMMRVarianceMin:  0,
		NormalTeamWaitTimeSec:        0,
		UnfriendlyTeamWaitTimeSec:    0,
		MaliciousTeamWaitTimeSec:     0,
		MatchRanges: []glicko2.MatchRange{
			{
				MaxMatchSec:   0,
				MMRGapPercent: 0,
				CanJoinTeam:   true,
				StarGap:       0,
			},
		},
	}
}

func (m *Matcher) newGoatGameTeam(g glicko2.Group) glicko2.Team {
	t, err := m.teamMgr.CreateTeam(g.(entry.Group))
	if err != nil {
		panic(fmt.Sprintf("create team error: %s", err.Error()))
	}
	return t.(glicko2.Team)
}

func (m *Matcher) newGoatGameRoom(t glicko2.Team) glicko2.Room {
	r, err := m.roomMgr.CreateRoom(t.(entry.Team))
	if err != nil {
		panic(fmt.Sprintf("create room error: %s", err.Error()))
	}
	return r.(glicko2.Room)
}

func (m *Matcher) newGoatGameRoomWithAI(t glicko2.Team) glicko2.Room {
	return m.newGoatGameRoom(t)
}
