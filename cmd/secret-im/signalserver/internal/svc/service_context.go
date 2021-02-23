/// logic所依赖的资源池

package svc

import (
	"github.com/tal-tech/go-zero/rest"
	"secret-im/model"
	"secret-im/signalserver/config"
	"secret-im/signalserver/internal/middleware"

	"github.com/tal-tech/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config config.Config
	UserModel model.UserModel  //CRUD
	UserCheck rest.Middleware  // middleware

}

func NewServiceContext(c config.Config) *ServiceContext {

	mysqlConn := sqlx.NewMysql(c.Mysql.DataSource)
	um:=model.NewUserModel(mysqlConn)

	return &ServiceContext{
		Config: c,
		UserModel: um,
		UserCheck: middleware.NewUsercheckMiddleware().Handle,
	}
}
