# net/trace

## 项目简介
1. 提供Trace的接口规范
2. 提供 trace 对Tracer接口的实现，供业务接入使用

## 接入示例
1. 启动接入示例
    ```go
    trace.Init(traceConfig) // traceConfig is Config object with value.
    ```
2. 配置参考
    ```toml
    [tracer]
    network = "unixgram"
    addr = "/var/run/dapper-collect/dapper-collect.sock"
    ```

## 测试
1. 执行当前目录下所有测试文件，测试所有功能

## 思考???
1. 每个span需要的基本信息何时生成？
2. 哪些信息需要随着服务调用传递给服务提供方？
3. 什么时候发送span至zipkin 服务端？
4. 以何种方式发送span?


## Trace
>一次服务调用追踪链路，由一组Span组成。需在web总入口处生成TraceID，并确保在当前请求上下文里能访问到

## Annotation
表示某个时间点发生的Event
> Event类型

> cs：Client Send 请求
> sr：Server Receive到请求
> ss：Server 处理完成、并Send Response
> cr：Client Receive 到响应
> 什么时候生成
客户端发送Request、接受到Response、服务器端接受到Request、发送 Response时生成。
>Annotation属于某个Span，需把新生成的Annotation添加到当前上下文里Span的annotations数组里

## Span
   表示一次完整RPC调用，是由一组Annotation和BinaryAnnotation组成。
   是追踪服务调用的基本结构，
   多span形成树形结构组合成一次Trace追踪记录。Span是有父子关系的，
   比如：Client A、Client A -> B、B ->C、C -> D、分别会产生4个Span。
   Client A接收到请求会时生成一个Span A、Client A -> B发请求时会再生成一个Span A-B，
   并且Span A是 Span A-B的父节点
   
   什么时候生成
   
   服务接受到 Request时，若当前Request没有关联任何Span，便生成一个Span，包括：Span ID、TraceID
   向下游服务发送Request时，需生成一个Span，并把新生成的Span的父节点设置成上一步生成的Span
   
## 服务之间需传递的信息
   Trace的基本信息需在上下游服务之间传递，如下信息是必须的：
   
   Trace ID：起始(根)服务生成的TraceID
   Span ID：调用下游服务时所生成的Span ID
   Parent Span ID：父Span ID
   Is Sampled：是否需要采样
   Flags：告诉下游服务，是否是debug Reqeust
   
## Trace Tree组成
      一个完整Trace 由一组Span组成，这一组Span必须具有相同的TraceID；
      Span具有父子关系，处于子节点的Span必须有parent_id，
      Span由一组 Annotation和BinaryAnnotation组成。
      整个Trace Tree通过Trace Id、Span ID、parent Span ID串起来的。
      

## 设计
1.先全局生成一个tracer,其中tracer里面有reporter,负责把span 发走
2.生成一个span,收集信息,然后调用span.Finish(),把span通过reporter发走
