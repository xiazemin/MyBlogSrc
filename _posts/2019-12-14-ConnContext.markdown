---
title: ConnContext Go1.13 标准库的 http 内存泄漏
layout: post
category: golang
author: 夏泽民
---
2019 年 11 月 21 日，golang 的官方 github 仓库提交了一个 https://github.com/golang/go/issues/35750 ,该 issue 指出如果初始化 http.Server 结构体时指定了一个非空的 ConnContext 成员函数，且如果在该函数内使用了 context.Value 方法写入数据到上下文并返回，则 Context 将会以链表的方式泄漏。

这是一个很恐怖的 bug，因为一方面 http.Server 是几乎所有市面上流行 web 框架如 gin，beego 的底层依赖，一旦发生问题则全部中招，一个都跑不了；另一方面则 ConnContext 函数在每一个请求初始化时都会被调用，这意味着如果一旦发生泄漏，则服务端程序几乎一定会溢出。

据官方开发人员表示，该 bug 于 1.13 版本引进，目前已经在 1.13.5 修复。
<!-- more -->

影响范围
所有 1.13~1.13.4 版本，使用原生 http.Server 指定了 ConnContext 成员函数，且在该函数中使用 With*方法写入数据并返回新 Context；或者使用了上层框架的相应功能。

现象
内存根据访问量持续上升，且 pprof 分析发现 cpu 大量耗费在 Context 的底层方法上。

故障原理
问题出在 Server.Serve 方法，该方法是 http.Server 启动的底层唯一入口，负责循环 Accept 连接以及为每个新连接开启 goroutine 做下一步处理。

来看看出问题的代码，为了简洁这里略去不必要代码:

type Server struct {
    ...
	// ConnContext optionally specifies a function that modifies
	// the context used for a new connection c. The provided ctx
	// is derived from the base context and has a ServerContextKey
	// value.
	ConnContext func(ctx context.Context, c net.Conn) context.Context
    ...
}

func (srv *Server) Serve(l net.Listener) error {
        ...
	baseCtx := context.Background()
        ...
	ctx := context.WithValue(baseCtx, ServerContextKey, srv)
	for {
		rw, e := l.Accept()
	        ...
		if cc := srv.ConnContext; cc != nil {
			ctx = cc(ctx, rw)
			if ctx == nil {
				panic("ConnContext returned nil")
			}
		}
                ...
		go c.serve(ctx)
	}
}
我们知道 Context 是一个反向链表结构，从最初的 Background 通过各种 With 方法推入表头节点，而 With 方法返回的则是新的表头节点。

从上边的代码中我们看到，如果 srv.ConnContext 不为空，则每次 Accept 连接后都会调用此函数并传入 ctx，然后将返回的结果存入 ctx 中，这意味着如果在此函数中使用 With 函数写入节点并返回，则该节点将被缓存到全局的 ctx，从而造成泄漏。

复现
这个 bug 非常容易复现，下面我们复现一下：

// go version:1.13.4

func main() {
	var count int32 = 0
	server := &http.Server{
		Addr: ":4444",
		Handler: http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.Header().Set("Connection", "close")
		}),
		ConnContext: func(ctx context.Context, c net.Conn) context.Context {
			atomic.AddInt32(&count, 1)
			if c2 := ctx.Value("count"); c2 != nil {
				fmt.Printf("发现了遗留数据: %d\n", c2.(int32))
			}
			fmt.Printf("本次数据: %d\n", count)
			return context.WithValue(ctx, "count", count)
		},
	}
	gofunc() {
		panic(server.ListenAndServe())
	}()

	var err error

	fmt.Println("第一次请求")
	_, err = http.Get("http://localhost:4444/")
	if err != nil {
		panic(err)
	}
	fmt.Println("\n第二次请求")

	_, err = http.Get("http://localhost:4444/")
	if err != nil {
		panic(err)
	}
}
结果：

第一次请求
本次数据: 1

第二次请求
发现了遗留数据: 1
本次数据: 2
可以看到，第二个从请求的 Context 中能读取到第一个请求的 Context 中写入的数据，确实发生了泄漏。

修复
我们首先要理解 ConnContext 这个函数的作用，按照设计它应该是为每个请求的 Context 做一些初始化处理，然后将这个处理后的 Context 链传入go c.serve(ctx)，而不应该缓存到全局；下一个请求过来后应该将原始的 Context 传入 ConnContext 进行处理，从而得到新的 Context 链。

明白了目的，再看看问题代码，我们发现罪魁祸首在这里

ctx = cc(ctx, rw)
这一行错误地将 cc 方法生成的新链缓存到了全局，导致泄漏(ps:实在是搞不懂 google 的大神居然会犯这么低级且致命的错误...)。

修复后的代码如下：

func (srv *Server) Serve(l net.Listener) error {
	...
	baseCtx := context.Background()
	...
	ctx := context.WithValue(baseCtx, ServerContextKey, srv)
	for {
		rw, e := l.Accept()
		...
		connCtx = ctx
		if cc := srv.ConnContext; cc != nil {
			connCtx = cc(connCtx, rw)
			if connCtx == nil {
				panic("ConnContext returned nil")
			}
		}
		...
		go c.serve(connCtx)
	}
}


https://github.com/golang/go/issues/35750