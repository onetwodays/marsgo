package handler

import (
	"net/http"
	"secret-im/service/signalserver/cmd/api/internal/logic/textsecret"
	"secret-im/service/signalserver/cmd/api/internal/svc"
	"secret-im/service/signalserver/cmd/api/internal/types"
	shared "secret-im/service/signalserver/cmd/api/shared"

	"github.com/tal-tech/go-zero/rest/httpx"
)

func AdxUserLoginHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AdxUserLoginReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}
		userAgent := r.Header.Get("User-Agent")

		l := logic.NewAdxUserLoginLogic(r.Context(),ctx)
		resp, err := l.AdxUserLogin(req,userAgent)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, shared.NewOkResponse(resp))
		}
	}
}
