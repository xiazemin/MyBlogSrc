---
title: Singleflight
layout: post
category: golang
author: 夏泽民
---
https://github.com/golang/groupcache
缓存更新问题
当缓存失效时，需要去数据存储层获取数据，然后存储到缓存中。

通常缓存更新方案：

业务代码中，根据key从缓存拿不到数据，访问存储层获取数据后更新缓存
由专门的定时脚本在缓存失效前对其进行更新
通过分布式锁，实现只有一个请求负责缓存更新，其他请求等待：一种基于哨兵的缓存访问策略
服务中某个接口请求量暴增问题
比如某个帖子突然很火，帖子下有非常多的跟帖回复，负责提供帖子内容、回帖内容的接口，对于该帖子的请求量就会非常多。

如果每个请求都落到下游服务，通常会导致下游服务瞬时负载升高。如果使用缓存，如何判断当前接口请求的内容需要缓存下来？缓存的过期、更新问题？

golang singleflight
该库提供了一个简单有效的方案应对上面提到的问题，初次见识到 singleflight 是在 golang/groupcache 中。

groupcache 缓存更新能够做到对同一个失效key的多个请求，只有一个请求执行对key的更新操作，其文档相关描述如下：

comes with a cache filling mechanism. Whereas memcached just says “Sorry, cache miss”, often resulting in a thundering herd of database (or whatever) loads from an unbounded number of clients (which has resulted in several fun outages), groupcache coordinates cache fills such that only one load in one process of an entire replicated set of processes populates the cache, then multiplexes the loaded value to all callers.

从 singleflight 的 test 可以了解到其用法：

func TestDoDupSuppress(t *testing.T) {
	var g Group
	c := make(chan string)
	var calls int32
	fn := func() (interface{}, error) {
		atomic.AddInt32(&calls, 1)
		return <-c, nil
	}

	const n = 10
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() { // n个协程同时调用了g.Do，fn中的逻辑只会被一个协程执行
			v, err := g.Do("key", fn)
			if err != nil {
				t.Errorf("Do error: %v", err)
			}
			if v.(string) != "bar" {
				t.Errorf("got %q; want %q", v, "bar")
			}
			wg.Done()
		}()
	}
	time.Sleep(100 * time.Millisecond) // let goroutines above block
	c <- "bar"
	wg.Wait()
	if got := atomic.LoadInt32(&calls); got != 1 {
		t.Errorf("number of calls = %d; want 1", got)
	}
}
该测试用例中，只有1个协程执行了fn，其他9个协程能拿到fn执行后的返回结果。即fn只执行了1次，但其结果会返回给多个协程。

看下 singleflight 是如何做到这一点的：

// call is an in-flight or completed Do call
type call struct {
	wg  sync.WaitGroup
	val interface{}
	err error
}
call 用来表示一个正在执行或已完成的函数调用。

// Group represents a class of work and forms a namespace in which
// units of work can be executed with duplicate suppression.
type Group struct {
	mu sync.Mutex       // protects m
	m  map[string]*call // lazily initialized
}
Group 可以看做是任务的分类。

// Do executes and returns the results of the given function, making
// sure that only one execution is in-flight for a given key at a
// time. If a duplicate comes in, the duplicate caller waits for the
// original to complete and receives the same results.
func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call)
	}
	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		c.wg.Wait()
		return c.val, c.err
	}
	c := new(call)
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()

	c.val, c.err = fn()
	c.wg.Done()

	g.mu.Lock()
	delete(g.m, key)
	g.mu.Unlock()

	return c.val, c.err
}
在 Do 函数的源码中，g.m 的读写被 g.mu 互斥锁保护，fn 的返回结果存储在 call.val、call.err 中，通过 sync.WaitGroup 实现等待 fn 执行结束。

回到本文开头提到的问题，对于缓存的更新，可以这样实现：

if (cacheMiss) {
    fn = func() (interface{}, error) {
        // 缓存更新逻辑
    }
    data, err = g.Do(cacheKey, fn)
}
对于防止暴增的接口请求对下游服务造成瞬时高负载，可以这样实现：

fn = func() (interface{}, error) {
    // 发送请求到其他服务接口
}
data, err = g.Do(apiNameWithParams, fn)

缓存击穿
    给缓存加一个过期时间，下次未命中缓存时再去从数据源获取结果写入新的缓存，这个是后端开发人员再熟悉不过的基操。本人之前在做直播平台活动业务的时候，当时带着这份再熟练不过的自信，把复杂的数据库链表语句写好，各种微服务之间调用捞数据最后算好的结果，丢进了缓存然后设了一个过期时间，当时噼里啪啦两下写完代码觉得稳如铁蛋，结果在活动快结束之前，数据库很友好的挂掉了。当时回去查看监控后发现，是在活动快结束前，大量用户都在疯狂的刷活动页，导致缓存过期的瞬间有大量未命中缓存的请求直接打到数据库上所导致的，所以这个经典的问题稍不注意还是害死人

    防缓存击穿的方式有很多种，比如通过计划任务来跟新缓存使得从前端过来的所有请求都是从缓存读取等等。之前读过 groupCache的源码，发现里面有一个很有意思的库，叫singleFlight, 因为groupCache从节点上获取缓存如果未命中，则会去其他节点寻找，其他节点还没有的话再从数据源获取，所以这个步骤对于防击穿非常有必要。singleFlight使得groupCache在多个并发请求对一个失效的key进行源数据获取时，只让其中一个得到执行，其余阻塞等待到执行的那个请求完成后，将结果传递给阻塞的其他请求达到防止击穿的效果。
    
    https://mp.weixin.qq.com/s/lSBIPbMXcjoWgrqMVoDanA
