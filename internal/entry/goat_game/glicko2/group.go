package goatgameglicko2

import (
	"fmt"
	"log/slog"

	"github.com/hedon954/go-matcher/internal/entry"
	"github.com/hedon954/go-matcher/pkg/algorithm/glicko2"

	"github.com/spf13/cast"
)

type Group struct {
	entry.Group
}

func (g *Group) GetID() string {
	return cast.ToString(g.ID())
}

func (g *Group) GetPlayers() []glicko2.Player {
	players := g.Base().GetPlayers()
	res := make([]glicko2.Player, len(players))
	for i := 0; i < len(players); i++ {
		res[i] = players[i].(*Player)
	}
	return res
}

func (g *Group) PlayerCount() int {
	return len(g.Base().GetPlayers())
}

func (g *Group) GetMMR() float64 {
	total := 0.0
	players := g.GetPlayers()
	if len(players) == 0 {
		return 0.0
	}
	for _, p := range players {
		total += p.GetMMR()
	}
	return total / float64(len(players))
}

func (g *Group) GetStar() int {
	total := 0
	players := g.GetPlayers()
	if len(players) == 0 {
		return 0.0
	}
	for _, p := range players {
		total += p.GetStar()
	}
	return total / len(players)
}

func (g *Group) GetState() glicko2.GroupState {
	g.Base().RLock()
	defer g.Base().RUnlock()
	switch g.Base().GetState() {
	case entry.GroupStateDissolved, entry.GroupStateInvite:
		return glicko2.GroupStateUnready
	case entry.GroupStateMatch:
		return glicko2.GroupStateQueuing
	case entry.GroupStateGame:
		return glicko2.GroupStateMatched
	}
	panic(fmt.Sprintf("unreachable, state: %d", g.Base().GetState()))
}

func (g *Group) SetState(state glicko2.GroupState) {
	g.Base().Lock()
	defer g.Base().Unlock()
	switch state {
	case glicko2.GroupStateUnready:
		g.Base().SetState(entry.GroupStateInvite)
	case glicko2.GroupStateQueuing:
		g.Base().SetState(entry.GroupStateMatch)
	case glicko2.GroupStateMatched:
		g.Base().SetState(entry.GroupStateGame)
	}
}

func (g *Group) GetStartMatchTimeSec() int64 {
	return g.GetPlayers()[0].GetStartMatchTimeSec()
}

func (g *Group) SetStartMatchTimeSec(t int64) {
	for _, p := range g.GetPlayers() {
		p.SetStartMatchTimeSec(t)
	}
}

func (g *Group) GetFinishMatchTimeSec() int64 {
	return g.GetPlayers()[0].GetFinishMatchTimeSec()
}

func (g *Group) SetFinishMatchTimeSec(t int64) {
	for _, p := range g.GetPlayers() {
		p.SetFinishMatchTimeSec(t)
	}
}

func (g *Group) Type() glicko2.GroupType {
	return glicko2.GroupTypeNormalTeam
}

func (g *Group) CanFillAi() bool {
	return false
}

func (g *Group) ForceCancelMatch(reason string, waitSec int64) {
	slog.Info("force cancel match",
		slog.Any("group", g),
		slog.String("reason", reason),
		slog.Int64("wait_sec", waitSec),
	)
}

func (g *Group) IsNewer() bool {
	return false
}
