---
title: connectpool
layout: post
category: golang
author: 夏泽民
---
net/http 长连接验证 ?

默认是长连接，毋庸置疑. 客户端发起的时候会在header里标记HTTP/1.1 。

net/http 连接复用 ?

连接可复用. 只要匹配到目标ip及端口就可以服用到该维度的连接池.

直接注释 response.Body.Close() 会出现什么?

各种循环测试，不仅长连接，而且连接还是会被复用. 社区里有人反映说close注释掉会出现连接不停创建的情况.

如果连接池中，某个主机的连接被占用，上层并发请求会发生什么?

net/http在池中找不到有用的连接，就会不断的重新new一个连接，不会阻塞等待一个连接。
<!-- more -->
对端关闭，上层代码如果不管不问的会出现什么?

go runtime 会在底层一直帮你epoll wait, 监听读事件的close报文 (空值报文) ，接着自动帮你做close fd相关, 然后在上层标记出网络连接是否发生关闭. 哪怕你的逻辑只是成功发起请求后，一直等待的sleep下去。

通过Linux strace系统调用、tcpdump能看到fin的过程，可以用tcpdump把包导出到wireshark查看.

http://xiaorui.cc/archives/5056

长连接和短链接的区别
客户端和服务端响应的次数
长连接：可以多次。
短链接：一次。
传输数据的方式
长连接：连接--数据传输--保持连接
短连接：连接--数据传输--关闭连接
长连接和短链接的优缺点
长连接
优点
省去较多的TCP建立和关闭的操作，从而节约时间。
性能比较好。（因为客户端一直和服务端保持联系）
缺点
当客户端越来越多的时候，会将服务器压垮。
连接管理难。
安全性差。（因为会一直保持着连接，可能会有些无良的客户端，随意发送数据等）
短链接
优点
服务管理简单。存在的连接都是有效连接
缺点
请求频繁，在TCP的建立和关闭操作上浪费时间
长连接和短连接使用情况举例
长连接
微信/qq
一些游戏
短连接
普通的web网站

golang的client实现长连接的方式
web Server 支持长连接。（golang默认支持长连接）.client要和服务端响应之后，保持连接
根据需求，加大：DefaultMaxIdleConnsPerHost或设置MaxIdleConnsPerHost
读完Response Body再Close

https://www.zhihu.com/question/22925358
https://github.com/henrylee2cn/erpc
https://github.com/davyxu/cellnet


Websocket简介
WebSocket可以实现客户端与服务器间双向、基于消息的文本或二进制数据传输。它是浏览器中最靠近套接字的API。但WebSocket连接远远不是一个网络套接字，因为浏览器在这个简单的API之后隐藏了所有的复杂性，而且还提供了更多服务：

连接协商和同源策略；
与既有HTTP基础设施的互操作；
基于消息的通信和高效消息分帧；
子协议协商及可扩展能力。
WebSocket资源URL采用了自定义模式：ws表示纯文本通信（如ws://example.com/socket），wss表示使用加密信道通信（TCP+TLS）。

长连接实现
本例用gin框架，引入github.com/gorilla/websocket包，项目源码可到https://github.com/shidawuhen/asap/tree/feature_pzq_longconnect查看

服务端核心代码:

ping函数，将请求升级为websocket。
只有Get请求才能进行升级，具体原因可以查看websocket源码
当客户端请求ping时，便建立了长连接，服务端读取客户端数据，如果数据为ping，则返回pong，如果不为ping，则把输入的内容返回。
package main

import (
   "net/http"

   "github.com/gin-gonic/gin"
   "github.com/gorilla/websocket"
)

var upGrader = websocket.Upgrader{
   CheckOrigin: func (r *http.Request) bool {
      return true
   },
}

func setupRouter() *gin.Engine {
   r := gin.Default()
   r.LoadHTMLGlob("templates/*")
   // Ping test
   r.GET("/ping", ping)
   r.GET("/longconnecthtml",longconnecthtml)

   return r
}

func longconnecthtml(c *gin.Context)  {
   c.HTML(http.StatusOK, "longconnect.tmpl",gin.H{})
}

func ping(c *gin.Context) {
   //c.String(http.StatusOK, "ok")
   //升级get请求为webSocket协议
   ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
   if err != nil {
      return
   }
   defer ws.Close()
   for {
      //读取ws中的数据
      mt, message, err := ws.ReadMessage()
      if err != nil {
         break
      }
      if string(message) == "ping" {
         message = []byte("pong")
      }
      //写入ws数据
      err = ws.WriteMessage(mt, message)
      if err != nil {
         break
      }

   }
}

