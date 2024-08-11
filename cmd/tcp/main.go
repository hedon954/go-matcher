package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/hedon954/go-matcher/cmd"
	"github.com/hedon954/go-matcher/internal/api/apitcp"
)

func main() {
	defer cmd.StopSafe()
	server := apitcp.SetupTCPServer("zinx.yml")
	defer server.Stop()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
}
