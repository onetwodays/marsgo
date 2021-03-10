package logic

import (
	"context"
	"fmt"

	"secret-im/service/signalserver/cmd/rpc/hello/hello_rpc"
	"secret-im/service/signalserver/cmd/rpc/hello/internal/svc"

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

func (l *PingLogic) Ping(in *hello_rpc.Request) (*hello_rpc.Response, error) {
	// todo: add your logic here and delete this line

	return &hello_rpc.Response{
		Pong: fmt.Sprint("hello %s",in.Ping),
	}, nil
}
