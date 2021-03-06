
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
