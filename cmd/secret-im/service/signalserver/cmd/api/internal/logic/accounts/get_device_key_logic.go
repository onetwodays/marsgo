package logic

import (
	"context"

	"secret-im/service/signalserver/cmd/api/internal/svc"
	"secret-im/service/signalserver/cmd/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetDeviceKeyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetDeviceKeyLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetDeviceKeyLogic {
	return GetDeviceKeyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetDeviceKeyLogic) GetDeviceKey(req types.GetDeviceKeyReq) (*types.GetDeviceKeyRes, error) {
	// todo: add your logic here and delete this line
	_,account,err:=l.svcCtx.GetOneAccountByNumber(req.Number)
	if err!=nil{
		return nil,err
	}
	device:=account.Devices[0]
	resp,err:=l.svcCtx.KeysModel.FindManyFirst(req.Number,device.Id)
	if err!=nil{
		return nil, err
	}
	devices:=[]types.KeyDevices{types.KeyDevices{
		DeviceId: device.Id,
		RegistrationId: device.AccountAttributes.RegistrationId,
		SignedPreKey: device.SignedPreKey,
		PreKey: types.PreKey{
			KeyId: resp.Keyid,
			PublicKey: resp.Publickey,
		},

	}}



	return &types.GetDeviceKeyRes{
		IdentityKey: account.IdentityKey,
		Devices: devices,

	}, nil
}
