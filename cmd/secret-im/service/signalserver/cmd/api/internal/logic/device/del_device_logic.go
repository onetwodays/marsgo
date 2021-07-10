package logic

import (
	"context"
	"net/http"
	"secret-im/service/signalserver/cmd/api/internal/logic"
	"secret-im/service/signalserver/cmd/api/internal/storage"
	"secret-im/service/signalserver/cmd/api/shared"
	"strconv"

	"secret-im/service/signalserver/cmd/api/internal/svc"
	"secret-im/service/signalserver/cmd/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type DelDeviceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDelDeviceLogic(ctx context.Context, svcCtx *svc.ServiceContext) DelDeviceLogic {
	return DelDeviceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DelDeviceLogic) DelDevice(r *http.Request, req types.DelDeviceReq) error {
	account, err := logic.GetSourceAccount(r, l.svcCtx.AccountsModel)
	if err != nil {
		return shared.Status(http.StatusUnauthorized, err.Error())
	}
	devicdID, err := strconv.ParseInt(req.DeviceId, 10, 64)
	if err != nil {
		return shared.Status(http.StatusBadRequest, err.Error())
	}
	account.RemoveDevice(devicdID)
	err = storage.AccountManager{}.Update(account)
	if err != nil {
		return shared.Status(http.StatusInternalServerError, err.Error())
	}

	err = storage.MessagesManager{}.ClearDevice(account.Number, devicdID)
	if err != nil {
		logx.Error("[Device::removeDevice] failed to clear device messages,",
			"  number:", account.Number,
			" device_id:", devicdID,
			" reason:", err,
		)
	}

	return nil
}
