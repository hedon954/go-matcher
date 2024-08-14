package apitcp

import (
	"time"

	"github.com/hedon954/go-matcher/internal/pb"
	"github.com/hedon954/go-matcher/pkg/zinx/zconfig"
	"github.com/hedon954/go-matcher/pkg/zinx/ziface"
	"github.com/hedon954/go-matcher/pkg/zinx/znet"
)

func SetupTCPServer(conf *zconfig.ZConfig) (*API, ziface.IServer) {
	server := znet.NewServer(conf)

	api := NewAPI(1, time.Second)
	api.setupRouter(server)

	go server.Serve()
	return api, server
}

func (api *API) setupRouter(s ziface.IServer) {
	s.AddRouter(uint32(pb.ReqType_REQ_TYPE_CREATE_GROUP), api.CreateGroup)
	s.AddRouter(uint32(pb.ReqType_REQ_TYPE_ENTER_GROUP), api.EnterGroup)
	s.AddRouter(uint32(pb.ReqType_REQ_TYPE_EXIT_GROUP), api.ExitGroup)
	s.AddRouter(uint32(pb.ReqType_REQ_TYPE_DISSOLVE_GROUP), api.DissolveGroup)
	s.AddRouter(uint32(pb.ReqType_REQ_TYPE_KICK_PLAYER), api.KickPlayer)
	s.AddRouter(uint32(pb.ReqType_REQ_TYPE_CHANGE_ROLE), api.ChangeRole)
	s.AddRouter(uint32(pb.ReqType_REQ_TYPE_INVITE), api.Invite)
	s.AddRouter(uint32(pb.ReqType_REQ_TYPE_ACCEPT_INVITE), api.AcceptInvite)
	s.AddRouter(uint32(pb.ReqType_REQ_TYPE_REFUSE_INVITE), api.RefuseInvite)
	s.AddRouter(uint32(pb.ReqType_REQ_TYPE_SET_NEARBY_JOIN_GROUP), api.SetNearbyJoinGroup)
	s.AddRouter(uint32(pb.ReqType_REQ_TYPE_SET_RECENT_JOIN_GROUP), api.SetRecentJoinGroup)
	s.AddRouter(uint32(pb.ReqType_REQ_TYPE_SET_VOICE_STATE), api.SetVoiceState)
	s.AddRouter(uint32(pb.ReqType_REQ_TYPE_START_MATCH), api.StartMatch)
	s.AddRouter(uint32(pb.ReqType_REQ_TYPE_CANCEL_MATCH), api.CancelMatch)
	s.AddRouter(uint32(pb.ReqType_REQ_TYPE_READY), api.Ready)
	s.AddRouter(uint32(pb.ReqType_REQ_TYPE_UNREADY), api.Unready)
	s.AddRouter(uint32(pb.ReqType_REQ_TYPE_UPLOAD_PLAYER_ATTR), api.UploadPlayerAttr)
	s.AddRouter(uint32(pb.ReqType_REQ_TYPE_EXIT_GAME), api.ExitGame)
}
