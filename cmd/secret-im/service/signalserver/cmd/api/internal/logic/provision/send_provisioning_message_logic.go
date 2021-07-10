package logic

import (
	"context"
	"encoding/base64"
	"net/http"
	"secret-im/service/signalserver/cmd/api/internal/logic"
	"secret-im/service/signalserver/cmd/api/internal/svc/push"
	"secret-im/service/signalserver/cmd/api/shared"

	"secret-im/service/signalserver/cmd/api/internal/svc"
	"secret-im/service/signalserver/cmd/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type SendProvisioningMessageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSendProvisioningMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) SendProvisioningMessageLogic {
	return SendProvisioningMessageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendProvisioningMessageLogic) SendProvisioningMessage(r *http.Request,req types.ProvisioningMessage) error {
	account,err:= logic.GetSourceAccount(r,l.svcCtx.AccountsModel)
	if err!=nil{
		return shared.Status(http.StatusUnauthorized,err.Error())
	}
	body,err:=base64.StdEncoding.DecodeString(req.Body)
	if err!=nil{
		return shared.Status(http.StatusBadRequest,err.Error())
	}
	//todo limit
	address:=push.ProvisioningAddress{}
	address.Number=req.Destination
	delivered, err := l.svcCtx.PushSender.GetWebsocketSender().SendProvisioningMessage(address, body)
	if err!=nil || !delivered{
		logx.Error("[Provisioning] failed to send provisioning message",
			" source:",      account.Number,
			" destination:", req.Destination,
			" address:",     address.Serialize())
		return shared.Status(http.StatusNotFound,"[Provisioning] failed to send provisioning message")
	}

	return nil
}
