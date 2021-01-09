---
title: tcp http WebSocket长连接区别
layout: post
category: web
author: 夏泽民
---
一、HTTP的长连接和短连接区别
首先需要消除一个误解：HTTP协议是基于请求/响应模式的，因此客户端请求后只要服务端给了响应，本次HTTP请求就结束了，没有长连接这一说。那么自然也就没有短连接这一说了。

所谓的HTTP分为长连接和短连接，其实本质上是说的TCP连接。TCP连接是一个双向的通道，它是可以保持一段时间不关闭的，因此TCP连接才有真正的长连接和短连接这一说。

HTTP协议是应用层的协议，而TCP才是真正的传输层协议，只有负责传输的这一层才需要建立连接。

1、短连接
过程：连接->传输数据->关闭连接 
短链接就是浏览器和服务器每进行一次HTTP操作，就建立一次连接，但任务结束就中断连接。 比如HTTP1.0。
具体就是 浏览器client发起并建立TCP连接 -> client发送HttpRequest报文 -> server接收到报文->server handle并发送HttpResponse报文给前端,发送完毕之后立即调用socket.close方法->client接收response报文->client最终会收到server端断开TCP连接的信号->client 端断开TCP连接，具体就是调用close方法。 
也就是说，短连接是指SOCKET连接后，发送接收完数据后马上断开连接。 因为连接后接收了数据就断开了，所以每次数据接受处理不会有联系。 这也是HTTP协议无状态的原因之一。

2、长连接
过程：连接->传输数据->保持连接 -> 传输数据-> ………..->一方关闭连接

长连接指建立SOCKET连接后不管是否使用都保持TCP连接。

HTTP1.1默认是长连接，也就是默认Connection的值就是keep-alive，本次请求响应结束后，TCP连接将仍然保持打开状态，所以浏览器可以继续通过相同的连接发送请求，节省了很多TCP连接建立和断开的消耗，还节约了带宽。

长连接并不是永久连接的。如果一段时间内（具体的时间可以在header当中进行设置，也就是所谓的超时时间），这个连接没有HTTP请求发出的话，那么这个长连接就会被断掉。这一点其实很容易理解，否则的话，TCP连接将会越来越多，直到把服务器的TCP连接数量撑爆为止。
<!-- more -->
{% raw %}
二、HTTP长连接和WebSocket长连接的区别
HTTP1.1中，Connection默认为Keep-alive参数，官方的说法是可以用这个来作为长连接。那么问题来了，既然HTTP1.1支持长连接，为什么还要搞出一个WebSocket呢？

1、HTTP1.1
Keep-alive的确可以实现长连接，但是这个长连接是有问题的，本质上依然是客户端主动发起-服务端应答的模式，是没法做到服务端主动发送通知给客户端的。也就是说，在一个HTTP连接中，可以发送多个Request，接收多个Response。但是一个request只能有一个response。而且这个response也是被动的，不能主动发起。开启了Keep-alive，可以看出依然是一问一答的模式，只是省略了每次的关闭和打开操作。

2、WebSocket
WebSocket是可以互相主动发起的。相对于传统 HTTP 每次请求-应答都需要客户端与服务端建立连接的模式，WebSocket 是类似 TCP 长连接的通讯模式，一旦 WebSocket 连接建立后，后续数据都以帧序列的形式传输。在客户端断开 WebSocket 连接或 Server 端断掉连接前，不需要客户端和服务端重新发起连接请求。在海量并发及客户端与服务器交互负载流量大的情况下，极大的节省了网络带宽资源的消耗，有明显的性能优势，且客户端发送和接受消息是在同一个持久连接上发起，实时性优势明显。

WebSocket API 是 HTML5 标准的一部分， 但这并不代表 WebSocket 一定要用在 HTML 中，或者只能在基于浏览器的应用程序中使用。
在WebSocket中，只需要服务器和浏览器通过HTTP协议进行一个握手的动作，然后单独建立一条TCP的通信通道进行数据的传送。WebSocket同HTTP一样也是应用层的协议，但是它是一种双向通信协议，是建立在TCP之上的。

WebSocket的流程大概是以下几步：

1、浏览器、服务器建立TCP连接，三次握手。这是通信的基础，传输控制层，若失败后续都不执行。
2、TCP连接成功后，浏览器通过HTTP协议向服务器传送WebSocket支持的版本号等信息。（开始前的HTTP握手）服务器收到客户端的握手请求后，同样采用HTTP协议回馈数据。
3、连接成功后，双方通过TCP通道进行数据传输，不需要HTTP协议。
也就是说WebSocket在建立握手时，数据是通过HTTP传输的。但是建立之后，在真正传输时候是不需要HTTP协议的。

WebSocket 客户端连接报文

GET /webfin/websocket/ HTTP/1.1
Host: localhost
Upgrade: websocket
Connection: Upgrade
Sec-WebSocket-Key: xqBt3ImNzJbYqRINxEFlkg==
Origin: 
http://localhost
:8080
Sec-WebSocket-Version: 13
客户端发起的 WebSocket 连接报文类似传统 HTTP 报文，”Upgrade：websocket”参数值表明这是 WebSocket 类型请求，“Sec-WebSocket-Key”是 WebSocket 客户端发送的一个 base64 编码的密文，要求服务端必须返回一个对应加密的“Sec-WebSocket-Accept”应答，否则客户端会抛出“Error during WebSocket handshake”错误，并关闭连接。

服务端收到报文后返回的数据格式类似：
WebSocket 服务端响应报文：

HTTP/1.1 101 Switching Protocols
Upgrade: websocket
Connection: Upgrade
Sec-WebSocket-Accept: K7DJLdLooIwIG/MOpvWFB3y3FE8=
“Sec-WebSocket-Accept”的值是服务端采用与客户端一致的密钥计算出来后返回客户端的,“HTTP/1.1 101 Switching Protocols”表示服务端接受 WebSocket 协议的客户端连接，经过这样的请求-响应处理后，客户端服务端的 WebSocket 连接握手成功, 后续就可以进行 TCP 通讯了。
{% endraw %}
HTTP是一个应用层协议，无状态的,端口号为80。主要的版本有1.0/1.1/2.0.

HTTP/1.* 一次请求-响应，建立一个连接，用完关闭；

HTTP/1.1 串行化单线程处理，可以同时在同一个tcp链接上发送多个请求，但是只有响应是有顺序的，只有上一个请求完成后，下一个才能响应。一旦有任务处理超时等，后续任务只能被阻塞(线头阻塞)；

HTTP/2 并行执行。某任务耗时严重，不会影响到任务正常执行

 

什么是websocket
Websocket是html5提出的一个协议规范，是为解决客户端与服务端实时通信。本质上是一个基于tcp，先通过HTTP/HTTPS协议发起一条特殊的http请求进行握手后创建一个用于交换数据的TCP连接。

 

WebSocket优势： 浏览器和服务器只需要要做一个握手的动作，在建立连接之后，双方可以在任意时刻，相互推送信息。同时，服务器与客户端之间交换的头信息很小。

 

 

什么是长连接、短连接
短连接：

连接->传输数据->关闭连接

HTTP是无状态的，浏览器和服务器每进行一次HTTP操作，就建立一次连接，但任务结束就中断连接。

也可以这样说：短连接是指SOCKET连接后发送后接收完数据后马上断开连接。

 

长连接:

连接->传输数据->保持连接 -> 传输数据-> 。。。 ->关闭连接。

长连接指建立SOCKET连接后不管是否使用都保持连接，但安全性较差。

 

 

http和websocket的长连接区别
HTTP1.1通过使用Connection:keep-alive进行长连接，HTTP 1.1默认进行持久连接。在一次 TCP 连接中可以完成多个 HTTP 请求，但是对每个请求仍然要单独发 header，Keep-Alive不会永久保持连接，它有一个保持时间，可以在不同的服务器软件（如Apache）中设定这个时间。这种长连接是一种“伪链接”

websocket的长连接，是一个真的全双工。长连接第一次tcp链路建立之后，后续数据可以双方都进行发送，不需要发送请求头。

 

keep-alive双方并没有建立正真的连接会话，服务端可以在任何一次请求完成后关闭。WebSocket 它本身就规定了是正真的、双工的长连接，两边都必须要维持住连接的状态。

在http1.1中，Connection默认为Keep-alive参数，官方的说法是可以用这个来作为长连接。那么问题来了，既然http1.1支持长连接，为什么还要搞出一个WebSocket呢？

关于Keep-alive的缺点
Keep-alive的确可以实现长连接，但是这个长连接是有问题的，本质上依然是客户端主动发起-服务端应答的模式，是没法做到服务端主动发送通知给客户端的。也就是说，在一个HTTP连接中，可以发送多个Request，接收多个Response。但是一个request只能有一个response。而且这个response也是被动的，不能主动发起。放上一张图，图左为没开启Keep-alive，图右为开启了Keep-alive，可以看出依然是一问一答的模式，相较左边只是省略了每次的关闭和打开操作。

https://www.jianshu.com/p/86a550a521c5
https://www.cnblogs.com/ricklz/p/11108320.html

https://www.cnblogs.com/show58/p/12362480.html

websocket连接过程概述
WebSocket 建立连接需要先通过一个 http 请求进行和服务端握手。握手通过后连接就建立并保持了。
浏览器先发送请求：

