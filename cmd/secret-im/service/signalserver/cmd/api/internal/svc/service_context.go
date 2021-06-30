/// logic所依赖的资源池

package svc

import (
	"database/sql"
	"encoding/base64"
	eos "github.com/marsofsnow/eos-go"
	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/core/stores/sqlx"
	"github.com/tal-tech/go-zero/rest"
	"secret-im/service/signalserver/cmd/api/config"
	"secret-im/service/signalserver/cmd/api/internal/auth"
	"secret-im/service/signalserver/cmd/api/internal/crypto"
	"secret-im/service/signalserver/cmd/api/internal/middleware"
	"secret-im/service/signalserver/cmd/model"
	"secret-im/service/signalserver/cmd/rpc/bookstore/bookstoreclient"
)

var preKeysInsertStmt string = `insert into t_keys (number,device_id,key_id,public_key) values(?,?,?,?) `

type ServiceContext struct {
	Config    config.Config
	//Hub       *chat.Hub


	UserModel model.UserModel //CRUD

	//preKey insertor
	PreKeysInsertor *sqlx.BulkInserter


	//-----------------
	PendAccountsModel model.TPendingAccountsModel
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

	BackupCredentialsGenerator *auth.ExternalServiceCredentialGenerator

	CertificateGenerator    *auth.CertificateGenerator


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


	logx.Infof("%+v",config.AppConfig.ServerCertificate)

	//创建证书生成器
	var key [32]byte
	cert := config.AppConfig.ServerCertificate
	privateKey,err:= base64.StdEncoding.DecodeString(cert.PrivateKey)
	if err!=nil{
		logx.Errorf("base64.StdEncoding.DecodeString(cert.PrivateKey):%s",err)
	}
	certificate,err:= base64.StdEncoding.DecodeString(cert.Certificate)
	if err!=nil{
		logx.Errorf("base64.StdEncoding.DecodeString(cert.Certificate):%s",err)
	}
	copy(key[:],privateKey)

	certificateGenerator, err := auth.NewCertificateGenerator(
		certificate, crypto.NewDjbECPrivateKey(key), cert.ExpiresDays)
	if err!=nil{
		logx.Error("auth.NewCertificateGenerator error:",err.Error())
	}

	preKeysInsertor,err:=sqlx.NewBulkInserter(mysqlConn,preKeysInsertStmt)
	if err!=nil {
		logx.Error("sqlx.NewBulkInserter(mysqlConn,preKeysInsertStmt) error:",err)
	}
	preKeysInsertor.SetResultHandler(func (result sql.Result,err error){
		if err!=nil{
			logx.Error("bulk insert prekeys error:",err)
		}

	})


	return &ServiceContext{
		Config:    c,
		PreKeysInsertor: preKeysInsertor,

		//Hub: hub,
		UserModel: um,

		PendAccountsModel: model.NewTPendingAccountsModel(mysqlConn),
		AccountsModel: model.NewTAccountsModel(mysqlConn),
		KeysModel: model.NewTKeysModel(mysqlConn),
		MsgsModel: model.NewTMessagesModel(mysqlConn),
		ProfileKeyModel: model.NewTProfilekeyModel(mysqlConn),

		UserCheck: middleware.NewUsercheckMiddleware().Handle,
		CheckBasicAuth: middleware.NewCheckBasicAuthMiddleware(model.NewTAccountsModel(mysqlConn)).Handle,
		UserNameCheck: middleware.NewUserNameCheckMiddleware().Handle,
		//BookStoreClient: bookstoreclient.NewBookstore(zrpc.MustNewClient(c.BookStore,zrpc.WithUnaryClientInterceptor(interceptor.TimeInterceptor))),
		EosApi: eosApi,
		ProfileKeyMap:make(map[string]string),

		BackupCredentialsGenerator: auth.NewExternalServiceCredentialGenerator([]byte(c.BackupService.UserAuthenticationTokenSharedSecret),nil,false),
		CertificateGenerator: certificateGenerator,
	}
}


