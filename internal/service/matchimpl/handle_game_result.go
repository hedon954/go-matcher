package matchimpl

import (
	"slices"

	"github.com/hedon954/go-matcher/internal/constant"
	"github.com/hedon954/go-matcher/internal/entry"
	"github.com/hedon954/go-matcher/internal/log"
	"github.com/hedon954/go-matcher/internal/pto"
)

func (impl *Impl) handleGameResult(result *pto.GameResult) {
	impl.result[result.RoomID] = result
	impl.removeClearRoomTimer(result.RoomID)

	log.Info().
		Int64("room_id", result.RoomID).
		Int("game_mode", int(result.GameMode)).
		Int64("mode_version", result.ModeVersion).
		Int("match_strategy", int(result.MatchStrategy)).
		Any("player_meta_infos", result.PlayerMetaInfo).
		Msg("handle game result")

	r := impl.roomMgr.Get(result.RoomID)
	if r == nil {
		log.Error().
			Int64("room_id", result.RoomID).
			Any("result", result).
			Msg("can not find room when handle game result")
		return
	}

	r.Base().Lock()
	defer r.Base().Unlock()

	escapePlayers := r.Base().GetEscapePlayers()
	impl.updateStateToSettle(r, escapePlayers)
	impl.clearMatchStrategy(r, escapePlayers) // do not worry about performance, just make it readable

	// ... do something to punish escape players
	// ... do something to handle result
}

func (impl *Impl) updateStateToSettle(r entry.Room, escapePlayers []string) {
	for _, t := range r.Base().GetTeams() {
		impl.updateTeamStateToSettle(t, escapePlayers)
	}
}

func (impl *Impl) updateTeamStateToSettle(team entry.Team, escapePlayers []string) {
	team.Base().Lock()
	defer team.Base().Unlock()
	for _, g := range team.Base().GetGroups() {
		impl.updateGroupStateToSettle(g, escapePlayers)
	}
}

func (impl *Impl) updateGroupStateToSettle(g entry.Group, escapePlayers []string) {
	g.Base().Lock()
	defer g.Base().Unlock()
	g.Base().SetState(entry.GroupStateInvite)
	for _, p := range g.Base().GetPlayers() {
		if slices.Index(escapePlayers, p.UID()) < 0 {
			continue
		}
		p.Base().SetOnlineStateWithLock(entry.PlayerOnlineStateInSettle)
		p.Base().SetMatchStrategyWithLock(constant.MatchStrategy(0))
	}
}

func (impl *Impl) clearMatchStrategy(r entry.Room, escapePlayers []string) {
	for _, t := range r.Base().GetTeams() {
		impl.clearTeamMatchStrategy(t, escapePlayers)
	}
}

func (impl *Impl) clearTeamMatchStrategy(team entry.Team, escapePlayers []string) {
	team.Base().Lock()
	defer team.Base().Unlock()
	for _, g := range team.Base().GetGroups() {
		impl.clearGroupMatchStrategy(g, escapePlayers)
	}
}

func (impl *Impl) clearGroupMatchStrategy(g entry.Group, escapePlayers []string) {
	g.Base().Lock()
	defer g.Base().Unlock()
	for _, p := range g.Base().GetPlayers() {
		if slices.Index(escapePlayers, p.UID()) < 0 {
			continue
		}
		p.Base().SetMatchStrategyWithLock(constant.MatchStrategy(0))
	}
}
