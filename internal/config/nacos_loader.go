package config

import (
	"fmt"
	"sync"

	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"gopkg.in/yaml.v3"

	"github.com/hedon954/go-matcher/internal/log"
	"github.com/hedon954/go-matcher/thirdparty"
)

type NacosLoader struct {
	group       string
	dataID      string
	namespaceID string
	nacosClient config_client.IConfigClient

	sync.RWMutex
	mc *MatchConfig
}

func NewNacosLoader(namespaceID, group, dataID string, scs []*NacosServerConfig) *NacosLoader {
	nl := &NacosLoader{
		group:       group,
		dataID:      dataID,
		namespaceID: namespaceID,
		nacosClient: newNacosClient(namespaceID, scs)}
	nl.load()
	return nl
}

func (nl *NacosLoader) Get() *MatchConfig {
	nl.RLock()
	defer nl.RUnlock()
	return nl.mc
}

func (nl *NacosLoader) load() {
	nl.loadMatchConfig()
}

func (nl *NacosLoader) loadMatchConfig() {
	mc, err := nl.nacosClient.GetConfig(vo.ConfigParam{
		DataId:   nl.dataID,
		Group:    nl.group,
		OnChange: nl.updateMatchConfig,
	})
	if err != nil {
		panic(fmt.Errorf("get nacos match config error: %w", err))
	}

	nl.updateMatchConfig(nl.namespaceID, nl.dataID, nl.group, mc)
	if nl.Get().DelayTimerType == "" {
		panic("load match config from nacos failed")
	}
}

func (nl *NacosLoader) updateMatchConfig(namespace, group, dataID, data string) {
	mc := &MatchConfig{}
	if err := yaml.Unmarshal([]byte(data), mc); err != nil {
		log.Error().
			Str("namespace", namespace).
			Str("group", group).
			Str("data_id", dataID).
			Str("data", data).
			Err(err).
			Msg("unmarshal nacos basic config error when config update")
		return
	}

	nl.Lock()
	defer nl.Unlock()
	nl.mc = mc
}

func newNacosClient(namespaceID string, scs []*NacosServerConfig) config_client.IConfigClient {
	nc, err := thirdparty.NewNacosConfigClient(namespaceID, ToNacosServerConfigs(scs))
	if err != nil {
		panic(fmt.Errorf("new nacos client error: %w", err))
	}
	return nc
}
