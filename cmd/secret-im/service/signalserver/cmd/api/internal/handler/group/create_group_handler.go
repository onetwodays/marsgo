package handler

import (
	"net/http"
	"secret-im/service/signalserver/cmd/api/shared"
	"secret-im/service/signalserver/cmd/api/textsecure"
	"secret-im/service/signalserver/cmd/api/util"

	"github.com/tal-tech/go-zero/rest/httpx"
	"secret-im/service/signalserver/cmd/api/internal/logic/group"
	"secret-im/service/signalserver/cmd/api/internal/svc"
)

func CreateGroupHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req textsecure.Group
		if err:=util.PFParse(r,&req);err!=nil{
			httpx.Error(w,shared.Status(http.StatusInternalServerError,err.Error()))
			return
		}


		l := logic.NewCreateGroupLogic(r.Context(), ctx)
		resp, err := l.CreateGroup()
		if err != nil {
			httpx.Error(w, err)
		} else {
			err:=util.Render(w,resp.Obj)
			if err!=nil{
				httpx.Error(w,shared.Status(http.StatusInternalServerError,err.Error()))
			}
		}
	}
}
