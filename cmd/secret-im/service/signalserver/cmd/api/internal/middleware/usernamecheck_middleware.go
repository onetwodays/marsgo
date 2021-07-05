package middleware

import (
	"github.com/tal-tech/go-zero/rest/httpx"
	"net/http"
	shared "secret-im/service/signalserver/cmd/api/shared"
)


type UserNameCheckMiddleware struct {
}

func NewUserNameCheckMiddleware() *UserNameCheckMiddleware {
	return &UserNameCheckMiddleware{}
}

func (m *UserNameCheckMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO generate middleware implement function, delete after code implementation

		adxName:= r.Header.Get(shared.HEADADXUSERNAME)
		if len(adxName)==0{
			httpx.Error(w, shared.ErrAdxHeadInvalid)
			return
		}

		// go-zero从jwt token解析后会将用户生成token时传入的kv原封不动的放在http.Request的Context中，因此我们可以通过Context就可以拿到你想要的值
		jwtName:=r.Context().Value(shared.JWTADXUSERNAME)


		if jwtName != adxName {
			httpx.Error(w, shared.ErrAdxCheck)
			return
		}



		// Passthrough to next handler if need
		next(w, r)
	}
}
