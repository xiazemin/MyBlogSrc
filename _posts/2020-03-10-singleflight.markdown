---
title: Go防缓存击穿
layout: post
category: golang
author: 夏泽民
---
我们在开发时，有时会碰到一个接口的访问量突然上升，导致服务响应延迟或者宕机的情况。这时，除了利用缓存之外，也可以用到singlefilght来解决
package main

import (
    "fmt"
    "sync"
    "sync/atomic"
    "time"

    "golang.org/x/sync/singleflight"
)

func main() {
    g := singleflight.Group{}

    wg := sync.WaitGroup{}

    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func(j int) {
            defer wg.Done()
            val, err, shared := g.Do("a", a)
            if err != nil {
                fmt.Println(err)
                return
            }
            fmt.Printf("index: %d, val: %d, shared: %v\n", j, val, shared)
        }(i)
    }

    wg.Wait()

}

var (
    count = int64(0)
)

// 模拟接口方法
func a() (interface{}, error) {
    time.Sleep(time.Millisecond * 500)
    return atomic.AddInt64(&count, 1), nil
}

// 部分输出，shared表示是否共享了其他请求的返回结果
index: 2, val: 1, shared: false
index: 71, val: 1, shared: true
index: 69, val: 1, shared: true
index: 73, val: 1, shared: true
index: 8, val: 1, shared: true
index: 24, val: 1, shared: true
<!-- more -->
val这里绝大部分都为1是因为程序运行时间太快了，可以试着把time.Sleep时间缩短一点看看效果

singleflight核心代码非常简单

// 成员非常少，就两个
type Group struct {
    mu sync.Mutex      
    m  map[string]*call
}

func (g *Group) Do(key string, fn func() (interface{}, error)) (v interface{}, err error, shared bool) {
    g.mu.Lock()
    if g.m == nil {
        g.m = make(map[string]*call)
    }
    // 判断key是否存在，存在则表示有其他请求先一步进来，
    // 直接等待其他请求返回就行
    if c, ok := g.m[key]; ok {
        c.dups++
        g.mu.Unlock()
        c.wg.Wait()
        return c.val, c.err, true
    }
    // 不存在就创建一个新的call对象，然后去执行
    c := new(call)
    c.wg.Add(1)
    g.m[key] = c
    g.mu.Unlock()

    g.doCall(c, key, fn)
    return c.val, c.err, c.dups > 0
}
doCall方法很简单，这里就不展开了，除了Do方法之外，还有一个异步的DoChan方法，原理一模一样。

我们一般可以在一些类似于幂等的接口上用singleflight


singleflight.go文件中是singleflight模块的代码，这主要是进行相同访问的一个合并操作。也就是说，如果对于某个key的请求已经存在并且正在进行，则对该key的新的请求会堵塞在这里，等原来的请求结束后，将请求得到的结果同时返回给堵塞中的请求。

该部分就封装了一个接口：

func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error)
1
首先，先定义了下面两个结构体：

//实际请求函数的封装结构体
// call is an in-flight or completed Do call
type call struct {
	wg  sync.WaitGroup
	//实际的请求函数
	val interface{}
	err error
}

//主要是用来组织已经存在的对某key的请求和对应的实际请求函数映射
// Group represents a class of work and forms a namespace in which
// units of work can be executed with duplicate suppression.
type Group struct {
	//用于对m上锁，保护m
	mu sync.Mutex       // protects m
	m  map[string]*call // lazily initialized
}

下面具体讲解一下这个函数。

该函数入参是一个key和一个实际请求函数，出参是一个接口类型和一个错误类型。

这里利用了go的锁机制，比如Metux、WaitGroup等。

// Do executes and returns the results of the given function, making
// sure that only one execution is in-flight for a given key at a
// time. If a duplicate comes in, the duplicate caller waits for the
// original to complete and receives the same results.
func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	//有可能要修改m，所以先上锁进行保护
	g.mu.Lock()
	//如果m为nil，则初始化一个
	if g.m == nil {
		g.m = make(map[string]*call)
	}
	//如果m中存在对该key的请求，则该线程不会直接再次访问key，所以释放锁
	//然后堵塞等待已经存在的请求得到的结果
	if c, ok := g.m[key]; ok {
		//解锁
		g.mu.Unlock()
		//堵塞
		c.wg.Wait()
		//如果已经存在的请求完成，则堵塞状态会解除，继续向下执行，得到正确结果
		return c.val, c.err
	}
	//如果不存在对该key的请求，则本线程要进行实际的请求，保持m的锁定状态
	//创建一个实际请求结构体
	c := new(call)
	//为了保证其他的相同请求的堵塞
	c.wg.Add(1)
	//组织好映射关系
	g.m[key] = c
	//解锁m
	g.mu.Unlock()
	
	//执行真正的请求函数，得到对该key请求的结果
	c.val, c.err = fn()
	//得到结果后取消其他请求的堵塞
	c.wg.Done()

	//该次请求完成后，要从已存在请求map中删掉
	g.mu.Lock()
	delete(g.m, key)
	g.mu.Unlock()
	
	//返回请求结果
	return c.val, c.err
}

package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/golang/groupcache/singleflight"
)

func NewDelayReturn(dur time.Duration, n int) func() (interface{}, error) {
	return func() (interface{}, error) {
		time.Sleep(dur)
		return n, nil
	}
}

