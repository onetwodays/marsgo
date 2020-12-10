package handler

import (
	"net/http"

	"privatedb/api/internal/logic"
	"privatedb/api/internal/svc"
	"privatedb/api/internal/types"

	"github.com/tal-tech/go-zero/rest/httpx"
)

func registerHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RegisterReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewRegisterLogic(r.Context(), ctx)
		err := l.Register(req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
