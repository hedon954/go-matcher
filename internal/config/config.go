package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/hedon954/go-matcher/internal/constant"
	"github.com/hedon954/go-matcher/pkg/algorithm/glicko2"
)

type DelayTimerType string

const (
	DelayTimerTypeAsynq  DelayTimerType = "asynq"
	DelayTimerTypeNative DelayTimerType = "native"
)

type Config struct {
	GroupPlayerLimit int                `yaml:"group_player_limit"`
	MatchIntervalMs  int64              `yaml:"match_interval_ms"`
	Glicko2          *glicko2.QueueArgs `yaml:"glicko2"`
	AsynqRedis       *Redis             `yaml:"asynq_redis"`
	DelayTimerType   DelayTimerType     `yaml:"delay_timer_type"`
	DelayTimerConfig *DelayTimerConfig  `yaml:"delay_timer_config"`
}

// Load loads the config from file.
func Load(path string) *Config {
	bs, err := os.ReadFile(path)
	if err != nil {
		panic("read config file error: " + err.Error())
	}
	c := &Config{}
	err = yaml.Unmarshal(bs, c)
	if err != nil {
		panic("unmarshal config error: " + err.Error())
	}
	return c
}

func (c *Config) GetQueueArgs(_ constant.GameMode) *glicko2.QueueArgs {
	return c.Glicko2
}

func (c *Config) MatchInterval() time.Duration {
	return time.Duration(c.MatchIntervalMs) * time.Millisecond
}

func (c *Config) GetConfig(_ constant.GameMode) DelayTimerConfig {
	return *c.DelayTimerConfig
}
