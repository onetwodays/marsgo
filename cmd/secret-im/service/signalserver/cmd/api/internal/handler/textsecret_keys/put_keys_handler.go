package handler

import (
	"net/http"
	"secret-im/service/signalserver/cmd/shared"

	"secret-im/service/signalserver/cmd/api/internal/logic/textsecret_keys"
	"secret-im/service/signalserver/cmd/api/internal/svc"
	"secret-im/service/signalserver/cmd/api/internal/types"

	"github.com/tal-tech/go-zero/rest/httpx"
)

func PutKeysHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PutKeysReqx
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}
		adxName:= r.Header.Get(shared.HEADADXUSERNAME)

		l := logic.NewPutKeysLogic(r.Context(), ctx)
		err := l.PutKeys(adxName,req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w,shared.NewOkResponse(nil))
		}
	}
}
