package handler

import (
	"net/http"
	"secret-im/service/signalserver/cmd/shared"

	"secret-im/service/signalserver/cmd/api/internal/logic/accounts"
	"secret-im/service/signalserver/cmd/api/internal/svc"
	"secret-im/service/signalserver/cmd/api/internal/types"

	"github.com/tal-tech/go-zero/rest/httpx"
)

func GetSmsCodeHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetSmsCodeReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewGetSmsCodeLogic(r.Context(), ctx)
		resp, err := l.GetSmsCode(req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, shared.NewOkResponse(resp))
		}
	}
}
