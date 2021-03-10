package logic

import (
	"context"
	"secret-im/service/signalserver/cmd/rpc/bookstore/bookstore"

	"secret-im/service/signalserver/cmd/api/internal/svc"
	"secret-im/service/signalserver/cmd/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type AddLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) AddLogic {
	return AddLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddLogic) Add(req types.AddReq) (*types.AddResp, error) {
	// todo: add your logic here and delete this line
	resp,err:=l.svcCtx.BookStoreClient.Add(l.ctx,&bookstore.AddReq{
		Book: req.Book,
		Price: int32(req.Price),
	})
	if err!=nil{
		return nil, err
	}

	return &types.AddResp{
		Ok: resp.Ok,
	}, nil
}
