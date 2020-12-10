

```
mkdir hello-rpc
goctl rpc template -o hello.proto
goctl rpc proto -src hello.proto -dir . --style goZero //protoc  -I=/home/zh/svn/eosio.contracts/golang/privatedb/hello-rpc hello.proto --go_out=plugins=grpc:/home/zh/svn/eosio.contracts/golang/privatedb/hello-rpc/hello

```