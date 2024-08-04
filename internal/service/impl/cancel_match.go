package impl

import "github.com/hedon954/go-matcher/internal/entry"

func (impl *Impl) cancelMatch(cancelUID string, g entry.Group) {
	g.Base().SetState(entry.GroupStateInvite)

	// add dissolve group timer
	impl.connectorClient.PushCancelMatch(g.Base().UIDs(), cancelUID)
}
