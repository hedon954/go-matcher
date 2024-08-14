package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/hedon954/go-matcher/cmd"
	"github.com/hedon954/go-matcher/internal/api/apitcp"
	"github.com/hedon954/go-matcher/pkg/zinx/zconfig"
)

func main() {
	defer cmd.StopSafe()
	conf := zconfig.Load("cmd/tcp/zinx.yml")
	_, server := apitcp.SetupTCPServer(1, conf)
	defer server.Stop()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
}
