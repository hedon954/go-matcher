package config

import (
	"reflect"
	"testing"

	"github.com/r3labs/diff/v3"

	"github.com/hedon954/go-matcher/internal/constant"
	"github.com/hedon954/go-matcher/pkg/algorithm/glicko2"
)

func Test_FileLoader_MatchConfig(t *testing.T) {
	loader := NewFileLoader[MatchConfig]("../../fixtures/match_conf_test.yml")
	expected := &MatchConfig{
		GroupPlayerLimit: 2,
		MatchIntervalMs:  1000,
		Glicko2: map[constant.GameMode]*glicko2.QueueArgs{
			constant.GameModeGoatGame: {
				MatchTimeoutSec: 300,
				TeamPlayerLimit: 2,
				RoomTeamLimit:   2,
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

	got := loader.Get()
	if !reflect.DeepEqual(expected, got) {
		d, _ := diff.Diff(expected, got)
		t.Errorf("diff: \n%+v", d)
	}
}

func Test_FileLoader_ServerConfig(t *testing.T) {
	loader := NewFileLoader[ServerConfig]("../../fixtures/server_conf_test.yml")

	expected := &ServerConfig{
		AsynqRedis: &RedisOpt{
			Addr:     "127.0.0.1:6379",
			Password: "",
			DB:       0,
		},
		NacosServers: []*NacosServerConfig{
			{
				Addr:        "127.0.0.1",
				Port:        8848,
				ContextPath: "/nacos",
				Schema:      "http",
			},
		},
		NacosNamespaceID: "7d638262-9e51-4822-9333-c3bcca838b7d",
	}

	got := loader.Get()

	if !reflect.DeepEqual(expected, got) {
		d, _ := diff.Diff(expected, got)
		t.Errorf("diff: \n%+v", d)
	}
}
