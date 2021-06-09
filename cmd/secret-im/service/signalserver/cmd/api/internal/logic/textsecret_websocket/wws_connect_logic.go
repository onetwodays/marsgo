package logic

import (
	"context"

	"secret-im/service/signalserver/cmd/api/internal/svc"
	"secret-im/service/signalserver/cmd/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type WwsConnectLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWwsConnectLogic(ctx context.Context, svcCtx *svc.ServiceContext) WwsConnectLogic {
	return WwsConnectLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WwsConnectLogic) WwsConnect(req types.WriteWsConnReq) error {
	// todo: add your logic here and delete this line

	return nil
}
