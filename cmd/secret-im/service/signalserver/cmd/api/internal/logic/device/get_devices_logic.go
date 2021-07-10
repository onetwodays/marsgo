package logic

import (
	"context"
	"net/http"
	"secret-im/service/signalserver/cmd/api/internal/logic"
	"secret-im/service/signalserver/cmd/api/shared"

	"secret-im/service/signalserver/cmd/api/internal/svc"
	"secret-im/service/signalserver/cmd/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetDevicesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetDevicesLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetDevicesLogic {
	return GetDevicesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetDevicesLogic) GetDevices(r *http.Request) (*types.DeviceInfoList, error) {
	account, err := logic.GetSourceAccount(r, l.svcCtx.AccountsModel)
	if err != nil {
		return nil, shared.Status(http.StatusUnauthorized, err.Error())
	}
	// 返回设备列表
	devices := make([]types.DeviceInfo, 0)
	for _, device := range account.Devices {
		devices = append(devices, types.DeviceInfo{
			ID:       device.ID,
			Name:     device.Name,
			LastSeen: device.LastSeen,
			Created:  device.Created,
		})
	}

	return &types.DeviceInfoList{Devices: devices}, nil
}
