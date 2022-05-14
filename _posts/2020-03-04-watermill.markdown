---
title: watermill
layout: post
category: golang
author: 夏泽民
---
https://github.com/ThreeDotsLabs/watermill/
Go 每日一库之 message-bus中，我们介绍了一款小巧、实现简单的异步通信库。作为学习，message-bus确实不错。但是在实际使用上，message-bus的功能就有点捉襟见肘了。例如，message-bus将消息发送到订阅者管道之后就不管了，这样如果订阅者处理压力较大，会在管道中堆积太多消息，一旦订阅者异常退出，这些消息将会全部丢失！另外，message-bus不负责保存消息，如果订阅者后启动，之前发布的消息，这个订阅者是无法收到的。这些问题，我们将要介绍的watermill都能解决！

watermill是 Go 语言的一个异步消息解决方案，它支持消息重传、保存消息，后启动的订阅者也能收到前面发布的消息。watermill内置了多种订阅-发布实现，包括Kafka/RabbitMQ，甚至还支持HTTP/MySQL binlog。当然也可以编写自己的订阅-发布实现。此外，它还提供了监控、限流等中间件。
https://segmentfault.com/a/1190000021881252
<!-- more -->
快速使用
watermill内置了很多订阅-发布实现，最简单、直接的要属GoChannel。我们就以这个实现为例介绍watermill的特性。

安装：

$ go get github.com/ThreeDotsLabs/watermill
使用：

package main

import (
  "context"
  "log"
  "time"

  "github.com/ThreeDotsLabs/watermill"
  "github.com/ThreeDotsLabs/watermill/message"
  "github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
)

func main() {
  pubSub := gochannel.NewGoChannel(
    gochannel.Config{},
    watermill.NewStdLogger(false, false),
  )

  messages, err := pubSub.Subscribe(context.Background(), "example.topic")
  if err != nil {
    panic(err)
  }

  go process(messages)

  publishMessages(pubSub)
}

func publishMessages(publisher message.Publisher) {
  for {
    msg := message.NewMessage(watermill.NewUUID(), []byte("Hello, world!"))

    if err := publisher.Publish("example.topic", msg); err != nil {
      panic(err)
    }

    time.Sleep(time.Second)
  }
}

func process(messages <-chan *message.Message) {
  for msg := range messages {
    log.Printf("received message: %s, payload: %s", msg.UUID, string(msg.Payload))
    msg.Ack()
  }
}
首先，我们创建一个GoChannel对象，它是一个消息管理器。可以调用其Subscribe订阅某个主题（topic）的消息，调用其Publish()以某个主题发布消息。Subscribe()方法会返回一个<-chan *message.Message，一旦该主题有消息发布，GoChannel就会将消息发送到该管道中。订阅者只需监听此管道，接收消息进行处理。在上面的例子中，我们启动了一个消息处理的goroutine，持续从管道中读取消息，然后打印输出。主goroutine在一个死循环中每隔 1s 发布一次消息。

message.Message这个结构是watermill库的核心，每个消息都会封装到该结构中发送。Message保存的是原始的字节流（[]byte），所以可以将 JSON/protobuf/XML 等等格式的序列化结果保存到Message中。

有两点注意：

收到的每个消息都需要调用Message的Ack() 方法确认，否则GoChannel会重发当前消息；
Message有一个UUID字段，建议设置为唯一的，方便定位问题。watermill提供方法NewUUID()生成唯一 id。
下面看示例运行：



路由
上面的发布和订阅实现是非常底层的模式。在实际应用中，我们通常想要监控、重试、统计等一些功能。而且上面的例子中，每个消息处理结束需要手动调用Ack()方法，消息管理器才会下发后面一条信息，很容易遗忘。还有些时候，我们有这样的需求，处理完某个消息后，重新发布另外一些消息。

这些功能都是比较通用的，为此watermill提供了路由（Router）功能。直接拿来官网的图：



路由其实管理多个订阅者，每个订阅者在一个独立的goroutine中运行，彼此互不干扰。订阅者收到消息后，交由注册时指定的处理函数（HandlerFunc）。路由还可以设置插件（plugin）和中间件（middleware），插件是定制路由的行为，而中间件是定制处理器的行为。处理器处理消息后会返回若干消息，这些消息会被路由重新发布到（另一个）管理器中。

var (
  logger = watermill.NewStdLogger(false, false)
)

