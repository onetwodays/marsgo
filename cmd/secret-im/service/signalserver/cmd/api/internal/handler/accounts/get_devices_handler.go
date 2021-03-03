package handler

import (
	"net/http"
	"secret-im/service/signalserver/cmd/shared"

	"github.com/tal-tech/go-zero/rest/httpx"
	"secret-im/service/signalserver/cmd/api/internal/logic/accounts"
	"secret-im/service/signalserver/cmd/api/internal/svc"
)

func GetDevicesHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		phoneNumber:=r.Header.Get(shared.PhoneNumberKey)
		l := logic.NewGetDevicesLogic(r.Context(), ctx)
		resp, err := l.GetDevices(phoneNumber)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
