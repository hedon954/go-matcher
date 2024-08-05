package config

import (
	"github.com/hedon954/go-matcher/internal/constant"
	"github.com/hedon954/go-matcher/pkg/algorithm/glicko2"
)

type Glicko2 interface {
	GetQueueArgs(mode constant.GameMode) *glicko2.QueueArgs
}

type Glicko2Config struct {
}

func (gc *Glicko2Config) GetQueueArgs(mode constant.GameMode) *glicko2.QueueArgs {
	if mode == constant.GameModeGoatGame {

	}
	return nil
}
