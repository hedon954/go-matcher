package entry

import (
	"sync"

	"github.com/hedon954/go-matcher/internal/constant"
	"github.com/hedon954/go-matcher/internal/pto"
	"github.com/hedon954/go-matcher/pkg/typeconv"
)

type Room interface {
	Coder
	Base() *RoomBase
	ID() int64
	NeedAI() bool
	GetMatchInfo() *pto.MatchInfo
}

type RoomBase struct {
	lock      sync.RWMutex
	RoomID    int64
	Teams     map[int64]struct{}
	TeamLimit int

	GameMode       constant.GameMode
	MatchStrategy  constant.MatchStrategy
	ModeVersion    int64
	FinishMatchSec int64

	EscapePlayer []string

	GameServerInfo pto.GameServerInfo
}

func NewRoomBase(id int64, teamLimit int, t Team) *RoomBase {
	r := &RoomBase{
		RoomID:        id,
		TeamLimit:     teamLimit,
		Teams:         make(map[int64]struct{}),
		EscapePlayer:  make([]string, 0),
		GameMode:      t.Base().GameMode,
		MatchStrategy: t.Base().MatchStrategy,
		ModeVersion:   t.Base().ModeVersion,
	}
	return r
}

func (r *RoomBase) Base() *RoomBase {
	return r
}

func (r *RoomBase) ID() int64 {
	return r.RoomID
}

func (r *RoomBase) GetTeams() []int64 {
	return typeconv.MapToSlice(r.Teams)
}

func (r *RoomBase) AddTeam(t Team) {
	r.Teams[t.ID()] = struct{}{}
}

func (r *RoomBase) RemoveTeam(id int64) {
	delete(r.Teams, id)
}

func (r *RoomBase) GetMatchInfo() *pto.MatchInfo {
	// TODO
	return nil
}

func (r *RoomBase) NeedAI() bool {
	return false
}

func (r *RoomBase) AddEscapePlayer(uid string) {
	r.EscapePlayer = append(r.EscapePlayer, uid)
}

func (r *RoomBase) GetEscapePlayers() []string {
	return r.EscapePlayer
}

func (r *RoomBase) Lock() {
	r.lock.Lock()
}

func (r *RoomBase) Unlock() {
	r.lock.Unlock()
}

func (r *RoomBase) RLock() {
	r.lock.RLock()
}

func (r *RoomBase) RUnlock() {
	r.lock.RUnlock()
}
