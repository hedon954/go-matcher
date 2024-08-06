package main

import (
	"github.com/hedon954/go-matcher/internal/api"
	"github.com/hedon954/go-matcher/pkg/safe"
)

// swag init --generalInfo  ../../internal/api/http.go --dir ../../internal/api
func main() {
	api.SetupHTTPServer()
	safe.Wait()
}
