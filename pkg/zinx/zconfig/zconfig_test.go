package zconfig

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	t.Run("empty config path should return default", func(t *testing.T) {
		conf := Load("")
		assert.Equal(t, conf, DefaultConfig)
	})

	t.Run("file not exists should panic", func(t *testing.T) {
		assert.Panics(t, func() {
			_ = Load("not_exists.yml")
		})
	})

	t.Run("invalid config should panic", func(t *testing.T) {
		assert.Panics(t, func() {
			_ = Load("zconfig_test.go")
		})
	})

	t.Run("valid config should work", func(t *testing.T) {
		conf := Load("../../../fixtures/zconfig_test.yml")
		assert.Equal(t, conf, &ZConfig{
			Name:             "ZinxServerApp",
			Host:             defaultHost,
			TCPPort:          defaultTCPPort,
			IPVersion:        defaultIPVersion,
			Version:          "unknown",
			MaxPacketSize:    defaultMaxPacketSize,
			MaxConn:          128,
			WorkPoolSize:     0,
			MaxWorkerTaskLen: 2024,
			MaxMsgChanLen:    defaultMaxMsgChanLen,
		})
	})
}
