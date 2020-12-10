package handler

import (
	"net/http"

	"github.com/tal-tech/go-zero/rest/httpx"
	"privatedb/api/internal/logic"
	"privatedb/api/internal/svc"
)

func GetPaymentTypeHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		l := logic.NewGetPaymentTypeLogic(r.Context(), ctx)
		resp, err := l.GetPaymentType()
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
