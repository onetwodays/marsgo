package logic

import (
	"context"
	"encoding/json"
	"net/http"
	"secret-im/service/signalserver/cmd/api/internal/entities"
	"secret-im/service/signalserver/cmd/api/internal/logic"
	"secret-im/service/signalserver/cmd/api/shared"

	"secret-im/service/signalserver/cmd/api/internal/svc"
	"secret-im/service/signalserver/cmd/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)
// 最大设备数量
const MaxDevices = 6

// 超出设备限制
type LimitExceeded struct {
	Current int `json:"current"`
	Max     int `json:"max"`
}

type CreateDeviceTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateDeviceTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) CreateDeviceTokenLogic {
	return CreateDeviceTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateDeviceTokenLogic) CreateDeviceToken(r *http.Request) (*types.VerificationCode, error) {
	account, err := logic.GetSourceAccount(r, l.svcCtx.AccountsModel)
	if err != nil {
		return nil,shared.Status(http.StatusUnauthorized, err.Error())
	}
	//todo limit
	if account.GetEnabledDeviceCount()>=MaxDevices{
		l:=LimitExceeded{
			Current: len(account.Devices),
			Max:     MaxDevices,
		}
		jsb,_:=json.Marshal(&l)
		return nil,shared.Status(http.StatusLengthRequired,string(jsb))
	}
	if account.AuthenticatedDevice.ID!=entities.DeviceMasterID{
		return nil, shared.Status(http.StatusUnauthorized,"account.AuthenticatedDevice.ID != storage.DeviceMasterID")
	}

	// 生成设备令牌
	/*
	verificationCode := generateVerificationCode()
	storedVerificationCode := types.VerificationCode{
		Code:      verificationCode,
		Timestamp: pkgUtils.CurrentTimeMillis(),
	}*/



	return &types.VerificationCode{}, nil
}



