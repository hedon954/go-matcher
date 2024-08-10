package utils

import (
	"os"

	"github.com/hedon954/go-matcher/pkg/zinx/ziface"
	"gopkg.in/yaml.v3"
)

const (
	defaultTCPPort          = 7777
	defaultMaxPacketSize    = 12000
	defaultMaxConn          = 4096
	defaultWorkPoolSize     = 10
	defaultMaxWorkerTaskLen = 1024
	defaultMaxMsgChanLen    = 64
)

type GlobalObj struct {
	TCPServer        ziface.IServer `yaml:"-"`
	Host             string         `yaml:"host"`
	TCPPort          int            `yaml:"tcp_port"`
	Name             string         `yaml:"name"`
	Version          string         `yaml:"version"`
	MaxPacketSize    uint32         `yaml:"max_packet_size"`
	MaxConn          int            `yaml:"max_conn"`
	WorkPoolSize     uint32         `yaml:"work_pool_size"`
	MaxWorkerTaskLen uint32         `yaml:"max_worker_task_len"`
	MaxMsgChanLen    uint32         `yaml:"max_msg_chan_len"`
}

var GlobalObject *GlobalObj

func init() {
	GlobalObject = &GlobalObj{
		Host:             "0.0.0.0",
		Name:             "ZinxServerApp",
		Version:          "v0.4",
		TCPPort:          defaultTCPPort,
		MaxPacketSize:    defaultMaxPacketSize,
		MaxConn:          defaultMaxConn,
		WorkPoolSize:     defaultWorkPoolSize,
		MaxWorkerTaskLen: defaultMaxWorkerTaskLen,
		MaxMsgChanLen:    defaultMaxMsgChanLen,
	}
}

func (g *GlobalObj) Reload(conf string) {
	bs, err := os.ReadFile(conf)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(bs, &GlobalObject)
	if err != nil {
		panic(nil)
	}
}
