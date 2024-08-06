package api

import (
	"log"
	"time"

	"github.com/hedon954/go-matcher/docs"
	"github.com/hedon954/go-matcher/internal/config/mock"
	"github.com/hedon954/go-matcher/internal/entry"
	"github.com/hedon954/go-matcher/internal/matcher"
	"github.com/hedon954/go-matcher/internal/matcher/glicko2"
	"github.com/hedon954/go-matcher/internal/pto"
	"github.com/hedon954/go-matcher/internal/repository"
	"github.com/hedon954/go-matcher/internal/service"
	"github.com/hedon954/go-matcher/internal/service/impl"
	"github.com/hedon954/go-matcher/pkg/response"
	"github.com/hedon954/go-matcher/pkg/safe"

	"github.com/gin-gonic/gin"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Match Service Swagger API
// @version         1.0
// @description     This is the open api doc for match sergvice

// @host      :5050
// @BasePath  /

func SetupHTTPServer() {
	api := NewAPI(1, time.Second)
	r := api.setupRouter()
	if err := r.Run(":5050"); err != nil {
		log.Fatal(err)
	}
}

func (api *API) setupRouter() *gin.Engine {
	r := gin.Default()
	mg := r.Group("/match")
	{
		mg.POST("/create_group", api.CreateGroup)
		mg.POST("/enter_group", api.EnterGroup)
		mg.POST("/exit_group/:uid", api.ExitGroup)
		mg.POST("/dissolve_group/:uid", api.DissolveGroup)
		mg.POST("/kick_player", api.KickPlayer)
		mg.POST("/change_role", api.ChangeRole)
		mg.POST("/invite", api.Invite)
		mg.POST("/accept_invite", api.AcceptInvite)
		mg.POST("/refuse_invite", api.RefuseInvite)
		mg.POST("/set_nearby_join_group", api.SetNearbyJoinGroup)
		mg.POST("/set_recent_join_group", api.SetRecentJoinGroup)
		mg.POST("/set_voice_state", api.SetVoiceState)
		mg.POST("/start_match/:uid", api.StartMatch)
		mg.POST("/cancel_match/:uid", api.CancelMatch)
	}

	// sg := r.Group("/stat")
	// {
	// 	sg.POST("/group/:group_id", api.StatGroup)
	// }

	docs.SwaggerInfo.BasePath = "/"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	return r
}

type API struct {
	ms service.Service
	m  *matcher.Matcher
	pm *repository.PlayerMgr
	gm *repository.GroupMgr
	tm *repository.TeamMgr
	rm *repository.RoomMgr
}

func NewAPI(groupPlayerLimit int, matchInterval time.Duration) *API {
	var (
		groupChannel = make(chan entry.Group, 1024)
		roomChannel  = make(chan entry.Room, 1024)
	)

	var (
		playerMgr = repository.NewPlayerMgr()
		groupMgr  = repository.NewGroupMgr(0)
		teamMgr   = repository.NewTeamMgr(0)
		roomMgr   = repository.NewRoomMgr(0)
		configer  = &glicko2.Configer{
			Glicko2: new(mock.Glicko2Mock),
		}
		glicko2Matcher = glicko2.New(roomChannel, configer, matchInterval, playerMgr, groupMgr, teamMgr, roomMgr)
	)

	api := &API{
		ms: impl.NewDefault(groupPlayerLimit, playerMgr, groupMgr, groupChannel, roomChannel),
		m:  matcher.New(groupChannel, glicko2Matcher),
		pm: playerMgr,
		gm: groupMgr,
		tm: teamMgr,
		rm: roomMgr,
	}

	safe.Go(api.m.Start)
	return api
}

// CreateGroup godoc
// @Summary create a new group
// @Description create a new group based on the request
// @Tags match service
// @Accept json
// @Produce json
// @Param CreateGroup body pto.CreateGroup true "Create Group Request Body"
// @Success 200 {object} CreateGroupRsp
// @Failure 200 {object} string
// @Router /match/create_group [post]
func (api *API) CreateGroup(c *gin.Context) {
	var req pto.CreateGroup
	if err := c.ShouldBindJSON(&req); err != nil {
		response.GinParamError(c, err)
		return
	}
	g, err := api.ms.CreateGroup(&req)
	if err != nil {
		response.GinError(c, err)
		return
	}
	response.GinSuccess(c, CreateGroupRsp{GroupID: g.ID()})
}

// EnterGroup godoc
// @Summary enter a group
// @Description enter a group based on the request
// @Tags match service
// @Accept json
// @Produce json
// @Param EnterGroupReq body EnterGroupReq true "Enter Group Request Body"
// @Success 200 {object} string "ok"
// @Failure 400 {object} string "Bad Request"
// @Failure 200 {object} string "Concrete Error Msg"
// @Router /match/enter_group [post]
func (api *API) EnterGroup(c *gin.Context) {
	var req EnterGroupReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.GinParamError(c, err)
		return
	}
	if err := api.ms.EnterGroup(&req.PlayerInfo, req.GroupID); err != nil {
		response.GinError(c, err)
		return
	}
	response.GinSuccess(c, nil)
}

