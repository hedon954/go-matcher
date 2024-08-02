package impl

import (
	"github.com/hedon954/go-matcher/internal/entry"
	"github.com/hedon954/go-matcher/internal/merr"
	"github.com/hedon954/go-matcher/internal/pto"
)

func (impl *Impl) checkEnterSourceValidation(g entry.Group, source pto.InvitationSrcType) error {
	switch source {
	case pto.InvitationSrcNearBy:
		if !g.Base().AllowNearbyJoin() {
			return merr.ErrGroupDenyNearbyJoin
		}
	case pto.InvitationSrcRecent:
		if !g.Base().AllowRecentJoin() {
			return merr.ErrGroupDenyRecentJoin
		}
		// TODO: check other
	}

	return nil
}

func (impl *Impl) enterGroup(p entry.Player, g entry.Group) error {
	if err := g.Base().AddPlayer(p); err != nil {
		return err
	}
	p.Base().GroupID = g.GroupID()
	p.Base().SetOnlineState(entry.PlayerOnlineStateInGroup)
	impl.connectorClient.UpdateOnlineState([]string{p.UID()}, int(entry.PlayerOnlineStateInGroup))
	impl.connectorClient.PushGroupUsers(g.Base().UIDs(), g.GetPlayerInfos())
	return nil
}
