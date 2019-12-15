---
title: ratelimit
layout: post
category: golang
author: 夏泽民
---
https://github.com/juju/ratelimit
https://github.com/juju/utils
限流算法
计数器法
计数器法是限流算法里最简单也是最容易实现的一种算法。维护一个单位时间内的Counter，当单位时间已经过去则将Counter重置零。这个算法虽然简单，但是有一个十分致命的问题，那就是临界问题。在临界时间的前一毫秒和后一毫秒都触发了最大的请求数，就会在两毫秒内发生了两倍单位时间的最大请求数量。

滑动窗口
如果接触过TCP协议的话，那么一定对滑动窗口这个名词不会陌生。在时间窗划分多个格子，每个格子都单独维护一个Counter，窗口每次滑动一个格子。指定时间窗最大请求数，也就是限制的时间范围内允许的最大请求数。计数器算法其实就是滑动窗口算法，只是它没有对时间窗口做进一步地划分，所以只有1格。当维护当滑动窗口的格子划分的越多，限流就会越精确。可是这种方式没有完全解决临界问题，时间窗内一小段流量可能占比特别大。

漏桶算法
首先，我们有一个固定容量的桶，有水流进来，也有水流出去。对于流进来的水来说，我们无法预计一共有多少水会流进来，也无法预计水流的速度。但是对于流出去的水来说，这个桶可以固定水流出的速率。而且，当桶满了之后，多余的水将会溢出。该算法保证以一个常速速率来处理请求，所以不会出现临界问题。

令牌桶算法
和漏桶算法效果类似但方向相反的算法。桶一开始是空的，token（令牌）以一个固定的速率r往桶里填充，直到达到桶的容量，多余的token将会被丢弃。每当一个请求过来时，就会尝试从桶里移除一个token，如果没有token的话，请求无法通过。令牌桶还可以方便的改变速度。 一旦需要提高速率,只要按需提高放入桶中的token的速率就行了。令牌桶算法允许流量一定程度的突发，取走token是不需要耗费时间的，也就是说，假设桶内有100个token时，那么可以瞬间允许100个请求通过。

算法总结
令牌桶算法由于实现简单，且允许某些流量的突发，对用户友好，所以被业界采用地较多。当然我们需要根据具体场景选择合适的算法。
<!-- more -->
在go中的使用
Go提供了一个package(golang.org/x/time/rate)，采用令牌桶的算法实现，用来方便的对速度进行限制。

type Limiter struct {
    limit Limit
    burst int
 
    mu     sync.Mutex
    tokens float64
    last time.Time
    lastEvent time.Time
}
 
func NewLimiter(r Limit, b int) *Limiter
首先创建一个rate.Limiter，其有两个参数，第一个参数为允许每秒发生多少次事件，第二个参数是其缓存最大可存多少个事件。这个桶一开始容量为b，装满b个token，然后每秒往里面填充r个token。由于令牌桶中最多有b个token，所以一次最多只能允许b个事件发生，一个事件花费掉一个token。

rate.Limiter提供三种主要的函数。

Wait/WaitN
func (lim *Limiter) Wait(ctx context.Context) (err error)
func (lim *Limiter) WaitN(ctx context.Context, n int) (err error)
Wait是WaitN(ctx, 1)的简化形式。WaitN阻塞当前直到lim允许n个事件的发生。当没有可用或足够的事件时，将阻塞等待，推荐实际程序中使用这个方法。

Allow/AllowN
func (lim *Limiter) Allow() bool
func (lim *Limiter) AllowN(now time.Time, n int) bool
Allow是函数AllowN(time.Now(), 1)的简化函数。AllowN标识在时间now的时候，n个事件是否可以同时发生(也意思就是now的时候是否可以从令牌桶中取n个token)。适合在超出频率的时候丢弃或跳过事件的场景。

Reserve/ReserveN
func (lim *Limiter) Reserve() *Reservation
func (lim *Limiter) ReserveN(now time.Time, n int) *Reservation
Reserve是ReserveN(time.Now(), 1)的简化形式。ReserveN返回对象Reservation，用于标识调用者需要等多久才能等到n个事件发生(意思就是等多久令牌桶中至少含有n个token)。Wait/WaitN和Allow/AllowN其实就是基于其之上实现的，通过sleep等待时间和直接返回状态。如果想对事件发生的频率和等待处理逻辑更加精细的话就可以使用它。

example
package main
 
import (
    "context"
    "fmt"
    "time"
 
    "golang.org/x/time/rate"
)
 
func main() {
    l := rate.NewLimiter(2, 5)
    ctx := context.Background()
    start := time.Now()
    // 要处理二十个事件
    for i := 0; i < 20; i++ {
        l.Wait(ctx)
        // dosomething
    }
    fmt.Println(time.Since(start)) // output: 7.501262697s （初始桶内5个和每秒2个token）
}
ratelimit
--导入"github.com/juju/ratelimit"

ratelimit软件包提供了高效的令牌桶实现。