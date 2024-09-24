package thirdparty

import (
	"testing"
	"time"

	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

var (
	host     = "127.0.0.1"
	port     = uint64(8848)
	grpcPort = uint64(9848)
	dataID   = "test-nacos-data-id"
	group    = "test-nacos-group"
	expected = nacosConfig{
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

func TestNewNacosClient(t *testing.T) {
	if testing.Short() {
		t.Skip("because skywalking-go would panic in this case, should just skip it")
	}

	namespaceID := PrepareNacosConfig(host, dataID, group, port, expected)
	defer ClearNacosConfig(namespaceID, host, port)

	client, err := NewNacosConfigClient(namespaceID, []constant.ServerConfig{
		{
			IpAddr:      host,
			Port:        port,
			GrpcPort:    grpcPort,
			ContextPath: "/nacos",
			Scheme:      "http",
		},
	})
	assert.Nil(t, err)
	assert.NotNil(t, client)

	time.Sleep(time.Second)
	c, err := client.GetConfig(vo.ConfigParam{
		DataId: dataID,
		Group:  group,
		OnChange: func(namespace, group, dataId, data string) {
		},
	})
	assert.Nil(t, err)

	var c1 nacosConfig
	err = yaml.Unmarshal([]byte(c), &c1)
	assert.Nil(t, err)
	assert.Equal(t, expected, c1)
}
