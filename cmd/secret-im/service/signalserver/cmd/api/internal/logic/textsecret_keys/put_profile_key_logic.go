package logic

import (
	"context"
	"secret-im/service/signalserver/cmd/model"

	"secret-im/service/signalserver/cmd/api/internal/svc"
	"secret-im/service/signalserver/cmd/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)


type PutProfileKeyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPutProfileKeyLogic(ctx context.Context, svcCtx *svc.ServiceContext) PutProfileKeyLogic {
	return PutProfileKeyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PutProfileKeyLogic) PutProfileKey(req types.PutProfileKeyReq) error {
	// todo: add your logic here and delete this line
	l.svcCtx.ProfileKeyMap[req.AccountName]=req.Profilekey
	data:=&model.TProfilekey{
		AccountName: req.AccountName,
		ProfileKey: req.Profilekey,
	}

	_,err:=l.svcCtx.ProfileKeyModel.Insert(*data)
	return err
}
