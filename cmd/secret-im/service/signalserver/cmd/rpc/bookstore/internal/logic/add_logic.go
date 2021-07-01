package logic

import (
	"context"
	"secret-im/service/signalserver/cmd/api/internal/model"

	"secret-im/service/signalserver/cmd/rpc/bookstore/bookstore"
	"secret-im/service/signalserver/cmd/rpc/bookstore/internal/svc"

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

func (l *AddLogic) Add(in *bookstore.AddReq) (*bookstore.AddResp, error) {
	// todo: add your logic here and delete this line
	_,err:=l.svcCtx.BookModel.Insert(model.Book{
		Book: in.Book,
		Price: int64(in.Price),
	})
	if err!=nil{
		return nil, err
	}

	return &bookstore.AddResp{
		Ok:true,
	}, nil
}
