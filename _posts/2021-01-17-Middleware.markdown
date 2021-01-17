---
title: Middleware
layout: post
category: golang
author: 夏泽民
---
微服务架构的一些常用组件
1. 服务治理。通常是采用注册发现的机制。有一个注册中心
2. 集中式配置
3. 反向代理
4. ADN， CDN
5. 分布式存储
6. 分布式日志
7. 分布式锁
8. 消息队列
9. 分布式文件存储
10. 断路器
11. 分布式数据库
12.路由与web服务器
13.RPC框架
14.缓存服务
15.分布式事务
16.任务调度
17.安全
https://thinkjs.org/zh-cn/doc/3.0/middleware.html
https://www.zhihu.com/question/19730582

AOP已经实现了业务隔离。但却带来了一串长长的链式调用，如果处理不当很容易掉链子。另外，这种结构实现异步操作较为麻烦。

多个中间件会形成一个栈结构（middle stack），以"先进后出"（first-in-last-out）的顺序执行，被称为洋葱结构。
<!-- more -->
Koa的中间件使用。
const logger = (ctx, next) => {
  console.log(`${Date.now()} ${ctx.request.method} ${ctx.request.url}`);
  next();
}
app.use(logger);

中间件一个奇妙的点在于next函数。如果中间件内部没有调用next函数，那么执行权就不会传递下去

 Egg 是基于 Koa 实现的，所以 Egg 的中间件形式和 Koa 的中间件形式是一样的，都是基于洋葱圈模型。
 https://eggjs.org/zh-cn/basics/middleware.html
 
 http://c.biancheng.net/view/3860.html
 
 Middleware module	Description	Replaces built-in function (Express 3)
body-parser	Parse HTTP request body. See also: body, co-body, and raw-body.	express.bodyParser
compression	Compress HTTP responses.	express.compress
connect-rid	Generate unique request ID.	NA
cookie-parser	Parse cookie header and populate req.cookies. See also cookies and keygrip.	express.cookieParser
cookie-session	Establish cookie-based sessions.	express.cookieSession
cors	Enable cross-origin resource sharing (CORS) with various options.	NA
csurf	Protect from CSRF exploits.	express.csrf
errorhandler	Development error-handling/debugging.	express.errorHandler
method-override	Override HTTP methods using header.	express.methodOverride
morgan	HTTP request logger.	express.logger
multer	Handle multi-part form data.	express.bodyParser
response-time	Record HTTP response time.	express.responseTime
serve-favicon	Serve a favicon.	express.favicon
serve-index	Serve directory listing for a given path.	express.directory
serve-static	Serve static files.	express.static
session	Establish server-based sessions (development only).	express.session
timeout	Set a timeout period for HTTP request processing.	express.timeout
vhost	Create virtual domains.	express.vhost
 http://expressjs.com/en/resources/middleware.html
 
 Middleware module	Description
cls-rtracer	Middleware for CLS-based request id generation. An out-of-the-box solution for adding request ids into your logs.
connect-image-optimus	Optimize image serving. Switches images to .webp or .jxr, if possible.
express-debug	Development tool that adds information about template variables (locals), current session, and so on.
express-partial-response	Filters out parts of JSON responses based on the fields query-string; by using Google API’s Partial Response.
express-simple-cdn	Use a CDN for static assets, with multiple host support.
express-slash	Handles routes with and without trailing slashes.
express-stormpath	User storage, authentication, authorization, SSO, and data security.
express-uncapitalize	Redirects HTTP requests containing uppercase to a canonical lowercase form.
helmet	Helps secure your apps by setting various HTTP headers.
join-io	Joins files on the fly to reduce the requests count.
passport	Authentication using “strategies” such as OAuth, OpenID and many others. See http://passportjs.org/ for more information.
static-expiry	Fingerprint URLs or caching headers for static assets.
view-helpers	Common helper methods for views.
sriracha-admin	Dynamically generate an admin site for Mongoose.

纵观GO语言，中间件应用比较普遍，主要应用：

记录对服务器发送的请求（request）
处理服务器响应（response ）
请求和处理之间做一个权限认证工作
远程调用
安全

1,单个中间件
type Middleware func(http.HandlerFunc) http.HandlerFunc

2，多个中间件，递归方式使用，洋葱结构，koa、express、echo实现方式
func Chain(f http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
   for _, m := range middlewares {
      f = m(f)
   }
   return f
}

3，多个中间件，全局列表结构，gin框架实现
// HandlersChain defines a HandlerFunc array.
type HandlersChain []HandlerFunc
//模拟的调用堆栈
func (c *Context) Next() {
	c.index++
	for c.index < int8(len(c.handlers)) {
		//按顺序执行HandlersChain内的函数
		//如果函数内无c.Next()方法调用则函数顺序执行完
		//如果函数内有c.Next()方法调用则代码执行到c.Next()方法处压栈，等待后面的函数执行完在回来执行c.Next()后的命令
		c.handlers[c.index](c)
		c.index++
	}
}

