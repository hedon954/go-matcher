package mock

import (
	"github.com/hedon954/go-matcher/internal/config"
)

type ServerConfigerMock struct {
	sc *config.ServerConfig
}

//nolint:all
func NewServerConfigerMock() *ServerConfigerMock {
	return &ServerConfigerMock{sc: &config.ServerConfig{
		AsynqRedis: &config.RedisOpt{
			Addr:     "127.0.0.1:6379",
			Password: "",
			DB:       0,
		},
		NacosServers: []*config.NacosServerConfig{
			{
				Addr:        "127.0.0.1",
				Port:        8848,
				ContextPath: "/nacos",
				Schema:      "http",
			},
		},
	}}
}

func (sc *ServerConfigerMock) Get() *config.ServerConfig {
	return sc.sc
}
