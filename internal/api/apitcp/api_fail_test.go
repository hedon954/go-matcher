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

	"github.com/hedon954/go-matcher/internal/constant"
	"github.com/hedon954/go-matcher/internal/entry"
	"github.com/hedon954/go-matcher/internal/merr"
	"github.com/hedon954/go-matcher/pkg/zinx/zconfig"
	"github.com/hedon954/go-matcher/pkg/zinx/ziface"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
	log.Logger = zerolog.New(io.Discard)
}

func TestAPI_CancelMatch_StateNotMatching(t *testing.T) {
	api, server, client, _ := initServerClientAndCreateGroup("uid", 1, entry.GroupStateInvite, t)
	defer server.Stop()
	defer api.m.Stop()
	defer func() { _ = client.Close() }()

	_, errMsg := requestCancelMatch(client, "uid", t)
	assert.Equal(t, merr.ErrGroupInInvite.Error(), errMsg)
}

func TestAPI_StartMatch_StateNotInvite(t *testing.T) {
	api, server, client, _ := initServerClientAndCreateGroup("uid", 1, entry.GroupStateMatch, t)
	defer server.Stop()
	defer api.m.Stop()
	defer func() { _ = client.Close() }()

	errMsg := requestStartMatch(client, "uid", t)
	assert.Equal(t, merr.ErrGroupInMatch.Error(), errMsg)
}

func TestAPI_UploadPlayerAttr_LackOfUID(t *testing.T) {
	api, server, client, _ := initServerClientAndCreateGroup("uid", 1, entry.GroupStateInvite, t)
	defer server.Stop()
	defer api.m.Stop()
	defer func() { _ = client.Close() }()

	errMsg := requestUnloadPlayerAttr(client, "", true, true, t)
	assert.Equal(t, merr.ErrPlayerNotExists.Error(), errMsg)
}

func TestAPI_UploadPlayerAttr_AttrInvalid(t *testing.T) {
	api, server, client, _ := initServerClientAndCreateGroup("uid", 1, entry.GroupStateInvite, t)
	defer server.Stop()
	defer api.m.Stop()
	defer func() { _ = client.Close() }()

	errMsg := requestUnloadPlayerAttr(client, "uid", false, true, t)
	assert.Equal(t, "lack of basic attr", errMsg)
}

func TestAPI_UploadPlayerAttr_Glicko2AttrInvalid(t *testing.T) {
	api, server, client, _ := initServerClientAndCreateGroup("uid", 1, entry.GroupStateInvite, t)
	defer server.Stop()
	defer api.m.Stop()
	defer func() { _ = client.Close() }()

	errMsg := requestUnloadPlayerAttr(client, "uid", true, false, t)
	assert.Equal(t, "invalid glicko2 attribute: protobuf data is empty", errMsg)
}

//nolint:dupl
func TestAPI_Unready_StateNotInvite(t *testing.T) {
	api, server, client, groupID := initServerClientAndCreateGroup("uid", 1, entry.GroupStateInvite, t)
	defer server.Stop()
	defer api.m.Stop()
	defer func() { _ = client.Close() }()

	// in game can not unready
	api.gm.Get(groupID).Base().SetStateWithLock(entry.GroupStateGame)

	errMsg := requestUnready(client, "uid", t)
	assert.Equal(t, merr.ErrGroupInGame.Error(), errMsg)
}

//nolint:dupl
func TestAPI_Ready_StateNotInvite(t *testing.T) {
	api, server, client, groupID := initServerClientAndCreateGroup("uid", 1, entry.GroupStateMatch, t)
	defer server.Stop()
	defer api.m.Stop()
	defer func() { _ = client.Close() }()

	// in match can not ready
	errMsg := requestReady(client, "uid", t)
	assert.Equal(t, merr.ErrGroupInMatch.Error(), errMsg)

	// in game can not ready
	api.gm.Get(groupID).Base().SetStateWithLock(entry.GroupStateGame)
	errMsg = requestReady(client, "uid", t)
	assert.Equal(t, merr.ErrGroupInGame.Error(), errMsg)
}

