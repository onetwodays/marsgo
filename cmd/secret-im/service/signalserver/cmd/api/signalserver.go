package main

import (
	"fmt"
	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/rest/httpx"
	"net/http"
	"os"
	"path/filepath"
	"secret-im/common"
	"secret-im/service/signalserver/cmd/api/chat"
	"secret-im/service/signalserver/cmd/api/middleware"
	"secret-im/service/signalserver/cmd/shared"
	"strings"

	"github.com/tal-tech/go-zero/rest"
	"github.com/tal-tech/go-zero/rest/router"

	"secret-im/service/signalserver/cmd/api/config"
	"secret-im/service/signalserver/cmd/api/internal/handler"
	"secret-im/service/signalserver/cmd/api/internal/svc"
)



var isStartWss = true
func main() {
	//url 不存在时,个性化提示
	rt := router.NewRouter()
	rt.SetNotFoundHandler(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		//这里内容可以定制
		w.Write([]byte("服务器开小差了,这里可定制"))
	}))


	ctx := svc.NewServiceContext(config.AppConfig)
	server := rest.MustNewServer(config.AppConfig.RestConf, rest.WithRouter(rt)) //url 不存在时会报 服务器开小差了,这里可定制
	//server := rest.MustNewServer(c.RestConf) //url 不存在时默认会报 404 page not found
	defer server.Stop()
	// 全局中间件
	server.Use(middleware.GlobalMWLogFunc) //global middleware

	//静态文件服务
	registerDirHandlers(server,ctx)

	// 注册路由组件
	handler.RegisterHandlers(server, ctx) // handle api



	/*
	httpx.Error(...) 调用会先经过自定义的 errorHandler 处理再返回。
	• errorHandler 返回的 int 作为 http status code 返回客户端
	• 如果 errorHandler 返回的 interface{} 是 error 类型的话，
	那么会直接用 err.Error() 的内容以非 json 的格式返回客户端，
	不是 error 的话，那么会 marshal 成 json 再返回
	*/

	httpx.SetErrorHandler(func(err error) (int, interface{}) {
		switch e:=err.(type) {
		case *shared.CodeError:
			return http.StatusOK,e.Data()
		default:
			return http.StatusInternalServerError,err
		}
	})


    //websocket server.调试使用，可以通过网页看到ws结果，生产环境要关闭
    if isStartWss{

		server.AddRoute(rest.Route{
			Method: http.MethodGet,
			Path: "/ws",
			Handler: chat.WsConnectHandler(ctx),
		})
	}


	fmt.Printf("Starting server at %s:%d...\n", config.AppConfig.Host, config.AppConfig.Port)
	server.Start()


}


func registerDirHandlers(engine *rest.Server, serverCtx *svc.ServiceContext) {

	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	logx.Info("current dir:",exPath)

	//这里注册
	dirlevel := []string{":1", ":2", ":3", ":4", ":5", ":6", ":7", ":8"}
	patern := "/deploy/" //url用的
	dirpath := filepath.Join(exPath,"static") //服务所在目录的static目录
	logx.Info("static file dir:",dirpath)
	for i := 1; i < len(dirlevel); i++ {
		path := patern + strings.Join(dirlevel[:i], "/")
		//最后生成 /asset
		engine.AddRoute(
			rest.Route{
				Method:  http.MethodGet,
				Path:    path,
				Handler: common.Dirhandler(patern,dirpath),
			})

		logx.Infof("register dir  %s  %s", path,dirpath)
	}
}