func main() {
   r := setupRouter()
   // Listen and Server in 0.0.0.0:8080
   r.Run(":9090")
}
客户端核心代码：

要使用websocket，客户端必须发起websocket请求
websocket的建立和使用也很方便，主要涉及
new WebSocket：创建websocket对象
onopen：连接建立时触发
onmessage：客户端接收服务端数据时触发
onerror：通信发生错误时触发
onclose：连接关闭时触发
该代码块的主要功能是和服务端建立连接，在文本框输入信息，将信息发送给服务端，并将服务端返回内容显示出来

<!DOCTYPE HTML>
<html>
   <head>
   <meta charset="utf-8">
   <title>长连接测试</title>
   <textarea id="inp_send" class="form-control" style="height:100px;" placeholder="发送的内容"></textarea>
   <button type="button" id="btn_send" class="btn btn-info" onclick="fun_sendto();">发送（ctrl+回车）</button>
<script type="text/javascript" src="https://dss0.bdstatic.com/5aV1bjqh_Q23odCf/static/superman/js/lib/jquery-1-edb203c114.10.2.js"></script>
      <script type="text/javascript">
        var ws = new WebSocket("ws://localhost:9090/ping");
        //连接打开时触发
        ws.onopen = function(evt) {
            console.log("Connection open ...");
            ws.send("Hello WebSockets!");
        };
        //接收到消息时触发
        ws.onmessage = function(evt) {
            console.log("Received Message: " + evt.data);
        };
        //连接关闭时触发
        ws.onclose = function(evt) {
            console.log("Connection closed.");
        };

        function fun_sendto(){
            var content = $("#inp_send").val();
            ws.send(content);
        }
      </script>

   </head>
   <body>
   </body>
</html>
服务端运行起来之后，请求http://localhost:9090/longconnecthtml即可查看效果

长连接展示
浏览器，请求ping接口，便可看到该请求为websocket请求，而且为pending状态

也可以查看到在发起请求时，request header中有很多新的参数，具体含义大家可以看我提供的参考资料

console中显示的是服务端推送到客户端的数据

注意事项
长连接可以实现双向通信，但是以占用连接为前提的，如果请求量较大，需要考虑资源问题
服务部署情况对效果影响很大。例如，如果只部署北京机房供全国使用，当服务端向客户端推送数据时，不同地区可能会有很大延迟。
参考资料
golang长连接和短连接的学习
https://www.zhihu.com/question/22925358
Golang-长连接-状态推送
Golang实现单机百万长连接服务 - 美图的三年优化经验
golang 长连接web socket原理
golang Gin建立长连接web socket
HTML5 WebSocket

https://blog.csdn.net/shida219/article/details/106762944


CP 相关
长连接的概念包括TCP长连接和HTTP长连接。首先得保证TCP是长连接。我们就从它说起。

func (c *TCPConn) SetKeepAlive(keepalive bool) error
SetKeepAlive sets whether the operating system should send keepalive messages on the connection. 这个方法比较简单，设置是否开启长连接。

func (c *TCPConn) SetReadDeadline(t time.Time) error
SetReadDeadline sets the deadline for future Read calls and any currently-blocked Read call. A zero value for t means Read will not time out.这个函数就很讲究了。我之前的理解是设置读取超时时间，这个方法也有这个意思，但是还有别的内容。它设置的是读取超时的绝对时间。

func (c *TCPConn) SetWriteDeadline(t time.Time) error
SetWriteDeadline sets the deadline for future Write calls and any currently-blocked Write call. Even if write times out, it may return n > 0, indicating that some of the data was successfully written. A zero value for t means Write will not time out. 这个方法是设置写超时，同样是绝对时间。

HTTP 包如何使用 TCP 长连接？
http 服务器启动之后，会循环接受新请求，为每一个请求（连接）创建一个协程。

// net/http/server.go L1892
for {
    rw, e := l.Accept()
    go c.serve()
}
下面是每个协程的执行的代码，我只摘录了一部分关键的逻辑。可以发现，serve方法里面还有一个for循环。

// net/http/server.go L1320
func (c *conn) serve() {
    defer func() {
        if !c.hijacked() {
            c.close()
        }
    }()

    for {
        w, err := c.readRequest()
        
        if err != nil {
        }
        
        serverHandler{c.server}.ServeHTTP(w, w.req)
    }
}
这个循环是用来做什么的？其实也容易理解，如果是长连接，一个协程可以执行多次响应。如果只执行了一次，那就是短连接。长连接会在超时或者出错后退出循环，也就是关闭长连接。defer函数可以让协程结束之后关闭 TCP 连接。

