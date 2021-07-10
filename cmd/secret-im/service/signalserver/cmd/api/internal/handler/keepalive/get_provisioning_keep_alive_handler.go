package handler

import (
	"net/http"

	"github.com/tal-tech/go-zero/rest/httpx"
	"secret-im/service/signalserver/cmd/api/internal/logic/keepalive"
	"secret-im/service/signalserver/cmd/api/internal/svc"
)

func GetProvisioningKeepAliveHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		l := logic.NewGetProvisioningKeepAliveLogic(r.Context(), ctx)
		err := l.GetProvisioningKeepAlive(r)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
