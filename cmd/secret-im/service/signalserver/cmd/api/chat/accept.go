package chat

import (
	"github.com/tal-tech/go-zero/rest/httpx"
	"net/http"
	"secret-im/service/signalserver/cmd/api/internal/entities"
	"secret-im/service/signalserver/cmd/api/internal/middleware"
	"secret-im/service/signalserver/cmd/api/internal/svc"
	"strings"
)



type WsConnReq struct {

	Login string `form:"login,optional"`
	Password string `form:"password,optional"`
}

func authenticate(r *http.Request,ctx *svc.ServiceContext) (*entities.Account,bool){

	var ba  WsConnReq
	if err := httpx.Parse(r, &ba); err != nil {
		return nil,true
	}
	if len(ba.Login)==0 && len(ba.Password)==0{
		return nil,true
	}

	ba.Login = strings.ReplaceAll(ba.Login," ","+")
	ba.Password = strings.ReplaceAll(ba.Password," ","+")

	// 帐号鉴权
	enabledRequired :=true
	/*
	if true{
		enabledRequired =false
	}
	 */
	checkBasicAuth := middleware.NewCheckBasicAuthMiddleware(ctx.AccountsModel)
	appAccount,err,ok:=checkBasicAuth.BasicAuthForHeader(r,enabledRequired)
	if err!=nil{
		return nil, false
	}
	return appAccount,ok

}






