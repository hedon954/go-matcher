package matcher

import (
	"errors"
	"time"

	"github.com/hedon954/go-matcher/common"
	"github.com/hedon954/go-matcher/config"
	"github.com/hedon954/go-matcher/def"
	"github.com/hedon954/go-matcher/group"
	"github.com/hedon954/go-matcher/merr"
	"github.com/hedon954/go-matcher/pkg/rdstimer"
	"github.com/hedon954/go-matcher/player"
	"github.com/hedon954/go-matcher/pto"
	"github.com/hedon954/go-matcher/rpc/rpcclient/connector"

	"github.com/samborkent/uuidv7"
)

// TODO: 应该支持多种实现
// TODO: 先不考虑状态的流转，先把整个流程完成再说，后续尝试使用状态机来进行改进管理

type Impl struct {
	rtm       *rdstimer.TimerManager
	playerMgr *player.Manager
	groupMgr  *group.Manager

	connectorRpc *connector.Manager

	helper GameHelper
	config config.GroupConfig
}

func New(h GameHelper, c config.GroupConfig) *Impl {
	impl := &Impl{
		rtm:          rdstimer.NewRdsTimerManager(),
		playerMgr:    player.NewManager(),
		groupMgr:     group.NewManager(0), // TODO: read groupID from file
		connectorRpc: connector.NewMgr(),
		helper:       h,
		config:       c,
	}
	return impl
}

func (impl *Impl) BindRestore(i int, s string) (state int, matchInfo interface{}) {
	// TODO implement me
	panic("implement me")
}

func (impl *Impl) CreateGroup(info *pto.PlayerInfo) (int64, error) {
	g, err := impl.createGroup(info)
	if err != nil {
		return 0, err
	}
	return g.GroupID(), nil
}

func (impl *Impl) createGroup(info *pto.PlayerInfo) (common.Group, error) {
	// 1. 如果玩家已经在组里，则尝试返回之前的组
	if g := impl.inOwnedGroup(info); g != nil {
		return g, nil
	}

	// 2. 创建新的组
	g, err := impl.groupMgr.CreateGroup(info)
	if err != nil {
		return nil, err
	}
	return g, nil
}

// inOwnedGroup 判断玩家是否在自己的房间里，如果是，则直接返回，如果不是，则退出之前的队伍（如果在的话）。
func (impl *Impl) inOwnedGroup(info *pto.PlayerInfo) common.Group {
	p := impl.playerMgr.GetPlayer(info.Uid)
	if p == nil {
		return nil
	}
	g := impl.groupMgr.GetGroup(p.Base().GroupID)
	if g == nil {
		return nil
	}
	g.Base().Lock()
	defer g.Base().Unlock()

	// 只有自己是房主的时候，才返回原来的房间
	if g.Base().IsOwner(info.Uid) {
		return g
	}

	// 不是房主的时候，需要先退出之前的房间
	g.Base().DelPlayer(info.Uid)
	return nil
}

func (impl *Impl) InviteFriend(param *pto.InviteFriend) error {
	// 自己首先要在组里
	_, g, err := impl.getPlayerAndGroup(param.InviteUid)
	if err != nil {
		return err
	}

	g.Base().Lock()
	defer g.Base().Unlock()

	// 分享邀请提前返回
	if param.Source == pto.InvitationSrcShare {
		g.Base().MarkSrcShare()
		return nil
	}

	// 如果被邀请者已经在组里了，直接返回
	invitee := impl.playerMgr.GetPlayer(param.FriendUid)
	if invitee != nil {
		if g.Base().PlayerExists(invitee) {
			impl.BroadcastGroupUsers(g)
			return nil
		}
	}

	if g.IsFull() {
		return merr.ErrGroupFull
	}

	// 发送邀请
	g.AddInvitedPlayer(param.FriendUid)
	impl.connectorRpc.PushInviteFriend(param)
	return nil
}

func (impl *Impl) getPlayerAndGroup(uid string) (common.Player, common.Group, error) {
	p := impl.playerMgr.GetPlayer(uid)
	if p == nil {
		return nil, nil, merr.ErrPlayerNotInGroup
	}
	g := impl.groupMgr.GetGroup(p.Base().GroupID)
	if g == nil {
		return nil, nil, merr.ErrPlayerNotInGroup
	}
	return p, g, nil
}

