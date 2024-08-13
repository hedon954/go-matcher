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
	lock          sync.RWMutex
	id            int64 // id is the zutils unique team id.
	TeamID        int   // TeamID is the unique team id in one room, start from 1.
	IsAI          bool
	groups        map[int64]Group
	GameMode      constant.GameMode
	MatchStrategy constant.MatchStrategy
	ModeVersion   int64
}

func NewTeamBase(id int64, g Group) *TeamBase {
	t := &TeamBase{
		id:            id,
		groups:        make(map[int64]Group),
		GameMode:      g.Base().GameMode,
		MatchStrategy: g.Base().MatchStrategy,
		ModeVersion:   g.Base().ModeVersion,
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

func (t *TeamBase) Lock() {
	t.lock.Lock()
}

func (t *TeamBase) Unlock() {
	t.lock.Unlock()
}

func (t *TeamBase) RLock() {
	t.lock.RLock()
}

func (t *TeamBase) RUnlock() {
	t.lock.RUnlock()
}
