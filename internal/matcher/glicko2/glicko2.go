package glicko2

import (
	"log/slog"
	"sync"

	"github.com/hedon954/go-matcher/pkg/algorithm/glicko2"
)

type Matcher struct {
	mLock    sync.RWMutex
	matchers map[string]*glicko2.Matcher
	argsFunc map[string]func() *glicko2.QueueArgs
	errChan  chan error
	roomChan chan glicko2.Room

	// for debug
	ErrCount  int
	RoomCount int
}

func New() *Matcher {
	m := &Matcher{
		matchers: make(map[string]*glicko2.Matcher, 8),
		argsFunc: make(map[string]func() *glicko2.QueueArgs, 8),
		errChan:  make(chan error),
		roomChan: make(chan glicko2.Room),
	}

	go m.handleMatchResult()
	return m
}

func (m *Matcher) Stop() {
	for _, matcher := range m.matchers {
		matcher.Stop()
	}
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
func (m *Matcher) NewMatcher(key string,
	argsFunc func() *glicko2.QueueArgs,
	newTeamFunc func(group glicko2.Group) glicko2.Team,
	newRoomFunc, newRoomWithAIFunc func(team glicko2.Team) glicko2.Room) (matcher *glicko2.Matcher, err error) {
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
	go matcher.Match()
	return matcher, nil
}

func (m *Matcher) newMatcher(
	argsFunc func() *glicko2.QueueArgs,
	newTeamFunc func(group glicko2.Group) glicko2.Team,
	newRoomFunc, newRoomWithAIFunc func(team glicko2.Team) glicko2.Room) (*glicko2.Matcher, error) {
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
	slog.Error("glicko2 matcher occurs error", slog.Any("err", err))
}

func (m *Matcher) handleSuccess(room glicko2.Room) {
	m.RoomCount++
	slog.Info("match success", slog.Any("room", room))
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
	matcher, err := m.NewMatcher("", nil, nil, nil, nil)
	if err != nil {
		slog.Error("match by glicko2 error", slog.Any("group", g), slog.String("err", err.Error()))
		return
	}
	if err = matcher.AddGroups(g); err != nil {
		slog.Error("add group to glicko2 error", slog.Any("group", g), slog.String("err", err.Error()))
		return
	}
}
