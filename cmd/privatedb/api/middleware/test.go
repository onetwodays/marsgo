package middleware

import (
    "github.com/tal-tech/go-zero/core/logx"
    "net/http"
)

func MidderwareDemoFunc(next  http.HandlerFunc) http.HandlerFunc  {
    return func(w http.ResponseWriter, r *http.Request) {
        logx.Info("全局中间件request ... ",r)
        next(w, r)
        logx.Info("全局中间件reponse ... ")
    }
}