GET / HTTP/1.1
Host: localhost:8080
Origin: [url=http://127.0.0.1:3000]http://127.0.0.1:3000[/url]
Connection: Upgrade
Upgrade: WebSocket
Sec-WebSocket-Version: 13
Sec-WebSocket-Key: w4v7O6xFTi36lq3RNcgctw==
服务端返回一个请求：

HTTP/1.1 101 Switching Protocols
Connection:Upgrade
Upgrade: WebSocket
Sec-WebSocket-Accept: Oy4NRAQ13jhfONC7bP8dTKb4PTU=
这样握手就完成了（具体 http 请求头和返回头各个字段的含义我就懒得写了，网上一搜一大把，譬如http://www.52im.net/thread-13...）。此时这个连接并不会断掉，而浏览器和服务端可以用这个连接相互发消息。（但是这个时候连接就不是 http 连接而是升级成了 WebSocket 连接。浏览器和服务端相互发送的不是 http 请求。这里先说明下，接下来我们来看下 http 长连接是怎么回事）。

http长连接类型
keep-alive
http1.1 出了新头，如果请求头中包含 keep-alive，那么这个 http 请求发送收到返回之后，底层的 tcp 连接不会立马断掉，如果后续有 http 请求还是会利用。但是这个连接保持一来是没有硬性规定时间的，由浏览器和服务端实现来控制。二来这个连接不断是指底层 tcp 连接，不是说一次 http 请求收到返回之后不会断掉，还能再收服务端的返回（如果服务端对这次 http 请求立马返回，那么这次 http 请求就结束了。http 请求和底层 tcp 连接的关系后面再说）。这种不是应用层面的长连接，其实和模拟 WebSocket 没啥关系。

comet
这种技术是一种 hack 技术，即浏览器发送一个 http 请求，但是服务端不是立马返回，服务端一直不返回直到有浏览器需要的内容了在返回。期间这个 http 请求可以连着维持比较长的时间（在服务端返回之前）。这样模拟一种服务端推送机制。因为浏览器请求的时候等于先把连接建立好，等服务端有消息需要返回时再返回给浏览器。

websocket和http长连接的区别
先说 comet 和 WebSocket 表现的区别：
comet 发送 http 请求后服务端如果没有返回则连接是一直连着的，等服务端有东西要“推送”给浏览器时，相当于给之前发送的这个 http 请求回了一个 http 响应。然后这个保持的时间比较长的 http 连接就断了。然后浏览器再次发送一个 http 请求，服务器端再 hold 住不返回，等待有东西需要“推送”给浏览器时，再给这个 http 请求一个响应，然后断开连接。循环往复。一旦浏览器不给服务器发送 http 请求，那么服务器是不能主动给浏览器推送消息的，因为根本没有连着的连接给你推。

WebSocket 则不同，它握手后建立的连接是不会断的（除了意外情况和程序主动掐断）。不需要浏览器在每次收到服务器推送的消息后再发起请求。而且服务器端可以随时给浏览器推送消息，不需要等浏览器发 http 请求，因为 WebSocket 的连接一直在没断。

为什么会有这样的区别？

这是协议层面的区别。http 协议规定了 http 连接是一个一来（request）一回（response）的过程。一个请求获得一个响应后必须断掉。而且只有先有请求才会有响应。拿 http1.1 keep-alive 来说，即使底层 tcp 连接没有断，服务端无缘无故给浏览器发一个 http 响应，浏览器是不收的，他找不到收的人啊，因为这个响应没有对应的请求。你看 ajax 必须先发请求才会有一个 onsuccess 回调来响应这个请求。这个 onsuccess 的回调会在你 ajax 不发送的情况下被调用到吗？

而 WebSocket 协议不同，他通过握手之后规定说你连接给我保持着，别断咯。所以浏览器服务器在这种情况下可以相互的发送消息。浏览器端 new 一个 WebSocket 之后注册 onmessage 回调，那么这个 onmessage 可以被反复调用，只要服务器端有消息过来。而不会说是 new 一个 WebSocket onmessage 只会被调用一次，下次还得再 new 一个 websocket。

上面说到 http 连接，tcp 连接，websockt 连接到底啥区别。其实这是新人最容易搞不懂的地方。接下来我就要胡诌了，为啥说胡诌，因为我只是看了个皮毛，然后按我自己的理解说下区别。网络5层分层（自下而上）：

物理层
数据链路层
网络层
传输层
应用层
http，websocket都是应用层协议，他们规定的是数据怎么封装，而他们传输的通道是下层提供的。就是说无论是 http 请求，还是 WebSocket 请求，他们用的连接都是传输层提供的，即 tcp 连接（传输层还有 udp 连接）。只是说 http1.0 协议规定，你一个请求获得一个响应后，你要把连接关掉。所以你用 http 协议发送的请求是无法做到一直连着的（如果服务器一直不返回也可以保持相当一段时间，但是也会有超时而被断掉）。而 WebSocket 协议规定说等握手完成后我们的连接不能断哈。虽然 WebSocket 握手用的是 http 请求，但是请求头和响应头里面都有特殊字段，当浏览器或者服务端收到后会做相应的协议转换。所以 http 请求被 hold 住不返回的长连接和 WebSocket 的连接是有本质区别的。

https://segmentfault.com/a/1190000015122195

https://tools.ietf.org/html/rfc6455
https://blog.csdn.net/done58/article/details/50996680

https://www.cnblogs.com/show58/p/12362480.html

一、HTTP协议和TCP协议
HTTP的长连接和短连接本质上是TCP长连接和短连接。HTTP属于应用层协议，在传输层使用TCP协议，在网络层使用IP协议。IP协议主要解决网络路由和寻址问题，TCP协议主要解决如何在IP层之上可靠的传递数据包，使在网络上的另一端收到发端发出的所有包，并且顺序与发出顺序一致。TCP有可靠，面向连接的特点。

二、HTTP协议的长连接和短连接
在HTTP/1.0中，默认使用的是短连接。也就是说，浏览器和服务器每进行一次HTTP操作，就建立一次连接，但任务结束就中断连接。如果客户端浏览器访问的某个HTML或其他类型的 Web页中包含有其他的Web资源，如JavaScript文件、图像文件、CSS文件等；当浏览器每遇到这样一个Web资源，就会建立一个HTTP会话。

但从 HTTP/1.1起，默认使用长连接，用以保持连接特性。使用长连接的HTTP协议，会在响应头有加入这行代码：

Connection:keep-alive

在使用长连接的情况下，当一个网页打开完成后，客户端和服务器之间用于传输HTTP数据的 TCP连接不会关闭，如果客户端再次访问这个服务器上的网页，会继续使用这一条已经建立的连接。Keep-Alive不会永久保持连接，它有一个保持时间，可以在不同的服务器软件（如Apache）中设定这个时间。实现长连接要客户端和服务端都支持长连接。

HTTP协议的长连接和短连接，实质上是TCP协议的长连接和短连接。

三、TCP长连接和短连接：
我们模拟一下TCP短连接的情况，client向server发起连接请求，server接到请求，然后双方建立连接。client向server 发送消息，server回应client，然后一次读写就完成了，这时候双方任何一个都可以发起close操作，不过一般都是client先发起 close操作。为什么呢，一般的server不会回复完client后立即关闭连接的，当然不排除有特殊的情况。从上面的描述看，短连接一般只会在 client/server间传递一次读写操作

短连接的优点是：管理起来比较简单，存在的连接都是有用的连接，不需要额外的控制手段

接下来我们再模拟一下长连接的情况，client向server发起连接，server接受client连接，双方建立连接。Client与server完成一次读写之后，它们之间的连接并不会主动关闭，后续的读写操作会继续使用这个连接。

首先说一下TCP/IP详解上讲到的TCP保活功能，保活功能主要为服务器应用提供，服务器应用希望知道客户主机是否崩溃，从而可以代表客户使用资源。如果客户已经消失，使得服务器上保留一个半开放的连接，而服务器又在等待来自客户端的数据，则服务器将应远等待客户端的数据，保活功能就是试图在服务 器端检测到这种半开放的连接。

四、长连接和短连接的生命周期
短连接在建立连接后，完成一次读写就会自动关闭了。

正常情况下，一条TCP长连接建立后，只要双不提出关闭请求并且不出现异常情况，这条连接是一直存在的，操作系统不会自动去关闭它，甚至经过物理网络拓扑的改变之后仍然可以使用。所以一条连接保持几天、几个月、几年或者更长时间都有可能，只要不出现异常情况或由用户（应用层）主动关闭。

在编程中，往往需要建立一条TCP连接，并且长时间处于连接状态。所谓的TCP长连接并没有确切的时间限制，而是说这条连接需要的时间比较长。

五、怎样维护长连接或者检测中断
1、在应用层使用heartbeat来主动检测。
对于实时性要求较高的网络通信程序，往往需要更加及时的获取已经中断的连接，从而进行及时的处理。但如果对方的连接异常中断，往往是不能及时的得到对方连接已经中断的信息，操作系统检测连接是否中断的时间间隔默认是比较长的，即便它能够检测到，但却不符合我们的实时性需求，所以需要我们进行手工去不断探测。

2、改变socket的keepalive选项，以使socket检测连接是否中断的时间间隔更小，以满足我们的及时性需求。有关的几个选项使用和解析如下：
A、我们在检测对端以一种非优雅的方式断开连接的时候，可以设置SO_KEEPALIVE属性使得我们在2小时以后发现对方的TCP连接是否依然存在。用法如下：

keepAlive = 1；

setsockopt(listenfd, SOL_SOCKET, SO_KEEPALIVE, (void*)&keepAlive, sizeof(keepAlive));

B、如果我们不想使用这么长的等待时间，可以修改内核关于网络方面的配置参数，也可设置SOCKET的TCP层（SOL_TCP）选项TCP_KEEPIDLE、TCP_KEEPINTVL和TCP_KEEPCNT。

TCP_KEEPIDLE：开始首次KeepAlive探测前的TCP空闭时间

The tcp_keepidle parameter specifies the interval of inactivity that causes TCP to generate a KEEPALIVE transmission for an application that requests them. tcp_keepidle defaults to 14400 (two hours).

TCP_KEEPINTVL：两次KeepAlive探测间的时间间隔

The tcp_keepintvl parameter specifies the interval between the nine retries that are attempted if a KEEPALIVE transmission is not acknowledged. tcp_keepintvl defaults to 150 (75 seconds).

TCP_KEEPCNT：断开前的KeepAlive探测次数

The TCP_KEEPCNT option specifies the maximum number of keepalive probes to be sent. The value of TCP_KEEPCNT is an integer value between 1 and n, where n is the value of the systemwide tcp_keepcnt parameter.

如果心搏函数要维护客户端的存活，即服务器必须每隔一段时间必须向客户段发送一定的数据，那么使用SO_KEEPALIVE是有很大的不足的。因为SO_KEEPALIVE选项指"此套接口的任一方向都没有数据交换"。在Linux 2.6系列上，上面话的理解是只要打开SO_KEEPALIVE选项的套接口端检测到数据发送或者数据接受就认为是数据交换。因此在这种情况下使用 SO_KEEPALIVE选项 检测对方是否非正常连接是完全没有作用的，在每隔一段时间发包的情况， keep-alive的包是不可能被发送的。上层程序在非正常断开的情况下是可以正常发送包到缓冲区的。非正常端开的情况是指服务器没有收到"FIN" 或者 "RST"包。

什么时候用长连接，短连接？
长连接多用于操作频繁，点对点的通讯，而且连接数不能太多情况。
每个TCP连接都需要三步握手，这需要时间，如果每个操作都是先连接，再操作的话那么处理速度会降低很多，所以每个操作完后都不断开，次处理时直接发送数据包就OK了，不用建立TCP连接。例如：数据库的连接用长连接， 如果用短连接频繁的通信会造成socket错误，而且频繁的socket 创建也是对资源的浪费。

而像WEB网站的http服务一般都用短链接，因为长连接对于服务端来说会耗费一定的资源，而像WEB网站这么频繁的成千上万甚至上亿客户端的连接用短连接会更省一些资源，如果用长连接，而且同时有成千上万的用户，如果每个用户都占用一个连接的话，那可想而知吧。所以并发量大，但每个用户无需频繁操作情况下需用短连好。

总之，长连接和短连接的选择要视情况而定。

具体网络中的应用的话：

http 1.0一般就指短连接，smtp,pop3,telnet这种就可以认为是长连接。
一般的网络游戏应用都是长连接

https://blog.csdn.net/qq_41181857/article/details/107729857

但从 HTTP/1.1起，默认使用长连接，用以保持连接特性。使用长连接的HTTP协议，会在响应头有加入这行代码：

Connection:keep-alive 
在使用长连接的情况下，当一个网页打开完成后，客户端和服务器之间用于传输HTTP数据的 TCP连接不会关闭，如果客户端再次访问这个服务器上的网页，会继续使用这一条已经建立的连接。Keep-Alive不会永久保持连接，它有一个保持时间，可以在不同的服务器软件（如Apache）中设定这个时间。实现长连接要客户端和服务端都支持长连接。

HTTP协议的长连接和短连接，实质上是TCP协议的长连接和短连接。

三、TCP长连接和短连接：
我们模拟一下TCP短连接的情况，client向server发起连接请求，server接到请求，然后双方建立连接。client向server 发送消息，server回应client，然后一次读写就完成了，这时候双方任何一个都可以发起close操作，不过一般都是client先发起 close操作。为什么呢，一般的server不会回复完client后立即关闭连接的，当然不排除有特殊的情况。从上面的描述看，短连接一般只会在 client/server间传递一次读写操作

短连接的优点是：管理起来比较简单，存在的连接都是有用的连接，不需要额外的控制手段

接下来我们再模拟一下长连接的情况，client向server发起连接，server接受client连接，双方建立连接。Client与server完成一次读写之后，它们之间的连接并不会主动关闭，后续的读写操作会继续使用这个连接。

首先说一下TCP/IP详解上讲到的TCP保活功能，保活功能主要为服务器应用提供，服务器应用希望知道客户主机是否崩溃，从而可以代表客户使用资源。如果客户已经消失，使得服务器上保留一个半开放的连接，而服务器又在等待来自客户端的数据，则服务器将应远等待客户端的数据，保活功能就是试图在服务 器端检测到这种半开放的连接。

四、长连接和短连接的生命周期
短连接在建立连接后，完成一次读写就会自动关闭了。

正常情况下，一条TCP长连接建立后，只要双不提出关闭请求并且不出现异常情况，这条连接是一直存在的，操作系统不会自动去关闭它，甚至经过物理网络拓扑的改变之后仍然可以使用。所以一条连接保持几天、几个月、几年或者更长时间都有可能，只要不出现异常情况或由用户（应用层）主动关闭。

在编程中，往往需要建立一条TCP连接，并且长时间处于连接状态。所谓的TCP长连接并没有确切的时间限制，而是说这条连接需要的时间比较长。

五、怎样维护长连接或者检测中断
1、在应用层使用heartbeat来主动检测。
对于实时性要求较高的网络通信程序，往往需要更加及时的获取已经中断的连接，从而进行及时的处理。但如果对方的连接异常中断，往往是不能及时的得到对方连接已经中断的信息，操作系统检测连接是否中断的时间间隔默认是比较长的，即便它能够检测到，但却不符合我们的实时性需求，所以需要我们进行手工去不断探测。

探测的方式有两种：

2、改变socket的keepalive选项，以使socket检测连接是否中断的时间间隔更小，以满足我们的及时性需求。有关的几个选项使用和解析如下：
A、我们在检测对端以一种非优雅的方式断开连接的时候，可以设置SO_KEEPALIVE属性使得我们在2小时以后发现对方的TCP连接是否依然存在。用法如下：

keepAlive = 1；

setsockopt(listenfd, SOL_SOCKET, SO_KEEPALIVE, (void*)&keepAlive, sizeof(keepAlive));

B、如果我们不想使用这么长的等待时间，可以修改内核关于网络方面的配置参数，也可设置SOCKET的TCP层（SOL_TCP）选项TCP_KEEPIDLE、TCP_KEEPINTVL和TCP_KEEPCNT。

TCP_KEEPIDLE：开始首次KeepAlive探测前的TCP空闭时间

The tcp_keepidle parameter specifies the interval of inactivity that causes TCP to generate a KEEPALIVE transmission for an application that requests them. tcp_keepidle defaults to 14400 (two hours).

TCP_KEEPINTVL：两次KeepAlive探测间的时间间隔

The tcp_keepintvl parameter specifies the interval between the nine retries that are attempted if a KEEPALIVE transmission is not acknowledged. tcp_keepintvl defaults to 150 (75 seconds).

TCP_KEEPCNT：断开前的KeepAlive探测次数

The TCP_KEEPCNT option specifies the maximum number of keepalive probes to be sent. The value of TCP_KEEPCNT is an integer value between 1 and n, where n is the value of the systemwide tcp_keepcnt parameter. 

如果心搏函数要维护客户端的存活，即服务器必须每隔一段时间必须向客户段发送一定的数据，那么使用SO_KEEPALIVE是有很大的不足的。因为SO_KEEPALIVE选项指"此套接口的任一方向都没有数据交换"。在Linux 2.6系列上，上面话的理解是只要打开SO_KEEPALIVE选项的套接口端检测到数据发送或者数据接受就认为是数据交换。因此在这种情况下使用 SO_KEEPALIVE选项 检测对方是否非正常连接是完全没有作用的，在每隔一段时间发包的情况， keep-alive的包是不可能被发送的。上层程序在非正常断开的情况下是可以正常发送包到缓冲区的。非正常端开的情况是指服务器没有收到"FIN" 或者 "RST"包。

https://www.cnblogs.com/gotodsp/p/6366163.html

https://segmentfault.com/a/1190000015122195

http://www.52im.net/thread-331-1-1.html
这是一个http get请求报文，注意该报文中有一个upgrade首部，它的作用是告诉服务端需要将通信协议切换到websocket,如果服务端支持websocket协议，那么它就会将自己的通信协议切换到websocket,同时发给客户端类似于以下的一个响应报文头：

返回的状态码为101，表示同意客户端协议转换请求，并将它转换为websocket协议。以上过程都是利用http通信完成的，称之为websocket协议握手(websocket Protocol handshake)，进过这握手之后，客户端和服务端就建立了websocket连接，以后的通信走的都是websocket协议了。所以总结为websocket握手需要借助于http协议，建立连接后通信过程使用websocket协议。同时需要了解的是，该websocket连接还是基于我们刚才发起http连接的那个TCP连接。一旦建立连接之后，我们就可以进行数据传输了，websocket提供两种数据传输：文本数据和二进制数据。

基于以上分析，我们可以看到，websocket能够提供低延迟，高性能的客户端与服务端的双向数据通信。它颠覆了之前web开发的请求处理响应模式，并且提供了一种真正意义上的客户端请求，服务器推送数据的模式，特别适合实时数据交互应用开发。

http://www.52im.net/thread-331-1-1.html

w3c规范中关于HTML5 websocket API的原生API，这些api很简单，就是利用new WebSocket创建一个指定连接服务端地址的ws实例，然后为该实例注册onopen(连接服务端),onmessage(接受服务端数据)，onclose(关闭连接)以及ws.send(建立连接后)发送请求。上面说了那么多，事实上可以看到html5 websocket API本身是很简单的一个对象和它的几个方法而已。

https://blog.csdn.net/zhaohong_bo/article/details/89505528

Websocket其实是一个新协议，跟HTTP协议基本没有关系，只是为了兼容现有浏览器的握手规范而已，也就是说它是HTTP协议上的一种补充
作者：Ovear
链接：https://www.zhihu.com/question/20215561/answer/40316953
来源：知乎
著作权归作者所有。商业转载请联系作者获得授权，非商业转载请注明出处。

另外Html5是指的一系列新的API，或者说新规范，新技术。Http协议本身只有1.0和1.1，而且跟Html本身没有直接关系。。通俗来说，你可以用HTTP协议传输非Html数据，就是这样=。=再简单来说，层级不一样。二、Websocket是什么样的协议，具体有什么优点首先，Websocket是一个持久化的协议，相对于HTTP这种非持久的协议来说。简单的举个例子吧，用目前应用比较广泛的PHP生命周期来解释。1) HTTP的生命周期通过Request来界定，也就是一个Request 一个Response，那么在HTTP1.0中，这次HTTP请求就结束了。在HTTP1.1中进行了改进，使得有一个keep-alive，也就是说，在一个HTTP连接中，可以发送多个Request，接收多个Response。但是请记住 Request = Response ， 在HTTP中永远是这样，也就是说一个request只能有一个response。而且这个response也是被动的，不能主动发起。教练，你BB了这么多，跟Websocket有什么关系呢？_(:з」∠)_好吧，我正准备说Websocket呢。。首先Websocket是基于HTTP协议的，或者说借用了HTTP的协议来完成一部分握手。在握手阶段是一样的-------以下涉及专业技术内容，不想看的可以跳过lol:，或者只看加黑内容--------首先我们来看个典型的Websocket握手（借用Wikipedia的。。）GET /chat HTTP/1.1
Host: server.example.com
Upgrade: websocket
Connection: Upgrade
Sec-WebSocket-Key: x3JJHMbDL1EzLkh9GBhXDw==
Sec-WebSocket-Protocol: chat, superchat
Sec-WebSocket-Version: 13
Origin: http://example.com熟悉HTTP的童鞋可能发现了，这段类似HTTP协议的握手请求中，多了几个东西。我会顺便讲解下作用。Upgrade: websocket
Connection: Upgrade
这个就是Websocket的核心了，告诉Apache、Nginx等服务器：注意啦，窝发起的是Websocket协议，快点帮我找到对应的助理处理~不是那个老土的HTTP。Sec-WebSocket-Key: x3JJHMbDL1EzLkh9GBhXDw==
Sec-WebSocket-Protocol: chat, superchat
Sec-WebSocket-Version: 13首先，Sec-WebSocket-Key 是一个Base64 encode的值，这个是浏览器随机生成的，告诉服务器：泥煤，不要忽悠窝，我要验证尼是不是真的是Websocket助理。然后，Sec_WebSocket-Protocol 是一个用户定义的字符串，用来区分同URL下，不同的服务所需要的协议。简单理解：今晚我要服务A，别搞错啦~最后，Sec-WebSocket-Version 是告诉服务器所使用的Websocket Draft（协议版本），在最初的时候，Websocket协议还在 Draft 阶段，各种奇奇怪怪的协议都有，而且还有很多期奇奇怪怪不同的东西，什么Firefox和Chrome用的不是一个版本之类的，当初Websocket协议太多可是一个大难题。。不过现在还好，已经定下来啦~大家都使用的一个东西~ 脱水：服务员，我要的是13岁的噢→_→然后服务器会返回下列东西，表示已经接受到请求， 成功建立Websocket啦！HTTP/1.1 101 Switching Protocols
Upgrade: websocket
Connection: Upgrade
Sec-WebSocket-Accept: HSmrc0sMlYUkAGmm5OPpG2HaGWk=
Sec-WebSocket-Protocol: chat这里开始就是HTTP最后负责的区域了，告诉客户，我已经成功切换协议啦~Upgrade: websocket
Connection: Upgrade依然是固定的，告诉客户端即将升级的是Websocket协议，而不是mozillasocket，lurnarsocket或者shitsocket。然后，Sec-WebSocket-Accept 这个则是经过服务器确认，并且加密过后的 Sec-WebSocket-Key。服务器：好啦好啦，知道啦，给你看我的ID CARD来证明行了吧。。后面的，Sec-WebSocket-Protocol 则是表示最终使用的协议。至此，HTTP已经完成它所有工作了，接下来就是完全按照Websocket协议进行了。具体的协议就不在这阐述了。