<!-- more -->
如果你曾经使用过 Go 一段时间，那么你可能了解一些 Go 中的并发原语：

go 关键字用来生成 goroutines ；channel 用于 goroutines 之间通信 ；context 用于传播取消 ；sync 和 sync/atomic 包用于低级别的原语，例如互斥锁和内存的原子操作 。这些语言特性和包组合在一起，为构建高并发的应用程序提供了丰富的工具集。你可能还没有发现在扩展库 golang.org/x/sync 中，提供了一系列更高级别的并发原语。我们将在本文中来谈谈这些内容。

Singleflight包
正如文档中所描述，这个包提供了一个重复函数调用抑制的机制。

如果你正在处理计算量大（也可能仅仅是慢，比如网络访问）的用户请求时，这个包就很有用。例如，你的数据库中包含每个城市的天气信息，并且你想将这些数据以 API 的形式提供服务。在某些情况下，可能同时有多个用户想查询同一城市的天气。

在这种场景下，如果你只查询一次数据库然后将结果共享给所有等待的请求，这样不是更好吗？这就是 singleflight 提供的功能。

在使用的时候，我们需要要创建一个 singleflight.Group。它需要在所有请求中共享才能工作。然后将缓慢或者开销大的操作包装到 group.Do(key, fn) 的调用中。对同一个 key 的多个并发请求将仅调用 fn 一次，并且将 fn 的结果返回给所有调用者。

实际中的使用如下:

package weather
type Info struct {
    TempC, TempF int // temperature in Celsius and Farenheit
    Conditions string // "sunny", "snowing", etc
}

var group singleflight.Group

func City(city string) (*Info, error) {
    results, err, _ := group.Do(city, func() (interface{}, error) {
        info, err := fetchWeatherFromDB(city) // 慢操作
        return info, err
    })
    if err != nil {
        return nil, fmt.Errorf("weather.City %s: %w", city, err)
    }
    return results.(*Info), nil
}
需要注意的是，我们传递给 group.Do 的闭包必须返回 (interface{}, error) 才能和 Go 类型系统一起使用。上面的例子中忽略了 group.Do 的第三个返回值，该值是用来表示结果是否在多个调用方之间共享。

如果需要查看更多完整的例子，可以查看 Encore Playground 中的代码。

errgroup 包
另一个有用的包是 errgroup package。它和 sync.WaitGroup 比较相似，但是会将任务返回的错误回传给阻塞的调用方。

当你有多个等待的操作，但又想知道它们是否都已经成功完成时，这个包就很有用。还是以上面的天气为例，假如你要一次查询多个城市的天气，并且要确保其中所有的查询都成功返回。

首先定义一个 errgroup.Group，然后为每个城市都使用 group.Go(fn func() error) 方法。该方法会生成一个 goroutine 来执行这个任务。当生成你想执行的所有任务时，使用 group.Wait() 等待它们完成。需要注意和 sync.WaitGroup 有一点不同的是，该方法会返回错误。当且仅当所有任务都返回 nil 时，才会返回一个 nil 错误。

实际中的使用如下:

func Cities(cities ...string) ([]*Info, error) {
    var g errgroup.Group
    var mu sync.Mutex
    res := make([]*Info, len(cities)) // res[i] corresponds to cities[i]

    for i, city := range cities {
        i, city := i, city // 为下面的闭包创建局部变量
        g.Go(func() error {
            info, err := City(city)
            mu.Lock()
            res[i] = info
            mu.Unlock()
            return err
        })
    }
    if err := g.Wait(); err != nil {
        return nil, err
    }
    return res, nil
}
这里我们使用一个 res 切片来存储每个 goroutine 执行的结果。尽管上面的代码没有使用 mu 互斥锁也是线程安全的，但是每个 goroutine 都是在切片中自己的位置写入结果，因此我们不得不使用一个切片，以防代码变化。

限制并发
上面的代码将同时查找给定城市的天气信息。如果城市数量比较少，那还不错，但是如果城市数量很多，可能会导致性能问题。在这种情况下，就应该引入限制并发了。

在 Go 中使用 semaphores 信号量让实现限制并发变得非常简单。信号量是你学习计算机科学中可能已经遇到过的并发原语，如果没有遇到也不用担心。你可以出于多种目的来使用信号量，但是这里我们只使用它来追踪运行中的任务的数量，并阻塞直到有空间可以执行其他任务。

