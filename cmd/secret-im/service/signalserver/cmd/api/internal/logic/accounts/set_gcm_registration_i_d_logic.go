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

type SetGcmRegistrationIDLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSetGcmRegistrationIDLogic(ctx context.Context, svcCtx *svc.ServiceContext) SetGcmRegistrationIDLogic {
	return SetGcmRegistrationIDLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SetGcmRegistrationIDLogic) SetGcmRegistrationID(r *http.Request,req types.GcmRegistrationID) error {
	account,err:= logic.GetSourceAccount(r,l.svcCtx.AccountsModel)
	if err!=nil{
		return shared.Status(http.StatusUnauthorized,err.Error())
	}
	// 设置GCM注册ID
	device := account.AuthenticatedDevice
	if device.GcmID ==req.GcmRegistrationID {
		return nil
	}
	device.SetGcmID(req.GcmRegistrationID)
	device.SetApnID("")
	device.VoipApnID = ""
	device.FetchesMessages = false
	if err := new(storage.AccountManager).Update(account); err != nil {
		return shared.Status(http.StatusInternalServerError,err.Error())

	}

	return nil
}
