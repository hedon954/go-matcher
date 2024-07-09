package glicko2

import (
	"errors"

	"matcher/common"
)

type Player struct {
	*common.Base
}

func (p *Player) Inner() *common.Base {
	return p.Base
}

func (p *Player) UID() string {
	return p.Base.UID()
}

func NewPlayer(base *common.Base) (*Player, error) {
	if base == nil {
		return nil, errors.New("base player must be inited")
	}
	p := &Player{Base: base}
	// TODO: other op...
	return p, nil
}
