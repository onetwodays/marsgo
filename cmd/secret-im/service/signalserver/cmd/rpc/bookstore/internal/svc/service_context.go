package svc

import (
	//"secret-im/service/signalserver/cmd/api/internal/model"
	"secret-im/service/signalserver/cmd/rpc/bookstore/internal/config"
)

type ServiceContext struct {
	Config    config.Config
	//BookModel model.BookModel // add
}

func NewServiceContext(c config.Config) *ServiceContext {

	return &ServiceContext{
		Config:    c,
		//BookModel: model.NewBookModel(sqlx.NewMysql(c.Mysql.DataSource),c.CacheRedis),
	}
}
