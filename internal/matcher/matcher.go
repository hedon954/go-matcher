package matcher

import (
	"log/slog"

	"github.com/hedon954/go-matcher/internal/entry"
)

type Matcher struct {
	groupChannel chan entry.Group
}

func New(groupChannel chan entry.Group) *Matcher {
	m := &Matcher{
		groupChannel: groupChannel,
	}
	return m
}

func (m *Matcher) Start() {
	go func() {
		for g := range m.groupChannel {
			m.handle(g)
		}
	}()
}

func (m *Matcher) handle(g entry.Group) {
	defer func() {
		if err := recover(); err != nil {
			slog.Error("handle group match error", slog.Any("group", g), slog.Any("err", err))
		}
	}()

	// TODO
}