在 Go 中，我们可以使用 channel 来实现信号量的功能。如果我们一次需要最多执行 10 个任务，则需要创建一个容量为 10 的 channel：semaphore := make(chan struct{}, 10)。你可以想象它为一个可以容纳 10 个球的管道。

如果想执行一个新的任务，我们只需要给 channel 发送一个值：semaphore <- struct{}{}，如果已经有很多任务在运行的话，将会阻塞。这类似于将一个球推入管道，如果管道已满，则需要等待直到有空间为止。

当通过 <-semaphore 能从该 channel 中取出一个值时，这表示一个任务完成了。这类似于在管道另一端拿出一个球，这将为塞入下一个球提供了空间。

如描述一样，我们修改后的 Cities 代码如下：

func Cities(cities ...string) ([]*Info, error) {
    var g errgroup.Group
    var mu sync.Mutex
    res := make([]*Info, len(cities)) // res[i] corresponds to cities[i]
    sem := make(chan struct{}, 10)
    for i, city := range cities {
        i, city := i, city // create locals for closure below
        sem <- struct{}{}
        g.Go(func() error {
            info, err := City(city)
            mu.Lock()
            res[i] = info
            mu.Unlock()
            <-sem
            return err
        })
    }
    if err := g.Wait(); err != nil {
        return nil, err
    }
    return res, nil
}
加权限制并发
最后，当你想要限制并发的时候，并不是所有任务优先级都一样。在这种情况下，我们消耗的资源将依据高、低优先级任务的分布以及它们如何开始运行而发生变化。

在这种场景下使用加权限制并发是一种不错的解决方式。它的工作原理很简单：我们不需要为同时运行的任务数量做预估，而是为每个任务提供一个 "cost"，并从信号量中获取和释放它。

我们不再使用 channel 来做这件事，因为我们需要立即获取并释放 "cost"。幸运的是，"扩展库" golang.org/x/sync/sempahore 实现了加权信号量。

sem <- struct{}{} 操作叫 "获取"，<-sem 操作叫 "释放"。你可能会注意到 semaphore.Acquire 方法会返回错误，那是因为它可以和 context 包一起使用来控制提前结束。在这个例子中，我们将忽略它。

实际上，天气查询的例子比较简单，不适用加权信号量，但是为了简单起见，我们假设 cost 变量随城市名称长度而变化。然后，我们修改如下：

func Cities(cities ...string) ([]*Info, error) {
    ctx := context.TODO() // 需要的时候，可以用 context 替换 
    var g errgroup.Group
    var mu sync.Mutex
    res := make([]*Info, len(cities)) // res[i] 对应 cities[i]
    sem := semaphore.NewWeighted(100) // 并发处理 100 个字符
    for i, city := range cities {
        i, city := i, city // 为闭包创建局部变量
        cost := int64(len(city))
        if err := sem.Acquire(ctx, cost); err != nil {
            break
        }
        g.Go(func() error {
            info, err := City(city)
            mu.Lock()
            res[i] = info
            mu.Unlock()
            sem.Release(cost)
            return err
        })
    }
    if err := g.Wait(); err != nil {
        return nil, err
    } else if err := ctx.Err(); err != nil {
        return nil, err
    }
    return res, nil
}
socketFunc创建了 socket，通知将 socket 设置非阻塞（SOCK_NONBLOCK）以及 fork 时关闭（SOCK_CLOEXEC），这两个标志是在 linux 内核版本 2.6.27 之后添加，在此之前的版本代码将会走到syscall.ForkLock.RLock()，主要是为了防止在 fork 时导致文件描述符泄露。

当 socket 创建之后进入新建 fd 流程，在 Go 的包装层面，fd 均以netFD结构表示，该接口描述原始 socket 的地址信息、协议类型、协议族以及 option，netFD在整个包装结构中居于用户接口的下一层。最后进入监听逻辑，逻辑走向区分 TCP 和 UDP，而监听逻辑比较简单，即调用系统 bind 和 listen 接口 (net/sock_posix.go)：

func (fd *netFD) listenStream(laddr sockaddr, backlog int, ctrlFn func(string, string, syscall.RawConn) error) error {
    // ...

    if ctrlFn != nil {
        c, err := newRawConn(fd)
        if err != nil {
            return err
        }
        if err := ctrlFn(fd.ctrlNetwork(), laddr.String(), c); err != nil {
            return err
        }
    }
    if err = syscall.Bind(fd.pfd.Sysfd, lsa); err != nil {
        return os.NewSyscallError("bind", err)
    }
    if err = listenFunc(fd.pfd.Sysfd, backlog); err != nil {
        return os.NewSyscallError("listen", err)
    }
    if err = fd.init(); err != nil {
        return err
    }
    lsa, _ = syscall.Getsockname(fd.pfd.Sysfd)
    fd.setAddr(fd.addrFunc()(lsa), nil)
    return nil
}
结论
上面的例子展示了在 Go 中通过微调来实现需要的并发模式是多么简单。
