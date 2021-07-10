/// logic所依赖的资源池

package svc

import (
	"database/sql"
	"encoding/base64"
	"github.com/go-redis/redis"
	eos "github.com/marsofsnow/eos-go"
	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/core/stores/sqlx"
	"github.com/tal-tech/go-zero/rest"
	"regexp"
	"secret-im/service/signalserver/cmd/api/config"
	"secret-im/service/signalserver/cmd/api/internal/auth"
	"secret-im/service/signalserver/cmd/api/internal/crypto"
	"secret-im/service/signalserver/cmd/api/internal/middleware"
	"secret-im/service/signalserver/cmd/api/internal/model"
	"secret-im/service/signalserver/cmd/api/internal/svc/pubsub"
	"secret-im/service/signalserver/cmd/api/internal/svc/pubsub/dispatch"
	"secret-im/service/signalserver/cmd/api/internal/svc/push"
	"secret-im/service/signalserver/cmd/rpc/bookstore/bookstoreclient"
)

var preKeysInsertStmt string = `insert into t_keys (number,device_id,key_id,public_key) values(?,?,?,?) `

type ServiceContext struct {
	Config config.Config
	//Hub       *chat.Hub
	MatchUserNameRegx *regexp.Regexp


	UserModel model.UserModel //CRUD

	//preKey insertor
	PreKeysInsertor *sqlx.BulkInserter

	//-----------------
	PendAccountsModel model.TPendingAccountsModel
	AccountsModel     model.TAccountsModel
	KeysModel         model.TKeysModel
	MsgsModel         model.TMessagesModel
	ProfileKeyModel   model.TProfilekeyModel
	ProfilesModel     model.TProfilesModel
	UsernameModel     model.TUsernamesModel

	// --------------------
	UserCheck      rest.Middleware // jwt
	CheckBasicAuth rest.Middleware //basic auth
	UserNameCheck  rest.Middleware

	//---------------------
	BookStoreClient bookstoreclient.Bookstore //这个是rpc客户端，发起rpc请求的
	EosApi          *eos.API

	ProfileKeyMap map[string]string

	BackupCredentialsGenerator   *auth.ExternalServiceCredentialGenerator
	DirectoryCredentialsGenerator *auth.ExternalServiceCredentialGenerator

	CertificateGenerator *auth.CertificateGenerator

	RedisClient *redis.Client

	RedisOperation *push.RedisOperation
	Dispatcher     *dispatch.RedisDispatchManager
	PubSubManager  *pubsub.Manager
	PushSender     *push.Sender
}

func NewServiceContext(c config.Config) *ServiceContext {

	//连接redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     config.AppConfig.CacheRedis.Addr,
		Password: config.AppConfig.CacheRedis.Password,
		DB:       config.AppConfig.CacheRedis.DB,
		// PoolSize: 0,
		//MinIdleConns: 32,
	})
	pingRes := redisClient.Ping()
	if pingRes.Err() != nil {
		logx.Error("ping redis error:", pingRes.Err())
		return nil
	} else {
		logx.Info(pingRes.String())
	}

	redisOperation, err := push.NewRedisOperation(redisClient)
	if err != nil {
		logx.Errorf("NewRedisOperation fail,reason:%s", err.Error())
		return nil
	}

	dispatcher := dispatch.NewRedisDispatchManager(redisClient, 128, new(DeadLetterHandler))
	pubSubManager := pubsub.NewManager(dispatcher)
	pushSender := push.NewPushSender(pubSubManager, redisOperation)

	if err := redisClient.Ping().Err(); err != nil {
		logx.Errorf("ping redis fail,reason:%s", err.Error())
		return nil
	}

	mysqlConn := sqlx.NewMysql(c.Mysql.DataSource)

	um := model.NewUserModel(mysqlConn)

	eosApi := eos.New(c.EOSChainUrls[0])
	eosApi.EnableKeepAlives()
	eosApi.Debug = true

	logx.Infof("%+v", config.AppConfig.ServerCertificate)

	//创建证书生成器
	var key [32]byte
	cert := config.AppConfig.ServerCertificate
	privateKey, err := base64.StdEncoding.DecodeString(cert.PrivateKey)
	if err != nil {
		logx.Errorf("base64.StdEncoding.DecodeString(cert.PrivateKey):%s", err)
	}
	certificate, err := base64.StdEncoding.DecodeString(cert.Certificate)
	if err != nil {
		logx.Errorf("base64.StdEncoding.DecodeString(cert.Certificate):%s", err)
	}
	copy(key[:], privateKey)

	certificateGenerator, err := auth.NewCertificateGenerator(
		certificate, crypto.NewDjbECPrivateKey(key), cert.ExpiresDays)
	if err != nil {
		logx.Error("auth.NewCertificateGenerator error:", err.Error())
	}

	preKeysInsertor, err := sqlx.NewBulkInserter(mysqlConn, preKeysInsertStmt)
	if err != nil {
		logx.Error("sqlx.NewBulkInserter(mysqlConn,preKeysInsertStmt) error:", err)
	}
	preKeysInsertor.SetResultHandler(func(result sql.Result, err error) {
		if err != nil {
			logx.Error("bulk insert prekeys error:", err)

		}

	})

	matchUsername, _:= regexp.Compile("^[a-z_][a-z0-9_]+$")

	return &ServiceContext{
		Config:          c,
		MatchUserNameRegx: matchUsername,

		PreKeysInsertor: preKeysInsertor,

		//Hub: hub,
		UserModel: um,

		PendAccountsModel: model.NewTPendingAccountsModel(mysqlConn),
		AccountsModel:     model.NewTAccountsModel(mysqlConn),
		KeysModel:         model.NewTKeysModel(mysqlConn),
		MsgsModel:         model.NewTMessagesModel(mysqlConn),
		ProfileKeyModel:   model.NewTProfilekeyModel(mysqlConn),
		ProfilesModel:     model.NewTProfilesModel(mysqlConn),
		UsernameModel:     model.NewTUsernamesModel(mysqlConn),

		UserCheck:      middleware.NewUsercheckMiddleware().Handle,
		CheckBasicAuth: middleware.NewCheckBasicAuthMiddleware(model.NewTAccountsModel(mysqlConn)).Handle,
		UserNameCheck:  middleware.NewUserNameCheckMiddleware().Handle,
		//BookStoreClient: bookstoreclient.NewBookstore(zrpc.MustNewClient(c.BookStore,zrpc.WithUnaryClientInterceptor(interceptor.TimeInterceptor))),
		EosApi:        eosApi,
		ProfileKeyMap: make(map[string]string),

		BackupCredentialsGenerator: auth.NewExternalServiceCredentialGenerator([]byte(c.BackupService.UserAuthenticationTokenSharedSecret), nil, false),
		DirectoryCredentialsGenerator: auth.NewExternalServiceCredentialGenerator(
			[]byte(c.DirectoryClient.UserAuthenticationTokenSharedSecret),
			[]byte(c.DirectoryClient.UserAuthenticationTokenUserIdSecret),
			true),
		CertificateGenerator:       certificateGenerator,
		RedisClient:                redisClient,
		RedisOperation:             redisOperation,
		Dispatcher:                 dispatcher,
		PubSubManager:              pubSubManager,
		PushSender:                 pushSender,
	}
}
