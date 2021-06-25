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

func (l *GetKeysLogic) GetKeys(req types.GetKeysReq) (*types.GetKeysResx, error) {
	// todo: add your logic here and delete this line

	return &types.GetKeysResx{}, nil
}
