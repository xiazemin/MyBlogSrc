---
title: groupcache 源码分析
layout: post
category: golang
author: 夏泽民
---
https://github.com/golang/groupcache

雪崩就是指缓存中大批量热点数据过期 或 缓存机器意外发生了全盘宕机后系统涌入大量查询请求，因为大部分数据在Redis层已经失效，请求渗透到数据库层，大批量请求犹如洪水一般涌入，引起数据库压力造成查询堵塞甚至宕机。

解决办法：

将缓存失效时间分散开，比如每个key的过期时间是随机，防止同一时间大量数据过期现象发生，这样不会出现同一时间全部请求都落在数据库层，如果缓存数据库是分布式部署，将热点数据均匀分布在不同Redis和数据库中，有效分担压力，别一个人扛。
简单粗暴，让Redis数据永不过期（如果业务准许，比如不用更新的名单类）。当然，如果业务数据准许的情况下可以，比如中奖名单用户，每期用户开奖后，名单不可能会变了，无需更新。

缓存雪崩的事前事中事后的解决方案如下。 - 事前：redis 高可用，主从+哨兵，redis cluster，避免全盘崩溃。 - 事中：本地 ehcache 缓存 + hystrix 限流&降级，避免 MySQL 被打死。 - 事后：redis 持久化，一旦重启，自动从磁盘上加载数据，快速恢复缓存数据。

缓存击穿，就是说某个 key 非常热点，访问非常频繁，处于集中式高并发访问的情况，当这个 key 在失效的瞬间，大量的请求就击穿了缓存，直接请求数据库，就像是在一道屏障上凿开了一个洞。

解决方式也很简单，可以将热点数据设置为永远不过期；或者基于 redis or zookeeper 实现互斥锁，等待第一个请求构建完缓存之后，再释放锁，进而其它请求才能通过该 key 访问数据。

缓存穿透 
假设一秒 5000 个请求，结果其中 4000 个请求是黑客发出的恶意攻击。
　　黑客发出的那 4000 个攻击，缓存中查不到，每次你去数据库里查，也查不到。
　　举个栗子。数据库 id 是从 1 开始的，结果黑客发过来的请求 id 全部都是负数。这样的话，缓存中不会有，请求每次都“视缓存于无物”，直接查询数据库。这种恶意攻击场景的缓存穿透就会直接把数据库给打死。

解决方式很简单，每次系统 A 从数据库中只要没查到，就写一个空值到缓存里去，比如 set -999 UNKNOWN。然后设置一个过期时间，这样的话，下次有相同的 key 来访问的时候，在缓存失效之前，都可以直接从缓存中取数据。

https://www.cnblogs.com/myseries/p/12853369.html
https://blog.csdn.net/kongtiao5/article/details/82771694
https://www.jianshu.com/p/b7f822935e28

Package singleflight provides a duplicate function call suppression mechanism.

翻译过来就是：singleflight包提供了一种抑制重复函数调用的机制。

具体到Go程序运行的层面来说，SingleFlight的作用是在处理多个goroutine同时调用同一个函数的时候，只让一个goroutine去实际调用这个函数，等到这个goroutine返回结果的时候，再把结果返回给其他几个同时调用了相同函数的goroutine，这样可以减少并发调用的数量。在实际应用中也是，它能够在一个服务中减少对下游的并发重复请求。还有一个比较常见的使用场景是用来防止缓存击穿。

Go提供的SingleFlight
Go扩展库里用singleflight.Group结构体类型提供了SingleFlight并发原语的功能。

singleflight.Group类型提供了三个方法：

func (g *Group) Do(key string, fn func() (interface{}, error)) (v interface{}, err error, shared bool)
 
func (g *Group) DoChan(key string, fn func() (interface{}, error)) <-chan Result
 
func (g *Group) Forget(key string)
Do方法，接受一个字符串Key和一个待调用的函数，会返回调用函数的结果和错误。使用Do方法的时候，它会根据提供的Key判断是否去真正调用fn函数。同一个 key，在同一时间只有第一次调用Do方法时才会去执行fn函数，其他并发的请求会等待调用的执行结果。

DoChan方法：类似Do方法，只不过是一个异步调用。它会返回一个通道，等fn函数执行完，产生了结果以后，就能从这个 chan 中接收这个结果。

Forget方法：在SingleFlight中删除一个Key。这样一来，之后这个Key的Do方法调用会执行fn函数，而不是等待前一个未完成的fn 函数的结果。

