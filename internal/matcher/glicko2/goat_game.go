package glicko2

import (
	"fmt"

	"github.com/hedon954/go-matcher/internal/constant"
	"github.com/hedon954/go-matcher/internal/entry"
	"github.com/hedon954/go-matcher/pkg/algorithm/glicko2"
)

func (m *Matcher) registerGoatGame() {
	m.AddMode(constant.GameModeGoatGame, &Funcs{
		ArgsFunc:          m.goatGameArgs,
		NewTeamFunc:       m.newGoatGameTeam,
		NewRoomFunc:       m.newGoatGameRoom,
		NewRoomWithAIFunc: m.newGoatGameRoomWithAI,
	})
}

func (m *Matcher) goatGameArgs() *glicko2.QueueArgs {
	return m.configer.Glicko2.GetQueueArgs(constant.GameModeGoatGame)
}

func (m *Matcher) newGoatGameTeam(g glicko2.Group) glicko2.Team {
	t, err := m.teamMgr.CreateTeam(g.(entry.Group))
	if err != nil {
		panic(fmt.Sprintf("create team error: %s", err.Error()))
	}
	return t.(glicko2.Team)
}

func (m *Matcher) newGoatGameRoom(t glicko2.Team) glicko2.Room {
	teamLimit := m.configer.Glicko2.GetQueueArgs(constant.GameModeGoatGame).RoomTeamLimit
	r, err := m.roomMgr.CreateRoom(teamLimit, t.(entry.Team))
	if err != nil {
		panic(fmt.Sprintf("create room error: %s", err.Error()))
	}
	return r.(glicko2.Room)
}

func (m *Matcher) newGoatGameRoomWithAI(t glicko2.Team) glicko2.Room {
	return m.newGoatGameRoom(t)
}
