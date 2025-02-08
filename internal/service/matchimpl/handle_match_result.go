package matchimpl

import (
	"context"
	"fmt"

	"github.com/hedon954/go-matcher/internal/entry"
	"github.com/hedon954/go-matcher/internal/matcher/common"
)

func (impl *Impl) waitForMatchResult() {
	for r := range impl.roomChannel {
		fmt.Println("new room: ", r.Room.ID())
		impl.HandleMatchResult(r)
	}
}

func (impl *Impl) handleMatchResult(ctx context.Context, result common.Result) (err error) {
	r := result.Room
	teams := result.Teams

	// add room to manager
	defer func() {
		if err == nil {
			fmt.Println("add room to manager: ", r.ID())
			impl.roomMgr.Add(r.ID(), r)
		}
	}()
	// add teams to managers
	for _, team := range teams {
		fmt.Printf("add team %d to room: %d\n", team.ID(), r.ID())
		impl.teamMgr.Add(team.ID(), team)
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
	impl.pushService.PushMatchInfo(ctx, impl.getRoomUIDs(r), r.GetMatchInfo())
	impl.addClearRoomTimer(r.ID(), r.Base().GameMode)
	return nil
}

func (impl *Impl) getRoomUIDs(r entry.Room) []string {
	res := make([]string, 0)
	for _, teamID := range r.Base().GetTeams() {
		t := impl.teamMgr.Get(teamID)
		for _, groupID := range t.Base().GetGroups() {
			res = append(res, impl.groupMgr.Get(groupID).Base().UIDs()...)
		}
	}
	return res
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

	return nil
}

func (impl *Impl) fillRoomWithAI(r entry.Room) error {
	if !r.NeedAI() {
		return nil
	}

	for _, teamID := range r.Base().GetTeams() {
		t := impl.teamMgr.Get(teamID)
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

func (impl *Impl) fillTeamWithAI(_ entry.Team) error {
	// TODO implement AI generator
	return nil
}

func (impl *Impl) createAITeam() entry.Team {
	// TODO implement AI generator
	return nil
}

func (impl *Impl) clearDelayTimer(r entry.Room) {
	for _, teamID := range r.Base().GetTeams() {
		t := impl.teamMgr.Get(teamID)
		t.Base().Lock()
		for _, groupID := range t.Base().GetGroups() {
			impl.removeWaitAttrTimer(groupID)
			impl.removeWaitAttrTimer(groupID)
			impl.removeCancelMatchTimer(groupID)
		}
		t.Base().Unlock()
	}
}

func (impl *Impl) updateStateToGame(ctx context.Context, r entry.Room) {
	for _, teamID := range r.Base().GetTeams() {
		t := impl.teamMgr.Get(teamID)
		impl.updateTeamStateToGame(ctx, t)
	}
}

func (impl *Impl) updateTeamStateToGame(ctx context.Context, t entry.Team) {
	t.Base().Lock()
	defer t.Base().Unlock()
	for _, groupID := range t.Base().GetGroups() {
		impl.updateGroupStateToGame(ctx, impl.groupMgr.Get(groupID))
	}
}

func (impl *Impl) updateGroupStateToGame(ctx context.Context, g entry.Group) {
	g.Base().Lock()
	defer g.Base().Unlock()
	g.Base().SetState(entry.GroupStateGame)
	for _, puid := range g.Base().GetPlayers() {
		p := impl.playerMgr.Get(puid)
		p.Base().SetOnlineStateWithLock(entry.PlayerOnlineStateInGame)
	}
	impl.pushService.PushPlayerOnlineState(ctx, g.Base().UIDs(), entry.PlayerOnlineStateInGame)
}
