---
title: grpc websocket WebRTC
layout: post
category: golang
author: 夏泽民
---
从本质上讲,HTTP/2是在后台运行服务器推送的客户端/服务器,因此您可以发出请求,只需在该连接上继续监听更新,而无需轮询.

虽然WebSockets不会因为HTTP/2而消失,但对于以"让我知道何时更新发生与我刚刚做过的事情有关"的用例来说,它们可能不被认为是必要的.


<!-- more -->
https://qa.1r1g.com/sf/ask/3283327211/

	
HTTP/2

WebSocket

Headers

Compressed (HPACK)

None

Binary

Yes

Binary or Textual

Multiplexing

Yes

Yes

Prioritization

Yes

No

Compression

Yes

Yes

Direction

Client/Server + Server Push

Bidirectional

Full-duplex

Yes

Yes


https://www.infoq.com/articles/websocket-and-http2-coexist/

WebSocket 是一个双向通信协议，它在握手阶段采用 HTTP/1.1 协议（暂时不支持 HTTP/2）。

握手过程如下：

首先客户端向服务端发起一个特殊的 HTTP 请求，其消息头如下：
GET /chat HTTP/1.1  // 请求行
Host: server.example.com
Upgrade: websocket  // required
Connection: Upgrade // required
Sec-WebSocket-Key: dGhlIHNhbXBsZSBub25jZQ== // required，一个 16bits 编码得到的 base64 串
Origin: http://example.com  // 用于防止未认证的跨域脚本使用浏览器 websocket api 与服务端进行通信
Sec-WebSocket-Protocol: chat, superchat  // optional, 子协议协商字段
Sec-WebSocket-Version: 13
如果服务端支持该版本的 WebSocket，会返回 101 响应，响应标头如下：
HTTP/1.1 101 Switching Protocols  // 状态行
Upgrade: websocket   // required
Connection: Upgrade  // required
Sec-WebSocket-Accept: s3pPLMBiTxaQ9kYGzzhZRbK+xOo= // required，加密后的 Sec-WebSocket-Key
Sec-WebSocket-Protocol: chat // 表明选择的子协议
握手完成后，接下来的 TCP 数据包就都是 WebSocket 协议的帧了。

可以看到，这里的握手不是 TCP 的握手，而是在 TCP 连接内部，从 HTTP/1.1 upgrade 到 WebSocket 的握手。

WebSocket 提供两种协议：不加密的 ws:// 和 加密的 wss://. 因为是用 HTTP 握手，它和 HTTP 使用同样的端口：ws 是 80（HTTP），wss 是 443（HTTPS）

在 Python 编程中，可使用 websockets 实现的异步 WebSocket 客户端与服务端。此外 aiohttp 也提供了 WebSocket 支持。

Note：如果你搜索 Flask 的 WebScoket 插件，得到的第一个结果很可能是 Flask-SocketIO。但是 Flask-ScoektIO 使用的是它独有的 SocketIO 协议，并不是标准的 WebSocket。只是它刚好提供与 WebSocket 相同的功能而已。

SocketIO 的优势在于只要 Web 端使用了 SocketIO.js，就能支持该协议。而纯 WS 协议，只有较新的浏览器才支持。对于客户端非 Web 的情况，更好的选择可能是使用 Flask-Sockets。

https://www.cnblogs.com/kirito-c/p/10360309.html
https://www.infoq.com/articles/websocket-and-http2-coexist/

一、基础

1、HTTP协议是无状态的，服务器只会响应来自客户端的请求，但是它与客户端之间不具备持续连接；且只能从客户端主动请求服务端，服务端不能主动通知客户端。

     对于实时通信系统（聊天室或监控系统）这样显然是不合理的。传统的方法有：长轮询（客户端每隔很短的时间，都对服务器发出请求，当时间足够小就能实现实时的效果）、长连接（客户端只请求一次，但是服务器会将连接保持，当有数据时就返回结果给客户端）。这两种方式，都对客户端和服务器都造成了大量的性能浪费，于是WebSocket应运而生。WebSocket协议能够让浏览器和服务器全双工实时通信，互相的，服务器也能主动通知客户端。

2、 WebSocket的原理非常的简单：利用HTTP请求产生握手，HTTP头部中含有WebSocket协议的请求，所以握手之后，二者转用TCP协议进行交流。

       Socket.IO是业界良心，新手福音。它屏蔽了所有底层细节，让顶层调用非常简单。并且还为不支持WebSocket协议的浏览器，提供了长轮询的透明模拟机制。

