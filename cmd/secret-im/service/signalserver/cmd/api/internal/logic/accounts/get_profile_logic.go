package logic

import (
	"context"

	"secret-im/service/signalserver/cmd/api/internal/svc"
	"secret-im/service/signalserver/cmd/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetProfileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetProfileLogic {
	return GetProfileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetProfileLogic) GetProfile(req types.GetProfileReq) (*types.GetProfileResp, error) {
	// todo: add your logic here and delete this line
	_,account,err:=l.svcCtx.GetOneAccountByNumber(req.Number)
	if err!=nil{
		return nil, err
	}

	return &types.GetProfileResp{
		Name: account.Name,
		Avatar: account.Avatar,
		IdentityKey: account.IdentityKey,
	}, nil
}
