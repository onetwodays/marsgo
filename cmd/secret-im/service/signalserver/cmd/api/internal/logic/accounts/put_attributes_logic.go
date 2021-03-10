package logic

import (
	"context"
	"encoding/json"
	"secret-im/service/signalserver/cmd/model"

	"secret-im/service/signalserver/cmd/api/internal/svc"
	"secret-im/service/signalserver/cmd/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type PutAttributesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPutAttributesLogic(ctx context.Context, svcCtx *svc.ServiceContext) PutAttributesLogic {
	return PutAttributesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PutAttributesLogic) PutAttributes(req types.PutAccountAttributesReq) (*types.PutAccountAttributesRes, error) {
	// todo: add your logic here and delete this line
	accountModel,err:= l.svcCtx.AccountsModel.FindOneByNumber(req.PhoneNumber)
	if err!=nil {
		return nil, err
	}
	data:=accountModel.Data
	account:=types.Account{}
	err = json.Unmarshal([]byte(data),&account)
	account.Devices[0].AccountAttributes = req.AccountAttributes
	accountStr,err:=json.Marshal(account)
	if err!=nil{
		return nil, err
	}
	err=l.svcCtx.AccountsModel.Update(model.TAccounts{
		Id:accountModel.Id,
		Number: accountModel.Number,
		Data: string(accountStr),
	})
	if err!=nil{
		return nil, err
	}

	return &types.PutAccountAttributesRes{Account: string(accountStr)}, nil
}
