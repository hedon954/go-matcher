package glicko2

import (
	"github.com/hedon954/go-matcher/internal/entry"
	"github.com/hedon954/go-matcher/pkg/algorithm/glicko2"
)

type TeamBaseGlicko2 struct {
	*entry.TeamBase
	groupMgr GroupMgr `msgpack:"-"`
}

type GroupMgr interface {
	Get(id int64) entry.Group
}

func CreateTeamBase(base *entry.TeamBase, mgr GroupMgr) *TeamBaseGlicko2 {
	t := &TeamBaseGlicko2{
		TeamBase: base,
		groupMgr: mgr,
	}
	return t
}

func (t *TeamBaseGlicko2) GetGroups() []glicko2.Group {
	t.RLock()
	defer t.RUnlock()
	groups := t.Base().GetGroups()
	res := make([]glicko2.Group, len(groups))
	for i := 0; i < len(groups); i++ {
		g := t.groupMgr.Get(groups[i])
		res[i] = g.(glicko2.Group)
	}
	return res
}

func (t *TeamBaseGlicko2) AddGroup(g glicko2.Group) {
	t.Lock()
	defer t.Unlock()
	t.Base().AddGroup(g.(entry.Group))
}

func (t *TeamBaseGlicko2) PlayerCount() int {
	count := 0
	for _, g := range t.GetGroups() {
		count += g.PlayerCount()
	}
	return count
}

func (t *TeamBaseGlicko2) GetMMR() float64 {
	groups := t.GetGroups()
	if len(groups) == 0 {
		return 0.0
	}
	total := 0.0
	for _, g := range t.GetGroups() {
		total += g.GetMMR()
	}
	return total / float64(len(t.GetGroups()))
}

func (t *TeamBaseGlicko2) GetStar() int {
	groups := t.GetGroups()
	if len(groups) == 0 {
		return 0
	}
	total := 0
	for _, g := range t.GetGroups() {
		total += g.GetStar()
	}
	return total / len(t.GetGroups())
}

func (t *TeamBaseGlicko2) GetStartMatchTimeSec() int64 {
	return t.GetGroups()[0].GetStartMatchTimeSec()
}

func (t *TeamBaseGlicko2) GetFinishMatchTimeSec() int64 {
	return t.GetGroups()[0].GetFinishMatchTimeSec()
}

func (t *TeamBaseGlicko2) SetFinishMatchTimeSec(unix int64) {
	for _, g := range t.GetGroups() {
		g.SetFinishMatchTimeSec(unix)
	}
}

func (t *TeamBaseGlicko2) IsAi() bool {
	for _, g := range t.GetGroups() {
		for _, p := range g.GetPlayers() {
			if p.IsAi() {
				return true
			}
		}
	}
	return false
}

func (t *TeamBaseGlicko2) CanFillAi() bool {
	return false
}

func (t *TeamBaseGlicko2) IsFull(teamPlayerLimit int) bool {
	return t.PlayerCount() >= teamPlayerLimit
}

func (t *TeamBaseGlicko2) IsNewer() bool {
	for _, g := range t.GetGroups() {
		if g.IsNewer() {
			return true
		}
	}
	return false
}

func (t *TeamBaseGlicko2) SetGroupMgr(mgr GroupMgr) {
	t.groupMgr = mgr
}
