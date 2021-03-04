package handler

import (
	"net/http"
	"secret-im/service/signalserver/cmd/shared"

	"secret-im/service/signalserver/cmd/api/internal/logic/msgs"
	"secret-im/service/signalserver/cmd/api/internal/svc"
	"secret-im/service/signalserver/cmd/api/internal/types"

	"github.com/tal-tech/go-zero/rest/httpx"
)

func PutMsgsSendHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PutMsgsSendReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		number:=r.Header.Get(shared.PhoneNumberKey)

		l := logic.NewPutMsgsSendLogic(r.Context(), ctx)
		resp, err := l.PutMsgsSend(req,number)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
