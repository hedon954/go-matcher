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

	"github.com/hedon954/go-matcher/internal/constant"
	"github.com/hedon954/go-matcher/internal/entry"
	"github.com/hedon954/go-matcher/internal/pb"
	"github.com/hedon954/go-matcher/pkg/typeconv"
	"github.com/hedon954/go-matcher/pkg/zinx/zconfig"
	"github.com/hedon954/go-matcher/pkg/zinx/ziface"
	"github.com/hedon954/go-matcher/pkg/zinx/znet"

	internalapi "github.com/hedon954/go-matcher/internal/api"

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

var dp = znet.NewDataPack(zconfig.DefaultConfig)

func Test_TCP_ShouldWork(t *testing.T) {
	inner, shutdown := internalapi.Start(newConf(2))
	defer shutdown()
	api := &API{inner}
	server, p := startServer()
	api.setupRouter(server)
	defer server.Stop()
	conn := startClient(p)
	defer func() { _ = conn.Close() }()

	// 1. 'a' create a group 'g1'
	rsp, errMsg := requestCreateGroup(conn, UIDA, t)
	assert.Equal(t, "", errMsg)
	assert.Equal(t, int64(1), rsp.GroupId)
	fmt.Println(rsp)
	g1 := api.GM.Get(rsp.GroupId)
	assert.NotNil(t, g1)
	assert.Equal(t, 1, len(g1.Base().GetPlayers()))
	assert.Equal(t, G1, g1.ID())
	requestEnterGroup(conn, UIDA, rsp.GroupId, t) // re enter, no influence
	assert.Equal(t, 1, len(g1.Base().GetPlayers()))

	// 2. 'a' dissolve group 'g1'
	requestDissolveGroup(conn, UIDA, t)
	assert.Nil(t, api.GM.Get(rsp.GroupId))

	// 3. 'a' create a group again 'g2'
	rsp, errMsg = requestCreateGroup(conn, UIDA, t)
	assert.Equal(t, "", errMsg)
	assert.Equal(t, int64(2), rsp.GroupId)
	g2 := api.GM.Get(rsp.GroupId)
	assert.NotNil(t, g2)
	assert.Equal(t, 1, len(g2.Base().GetPlayers()))
	assert.Equal(t, G2, g2.ID())
	ua := api.PM.Get(UIDA)
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
	ub := api.PM.Get(UIDB)
	assert.Equal(t, entry.PlayerOnlineStateInGroup, ub.Base().GetOnlineStateWithLock())
	assert.Equal(t, g2.ID(), ub.Base().GroupID)

	// 9. 'a' change role to 'b'
	assert.Equal(t, ua.UID(), g2.GetCaptain())
	requestChangeRole(conn, UIDA, UIDB, entry.GroupRoleCaptain, t)
	assert.Equal(t, ub.UID(), g2.GetCaptain())
	assert.Equal(t, 2, len(g2.Base().GetPlayers()))

	// 10. 'b' exit group 'g2'
	requestExitGroup(conn, UIDB, t)
	assert.Equal(t, 1, len(g2.Base().GetPlayers()))
	assert.Equal(t, ua.UID(), g2.GetCaptain())
	assert.Nil(t, api.PM.Get(UIDB))

	// 11. 'a' invite friend 'b'
	requestInvite(conn, UIDA, UIDB, t)
	assert.Equal(t, 1, len(g2.Base().GetInviteRecords()))

	// 12. 'b' accept invite
	requestAcceptInvite(conn, UIDA, UIDB, g2.ID(), t)
	assert.Equal(t, 0, len(g2.Base().GetInviteRecords()))

	// 13. 'b' enter group
	requestEnterGroup(conn, UIDB, g2.ID(), t)
	assert.Equal(t, 2, len(g2.Base().GetPlayers()))
	ub = api.PM.Get(UIDB)
	assert.NotNil(t, ub)
	assert.Equal(t, entry.PlayerOnlineStateInGroup, ub.Base().GetOnlineStateWithLock())

	// 14. 'a' kick player 'b'
	requestKick(conn, UIDA, UIDB, t)
	assert.Equal(t, 1, len(g2.Base().GetPlayers()))
	assert.Nil(t, api.PM.Get(UIDB))

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
	requestUnloadPlayerAttr(conn, UIDA, true, true, t)
	assert.Equal(t, "hedon2", api.PM.Get(UIDA).Base().Attribute.Nickname)

	// 21. 'a' start match 'g2'
	requestStartMatch(conn, UIDA, t)
	time.Sleep(11 * time.Millisecond) // for at least get one round match
	assert.Equal(t, entry.GroupStateMatch, g2.Base().GetStateWithLock())
	assert.Equal(t, entry.PlayerOnlineStateInMatch, ua.Base().GetOnlineStateWithLock())

	// 22. 'a' cancel match 'g2'
	_, errMsg = requestCancelMatch(conn, UIDA, t)
	assert.Equal(t, "", errMsg)
	time.Sleep(11 * time.Millisecond) // for at least get one round match
	assert.Equal(t, entry.GroupStateInvite, g2.Base().GetStateWithLock())
	assert.Equal(t, entry.PlayerOnlineStateInGroup, ua.Base().GetOnlineStateWithLock())

	// 23. 'a' start match 'g2'
	requestStartMatch(conn, UIDA, t)
	assert.Equal(t, entry.GroupStateMatch, g2.Base().GetStateWithLock())
	assert.Equal(t, entry.PlayerOnlineStateInMatch, ua.Base().GetOnlineStateWithLock())

	// 24. 'b' create a full group 'g3'
	requestCreateFullGroup(conn, UIDB, UIDBB, t)
	g3 := api.GM.Get(G3)
	assert.Equal(t, 2, len(g3.Base().GetPlayers()))

	// 25. 'c' create a full group 'g4'
	requestCreateFullGroup(conn, UIDC, UIDCC, t)
	g4 := api.GM.Get(G4)
	assert.Equal(t, 2, len(g4.Base().GetPlayers()))

	// 26. 'd' create group 'g5'
	requestCreateGroup(conn, UIDD, t)
	g5 := api.GM.Get(G5)
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
	success := false
	for i := 0; i <= 10; i++ {
		if api.M.Glicko2Matcher.RoomCount.Load() == 1 {
			assert.Equal(t, entry.GroupStateGame, g2.Base().GetStateWithLock())
			assert.Equal(t, entry.GroupStateGame, g3.Base().GetStateWithLock())
			assert.Equal(t, entry.GroupStateGame, g4.Base().GetStateWithLock())
			assert.Equal(t, entry.GroupStateGame, g5.Base().GetStateWithLock())
			success = true
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	assert.True(t, success)

	rooms := []entry.Room{}
	api.RM.Range(func(_ int64, room entry.Room) bool {
		rooms = append(rooms, room)
		return true
	})

	// 29. 'd' exit game, 'g5' should be dissolved
	requestExitGame(conn, UIDD, rooms[0].ID(), t)
	assert.Nil(t, api.PM.Get(UIDD))
	assert.Nil(t, api.GM.Get(G5))
	assert.Equal(t, entry.GroupStateDissolved, g5.Base().GetStateWithLock())
}

func requestExitGame(conn net.Conn, uid string, roomID int64, t *testing.T) string {
	var req = &pb.ExitGameReq{
		Uid:    uid,
		RoomId: roomID,
	}
	bs, _ := proto.Marshal(req)
	msg, err := dp.Pack(znet.NewMsgPackage(uint32(pb.ReqType_REQ_TYPE_EXIT_GAME), bs))
	assert.Nil(t, err)
	_, err = conn.Write(msg)
	assert.Nil(t, err)

	rsp, em := readFromServer(conn, t)
	_ = rsp.(*pb.ExitGameRsp)
	return em
}

func requestCreateFullGroup(conn net.Conn, uid1, uid2 string, t *testing.T) {
	rsp, errMsg := requestCreateGroup(conn, uid1, t)
	assert.Equal(t, "", errMsg)
	requestEnterGroup(conn, uid2, rsp.GroupId, t)
}

func requestCancelMatch(conn net.Conn, uid string, t *testing.T) (*pb.CancelMatchRsp, string) {
	var req = &pb.CancelMatchReq{
		Uid: uid,
	}
	bs, _ := proto.Marshal(req)
	msg, err := dp.Pack(znet.NewMsgPackage(uint32(pb.ReqType_REQ_TYPE_CANCEL_MATCH), bs))
	assert.Nil(t, err)
	_, err = conn.Write(msg)
	assert.Nil(t, err)

	rsp, em := readFromServer(conn, t)
	if em == "" {
		return rsp.(*pb.CancelMatchRsp), ""
	}
	return nil, em
}

func requestStartMatch(conn net.Conn, uid string, t *testing.T) string {
	var req = &pb.StartMatchReq{
		Uid: uid,
	}
	bs, _ := proto.Marshal(req)
	msg, err := dp.Pack(znet.NewMsgPackage(uint32(pb.ReqType_REQ_TYPE_START_MATCH), bs))
	assert.Nil(t, err)
	_, err = conn.Write(msg)
	assert.Nil(t, err)

	rsp, em := readFromServer(conn, t)
	if em == "" {
		_ = rsp.(*pb.StartMatchRsp)
	}
	return em
}

func requestUnloadPlayerAttr(conn net.Conn, uid string, attr, goat bool, t *testing.T) string {
	var req = &pb.UploadPlayerAttrReq{
		Uid: uid,
	}
	if attr {
		req.Attr = &pb.UserAttribute{
			Nickname: "hedon2",
		}
	}
	if goat {
		req.Type = &pb.UploadPlayerAttrReq_GoatGameAttr{
			GoatGameAttr: &pb.GoatGameAttribute{Mmr: 1.0},
		}
	}
	bs, _ := proto.Marshal(req)
	msg, err := dp.Pack(znet.NewMsgPackage(uint32(pb.ReqType_REQ_TYPE_UPLOAD_PLAYER_ATTR), bs))
	assert.Nil(t, err)
	_, err = conn.Write(msg)
	assert.Nil(t, err)

	rsp, em := readFromServer(conn, t)
	if em == "" {
		_ = rsp.(*pb.UploadPlayerAttrRsp)
	}
	return em
}

func requestReady(conn net.Conn, uid string, t *testing.T) string {
	var req = &pb.ReadyReq{
		Uid: uid,
	}
	bs, _ := proto.Marshal(req)
	msg, err := dp.Pack(znet.NewMsgPackage(uint32(pb.ReqType_REQ_TYPE_READY), bs))
	assert.Nil(t, err)
	_, err = conn.Write(msg)
	assert.Nil(t, err)

	rsp, em := readFromServer(conn, t)
	if em == "" {
		_ = rsp.(*pb.ReadyRsp)
	}
	return em
}

func requestUnready(conn net.Conn, uid string, t *testing.T) string {
	var req = &pb.UnreadyReq{
		Uid: uid,
	}
	bs, _ := proto.Marshal(req)
	msg, err := dp.Pack(znet.NewMsgPackage(uint32(pb.ReqType_REQ_TYPE_UNREADY), bs))
	assert.Nil(t, err)
	_, err = conn.Write(msg)
	assert.Nil(t, err)

	rsp, em := readFromServer(conn, t)
	if em == "" {
		_ = rsp.(*pb.UnreadyRsp)
	}
	return em
}

func requestSetRecentJoinGroup(conn net.Conn, uid string, allow bool, t *testing.T) string {
	var req = &pb.SetRecentJoinGroupReq{
		Uid:   uid,
		Allow: allow,
	}
	bs, _ := proto.Marshal(req)
	msg, err := dp.Pack(znet.NewMsgPackage(uint32(pb.ReqType_REQ_TYPE_SET_RECENT_JOIN_GROUP), bs))
	assert.Nil(t, err)
	_, err = conn.Write(msg)
	assert.Nil(t, err)

	rsp, em := readFromServer(conn, t)
	if em == "" {
		_ = rsp.(*pb.SetRecentJoinGroupRsp)
	}
	return em
}

func requestSetNearbyJoinGroup(conn net.Conn, uid string, allow bool, t *testing.T) string {
	var req = &pb.SetNearbyJoinGroupReq{
		Uid:   uid,
		Allow: allow,
	}
	bs, _ := proto.Marshal(req)
	msg, err := dp.Pack(znet.NewMsgPackage(uint32(pb.ReqType_REQ_TYPE_SET_NEARBY_JOIN_GROUP), bs))
	assert.Nil(t, err)
	_, err = conn.Write(msg)
	assert.Nil(t, err)

	rsp, em := readFromServer(conn, t)
	if em == "" {
		_ = rsp.(*pb.SetNearbyJoinGroupRsp)
	}
	return em
}

func requestSetVoiceState(conn net.Conn, uid string, state entry.PlayerVoiceState, t *testing.T) string {
	var req = &pb.SetVoiceStateReq{
		Uid:   uid,
		State: pb.PlayerVoiceState(state),
	}
	bs, _ := proto.Marshal(req)
	msg, err := dp.Pack(znet.NewMsgPackage(uint32(pb.ReqType_REQ_TYPE_SET_VOICE_STATE), bs))
	assert.Nil(t, err)
	_, err = conn.Write(msg)
	assert.Nil(t, err)

	rsp, em := readFromServer(conn, t)
	if em == "" {
		_ = rsp.(*pb.SetVoiceStateRsp)
	}
	return em
}

func requestKick(conn net.Conn, captain, kicked string, t *testing.T) string {
	var req = &pb.KickPlayerReq{
		CaptainUid: captain,
		KickedUid:  kicked,
	}
	bs, _ := proto.Marshal(req)
	msg, err := dp.Pack(znet.NewMsgPackage(uint32(pb.ReqType_REQ_TYPE_KICK_PLAYER), bs))
	assert.Nil(t, err)
	_, err = conn.Write(msg)
	assert.Nil(t, err)

	rsp, em := readFromServer(conn, t)
	if em == "" {
		_ = rsp.(*pb.KickPlayerRsp)
	}
	return em
}

func requestExitGroup(conn net.Conn, uid string, t *testing.T) string {
	var req = &pb.ExitGroupReq{
		Uid: uid,
	}
	bs, _ := proto.Marshal(req)
	msg, err := dp.Pack(znet.NewMsgPackage(uint32(pb.ReqType_REQ_TYPE_EXIT_GROUP), bs))
	assert.Nil(t, err)
	_, err = conn.Write(msg)
	assert.Nil(t, err)

	rsp, em := readFromServer(conn, t)
	if em == "" {
		_ = rsp.(*pb.ExitGroupRsp)
	}
	return em
}

func requestChangeRole(conn net.Conn, captain, target string, role entry.GroupRole, t *testing.T) string {
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

	ret, em := readFromServer(conn, t)
	if em == "" {
		_ = ret.(*pb.ChangeRoleRsp)
	}
	return em
}

func requestAcceptInvite(conn net.Conn, inviter, invitee string, groupID int64, t *testing.T) string {
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

	rsp, em := readFromServer(conn, t)
	if em == "" {
		_ = rsp.(*pb.AcceptInviteRsp)
	}
	return em
}

func requestRefuseInvite(conn net.Conn, inviter, invitee string, groupID int64, t *testing.T) string {
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

	rsp, em := readFromServer(conn, t)
	if em == "" {
		_ = rsp.(*pb.RefuseInviteRsp)
	}
	return em
}

func requestInvite(conn net.Conn, inviter, invitee string, t *testing.T) string {
	var req = &pb.InviteReq{
		InviterUid: inviter,
		InviteeUid: invitee,
	}
	bs, _ := proto.Marshal(req)
	msg, err := dp.Pack(znet.NewMsgPackage(uint32(pb.ReqType_REQ_TYPE_INVITE), bs))
	assert.Nil(t, err)
	_, err = conn.Write(msg)
	assert.Nil(t, err)

	rsp, em := readFromServer(conn, t)
	if em == "" {
		_ = rsp.(*pb.InviteRsp)
	}
	return em
}

func requestDissolveGroup(conn net.Conn, uid string, t *testing.T) string {
	var req = &pb.DissolveGroupReq{
		Uid: uid,
	}
	bs, _ := proto.Marshal(req)
	msg, err := dp.Pack(znet.NewMsgPackage(uint32(pb.ReqType_REQ_TYPE_DISSOLVE_GROUP), bs))
	assert.Nil(t, err)
	_, err = conn.Write(msg)
	assert.Nil(t, err)

	rsp, em := readFromServer(conn, t)
	if em == "" {
		_ = rsp.(*pb.DissolveGroupRsp)
	}
	return em
}

func requestEnterGroup(conn net.Conn, uid string, groupID int64, t *testing.T) string {
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

	rsp, em := readFromServer(conn, t)
	if em == "" {
		_ = rsp.(*pb.EnterGroupRsp)
	}
	return em
}

func requestCreateGroup(conn net.Conn, uid string, t *testing.T) (*pb.CreateGroupRsp, string) {
	return requestCreateGroupWithMode(conn, uid, constant.GameModeGoatGame, t)
}

func requestCreateGroupWithMode(conn net.Conn, uid string, mode constant.GameMode, t *testing.T) (resp *pb.CreateGroupRsp, errMsg string) {
	return requestCreateGroupWithModeAndModeVersion(conn, uid, mode, 1, t)
}

func requestCreateGroupWithModeAndModeVersion(conn net.Conn, uid string, mode constant.GameMode, modeVersion int64,
	t *testing.T) (resp *pb.CreateGroupRsp, errMsg string) {
	var req = &pb.CreateGroupReq{PlayerInfo: newPlayerInfoWithModeAndModeVersion(uid, mode, modeVersion)}
	bs, _ := proto.Marshal(req)
	msg, err := dp.Pack(znet.NewMsgPackage(uint32(pb.ReqType_REQ_TYPE_CREATE_GROUP), bs))
	assert.Nil(t, err)
	_, err = conn.Write(msg)
	assert.Nil(t, err)

	rsp, em := readFromServer(conn, t)
	if em == "" {
		return rsp.(*pb.CreateGroupRsp), ""
	}
	return nil, em
}

func newPlayerInfo(uid string) *pb.PlayerInfo {
	return newPlayerInfoWithMode(uid, constant.GameModeGoatGame)
}

func newPlayerInfoWithMode(uid string, mode constant.GameMode) *pb.PlayerInfo {
	return newPlayerInfoWithModeAndModeVersion(uid, mode, 1)
}

func newPlayerInfoWithModeAndModeVersion(uid string, mode constant.GameMode, modeVersion int64) *pb.PlayerInfo {
	return &pb.PlayerInfo{
		Uid:         uid,
		GameMode:    pb.GameMode(mode),
		ModeVersion: modeVersion,
		Star:        0,
		Rank:        0,
		Glicko2Info: &pb.Glicko2Info{
			Mmr:  0.0,
			Star: 0,
			Rank: 0,
		},
	}
}

func startServer() (server ziface.IServer, p int) {
	conf := *zconfig.DefaultConfig
	conf.TCPPort = int(port.Add(1))
	s := znet.NewServer(&conf)
	go s.Start()
	time.Sleep(3 * time.Millisecond) // wait for server to start
	return s, conf.TCPPort
}

func startClient(port int) net.Conn {
	conn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		panic("client dial error: " + err.Error())
	}
	return conn
}

func readFromServer(conn net.Conn, t *testing.T) (res any, errMsg string) {
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
func resolveData(data []byte) (any, string) {
	var rsp = new(pb.CommonRsp)
	err := proto.Unmarshal(data, rsp)
	if err != nil {
		fmt.Println("unmarshal error", err)
		return nil, ""
	}

	if rsp.Code != pb.RspCode_RSP_CODE_SUCCESS {
		return nil, rsp.Message
	}
	switch rsp.ReqType {
	case pb.ReqType_REQ_TYPE_CREATE_GROUP:
		return typeconv.MustFromProto[pb.CreateGroupRsp](rsp.Data), ""
	case pb.ReqType_REQ_TYPE_ENTER_GROUP:
		return typeconv.MustFromProto[pb.EnterGroupRsp](rsp.Data), ""
	case pb.ReqType_REQ_TYPE_DISSOLVE_GROUP:
		return typeconv.MustFromProto[pb.DissolveGroupRsp](rsp.Data), ""
	case pb.ReqType_REQ_TYPE_INVITE:
		return typeconv.MustFromProto[pb.InviteRsp](rsp.Data), ""
	case pb.ReqType_REQ_TYPE_REFUSE_INVITE:
		return typeconv.MustFromProto[pb.RefuseInviteRsp](rsp.Data), ""
	case pb.ReqType_REQ_TYPE_ACCEPT_INVITE:
		return typeconv.MustFromProto[pb.AcceptInviteRsp](rsp.Data), ""
	case pb.ReqType_REQ_TYPE_CHANGE_ROLE:
		return typeconv.MustFromProto[pb.ChangeRoleRsp](rsp.Data), ""
	case pb.ReqType_REQ_TYPE_EXIT_GROUP:
		return typeconv.MustFromProto[pb.ExitGroupRsp](rsp.Data), ""
	case pb.ReqType_REQ_TYPE_KICK_PLAYER:
		return typeconv.MustFromProto[pb.KickPlayerRsp](rsp.Data), ""
	case pb.ReqType_REQ_TYPE_SET_VOICE_STATE:
		return typeconv.MustFromProto[pb.SetVoiceStateRsp](rsp.Data), ""
	case pb.ReqType_REQ_TYPE_SET_NEARBY_JOIN_GROUP:
		return typeconv.MustFromProto[pb.SetNearbyJoinGroupRsp](rsp.Data), ""
	case pb.ReqType_REQ_TYPE_SET_RECENT_JOIN_GROUP:
		return typeconv.MustFromProto[pb.SetRecentJoinGroupRsp](rsp.Data), ""
	case pb.ReqType_REQ_TYPE_UNREADY:
		return typeconv.MustFromProto[pb.UnreadyRsp](rsp.Data), ""
	case pb.ReqType_REQ_TYPE_READY:
		return typeconv.MustFromProto[pb.ReadyRsp](rsp.Data), ""
	case pb.ReqType_REQ_TYPE_UPLOAD_PLAYER_ATTR:
		return typeconv.MustFromProto[pb.UploadPlayerAttrRsp](rsp.Data), ""
	case pb.ReqType_REQ_TYPE_START_MATCH:
		return typeconv.MustFromProto[pb.StartMatchRsp](rsp.Data), ""
	case pb.ReqType_REQ_TYPE_CANCEL_MATCH:
		return typeconv.MustFromProto[pb.CancelMatchRsp](rsp.Data), ""
	case pb.ReqType_REQ_TYPE_EXIT_GAME:
		return typeconv.MustFromProto[pb.ExitGameRsp](rsp.Data), ""
	default:
		return nil, "unknown req type: " + rsp.ReqType.String()
	}
}
