package matchimpl

import (
	"context"
	"fmt"

	"github.com/hedon954/go-matcher/internal/entry"
)

func (impl *Impl) waitForMatchResult() {
	for r := range impl.roomChannel {
		fmt.Println("new room: ", r.ID())
		impl.HandleMatchResult(r)
	}
}

func (impl *Impl) handleMatchResult(ctx context.Context, r entry.Room) (err error) {
	// add room to manager
	defer func() {
		if err == nil {
			fmt.Println("add room to manager: ", r.ID())
			impl.roomMgr.Add(r.ID(), r)
		}
	}()
	// add teams to managers
	for _, t := range r.Base().GetTeams() {
		fmt.Printf("add team %d to room: %d\n", t.ID(), r.ID())
		impl.teamMgr.Add(t.ID(), t)
	}

	// ---------------------------
	// some operations without AI
	// ---------------------------
	r.Base().FinishMatchSec = impl.nowFunc()
	impl.clearDelayTimer(r)
	impl.updateStateToGame(ctx, r)

	// ----------------------------
	// some operations may need AI
	// ----------------------------
	err = impl.fillRoomInfo(r)
	if err != nil {
		return err
	}
	impl.pushService.PushMatchInfo(ctx, r.Base().UIDs(), r.GetMatchInfo())
	impl.addClearRoomTimer(r.ID(), r.Base().GameMode)
	return nil
}

func (impl *Impl) fillRoomInfo(r entry.Room) (err error) {
	// dispatch a game server address
	r.Base().GameServerInfo, err = impl.gameServerDispatch.Dispatch(context.Background(), r.Base().GameMode, r.Base().ModeVersion)
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

func (impl *Impl) updateStateToGame(ctx context.Context, r entry.Room) {
	for _, t := range r.Base().GetTeams() {
		impl.updateTeamStateToGame(ctx, t)
	}
}

func (impl *Impl) updateTeamStateToGame(ctx context.Context, t entry.Team) {
	t.Base().Lock()
	defer t.Base().Unlock()
	for _, g := range t.Base().GetGroups() {
		impl.updateGroupStateToGame(ctx, g)
	}
}

func (impl *Impl) updateGroupStateToGame(ctx context.Context, g entry.Group) {
	g.Base().Lock()
	defer g.Base().Unlock()
	g.Base().SetState(entry.GroupStateGame)
	for _, p := range g.Base().GetPlayers() {
		p.Base().SetOnlineStateWithLock(entry.PlayerOnlineStateInGame)
	}
	impl.pushService.PushPlayerOnlineState(ctx, g.Base().UIDs(), entry.PlayerOnlineStateInGame)
}
