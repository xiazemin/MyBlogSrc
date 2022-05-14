---
title: jaegertracing
layout: post
category: golang
author: 夏泽民
---
docker pull jaegertracing/all-in-one

docker run -d --name jaeger
-e COLLECTOR_ZIPKIN_HTTP_PORT=9411
-p 5775:5775/udp
-p 6831:6831/udp
-p 6832:6832/udp
-p 5778:5778
-p 16686:16686
-p 14268:14268
-p 9411:9411
jaegertracing/all-in-one:1.14
<!-- more -->
https://registry.hub.docker.com/r/jaegertracing/all-in-one
Opentracing是分布式链路追踪的一种规范标准，是CNCF（云原生计算基金会）下的项目之一。和一般的规范标准不同，Opentracing不是传输协议，消息格式层面上的规范标准，而是一种语言层面上的API标准。以Go语言为例，只要某链路追踪系统实现了Opentracing规定的接口（interface），符合Opentracing定义的表现行为，那么就可以说该应用符合Opentracing标准。这意味着开发者只需修改少量的配置代码，就可以在符合Opentracing标准的链路追踪系统之间自由切换。

所定义的数据模型。
Span
Span是一条追踪链路中的基本组成要素，一个span表示一个独立的工作单元，比如可以表示一次函数调用，一次http请求等等。span会记录如下基本要素:

服务名称(operation name)
服务的开始时间和结束时间
K/V形式的Tags
K/V形式的Logs
SpanContext
References：该span对一个或多个span的引用（通过引用SpanContext）。

Tags
Tags以K/V键值对的形式保存用户自定义标签，主要用于链路追踪结果的查询过滤。例如： http.method="GET",http.status_code=200。其中key值必须为字符串，value必须是字符串，布尔型或者数值型。
span中的tag仅自己可见，不会随着 SpanContext传递给后续span。
例如：
span.SetTag("http.method","GET")
span.SetTag("http.status_code",200)
复制代码Logs
Logs与tags类似，也是K/V键值对形式。与tags不同的是，logs还会记录写入logs的时间，因此logs主要用于记录某些事件发生的时间。logs的key值同样必须为字符串，但对value类型则没有限制。例如：
span.LogFields(
	log.String("event", "soft error"),
	log.String("type", "cache timeout"),
	log.Int("waited.millis", 1500),

SpanContext
SpanContext携带着一些用于跨服务通信的（跨进程）数据，主要包含：

足够在系统中标识该span的信息，比如：span_id,trace_id。
Baggage Items，为整条追踪连保存跨服务（跨进程）的K/V格式的用户自定义数据。

Baggage Items
Baggage Items与tags类似，也是K/V键值对。与tags不同的是：

其key跟value都只能是字符串格式
Baggage items不仅当前span可见，其会随着SpanContext传递给后续所有的子span。要小心谨慎的使用baggage items——因为在所有的span中传递这些K,V会带来不小的网络和CPU开销。

References
Opentracing定义了两种引用关系:ChildOf和FollowFrom。
ChildOf: 父span的执行依赖子span的执行结果时，此时子span对父span的引用关系是ChildOf。比如对于一次RPC调用，服务端的span（子span）与客户端调用的span（父span）是ChildOf关系。
FollowFrom：父span的执不依赖子span执行结果时，此时子span对父span的引用关系是FollowFrom。FollowFrom常用于异步调用的表示，例如消息队列中consumerspan与producerspan之间的关系。
Trace
Trace表示一次完整的追踪链路，trace由一个或多个span组成。下图示例表示了一个由8个span组成的trace:
        [Span A]  ←←←(the root span)
            |
     +------+------+
     |             |
 [Span B]      [Span C] ←←←(Span C is a `ChildOf` Span A)
     |             |
 [Span D]      +---+-------+
               |           |
           [Span E]    [Span F] >>> [Span G] >>> [Span H]
                                       ↑
                                       ↑
                                       ↑
                         (Span G `FollowsFrom` Span F)
复制代码时间轴的展现方式会更容易理解：
––|–––––––|–––––––|–––––––|–––––––|–––––––|–––––––|–––––––|–> time

 [Span A···················································]
   [Span B··············································]
      [Span D··········································]
    [Span C········································]
         [Span E·······]        [Span F··] [Span G··] [Span H··]

Jaeger
Jaeger\ˈyā-gər\ 是Uber开源的分布式追踪系统，是遵循Opentracing的系统之一，也是CNCF项目。本篇将使用Jaeger来演示如何在系统中引入分布式追踪。

下载客户端library:
go get github.com/jaegertracing/jaeger-client-go
复制代码初始化Jaeger tracer:
import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
)
// initJaeger 将jaeger tracer设置为全局tracer
func initJaeger(service string) io.Closer {
	cfg := jaegercfg.Configuration{
		// 将采样频率设置为1，每一个span都记录，方便查看测试结果
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans: true,
			// 将span发往jaeger-collector的服务地址
			CollectorEndpoint: "http://localhost:14268/api/traces",
		},
	}
	closer, err := cfg.InitGlobalTracer(service, jaegercfg.Logger(jaeger.StdLogger))
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}
	return closer
}

复制代码创建tracer，生成root span：
func main() {
	closer := initJaeger("in-process")
	defer closer.Close()
	// 获取jaeger tracer
	t := opentracing.GlobalTracer()
	// 创建root span
	sp := t.StartSpan("in-process-service")
	// main执行完结束这个span
	defer sp.Finish()
	// 将span传递给Foo
	ctx := opentracing.ContextWithSpan(context.Background(), sp)
	Foo(ctx)
}
复制代码上述代码创建了一个root span，并将该span通过context传递给Foo方法，以便在Foo方法中将追踪链继续延续下去：
func Foo(ctx context.Context) {
    // 开始一个span, 设置span的operation_name=Foo
	span, ctx := opentracing.StartSpanFromContext(ctx, "Foo")
	defer span.Finish()
	// 将context传递给Bar
	Bar(ctx)
	// 模拟执行耗时
	time.Sleep(1 * time.Second)
}
func Bar(ctx context.Context) {
    // 开始一个span，设置span的operation_name=Bar
	span, ctx := opentracing.StartSpanFromContext(ctx, "Bar")
	defer span.Finish()
	// 模拟执行耗时
	time.Sleep(2 * time.Second)

	// 假设Bar发生了某些错误
	err := errors.New("something wrong")
	span.LogFields(
		log.String("event", "error"),
		log.String("message", err.Error()),
	)
	span.SetTag("error", true)
}

复制代码Foo方法调用了Bar，假设在Bar中发生了一些错误，可以通过span.LogFields和span.SetTag将错误记录在追踪链中。
通过上面的例子可以发现，如果要确保追踪链在程序中不断开，需要将函数的第一个参数设置为context.Context，通过opentracing.ContextWithSpan将保存到context中，通过opentracing.StartSpanFromContext开始一个新的子span。
效果查看
执行完上面的程序后，打开Jaeger UI: http://localhost:16686/search，可以看到链路追踪的结果

https://juejin.cn/post/6844903942309019661
https://github.com/opentracing/opentracing-go
https://github.com/opentracing/specification/blob/master/semantic_conventions.md
https://www.jaegertracing.io/docs/1.33/getting-started/
https://github.com/yurishkuro/opentracing-tutorial/tree/master/go
https://github.com/opentracing/specification/blob/master/specification.md#the-opentracing-data-model
https://github.com/yurishkuro/opentracing-tutorial/tree/master/go
