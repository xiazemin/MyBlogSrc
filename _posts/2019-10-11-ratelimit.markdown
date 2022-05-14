---
title: ratelimit
layout: post
category: golang
author: 夏泽民
---
https://github.com/didip/tollbooth

https://github.com/uber-go/ratelimit
https://github.com/tidwall/limiter

https://github.com/yangwenmai/ratelimit

https://github.com/GuoZhaoran/rateLimit
https://github.com/GuoZhaoran/Go-RateLimit

golang.org/x/time
https://github.com/golang/time
一般有漏桶算法和令牌桶算法两种限流算法。

漏桶算法
漏桶算法(Leaky Bucket)是网络世界中流量整形（Traffic Shaping）或速率限制（Rate Limiting）时经常使用的一种算法，它的主要目的是控制数据注入到网络的速率，平滑网络上的突发流量。漏桶算法提供了一种机制，通过它，突发流量可以被整形以便为网络提供一个稳定的流量。

漏桶可以看作是一个带有常量服务时间的单服务器队列，如果漏桶（包缓存）溢出，那么数据包会被丢弃。 在网络中，漏桶算法可以控制端口的流量输出速率，平滑网络上的突发流量，实现流量整形，从而为网络提供一个稳定的流量。

把请求比作是水，水来了都先放进桶里，并以限定的速度出水，当水来得过猛而出水不够快时就会导致水直接溢出，即拒绝服务。

令牌桶算法
令牌桶算法是网络流量整形（Traffic Shaping）和速率限制（Rate Limiting）中最常使用的一种算法。典型情况下，令牌桶算法用来控制发送到网络上的数据的数目，并允许突发数据的发送。

令牌桶算法的原理是系统会以一个恒定的速度往桶里放入令牌，而如果请求需要被处理，则需要先从桶里获取一个令牌，当桶里没有令牌可取时，则拒绝服务。从原理上看，令牌桶算法和漏桶算法是相反的，一个“进水”，一个是“漏水”。

Google的Guava包中的RateLimiter类就是令牌桶算法的解决方案。

漏桶算法和令牌桶算法的选择
漏桶算法与令牌桶算法在表面看起来类似，很容易将两者混淆。但事实上，这两者具有截然不同的特性，且为不同的目的而使用。

漏桶算法与令牌桶算法的区别在于，漏桶算法能够强行限制数据的传输速率，令牌桶算法能够在限制数据的平均传输速率的同时还允许某种程度的突发传输。

需要注意的是，在某些情况下，漏桶算法不能够有效地使用网络资源，因为漏桶的漏出速率是固定的，所以即使网络中没有发生拥塞，漏桶算法也不能使某一个单独的数据流达到端口速率。因此，漏桶算法对于存在突发特性的流量来说缺乏效率。而令牌桶算法则能够满足这些具有突发特性的流量。通常，漏桶算法与令牌桶算法结合起来为网络流量提供更高效的控制。

两者主要区别在于“漏桶算法”能够强行限制数据的传输速率，而“令牌桶算法”在能够限制数据的平均传输速率外，还允许某种程度的突发传输。在“令牌桶算法”中，只要令牌桶中存在令牌，那么就允许突发地传输数据直到达到用户配置的门限，所以它适合于具有突发特性的流量。
<!-- more -->
golang.org/x/time/rate
type Limiter
type Limiter struct {
    // contains filtered or unexported fields
}
1
2
3
Limter限制时间的发生频率，采用令牌池的算法实现。这个池子一开始容量为b，装满b个令牌，然后每秒往里面填充r个令牌。
由于令牌池中最多有b个令牌，所以一次最多只能允许b个事件发生，一个事件花费掉一个令牌。

Limter提供三中主要的函数 Allow, Reserve, and Wait. 大部分时候使用Wait。

func NewLimiter
func NewLimiter(r Limit, b int) *Limiter
1
NewLimiter 返回一个新的Limiter。

func (*Limiter) [Allow]
func (lim *Limiter) Allow() bool
1
Allow 是函数 AllowN(time.Now(), 1)的简化函数。

func (*Limiter) AllowN
func (lim *Limiter) AllowN(now time.Time, n int) bool
1
AllowN标识在时间now的时候，n个事件是否可以同时发生(也意思就是now的时候是否可以从令牌池中取n个令牌)。如果你需要在事件超出频率的时候丢弃或跳过事件，就使用AllowN,否则使用Reserve或Wait.

func (*Limiter) Reserve
func (lim *Limiter) Reserve() *Reservation
1
Reserve是ReserveN(time.Now(), 1).的简化形式。

func (*Limiter) ReserveN
func (lim *Limiter) ReserveN(now time.Time, n int) *Reservation
1
ReserveN 返回对象Reservation ，标识调用者需要等多久才能等到n个事件发生(意思就是等多久令牌池中至少含有n个令牌)。

如果ReserveN 传入的n大于令牌池的容量b，那么返回false.
使用样例如下：

r := lim.ReserveN(time.Now(), 1)
if !r.OK() {
  // Not allowed to act! Did you remember to set lim.burst to be > 0 ?我只要1个事件发生仍然返回false，是不是b设置为了0？
  return
}
time.Sleep(r.Delay())
Act()

如果希望根据频率限制等待和降低事件发生的速度而不丢掉事件，就使用这个方法。
我认为这里要表达的意思就是如果事件发生的频率是可以由调用者控制的话，可以用ReserveN 来控制事件发生的速度而不丢掉事件。如果要使用context的截止日期或cancel方法的话，使用WaitN。

