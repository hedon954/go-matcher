package apitcp

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hedon954/go-matcher/internal/config/mock"
	"github.com/hedon954/go-matcher/internal/constant"
	"github.com/hedon954/go-matcher/internal/entry"
	"github.com/hedon954/go-matcher/internal/log"
	"github.com/hedon954/go-matcher/internal/matcher"
	"github.com/hedon954/go-matcher/internal/matcher/glicko2"
	"github.com/hedon954/go-matcher/internal/pb"
	"github.com/hedon954/go-matcher/internal/pto"
	"github.com/hedon954/go-matcher/internal/repository"
	"github.com/hedon954/go-matcher/internal/service"
	"github.com/hedon954/go-matcher/internal/service/matchimpl"
	"github.com/hedon954/go-matcher/pkg/typeconv"
	"github.com/hedon954/go-matcher/pkg/zinx/ziface"

	timermock "github.com/hedon954/go-matcher/pkg/timer/mock"

	"google.golang.org/protobuf/proto"
)

type API struct {
	ms service.Match
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
			Glicko2: new(mock.Glicko2Mock), // TODO: change
		}
		glicko2Matcher = glicko2.New(roomChannel, configer, matchInterval, playerMgr, groupMgr, teamMgr, roomMgr)
	)

	delayTimer := timermock.NewTimer() // TODO: get from param

	api := &API{
		pm: playerMgr,
		gm: groupMgr,
		tm: teamMgr,
		rm: roomMgr,
		m:  matcher.New(groupChannel, glicko2Matcher),
		ms: matchimpl.NewDefault(groupPlayerLimit, playerMgr, groupMgr, teamMgr, roomMgr, groupChannel, roomChannel, delayTimer),
	}

	go delayTimer.Start()
	go api.m.Start()
	return api
}

// TODO: generate a reqType -> msgStruct map

func (api *API) CreateGroup(request ziface.IRequest) {
	data := typeconv.MustFromProto[pb.CreateGroupReq](request.GetData())

	if err := checkPlayerInfo(data.PlayerInfo); err != nil {
		api.responseParamError(request, err)
		return
	}

	param := &pto.CreateGroup{
		PlayerInfo: playerInfoFromPBToPTO(data.PlayerInfo),
	}

	group, err := api.ms.CreateGroup(context.Background(), param)
	if err != nil {
		api.responseError(request, err)
		return
	}

	api.responseSuccess(request, &pb.CreateGroupRsp{GroupId: group.ID()})
	fmt.Println("data: ", data)
}

func (api *API) EnterGroup(request ziface.IRequest) {
	data := typeconv.MustFromProto[pb.EnterGroupReq](request.GetData())

	if err := checkPlayerInfo(data.PlayerInfo); err != nil {
		api.responseParamError(request, err)
		return
	}
	if data.GroupId == 0 {
		api.responseParamError(request, errors.New("lack of group id"))
		return
	}

	param := &pto.EnterGroup{
		PlayerInfo: playerInfoFromPBToPTO(data.PlayerInfo),
		Source:     pto.EnterGroupSourceType(data.Source),
	}

	if err := api.ms.EnterGroup(context.Background(), param, data.GroupId); err != nil {
		api.responseError(request, err)
		return
	}

	api.responseSuccess(request, &pb.EnterGroupRsp{})
}

func (api *API) ExitGroup(request ziface.IRequest) {
	param := typeconv.MustFromProto[pb.ExitGroupReq](request.GetData())

	if err := api.ms.ExitGroup(context.Background(), param.Uid); err != nil {
		api.responseError(request, err)
		return
	}

	api.responseSuccess(request, &pb.ExitGroupRsp{})
}

func (api *API) DissolveGroup(request ziface.IRequest) {
	param := typeconv.MustFromProto[pb.DissolveGroupReq](request.GetData())

	if err := api.ms.DissolveGroup(context.Background(), param.Uid); err != nil {
		api.responseError(request, err)
		return
	}

	api.responseSuccess(request, &pb.DissolveGroupRsp{})
}

