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
	_, _, shutdown := apitcp.SetupTCPServer(
		config.NewFileLoader[config.ServerConfig]("cmd/server_conf_tmp.yml"),
		config.NewFileLoader[config.MatchConfig]("cmd/match_conf_tmp.yml"),
		zconfig.Load("cmd/zinx_conf_tmp.yml"),
	)
	defer shutdown()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
}
