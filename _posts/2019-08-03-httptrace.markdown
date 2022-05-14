---
title: httptrace
layout: post
category: golang
author: 夏泽民
---
net/http/httptrace主要是用于追踪客户端的 Request 请求过程中发生的各种事件及行为，在标准库 net/http/httptrace/trace.go 中定义了一个叫 ClientTrace 的结构体，它包含了一系列的钩子函数 hooks 作为成员变量
<!-- more -->
// ClientTrace is a set of hooks to run at various stages of an outgoing HTTP request. 
type ClientTrace struct {
    GetConn func(hostPort string)
    GotConn func(GotConnInfo)
    PutIdleConn func(err error)

    GotFirstResponseByte func()
    Got100Continue func()

    DNSStart func(DNSStartInfo)
    DNSDone func(DNSDoneInfo)

    ConnectStart func(network, addr string)
    ConnectDone func(network, addr string, err error)

    TLSHandshakeStart func()
    TLSHandshakeDone func(tls.ConnectionState, error)
    WroteHeaders func()
    Wait100Continue func()
    WroteRequest func(WroteRequestInfo)
}


race.go 还提供了一个 WithClientTrace() 包函数，用来把 ClientTrace 结构体中的钩子都保存（注册）到 Context 中去（因为 Context 提供 key/value 存储嘛），
key 就是一个叫 clientEventContextKey 的空结构体，value 是 nettrace 包中的 Trace 结构体，这个结构体作用跟 ClientTrace 一样，都是包含了一堆 hook 函数作为成员，
在这里它的目的只是封装下 ClientTrace 中的 hook 函数。最终， WithClientTrace() 会返回一个 context，它保存了上述的 hook 函数。
type clientEventContextKey struct{}
func WithClientTrace(ctx context.Context, trace *ClientTrace) context.Context {
    if trace == nil {
        panic("nil trace")
    }
    old := ContextClientTrace(ctx)
    trace.compose(old)

    ctx = context.WithValue(ctx, clientEventContextKey{}, trace)
    if trace.hasNetHooks() {
        nt := &nettrace.Trace{
            ConnectStart: trace.ConnectStart,
            ConnectDone:  trace.ConnectDone,
        }
        if trace.DNSStart != nil {
            nt.DNSStart = func(name string) {
                trace.DNSStart(DNSStartInfo{Host: name})
            }
        }
        if trace.DNSDone != nil {
            ...
        }
        ctx = context.WithValue(ctx, nettrace.TraceKey{}, nt)
    }
    return ctx
}

通过 ContextClientTrace() 的函数，可以把 ClientTrace 从 Context 中取出来。
// ContextClientTrace returns the ClientTrace associated with the
// provided context. If none, it returns nil.
func ContextClientTrace(ctx context.Context) *ClientTrace {
    trace, _ := ctx.Value(clientEventContextKey{}).(*ClientTrace)
    return trace
}

现在，我们知道，有了 WithClientTrace()，我们就可以把钩子函数保存在 Context 中了，现在，我们要把这些钩子函数挂到 Request 中去，该怎么弄？
很简单，通过 Request.WithContext() 把刚才赋值好的 Context 保存到 Request 中就可以了。
现在 Request 有了这些钩子函数，那么什么时候会被调用呢？ 当然会 http.Client.Do(req) 的时候啦。
接下来我们通过一段实际的代码看看整个流程：
package main

import (
    "fmt"
    "log"
    "net/http"
    "net/http/httptrace"
)

// transport is an http.RoundTripper that keeps track of the in-flight
// request and implements hooks to report HTTP tracing events.
type transport struct {
    current *http.Request
}

// RoundTrip wraps http.DefaultTransport.RoundTrip to keep track
// of the current request.
func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
    t.current = req
    return http.DefaultTransport.RoundTrip(req)
}

// GotConn prints whether the connection has been used previously
// for the current request.
func (t *transport) GotConn(info httptrace.GotConnInfo) {
    fmt.Printf("Connection reused for %v? %v\n", t.current.URL, info.Reused)
}

func main() {
    t := &transport{}

    req, _ := http.NewRequest("GET", "https://google.com", nil)
    trace := &httptrace.ClientTrace{
        GotConn: t.GotConn,
    }
    req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))

    client := &http.Client{Transport: t}
    if _, err := client.Do(req); err != nil {
        log.Fatal(err)
    }
}

所有的钩子的调用，最终都会在 client.Do(req) 里面执行，我们看看是怎么执行的。
注意到这里的 transport 结构体，它其实是 RoundTripper 接口类型（在 client.go 中声明）的一个 implementer，这个 RoundTripper 实际只有一个方法：
// RoundTripper is an interface representing the ability to execute a
// single HTTP transaction, obtaining the Response for a given Request.
type RoundTripper interface {
    // RoundTrip executes a single HTTP transaction, returning
    // a Response for the provided Request.
    RoundTrip(*Request) (*Response, error)
}

在 client.Do() 中，会调用 client.send()，如下：
 resp, didTimeout, err = c.send(req, deadline)

c.send() 内部：
// didTimeout is non-nil only if err != nil.
func (c *Client) send(req *Request, deadline time.Time) (resp *Response, didTimeout func() bool, err error) {
    ...
    resp, didTimeout, err = send(req, c.transport(), deadline)
    ...
    return resp, nil, nil
}

send() 内部：
func send(ireq *Request, rt RoundTripper, deadline time.Time) (resp *Response, didTimeout func() bool, err error) {
    ...
    resp, err = rt.RoundTrip(req)
    ...
}

可见，最终调用了 rt.RoundTrip() 函数。也就是上述 main.go 中 transport 实现的 RoundTrip() 函数。
在 rt.RoundTrip() 里面，把 req 赋给了 DefaultTransport.RoundTrip(req)，
这个 DefaultTransport 是包提供的一个 RoundTripper 的默认实现，
var DefaultTransport RoundTripper = &Transport{
    Proxy: ProxyFromEnvironment,
    DialContext: (&net.Dialer{
        Timeout:   30 * time.Second,
        KeepAlive: 30 * time.Second,
        DualStack: true,
    }).DialContext,
    MaxIdleConns:          100,
    IdleConnTimeout:       90 * time.Second,
    TLSHandshakeTimeout:   10 * time.Second,
    ExpectContinueTimeout: 1 * time.Second,
}

然后，在它的 RoundTrip() 函数里面最终会调用上述的钩子函数。
// RoundTrip implements the RoundTripper interface.
//
// For higher-level HTTP client support (such as handling of cookies
// and redirects), see Get, Post, and the Client type.
func (t *Transport) RoundTrip(req *Request) (*Response, error) {
    t.nextProtoOnce.Do(t.onceSetNextProtoDefaults)
    ctx := req.Context()
    trace := httptrace.ContextClientTrace(ctx)
    
    for {
        treq := &transportRequest{Request: req, trace: trace}
        cm, err := t.connectMethodForRequest(treq)
        ...
        pconn, err := t.getConn(treq, cm)
    }
}

解析：
通过调用 httptrace.ContextClientTrace(ctx) 把 context 中的钩子函数都取出来，再在 t.getConn() 中调用钩子函数，

