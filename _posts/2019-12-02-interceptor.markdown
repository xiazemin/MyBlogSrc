---
title: interceptor
layout: post
category: golang
author: 夏泽民
---
gRPC-Go 增加了拦截器(interceptor)的功能， 就像Java Servlet中的 filter一样，可以对RPC的请求和响应进行拦截处理，而且既可以在客户端进行拦截，也可以对服务器端进行拦截。

利用拦截器，可以对gRPC进行扩展，利用社区的力量将gRPC发展壮大，也可以让开发者更灵活地处理gRPC流程中的业务逻辑。下面列出了利用拦截器实现的一些功能框架：

Go gRPC Middleware:提供了拦截器的interceptor链式的功能，可以将多个拦截器组合成一个拦截器链，当然它还提供了其它的功能，所以以gRPC中间件命名。
grpc-multi-interceptor: 是另一个interceptor链式功能的库，也可以将单向的或者流式的拦截器组合。
grpc_auth: 身份验证拦截器
grpc_ctxtags: 为上下文增加Tag map对象
grpc_zap: 支持zap日志框架
grpc_logrus: 支持logrus日志框架
grpc_prometheus: 支持 prometheus
otgrpc: 支持opentracing/zipkin
grpc_opentracing:支持opentracing/zipkin
grpc_retry: 为客户端增加重试的功能
grpc_validator: 为服务器端增加校验的功能
xrequestid: 将request id 设置到context中
go-grpc-interceptor: 解析Accept-Language并设置到context
requestdump: 输出request/response
也有其它一些文章介绍的利用拦截器的例子
https://texlution.com/post/oauth-and-grpc-go/
https://segmentfault.com/a/1190000007997759
<!-- more -->
服务器只能配置一个 unary interceptor和 stream interceptor，否则会报错，客户端也是，虽然不会报错，但是只有最后一个才起作用。 如果你想配置多个，可以使用前面提到的拦截器链或者自己实现一个。

实现拦截器麻烦吗？一点都不麻烦，相反，非常的简单。

对于服务器端的单向调用的拦截器，只需定义一个UnaryServerInterceptor方法:

type UnaryServerInterceptor func(ctx context.Context, req interface{}, info *UnaryServerInfo, handler UnaryHandler) (resp interface{}, err error)
对于服务器端stream调用的拦截器，只需定义一个StreamServerInterceptor方法:

type StreamServerInterceptor func(srv interface{}, ss ServerStream, info *StreamServerInfo, handler StreamHandler) error
方法的参数中包含了上下文，请求和stream以及要调用对象的信息。

对于客户端的单向的拦截，只需定义一个``方法：

type UnaryClientInterceptor func(ctx context.Context, method string, req, reply interface{}, cc *ClientConn, invoker UnaryInvoker, opts ...CallOption) error
对于客户端的stream的拦截，只需定义一个``方法：

type StreamClientInterceptor func(ctx context.Context, desc *StreamDesc, cc *ClientConn, method string, streamer Streamer, opts ...CallOption) (ClientStream, error)
你可以查看上面提到的一些开源的拦截器的实现，它们的实现都不是太复杂，下面我们以一个简单的例子来距离，在方法调用的前后打印一个log。

