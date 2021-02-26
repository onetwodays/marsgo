## 创建项目
cd rpc
mkdir -r bookstore/proto
cd   bookstore/proto
goctl rpc template -o bookstore.proto

## 编辑上步生成的bookstore.proto

## 编译bookstore.proto
goctl rpc proto -src bookstore.proto -dir .. -style go_zero

## 拦截器
https://github.com/grpc-ecosystem/go-grpc-middleware
