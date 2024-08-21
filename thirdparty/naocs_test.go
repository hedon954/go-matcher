package thirdparty

import (
	"fmt"
	"testing"

	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/stretchr/testify/assert"
)

func TestNewNacosClient(t *testing.T) {
	client, err := NewNacosClient("7d638262-9e51-4822-9333-c3bcca838b7d")
	assert.Nil(t, err)
	c, err := client.GetConfig(vo.ConfigParam{
		DataId: "basic_config",
		Group:  "DEFAULT_GROUP",
		OnChange: func(namespace, group, dataId, data string) {
			fmt.Printf("group:%s  dataId:%s  data:%s\n", group, dataId, data)
		},
	})
	assert.Nil(t, err)
	fmt.Println(c)
}

func prepareNacosConfig(client config_client.IConfigClient) {

}

func clearNacosConfig(client config_client.IConfigClient) {

}
