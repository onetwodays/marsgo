package logic

import (
    "context"
    "privatedb/api/model"

    "privatedb/api/internal/svc"
    "privatedb/api/internal/types"

    "github.com/tal-tech/go-zero/core/logx"
)

type AddmsgLogic struct {
    logx.Logger
    ctx    context.Context
    svcCtx *svc.ServiceContext
}

func NewAddmsgLogic(ctx context.Context, svcCtx *svc.ServiceContext) AddmsgLogic {
    return AddmsgLogic{
        Logger: logx.WithContext(ctx),
        ctx:    ctx,
        svcCtx: svcCtx,
    }
}

func (l *AddmsgLogic) Addmsg(req types.AddmsgReq) (*types.AddmsgReply, error) {
    // todo: add your logic here and delete this line
    logx.Infof("context:",req.Content)
    _, err := l.svcCtx.TMsgModel.Insert(model.TMsg{
        Pair:     req.Pair,
        DealId:   req.Dealid,
        Sender:   req.Sender,
        Receiver: req.Receiver,
        Context:  req.Content,
        Status: 1,
    })
    ok := true
    if err != nil {
        ok = false

    }

    return &types.AddmsgReply{Ok: ok}, err
}