readRequest函数用来解析 HTTP 协议。

// net/http/server.go
func (c *conn) readRequest() (w *response, err error) {
    if d := c.server.ReadTimeout; d != 0 {
        c.rwc.SetReadDeadline(time.Now().Add(d))
    }
    if d := c.server.WriteTimeout; d != 0 {
        defer func() {
            c.rwc.SetWriteDeadline(time.Now().Add(d))
        }()
    }
    
    if req, err = ReadRequest(c.buf.Reader); err != nil {
        if c.lr.N == 0 {
            return nil, errTooLarge
        }
        return nil, err
    }
}

func ReadRequest(b *bufio.Reader) (req *Request, err error) {
    // First line: GET /index.html HTTP/1.0
    var s string
    if s, err = tp.ReadLine(); err != nil {
        return nil, err
    }
    
    req.Method, req.RequestURI, req.Proto, ok = parseRequestLine(s)
    
    mimeHeader, err := tp.ReadMIMEHeader()
}
具体参与解析 HTTP 协议的部分是ReadRequest方法，而调用它之前，设置了读写超时时间。根据前面的描述，超时时间设置的是绝对时间。所以这里都是通过time.Now().Add(d)来设置的。不同的是写超时是defer执行，也就是函数返回后才执行。

我们的程序为啥长连接失效？
通过源码我们能大概知道程序流程了，按道理是支持长连接的。为啥我们的程序不行呢？

我们的程序使用的是 beego 框架，它支持的超时是同时设置读写超时。而我们的设置是1秒。

beego.HttpServerTimeOut = 1
我对读写超时的理解，读超时是收到数据到读取完毕的时间；写超时是从一开始写到写完的时间。我对这两个超时的理解都不对。

实际上，从上面的源码可以发现，写超时是读取完毕之后设置的超时时间。也就是读取完毕之后的时间，加上逻辑执行时间，加上内容返回时间的总和。按照我们的设置，超过1秒就算超时。

下面详细说说读超时。ReadRequest是堵塞执行的，如果没有用户请求，它会一直等待着。而读超时是ReadRequest之前设置的，它除了读取数据之外，还有一部分耗时，那就是等待时间。假如一直没有用户请求，此时读超时已经被设置成1秒后了，超过1秒之后，这个连接还是会被断开。

如何解决问题？
原因已经说明白了。大量TIME_WAIT是超时引起的，有可能是等待时间过长引起的读超时；也有可能是程序在压测情况下出现一部分执行超时，这样会导致写超时。

我们目前使用的是 beego 框架，它并不支持单独设置读写超时，所以我目前的解决方式是将读写超时调整得大一些。

从1.6版本开始，Golang 能够支持空闲超时IdleTimeout，可以认为读超时就是读取数据的时间，空闲超时来控制等待时间。但是它有一个问题，如果空闲超时没有设置，而读超时设置了，那么读超时还是会作为空闲超时时间来使用。我估计这么做的原因是为了向前兼容。再一个问题就是 beego 并不支持这个时间的设置，所以我目前也没有别的太好的方法来控制超时时间。

后续
其实服务端最合理的超时控制需要这几个方面：

读超时。就是单纯的读超时，不要包括等待时间，否则无法区分超时是读数据引起的还是等待引起的。

写超时。最好也是单纯的写数据超时。如果网络良好，因为逻辑执行慢就把连接断开，这样也不是很合适。读写超时都应该和目前逻辑设置的一样，设置得短一些。

空闲超时。这个可以根据实际情况配置，可以适当大一些。

逻辑超时。一般情况下是不会发生网络层面的读写超时的，压测情况下超时大部分都是由于逻辑超时引起的。Golang 原生包支持了TimeoutHandler。它可以控制逻辑的超时。可惜 beego 目前不支持设置逻辑超时。而我也没有想到太好的方法把 beego 中接入它。


 func TimeoutHandler(h Handler, dt time.Duration, msg string) Handler

https://www.jianshu.com/p/99832e6bab73

working directory is not part of a module 问题处理
go mod init


https://www.ctolib.com/mip/connsvr.html
https://blog.cyeam.com/golang/2017/05/31/go-http-keepalive
https://github.com/simplejia/connsvr/