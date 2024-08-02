package impl

import (
	"testing"

	"github.com/hedon954/go-matcher/internal/constant"
	"github.com/hedon954/go-matcher/internal/entry"
	"github.com/hedon954/go-matcher/internal/merr"
	"github.com/hedon954/go-matcher/internal/pto"

	"github.com/stretchr/testify/assert"
)

const GameMode = 10086
const ModeVersion = 1008611
const MatchStrategy = 10010
const PlayerLimit = 5
const UID = "uid"

func newCreateGroupParam(uid string) *pto.CreateGroup {
	return &pto.CreateGroup{
		PlayerInfo: pto.PlayerInfo{
			UID:           uid,
			GameMode:      GameMode,
			ModeVersion:   ModeVersion,
			MatchStrategy: MatchStrategy,
		},
	}
}

func newEnterGroupParam(uid string) *pto.PlayerInfo {
	return &pto.PlayerInfo{
		UID:           uid,
		GameMode:      GameMode,
		ModeVersion:   ModeVersion,
		MatchStrategy: MatchStrategy,
	}
}

func TestImpl_CreateGroup(t *testing.T) {
	impl := NewDefault(PlayerLimit)
	param := newCreateGroupParam(UID)

	// 1. no group and player, create group should be success
	g, err := impl.CreateGroup(param)
	assert.Nil(t, err)
	assert.NotNil(t, g)
	assert.Equal(t, g.GroupID(), int64(1))
	assert.Equal(t, impl.playerMgr.Get(UID), g.GetCaptain())
	assert.Equal(t, false, g.IsFull())
	assert.Equal(t, entry.GroupStateInvite, g.Base().GetState())
	assert.Equal(t, entry.PlayerOnlineStateInGroup, impl.playerMgr.Get(UID).Base().GetOnlineState())

	// 2. create group with same player info, should be success and return the origin group
	g2, err := impl.CreateGroup(param)
	assert.Nil(t, err)
	assert.NotNil(t, g2)
	assert.Equal(t, g2.GroupID(), int64(1))
	assert.Equal(t, g, g2)
	assert.Equal(t, impl.playerMgr.Get(UID), g2.GetCaptain())
	assert.Equal(t, false, g2.IsFull())

	// 3. change the game mode, should create a new group and dissolve the old group
	param.GameMode = 1
	g3, err := impl.CreateGroup(param)
	assert.Nil(t, err)
	assert.NotNil(t, g3)
	assert.Equal(t, int64(2), g3.GroupID())
	assert.NotEqual(t, g, g3)
	assert.Nil(t, impl.groupMgr.Get(g.GroupID()))
	assert.Equal(t, constant.GameMode(1), g3.GetCaptain().Base().GameMode)
	assert.Equal(t, constant.MatchStrategy(10010), g3.GetCaptain().Base().MatchStrategy)

	// 4. if the player state is not `online` or `group`, should return error
	p2, err := impl.playerMgr.CreatePlayer(newEnterGroupParam(UID + "2"))
	assert.Nil(t, err)
	p2.Base().SetOnlineState(entry.PlayerOnlineStateOffline)
	_, err = impl.CreateGroup(newCreateGroupParam(p2.Base().UID()))
	assert.Equal(t, merr.ErrPlayerOffline, err)
	p2.Base().SetOnlineState(entry.PlayerOnlineStateInMatch)
	_, err = impl.CreateGroup(newCreateGroupParam(p2.Base().UID()))
	assert.Equal(t, merr.ErrPlayerInMatch, err)
	p2.Base().SetOnlineState(entry.PlayerOnlineStateInGame)
	_, err = impl.CreateGroup(newCreateGroupParam(p2.Base().UID()))
	assert.Equal(t, merr.ErrPlayerInGame, err)
	p2.Base().SetOnlineState(entry.PlayerOnlineStateInSettle)
	_, err = impl.CreateGroup(newCreateGroupParam(p2.Base().UID()))
	assert.Equal(t, merr.ErrPlayerInSettle, err)

	// add another player to the group
	p2, err = impl.playerMgr.CreatePlayer(newEnterGroupParam(UID + "2"))
	assert.Nil(t, err)
	err = g3.Base().AddPlayer(p2)
	p2.Base().SetOnlineState(entry.PlayerOnlineStateInGroup)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(g3.Base().GetPlayers()))

	// 5. change the group captain, then create group, should create a new group and exit the old group
	g3.Base().SetCaptain(p2)
	assert.Equal(t, p2, g3.GetCaptain())
	g4, err := impl.CreateGroup(param)
	assert.Nil(t, err)
	assert.Equal(t, int64(3), g4.GroupID())
	assert.Equal(t, 1, len(g3.Base().GetPlayers()))
	assert.Equal(t, 1, len(g4.Base().GetPlayers()))
	assert.Equal(t, entry.GroupStateInvite, g4.Base().GetState())
	assert.Equal(t, entry.PlayerOnlineStateInGroup, impl.playerMgr.Get(UID).Base().GetOnlineState())
}