作者：Ovear
链接：https://www.zhihu.com/question/20215561/answer/40316953
来源：知乎
著作权归作者所有。商业转载请联系作者获得授权，非商业转载请注明出处。

三、Websocket的作用在讲Websocket之前，我就顺带着讲下 long poll 和 ajax轮询 的原理。首先是 ajax轮询 ，ajax轮询 的原理非常简单，让浏览器隔个几秒就发送一次请求，询问服务器是否有新信息。场景再现：客户端：啦啦啦，有没有新信息(Request)服务端：没有（Response）客户端：啦啦啦，有没有新信息(Request)服务端：没有。。（Response）客户端：啦啦啦，有没有新信息(Request)服务端：你好烦啊，没有啊。。（Response）客户端：啦啦啦，有没有新消息（Request）服务端：好啦好啦，有啦给你。（Response）客户端：啦啦啦，有没有新消息（Request）服务端：。。。。。没。。。。没。。。没有（Response） ---- looplong poll long poll 其实原理跟 ajax轮询 差不多，都是采用轮询的方式，不过采取的是阻塞模型（一直打电话，没收到就不挂电话），也就是说，客户端发起连接后，如果没消息，就一直不返回Response给客户端。直到有消息才返回，返回完之后，客户端再次建立连接，周而复始。场景再现客户端：啦啦啦，有没有新信息，没有的话就等有了才返回给我吧（Request）服务端：额。。   等待到有消息的时候。。来 给你（Response）客户端：啦啦啦，有没有新信息，没有的话就等有了才返回给我吧（Request） -loop从上面可以看出其实这两种方式，都是在不断地建立HTTP连接，然后等待服务端处理，可以体现HTTP协议的另外一个特点，被动性。何为被动性呢，其实就是，服务端不能主动联系客户端，只能有客户端发起。简单地说就是，服务器是一个很懒的冰箱（这是个梗）（不会、不能主动发起连接），但是上司有命令，如果有客户来，不管多么累都要好好接待。说完这个，我们再来说一说上面的缺陷（原谅我废话这么多吧OAQ）从上面很容易看出来，不管怎么样，上面这两种都是非常消耗资源的。ajax轮询 需要服务器有很快的处理速度和资源。（速度）long poll 需要有很高的并发，也就是说同时接待客户的能力。（场地大小）所以ajax轮询 和long poll 都有可能发生这种情况。客户端：啦啦啦啦，有新信息么？服务端：月线正忙，请稍后再试（503 Server Unavailable）客户端：。。。。好吧，啦啦啦，有新信息么？服务端：月线正忙，请稍后再试（503 Server Unavailable）客户端：<img src="https://pic1.zhimg.com/50/7c0cf075c7ee4cc6cf52f4572a4c1c10_hd.jpg?source=1940ef5c" data-rawwidth="143" data-rawheight="50" class="content_image" width="143"/>然后服务端在一旁忙的要死：冰箱，我要更多的冰箱！更多。。更多。。（我错了。。这又是梗。。）--------------------------言归正传，我们来说Websocket吧通过上面这个例子，我们可以看出，这两种方式都不是最好的方式，需要很多资源。一种需要更快的速度，一种需要更多的'电话'。这两种都会导致'电话'的需求越来越高。哦对了，忘记说了HTTP还是一个无状态协议。（感谢评论区的各位指出OAQ）通俗的说就是，服务器因为每天要接待太多客户了，是个健忘鬼，你一挂电话，他就把你的东西全忘光了，把你的东西全丢掉了。你第二次还得再告诉服务器一遍。所以在这种情况下出现了，Websocket出现了。他解决了HTTP的这几个难题。首先，被动性，当服务器完成协议升级后（HTTP->Websocket），服务端就可以主动推送信息给客户端啦。所以上面的情景可以做如下修改。客户端：啦啦啦，我要建立Websocket协议，需要的服务：chat，Websocket协议版本：17（HTTP Request）服务端：ok，确认，已升级为Websocket协议（HTTP Protocols Switched）客户端：麻烦你有信息的时候推送给我噢。。服务端：ok，有的时候会告诉你的。服务端：balabalabalabala服务端：balabalabalabala服务端：哈哈哈哈哈啊哈哈哈哈服务端：笑死我了哈哈哈哈哈哈哈就变成了这样，只需要经过一次HTTP请求，就可以做到源源不断的信息传送了。（在程序设计中，这种设计叫做回调，即：你有信息了再来通知我，而不是我傻乎乎的每次跑来问你）这样的协议解决了上面同步有延迟，而且还非常消耗资源的这种情况。那么为什么他会解决服务器上消耗资源的问题呢？其实我们所用的程序是要经过两层代理的，即HTTP协议在Nginx等服务器的解析下，然后再传送给相应的Handler（PHP等）来处理。简单地说，我们有一个非常快速的接线员（Nginx），他负责把问题转交给相应的客服（Handler）。本身接线员基本上速度是足够的，但是每次都卡在客服（Handler）了，老有客服处理速度太慢。，导致客服不够。Websocket就解决了这样一个难题，建立后，可以直接跟接线员建立持久连接，有信息的时候客服想办法通知接线员，然后接线员在统一转交给客户。这样就可以解决客服处理速度过慢的问题了。同时，在传统的方式上，要不断的建立，关闭HTTP协议，由于HTTP是非状态性的，每次都要重新传输identity info（鉴别信息），来告诉服务端你是谁。虽然接线员很快速，但是每次都要听这么一堆，效率也会有所下降的，同时还得不断把这些信息转交给客服，不但浪费客服的处理时间，而且还会在网路传输中消耗过多的流量/时间。但是Websocket只需要一次HTTP握手，所以说整个通讯过程是建立在一次连接/状态中，也就避免了HTTP的非状态性，服务端会一直知道你的信息，直到你关闭请求，这样就解决了接线员要反复解析HTTP协议，还要查看identity info的信息。同时由客户主动询问，转换为服务器（推送）有信息的时候就发送（当然客户端还是等主动发送信息过来的。。），没有信息的时候就交给接线员（Nginx），不需要占用本身速度就慢的客服（Handler）了--------------------至于怎么在不支持Websocket的客户端上使用Websocket。。答案是：不能但是可以通过上面说的 long poll 和 ajax 轮询来 模拟出类似的效果

