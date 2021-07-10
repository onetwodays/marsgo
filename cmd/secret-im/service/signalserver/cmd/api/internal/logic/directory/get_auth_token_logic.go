package logic

import (
	"context"
	"net/http"
	"secret-im/service/signalserver/cmd/api/internal/logic"
	"secret-im/service/signalserver/cmd/api/shared"

	"secret-im/service/signalserver/cmd/api/internal/svc"
	"secret-im/service/signalserver/cmd/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetAuthTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAuthTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetAuthTokenLogic {
	return GetAuthTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAuthTokenLogic) GetAuthToken(r *http.Request) (*types.ExternalServiceCredentials, error) {
	account,err:= logic.GetSourceAccount(r,l.svcCtx.AccountsModel)
	if err!=nil{
		return nil,shared.Status(http.StatusUnauthorized,err.Error())
	}

	return l.svcCtx.DirectoryCredentialsGenerator.GenerateFor(account.Number), nil
}
