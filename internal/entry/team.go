package entry

import (
	"sync"

	"github.com/hedon954/go-matcher/internal/constant"
	"github.com/hedon954/go-matcher/pkg/typeconv"
)

type Team interface {
	Coder
	Base() *TeamBase
	ID() int64
}

type TeamBase struct {
	lock          sync.RWMutex
	UniqueID      int64 // UniqueID is the global unique team id.
	TeamID        int   // TeamID is the unique team id in one room, start from 1.
	IsAI          bool
	Groups        map[int64]struct{}
	GameMode      constant.GameMode
	MatchStrategy constant.MatchStrategy
	ModeVersion   int64
}

func NewTeamBase(id int64, g Group) *TeamBase {
	t := &TeamBase{
		UniqueID:      id,
		Groups:        make(map[int64]struct{}),
		GameMode:      g.Base().GameMode,
		MatchStrategy: g.Base().MatchStrategy,
		ModeVersion:   g.Base().ModeVersion,
	}
	t.Groups[g.ID()] = struct{}{}
	return t
}

func (t *TeamBase) Base() *TeamBase {
	return t
}

func (t *TeamBase) ID() int64 {
	return t.UniqueID
}

func (t *TeamBase) GetGroups() []int64 {
	return typeconv.MapToSlice(t.Groups)
}

func (t *TeamBase) AddGroup(g Group) {
	t.Groups[g.ID()] = struct{}{}
}

func (t *TeamBase) RemoveGroup(id int64) {
	delete(t.Groups, id)
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