func (impl *Impl) HandleInvite(param *pto.HandleInvite) error {
	// 判断队伍在不在
	inviter, g, err := impl.getPlayerAndGroup(param.InviteUid)
	if err != nil {
		return merr.ErrGroupDissolved
	}

	// 判断队伍状态是否正确
	g.Base().Lock()
	defer g.Base().Unlock()
	if err := g.CheckState(common.GroupStateInvite); err != nil {
		return err
	}

	// 判断邀请是否已经过期
	if err := g.CheckHandleInviteExpired(param.Player.Uid, param.SrcType); err != nil {
		return err
	}

	// 拒绝
	if param.HandleType == pto.InviteHandleTypeRefuseMsg {
		return impl.refuseInvite(param)
	}

	// 判断队伍人数是否已满
	if g.IsFull() {
		return merr.ErrGroupFull
	}

	// 构建被邀请者对象
	invitee := impl.playerMgr.GetPlayer(param.Player.Uid)
	if invitee == nil {
		invitee, err = impl.playerMgr.CreatePlayer(param.Player)
		if err != nil {
			return err
		}
	}

	// TODO: 这样做不好，采用命令式，不要访问属性
	invitee.Base().GameMode = param.Player.GameMode
	invitee.Base().MatchStrategy = param.Player.MatchStrategy
	inviter.Base().ModeVersion = param.Player.ModeVersion

	// 判断版本是否一致
	if err := impl.checkVersion(inviter, invitee); err != nil {
		return err
	}

	// 接受请求
	impl.acceptInvite(g, inviter, invitee, param.SrcType)
	return nil
}

func (impl *Impl) checkVersion(inviter, invitee common.Player) error {
	if !invitee.Base().VersionMatched(inviter) {
		return merr.ErrVersionNotSame
	}
	return nil
}

func (impl *Impl) refuseInvite(param *pto.HandleInvite) error {
	if !impl.checkRefuseMsgIsValid(param.Message) {
		return merr.ErrInvalidMessage
	}
	return impl.connectorRpc.PushHandleInvite(param.InviteUid, param.Player.Uid, pto.InviteHandleTypeRefuseMsg,
		param.Message)
}

func (impl *Impl) checkRefuseMsgIsValid(msg string) bool {
	return true
}

func (impl *Impl) acceptInvite(g common.Group, inviter, invitee common.Player, srcType pto.InvitationSrcType) {
	// 如果之前就在队伍里，则需要退出
	preG := impl.groupMgr.GetGroup(invitee.Base().GroupID)
	if preG != nil {
		preG.Base().DelPlayer(invitee.UID())
	}

	// 加入队伍
	invitee.Base().GroupID = g.GroupID()
	invitee.Base().SetOnlineState(common.PlayerOnlineStateGroup)
	g.Base().AddPlayer(invitee)
	impl.BroadcastGroupUsers(g)

	if srcType == pto.InvitationSrcNearBy {
		g.Base().AddNearby(inviter.UID(), invitee.UID())
	}

	// 通知被邀请者
	_ = impl.connectorRpc.PushHandleInvite(inviter.UID(), inviter.UID(), pto.InviteHandleTypeAccept, "")
	impl.connectorRpc.UpdateOnlineState(g.Base().UIDs(), common.PlayerOnlineStateGroup)

	// 定时广播队伍状态
	fireTime := time.Now().UnixMilli() + impl.config.BroadcastUsersTimeoutMs
	impl.rtm.AddTimer(def.RdsTimerGroupBroadcast, int(g.GroupID()), int(fireTime))
}

func (impl *Impl) Kick(uid, kickUid string) error {
	if uid == kickUid {
		return merr.ErrCantKickSelf
	}

	p, g, err := impl.getPlayerAndGroup(uid)
	if err != nil {
		return err
	}

	// 判断队伍状态是否正确
	g.Base().Lock()
	defer g.Base().Unlock()
	if err := g.CheckState(common.GroupStateInvite); err != nil {
		return err
	}

	if !g.Base().IsOwner(uid) {
		return merr.ErrPermissionDeny
	}

	kp := impl.playerMgr.GetPlayer(kickUid)
	if !g.Base().PlayerExists(p) {
		return nil
	}

	impl.removePlayer(g, kp.UID())
	return nil
}

func (impl *Impl) removePlayer(g common.Group, uid string) {
	if !g.Base().DelPlayer(uid) {
		impl.BroadcastGroupUsers(g)
	}
}

func (impl *Impl) ExitGroup(uid string) error {
	p, g, err := impl.getPlayerAndGroup(uid)
	if err != nil {
		return nil
	}

	g.Base().Lock()
	defer g.Base().Unlock()
	if err := g.CheckState(common.GroupStateInvite); err != nil {
		return err
	}

	impl.exitGroup(p, g)
	return nil
}

