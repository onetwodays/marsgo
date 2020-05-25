git clone https://github.com/openzipkin/docker-zipkin
cd docker-zipkin
docker-compose up

> 在一次Trace中，每个服务的每一次调用，就是一个基本工作单元，
>就像上图中的每一个树节点，称之为span。每一个span都有一个id作为唯一标识，
>同样每一次Trace都会生成一个traceId在span中作为追踪标识，
>另外再通过一个parentId标明本次调用的发起者（就是发起者的span-id）。
>当span有了上面三个标识后，就可以很清晰的将多个span进行梳理串联，
>最终归纳出一条完整的跟踪链路。此外，span还会有其他数据，
>比如：名称、节点上下文、时间戳以及K-V结构的tag信息等等
>（Zipkin v1核心注解如“cs”和“sr”已被Span.Kind取代，详情查看zipkin-api，
>本文会在入门的demo介绍完后对具体的Span数据模型进行说明）。