func TestAPI_SetVoice_WrongVoiceState(t *testing.T) {
	api, server, client, _ := initServerClientAndCreateGroup("uid", 1, entry.GroupStateInvite, t)
	defer server.Stop()
	defer api.m.Stop()
	defer func() { _ = client.Close() }()

	errMsg := requestSetVoiceState(client, "uid", entry.PlayerVoiceState(100), t)
	assert.Equal(t, "invalid voice state", errMsg)
}

func TestAPI_SetVoice_PlayerNotExists(t *testing.T) {
	api, server, client, _ := initServerClientAndCreateGroup("uid", 1, entry.GroupStateInvite, t)
	defer server.Stop()
	defer api.m.Stop()
	defer func() { _ = client.Close() }()

	errMsg := requestSetVoiceState(client, "uid2", entry.PlayerVoiceStateMute, t)
	assert.Equal(t, merr.ErrPlayerNotExists.Error(), errMsg)
}

func TestAPI_SetRecentJoinGroup_BadRequest(t *testing.T) {
	api, server, client, _ := initServerClientAndCreateGroup("uid", 1, entry.GroupStateInvite, t)
	defer server.Stop()
	defer api.m.Stop()
	defer func() { _ = client.Close() }()

	errMsg := requestSetRecentJoinGroup(client, "", true, t)
	assert.Equal(t, merr.ErrPlayerNotExists.Error(), errMsg)
}

func TestAPI_SetRecentJoinGroup_NotCaptain(t *testing.T) {
	api, server, client, groupID := initServerClientAndCreateGroup("uid", 2, entry.GroupStateInvite, t)
	defer server.Stop()
	defer api.m.Stop()
	defer func() { _ = client.Close() }()

	errMsg := requestEnterGroup(client, "uid2", groupID, t)
	assert.Empty(t, errMsg)

	errMsg = requestSetRecentJoinGroup(client, "uid2", true, t)
	assert.Equal(t, merr.ErrPermissionDeny.Error(), errMsg)
}

func TestAPI_SetNearbyJoinGroup_BadRequest(t *testing.T) {
	api, server, client, _ := initServerClientAndCreateGroup("uid", 1, entry.GroupStateInvite, t)
	defer server.Stop()
	defer api.m.Stop()
	defer func() { _ = client.Close() }()

	errMsg := requestSetNearbyJoinGroup(client, "", true, t)
	assert.Equal(t, merr.ErrPlayerNotExists.Error(), errMsg)
}

func TestAPI_SetNearbyJoinGroup_NotCaptain(t *testing.T) {
	api, server, client, groupID := initServerClientAndCreateGroup("uid", 2, entry.GroupStateInvite, t)
	defer server.Stop()
	defer api.m.Stop()
	defer func() { _ = client.Close() }()

	errMsg := requestEnterGroup(client, "uid2", groupID, t)
	assert.Empty(t, errMsg)

	errMsg = requestSetNearbyJoinGroup(client, "uid2", true, t)
	assert.Equal(t, merr.ErrPermissionDeny.Error(), errMsg)
}

func TestAPI_RefuseInvite_BadRequest(t *testing.T) {
	api, server, client, _ := initServerClientAndCreateGroup("uid", 2, entry.GroupStateInvite, t)
	defer server.Stop()
	defer api.m.Stop()
	defer func() { _ = client.Close() }()

	errMsg := requestRefuseInvite(client, "uid1", "uid2", 0, t)
	assert.Equal(t, "lack of group id", errMsg)
}

func TestAPI_AcceptInvite_BadRequest(t *testing.T) {
	api, server, client, _ := initServerClientAndCreateGroup("uid", 2, entry.GroupStateInvite, t)
	defer server.Stop()
	defer api.m.Stop()
	defer func() { _ = client.Close() }()

	errMsg := requestAcceptInvite(client, "uid1", "uid2", 0, t)
	assert.Equal(t, "lack of group id", errMsg)
}

func TestAPI_AcceptInvite_GroupDissolved(t *testing.T) {
	api, server, client, _ := initServerClientAndCreateGroup("uid", 2, entry.GroupStateInvite, t)
	defer server.Stop()
	defer api.m.Stop()
	defer func() { _ = client.Close() }()

	errMsg := requestAcceptInvite(client, "uid1", "uid2", 1000, t)
	assert.Equal(t, merr.ErrGroupDissolved.Error(), errMsg)
}