func main() {
  router, err := message.NewRouter(message.RouterConfig{}, logger)
  if err != nil {
    panic(err)
  }

  pubSub := gochannel.NewGoChannel(gochannel.Config{}, logger)
  go publishMessages(pubSub)

  router.AddHandler("myhandler", "in_topic", pubSub, "out_topic", pubSub, myHandler{}.Handler)

  router.AddNoPublisherHandler("print_in_messages", "in_topic", pubSub, printMessages)
  router.AddNoPublisherHandler("print_out_messages", "out_topic", pubSub, printMessages)

  ctx := context.Background()
  if err := router.Run(ctx); err != nil {
    panic(err)
  }
}

func publishMessages(publisher message.Publisher) {
  for {
    msg := message.NewMessage(watermill.NewUUID(), []byte("Hello, world!"))
    if err := publisher.Publish("in_topic", msg); err != nil {
      panic(err)
    }

    time.Sleep(time.Second)
  }
}

func printMessages(msg *message.Message) error {
  fmt.Printf("\n> Received message: %s\n> %s\n>\n", msg.UUID, string(msg.Payload))
  return nil
}

type myHandler struct {
}

func (m myHandler) Handler(msg *message.Message) ([]*message.Message, error) {
  log.Println("myHandler received message", msg.UUID)

  msg = message.NewMessage(watermill.NewUUID(), []byte("message produced by myHandler"))
  return message.Messages{msg}, nil
}
首先，我们创建一个路由：

router, err := message.NewRouter(message.RouterConfig{}, logger)
然后为路由注册处理器。注册的处理器有两种类型，一种是：

router.AddHandler("myhandler", "in_topic", pubSub, "out_topic", pubSub, myHandler{}.Handler)
这个方法原型为：

func (r *Router) AddHandler(
  handlerName string,
  subscribeTopic string,
  subscriber Subscriber,
  publishTopic string,
  publisher Publisher,
  handlerFunc HandlerFunc,
) *Handler
该方法的作用是创建一个名为handlerName的处理器，监听subscriber中主题为subscribeTopic的消息，收到消息后调用handlerFunc处理，将返回的消息以主题publishTopic发布到publisher中。

另外一种处理器是下面这种形式：

router.AddNoPublisherHandler("print_in_messages", "in_topic", pubSub, printMessages)
router.AddNoPublisherHandler("print_out_messages", "out_topic", pubSub, printMessages)
从名字我们也可以看出，这种形式的处理器只处理接收到的消息，不发布新消息。

最后，我们调用router.Run()运行这个路由。

其中，创建GoChannel发布消息和上面的没什么不同。

使用路由还有个好处，处理器返回时，若无错误，路由会自动调用消息的Ack()方法；若发生错误，路由会调用消息的Nack()方法通知管理器重发这条消息。

上面只是路由的最基本用法，路由的强大之处在于中间件。

中间件
watermill中内置了几个比较常用的中间件：

IgnoreErrors：可以忽略指定的错误；
Throttle：限流，限制单位时间内处理的消息数量；
Poison：将处理失败的消息以另一个主题发布；
Retry：重试，处理失败可以重试；
Timeout：超时，如果消息处理时间超过给定的时间，直接失败。
InstantAck：直接调用消息的Ack()方法，不管后续成功还是失败；
RandomFail：随机抛出错误，测试时使用；
Duplicator：调用两次处理函数，两次返回的消息都重新发布出去，double~
Correlation：处理函数生成的消息都统一设置成原始消息中的correlation id，方便追踪消息来源；
Recoverer：捕获处理函数中的panic，包装成错误返回。
中间件的使用也是比较简单和直接的：调用router.AddMiddleware()。例如，我们想要把处理返回的消息 double 一下：

router.AddMiddleware(middleware.Duplicator)
想重试？可以：

router.AddMiddleware(middleware.Retry{
  MaxRetries:      3,
  InitialInterval: time.Millisecond * 100,
  Logger:          logger,
}.Middleware)
上面设置最大重试次数为 3，重试初始时间间隔为 100ms。

一般情况下，生产环境需要保证稳定性，某个处理异常不能影响后续的消息处理。故设置Recoverer是比较好的选择：

router.AddMiddleware(middleware.Recoverer)
也可以实现自己的中间件：

