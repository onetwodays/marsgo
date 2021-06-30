package handler

import (
	"net/http"
	"secret-im/service/signalserver/cmd/api/internal/entities"
	"secret-im/service/signalserver/cmd/shared"

	"github.com/tal-tech/go-zero/rest/httpx"
	"secret-im/service/signalserver/cmd/api/internal/logic/keys"
	"secret-im/service/signalserver/cmd/api/internal/svc"
)

func GetKeyCountHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		appAccount :=r.Context().Value(shared.HttpReqContextAccountKey)
		if appAccount==nil{
			httpx.Error(w, shared.Status(http.StatusUnauthorized,"check basic auth fail ,may by the handler not use middle"))
			return
		}

		l := logic.NewGetKeyCountLogic(r.Context(), ctx)
		resp, err := l.GetKeyCount(appAccount.(*entities.Account))
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
