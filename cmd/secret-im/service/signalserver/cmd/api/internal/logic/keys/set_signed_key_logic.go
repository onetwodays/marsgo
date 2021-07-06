package logic

import (
	"context"
	"net/http"
	"secret-im/service/signalserver/cmd/api/internal/entities"
	"secret-im/service/signalserver/cmd/api/internal/storage"
	"secret-im/service/signalserver/cmd/api/internal/svc"
	"secret-im/service/signalserver/cmd/api/internal/types"
	shared "secret-im/service/signalserver/cmd/api/shared"

	"github.com/tal-tech/go-zero/core/logx"
)

type SetSignedKeyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSetSignedKeyLogic(ctx context.Context, svcCtx *svc.ServiceContext) SetSignedKeyLogic {
	return SetSignedKeyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SetSignedKeyLogic) SetSignedKey(r *http.Request,req types.SignedPrekey) error {
	appAccount := r.Context().Value(shared.HttpReqContextAccountKey)
	if appAccount == nil {
		reason := "check basic auth fail ,may by the handler not use middle"
		logx.Error(reason)
		return shared.Status(http.StatusUnauthorized, reason)
	}
	account := appAccount.(*entities.Account)

	account.AuthenticatedDevice.SignedPreKey = & types.SignedPrekey{
		PublicKey:req.PublicKey,
		Signature: req.Signature,
		KeyId: req.KeyId,
	}
	err := storage.AccountManager{}.Update(account)
	if err !=nil{
		return shared.Status(http.StatusInternalServerError,err.Error())
	}

	return nil
}
