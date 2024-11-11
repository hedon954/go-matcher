package cmd

import (
	"os"

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

// GetConfigPath returns the config path from environment variable or default path
func GetConfigPath(envKey, defaultPath string) string {
	if path := os.Getenv(envKey); path != "" {
		return path
	}
	return defaultPath
}
