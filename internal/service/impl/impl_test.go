package impl

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hedon954/go-matcher/internal/constant"
	"github.com/hedon954/go-matcher/internal/entry"
	"github.com/hedon954/go-matcher/internal/merr"
	"github.com/hedon954/go-matcher/internal/pto"
	"github.com/hedon954/go-matcher/internal/repository"
	"github.com/stretchr/testify/assert"
)

const GameMode = constant.GameModeGoatGame
const ModeVersion = 1008611
const MatchStrategy = constant.MatchStrategyGlicko2
const PlayerLimit = 5
const UID = "uid"

func defaultImpl(playerLimit int, opts ...Option) *Impl {
	gc := make(chan entry.Group, 1024)
	rc := make(chan entry.Room, 1024)
	pm := repository.NewPlayerMgr()
	gm := repository.NewGroupMgr(0)
	return NewDefault(playerLimit, pm, gm, gc, rc, opts...)
}

func newCreateGroupParam(uid string) *pto.CreateGroup {
	return &pto.CreateGroup{
		PlayerInfo: pto.PlayerInfo{
			UID:           uid,
			GameMode:      GameMode,
			ModeVersion:   ModeVersion,
			MatchStrategy: MatchStrategy,
			Glicko2Info:   &pto.Glicko2Info{},
		},
	}
}

func newEnterGroupParam(uid string) *pto.EnterGroup {
	return &pto.EnterGroup{
		PlayerInfo: *newPlayerInfo(uid),
	}
}

func newEnterGroupParamWithSrc(uid string, source pto.EnterGroupSourceType) *pto.EnterGroup {
	return &pto.EnterGroup{
		PlayerInfo: *newPlayerInfo(uid),
		Source:     source,
	}
}

func newPlayerInfo(uid string) *pto.PlayerInfo {
	return &pto.PlayerInfo{
		UID:           uid,
		GameMode:      GameMode,
		ModeVersion:   ModeVersion,
		MatchStrategy: MatchStrategy,
		Glicko2Info:   &pto.Glicko2Info{},
	}
}

func createTempGroup(uid string, impl *Impl, t *testing.T) (entry.Player, entry.Group) {
	g, err := impl.CreateGroup(newCreateGroupParam(uid))
	assert.Nil(t, err)
	p := impl.playerMgr.Get(uid)
	assert.Equal(t, g.ID(), p.Base().GroupID)
	assert.Equal(t, 1, len(g.Base().GetPlayers()))
	assert.Equal(t, entry.PlayerOnlineStateInGroup, p.Base().GetOnlineState())
	assert.Equal(t, entry.GroupStateInvite, g.Base().GetState())
	return p, g
}

func createFullGroup(impl *Impl, t *testing.T) (entry.Player, entry.Group) {
	captainUID := uuid.NewString()
	p, g := createTempGroup(captainUID, impl, t)
	for i := 0; i < PlayerLimit-1; i++ {
		err := impl.EnterGroup(newEnterGroupParam(uuid.NewString()), g.ID())
		assert.Nil(t, err)
	}
	assert.Equal(t, true, g.IsFull())
	return p, g
}

func TestImpl_CreateGroup(t *testing.T) {
	impl := defaultImpl(PlayerLimit)
	param := newCreateGroupParam(UID)

	// 1. no group and player, create group should be success
	g, err := impl.CreateGroup(param)
	assert.Nil(t, err)
	assert.NotNil(t, g)
	assert.Equal(t, g.ID(), int64(1))
	assert.Equal(t, impl.playerMgr.Get(UID), g.GetCaptain())
	assert.Equal(t, false, g.IsFull())
	assert.Equal(t, entry.GroupStateInvite, g.Base().GetState())
	assert.Equal(t, entry.PlayerOnlineStateInGroup, impl.playerMgr.Get(UID).Base().GetOnlineState())

	// 2. create group with same player info, should be success and return the origin group
	g2, err := impl.CreateGroup(param)
	assert.Nil(t, err)
	assert.NotNil(t, g2)
	assert.Equal(t, g2.ID(), int64(1))
	assert.Equal(t, g, g2)
	assert.Equal(t, impl.playerMgr.Get(UID), g2.GetCaptain())
	assert.Equal(t, false, g2.IsFull())

	// 3. change the game mode, should create a new group and dissolve the old group
	param.GameMode = constant.GameModeTest
	g3, err := impl.CreateGroup(param)
	assert.Nil(t, err)
	assert.NotNil(t, g3)
	assert.Equal(t, int64(2), g3.ID())
	assert.NotEqual(t, g, g3)
	assert.Nil(t, impl.groupMgr.Get(g.ID()))
	assert.Equal(t, constant.GameModeTest, g3.GetCaptain().Base().GameMode)
	assert.Equal(t, constant.MatchStrategyGlicko2, g3.GetCaptain().Base().MatchStrategy)

	// 4. if the player state is not `online` or `group`, should return error
	p2, err := impl.playerMgr.CreatePlayer(newPlayerInfo(UID + "2"))
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
	p2, err = impl.playerMgr.CreatePlayer(newPlayerInfo(UID + "2"))
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
	assert.Equal(t, int64(3), g4.ID())
	assert.Equal(t, 1, len(g3.Base().GetPlayers()))
	assert.Equal(t, 1, len(g4.Base().GetPlayers()))
	assert.Equal(t, entry.GroupStateInvite, g4.Base().GetState())
	assert.Equal(t, entry.PlayerOnlineStateInGroup, impl.playerMgr.Get(UID).Base().GetOnlineState())
}

