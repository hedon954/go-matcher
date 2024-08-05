package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hedon954/go-matcher/internal/constant"
	"github.com/hedon954/go-matcher/internal/entry"
	"github.com/hedon954/go-matcher/internal/merr"
	"github.com/hedon954/go-matcher/internal/pto"
	"github.com/hedon954/go-matcher/pkg/response"
	"github.com/stretchr/testify/assert"
)

func Test_HTTP_ShouldWork(t *testing.T) {
	api := NewAPI(2, time.Millisecond*10)
	defer api.m.Stop()

	router := api.setupRouter()

	const (
		UIDA  = "a"
		UIDB  = "b"
		UIDBB = "bb"
		UIDC  = "c"
		UIDCC = "cc"
		UIDD  = "D"

		G1 int64 = 1
		G2 int64 = 2
		G3 int64 = 3
		G4 int64 = 4
		G5 int64 = 5
	)

	// 1. 'a' create a group 'g1'
	rsp := requestCreateGroup(router, UIDA, t)
	assert.Equal(t, int64(1), rsp.GroupID)
	g1 := api.gm.Get(rsp.GroupID)
	assert.NotNil(t, g1)
	assert.Equal(t, 1, len(g1.Base().GetPlayers()))
	assert.Equal(t, G1, g1.ID())
	requestEnterGroup(router, UIDA, rsp.GroupID, t) // re enter, no influence
	assert.Equal(t, 1, len(g1.Base().GetPlayers()))

	// 2. 'a' dissolve group 'g1'
	requestDissolveGroup(router, UIDA, t)
	assert.Nil(t, api.gm.Get(rsp.GroupID))

	// 3. 'a' create a group again 'g2'
	rsp = requestCreateGroup(router, UIDA, t)
	assert.Equal(t, int64(2), rsp.GroupID)
	g2 := api.gm.Get(rsp.GroupID)
	assert.NotNil(t, g2)
	assert.Equal(t, 1, len(g2.Base().GetPlayers()))
	assert.Equal(t, G2, g2.ID())
	ua := api.pm.Get(UIDA)
	assert.NotNil(t, ua)
	assert.Equal(t, entry.PlayerOnlineStateInGroup, getPlayerOnlineStateWithLock(ua))

	// 4. 'a' invite friend 'b'
	requestInvite(router, UIDA, UIDB, t)
	assert.Equal(t, 1, len(g2.Base().GetInviteRecords()))

	// 5. 'b' refuse invite
	requestRefuseInvite(router, UIDA, UIDB, g2.ID(), t)
	assert.Equal(t, 0, len(g2.Base().GetInviteRecords()))

	// 6. 'a' invite friend 'b' again
	requestInvite(router, UIDA, UIDB, t)
	assert.Equal(t, 1, len(g2.Base().GetInviteRecords()))

	// 7. 'b' accept invite
	requestAcceptInvite(router, UIDA, UIDB, g2.ID(), t)
	assert.Equal(t, 0, len(g2.Base().GetInviteRecords()))

	// 8. 'b' enter group
	requestEnterGroup(router, UIDB, g2.ID(), t)
	assert.Equal(t, 2, len(g2.Base().GetPlayers()))
	ub := api.pm.Get(UIDB)
	assert.Equal(t, entry.PlayerOnlineStateInGroup, getPlayerOnlineStateWithLock(ub))
	assert.Equal(t, g2.ID(), ub.Base().GroupID)

	// 9. 'a' change role to 'b'
	assert.Equal(t, ua, g2.GetCaptain())
	requestChangeRole(router, UIDA, UIDB, entry.GroupRoleCaptain, t)
	assert.Equal(t, ub, g2.GetCaptain())
	assert.Equal(t, 2, len(g2.Base().GetPlayers()))

	// 10. 'b' exit group 'g2'
	requestExitGroup(router, UIDB, t)
	assert.Equal(t, 1, len(g2.Base().GetPlayers()))
	assert.Equal(t, ua, g2.GetCaptain())
	assert.Nil(t, api.pm.Get(UIDB))

	// 11. 'a' invite friend 'b'
	requestInvite(router, UIDA, UIDB, t)
	assert.Equal(t, 1, len(g2.Base().GetInviteRecords()))

	// 12. 'b' accept invite
	requestAcceptInvite(router, UIDA, UIDB, g2.ID(), t)
	assert.Equal(t, 0, len(g2.Base().GetInviteRecords()))

	// 13. 'b' enter group
	requestEnterGroup(router, UIDB, g2.ID(), t)
	assert.Equal(t, 2, len(g2.Base().GetPlayers()))
	ub = api.pm.Get(UIDB)
	assert.NotNil(t, ub)
	assert.Equal(t, entry.PlayerOnlineStateInGroup, getPlayerOnlineStateWithLock(ub))

	// 14. 'a' kick player 'b'
	requestKick(router, UIDA, UIDB, t)
	assert.Equal(t, 1, len(g2.Base().GetPlayers()))
	assert.Nil(t, api.pm.Get(UIDB))

	// 15. 'a' set voice state
	requestSetVoiceState(router, UIDA, entry.PlayerVoiceStateUnmute, t)
	assert.Equal(t, entry.PlayerVoiceStateUnmute, ua.Base().VoiceState)

	// 16. 'a' set 'g2' allow nearby join
	assert.False(t, g2.Base().AllowNearbyJoin())
	requestSetNearbyJoinGroup(router, UIDA, true, t)
	assert.True(t, g2.Base().AllowNearbyJoin())

	// 17. 'a' set 'g2' allow recent join
	assert.False(t, g2.Base().AllowRecentJoin())
	requestSetRecentJoinGroup(router, UIDA, true, t)
	assert.True(t, g2.Base().AllowRecentJoin())

	// 18. 'a' start match 'g2'
	requestStartMatch(router, UIDA, t)
	time.Sleep(11 * time.Millisecond) // for at least get one round match
	assert.Equal(t, entry.GroupStateMatch, getGroupStateWithLock(g2))
	assert.Equal(t, entry.PlayerOnlineStateInMatch, getPlayerOnlineStateWithLock(ua))

	// 19. 'a' cancel match 'g2'
	requestCancelMatch(router, UIDA, t)
	time.Sleep(11 * time.Millisecond) // for at least get one round match
	assert.Equal(t, entry.GroupStateInvite, getGroupStateWithLock(g2))
	assert.Equal(t, entry.PlayerOnlineStateInGroup, getPlayerOnlineStateWithLock(ua))

	// 20. 'a' start match 'g2'
	requestStartMatch(router, UIDA, t)
	assert.Equal(t, entry.GroupStateMatch, getGroupStateWithLock(g2))
	assert.Equal(t, entry.PlayerOnlineStateInMatch, getPlayerOnlineStateWithLock(ua))

	// 21. 'b' create a full group 'g3'
	requestCreateFullGroup(router, UIDB, UIDBB, t)
	g3 := api.gm.Get(G3)
	assert.Equal(t, 2, len(g3.Base().GetPlayers()))

	// 22. 'c' create a full group 'g4'
	requestCreateFullGroup(router, UIDC, UIDCC, t)
	g4 := api.gm.Get(G4)
	assert.Equal(t, 2, len(g4.Base().GetPlayers()))

	// 23. 'd' create group 'g5'
	requestCreateGroup(router, UIDD, t)
	g5 := api.gm.Get(G5)
	assert.Equal(t, 1, len(g5.Base().GetPlayers()))

	// 24, start g3, g4 to match
	requestStartMatch(router, UIDB, t)
	requestStartMatch(router, UIDC, t)
	time.Sleep(20 * time.Millisecond)
	assert.Equal(t, entry.GroupStateMatch, getGroupStateWithLock(g3))
	assert.Equal(t, entry.GroupStateMatch, getGroupStateWithLock(g4))
	assert.Equal(t, entry.GroupStateMatch, getGroupStateWithLock(g2))

	// 25. start g5 to match, match success and create a new room
	requestStartMatch(router, UIDD, t)

	for i := 0; i <= 10; i++ {
		if api.m.Glicko2Matcher.RoomCount.Load() == 1 {
			assert.Equal(t, entry.GroupStateGame, getGroupStateWithLock(g2))
			assert.Equal(t, entry.GroupStateGame, getGroupStateWithLock(g3))
			assert.Equal(t, entry.GroupStateGame, getGroupStateWithLock(g4))
			assert.Equal(t, entry.GroupStateGame, getGroupStateWithLock(g5))
			assert.NotNil(t, api.rm.Get(1))
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func getGroupStateWithLock(g entry.Group) entry.GroupState {
	g.Base().Lock()
	defer g.Base().Unlock()
	return g.Base().GetState()
}

func getPlayerOnlineStateWithLock(p entry.Player) entry.PlayerOnlineState {
	p.Base().Lock()
	defer p.Base().Unlock()
	return p.Base().GetOnlineState()
}

func requestCreateFullGroup(router *gin.Engine, uid1, uid2 string, t *testing.T) {
	req, _ := http.NewRequest("POST", "/match/create_group", bytes.NewBuffer(createGroupParam(uid1)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assertRspOk(w, t)

	rsp := response.NewHTTPResponse(w.Body.Bytes())
	groupID := response.FromHTTPResponse[CreateGroupRsp](rsp).GroupID
	requestEnterGroup(router, uid2, groupID, t)
}

func requestCancelMatch(router *gin.Engine, uid string, t *testing.T) {
	req, _ := http.NewRequest("POST", "/match/cancel_match/"+uid, http.NoBody)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assertRspOk(w, t)
}

func requestStartMatch(router *gin.Engine, uid string, t *testing.T) {
	req, _ := http.NewRequest("POST", "/match/start_match/"+uid, http.NoBody)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assertRspOk(w, t)
}

func requestSetRecentJoinGroup(router *gin.Engine, uid string, allow bool, t *testing.T) {
	req, _ := http.NewRequest("POST", "/match/set_recent_join_group",
		bytes.NewBuffer(createSetRecentJoinParam(uid, allow)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assertRspOk(w, t)
}

func createSetRecentJoinParam(uid string, allow bool) []byte {
	param := &SetRecentJoinGroupReq{
		CaptainUID: uid,
		Allow:      allow,
	}
	bs, _ := json.Marshal(param)
	return bs
}

func requestSetNearbyJoinGroup(router *gin.Engine, uid string, allow bool, t *testing.T) {
	req, _ := http.NewRequest("POST", "/match/set_nearby_join_group",
		bytes.NewBuffer(createSetNearbyJoinParam(uid, allow)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assertRspOk(w, t)
}

func createSetNearbyJoinParam(uid string, allow bool) []byte {
	param := &SetNearbyJoinGroupReq{
		CaptainUID: uid,
		Allow:      allow,
	}
	bs, _ := json.Marshal(param)
	return bs
}

func requestSetVoiceState(router *gin.Engine, uid string, state entry.PlayerVoiceState, t *testing.T) {
	req, _ := http.NewRequest("POST", "/match/set_voice_state", bytes.NewBuffer(createSetVoiceStateParam(uid, state)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assertRspOk(w, t)
}

func createSetVoiceStateParam(uid string, state entry.PlayerVoiceState) []byte {
	param := &SetVoiceStateReq{
		UID:   uid,
		State: state,
	}
	bs, _ := json.Marshal(param)
	return bs
}

func requestKick(router *gin.Engine, captainUID, kickedUID string, t *testing.T) {
	req, _ := http.NewRequest("POST", "/match/kick_player",
		bytes.NewBuffer(createKickParam(captainUID, kickedUID)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assertRspOk(w, t)
}

func createKickParam(captainUID, kickedUID string) []byte {
	param := &KickPlayerReq{
		CaptainUID: captainUID,
		KickedUID:  kickedUID,
	}
	bs, _ := json.Marshal(param)
	return bs
}

func requestExitGroup(router *gin.Engine, uid string, t *testing.T) {
	req, _ := http.NewRequest("POST", "/match/exit_group/"+uid, http.NoBody)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assertRspOk(w, t)
}

func requestChangeRole(router *gin.Engine, captain, target string, role entry.GroupRole, t *testing.T) {
	req, _ := http.NewRequest("POST", "/match/change_role",
		bytes.NewBuffer(createChangeRoleParam(captain, target, role)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assertRspOk(w, t)
}

func createChangeRoleParam(captain, target string, role entry.GroupRole) []byte {
	param := &ChangeRoleReq{
		CaptainUID: captain,
		TargetUID:  target,
		Role:       role,
	}
	bs, _ := json.Marshal(param)
	return bs
}

func requestEnterGroup(router *gin.Engine, invitee string, groupID int64, t *testing.T) {
	req, _ := http.NewRequest("POST", "/match/enter_group", bytes.NewBuffer(createEnterGroupParam(invitee, groupID)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assertRspOk(w, t)
}

func createEnterGroupParam(invitee string, id int64) []byte {
	param := &EnterGroupReq{
		PlayerInfo: pto.EnterGroup{
			PlayerInfo: playerInfo(invitee),
		},
		GroupID: id,
	}
	bs, _ := json.Marshal(param)
	return bs
}

func requestAcceptInvite(router *gin.Engine, inviter, invitee string, groupID int64, t *testing.T) {
	req, _ := http.NewRequest("POST", "/match/accept_invite",
		bytes.NewBuffer(createAcceptInviteParam(inviter, invitee, groupID)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assertRspOk(w, t)
}

func createAcceptInviteParam(inviter, invitee string, groupID int64) []byte {
	info := playerInfo(invitee)
	param := &AcceptInviteReq{
		InviterUID:  inviter,
		InviteeInfo: &info,
		GroupID:     groupID,
	}
	bs, _ := json.Marshal(param)
	return bs
}

func requestRefuseInvite(router *gin.Engine, inviter, invitee string, groupID int64, t *testing.T) {
	req, _ := http.NewRequest("POST", "/match/refuse_invite",
		bytes.NewBuffer(createRefuseInviteParam(inviter, invitee, groupID)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assertRspOk(w, t)
}

func createRefuseInviteParam(inviter, invitee string, groupID int64) []byte {
	param := &RefuseInviteReq{
		InviterUID: inviter,
		InviteeUID: invitee,
		GroupID:    groupID,
		RefuseMsg:  "sorry",
	}
	bs, _ := json.Marshal(param)
	return bs
}

func requestInvite(router *gin.Engine, uid1, uid2 string, t *testing.T) {
	req, _ := http.NewRequest("POST", "/match/invite", bytes.NewBuffer(createInviteParam(uid1, uid2)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assertRspOk(w, t)
}

func createInviteParam(uid1, uid2 string) []byte {
	param := &InviteReq{
		InviterUID: uid1,
		InviteeUID: uid2,
	}
	bs, _ := json.Marshal(param)
	return bs
}

func requestDissolveGroup(router *gin.Engine, uid string, t *testing.T) {
	req, _ := http.NewRequest("POST", "/match/dissolve_group/"+uid, http.NoBody)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assertRspOk(w, t)
}

func requestCreateGroup(router *gin.Engine, uid string, t *testing.T) *CreateGroupRsp {
	req, _ := http.NewRequest("POST", "/match/create_group", bytes.NewBuffer(createGroupParam(uid)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	return creategGroupRsp(w)
}

func creategGroupRsp(w *httptest.ResponseRecorder) *CreateGroupRsp {
	var rsp response.HTTPResponse
	_ = json.Unmarshal(w.Body.Bytes(), &rsp)
	return response.FromHTTPResponse[CreateGroupRsp](&rsp)
}

func createGroupParam(uid string) []byte {
	param := &pto.CreateGroup{
		PlayerInfo: playerInfo(uid),
	}
	bs, _ := json.Marshal(param)
	return bs
}

func assertRspOk(w *httptest.ResponseRecorder, t *testing.T) {
	assert.Equal(t, http.StatusOK, w.Code)
	rsp := response.NewHTTPResponse(w.Body.Bytes())
	assert.Equal(t, http.StatusOK, rsp.Code)
	assert.Equal(t, "ok", rsp.Message)
}

func playerInfo(uid string) pto.PlayerInfo {
	return pto.PlayerInfo{
		UID:           uid,
		GameMode:      constant.GameModeGoatGame,
		ModeVersion:   1,
		MatchStrategy: constant.MatchStrategyGlicko2,
		Glicko2Info:   &pto.Glicko2Info{},
	}
}

func TestAPI_CreateGroup_BadRequest(t *testing.T) {
	api := NewAPI(1, time.Second)
	router := api.setupRouter()

	req, _ := http.NewRequest("POST", "/match/create_group", bytes.NewBuffer(createGroupParamBad("a")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func createGroupParamBad(uid string) []byte {
	param := &pto.CreateGroup{
		PlayerInfo: pto.PlayerInfo{
			UID:           uid,
			ModeVersion:   1,
			MatchStrategy: constant.MatchStrategyGlicko2,
			Glicko2Info:   &pto.Glicko2Info{},
		},
	}
	bs, _ := json.Marshal(param)
	return bs
}

func TestAPI_DissolveGroup_PlayerNotExists(t *testing.T) {
	api := NewAPI(1, time.Second)
	router := api.setupRouter()

	req, _ := http.NewRequest("POST", "/match/dissolve_group/uid", http.NoBody)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	rsp := response.NewHTTPResponse(w.Body.Bytes())
	assert.Equal(t, http.StatusOK, rsp.Code)
	assert.Equal(t, merr.ErrPlayerNotExists.Error(), rsp.Message)
}
