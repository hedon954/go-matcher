package repository

import (
	"fmt"

	"github.com/hedon954/go-matcher/internal/constant"
	"github.com/hedon954/go-matcher/internal/entry"
	"github.com/hedon954/go-matcher/internal/entry/goat_game"
	"github.com/hedon954/go-matcher/internal/pto"
	"github.com/hedon954/go-matcher/pkg/collection"
)

type PlayerMgr struct {
	*collection.Manager[string, entry.Player]
}

func NewPlayerMgr() *PlayerMgr {
	return &PlayerMgr{Manager: collection.New[string, entry.Player]()}
}

func (m *PlayerMgr) CreatePlayer(pInfo *pto.PlayerInfo) (p entry.Player, err error) {
	base := entry.NewPlayerBase(pInfo)

	switch base.GameMode {
	case constant.GameModeGoatGame:
		p = goat_game.CreatePlayer(base, pInfo)
	case constant.GameModeTest:
		p = base
	default:
		return nil, fmt.Errorf("unsupported game mode: %d", base.GameMode)
	}

	m.Add(p.UID(), p)
	return p, nil
}