func TestImpl_ExitGroup(t *testing.T) {
	impl := NewDefault(PlayerLimit)

	// 1. if player is not exists, should return error
	err := impl.ExitGroup(UID)
	assert.Equal(t, merr.ErrPlayerNotInGroup, err)

	// create a temp group
	g, err := impl.CreateGroup(newCreateGroupParam(UID))
	assert.Nil(t, err)
	p := impl.playerMgr.Get(UID)
	assert.Equal(t, g.GroupID(), p.Base().GroupID)
	assert.Equal(t, entry.PlayerOnlineStateInGroup, p.Base().GetOnlineState())
	assert.Equal(t, entry.GroupStateInvite, g.Base().GetState())

	// 2. if group is not existed, should return error
	p.Base().GroupID = 0
	err = impl.ExitGroup(UID)
	assert.Equal(t, merr.ErrPlayerNotInGroup, err)
	p.Base().GroupID = g.GroupID()

	// 3. if player state is not in group, should return error
	p.Base().SetOnlineState(entry.PlayerOnlineStateOffline)
	err = impl.ExitGroup(UID)
	assert.Equal(t, merr.ErrPlayerOffline, err)
	p.Base().SetOnlineState(entry.PlayerOnlineStateOnline)
	err = impl.ExitGroup(UID)
	assert.Equal(t, merr.ErrPlayerNotInGroup, err)
	p.Base().SetOnlineState(entry.PlayerOnlineStateInGame)
	err = impl.ExitGroup(UID)
	assert.Equal(t, merr.ErrPlayerInGame, err)
	p.Base().SetOnlineState(entry.PlayerOnlineStateInMatch)
	err = impl.ExitGroup(UID)
	assert.Equal(t, merr.ErrPlayerInMatch, err)
	p.Base().SetOnlineState(entry.PlayerOnlineStateInSettle)
	err = impl.ExitGroup(UID)
	assert.Equal(t, merr.ErrPlayerInSettle, err)
	p.Base().SetOnlineState(entry.PlayerOnlineStateInGroup)

	// 4. if group state is not in invite, should return error
	g.Base().SetState(entry.GroupStateMatch)
	err = impl.ExitGroup(UID)
	assert.Equal(t, merr.ErrGroupInMatch, err)
	g.Base().SetState(entry.GroupStateInvite)

	// 5. should success, and beacase the group only has one player,
	// both player and group instance should be deleted.
	err = impl.ExitGroup(UID)
	assert.Nil(t, err)
	assert.Nil(t, impl.groupMgr.Get(g.GroupID()))
	assert.Nil(t, impl.playerMgr.Get(p.UID()))

	// create another group
	g, err = impl.CreateGroup(newCreateGroupParam(UID))
	assert.Nil(t, err)
	p = impl.playerMgr.Get(UID)
	assert.Equal(t, g.GroupID(), p.Base().GroupID)
	assert.Equal(t, entry.PlayerOnlineStateInGroup, p.Base().GetOnlineState())
	assert.Equal(t, entry.GroupStateInvite, g.Base().GetState())

	// make the group to have two players
	enterInfo := newEnterGroupParam(UID + "2")
	err = impl.EnterGroup(enterInfo, g.GroupID())
	assert.Nil(t, err)
	assert.Equal(t, 2, len(g.Base().GetPlayers()))

	// 6. if group has multi-players, and the player no captain exit group,
	// should return success and the group capatin should not change,
	// also, the player should be removed from the manager.
	err = impl.ExitGroup(UID + "2")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(g.Base().GetPlayers()))
	assert.Equal(t, UID, g.GetCaptain().UID())
	assert.Nil(t, impl.playerMgr.Get(UID+"2"))
	assert.Equal(t, entry.PlayerOnlineStateInGroup, impl.playerMgr.Get(UID).Base().GetOnlineState())
	assert.Equal(t, entry.GroupStateInvite, g.Base().GetState())
	assert.Equal(t, false, g.Base().PlayerExists(UID+"2"))
	err = impl.EnterGroup(enterInfo, g.GroupID()) // add back
	assert.Nil(t, err)

	// 7. if group has multi-players, and the player is captain exit group,
	// should return success and the group capatin should change
	err = impl.ExitGroup(UID)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(g.Base().GetPlayers()))
	assert.Equal(t, UID+"2", g.GetCaptain().UID())
	assert.Nil(t, impl.playerMgr.Get(UID))
	assert.Equal(t, entry.PlayerOnlineStateInGroup, impl.playerMgr.Get(UID+"2").Base().GetOnlineState())
	assert.Equal(t, entry.GroupStateInvite, g.Base().GetState())
}

