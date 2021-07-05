#0. 安装更新
```
GOPROXY=https://goproxy.cn/,direct go get -u github.com/tal-tech/go-zero
GO111MODULE=on GOPROXY=https://goproxy.cn/,direct go get -u github.com/tal-tech/go-zero/tools/goctl
```

#1 .创建项目
https://www.yuque.com/tal-tech/go-zero/rslrhx
```
goctl api new greet
cd greet
go mod init
go mod tidy
```

#2.编写业务代码：

+ api 文件定义了服务对外暴露的路由，可参考 api 规范
+ 可以在 servicecontext.go 里面传递依赖给 logic，比如 mysql, redis 等
+ 在 api 定义的 get/post/put/delete 等请求对应的 logic 里增加业务处理逻辑
+ 可以根据 api 文件生成前端需要的 Java, TypeScript, Dart, JavaScript 代码
```
goctl api java -api greet.api -dir greet
goctl api dart -api greet.api -dir greet
```
+ 生成model
```
goctl model mysql datasource -url="hopexdev:devhopex@tcp(127.0.0.1:3306)/pb" -table="user" -dir ./model
```

# 3.启动
```
go run privatedb.go  -f ./etc/privatedb-api.yaml
```

#4.编译
``
go build -o pb.exe privatedb.go
``


#5.命令
goctl api [go/java/ts] [-api user/user.api] [-dir ./src]  -style [go_zero/GoZero/gozero]
api 后面接生成的语言，现支持go/java/typescript
-api 自定义api所在路径
-dir 自定义生成目录
https://github.com/tal-tech/go-zero/blob/master/tools/goctl/config/readme.md


{"@timestamp":"2021-07-05T20:58:04.810+08","level":"info","content":"消息已经推送至redis中，key is {liguozhen114}:1"}
{"@timestamp":"2021-07-05T20:58:04.810+08","level":"info","content":"[Message] send message success delivered:true online:false destination:liguozhen114 FetchesMessages:true timestamp:1625489273757"}
{"@timestamp":"2021-07-05T20:58:04.810+08","level":"info","content":"200 - /v1/messages/liguozhen114 - 192.168.0.211:53478 -  - 4.4ms","trace":"0d34de50f6014662","span":"0"}
{"@timestamp":"2021-07-05T20:58:04.810+08","level":"info","content":"websocket经http处理后的响应是:{\"needsSync\":false}"}
{"@timestamp":"2021-07-05T20:58:04.810+08","level":"info","content":"account:liguozhen114timestamp:1625489273757[Authenticated] deliver message"}
{"@timestamp":"2021-07-05T20:58:04.810+08","level":"info","content":"[Client] send request to liguozhen114 id:3694562916659513808 verb:PUT path:/api/v1/message"}
{"@timestamp":"2021-07-05T20:58:04.825+08","level":"info","content":"[Client] recv response from liguozhen114 WebSocketResponseMessage id:3694562916659513808"}
{"@timestamp":"2021-07-05T20:58:04.825+08","level":"info","content":"[Authenticated] deliver message ackaccount:liguozhen114timestamp:1625489273757"}
{"@timestamp":"2021-07-05T20:58:04.825+08","level":"info","content":" [Authenticated] send delivery receiptaccount:liguozhen114 timestamp:1625489273757"}
{"@timestamp":"2021-07-05T20:58:04.826+08","level":"info","duration":"1.0ms","content":"sql query: liguozhen115"}
