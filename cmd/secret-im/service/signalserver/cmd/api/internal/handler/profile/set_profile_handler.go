package handler

import (
	"net/http"
	"secret-im/service/signalserver/cmd/api/internal/logic/profile"
	"secret-im/service/signalserver/cmd/api/internal/svc"
	"secret-im/service/signalserver/cmd/api/internal/types"

	"github.com/tal-tech/go-zero/rest/httpx"
)

func SetProfileHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateProfileRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewSetProfileLogic(r.Context(), ctx)
		resp, err := l.SetProfile(r,req)
		if err != nil {
			httpx.Error(w, err)
		} else {


			httpx.OkJson(w, resp)
		}
	}
}
