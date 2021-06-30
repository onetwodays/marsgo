package logic

import (
	"context"

	"secret-im/service/signalserver/cmd/api/internal/svc"
	"secret-im/service/signalserver/cmd/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type WSConnectLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWSConnectLogic(ctx context.Context, svcCtx *svc.ServiceContext) WSConnectLogic {
	return WSConnectLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WSConnectLogic) WSConnect(req types.WsConnReq) error {
	// todo: add your logic here and delete this line

	return nil
}
