package matchimpl

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/hedon954/go-matcher/internal/entry"
)

func (impl *Impl) startMatch(g entry.Group) {
	base := g.Base()
	uids := base.UIDs()

	// update group state
	base.SetState(entry.GroupStateMatch)
	base.MatchID = uuid.NewString()
	impl.pushService.PushGroupState(uids, g.ID(), base.GetState())

	// update players state
	for _, p := range base.GetPlayers() {
		p.Base().Lock()
		p.Base().SetOnlineState(entry.PlayerOnlineStateInMatch)
		p.Base().Unlock()
	}
	impl.pushService.PushPlayerOnlineState(uids, entry.PlayerOnlineStateInMatch)

	impl.removeInviteTimer(g.ID())
	impl.addWaitAttrTimer(g.ID(), g.Base().GameMode)
	impl.addCancelMatchTimer(g.ID(), base.GameMode)
}

func (impl *Impl) sendGroupToChannel(g entry.Group) {
	fmt.Println("sendGroupToChannel: ", g)
	impl.groupChannel <- g
}
