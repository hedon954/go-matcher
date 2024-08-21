package thirdparty

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"

	"github.com/google/uuid"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"gopkg.in/yaml.v3"
)

func NewNacosConfigClient(namespaceID string, serverConfigs []constant.ServerConfig) (config_client.IConfigClient, error) {
	return clients.NewConfigClient(vo.NacosClientParam{
		ClientConfig:  newClientConfig(namespaceID),
		ServerConfigs: serverConfigs,
	})
}

//nolint:all
func newClientConfig(namespaceID string) *constant.ClientConfig {
	return &constant.ClientConfig{
		NamespaceId:         namespaceID,
		NotLoadCacheAtStart: true,
		LogDir:              "./nacos/log",
		CacheDir:            "./nacos/cache",
		LogLevel:            "debug",
	}
}

// PrepareNacosConfig prepares nacos config, just for test.
func PrepareNacosConfig(addr, dataID, group string, port uint64, config any) (namespaceID string) {
	// create namespace
	namespaceID = uuid.NewString()
	rsp, err := http.PostForm(fmt.Sprintf("http://%s:%d/nacos/v1/console/namespaces", addr, port),
		map[string][]string{
			"customNamespaceId": {namespaceID},
			"namespaceName":     {"tmp-ns"},
			"namespaceDesc":     {"tmp-ns-for-test"},
		})
	if err != nil {
		panic(err)
	}
	defer func() { _ = rsp.Body.Close() }()
	if rsp.StatusCode != http.StatusOK {
		panic(fmt.Errorf("create nacos namespace error: %d, msg: %s", rsp.StatusCode, rsp.Status))
	}

	// create config
	bs, _ := yaml.Marshal(config)
	rsp, err = http.PostForm("http://localhost:8848/nacos/v1/cs/configs", map[string][]string{
		"tenant":  {namespaceID},
		"dataId":  {dataID},
		"group":   {group},
		"type":    {"yaml"},
		"content": {string(bs)},
	})
	if err != nil {
		panic(err)
	}
	defer func() { _ = rsp.Body.Close() }()
	if rsp.StatusCode != http.StatusOK {
		panic(fmt.Errorf("create nacos config error: %d, msg: %s", rsp.StatusCode, rsp.Status))
	}
	return namespaceID
}

// ClearNacosConfig clears nacos config, just for test.
func ClearNacosConfig(namespaceID, addr string, port uint64) {
	formData := url.Values{}
	formData.Add("namespaceId", namespaceID)
	req, err := http.NewRequest("DELETE",
		fmt.Sprintf("http://%s:%d/nacos/v1/console/namespaces", addr, port),
		bytes.NewBufferString(formData.Encode()),
	)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer func() { _ = rsp.Body.Close() }()
	if rsp.StatusCode != http.StatusOK {
		panic(fmt.Errorf("delete nacos namespace error: %d, msg: %s", rsp.StatusCode, rsp.Status))
	}
}