// ExitGroup godoc
// @Summary exit a group
// @Description exit a group based on the request
// @Tags match service
// @Accept json
// @Produce json
// @Param uid path string true "User ID"
// @Success 200 {object} string "ok"
// @Failure 200 {object} string "Concrete Error Msg"
// @Router /match/exit_group/{uid} [post]
func (api *API) ExitGroup(c *gin.Context) {
	uid := c.Param("uid")
	if err := api.ms.ExitGroup(uid); err != nil {
		response.GinError(c, err)
		return
	}
	response.GinSuccess(c, nil)
}

// DissolveGroup godoc
// @Summary dissolve a group
// @Description dissolve a group based on the request
// @Tags match service
// @Accept json
// @Produce json
// @Param uid path string true "User ID"
// @Success 200 {object} string "ok"
// @Failure 200 {object} string "Concrete Error Msg"
// @Router /match/dissolve_group/{uid} [post]
func (api *API) DissolveGroup(c *gin.Context) {
	uid := c.Param("uid")
	if err := api.ms.DissolveGroup(uid); err != nil {
		response.GinError(c, err)
		return
	}
	response.GinSuccess(c, nil)
}

// KickPlayer godoc
// @Summary kick a player
// @Description kick a player based on the request
// @Tags match service
// @Accept json
// @Produce json
// @Param KickPlayerReq body KickPlayerReq true "Kick Player Request Body"
// @Success 200 {object} string "ok"
// @Failure 400 {object} string "Bad Request"
// @Failure 200 {object} string "Concrete Error Msg"
// @Router /match/kick_player [post]
func (api *API) KickPlayer(c *gin.Context) {
	var req KickPlayerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.GinParamError(c, err)
		return
	}
	if err := api.ms.KickPlayer(req.CaptainUID, req.KickedUID); err != nil {
		response.GinError(c, err)
		return
	}
	response.GinSuccess(c, nil)
}

// ChangeRole godoc
// @Summary change a player's role
// @Description change a player's role based on the request
// @Tags match service
// @Accept json
// @Produce json
// @Param ChangeRoleReq body ChangeRoleReq true "Change Role Request Body"
// @Success 200 {object} string "ok"
// @Failure 400 {object} string "Bad Request"
// @Failure 200 {object} string "Concrete Error Msg"
// @Router /match/change_role [post]
func (api *API) ChangeRole(c *gin.Context) {
	var req ChangeRoleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.GinParamError(c, err)
		return
	}
	if err := api.ms.ChangeRole(req.CaptainUID, req.TargetUID, req.Role); err != nil {
		response.GinError(c, err)
		return
	}
	response.GinSuccess(c, nil)
}

// Invite godoc
// @Summary invite a player
// @Description invite a player based on the request
// @Tags match service
// @Accept json
// @Produce json
// @Param InviteReq body InviteReq true "Invite Request Body"
// @Success 200 {object} string "ok"
// @Failure 400 {object} string "Bad Request"
// @Failure 200 {object} string "Concrete Error Msg"
// @Router /match/invite [post]
func (api *API) Invite(c *gin.Context) {
	var req InviteReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.GinParamError(c, err)
		return
	}
	if err := api.ms.Invite(req.InviterUID, req.InviteeUID); err != nil {
		response.GinError(c, err)
		return
	}
	response.GinSuccess(c, nil)
}

// AcceptInvite godoc
// @Summary accept an invitation
// @Description accept an invitation based on the request
// @Tags match service
// @Accept json
// @Produce json
// @Param AcceptInviteReq body AcceptInviteReq true "Accept Invite Request Body"
// @Success 200 {object} string "ok"
// @Failure 400 {object} string "Bad Request"
// @Failure 200 {object} string "Concrete Error Msg"
// @Router /match/accept_invite [post]
func (api *API) AcceptInvite(c *gin.Context) {
	var req AcceptInviteReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.GinParamError(c, err)
		return
	}
	if err := api.ms.AcceptInvite(req.InviterUID, req.InviteeInfo, req.GroupID); err != nil {
		response.GinError(c, err)
		return
	}
	response.GinSuccess(c, nil)
}

