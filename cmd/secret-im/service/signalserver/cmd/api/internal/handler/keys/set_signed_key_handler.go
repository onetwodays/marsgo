package handler

import (
	"net/http"
	"secret-im/service/signalserver/cmd/api/internal/logic/keys"
	"secret-im/service/signalserver/cmd/api/internal/svc"
	"secret-im/service/signalserver/cmd/api/internal/types"
	shared "secret-im/service/signalserver/cmd/api/shared"

	"github.com/tal-tech/go-zero/rest/httpx"
)

func SetSignedKeyHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SignedPrekey
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		appAccount :=r.Context().Value(shared.HttpReqContextAccountKey)
		if appAccount==nil{
			httpx.Error(w, shared.Status(http.StatusUnauthorized,"check basic auth fail ,may by the handler not use middle"))
			return
		}

		l := logic.NewSetSignedKeyLogic(r.Context(), ctx)
		err := l.SetSignedKey(r,req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