https://www.cnblogs.com/george93/p/7513334.html


　　WebRTC
　　WebRTC（Web Real-Time Communication）。Real-Time Communication，实时通讯。

　　WebRTC能让web应用和站点之间选择性地分享音视频流。在不安装其它应用和插件的情况下，完成点对点通信。 WebRTC背后的技术被实现为一个开放的Web标准，并在所有主要浏览器中均以常规JavaScript API的形式提供。对于客户端（例如Android和iOS），可以使用提供相同功能的库。 WebRTC是个开源项目，得到Google，Apple，Microsoft和Mozilla等等公司的支持。2011年6月1日开源并在Google、Mozilla、Opera支持下被纳入万维网联盟的W3C推荐标准。

 

　　WebSocket
　　WebSocket是一种在单个TCP连接上进行全双工通信的协议。WebSocket通信协议于2011年被IETF定为标准RFC 6455，并由RFC7936补充规范。WebSocket API也被W3C定为标准。
　　WebSocket使得客户端和服务器之间的数据交换变得更加简单，允许服务端主动向客户端推送数据。在WebSocket API中，浏览器和服务器只需要完成一次握手，两者之间就直接可以创建持久性的连接，并进行双向数据传输。

https://www.cnblogs.com/huanzi-qch/p/15716286.html

webrtc是p2p,去中心化的。而web socket还是有中心(服务端)的，webrtc在建立信道过程中，依赖web socket来传输sdp

Websocket 是浏览器 和 服务器之间传输数据。

Webrtc 是 能支持浏览器之间进行数据传输，当然浏览器之间的连接建立要依赖服务器，连接建立的过程中使用websocket传输协商数据。

https://www.zhihu.com/question/424264607/answer/1512337034

https://developer.mozilla.org/zh-CN/docs/Web/API
WebSocket的主要问题是队首阻塞，我在稍后将详细说明什么是队首堵塞，这是推广WebSocket协议使用的一个主要障碍，不过WebSocket能随时提供可靠的传送服务。WebRTC数据通道的问题是建立连接的负担很高

那么为什么不使用现有的Web协议WebRTC呢？这张图表就能解释其中的原因，因为这是一个非常复杂的协议。它最初被构建为p2p通信协议，并且在建立连接之前会要求会话描述协议来进行SDP消息传递。在客户端服务器通信模型中我们不需要这样做，因为服务器端也需要通过客户端进行寻址。

WebRTC也要求的ICE、DTL、SCTP协议也通常不会再CDN中大规模部署。因此在某些情况下我们可以使用WebRTC，但WebRTC不是为服务器-客户端模型的应用程序而专门设计的。

https://segmentfault.com/a/1190000039710193?utm_source=tag-newest

建立一个WebRTC端到端连接需要3步：

                  #1 信令

                  #2 发现

                  #3 连接建立
 
 
 第一步：信令

         信令是建立WebRTC端到端连接的第一步。信令是想要建立端到端连接的两方用来交换初始信息的通道。

         建立初始阶段需要交换下面几个信息：

# 参与各方的IP及可使用的端口号（ICE候选）

# 媒体功能

# 会话控制消息

         Websocket被广泛应用于信令中。像Kurento这样的著名的WebRTC媒体服务器在使用它们。如果你想要在信令过程中进行安全的数据传输，那么建议你使用安全websocket（wss://）
         
    
  以Kurento为例，Kurento在8888端口接收websocket连接，在8443端口接收安全websocket连接。这个默认配置允许Kurento在你的网页服务器上平行运行，但因为它们使用的不是像80或者443这样的惯用端口，所以那些处于受限网络的电脑或者设备，有很大的可能性并不能与你的信令服务器通过这些端口进行通信。

         在端口80或者443上运行信令是你确保WebRTC高连接率所能做的第一件事。

第二步：发现（STUN和TURN）

         一旦在WebRTC终端和信令服务器之间建立了信令连接之后，就可以进行信息交换了。

         其中一个重要的信息就是公共IP和每个终端可以使用的端口。对于直接连入互联网的电脑来说，想要寻找到IP并不是一个问题，因为它（OS）知道自身的公共IP地址，并且可以很简单地通过浏览器查询到。但是这对于处于本地网络中（路由器后方）的电脑和设备来说就会成为一个问题，也包括通过3G/4G连入网络的移动设备，因为它们的IP是本地网络分配的IP地址。

         这些设备只知道它们自己的本地网络IP，所以它们使用STUN协议来：

