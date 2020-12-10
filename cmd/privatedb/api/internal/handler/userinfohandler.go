package handler

import (
    "errors"
    "net/http"

    "github.com/tal-tech/go-zero/rest/httpx"
    "privatedb/api/internal/logic"
    "privatedb/api/internal/svc"
)

func userInfoHandler(ctx *svc.ServiceContext) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {

        l := logic.NewUserInfoLogic(r.Context(), ctx)
        //获取头
        userId := r.Header.Get("x-user-id")
        if len(userId) == 0 {
            httpx.Error(w, errors.New("lost head:x-user-id"))
        } else {
            resp, err := l.UserInfo(userId)
            if err != nil {
                httpx.Error(w, err)
            } else {
                httpx.OkJson(w, resp)
            }

        }

    }
}
