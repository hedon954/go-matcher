package glicko2

import (
	"math"

	"github.com/hedon954/go-matcher/internal/entry"
	"github.com/hedon954/go-matcher/pkg/algorithm/glicko2"
)

type RoomBaseGlicko2 struct {
	*entry.RoomBase
}

func CreateRoomBase(base *entry.RoomBase) *RoomBaseGlicko2 {
	r := &RoomBaseGlicko2{
		RoomBase: base,
	}
	return r
}

func (r *RoomBaseGlicko2) GetTeams() []glicko2.Team {
	r.RLock()
	defer r.RUnlock()
	teams := r.Base().GetTeams()
	res := make([]glicko2.Team, len(teams))
	for i := 0; i < len(res); i++ {
		res[i] = teams[i].(glicko2.Team)
	}
	return res
}

func (r *RoomBaseGlicko2) AddTeam(t glicko2.Team) {
	r.Lock()
	defer r.Unlock()
	r.Base().AddTeam(t.(entry.Team))
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
	teams := r.Base().GetTeams()
	if len(teams) == 0 {
		return 0
	}
	for _, t := range teams {
		groups := t.Base().GetGroups()
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
