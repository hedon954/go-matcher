package glicko2

import (
	"errors"

	"github.com/hedon954/go-matcher/common"
)

type Player struct {
	*common.PlayerBase
}

func NewPlayer(base *common.PlayerBase) (*Player, error) {
	if base == nil {
		return nil, errors.New("base player must be inited")
	}
	p := &Player{PlayerBase: base}
	// TODO: other op...
	return p, nil
}
