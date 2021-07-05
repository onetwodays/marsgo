package websocket

import (
	"github.com/gorilla/websocket"
	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/rest/httpx"
	"net/http"
	"secret-im/service/signalserver/cmd/api/internal/entities"
	"secret-im/service/signalserver/cmd/api/internal/middleware"
	"secret-im/service/signalserver/cmd/api/internal/svc"
	"strings"
)

// 连接升级
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WsConnReq struct {
	Login    string `form:"login,optional"`
	Password string `form:"password,optional"`
}

func authenticate(r *http.Request, ctx *svc.ServiceContext) (*entities.Account, bool) {

	var ba WsConnReq
	if err := httpx.Parse(r, &ba); err != nil {
		return nil, true
	}
	if len(ba.Login) == 0 && len(ba.Password) == 0 {
		return nil, true
	}

	ba.Login = strings.ReplaceAll(ba.Login, " ", "+")
	ba.Password = strings.ReplaceAll(ba.Password, " ", "+")

	// 帐号鉴权
	enabledRequired := true
	/*
		if true{
			enabledRequired =false
		}
	*/
	checkBasicAuth := middleware.NewCheckBasicAuthMiddleware(ctx.AccountsModel)
	logx.Info("ba.login=", ba.Login)
	logx.Info("ba.passwd=", ba.Password)
	appAccount, err := checkBasicAuth.BasicAuthByUserPasswd(ba.Login, ba.Password, enabledRequired)
	if err != nil {
		logx.Error("checkBasicAuth.BasicAuthByUserPasswd fail:",err.Error())
		return nil, false
	}
	return appAccount, true

}
