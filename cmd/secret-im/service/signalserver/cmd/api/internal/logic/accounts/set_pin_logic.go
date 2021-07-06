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

type SetPinLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSetPinLogic(ctx context.Context, svcCtx *svc.ServiceContext) SetPinLogic {
	return SetPinLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SetPinLogic) SetPin(r *http.Request,req types.DeprecatedPin) error {
	appAccount := r.Context().Value(shared.HttpReqContextAccountKey)
	if appAccount == nil  {
		reason := "check basic auth fail ,may by the handler not use middle"
		logx.Error(reason)
		return shared.Status(http.StatusUnauthorized, reason)
	}
	account := appAccount.(*entities.Account)
	account.Pin=req.Pin

	if err := new(storage.AccountManager).Update(account); err != nil {
		return shared.Status(http.StatusInternalServerError,err.Error())

	}

	return nil
}
