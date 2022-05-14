---
title: close_wait
layout: post
category: golang
author: 夏泽民
---
1 背景
线上有一个高并发的 HTTP Go 服务部署在 A 区域，能够正常提供服务。我们有天将 B 区域的流量也准备切到这个服务的时候，发生了一个很诡异的事情。从 A 区域来的流量全部 200，但是从 B 区域来的流量全部都 502、504。

2 排查过程
2.1 怀疑网络问题
这种从不同区域一部分流量正常，一部分流量异常，第一直觉觉得是异常流量那块应该是网络问题。所以我们让运维去异常流量 nginx 机器上，发 telnet 和 curl 到服务上。
https://mp.weixin.qq.com/s/sqBdycaClUixZQPgKy52Pw
<!-- more -->
可以看到 telnet 这个 ip port 端口是通的，说明四层以下网络是通的。但是当我们 curl 到对应服务的时候，发现服务被连接重置。这时候我们就献出我们的重启大法，对该实例进行重启。发现该实例的 curl 请求恢复。所以就可以直接排除我们的网络问题。这个时候，我们就要思考了，是啥问题让我们程序出现这种诡异问题。因为我们还留了一台问题机器。我们开始着手对他进行分析。

2.2 分析问题机器
2.2.1 查看网络问题
我们在出问题的机器上，查看网络情况。我们发现有大量的 CLOSE WAIT 现象。







然后我们查看下 CLOSE WAIT 详细信息。







可以看到上面有以下情况

第二列是 RECVQ，可以看到有很多都大于 1

Close Wait，来自于探活的客户端和 nginx 客户端

Close Wait 没有 pid 号（我这里有 root 权限，不要怀疑我是没权限看不到 pid 信息）

第一个说明了 RECVQ 大于 1 的都是 nginx 客户端，1 的都是探活的客户端。这一块说明了 nginx 是被堵住了。

第二个说明了里面探活的客户端有 8000 条信息，nginx 的 close wait 大概 100 多条

第三个说明了没有 pid 问题，我们待会在后面分析中会提及。

看到这么多 Close Wait。我们首先要拿出这个四次挥手图。







这里要注意，客户端和服务端是个相对概念。在我们这个问题里，我们的客户端是探活服务，服务端是我们的业务服务。通过这个图，可以知道我们的客户端是发了 close 包，我们的服务端可能没有正常响应 ack 或者 fin 包，导致 server 端出现 close wait。

这个 close wait 是没有 pid 的，并且是探活客户端导致的，

2.2.2 诡异问题
那不是业务导致，那为什么会产生这种现象。排查的时候一堆问号，并且还发生很多诡异的事情。

我们发现 CLOSE WAIT 的客户端的端口号，在服务端机器上有，但在客户端机器上早就消失。CLOSE WAIT 就死在了这台机器上。

RECVQ 这么大，为什么服务还能正常响应 200

为什么上面 curl 会超时，但是 telnet 可以成功

为什么 CLOSE WAIT 没有 pid

2.2.3 线上调试
我们使用了sudo strace -p pid，发现主进程卡住了，想当然的认为自己发现了问题。但实际过程中，这是正常现象，我们框架使用了 golang.org/x/net/netutil 里的 LimitListener ，当你的 goroutine 达到最大值，那么就会出现这个阻塞现象，因为没有可用的 goroutine 了，主进程就卡住了。其他线程会提供服务，查看全部线程 trace 指令为sudo strace -f -p pid ，当然也可以查看单个线程的 ps -T -p pid ，然后拿到 spid，在执行 sudo strace -p spid 。这个调试还是没有太大用处，所以就想怎么在线下进行复现。



3 线下复现
复现这个是经过线下抓包，调试出来的。以下我们写个简单代码进行复现

3.1 server 端代码
package main

import (
    "fmt"
    "golang.org/x/net/netutil"
    "net"
    "net/http"
)

