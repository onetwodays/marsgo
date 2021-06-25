package handler

import (
	"net/http"
	"secret-im/service/signalserver/cmd/api/internal/entities"
	"secret-im/service/signalserver/cmd/shared"

	"github.com/tal-tech/go-zero/rest/httpx"
	"secret-im/service/signalserver/cmd/api/internal/logic/certificate"
	"secret-im/service/signalserver/cmd/api/internal/svc"
)

func DeliveryHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//


		l := logic.NewDeliveryLogic(r.Context(), ctx)
		appAccount :=r.Context().Value(shared.HttpReqContextAccountKey)
		if appAccount==nil{
			httpx.Error(w, shared.Status(http.StatusUnauthorized,"check basic auth fail ,may by the handler not use middle"))
			return
		}


		resp, err := l.Delivery(appAccount.(*entities.Account))
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
