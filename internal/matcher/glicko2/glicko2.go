package glicko2

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/hedon954/go-matcher/internal/config"
	"github.com/hedon954/go-matcher/internal/constant"
	"github.com/hedon954/go-matcher/internal/entry"
	"github.com/hedon954/go-matcher/internal/repository"
	"github.com/hedon954/go-matcher/pkg/algorithm/glicko2"
)

// Matcher is the glicko2 matcher.
type Matcher struct {
	mLock sync.RWMutex

	configer *Configer

	playerMgr *repository.PlayerMgr
	groupMgr  *repository.GroupMgr
	teamMgr   *repository.TeamMgr
	roomMgr   *repository.RoomMgr

	// matchers is the map of glicko2 matchers.
	// `key` is used to separate different matching groups.
	// `value` is the glicko2 matcher.
	matchers map[string]*glicko2.Matcher

	// errChan is a channel for handle error form glicko2 matcher.
	errChan chan error

	// roomChan is a channel for handle room form glicko2 matcher.
	roomChan chan glicko2.Room

	// roomChannelToService is a channel for send room to service.
	roomChannelToService chan entry.Room

	// gameModes is the map of game modes, `value` is the funcs of the mode.
	gameModes map[constant.GameMode]*Funcs

	// for debug
	ErrCount      int
	RoomCount     atomic.Int64
	matchInterval time.Duration
}

type Configer struct {
	Glicko2 config.Glicko2
}

// Funcs is the funcs needed for glicko2 matcher.
type Funcs struct {
	ArgsFunc          func() *glicko2.QueueArgs
	NewTeamFunc       func(group glicko2.Group) glicko2.Team
	NewRoomFunc       func(team glicko2.Team) glicko2.Room
	NewRoomWithAIFunc func(team glicko2.Team) glicko2.Room
}

// New returns the new glicko2 matcher, and start it.
func New(
	roomChannelToService chan entry.Room,
	configer *Configer, matchInterval time.Duration,
	playerMgr *repository.PlayerMgr, groupMgr *repository.GroupMgr,
	teamMgr *repository.TeamMgr, roomMgr *repository.RoomMgr,
) *Matcher {
	m := &Matcher{
		matchers:             make(map[string]*glicko2.Matcher, 8),
		errChan:              make(chan error),
		roomChan:             make(chan glicko2.Room),
		roomChannelToService: roomChannelToService,
		gameModes:            make(map[constant.GameMode]*Funcs, 16),
		playerMgr:            playerMgr,
		groupMgr:             groupMgr,
		teamMgr:              teamMgr,
		roomMgr:              roomMgr,
		configer:             configer,
		matchInterval:        matchInterval,
	}

	// register funcs
	m.registerGoatGame()

	// start to handle match result
	go m.handleMatchResult()

	return m
}

// Stop stops all matchers.
func (m *Matcher) Stop() {
	for _, matcher := range m.matchers {
		matcher.Stop()
	}
}

// AddMode adds the funcs of the given mode.
func (m *Matcher) AddMode(mode constant.GameMode, funcs *Funcs) {
	m.Lock()
	defer m.Unlock()
	m.gameModes[mode] = funcs
}

// GetFuncs returns the funcs of the given mode.
func (m *Matcher) GetFuncs(mode constant.GameMode) *Funcs {
	m.RLock()
	defer m.RUnlock()
	return m.gameModes[mode]
}

// GetMatcher returns the matcher of the given key.
func (m *Matcher) GetMatcher(key string) *glicko2.Matcher {
	m.RLock()
	defer m.RUnlock()
	return m.matchers[key]
}

// NewMatcher returns the new matcher of the given key,
// if the key exists, it will return the existing one.
// `key` is used to separate different matching groups.
func (m *Matcher) NewMatcher(
	key string,
	argsFunc func() *glicko2.QueueArgs,
	newTeamFunc func(group glicko2.Group) glicko2.Team,
	newRoomFunc, newRoomWithAIFunc func(team glicko2.Team) glicko2.Room,
) (matcher *glicko2.Matcher, err error) {
	m.Lock()
	defer m.Unlock()
	matcher = m.matchers[key]
	if matcher != nil {
		return matcher, nil
	}

	matcher, err = m.newMatcher(argsFunc, newTeamFunc, newRoomFunc, newRoomWithAIFunc)
	if err != nil {
		return nil, err
	}

	m.matchers[key] = matcher
	go matcher.Match(m.matchInterval)
	return matcher, nil
}

func (m *Matcher) newMatcher(
	argsFunc func() *glicko2.QueueArgs,
	newTeamFunc func(group glicko2.Group) glicko2.Team,
	newRoomFunc, newRoomWithAIFunc func(team glicko2.Team) glicko2.Room,
) (*glicko2.Matcher, error) {
	return glicko2.NewMatcher(m.errChan, m.roomChan, argsFunc,
		newTeamFunc, newRoomFunc, newRoomWithAIFunc)
}

func (m *Matcher) handleMatchResult() {
	for {
		select {
		case err := <-m.errChan:
			m.handleError(err)
		case room := <-m.roomChan:
			m.handleSuccess(room)
		}
	}
}

func (m *Matcher) handleError(err error) {
	m.ErrCount++
	log.Error().Err(err).Msg("glicko2 matcher occurs error")
}

func (m *Matcher) handleSuccess(room glicko2.Room) {
	m.RoomCount.Add(1)
	log.Info().Any("room", room).Msg("glicko2 match success")
	m.roomChannelToService <- room.(entry.Room)
}

func (m *Matcher) Lock() {
	m.mLock.Lock()
}

func (m *Matcher) Unlock() {
	m.mLock.Unlock()
}

func (m *Matcher) RLock() {
	m.mLock.RLock()
}

func (m *Matcher) RUnlock() {
	m.mLock.RUnlock()
}

func (m *Matcher) Match(g glicko2.Group) {
	funcs := m.GetFuncs(g.(entry.Group).Base().GameMode)
	if funcs == nil {
		panic(fmt.Sprintf("game mode glicko2 funcs not register: %d", g.(entry.Group).Base().GameMode))
	}

	matcher, err := m.NewMatcher(g.QueueKey(),
		funcs.ArgsFunc, funcs.NewTeamFunc, funcs.NewRoomFunc, funcs.NewRoomWithAIFunc)
	if err != nil {
		log.Error().
			Any("group", g).
			Err(err).
			Msg("match by glicko2 error")
		return
	}
	fmt.Println("add a group to glicko matcher: ", g.GetPlayers())
	if err = matcher.AddGroups(g); err != nil {
		log.Error().
			Str("group_id", g.GetID()).
			Err(err).
			Msg("add group to glicko2 error")
		return
	}
}
