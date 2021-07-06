package logic

import (
	"context"
	"net/http"
	"secret-im/service/signalserver/cmd/api/internal/entities"
	"secret-im/service/signalserver/cmd/api/shared"

	"secret-im/service/signalserver/cmd/api/internal/svc"
	"secret-im/service/signalserver/cmd/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetWhoamiLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetWhoamiLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetWhoamiLogic {
	return GetWhoamiLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetWhoamiLogic) GetWhoami(r *http.Request) (*types.AccountCreationResult, error) {
	appAccount := r.Context().Value(shared.HttpReqContextAccountKey)
	if appAccount == nil {
		reason := "check basic auth fail ,may by the handler not use middle"
		logx.Error(reason)
		return nil, shared.Status(http.StatusUnauthorized, reason)
	}
	account := appAccount.(*entities.Account)

	return &types.AccountCreationResult{UUID: account.UUID}, nil
}
