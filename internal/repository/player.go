package repository

import (
	"github.com/hedon954/go-matcher/internal/entry"
	"github.com/hedon954/go-matcher/internal/pto"
	"github.com/hedon954/go-matcher/pkg/collection"
)

type PlayerMgr struct {
	*collection.Manager[string, entry.Player]
}

func NewPlayerMgr() *PlayerMgr {
	return &PlayerMgr{Manager: collection.New[string, entry.Player]()}
}

func (m *PlayerMgr) CreatePlayer(pInfo *pto.PlayerInfo) (entry.Player, error) {
	// TODO: factory method

	bp := entry.NewPlayerBase(pInfo)

	m.Add(bp.UID(), bp)

	return bp, nil
}
