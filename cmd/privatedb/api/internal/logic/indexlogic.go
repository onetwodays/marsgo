package logic

import (
	"context"
	"time"

	"privatedb/api/internal/svc"
	"privatedb/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type IndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) IndexLogic {
	return IndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *IndexLogic) Index() (*types.IndexReply, error) {
	// todo: add your logic here and delete this line

	return &types.IndexReply{
		Resp: "hello world@"+time.Now().Format("2006-01-02T15:04:05.999999999Z07:00"),
	}, nil
}
