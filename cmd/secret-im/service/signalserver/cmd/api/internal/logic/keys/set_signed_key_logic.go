package logic

import (
	"context"
	"net/http"
	"secret-im/service/signalserver/cmd/api/internal/entities"
	"secret-im/service/signalserver/cmd/api/internal/storage"
	"secret-im/service/signalserver/cmd/shared"

	"secret-im/service/signalserver/cmd/api/internal/svc"
	"secret-im/service/signalserver/cmd/api/internal/types"

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

func (l *SetSignedKeyLogic) SetSignedKey(req types.SignedPrekey,appAccount *entities.Account) error {
	// todo: add your logic here and delete this line

	appAccount.AuthenticatedDevice.SignedPreKey = & types.SignedPrekey{
		PublicKey:req.PublicKey,
		Signature: req.Signature,
		KeyId: req.KeyId,
	}
	err := storage.AccountManager{}.Update(appAccount,l.svcCtx.AccountsModel)
	if err !=nil{
		return shared.Status(http.StatusInternalServerError,err.Error())
	}

	return nil
}