func (api *API) KickPlayer(request ziface.IRequest) {
	param := typeconv.MustFromProto[pb.KickPlayerReq](request.GetData())

	if err := api.ms.KickPlayer(context.Background(), param.CaptainUid, param.KickedUid); err != nil {
		api.responseError(request, err)
		return
	}

	api.responseSuccess(request, &pb.KickPlayerRsp{})
}

func (api *API) ChangeRole(request ziface.IRequest) {
	param := typeconv.MustFromProto[pb.ChangeRoleReq](request.GetData())

	if err := api.ms.ChangeRole(context.Background(),
		param.CaptainUid, param.TargetUid, entry.GroupRole(param.Role)); err != nil {
		api.responseError(request, err)
		return
	}

	api.responseSuccess(request, &pb.ChangeRoleRsp{})
}

func (api *API) Invite(request ziface.IRequest) {
	param := typeconv.MustFromProto[pb.InviteReq](request.GetData())

	if err := api.ms.Invite(context.Background(), param.InviterUid, param.InviteeUid); err != nil {
		api.responseError(request, err)
		return
	}

	api.responseSuccess(request, &pb.InviteRsp{})
}

func (api *API) AcceptInvite(request ziface.IRequest) {
	param := typeconv.MustFromProto[pb.AcceptInviteReq](request.GetData())

	if err := checkPlayerInfo(param.InviteeInfo); err != nil {
		api.responseParamError(request, err)
		return
	}
	if param.GroupId == 0 {
		api.responseParamError(request, errors.New("lack of group id"))
		return
	}

	inviteInfo := playerInfoFromPBToPTO(param.InviteeInfo)

	if err := api.ms.AcceptInvite(context.Background(), param.InviterUid, &inviteInfo, param.GroupId); err != nil {
		api.responseError(request, err)
		return
	}

	api.responseSuccess(request, &pb.AcceptInviteRsp{})
}

func (api *API) RefuseInvite(request ziface.IRequest) {
	param := typeconv.MustFromProto[pb.RefuseInviteReq](request.GetData())

	if param.GroupId == 0 {
		api.responseParamError(request, errors.New("lack of group id"))
		return
	}

	api.ms.RefuseInvite(context.Background(), param.InviterUid, param.InviteeUid, param.GroupId, param.RefuseMsg)

	api.responseSuccess(request, &pb.RefuseInviteRsp{})
}

func (api *API) SetNearbyJoinGroup(request ziface.IRequest) {
	param := typeconv.MustFromProto[pb.SetNearbyJoinGroupReq](request.GetData())

	if err := api.ms.SetNearbyJoinGroup(context.Background(), param.Uid, param.Allow); err != nil {
		api.responseError(request, err)
		return
	}

	api.responseSuccess(request, &pb.SetNearbyJoinGroupRsp{})
}

func (api *API) SetRecentJoinGroup(request ziface.IRequest) {
	param := typeconv.MustFromProto[pb.SetRecentJoinGroupReq](request.GetData())

	if err := api.ms.SetRecentJoinGroup(context.Background(), param.Uid, param.Allow); err != nil {
		api.responseError(request, err)
		return
	}

	api.responseSuccess(request, &pb.SetRecentJoinGroupRsp{})
}

func (api *API) SetVoiceState(request ziface.IRequest) {
	param := typeconv.MustFromProto[pb.SetVoiceStateReq](request.GetData())

	state := entry.PlayerVoiceState(param.State)
	if state != entry.PlayerVoiceStateUnmute && state != entry.PlayerVoiceStateMute {
		api.responseParamError(request, errors.New("invalid voice state"))
		return
	}

	if err := api.ms.SetVoiceState(context.Background(), param.Uid, state); err != nil {
		api.responseError(request, err)
		return
	}

	api.responseSuccess(request, &pb.SetVoiceStateRsp{})
}

