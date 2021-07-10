package logic

import (
	"context"

	"github.com/tal-tech/go-zero/core/logx"
	"secret-im/service/signalserver/cmd/api/internal/svc"
)

type VerifyDeviceTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewVerifyDeviceTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) VerifyDeviceTokenLogic {
	return VerifyDeviceTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *VerifyDeviceTokenLogic) VerifyDeviceToken() error {
	// todo: add your logic here and delete this line

	return nil
}
