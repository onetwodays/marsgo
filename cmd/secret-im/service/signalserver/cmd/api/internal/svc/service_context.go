/// logic所依赖的资源池

package svc

import (
	"github.com/tal-tech/go-zero/rest"
	"github.com/tal-tech/go-zero/zrpc"
	"secret-im/service/signalserver/cmd/api/config"
	"secret-im/service/signalserver/cmd/api/internal/middleware"
	"secret-im/service/signalserver/cmd/model"
	"secret-im/service/signalserver/cmd/rpc/bookstore/bookstoreclient"
	"secret-im/service/signalserver/cmd/rpc/interceptor"

	"github.com/tal-tech/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config    config.Config
	UserModel model.UserModel //CRUD
	UserCheck rest.Middleware // middleware
	BookStoreClient bookstoreclient.Bookstore //这个是rpc客户端，发起rpc请求的

}

func NewServiceContext(c config.Config) *ServiceContext {

	mysqlConn := sqlx.NewMysql(c.Mysql.DataSource)
	um:= model.NewUserModel(mysqlConn)

	return &ServiceContext{
		Config:    c,
		UserModel: um,
		UserCheck: middleware.NewUsercheckMiddleware().Handle,
		BookStoreClient: bookstoreclient.NewBookstore(zrpc.MustNewClient(c.BookStore,zrpc.WithUnaryClientInterceptor(interceptor.TimeInterceptor))),
	}
}
