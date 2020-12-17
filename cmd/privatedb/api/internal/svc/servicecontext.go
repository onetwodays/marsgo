package svc

import (
    "github.com/tal-tech/go-zero/core/stores/sqlx"
    "github.com/tal-tech/go-zero/rest"
    "privatedb/api/internal/config"
    "privatedb/api/internal/middleware"
    "privatedb/api/model"
)

type ServiceContext struct {
    Config               config.Config
    EOSUserCheck         rest.Middleware
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
        EOSUserCheck:         middleware.NewEOSUserCheckMiddleware(c.EOSChainUrls[0]).Handle,
        Usercheck:            middleware.NewUsercheckMiddleware().Handle,
        UserModer:            um,
        TMsgModel:            tmsgm,
        TPayTypeModel:        tpaymenttypem,
        TPaymentAccountModel: tpaymentaccountm,
    }
}
