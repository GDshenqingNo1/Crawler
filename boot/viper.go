package boot

import (
	g "WebCrawler/app/global"
	"flag"
	"fmt"
	"github.com/spf13/viper"
	"os"
)

const (
	configEnv  = "CRAWLER_CONFIG_PATH"
	configFile = "manifest/config/config.yaml"
)

func ViperSetup(path ...string) {
	var configPath string

	if len(path) != 0 {
		configPath = path[0]
	} else {
		flag.StringVar(&configPath, "c", "", "设置配置文件路径")
		flag.Parse()
		if configPath == "" {
			if configPath = os.Getenv(configEnv); configPath != "" {
			} else {
				configPath = configFile
			}
		}
	}

	fmt.Printf("get config path: %v ", configPath)

	v := viper.New()
	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("get config file failed,err:%v", err))
	}
	if err = v.Unmarshal(&g.Config); err != nil {
		// 将配置文件反序列化到 Config 结构体
		panic(fmt.Errorf("unmarshal config failed, err: %v", err))
	}
}
