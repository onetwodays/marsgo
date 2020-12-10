package main

import (
    "flag"
    "fmt"
    "github.com/tal-tech/go-zero/rest/httpx"
    "net/http"

    "privatedb/api/internal/config"
    "privatedb/api/internal/handler"
    "privatedb/api/internal/svc"

    "github.com/tal-tech/go-zero/core/conf"
    "github.com/tal-tech/go-zero/rest"
    "privatedb/api/middleware"
)

// 返回的结构体，json格式的body
type Message struct {
    Code int    `json:"code"`
    Desc string `json:"desc"`
}

// 定义错误处理函数
/*
httpx.Error(...) 调用会先经过自定义的 errorHandler 处理再返回。
• errorHandler 返回的 int 作为 http status code 返回客户端
• 如果 errorHandler 返回的 interface{} 是 error 类型的话，
那么会直接用 err.Error() 的内容以非 json 的格式返回客户端，
不是 error 的话，那么会 marshal 成 json 再返回
 */
func errorHandler(err error) (int, interface{}) {
    return http.StatusConflict, Message{
        Code: -1,
        Desc: err.Error(),
    }
}

var configFile = flag.String("f", "etc/privatedb-api.yaml", "the config file")

func main() {
    flag.Parse()

    var c config.Config
    conf.MustLoad(*configFile, &c)

    ctx := svc.NewServiceContext(c)
    server := rest.MustNewServer(c.RestConf)
    defer server.Stop()

    server.Use(
        middleware.MidderwareDemoFunc, //全局的中间件
    )
    // 设置错误处理函数
    httpx.SetErrorHandler(errorHandler)

    handler.RegisterHandlers(server, ctx)

    fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
    server.Start()
}
