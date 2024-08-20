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
	conf, err := config.NewFileLoader("cmd/tcp/conf.yml").Load()
	if err != nil {
		panic(err)
	}
	_, _, shutdown := apitcp.SetupTCPServer(
		conf, zconfig.Load("cmd/tcp/zinx.yml"),
	)
	defer shutdown()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
}
