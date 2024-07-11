package matcher

import (
	"matcher/common"
	"matcher/group"
	"matcher/merr"
	"matcher/player"
	"matcher/pto"
)

// TODO: 应该支持多种实现
// TODO: 先不考虑状态的流转，先把整个流程完成再说，后续尝试使用状态机来进行改进管理

type Impl struct {
	playerMgr *player.Manager
	groupMgr  *group.Manager
}

func New() *Impl {
	impl := &Impl{
		playerMgr: player.NewManager(),
		groupMgr:  group.NewManager(),
	}
	return impl
}

func (impl *Impl) BindRestore(i, s string) (state int, matchInfo interface{}) {
	// TODO implement me
	panic("implement me")
}

func (impl *Impl) CreateGroup(info *pto.PlayerInfo) (int64, error) {
	// 1. 如果玩家已经在组里，则尝试返回之前的组
	if g := impl.inOwnedGroup(info); g != nil {
		return g.GroupID(), nil
	}

	// 2. 创建新的组
	g, err := impl.groupMgr.CreateGroup(info)
	if err != nil {
		return 0, err
	}
	return g.GroupID(), nil
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
	g.Inner().Lock()
	defer g.Inner().Unlock()

	// 只有自己是房主的时候，才返回原来的房间
	if g.Inner().IsOwner(info.Uid) {
		return g
	}

	// 不是房主的时候，需要先退出之前的房间
	g.Inner().DelPlayer(info.Uid)
	return nil
}

func (impl *Impl) InviteFriend(param *InviteFriend) error {
	// 自己首先要在组里
	_, g, err := impl.getPlayerAndGroup(param.InviteUid)
	if err != nil {
		return err
	}

	g.Inner().Lock()
	defer g.Inner().Unlock()

	// 分享邀请提前返回
	if param.Source == pto.InvitationSrcShare {
		g.Inner().MarkSrcShare()
		return nil
	}

	// 如果被邀请者已经在组里了，直接返回
	invitee := impl.playerMgr.GetPlayer(param.FriendUid)
	if invitee != nil {
		if g.Inner().PlayerExists(invitee) {
			g.Inner().BroadcastUsers()
			return nil
		}
	}

	// 发送邀请
	g.AddInvitedPlayer(param.FriendUid)
	// TODO: m.pm.PushInviteFriend(param)
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

func (impl *Impl) HandleInvite(param *HandleInvite) error {
	// 判断队伍在不在
	inviter, g, err := impl.getPlayerAndGroup(param.InviteUid)
	if err != nil {
		return merr.ErrGroupDissolved
	}

	// 判断队伍状态是否正确
	g.Inner().Lock()
	defer g.Inner().Unlock()
	if err := g.CheckState(common.GroupStateInvite); err != nil {
		return err
	}

	// 判断邀请是否已经过期
	if err := g.CheckHandleInviteExpired(param.Player.Uid, param.SrcType); err != nil {
		return err
	}

	// 判断队伍人数是否已满
	if g.IsFull() {
		return merr.ErrGroupFull
	}

	// 处理请求
	if param.HandleType == pto.InviteHandleTypeRefuseMsg {
		// TODO: send refuse msg to inviter
		return nil
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

func (impl *Impl) acceptInvite(g common.Group, inviter, invitee common.Player, srcType pto.InvitationSrcType) {
	// 如果之前就在队伍里，则需要退出
	preG := impl.groupMgr.GetGroup(invitee.Base().GroupID)
	if preG != nil {
		preG.Inner().DelPlayer(invitee.UID())
	}

	// 加入队伍
	g.Inner().AddPlayer(invitee)
	g.Inner().BroadcastUsers()

	if srcType == pto.InvitationSrcNearBy {
		g.Inner().AddNearby(inviter.UID(), invitee.UID())
	}
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
	g.Inner().Lock()
	defer g.Inner().Unlock()
	if err := g.CheckState(common.GroupStateInvite); err != nil {
		return err
	}

	if !g.Inner().IsOwner(uid) {
		return merr.ErrPermissionDeny
	}

	kp := impl.playerMgr.GetPlayer(kickUid)
	if !g.Inner().PlayerExists(p) {
		return nil
	}

	impl.removePlayer(g, kp.UID())
	return nil
}

func (impl *Impl) removePlayer(g common.Group, uid string) {
	if !g.Inner().DelPlayer(uid) {
		g.Inner().BroadcastUsers()
	}
}

func (impl *Impl) ExitGroup(uid string) error {
	p, g, err := impl.getPlayerAndGroup(uid)
	if err != nil {
		return nil
	}

	g.Inner().Lock()
	defer g.Inner().Unlock()
	if err := g.CheckState(common.GroupStateInvite); err != nil {
		return err
	}

	impl.exitGroup(p, g)
	return nil
}

func (impl *Impl) exitGroup(p common.Player, g common.Group) {
	groupIsEmpty := g.Inner().DelPlayer(p.UID())
	if groupIsEmpty {
		impl.groupMgr.DelGroup(g.GroupID())
	}
	impl.playerMgr.DelPlayer(p.UID())
}

func (impl *Impl) DissolveGroup(s string) error {
	// TODO implement me
	panic("implement me")
}

func (impl *Impl) ExitGameGroup(c common.Player) {
	// TODO implement me
	panic("implement me")
}

func (impl *Impl) Match(info *pto.PlayerInfo, reply *MatchReply) error {
	// TODO implement me
	panic("implement me")
}

func (impl *Impl) CancelMatch(match *CancelMatch) error {
	// TODO implement me
	panic("implement me")
}

func (impl *Impl) MatchSuccess(c common.Group) {
	// TODO implement me
	panic("implement me")
}

func (impl *Impl) UploadAttr(attr *UploadAttr) error {
	// TODO implement me
	panic("implement me")
}

func (impl *Impl) ChatInvite(invite *ChatInvite, rsp *ChatInviteRsp) error {
	// TODO implement me
	panic("implement me")
}

func (impl *Impl) BroadcastMessage(message *RpcChatMessage) {
	// TODO implement me
	panic("implement me")
}

func (impl *Impl) SetVoiceState(s string, voice int) error {
	// TODO implement me
	panic("implement me")
}

func (impl *Impl) JoinGroup(joinGroup *JoinGroup) error {
	// TODO implement me
	panic("implement me")
}

func (impl *Impl) SyncGroup(g *SyncGroup) error {
	// TODO implement me
	panic("implement me")
}

func (impl *Impl) SetNearbyJoinGroup(setting *UserSetting) error {
	// TODO implement me
	panic("implement me")
}

func (impl *Impl) SetRecentJoinGroup(setting *UserSetting) error {
	// TODO implement me
	panic("implement me")
}

func (impl *Impl) RenewGroup(c common.Group) []string {
	// TODO implement me
	panic("implement me")
}
