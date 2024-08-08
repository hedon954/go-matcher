package config

import (
	"github.com/hedon954/go-matcher/internal/constant"
)

type MatchStrategy interface {
	GetMatchStrategy(mode constant.GameMode) constant.MatchStrategy
}
