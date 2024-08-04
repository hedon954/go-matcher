package glicko2

import (
	"github.com/hedon954/go-matcher/internal/entry"
	"github.com/hedon954/go-matcher/pkg/algorithm/glicko2"
)

type Room struct {
	entry.Room
}

func (r *Room) GetID() int64 {
	return r.ID()
}

func (r *Room) GetTeams() []glicko2.Team {
	//TODO implement me
	panic("implement me")
}

func (r *Room) SortTeamByRank() []glicko2.Team {
	//TODO implement me
	panic("implement me")
}

func (r *Room) AddTeam(t glicko2.Team) {
	//TODO implement me
	panic("implement me")
}

func (r *Room) RemoveTeam(t glicko2.Team) {
	//TODO implement me
	panic("implement me")
}

func (r *Room) GetMMR() float64 {
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

func (r *Room) PlayerCount() int {
	count := 0
	for _, t := range r.GetTeams() {
		count += t.PlayerCount()
	}
	return count
}

func (r *Room) GetStartMatchTimeSec() int64 {
	return r.Base().GetTeams()[0].(*Team).Base().GetGroups()[0].(*Group).GetStartMatchTimeSec()
}

func (r *Room) GetFinishMatchTimeSec() int64 {
	return r.Base().GetTeams()[0].(*Team).Base().GetGroups()[0].(*Group).GetFinishMatchTimeSec()
}

func (r *Room) SetFinishMatchTimeSec(unix int64) {
	for _, t := range r.GetTeams() {
		for _, g := range t.GetGroups() {
			g.SetFinishMatchTimeSec(unix)
		}
	}
}

func (r *Room) HasAi() bool {
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
