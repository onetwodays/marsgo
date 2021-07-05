package logic

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"secret-im/service/signalserver/cmd/api/internal/entities"
	"secret-im/service/signalserver/cmd/api/internal/svc"
	"secret-im/service/signalserver/cmd/api/internal/types"
	shared "secret-im/service/signalserver/cmd/api/shared"

	"github.com/tal-tech/go-zero/core/logx"
)

type DeliveryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeliveryLogic(ctx context.Context, svcCtx *svc.ServiceContext) DeliveryLogic {
	return DeliveryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeliveryLogic) Delivery(appAccount *entities.Account) (*types.DeliveryRes, error) {
	// todo: add your logic here and delete this line


	includeUuid := false
	data,err := l.svcCtx.CertificateGenerator.CreateFor(appAccount.Number,
		                                                appAccount.IdentityKey,
		                                                appAccount.UUID,
		                                                appAccount.AuthenticatedDevice.ID,
		                                                includeUuid)
	if err != nil {
		reason := fmt.Sprintf("%s",err.Error())
		logx.Error(reason)
		return nil, shared.Status(http.StatusInternalServerError,reason)
	}

	return &types.DeliveryRes{
		Certificate: base64.StdEncoding.EncodeToString(data),
	}, nil
}
