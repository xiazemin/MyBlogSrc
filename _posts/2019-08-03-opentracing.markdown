---
title: opentracing
layout: post
category: golang
author: 夏泽民
---
在Go 1.7，我们有一个新包/ HTTP / httptrace提供了一个方便的机制，观察一个HTTP请求时会发生什么
分布式跟踪是监测和分析微服务架构系统，导出结果到为X-TRACE，如谷歌的Dapper和Twitter的Zipkin 。 它们的底层原理是分布式环境传播 ，其中涉及的某些元数据与进入系统的每个请求相关联，并且跨线程和进程边界传播元数据跟随请求进出各种微服务调用。 如果我们为每个入站请求分配一个唯一的ID并将其作为分布式上下文的一部分，那么我们可以将来自多个线程和多个进程的各种性能分析数据合并到统一的表示我们系统执行请求的“跟踪”中。
分布式跟踪需要使用Hook钩子和上下文传播机制来测试应用程序代码（或其使用的框架）。
没有良好的API为开发人员提供在编程语言之间内部一致性，那就无法绑定到指定的跟踪系统。
2015年10月一个新的社区形成，催生了OpenTracing API，一个开放的，厂商中立的，与语言无关的分布式跟踪标准。你可以阅读更多关于Ben Sigelman有关OpenTracing动机和设计原理背后的文章 。
https://opentracing.io/
https://zipkin.io/
https://github.com/openzipkin/zipkin
https://github.com/openzipkin/zipkin-go
https://github.com/opentracing/opentracing-go
<!-- more -->

import (
   "net/http"
   "net/http/httptrace"

   "github.com/opentracing/opentracing-go"
   "github.com/opentracing/opentracing-go/log"
   "golang.org/x/net/context"
)

// We will talk about this later
var tracer opentracing.Tracer

func AskGoogle(ctx context.Context) error {
   // retrieve current Span from Context
   var parentCtx opentracing.SpanContext
   parentSpan := opentracing.SpanFromContext(ctx); 
   if parentSpan != nil {
      parentCtx = parentSpan.Context()
   }

   // start a new Span to wrap HTTP request
   span := tracer.StartSpan(
      "ask google",
      opentracing.ChildOf(parentCtx),
   )

   // make sure the Span is finished once we're done
   defer span.Finish()

   // make the Span current in the context
   ctx = opentracing.ContextWithSpan(ctx, span)

   // now prepare the request
   req, err := http.NewRequest("GET", "http://google.com", nil)
   if err != nil {
      return err
   }

   // attach ClientTrace to the Context, and Context to request 
   trace := NewClientTrace(span)
   ctx = httptrace.WithClientTrace(ctx, trace)
   req = req.WithContext(ctx)

   // execute the request
   res, err := http.DefaultClient.Do(req)
   if err != nil {
      return err
   }
   
   // Google home page is not too exciting, so ignore the result
   res.Body.Close()
   return nil
}

注意：

1.回避了tracer变量初始化的问题。

2.AskGoogle函数接受context.Context对象。这是为开发分布式应用程序的Go推荐方式，因为上下文对象是要让分布式环境传播。

3.我们假设上下文已经包含父跟踪Span。 OpenTracing API中的Span是用于表示由微服务执行的工作单元。HTTP调用就是可以包裹在跟踪Span中的操作案例。 当我们运行处理入站请求的服务时，服务通常会为每个请求创建一个跟踪范围，并将其存储在上下文中，以便在对另一个服务进行下游调用时可用（在我们的示例中为google.com ）。

4.我们启动一个新的子Span来包装出站HTTP调用。 如果父Span缺失，这是好方法。

5.最后，在做出HTTP请求之前，我们实例化一个ClientTrace并将其附加到请求。

ClientTrace结构是httptrace的基本构建块 。它允许我们在HTTP请求的生命周期内注册将由HTTP客户端执行的回调函数。 例如，ClientTrace结构有这样的方法：

type ClientTrace struct {
...
        // DNSStart is called when a DNS lookup begins.
        DNSStart func(DNSStartInfo)
        // DNSDone is called when a DNS lookup ends.
        DNSDone func(DNSDoneInfo)
...
}
<p>
我们在NewClientTrace方法中创建这个结构的一个实例：

func NewClientTrace(span opentracing.Span) *httptrace.ClientTrace {
   trace := &clientTrace{span: span}
   return &httptrace.ClientTrace {
      DNSStart: trace.dnsStart,
      DNSDone:  trace.dnsDone,
   }
}

// clientTrace holds a reference to the Span and
// provides methods used as ClientTrace callbacks
type clientTrace struct {
   span opentracing.Span
}

func (h *clientTrace) dnsStart(info httptrace.DNSStartInfo) {
   h.span.LogKV(
      log.String("event", "DNS start"),
      log.Object("host", info.Host),
   )
}

func (h *clientTrace) dnsDone(httptrace.DNSDoneInfo) {
   h.span.LogKV(log.String("event", "DNS done"))
}

<p>
我们为DBBStart和DNSDone事件实现注册两个回调函数，通过私有结构clientTrace保有一个指向跟踪Span。 在回调方法中，我们使用Span的键值记录API来记录事件的信息，以及Span本身隐含捕获的时间戳。

你不是说关于UI的东西吗？

OpenTracing API的工作方式是，一旦调用跟踪Span的Finish（）方法，Span捕获的数据将发送到跟踪系统后端，通常在后台异步发送。然后，我们可以使用跟踪系统UI来查找跟踪并在时间线上将其可视化。 然而，上述例子只是为了说明使用OpenTracing与httptrace的原理。对于真正的工作示例，我们将使用现有的库https://github.com/opentracing-contrib/go-stdlib 。 使用这个库我们的客户端代码不需要担心跟踪实际的HTTP调用。但是，我们仍然希望创建一个顶层跟踪Span来表示客户端应用程序的整体执行，并记录任何错误。

