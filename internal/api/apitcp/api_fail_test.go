package apitcp

import (
	"io"
	"net"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"

	"github.com/hedon954/go-matcher/internal/config"
	"github.com/hedon954/go-matcher/internal/constant"
	"github.com/hedon954/go-matcher/internal/entry"
	"github.com/hedon954/go-matcher/internal/merr"
	"github.com/hedon954/go-matcher/pkg/algorithm/glicko2"
	"github.com/hedon954/go-matcher/pkg/zinx/zconfig"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
	log.Logger = zerolog.New(io.Discard)
}

func TestAPI_CancelMatch_StateNotMatching(t *testing.T) {
	_, client, _, shutdown := initServerClientAndCreateGroup("uid", 1, entry.GroupStateInvite, t)
	defer shutdown()

	_, errMsg := requestCancelMatch(client, "uid", t)
	assert.Equal(t, merr.ErrGroupInInvite.Error(), errMsg)
}

func TestAPI_StartMatch_StateNotInvite(t *testing.T) {
	_, client, _, shutdown := initServerClientAndCreateGroup("uid", 1, entry.GroupStateMatch, t)
	defer shutdown()

	errMsg := requestStartMatch(client, "uid", t)
	assert.Equal(t, merr.ErrGroupInMatch.Error(), errMsg)
}

func TestAPI_UploadPlayerAttr_LackOfUID(t *testing.T) {
	_, client, _, shutdown := initServerClientAndCreateGroup("uid", 1, entry.GroupStateInvite, t)
	defer shutdown()

	errMsg := requestUnloadPlayerAttr(client, "", true, true, t)
	assert.Equal(t, merr.ErrPlayerNotExists.Error(), errMsg)
}

func TestAPI_UploadPlayerAttr_AttrInvalid(t *testing.T) {
	_, client, _, shutdown := initServerClientAndCreateGroup("uid", 1, entry.GroupStateInvite, t)
	defer shutdown()

	errMsg := requestUnloadPlayerAttr(client, "uid", false, true, t)
	assert.Equal(t, "lack of basic attr", errMsg)
}

//nolint:dupl
func TestAPI_Unready_StateNotInvite(t *testing.T) {
	api, client, groupID, shutdown := initServerClientAndCreateGroup("uid", 1, entry.GroupStateMatch, t)
	defer shutdown()

	// in game can not unready
	api.GM.Get(groupID).Base().SetStateWithLock(entry.GroupStateGame)

	errMsg := requestUnready(client, "uid", t)
	assert.Equal(t, merr.ErrGroupInGame.Error(), errMsg)
}

//nolint:dupl
func TestAPI_Ready_StateNotInvite(t *testing.T) {
	api, client, groupID, shutdown := initServerClientAndCreateGroup("uid", 1, entry.GroupStateMatch, t)
	defer shutdown()

	// in match can not ready
	errMsg := requestReady(client, "uid", t)
	assert.Equal(t, merr.ErrGroupInMatch.Error(), errMsg)

	// in game can not ready
	api.GM.Get(groupID).Base().SetStateWithLock(entry.GroupStateGame)
	errMsg = requestReady(client, "uid", t)
	assert.Equal(t, merr.ErrGroupInGame.Error(), errMsg)
}

func TestAPI_SetVoice_WrongVoiceState(t *testing.T) {
	_, client, _, shutdown := initServerClientAndCreateGroup("uid", 1, entry.GroupStateInvite, t)
	defer shutdown()

	errMsg := requestSetVoiceState(client, "uid", entry.PlayerVoiceState(100), t)
	assert.Equal(t, "invalid voice state", errMsg)
}

func TestAPI_SetVoice_PlayerNotExists(t *testing.T) {
	_, client, _, shutdown := initServerClientAndCreateGroup("uid", 1, entry.GroupStateInvite, t)
	defer shutdown()

	errMsg := requestSetVoiceState(client, "uid2", entry.PlayerVoiceStateMute, t)
	assert.Equal(t, merr.ErrPlayerNotExists.Error(), errMsg)
}

func TestAPI_SetRecentJoinGroup_BadRequest(t *testing.T) {
	_, client, _, shutdown := initServerClientAndCreateGroup("uid", 1, entry.GroupStateInvite, t)
	defer shutdown()

	errMsg := requestSetRecentJoinGroup(client, "", true, t)
	assert.Equal(t, merr.ErrPlayerNotExists.Error(), errMsg)
}

func TestAPI_SetRecentJoinGroup_NotCaptain(t *testing.T) {
	_, client, groupID, shutdown := initServerClientAndCreateGroup("uid", 2, entry.GroupStateInvite, t)
	defer shutdown()

	errMsg := requestEnterGroup(client, "uid2", groupID, t)
	assert.Empty(t, errMsg)

	errMsg = requestSetRecentJoinGroup(client, "uid2", true, t)
	assert.Equal(t, merr.ErrPermissionDeny.Error(), errMsg)
}

