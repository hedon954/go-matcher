package config

import (
	"sync"
)

type NacosLoader struct {
	sync.RWMutex
	c *MatchConfig
}

func (nl *NacosLoader) Get() *MatchConfig {
	nl.RLock()
	defer nl.RUnlock()
	return nl.c
}