func main()  {
    mux := http.NewServeMux()
    mux.HandleFunc("/hello", func(w http.ResponseWriter,r *http.Request)  {
        fmt.Println("r.Body = ", r.Body)
        fmt.Fprintf(w,"HelloWorld!")
    })

    server := &http.Server{Addr: "", Handler: mux}
    listener, err := net.Listen("tcp4", "0.0.0.0:8880")
    if err != nil {
        fmt.Println("服务器错误")
    }
    // 这个地方要把限制变成1
    err = server.Serve(netutil.LimitListener(listener, 1))
    if err != nil {
        fmt.Println("服务器错误2")
    }
}


3.2 client 端代码
package main

import (
    "fmt"
    "io/ioutil"
    "net"
    "net/http"
    "strings"
    "time"
)

var httpclient *http.Client

func main() {
    for i := 0; i < 1000; i++ {
        req, err := http.NewRequest("POST", "http://127.0.0.1:8880/hello", strings.NewReader("{}") )
        if err != nil {
            fmt.Println(err)
            return
        }
        // 长连接
        req.Close = false
        httpclient = &http.Client{}
        resp, err := httpclient.Do(req)
        if err != nil {
            fmt.Println(err)
        }
        defer resp.Body.Close()
        body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            fmt.Println(string(body))
            fmt.Println(err)
        }

    }
    for i := 0; i < 1000; i++ {
        //time.Sleep(10 * time.Second)
        Dial("127.0.0.1:8880", time.Millisecond*500)
        //time.Sleep(100 * time.Second)
    }

}

// Dial dial 指定的端口，失败返回错误
func Dial(addr string, timeout time.Duration) error {
    conn, err := net.DialTimeout("tcp", addr, timeout)
    fmt.Println("dial err------>", err)
    if err != nil {
        return err
    }
    defer conn.Close()
    return err
}


3.3 描述
我们将服务端代码的 server.Serve(netutil.LimitListener(listener, 1)) 里面的限制设置为 1。然后客户端代码的 http 长连接开启，先做 http 请求，然后再做 tcp 请求。可以得到以下结果。

3.3.1 输出显示








3.3.2 命令行显示




















3.3.3 wireshark 显示


全部请求报文












一个 http 请求报文








一个 tcp 请求报文








3.4 分析
当我们在服务端设置了 limit 为 1 的时候，意味这我们服务端的最大连接数只能为 1。我们客户端在使用 http keepalive，会将这个连接占满，如果这个时候又进行了 tcp 探活，那么这个探活的请求就会被堵到 backlog 里，也就是上面我们看到 3.3.2 中第一个图，里面的 RECVQ 为 513。


我们先分析 3.3.3 中的前两个图。
黑色部分是 http keepalive，其客户端端口号为 43340，如图所示
![image.png]





灰色都是 tcp dial 操作，如图所示





http keepalive 的 close 报文，如图所示





红色的是客户端发现服务端有问题（客户端端口为 43346），进行了连接 RST 操作，如图所示






可以看到我们的 tcp 的 dial 操作，最后都是返回的 RST 操作，而且这个时间点刚好在 http keepalive 之后。接下来我们就看下 http keepalive 和 tcp dial 的详细报文。

3.3.3 第三个图，说明 http keepalive 是在他创建连接后，90s 后客户端发的 close 包。
3.3.3 第四个图，说明 tcp dial 发送了三次握手后，一次挥手后，客户端并没有立即收到接下来的两次挥手，而是过了 90s 后，才收到后两次挥手，导致了客户端 RST 操作。之所以等待 90s 的时间，是之前 tcp keepalive 的操作刚好进行了 close，goroutine 得到释放，可以使得 connection 被服务 accept，因此这个时候才能发送 fin 包。

我们这个时候如果在 RECVQ 不多的情况下。就可以复现一个场景，就是 curl 无法成功，但是 telnet 可以成功，如下所示。






