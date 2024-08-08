package mock

import (
	"github.com/hedon954/go-matcher/internal/constant"
)

type MatchStrategyMock struct{}

func (m *MatchStrategyMock) GetMatchStrategy(_ constant.GameMode) constant.MatchStrategy {
	return constant.MatchStrategyGlicko2
}
