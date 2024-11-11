package main

import (
	"github.com/hedon954/go-matcher/cmd"
	"github.com/hedon954/go-matcher/internal/api/apihttp"
	"github.com/hedon954/go-matcher/internal/config"
)

func main() {
	defer cmd.StopSafe()
	sc := config.NewFileLoader[config.ServerConfig](
		cmd.GetConfigPath("SERVER_CONFIG_PATH", "cmd/server_conf_tmp.yml"))
	mc := config.NewNacosLoader(
		sc.Get().NacosNamespaceID,
		"GO-MATCHER",
		"match_config",
		sc.Get().NacosServers)

	infra := apihttp.NewInfra(sc, mc)
	defer infra.Stop()
	infra.Start()
}
