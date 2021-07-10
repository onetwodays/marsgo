package logic

import (
	"context"
	"net/http"
	"secret-im/service/signalserver/cmd/api/internal/logic"
	"secret-im/service/signalserver/cmd/api/internal/storage"
	"secret-im/service/signalserver/cmd/api/shared"

	"secret-im/service/signalserver/cmd/api/internal/svc"
	"secret-im/service/signalserver/cmd/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type SetDeviceNameLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSetDeviceNameLogic(ctx context.Context, svcCtx *svc.ServiceContext) SetDeviceNameLogic {
	return SetDeviceNameLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SetDeviceNameLogic) SetDeviceName(r *http.Request,req types.DeviceName) error {
	account,err:= logic.GetSourceAccount(r,l.svcCtx.AccountsModel)
	if err!=nil{
		return shared.Status(http.StatusUnauthorized,err.Error())
	}
	account.AuthenticatedDevice.Name = req.DeviceName
	if err := new(storage.AccountManager).Update(account); err != nil {
		return shared.Status(http.StatusInternalServerError,err.Error())

	}
	return nil
}