func TestImpl_ExitGroup(t *testing.T) {
	impl := defaultImpl(PlayerLimit)

	// 1. if player is not exists, should return error
	err := impl.ExitGroup(UID)
	assert.Equal(t, merr.ErrPlayerNotExists, err)

	// create a temp group
	p, g := createTempGroup(UID, impl, t)

	// 2. if group is not existed, should return error
	p.Base().GroupID = 0
	err = impl.ExitGroup(UID)
	assert.Equal(t, merr.ErrGroupNotExists, err)
	p.Base().GroupID = g.ID()

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

	// 5. should success, and because the group only has one player,
	// both player and group instance should be deleted.
	err = impl.ExitGroup(UID)
	assert.Nil(t, err)
	assert.Nil(t, impl.groupMgr.Get(g.ID()))
	assert.Nil(t, impl.playerMgr.Get(p.UID()))

	// create another group
	g, err = impl.CreateGroup(newCreateGroupParam(UID))
	assert.Nil(t, err)
	p = impl.playerMgr.Get(UID)
	assert.Equal(t, g.ID(), p.Base().GroupID)
	assert.Equal(t, entry.PlayerOnlineStateInGroup, p.Base().GetOnlineState())
	assert.Equal(t, entry.GroupStateInvite, g.Base().GetState())

	// make the group to have two players
	enterInfo := newEnterGroupParam(UID + "2")
	err = impl.EnterGroup(enterInfo, g.ID())
	assert.Nil(t, err)
	assert.Equal(t, 2, len(g.Base().GetPlayers()))

	// 6. if group has multi-players, and the player no captain exit group,
	// should return success and the group captain should not change,
	// also, the player should be removed from the repository.
	err = impl.ExitGroup(UID + "2")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(g.Base().GetPlayers()))
	assert.Equal(t, UID, g.GetCaptain().UID())
	assert.Nil(t, impl.playerMgr.Get(UID+"2"))
	assert.Equal(t, entry.PlayerOnlineStateInGroup, impl.playerMgr.Get(UID).Base().GetOnlineState())
	assert.Equal(t, entry.GroupStateInvite, g.Base().GetState())
	assert.Equal(t, false, g.Base().PlayerExists(UID+"2"))
	err = impl.EnterGroup(enterInfo, g.ID()) // add back
	assert.Nil(t, err)

	// 7. if group has multi-players, and the player is captain exit group,
	// should return success and the group captain should change
	err = impl.ExitGroup(UID)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(g.Base().GetPlayers()))
	assert.Equal(t, UID+"2", g.GetCaptain().UID())
	assert.Nil(t, impl.playerMgr.Get(UID))
	assert.Equal(t, entry.PlayerOnlineStateInGroup, impl.playerMgr.Get(UID+"2").Base().GetOnlineState())
	assert.Equal(t, entry.GroupStateInvite, g.Base().GetState())
}