作者：董可人
链接：https://www.zhihu.com/question/20215561/answer/40250050
来源：知乎
著作权归作者所有。商业转载请联系作者获得授权，非商业转载请注明出处。

你可以把 WebSocket 看成是 HTTP 协议为了支持长连接所打的一个大补丁，它和 HTTP 有一些共性，是为了解决 HTTP 本身无法解决的某些问题而做出的一个改良设计。在以前 HTTP 协议中所谓的 keep-alive connection 是指在一次 TCP 连接中完成多个 HTTP 请求，但是对每个请求仍然要单独发 header；所谓的 polling 是指从客户端（一般就是浏览器）不断主动的向服务器发 HTTP 请求查询是否有新数据。这两种模式有一个共同的缺点，就是除了真正的数据部分外，服务器和客户端还要大量交换 HTTP header，信息交换效率很低。它们建立的“长连接”都是伪.长连接，只不过好处是不需要对现有的 HTTP server 和浏览器架构做修改就能实现。WebSocket 解决的第一个问题是，通过第一个 HTTP request 建立了 TCP 连接之后，之后的交换数据都不需要再发 HTTP request了，使得这个长连接变成了一个真.长连接。但是不需要发送 HTTP header就能交换数据显然和原有的 HTTP 协议是有区别的，所以它需要对服务器和客户端都进行升级才能实现。在此基础上 WebSocket 还是一个双通道的连接，在同一个 TCP 连接上既可以发也可以收信息。此外还有 multiplexing 功能，几个不同的 URI 可以复用同一个 WebSocket 连接。这些都是原来的 HTTP 不能做到的。另外说一点技术细节，因为看到有人提问 WebSocket 可能进入某种半死不活的状态。这实际上也是原有网络世界的一些缺陷性设计。上面所说的 WebSocket 真.长连接虽然解决了服务器和客户端两边的问题，但坑爹的是网络应用除了服务器和客户端之外，另一个巨大的存在是中间的网络链路。一个 HTTP/WebSocket 连接往往要经过无数的路由，防火墙。你以为你的数据是在一个“连接”中发送的，实际上它要跨越千山万水，经过无数次转发，过滤，才能最终抵达终点。在这过程中，中间节点的处理方法很可能会让你意想不到。比如说，这些坑爹的中间节点可能会认为一份连接在一段时间内没有数据发送就等于失效，它们会自作主张的切断这些连接。在这种情况下，不论服务器还是客户端都不会收到任何提示，它们只会一厢情愿的以为彼此间的红线还在，徒劳地一边又一边地发送抵达不了彼岸的信息。而计算机网络协议栈的实现中又会有一层套一层的缓存，除非填满这些缓存，你的程序根本不会发现任何错误。这样，本来一个美好的 WebSocket 长连接，就可能在毫不知情的情况下进入了半死不活状态。而解决方案，WebSocket 的设计者们也早已想过。就是让服务器和客户端能够发送 Ping/Pong Frame（RFC 6455 - The WebSocket Protocol）。这种 Frame 是一种特殊的数据包，它只包含一些元数据而不需要真正的 Data Payload，可以在不影响 Application 的情况下维持住中间网络的连接状态。

