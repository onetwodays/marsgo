/// logic所依赖的资源池

package svc

import (
	"encoding/json"
	//eos "github.com/eoscanada/eos-go"
	eos "github.com/eoscanada/eos-go"
	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/core/stores/sqlx"
	"github.com/tal-tech/go-zero/rest"
	"secret-im/service/signalserver/cmd/api/chat"
	"secret-im/service/signalserver/cmd/api/config"
	"secret-im/service/signalserver/cmd/api/db/redis"
	"secret-im/service/signalserver/cmd/api/internal/middleware"
	"secret-im/service/signalserver/cmd/api/internal/types"
	"secret-im/service/signalserver/cmd/api/util"
	"secret-im/service/signalserver/cmd/model"
	"secret-im/service/signalserver/cmd/rpc/bookstore/bookstoreclient"
)

type ServiceContext struct {
	Config    config.Config
	Hub       *chat.Hub
	UserModel model.UserModel //CRUD


	//-----------------
	PendAccountsModel model.TPendAccountsModel
	AccountsModel model.TAccountsModel
	KeysModel model.TKeysModel
	MsgsModel model.TMessagesModel
	// --------------------
	UserCheck rest.Middleware // jwt
	CheckBasicAuth rest.Middleware //basic auth
	BookStoreClient bookstoreclient.Bookstore //这个是rpc客户端，发起rpc请求的
	EosApi *eos.API


}

func NewServiceContext(c config.Config) *ServiceContext {

	mysqlConn := sqlx.NewMysql(c.Mysql.DataSource)
	um:= model.NewUserModel(mysqlConn)
	hub:=chat.NewHub()
	go func() {
		logx.Info("hub run start.....")
		hub.Run()
	}()

	eosApi:=eos.New(c.EOSChainUrls[0])
	eosApi.EnableKeepAlives()
	eosApi.Debug=true




	return &ServiceContext{
		Config:    c,
		Hub: hub,
		UserModel: um,

		PendAccountsModel: model.NewTPendAccountsModel(mysqlConn),
		AccountsModel: model.NewTAccountsModel(mysqlConn),
		KeysModel: model.NewTKeysModel(mysqlConn),
		MsgsModel: model.NewTMessagesModel(mysqlConn),
		UserCheck: middleware.NewUsercheckMiddleware().Handle,
		CheckBasicAuth: middleware.NewCheckBasicAuthMiddleware().Handle,
		//BookStoreClient: bookstoreclient.NewBookstore(zrpc.MustNewClient(c.BookStore,zrpc.WithUnaryClientInterceptor(interceptor.TimeInterceptor))),
		EosApi: eosApi,
	}
}

func (sc *ServiceContext) GetOneAccountByNumber(number string) (int64,*types.Account,error) {
	account:=types.Account{}
	accountModel,err:= sc.AccountsModel.FindOneByNumber(number)
	if err!=nil {
		return -1,nil, err
	}
	data:=accountModel.Data

	err = json.Unmarshal([]byte(data),&account)
	if err!=nil{
		return -1,nil, err
	}
	return accountModel.Id,&account, nil


}

func (sc *ServiceContext) UpdateDirectory(number string,voice,video bool) error {
	hs:= util.ContactToken(number)
	dirToken:=types.GetDirTokenRes{
		Voice: voice,
		Video:  video,
		Relay: "",
		Token: "",
	}
	v,err:=json.Marshal(dirToken)
	if err!=nil{
		return  err
	}
	_,err = redis.RedisDirectoryManager().HSet("directory",string(hs[:]),string(v)).Result()
	if err!=nil{
		return  err
	}
	return nil


}