func TestImpl_EnterGroup(t *testing.T) {
	impl := defaultImpl(PlayerLimit)

	info := newEnterGroupParam(UID)

	const GroupID = 1

	// 1. if the group not exists, should return error
	err := impl.EnterGroup(info, GroupID)
	assert.Equal(t, merr.ErrGroupDissolved, err)

	// create a temp group
	_, g := createTempGroup(UID+"2", impl, t)

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
	assert.Equal(t, g.ID(), impl.playerMgr.Get(UID).Base().GroupID)
	assert.Equal(t, entry.PlayerOnlineStateInGroup, impl.playerMgr.Get(UID).Base().GetOnlineState())

	// 4. if the player exists target group, and can not play together, should exit the group
	// 4.1 game mode not match
	info.GameMode = constant.GameMode(0)
	err = impl.EnterGroup(info, GroupID)
	assert.Equal(t, merr.ErrGameModeNotMatch, err)
	assert.Equal(t, UID+"2", g.GetCaptain().UID())
	assert.Equal(t, 1, len(g.Base().GetPlayers()))
	assert.Equal(t, entry.GroupStateInvite, g.Base().GetState())
	assert.Equal(t, int64(0), impl.playerMgr.Get(UID).Base().GroupID)
	assert.Equal(t, entry.PlayerOnlineStateOnline, impl.playerMgr.Get(UID).Base().GetOnlineState())
	info.GameMode = GameMode
	err = impl.EnterGroup(info, GroupID) // add back
	assert.Nil(t, err)
	// 4.2 group version to low
	info.ModeVersion = ModeVersion + 100
	err = impl.EnterGroup(info, GroupID)
	assert.Equal(t, merr.ErrGroupVersionTooLow, err)
	assert.Equal(t, UID+"2", g.GetCaptain().UID())
	assert.Equal(t, 1, len(g.Base().GetPlayers()))
	assert.Equal(t, entry.GroupStateInvite, g.Base().GetState())
	assert.Equal(t, int64(0), impl.playerMgr.Get(UID).Base().GroupID)
	assert.Equal(t, entry.PlayerOnlineStateOnline, impl.playerMgr.Get(UID).Base().GetOnlineState())
	info.ModeVersion = ModeVersion
	err = impl.EnterGroup(info, GroupID) // add back
	assert.Nil(t, err)
	// 4.3 player version to low
	info.ModeVersion = ModeVersion - 100
	err = impl.EnterGroup(info, GroupID)
	assert.Equal(t, merr.ErrPlayerVersionTooLow, err)
	assert.Equal(t, UID+"2", g.GetCaptain().UID())
	assert.Equal(t, 1, len(g.Base().GetPlayers()))
	assert.Equal(t, entry.GroupStateInvite, g.Base().GetState())
	assert.Equal(t, int64(0), impl.playerMgr.Get(UID).Base().GroupID)
	assert.Equal(t, entry.PlayerOnlineStateOnline, impl.playerMgr.Get(UID).Base().GetOnlineState())
	info.ModeVersion = ModeVersion
	err = impl.EnterGroup(info, GroupID) // add back
	assert.Nil(t, err)

	// 5. if the player exists target group, and can play together,
	// should success, then update the player info, and do not influence the group data.
	originPlayer := impl.playerMgr.Get(UID)
	originPlayer.Base().ModeVersion = ModeVersion - 1
	assert.NotEqual(t, originPlayer.GetPlayerInfo().ModeVersion, info.ModeVersion)
	err = impl.EnterGroup(info, GroupID)
	assert.Nil(t, err)
	assert.Equal(t, UID+"2", g.GetCaptain().UID())
	assert.Equal(t, 2, len(g.Base().GetPlayers()))
	assert.Equal(t, g.ID(), impl.playerMgr.Get(UID).Base().GroupID)
	assert.Equal(t, entry.PlayerOnlineStateInGroup, impl.playerMgr.Get(UID).Base().GetOnlineState())
	assert.Equal(t, info.ModeVersion, impl.playerMgr.Get(UID).GetPlayerInfo().ModeVersion)

	// 6. if player exists, and the player state is not `online` or `group`, should return error
	p := impl.playerMgr.Get(UID)
	p.Base().SetOnlineState(entry.PlayerOnlineStateInSettle)
	err = impl.EnterGroup(info, GroupID)
	assert.Equal(t, merr.ErrPlayerInSettle, err)
	p.Base().SetOnlineState(entry.PlayerOnlineStateInGroup) // set back

	// make UID into another group
	g2, err := impl.CreateGroup(newCreateGroupParam(UID))
	assert.Nil(t, err)
	assert.Equal(t, int64(2), g2.ID())
	assert.Equal(t, 1, len(g.Base().GetPlayers()))

	// 7. if the player exists in other group, and can not play together, should return error, and not exits the origin group
	info.GameMode = 2
	err = impl.EnterGroup(info, GroupID)
	assert.Equal(t, merr.ErrGameModeNotMatch, err)
	assert.Equal(t, g2.ID(), impl.playerMgr.Get(UID).Base().GroupID)
	assert.Equal(t, 1, len(g2.Base().GetPlayers()))
	assert.Equal(t, UID, g2.GetCaptain().UID())

	// 8. if the player exists in other group, and can play together, should success and exit the old group
	info.GameMode = GameMode
	err = impl.EnterGroup(info, GroupID)
	assert.Nil(t, err)
	assert.Equal(t, g.ID(), impl.playerMgr.Get(UID).Base().GroupID)
	assert.Equal(t, 2, len(g.Base().GetPlayers()))
	assert.Equal(t, UID+"2", g.GetCaptain().UID())
	assert.Equal(t, 0, len(g2.Base().GetPlayers()))
	assert.Nil(t, impl.groupMgr.Get(g2.ID()))

	// 9. if group denies to enter from nearby and source is `pto.EnterGroupSourceTypeNearby`, should return error
	g.Base().SetAllowNearbyJoin(false)
	err = impl.EnterGroup(newEnterGroupParamWithSrc(UID+"3", pto.EnterGroupSourceTypeNearby), g.ID())
	assert.Equal(t, merr.ErrGroupDenyNearbyJoin, err)

	// 10. if group denies to enter from recent and source is `pto.EnterGroupSourceTypeRecent`, should return error
	g.Base().SetAllowRecentJoin(false)
	err = impl.EnterGroup(newEnterGroupParamWithSrc(UID+"3", pto.EnterGroupSourceTypeRecent), g.ID())
	assert.Equal(t, merr.ErrGroupDenyRecentJoin, err)

	// 11. add another 4 players to group, the last one should return error
	err = impl.EnterGroup(newEnterGroupParam(UID+"3"), g.ID())
	assert.Nil(t, err)
	err = impl.EnterGroup(newEnterGroupParam(UID+"4"), g.ID())
	assert.Nil(t, err)
	err = impl.EnterGroup(newEnterGroupParam(UID+"5"), g.ID())
	assert.Nil(t, err)
	err = impl.EnterGroup(newEnterGroupParam(UID+"6"), g.ID())
	assert.Equal(t, merr.ErrGroupFull, err)
}

