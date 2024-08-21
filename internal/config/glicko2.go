package config

import (
	"github.com/hedon954/go-matcher/internal/constant"
	"github.com/hedon954/go-matcher/pkg/algorithm/glicko2"
)

type Glicko2 interface {
	GetGlicko2QueueArgs(mode constant.GameMode) *glicko2.QueueArgs
}
