package logic

import (
	"context"
	"net/http"
	"secret-im/service/signalserver/cmd/api/internal/logic"
	"secret-im/service/signalserver/cmd/api/internal/storage"
	"secret-im/service/signalserver/cmd/api/shared"

	"github.com/tal-tech/go-zero/core/logx"
	"secret-im/service/signalserver/cmd/api/internal/svc"
)

type DelGcmRegistrationIDLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDelGcmRegistrationIDLogic(ctx context.Context, svcCtx *svc.ServiceContext) DelGcmRegistrationIDLogic {
	return DelGcmRegistrationIDLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DelGcmRegistrationIDLogic) DelGcmRegistrationID(r *http.Request) error {
	account,err:= logic.GetSourceAccount(r,l.svcCtx.AccountsModel)
	if err!=nil{
		return shared.Status(http.StatusUnauthorized,err.Error())
	}

	// 删除GCM注册ID
	device := account.AuthenticatedDevice
	device.SetGcmID("")
	device.FetchesMessages = false

	if err := new(storage.AccountManager).Update(account); err != nil {
		return shared.Status(http.StatusInternalServerError, err.Error())

	}

	return nil
}
