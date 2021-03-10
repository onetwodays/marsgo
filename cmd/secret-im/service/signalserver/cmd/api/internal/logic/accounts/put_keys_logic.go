package logic

import (
	"context"
	"encoding/json"
	"secret-im/service/signalserver/cmd/model"

	"secret-im/service/signalserver/cmd/api/internal/svc"
	"secret-im/service/signalserver/cmd/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type PutKeysLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPutKeysLogic(ctx context.Context, svcCtx *svc.ServiceContext) PutKeysLogic {
	return PutKeysLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PutKeysLogic) PutKeys(req types.PutKeysReq,number string) (*types.PutKeysResp, error) {
	// todo: add your logic here and delete this line
	dbId,account,err:=l.svcCtx.GetOneAccountByNumber(number)
	if err!=nil{
		return nil, err
	}
	device :=&account.Devices[0]
	updateAccount := false
	if req.SignedPreKey!=device.SignedPreKey{
		device.SignedPreKey=req.SignedPreKey
		updateAccount=true
	}
	if req.IdentityKey != account.IdentityKey{
		updateAccount=true
		account.IdentityKey = req.IdentityKey
	}
	if updateAccount{
		accountStr,err:=json.Marshal(account)
		if err!=nil{
			return nil, err
		}

		err=l.svcCtx.UpdateDirectory(number,device.AccountAttributes.Voice,device.AccountAttributes.Video)
		if err!=nil{
			return nil, err
		}

		err=l.svcCtx.AccountsModel.Update(model.TAccounts{
			Id:dbId,
			Number: number,
			Data: string(accountStr),
		})
		if err!=nil{
			return nil, err
		}

	}

	//更新keys
	l.svcCtx.KeysModel.DeleteMany(number,device.Id)
	for i:=range req.PreKeys{

		_,err:=l.svcCtx.KeysModel.Insert(model.TKeys{
			Number: number,
			Deviceid: device.Id,
			Publickey: req.PreKeys[i].PublicKey,
			Keyid: req.PreKeys[i].KeyId,
			LastResort: 0,

		})
		if err!=nil{
			return nil, err
		}
	}

	return &types.PutKeysResp{}, nil
}