func (impl *Impl) exitGroup(p common.Player, g common.Group) {
	groupIsEmpty := g.Base().DelPlayer(p.UID())
	if groupIsEmpty {
		impl.groupMgr.DelGroup(g.GroupID())
	}
	impl.playerMgr.DelPlayer(p.UID())
	impl.BroadcastGroupUsers(g)

	p.Base().GroupID = 0
	p.Base().SetOnlineState(common.PlayerOnlineStateOnline)
	impl.connectorRpc.UpdateOnlineState([]string{p.UID()}, common.PlayerOnlineStateOnline)
}

func (impl *Impl) DissolveGroup(uid string) error {
	_, g, err := impl.getPlayerAndGroup(uid)
	if err != nil {
		return err
	}
	g.Base().Lock()
	defer g.Base().Unlock()
	if !g.Base().IsOwner(uid) {
		return merr.ErrPermissionDeny
	}
	if err := g.CheckState(common.GroupStateInvite); err != nil {
		return err
	}

	impl.inviteCardExpire(uid, g)
	return impl.dissolve(g)
}

func (impl *Impl) inviteCardExpire(uid string, g common.Group) {
	_, ok := g.Base().InvitedPlayers[string(common.SrcClanChat)]
	if ok {
		impl.connectorRpc.UpdateInviteCard(uid, common.ChatGroupExpire, 0, common.SrcClanChat)
		delete(g.Base().InvitedPlayers, string(common.SrcClanChat))
	}
	_, ok = g.Base().InvitedPlayers[string(common.SrcChannelChat)]
	if ok {
		impl.connectorRpc.UpdateInviteCard(uid, common.ChatGroupExpire, 0, common.SrcChannelChat)
		delete(g.Base().InvitedPlayers, string(common.SrcChannelChat))
	}
	_, ok = g.Base().InvitedPlayers[string(common.SrcSingleChat)]
	if ok {
		impl.connectorRpc.UpdateInviteCard(uid, common.ChatGroupExpire, 0, common.SrcSingleChat)
		delete(g.Base().InvitedPlayers, string(common.SrcSingleChat))
	}
}

func (impl *Impl) dissolve(g common.Group) error {
	g.Base().SetState(common.GroupStateDissolve)

	impl.deleteInviteTimer(g.GroupID())
	impl.groupMgr.DelGroup(g.GroupID())
	uids := g.Base().UIDs()
	for _, uid := range uids {
		impl.playerMgr.DelPlayer(uid)
	}

	impl.connectorRpc.GroupDissolved(uids, g.GroupID())
	return nil
}

func (impl *Impl) Match(info *pto.PlayerInfo, reply *MatchReply) error {
	if impl.helper == nil {
		return merr.ErrGameOffline
	}
	if impl.helper.IsBanGame(info.Uid) {
		return merr.ErrUserBanned
	}
	if !impl.helper.IsInGameTime(time.Now().Unix()) {
		return merr.ErrGameNotOpen
	}

	g, err := impl.createGroup(info)
	if err != nil {
		return err
	}

	g.Base().Lock()
	defer g.Base().Unlock()
	if err := g.CheckState(common.GroupStateInvite); err != nil {
		return err
	}

	if len(g.Base().UnReadyUsers) > 0 {
		return merr.ErrGroupPlayerNotReady
	}

	// 刷新基本信息
	g.Base().Name = info.GroupName
	g.Base().MatchID = uuidv7.New().String()
	g.Base().SetCoupleInfo()
	impl.inviteCardExpire(info.Uid, g)

	// 更新状态
	g.Base().SetState(common.GroupStateQueuing)
	for _, p := range g.Base().GetPlayers() {
		p.Base().SetOnlineState(common.PlayerOnlineStateQueuing)
	}
	uids := g.Base().UIDs()
	impl.connectorRpc.UpdateOnlineState(uids, common.PlayerOnlineStateQueuing)
	impl.broadGroupState(g)

	// 更新定时器
	impl.deleteInviteTimer(g.GroupID())
	impl.addWaitAttrTimer(g.GroupID())
	impl.addQueueTimer(g.GroupID())
	return nil
}

