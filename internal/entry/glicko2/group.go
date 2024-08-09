package glicko2

import (
	"fmt"
	"log/slog"

	"github.com/hedon954/go-matcher/internal/constant"
	"github.com/hedon954/go-matcher/internal/entry"
	"github.com/hedon954/go-matcher/pkg/algorithm/glicko2"

	"github.com/spf13/cast"
)

type GroupBaseGlicko2 struct {
	*entry.GroupBase
}

func NewGroup(base *entry.GroupBase) *GroupBaseGlicko2 {
	base.SupportMatchStrategies = append(base.SupportMatchStrategies, constant.MatchStrategyGlicko2)

	g := &GroupBaseGlicko2{
		GroupBase: base,
	}
	return g
}

func (g *GroupBaseGlicko2) GetID() string {
	return cast.ToString(g.ID())
}

func (g *GroupBaseGlicko2) QueueKey() string {
	return fmt.Sprintf("%d-%d", g.GameMode, g.ModeVersion)
}

func (g *GroupBaseGlicko2) GetPlayers() []glicko2.Player {
	players := g.Base().GetPlayers()
	res := make([]glicko2.Player, len(players))
	for i := 0; i < len(players); i++ {
		res[i] = players[i].(glicko2.Player)
	}
	return res
}

func (g *GroupBaseGlicko2) PlayerCount() int {
	return len(g.GetPlayers())
}

func (g *GroupBaseGlicko2) GetMMR() float64 {
	players := g.GetPlayers()
	if len(players) == 0 {
		return 0.0
	}
	total := 0.0
	for _, p := range players {
		total += p.GetMMR()
	}
	return total / float64(len(players))
}

func (g *GroupBaseGlicko2) GetStar() int {
	players := g.GetPlayers()
	if len(players) == 0 {
		return 0.0
	}
	total := 0
	for _, p := range players {
		total += p.GetStar()
	}
	return total / len(players)
}

func (g *GroupBaseGlicko2) GetState() glicko2.GroupState {
	g.Lock()
	defer g.Unlock()
	switch g.Base().GetState() {
	case entry.GroupStateDissolved, entry.GroupStateInvite:
		return glicko2.GroupStateUnready
	case entry.GroupStateMatch:
		return glicko2.GroupStateQueuing
	case entry.GroupStateGame:
		return glicko2.GroupStateMatched
	}
	panic(fmt.Sprintf("unreachable, state: %d", g.GetState()))
}

func (g *GroupBaseGlicko2) SetState(state glicko2.GroupState) {
	g.Lock()
	defer g.Unlock()
	switch state {
	case glicko2.GroupStateUnready:
		g.Base().SetState(entry.GroupStateInvite)
	case glicko2.GroupStateQueuing:
		g.Base().SetState(entry.GroupStateMatch)
	case glicko2.GroupStateMatched:
		g.Base().SetState(entry.GroupStateGame)
	}
}

func (g *GroupBaseGlicko2) GetStartMatchTimeSec() int64 {
	return g.GetPlayers()[0].GetStartMatchTimeSec()
}

func (g *GroupBaseGlicko2) SetStartMatchTimeSec(t int64) {
	for _, p := range g.GetPlayers() {
		p.SetStartMatchTimeSec(t)
	}
}

func (g *GroupBaseGlicko2) GetFinishMatchTimeSec() int64 {
	return g.GetPlayers()[0].GetFinishMatchTimeSec()
}

func (g *GroupBaseGlicko2) SetFinishMatchTimeSec(t int64) {
	for _, p := range g.GetPlayers() {
		p.SetFinishMatchTimeSec(t)
	}
}

func (g *GroupBaseGlicko2) Type() glicko2.GroupType {
	return glicko2.GroupTypeNormalTeam
}

func (g *GroupBaseGlicko2) CanFillAi() bool {
	return false
}

func (g *GroupBaseGlicko2) ForceCancelMatch(reason string, waitSec int64) {
	slog.Info("force cancel match",
		slog.Any("group", g),
		slog.String("reason", reason),
		slog.Int64("wait_sec", waitSec),
	)
	g.Lock()
	defer g.Unlock()
	g.Base().SetState(entry.GroupStateInvite)
	for _, p := range g.Base().GetPlayers() {
		p.Base().Lock()
		p.Base().SetOnlineState(entry.PlayerOnlineStateInGroup)
		p.Base().Unlock()
	}
	// TODO: push cancel match to users
}

func (g *GroupBaseGlicko2) IsNewer() bool {
	return false
}
