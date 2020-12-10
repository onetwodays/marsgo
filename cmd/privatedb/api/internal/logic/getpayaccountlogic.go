package logic

import (
	"context"

	"privatedb/api/internal/svc"
	"privatedb/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetPayAccountLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPayAccountLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetPayAccountLogic {
	return GetPayAccountLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPayAccountLogic) GetPayAccount(req types.GetPayAccountReq) (*types.GetPayAccountReply, error) {
	// todo: add your logic here and delete this line
	dbresult,err:=l.svcCtx.TPaymentAccountModel.FindMany(req)
	if err!=nil{
		return nil, err
	}

	return &types.GetPayAccountReply{
		Total: len(dbresult),
		List: dbresult,
	}, nil
}
