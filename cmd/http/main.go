package main

import (
	"github.com/hedon954/go-matcher/internal/api"
)

// swag init --generalInfo  ../../internal/api/http.go --dir ../../internal/api
func main() {
	api.SetupHTTPServer()
}
