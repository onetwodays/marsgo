package handler

import (
	"net/http"
	"secret-im/shared"

	"github.com/tal-tech/go-zero/rest/httpx"
	"secret-im/signalserver/internal/logic"
	"secret-im/signalserver/internal/svc"
)

func userInfoHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		l := logic.NewUserInfoLogic(r.Context(), ctx)
		userId:=r.Header.Get("x-user-id")
		resp, err := l.UserInfo(userId)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, shared.NewOkResponse(resp))
		}
	}
}