func (*Limiter) Wait
func (lim *Limiter) Wait(ctx context.Context) (err error)
1
Wait是WaitN(ctx, 1)的简化形式。

func (*Limiter) WaitN
func (lim *Limiter) WaitN(ctx context.Context, n int) (err error)
1
WaitN 阻塞当前直到lim允许n个事件的发生。
- 如果n超过了令牌池的容量大小则报错。
- 如果Context被取消了则报错。
- 如果lim的等待时间超过了Context的超时时间则报错。

golang.org/x/time/rate 提对速度进行限制的算法

l := rate.NewLimiter(1, 3) // 一个参数为每秒发生多少次事件，第二个参数是最大可运行多少个事件(burst)
1
Limter提供三中主要的函数 Allow, Reserve, Wait. 大部分时候使用Wait

Wait/WaitN 当没有可用事件时，将阻塞等待
c, _ := context.WithCancel(context.TODO())
for {
    l.Wait(c)
    fmt.Println(time.Now().Format("04:05.000"))
}
输出
07:35.055
07:35.055
07:35.055
07:36.060
07:37.059
07:38.059
缓存3次后，每秒执行一次
Allow/AllowN 当没有可用事件时，返回false
for {
    if (l.AllowN(time.Now(), 1)) {
        fmt.Println(time.Now().Format("04:05.000"))
    } else {
        time.Sleep(1 * time.Second / 10)
        fmt.Println(time.Now().Format("Second 04:05.000"))
    }

}

Reserve/ReserveN 当没有可用事件时，返回 Reservation，和要等待多久才能获得足够的事件
for {
        r := l.ReserveN(time.Now(), 1)
        s := r.Delay()
        time.Sleep(s)
            fmt.Println(s, time.Now().Format("04:05.000"))

    }


Go rate limiter
This package provides a Golang implementation of the leaky-bucket rate limit algorithm. This implementation refills the bucket based on the time elapsed between requests instead of requiring an interval clock to fill the bucket discretely.

Create a rate limiter with a maximum number of operations to perform per second. Call Take() before each operation. Take will sleep until you can continue.

import (
	"fmt"
	"time"

	"go.uber.org/ratelimit"
)

func main() {
    rl := ratelimit.New(100) // per second

    prev := time.Now()
    for i := 0; i < 10; i++ {
        now := rl.Take()
        fmt.Println(i, now.Sub(prev))
        prev = now
    }

    // Output:
    // 0 0
    // 1 10ms
    // 2 10ms
    // 3 10ms
    // 4 10ms
    // 5 10ms
    // 6 10ms
    // 7 10ms
    // 8 10ms
    // 9 10ms
}


Go 提供了一个package(golang.org/x/time/rate) 用来方便的对速度进行限制,
首先创建一个rate.Limiter,其有两个参数，第一个参数为每秒发生多少次事件，第二个参数是其缓存最大可存多少个事件。

rate.Limiter提供了三类方法用来限速

Wait/WaitN 当没有可用或足够的事件时，将阻塞等待 推荐实际程序中使用这个方法
Allow/AllowN 当没有可用或足够的事件时，返回false
Reserve/ReserveN 当没有可用或足够的事件时，返回 Reservation，和要等待多久才能获得足够的事件。


1,简单的并发控制
利用 channel 的缓冲设定，我们就可以来实现并发的限制。我们只要在执行并发的同时，往一个带有缓冲的 channel 里写入点东西（随便写啥，内容不重要）。让并发的 goroutine在执行完成后把这个 channel 里的东西给读走。这样整个并发的数量就讲控制在这个 channel的缓冲区大小上。


2,使用计数器实现请求限流
限流的要求是在指定的时间间隔内，server 最多只能服务指定数量的请求。实现的原理是我们启动一个计数器，每次服务请求会把计数器加一，同时到达指定的时间间隔后会把计数器清零；

3,使用golang官方包实现httpserver频率限制
使用golang来编写httpserver时，可以使用官方已经有实现好的包：
import(
    "fmt"
    "net"
    "golang.org/x/net/netutil"
)
 
func main() {
    l, err := net.Listen("tcp", "127.0.0.1:0")
    if err != nil {
        fmt.Fatalf("Listen: %v", err)
    }
    defer l.Close()
    l = LimitListener(l, max)
    
    http.Serve(l, http.HandlerFunc())
    
    //bla bla bla.................
}

源码如下(url : https://github.com/golang/net/blob/master/netutil/listen.go)，基本思路就是为连接数计数，通过make chan来建立一个最大连接数的channel, 每次accept就+1，close时候就-1. 当到达最大连接数时，就等待空闲连接出来之后再accept。

4,使用Token Bucket（令牌桶算法）实现请求限流
在开发高并发系统时有三把利器用来保护系统：缓存、降级和限流!为了保证在业务高峰期，线上系统也能保证一定的弹性和稳定性，最有效的方案就是进行服务降级了，而限流就是降级系统最常采用的方案之一。

这里为大家推荐一个开源库https://github.com/didip/tollbooth，但是，如果您想要一些简单的、轻量级的或者只是想要学习的东西，实现自己的中间件来处理速率限制并不困难。今天我们就来聊聊如何实现自己的一个限流中间件

首先我们需要安装一个提供了 Token bucket (令牌桶算法)的依赖包，上面提到的toolbooth 的实现也是基于它实现的：

$ go get golang.org/x/time/rate