https://www.zhihu.com/question/20215561

https://tools.ietf.org/html/rfc6455#section-5.5.2

WebSocket 与 HTTP
WebSocket 协议在2008年诞生，2011年成为国际标准。现在所有浏览器都已经支持了。WebSocket 的最大特点就是，服务器可以主动向客户端推送信息，客户端也可以主动向服务器发送信息，是真正的双向平等对话。

HTTP 有 1.1 和 1.0 之说，也就是所谓的 keep-alive ，把多个 HTTP 请求合并为一个，但是 Websocket 其实是一个新协议，跟 HTTP 协议基本没有关系，只是为了兼容现有浏览器，所以在握手阶段使用了 HTTP 。

https://www.cnblogs.com/nnngu/p/9347635.html


https://www.jianshu.com/p/3444ea70b6cb

WebSocket的实现原理
一、什么是websocket
Websocket是应用层第七层上的一个应用层协议，它必须依赖 HTTP 协议进行一次握手 ，握手成功后，数据就直接从 TCP 通道传输，与 HTTP 无关了。即：websocket分为握手和数据传输阶段，即进行了HTTP握手 + 双工的TCP连接。

下面我们分别来看一下这两个阶段的具体实现原理：

１、握手阶段
客户端发送消息：

GET /chat HTTP/1.1
Host: server.example.com
Upgrade: websocket
Connection: Upgrade
Sec-WebSocket-Key: dGhlIHNhbXBsZSBub25jZQ==
Origin: http://example.com
Sec-WebSocket-Version: 13
服务端返回消息：

HTTP/1.1 101 Switching Protocols
Upgrade: websocket
Connection: Upgrade
Sec-WebSocket-Accept: s3pPLMBiTxaQ9kYGzzhZRbK+xOo=
这里值得注意的是Sec-WebSocket-Accept的计算方法：
base64(hsa1(sec-websocket-key + 258EAFA5-E914-47DA-95CA-C5AB0DC85B11))
如果这个Sec-WebSocket-Accept计算错误浏览器会提示：Sec-WebSocket-Accept dismatch
如果返回成功，Websocket就会回调onopen事件

２、传输阶段
Websocket的数据传输是frame形式传输的，比如会将一条消息分为几个frame，按照先后顺序传输出去。这样做会有几个好处：

a、大数据的传输可以分片传输，不用考虑到数据大小导致的长度标志位不足够的情况。
b、和http的chunk一样，可以边生成数据边传递消息，即提高传输效率。
websocket传输使用的协议如下图：

1464250745194141.png
参数说明如下：

FIN：1位，用来表明这是一个消息的最后的消息片断，当然第一个消息片断也可能是最后的一个消息片断；

RSV1, RSV2, RSV3: 分别都是1位，如果双方之间没有约定自定义协议，那么这几位的值都必须为0,否则必须断掉WebSocket连接；

Opcode: 4位操作码，定义有效负载数据，如果收到了一个未知的操作码，连接也必须断掉，以下是定义的操作码：

*  %x0 表示连续消息片断
*  %x1 表示文本消息片断
*  %x2 表未二进制消息片断
*  %x3-7 为将来的非控制消息片断保留的操作码
*  %x8 表示连接关闭
*  %x9 表示心跳检查的ping
*  %xA 表示心跳检查的pong
*  %xB-F 为将来的控制消息片断的保留操作码
Mask: 1位，定义传输的数据是否有加掩码,如果设置为1,掩码键必须放在masking-key区域，客户端发送给服务端的所有消息，此位的值都是1；

Payload length: 传输数据的长度，以字节的形式表示：7位、7+16位、或者7+64位。如果这个值以字节表示是0-125这个范围，那这个值就表示传输数据的长度；如果这个值是126，则随后的两个字节表示的是一个16进制无符号数，用来表示传输数据的长度；如果这个值是127,则随后的是8个字节表示的一个64位无符合数，这个数用来表示传输数据的长度。多字节长度的数量是以网络字节的顺序表示。负载数据的长度为扩展数据及应用数据之和，扩展数据的长度可能为0,因而此时负载数据的长度就为应用数据的长度。

Masking-key: 0或4个字节，客户端发送给服务端的数据，都是通过内嵌的一个32位值作为掩码的；掩码键只有在掩码位设置为1的时候存在。

Payload data: (x+y)位，负载数据为扩展数据及应用数据长度之和。

Extension data: x位，如果客户端与服务端之间没有特殊约定，那么扩展数据的长度始终为0，任何的扩展都必须指定扩展数据的长度，或者长度的计算方式，以及在握手时如何确定正确的握手方式。如果存在扩展数据，则扩展数据就会包括在负载数据的长度之内。

Application data: y位，任意的应用数据，放在扩展数据之后，应用数据的长度=负载数据的长度-扩展数据的长度。

二、golang的websocket实现
我们了解了websocket的实现原理之后，请看以下golang的实现案例：

package main
import(
    "net"
    "log"
    "strings"
    "crypto/sha1"
    "io"
    "encoding/base64"
    "errors"
)
func main() {
    ln, err := net.Listen("tcp", ":8000")
    if err != nil {
        log.Panic(err)
    }
    for {
        conn, err := ln.Accept()
        if err != nil {
            log.Println("Accept err:", err)
        }
        for {
            handleConnection(conn)
        }
    }
}

func handleConnection(conn net.Conn) {
    content := make([]byte, 1024)
    _, err := conn.Read(content)
    log.Println(string(content))
    if err != nil {
        log.Println(err)
    }
    isHttp := false
    // 先暂时这么判断
    if string(content[0:3]) == "GET" {
        isHttp = true;
    }
    log.Println("isHttp:", isHttp)
    if isHttp {
        headers := parseHandshake(string(content))
        log.Println("headers", headers)
        secWebsocketKey := headers["Sec-WebSocket-Key"]
        // NOTE：这里省略其他的验证
        guid := "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"
        // 计算Sec-WebSocket-Accept
        h := sha1.New()
        log.Println("accept raw:", secWebsocketKey + guid)
        io.WriteString(h, secWebsocketKey + guid)
        accept := make([]byte, 28)
        base64.StdEncoding.Encode(accept, h.Sum(nil))
        log.Println(string(accept))
        response := "HTTP/1.1 101 Switching Protocols\r\n"
        response = response + "Sec-WebSocket-Accept: " + string(accept) + "\r\n"
        response = response + "Connection: Upgrade\r\n"
        response = response + "Upgrade: websocket\r\n\r\n"
        log.Println("response:", response)
        if lenth, err := conn.Write([]byte(response)); err != nil {
            log.Println(err)
        } else {
            log.Println("send len:", lenth)
        }
        wssocket := NewWsSocket(conn)
        for {
            data, err := wssocket.ReadIframe()
            if err != nil {
                log.Println("readIframe err:" , err)
            }
            log.Println("read data:", string(data))
            err = wssocket.SendIframe([]byte("good"))
            if err != nil {
                log.Println("sendIframe err:" , err)
            }
            log.Println("send data")
        }
    } else {
        log.Println(string(content))
        // 直接读取
    }
}

type WsSocket struct {
    MaskingKey []byte
    Conn net.Conn
}

func NewWsSocket(conn net.Conn) *WsSocket {
    return &WsSocket{Conn: conn}
}

func (this *WsSocket)SendIframe(data []byte) error {
    // 这里只处理data长度<125的
    if len(data) >= 125 {
        return errors.New("send iframe data error")
    }
    lenth := len(data)
    maskedData := make([]byte, lenth)
    for i := 0; i < lenth; i++ {
        if this.MaskingKey != nil {
            maskedData[i] = data[i] ^ this.MaskingKey[i % 4]
        } else {
            maskedData[i] = data[i]
        }
    }
    this.Conn.Write([]byte{0x81})
    var payLenByte byte
    if this.MaskingKey != nil && len(this.MaskingKey) != 4 {
        payLenByte = byte(0x80) | byte(lenth)
        this.Conn.Write([]byte{payLenByte})
        this.Conn.Write(this.MaskingKey)
    } else {
        payLenByte = byte(0x00) | byte(lenth)
        this.Conn.Write([]byte{payLenByte})
    }
    this.Conn.Write(data)
    return nil
}

