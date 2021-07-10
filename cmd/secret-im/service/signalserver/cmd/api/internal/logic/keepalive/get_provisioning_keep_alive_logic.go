package logic

import (
	"context"
	"net/http"

	"github.com/tal-tech/go-zero/core/logx"
	"secret-im/service/signalserver/cmd/api/internal/svc"
)

type GetProvisioningKeepAliveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetProvisioningKeepAliveLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetProvisioningKeepAliveLogic {
	return GetProvisioningKeepAliveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetProvisioningKeepAliveLogic) GetProvisioningKeepAlive(r *http.Request) error {
	// todo: add your logic here and delete this line

	return nil
}
