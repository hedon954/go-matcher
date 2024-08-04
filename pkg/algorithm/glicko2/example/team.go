package example

import (
	"sort"
	"sync"

	"github.com/hedon954/go-matcher/pkg/algorithm/glicko2"
)

type Team struct {
	sync.RWMutex

	Id int

	groups            map[string]glicko2.Group
	StartMatchTimeSec int64
	rank              int

	WinWeight int
}

func (t *Team) IsFull(teamPlayerLimit int) bool {
	return t.PlayerCount() >= teamPlayerLimit
}

func (t *Team) IsNewer() bool {
	return false
}

func NewTeam(group glicko2.Group) glicko2.Team {
	t := &Team{
		RWMutex: sync.RWMutex{},
		groups:  make(map[string]glicko2.Group),
	}
	t.AddGroup(group)
	return t
}

func (t *Team) GetModeVersion() int {
	return 0
}

func (t *Team) GetGroups() []glicko2.Group {
	res := make([]glicko2.Group, len(t.groups))
	i := 0
	for _, g := range t.groups {
		res[i] = g
		i++
	}
	return res
}

func (t *Team) AddGroup(g glicko2.Group) {
	t.groups[g.GetID()] = g
	gmst := g.GetStartMatchTimeSec()
	if gmst == 0 {
		return
	}
	if t.StartMatchTimeSec == 0 || t.StartMatchTimeSec > gmst {
		t.StartMatchTimeSec = gmst
	}
}

func (t *Team) RemoveGroup(groupId string) {
	delete(t.groups, groupId)
}

func (t *Team) PlayerCount() int {
	count := 0
	for _, group := range t.groups {
		count += group.PlayerCount()
	}
	return count
}

func (t *Team) GetMMR() float64 {
	if len(t.groups) == 0 {
		return 0
	}
	total := 0.0
	for _, group := range t.groups {
		total += group.GetMMR()
	}
	return total / float64(len(t.groups))
}

func (t *Team) GetStar() int {
	if len(t.groups) == 0 {
		return 0
	}
	rank := 0
	for _, g := range t.groups {
		rank += g.GetStar()
	}
	return rank / len(t.groups)
}

func (t *Team) SetFinishMatchTimeSec(t2 int64) {
	for _, g := range t.groups {
		g.SetFinishMatchTimeSec(t2)
	}
}

func (t *Team) GetStartMatchTimeSec() int64 {
	return t.StartMatchTimeSec
}

func (t *Team) GetFinishMatchTimeSec() int64 {
	for _, g := range t.groups {
		return g.GetFinishMatchTimeSec()
	}
	return 0
}

func (t *Team) IsAi() bool {
	for _, g := range t.groups {
		for _, p := range g.GetPlayers() {
			if p.IsAi() {
				return true
			}
		}
	}
	return false
}

func (t *Team) GetRank() int {
	return t.rank
}

func (t *Team) SetRank(rank int) {
	t.rank = rank
}

func (t *Team) SortPlayerByRank() []glicko2.Player {
	players := make([]glicko2.Player, 0, 5)
	for _, g := range t.groups {
		players = append(players, g.GetPlayers()...)
	}
	sort.SliceStable(players, func(i, j int) bool {
		return players[i].GetRank() < players[j].GetRank()
	})
	return players
}

func (t *Team) CanFillAi() bool {
	for _, g := range t.groups {
		if !g.CanFillAi() {
			return false
		}
	}
	return true
}
