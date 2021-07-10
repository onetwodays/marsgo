package handler

import (
	"net/http"

	"github.com/tal-tech/go-zero/rest/httpx"
	"secret-im/service/signalserver/cmd/api/internal/logic/device"
	"secret-im/service/signalserver/cmd/api/internal/svc"
)

func CreateDeviceTokenHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		l := logic.NewCreateDeviceTokenLogic(r.Context(), ctx)
		resp, err := l.CreateDeviceToken(r)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
