package main

import (
	"github.com/hedon954/go-matcher/cmd"
	"github.com/hedon954/go-matcher/internal/api/apihttp"
	"github.com/hedon954/go-matcher/internal/config"
)

func main() {
	defer cmd.StopSafe()
	configer := config.NewFileLoader("conf.yml")
	apihttp.SetupHTTPServer(configer)
}
