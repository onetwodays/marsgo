package main

import (
	"fmt"
	"github.com/tal-tech/go-zero/core/logx"
	"net/http"
	"strings"

	"github.com/tal-tech/go-zero/rest"
	"github.com/tal-tech/go-zero/rest/router"

	"secret-im/signalserver/chat"
	"secret-im/signalserver/config"
	"secret-im/signalserver/internal/handler"
	"secret-im/signalserver/internal/svc"
	"secret-im/signalserver/middleware"
)



var isStartWss = true
func main() {

	rt := router.NewRouter()
	rt.SetNotFoundHandler(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		//这里内容可以定制
		w.Write([]byte("服务器开小差了,这里可定制"))
	}))


	ctx := svc.NewServiceContext(config.AppConfig)
	server := rest.MustNewServer(config.AppConfig.RestConf, rest.WithRouter(rt)) //url 不存在时会报 服务器开小差了,这里可定制
	//server := rest.MustNewServer(c.RestConf) //url 不存在时默认会报 404 page not found
	defer server.Stop()
	server.Use(middleware.GlobalMWWebsocketFunc) //global middleware

	handler.RegisterHandlers(server, ctx) // handle api
	registerDirHandlers(server,ctx)


    //websocket server
    if isStartWss{
    	go chat.HubRun()
		http.HandleFunc("/", filehandler("./static/chat.html"))
		http.HandleFunc("/ws",wshandler())
		go func() {
			err:=http.ListenAndServe(config.AppConfig.WssAddress,nil)
			if err!=nil{
				fmt.Printf("\"Starting WebsocketServer at %s happend error:",err.Error())
			}else{
				fmt.Printf("Starting WebsocketServer at %s...\n", config.AppConfig.WssAddress)
			}
		}()
	}




	server.AddRoute(rest.Route{
		Method: http.MethodGet,
		Path: "/ws",
		Handler: wshandler(),
	})
	fmt.Printf("Starting server at %s:%d...\n", config.AppConfig.Host, config.AppConfig.Port)
	server.Start()


}
func wshandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		chat.ServeWs(w,req)
	}
}

func registerDirHandlers(engine *rest.Server, serverCtx *svc.ServiceContext) {

	//这里注册
	dirlevel := []string{":1", ":2", ":3", ":4", ":5", ":6", ":7", ":8"}
	patern := "/deploy/" //url用的
	dirpath := "./static" //服务所在目录的static目录
	for i := 1; i < len(dirlevel); i++ {
		path := patern + strings.Join(dirlevel[:i], "/")
		//最后生成 /asset
		engine.AddRoute(
			rest.Route{
				Method:  http.MethodGet,
				Path:    path,
				Handler: dirhandler(patern,dirpath),
			})

		logx.Infof("register dir  %s  %s", path,dirpath)
	}
}





