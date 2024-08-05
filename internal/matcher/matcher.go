package matcher

import (
	"log/slog"
	"runtime/debug"

	"github.com/hedon954/go-matcher/internal/constant"
	"github.com/hedon954/go-matcher/internal/entry"
	"github.com/hedon954/go-matcher/internal/matcher/glicko2"

	glicko2Algo "github.com/hedon954/go-matcher/pkg/algorithm/glicko2"
)

type Matcher struct {
	glicko2Matcher *glicko2.Matcher
	groupChannel   chan entry.Group
}

func New(groupChannel chan entry.Group, glicko2Match *glicko2.Matcher) *Matcher {
	m := &Matcher{
		groupChannel:   groupChannel,
		glicko2Matcher: glicko2Match,
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
			slog.Error("handle group match error: \n"+string(debug.Stack()), slog.Any("group", g), slog.Any("err", err))
		}
	}()

	switch g.Base().MatchStrategy {
	case constant.MatchStrategyGlicko2:
		m.glicko2Matcher.Match(g.(glicko2Algo.Group))
	default:
		slog.Error("unknown match strategy", slog.Any("group", g), slog.Int("strategy", int(g.Base().MatchStrategy)))
	}
}
