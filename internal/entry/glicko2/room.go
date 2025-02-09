package glicko2

import (
	"math"

	"github.com/hedon954/go-matcher/internal/entry"
	"github.com/hedon954/go-matcher/pkg/algorithm/glicko2"
)

type RoomBaseGlicko2 struct {
	*entry.RoomBase
	glicko2Teams map[int64]glicko2.Team `msgpack:"-"`
	teamMgr      TeamMgr                `msgpack:"-"`
}

type TeamMgr interface {
	Get(id int64) entry.Team
}

func CreateRoomBase(base *entry.RoomBase, mgr TeamMgr) *RoomBaseGlicko2 {
	r := &RoomBaseGlicko2{
		RoomBase:     base,
		teamMgr:      mgr,
		glicko2Teams: make(map[int64]glicko2.Team, base.TeamLimit),
	}
	return r
}

func (r *RoomBaseGlicko2) GetTeams() []glicko2.Team {
	r.RLock()
	defer r.RUnlock()
	teams := make([]glicko2.Team, 0, len(r.glicko2Teams))
	for _, t := range r.glicko2Teams {
		teams = append(teams, t)
	}
	return teams
}

func (r *RoomBaseGlicko2) AddTeam(t glicko2.Team) {
	r.Lock()
	defer r.Unlock()
	r.Base().AddTeam(t.(entry.Team))
	r.glicko2Teams[t.(entry.Team).ID()] = t
}

func (r *RoomBaseGlicko2) GetMMR() float64 {
	total := 0.0
	teams := r.GetTeams()
	if len(teams) == 0 {
		return 0.0
	}
	for _, t := range teams {
		total += t.GetMMR()
	}
	return total / float64(len(teams))
}

func (r *RoomBaseGlicko2) GetStartMatchTimeSec() int64 {
	res := int64(math.MaxInt64)
	for _, t := range r.glicko2Teams {
		groups := t.GetGroups()
		for _, g := range groups {
			if g.GetStartMatchTimeSec() < res {
				res = g.GetStartMatchTimeSec()
			}
		}
	}
	return res
}

func (r *RoomBaseGlicko2) HasAi() bool {
	for _, t := range r.GetTeams() {
		for _, g := range t.GetGroups() {
			for _, p := range g.GetPlayers() {
				if p.IsAi() {
					return true
				}
			}
		}
	}
	return false
}

func (r *RoomBaseGlicko2) SetTeamMgr(mgr TeamMgr) {
	r.teamMgr = mgr
}

func (r *RoomBaseGlicko2) FillGlicko2Teams() {
	r.glicko2Teams = make(map[int64]glicko2.Team, len(r.Base().Teams))
	for id := range r.Base().Teams {
		r.glicko2Teams[id] = r.teamMgr.Get(id).(glicko2.Team)
	}
}
