/// logic所依赖的资源池

package svc

import (
	"github.com/tal-tech/go-zero/core/stores/sqlx"
	"secret-im/signalserver/config"
	"secret-im/model"
)

type ServiceContext struct {
	Config config.Config
	UserModel model.UserModel

}

func NewServiceContext(c config.Config) *ServiceContext {

	mysqlConn := sqlx.NewMysql(c.Mysql.DataSource)
	um:=model.NewUserModel(mysqlConn)

	return &ServiceContext{
		Config: c,
		UserModel: um,
	}
}
