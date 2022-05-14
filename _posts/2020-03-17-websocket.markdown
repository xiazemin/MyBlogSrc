---
title: websocket
layout: post
category: golang
author: 夏泽民
---
https://www.gorillatoolkit.org/pkg/websocket
Overview
The Conn type represents a WebSocket connection. A server application calls the Upgrader.Upgrade method from an HTTP request handler to get a *Conn:
<!-- more -->
使用Go语言创建WebSocket服务
今天介绍如何用Go语言创建WebSocket服务，文章的前两部分简要介绍了WebSocket协议以及用Go标准库如何创建WebSocket服务。第三部分实践环节我们使用了gorilla/websocket库帮助我们快速构建WebSocket服务，它帮封装了使用Go标准库实现WebSocket服务相关的基础逻辑，让我们能从繁琐的底层代码中解脱出来，根据业务需求快速构建WebSocket服务。

Go Web 编程系列的每篇文章的源代码都打了对应版本的软件包，供大家参考。公众号中回复gohttp10获取本文源代码

WebSocket介绍
WebSocket通信协议通过单个TCP连接提供全双工通信通道。与HTTP相比，WebSocket不需要你为了获得响应而发送请求。它允许双向数据流，因此您只需等待服务器发送的消息即可。当Websocket可用时，它将向您发送一条消息。 对于需要连续数据交换的服务（例如即时通讯程序，在线游戏和实时交易系统），WebSocket是一个很好的解决方案。 WebSocket连接由浏览器请求，并由服务器响应，然后建立连接，此过程通常称为握手。 WebSocket中的特殊标头仅需要浏览器与服务器之间的一次握手即可建立连接，该连接将在其整个生命周期内保持活动状态。 WebSocket解决了许多实时Web开发的难题，并且与传统的HTTP相比，具有许多优点：

轻量级报头减少了数据传输开销。
单个Web客户端仅需要一个TCP连接。
WebSocket服务器可以将数据推送到Web客户端。
WebSocket协议实现起来相对简单。它使用HTTP协议进行初始握手。握手成功后即建立连接，WebSocket实质上使用原始TCP读取/写入数据。


客户端请求如下所示：

GET /chat HTTP/1.1
    Host: server.example.com
    Upgrade: websocket
    Connection: Upgrade
    Sec-WebSocket-Key: x3JJHMbDL1EzLkh9GBhXDw==
    Sec-WebSocket-Protocol: chat, superchat
    Sec-WebSocket-Version: 13
    Origin: http://example.com
复制代码
这是服务器响应：

HTTP/1.1 101 Switching Protocols
    Upgrade: websocket
    Connection: Upgrade
    Sec-WebSocket-Accept: HSmrc0sMlYUkAGmm5OPpG2HaGWk=
    Sec-WebSocket-Protocol: chat
复制代码
如何在Go中创建WebSocket应用
要基于Go 语言内置的net/http 库编写WebSocket服务器，你需要：

发起握手
从客户端接收数据帧
发送数据帧给客户端
关闭握手
发起握手
首先，让我们创建一个带有WebSocket端点的HTTP处理程序：

