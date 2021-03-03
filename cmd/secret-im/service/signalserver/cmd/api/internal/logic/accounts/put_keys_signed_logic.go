package logic

import (
	"context"
	"encoding/json"
	"secret-im/service/signalserver/cmd/model"

	"secret-im/service/signalserver/cmd/api/internal/svc"
	"secret-im/service/signalserver/cmd/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type PutKeysSignedLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPutKeysSignedLogic(ctx context.Context, svcCtx *svc.ServiceContext) PutKeysSignedLogic {
	return PutKeysSignedLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PutKeysSignedLogic) PutKeysSigned(req types.PutKeysSignedReq,number string) (*types.PutKeysSignedRes, error) {
	// todo: add your logic here and delete this line
	dbId,account,err:=l.svcCtx.GetOneAccountByNumber(number)
	if err!=nil{
		return nil, err
	}
	device:=&account.Devices[0]
	device.SignedPreKey.Prekey.KeyId = req.SignedPreKey.Prekey.KeyId
	device.SignedPreKey.Prekey.PublicKey=req.SignedPreKey.Prekey.PublicKey
	device.SignedPreKey.Signature = req.SignedPreKey.Signature
	accountStr,err:=json.Marshal(account)
	if err!=nil{
		return nil, err
	}
	err=l.svcCtx.UpdateDirectory(number,device.AccountAttributes.Voice,device.AccountAttributes.Video)
	if err!=nil{
		return nil, err
	}
	err= l.svcCtx.AccountsModel.Update(model.TAccounts{
		Id:dbId,
		Number: number,
		Data: string(accountStr),
	})
	if err!=nil{
		return nil, err
	}

	return &types.PutKeysSignedRes{}, nil
}
