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

func (l *GetKeyCountLogic) GetKeyCount() (*types.PreKeyCountx, error) {
	// todo: add your logic here and delete this line

	return &types.PreKeyCountx{}, nil
}
