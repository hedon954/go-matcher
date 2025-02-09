package goat_game

import (
	"fmt"

	"github.com/hedon954/go-matcher/internal/constant"
	"github.com/hedon954/go-matcher/internal/entry"
	"github.com/hedon954/go-matcher/internal/pto"
)

func RegisterFactory() {
	entry.RegisterFactory(constant.GameModeGoatGame, &factory{})
}

type factory struct{}

func (f *factory) CreatePlayer(mgr *entry.Mgrs, base *entry.PlayerBase, pInfo *pto.PlayerInfo) (entry.Player, error) {
	p := &Player{}
	// ... other common fields

	if pInfo.Glicko2Info == nil {
		return nil, fmt.Errorf("game[%d] need glicko2 info", base.GameMode)
	}
	p.withMatchStrategy(base, pInfo.Glicko2Info)
	return p, nil
}

func (f *factory) CreateGroup(mgr *entry.Mgrs, base *entry.GroupBase) (entry.Group, error) {
	g := &Group{}
	// ... other common fields

	g.withMatchStrategy(base, mgr.PlayerMgr)
	return g, nil
}

func (f *factory) CreateTeam(mgr *entry.Mgrs, base *entry.TeamBase) (entry.Team, error) {
	t := &Team{}

	t.withMatchStrategy(base, mgr.GroupMgr)
	return t, nil
}

func (f *factory) CreateRoom(mgr *entry.Mgrs, base *entry.RoomBase) (entry.Room, error) {
	room := &Room{}

	room.withMatchStrategy(base, mgr.TeamMgr)
	return room, nil
}
