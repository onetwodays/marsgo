package logic

import (
	"context"

	"github.com/tal-tech/go-zero/core/logx"
	"secret-im/service/signalserver/cmd/api/internal/svc"
)

type RwsConnectLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRwsConnectLogic(ctx context.Context, svcCtx *svc.ServiceContext) RwsConnectLogic {
	return RwsConnectLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RwsConnectLogic) RwsConnect() error {
	// todo: add your logic here and delete this line

	return nil
}
