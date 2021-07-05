package logic

import (
	"context"
	"net/http"
	"secret-im/service/signalserver/cmd/api/internal/entities"
	"secret-im/service/signalserver/cmd/api/internal/svc"
	"secret-im/service/signalserver/cmd/api/internal/types"
	shared "secret-im/service/signalserver/cmd/api/shared"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetKeyCountLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetKeyCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetKeyCountLogic {
	return GetKeyCountLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetKeyCountLogic) GetKeyCount(appAccount *entities.Account) (*types.PreKeyCountx, error) {
	// todo: add your logic here and delete this line
	count,err:=l.svcCtx.KeysModel.CountKey(appAccount.Number,appAccount.AuthenticatedDevice.ID)
	if err!=nil{
		return nil, shared.Status(http.StatusInternalServerError,err.Error())
	}
	return &types.PreKeyCountx{Count: int(*count)}, nil
}
