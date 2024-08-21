package thirdparty

import (
	"fmt"
	"testing"

	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

var (
	namespaceID string
	host        = "127.0.0.1"
	port        = uint64(8848)
	dataID      = "test-nacos-data-id"
	group       = "test-nacos-group"
	expected    = nacosConfig{
		Name:      "hedon",
		Addr:      "home",
		IsMarried: true,
		Age:       18,
		Extra: struct {
			Company string `yaml:"company"`
			Salary  int    `yaml:"salary"`
		}{
			Company: "nacos",
			Salary:  10000,
		},
	}
)

type nacosConfig struct {
	Name      string `yaml:"name"`
	Addr      string `yaml:"addr"`
	IsMarried bool   `yaml:"is_married"`
	Age       int    `yaml:"age"`
	Extra     struct {
		Company string `yaml:"company"`
		Salary  int    `yaml:"salary"`
	} `yaml:"extra"`
}

func TestMain(m *testing.M) {
	namespaceID = PrepareNacosConfig(host, dataID, group, port, expected)
	m.Run()
	ClearNacosConfig(namespaceID, host, port)
}

func TestNewNacosClient(t *testing.T) {
	client, err := NewNacosConfigClient(namespaceID, []constant.ServerConfig{
		{
			IpAddr:      host,
			Port:        port,
			ContextPath: "/nacos",
			Scheme:      "http",
		},
	})
	assert.Nil(t, err)
	c, err := client.GetConfig(vo.ConfigParam{
		DataId: dataID,
		Group:  group,
		OnChange: func(namespace, group, dataId, data string) {
			fmt.Printf("group:%s  dataId:%s  data:%s\n", group, dataId, data)
		},
	})
	assert.Nil(t, err)

	var c1 nacosConfig
	err = yaml.Unmarshal([]byte(c), &c1)
	assert.Nil(t, err)
	assert.Equal(t, expected, c1)
}
