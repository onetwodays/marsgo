package middleware

import (
	"fmt"
	"net/http"
	"net/http/httputil"
)

func GlobalMWLogFunc(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//logx.Info("request ... ",r.URL.String())

		dump, _ := httputil.DumpRequest(r, true)
		fmt.Println("===============")
		fmt.Println(string(dump))
		fmt.Println("===============")


		next(w,r)


		//logx.Info("reponse ... ")
	}
}
