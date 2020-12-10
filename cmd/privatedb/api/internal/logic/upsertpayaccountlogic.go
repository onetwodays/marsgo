package logic

import (
    "context"
    "privatedb/api/model"

    "privatedb/api/internal/svc"
    "privatedb/api/internal/types"

    "github.com/tal-tech/go-zero/core/logx"
)

type UpsertPayAccountLogic struct {
    logx.Logger
    ctx    context.Context
    svcCtx *svc.ServiceContext
}

func NewUpsertPayAccountLogic(ctx context.Context, svcCtx *svc.ServiceContext) UpsertPayAccountLogic {
    return UpsertPayAccountLogic{
        Logger: logx.WithContext(ctx),
        ctx:    ctx,
        svcCtx: svcCtx,
    }
}

func (l *UpsertPayAccountLogic) UpsertPayAccount(req types.UpsertPayAccountReq) (*types.UpsertPayAccountReply, error) {
    // todo: add your logic here and delete this line

    paytype,err:=l.svcCtx.TPayTypeModel.FindOne(req.PaymentTypeId)
    if err!=nil{
        return nil, err
    }
    data := model.TPaymentAccount{
        Id:              req.Id,
        Username:        req.Username,
        PaymentTypeId:   req.PaymentTypeId,
        PaymentTypeName: paytype.Name,
        IsActived:       req.IsActived,
        Detail:          req.Detail,
    }

    ok := true

    //插入一条记录
    if req.Id == 0 {
        _, err = l.svcCtx.TPaymentAccountModel.Insert(data)

    } else {
        err = l.svcCtx.TPaymentAccountModel.Update(data)
    }
    if err != nil {
        ok = false
    }

    return &types.UpsertPayAccountReply{Ok: ok}, err
}
