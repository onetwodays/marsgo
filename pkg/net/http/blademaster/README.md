#### net/http/blademaster

##### 项目简介

http 框架，带来如飞一般的体验。
我们自己在使用内置的net/http的默认路径处理HTTP请求的时候，会发现很多不足，比如
不能单独的对请求方法(POST,GET等)注册特定的处理函数
不支持Path变量参数
不能自动对Path进行校准
性能一般
扩展性不足
