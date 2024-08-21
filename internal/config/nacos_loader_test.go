package config

import (
	"reflect"
	"testing"

	"github.com/r3labs/diff/v3"

	"github.com/hedon954/go-matcher/internal/constant"
	"github.com/hedon954/go-matcher/pkg/algorithm/glicko2"
	"github.com/hedon954/go-matcher/thirdparty"
)

var (
	namespaceID   string
	group         = "GO-MATCHER"
	dataID        = "match_config"
	addr          = "127.0.0.1"
	port          = uint64(8848)
	serverConfigs = []*NacosServerConfig{
		{
			Addr:        addr,
			Port:        port,
			ContextPath: "/nacos",
			Schema:      "http",
		},
	}

	defaultMC = &MatchConfig{
		GroupPlayerLimit: 2,
		MatchIntervalMs:  1000,
		Glicko2: map[constant.GameMode]*glicko2.QueueArgs{
			constant.GameModeGoatGame: {
				MatchTimeoutSec: 300,
				TeamPlayerLimit: 2,
				RoomTeamLimit:   2,
				MatchRanges:     make([]glicko2.MatchRange, 0),
			},
		},
		DelayTimerType: DelayTimerTypeNative,
		DelayTimerConfig: &DelayTimerConfig{
			InviteTimeoutMs:    300000,
			MatchTimeoutMs:     60000,
			WaitAttrTimeoutMs:  1,
			ClearRoomTimeoutMs: 1800000,
		},
	}
)

func TestMain(m *testing.M) {
	namespaceID = thirdparty.PrepareNacosConfig(addr, dataID, group, port, defaultMC)
	m.Run()
	thirdparty.ClearNacosConfig(namespaceID, addr, port)
}

func Test_NacosLoader_MatchConfig(t *testing.T) {
	loader := NewNacosLoader(namespaceID, group, dataID, serverConfigs)
	got := loader.Get()
	if !reflect.DeepEqual(defaultMC, got) {
		d, _ := diff.Diff(defaultMC, got)
		t.Errorf("diff: \n%+v", d)
	}
}