func (this *WsSocket)ReadIframe() (data []byte, err error){
    err = nil
    //第一个字节：FIN + RSV1-3 + OPCODE
    opcodeByte := make([]byte, 1)
    this.Conn.Read(opcodeByte)
    FIN := opcodeByte[0] >> 7
    RSV1 := opcodeByte[0] >> 6 & 1
    RSV2 := opcodeByte[0] >> 5 & 1
    RSV3 := opcodeByte[0] >> 4 & 1
    OPCODE := opcodeByte[0] & 15
    log.Println(RSV1,RSV2,RSV3,OPCODE)

    payloadLenByte := make([]byte, 1)
    this.Conn.Read(payloadLenByte)
    payloadLen := int(payloadLenByte[0] & 0x7F)
    mask := payloadLenByte[0] >> 7
    if payloadLen == 127 {
        extendedByte := make([]byte, 8)
        this.Conn.Read(extendedByte)
    }
    maskingByte := make([]byte, 4)
    if mask == 1 {
        this.Conn.Read(maskingByte)
        this.MaskingKey = maskingByte
    }

    payloadDataByte := make([]byte, payloadLen)
    this.Conn.Read(payloadDataByte)
    log.Println("data:", payloadDataByte)
    dataByte := make([]byte, payloadLen)
    for i := 0; i < payloadLen; i++ {
        if mask == 1 {
            dataByte[i] = payloadDataByte[i] ^ maskingByte[i % 4]
        } else {
            dataByte[i] = payloadDataByte[i]
        }
    }
    if FIN == 1 {
        data = dataByte
        return
    }
    nextData, err := this.ReadIframe()
    if err != nil {
        return
    }
    data = append(data,  nextData...)
    return
}

func parseHandshake(content string) map[string]string {
    headers := make(map[string]string, 10)
    lines := strings.Split(content, "\r\n")
    for _,line := range lines {
        if len(line) >= 0 {
            words := strings.Split(line, ":")
            if len(words) == 2 {
                headers[strings.Trim(words[0]," ")] = strings.Trim(words[1], " ")
            }
        }
    }
    return headers
}
html和ｊs使用案例：


<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Title</title>
</head>
<body>
<html>
<head>
    <script type="text/javascript" src="./jquery.min.js"></script>
</head>
<body>
<input type="button" id="connect" value="websocket connect" />
<input type="button" id="send" value="websocket send" />
<input type="button" id="close" value="websocket close" />
</body>
<script type="text/javascript" src="./websocket.js"></script>
</html>
</body>
</html>

var socket;
$("#connect").click(function(event){
    socket = new WebSocket("ws://127.0.0.1:8000/chat");
    socket.onopen = function(){
    alert("Socket has been opened");
    }
    socket.onmessage = function(msg){
    alert(msg.data);
    }
    socket.onclose = function() {
    alert("Socket has been closed");
    }
});
$("#send").click(function(event){
    socket.send("send from client");
});
$("#close").click(function(event){
    socket.close();
})

https://zhuanlan.zhihu.com/p/32845970

https://blog.csdn.net/libaineu2004/article/details/81263689
https://github.com/owenliang/go-push
https://zhuanlan.zhihu.com/p/44711104
Golang 实现的连接池

功能：

* 连接池中连接类型为interface{}，使得更加通用

* 链接的最大空闲时间，超时的链接将关闭丢弃，可避免空闲时链接自动失效问题

* 使用channel处理池中的链接，高效

基本用法
//factory 创建连接的方法
factory := func() (interface{}, error) { return net.Dial("tcp", "127.0.0.1:4000") }
 
//close 关闭链接的方法
close := func(v interface{}) error { return v.(net.Conn).Close() }
 
//创建一个连接池： 初始化5，最大链接30
poolConfig := &pool.PoolConfig{
    InitialCap: 5,
    MaxCap:     30,
    Factory:    factory,
    Close:      close,
    //链接最大空闲时间，超过该时间的链接 将会关闭，可避免空闲时链接EOF，自动失效的问题
    IdleTimeout: 15 * time.Second,
}
p, err := pool.NewChannelPool(poolConfig)
if err != nil {
    fmt.Println("err=", err)
}
 
//从连接池中取得一个链接
v, err := p.Get()
 
//do something
//conn=v.(net.Conn)
 
//将链接放回连接池中
p.Put(v)
 
//释放连接池中的所有链接
p.Release()
 
//查看当前链接中的数量
current := p.Len()
https://github.com/silenceper/pool

在设计与实现连接池时，我们通常需要考虑以下几个问题：

 连接池的连接数目是否有限制，最大可以建立多少个连接？
 当连接长时间没有使用，需要回收该连接吗？
 业务请求需要获取连接时，此时若连接池无空闲连接且无法新建连接，业务需要排队等待吗？
 排队的话又存在另外的问题，队列长度有无限制，排队时间呢？
Golang连接池实现原理

我们以Golang HTTP连接池为例，分析连接池的实现原理。

结构体Transport

Transport结构定义如下：
type Transport struct {
  //操作空闲连接需要获取锁
  idleMu    sync.Mutex
  //空闲连接池，key为协议目标地址等组合
  idleConn   map[connectMethodKey][]*persistConn // most recently used at end
  //等待空闲连接的队列，基于切片实现，队列大小无限制
  idleConnWait map[connectMethodKey]wantConnQueue // waiting getConns
  
  //排队等待建立连接需要获取锁
  connsPerHostMu  sync.Mutex
  //每个host建立的连接数
  connsPerHost   map[connectMethodKey]int
  //等待建立连接的队列，同样基于切片实现，队列大小无限制
  connsPerHostWait map[connectMethodKey]wantConnQueue // waiting getConns
  
  //最大空闲连接数
  MaxIdleConns int
  //每个目标host最大空闲连接数；默认为2（注意默认值）
  MaxIdleConnsPerHost int
  //每个host可建立的最大连接数
  MaxConnsPerHost int
  //连接多少时间没有使用则被关闭
  IdleConnTimeout time.Duration
  
  //禁用长连接，使用短连接
  DisableKeepAlives bool
}
可以看到，连接护着队列，都是一个map结构，而key为协议目标地址等组合，即同一种协议与同一个目标host可建立的连接或者空闲连接是有限制的。

需要特别注意的是，MaxIdleConnsPerHost默认等于2，即与目标主机最多只维护两个空闲连接。这会导致什么呢？

如果遇到突发流量，瞬间建立大量连接，但是回收连接时，由于最大空闲连接数的限制，该联机不能进入空闲连接池，只能直接关闭。结果是，一直新建大量连接，又关闭大量连，业务机器的TIME_WAIT连接数随之突增。

线上有些业务架构是这样的：客户端 ===> LVS ===> Nginx ===> 服务。LVS负载均衡方案采用DR模式，LVS与Nginx配置统一VIP。此时在客户端看来，只有一个IP地址，只有一个Host。上述问题更为明显。

最后，Transport也提供了配置DisableKeepAlives，禁用长连接，使用短连接访问第三方资源或者服务。

连接获取与回收

Transport结构提供下面两个方法实现连接的获取与回收操作。
func (t *Transport) getConn(treq *transportRequest, cm connectMethod) (pc *persistConn, err error) {}
 
func (t *Transport) tryPutIdleConn(pconn *persistConn) error {}
连接的获取主要分为两步走：1）尝试获取空闲连接；2）尝试新建连接：
//getConn方法内部实现
 
if delivered := t.queueForIdleConn(w); delivered {
  return pc, nil
}
  
t.queueForDial(w)
当然，可能获取不到连接而需要排队，此时怎么办呢？当前会阻塞当前协程了，直到获取连接为止，或者httpclient超时取消请求：

select {
  case <-w.ready:
    return w.pc, w.err
    
  //超时被取消
  case <-req.Cancel:
    return nil, errRequestCanceledConn
  ……
}
 
var errRequestCanceledConn = errors.New("net/http: request canceled while waiting for connection") // TODO: unify?
排队等待空闲连接的逻辑如下：

func (t *Transport) queueForIdleConn(w *wantConn) (delivered bool) {
  //如果配置了空闲超时时间，获取到连接需要检测，超时则关闭连接
  if t.IdleConnTimeout > 0 {
    oldTime = time.Now().Add(-t.IdleConnTimeout)
  }
  
  if list, ok := t.idleConn[w.key]; ok {
    for len(list) > 0 && !stop {
      pconn := list[len(list)-1]
      tooOld := !oldTime.IsZero() && pconn.idleAt.Round(0).Before(oldTime)
      //超时了，关闭连接
      if tooOld {
        go pconn.closeConnIfStillIdle()
      }
      
      //分发连接到wantConn
      delivered = w.tryDeliver(pconn, nil)
    }
  }
  
  //排队等待空闲连接
  q := t.idleConnWait[w.key]
  q.pushBack(w)
  t.idleConnWait[w.key] = q
}
排队等待新建连接的逻辑如下：

func (t *Transport) queueForDial(w *wantConn) {
  //如果没有限制最大连接数，直接建立连接
  if t.MaxConnsPerHost <= 0 {
    go t.dialConnFor(w)
    return
  }
  
  //如果没超过连接数限制，直接建立连接
  if n := t.connsPerHost[w.key]; n < t.MaxConnsPerHost {
    go t.dialConnFor(w)
    return
  }
  
  //排队等待连接建立
  q := t.connsPerHostWait[w.key]
  q.pushBack(w)
  t.connsPerHostWait[w.key] = q
}
连接建立完成后，同样会调用tryDeliver分发连接到wantConn，同时关闭通道w.ready，这样主协程纠接触阻塞了。
func (w *wantConn) tryDeliver(pc *persistConn, err error) bool {
  w.pc = pc
  close(w.ready)
}
请求处理完成后，通过tryPutIdleConn将连接放回连接池；这时候如果存在等待空闲连接的协程，则需要分发复用该连接。另外，在回收连接时，还需要校验空闲连接数目是否超过限制：
func (t *Transport) tryPutIdleConn(pconn *persistConn) error {
  //禁用长连接；或者最大空闲连接数不合法
  if t.DisableKeepAlives || t.MaxIdleConnsPerHost < 0 {
    return errKeepAlivesDisabled
  }
  
  if q, ok := t.idleConnWait[key]; ok {
    //如果等待队列不为空，分发连接
    for q.len() > 0 {
      w := q.popFront()
      if w.tryDeliver(pconn, nil) {
        done = true
        break
      }
    }
  }
  
  //空闲连接数目超过限制，默认为DefaultMaxIdleConnsPerHost=2
  idles := t.idleConn[key]
  if len(idles) >= t.maxIdleConnsPerHost() {
    return errTooManyIdleHost
  }
 
}
空闲连接超时关闭

Golang HTTP连接池如何实现空闲连接的超时关闭逻辑呢？从上述queueForIdleConn逻辑可以看到，每次在获取到空闲连接时，都会检测是否已经超时，超时则关闭连接。

