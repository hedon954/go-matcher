package matcher

import (
	"fmt"
	"runtime/debug"

	"github.com/rs/zerolog/log"

	"github.com/hedon954/go-matcher/internal/constant"
	"github.com/hedon954/go-matcher/internal/entry"
	"github.com/hedon954/go-matcher/internal/matcher/glicko2"

	glicko2Algo "github.com/hedon954/go-matcher/pkg/algorithm/glicko2"
)

type Matcher struct {
	Glicko2Matcher *glicko2.Matcher
	groupChannel   chan entry.Group
}

func New(groupChannel chan entry.Group, glicko2Match *glicko2.Matcher) *Matcher {
	m := &Matcher{
		groupChannel:   groupChannel,
		Glicko2Matcher: glicko2Match,
	}
	return m
}

func (m *Matcher) Start() {
	log.Info().Msg("start matcher")
	fmt.Println("start matcher")
	defer log.Info().Msg("stop matcher")
	defer fmt.Println("stop matcher")

	for g := range m.groupChannel {
		m.handle(g)
	}
}

func (m *Matcher) Stop() {
	m.Glicko2Matcher.Stop()
}

func (m *Matcher) handle(g entry.Group) {
	defer func() {
		if err := recover(); err != nil {
			log.Error().
				Any("err", err).
				Str("group", g.Json()).
				Str("stack", string(debug.Stack())).
				Msg("handle group match error")
		}
	}()

	switch g.Base().MatchStrategy {
	case constant.MatchStrategyGlicko2:
		m.Glicko2Matcher.Match(g.(glicko2Algo.Group))
	default:
		log.Error().
			Int64("group_id", g.ID()).
			Int("game_mode", int(g.Base().GameMode)).
			Any("group", g).
			Int("strategy", int(g.Base().MatchStrategy)).
			Msg("unknown match strategy")
	}
}