func (impl *Impl) CancelMatch(param *CancelMatch) error {
	_, g, err := impl.getPlayerAndGroup(param.Uid)
	if err != nil {
		return err
	}
	g.Base().Lock()
	defer g.Base().Unlock()
	if err := g.CheckState(common.GroupStateQueuing); err != nil {
		return err
	}

	// 更新状态
	g.Base().SetState(common.GroupStateInvite)
	for _, p := range g.Base().GetPlayers() {
		p.Base().SetOnlineState(common.PlayerOnlineStateOnline)
	}
	uids := g.Base().UIDs()
	impl.connectorRpc.UpdateOnlineState(uids, common.PlayerOnlineStateOnline)
	impl.broadGroupState(g)

	// TODO: 需不需要从匹配队列中移除，还是仅通过状态管理即可？

	// 更新定时器
	impl.deleteQueueTimer(g.GroupID())
	impl.deleteWaitAttrTimer(g.GroupID())
	impl.addInviteTimer(g.GroupID())
	return nil
}

func (impl *Impl) UploadAttr(attr *pto.UploadAttr) error {
	p, g, err := impl.getPlayerAndGroup(attr.Uid)
	if err != nil {
		return err
	}
	g.Base().Lock()
	defer g.Base().Unlock()
	if err := g.CheckState(common.GroupStateInvite, common.GroupStateQueuing); err != nil {
		return err
	}
	p.Base().Lock()
	defer p.Base().Unlock()
	return p.SetAttr(attr)
}

func (impl *Impl) ChatInvite(param *ChatInvite, rsp *ChatInviteRsp) error {
	p, g, err := impl.getPlayerAndGroup(param.InviteeUid)
	if err != nil {
		return err
	}
	g.Base().Lock()
	defer g.Base().Unlock()
	if err := g.CheckState(common.GroupStateInvite); err != nil {
		return err
	}

	if param.ChatType == def.ChatTypeSingle {
		if err := impl.connectorRpc.CheckInviteFriend(&pto.CheckInviteFriend{
			InviteUid:   param.InviteeUid,
			FriendUid:   param.FriendUid,
			Platform:    p.Base().Platform,
			GameMode:    p.Base().GameMode,
			ModeVersion: p.Base().ModeVersion,
			Source:      pto.InvitationSrcSingleChat,
		}); err != nil {
			return err
		}
	}

	g.Base().InvitedPlayers[param.FriendUid] = true
	switch param.ChatType {
	case def.ChatTypeSingle:
		g.Base().InvitedPlayers[string(common.SrcSingleChat)] = true
	case def.ChatTypeChannel:
		g.Base().InvitedPlayers[string(common.SrcChannelChat)] = true
	case def.ChatTypeClan:
		g.Base().InvitedPlayers[string(common.SrcClanChat)] = true
	}

	rsp.GroupMemberCount = g.Base().PlayerCount()
	return nil
}

func (impl *Impl) BroadcastMessage(message *RpcChatMessage) {
	// TODO implement me
	panic("implement me")
}

func (impl *Impl) SetVoiceState(uid string, voiceState common.PlayerVoiceState) error {
	p, g, err := impl.getPlayerAndGroup(uid)
	if err != nil {
		return err
	}
	p.Base().SetVoiceState(voiceState)

	g.Base().Lock()
	defer g.Base().Unlock()
	if err := g.CheckState(common.GroupStateInvite, common.GroupStateQueuing); err != nil {
		return nil
	}
	impl.broadcastVoiceState(g)
	return nil
}

func (impl *Impl) JoinGroup(param *JoinGroup) error {
	// TODO implement me
	panic("implement me")
}

func (impl *Impl) SyncGroup(g *SyncGroup) error {
	// TODO implement me
	panic("implement me")
}

func (impl *Impl) SetNearbyJoinGroup(setting *UserSetting) error {
	_, g, err := impl.getPlayerAndGroup(setting.Uid)
	if err != nil {
		return err
	}
	g.Base().Lock()
	defer g.Base().Unlock()
	g.Base().NearbyJoinAllowed = setting.Setting
	return nil
}

func (impl *Impl) SetRecentJoinGroup(setting *UserSetting) error {
	_, g, err := impl.getPlayerAndGroup(setting.Uid)
	if err != nil {
		return err
	}
	g.Base().Lock()
	defer g.Base().Unlock()
	g.Base().RecentJoinAllowed = setting.Setting
	return nil
}

func (impl *Impl) RenewGroup(g common.Group) []string {
	g.Base().Lock()
	defer g.Base().Unlock()

	// 元信息
	g.InitUnReadyMap()
	g.Base().MatchID = ""

	// 状态
	g.Base().SetState(common.GroupStateInvite)
	for _, p := range g.Base().GetPlayers() {
		p.Base().SetOnlineState(common.PlayerOnlineStateGroup)
	}
	impl.broadGroupState(g)
	impl.connectorRpc.UpdateOnlineState(g.Base().UIDs(), common.PlayerOnlineStateGroup)

	// 定时器
	impl.deleteQueueTimer(g.GroupID())
	impl.deleteWaitAttrTimer(g.GroupID())
	impl.addInviteTimer(g.GroupID())
	impl.addGroupBroadcastTimer(g.GroupID())
	return g.Base().UIDs()
}

