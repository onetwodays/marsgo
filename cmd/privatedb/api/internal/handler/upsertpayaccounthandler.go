package handler

import (
	"net/http"

	"privatedb/api/internal/logic"
	"privatedb/api/internal/svc"
	"privatedb/api/internal/types"

	"github.com/tal-tech/go-zero/rest/httpx"
)

func UpsertPayAccountHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpsertPayAccountReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewUpsertPayAccountLogic(r.Context(), ctx)
		resp, err := l.UpsertPayAccount(req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
