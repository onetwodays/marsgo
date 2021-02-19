package config

import (
	"github.com/spf13/viper"

	"fmt"
)

var AppConfig *viper.Viper  //构造一个导出的包级全局变量，在init()进行初始化

func init() {
	AppConfig = viper.New()
	AppConfig.SetConfigName("config")
	AppConfig.SetConfigType("toml")
	AppConfig.AddConfigPath("config")
	err := AppConfig.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}

// config.AppConfig.GetString("system.port"))