func TestImpl_DissolveGroup(t *testing.T) {
	impl := defaultImpl(PlayerLimit)

	// create a group and make it have multi-players
	_, g := createTempGroup(UID, impl, t)
	err := impl.EnterGroup(newEnterGroupParam(UID+"2"), g.ID())
	assert.Nil(t, err)
	assert.Equal(t, 2, len(g.Base().GetPlayers()))

	// 1. if the player not exists, should return error
	err = impl.DissolveGroup(UID + "3")
	assert.Equal(t, merr.ErrPlayerNotExists, err)

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
	// and all the group's players should be deleted from the player repository.
	err = impl.DissolveGroup(UID)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(g.Base().GetPlayers()))
	assert.Equal(t, entry.GroupStateDissolved, g.Base().GetState())
	assert.Nil(t, impl.groupMgr.Get(g.ID()))
	assert.Nil(t, impl.playerMgr.Get(UID))
	assert.Nil(t, impl.playerMgr.Get(UID+"2"))
}

func TestImpl_KickPlayer(t *testing.T) {
	impl := defaultImpl(PlayerLimit)

	// create a temp group
	_, g := createTempGroup(UID, impl, t)

	// 1. if captain equals to kicked player, should return error
	err := impl.KickPlayer(UID, UID)
	assert.Equal(t, merr.ErrKickSelf, err)

	// 2. if the captain not exists, should return error
	err = impl.KickPlayer(UID+"1", UID+"2")
	assert.Equal(t, merr.ErrPlayerNotExists, err)

	// 3. if the kicked player not exists, should return error
	err = impl.KickPlayer(UID, UID+"2")
	assert.Equal(t, merr.ErrPlayerNotExists, err)

	// add UID+"2" to group
	err = impl.EnterGroup(newEnterGroupParam(UID+"2"), g.ID())
	assert.Nil(t, err)
	assert.Equal(t, 2, len(g.Base().GetPlayers()))
	assert.Equal(t, entry.PlayerOnlineStateInGroup, impl.playerMgr.Get(UID+"2").Base().GetOnlineState())

	// 4. if the group not exists, should return error
	impl.groupMgr.Delete(g.ID())
	err = impl.KickPlayer(UID, UID+"2")
	assert.Equal(t, merr.ErrGroupNotExists, err)
	impl.groupMgr.Add(g.ID(), g) // add back

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

func TestImpl_ChangeRole(t *testing.T) {
	impl := defaultImpl(PlayerLimit)

	// create two temp group
	_, g := createTempGroup(UID, impl, t)
	_, _ = createTempGroup(UID+"3", impl, t)
	err := impl.EnterGroup(newEnterGroupParam(UID+"2"), g.ID())
	assert.Nil(t, err)

	// 1. if the captain not exists, should return error
	err = impl.ChangeRole(UID+"1", UID+"2", entry.GroupRoleCaptain)
	assert.Equal(t, merr.ErrPlayerNotExists, err)

	// 2. if the target player not exists, should return error
	err = impl.ChangeRole(UID, UID+"1", entry.GroupRoleCaptain)
	assert.Equal(t, merr.ErrPlayerNotExists, err)

	// 3. if the group not exists, should return error
	impl.groupMgr.Delete(g.ID()) // delete temp
	err = impl.ChangeRole(UID, UID+"2", entry.GroupRoleCaptain)
	assert.Equal(t, merr.ErrGroupNotExists, err)
	impl.groupMgr.Add(g.ID(), g) // add back

	// 4. if the group state is not `invite`, should return error
	g.Base().SetState(entry.GroupStateMatch) // set temp
	err = impl.ChangeRole(UID, UID+"2", entry.GroupRoleCaptain)
	assert.Equal(t, merr.ErrGroupInMatch, err)
	g.Base().SetState(entry.GroupStateGame)
	err = impl.ChangeRole(UID, UID+"2", entry.GroupRoleCaptain)
	assert.Equal(t, merr.ErrGroupInGame, err)
	g.Base().SetState(entry.GroupStateDissolved)
	err = impl.ChangeRole(UID, UID+"2", entry.GroupRoleCaptain)
	assert.Equal(t, merr.ErrGroupDissolved, err)
	g.Base().SetState(entry.GroupStateInvite) // set back

	// 5. if the target player not in group, should return error
	err = impl.ChangeRole(UID, UID+"3", entry.GroupRoleCaptain)
	assert.Equal(t, merr.ErrPlayerNotInGroup, err)

	// 6. if current player not the captain, should return error
	err = impl.ChangeRole(UID+"2", UID, entry.GroupRoleCaptain)
	assert.Equal(t, merr.ErrNotCaptain, err)

	// 7. failed to change to unknown role
	err = impl.ChangeRole(UID, UID+"2", entry.GroupRole(-1))
	assert.Equal(t, errors.New("unsupported role: -1"), err)

	// 8. success to change role to entry.GroupRoleCaptain
	testImplChangeRoleCaptain(impl, g, t)
}

func testImplChangeRoleCaptain(impl *Impl, g entry.Group, t *testing.T) {
	// 7. success
	err := impl.ChangeRole(UID, UID+"2", entry.GroupRoleCaptain)
	assert.Nil(t, err)
	assert.Equal(t, UID+"2", g.GetCaptain().UID())
	assert.Equal(t, 2, len(g.Base().GetPlayers()))
	assert.Equal(t, entry.GroupStateInvite, g.Base().GetState())
}

func TestImpl_SetNearbyJoinGroup(t *testing.T) {
	impl := defaultImpl(PlayerLimit)

	// create a temp group
	_, g := createTempGroup(UID, impl, t)
	assert.False(t, g.Base().AllowNearbyJoin())

	// 1. if the player is not exists, should return error
	err := impl.SetNearbyJoinGroup(UID+"1", true)
	assert.Equal(t, merr.ErrPlayerNotExists, err)

	// 2. if the group is not exists, should return error
	impl.groupMgr.Delete(g.ID()) // delete temp
	err = impl.SetNearbyJoinGroup(UID, true)
	assert.Equal(t, merr.ErrGroupNotExists, err)
	impl.groupMgr.Add(g.ID(), g) // add back

	// 3. if the player is not the capatin, should return error
	err = impl.EnterGroup(newEnterGroupParam(UID+"2"), g.ID())
	assert.Nil(t, err)
	err = impl.SetNearbyJoinGroup(UID+"2", true)
	assert.Equal(t, merr.ErrPermissionDeny, err)

	// 4. set true success
	err = impl.SetNearbyJoinGroup(UID, true)
	assert.Nil(t, err)
	assert.True(t, g.Base().AllowNearbyJoin())

	// 5. set false success
	err = impl.SetNearbyJoinGroup(UID, false)
	assert.Nil(t, err)
	assert.False(t, g.Base().AllowNearbyJoin())
}

func TestImpl_SetRecentJoinGroup(t *testing.T) {
	impl := defaultImpl(PlayerLimit)

	// create a temp group
	_, g := createTempGroup(UID, impl, t)
	assert.False(t, g.Base().AllowRecentJoin())

	// 1. if the player is not exists, should return error
	err := impl.SetRecentJoinGroup(UID+"1", true)
	assert.Equal(t, merr.ErrPlayerNotExists, err)

	// 2. if the group is not exists, should return error
	impl.groupMgr.Delete(g.ID()) // delete temp
	err = impl.SetRecentJoinGroup(UID, true)
	assert.Equal(t, merr.ErrGroupNotExists, err)
	impl.groupMgr.Add(g.ID(), g) // add back

	// 3. if the player is not the capatin, should return error
	err = impl.EnterGroup(newEnterGroupParam(UID+"2"), g.ID())
	assert.Nil(t, err)
	err = impl.SetRecentJoinGroup(UID+"2", true)
	assert.Equal(t, merr.ErrPermissionDeny, err)

	// 4. set true success
	err = impl.SetRecentJoinGroup(UID, true)
	assert.Nil(t, err)
	assert.True(t, g.Base().AllowRecentJoin())

	// 5. set false success
	err = impl.SetRecentJoinGroup(UID, false)
	assert.Nil(t, err)
	assert.False(t, g.Base().AllowRecentJoin())
}

func TestImpl_Invite(t *testing.T) {
	const nowSec = 100
	nowFunc := func() int64 { return nowSec }
	impl := defaultImpl(PlayerLimit, WithNowFunc(nowFunc))

	// create a temp group
	_, g := createTempGroup(UID, impl, t)

	// 1. if the inviter is not exists, should return error
	err := impl.Invite(UID+"1", UID+"2")
	assert.Equal(t, merr.ErrPlayerNotExists, err)

	// 2. if the group is not exists, should return error
	impl.groupMgr.Delete(g.ID()) // delete temp
	err = impl.Invite(UID, UID+"2")
	assert.Equal(t, merr.ErrGroupNotExists, err)
	impl.groupMgr.Add(g.ID(), g) // add back

	// 3. if the group state is not `invite`, should return error
	g.Base().SetState(entry.GroupStateMatch)
	err = impl.Invite(UID, UID+"2")
	assert.Equal(t, merr.ErrGroupInMatch, err)
	g.Base().SetState(entry.GroupStateGame)
	err = impl.Invite(UID, UID+"2")
	assert.Equal(t, merr.ErrGroupInGame, err)
	g.Base().SetState(entry.GroupStateDissolved)
	err = impl.Invite(UID, UID+"2")
	assert.Equal(t, merr.ErrGroupDissolved, err)
	g.Base().SetState(entry.GroupStateInvite) // set back

	// 4. if the group is full, should return error
	p2, _ := createFullGroup(impl, t)
	err = impl.Invite(p2.UID(), UID)
	assert.Equal(t, merr.ErrGroupFull, err)

	// 5. invite success and save invite record
	err = impl.Invite(UID, UID+"2")
	assert.Nil(t, err)
	assert.Equal(t, int64(nowSec+entry.InviteExpireSec), g.Base().GetInviteExpireTimeStamp(UID+"2"))
}

func TestImpl_RefuseInvite(t *testing.T) {
	impl := defaultImpl(PlayerLimit)

	// 1. if the group not exists, just return nil
	impl.RefuseInvite(UID, UID+"1", 1, "")

	// 2. return the group state is dissolved, just return nil
	p, g := createTempGroup(UID+"1", impl, t)
	g.Base().SetState(entry.GroupStateDissolved) // set temp
	impl.RefuseInvite(UID, UID+"1", 1, "")
	g.Base().SetState(entry.GroupStateInvite) // set back

	// 3. push refuse msg to the inviter and the invite record should be deleted
	err := impl.Invite(p.UID(), UID)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(g.Base().GetInviteRecords()))
	impl.RefuseInvite(p.UID(), UID, g.ID(), "")
	assert.Equal(t, 0, len(g.Base().GetInviteRecords()))
}

