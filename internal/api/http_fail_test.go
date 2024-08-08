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

func init() {
	gin.SetMode(gin.ReleaseMode)
}

func TestAPI_StartMatch_StateNotMatching(t *testing.T) {
	api := NewAPI(2, time.Second)
	router := api.setupRouter()

	_ = requestCreateGroup(router, "uid", t)

	req, _ := http.NewRequest("POST", "/match/cancel_match/uid", http.NoBody)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, merr.ErrGroupInInvite.Error(), assertRspNotOk(w, t))
}

func TestAPI_StartMatch_StateNotInvite(t *testing.T) {
	api := NewAPI(2, time.Second)
	router := api.setupRouter()

	g := requestCreateGroup(router, "uid", t)
	api.gm.Get(g.GroupID).Base().SetState(entry.GroupStateMatch)

	req, _ := http.NewRequest("POST", "/match/start_match/uid", http.NoBody)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, merr.ErrGroupInMatch.Error(), assertRspNotOk(w, t))
}

func TestAPI_SetVoice_WrongVoiceState(t *testing.T) {
	api := NewAPI(2, time.Second)
	router := api.setupRouter()

	_ = requestCreateGroup(router, "uid", t)

	req, _ := http.NewRequest("POST", "/match/set_voice_state",
		bytes.NewBuffer(createSetVoiceStateParam("uid", entry.PlayerVoiceState(100))))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAPI_SetVoice_PlayerNotInGroup(t *testing.T) {
	api := NewAPI(2, time.Second)
	router := api.setupRouter()

	req, _ := http.NewRequest("POST", "/match/set_voice_state",
		bytes.NewBuffer(createSetVoiceStateParam("uid", entry.PlayerVoiceStateUnmute)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, merr.ErrPlayerNotExists.Error(), assertRspNotOk(w, t))
}

func TestAPI_SetVoice_BadRequest(t *testing.T) {
	api := NewAPI(2, time.Second)
	router := api.setupRouter()

	req, _ := http.NewRequest("POST", "/match/set_voice_state",
		bytes.NewBuffer(createSetVoiceStateParam("", entry.PlayerVoiceState(10))))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAPI_SetRecentJoinGroup_BadRequest(t *testing.T) {
	api := NewAPI(2, time.Second)
	router := api.setupRouter()

	req, _ := http.NewRequest("POST", "/match/set_recent_join_group",
		bytes.NewBuffer(createSetRecentJoinParam("", true)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAPI_SetRecentJoinGroup_NotCaptain(t *testing.T) {
	api := NewAPI(2, time.Second)
	router := api.setupRouter()

	g := requestCreateGroup(router, "uid1", t)
	requestEnterGroup(router, "uid2", g.GroupID, t)

	req, _ := http.NewRequest("POST", "/match/set_recent_join_group",
		bytes.NewBuffer(createSetRecentJoinParam("uid2", true)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, merr.ErrPermissionDeny.Error(), assertRspNotOk(w, t))
}

func TestAPI_SetNearbyJoinGroup_BadRequest(t *testing.T) {
	api := NewAPI(2, time.Second)
	router := api.setupRouter()

	req, _ := http.NewRequest("POST", "/match/set_nearby_join_group",
		bytes.NewBuffer(createSetNearbyJoinParam("", true)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAPI_SetNearbyJoinGroup_NotCaptain(t *testing.T) {
	api := NewAPI(2, time.Second)
	router := api.setupRouter()

	g := requestCreateGroup(router, "uid1", t)
	requestEnterGroup(router, "uid2", g.GroupID, t)

	req, _ := http.NewRequest("POST", "/match/set_nearby_join_group",
		bytes.NewBuffer(createSetNearbyJoinParam("uid2", true)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, merr.ErrPermissionDeny.Error(), assertRspNotOk(w, t))
}

func TestAPI_RefuseInvite_BadRequest(t *testing.T) {
	api := NewAPI(2, time.Second)
	router := api.setupRouter()

	req, _ := http.NewRequest("POST", "/match/refuse_invite",
		bytes.NewBuffer(createRefuseInviteParam("uid1", "uid2", 0)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAPI_AcceptInvite_BadRequest(t *testing.T) {
	api := NewAPI(2, time.Second)
	router := api.setupRouter()

	req, _ := http.NewRequest("POST", "/match/accept_invite",
		bytes.NewBuffer(createAcceptInviteParam("uid1", "uid2", 0)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAPI_AcceptInvite_GroupDissolved(t *testing.T) {
	api := NewAPI(2, time.Second)
	router := api.setupRouter()

	req, _ := http.NewRequest("POST", "/match/accept_invite",
		bytes.NewBuffer(createAcceptInviteParam("uid1", "uid2", 1)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, merr.ErrGroupDissolved.Error(), assertRspNotOk(w, t))
}

func TestAPI_Invite_PlayerNotExists(t *testing.T) {
	api := NewAPI(2, time.Second)
	router := api.setupRouter()

	req, _ := http.NewRequest("POST", "/match/invite", bytes.NewBuffer(createInviteParam("uid1", "uid2")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, merr.ErrPlayerNotExists.Error(), assertRspNotOk(w, t))
}

func TestAPI_Invite_BadRequest(t *testing.T) {
	api := NewAPI(2, time.Second)
	router := api.setupRouter()

	req, _ := http.NewRequest("POST", "/match/invite", bytes.NewBuffer(createInviteParam("uid1", "")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAPI_ChangeRole_BadRequest(t *testing.T) {
	api := NewAPI(2, time.Second)
	router := api.setupRouter()

	g := requestCreateGroup(router, "uid1", t)
	requestEnterGroup(router, "uid2", g.GroupID, t)

	req, _ := http.NewRequest("POST", "/match/change_role",
		bytes.NewBuffer(createChangeRoleParam("uid2", "uid1", entry.GroupRole(0))))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAPI_ChangeRole_RoleNotExists(t *testing.T) {
	api := NewAPI(2, time.Second)
	router := api.setupRouter()

	g := requestCreateGroup(router, "uid1", t)
	requestEnterGroup(router, "uid2", g.GroupID, t)

	req, _ := http.NewRequest("POST", "/match/change_role",
		bytes.NewBuffer(createChangeRoleParam("uid2", "uid1", entry.GroupRole(127))))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, "unsupported role: 127", assertRspNotOk(w, t))
}

func TestAPI_ChangeRole_NotCaptain(t *testing.T) {
	api := NewAPI(2, time.Second)
	router := api.setupRouter()

	g := requestCreateGroup(router, "uid1", t)
	requestEnterGroup(router, "uid2", g.GroupID, t)

	req, _ := http.NewRequest("POST", "/match/change_role",
		bytes.NewBuffer(createChangeRoleParam("uid2", "uid1", entry.GroupRoleCaptain)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, merr.ErrNotCaptain.Error(), assertRspNotOk(w, t))
}

func TestAPI_KickPlayer_BadRequest(t *testing.T) {
	api := NewAPI(2, time.Second)
	router := api.setupRouter()
	req, _ := http.NewRequest("POST", "/match/kick_player", bytes.NewBuffer(createKickParam("uid1", "")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAPI_KickPlayer_NotCaptain(t *testing.T) {
	api := NewAPI(2, time.Second)
	router := api.setupRouter()

	g := requestCreateGroup(router, "uid1", t)
	requestEnterGroup(router, "uid2", g.GroupID, t)

	req, _ := http.NewRequest("POST", "/match/kick_player", bytes.NewBuffer(createKickParam("uid2", "uid1")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, merr.ErrOnlyCaptainCanKickPlayer.Error(), assertRspNotOk(w, t))
}

func TestAPI_DissolveGroup_NotCaptain(t *testing.T) {
	api := NewAPI(2, time.Second)
	router := api.setupRouter()

	g := requestCreateGroup(router, "uid1", t)
	requestEnterGroup(router, "uid2", g.GroupID, t)

	req, _ := http.NewRequest("POST", "/match/dissolve_group/uid2", http.NoBody)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, merr.ErrOnlyCaptainCanDissolveGroup.Error(), assertRspNotOk(w, t))
}

func TestAPI_ExitGroup_NotInGroup(t *testing.T) {
	api := NewAPI(1, time.Second)
	router := api.setupRouter()

	req, _ := http.NewRequest("POST", "/match/exit_group/uid", http.NoBody)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, merr.ErrPlayerNotExists.Error(), assertRspNotOk(w, t))
}

func TestAPI_EnterGroup_BadRequest(t *testing.T) {
	api := NewAPI(1, time.Second)
	router := api.setupRouter()

	req, _ := http.NewRequest("POST", "/match/enter_group", bytes.NewBuffer(createEnterGroupParam("a", 0)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAPI_EnterGroup_GroupNotExists(t *testing.T) {
	api := NewAPI(1, time.Second)
	router := api.setupRouter()

	req, _ := http.NewRequest("POST", "/match/enter_group", bytes.NewBuffer(createEnterGroupParam("a", 1)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, merr.ErrGroupDissolved.Error(), assertRspNotOk(w, t))
}

func TestAPI_CreateGroup_BadRequest(t *testing.T) {
	api := NewAPI(1, time.Second)
	router := api.setupRouter()

	req, _ := http.NewRequest("POST", "/match/create_group", bytes.NewBuffer(createGroupParamBad("a", 0)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAPI_CreateGroup_UnsupportedMode(t *testing.T) {
	api := NewAPI(1, time.Second)
	router := api.setupRouter()

	req, _ := http.NewRequest("POST", "/match/create_group",
		bytes.NewBuffer(createGroupParamBad("a", constant.GameMode(10010))))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, "unsupported game mode: 10010", assertRspNotOk(w, t))
}

func createGroupParamBad(uid string, mode constant.GameMode) []byte {
	param := &pto.CreateGroup{
		PlayerInfo: pto.PlayerInfo{
			UID:         uid,
			GameMode:    mode,
			ModeVersion: 1,
			Glicko2Info: &pto.Glicko2Info{},
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

func assertRspNotOk(w *httptest.ResponseRecorder, t *testing.T) string {
	assert.Equal(t, http.StatusOK, w.Code)
	rsp := response.NewHTTPResponse(w.Body.Bytes())
	assert.Equal(t, http.StatusOK, rsp.Code)
	return rsp.Message
}
