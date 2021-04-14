package handler

import (
	"github.com/gorilla/websocket"
	"github.com/tal-tech/go-zero/rest/httpx"
	"net/http"
	"secret-im/service/signalserver/cmd/api/chat"
	logic "secret-im/service/signalserver/cmd/api/internal/logic/textsecret"
	"secret-im/service/signalserver/cmd/api/internal/svc"
	"secret-im/service/signalserver/cmd/shared"
)

func AdxUserWSHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if websocket.IsWebSocketUpgrade(r){
			adxName:= r.Header.Get(shared.HEADADXUSERNAME)
			chat.AdxWsConnectHandler(adxName,w,r)
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
