---
title: goroutine 泄露
layout: post
category: golang
author: 夏泽民
---
如果有goroutine因为ch没有接收动作而被一直阻塞在发送处，无法被回收。
发现goroutine 泄露工具：https://github.com/uber-go/goleak
产生原因分析
产生goroutine leak（协程泄露）的原因可能有以下几种： * goroutine由于channel的读/写端退出而一直阻塞，导致goroutine一直占用资源，而无法退出 * goroutine进入死循环中，导致资源一直无法释放

goroutine终止的场景
一个goroutine终止有以下几种情况： * 当一个goroutine完成它的工作 * 由于发生了没有处理的错误 * 有其他的协程告诉它终止

如何调试和发现goroutine leak
runtime
可以通过runtime.NumGoroutine()函数来获取后台服务的协程数量。通过查看每次的协程数量的变化和增减，我们可以判断是否有goroutine泄露发生。

...
	fmt.Fprintf(os.Stderr, "%d\n", runtime.NumGoroutine())
	time.Sleep(10e9) //等一会，查看协程数量的变化
	fmt.Fprintf(os.Stderr, "%d\n", runtime.NumGoroutine())
...
pprof来确认泄露的地方
一旦我们发现了goroutein leak，我们就需要确认泄露的出处。

import (
  "runtime/debug"
  "runtime/pprof"
)

func getStackTraceHandler(w http.ResponseWriter, r *http.Request) {
    stack := debug.Stack()
    w.Write(stack)
    pprof.Lookup("goroutine").WriteTo(w, 2)
}
func main() {
    http.HandleFunc("/_stack", getStackTraceHandler)
}

http://localhost:11181/，我们就可以得到整个goroutine的信息，仅列出关键信息：
goroutine profile: total 10004
 
10000 @ 0x186f6 0x616b 0x6298 0x2033 0x188c0
#   0x2033  main.f+0x33 /Users/siddontang/test/pprof.go:11

<!-- more -->
实际的goroutine leak
生产者消费者场景
func main() {
	newRandStream := func() <-chan int {
		randStream := make(chan int)

		go func() {
			defer fmt.Println("newRandStream closure exited.")
			defer close(randStream)
			// 死循环：不断向channel中放数据，直到阻塞
			for {
				randStream <- rand.Int()
			}
		}()

		return randStream
	}

	randStream := newRandStream()
	fmt.Println("3 random ints:")

	// 只消耗3个数据，然后去做其他的事情，此时生产者阻塞，
	// 若主goroutine不处理生产者goroutine，则就产生了泄露
	for i := 1; i <= 3; i++ {
		fmt.Printf("%d: %d\n", i, <-randStream)
	}

	fmt.Fprintf(os.Stderr, "%d\n", runtime.NumGoroutine())
	time.Sleep(10e9)
	fmt.Fprintf(os.Stderr, "%d\n", runtime.NumGoroutine())
}
生产协程进入死循环，不断产生数据，消费协程，也就是主协程只消费其中的3个值，然后主协程就再也不消费channel中的数据，去做其他的事情了。此时生产协程放了一个数据到channel中，但已经不会有协程消费该数据了，所以，生产协程阻塞。 此时，若没有人再消费channel中的数据，生成协程是被泄露的协程。

解决办法
总的来说，要解决channel引起的goroutine leak问题，主要是看在channel阻塞goroutine时，该goroutine的阻塞是正常的，还是可能导致协程永远没有机会执行。若可能导致协程永远没有机会执行，则可能会导致协程泄露。 所以，在创建协程时就要考虑到它该如何终止。

解决一般问题的办法就是，当主线程结束时，告知生产线程，生产线程得到通知后，进行清理工作：或退出，或做一些清理环境的工作。

