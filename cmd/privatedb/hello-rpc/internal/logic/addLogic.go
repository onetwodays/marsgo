package logic

import (
	"context"

	"privatedb/hello-rpc/hello"
	"privatedb/hello-rpc/internal/svc"

	"github.com/tal-tech/go-zero/core/logx"
)

type AddLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddLogic {
	return &AddLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddLogic) Add(in *hello.AddReq) (*hello.AddResp, error) {
	// todo: add your logic here and delete this line

	return &hello.AddResp{}, nil
}