func TestAPI_Invite_PlayerNotExists(t *testing.T) {
	api, server, client, _ := initServerClientAndCreateGroup("uid", 2, entry.GroupStateInvite, t)
	defer server.Stop()
	defer api.m.Stop()
	defer func() { _ = client.Close() }()

	errMsg := requestInvite(client, "uid1", "uid2", t)
	assert.Equal(t, merr.ErrPlayerNotExists.Error(), errMsg)
}

func TestAPI_Invite_BadRequest(t *testing.T) {
	api, server, client, _ := initServerClientAndCreateGroup("uid", 2, entry.GroupStateInvite, t)
	defer server.Stop()
	defer api.m.Stop()
	defer func() { _ = client.Close() }()

	errMsg := requestInvite(client, "uid", "", t)
	assert.Equal(t, "lack of invitee uid", errMsg)
}

func TestAPI_ChangeRole_BadRequest(t *testing.T) {
	api, server, client, groupID := initServerClientAndCreateGroup("uid", 2, entry.GroupStateInvite, t)
	defer server.Stop()
	defer api.m.Stop()
	defer func() { _ = client.Close() }()

	errMsg := requestEnterGroup(client, "uid2", groupID, t)
	assert.Empty(t, errMsg)

	errMsg = requestChangeRole(client, "uid2", "uid1", entry.GroupRole(0), t)
	assert.Equal(t, "unsupported role: 0", errMsg)
}

func TestAPI_ChangeRole_RoleNotExists(t *testing.T) {
	api, server, client, groupID := initServerClientAndCreateGroup("uid", 2, entry.GroupStateInvite, t)
	defer server.Stop()
	defer api.m.Stop()
	defer func() { _ = client.Close() }()

	errMsg := requestEnterGroup(client, "uid2", groupID, t)
	assert.Empty(t, errMsg)

	errMsg = requestChangeRole(client, "uid2", "uid1", entry.GroupRole(127), t)
	assert.Equal(t, "unsupported role: 127", errMsg)
}

func TestAPI_ChangeRole_NotCaptain(t *testing.T) {
	api, server, client, groupID := initServerClientAndCreateGroup("uid1", 2, entry.GroupStateInvite, t)
	defer server.Stop()
	defer api.m.Stop()
	defer func() { _ = client.Close() }()

	errMsg := requestEnterGroup(client, "uid2", groupID, t)
	assert.Empty(t, errMsg)

	errMsg = requestChangeRole(client, "uid2", "uid1", entry.GroupRoleCaptain, t)
	assert.Equal(t, merr.ErrNotCaptain.Error(), errMsg)
}

func TestAPI_KickPlayer_BadRequest(t *testing.T) {
	api, server, client, _ := initServerClientAndCreateGroup("uid1", 2, entry.GroupStateInvite, t)
	defer server.Stop()
	defer api.m.Stop()
	defer func() { _ = client.Close() }()

	errMsg := requestKick(client, "uid1", "", t)
	assert.Equal(t, "lack of kicked uid", errMsg)
}

func TestAPI_KickPlayer_NotCaptain(t *testing.T) {
	api, server, client, groupID := initServerClientAndCreateGroup("uid1", 2, entry.GroupStateInvite, t)
	defer server.Stop()
	defer api.m.Stop()
	defer func() { _ = client.Close() }()
	errMsg := requestEnterGroup(client, "uid2", groupID, t)
	assert.Empty(t, errMsg)

	errMsg = requestKick(client, "uid2", "uid1", t)
	assert.Equal(t, merr.ErrOnlyCaptainCanKickPlayer.Error(), errMsg)
}

func TestAPI_DissolveGroup_NotCaptain(t *testing.T) {
	api, server, client, groupID := initServerClientAndCreateGroup("uid1", 2, entry.GroupStateInvite, t)
	defer server.Stop()
	defer api.m.Stop()
	defer func() { _ = client.Close() }()
	errMsg := requestEnterGroup(client, "uid2", groupID, t)
	assert.Empty(t, errMsg)

	errMsg = requestDissolveGroup(client, "uid2", t)
	assert.Equal(t, merr.ErrOnlyCaptainCanDissolveGroup.Error(), errMsg)
}