func (api *API) StartMatch(request ziface.IRequest) {
	param := typeconv.MustFromProto[pb.StartMatchReq](request.GetData())

	if err := api.ms.StartMatch(context.Background(), param.Uid); err != nil {
		api.responseError(request, err)
		return
	}

	api.responseSuccess(request, &pb.StartMatchRsp{})
}

func (api *API) CancelMatch(request ziface.IRequest) {
	param := typeconv.MustFromProto[pb.CancelMatchReq](request.GetData())

	if err := api.ms.CancelMatch(context.Background(), param.Uid); err != nil {
		api.responseError(request, err)
		return
	}

	api.responseSuccess(request, &pb.CancelMatchRsp{})
}

func (api *API) Ready(request ziface.IRequest) {
	param := typeconv.MustFromProto[pb.ReadyReq](request.GetData())

	if err := api.ms.Ready(context.Background(), param.Uid); err != nil {
		api.responseError(request, err)
		return
	}

	api.responseSuccess(request, &pb.ReadyRsp{})
}

func (api *API) Unready(request ziface.IRequest) {
	param := typeconv.MustFromProto[pb.UnreadyReq](request.GetData())

	if err := api.ms.Unready(context.Background(), param.Uid); err != nil {
		api.responseError(request, err)
		return
	}

	api.responseSuccess(request, &pb.UnreadyRsp{})
}

func (api *API) createAndSendResponse(req ziface.IRequest, code pb.RspCode, err error) {
	rsp := &pb.CommonRsp{
		Code:      code,
		Message:   err.Error(),
		ReqType:   pb.ReqType(req.GetMsgID()),
		RequestId: "",
		TraceId:   "",
		Data:      nil,
	}
	bs, err := proto.Marshal(rsp)
	if err != nil {
		panic(err)
	}

	err = req.GetConnection().SendMsg(uint32(pb.ReqType_REQ_TYPE_MATCH_RESPONSE), bs)
	if err != nil {
		log.Error().
			Err(err).
			Any("rsp", rsp).
			Msg("send response error")
	}
}

func (api *API) responseParamError(req ziface.IRequest, err error) {
	api.createAndSendResponse(req, pb.RspCode_RSP_CODE_BAD_REQUEST, err)
}

func (api *API) responseError(req ziface.IRequest, err error) {
	api.createAndSendResponse(req, pb.RspCode_RSP_CODE_USER_ERROR, err)
}

func (api *API) responseSuccess(request ziface.IRequest, p proto.Message) {
	bs, err := proto.Marshal(p)
	if err != nil {
		panic(err)
	}
	rsp := &pb.CommonRsp{
		Code:      pb.RspCode_RSP_CODE_SUCCESS,
		Message:   "",
		ReqType:   pb.ReqType(request.GetMsgID()),
		RequestId: "",
		TraceId:   "",
		Data:      bs,
	}

	bs, err = proto.Marshal(rsp)
	if err != nil {
		panic(err)
	}

	err = request.GetConnection().SendMsg(uint32(pb.ReqType_REQ_TYPE_MATCH_RESPONSE), bs)
	if err != nil {
		log.Error().
			Err(err).
			Any("rsp", rsp).
			Msg("send response error")
	}
}

func playerInfoFromPBToPTO(pInfo *pb.PlayerInfo) pto.PlayerInfo {
	return pto.PlayerInfo{
		UID:         pInfo.Uid,
		GameMode:    constant.GameMode(pInfo.GameMode),
		ModeVersion: pInfo.ModeVersion,
		Star:        pInfo.Star,
		Rank:        pInfo.Rank,
		Glicko2Info: &pto.Glicko2Info{
			MMR:  pInfo.Glicko2Info.Mmr,
			Star: int(pInfo.Glicko2Info.Star),
			Rank: int(pInfo.Glicko2Info.Rank),
		},
	}
}

func checkPlayerInfo(info *pb.PlayerInfo) error {
	if info == nil {
		return errors.New("lack of player info")
	}
	if info.Glicko2Info == nil {
		return errors.New("lack of glicko2 info")
	}
	return nil
}
