package logic

import (
	"context"
	"net/http"
	"secret-im/service/signalserver/cmd/api/internal/auth"
	"secret-im/service/signalserver/cmd/api/internal/auth/helper"
	"secret-im/service/signalserver/cmd/api/internal/logic"
	"secret-im/service/signalserver/cmd/api/internal/storage"
	"secret-im/service/signalserver/cmd/api/shared"
	"strconv"

	"secret-im/service/signalserver/cmd/api/internal/svc"
	"secret-im/service/signalserver/cmd/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetDeviceKeysLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetDeviceKeysLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetDeviceKeysLogic {
	return GetDeviceKeysLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetDeviceKeysLogic) GetDeviceKeys(r *http.Request, req types.GetKeysReqx) (*types.GetKeysRes, error) {



	header := r.Header.Get(helper.UNIDENTIFIED)
	accessKey, _ := auth.NewAnonymous(header)
	appAccount,err:= logic.GetSourceAccount(r,l.svcCtx.AccountsModel)
	if err!=nil{
		return nil,shared.Status(http.StatusUnauthorized,err.Error())
	}
	if accessKey == nil && appAccount==nil {
		return nil, shared.Status(http.StatusUnauthorized, "")
	}
	// 获取目标用户
	targetName := auth.NewAmbiguousIdentifier(req.Identifier)
	_, target, err := storage.AccountManager{}.Get(targetName)
	if  err != nil {
		return nil, shared.Status(http.StatusNotFound, err.Error())
	}
	code, ok := helper.OptionalAccess{}.VerifyDevices(appAccount, accessKey, target, req.DeviceId)
	if !ok {
		logx.Error("VerifyDevices fail,code=",code)
		//return nil, shared.Status(code, "VerifyDevices fail")
	}

	if appAccount != nil {
		/*
			key := account.Number + "__" + target.Number + "." + deviceID
				if !c.rateLimiters.PreKeysLimiter.Validate(key, 1) {
					respond.Status(w, http.StatusRequestEntityTooLarge)
					return
				}
		*/
	}

	// 获取设备密钥
	var deviceId int64 = 0
	if req.DeviceId != "*" {
		id, _ := strconv.ParseInt(req.DeviceId, 10, 64)
		deviceId = id
	}

	keys, err := l.svcCtx.KeysModel.FindMany(target.Number, deviceId)
	if err != nil {
		return nil, shared.Status(http.StatusNotFound, err.Error())
	}

	deleteIds := make([]int64, 0)
	devices := make([]types.PreKeyResponseItemx, 0)

	for _, device := range target.Devices {
		//if device.IsEnabled() && (req.DeviceId == "*" || device.ID == deviceId) {
		if  req.DeviceId == "*" || device.ID == deviceId {
			signedPreKey := device.SignedPreKey
			for _, keyRecode := range keys {
				if keyRecode.DeviceId == device.ID {
					if signedPreKey != nil  {
						devices = append(devices, types.PreKeyResponseItemx{
							DeviceId:       device.ID,
							RegistrationId: int64(device.RegistrationID),
							SignedPrekey:   *signedPreKey,
							PreKey:         types.PreKey{
								KeyId:     keyRecode.KeyId,
								PublicKey: keyRecode.PublicKey,
							},
						})
						deleteIds = append(deleteIds, keyRecode.Id)
					}
				}
			}
		}
	}

	if len(devices) == 0 {
		return nil, shared.Status(http.StatusOK, "")
	}

	//删除已经查到的数据
	err = l.svcCtx.KeysModel.DeleteMany(deleteIds)
	if err != nil {
		return nil, shared.Status(http.StatusInternalServerError, err.Error())
	}
	return &types.GetKeysRes{
		IdentityKey: target.IdentityKey,
		Devices:     devices,
	}, nil
}
