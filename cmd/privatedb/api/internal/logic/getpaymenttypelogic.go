package logic

import (
	"context"

	"privatedb/api/internal/svc"
	"privatedb/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetPaymentTypeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPaymentTypeLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetPaymentTypeLogic {
	return GetPaymentTypeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPaymentTypeLogic) GetPaymentType() (*types.GetPaymentTypeReply, error) {
	// todo: add your logic here and delete this line
	resp,err:=l.svcCtx.TPayTypeModel.FindAll()
	if err!=nil{
		return nil, err
	}

	return &types.GetPaymentTypeReply{
		Total: len(resp),
		List: resp,
	}, nil
}
