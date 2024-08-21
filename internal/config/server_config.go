package config

import (
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
)

// ServerConfig defines the server config.
type ServerConfig struct {
	AsynqRedis       *RedisOpt            `yaml:"asynq_redis"`
	NacosNamespaceID string               `yaml:"nacos_namespace_id"`
	NacosServers     []*NacosServerConfig `yaml:"nacos_servers"`
}

type RedisOpt struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type NacosServerConfig struct {
	Addr        string `yaml:"addr"`
	Port        uint64 `yaml:"port"`
	GRPCPort    uint64 `yaml:"grpc_port"`
	ContextPath string `yaml:"context_path"`
	Schema      string `yaml:"schema"`
}

func ToNacosServerConfigs(scs []*NacosServerConfig) []constant.ServerConfig {
	serverConfigs := make([]constant.ServerConfig, len(scs))
	for i := 0; i < len(scs); i++ {
		schema := "http"
		if scs[i].Schema != "" {
			schema = scs[i].Schema
		}
		contextPath := "/nacos"
		if scs[i].ContextPath != "" {
			contextPath = scs[i].ContextPath
		}
		serverConfigs[i] = constant.ServerConfig{
			Scheme:      schema,
			ContextPath: contextPath,
			IpAddr:      scs[i].Addr,
			Port:        scs[i].Port,
			GrpcPort:    scs[i].GRPCPort,
		}
	}
	return serverConfigs
}