Server端的拦截器
package main
import (
	"log"
	"net"
	"flag"
	pb "github.com/smallnest/grpc/a/pb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)
var (
	port = flag.String("p", ":8972", "port")
)
type server struct{}
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}
func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", *port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer(grpc.StreamInterceptor(StreamServerInterceptor),
		grpc.UnaryInterceptor(UnaryServerInterceptor))
	pb.RegisterGreeterServer(s, &server{})
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
func UnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Printf("before handling. Info: %+v", info)
	resp, err := handler(ctx, req)
	log.Printf("after handling. resp: %+v", resp)
	return resp, err
}
// StreamServerInterceptor is a gRPC server-side interceptor that provides Prometheus monitoring for Streaming RPCs.
func StreamServerInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	log.Printf("before handling. Info: %+v", info)
	err := handler(srv, ss)
	log.Printf("after handling. err: %v", err)
	return err
}
grpc.NewServer可以将拦截器作为参数传入，在提供服务的时候，我们可以看到拦截器打印出log:
2017/04/17 23:34:20 before handling. Info: &{Server:0x17309c8 FullMethod:/pb.Greeter/SayHello}
2017/04/17 23:34:20 after handling. resp: &HelloReply{Message:Hello world,}
客户端的拦截器
package main
import (
	//"context"
	"flag"
	"log"
	"golang.org/x/net/context"
	pb "github.com/smallnest/grpc/a/pb"
	"google.golang.org/grpc"
)
var (
	address = flag.String("addr", "localhost:8972", "address")
	name    = flag.String("n", "world", "name")
)
func main() {
	flag.Parse()
	// 连接服务器
	conn, err := grpc.Dial(*address, grpc.WithInsecure(), grpc.WithUnaryInterceptor(UnaryClientInterceptor),
		grpc.WithStreamInterceptor(StreamClientInterceptor))
	if err != nil {
		log.Fatalf("faild to connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)
	r, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: *name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.Message)
}
func UnaryClientInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	log.Printf("before invoker. method: %+v, request:%+v", method, req)
	err := invoker(ctx, method, req, reply, cc, opts...)
	log.Printf("after invoker. reply: %+v", reply)
	return err
}
func StreamClientInterceptor(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	log.Printf("before invoker. method: %+v, StreamDesc:%+v", method, desc)
	clientStream, err := streamer(ctx, desc, cc, method, opts...)
	log.Printf("before invoker. method: %+v", method)
	return clientStream, err
}
通过grpc.WithUnaryInterceptor、grpc.WithStreamInterceptor可以将拦截器传递给Dial做参数。在客户端调用的时候，可以查看拦截器输出的日志:

2017/04/17 23:34:20 before invoker. method: /pb.Greeter/SayHello, request:&HelloRequest{Name:world,}
2017/04/17 23:34:20 after invoker. reply: &HelloReply{Message:Hello world,}
2017/04/17 23:34:20 Greeting: Hello world
通过这个简单的例子，你可以很容易的了解拦截器的开发。unary和stream两种类型的拦截器可以根据你的gRPC server/client实现的不同，有选择的实现。

业务线的活动，每一次新活动都做独立项目开发，有大量重复代码，并且浪费数据服务的连接资源；排序服务也许要经常添加业务代码，目前是停服务发布……这些场景为了开发维护效率、稳定性、安全性和性能都使用了Go语言。Go是静态编译语言，在具体的动态场景该如何实现应用级别的持续交付呢？

基于k8s，nginx网关，队列回溯消费等工具的实现也可以实现不同程度的持续交付，但是持续交付的要求越高，搭建平台和维护的成本也越高。

从应用开发本身出发，可以考虑插件化。

插件使用场景特点
可以热更新式扩展应用程序的功能列表
应对多变的业务需求，方便功能上下线
对于任意的go应用，能进行增量架构、代码分发以及代码上下线
插件设计标准
性能：调用插件要尽可能的快；对于任务插件，使用单独的工作空间（协程、线程、进程的池子化处理），大的、慢的、长期运行的插件，要少调用
稳定性：插件依赖的发布平台要少发布，交互API的设计要做好抽象，上下文的环境变量非必须不添加，减少升级需求，甚至能支持多个实例互备热升级
可靠性：如果有失效、崩溃的可能，必须有快速、简单、完整的恢复机制；业务插件的执行不能影响依赖的发布平台的守护进程或者线程的稳定
安全性：应该通过代码签名之类的手段防篡改
扩展性：支持插件热更新和上下线，下线需要健康检查，公共库插件至少能热加载
复用性：业务插件不要太多一次性的上下线
易用性：提供使用简单、功能正交的API，业务插件能够获取依赖的发布平台的上下文和调用公共库
2. Go的插件方式
动态链接库plugin，官方文档
语言本身支持，插件和主程序原生语法交互

进程隔离：无，单进程
主程序调用插件：一切预协定object（包括function、channel）
插件感知主程序上下文：主程序预定义类型参数object（包括function、channel）
stream支持：单向，基于channel
插件发现：主程序循环扫描插件目录并维护状态；通过第三方文件diff工具维护，例如git
上线：能
下线：不能
更新：不能
通信：进程内
序列化：不需要
性能：高
Go plugin判断两个插件是否相同是通过比较pluginpath实现的，如果没有指定pluginpath，则由内部的算法生成, 生成的格式为plugin/unnamed-“ + root.Package.Internal.BuildID 。这种情况下，如果两个插件的文件名不同，引用包不同，或者引用的cgo不同，则会生成不同的插件，同时加载不会有问题。但是如果两个插件的文件名相同，相关的引用包也相同，则可能生成相同的插件，即使插件内包含的方法和变量不同，实现也不同。判断插件相同，热加载不会成功，也就意味着老插件不支持覆盖更新。