func main() {
	newRandStream := func(done <-chan interface{}) <-chan int {
		randStream := make(chan int)

		go func() {
			defer fmt.Println("newRandStream closure exited.")
			defer close(randStream)

			for {
				select {
				case randStream <- rand.Int():
				case <-done:  // 得到通知，结束自己
					return
				}
			}
		}()

		return randStream
	}


	done := make(chan interface{})
	randStream := newRandStream(done)
	fmt.Println("3 random ints:")

	for i := 1; i <= 3; i++ {
		fmt.Printf("%d: %d\n", i, <-randStream)
	}

    // 通知子协程结束自己
    // done <- struct{}{}
	close(done)
	// Simulate ongoing work
	time.Sleep(1 * time.Second)
}
上面的代码中，协程通过一个channel来得到结束的通知，这样它就可以清理现场。防止协程泄露。 通知协程结束的方式，可以是发送一个空的struct，更加简单的方式是直接close channel。如上图所示。

master work场景
在该场景下，我们一般是把工作划分成多个子工作，把每个子工作交给每个goroutine来完成。此时若处理不当，也是有可能发生goroutine泄漏的。我们来看一下实际的例子：

代码
// function to add an array of numbers.
func worker_adder(s []int, c chan int) {
	sum := 0
	for _, v := range s {
		sum += v
	}
	// writes the sum to the go routines.
	c <- sum // send sum to c
	fmt.Println("end")
}

func main() {
	s := []int{7, 2, 8, -9, 4, 0}

	c1 := make(chan int)
	c2 := make(chan int)

	// spin up a goroutine.
	go worker_adder(s[:len(s)/2], c1)
	// spin up a goroutine.
	go worker_adder(s[len(s)/2:], c2)

	//x, y := <-c1, <-c2 // receive from c1 aND C2
	x, _:= <-c1
	// 输出从channel获取到的值
	fmt.Println(x)

	fmt.Println(runtime.NumGoroutine())
	time.Sleep(10e9)
	fmt.Println(runtime.NumGoroutine())
}
以上代码在主协程中，把一个数组分成两个部分，分别交给两个worker协程来计算其值，这两个协程通过channel把结果传回给主协程。 但，在以上代码中，我们只接收了一个channel的数据，导致另一个协程在写channel时阻塞，再也没有执行的机会。 要是我们把这段代码放入一个常驻服务中，看的更加明显：

http server 场景
代码
// 把数组s中的数字加起来
func sumInt(s []int, c chan int) {
	sum := 0
	for _, v := range s {
		sum += v
	}
	c <- sum
}

// HTTP handler for /sum
func sumConcurrent2(w http.ResponseWriter, r *http.Request) {
	s := []int{7, 2, 8, -9, 4, 0}

	c1 := make(chan int)
	c2 := make(chan int)

	go sumInt(s[:len(s)/2], c1)
	go sumInt(s[len(s)/2:], c2)

	// 这里故意不在c2中读取数据，导致向c2写数据的协程阻塞。
	x := <-c1

	// write the response.
	fmt.Fprintf(w, strconv.Itoa(x))
}