func TestImpl_AcceptInvite(t *testing.T) {
	impl := defaultImpl(PlayerLimit)

	inviteeInfo := newPlayerInfo(UID + "1")

	// 1. if the group not exists, should return err
	err := impl.AcceptInvite(UID, inviteeInfo, 1)
	assert.Equal(t, merr.ErrGroupDissolved, err)

	inviter, g := createTempGroup(UID, impl, t)
	err = impl.Invite(inviter.UID(), inviteeInfo.UID)
	assert.Nil(t, err)

	// 2. if the inviter is not in the group, should return err
	err = impl.AcceptInvite(UID+"2", inviteeInfo, g.ID())
	assert.Equal(t, merr.ErrInvitationExpired, err)

	// 3. if the group state is not `entry.GroupStateInvite`, should return err
	g.Base().SetState(entry.GroupStateMatch)
	err = impl.AcceptInvite(inviter.UID(), inviteeInfo, g.ID())
	assert.Equal(t, merr.ErrGroupInMatch, err)
	g.Base().SetState(entry.GroupStateGame)
	err = impl.AcceptInvite(inviter.UID(), inviteeInfo, g.ID())
	assert.Equal(t, merr.ErrGroupInGame, err)
	g.Base().SetState(entry.GroupStateDissolved)
	err = impl.AcceptInvite(inviter.UID(), inviteeInfo, g.ID())
	assert.Equal(t, merr.ErrGroupDissolved, err)
	g.Base().SetState(entry.GroupStateInvite) // set back

	// 4. if the invitee's state is not either `online` or `group`, should return err
	impl.playerMgr.Add(UID+"1", entry.NewPlayerBase(newPlayerInfo(UID+"1"))) // add temp player
	invitee := impl.playerMgr.Get(UID + "1")
	assert.NotNil(t, invitee)
	invitee.Base().SetOnlineState(entry.PlayerOnlineStateOffline)
	err = impl.AcceptInvite(inviter.UID(), invitee.GetPlayerInfo(), g.ID())
	assert.Equal(t, merr.ErrPlayerOffline, err)
	invitee.Base().SetOnlineState(entry.PlayerOnlineStateInMatch)
	err = impl.AcceptInvite(inviter.UID(), invitee.GetPlayerInfo(), g.ID())
	assert.Equal(t, merr.ErrPlayerInMatch, err)
	invitee.Base().SetOnlineState(entry.PlayerOnlineStateInGame)
	err = impl.AcceptInvite(inviter.UID(), invitee.GetPlayerInfo(), g.ID())
	assert.Equal(t, merr.ErrPlayerInGame, err)
	invitee.Base().SetOnlineState(entry.PlayerOnlineStateInSettle)
	err = impl.AcceptInvite(inviter.UID(), invitee.GetPlayerInfo(), g.ID())
	assert.Equal(t, merr.ErrPlayerInSettle, err)
	impl.playerMgr.Delete(UID + "1") // delete temp player

	// 5. if invitee can not player with group, should return err
	// 5.1 player version too low
	err = impl.Invite(inviter.UID(), invitee.UID())
	assert.Nil(t, err)
	invitee.Base().ModeVersion = ModeVersion - 1
	err = impl.AcceptInvite(inviter.UID(), invitee.GetPlayerInfo(), g.ID())
	assert.Equal(t, merr.ErrPlayerVersionTooLow, err)
	// 5.2 group version too low
	err = impl.Invite(inviter.UID(), invitee.UID())
	assert.Nil(t, err)
	invitee.Base().ModeVersion = ModeVersion + 1
	err = impl.AcceptInvite(inviter.UID(), invitee.GetPlayerInfo(), g.ID())
	assert.Equal(t, merr.ErrGroupVersionTooLow, err)
	// 5.3 game mode not match
	err = impl.Invite(inviter.UID(), invitee.UID())
	assert.Nil(t, err)
	invitee.Base().GameMode = -1
	err = impl.AcceptInvite(inviter.UID(), invitee.GetPlayerInfo(), g.ID())
	assert.Equal(t, merr.ErrGameModeNotMatch, err)
	invitee.Base().ModeVersion = ModeVersion // set back
	invitee.Base().GameMode = GameMode       // set back

	// 6. if the invitation has expired, should return err
	// 6.1 not invitation record
	g.Base().DelInviteRecord(invitee.UID()) // delete temp
	err = impl.AcceptInvite(inviter.UID(), invitee.GetPlayerInfo(), g.ID())
	assert.Equal(t, merr.ErrInvitationExpired, err)
	// 6.2 invitation expired
	g.Base().AddInviteRecord(invitee.UID(), time.Now().Unix()-entry.InviteExpireSec-1)
	err = impl.AcceptInvite(inviter.UID(), invitee.GetPlayerInfo(), g.ID())
	assert.Equal(t, merr.ErrInvitationExpired, err)
	g.Base().AddInviteRecord(invitee.UID(), time.Now().Unix()) // set back

	// 7. if the group is full, should return err
	fullGroupInviter, fullGroup := createFullGroup(impl, t) // create a temp full group
	err = impl.Invite(fullGroup.GetCaptain().UID(), inviteeInfo.UID)
	assert.Equal(t, merr.ErrGroupFull, err)
	err = impl.ExitGroup(fullGroupInviter.UID()) // exit one player
	assert.Nil(t, err)
	err = impl.Invite(fullGroup.GetCaptain().UID(), inviteeInfo.UID) // send invite
	assert.Nil(t, err)
	err = impl.EnterGroup(newEnterGroupParam(fullGroupInviter.UID()),
		fullGroup.ID()) // another player enter the group make group full
	assert.Nil(t, err)
	err = impl.AcceptInvite(fullGroup.GetCaptain().UID(), inviteeInfo,
		fullGroup.ID()) // now the group is full, so can not accept invite
	assert.Equal(t, merr.ErrGroupFull, err)

	// 8. return success and delete the invite record
	err = impl.AcceptInvite(inviter.UID(), invitee.GetPlayerInfo(), g.ID())
	assert.Nil(t, err)
	assert.Equal(t, int64(0), g.Base().GetInviteExpireTimeStamp(invitee.UID()))
}