那如果没有业务请求到达，一直不需要获取连接，空闲连接就不会超时关闭吗？其实在将空闲连接添加到连接池时，Golang同时还设置了定时器，定时器到期后，自然会关闭该连接。
pconn.idleTimer = time.AfterFunc(t.IdleConnTimeout, pconn.closeConnIfStillIdle)
排队队列怎么实现

怎么实现队列模型呢？很简单，可以基于切片：

queue  []*wantConn
 
//入队
queue = append(queue, w)
 
//出队
v := queue[0]
queue[0] = nil
queue = queue[1:]
这样有什么问题吗？随着频繁的入队与出队操作，切片queue的底层数组，会有大量空间无法复用而造成浪费。除非该切片执行了扩容操作。

Golang在实现队列时，使用了两个切片head和tail；head切片用于出队操作，tail切片用于入队操作；出队时，如果head切片为空，则交换head与tail。通过这种方式，Golang实现了底层数组空间的复用。
func (q *wantConnQueue) pushBack(w *wantConn) {
  q.tail = append(q.tail, w)
}
 
func (q *wantConnQueue) popFront() *wantConn {
  if q.headPos >= len(q.head) {
    if len(q.tail) == 0 {
      return nil
    }
    // Pick up tail as new head, clear tail.
    q.head, q.headPos, q.tail = q.tail, 0, q.head[:0]
  }
  w := q.head[q.headPos]
  q.head[q.headPos] = nil
  q.headPos++
  return w
}

grpc长连接池 
https://github.com/0x5010/grpcp
https://www.jb51.net/article/193675.htm
https://www.cnblogs.com/smallleiit/articles/12632926.html
https://segmentfault.com/a/1190000013089363
https://github.com/goctx/generic-pool

https://www.dazhuanlan.com/2019/12/13/5df2e64445102/
http://www.manongjc.com/detail/5-gijopfmhpkglfmm.html
https://github.com/tRavAsty/fdfs_client

https://blog.csdn.net/fjslovejhl/article/details/50355691
http://www.voidcn.com/article/p-pwqvdjdt-pe.html

1. TCP 连接本身并没有长短的区分， 长或短只是在描述我们使用它的方式

2. 长/短是指多次数据交换能否复用同一个连接， 而不是指连接的持续时间

3. TCP 的 keepalive 仅起到保活探测的作用， 和连接的长短并没有因果关系

https://zhuanlan.zhihu.com/p/245693513?utm_source=wechat_session

长连接的优势
相比于短连接，长连接具有：

1. 较低的延时。 由于跳过了三次握手的过程，长连接比短连接有更低的延迟。

2. 较低的带宽占用。由于不用为每个请求建立和关闭连接，长连接交换效率更高，网络带宽占用更少。

3. 较少的系统资源占用。server 为了维持连接，会为每个连接创建 socket，分配文件句柄， 在内存中分配读写 buffer，设置定时器进行keepalive。 因此更少的连接数也意味着更少的资源占用。

另外， gRPC 使用 HTTP/2.0 作为传输协议， 从该协议的设计来讲， 长连接也是更推荐的使用方式， 原因如下：

1. HTTP/2.0 的多路复用， 使得连接的复用效率得到了质的提升。

HTTP/1.0 开始支持长连接， 如下图1， 请求会在 client 排队(request queuing)， 当响应返回之后再发送下一个请求。而这个过程中， 任何一个请求处理过慢都会阻塞整个流程， 这个问题被称为线头阻塞问题， 即 Head-of-line blocking。

HTTP/1.1 做出了改进， 允许client 可以连续发送多个请求， 但 server 的响应必须按照请求发送的顺序依次返回， 称为Pipelining (server 端响应排队)， 如下图2。这在一定程度上提高了复用效率， 但并没能解决线头阻塞的问题。

HTTP/2.0 引入了分帧分流的机制， 实现了多路复用(乱序发送乱序接受)， 彻底的解决了线头阻塞， 极大提高了连接复用的效率。

长连接不是银弹
虽然长连接有很多优势， 但并不是所有的场景都适用。 在使用长连接之前， 至少有以下两个点需要考虑。

1. client 和 server 的数量
长连接模式下， server 要和每一个 client都保持连接。 如果 client 数量远远超过 server 数量， 与每个 client 都维持一个长连接， 对 server 来说会是一个极大的负担。 好在这种场景中， 连接的利用率和复用率往往不高，使用简单且易于管理的短连接是更好的选择。 即使用长连接， 也必须设置一个合理的超时机制， 如在空闲时间过长时断开连接， 释放 server 资源。

