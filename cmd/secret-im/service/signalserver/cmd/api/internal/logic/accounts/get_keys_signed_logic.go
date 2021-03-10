package logic

import (
	"context"

	"secret-im/service/signalserver/cmd/api/internal/svc"
	"secret-im/service/signalserver/cmd/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetKeysSignedLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetKeysSignedLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetKeysSignedLogic {
	return GetKeysSignedLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetKeysSignedLogic) GetKeysSigned(number string) (*types.GetKeysSignedRes, error) {
	// todo: add your logic here and delete this line
	_,account,err:=l.svcCtx.GetOneAccountByNumber(number)
	if err!=nil{
		return nil, err
	}
	device:=&account.Devices[0]



	return &types.GetKeysSignedRes{
		SignedPreKey: device.SignedPreKey,
	}, nil
}
