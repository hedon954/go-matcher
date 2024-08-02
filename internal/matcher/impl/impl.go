package impl

import (
	"github.com/hedon954/go-matcher/internal/entry"
	"github.com/hedon954/go-matcher/internal/manager"
	"github.com/hedon954/go-matcher/internal/merr"
	"github.com/hedon954/go-matcher/internal/pto"
	"github.com/hedon954/go-matcher/internal/rpc/rpcclient/connector"
)

// Impl implements a default matcher,
// in most cases, you don't need to implement your own matcher.
type Impl struct {
	playerMgr *manager.PlayerMgr
	groupMgr  *manager.GroupMgr

	connectorClient *connector.Client

	playerLimit int
}

func NewDefault(playerLimit int) *Impl {
	impl := &Impl{
		playerMgr:       manager.NewPlayerMgr(),
		groupMgr:        manager.NewGroupMgr(0), // TODO: confirm the groupIDStart
		connectorClient: connector.New(),        // TODO: DI
		playerLimit:     playerLimit,
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

	if err = p.Base().CheckOnlineState(entry.PlayerOnlineStateOnline, entry.PlayerOnlineStateInGroup); err != nil {
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

func (impl *Impl) EnterGroup(info *pto.PlayerInfo, groupID int64) error {
	g := impl.groupMgr.Get(groupID)
	if g == nil {
		return merr.ErrGroupDissolved
	}

	g.Base().Lock()
	defer g.Base().Unlock()
	if err := g.Base().CheckState(entry.GroupStateInvite); err != nil {
		return err
	}

	p, err := impl.getPlayer(info)
	if err != nil {
		return err
	}

	p.Base().Lock()
	defer p.Base().Unlock()
	if err = p.Base().CheckOnlineState(entry.PlayerOnlineStateOnline, entry.PlayerOnlineStateInGroup); err != nil {
		return err
	}

	// player is already in a group
	if p.Base().GroupID != 0 {
		// already in targeted group
		if p.Base().GroupID == groupID {
			// if p can not play together, should exit the origin group
			if !g.CanPlayTogether(p) {
				if err = impl.exitGroup(p, g); err != nil {
					return err
				}
			} else {
				// can play together, just broadcast the group player infos
				impl.connectorClient.PushGroupUsers(g.Base().UIDs(), g.GetPlayerInfos())
				return nil
			}
		} else {
			// check if player can play together with the group's players
			if !g.CanPlayTogether(p) {
				return merr.ErrVersionNotMatch
			}

			// not in targeted group, should exit the origin group
			originGroup := impl.groupMgr.Get(p.Base().GroupID)
			if originGroup != nil {
				originGroup.Base().Lock()
				defer originGroup.Base().Unlock()
				if err = impl.exitGroup(p, originGroup); err != nil {
					return err
				}
			}
		}
	}

	// check if player can play together with the group's players
	if !g.CanPlayTogether(p) {
		return merr.ErrVersionNotMatch
	}

	// enter the targeted group
	return impl.enterGroup(p, g)
}

func (impl *Impl) ExitGroup(uid string) error {
	p := impl.playerMgr.Get(uid)
	if p == nil {
		return merr.ErrPlayerNotInGroup
	}

	g := impl.groupMgr.Get(p.Base().GroupID)
	if g == nil {
		return merr.ErrPlayerNotInGroup
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

func (impl *Impl) Invite(inviterUID, inviteeUID string) error {
	// TODO implement me
	panic("implement me")
}

func (impl *Impl) AcceptInvite(inviteeUID string, groupID int64) error {
	// TODO implement me
	panic("implement me")
}

func (impl *Impl) RefuseInvite(inviteeUID string, groupID int64, refuseMsg string) error {
	// TODO implement me
	panic("implement me")
}

func (impl *Impl) CancelMatch(uid string) error {
	// TODO implement me
	panic("implement me")
}

func (impl *Impl) ReadyToMatch(uid string) error {
	// TODO implement me
	panic("implement me")
}

func (impl *Impl) DissolveGroup(uid string) error {
	p := impl.playerMgr.Get(uid)
	if p == nil {
		return merr.ErrPlayerNotInGroup
	}

	g := impl.groupMgr.Get(p.Base().GroupID)
	if g == nil {
		return merr.ErrGroupNotExists
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

	captain := impl.playerMgr.Get(captainUID)
	if captain == nil {
		return merr.ErrPlayerNotExists
	}
	kicked := impl.playerMgr.Get(kickedUID)
	if kicked == nil {
		return merr.ErrPlayerNotExists
	}
	g := impl.groupMgr.Get(captain.Base().GroupID)
	if g == nil {
		return merr.ErrGroupNotExists
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

func (impl *Impl) StartMatch(captainUID string) error {
	// TODO implement me
	panic("implement me")
}

func (impl *Impl) HandoverCaptain(captainUID, targetUID string) error {
	if captainUID == targetUID {
		return merr.ErrHandoverSelf
	}

	captain := impl.playerMgr.Get(captainUID)
	if captain == nil {
		return merr.ErrPlayerNotExists
	}
	target := impl.playerMgr.Get(targetUID)
	if target == nil {
		return merr.ErrPlayerNotExists
	}
	g := impl.groupMgr.Get(captain.Base().GroupID)
	if g == nil {
		return merr.ErrGroupNotExists
	}

	g.Base().Lock()
	defer g.Base().Unlock()

	if g.GetCaptain() != captain {
		return merr.ErrNotCaptain
	}

	if err := g.Base().CheckState(entry.GroupStateInvite); err != nil {
		return err
	}

	return impl.handoverCaptain(captain, target, g)
}
