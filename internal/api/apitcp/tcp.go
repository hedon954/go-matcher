package apitcp

import (
	"github.com/hedon954/go-matcher/pkg/safe"
	"github.com/hedon954/go-matcher/pkg/zinx/ziface"
	"github.com/hedon954/go-matcher/pkg/zinx/znet"
)

func SetupTCPServer(conf string) ziface.IServer {
	server := znet.NewServer(conf)
	registerRoute(server)

	safe.Go(server.Serve)
	return server
}

func registerRoute(_ ziface.IServer) {}
