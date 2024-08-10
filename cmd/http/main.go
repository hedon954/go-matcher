package main

import (
	"github.com/hedon954/go-matcher/internal/api"
	"github.com/hedon954/go-matcher/pkg/safe"
	"github.com/rs/zerolog/log"
)

func main() {
	startSafe()
	defer stopSafe()
	api.SetupHTTPServer()
}

func startSafe() {
	safe.Callback(func(err any, stack []byte) {
		log.Error().
			Any("err", err).
			Str("stack", string(stack)).
			Msg("safe occurs panic")
	})
}

func stopSafe() {
	safe.Wait()
}

func init() {
	startSafe()
}
