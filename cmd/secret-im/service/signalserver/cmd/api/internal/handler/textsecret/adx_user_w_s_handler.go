package handler

import (
	"net/http"
	"secret-im/service/signalserver/cmd/api/chat"
	"secret-im/service/signalserver/cmd/shared"

	"secret-im/service/signalserver/cmd/api/internal/svc"
)

func AdxUserWSHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		adxName:= r.Header.Get(shared.HEADADXUSERNAME)

		chat.AdxWsConnectHandler(adxName,w,r)


		/**
		l := logic.NewAdxUserWSLogic(r.Context(), ctx)
		err := l.AdxUserWS()
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.Ok(w)
		}
		*/

	}
}
