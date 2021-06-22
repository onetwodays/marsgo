package logic

import (
	"context"

	"secret-im/service/signalserver/cmd/api/internal/svc"
	"secret-im/service/signalserver/cmd/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetKeyCountLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetKeyCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetKeyCountLogic {
	return GetKeyCountLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetKeyCountLogic) GetKeyCount(adxName string) (*types.PreKeyCount, error) {
	// todo: add your logic here and delete this line
	count,err:=l.svcCtx.KeysModel.CountKey(adxName,1)
	println("mmmmmmmmmmmmmmmm",count)
	if err!=nil{
		return nil, err
	}
	return &types.PreKeyCount{Count: int(*count)}, nil
}
