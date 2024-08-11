package main

import (
	"github.com/hedon954/go-matcher/cmd"
	"github.com/hedon954/go-matcher/internal/api/apihttp"
)

func main() {
	defer cmd.StopSafe()
	apihttp.SetupHTTPServer()
}