func TestImpl_EnterGroup(t *testing.T) {
	impl := NewDefault(PlayerLimit)

	info := newEnterGroupParam(UID)

	const GroupID = 1

	// 1. if the group not exists, should return error
	err := impl.EnterGroup(info, GroupID)
	assert.Equal(t, merr.ErrGroupDissolved, err)

	// create a temp group
	g, err := impl.CreateGroup(newCreateGroupParam(UID + "2"))
	assert.Nil(t, err)
	assert.Equal(t, int64(1), g.GroupID())
	assert.Equal(t, entry.GroupStateInvite, g.Base().GetState())
	assert.Equal(t, UID+"2", g.GetCaptain().UID())

	// 2. if the group state is not invite, should return error
	g.Base().SetState(entry.GroupStateMatch)
	err = impl.EnterGroup(info, GroupID)
	assert.Equal(t, merr.ErrGroupInMatch, err)
	g.Base().SetState(entry.GroupStateInvite)
	assert.Equal(t, UID+"2", g.GetCaptain().UID())
	assert.Equal(t, 1, len(g.Base().GetPlayers()))

	// 3. if the player not exists, should success and create a new player
	err = impl.EnterGroup(info, GroupID)
	assert.Nil(t, err)
	assert.Equal(t, UID+"2", g.GetCaptain().UID())
	assert.Equal(t, 2, len(g.Base().GetPlayers()))
	assert.Equal(t, g.GroupID(), impl.playerMgr.Get(UID).Base().GroupID)
	assert.Equal(t, entry.PlayerOnlineStateInGroup, impl.playerMgr.Get(UID).Base().GetOnlineState())

	// 4. if the player exists target group, and can not play together, should exit the group
	info.GameMode = constant.GameMode(0)
	err = impl.EnterGroup(info, GroupID)
	assert.Equal(t, merr.ErrVersionNotMatch, err)
	assert.Equal(t, UID+"2", g.GetCaptain().UID())
	assert.Equal(t, 1, len(g.Base().GetPlayers()))
	assert.Equal(t, entry.GroupStateInvite, g.Base().GetState())
	assert.Equal(t, int64(0), impl.playerMgr.Get(UID).Base().GroupID)
	assert.Equal(t, entry.PlayerOnlineStateOnline, impl.playerMgr.Get(UID).Base().GetOnlineState())
	info.GameMode = GameMode
	err = impl.EnterGroup(info, GroupID) // add back
	assert.Nil(t, err)

	// 5. if the player exists target group, and can play together should success and not influence the group data
	err = impl.EnterGroup(info, GroupID)
	assert.Nil(t, err)
	assert.Equal(t, UID+"2", g.GetCaptain().UID())
	assert.Equal(t, 2, len(g.Base().GetPlayers()))
	assert.Equal(t, g.GroupID(), impl.playerMgr.Get(UID).Base().GroupID)
	assert.Equal(t, entry.PlayerOnlineStateInGroup, impl.playerMgr.Get(UID).Base().GetOnlineState())

	// 6. if player exists, and the player state is not `online` or `group`, should return error
	p := impl.playerMgr.Get(UID)
	p.Base().SetOnlineState(entry.PlayerOnlineStateInSettle)
	err = impl.EnterGroup(info, GroupID)
	assert.Equal(t, merr.ErrPlayerInSettle, err)
	p.Base().SetOnlineState(entry.PlayerOnlineStateInGroup) // set back

	// make UID into another group
	g2, err := impl.CreateGroup(newCreateGroupParam(UID))
	assert.Nil(t, err)
	assert.Equal(t, int64(2), g2.GroupID())
	assert.Equal(t, 1, len(g.Base().GetPlayers()))

	// 7. if the player exists in other group, and can not play together, should return error, and not exits the origin group
	info.GameMode = 2
	err = impl.EnterGroup(info, GroupID)
	assert.Equal(t, merr.ErrVersionNotMatch, err)
	assert.Equal(t, g2.GroupID(), impl.playerMgr.Get(UID).Base().GroupID)
	assert.Equal(t, 1, len(g2.Base().GetPlayers()))
	assert.Equal(t, UID, g2.GetCaptain().UID())

	// 8. if the player exists in other group, and can play together, should success and exit the old group
	info.GameMode = GameMode
	err = impl.EnterGroup(info, GroupID)
	assert.Nil(t, err)
	assert.Equal(t, g.GroupID(), impl.playerMgr.Get(UID).Base().GroupID)
	assert.Equal(t, 2, len(g.Base().GetPlayers()))
	assert.Equal(t, UID+"2", g.GetCaptain().UID())
	assert.Equal(t, 0, len(g2.Base().GetPlayers()))
	assert.Nil(t, impl.groupMgr.Get(g2.GroupID()))

	// 9. add another 4 players to group, the last one should return error
	err = impl.EnterGroup(newEnterGroupParam(UID+"3"), g.GroupID())
	assert.Nil(t, err)
	err = impl.EnterGroup(newEnterGroupParam(UID+"4"), g.GroupID())
	assert.Nil(t, err)
	err = impl.EnterGroup(newEnterGroupParam(UID+"5"), g.GroupID())
	assert.Nil(t, err)
	err = impl.EnterGroup(newEnterGroupParam(UID+"6"), g.GroupID())
	assert.Equal(t, merr.ErrGroupFull, err)
}

