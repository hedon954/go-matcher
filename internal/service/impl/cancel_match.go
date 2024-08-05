package impl

import "github.com/hedon954/go-matcher/internal/entry"

func (impl *Impl) cancelMatch(cancelUID string, g entry.Group) {
	base := g.Base()

	base.SetState(entry.GroupStateInvite)
	base.MatchID = ""

	uids := base.UIDs()
	impl.connectorClient.PushGroupState(uids, g.ID(), base.GetState())
	impl.connectorClient.PushCancelMatch(base.UIDs(), cancelUID)
	// TODO: add dissolve group timer
}
