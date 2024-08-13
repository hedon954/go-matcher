package matchimpl

import (
	"context"

	"github.com/google/uuid"

	"github.com/hedon954/go-matcher/internal/entry"
	"github.com/hedon954/go-matcher/internal/log"
)

func (impl *Impl) startMatch(ctx context.Context, g entry.Group) {
	base := g.Base()
	uids := base.UIDs()

	// update group state
	base.SetState(entry.GroupStateMatch)
	base.MatchID = uuid.NewString()
	impl.pushService.PushGroupState(ctx, uids, g.ID(), base.GetState())

	// update players state
	for _, p := range base.GetPlayers() {
		p.Base().Lock()
		p.Base().SetOnlineState(entry.PlayerOnlineStateInMatch)
		p.Base().SetMatchStrategy(base.MatchStrategy)
		p.Base().Unlock()
	}
	impl.pushService.PushPlayerOnlineState(ctx, uids, entry.PlayerOnlineStateInMatch)

	impl.removeInviteTimer(g.ID())
	impl.addWaitAttrTimer(g.ID(), g.Base().GameMode)
	impl.addCancelMatchTimer(g.ID(), base.GameMode)
}

func (impl *Impl) sendGroupToChannel(g entry.Group) {
	log.Debug().Int64("group_id", g.ID()).Msg("send group to channel")
	impl.groupChannel <- g
}
