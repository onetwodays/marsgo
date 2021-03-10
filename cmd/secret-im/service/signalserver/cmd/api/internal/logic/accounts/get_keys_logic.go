package logic

import (
	"context"

	"secret-im/service/signalserver/cmd/api/internal/svc"
	"secret-im/service/signalserver/cmd/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetKeysLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetKeysLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetKeysLogic {
	return GetKeysLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetKeysLogic) GetKeys(number string) (*types.GetKeysRes, error) {
	// todo: add your logic here and delete this line
	_,account,err:=l.svcCtx.GetOneAccountByNumber(number)
	if err!=nil{
		return nil,err
	}
	device:=&account.Devices[0]

	resp,err:=l.svcCtx.KeysModel.FindMany(number,device.Id)
	if err!=nil{
		return nil, err
	}

	return &types.GetKeysRes{
		List: resp,
		Total: len(resp),
	}, nil
}
