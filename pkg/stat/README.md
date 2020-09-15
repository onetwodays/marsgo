# stat

## 项目简介

数据统计、监控采集等

this is the Go client library for Prometheus. 
It has two separate parts, one for instrumenting application code, 
and one for creating clients that talk to the Prometheus HTTP API.

任意组件只要提供对应的HTTP接口并符合Prometheus定义的数据格式，就可以介入Prometheus监控
Prometheus Server负载定时在目标上抓取metrics(指标)数据，
每个抓取目标都需要暴露一个HTTP服务接口用于Prometheus定时抓取。
这种调用被监控对象获取监控数据的方式被称为Pull(拉)。
Pull方式体现了Prometheus独特的设计哲学与大多数采用Push(推)方式的监控不同

>
>度量指标和标签
 每个时间序列（Time Serie，简称时序）由度量指标和一组标签键值对唯一确定。
>标签值可以包含任意 Unicode 字符，包括中文。
>采样值（Sample）
 时序数据其实就是一系列采样值。每个采样值包括2部分：
 
 1. 一个 64 位的浮点数值
 2. 一个精确到毫秒的时间戳
 
 >注解（Notation）
  一个注解由一个度量指标和一组标签键值对构成 api_http_requests_total{method="POST", handler="/messages"}
  度量指标为 api_http_requests_total，标签为 method="POST"、handler="/messages"
>
>Prometheus 的 histogram 是一种累积直方图，与上面的区间划分方式是有差别的，
>它的划分方式如下：还假设每个 bucket 的宽度是 0.2s，
>那么第一个 bucket 表示响应时间小于等于 0.2s 的请求数量，
>第二个 bucket 表示响应时间小于等于 0.4s 的请求数量，以此类推。
>也就是说，每一个 bucket 的样本包含了之前所有 bucket 的样本，所以叫累积直方图