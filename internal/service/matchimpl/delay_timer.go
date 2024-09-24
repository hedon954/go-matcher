package matchimpl

import (
	"context"

	"github.com/rs/zerolog/log"

	"github.com/hedon954/go-matcher/internal/constant"
	"github.com/hedon954/go-matcher/internal/entry"
	"github.com/hedon954/go-matcher/pkg/timer"
)

const (
	// TimerOpTypeGroupInvite used to dissolve the group if it not starts game after delay.
	TimerOpTypeGroupInvite timer.OpType = "match:timer_group_invite"

	// TimerOpTypeGroupMatch if used to cancel match if the group not matched after delay.
	TimerOpTypeGroupMatch timer.OpType = "match:timer_group_match"

	// TimerOpTypeGroupWaitAttr used to wait for players to upload attributes after client clicks `StartMatch`.
	// If all players upload attributes, the group would start to match.
	// If timeout. the group would auto start to match.
	TimerOpTypeGroupWaitAttr timer.OpType = "match:timer_group_wait_attr" // nolint

	// TimeOpTypeClearRoom used to clear room in some unexpected cases like client do not settle game.
	// We use this optype to force clear the room info.
	TimeOpTypeClearRoom timer.OpType = "match:timer_clear_room"
)

func (impl *Impl) initDelayTimer() {
	impl.delayTimer.Register(TimerOpTypeGroupInvite, impl.inviteTimeoutHandler)
	impl.delayTimer.Register(TimerOpTypeGroupMatch, impl.matchTimeoutHandler)
	impl.delayTimer.Register(TimerOpTypeGroupWaitAttr, impl.waitAttrTimeoutHandler)
	impl.delayTimer.Register(TimeOpTypeClearRoom, impl.clearRoomTimeoutHandler)
}

func (impl *Impl) inviteTimeoutHandler(groupID int64) {
	g := impl.groupMgr.Get(groupID)
	if g != nil {
		log.Debug().
			Int64("group_id", g.Base().GroupID).
			Msg("invite timeout")
		g.Base().Lock()
		defer g.Base().Unlock()
		if err := impl.dissolveGroup(context.Background(), g); err != nil {
			log.Error().
				Int64("group_id", g.Base().GroupID).
				Any("group", g).
				Err(err).
				Msg("dissolve group error")
		}
	}
}

func (impl *Impl) matchTimeoutHandler(groupID int64) {
	g := impl.groupMgr.Get(groupID)
	if g != nil {
		log.Debug().Int64("group_id", g.Base().GroupID).Msg("match timeout")

		g.Base().Lock()
		defer g.Base().Unlock()
		if g.Base().GetState() == entry.GroupStateMatch {
			impl.cancelMatch(context.Background(), "", g)
			log.Debug().Int64("group_id", g.Base().GroupID).Msg("cancel match by timeout")
		}
	}
}

func (impl *Impl) waitAttrTimeoutHandler(groupID int64) {
	g := impl.groupMgr.Get(groupID)
	if g != nil {
		log.Debug().
			Int64("group_id", g.Base().GroupID).
			Msg("wait attr timeout")
		g.Base().Lock()
		defer g.Base().Unlock()
		if g.Base().GetState() == entry.GroupStateMatch {
			impl.sendGroupToChannel(g)
		}
	}
}

func (impl *Impl) clearRoomTimeoutHandler(roomID int64) {
	r := impl.roomMgr.Get(roomID)
	if r != nil {
		impl.roomMgr.Delete(roomID)
		log.Warn().
			Int64("room_id", roomID).
			Any("room_info", r).
			Msg("clear room timeout")
	} else {
		log.Debug().
			Int64("room_id", roomID).
			Msg("clear room timeout, but room not found")
	}
}

func (impl *Impl) addInviteTimer(groupID int64, mode constant.GameMode) {
	err := impl.delayTimer.Add(TimerOpTypeGroupInvite, groupID,
		impl.Configer.Get().DelayTimerConfig.InviteTimeout())
	if err != nil {
		log.Error().
			Int64("group_id", groupID).
			Int("mode", int(mode)).
			Err(err).
			Msg("add invite timer error")
	} else {
		log.Debug().
			Int64("group_id", groupID).
			Int("mode", int(mode)).
			Msg("add invite timer successfully")
	}
}

func (impl *Impl) removeInviteTimer(groupID int64) {
	_ = impl.delayTimer.Remove(TimerOpTypeGroupInvite, groupID)
	log.Debug().
		Int64("group_id", groupID).
		Msg("remove invite timer successfully")
}

func (impl *Impl) addCancelMatchTimer(groupID int64, mode constant.GameMode) {
	err := impl.delayTimer.Add(TimerOpTypeGroupMatch, groupID,
		impl.Configer.Get().DelayTimerConfig.MatchTimeout())
	if err != nil {
		log.Error().
			Int64("group_id", groupID).
			Int("mode", int(mode)).
			Err(err).
			Msg("add cancel match timer error")
	} else {
		log.Debug().
			Int64("group_id", groupID).
			Int("mode", int(mode)).
			Msg("add cancel match timer successfully")
	}
}

func (impl *Impl) removeCancelMatchTimer(groupID int64) {
	_ = impl.delayTimer.Remove(TimerOpTypeGroupMatch, groupID)
	log.Debug().
		Int64("group_id", groupID).
		Msg("remove cancel match timer successfully")
}

func (impl *Impl) addWaitAttrTimer(groupID int64, mode constant.GameMode) {
	err := impl.delayTimer.Add(TimerOpTypeGroupWaitAttr, groupID,
		impl.Configer.Get().DelayTimerConfig.WaitAttrTimeout())
	if err != nil {
		log.Error().
			Int64("group_id", groupID).
			Int("mode", int(mode)).
			Err(err).
			Msg("add wait attr timer error")
	} else {
		log.Debug().
			Int64("group_id", groupID).
			Int("mode", int(mode)).
			Msg("add wait attr timer successfully")
	}
}

func (impl *Impl) removeWaitAttrTimer(groupID int64) {
	_ = impl.delayTimer.Remove(TimerOpTypeGroupWaitAttr, groupID)
	log.Debug().
		Int64("group_id", groupID).
		Msg("remove wait attr timer successfully")
}

func (impl *Impl) addClearRoomTimer(roomID int64, mode constant.GameMode) {
	err := impl.delayTimer.Add(TimeOpTypeClearRoom, roomID,
		impl.Configer.Get().DelayTimerConfig.ClearRoomTimeout())
	if err != nil {
		log.Error().
			Int64("room_id", roomID).
			Int("mode", int(mode)).
			Err(err).
			Msg("add clear room timer error")
	} else {
		log.Debug().
			Int64("room_id", roomID).
			Int("mode", int(mode)).
			Msg("add clear room timer successfully")
	}
}

func (impl *Impl) removeClearRoomTimer(roomID int64) {
	_ = impl.delayTimer.Remove(TimeOpTypeClearRoom, roomID)
	log.Debug().
		Int64("room_id", roomID).
		Msg("remove clear room timer successfully")
}