func TestImpl_SetVoiceState(t *testing.T) {
	impl := defaultImpl(PlayerLimit)

	// 1. if player not exists, should return err
	err := impl.SetVoiceState(UID, entry.PlayerVoiceStateMute)
	assert.Equal(t, merr.ErrPlayerNotExists, err)

	p, g := createTempGroup(UID, impl, t)
	assert.Equal(t, entry.PlayerVoiceStateMute, p.Base().GetVoiceState())

	// 2. if group not exists, should return err
	impl.groupMgr.Delete(g.ID()) // delete temp
	err = impl.SetVoiceState(UID, entry.PlayerVoiceStateMute)
	assert.Equal(t, merr.ErrGroupNotExists, err)
	impl.groupMgr.Add(g.ID(), g) // add back

	// 3. if state no change, should do nothing and success
	err = impl.SetVoiceState(UID, entry.PlayerVoiceStateMute)
	assert.Nil(t, err)
	assert.Equal(t, entry.PlayerVoiceStateMute, p.Base().GetVoiceState())

	// 4. change state successfully
	err = impl.SetVoiceState(UID, entry.PlayerVoiceStateUnmute)
	assert.Nil(t, err)
	assert.Equal(t, entry.PlayerVoiceStateUnmute, p.Base().GetVoiceState())
}

