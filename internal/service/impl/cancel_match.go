package impl

import "github.com/hedon954/go-matcher/internal/entry"

func (impl *Impl) cancelMatch(cancelUID string, g entry.Group) {
	base := g.Base()

	base.SetState(entry.GroupStateInvite)
	base.MatchID = ""

	uids := base.UIDs()
	for _, uid := range uids {
		p := impl.playerMgr.Get(uid)
		p.Base().Lock()
		p.Base().SetOnlineState(entry.PlayerOnlineStateInGroup)
		p.Base().Unlock()
	}
	impl.connectorClient.PushGroupState(uids, g.ID(), base.GetState())
	impl.connectorClient.PushCancelMatch(base.UIDs(), cancelUID)

	impl.removeCancelMatchTimer(g.ID())
	impl.removeWaitAttrTimer(g.ID())
	impl.addInviteTimer(g.ID(), base.GameMode)
}
