package entry

import (
	"fmt"

	"github.com/hedon954/go-matcher/internal/pto"
)

type Mgrs struct {
	PlayerMgr *PlayerMgr
	GroupMgr  *GroupMgr
	TeamMgr   *TeamMgr
	RoomMgr   *RoomMgr
}

func (m *Mgrs) CreatePlayer(pInfo *pto.PlayerInfo) (p Player, err error) {
	base := NewPlayerBase(pInfo)

	factory := GetFactory(base.GameMode)
	if factory == nil {
		return nil, fmt.Errorf("unsupported game mode: %d", base.GameMode)
	}
	p, err = factory.CreatePlayer(m, base, pInfo)
	if err != nil {
		return nil, err
	}

	m.PlayerMgr.Add(p.UID(), p)
	return p, nil
}

func (m *Mgrs) CreateGroup(playerLimit int, p Player) (g Group, err error) {
	base := NewGroupBase(m.GroupMgr.GenGroupID(), playerLimit, p.Base())

	factory := GetFactory(base.GameMode)
	if factory == nil {
		return nil, fmt.Errorf("unsupported game mode: %d", base.GameMode)
	}
	g, err = factory.CreateGroup(m, base)
	if err != nil {
		return nil, err
	}

	_ = g.Base().AddPlayer(p)
	m.GroupMgr.Add(g.ID(), g)
	return g, nil
}

func (m *Mgrs) CreateTeam(g Group) (t Team, err error) {
	base := NewTeamBase(m.TeamMgr.teamIDIter.Add(1), g)

	factory := GetFactory(base.GameMode)
	if factory == nil {
		return nil, fmt.Errorf("unsupported game mode: %d", base.GameMode)
	}
	t, err = factory.CreateTeam(m, base)
	if err != nil {
		return nil, err
	}

	t.Base().AddGroup(g)

	// NOTE: don't add team to manager here.
	// because it may be created in match process for temp,
	// only add it after match success.
	// m.Add(t.Base().ID(), t)
	return t, nil
}

func (m *Mgrs) CreateAITeam(g Group) (t Team, err error) {
	base := NewTeamBase(m.TeamMgr.teamIDIter.Add(1), g)
	base.Base().RemoveGroup(g.ID())
	return t, nil
}

func (m *Mgrs) CreateRoom(teamLimit int, t Team) (r Room, err error) {
	base := NewRoomBase(m.RoomMgr.roomIDIter.Add(1), teamLimit, t)

	factory := GetFactory(base.GameMode)
	if factory == nil {
		return nil, fmt.Errorf("unsupported game mode: %d", base.GameMode)
	}
	r, err = factory.CreateRoom(m, base)
	if err != nil {
		return nil, err
	}

	r.Base().AddTeam(t)

	// NOTE: don't add room to manager here.
	// because it may be created in match process for temp,
	// only add it after match success.
	// m.Add(r.Base().ID(), r)
	return r, nil
}