https://www.cnblogs.com/maji233/p/11237349.html
https://blog.csdn.net/u014270740/article/details/91401698
https://blog.csdn.net/weixin_30339457/article/details/97592881


1,log
2,header
3,compress
4,auth
5,crsf
https://github.com/clevergo/middleware

服务描述
服务调用首先要解决的问题就是服务如何对外描述。比如服务名、调用这个服务需要提供哪些信息、返回的结果是什么格式的、如何解析等问题。

常用的服务描述方式包括 RESTful API、XML 配置以及 IDL 文件三种。

其中，RESTful API 方式通常用于 HTTP 协议的服务描述，并且常用 Wiki 或者Swagger来进行管理。下面是一个 RESTful API 方式的服务描述的例子。

XML 配置方式多用作 RPC 协议的服务描述，通过 *.xml 配置文件来定义接口名、参数以及返回值类型等。

如motan_server.xml中

<motan:basicService id="serviceBasicConfig" export="demoMotan:8002" group="motan-demo-rpc" module="motan-demo-rpc" registry="registry"/>
IDL 文件方式通常用作 Thrift 和 gRPC 这类跨语言服务调用框架中，比如 gRPC 就是通过 Protobuf 文件来定义服务的接口名、参数以及返回值的数据结构。

服务描述方式                           使用场景                                                缺点

RESTFUL API                          跨语言平台，组织内外皆可                   使用了HTP作为通信协议，相比TCP协议，性能较差

XML配置                                  Java平台，一般用作组织内部                不支持跨语言平台

IDL文件                                     跨语言平台，组织内外皆可                   修改或者除PB字段不能向前兼容

 

注册中心
有了服务的接口描述，下一步要解决的问题就是服务的发布和订阅，就是说你提供了一个服务，如何让外部想调用你的服务的人知道。这个时候就需要一个类似注册中心的角色，服务提供者将自己提供的服务以及地址登记到注册中心，服务消费者则从注册中心查询所需要调用的服务的地址，然后发起请求。

一般来讲，注册中心的工作流程是：

服务提供者在启动时，根据服务发布文件中配置的发布信息向注册中心注册自己的服务。

服务消费者在启动时，根据消费者配置文件中配置的服务信息向注册中心订阅自己所需要的服务。

注册中心返回服务提供者地址列表给服务消费者。

当服务提供者发生变化，比如有节点新增或者销毁，注册中心将变更通知给服务消费者。

 

服务框架
通过注册中心，服务消费者就可以获取到服务提供者的地址，有了地址后就可以发起调用。但在发起调用之前你还需要解决以下几个问题。

服务通信采用什么协议？就是说服务提供者和服务消费者之间以什么样的协议进行网络通信，是采用四层 TCP、UDP 协议，还是采用七层 HTTP 协议，还是采用其他协议？

数据传输采用什么方式？就是说服务提供者和服务消费者之间的数据传输采用哪种方式，是同步还是异步，是在单连接上传输，还是多路复用。

数据压缩采用什么格式？通常数据传输都会对数据进行压缩，来减少网络传输的数据量，从而减少带宽消耗和网络传输时间，比如常见的 JSON 序列化、Java 对象序列化以及 Protobuf 序列化等。

 

服务监控
一旦服务消费者与服务提供者之间能够正常发起服务调用，你就需要对调用情况进行监控，以了解服务是否正常。通常来讲，服务监控主要包括三个流程。

指标收集。就是要把每一次服务调用的请求耗时以及成功与否收集起来，并上传到集中的数据处理中心。

数据处理。有了每次调用的请求耗时以及成功与否等信息，就可以计算每秒服务请求量、平均耗时以及成功率等指标。

数据展示。数据收集起来，经过处理之后，还需要以友好的方式对外展示，才能发挥价值。通常都是将数据展示在 Dashboard 面板上，并且每隔 10s 等间隔自动刷新，用作业务监控和报警等。

 

服务追踪
除了需要对服务调用情况进行监控之外，你还需要记录服务调用经过的每一层链路，以便进行问题追踪和故障定位。

服务追踪的工作原理大致如下：

服务消费者发起调用前，会在本地按照一定的规则生成一个 requestid，发起调用时，将 requestid 当作请求参数的一部分，传递给服务提供者。

服务提供者接收到请求后，记录下这次请求的 requestid，然后处理请求。如果服务提供者继续请求其他服务，会在本地再生成一个自己的 requestid，然后把这两个 requestid 都当作请求参数继续往下传递。

以此类推，通过这种层层往下传递的方式，一次请求，无论最后依赖多少次服务调用、经过多少服务节点，都可以通过最开始生成的 requestid 串联所有节点，从而达到服务追踪的目的。

 

服务治理
服务监控能够发现问题，服务追踪能够定位问题所在，而解决问题就得靠服务治理了。服务治理就是通过一系列的手段来保证在各种意外情况下，服务调用仍然能够正常进行。

常见故障：

单机故障：自动摘除故障节点。

单 IDC 故障：自动切换故障 IDC 的流量到其他正常 

IDC依赖服务不可用：熔断。

https://blog.csdn.net/haponchang/article/details/90746408



