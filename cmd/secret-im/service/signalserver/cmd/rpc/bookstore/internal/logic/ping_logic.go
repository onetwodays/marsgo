package logic

import (
	"context"

	"secret-im/service/signalserver/cmd/rpc/bookstore/bookstore"
	"secret-im/service/signalserver/cmd/rpc/bookstore/internal/svc"

	"github.com/tal-tech/go-zero/core/logx"
)

type PingLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PingLogic {
	return &PingLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PingLogic) Ping(in *bookstore.Request) (*bookstore.Response, error) {
	// todo: add your logic here and delete this line

	return &bookstore.Response{
		Pong: "successful",
	}, nil
}
