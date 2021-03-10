package logic

import (
	"context"
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

func (l *GetDevicesLogic) GetDevices(number string) (*types.GetDevicesRes, error) {
	// todo: add your logic here and delete this line
	_,account,err:=l.svcCtx.GetOneAccountByNumber(number)
	if err!=nil{
		return nil, err
	}
	total:=len(account.Devices)
	list:=make([]types.SimpleDevice,total,total)
	for i:=range account.Devices{
		simpleDevice:=types.SimpleDevice{
			Id:account.Devices[i].Id,
			Name:account.Devices[i].AccountAttributes.Name,
			LastSeen: account.Devices[i].LastSeen,
			Created: account.Devices[i].Created,
		}
		list = append(list,simpleDevice)
	}


	return &types.GetDevicesRes{
		Total: total,
		List: list,

	}, nil
}