func TestAPI_SetNearbyJoinGroup_BadRequest(t *testing.T) {
	_, client, _, shutdown := initServerClientAndCreateGroup("uid", 1, entry.GroupStateInvite, t)
	defer shutdown()

	errMsg := requestSetNearbyJoinGroup(client, "", true, t)
	assert.Equal(t, merr.ErrPlayerNotExists.Error(), errMsg)
}

func TestAPI_SetNearbyJoinGroup_NotCaptain(t *testing.T) {
	_, client, groupID, shutdown := initServerClientAndCreateGroup("uid", 2, entry.GroupStateInvite, t)
	defer shutdown()

	errMsg := requestEnterGroup(client, "uid2", groupID, t)
	assert.Empty(t, errMsg)

	errMsg = requestSetNearbyJoinGroup(client, "uid2", true, t)
	assert.Equal(t, merr.ErrPermissionDeny.Error(), errMsg)
}

func TestAPI_RefuseInvite_BadRequest(t *testing.T) {
	_, client, _, shutdown := initServerClientAndCreateGroup("uid", 2, entry.GroupStateInvite, t)
	defer shutdown()

	errMsg := requestRefuseInvite(client, "uid1", "uid2", 0, t)
	assert.Equal(t, "lack of group id", errMsg)
}

func TestAPI_AcceptInvite_BadRequest(t *testing.T) {
	_, client, _, shutdown := initServerClientAndCreateGroup("uid", 2, entry.GroupStateInvite, t)
	defer shutdown()

	errMsg := requestAcceptInvite(client, "uid1", "uid2", 0, t)
	assert.Equal(t, "lack of group id", errMsg)
}

func TestAPI_AcceptInvite_GroupDissolved(t *testing.T) {
	_, client, _, shutdown := initServerClientAndCreateGroup("uid", 2, entry.GroupStateInvite, t)
	defer shutdown()

	errMsg := requestAcceptInvite(client, "uid1", "uid2", 1000, t)
	assert.Equal(t, merr.ErrGroupDissolved.Error(), errMsg)
}

func TestAPI_Invite_PlayerNotExists(t *testing.T) {
	_, client, _, shutdown := initServerClientAndCreateGroup("uid", 2, entry.GroupStateInvite, t)
	defer shutdown()

	errMsg := requestInvite(client, "uid1", "uid2", t)
	assert.Equal(t, merr.ErrPlayerNotExists.Error(), errMsg)
}

func TestAPI_Invite_BadRequest(t *testing.T) {
	_, client, _, shutdown := initServerClientAndCreateGroup("uid", 2, entry.GroupStateInvite, t)
	defer shutdown()

	errMsg := requestInvite(client, "uid", "", t)
	assert.Equal(t, "lack of invitee uid", errMsg)
}

func TestAPI_ChangeRole_BadRequest(t *testing.T) {
	_, client, groupID, shutdown := initServerClientAndCreateGroup("uid", 2, entry.GroupStateInvite, t)
	defer shutdown()

	errMsg := requestEnterGroup(client, "uid2", groupID, t)
	assert.Empty(t, errMsg)

	errMsg = requestChangeRole(client, "uid2", "uid1", entry.GroupRole(0), t)
	assert.Equal(t, "unsupported role: 0", errMsg)
}

func TestAPI_ChangeRole_RoleNotExists(t *testing.T) {
	_, client, groupID, shutdown := initServerClientAndCreateGroup("uid", 2, entry.GroupStateInvite, t)
	defer shutdown()

	errMsg := requestEnterGroup(client, "uid2", groupID, t)
	assert.Empty(t, errMsg)

	errMsg = requestChangeRole(client, "uid2", "uid1", entry.GroupRole(127), t)
	assert.Equal(t, "unsupported role: 127", errMsg)
}

func TestAPI_ChangeRole_NotCaptain(t *testing.T) {
	_, client, groupID, shutdown := initServerClientAndCreateGroup("uid1", 2, entry.GroupStateInvite, t)
	defer shutdown()

	errMsg := requestEnterGroup(client, "uid2", groupID, t)
	assert.Empty(t, errMsg)

	errMsg = requestChangeRole(client, "uid2", "uid1", entry.GroupRoleCaptain, t)
	assert.Equal(t, merr.ErrNotCaptain.Error(), errMsg)
}

func TestAPI_KickPlayer_BadRequest(t *testing.T) {
	_, client, _, shutdown := initServerClientAndCreateGroup("uid1", 2, entry.GroupStateInvite, t)
	defer shutdown()

	errMsg := requestKick(client, "uid1", "", t)
	assert.Equal(t, "lack of kicked uid", errMsg)
}

func TestAPI_KickPlayer_NotCaptain(t *testing.T) {
	_, client, groupID, shutdown := initServerClientAndCreateGroup("uid1", 2, entry.GroupStateInvite, t)
	defer shutdown()
	errMsg := requestEnterGroup(client, "uid2", groupID, t)
	assert.Empty(t, errMsg)

	errMsg = requestKick(client, "uid2", "uid1", t)
	assert.Equal(t, merr.ErrOnlyCaptainCanKickPlayer.Error(), errMsg)
}

