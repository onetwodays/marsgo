package handler

import (
	"net/http"
	"secret-im/service/signalserver/cmd/api/chat"

	"secret-im/service/signalserver/cmd/api/internal/svc"
)

func WwsConnectHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return chat.WsConnectHandler(ctx)
}
