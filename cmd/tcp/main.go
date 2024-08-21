package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/hedon954/go-matcher/cmd"
	"github.com/hedon954/go-matcher/internal/api/apitcp"
	"github.com/hedon954/go-matcher/internal/config"
	"github.com/hedon954/go-matcher/pkg/zinx/zconfig"
)

func main() {
	defer cmd.StopSafe()
	mc := config.NewFileLoader[config.ServerConfig]("cmd/server_conf_tmp.yml")
	_, _, shutdown := apitcp.SetupTCPServer(
		mc,
		config.NewNacosLoader(
			mc.Get().NacosNamespaceID,
			"GO-MATCHER",
			"match_config",
			mc.Get().NacosServers,
		),
		zconfig.Load("cmd/zinx_conf_tmp.yml"),
	)
	defer shutdown()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
}
