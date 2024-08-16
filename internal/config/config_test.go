package config

import (
	"reflect"
	"testing"

	"github.com/hedon954/go-matcher/pkg/algorithm/glicko2"
)

func TestLoad(t *testing.T) {
	conf := Load("../../fixtures/conf_test.yml")
	reflect.DeepEqual(&Config{
		GroupPlayerLimit: 2,
		MatchIntervalMs:  1000,
		Glicko2: &glicko2.QueueArgs{
			MatchTimeoutSec: 300,
			TeamPlayerLimit: 2,
			RoomTeamLimit:   2,
		},
		AsynqRedis: &Redis{
			Addr:     "127.0.0.1:6379",
			Password: "",
			DB:       0,
		},
		DelayTimerType: DelayTimerTypeNative,
		DelayTimerConfig: &DelayTimerConfig{
			InviteTimeoutMs:    300000,
			MatchTimeoutMs:     60000,
			WaitAttrTimeoutMs:  1,
			ClearRoomTimeoutMs: 1800000,
		},
	}, conf)
}
