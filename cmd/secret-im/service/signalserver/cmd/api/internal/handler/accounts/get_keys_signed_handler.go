package handler

import (
	"net/http"
	"secret-im/service/signalserver/cmd/shared"

	"github.com/tal-tech/go-zero/rest/httpx"
	"secret-im/service/signalserver/cmd/api/internal/logic/accounts"
	"secret-im/service/signalserver/cmd/api/internal/svc"
)

func GetKeysSignedHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		number:=r.Header.Get(shared.PhoneNumberKey)
		l := logic.NewGetKeysSignedLogic(r.Context(), ctx)
		resp, err := l.GetKeysSigned(number)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