func TestAPI_DissolveGroup_NotCaptain(t *testing.T) {
	_, client, groupID, shutdown := initServerClientAndCreateGroup("uid1", 2, entry.GroupStateInvite, t)
	defer shutdown()
	errMsg := requestEnterGroup(client, "uid2", groupID, t)
	assert.Empty(t, errMsg)

	errMsg = requestDissolveGroup(client, "uid2", t)
	assert.Equal(t, merr.ErrOnlyCaptainCanDissolveGroup.Error(), errMsg)
}

func TestAPI_DissolveGroup_PlayerNotExists(t *testing.T) {
	_, client, _, shutdown := initServerClientAndCreateGroup("uid1", 2, entry.GroupStateInvite, t)
	defer shutdown()

	errMsg := requestDissolveGroup(client, "uid2", t)
	assert.Equal(t, merr.ErrPlayerNotExists.Error(), errMsg)
}

func TestAPI_ExitGroup_NotInGroup(t *testing.T) {
	_, client, _, shutdown := initServerClientAndCreateGroup("uid1", 2, entry.GroupStateInvite, t)
	defer shutdown()
	errMsg := requestExitGroup(client, "uid2", t)
	assert.Equal(t, merr.ErrPlayerNotExists.Error(), errMsg)
}

func TestAPI_EnterGroup_BadRequest(t *testing.T) {
	_, client, _, shutdown := initServerClientAndCreateGroup("uid1", 2, entry.GroupStateInvite, t)
	defer shutdown()
	errMsg := requestEnterGroup(client, "uid2", 0, t)
	assert.Equal(t, "lack of group id", errMsg)
}

func TestAPI_EnterGroup_GroupNotExists(t *testing.T) {
	_, client, _, shutdown := initServerClientAndCreateGroup("uid1", 2, entry.GroupStateInvite, t)
	defer shutdown()
	errMsg := requestEnterGroup(client, "uid2", 100, t)
	assert.Equal(t, merr.ErrGroupDissolved.Error(), errMsg)
}

func TestAPI_CreateGroup_BadRequest(t *testing.T) {
	_, client, shutdown := initServerClient(1)
	defer shutdown()

	_, errMsg := requestCreateGroupWithMode(client, "uid1", 0, t)
	assert.Equal(t, "lack of game mode", errMsg)

	_, errMsg = requestCreateGroupWithModeAndModeVersion(client, "uid1", constant.GameModeGoatGame, 0, t)
	assert.Equal(t, "lack of mode version", errMsg)
}

func TestAPI_CreateGroup_UnsupportedMode(t *testing.T) {
	_, client, shutdown := initServerClient(1)
	defer shutdown()

	_, errMsg := requestCreateGroupWithMode(client, "uid1", constant.GameMode(10000), t)
	assert.Equal(t, "unsupported game mode: 10000", errMsg)
}

func initServerClientAndCreateGroup(uid string, groupPlayerLimit int, state entry.GroupState,
	t *testing.T) (api *API, client net.Conn, p int64, shutdown func()) {
	api, client, shutdown = initServerClient(groupPlayerLimit)
	rsp, errMsg := requestCreateGroup(client, uid, t)
	assert.Equal(t, "", errMsg)
	assert.Equal(t, int64(1), rsp.GroupId)
	api.GM.Get(rsp.GroupId).Base().SetStateWithLock(state)
	api.PM.Get(uid).Base().SetMatchStrategyWithLock(constant.MatchStrategyGlicko2)
	return api, client, rsp.GroupId, shutdown
}

func initServerClient(groupPlayerLimit int) (*API, net.Conn, func()) {
	conf := *zconfig.DefaultConfig
	conf.TCPPort = int(port.Add(1))
	api, server, shutdown := SetupTCPServer(newConf(groupPlayerLimit), &conf)
	time.Sleep(3 * time.Millisecond)
	client := startClient(conf.TCPPort)
	return api, client, func() {
		shutdown()
		_ = client.Close()
		server.Stop()
	}
}

var port = atomic.Int64{}

func init() {
	port.Store(9000)
}

func newConf(groupPlayerLimit int) *config.Config {
	return &config.Config{
		GroupPlayerLimit: groupPlayerLimit,
		MatchIntervalMs:  10,
		Glicko2: &glicko2.QueueArgs{
			MatchTimeoutSec: 300,
			TeamPlayerLimit: groupPlayerLimit,
			RoomTeamLimit:   3,
		},
		DelayTimerType: config.DelayTimerTypeNative,
		DelayTimerConfig: &config.DelayTimerConfig{
			InviteTimeoutMs:    300000,
			MatchTimeoutMs:     60000,
			WaitAttrTimeoutMs:  1,
			ClearRoomTimeoutMs: 1800000,
		},
	}
}