func (impl *Impl) GroupReady(uid string, groupID int, opt int) error {
	_, g, err := impl.getPlayerAndGroup(uid)
	if err != nil {
		return err
	}
	if g.GroupID() != int64(groupID) {
		return errors.New("group id not match")
	}
	g.Base().Lock()
	defer g.Base().Unlock()

	if err := g.CheckState(common.GroupStateInvite); err != nil {
		return err
	}
	delete(g.Base().UnReadyUsers, uid)
	impl.addGroupBroadcastTimer(g.GroupID())
	return nil
}

func (impl *Impl) ChangeRole(uid string, groupID int, role int) error {
	_, g, err := impl.getPlayerAndGroup(uid)
	if err != nil {
		return err
	}
	if g.GroupID() != int64(groupID) {
		return errors.New("group id not match")
	}
	g.Base().Lock()
	defer g.Base().Unlock()

	if err := g.CheckState(common.GroupStateInvite); err != nil {
		return err
	}
	g.Base().Roles[uid] = role
	impl.addGroupBroadcastTimer(g.GroupID())
	return nil
}

// BroadcastGroupUsers 将队伍信息同步给玩家
func (impl *Impl) BroadcastGroupUsers(g common.Group) {
	state := g.Base().GetState()
	if state != common.GroupStateInvite && state != common.GroupStateQueuing {
		return
	}

	impl.connectorRpc.PushGroupUsers(g.Base().UIDs(), g.Base().GetGroupUsers())
	impl.broadGroupState(g)
}

func (impl *Impl) broadGroupState(g common.Group, cancelUID ...string) {
	uids := g.Base().UIDs()
	if len(uids) == 0 {
		return
	}
	if len(cancelUID) > 0 {
		impl.connectorRpc.PushGroupState(uids, g.GroupID(), g.Base().GetState(), g.Base().Name, cancelUID[0])
	} else {
		impl.connectorRpc.PushGroupState(uids, g.GroupID(), g.Base().GetState(), g.Base().Name, "")
	}
}

func (impl *Impl) addGroupBroadcastTimer(groupID int64) {
	fireTime := time.Now().UnixMilli() + impl.config.GroupBroadcastTimeoutMs
	impl.rtm.AddTimer(def.RdsTimerGroupBroadcast, int(groupID), int(fireTime))
}

func (impl *Impl) deleteGroupBroadcastTimer(groupID int64) {
	impl.rtm.DelTimer(def.RdsTimerGroupBroadcast, int(groupID))
}

func (impl *Impl) addInviteTimer(groupID int64) {
	fireTime := time.Now().UnixMilli() + impl.config.InviteTimeoutMs
	impl.rtm.AddTimer(def.RdsTimerGroupInvite, int(groupID), int(fireTime))
}

func (impl *Impl) deleteInviteTimer(groupID int64) {
	impl.rtm.DelTimer(def.RdsTimerGroupInvite, int(groupID))
}

func (impl *Impl) addQueueTimer(groupID int64) {
	fireTime := time.Now().UnixMilli() + impl.config.QueueTimeoutMs
	impl.rtm.AddTimer(def.RdsTimerGroupQueue, int(groupID), int(fireTime))
}

func (impl *Impl) deleteQueueTimer(groupID int64) {
	impl.rtm.DelTimer(def.RdsTimerGroupQueue, int(groupID))
}

func (impl *Impl) addWaitAttrTimer(groupID int64) {
	fireTime := time.Now().UnixMilli() + impl.config.WaitAttrTimeoutMs
	impl.rtm.AddTimer(def.RdsTimerGroupWaitAttr, int(groupID), int(fireTime))
}

func (impl *Impl) deleteWaitAttrTimer(groupID int64) {
	impl.rtm.DelTimer(def.RdsTimerGroupWaitAttr, int(groupID))
}

func (impl *Impl) broadcastVoiceState(g common.Group) {
	voiceStates := make([]*pto.UserVoiceState, 0, g.Base().PlayerCount())
	for _, p := range g.Base().GetPlayers() {
		userVoiceState := &pto.UserVoiceState{
			Uid:   p.UID(),
			State: int(p.Base().GetVoiceState()),
		}
		voiceStates = append(voiceStates, userVoiceState)
	}
	impl.connectorRpc.PushGroupVoiceState(g.Base().UIDs(), voiceStates)
}
