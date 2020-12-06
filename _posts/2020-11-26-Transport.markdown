---
title: Transport
layout: post
category: golang
author: 夏泽民
---
HTTP的连接
1. 串行连接
每次发起一次HTTP请求都需要一次链接，使用TCP的三次握手完成一次链接，发起请求获取数据。

必须等待上一次连接结束后在发起新的HTTP请求，建立新的TCP连接。
2. 持久连接
HTTP/1.1 允许客户端在发起请求结束后仍然保持在打开状态的TCP连接，便于后续请求继续使用。

减少TCP连接建立握手的时间延迟
减少了打开连接的潜在数量
3. 并行连接
通过多条TCP连接发起并发的HTTP请求。

每个连接之间有较小的时间延迟
每次发起的HTTP请求也是一个独立的TCP连接
并行连接数不能太多，占用本地CPU、内存、端口等各类资源，目前浏览器基本也都支持并行连接，一般限制连接数的值，比如4个
4. 管道化连接
通过共享的TCP连接发起并发的HTTP请求。

建立在持久化连接的基础上，将多条请求放入队列，一同发往服务端
降低网络的环回时间，提高性能
必须确保服务端支持持久化连接
做好连接会在任意时间关闭的准备，准备好重发所有未完成的管道化请求
Golang中持久化连接
在Golang中使用持久化连接发起HTTP请求，主要依赖Transport，官方封装的net库中已经支持。

Transport实现了RoundTripper接口，该接口只有一个方法RoundTrip()，故Transport的入口函数就是RoundTrip()。
Transport的主要功能：

缓存了长连接，用于大量http请求场景下的连接复用
对连接做一些限制，连接超时时间，每个host的最大连接数
在实际应用中，需要在初始化HTTP的client时传入transport，以进行保持连接，Transport的主要结构为：


type Transport struct {
    // DialContext specifies the dial function for creating unencrypted TCP connections.
    // If DialContext is nil (and the deprecated Dial below is also nil),
    // then the transport dials using package net.
    //
    // DialContext runs concurrently with calls to RoundTrip.
    // A RoundTrip call that initiates a dial may end up using
    // a connection dialed previously when the earlier connection
    // becomes idle before the later DialContext completes.
    DialContext func(ctx context.Context, network, addr string) (net.Conn, error)

    // MaxIdleConns controls the maximum number of idle (keep-alive)
    // connections across all hosts. Zero means no limit.
    MaxIdleConns int

    // MaxIdleConnsPerHost, if non-zero, controls the maximum idle
    // (keep-alive) connections to keep per-host. If zero,
    // DefaultMaxIdleConnsPerHost is used.
    MaxIdleConnsPerHost int

    // MaxConnsPerHost optionally limits the total number of
    // connections per host, including connections in the dialing,
    // active, and idle states. On limit violation, dials will block.
    //
    // Zero means no limit.
    MaxConnsPerHost int

    // IdleConnTimeout is the maximum amount of time an idle
    // (keep-alive) connection will remain idle before closing
    // itself.
    // Zero means no limit.
    IdleConnTimeout time.Duration

}
<!-- more -->
{% raw %}
具体使用Demo如下

package main

import (
    "fmt"
    "io/ioutil"
    "net"
    "net/http"
    "time"
)

var HTTPTransport = &http.Transport{
    DialContext: (&net.Dialer{
        Timeout:   30 * time.Second, // 连接超时时间
        KeepAlive: 60 * time.Second, // 保持长连接的时间
    }).DialContext, // 设置连接的参数
    MaxIdleConns:          500, // 最大空闲连接
    IdleConnTimeout:       60 * time.Second, // 空闲连接的超时时间
    ExpectContinueTimeout: 30 * time.Second, // 等待服务第一个响应的超时时间
    MaxIdleConnsPerHost:   100, // 每个host保持的空闲连接数
}

func main() {
    times := 50
    uri := "http://local.test.com/t.php"
    // uri := "http://www.baidu.com"
    

    // 短连接的情况

    start := time.Now()
    client := http.Client{} // 初始化http的client
    for i := 0; i < times; i++ {
        req, err := http.NewRequest(http.MethodGet, uri, nil)
        if err != nil {
            panic("Http Req Failed " + err.Error())
        }
        resp, err := client.Do(req) // 发起请求
        if err != nil {
            panic("Http Request Failed " + err.Error())
        }
        defer resp.Body.Close()
        ioutil.ReadAll(resp.Body)
    }
    fmt.Println("Orig GoNet Short Link", time.Since(start))
    

    // 长连接的情况

    start2 := time.Now()
    client2 := http.Client{Transport: HTTPTransport} // 初始化一个带有transport的http的client
    for i := 0; i < times; i++ {
        req, err := http.NewRequest(http.MethodGet, uri, nil)
        if err != nil {
            panic("Http Req Failed " + err.Error())
        }
        resp, err := client2.Do(req)
        if err != nil {
            panic("Http Request Failed " + err.Error())
        }
        defer resp.Body.Close()
        ioutil.ReadAll(resp.Body) // 如果不及时从请求中获取结果，此连接会占用，其他请求服务复用连接
    }
    fmt.Println("Orig GoNet Long Link", time.Since(start2))
}

经过本地测试，使用transport确实能控制客户端的连接数，使得本地资源使用得到大幅度的降低。通过netstat可以查看具体的连接情况：

{% endraw %}

https://github.com/aceld/zinx
https://www.bilibili.com/video/av71067087/
https://gitee.com/caoxy/zinx
