package apihttp

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/hedon954/go-matcher/docs"
	internalapi "github.com/hedon954/go-matcher/internal/api"
	"github.com/hedon954/go-matcher/internal/config"
	"github.com/hedon954/go-matcher/internal/middleware"
	"github.com/hedon954/go-matcher/internal/pto"
	"github.com/hedon954/go-matcher/pkg/response"
)

// @title           Match Service Swagger API
// @version         1.0
// @description     This is the open api doc for match service

// @host      :5050
// @BasePath  /

func SetupHTTPServer(
	sc config.Configer[config.ServerConfig],
	mc config.Configer[config.MatchConfig],
) {
	api, shutdown := internalapi.Start(sc, mc)
	defer shutdown()

	server := API{api}
	r := server.setupRouter()
	srv := &http.Server{
		Addr:              ":5050",
		Handler:           r.Handler(),
		ReadHeaderTimeout: time.Minute,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Error().Err(err).Msg("error occurs in http server")
	}
}

func (api *API) setupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.WithRequestAndTrace())

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
		mg.POST("/upload_player_attr", api.UploadPlayerAttr)
		mg.POST("/ready/:uid", api.Ready)
		mg.POST("/unready/:uid", api.Unready)
		mg.POST("/exit_game", api.ExitGame)
	}

	docs.SwaggerInfo.BasePath = "/"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	return r
}

type API struct {
	*internalapi.API
}

