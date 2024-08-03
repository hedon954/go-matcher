package impl

import (
	"time"

	"github.com/hedon954/go-matcher/internal/entry"
	"github.com/hedon954/go-matcher/internal/merr"
	"github.com/hedon954/go-matcher/internal/pto"
	"github.com/hedon954/go-matcher/internal/repository"
	"github.com/hedon954/go-matcher/internal/rpc/rpcclient/connector"
)

// Impl implements a default service,
// in most cases, you don't need to implement your own service.
type Impl struct {
	playerMgr *repository.PlayerMgr
	groupMgr  *repository.GroupMgr

	connectorClient *connector.Client

	playerLimit int
	nowFunc     func() int64
}

type Option func(*Impl)

func WithNowFunc(f func() int64) Option {
	return func(impl *Impl) {
		impl.nowFunc = f
	}
}

func NewDefault(playerLimit int, options ...Option) *Impl {
	impl := &Impl{
		playerMgr:       repository.NewPlayerMgr(),
		groupMgr:        repository.NewGroupMgr(0), // TODO: confirm the groupIDStart
		connectorClient: connector.New(),           // TODO: DI
		playerLimit:     playerLimit,
		nowFunc:         time.Now().Unix,
	}

	for _, opt := range options {
		opt(impl)
	}
	return impl
}

func (impl *Impl) CreateGroup(param *pto.CreateGroup) (entry.Group, error) {
	p, err := impl.getPlayer(&param.PlayerInfo)
	if err != nil {
		return nil, err
	}

	p.Base().Lock()
	defer p.Base().Unlock()

	if err := p.Base().CheckOnlineState(entry.PlayerOnlineStateOnline, entry.PlayerOnlineStateInGroup); err != nil {
		return nil, err
	}

	g := impl.groupMgr.Get(p.Base().GroupID)
	if g == nil {
		// create a group
		g, err = impl.createGroup(param, p)
		if err != nil {
			return nil, err
		}
	} else {
		g.Base().Lock()
		defer g.Base().Unlock()

		// check game mode
		// if game mode is not the same, exits the group and create a new one
		// if game mode is the same, check if the player is the captain of the group
		//  if not, exits the group and create a new one
		//  if yes, return current group
		if g.GetCaptain() != p || g.Base().GameMode != param.GameMode {
			if g.Base().RemovePlayer(p) {
				impl.groupMgr.Delete(g.GroupID())
			}
			g, err = impl.createGroup(param, p)
			if err != nil {
				return nil, err
			}
		}
	}

	return g, nil
}

func (impl *Impl) EnterGroup(info *pto.EnterGroup, groupID int64) error {
	g := impl.groupMgr.Get(groupID)
	if g == nil {
		return merr.ErrGroupDissolved
	}

	g.Base().Lock()
	defer g.Base().Unlock()
	if err := g.Base().CheckState(entry.GroupStateInvite); err != nil {
		return err
	}

	p, err := impl.getPlayer(&info.PlayerInfo)
	if err != nil {
		return err
	}

	p.Base().Lock()
	defer p.Base().Unlock()
	if err := p.Base().CheckOnlineState(entry.PlayerOnlineStateOnline, entry.PlayerOnlineStateInGroup); err != nil {
		return err
	}

	// check source validation
	if err := impl.checkEnterSourceValidation(g, info.Source); err != nil {
		return err
	}

	// player is already in a group
	if p.Base().GroupID != 0 {
		// already in targeted group
		if p.Base().GroupID == groupID {
			// if p can not play together, should exit the origin group
			if err := g.CanPlayTogether(&info.PlayerInfo); err != nil {
				if err := impl.exitGroup(p, g); err != nil {
					return err
				}
			} else {
				// can play together, refresh the player info and  broadcast the group player infos
				p.Base().PlayerInfo = info.PlayerInfo
				impl.connectorClient.PushGroupUsers(g.Base().UIDs(), g.GetPlayerInfos())
				return nil
			}
		} else {
			// check if player can play together with the group's players
			if err := g.CanPlayTogether(&info.PlayerInfo); err != nil {
				return err
			}

			// not in targeted group, should exit the origin group
			originGroup := impl.groupMgr.Get(p.Base().GroupID)
			if originGroup != nil {
				originGroup.Base().Lock()
				defer originGroup.Base().Unlock()
				if err := impl.exitGroup(p, originGroup); err != nil {
					return err
				}
			}
		}
	}

	// check if player can play together with the group's players
	if err := g.CanPlayTogether(&info.PlayerInfo); err != nil {
		return err
	}

	// refresh the player info
	p.Base().PlayerInfo = info.PlayerInfo

	// enter the targeted group
	return impl.enterGroup(p, g)
}

func (impl *Impl) ExitGroup(uid string) error {
	p, g, err := impl.getPlayerAndGroup(uid)
	if err != nil {
		return err
	}

	p.Base().Lock()
	defer p.Base().Unlock()
	if err := p.Base().CheckOnlineState(entry.PlayerOnlineStateInGroup); err != nil {
		return err
	}

	g.Base().Lock()
	defer g.Base().Unlock()
	if !g.Base().PlayerExists(p.UID()) {
		return nil
	}

	if err := g.Base().CheckState(entry.GroupStateInvite); err != nil {
		return err
	}

	impl.playerMgr.Delete(p.UID())
	return impl.exitGroup(p, g)
}

func (impl *Impl) DissolveGroup(uid string) error {
	p, g, err := impl.getPlayerAndGroup(uid)
	if err != nil {
		return err
	}

	g.Base().Lock()
	defer g.Base().Unlock()

	if g.GetCaptain() != p {
		return merr.ErrOnlyCaptainCanDissolveGroup
	}

	if err := g.Base().CheckState(entry.GroupStateInvite); err != nil {
		return err
	}

	return impl.dissolveGroup(g)
}

