package glicko2

import (
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

func (r *RoomBaseGlicko2) GetID() int64 {
	return r.ID()
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

func (r *RoomBaseGlicko2) SortTeamByRank() []glicko2.Team {
	r.RLock()
	defer r.RUnlock()
	return r.GetTeams()
}

func (r *RoomBaseGlicko2) AddTeam(t glicko2.Team) {
	r.Lock()
	defer r.Unlock()
	r.Base().AddTeam(t.(entry.Team))
}

func (r *RoomBaseGlicko2) RemoveTeam(t glicko2.Team) {
	r.Lock()
	defer r.Unlock()
	r.Base().RemoveTeam(t.(entry.Team).ID())
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

func (r *RoomBaseGlicko2) PlayerCount() int {
	count := 0
	for _, t := range r.GetTeams() {
		count += t.PlayerCount()
	}
	return count
}

func (r *RoomBaseGlicko2) GetStartMatchTimeSec() int64 {
	return r.Base().GetTeams()[0].(*TeamBaseGlicko2).Base().GetGroups()[0].(*GroupBaseGlicko2).GetStartMatchTimeSec()
}

func (r *RoomBaseGlicko2) GetFinishMatchTimeSec() int64 {
	return r.Base().GetTeams()[0].(*TeamBaseGlicko2).Base().GetGroups()[0].(*GroupBaseGlicko2).GetFinishMatchTimeSec()
}

func (r *RoomBaseGlicko2) SetFinishMatchTimeSec(unix int64) {
	for _, t := range r.GetTeams() {
		for _, g := range t.GetGroups() {
			g.SetFinishMatchTimeSec(unix)
		}
	}
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