这个和线上情况一样。

然后我们再来看下为什么线上 CLOSE_WAIT 在服务端上产生的过程。我们先抓了一个端口为 57688 的请求。可以看到 57688 挥手的时候发给 8880 的 fin 包，但是 8880 只响应了 ack。





所以现场抓包这两个状态变为 FIN_WAIT2 和 CLOSE_WAIT






我机器上 linux 配置了 FIN_TIMEOUT 为 40s





过 40s 后，可以看到客户端的 FIN_WAIT 被回收，只留下了服务端的 CLOSE_WAIT






再过 50s，http 的 keepalive 关闭连接，释放出 goroutine，tcp 的 fin 包就会发出。最终 CLOSE_WAIT 消失






以上我们把线上大部分场景都复现了，如下所示：

（复现了）服务端监听的 RECVQ 比较大

（复现了）服务端出现了 CLOSE WAIT，带没有 pid

（复现了）服务端出现了 CLOSE WAIT 的端口号，在客户端找不到

（复现了）服务端出现了 CLOSE WAIT，telnet 可以成功，但是没办法 curl 成功

（未复现）服务端的 CLOSE WAIT 一直不消失，下文中会解释，但并不确定



3.5 CLOSE WAIT 不消失的情况
出现这个情况比较极端，而且跟我们业务结合起来，有点麻烦，所以不在处理，请大家仔细阅读这个文章。https://blog.cloudflare.com/this-is-strictly-a-violation-of-the-tcp-specification/

















以上是 CLOSE WAIT 出现的一种场景。并没有完全验证这种情况。



4 产生线上问题的可能原因
线上的 nginx 到后端 go 配置的 keepalive，当 GO 的 HTTP 连接数达到系统默认 1024 值，那么就会出现 Goroutine 无法让出，这个时候使用 TCP 的探活，将会堵在队列里，服务端出现 CLOSE WAIT 情况，然后由于一直没有释放 goroutine，导致服务端无法发出 fin 包，一直保持 CLOSE WAIT。而客户端的端口在此期间可能会被重用，一旦重用后，就造成了混乱。（如果在混乱后，Goroutine 恢复后，服务端过了好久响应了 fin 包，客户端被重用的端口收到这个包是返回 RST，还是丢弃？这个太不好验证）。个人猜测是丢弃了，导致服务端的 CLOSE WAIT 一直无法关闭，造成 RECVQ 的一直阻塞。



5 其他问题
1. 因为 GO 服务端的 HTTP Keepalive 是使用的客户端的 Keepalive 的时间，如果客户端的 Keepalive 存在问题，比如客户端的 http keepalive 泄露，也会导致服务端无法关闭 Keepalive，从而无法回收 goroutine，当然 go 前面挡了一层 nginx，所以应该不会有这种泄露问题。但保险起见，go 的服务端应该加一个 keepalive 的最大值。例如 120s，避免这种问题。

2.GO 服务端的 HTTP chucked 编码时候，如果客户端没有正确将 response 的 body 内容取走，会导致数据仍然在服务端的缓冲区，从而导致无法响应 fin 包，但这个理论上不会出现，并且客户端会自动的进行 rst，不会影响业务。不过最好避免这种编码。

3 linux keepalive 参数可能需要配置，但不敢完全确信他的作用，如下所示。









6 其他 CLOSE WAIT 延伸
如过写的 http 服务存在业务问题，例如里面有个死循环，无法响应客户端，也会导致 http 服务出现 CLOSE WAIT 问题，但这个是有 pid 号的。

如果我们的业务调用某个服务的时候，由于没发心跳包给服务，会被服务关闭，但我们这个时候没正确处理这个场景，那么我们的业务处就会出现 CLOSE WAIT。

资料
线上大量 close wait 问题 

网络编程 tcp 链接 

tcp close wait 问题 

tcp keepalive 设置
