package matchimpl

import (
	"github.com/hedon954/go-matcher/internal/entry"
)

func (impl *Impl) waitForMatchResult() {
	for r := range impl.roomChannel {
		impl.HandleMatchResult(r)
	}
}

func (impl *Impl) handleMatchResult(r entry.Room) error {
	// ---------------------------
	// some operations without AI
	// ---------------------------
	r.Base().FinishMatchSec = impl.nowFunc()
	impl.clearDelayTimer(r)
	impl.updateStateToGame(r)

	// ----------------------------
	// some operations may need AI
	// ----------------------------
	if err := impl.fillRoomInfo(r); err != nil {
		return err
	}
	impl.pushService.PushMatchInfo(r.Base().UIDs(), r.GetMatchInfo())
	impl.addClearRoomTimer(r.ID(), r.Base().GameMode)
	return nil
}

func (impl *Impl) fillRoomInfo(r entry.Room) (err error) {
	// dispatch a game server address
	r.Base().GameServerInfo, err = impl.gameServerDispatch.Dispatch(r.Base().GameMode, r.Base().ModeVersion)
	if err != nil {
		return err
	}

	// fill room with AI
	if err := impl.fillRoomWithAI(r); err != nil {
		return err
	}

	// shuffle camp order randomly
	r.Base().ShuffleTeamOrder()
	return nil
}

func (impl *Impl) fillRoomWithAI(r entry.Room) error {
	if !r.NeedAI() {
		return nil
	}

	for _, t := range r.Base().GetTeams() {
		if err := impl.fillTeamWithAI(t); err != nil {
			return err
		}
	}

	teamCount := len(r.Base().GetTeams())
	for i := teamCount; i < r.Base().TeamLimit; i++ {
		r.Base().AddTeam(impl.createAITeam())
	}

	return nil
}

func (impl *Impl) fillTeamWithAI(t entry.Team) error {
	// TODO implement AI generator
	return nil
}

func (impl *Impl) createAITeam() entry.Team {
	// TODO implement AI generator
	return nil
}

func (impl *Impl) clearDelayTimer(r entry.Room) {
	for _, t := range r.Base().GetTeams() {
		t.Base().Lock()
		for _, g := range t.Base().GetGroups() {
			impl.removeWaitAttrTimer(g.ID())
			impl.removeWaitAttrTimer(g.ID())
			impl.removeCancelMatchTimer(g.ID())
		}
		t.Base().Unlock()
	}
}

func (impl *Impl) updateStateToGame(r entry.Room) {
	for _, t := range r.Base().GetTeams() {
		impl.updateTeamStateToGame(t)
	}
}

func (impl *Impl) updateTeamStateToGame(t entry.Team) {
	t.Base().Lock()
	defer t.Base().Unlock()
	for _, g := range t.Base().GetGroups() {
		impl.updateGroupStateToGame(g)
	}
}

func (impl *Impl) updateGroupStateToGame(g entry.Group) {
	g.Base().Lock()
	defer g.Base().Unlock()
	g.Base().SetState(entry.GroupStateGame)
	for _, p := range g.Base().GetPlayers() {
		p.Base().SetOnlineStateWithLock(entry.PlayerOnlineStateInGame)
	}
	impl.pushService.PushPlayerOnlineState(g.Base().UIDs(), entry.PlayerOnlineStateInGame)
}