func main() {
	g := singleflight.Group{}
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		ret, err := g.Do("key", NewDelayReturn(time.Second*1, 1))
		if err != nil {
			panic(err)
		}
		fmt.Printf("key-1 get %v\n", ret)
		wg.Done()
	}()
	go func() {
		time.Sleep(100 * time.Millisecond) // make sure this is call is later
		ret, err := g.Do("key", NewDelayReturn(time.Second*2, 2))
		if err != nil {
			panic(err)
		}
		fmt.Printf("key-2 get %v\n", ret)
		wg.Done()
	}()
	wg.Wait()
}
执行结果(耗时： 1.019s)

key-2 get 1
key-1 get 1

背景
先来了解一下缓存问题的几种场景，以redis为例

缓存穿透
缓存穿透，是指查询一个数据库一定不存在的数据。正常的使用缓存流程大致是，数据查询先进行缓存查询，如果key不存在或者key已经过期，再对数据库进行查询，并把查询到的对象，放进缓存。如果数据库查询对象为空，则不放进缓存。

代码流程
1.参数传入对象主键ID
2.根据key从缓存中获取对象
3.如果对象不为空，直接返回
4.如果对象为空，进行数据库查询
5.如果从数据库查询出的对象不为空，则放入缓存（设定过期时间）
想象一下这个情况，如果传入的参数为-1，会是怎么样？这个-1，就是一定不存在的对象。就会每次都去查询数据库，而每次查询都是空，每次又都不会进行缓存。假如有恶意攻击，就可以利用这个漏洞，对数据库造成压力，甚至压垮数据库。即便是采用UUID，也是很容易找到一个不存在的KEY，进行攻击。

小编在工作中，会采用缓存空值的方式，也就是【代码流程】中第5步，如果从数据库查询的对象为空，也放入缓存，只是设定的缓存过期时间较短，比如设置为60秒。

缓存雪崩
缓存雪崩，是指在某一个时间段，缓存集中过期失效。

产生雪崩的原因之一，比如在写本文的时候，马上就要到双十二零点，很快就会迎来一波抢购，这波商品时间比较集中的放入了缓存，假设缓存一个小时。那么到了凌晨一点钟的时候，这批商品的缓存就都过期了。而对这批商品的访问查询，都落到了数据库上，对于数据库而言，就会产生周期性的压力波峰。

比如电商项目，一般是采取不同分类商品，缓存不同周期。在同一分类中的商品，加上一个随机因子。这样能尽可能分散缓存过期时间，而且，热门类目的商品缓存时间长一些，冷门类目的商品缓存时间短一些，也能节省缓存服务的资源。

其实集中过期，倒不是非常致命，比较致命的缓存雪崩，是缓存服务器某个节点宕机或断网。因为自然形成的缓存雪崩，一定是在某个时间段集中创建缓存，那么那个时候数据库能顶住压力，这个时候，数据库也是可以顶住压力的。无非就是对数据库产生周期性的压力而已。而缓存服务节点的宕机，对数据库服务器造成的压力是不可预知的，很有可能瞬间就把数据库压垮。

缓存击穿
缓存击穿，是指一个key非常热点，在不停的扛着大并发，大并发集中对这一个点进行访问，当这个key在失效的瞬间，持续的大并发就穿破缓存，直接请求数据库，就像在一个屏障上凿开了一个洞。相当于电商项目中的“爆款”。

其实，大多数情况下这种爆款很难对数据库服务器造成压垮性的压力。达到这个级别的公司没有几家的。如果流量不大，务实一点就是直接设置缓存永不过期，即便某些商品自己发酵成了爆款，也是直接设为永不过期就好了。

但是万一流量很大，遇到这种缓存击穿怎么办，这个时候推荐一款golang的包singleflight

原理：
多个并发请求对一个失效的key进行源数据获取时，只让其中一个得到执行，其余阻塞等待到执行的那个请求完成后，将结果传递给阻塞的其他请求达到防止击穿的效果。

demo：
模拟一百个并发请求在缓存失效的瞬间同时调用rpc访问源数据

package backup

import (
    "fmt"
    "github.com/golang/groupcache/singleflight"
    "net/http"
    "net/rpc"
    "sync"
    "time"
)

type (
    Arg struct {
        Caller int
    }
    Data struct{}
)

// 模拟从数据源获取数据
func (d *Data) GetData(arg *Arg, replay *string) error {
    fmt.Printf("request from client %d\n", arg.Caller)
    time.Sleep(1 * time.Second)
    *replay = "source data from rpcServer"
    return nil
}

func main() {
    d := new(Data)
    rpc.Register(d)
    rpc.HandleHTTP()
    fmt.Println("start rpc server")
    if err := http.ListenAndServe(":8976", nil); err != nil {
        panic(err)
    }

    client, err := rpc.DialHTTP("tcp", ":8976")
    if (err != nil) {
        panic(err)
    }

    singleFlight := new(singleflight.Group)
    wg := sync.WaitGroup{}
    wg.Add(100)

    //然后再模拟一百个并发请求在缓存失效的瞬间同时调用rpc访问源数据
    for i := 0; i < 100; i++ {
        fn := func() (interface{}, error) {
            var replay string
            err = client.Call("Data.GetData", Arg{Caller: i}, &replay)
            //从数据源拿到数据后再更新缓存等
            //...
            //...
            return replay, err
        }

        go func(i int) {
            result, _ := singleFlight.Do("foo", fn)
            fmt.Printf("caller %d get result '%s'\n", i, result)
            wg.Done()
        }(i)
    }

}

效果：

可以看到100个并发请求从源数据获取时，rpcServer端只收到了来自client 17的请求，而其余99个最后也都得到了正确的返回值。

其实singleflight就几十行代码,主要就用到sync包的两个特性Mutex和WaitGroup，所以说看优秀开源代码是学习最快的方式，没有之一。

package singleflight

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



