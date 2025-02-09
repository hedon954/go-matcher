package test_game

import (
	"github.com/hedon954/go-matcher/internal/constant"
	"github.com/hedon954/go-matcher/internal/entry"
	"github.com/hedon954/go-matcher/internal/pto"
)

func RegisterFactory() {
	entry.RegisterFactory(constant.GameModeTest, &factory{})
}

type factory struct{}

func (f *factory) CreatePlayer(mgr *entry.Mgrs, base *entry.PlayerBase, pInfo *pto.PlayerInfo) (entry.Player, error) {
	return CreatePlayer(base), nil
}

func (f *factory) CreateGroup(mgr *entry.Mgrs, base *entry.GroupBase) (entry.Group, error) {
	return CreateGroup(base), nil
}

func (f *factory) CreateTeam(mgr *entry.Mgrs, base *entry.TeamBase) (entry.Team, error) {
	return CreateTeam(base), nil
}

func (f *factory) CreateRoom(mgr *entry.Mgrs, base *entry.RoomBase) (entry.Room, error) {
	return CreateRoom(base), nil
}
