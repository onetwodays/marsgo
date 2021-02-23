package logic

import (
	"context"

	"secret-im/signalserver/internal/svc"
	"secret-im/signalserver/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type SignalserverLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSignalserverLogic(ctx context.Context, svcCtx *svc.ServiceContext) SignalserverLogic {
	return SignalserverLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SignalserverLogic) Signalserver(req types.Request) (*types.Response, error) {
	// todo: add your logic here and delete this line

	return &types.Response{}, nil
}
