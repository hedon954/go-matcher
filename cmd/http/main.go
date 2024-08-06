package main

import (
	"log/slog"

	"github.com/hedon954/go-matcher/internal/api"
	"github.com/hedon954/go-matcher/pkg/safe"
)

// swag init --generalInfo  ../../internal/api/http.go --dir ../../internal/api
func main() {
	startSafe()
	defer stopSafe()
	api.SetupHTTPServer()
}

func startSafe() {
	safe.Callback(func(err any, stack []byte) {
		slog.Error("safe occurs panic",
			slog.Any("err", err),
			slog.String("stack", string(stack)),
		)
	})
}

func stopSafe() {
	safe.Wait()
}

func init() {
	startSafe()
}