func TestAPI_DissolveGroup_PlayerNotExists(t *testing.T) {
	api, server, client, _ := initServerClientAndCreateGroup("uid1", 2, entry.GroupStateInvite, t)
	defer server.Stop()
	defer api.m.Stop()
	defer func() { _ = client.Close() }()

	errMsg := requestDissolveGroup(client, "uid2", t)
	assert.Equal(t, merr.ErrPlayerNotExists.Error(), errMsg)
}

func TestAPI_ExitGroup_NotInGroup(t *testing.T) {
	api, server, client, _ := initServerClientAndCreateGroup("uid1", 2, entry.GroupStateInvite, t)
	defer server.Stop()
	defer api.m.Stop()
	defer func() { _ = client.Close() }()
	errMsg := requestExitGroup(client, "uid2", t)
	assert.Equal(t, merr.ErrPlayerNotExists.Error(), errMsg)
}

func TestAPI_EnterGroup_BadRequest(t *testing.T) {
	api, server, client, _ := initServerClientAndCreateGroup("uid1", 2, entry.GroupStateInvite, t)
	defer server.Stop()
	defer api.m.Stop()
	defer func() { _ = client.Close() }()
	errMsg := requestEnterGroup(client, "uid2", 0, t)
	assert.Equal(t, "lack of group id", errMsg)
}

func TestAPI_EnterGroup_GroupNotExists(t *testing.T) {
	api, server, client, _ := initServerClientAndCreateGroup("uid1", 2, entry.GroupStateInvite, t)
	defer server.Stop()
	defer api.m.Stop()
	defer func() { _ = client.Close() }()
	errMsg := requestEnterGroup(client, "uid2", 100, t)
	assert.Equal(t, merr.ErrGroupDissolved.Error(), errMsg)
}

func TestAPI_CreateGroup_BadRequest(t *testing.T) {
	api, server, client := initServerClient(1)
	defer server.Stop()
	defer api.m.Stop()
	defer func() { _ = client.Close() }()

	_, errMsg := requestCreateGroupWithMode(client, "uid1", 0, t)
	assert.Equal(t, "lack of game mode", errMsg)

	_, errMsg = requestCreateGroupWithModeAndModeVersion(client, "uid1", constant.GameModeGoatGame, 0, t)
	assert.Equal(t, "lack of mode version", errMsg)
}

func TestAPI_CreateGroup_UnsupportedMode(t *testing.T) {
	api, server, client := initServerClient(1)
	defer server.Stop()
	defer api.m.Stop()
	defer func() { _ = client.Close() }()

	_, errMsg := requestCreateGroupWithMode(client, "uid1", constant.GameMode(10000), t)
	assert.Equal(t, "unsupported game mode: 10000", errMsg)
}

func initServerClientAndCreateGroup(uid string, groupPlayerLimit int, state entry.GroupState,
	t *testing.T) (api *API, server ziface.IServer, client net.Conn, p int64) {
	api, server, client = initServerClient(groupPlayerLimit)
	rsp, errMsg := requestCreateGroup(client, uid, t)
	assert.Equal(t, "", errMsg)
	assert.Equal(t, int64(1), rsp.GroupId)
	api.gm.Get(rsp.GroupId).Base().SetStateWithLock(state)
	api.pm.Get(uid).Base().SetMatchStrategyWithLock(constant.MatchStrategyGlicko2)
	return api, server, client, rsp.GroupId
}

func initServerClient(groupPlayerLimit int) (*API, ziface.IServer, net.Conn) {
	conf := *zconfig.DefaultConfig
	conf.TCPPort = int(port.Add(1))
	api, server := SetupTCPServer(groupPlayerLimit, &conf)
	time.Sleep(3 * time.Millisecond)
	client := startClient(conf.TCPPort)
	return api, server, client
}

var port = atomic.Int64{}

func init() {
	port.Store(9000)
}
