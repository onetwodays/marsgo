# 在安装go get go.etcd.io/etcd/clientv3时出错

>> github.com/coreos/etcd/clientv3/balancer/resolver/endpoint
>> ../../pkg/mod/github.com/coreos/etcd@v3.3.18+incompatible/clientv3/balancer/resolver/endpoint/endpoint.go:114:78: undefined: resolver.BuildOption
>> ../../pkg/mod/github.com/coreos/etcd@v3.3.18+incompatible/clientv3/balancer/resolver/endpoint/endpoint.go:182:31: undefined: resolver.ResolveNowOption

解决方法
<br>
1.将grpc版本替换成v1.26.0版本

 1. 修改依赖为v1.26.0
go mod edit -require=google.golang.org/grpc@v1.26.0

 2.下载v1.26.0版本的grpc

go get -u -x google.golang.org/grpc@v1.26.0

2.执行下面的命令也可以<br>
go mod edit -replace github.com/coreos/bbolt@v1.3.4=go.etcd.io/bbolt@v1.3.4
go mod edit -replace google.golang.org/grpc@v1.29.1=google.golang.org/grpc@v1.26.0
go mod tidy

# 使用 go.etcd.io/etcd/mvcc/mvccpb 出错
> invalid operation: ev.Type == "go.etcd.io/etcd/mvcc/mvccpb".DELETE 
(mismatched types "github.com/coreos/etcd/mvcc/mvccpb".Event_EventType 
and "go.etcd.io/etcd/mvcc/mvccpb".Event_EventType)

> invalid operation: ev.Type == "go.etcd.io/etcd/mvcc/mvccpb".PUT
 (mismatched types "github.com/coreos/etcd/mvcc/mvccpb".Event_EventType 
and "go.etcd.io/etcd/mvcc/mvccpb".Event_EventType)

1. 问题原因
>import 的 go.etcd.io/etcd/clientv3 引用的是 go.etcd.io/etcd/mvcc/mvccpb, 
但是在  
（if ev.Type == mvccpb.DELETE） 和
（if ev.Type == mvccpb.PUT && string(ev.Kv.Key) == key）
应用了 github.com/coreos/etcd/mvcc/mvccpb

2. 解决
导入 github.com/coreos/etcd/mvcc/mvccpb 不要导入go.etcd.io/etcd/mvcc/mvccpb