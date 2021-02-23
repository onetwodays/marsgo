package middleware
// 路由组 Middleware 配置
import (
	"encoding/json"
	"fmt"

	"github.com/tal-tech/go-zero/rest/httpx"
	"net/http"
	"secret-im/shared"
)

var (
	errorUserInfo = shared.NewCodeError(shared.GetUserInfoFailed,shared.CodeErrorMap[shared.GetUserInfoFailed])
	errorAuthDeny = shared.NewCodeError(shared.AuthDeny,shared.CodeErrorMap[shared.AuthDeny])
)

const (
	userKey=`x-user-id`
)

type UsercheckMiddleware struct {
}

func NewUsercheckMiddleware() *UsercheckMiddleware {
	return &UsercheckMiddleware{}
}

func (m *UsercheckMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO generate middleware implement function, delete after code implementation
		userId:= r.Header.Get(userKey)
		userInt,err:=json.Number(userId).Int64()
		if err!=nil{
			httpx.Error(w,errorUserInfo)
			return
		}
        // go-zero从jwt token解析后会将用户生成token时传入的kv原封不动的放在http.Request的Context中，因此我们可以通过Context就可以拿到你想要的值
		jwtUserId:=r.Context().Value("userId")
		jwtInt, err := json.Number(fmt.Sprintf("%v", jwtUserId)).Int64()
		if err != nil {
			httpx.Error(w, errorUserInfo)
			return
		}

		if jwtInt != userInt {
			httpx.Error(w, errorAuthDeny)
			return
		}

		// Passthrough to next handler if need
		next(w, r)
	}
}
