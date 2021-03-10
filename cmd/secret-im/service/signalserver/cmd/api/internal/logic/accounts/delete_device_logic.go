package logic

import (
	"context"
	"encoding/json"
	"errors"
	"secret-im/service/signalserver/cmd/api/internal/svc"
	"secret-im/service/signalserver/cmd/api/internal/types"
	"secret-im/service/signalserver/cmd/model"

	"github.com/tal-tech/go-zero/core/logx"
)

type DeleteDeviceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteDeviceLogic(ctx context.Context, svcCtx *svc.ServiceContext) DeleteDeviceLogic {
	return DeleteDeviceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteDeviceLogic) DeleteDevice(req types.DelDeviceReq,number string) (*types.DelDeviceRes, error) {
	// todo: add your logic here and delete this line

    dbId,account,err:=l.svcCtx.GetOneAccountByNumber(number)
	if err!=nil{
		return nil, err
	}
	if req.DeviceId >=int64(len(account.Devices)){
		return nil, errors.New("device id invaild")
	}
	simpleDevice:=types.SimpleDevice{}

	for i:=range account.Devices{
		if account.Devices[i].Id == req.DeviceId{
			simpleDevice.Id=account.Devices[i].Id
			simpleDevice.Name=account.Devices[i].AccountAttributes.Name
			simpleDevice.LastSeen= account.Devices[i].LastSeen
			simpleDevice.Created= account.Devices[i].Created
			account.Devices = append(account.Devices[:i],account.Devices[i+1:]...)
			break
		}
	}
	accountStr, _ := json.Marshal(account)
	//更新之
	err= l.svcCtx.AccountsModel.Update(	model.TAccounts{
		Id:dbId,
		Number: account.Number,
		Data: string(accountStr),
	})
	if err!=nil{
		return nil, err
	}

	return &types.DelDeviceRes{
		Device:simpleDevice,

	}, nil
}
