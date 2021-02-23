package config

import (
	"flag"
	"github.com/tal-tech/go-zero/core/conf"
	"github.com/tal-tech/go-zero/rest"
)

var configFile = flag.String("f", "etc/signalserver-api.yaml", "the config file")
var AppConfig Config

func init()  {
	flag.Parse()
	conf.MustLoad(*configFile, &AppConfig)
}

type Config struct {
	rest.RestConf
	Auth struct {
		AccessSecret string
		AccessExpire int64
	}
	WssAddress string
	CacheRedis struct{
		Addr string
		Password string
		DB int
	}
	DirectoryRedis struct{
		Addr string
		Password string
		DB int
	}
	Mysql struct{
		DataSource string
	}

}
