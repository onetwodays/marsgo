package handler

import (
	"net/http"
	"secret-im/service/signalserver/cmd/shared"

	"secret-im/service/signalserver/cmd/api/internal/logic/textsecret_messages"
	"secret-im/service/signalserver/cmd/api/internal/svc"
	"secret-im/service/signalserver/cmd/api/internal/types"

	"github.com/tal-tech/go-zero/rest/httpx"
)

func PutMsgsHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PutMessagesReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewPutMsgsLogic(r.Context(), ctx)
		resp, err := l.PutMsgs(req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, shared.NewOkResponse(resp))
		}
	}
}
