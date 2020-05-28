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