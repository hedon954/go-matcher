package cmd

import (
	"github.com/hedon954/go-matcher/internal/log"
	"github.com/hedon954/go-matcher/pkg/safe"
)

func init() {
	startSafe()
}

func startSafe() {
	safe.Callback(func(err any, stack []byte) {
		log.Error().
			Any("err", err).
			Str("stack", string(stack)).
			Msg("safe occurs panic")
	})
}

func StopSafe() {
	safe.Wait()
}
