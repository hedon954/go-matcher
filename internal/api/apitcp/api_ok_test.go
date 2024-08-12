package apitcp

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"

	"github.com/hedon954/go-matcher/internal/entry"
	"github.com/hedon954/go-matcher/internal/pb"
	"github.com/hedon954/go-matcher/pkg/typeconv"
	"github.com/hedon954/go-matcher/pkg/zinx/ziface"
	"github.com/hedon954/go-matcher/pkg/zinx/znet"

	"google.golang.org/protobuf/proto"
)

func init() {
	log.Logger = zerolog.New(io.Discard)
}

const (
	UIDA  = "a"
	UIDB  = "b"
	UIDBB = "bb"
	UIDC  = "c"
	UIDCC = "cc"
	UIDD  = "d"

	G1 int64 = 1
	G2 int64 = 2
	G3 int64 = 3
	G4 int64 = 4
	G5 int64 = 5
)

var dp = znet.NewDataPack()

func Test_TCP_ShouldWork(t *testing.T) {
	api := NewAPI(2, time.Millisecond*10)
	server := startServer()
	api.setupRouter(server)
	defer server.Stop()
	defer api.m.Stop()

	conn := startClient()
	defer func() {
		time.Sleep(10 * time.Millisecond) // TODO: refact it
		_ = conn.Close()
	}()

	// 1. 'a' create a group 'g1'
	rsp := requestCreateGroup(conn, UIDA, t)
	assert.Equal(t, int64(1), rsp.GroupId)
	fmt.Println(rsp)
	g1 := api.gm.Get(rsp.GroupId)
	assert.NotNil(t, g1)
	assert.Equal(t, 1, len(g1.Base().GetPlayers()))
	assert.Equal(t, G1, g1.ID())
	requestEnterGroup(conn, UIDA, rsp.GroupId, t) // re enter, no influence
	assert.Equal(t, 1, len(g1.Base().GetPlayers()))

	// 2. 'a' dissolve group 'g1'
	requestDissolveGroup(conn, UIDA, t)
	assert.Nil(t, api.gm.Get(rsp.GroupId))

	// 3. 'a' create a group again 'g2'
	rsp = requestCreateGroup(conn, UIDA, t)
	assert.Equal(t, int64(2), rsp.GroupId)
	g2 := api.gm.Get(rsp.GroupId)
	assert.NotNil(t, g2)
	assert.Equal(t, 1, len(g2.Base().GetPlayers()))
	assert.Equal(t, G2, g2.ID())
	ua := api.pm.Get(UIDA)
	assert.NotNil(t, ua)
	assert.Equal(t, entry.PlayerOnlineStateInGroup, ua.Base().GetOnlineStateWithLock())

	// 4. 'a' invite friend 'b'
	requestInvite(conn, UIDA, UIDB, t)
	assert.Equal(t, 1, len(g2.Base().GetInviteRecords()))

	// 5. 'b' refuse invite
	requestRefuseInvite(conn, UIDA, UIDB, g2.ID(), t)
	assert.Equal(t, 0, len(g2.Base().GetInviteRecords()))

	// 6. 'a' invite friend 'b' again
	requestInvite(conn, UIDA, UIDB, t)
	assert.Equal(t, 1, len(g2.Base().GetInviteRecords()))

	// 7. 'b' accept invite
	requestAcceptInvite(conn, UIDA, UIDB, g2.ID(), t)
	assert.Equal(t, 0, len(g2.Base().GetInviteRecords()))

	// 8. 'b' enter group
	requestEnterGroup(conn, UIDB, g2.ID(), t)
	assert.Equal(t, 2, len(g2.Base().GetPlayers()))
	ub := api.pm.Get(UIDB)
	assert.Equal(t, entry.PlayerOnlineStateInGroup, ub.Base().GetOnlineStateWithLock())
	assert.Equal(t, g2.ID(), ub.Base().GroupID)

	// 9. 'a' change role to 'b'
	assert.Equal(t, ua, g2.GetCaptain())
	requestChangeRole(conn, UIDA, UIDB, entry.GroupRoleCaptain, t)
	assert.Equal(t, ub, g2.GetCaptain())
	assert.Equal(t, 2, len(g2.Base().GetPlayers()))

	// 10. 'b' exit group 'g2'
	requestExitGroup(conn, UIDB, t)
	assert.Equal(t, 1, len(g2.Base().GetPlayers()))
	assert.Equal(t, ua, g2.GetCaptain())
	assert.Nil(t, api.pm.Get(UIDB))

	// 11. 'a' invite friend 'b'
	requestInvite(conn, UIDA, UIDB, t)
	assert.Equal(t, 1, len(g2.Base().GetInviteRecords()))

	// 12. 'b' accept invite
	requestAcceptInvite(conn, UIDA, UIDB, g2.ID(), t)
	assert.Equal(t, 0, len(g2.Base().GetInviteRecords()))

	// 13. 'b' enter group
	requestEnterGroup(conn, UIDB, g2.ID(), t)
	assert.Equal(t, 2, len(g2.Base().GetPlayers()))
	ub = api.pm.Get(UIDB)
	assert.NotNil(t, ub)
	assert.Equal(t, entry.PlayerOnlineStateInGroup, ub.Base().GetOnlineStateWithLock())

	// 14. 'a' kick player 'b'
	requestKick(conn, UIDA, UIDB, t)
	assert.Equal(t, 1, len(g2.Base().GetPlayers()))
	assert.Nil(t, api.pm.Get(UIDB))

	// 15. 'a' set voice state
	requestSetVoiceState(conn, UIDA, entry.PlayerVoiceStateUnmute, t)
	assert.Equal(t, entry.PlayerVoiceStateUnmute, ua.Base().GetVoiceState())

	// 16. 'a' set 'g2' allow nearby join
	assert.False(t, g2.Base().AllowNearbyJoin())
	requestSetNearbyJoinGroup(conn, UIDA, true, t)
	assert.True(t, g2.Base().AllowNearbyJoin())

	// 17. 'a' set 'g2' allow recent join
	assert.False(t, g2.Base().AllowRecentJoin())
	requestSetRecentJoinGroup(conn, UIDA, true, t)
	assert.True(t, g2.Base().AllowRecentJoin())

	// 18. 'a' unready
	requestUnready(conn, UIDA, t)
	assert.Equal(t, 1, len(g2.Base().UnReadyPlayer))

	// 19. 'a' ready
	requestReady(conn, UIDA, t)
	assert.Equal(t, 0, len(g2.Base().UnReadyPlayer))

	// 20. 'a' upload player attr
	requestUnloadPlayerAttr(conn, UIDA, t)
	assert.Equal(t, "hedon2", api.pm.Get(UIDA).Base().Attribute.Nickname)

	// 21. 'a' start match 'g2'
	requestStartMatch(conn, UIDA, t)
	time.Sleep(11 * time.Millisecond) // for at least get one round match
	assert.Equal(t, entry.GroupStateMatch, g2.Base().GetStateWithLock())
	assert.Equal(t, entry.PlayerOnlineStateInMatch, ua.Base().GetOnlineStateWithLock())

	// 22. 'a' cancel match 'g2'
	requestCancelMatch(conn, UIDA, t)
	time.Sleep(11 * time.Millisecond) // for at least get one round match
	assert.Equal(t, entry.GroupStateInvite, g2.Base().GetStateWithLock())
	assert.Equal(t, entry.PlayerOnlineStateInGroup, ua.Base().GetOnlineStateWithLock())

	// 23. 'a' start match 'g2'
	requestStartMatch(conn, UIDA, t)
	assert.Equal(t, entry.GroupStateMatch, g2.Base().GetStateWithLock())
	assert.Equal(t, entry.PlayerOnlineStateInMatch, ua.Base().GetOnlineStateWithLock())

	// 24. 'b' create a full group 'g3'
	requestCreateFullGroup(conn, UIDB, UIDBB, t)
	g3 := api.gm.Get(G3)
	assert.Equal(t, 2, len(g3.Base().GetPlayers()))

	// 25. 'c' create a full group 'g4'
	requestCreateFullGroup(conn, UIDC, UIDCC, t)
	g4 := api.gm.Get(G4)
	assert.Equal(t, 2, len(g4.Base().GetPlayers()))

	// 26. 'd' create group 'g5'
	requestCreateGroup(conn, UIDD, t)
	g5 := api.gm.Get(G5)
	assert.Equal(t, 1, len(g5.Base().GetPlayers()))

	// 27, start g3, g4 to match
	requestStartMatch(conn, UIDB, t)
	requestStartMatch(conn, UIDC, t)
	time.Sleep(20 * time.Millisecond)
	assert.Equal(t, entry.GroupStateMatch, g3.Base().GetStateWithLock())
	assert.Equal(t, entry.GroupStateMatch, g4.Base().GetStateWithLock())
	assert.Equal(t, entry.GroupStateMatch, g2.Base().GetStateWithLock())

	// 28. start g5 to match, match success and create a new room
	requestStartMatch(conn, UIDD, t)

	// TeamPlayerLimit: 2
	// RoomTeamLimit:   3
	// RoomPlayerCount = TeamPlayerLimit * RoomTeamLimit = 2 * 3 = 6
	// 'a' -> g2 -> 1
	// 'b' -> g3 -> 2
	// 'c' -> g4 -> 2
	// 'd' -> g5 -> 1
	for i := 0; i <= 10; i++ {
		if api.m.Glicko2Matcher.RoomCount.Load() == 1 {
			assert.Equal(t, entry.GroupStateGame, g2.Base().GetStateWithLock())
			assert.Equal(t, entry.GroupStateGame, g3.Base().GetStateWithLock())
			assert.Equal(t, entry.GroupStateGame, g4.Base().GetStateWithLock())
			assert.Equal(t, entry.GroupStateGame, g5.Base().GetStateWithLock())
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	rooms := []entry.Room{}
	api.rm.Range(func(_ int64, room entry.Room) bool {
		rooms = append(rooms, room)
		return true
	})

	// 29. 'd' exit game, 'g5' should be dissovled
	requestExitGame(conn, UIDD, rooms[0].ID(), t)
	assert.Nil(t, api.pm.Get(UIDD))
	assert.Nil(t, api.gm.Get(G5))
	assert.Equal(t, entry.GroupStateDissolved, g5.Base().GetStateWithLock())
}

func requestExitGame(conn net.Conn, uid string, roomID int64, t *testing.T) {
	var req = &pb.ExitGameReq{
		Uid:    uid,
		RoomId: roomID,
	}
	bs, _ := proto.Marshal(req)
	msg, err := dp.Pack(znet.NewMsgPackage(uint32(pb.ReqType_REQ_TYPE_EXIT_GAME), bs))
	assert.Nil(t, err)
	_, err = conn.Write(msg)
	assert.Nil(t, err)

	_ = readFromServer(conn, t).(*pb.ExitGameRsp)
}

func requestCreateFullGroup(conn net.Conn, uid1, uid2 string, t *testing.T) {
	rsp := requestCreateGroup(conn, uid1, t)
	requestEnterGroup(conn, uid2, rsp.GroupId, t)
}

func requestCancelMatch(conn net.Conn, uid string, t *testing.T) {
	var req = &pb.CancelMatchReq{
		Uid: uid,
	}
	bs, _ := proto.Marshal(req)
	msg, err := dp.Pack(znet.NewMsgPackage(uint32(pb.ReqType_REQ_TYPE_CANCEL_MATCH), bs))
	assert.Nil(t, err)
	_, err = conn.Write(msg)
	assert.Nil(t, err)

	_ = readFromServer(conn, t).(*pb.CancelMatchRsp)
}

func requestStartMatch(conn net.Conn, uid string, t *testing.T) {
	var req = &pb.StartMatchReq{
		Uid: uid,
	}
	bs, _ := proto.Marshal(req)
	msg, err := dp.Pack(znet.NewMsgPackage(uint32(pb.ReqType_REQ_TYPE_START_MATCH), bs))
	assert.Nil(t, err)
	_, err = conn.Write(msg)
	assert.Nil(t, err)

	_ = readFromServer(conn, t).(*pb.StartMatchRsp)
}

func requestUnloadPlayerAttr(conn net.Conn, uid string, t *testing.T) {
	var req = &pb.UploadPlayerAttrReq{
		Uid: uid,
		Attr: &pb.UserAttribute{
			Nickname: "hedon2",
		},
	}
	bs, _ := proto.Marshal(req)
	msg, err := dp.Pack(znet.NewMsgPackage(uint32(pb.ReqType_REQ_TYPE_UPLOAD_PLAYER_ATTR), bs))
	assert.Nil(t, err)
	_, err = conn.Write(msg)
	assert.Nil(t, err)

	_ = readFromServer(conn, t).(*pb.UploadPlayerAttrRsp)
}

func requestReady(conn net.Conn, uid string, t *testing.T) {
	var req = &pb.ReadyReq{
		Uid: uid,
	}
	bs, _ := proto.Marshal(req)
	msg, err := dp.Pack(znet.NewMsgPackage(uint32(pb.ReqType_REQ_TYPE_READY), bs))
	assert.Nil(t, err)
	_, err = conn.Write(msg)
	assert.Nil(t, err)

	_ = readFromServer(conn, t).(*pb.ReadyRsp)
}

func requestUnready(conn net.Conn, uid string, t *testing.T) {
	var req = &pb.UnreadyReq{
		Uid: uid,
	}
	bs, _ := proto.Marshal(req)
	msg, err := dp.Pack(znet.NewMsgPackage(uint32(pb.ReqType_REQ_TYPE_UNREADY), bs))
	assert.Nil(t, err)
	_, err = conn.Write(msg)
	assert.Nil(t, err)

	_ = readFromServer(conn, t).(*pb.UnreadyRsp)
}

func requestSetRecentJoinGroup(conn net.Conn, uid string, allow bool, t *testing.T) {
	var req = &pb.SetRecentJoinGroupReq{
		Uid:   uid,
		Allow: allow,
	}
	bs, _ := proto.Marshal(req)
	msg, err := dp.Pack(znet.NewMsgPackage(uint32(pb.ReqType_REQ_TYPE_SET_RECENT_JOIN_GROUP), bs))
	assert.Nil(t, err)
	_, err = conn.Write(msg)
	assert.Nil(t, err)

	_ = readFromServer(conn, t).(*pb.SetRecentJoinGroupRsp)
}

func requestSetNearbyJoinGroup(conn net.Conn, uid string, allow bool, t *testing.T) {
	var req = &pb.SetNearbyJoinGroupReq{
		Uid:   uid,
		Allow: allow,
	}
	bs, _ := proto.Marshal(req)
	msg, err := dp.Pack(znet.NewMsgPackage(uint32(pb.ReqType_REQ_TYPE_SET_NEARBY_JOIN_GROUP), bs))
	assert.Nil(t, err)
	_, err = conn.Write(msg)
	assert.Nil(t, err)

	_ = readFromServer(conn, t).(*pb.SetNearbyJoinGroupRsp)
}

func requestSetVoiceState(conn net.Conn, uid string, state entry.PlayerVoiceState, t *testing.T) {
	var req = &pb.SetVoiceStateReq{
		Uid:   uid,
		State: pb.PlayerVoiceState(state),
	}
	bs, _ := proto.Marshal(req)
	msg, err := dp.Pack(znet.NewMsgPackage(uint32(pb.ReqType_REQ_TYPE_SET_VOICE_STATE), bs))
	assert.Nil(t, err)
	_, err = conn.Write(msg)
	assert.Nil(t, err)

	_ = readFromServer(conn, t).(*pb.SetVoiceStateRsp)
}

func requestKick(conn net.Conn, captain, kicked string, t *testing.T) {
	var req = &pb.KickPlayerReq{
		CaptainUid: captain,
		KickedUid:  kicked,
	}
	bs, _ := proto.Marshal(req)
	msg, err := dp.Pack(znet.NewMsgPackage(uint32(pb.ReqType_REQ_TYPE_KICK_PLAYER), bs))
	assert.Nil(t, err)
	_, err = conn.Write(msg)
	assert.Nil(t, err)

	_ = readFromServer(conn, t).(*pb.KickPlayerRsp)
}

func requestExitGroup(conn net.Conn, uid string, t *testing.T) {
	var req = &pb.ExitGroupReq{
		Uid: uid,
	}
	bs, _ := proto.Marshal(req)
	msg, err := dp.Pack(znet.NewMsgPackage(uint32(pb.ReqType_REQ_TYPE_EXIT_GROUP), bs))
	assert.Nil(t, err)
	_, err = conn.Write(msg)
	assert.Nil(t, err)

	_ = readFromServer(conn, t).(*pb.ExitGroupRsp)
}

func requestChangeRole(conn net.Conn, captain, target string, role entry.GroupRole, t *testing.T) {
	var req = &pb.ChangeRoleReq{
		CaptainUid: captain,
		TargetUid:  target,
		Role:       pb.GroupRole(role),
	}
	bs, _ := proto.Marshal(req)
	msg, err := dp.Pack(znet.NewMsgPackage(uint32(pb.ReqType_REQ_TYPE_CHANGE_ROLE), bs))
	assert.Nil(t, err)
	_, err = conn.Write(msg)
	assert.Nil(t, err)

	_ = readFromServer(conn, t).(*pb.ChangeRoleRsp)
}

func requestAcceptInvite(conn net.Conn, inviter, invitee string, groupID int64, t *testing.T) {
	var req = &pb.AcceptInviteReq{
		InviterUid:  inviter,
		InviteeInfo: newPlayerInfo(invitee),
		GroupId:     groupID,
	}
	bs, _ := proto.Marshal(req)
	msg, err := dp.Pack(znet.NewMsgPackage(uint32(pb.ReqType_REQ_TYPE_ACCEPT_INVITE), bs))
	assert.Nil(t, err)
	_, err = conn.Write(msg)
	assert.Nil(t, err)

	_ = readFromServer(conn, t).(*pb.AcceptInviteRsp)
}

func requestRefuseInvite(conn net.Conn, inviter, invitee string, groupID int64, t *testing.T) {
	var req = &pb.RefuseInviteReq{
		InviterUid: inviter,
		InviteeUid: invitee,
		GroupId:    groupID,
	}
	bs, _ := proto.Marshal(req)
	msg, err := dp.Pack(znet.NewMsgPackage(uint32(pb.ReqType_REQ_TYPE_REFUSE_INVITE), bs))
	assert.Nil(t, err)
	_, err = conn.Write(msg)
	assert.Nil(t, err)

	_ = readFromServer(conn, t).(*pb.RefuseInviteRsp)
}

func requestInvite(conn net.Conn, inviter, invitee string, t *testing.T) {
	var req = &pb.InviteReq{
		InviterUid: inviter,
		InviteeUid: invitee,
	}
	bs, _ := proto.Marshal(req)
	msg, err := dp.Pack(znet.NewMsgPackage(uint32(pb.ReqType_REQ_TYPE_INVITE), bs))
	assert.Nil(t, err)
	_, err = conn.Write(msg)
	assert.Nil(t, err)

	_ = readFromServer(conn, t).(*pb.InviteRsp)
}

func requestDissolveGroup(conn net.Conn, uid string, t *testing.T) {
	var req = &pb.DissolveGroupReq{
		Uid: uid,
	}
	bs, _ := proto.Marshal(req)
	msg, err := dp.Pack(znet.NewMsgPackage(uint32(pb.ReqType_REQ_TYPE_DISSOLVE_GROUP), bs))
	assert.Nil(t, err)
	_, err = conn.Write(msg)
	assert.Nil(t, err)

	_ = readFromServer(conn, t).(*pb.DissolveGroupRsp)
}

func requestEnterGroup(conn net.Conn, uid string, groupID int64, t *testing.T) {
	var req = &pb.EnterGroupReq{
		PlayerInfo: newPlayerInfo(uid),
		Source:     pb.EnterGroupSource_ENTER_GROUP_SOURCE_INVITATION,
		GroupId:    groupID,
	}
	bs, _ := proto.Marshal(req)
	msg, err := dp.Pack(znet.NewMsgPackage(uint32(pb.ReqType_REQ_TYPE_ENTER_GROUP), bs))
	assert.Nil(t, err)
	_, err = conn.Write(msg)
	assert.Nil(t, err)

	_ = readFromServer(conn, t).(*pb.EnterGroupRsp)
}

func requestCreateGroup(conn net.Conn, uid string, t *testing.T) *pb.CreateGroupRsp {
	var req = &pb.CreateGroupReq{PlayerInfo: newPlayerInfo(uid)}
	bs, _ := proto.Marshal(req)
	msg, err := dp.Pack(znet.NewMsgPackage(uint32(pb.ReqType_REQ_TYPE_CREATE_GROUP), bs))
	assert.Nil(t, err)
	_, err = conn.Write(msg)
	assert.Nil(t, err)

	return readFromServer(conn, t).(*pb.CreateGroupRsp)
}

func newPlayerInfo(uid string) *pb.PlayerInfo {
	return &pb.PlayerInfo{
		Uid:         uid,
		GameMode:    pb.GameMode_GAME_MODE_GOAT_GAME,
		ModeVersion: 1,
		Star:        0,
		Rank:        0,
		Glicko2Info: &pb.Glicko2Info{
			Mmr:  0.0,
			Star: 0,
			Rank: 0,
		},
	}
}

func startServer() ziface.IServer {
	s := znet.NewServer("")
	s.Start()
	return s
}

func startClient() net.Conn {
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		panic("client dial error: " + err.Error())
	}
	return conn
}

func readFromServer(conn net.Conn, t *testing.T) (res any) {
	headData := make([]byte, dp.GetHeadLen())
	_, err := io.ReadFull(conn, headData)
	if errors.Is(err, os.ErrDeadlineExceeded) {
		return
	}
	if err != nil {
		return
	}

	// read msg head
	msgHead, err := dp.Unpack(headData)
	assert.Nil(t, err)

	if msgHead.GetDataLen() > 0 {
		// read msg body
		msg := msgHead.(*znet.Message)
		msg.Data = make([]byte, msg.GetDataLen())
		_, err = io.ReadFull(conn, msg.Data)
		if errors.Is(err, os.ErrDeadlineExceeded) {
			return
		}
		if err != nil {
			return
		}
		return resolveData(msg.Data)
	}
	return
}

//nolint:gocyclo
func resolveData(data []byte) any {
	var rsp = new(pb.CommonRsp)
	err := proto.Unmarshal(data, rsp)
	if err != nil {
		fmt.Println("unmarshal error", err)
		return ""
	}

	if rsp.Code != pb.RspCode_RSP_CODE_SUCCESS {
		return rsp.Message
	}
	switch rsp.ReqType {
	case pb.ReqType_REQ_TYPE_CREATE_GROUP:
		return typeconv.MustFromProto[pb.CreateGroupRsp](rsp.Data)
	case pb.ReqType_REQ_TYPE_ENTER_GROUP:
		return typeconv.MustFromProto[pb.EnterGroupRsp](rsp.Data)
	case pb.ReqType_REQ_TYPE_DISSOLVE_GROUP:
		return typeconv.MustFromProto[pb.DissolveGroupRsp](rsp.Data)
	case pb.ReqType_REQ_TYPE_INVITE:
		return typeconv.MustFromProto[pb.InviteRsp](rsp.Data)
	case pb.ReqType_REQ_TYPE_REFUSE_INVITE:
		return typeconv.MustFromProto[pb.RefuseInviteRsp](rsp.Data)
	case pb.ReqType_REQ_TYPE_ACCEPT_INVITE:
		return typeconv.MustFromProto[pb.AcceptInviteRsp](rsp.Data)
	case pb.ReqType_REQ_TYPE_CHANGE_ROLE:
		return typeconv.MustFromProto[pb.ChangeRoleRsp](rsp.Data)
	case pb.ReqType_REQ_TYPE_EXIT_GROUP:
		return typeconv.MustFromProto[pb.ExitGroupRsp](rsp.Data)
	case pb.ReqType_REQ_TYPE_KICK_PLAYER:
		return typeconv.MustFromProto[pb.KickPlayerRsp](rsp.Data)
	case pb.ReqType_REQ_TYPE_SET_VOICE_STATE:
		return typeconv.MustFromProto[pb.SetVoiceStateRsp](rsp.Data)
	case pb.ReqType_REQ_TYPE_SET_NEARBY_JOIN_GROUP:
		return typeconv.MustFromProto[pb.SetNearbyJoinGroupRsp](rsp.Data)
	case pb.ReqType_REQ_TYPE_SET_RECENT_JOIN_GROUP:
		return typeconv.MustFromProto[pb.SetRecentJoinGroupRsp](rsp.Data)
	case pb.ReqType_REQ_TYPE_UNREADY:
		return typeconv.MustFromProto[pb.UnreadyRsp](rsp.Data)
	case pb.ReqType_REQ_TYPE_READY:
		return typeconv.MustFromProto[pb.ReadyRsp](rsp.Data)
	case pb.ReqType_REQ_TYPE_UPLOAD_PLAYER_ATTR:
		return typeconv.MustFromProto[pb.UploadPlayerAttrRsp](rsp.Data)
	case pb.ReqType_REQ_TYPE_START_MATCH:
		return typeconv.MustFromProto[pb.StartMatchRsp](rsp.Data)
	case pb.ReqType_REQ_TYPE_CANCEL_MATCH:
		return typeconv.MustFromProto[pb.CancelMatchRsp](rsp.Data)
	case pb.ReqType_REQ_TYPE_EXIT_GAME:
		return typeconv.MustFromProto[pb.ExitGameRsp](rsp.Data)
	default:
		return nil
	}
}