1 与STUN服务器进行通信，来找到它们网络的公共IP以及可到达的端口。

2 打穿一个双向的通道，通过网络路由器的隐性NAT功能。

         此外，对称NAT后的设备只能与之前进行过连接的设备通信，所以还需要一个TURN服务器传输由一端发送的数据，因为另一端终端无法穿过对称NAT与我们的设备直接连接。

         每个WebRTC端点都要询问STUN/TURN服务器它们自己的公共IP和可连接的端口是什么。一旦接受到一个响应，WebRTC终端就会通过信令通道发送一个数据对给对方。这个IP:端口号对被叫做ICE候选。
         
与STUN/TURN服务器的通信是第二个可能造成WebRTC连接失败的地方。我们已经遇到过3个可能造成失败的原因：

1 默认的STUN/TURN端口被屏蔽

2 所有UDP端口都被屏蔽

3 STUN/TURN协议被禁止使用

 

端口屏蔽

         还记得我们建议通过端口80或者端口443建立信令连接吗？STUN和TURN有它们自己默认的端口（互不相同）：

# 发送（或接听）STUN/TURN请求的默认端口为3478。

# 通过TLS发送（或接听）STURN/TURN的默认端口5349。

# 一些服务器，像Google的通用STUN服务器，使用其他端口，如19305和19307。

         上述任何端口都可以被想要进行连接的两端之一所屏蔽。这些情况中，都无法连接到STUN/TURN服务器。

         为了避免这些问题的发生，可以对STUN/TURN服务器使用惯用端口（443/80），但是UDP和协议屏蔽的问题还是没法解决。

UDP屏蔽

         默认STUN/TURN消息通过UDP传输，意味防火墙不允许DNS询问机制使用53端口，也就不允许STUN/TURN消息通过防火墙。

         幸运的是，STUN/TURN服务器还可以通过TCP进行连接，通过指定URL中的transport参数来实现：turn:myTurnServer.com?transport=tcp

         意思基本是告诉WebRTC客户“对于TURN/STUN服务器，通过TCP连接而不是UDP”。你也可以指定udp或者tls。

STUN/TURN屏蔽

         开发者可能遇到的一个更严重的情况是STUN/TURN协议消息同时被屏蔽。举个例子，我们已经发现Tunnel Bear VPN屏蔽了STUN/TURN数据包，因为即便你通过VPN连入网络它们也会暴露你的真实IP。在这种情况中，你除了让用户在进行WebRTC通话的时候关闭这类app以外并没有其他解决办法。

 

第三步：建立连接

         在每个WebRTC终端知道了对方的ICE候选之后，就可以建立端到端连接了
         
         在一些WebRTC用例中，比如视频录制，终端会同时作为信令服务器和WebRTC终端。

         每个用户都会通过UDP向另一个终端发送数据：

# 如果是直接发送给对方的，那么就可能将数据发送给0-65535之间任一端口

# 如果是发送给TURN服务器的，那么数据会发送给49152-65535之间的一个端口

         没有一种方法可以控制这些端口，它们会在“发现”阶段中被分配并作为ICE候选的一部分。
         
         https://webrtc.org.cn/troubleshooting-connection/
 
 https://webrtc.org.cn/ice-restarts/
 
 实现思路：
两个浏览器打开同一页面，连接到同一个socket。
此时由一端点击建立连接，发起建立连接的一端就是offer(携带信号源信息)，发给另外一个端，另外一个端收到offer之后，发出响应answer(携带信号源信息)，offer端收到answer端信息进行存储；这样每个端都有了自己的信息和对方的信息，offer发出answer发出后设置了localDescription和remoteDescription后就会触发onicecandidate，如此一来，双方都有了对方的localDescription、remoteDescription和candidata；有了这三种数据之后，就可以触发Connection.onaddstream函数,然后通过theirVideo.srcObject = e.stream这个方法，把流写到video标签内，然后video标签里就会有对方的视频画面了。

注意：这样实现之后不能直接调用NAVIGATOR.GETUSERMEDIA()函数，浏览器会默认GETUSERMEDIA为UNDIFINED，如果是互联网模式这里需要设置浏览(如果只是本地测试则不需要)（谷歌浏览器配置方法）；WEBRTC如只是P2P不需要特别服务器，自已开发信令服务就可以啦，当要安装TURN SERVER 国内常有打洞不成功需要转发。

https://www.freesion.com/article/62371393961/