func main() {
	StasticGroutine := func() {
		for {
			time.Sleep(1e9)
			total := runtime.NumGoroutine()
			fmt.Println(total)
		}
	}

	go StasticGroutine()

	http.HandleFunc("/sum", sumConcurrent2)
	err := http.ListenAndServe(":8001", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
这个输出是我们的http server的协程数量，可以看到：每请求一次，协程数就增加一个，而且不会减少。说明已经发生了协程泄露(goroutine leak)。

解决办法
解决的办法就是不管在任何情况下，都必须要有协程能够读写channel，让协程不会阻塞。 代码修改如下：

...
	x,y := <-c1,<-c2

	// write the response.
	fmt.Fprintf(w, strconv.Itoa(x+y))
...


func main() {
	ch := make(chan int)
	go func(i int) {
		ch <- 1
		fmt.Println("send ", i)
	}(0)
	go func(i int) {
		ch <- 1
		fmt.Println("send ", i)
	}(1)
    <-ch
    fmt.Println("exit!")
}
会泄露
func main() {
    ch := make(chan int)
    for i := 0; i < 3; i++ {
        go func(i int) {
            ch <- 1
            fmt.Println("send ", i)
        }(i)
    }
    
    <-ch
    fmt.Println("exit!")
}
不会泄露

PouchContainer 是阿里巴巴集团开源的一款容器运行时产品，它具备强隔离和可移植性等特点，可用来帮助企业快速实现存量业务容器化，以及提高企业内部物理资源的利用率。

PouchContainer 同时还是一款 golang 项目。在此项目中，大量运用了 goroutine 来实现容器管理、镜像管理和日志管理等模块。goroutine 是 golang 在语言层面就支持的用户态 “线程”，这种原生支持并发的特性能够帮助开发者快速构建高并发的服务。

虽然 goroutine 容易完成并发或者并行的操作，但如果出现 channel 接收端长时间阻塞却无法唤醒的状态，那么将会出现 goroutine leak 。 goroutine leak 同内存泄漏一样可怕，这样的 goroutine 会不断地吞噬资源，导致系统运行变慢，甚至是崩溃。为了让系统能健康运转，需要开发者保证 goroutine 不会出现泄漏的情况。 接下来本文将从什么是 goroutine leak, 如何检测以及常用的分析工具来介绍 PouchContainer 在 goroutine leak 方面的检测实践。

2.2 如何检测 goroutine leak？
对于 HTTP Server 而言，我们通常会通过引入包 net/http/pprof 来查看当前进程运行的状态，其中有一项就是查看 goroutine stack 的信息，{ip}:{port}/debug/pprof/goroutine?debug=2 。我们来看看调用者主动断开链接之后的 goroutine stack 信息。

# step 1: create background jobpouch run -d busybox sh -c "while true; do sleep 1; done"# step 2: follow the log and stop it after 3 secondscurl -m 3 {ip}:{port}/v1.24/containers/{container_id}/logs?stdout=1&follow=1# step 3: after 3 seconds, dump the stack infocurl -s "{ip}:{port}/debug/pprof/goroutine?debug=2" | grep -A 10 logsContainergithub.com/alibaba/pouch/apis/server.(*Server).logsContainer(0xc420330b80, 0x251b3e0, 0xc420d93240, 0x251a1e0, 0xc420432c40, 0xc4203f7a00, 0x3, 0x3)        /tmp/pouchbuild/src/github.com/alibaba/pouch/apis/server/container_bridge.go:339 +0x347github.com/alibaba/pouch/apis/server.(*Server).(github.com/alibaba/pouch/apis/server.logsContainer)-fm(0x251b3e0, 0xc420d93240, 0x251a1e0, 0xc420432c40, 0xc4203f7a00, 0x3, 0x3)        /tmp/pouchbuild/src/github.com/alibaba/pouch/apis/server/router.go:53 +0x5cgithub.com/alibaba/pouch/apis/server.withCancelHandler.func1(0x251b3e0, 0xc420d93240, 0x251a1e0, 0xc420432c40, 0xc4203f7a00, 0xc4203f7a00, 0xc42091dad0)        /tmp/pouchbuild/src/github.com/alibaba/pouch/apis/server/router.go:114 +0x57github.com/alibaba/pouch/apis/server.filter.func1(0x251a1e0, 0xc420432c40, 0xc4203f7a00)        /tmp/pouchbuild/src/github.com/alibaba/pouch/apis/server/router.go:181 +0x327net/http.HandlerFunc.ServeHTTP(0xc420a84090, 0x251a1e0, 0xc420432c40, 0xc4203f7a00)        /usr/local/go/src/net/http/server.go:1918 +0x44github.com/alibaba/pouch/vendor/github.com/gorilla/mux.(*Router).ServeHTTP(0xc4209fad20, 0x251a1e0, 0xc420432c40, 0xc4203f7a00)        /tmp/pouchbuild/src/github.com/alibaba/pouch/vendor/github.com/gorilla/mux/mux.go:133 +0xednet/http.serverHandler.ServeHTTP(0xc420a18d00, 0x251a1e0, 0xc420432c40, 0xc4203f7800)
我们会发现当前进程中还存留着 logsContainer goroutine。因为这个容器没有输出任何日志的机会，所以这个 goroutine 没办法通过 write: broken pipe 的错误退出，它会一直占用着系统资源。那我们该怎么解决这个问题呢？

2.3 怎么解决？
golang 提供的包 net/http 有监控链接断开的功能:

// HTTP Handler Interceptorsfunc withCancelHandler(h handler) handler {        return func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {                // https://golang.org/pkg/net/http/#CloseNotifier                if notifier, ok := rw.(http.CloseNotifier); ok {                        var cancel context.CancelFunc                        ctx, cancel = context.WithCancel(ctx)                        waitCh := make(chan struct{})                        defer close(waitCh)                        closeNotify := notifier.CloseNotify()                        go func() {                                select {                                case <-closeNotify:                                        cancel()                                case <-waitCh:                                }                        }()                }                return h(ctx, rw, req)        }}
当请求还没执行完毕时，客户端主动退出了，那么 CloseNotify() 将会收到相应的消息，并通过  context.Context 来取消，这样我们就可以很好地处理 goroutine leak 的问题了。在 golang 的世界里，你会经常看到 读_ 和 _写 的 goroutine，它们这种函数的第一个参数一般会带有 context.Context , 这样就可以通过 WithTimeout 和  WithCancel 来控制 goroutine 的回收，避免出现泄漏的情况。

3. 常用的分析工具
3.1 net/http/pprof
在开发 HTTP Server 的时候，我们可以引入包 net/http/pprof 来打开 debug 模式，然后通过 /debug/pprof/goroutine 来访问 goroutine stack 信息。一般情况下，goroutine stack 会具有以下样式。

goroutine 93 [chan receive]:github.com/alibaba/pouch/daemon/mgr.NewContainerMonitor.func1(0xc4202ce618)        /tmp/pouchbuild/src/github.com/alibaba/pouch/daemon/mgr/container_monitor.go:62 +0x45created by github.com/alibaba/pouch/daemon/mgr.NewContainerMonitor        /tmp/pouchbuild/src/github.com/alibaba/pouch/daemon/mgr/container_monitor.go:60 +0x8dgoroutine 94 [chan receive]:github.com/alibaba/pouch/daemon/mgr.(*ContainerManager).execProcessGC(0xc42037e090)        /tmp/pouchbuild/src/github.com/alibaba/pouch/daemon/mgr/container.go:2177 +0x1a5created by github.com/alibaba/pouch/daemon/mgr.NewContainerManager        /tmp/pouchbuild/src/github.com/alibaba/pouch/daemon/mgr/container.go:179 +0x50b
goroutine stack 通常第一行包含着 Goroutine ID，接下来的几行是具体的调用栈信息。有了调用栈信息，我们就可以通过 关键字匹配 的方式来检索是否存在泄漏的情况了。

在 Pouch 的集成测试里，Pouch Logs API 对包含 (*Server).logsContainer 的 goroutine stack 比较感兴趣。因此在测试跟随模式完毕后，会调用 debug 接口检查是否包含  (*Server).logsContainer 的调用栈。一旦发现包含便说明该 goroutine 还没有被回收，存在泄漏的风险。

总的来说，debug 接口的方式适用于  集成测试 ，因为测试用例和目标服务不在同一个进程里，需要 dump 目标进程的 goroutine stack 来获取泄漏信息。

3.2 runtime.NumGoroutine
当测试用例和目标函数／服务在同一个进程里时，可以通过 goroutine 的数目变化来判断是否存在泄漏问题。

func TestXXX(t *testing.T) {    orgNum := runtime.NumGoroutine()    defer func() {        if got := runtime.NumGoroutine(); orgNum != got {            t.Fatalf("xxx", orgNum, got)        }    }()    ...}
3.3 github.com/google/gops
gops 与包 net/http/pprof 相似，它是在你的进程内放入了一个 agent ，并提供命令行接口来查看进程运行的状态，其中 gops stack ${PID} 可以查看当前 goroutine stack 状态