最好在编译的指定pluginpath，同时方便版本跟踪。目前生产环境建议一些公共库无服务依赖的函数，例如算法库之类的。

go build -ldflags "-pluginpath=plugin/hot-$(date +%s)" -buildmode=plugin -o so/Eng.so eng/greeter.go
通信+序列化
natefinch/pie，github仓库
进程隔离：有，多进程，provider+comsumer
主程序调用插件：provider模式调用插件进程中预协定method；consumer模式消费插件进程中的预协定参数object（包括function、除了channel）
插件感知主程序上下文：provider模式消费主程序的预定义参数object（包括function、除了channel）；consumer模式调用主程序中预定义method
stream支持：不支持
插件发现：主程序循环扫描插件目录并维护状态；通过第三方文件diff工具维护，例如git
上线：能
下线：能
更新：能
通信：支持stdin/stdout、pipe、unix socket、tcp、http、jsonrpc
序列化：gob，protobuf, json, xml
性能：中/偏高
基于Go的net/rpc库，无法支持主程序和插件之间的streaming数据交互，有golang的官方包[issue1]和[issue2]直接建议。另外，每一个插件都要开一个进程，因此要注意通信序列化的性能消耗和进程管理，默认使用stdin/stdout建立连接，如下图，一个plugin和主程序之间有两条单向连接。

可以上成产环境，要做好资源管理。

plugin

hashicorp/go-plugin，github仓库
进程隔离：有，多进程，server+client
主程序调用插件：一切协议预协定object
插件感知主程序上下文：一切协议预协定object
stream支持：单向和双向，基于http/2
插件发现：主程序循环扫描插件目录并维护状态；通过第三方文件diff工具维护，例如git
上线：能
下线：能
更新：能
通信：支持grpc
序列化：protobuf
性能：中/偏高
基于Google的grpc库，按照微服务的流程定义proto文件，能通信就能互相调用。知名团队出品。可以上成产环境，要做好资源管理。

grpc

go-mangos/mangos，github仓库
进程隔离：有，多进程，provider+comsumer
主程序调用插件：一切预协定object
插件感知主程序上下文：一切预协定object
stream支持：单向，基于mq
插件发现：主程序循环扫描插件目录并维护状态；通过第三方文件diff工具维护，例如git
上线：能
下线：能
更新：能
通信：支持mq
序列化：未知
性能：中/偏高
基于消息队列协议通信，nanomsg和ZeroMQ一类的规范包含一组预定义的通信拓扑（称为“可扩展性协议”），涵盖许多不同的场景：Pair，PubSub，Bus，Survey，Pipeline和ReqRep。Mangos是该协议的golang实现，能够灵活方便支地持两个插件交流。

可以上成产环境，要走大量的基础建设开发。

mq

嵌入式脚本语言
一般都是进程内内嵌第三方语言的解释器，需要考虑解释器的工作线程资源的重复利用。

embedscript

进程隔离：无，单进程，解释器有goroutine开销
主程序调用插件：一切语言协定object
插件感知主程序上下文：一切语言协定object
stream支持：看语言是否支持channel互通
插件发现：主程序循环扫描插件目录并维护状态；通过第三方文件diff工具维护，例如git
上线：能
下线：能
更新：能
通信：无
序列化：无
性能：中
go-like脚本语言，agora和七牛qlang
agora和qlang都是go语法的动态脚本语言，都好几年没维护了，建议不要用在生产环境，其中Qlang还有用户提[issue]觉得不稳定。

其他脚本语言，js-otto、go-lua5.1、go-lua5.2
otta支持目前受欢迎的js语法，star比较多，协定了大部分go原生支持的类型，不包括channel和goroutine，没有提供解释器的工作空间池子化管理，需要开发者使用goroutine和解释器的interrupt接口自行实现，但是从issue和TODO来看，也不适合生产环境。

gopher-lua支持lua5.1语法，和go交互的object类型比较完备，协定了大部分go原生支持的类型，包括channel和goroutine，有提供解释器的工作空间池子化管理，可以上生产环境。

go-lua支持lua5.2语法，目前不建议上生产环境。