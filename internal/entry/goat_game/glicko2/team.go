package glicko2

import (
	"encoding/json"
	"log/slog"

	"github.com/hedon954/go-matcher/internal/entry"
	"github.com/hedon954/go-matcher/pkg/algorithm/glicko2"
	"github.com/spf13/cast"
)

type Team struct {
	entry.Team
}

func CreateTeam(team entry.Team, group *Group) entry.Team {
	t := &Team{
		Team: team,
	}
	bs, _ := json.Marshal(group)
	slog.Info("create team", slog.Any("group", string(bs)))
	return t
}

func (t *Team) GetGroups() []glicko2.Group {
	t.Base().RLock()
	defer t.Base().RUnlock()
	groups := t.Base().GetGroups()
	res := make([]glicko2.Group, len(groups))
	for i := 0; i < len(groups); i++ {
		res[i] = groups[i].(glicko2.Group)
	}
	return res
}

func (t *Team) AddGroup(g glicko2.Group) {
	t.Base().Lock()
	defer t.Base().Unlock()
	t.Base().AddGroup(g.(entry.Group))
}

func (t *Team) RemoveGroup(groupId string) {
	t.Base().Lock()
	defer t.Base().Unlock()
	t.Base().RemoveGroup(cast.ToInt64(groupId))
}

func (t *Team) PlayerCount() int {
	count := 0
	for _, g := range t.GetGroups() {
		count += g.PlayerCount()
	}
	return count
}

func (t *Team) GetMMR() float64 {
	if t.PlayerCount() == 0 {
		return 0.0
	}
	total := 0.0
	for _, g := range t.GetGroups() {
		total += g.GetMMR()
	}
	return total / float64(len(t.GetGroups()))
}

func (t *Team) GetStar() int {
	if t.PlayerCount() == 0 {
		return 0
	}
	total := 0
	for _, g := range t.GetGroups() {
		total += g.GetStar()
	}
	return total / len(t.GetGroups())
}

func (t *Team) GetStartMatchTimeSec() int64 {
	return t.GetGroups()[0].GetStartMatchTimeSec()
}

func (t *Team) GetFinishMatchTimeSec() int64 {
	return t.GetGroups()[0].GetFinishMatchTimeSec()
}

func (t *Team) SetFinishMatchTimeSec(unix int64) {
	for _, g := range t.GetGroups() {
		g.SetFinishMatchTimeSec(unix)
	}
}

func (t *Team) IsAi() bool {
	for _, g := range t.GetGroups() {
		for _, p := range g.GetPlayers() {
			if p.IsAi() {
				return true
			}
		}
	}
	return false
}

func (t *Team) CanFillAi() bool {
	return false
}

func (t *Team) IsFull(teamPlayerLimit int) bool {
	return t.PlayerCount() >= teamPlayerLimit
}

func (t *Team) IsNewer() bool {
	for _, g := range t.GetGroups() {
		if g.IsNewer() {
			return true
		}
	}
	return false
}
