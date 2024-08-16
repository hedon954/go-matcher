package zconfig

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

const (
	defaultHost             = "0.0.0.0"
	defaultTCPPort          = 7777
	defaultIPVersion        = "tcp4"
	defaultMaxPacketSize    = 12000
	defaultMaxConn          = 4096
	defaultWorkPoolSize     = 10
	defaultMaxWorkerTaskLen = 1024
	defaultMaxMsgChanLen    = 64
)

type ZConfig struct {
	Host             string `yaml:"host"`
	TCPPort          int    `yaml:"tcp_port"`
	IPVersion        string `yaml:"ip_version"`
	Name             string `yaml:"name"`
	Version          string `yaml:"version"`
	MaxPacketSize    uint32 `yaml:"max_packet_size"`
	MaxConn          int    `yaml:"max_conn"`
	WorkPoolSize     uint32 `yaml:"work_pool_size"`
	MaxWorkerTaskLen uint32 `yaml:"max_worker_task_len"`
	MaxMsgChanLen    uint32 `yaml:"max_msg_chan_len"`
}

var DefaultConfig = &ZConfig{
	Host:             defaultHost,
	Name:             "ZinxServerApp",
	Version:          "v1.0",
	TCPPort:          defaultTCPPort,
	IPVersion:        defaultIPVersion,
	MaxPacketSize:    defaultMaxPacketSize,
	MaxConn:          defaultMaxConn,
	WorkPoolSize:     defaultWorkPoolSize,
	MaxWorkerTaskLen: defaultMaxWorkerTaskLen,
	MaxMsgChanLen:    defaultMaxMsgChanLen,
}

func Load(conf string) *ZConfig {
	if conf == "" {
		return DefaultConfig
	}
	bs, err := os.ReadFile(conf)
	if err != nil {
		panic(fmt.Sprintf("read config file error: %v", err))
	}

	c := &ZConfig{}
	err = yaml.Unmarshal(bs, c)
	if err != nil {
		panic(fmt.Sprintf("unmarshal config error: %v", err))
	}

	if c.Host == "" {
		c.Host = defaultHost
	}
	if c.Name == "" {
		c.Name = "ZinxServerApp"
	}
	if c.Version == "" {
		c.Version = "unknown"
	}
	if c.TCPPort == 0 {
		c.TCPPort = defaultTCPPort
	}
	if c.IPVersion == "" {
		c.IPVersion = defaultIPVersion
	}
	if c.MaxPacketSize == 0 {
		c.MaxPacketSize = defaultMaxPacketSize
	}
	if c.MaxMsgChanLen == 0 {
		c.MaxMsgChanLen = defaultMaxMsgChanLen
	}
	return c
}
