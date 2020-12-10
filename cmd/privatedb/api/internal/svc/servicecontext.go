package svc

import (
    "github.com/tal-tech/go-zero/core/logx"
    "github.com/tal-tech/go-zero/core/stores/sqlx"
    "github.com/tal-tech/go-zero/rest"
    "net/http"
    "privatedb/api/internal/config"
    "privatedb/api/internal/middleware"
    "privatedb/api/model"
)

type ServiceContext struct {
    Config               config.Config
    GreetMiddleware1     rest.Middleware
    GreetMiddleware2     rest.Middleware
    Usercheck            rest.Middleware
    UserModer            model.UserModel
    TMsgModel            model.TMsgModel
    TPayTypeModel        model.TPayTypeModel
    TPaymentAccountModel model.TPaymentAccountModel
}

func NewServiceContext(c config.Config) *ServiceContext {
    conn := sqlx.NewMysql(c.Mysql.DataSource)
    um := model.NewUserModel(conn)
    tmsgm := model.NewTMsgModel(conn)
    tpaymenttypem := model.NewTPayTypeModel(conn)
    tpaymentaccountm := model.NewTPaymentAccountModel(conn)
    return &ServiceContext{
        Config:               c,
        GreetMiddleware1:     greetMiddleware1,
        GreetMiddleware2:     greetMiddleware2,
        Usercheck:            middleware.NewUsercheckMiddleware().Handle,
        UserModer:            um,
        TMsgModel:            tmsgm,
        TPayTypeModel:        tpaymenttypem,
        TPaymentAccountModel: tpaymentaccountm,
    }
}

func greetMiddleware1(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        logx.Info("greetMiddleware1 request ... ")
        next(w, r)
        logx.Info("greetMiddleware1 reponse ... ")
    }
}

func greetMiddleware2(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        logx.Info("greetMiddleware2 request ... ")
        next(w, r)
        logx.Info("greetMiddleware2 reponse ... ")
    }
}
