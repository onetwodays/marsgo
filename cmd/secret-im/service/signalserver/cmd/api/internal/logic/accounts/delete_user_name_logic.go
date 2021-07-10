package logic

import (
	"context"
	"net/http"
	"secret-im/service/signalserver/cmd/api/internal/logic"
	"secret-im/service/signalserver/cmd/api/internal/storage"
	"secret-im/service/signalserver/cmd/api/shared"

	"github.com/tal-tech/go-zero/core/logx"
	"secret-im/service/signalserver/cmd/api/internal/svc"
)

type DeleteUserNameLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteUserNameLogic(ctx context.Context, svcCtx *svc.ServiceContext) DeleteUserNameLogic {
	return DeleteUserNameLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteUserNameLogic) DeleteUserName(r *http.Request) error {
	account,err:= logic.GetSourceAccount(r,l.svcCtx.AccountsModel)
	if err!=nil{
		return shared.Status(http.StatusUnauthorized,err.Error())
	}

	err=storage.UsernamesManager{}.Delete(account.UUID)
	if err!=nil{
		return shared.Status(http.StatusInternalServerError,err.Error())
	}

	return nil
}
