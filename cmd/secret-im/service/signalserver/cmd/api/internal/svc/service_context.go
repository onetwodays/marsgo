/// logic所依赖的资源池

package svc

import (
	eos "github.com/marsofsnow/eos-go"
	"github.com/tal-tech/go-zero/core/stores/sqlx"
	"github.com/tal-tech/go-zero/rest"
	"secret-im/service/signalserver/cmd/api/config"
	"secret-im/service/signalserver/cmd/api/internal/middleware"
	"secret-im/service/signalserver/cmd/model"
	"secret-im/service/signalserver/cmd/rpc/bookstore/bookstoreclient"
)

type ServiceContext struct {
	Config    config.Config
	//Hub       *chat.Hub
	UserModel model.UserModel //CRUD


	//-----------------
	PendAccountsModel model.TPendAccountsModel
	AccountsModel model.TAccountsModel
	KeysModel model.TKeysModel
	MsgsModel model.TMessagesModel
	ProfileKeyModel model.TProfilekeyModel
	// --------------------
	UserCheck rest.Middleware // jwt
	CheckBasicAuth rest.Middleware //basic auth
	UserNameCheck    rest.Middleware

	//---------------------
	BookStoreClient bookstoreclient.Bookstore //这个是rpc客户端，发起rpc请求的
	EosApi *eos.API

	ProfileKeyMap  map[string]string


}

func NewServiceContext(c config.Config) *ServiceContext {

	mysqlConn := sqlx.NewMysql(c.Mysql.DataSource)
	um:= model.NewUserModel(mysqlConn)
	/*
	hub:=chat.NewHub()
	go func() {
		logx.Info("hub run start.....")
		hub.Run()
	}()
	chat.SetHub(hub)

	 */



	eosApi:=eos.New(c.EOSChainUrls[0])
	eosApi.EnableKeepAlives()
	eosApi.Debug=true




	return &ServiceContext{
		Config:    c,
		//Hub: hub,
		UserModel: um,

		PendAccountsModel: model.NewTPendAccountsModel(mysqlConn),
		AccountsModel: model.NewTAccountsModel(mysqlConn),
		KeysModel: model.NewTKeysModel(mysqlConn),
		MsgsModel: model.NewTMessagesModel(mysqlConn),
		ProfileKeyModel: model.NewTProfilekeyModel(mysqlConn),

		UserCheck: middleware.NewUsercheckMiddleware().Handle,
		CheckBasicAuth: middleware.NewCheckBasicAuthMiddleware().Handle,
		UserNameCheck: middleware.NewUserNameCheckMiddleware().Handle,
		//BookStoreClient: bookstoreclient.NewBookstore(zrpc.MustNewClient(c.BookStore,zrpc.WithUnaryClientInterceptor(interceptor.TimeInterceptor))),
		EosApi: eosApi,
		ProfileKeyMap:make(map[string]string),
	}
}


