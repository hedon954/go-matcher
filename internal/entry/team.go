package entry

import (
	"sync"

	"github.com/hedon954/go-matcher/internal/constant"
)

type Team interface {
	Base() *TeamBase
	ID() int64
}

type TeamBase struct {
	sync.RWMutex
	id     int64
	groups map[int64]Group
}

func NewTeamBase(id int64, g Group) *TeamBase {
	t := &TeamBase{
		id:     id,
		groups: make(map[int64]Group),
	}
	t.groups[g.ID()] = g
	return t
}

func (t *TeamBase) Base() *TeamBase {
	return t
}

func (t *TeamBase) ID() int64 {
	return t.id
}

func (t *TeamBase) MatchStrategy() constant.MatchStrategy {
	return t.randGroup().Base().MatchStrategy
}

func (t *TeamBase) GameMode() constant.GameMode {
	return t.randGroup().Base().GameMode
}

func (t *TeamBase) GetGroups() []Group {
	res := make([]Group, len(t.groups))
	i := 0
	for _, g := range t.groups {
		res[i] = g
		i++
	}
	return res
}

func (t *TeamBase) AddGroup(g Group) {
	t.groups[g.ID()] = g
}

func (t *TeamBase) RemoveGroup(id int64) {
	delete(t.groups, id)
}

func (t *TeamBase) randGroup() Group {
	for _, g := range t.groups {
		return g
	}
	return nil
}
