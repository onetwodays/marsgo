package handler

import (
	"github.com/gorilla/websocket"
	"github.com/tal-tech/go-zero/rest/httpx"
	"net/http"
	"secret-im/service/signalserver/cmd/api/chat"
	logic "secret-im/service/signalserver/cmd/api/internal/logic/textsecret"
	"secret-im/service/signalserver/cmd/api/internal/svc"
)

func AdxUserWSHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if websocket.IsWebSocketUpgrade(r){
			chat.WsConnectHandler(ctx)
		}else{
			l := logic.NewAdxUserWSLogic(r.Context(), ctx)
			err := l.AdxUserWS()
			if err != nil {
				httpx.Error(w, err)
			} else {
				httpx.Ok(w)
			}
		}

	}
}
