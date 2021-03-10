package logic

import (
	"context"

	"secret-im/service/signalserver/cmd/rpc/bookstore/bookstore"
	"secret-im/service/signalserver/cmd/rpc/bookstore/internal/svc"

	"github.com/tal-tech/go-zero/core/logx"
)

type CheckLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCheckLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckLogic {
	return &CheckLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CheckLogic) Check(in *bookstore.CheckReq) (*bookstore.CheckResp, error) {
	// todo: add your logic here and delete this line
	resp,err:=l.svcCtx.BookModel.FindOne(in.Book)
	if err!=nil{
		return nil, err
	}

	return &bookstore.CheckResp{
		Found: true,
		Price: int32(resp.Price),
	}, nil
}
