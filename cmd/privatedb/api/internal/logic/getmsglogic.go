package logic

import (
	"context"

	"privatedb/api/internal/svc"
	"privatedb/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetmsgLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetmsgLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetmsgLogic {
	return GetmsgLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetmsgLogic) Getmsgs(req types.GetmsgsReq) (*types.GetmsgsReply, error) {
	// todo: add your logic here and delete this line
	resp,err:=l.svcCtx.TMsgModel.FindMany(req)
	if err!=nil{
		return nil, err
	}

	return &types.GetmsgsReply{
		Total: len(resp),
		PageSize: req.PageSize,
		Data: resp,
	}, nil
}
