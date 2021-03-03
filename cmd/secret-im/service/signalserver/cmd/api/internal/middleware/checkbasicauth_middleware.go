package middleware

import (
	"encoding/base64"
	"github.com/tal-tech/go-zero/rest/httpx"
	"net/http"
	"secret-im/service/signalserver/cmd/shared"
	"strings"
)

type CheckBasicAuthMiddleware struct {
}

func NewCheckBasicAuthMiddleware() *CheckBasicAuthMiddleware {
	return &CheckBasicAuthMiddleware{}
}

func (m *CheckBasicAuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO generate middleware implement function, delete after code implementation
		s := strings.SplitN(r.Header.Get("x-Authorization"), " ", 2)
		if len(s) != 2 {
			httpx.Error(w, errorUserInfo) //todo
			return
		}

		b, err := base64.StdEncoding.DecodeString(s[1])
		if err != nil {
			httpx.Error(w, errorUserInfo) //todo
			return
		}

		pair := strings.SplitN(string(b), ":", 2)
		if len(pair) != 2 {
			httpx.Error(w, errorUserInfo) //todo
			return
		}
		r.Header.Set(shared.PhoneNumberKey,pair[0]) //设置一个新头
		r.Header.Set(shared.PasswordKey,pair[1]) //设置一个新头

		// Passthrough to next handler if need
		next(w, r)
	}
}
