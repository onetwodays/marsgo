package logic

import (
	"context"

	"secret-im/service/signalserver/cmd/api/internal/svc"
	"secret-im/service/signalserver/cmd/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetProfileKeyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetProfileKeyLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetProfileKeyLogic {
	return GetProfileKeyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetProfileKeyLogic) GetProfileKey(req types.GetProfileKeyReq) (*types.GetProfileKeyRes, error) {
	// todo: add your logic here and delete this line

	res,err:=l.svcCtx.ProfileKeyModel.FindOneByAccountName(req.AccountName)
	if err!=nil {
		return nil, err
	}

	dbAccount,err:=l.svcCtx.AccountsModel.FindOneByNumber(req.AccountName)
	if err!=nil{
		return nil,err
	}

	return &types.GetProfileKeyRes{
		Profilekey: res.ProfileKey,
		AccountName: req.AccountName,
		Uuid: dbAccount.Uuid,
	}, nil
}
