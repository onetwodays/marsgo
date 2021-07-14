package handler

import (
	"net/http"

	"secret-im/service/signalserver/cmd/api/internal/logic/channel"
	"secret-im/service/signalserver/cmd/api/internal/svc"
	"secret-im/service/signalserver/cmd/api/internal/types"

	"github.com/tal-tech/go-zero/rest/httpx"
)

func CreateChannelHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ChannelCreationInfo
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewCreateChannelLogic(r.Context(), ctx)
		resp, err := l.CreateChannel(r,req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
