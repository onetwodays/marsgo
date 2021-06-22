package logic

import (
	"context"

	"github.com/tal-tech/go-zero/core/logx"
	"secret-im/service/signalserver/cmd/api/internal/svc"
)

type GetKeepAliveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetKeepAliveLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetKeepAliveLogic {
	return GetKeepAliveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetKeepAliveLogic) GetKeepAlive() error {
	// todo: add your logic here and delete this line

	return nil
}
