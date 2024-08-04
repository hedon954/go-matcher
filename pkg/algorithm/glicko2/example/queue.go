package example

import "github.com/hedon954/go-matcher/pkg/algorithm/glicko2"

func GetQueueArgs() *glicko2.QueueArgs {
	return &glicko2.QueueArgs{
		TeamPlayerLimit:              TeamPlayerLimit,
		RoomTeamLimit:                RoomTeamLimit,
		UnfriendlyTeamMMRVarianceMin: UnfriendlyTeamVarianceMin,
		MaliciousTeamMMRVarianceMin:  MaliciousTeamVarianceMin,
		NormalTeamWaitTimeSec:        NormalTeamWaitTimeSec,
		UnfriendlyTeamWaitTimeSec:    UnfriendlyTeamWaitTimeSec,
		MaliciousTeamWaitTimeSec:     MaliciousTeamWaitTimeSec,
		MatchRanges: []glicko2.MatchRange{
			{
				MaxMatchSec:   15,
				MMRGapPercent: 10,
				CanJoinTeam:   false,
				StarGap:       0,
			},
			{
				MaxMatchSec:   30,
				MMRGapPercent: 20,
				CanJoinTeam:   false,
				StarGap:       0,
			},
			{
				MaxMatchSec:   60,
				MMRGapPercent: 40,
				CanJoinTeam:   false,
				StarGap:       0,
			},
			{
				MaxMatchSec:   90,
				MMRGapPercent: 60,
				CanJoinTeam:   true,
				StarGap:       0,
			},
			{
				MaxMatchSec:   120,
				MMRGapPercent: 60,
				CanJoinTeam:   true,
				StarGap:       0,
			},
		},
	}
}
