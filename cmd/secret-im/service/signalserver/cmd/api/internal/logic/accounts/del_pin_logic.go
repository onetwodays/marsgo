package logic

import (
	"context"
	"net/http"
	"secret-im/service/signalserver/cmd/api/internal/entities"
	"secret-im/service/signalserver/cmd/api/internal/storage"
	"secret-im/service/signalserver/cmd/api/shared"

	"github.com/tal-tech/go-zero/core/logx"
	"secret-im/service/signalserver/cmd/api/internal/svc"
)

type DelPinLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDelPinLogic(ctx context.Context, svcCtx *svc.ServiceContext) DelPinLogic {
	return DelPinLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DelPinLogic) DelPin(r *http.Request) error {
	appAccount := r.Context().Value(shared.HttpReqContextAccountKey)
	if appAccount == nil {
		reason := "check basic auth fail ,may by the handler not use middle"
		logx.Error(reason)
		return shared.Status(http.StatusUnauthorized, reason)
	}
	account := appAccount.(*entities.Account)
	account.Pin = ""

	if err := new(storage.AccountManager).Update(account); err != nil {
		return shared.Status(http.StatusInternalServerError, err.Error())

	}

	return nil
}