func MyMiddleware(h message.HandlerFunc) message.HandlerFunc {
  return func(message *message.Message) ([]*message.Message, error) {
    fields := watermill.LogFields{"name": m.Name}
    logger.Info("myMiddleware before", fields)
    ms, err := h(message)
    logger.Info("myMiddleware after", fields)
    return ms, err
  }
}
中间件有两种实现方式，如果不需要参数或依赖，那么直接实现为函数即可，像上面这样。如果需要有参数，那么可以实现为一个结构：

type myMiddleware struct {
  Name string
}

func (m myMiddleware) Middleware(h message.HandlerFunc) message.HandlerFunc {
  return func(message *message.Message) ([]*message.Message, error) {
    fields := watermill.LogFields{"name": m.Name}
    logger.Info("myMiddleware before", fields)
    ms, err := h(message)
    logger.Info("myMiddleware after", fields)
    return ms, err
  }
}
这两种中间件的添加方式有所不同，第一种直接添加：

router.AddMiddleware(MyMiddleware)
第二种要构造一个对象，然后将其Middleware方法传入，在该方法中可以访问MyMiddleware对象的字段：

router.AddMiddleware(MyMiddleware{Name:"dj"}.Middleware)
设置
如果运行上面程序，你很可能会看到这样一条日志：

No subscribers to send message
因为发布消息是在另一个goroutine，我们没有控制何时发布，可能发布消息时，我们还未订阅。我们观察后面的处理日志，对比 uuid 发现这条消息直接被丢弃了。watermill提供了一个选项，可以将消息都保存下来，订阅某个主题时将该主题之前的消息也发送给它：

pubSub := gochannel.NewGoChannel(
  gochannel.Config{
    Persistent: true,
  }, logger)
创建GoChannel时将Config中Persistent字段设置为true即可。此时运行，我们仔细观察一下，出现No subscribers to send message信息的消息后续确实被处理了。

RabbitMQ
除了GoChannel，watermill还内置了其他的发布-订阅实现。这些实现除了发布-订阅器创建的方式不同，其他与我们之前介绍的基本一样。这里我们简单介绍一下RabbitMQ，其他的可自行研究。

使用RabbitMQ需要先运行RabbitMQ程序，RabbitMQ采用Erlang开发。我们之前很多文章也介绍过 windows 上的软件安装神器choco。使用choco安装RabbitMQ：

$ choco install rabbitmq
启动RabbitMQ服务器：

$ rabbitmq-server.bat
watermill对RabbitMQ的支持使用独立库的形式，需要另行安装：

$ go get -u github.com/ThreeDotsLabs/watermill-amqp/pkg/amqp
发布订阅：

var amqpURI = "amqp://localhost:5672/"

func main() {
  amqpConfig := amqp.NewDurableQueueConfig(amqpURI)

  subscriber, err := amqp.NewSubscriber(
    amqpConfig,
    watermill.NewStdLogger(false, false),
  )
  if err != nil {
    panic(err)
  }

  messages, err := subscriber.Subscribe(context.Background(), "example.topic")
  if err != nil {
    panic(err)
  }

  go process(messages)

  publisher, err := amqp.NewPublisher(amqpConfig, watermill.NewStdLogger(false, false))
  if err != nil {
    panic(err)
  }

  publishMessages(publisher)
}

func publishMessages(publisher message.Publisher) {
  for {
    msg := message.NewMessage(watermill.NewUUID(), []byte("Hello, world!"))

    if err := publisher.Publish("example.topic", msg); err != nil {
      panic(err)
    }

    time.Sleep(time.Second)
  }
}

func process(messages <-chan *message.Message) {
  for msg := range messages {
    log.Printf("received message: %s, payload: %s", msg.UUID, string(msg.Payload))
    msg.Ack()
  }
}
如果有自定义发布-订阅实现的需求，可以参考RabbitMQ的实现：github.com/ThreeDotsLabs/watermill-amqp/pkg/amqp。

总结
watermill提供丰富的功能，且预留了扩展点，可自行扩展。另外，源码中处理goroutine创建和通信、多种并发模式的应用都是值得一看的。官方 GitHub 上还有一个事件驱动示例：https://github.com/ThreeDotsLabs/event-driven-example。

大家如果发现好玩、好用的 Go 语言库，欢迎到 Go 每日一库 GitHub 上提交 issue😄

参考
watermill 官方文档：https://watermill.io/
Go 每日一库 GitHub：https://github.com/darjun/go-daily-lib
