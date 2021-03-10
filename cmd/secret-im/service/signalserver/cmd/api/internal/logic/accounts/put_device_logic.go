package logic

import (
	"context"
	"encoding/json"
	"secret-im/service/signalserver/cmd/api/util"
	"secret-im/service/signalserver/cmd/model"

	"secret-im/service/signalserver/cmd/api/internal/svc"
	"secret-im/service/signalserver/cmd/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type PutDeviceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPutDeviceLogic(ctx context.Context, svcCtx *svc.ServiceContext) PutDeviceLogic {
	return PutDeviceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PutDeviceLogic) PutDevice(req types.PutDeviceReq,number,passwd string) (*types.PutDeviceRes, error) {
	// todo: add your logic here and delete this line
	dbId,account,err:=l.svcCtx.GetOneAccountByNumber(number)
	if err!=nil{
		return nil, err
	}
	salt := util.GenSalt()
	device:=types.Device{
		Id:int64(len(account.Devices)),
		Salt:salt,
		AccountAttributes:req.AccountAttributes,
		AuthToken: util.GenAuthKey(salt,passwd),
	}
	account.Devices=append(account.Devices,device)
	accountStr,err:=json.Marshal(account)
	if err!=nil{
		return nil, err
	}
	data:= string(accountStr)
	err = l.svcCtx.AccountsModel.Update(model.TAccounts{
		Id:dbId,
		Number: number,
		Data:data,
	})
	if err!=nil{
		return nil, err
	}


	return &types.PutDeviceRes{
		Account: data,
	}, nil
}