func TestImpl_StartMatch(t *testing.T) {
	impl := defaultImpl(PlayerLimit)

	// 1. if player not exists, should return err
	err := impl.StartMatch(UID)
	assert.Equal(t, merr.ErrPlayerNotExists, err)

	p, g := createTempGroup(UID, impl, t)
	err = impl.StartMatch(p.UID())
	assert.Nil(t, err)
	assert.Equal(t, entry.GroupStateMatch, g.Base().GetState())

	// 2. if group not exists, should return err
	impl.groupMgr.Delete(g.ID()) // delete temp
	err = impl.StartMatch(p.UID())
	assert.Equal(t, merr.ErrGroupNotExists, err)
	impl.groupMgr.Add(g.ID(), g) // add back

	// 3. if group state is not `invite`, should return err
	g.Base().SetState(entry.GroupStateGame)
	err = impl.StartMatch(p.UID())
	assert.Equal(t, merr.ErrGroupInGame, err)
	g.Base().SetState(entry.GroupStateDissolved)
	err = impl.StartMatch(p.UID())
	assert.Equal(t, merr.ErrGroupDissolved, err)
	g.Base().SetState(entry.GroupStateMatch)
	err = impl.StartMatch(p.UID())
	assert.Equal(t, merr.ErrGroupInMatch, err)
	g.Base().SetState(entry.GroupStateInvite) // set back

	// 4. if the player is not captain, should return err
	err = impl.EnterGroup(newEnterGroupParam(UID+"1"), g.ID())
	assert.Nil(t, err)
	assert.Equal(t, 2, len(g.Base().GetPlayers()))
	err = impl.StartMatch(UID + "1")
	assert.Equal(t, merr.ErrNotCaptain, err)

	// 5. if the player is captain, should success and the group state and players' state should change to `match`
	err = impl.StartMatch(UID)
	assert.Nil(t, err)
	assert.Equal(t, entry.GroupStateMatch, g.Base().GetState())
	assert.Equal(t, 2, len(g.Base().GetPlayers()))
	assert.Equal(t, entry.PlayerOnlineStateInMatch, p.Base().GetOnlineState())
	assert.Equal(t, entry.PlayerOnlineStateInMatch, impl.playerMgr.Get(UID+"1").Base().GetOnlineState())
}

