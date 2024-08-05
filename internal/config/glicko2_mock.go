package config

import (
	"github.com/hedon954/go-matcher/internal/constant"
	"github.com/hedon954/go-matcher/pkg/algorithm/glicko2"
)

type Glicko2Mock struct {
}

func (gc *Glicko2Mock) GetQueueArgs(mode constant.GameMode) *glicko2.QueueArgs {
	if mode == constant.GameModeGoatGame {
		return &glicko2.QueueArgs{
			MatchTimeoutSec:              50,
			TeamPlayerLimit:              1,
			RoomTeamLimit:                3,
			NewerWithNewer:               false,
			UnfriendlyTeamMMRVarianceMin: 0,
			MaliciousTeamMMRVarianceMin:  0,
			NormalTeamWaitTimeSec:        0,
			UnfriendlyTeamWaitTimeSec:    0,
			MaliciousTeamWaitTimeSec:     0,
			MatchRanges: []glicko2.MatchRange{
				{
					MaxMatchSec:   0,
					MMRGapPercent: 0,
					CanJoinTeam:   true,
					StarGap:       0,
				},
			},
		}
	}
	return nil
}
