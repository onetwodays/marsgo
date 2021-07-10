package handler

import (
	"github.com/tal-tech/go-zero/rest/httpx"
	"net/http"
	"secret-im/service/signalserver/cmd/api/internal/logic/messages"
	"secret-im/service/signalserver/cmd/api/internal/svc"
	shared "secret-im/service/signalserver/cmd/api/shared"
)

func GetMsgsHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		l := logic.NewGetMsgsLogic(r.Context(), ctx)
		adxName := r.Header.Get(shared.HEADADXUSERNAME)
		resp, err := l.GetMsgs(adxName, 1)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