func TestImpl_DissolveGroup(t *testing.T) {
	impl := NewDefault(PlayerLimit)

	// create a group and make it have multi-players
	g, err := impl.CreateGroup(newCreateGroupParam(UID))
	assert.Nil(t, err)
	assert.Equal(t, 1, len(g.Base().GetPlayers()))
	assert.Equal(t, UID, g.GetCaptain().UID())
	assert.Equal(t, int64(1), g.GroupID())
	assert.Equal(t, entry.GroupStateInvite, g.Base().GetState())
	assert.Equal(t, entry.PlayerOnlineStateInGroup, impl.playerMgr.Get(UID).Base().GetOnlineState())

	err = impl.EnterGroup(newEnterGroupParam(UID+"2"), g.GroupID())
	assert.Nil(t, err)
	assert.Equal(t, 2, len(g.Base().GetPlayers()))

	// 1. if the player not exists, should return error
	err = impl.DissolveGroup(UID + "3")
	assert.Equal(t, merr.ErrPlayerNotInGroup, err)

	// 2. if the group not exists, should return error
	impl.playerMgr.Add(UID+"3", entry.NewPlayerBase(&pto.PlayerInfo{}))
	err = impl.DissolveGroup(UID + "3")
	assert.Equal(t, merr.ErrGroupNotExists, err)
	impl.playerMgr.Delete(UID + "3") // delete back

	// 3. if the group state is not `invite`, should return error
	g.Base().SetState(entry.GroupStateMatch)
	err = impl.DissolveGroup(UID)
	assert.Equal(t, merr.ErrGroupInMatch, err)
	g.Base().SetState(entry.GroupStateGame)
	err = impl.DissolveGroup(UID)
	assert.Equal(t, merr.ErrGroupInGame, err)
	g.Base().SetState(entry.GroupStateInvite) // set back

	// 4. if the player is not the captain, should return error
	err = impl.DissolveGroup(UID + "2")
	assert.Equal(t, merr.ErrOnlyCaptainCanDissolveGroup, err)

	// 5. if the player is the captain, should success,
	// and all the group's players should be deleted from the player manager.
	err = impl.DissolveGroup(UID)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(g.Base().GetPlayers()))
	assert.Equal(t, entry.GroupStateDissolved, g.Base().GetState())
	assert.Nil(t, impl.groupMgr.Get(g.GroupID()))
	assert.Nil(t, impl.playerMgr.Get(UID))
	assert.Nil(t, impl.playerMgr.Get(UID+"2"))
}

