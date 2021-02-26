## goctl rpc 命令介绍
https://www.yuque.com/tal-tech/go-zero/hlxlbt

##生成项目
cd cmd/secret-im/rpc/
goctl rpc new hello_rpc


## 使用goctl工具创建proto文件模板
goctl rpc template -o=hello.proto

## gotcl生成rpc服务
指定hello.proto生成rpc服务
goctl rpc proto -src hello.proto -dir .

protoc  -I=/home/zh/go/marsgo/cmd/secret-im/rpc/hello_rpc hello_rpc.proto --go_out=plugins=grpc:/home/zh/go/marsgo/cmd/secret-im/rpc/hello_rpc/hello_rpc