// HTTP server with WebSocket endpoint
func Server() {
        http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
            ws, err := NewHandler(w, r)
            if err != nil {
                 // handle error
            }
            if err = ws.Handshake(); err != nil {
                // handle error
            }
        …
复制代码
然后初始化WebSocket结构。

初始握手请求始终来自客户端。服务器确定了WebSocket请求后，需要使用握手响应进行回复。

请记住，你无法使用http.ResponseWriter编写响应，因为一旦开始发送响应，它将关闭其基础的TCP连接（这是HTTP 协议的运行机制决定的，发送响应后即关闭连接）。

因此，您需要使用HTTP劫持(hijack)。通过劫持，可以接管基础的TCP连接处理程序和bufio.Writer。这使可以在不关闭TCP连接的情况下读取和写入数据。

// NewHandler initializes a new handler
func NewHandler(w http.ResponseWriter, req *http.Request) (*WS, error) {
        hj, ok := w.(http.Hijacker)
        if !ok {
            // handle error
        }                  .....
}
复制代码
要完成握手，服务器必须使用适当的头进行响应。

// Handshake creates a handshake header
    func (ws *WS) Handshake() error {

        hash := func(key string) string {
            h := sha1.New()
            h.Write([]byte(key))
            h.Write([]byte("258EAFA5-E914-47DA-95CA-C5AB0DC85B11"))

        return base64.StdEncoding.EncodeToString(h.Sum(nil))
        }(ws.header.Get("Sec-WebSocket-Key"))
      .....
}
复制代码
客户端发起WebSocket连接请求时用的Sec-WebSocket-key是随机生成的，并且是Base64编码的。接受请求后，服务器需要将此密钥附加到固定字符串。假设秘钥是x3JJHMbDL1EzLkh9GBhXDw==。在这个例子中，可以使用SHA-1计算二进制值，并使用Base64对其进行编码。得到HSmrc0sMlYUkAGmm5OPpG2HaGWk=。然后使用它作为Sec-WebSocket-Accept 响应头的值。

传输数据帧
握手成功完成后，您的应用程序可以从客户端读取数据或向客户端写入数据。WebSocket规范定义了的一个客户机和服务器之间使用的特定帧格式。这是框架的位模式：

img{512x368}
图:传输数据帧的位模式
使用以下代码对客户端有效负载进行解码：

// Recv receives data and returns a Frame
    func (ws *WS) Recv() (frame Frame, _ error) {
        frame = Frame{}
        head, err := ws.read(2)
        if err != nil {
            // handle error
        }
复制代码
反过来，这些代码行允许对数据进行编码：

// Send sends a Frame
    func (ws *WS) Send(fr Frame) error {
        // make a slice of bytes of length 2
        data := make([]byte, 2)

        // Save fragmentation & opcode information in the first byte
        data[0] = 0x80 | fr.Opcode
        if fr.IsFragment {
            data[0] &= 0x7F
        }
        .....
复制代码
关闭握手
当各方之一发送状态为关闭的关闭帧作为有效负载时，握手将关闭。可选的，发送关闭帧的一方可以在有效载荷中发送关闭原因。如果关闭是由客户端发起的，则服务器应发送相应的关闭帧作为响应。

// Close sends a close frame and closes the TCP connection
func (ws *Ws) Close() error {
    f := Frame{}
    f.Opcode = 8
    f.Length = 2
    f.Payload = make([]byte, 2)
    binary.BigEndian.PutUint16(f.Payload, ws.status)
    if err := ws.Send(f); err != nil {
        return err
    }
    return ws.conn.Close()
}
复制代码
使用第三方库快速构建WebSocket服务
通过上面的章节可以看到用Go自带的net/http库实现WebSocket服务还是太复杂了。好在有很多对WebSocket支持良好的第三方库，能减少我们很多底层的编码工作。这里我们使用gorilla web toolkit家族的另外一个库gorilla/websocket来实现我们的WebSocket服务，构建一个简单的Echo服务（echo意思是回音，就是客户端发什么，服务端再把消息发回给客户端）。

我们在http_demo项目的handler目录下新建一个ws子目录用来存放WebSocket服务相关的路由对应的请求处理程序。

增加两个路由：

/ws/echo echo应用的WebSocket 服务的路由。
/ws/echo_display echo应用的客户端页面的路由。
创建WebSocket服务端
// handler/ws/echo.go
package ws

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func EchoMessage(w http.ResponseWriter, r *http.Request) {
	conn, _ := upgrader.Upgrade(w, r, nil) // 实际应用时记得做错误处理

	for {
		// 读取客户端的消息
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}

		// 把消息打印到标准输出
		fmt.Printf("%s sent: %s\n", conn.RemoteAddr(), string(msg))

		// 把消息写回客户端，完成回音
		if err = conn.WriteMessage(msgType, msg); err != nil {
			return
		}
	}
}
复制代码
conn变量的类型是*websocket.Conn, websocket.Conn类型用来表示WebSocket连接。服务器应用程序从HTTP请求处理程序调用Upgrader.Upgrade方法以获取*websocket.Conn

调用连接的WriteMessage和ReadMessage方法发送和接收消息。上面的msg接收到后在下面又回传给了客户端。msg的类型是[]byte。

创建WebSocket客户端
前端页面路由对应的请求处理程序如下，直接返回views/websockets.html给到浏览器渲染页面即可。

// handler/ws/echo_display.go
package ws

import "net/http"

func DisplayEcho(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "views/websockets.html")
}
复制代码
websocket.html里我们需要用JavaScript连接WebScoket服务进行收发消息，篇幅原因我就只贴JS代码了，完整的代码通过本节的口令去公众号就能获取到下载链接。

<form>
    <input id="input" type="text" />
    <button onclick="send()">Send</button>
    <pre id="output"></pre>
</form>
...
<script>
    var input = document.getElementById("input");
    var output = document.getElementById("output");
    var socket = new WebSocket("ws://localhost:8000/ws/echo");

    socket.onopen = function () {
        output.innerHTML += "Status: Connected\n";
    };

    socket.onmessage = function (e) {
        output.innerHTML += "Server: " + e.data + "\n";
    };

    function send() {
        socket.send(input.value);
        input.value = "";
    }
</script>
...
复制代码
注册路由
服务端和客户端的程序都准备好后，我们按照之前约定好的路径为他们注册路由和对应的请求处理程序：

// router/router.go
func RegisterRoutes(r *mux.Router) {
    ...
    wsRouter := r.PathPrefix("/ws").Subrouter()
    wsRouter.HandleFunc("/echo", ws.EchoMessage)
    wsRouter.HandleFunc("/echo_display", ws.DisplayEcho)
}
复制代码
测试验证
重启服务后访问http://localhost:8000/ws/echo_display，在输入框中输入任何消息都能再次回显到浏览器中。

图片
服务端则是把收到的消息打印到终端中然后把调用writeMessage把消息再回传给客户端，可以在终端中查看到记录。

image-20200316142506287
总结
WebSocket在现在更新频繁的应用中使用非常广泛，进行WebSocket编程也是我们需要掌握的一项必备技能。文章的实践练习稍微简单了一些，也没有做错误和安全性检查。主要是为了讲清楚大概的流程。关于gorilla/websocket更多的细节在使用时还需要查看官方文档才行。

参考链接：

yalantis.com/blog/how-to…

www.gorillatoolkit.org/pkg/websock…