func TestImpl_KickPlayer(t *testing.T) {
	impl := NewDefault(PlayerLimit)

	// create a temp group
	g, err := impl.CreateGroup(newCreateGroupParam(UID))
	assert.Nil(t, err)
	assert.Equal(t, 1, len(g.Base().GetPlayers()))
	assert.Equal(t, UID, g.GetCaptain().UID())

	// 1. if captain equals to kicked player, should return error
	err = impl.KickPlayer(UID, UID)
	assert.Equal(t, merr.ErrKickSelf, err)

	// 2. if the captain not exists, should return error
	err = impl.KickPlayer(UID+"1", UID+"2")
	assert.Equal(t, merr.ErrPlayerNotExists, err)

	// 3. if the kicked player not exists, should return error
	err = impl.KickPlayer(UID, UID+"2")
	assert.Equal(t, merr.ErrPlayerNotExists, err)

	// add UID+"2" to group
	err = impl.EnterGroup(newEnterGroupParam(UID+"2"), g.GroupID())
	assert.Nil(t, err)
	assert.Equal(t, 2, len(g.Base().GetPlayers()))
	assert.Equal(t, entry.PlayerOnlineStateInGroup, impl.playerMgr.Get(UID+"2").Base().GetOnlineState())

	// 4. if the group not exists, should return error
	impl.groupMgr.Delete(g.GroupID())
	err = impl.KickPlayer(UID, UID+"2")
	assert.Equal(t, merr.ErrGroupNotExists, err)
	impl.groupMgr.Add(g.GroupID(), g) // add back

	// 5. if the kicked player not in group, should return false
	g2, err := impl.CreateGroup(newCreateGroupParam(UID + "3"))
	assert.Nil(t, err)
	err = impl.KickPlayer(UID, UID+"3")
	assert.Equal(t, merr.ErrPlayerNotInGroup, err)
	assert.Equal(t, UID+"3", g2.GetCaptain().UID())

	// 5. if the group state is not `invite`, should return error
	g.Base().SetState(entry.GroupStateMatch)
	err = impl.KickPlayer(UID, UID+"2")
	assert.Equal(t, merr.ErrGroupInMatch, err)
	g.Base().SetState(entry.GroupStateGame)
	err = impl.KickPlayer(UID, UID+"2")
	assert.Equal(t, merr.ErrGroupInGame, err)
	g.Base().SetState(entry.GroupStateDissolved)
	err = impl.KickPlayer(UID, UID+"2")
	assert.Equal(t, merr.ErrGroupDissolved, err)
	g.Base().SetState(entry.GroupStateInvite) // set back

	// 6. if the player is not the captain, should return error
	err = impl.KickPlayer(UID+"2", UID)
	assert.Equal(t, merr.ErrOnlyCaptainCanKickPlayer, err)

	// 7. if the player is the captain, should success,
	err = impl.KickPlayer(UID, UID+"2")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(g.Base().GetPlayers()))
	assert.Equal(t, entry.PlayerOnlineStateInGroup, impl.playerMgr.Get(UID).Base().GetOnlineState())
	assert.Equal(t, UID, g.GetCaptain().UID())
	assert.Nil(t, impl.playerMgr.Get(UID+"2"))
}

func TestImpl_HandoverCaptain(t *testing.T) {
}