func TestImpl_CancelMatch(t *testing.T) {
	impl := defaultImpl(PlayerLimit)

	// 1. if player not exists, should return err
	err := impl.CancelMatch(UID)
	assert.Equal(t, merr.ErrPlayerNotExists, err)

	p, g := createTempGroup(UID, impl, t)
	err = impl.StartMatch(p.UID())
	assert.Nil(t, err)
	assert.Equal(t, entry.GroupStateMatch, g.Base().GetState())

	// 2. if group not exists, should return err
	impl.groupMgr.Delete(g.ID()) // delete temp
	err = impl.CancelMatch(UID)
	assert.Equal(t, merr.ErrGroupNotExists, err)
	impl.groupMgr.Add(g.ID(), g) // add  back

	// 4. if group state is not `match`, should return err
	g.Base().SetState(entry.GroupStateInvite)
	err = impl.CancelMatch(UID)
	assert.Equal(t, merr.ErrGroupInInvite, err)
	g.Base().SetState(entry.GroupStateGame)
	err = impl.CancelMatch(UID)
	assert.Equal(t, merr.ErrGroupInGame, err)
	g.Base().SetState(entry.GroupStateDissolved)
	err = impl.CancelMatch(UID)
	assert.Equal(t, merr.ErrGroupDissolved, err)
	g.Base().SetState(entry.GroupStateMatch) // set back

	// 5. return success and group state changes to `invite`
	err = impl.CancelMatch(UID)
	assert.Nil(t, err)
	assert.Equal(t, entry.GroupStateInvite, g.Base().GetState())
}
