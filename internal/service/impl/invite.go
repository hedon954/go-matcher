package impl

import (
	"github.com/hedon954/go-matcher/internal/entry"
	"github.com/hedon954/go-matcher/internal/pto"
)

func (impl *Impl) checkInviteeState(inviteeUID string) error {
	invitee := impl.playerMgr.Get(inviteeUID)
	if invitee == nil {
		return nil
	}
	invitee.Base().Lock()
	defer invitee.Base().Unlock()

	if err := invitee.Base().CheckOnlineState(
		entry.PlayerOnlineStateOnline,
		entry.PlayerOnlineStateInGroup); err != nil {
		return err
	}
	return nil
}

func (impl *Impl) invite(inviter entry.Player, inviteeUID string, g entry.Group) {
	g.Base().AddInviteRecord(inviteeUID, impl.nowFunc())
	impl.connectorClient.PushInviteMsg(&pto.InviteMsg{
		InviterUID:  inviter.UID(),
		InviteeUID:  inviteeUID,
		Source:      pto.EnterGroupSourceTypeInvite,
		GameMode:    g.Base().GameMode,
		ModeVersion: g.Base().ModeVersion,
	})
}
