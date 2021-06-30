package logic

import (
	"context"
	"secret-im/service/signalserver/cmd/api/internal/entities"

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

func (l *GetSignedKeyLogic) GetSignedKey(appAccount *entities.Account) (*types.SignedPrekey, error) {
	// todo: add your logic here and delete this line

	device :=appAccount.AuthenticatedDevice
	signedPreKey:=device.SignedPreKey

	return &types.SignedPrekey{
		PublicKey: signedPreKey.PublicKey,
		KeyId: signedPreKey.KeyId,
		Signature: signedPreKey.Signature,
	}, nil
}
