package middleware

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"strings"
)

var exinclude = map[string]struct{}{"/v1/keepalive": {}}

func GlobalMWLogFunc(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, ok := exinclude[strings.ToLower(r.RequestURI)]; !ok {
			dump, _ := httputil.DumpRequest(r, true)
			fmt.Println("===============")
			fmt.Println(string(dump))
			fmt.Println("===============")
		}
		next(w, r)


		//logx.Info("reponse ... ")
	}
}