2. 负载均衡机制
现代后端服务端架构中， 为了实现高可用和可伸缩， 一般都会引入单独的模块来提供负载均衡的功能， 称为负载均衡器。 根据工作在 OSI 不同的层级， 不同的负载均衡器会提供不同的转发功能。 接下来就最常见的 L4 (工作在TCP层）和 L7 (工作在应用层， 如 HTTP） 两种负载均衡器来分析。

L4负载均衡器:原理是将收到的 TCP 报文， 以一定的规则转发给后端的某一个 server。 这个转发规则其实是<client IP， Port> 到某个 server 地址的映射。 由于它只转发， 而不会进行报文解析， 因此这种场景下 client 会和 server 端握手后直接建立连接， 并且所有的数据报文都只会转发给同一个 server。 如下图所示， L4 会将 10.0.0.1:3001 的流量全部转发给 11.0.0.2:3110。

在短连接模式下， 由于连接会不断的建立和关闭， 同一个 client 的流量会被分发到不同的 server。
在长连接模式下， 由于连接一旦建立便不会断开， 就会导致流量会被分发到同一个 server。 在 client 与 server数量差距不大甚至client 少于 server 的情况下， 就会导致流量分发不均。 如下图中， 第三个 server 会一直处于空闲的状态。

为了避免这种场景中负载均衡失效的情况， L7 负载均衡器便成了一个更好的选择。

L7负载均衡器：相比 L4 只能基于连接进行负载均衡， L7 可以进行 HTTP 协议的解析。 当 client 发送请求的， client 会先和 L7 握手， L7 再和后端的一个或几个 server 握手， 进而实现基于请求的负载均衡。 如下图所示， 10.0.0.1 通过长连接发出的多个请求会根据 url， cookies 或 header 被 L7 分发到后端不同的 server。

https://www.cnblogs.com/gao88/p/12010917.html

nginx设置响应连接是长连接或者短连接
https://blog.csdn.net/qq_21127151/article/details/106880632

配置思路
长连接：

   http {
 ---------------------------
    keepalive_requests  100000;  //这里实际只需要大于1就可以
--------------------
    }
短连接：

   http {
 ---------------------------
    keepalive_requests  1;  //这里必须配置为1
--------------------
    }
验证是否配置成功
自己写客户端，每个链接发送多笔请求
通过curl 工具，发送请求,多个请求使用空格隔开，
curl http://10.9.2.111:80/ http://10.9.2.111:80/ http://10.9.2.111:80/
使用tcpdump抓取发往10.9.2.111的请求源端口是否改变，如果每一笔都变则为短连接；否则是长连接


HTTP1.1之后，HTTP协议支持持久连接，也就是长连接，优点在于在一个TCP连接上可以传送多个HTTP请求和响应，减少了建立和关闭连接的消耗和延迟。

如果我们使用了nginx去作为反向代理或者负载均衡，从客户端过来的长连接请求就会被转换成短连接发送给服务器端。

为了支持长连接，我们需要在nginx服务器上做一些配置。

   

·【要求】

使用nginx时，想要做到长连接，我们必须做到以下两点：

从client到nginx是长连接
从nginx到server是长连接
   

对于客户端而言，nginx其实扮演着server的角色，反之，之于server，nginx就是一个client。

   

·【保持和 Client 的长连接】

我们要想做到Client与Nginx之间保持长连接，需要：

Client发送过来的请求携带"keep-alive"header。
Nginx设置支持keep-alive
   

【HTTP配置】

默认情况下，nginx已经开启了对client连接的 keepalive 支持。对于特殊场景，可以调整相关参数。

http {

keepalive_timeout 120s;        #客户端链接超时时间。为0的时候禁用长连接。

keepalive_requests 10000;    #在一个长连接上可以服务的最大请求数目。

                                                  #当达到最大请求数目且所有已有请求结束后，连接被关闭。

                                                  #默认值为100

}

   

大多数情况下，keepalive_requests = 100也够用，但是对于 QPS 较高的场景，非常有必要加大这个参数，以避免出现大量连接被生成再抛弃的情况，减少TIME_WAIT。
 

QPS=10000 时，客户端每秒发送 10000 个请求 (通常建立有多个长连接)，每个连接只能最多跑 100 次请求，意味着平均每秒钟就会有 100 个长连接因此被 nginx 关闭。

同样意味着为了保持 QPS，客户端不得不每秒中重新新建 100 个连接。

   

因此，如果用netstat命令看客户端机器，就会发现有大量的TIME_WAIT的socket连接 (即使此时keep alive已经在 Client 和 NGINX 之间生效)。

   

·【保持和Server的长连接】

想让Nginx和Server之间维持长连接，最朴素的设置如下：

http {

upstream backend {

  server 192.168.0.1：8080 weight=1 max_fails=2 fail_timeout=30s;

  server 192.168.0.2：8080 weight=1 max_fails=2 fail_timeout=30s;

  keepalive 300; // 这个很重要！

}   

server {

listen 8080 default_server;

server_name "";

   

location / {

proxy_pass http://backend;

proxy_http_version 1.1;                         # 设置http版本为1.1

proxy_set_header Connection "";      # 设置Connection为长连接（默认为no）}

}

}

}

   

【upstream配置】

upstream中，有一个参数特别的重要，就是keepalive。

这个参数和之前http里面的 keepalive_timeout 不一样。

这个参数的含义是，连接池里面最大的空闲连接数量。

   

不理解？没关系，我们来举个例子：

场景：

有一个HTTP服务，作为upstream服务器接收请求，响应时间为100毫秒。

要求性能达到10000 QPS，我们需要在nginx与upstream服务器之间建立大概1000条HTTP请求。（1000/0.1s=10000）

   

最优情况：

假设请求非常的均匀平稳，每一个请求都是100ms，请求结束会被马上放入连接池并置为idle（空闲）状态。

我们以0.1s为单位：

1. 我们现在keepalive的值设置为10，每0.1s钟有1000个连接

2. 第0.1s的时候，我们一共有1000个请求收到并释放

3. 第0.2s的时候，我们又来了1000个请求，在0.2s结束的时候释放

   

请求和应答都比较均匀，0.1s释放的连接正好够用，不需要建立新连接，且连接池中没有idle状态的连接。

   

第一种情况：

应答非常平稳，但是请求不平稳的时候

4. 第0.3s的时候，我们只有500个请求收到，有500个请求因为网络延迟等原因没有进来

这个时候，Nginx检测到连接池中有500个idle状态的连接，就直接关闭了（500-10）个连接

5. 第0.4s的时候，我们收到了1500个请求，但是现在池里面只有（500+10）个连接，所以Nginx不得不重新建立了（1500-510）个连接。

如果在第4步的时候，没有关闭那490个连接的话，只需要重新建立500个连接。

   

第二种情况：

请求非常平稳，但是应答不平稳的时候

4. 第0.3s的时候，我们一共有1500个请求收到

但是池里面只有1000个连接，这个时候，Nginx又创建了500个连接，一共1500个连接

5. 第0.3s的时候，第0.3s的连接全部被释放，我们收到了500个请求

Nginx检测到池里面有1000个idle状态的连接，所以不得不释放了（1000-10）个连接

   

造成连接数量反复震荡的一个推手，就是这个keepalive 这个最大空闲连接数。

上面的两种情况说的都是 keepalive 设置的不合理导致Nginx有多次释放与创建连接的过程，造成资源浪费。

   

keepalive 这个参数设置一定要小心，尤其是对于 QPS 要求比较高或者网络环境不稳定的场景，一般根据 QPS 值和 平均响应时间能大致推算出需要的长连接数量。

然后将keepalive设置为长连接数量的10%到30%。

   

【location配置】

http {

server {

location / {

proxy_pass http://backend;

proxy_http_version 1.1;                         # 设置http版本为1.1

proxy_set_header Connection "";      # 设置Connection为长连接（默认为no）

}

}

}

HTTP 协议中对长连接的支持是从 1.1 版本之后才有的，因此最好通过 proxy_http_version 指令设置为 1.1。

HTTP1.0不支持keepalive特性，当没有使用HTTP1.1的时候，后端服务会返回101错误，然后断开连接。

   

而 "Connection" header 可以选择被清理，这样即便是 Client 和 Nginx 之间是短连接，Nginx 和 upstream 之间也是可以开启长连接的。

   

【另外一种高级方式】

http {

map $http_upgrade $connection_upgrade {

default upgrade;

'' close;

}   

upstream backend {

server 192.168.0.1：8080 weight=1 max_fails=2 fail_timeout=30s;

server 192.168.0.2：8080 weight=1 max_fails=2 fail_timeout=30s;

keepalive 300;

}   

server {

listen 8080 default_server;

server_name "";

location / {

proxy_pass http://backend;

   

proxy_connect_timeout 15;       #与upstream server的连接超时时间（没有单位，最大不可以超过75s）

proxy_read_timeout 60s;           #nginx会等待多长时间来获得请求的响应

proxy_send_timeout 12s;           #发送请求给upstream服务器的超时时间   

proxy_http_version 1.1;

proxy_set_header Upgrade $http_upgrade;

proxy_set_header Connection $connection_upgrade;

}

}

}

   

http里面的map的作用是：

让转发到代理服务器的 "Connection" 头字段的值，取决于客户端请求头的 "Upgrade" 字段值。

如果 $http_upgrade没有匹配，那 "Connection" 头字段的值会是upgrade。

如果 $http_upgrade为空字符串的话，那 "Connection" 头字段的值会是 close。

   

【补充】

NGINX支持WebSocket。

对于NGINX将升级请求从客户端发送到后台服务器，必须明确设置Upgrade和Connection标题。

这也算是上面情况所非常常用的场景。

HTTP的Upgrade协议头机制用于将连接从HTTP连接升级到WebSocket连接，Upgrade机制使用了Upgrade协议头和Connection协议头。

为了让Nginx可以将来自客户端的Upgrade请求发送到后端服务器，Upgrade和Connection的头信息必须被显式的设置。

   

【注意】

在nginx的配置文件中，如果当前模块中没有proxy_set_header的设置，则会从上级别继承配置。

继承顺序为：http, server, location。

   

如果在下一层使用proxy_set_header修改了header的值，则所有的header值都可能会发生变化，之前继承的所有配置将会被丢弃。

所以，尽量在同一个地方进行proxy_set_header，否则可能会有别的问题。

https://www.cnblogs.com/liufarui/p/11075630.html

用nginx做grpc反向代理，nginx到后端server不能维持长连接问题
问题描述
公司内部容器平台，接入层用nginx做LB，用户有grpc协议需求，所以在lb层支持grcp反向代理，nginx从1.13开始支持grpc反向代理，将公司使用的nginx包从1.12升级到1.14.0后，增加grpc反向代理配置。配置完成后，打压力测试时，发现接入层机器端口占满而导致服务异常，开始追查问题。

追查方向
深入了解grpc协议
gRPC是一个高性能、通用的开源 RPC 框架，其由 Google 主要面向移动应用开发并基于HTTP/2协议标准而设计，基于ProtoBuf(Protocol Buffers) 序列化协议开发，且支持众多开发语言。gRPC 提供了一种简单的方法来精确地定义服务和为 iOS、Android 和后台支持服务自动生成可靠性很强的客户端功能库。客户端充分利用高级流和链接功能，从而有助于节省带宽、降低的 TCP 链接次数、节省 CPU 使用、和电池寿命。
从上述描述可以看出grpc基于http2，client到server应该保持长连接，理论上不应该出现端口占满的问题

抓包
client到接入层抓包看请求状态，发现client到接入层确实是长连接状态。接入层到后端server抓包发现，请求并没有保持长连接，一个请求处理完之后，链接就断开了。

查nginx跟长连接相关配置
nginx长连接相关说明参考以下文档：nginx长连接
nginx跟grpc长连接相关配置，发现在1.15.6版本引入了"grpc_socket_keepalive"跟grpc直接相关长连接配置
Configures the “TCP keepalive” behavior for outgoing connections to a gRPC server. By default, the operating system’s settings are in effect for the socket. If the directive is set to the value “on”, the SO_KEEPALIVE socket option is turned on for the socket.

参考nginx长连接相关文档做了配置调整之后，发现nginx到server依然是短连接

将nginx升级到nginx1.15.6（升级过程中由于线上的nginx用到了lua模块，碰到了lua模块不适配问题，解决方案见链接lua模块适配问题)配置grpc_socket_keepalive on，抓包发现，会有少量处理多个请求的长连接存在，但大部分依然是短连接。

开启nginx debug模式

从debug日志来看nginx确实尝试重用链接，但是从实际抓包看，nginx的链接重用的情况非常少，大部分都是请求处理完之后链接断开，怀疑nginx对grpc反向代理支持的不够理想。

调整端口回收策略
回到问题本身，要解决的问题是接入层端口占满，各种调整nginx的长连接配置并不能解决这个问题，就尝试从tcp链接端口回收方面解决，大部分的tcp链接都处于TIME_WAIT状态。

TIME_WAIT状态的时间是2倍的MSL(linux里一个MSL为30s，是不可配置的)，在TIME_WAIT状态TCP连接实际上已经断掉，但是该端口又不能被新的连接实例使用。这种情况一般都是程序中建立了大量的短连接，而操作系统中对使用端口数量做了限制最大能支持65535个，TIME_WAIT过多非常容易出现连接数占满的异常。对TIME_WATI的优化有两个系统参数可配置：tcp_tw_reuse，tcp_tw_recycle 这两个参数在网上都有详细介绍，这是就不展开来讲

tcp_tw_reuse参考链接

tcp_tw_recycle（Enable fast recycling TIME-WAIT sockets）参考链接

测试来看，开启tcp_tw_reuse并没有解决端口被占满的问题，所以开启了更激进的tcp_tw_recycle，至此端口占用显著降低，问题解决，由于我们接入层并不是用的NAT，所以这个配置不会影响服务。

结论
通过测试发现nginx对grpc的反向代理支持的不够理想，从nginx到后端server大部分请求不能保持长连接
当用nginx做接入层反向代理时，对tcp参数调优可以避免端口占满等问题的发生

https://juejin.cn/post/6844903809534148622

https://skyao.gitbooks.io/learning-nginx/content/documentation/keep_alive.html

https://www.cnblogs.com/smileice/p/11156403.html

https://blog.csdn.net/bocai_xiaodaidai/article/details/90778227

https://www.cnblogs.com/arxive/p/7496489.html


https://www.rongcloud.cn/product/im?utm_source=baiduPay&utm_term=a0106014%23im%e7%b3%bb%e7%bb%9f%e5%ae%9e%e7%8e%b0&bd_vid=10957833802867095877#demo

https://blog.csdn.net/aaa31203/article/details/103917801
https://blog.csdn.net/Johnable/article/details/101760441

https://github.com/mpusher/mpush

https://www.cnblogs.com/crossoverJie/p/10206724.html

https://github.com/crossoverJie/cim

https://zhuanlan.zhihu.com/p/74832965
https://ifeve.com/%E4%B8%BA%E8%87%AA%E5%B7%B1%E6%90%AD%E5%BB%BA%E4%B8%80%E4%B8%AA%E5%88%86%E5%B8%83%E5%BC%8F-im%E5%8D%B3%E6%97%B6%E9%80%9A%E8%AE%AF-%E7%B3%BB%E7%BB%9F/

https://www.pianshen.com/article/63131302436/

https://blog.csdn.net/ytawu/article/details/85089068

https://studygolang.com/articles/17072?fr=sidebar
https://github.com/Terry-Ye/im
https://segmentfault.com/a/1190000017431165
https://studygolang.com/articles/19974

http://developer.gobelieve.io/

https://my.oschina.net/u/2397958/blog/467812

https://github.com/GoBelieveIO

https://www.cnblogs.com/parse-code/p/6160070.html

https://cloud.tencent.com/developer/article/1189548

https://github.com/lubanproj/grpc-read

https://www.bilibili.com/video/BV1VJ411a7K8?from=search&seid=2894220599736947308

https://github.com/grpc/grpc-go/

https://github.com/zhiqiangxu/qrpc

https://github.com/link1st/gowebsocket
