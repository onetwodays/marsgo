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

type GetSignedKeyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetSignedKeyLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetSignedKeyLogic {
	return GetSignedKeyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSignedKeyLogic) GetSignedKey(r *http.Request) (*types.SignedPrekey, error) {
	account,err:= logic.GetSourceAccount(r,l.svcCtx.AccountsModel)
	if err!=nil{
		return nil, shared.Status(http.StatusUnauthorized,err.Error())
	}

	spk := account.AuthenticatedDevice.SignedPreKey

	return &types.SignedPrekey{
		PublicKey: spk.PublicKey,
		KeyId:     spk.KeyId,
		Signature: spk.Signature,
	}, nil
}
