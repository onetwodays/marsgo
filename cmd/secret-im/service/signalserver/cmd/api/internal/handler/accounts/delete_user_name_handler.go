package handler

import (
	"net/http"

	"github.com/tal-tech/go-zero/rest/httpx"
	"secret-im/service/signalserver/cmd/api/internal/logic/accounts"
	"secret-im/service/signalserver/cmd/api/internal/svc"
)

func DeleteUserNameHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		l := logic.NewDeleteUserNameLogic(r.Context(), ctx)
		err := l.DeleteUserName(r)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
