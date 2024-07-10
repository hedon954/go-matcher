package glicko2

import (
	"errors"

	"matcher/common"
)

type Player struct {
	*common.Base
}

func NewPlayer(base *common.Base) (*Player, error) {
	if base == nil {
		return nil, errors.New("base player must be inited")
	}
	p := &Player{Base: base}
	// TODO: other op...
	return p, nil
}