<!-- more -->
应用场景
了解了Go语言提供的 SingleFlight并发原语都有哪些方法可以调用后 ，下面介绍两个它的应用场景。

查询DNS记录
Go语言的net标准库里使用的lookupGroup结构，就是Go扩展库提供的原语singleflight.Group

type Resolver struct {
  ......
 // 源码地址 https://github.com/golang/go/blob/master/src/net/lookup.go#L151
 // lookupGroup merges LookupIPAddr calls together for lookups for the same
 // host. The lookupGroup key is the LookupIPAddr.host argument.
 // The return values are ([]IPAddr, error).
 lookupGroup singleflight.Group
}

防止缓存击穿
在项目里使用缓存时，一个常见的用法是查询一个数据先去查询缓存，如果没有就去数据库里查到数据并缓存到Redis里。那么缓存击穿问题是指，高并发的系统中，大量的请求同时查询一个缓存Key 时，如果这个 Key 正好过期失效，就会导致大量的请求都打到数据库上，这就是缓存击穿。用 SingleFlight 来解决缓存击穿问题再合适不过，这个时候只要这些对同一个 Key 的并发请求的其中一个到数据库中查询就可以了，这些并发的请求可以共享同一个结果

Do方法
SingleFlight 定义一个call结构体，每个结构体都保存了fn调用对应的信息。

Do方法的执行逻辑是每次调用Do方法都会先去获取互斥锁，随后判断在映射表里是否已经有Key对应的fn函数调用信息的call结构体。

当不存在时，证明是这个Key的第一次请求，那么会初始化一个call结构体指针，增加SingleFlight内部持有的sync.WaitGroup计数器到1。释放互斥锁，然后阻塞的等待doCall方法执行fn函数的返回结果

当存在时，增加call结构体内代表fn重复调用次数的计数器dups，释放互斥锁，然后使用WaitGroup等待fn函数执行完成。

call结构体的val 和 err 两个字段只会在 doCall方法中执行fn有返回结果后才赋值，所以当 doCall方法 和 WaitGroup.Wait返回时，函数调用的结果和错误会返回给Do方法的所有调用者。

https://blog.csdn.net/kevin_tech/article/details/111878251

https://www.jianshu.com/p/7f3792549346

缓存击穿
什么是缓存击穿
平常在高并发系统中，会出现大量的请求同时查询一个key的情况，假如此时这个热key刚好失效了，就会导致大量的请求都打到数据库上面去，这种现象就是缓存击穿。缓存击穿和缓存雪崩有点像，但是又有一点不一样，缓存雪崩是因为大面积的缓存失效，打崩了DB，而缓存击穿则是指一个key非常热点，在不停的扛着高并发，高并发集中对着这一个点进行访问，如果这个key在失效的瞬间，持续的并发到来就会穿破缓存，直接请求到数据库，就像一个完好无损的桶上凿开了一个洞，造成某一时刻数据库请求量过大，压力剧增!

如何解决
方法一
我们简单粗暴点，直接让热点数据永远不过期，定时任务定期去刷新数据就可以了。不过这样设置需要区分场景，比如某宝首页可以这么做。

方法二
为了避免出现缓存击穿的情况，我们可以在第一个请求去查询数据库的时候对他加一个互斥锁，其余的查询请求都会被阻塞住，直到锁被释放，后面的线程进来发现已经有缓存了，就直接走缓存，从而保护数据库。但是也是由于它会阻塞其他的线程，此时系统吞吐量会下降。需要结合实际的业务去考虑是否要这么做。

方法三
方法三就是singleflight的设计思路，也会使用互斥锁，但是相对于方法二的加锁粒度会更细，这里先简单总结一下singleflight的设计原理，后面看源码在具体分析。

singleflightd的设计思路就是将一组相同的请求合并成一个请求，使用map存储，只会有一个请求到达mysql，使用sync.waitgroup包进行同步，对所有的请求返回相同的结果。

ype call struct { 
 wg sync.WaitGroup 
 // 存储返回值，在wg done之前只会写入一次 
 val interface{} 
  // 存储返回的错误信息 
 err error 
 
 // 标识别是否调用了Forgot方法 
 forgotten bool 
 
 // 统计相同请求的次数，在wg done之前写入 
 dups  int 
  // 使用DoChan方法使用，用channel进行通知 
 chans []chan<- Result 
} 
// Dochan方法时使用 
type Result struct { 
 Val    interface{} // 存储返回值 
 Err    error // 存储返回的错误信息 
 Shared bool // 标示结果是否是共享结果 
} 

