package handler

import (
	"github.com/tal-tech/go-zero/rest/httpx"
	"net/http"
	"secret-im/service/signalserver/cmd/api/internal/logic/keys"
	"secret-im/service/signalserver/cmd/api/internal/svc"
	shared "secret-im/service/signalserver/cmd/api/shared"
)

func GetSignedKeyHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		l := logic.NewGetSignedKeyLogic(r.Context(), ctx)
		appAccount :=r.Context().Value(shared.HttpReqContextAccountKey)
		if appAccount==nil{
			httpx.Error(w, shared.Status(http.StatusUnauthorized,"check basic auth fail ,may by the handler not use middle"))
			return
		}

		resp, err := l.GetSignedKey(r)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