package main

import (
   "fmt"
   "io/ioutil"
   "log"
   "net/http"

   "github.com/opentracing-contrib/go-stdlib/nethttp"
   "github.com/opentracing/opentracing-go"
   "github.com/opentracing/opentracing-go/ext"
   otlog "github.com/opentracing/opentracing-go/log"
   "golang.org/x/net/context"
)

func runClient(tracer opentracing.Tracer) {
   // nethttp.Transport from go-stdlib will do the tracing
   c := &http.Client{Transport: &nethttp.Transport{}}

   // create a top-level span to represent full work of the client
   span := tracer.StartSpan(client)
   span.SetTag(string(ext.Component), client)
   defer span.Finish()
   ctx := opentracing.ContextWithSpan(context.Background(), span)

   req, err := http.NewRequest(
      "GET",
      fmt.Sprintf("http://localhost:%s/", *serverPort),
      nil,
   )
   if err != nil {
      onError(span, err)
      return
   }

   req = req.WithContext(ctx)
   // wrap the request in nethttp.TraceRequest
   req, ht := nethttp.TraceRequest(tracer, req)
   defer ht.Finish()

   res, err := c.Do(req)
   if err != nil {
      onError(span, err)
      return
   }
   defer res.Body.Close()
   body, err := ioutil.ReadAll(res.Body)
   if err != nil {
      onError(span, err)
      return
   }
   fmt.Printf("Received result: %s\n", string(body))
}

func onError(span opentracing.Span, err error) {
   // handle errors by recording them in the span
   span.SetTag(string(ext.Error), true)
   span.LogKV(otlog.Error(err))
   log.Print(err)
}
<p>
上面的客户端代码调用本地服务器。 让我们实现它。

package main

import (
   "fmt"
   "io"
   "log"
   "net/http"
   "time"

   "github.com/opentracing-contrib/go-stdlib/nethttp"
   "github.com/opentracing/opentracing-go"
)

func getTime(w http.ResponseWriter, r *http.Request) {
   log.Print("Received getTime request")
   t := time.Now()
   ts := t.Format("Mon Jan _2 15:04:05 2006")
   io.WriteString(w, fmt.Sprintf("The time is %s", ts))
}

func redirect(w http.ResponseWriter, r *http.Request) {
   http.Redirect(w, r,
      fmt.Sprintf("http://localhost:%s/gettime", *serverPort), 301)
}

func runServer(tracer opentracing.Tracer) {
   http.HandleFunc("/gettime", getTime)
   http.HandleFunc("/", redirect)
   log.Printf("Starting server on port %s", *serverPort)
   http.ListenAndServe(
      fmt.Sprintf(":%s", *serverPort),
      // use nethttp.Middleware to enable OpenTracing for server
      nethttp.Middleware(tracer, http.DefaultServeMux))
}
 
注意，客户端向根端点“/”发出请求，但服务器将其重定向到“/ gettime”端点。 这样做允许我们更好地说明如何在跟踪系统中捕获跟踪。

开发和工程团队因为系统组件水平扩展、开发团队小型化、敏捷开发、CD（持续集成）、解耦等各种需求，正在使用现代的微服务架构替换老旧的单片机系统。也就是说，当一个生产系统面对真正的高并发，或者解耦成大量微服务时，以前很容易实现的重点任务变得困难了。过程中需要面临一系列问题：用户体验优化、后台真是错误原因分析，分布式系统内各组件的调用情况等。当代分布式跟踪系统（例如，Zipkin, Dapper, HTrace, X-Trace等）旨在解决这些问题，但是他们使用不兼容的API来实现各自的应用需求。尽管这些分布式追踪系统有着相似的API语法，但各种语言的开发人员依然很难将他们各自的系统（使用不同的语言和技术）和特定的分布式追踪系统进行整合，

为什么需要OpenTracing？
OpenTracing通过提供平台无关、厂商无关的API，使得开发人员能够方便的添加（或更换）追踪系统的实现。 OpenTracing提供了用于运营支撑系统的和针对特定平台的辅助程序库。程序库的具体信息请参考详细的规范。

什么是一个Trace?
在广义上，一个trace代表了一个事务或者流程在（分布式）系统中的执行过程。在OpenTracing标准中，trace是多个span组成的一个有向无环图（DAG），每一个span代表trace中被命名并计时的连续性的执行片段。


分布式追踪中的每个组件都包含自己的一个或者多个span。例如，在一个常规的RPC调用过程中，OpenTracing推荐在RPC的客户端和服务端，至少各有一个span，用于记录RPC调用的客户端和服务端信息。

一个父级的span会显示的并行或者串行启动多个子span。在OpenTracing标准中，甚至允许一个子span有个多父span（例如：并行写入的缓存，可能通过一次刷新操作写入动作）。

这种展现方式增加显示了执行时间的上下文，相关服务间的层次关系，进程或者任务的串行或并行调用关系。这样的视图有助于发现系统调用的关键路径。通过关注关键路径的执行过程，项目团队可能专注于优化路径中的关键位置，最大幅度的提升系统性能。例如：可以通过追踪一个资源定位的调用情况，明确底层的调用情况，发现哪些操作有阻塞的情况。

http://bigbully.github.io/Dapper-translation/

https://github.com/bigbully/Dapper-translation