// 入参：key：标识相同请求，fn：要执行的函数 
// 返回值：v: 返回结果 err: 执行的函数错误信息 shard: 是否是共享结果 
func (g *Group) Do(key string, fn func() (interface{}, error)) (v interface{}, err error, shared bool) { 
 // 代码块加锁 
 g.mu.Lock() 
 // map进行懒加载 
 if g.m == nil { 
   // map初始化 
  g.m = make(map[string]*call) 
 } 
 // 判断是否有相同请求 
 if c, ok := g.m[key]; ok { 
   // 相同请求次数+1 
  c.dups++ 
  // 解锁就好了，只需要等待执行结果了，不会有写入操作了 
  g.mu.Unlock() 
  // 已有请求在执行，只需要等待就好了 
  c.wg.Wait() 
  // 区分panic错误和runtime错误 
  if e, ok := c.err.(*panicError); ok { 
   panic(e) 
  } else if c.err == errGoexit { 
   runtime.Goexit() 
  } 
  return c.val, c.err, true 
 } 
 // 之前没有这个请求，则需要new一个指针类型 
 c := new(call) 
 // sync.waitgroup的用法，只有一个请求运行，其他请求等待，所以只需要add(1) 
 c.wg.Add(1) 
 // m赋值 
 g.m[key] = c 
 // 没有写入操作了，解锁即可 
 g.mu.Unlock() 
 // 唯一的请求该去执行函数了 
 g.doCall(c, key, fn) 
 return c.val, c.err, c.dups > 0 
}

https://developer.51cto.com/art/202107/672248.htm

https://www.cnblogs.com/xiaozhe97/p/13702010.html

https://zhuanlan.zhihu.com/p/343761986

```
import "sync"

// call is an in-flight or completed Do call
type call struct {
    wg  sync.WaitGroup
    val interface{}
    err error
}

// Group represents a class of work and forms a namespace in which
// units of work can be executed with duplicate suppression.
type Group struct {
    mu sync.Mutex       // protects m
    m  map[string]*call // lazily initialized
}

// Do executes and returns the results of the given function, making
// sure that only one execution is in-flight for a given key at a
// time. If a duplicate comes in, the duplicate caller waits for the
// original to complete and receives the same results.
//同一个对象多次同时多次调用这个逻辑的时候，可以使用其中的一个去执行
func (g *Group) Do(key string, fn func()(interface{},error)) (interface{}, error ){
    g.mu.Lock() //加锁保护存放key的map，因为要并发执行
    if g.m == nil { //lazing make 方式建立
        g.m = make(map[string]*call)
    }
    if c, ok := g.m[key]; ok { //如果map中已经存在对这个key的处理那就等着吧
        g.mu.Unlock() //解锁，对map的操作已经完毕
        c.wg.Wait()
        return c.val,c.err //map中只有一份key，所以只有一个c
    }
    c := new(call) //创建一个工作单元，只负责处理一种key
    c.wg.Add(1)
    g.m[key] = c //将key注册到map中
    g.mu.Unlock() //map的操做完成，解锁
    
    c.val, c.err = fn()//第一个注册者去执行
    c.wg.Done()
    
    g.mu.Lock()
    delete(g.m,key) //对map进行操作，需要枷锁
    g.mu.Unlock()
    
    return c.val, c.err //给第一个注册者返回结果
}
```
https://www.jianshu.com/p/f47feb4720f9
https://www.cnblogs.com/softlin/p/14133635.html

https://segmentfault.com/a/1190000039712358?utm_source=sf-similar-article


代码的结构：

consistanthash        实现一致性hash功能

lru                   实现缓存的置换算法（最近最少使用）

singleflight          实现多个同请求的合并，保证“同时”多个同参数的get请求只执行一次操作功能

groupcachepb          grpc生成的代码，用于远程调用

byteview.go           将byte于string进行了一次封装，对外提供不区分两者的接口

groupcache.go         groupcache的主API函数

http.go               http相关的代码

peers.go              单个节点的一些接口实现

sinks.go              暂时没太搞明白


https://blog.csdn.net/mrbuffoon/article/details/83510799

https://www.cnblogs.com/B0-1/p/5799094.html