// CreateGroup godoc
// @Summary create a new group
// @Description create a new group based on the request
// @Tags match service
// @Accept json
// @Produce json
// @Param x-request-id header string false "Request ID"
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
	g, err := api.MS.CreateGroup(c, &req)
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
// @Param x-request-id header string false "Request ID"
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
	if err := api.MS.EnterGroup(c, &req.PlayerInfo, req.GroupID); err != nil {
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
// @Param x-request-id header string false "Request ID"
// @Param uid path string true "User ID"
// @Success 200 {object} string "ok"
// @Failure 200 {object} string "Concrete Error Msg"
// @Router /match/exit_group/{uid} [post]
func (api *API) ExitGroup(c *gin.Context) {
	uid := c.Param("uid")
	if err := api.MS.ExitGroup(c, uid); err != nil {
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
// @Param x-request-id header string false "Request ID"
// @Param uid path string true "User ID"
// @Success 200 {object} string "ok"
// @Failure 200 {object} string "Concrete Error Msg"
// @Router /match/dissolve_group/{uid} [post]
func (api *API) DissolveGroup(c *gin.Context) {
	uid := c.Param("uid")
	if err := api.MS.DissolveGroup(c, uid); err != nil {
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
// @Param x-request-id header string false "Request ID"
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
	if err := api.MS.KickPlayer(c, req.CaptainUID, req.KickedUID); err != nil {
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
// @Param x-request-id header string false "Request ID"
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
	if err := api.MS.ChangeRole(c, req.CaptainUID, req.TargetUID, req.Role); err != nil {
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
// @Param x-request-id header string false "Request ID"
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
	if err := api.MS.Invite(c, req.InviterUID, &req.InviteeInfo); err != nil {
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
// @Param x-request-id header string false "Request ID"
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
	if err := api.MS.AcceptInvite(c, req.InviterUID, req.InviteeInfo, req.GroupID); err != nil {
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
// @Param x-request-id header string false "Request ID"
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
	api.MS.RefuseInvite(c, req.InviterUID, req.InviteeUID, req.GroupID, req.RefuseMsg)
	response.GinSuccess(c, nil)
}

// SetNearbyJoinGroup godoc
// @Summary set nearby join group
// @Description set whether group can be entered from nearby players list
// @Tags match service
// @Accept json
// @Produce json
// @Param x-request-id header string false "Request ID"
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
	if err := api.MS.SetNearbyJoinGroup(c, req.CaptainUID, req.Allow); err != nil {
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
// @Param x-request-id header string false "Request ID"
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
	if err := api.MS.SetRecentJoinGroup(c, req.CaptainUID, req.Allow); err != nil {
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
// @Param x-request-id header string false "Request ID"
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
	if err := api.MS.SetVoiceState(c, req.UID, req.State); err != nil {
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
// @Param x-request-id header string false "Request ID"
// @Param uid path string true "player uid"
// @Success 200 {object} string "ok"
// @Failure 400 {object} string "Bad Request"
// @Failure 200 {object} string "Concrete Error Msg"
// @Router /match/start_match/{uid} [post]
func (api *API) StartMatch(c *gin.Context) {
	uid := c.Param("uid")
	if err := api.MS.StartMatch(c, uid); err != nil {
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
// @Param x-request-id header string false "Request ID"
// @Param uid path string true "player uid"
// @Success 200 {object} string "ok"
// @Failure 400 {object} string "Bad Request"
// @Failure 200 {object} string "Concrete Error Msg"
// @Router /match/cancel_match/{uid} [post]
func (api *API) CancelMatch(c *gin.Context) {
	uid := c.Param("uid")
	if err := api.MS.CancelMatch(c, uid); err != nil {
		response.GinError(c, err)
		return
	}
	response.GinSuccess(c, nil)
}

// UploadPlayerAttrReq godoc
// @Summary upload player attr
// @Description upload player attr
// @Tags match service
// @Accept json
// @Produce json
// @Param x-request-id header string false "Request ID"
// @Param UploadPlayerAttrReq body UploadPlayerAttrReq true "Upload Player Attr Request Body"
// @Success 200 {object} string "ok"
// @Failure 400 {object} string "Bad Request"
// @Failure 200 {object} string "Concrete Error Msg"
// @Router /match/upload_player_attr [post]
func (api *API) UploadPlayerAttr(c *gin.Context) {
	var req UploadPlayerAttrReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.GinParamError(c, err)
		return
	}
	if err := api.MS.UploadPlayerAttr(c, req.UID, &req.UploadPlayerAttr); err != nil {
		response.GinError(c, err)
		return
	}
	response.GinSuccess(c, nil)
}

// Ready godoc
// @Summary ready
// @Description ready
// @Tags match service
// @Accept json
// @Produce json
// @Param x-request-id header string false "Request ID"
// @Param uid path string true "player uid"
// @Success 200 {object} string "ok"
// @Failure 400 {object} string "Bad Request"
// @Failure 200 {object} string "Concrete Error Msg"
// @Router /match/ready/{uid} [post]
func (api *API) Ready(c *gin.Context) {
	uid := c.Param("uid")
	if err := api.MS.Ready(c, uid); err != nil {
		response.GinError(c, err)
		return
	}
	response.GinSuccess(c, nil)
}

// Unready godoc
// @Summary unready
// @Description unready
// @Tags match service
// @Accept json
// @Produce json
// @Param x-request-id header string false "Request ID"
// @Param uid path string true "player uid"
// @Success 200 {object} string "ok"
// @Failure 400 {object} string "Bad Request"
// @Failure 200 {object} string "Concrete Error Msg"
// @Router /match/unready/{uid} [post]
func (api *API) Unready(c *gin.Context) {
	uid := c.Param("uid")
	if err := api.MS.Unready(c, uid); err != nil {
		response.GinError(c, err)
		return
	}
	response.GinSuccess(c, nil)
}

// ExitGame godoc
// @Summary exit game
// @Description exit game
// @Tags match service
// @Accept json
// @Produce json
// @Param x-request-id header string false "Request ID"
// @Param ExitGameReq body ExitGameReq true "Exit Game Request Body"
// @Success 200 {object} string "ok"
// @Failure 400 {object} string "Bad Request"
// @Failure 200 {object} string "Concrete Error Msg"
// @Router /match/exit_game [post]
func (api *API) ExitGame(c *gin.Context) {
	var req ExitGameReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.GinParamError(c, err)
		return
	}
	if err := api.MS.ExitGame(c, req.UID, req.RoomID); err != nil {
		response.GinError(c, err)
		return
	}
	response.GinSuccess(c, nil)
}
