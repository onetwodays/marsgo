package handler

import (
	"net/http"
	"secret-im/service/signalserver/cmd/shared"

	"github.com/tal-tech/go-zero/rest/httpx"
	"secret-im/service/signalserver/cmd/api/internal/logic/website"
	"secret-im/service/signalserver/cmd/api/internal/svc"
)

func UserInfoHandler(ctx *svc.ServiceContext) http.HandlerFunc {
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
