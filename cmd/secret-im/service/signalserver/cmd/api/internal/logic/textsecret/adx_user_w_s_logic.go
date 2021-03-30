package logic

import (
	"context"
	"github.com/tal-tech/go-zero/core/logx"
	"secret-im/service/signalserver/cmd/api/internal/svc"
)

type AdxUserWSLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdxUserWSLogic(ctx context.Context, svcCtx *svc.ServiceContext) AdxUserWSLogic {
	return AdxUserWSLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdxUserWSLogic) AdxUserWS() error {
	// todo: add your logic here and delete this line
	// 完成ws握手

	return nil
}