// RefuseInvite godoc
// @Summary refuse an invitation
// @Description refuse an invitation based on the request
// @Tags match service
// @Accept json
// @Produce json
// @Param RefuseInviteReq body RefuseInviteReq true "Refuse Invite Request Body"
// @Success 200 {object} string "ok"
// @Failure 400 {object} string "Bad Request"
// @Failure 200 {object} string "Concrete Error Msg"
// @Router /match/refuse_invite [post]
func (api *API) RefuseInvite(c *gin.Context) {
	var req RefuseInviteReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.GinParamError(c, err)
		return
	}
	api.ms.RefuseInvite(req.InviterUID, req.InviteeUID, req.GroupID, req.RefuseMsg)
	response.GinSuccess(c, nil)
}

// SetNearbyJoinGroup godoc
// @Summary set nearby join group
// @Description set whether group can be entered from nearby players list
// @Tags match service
// @Accept json
// @Produce json
// @Param SetNearbyJoinGroupReq body SetNearbyJoinGroupReq true "Set Nearby Join Group Request Body"
// @Success 200 {object} string "ok"
// @Failure 400 {object} string "Bad Request"
// @Failure 200 {object} string "Concrete Error Msg"
// @Router /match/set_nearby_join_group [post]
func (api *API) SetNearbyJoinGroup(c *gin.Context) {
	var req SetNearbyJoinGroupReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.GinParamError(c, err)
		return
	}
	if err := api.ms.SetNearbyJoinGroup(req.CaptainUID, req.Allow); err != nil {
		response.GinError(c, err)
		return
	}
	response.GinSuccess(c, nil)
}

// SetRecentJoinGroup godoc
// @Summary set recent join group
// @Description set whether group can be entered from recent players list
// @Tags match service
// @Accept json
// @Produce json
// @Param SetRecentJoinGroupReq body SetRecentJoinGroupReq true "Set Recent Join Group Request Body"
// @Success 200 {object} string "ok"
// @Failure 400 {object} string "Bad Request"
// @Failure 200 {object} string "Concrete Error Msg"
// @Router /match/set_recent_join_group [post]
func (api *API) SetRecentJoinGroup(c *gin.Context) {
	var req SetRecentJoinGroupReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.GinParamError(c, err)
		return
	}
	if err := api.ms.SetRecentJoinGroup(req.CaptainUID, req.Allow); err != nil {
		response.GinError(c, err)
		return
	}
	response.GinSuccess(c, nil)
}

// SetVoiceState godoc
// @Summary set voice state
// @Description set player voice state
// @Tags match service
// @Accept json
// @Produce json
// @Param SetVoiceStateReq body SetVoiceStateReq true "Set Voice State Request Body"
// @Success 200 {object} string "ok"
// @Failure 400 {object} string "Bad Request"
// @Failure 200 {object} string "Concrete Error Msg"
// @Router /match/set_voice_state [post]
func (api *API) SetVoiceState(c *gin.Context) {
	var req SetVoiceStateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.GinParamError(c, err)
		return
	}
	if err := api.ms.SetVoiceState(req.UID, req.State); err != nil {
		response.GinError(c, err)
		return
	}
	response.GinSuccess(c, nil)
}

// StartMatch godoc
// @Summary start match
// @Description start to match
// @Tags match service
// @Accept json
// @Produce json
// @Param uid path string true "player uid"
// @Success 200 {object} string "ok"
// @Failure 400 {object} string "Bad Request"
// @Failure 200 {object} string "Concrete Error Msg"
// @Router /match/start_match/{uid} [post]
func (api *API) StartMatch(c *gin.Context) {
	uid := c.Param("uid")
	if err := api.ms.StartMatch(uid); err != nil {
		response.GinError(c, err)
		return
	}
	response.GinSuccess(c, nil)
}

// CancelMatch godoc
// @Summary cancel match
// @Description cancel match
// @Tags match service
// @Accept json
// @Produce json
// @Param uid path string true "player uid"
// @Success 200 {object} string "ok"
// @Failure 400 {object} string "Bad Request"
// @Failure 200 {object} string "Concrete Error Msg"
// @Router /match/cancel_match/{uid} [post]
func (api *API) CancelMatch(c *gin.Context) {
	uid := c.Param("uid")
	if err := api.ms.CancelMatch(uid); err != nil {
		response.GinError(c, err)
		return
	}
	response.GinSuccess(c, nil)
}