func (impl *Impl) KickPlayer(captainUID, kickedUID string) error {
	if captainUID == kickedUID {
		return merr.ErrKickSelf
	}

	captain, g, err := impl.getPlayerAndGroup(captainUID)
	if err != nil {
		return err
	}
	kicked := impl.playerMgr.Get(kickedUID)
	if kicked == nil {
		return merr.ErrPlayerNotExists
	}

	g.Base().Lock()
	defer g.Base().Unlock()

	if g.GetCaptain() != captain {
		return merr.ErrOnlyCaptainCanKickPlayer
	}

	if err := g.Base().CheckState(entry.GroupStateInvite); err != nil {
		return err
	}

	if !g.Base().PlayerExists(kickedUID) {
		return merr.ErrPlayerNotInGroup
	}

	kicked.Base().Lock()
	defer kicked.Base().Unlock()

	return impl.kickPlayer(kicked, g)
}

func (impl *Impl) ChangeRole(captainUID, targetUID string, role entry.GroupRole) error {
	if captainUID == targetUID {
		return merr.ErrChangeSelfRole
	}

	if err := impl.checkRole(role); err != nil {
		return err
	}

	captain, g, err := impl.getPlayerAndGroup(captainUID)
	if err != nil {
		return err
	}
	target := impl.playerMgr.Get(targetUID)
	if target == nil {
		return merr.ErrPlayerNotExists
	}

	g.Base().Lock()
	defer g.Base().Unlock()

	if !g.Base().PlayerExists(targetUID) {
		return merr.ErrPlayerNotInGroup
	}

	if g.GetCaptain() != captain {
		return merr.ErrNotCaptain
	}

	if err := g.Base().CheckState(entry.GroupStateInvite); err != nil {
		return err
	}

	return impl.handoverCaptain(captain, target, g)
}

func (impl *Impl) SetNearbyJoinGroup(captainUID string, allow bool) error {
	p, g, err := impl.getPlayerAndGroup(captainUID)
	if err != nil {
		return err
	}

	g.Base().Lock()
	defer g.Base().Unlock()

	if g.GetCaptain() != p {
		return merr.ErrPermissionDeny
	}

	g.Base().SetAllowNearbyJoin(allow)
	return nil
}

func (impl *Impl) SetRecentJoinGroup(captainUID string, allow bool) error {
	p, g, err := impl.getPlayerAndGroup(captainUID)
	if err != nil {
		return err
	}

	g.Base().Lock()
	defer g.Base().Unlock()

	if g.GetCaptain() != p {
		return merr.ErrPermissionDeny
	}

	g.Base().SetAllowRecentJoin(allow)
	return nil
}

func (impl *Impl) Invite(inviterUID, inviteeUID string) error {
	if err := impl.checkInviteeState(inviteeUID); err != nil {
		return err
	}

	inviter, g, err := impl.getPlayerAndGroup(inviterUID)
	if err != nil {
		return err
	}

	g.Base().Lock()
	defer g.Base().Unlock()

	if err := g.Base().CheckState(entry.GroupStateInvite); err != nil {
		return err
	}

	if g.IsFull() {
		return merr.ErrGroupFull
	}

	// TODO: how to check if can play together?
	impl.invite(inviter, inviteeUID, g)
	return nil
}

func (impl *Impl) AcceptInvite(inviterUID string, inviteeInfo *pto.PlayerInfo, groupID int64) error {
	g := impl.groupMgr.Get(groupID)
	if g == nil {
		return merr.ErrGroupDissolved
	}

	invitee := impl.playerMgr.Get(inviteeInfo.UID)
	if invitee != nil {
		invitee.Base().Lock()
		defer invitee.Base().Unlock()
		if err := invitee.Base().CheckOnlineState(entry.PlayerOnlineStateOnline, entry.PlayerOnlineStateInGroup); err != nil {
			return err
		}
	}

	g.Base().Lock()
	defer g.Base().Unlock()

	if !g.Base().PlayerExists(inviterUID) {
		return merr.ErrInvitationExpired
	}

	if err := g.Base().CheckState(entry.GroupStateInvite); err != nil {
		return err
	}

	if err := g.CanPlayTogether(inviteeInfo); err != nil {
		return err
	}
	return nil
}

func (impl *Impl) RefuseInvite(inviterUID, inviteeUID string, groupID int64, refuseMsg string) error {
	const defaultRefuseMsg = "Sorry, I'm not available at the moment."

	g := impl.groupMgr.Get(groupID)
	if g == nil {
		return nil
	}

	g.Base().Lock()
	defer g.Base().Unlock()

	if g.Base().GetState() == entry.GroupStateDissolved {
		return nil
	}
	g.Base().DelInviteRecord(inviteeUID)
	if refuseMsg == "" {
		refuseMsg = defaultRefuseMsg
	}
	impl.connectorClient.PushHandleInvite(inviterUID, inviteeUID, false, refuseMsg)
	return nil
}

func (impl *Impl) StartMatch(captainUID string) error {
	// TODO implement me
	panic("implement me")
}

func (impl *Impl) CancelMatch(uid string) error {
	// TODO implement me
	panic("implement me")
}

// getPlayerAndGroup returns the player and group of the given uid.
func (impl *Impl) getPlayerAndGroup(uid string) (entry.Player, entry.Group, error) {
	p := impl.playerMgr.Get(uid)
	if p == nil {
		return nil, nil, merr.ErrPlayerNotExists
	}
	g := impl.groupMgr.Get(p.Base().GroupID)
	if g == nil {
		return nil, nil, merr.ErrGroupNotExists
	}
	return p, g, nil
}
