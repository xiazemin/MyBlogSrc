---
title: NewIncomingContext
layout: post
category: golang
author: 夏泽民
---
在 gRPC 中对于 metadata 进行了区别，分为了传入和传出用的 metadata，这是为了防止 metadata 从入站 RPC 转发到其出站 RPC 的情况（详见 issues #1148），针对此提供了两种方法来分别进行设置，如下：

NewIncomingContext：创建一个附加了所传入的 md 新上下文，仅供自身的 gRPC 服务端内部使用。
NewOutgoingContext：创建一个附加了传出 md 的新上下文，可供外部的 gRPC 客户端、服务端使用。
因此相对的在 metadata 的获取上，也区分了两种方法，分别是 FromIncomingContext 和 NewOutgoingContext，与设置的方法所相对应的含义，如下：

md1, _ := metadata.FromIncomingContext(ctx)
md2, _ := metadata.FromOutgoingContext(ctx)
那么总的来说，这两种方法在实现上有没有什么区别呢，我们可以一起深入看看：

type mdIncomingKey struct{}
type mdOutgoingKey struct{}

func NewIncomingContext(ctx context.Context, md MD) context.Context {
	return context.WithValue(ctx, mdIncomingKey{}, md)
}

func NewOutgoingContext(ctx context.Context, md MD) context.Context {
	return context.WithValue(ctx, mdOutgoingKey{}, rawMD{md: md})
}
实际上主要是在内部进行了 Key 的区分，以所指定的 Key 来读取相对应的 metadata，以防造成脏读，其在实现逻辑上本质上并没有太大的区别。另外大家可以看到，其对 Key 的设置，是用一个结构体去定义的，这是 Go 语言官方一直在推荐的写法，建议大家也这么写。
<!-- more -->
https://golang2.eddycjy.com/posts/ch3/09-grpc-metadata-creds/
