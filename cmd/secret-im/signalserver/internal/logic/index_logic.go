package logic

import (
	"context"
	"time"

	"secret-im/signalserver/internal/svc"
	"secret-im/signalserver/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type IndexLogic struct {
	logx.Logger
	ctx    context.Context  // 来自HttpRequest的context
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

	return &types.IndexReply{Resp: time.Now().Local().Format(`2006-01-02 15:04:05`)}, nil
}
