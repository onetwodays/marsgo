package logic

import (
	"context"
	"net/http"
	"secret-im/service/signalserver/cmd/api/internal/entities"
	"secret-im/service/signalserver/cmd/api/internal/storage"
	"secret-im/service/signalserver/cmd/api/shared"

	"secret-im/service/signalserver/cmd/api/internal/svc"
	"secret-im/service/signalserver/cmd/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type SetNameLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSetNameLogic(ctx context.Context, svcCtx *svc.ServiceContext) SetNameLogic {
	return SetNameLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SetNameLogic) SetName(r *http.Request,req types.SetNameReq) error {
	appAccount := r.Context().Value(shared.HttpReqContextAccountKey)
	if appAccount == nil  {
		reason := "check basic auth fail ,may by the handler not use middle"
		logx.Error(reason)
		return shared.Status(http.StatusUnauthorized, reason)
	}
	account := appAccount.(*entities.Account)
	account.Name=req.Name
	if err:=new(storage.AccountManager).Update(account);err!=nil{
		return shared.Status(http.StatusInternalServerError,err.Error())
	}

	return nil
}
