package main

import (
	"github.com/sirupsen/logrus"

	"github.com/hedon954/go-matcher/cmd"
	"github.com/hedon954/go-matcher/internal/api/apihttp"
	"github.com/hedon954/go-matcher/internal/config"
)

func main() {
	defer cmd.StopSafe()
	sc := config.NewFileLoader[config.ServerConfig]("cmd/server_conf_tmp.yml")
	mc := config.NewNacosLoader(
		sc.Get().NacosNamespaceID,
		"GO-MATCHER",
		"match_config",
		sc.Get().NacosServers)

	// mc := config.NewFileLoader[config.MatchConfig]("cmd/match_conf_tmp.yml")

	cmd.InitLogger(sc.Get().IsOnline())

	logrus.WithFields(logrus.Fields{
		"match_config":  mc.Get(),
		"server_config": sc.Get(),
	}).Info("load config success")
	apihttp.SetupHTTPServer(sc, mc)
}
