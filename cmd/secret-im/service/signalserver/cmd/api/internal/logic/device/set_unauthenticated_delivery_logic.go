package logic

import (
	"context"

	"github.com/tal-tech/go-zero/core/logx"
	"secret-im/service/signalserver/cmd/api/internal/svc"
)

type SetUnauthenticatedDeliveryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSetUnauthenticatedDeliveryLogic(ctx context.Context, svcCtx *svc.ServiceContext) SetUnauthenticatedDeliveryLogic {
	return SetUnauthenticatedDeliveryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SetUnauthenticatedDeliveryLogic) SetUnauthenticatedDelivery() error {
	// todo: add your logic here and delete this line

	return nil
}
