package cmd

import (
	"github.com/rs/zerolog"
	"github.com/sirupsen/logrus"

	"github.com/hedon954/go-matcher/internal/log"
	"github.com/hedon954/go-matcher/pkg/safe"
)

func init() {
	startSafe()
}

func InitLogger(online bool) {
	initZeroLog(online)
	initLogrus(online)
}

func initZeroLog(online bool) {
	if online {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
}

func initLogrus(online bool) {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	if online {
		logrus.SetLevel(logrus.InfoLevel)
	} else {
		logrus.SetLevel(logrus.DebugLevel)
	}
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
