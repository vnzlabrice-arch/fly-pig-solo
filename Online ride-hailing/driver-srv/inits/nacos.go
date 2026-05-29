package init

import (
	"driver-srv/global"
	"fmt"
	"strings"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/spf13/viper"
)

func NoCosInit() {
	viper.SetConfigFile("../user-srv/nacos.yaml")

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	var NaCosData global.NaCos
	err = viper.UnmarshalKey("NaCos", &NaCosData)
	if err != nil {
		panic(err)
	}
	fmt.Println(NaCosData)
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr: NaCosData.Host,
			Port:   uint64(NaCosData.Port),
		},
	}
	// 客户端配置
	clientConfig := constant.ClientConfig{
		NamespaceId:         "", // 如果不需要命名空间，可以留空
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "../nacos/log",
		CacheDir:            "../nacos/cache",
		LogLevel:            "debug",
	}

	// 创建配置客户端
	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig":  clientConfig,
	})
	config, err := configClient.GetConfig(vo.ConfigParam{
		DataId: NaCosData.DataId,
		Group:  NaCosData.Group,
	})
	viper.Reset()

	viper.SetConfigType("yaml")

	err = viper.ReadConfig(strings.NewReader(config))
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&global.ConfigData)
	if err != nil {
		panic(err)
	}
	fmt.Println("配置成功", global.ConfigData)

}
