package handler

import (
	"github.com/tal-tech/go-zero/rest/httpx"
	"net/http"
	"secret-im/signalserver/internal/logic"
	"secret-im/signalserver/internal/svc"
)

func IndexHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		l := logic.NewIndexLogic(r.Context(), ctx)
		resp, err := l.Index()
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
