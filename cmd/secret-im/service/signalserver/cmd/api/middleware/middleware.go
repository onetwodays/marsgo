package middleware

import (
	"github.com/tal-tech/go-zero/core/logx"
	"net/http"
)

func GlobalMWLogFunc(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logx.Info("request ... ",r.URL.String())


		next(w,r)


		logx.Info("reponse ... ")
	}
}
