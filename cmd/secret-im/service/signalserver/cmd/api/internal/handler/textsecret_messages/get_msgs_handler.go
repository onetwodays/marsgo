package handler

import (
	"net/http"
	"secret-im/service/signalserver/cmd/shared"

	"github.com/tal-tech/go-zero/rest/httpx"
	"secret-im/service/signalserver/cmd/api/internal/logic/textsecret_messages"
	"secret-im/service/signalserver/cmd/api/internal/svc"
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
