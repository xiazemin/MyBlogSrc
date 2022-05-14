---
title: httpCache
layout: post
category: linux
author: 夏泽民
---
HTTP Caching 用好了，可以极大的减小服务器负载和减少网络带宽。十分有必要深入了解下 http 的 caching 协议。

　　先来看下请求/响应过程：

http 请求/响应

http 请求/响应

　　1、用 Last-Modified 头

　　在第一次请求的响应头返回 Last-Modified 内容，时间格式如：Wed, 22 Jul 2009 07:08:07 GMT。是零时区的 GMT 时间，servlet 中可以用 response.addDateHeader ("Last-Modified"， date.getTime ()); 加入响应头。如图：

last-modified 和 If-Modified-Since

last-modified 和 If-Modified-Since

　　Last-Modified 与 If-Modified-Since 对应的，前者是响应头，后者是请求头。服务器要处理 If-Modified-Since 请求头与 Last-Modified 对比看是否有更新，如果没有更新就返回 304 响应，否则按正常请求处理。如果要在动态内容中使用它们，那就要程序来处理了。

　　ps：servlet 取 If-Modified-Since 可以用 long last = requst.getDateHeader ("If-Modified-Since");

　　2、用 Etag 头

　　很多时间可能不能用时间来确定内容是否有更新。那可以用 Etag 头，etag 是以内容计算一个标识。计算的方式可以自己决定，比如可以用 crc32、md5等。

Etag 和 If-None-Match

Etag 和 If-None-Match

　　Etag 与 If-None-Match 是对应的，前者是响应头，后者是请求头。服务器要判断请求内容计算得到的 etag 是否与请求头 If-None-Match 是否一致，如果一致就表示没有更新，返回 304 就可，否则按正常请求处理。可以参考：用 HttpServletResponseWrapper 实现 Etag 过滤器

　　3、用 Expires 头，过期时间

　　当请求的内容有 Expires 头的时候，浏览器会在这个时间内不去下载这个请求的内容（这个行为对 F5 或 Ctrl+F2 无效，用 IE7，Firefox 3.5 试了，有效的比如：在地址输入后回车）。

expires 过期时间

expires 过期时间

　　在 servlet 中可以用 response.addDateHeader ("Expires"， date.getTime ()); 添加过期内容。

　　ps：在 httpwatch 中可以看到 Result 为 (Cached) 状态的。

　　4、用 max-age 的 Cache-Control 头

　　max-age 的值表示，多少秒后失效，在失效之前，浏览器不会去下载请求的内容（当然，这个行为对 F5 或 Ctrl+F2 无效）。比如：服务器写 max-age 响应：response.addHeader ("Cache-Control"， "max-age=10");

　　ps：如果你还要加一些 Cache-Control 的内容，比如：private，最好不要写两个 addHeader，而是一个 response.addHeader ("Cache-Control"， "private, max-age=10"); 否则 ie 可能对 max-age 无效，原因它只读第一个 Cache-Control 头。

　　小结：

　　Last-Modified 与 Etag 头（即是方式 1 和2）还是要请求服务器的，只是仅返回 304  头，不返回内容。所以浏览怎么 F5 ，304 都是有效的。但用 Ctrl+F5 是全新请求的（这是浏览器行为，不发送缓存相关的头）。

　　Expires 头与 max-age 缓存是不需要请求服务器的，直接从本地缓存中取。但 F5 会忽视缓存（所以使用 httpwatch 之类的 http 协议监察工具时，不要 F5 误认为 Expires 和 max-age 是无效的）。
　　
如果客户端发送了一个带条件的GET 请求且该请求已被允许，而文档的内容（自上次访问以来或者根据请求的条件）并没有改变，则服务器应当返回这个304状态码。简单的表达就是：服务端已经执行了GET，但文件未变化。
<!-- more -->
通过使用缓存web网站和web应用的性能能够得到显著的提升。Web caches能够减小延迟和网络流量，从而缩短展示资源所花费的时间。

在http中控制缓存行为的首部字段是Cache-Control，Cache-Control可以有多个指令，指令之间用逗号分割。Cache-Control是通用首部字段，它即能出现在响应头中也能出现在请求头中

一.缓存请求指令
Cache-Control: max-age=<seconds>
Cache-Control: max-stale[=<seconds>]
Cache-Control: min-fresh=<seconds>
Cache-Control: no-cache
Cache-Control: no-store
Cache-Control: no-transform
Cache-Control: only-if-cached
二.缓存响应指令
Cache-Control: must-revalidate
Cache-Control: no-cache
Cache-Control: no-store
Cache-Control: no-transform
Cache-Control: public
Cache-Control: private
Cache-Control: proxy-revalidate
Cache-Control: max-age=<seconds>
Cache-Control: s-maxage=<seconds>
三.指令介绍
1.与缓存能力相关的指令
1.public：指明任何缓存区都能够缓存响应

2.private：指明响应是针对单一用户的，不能存储在共享缓存区中，只有私有缓存区能存储响应

3.no-cache：在使用缓存资源之前要向源服务器验证缓存的有效性

4.only-if-cached：指明客户端只想从缓存中获取响应，不需要与源服务器联系确定资源是否更新



2.与到期时间相关的指令
1.max-age=<seconds>：指定资源被视为有效的最大时间量，单位为秒

2.s-maxage=<seconds>:与max-age=<seconds>指令的作用相似，但是它只对共享缓存有效，对私有缓存无效

3.max-stale[=<seconds>]：即使缓存资源过期客户端还是接受缓存区中的资源。如果未指定数值，那么无论过期多久，客户端都接受缓存区中的响应，如果指定了具体数值，那么即使过期，只要处于max-stale指定的时间内，客户端还是接受缓存区中的资源

4.min-fresh=<seconds>：表明从缓存区中返回未过指定时间的缓存资源

5.stale-while-revalidate=<seconds>：指示客户端愿意接受一个过期的响应，同时在后台异步检查一个新的响应。秒值表示客户机愿意接受过期响应的时间。

6.stale-if-error=<seconds>:如果检查新资源失败，客户端愿意接受过期的资源。秒值指示客户端在初始过期后愿意接受过期响应的时间。



3.与重新验证和重新加载相关的指令
1.must-revalidate：在使用过期的缓存资源之前，必须向源服务器确认缓存资源的状态并且不会使用过期的资源。它会忽略max-stale[=<seconds>]指令

2.proxy-revalidate：和must-revalidate作用类似，但是它只应用于共享缓存，会被私有缓存忽略

3.immutable：指明在任何时候都不能改变响应体

4.其他指令
1.no-store：不缓存请求和响应中的任何内容

2.no-transform：缓存不能改变请求和响应中的任何实体主体(The Content-Encoding, Content-Range, Content-Type 头部字段不能被修改)

e-Control标头被定义为HTTP / 1.1规范的一部分，并且取代了原先用于定义response缓存策略的Expire等，现在所有的现代浏览器都支持HTTP Cache-Control。

2.1 “no-cache”和“no-store”
在这里插入图片描述
no-cache：有缓存，但是不直接使用缓存，需要经过校验。
如果资源已经发生更改，在没有与服务端进行校验前，浏览器不能用前面返回的response用于满足对同一URL的后续请求。如果浏览器提供了正确的ETag，当配置了no-cache时，客户端会先发出请求向后端验证资源是否发生更改，如果资源未更改，则可以取消下载。
no-store：完全没有缓存，所有的资源都需要重新发请求
当服务端对资源设置了no-store时，不允许浏览器和所有中间层缓存response。。例如：一些类似于银行及其他隐私数据不适合缓存，每次的用户请求都需要发送到服务端，下载全部的response。

2.2 “public” 和 “private”
public
如果response被标记为“public”，那么，即使有与其相关的HTTP认证信息或者返回的response是不可缓存的status code，它依然可以被缓存。大多数情况下，public并不是必需的，因为其他具体指示缓存的信息，如max-age会表明当前的response在任何情况下都是要缓存的。
private
相比之下，浏览器可以缓存private的response。但是，这些响应通常只用于单个用户，因此不允许其他中间缓存对齐进行缓存。例如：一个用户的浏览器可以带有用户私有信息的HTML页面，但是CDN无法缓存页面。

2.3 max-age
max-age指令指定了允许重用缓存的response的最长时间（以秒为单位）。例如，max-age=60表示response可以缓存，并且在接下来的60s内可以被重用，无需发出新的request请求。

三、 定义最佳的Cache-Control策略
在这里插入图片描述
遵循上面的流程图为应用程序中使用的特定资源或一组资源制定最佳缓存策略。理想的情况下，我们应该在客户端上在尽可能长的时间中缓存尽可能多的响应，并为每个响应提供验证令牌（ETag），从而实现有效的重新验证。

Cache-Control指令和解释

指令	解释
max-age=86400	Response can be cached by browser and any intermediary caches (that is, it’s “public”) for up to 1 day (60 seconds x 60 minutes x 24 hours).
private, max-age=600	Response can be cached by the client’s browser only for up to 10 minutes (60 seconds x 10 minutes). 不能被代理服务器缓存。
public	响应可以被任何缓存区缓存。
no-cache	浏览器对请求回来的response做缓存，但是每次在向客户端（浏览器）提供响应数据时，缓存都要向服务器评估缓存响应的有效性
no-store	Response is not allowed to be cached and must be fetched in full on every request. 禁止一切缓存。
根据HTTP档案，排名前300,000（按照Alexa）的网站中，浏览器可以缓存几乎一半的下载响应，这对于网页的重复浏览和访问来说是个巨大的节省。当然，这并不意味着你的应用程序可以缓存50%的资源。有些网站的静态资源几乎不会变动，可能可以缓存超过90%的资源；其他网站可能有很多私有的或者是时间敏感的数据完全不能使用缓存。

审核自己的页面，确认哪些资源是可以缓存的并且确保它们返回了合适的Cache-Control和Etag。

四、作废和更新缓存的response
本地缓存的response可以一直被使用，直到资源expires；
在URL中嵌入一个文件内容的指纹可以强制客户端不使用缓存，更新response；
每个应用程序都需要定义自己的缓存层次结构来获得最佳性能。
所有从浏览器发出去的请求首先要路由到浏览器缓存中，检验是否存在可以用来完成该请求的有效缓存。如果有匹配的缓存，则response从缓存中读取，从而消除网络延迟和传输引起的数据成本。

如何作废或者更新一个缓存的response？
例如： 你已经告知访问者缓存一个css样式表24h(max-age = 86400），但是开发者刚刚提交了你希望向所有用户提供的更新。这时该如何通知访问者更新原有的缓存呢？此时，如果不改变资源的URL，无法更新资源。因为浏览器认为当前的缓存尚未过期。

浏览器缓存响应之后，直到根据max-age或者expire指示，缓存已过期或者缓存被清理调，都会使用缓存的资源。因此，在页面构建时，或者说在访问某个网站时，不同的用户可能使用不同版本的资源。刚刚获取资源的用户使用的是最新版本，但缓存了早期资源（仍然有效）的用户依然使用旧版本的资源。

如何充分利用“客户端缓存”和“快速更新”？
在资源内容更改时，改变资源的URL，强制用户下载新的response。通常，我们可以通过在文件名中嵌入版本号或者其他文件指纹来实现改变URL，例如：style.x234dff。

定义每个资源的缓存策略（定义缓存结构层次），不仅可以控制每个资源的缓存时长，还可以控制访问者获取和查看新版本的速度。
在这里插入图片描述
如上图中的例子：

HTML文件被标记为no-cache，意味着浏览器对于每个请求都会重新验证文档，如果资源内容发生改变，就会拉取最新版本。 此外，在HTML标记中，在css和javascript资源的URL中嵌入了指纹：如果这些文件的内容发生更改，HTML文件也会更改，从而加载一个HTML 响应的新副本。
CSS允许浏览器和其他中间缓存（例如：CDN）进行缓存，过期时间为1年。其实，我们也可以使用far future expires of 1 year， 因为我们在文件名中嵌入了文件指纹，如果CSS文件发生更新，其请求的URL也会发生改变。
js文件的过期时间也设置为1年，但是标记为private，可能是因为它里面包含了一些CDN不应缓存的私有用户数据。
图片文件中没有加版本和独一无二的hash指纹，直接进行缓存。缓存时间设置为1天。
结合ETag、Cache-Control和独一无二的URL，可以实现最好的缓存策略：更长的过期时间、控制特定资源的缓存以及缓存位置、按需更新。
五、缓存设置清单
世界上没有最好的缓存策略。根据你的流量模式、所提供的数据类型以及应用程序对于资源新鲜度的特定需求，我们需要对每个资源以及整体的“缓存层次结构”制定适当的缓存策略。

制定缓存策略时的一些提示和技巧：

使用一致的URL（hash/version）：如果对于内容完全相同的资源使用不同的URL，那么这个资源会不停地被获取和存储。Tips：URL需要区分大小写。
确保服务端提供了验证令牌（ETag）： 验证令牌的存在避免了在服务端资源没有任何修改时传输完全相同的字节内容。
确定哪些资源可以被中介缓存（public/private）：对于所有用户都相同的response可以在CDN和其他中介缓存中缓存。
确定每个资源的最佳缓存时长(max-age)：不同的资源可能具有不同的新鲜度要求。审核并确认每个资源合适的max-age。
确定当前网站的最佳缓存层次结构：结合资源的URL 和资源内容的指纹以及HTML的短缓存或者no-cache， 我们可以控制客户端获取和更新资源的速度。
减少混乱：有些资源的更新频率要高于其他资源。如果有一部分资源（比如：一个javascript函数或者一组CSS样式）需要经常更新，可以考虑将这一部分代码作为单独的文件。这样做，其余部分的资源（例如：不经常更新的库代码）就可以从缓存中提取。当获取更新时，使得需要下载内容的数量最小化

1、什么是Keep-Alive模式？

我们知道HTTP协议采用“请求-应答”模式，当使用普通模式，即非KeepAlive模式时，每个请求/应答客户和服务器都要新建一个连接，完成 之后立即断开连接（HTTP协议为无连接的协议）；当使用Keep-Alive模式（又称持久连接、连接重用）时，Keep-Alive功能使客户端到服 务器端的连接持续有效，当出现对服务器的后继请求时，Keep-Alive功能避免了建立或者重新建立连接。

http 1.0中默认是关闭的，需要在http头加入"Connection: Keep-Alive"，才能启用Keep-Alive；http 1.1中默认启用Keep-Alive，如果加入"Connection: close "，才关闭。目前大部分浏览器都是用http1.1协议，也就是说默认都会发起Keep-Alive的连接请求了，所以是否能完成一个完整的Keep- Alive连接就看服务器设置情况。

2、启用Keep-Alive的优点

从上面的分析来看，启用Keep-Alive模式肯定更高效，性能更高。因为避免了建立/释放连接的开销。下面是RFC 2616 上的总结：
By opening and closing fewer TCP connections, CPU time is saved in routers and hosts (clients, servers, proxies, gateways, tunnels, or caches), and memory used for TCP protocol control blocks can be saved in hosts.
HTTP requests and responses can be pipelined on a connection. Pipelining allows a client to make multiple requests without waiting for each response, allowing a single TCP connection to be used much more efficiently, with much lower elapsed time.
Network congestion is reduced by reducing the number of packets caused by TCP opens, and by allowing TCP sufficient time to determine the congestion state of the network.
Latency on subsequent requests is reduced since there is no time spent in TCP's connection opening handshake.
HTTP can evolve more gracefully, since errors can be reported without the penalty of closing the TCP connection. Clients using future versions of HTTP might optimistically try a new feature, but if communicating with an older server, retry with old semantics after an error is reported.
RFC 2616 （P47）还指出：单用户客户端与任何服务器或代理之间的连接数不应该超过2个。一个代理与其它服务器或代码之间应该使用不超过2 * N的活跃并发连接。这是为了提高HTTP响应时间，避免拥塞（冗余的连接并不能代码执行性能的提升）。

3、回到我们的问题（即如何判断消息内容/长度的大小？）

Keep-Alive模式，客户端如何判断请求所得到的响应数据已经接收完成（或者说如何知道服务器已经发生完了数据）？我们已经知道 了，Keep-Alive模式发送玩数据HTTP服务器不会自动断开连接，所有不能再使用返回EOF（-1）来判断（当然你一定要这样使用也没有办法，可 以想象那效率是何等的低）！下面我介绍两种来判断方法。
3.1、使用消息首部字段Conent-Length

故名思意，Conent-Length表示实体内容长度，客户端（服务器）可以根据这个值来判断数据是否接收完成。但是如果消息中没有Conent-Length，那该如何来判断呢？又在什么情况下会没有Conent-Length呢？请继续往下看……

3.2、使用消息首部字段Transfer-Encoding

当客户端向服务器请求一个静态页面或者一张图片时，服务器可以很清楚的知道内容大小，然后通过Content-length消息首部字段告诉客户端 需要接收多少数据。但是如果是动态页面等时，服务器是不可能预先知道内容大小，这时就可以使用Transfer-Encoding：chunk模式来传输 数据了。即如果要一边产生数据，一边发给客户端，服务器就需要使用"Transfer-Encoding: chunked"这样的方式来代替Content-Length。

chunk编码将数据分成一块一块的发生。Chunked编码将使用若干个Chunk串连而成，由一个标明长度为0 的chunk标示结束。每个Chunk分为头部和正文两部分，头部内容指定正文的字符总数（十六进制的数字 ）和数量单位（一般不写），正文部分就是指定长度的实际内容，两部分之间用回车换行(CRLF) 隔开。在最后一个长度为0的Chunk中的内容是称为footer的内容，是一些附加的Header信息（通常可以直接忽略）。 Chunk编码的格式如下：

复制代码

代码如下:

Chunked-Body = *<strong>chunk </strong>
"0" CRLF
footer
CRLF
chunk = chunk-size [ chunk-ext ] CRLF
chunk-data CRLF</p><p>hex-no-zero = &lt;HEX excluding "0"&gt;</p><p>chunk-size = hex-no-zero *HEX
chunk-ext = *( ";" chunk-ext-name [ "=" chunk-ext-value ] )
chunk-ext-name = token
chunk-ext-val = token | quoted-string
chunk-data = chunk-size(OCTET)</p><p>footer = *entity-header

即Chunk编码由四部分组成： 1、<strong>0至多个chunk块</strong> ，2、<strong>"0" CRLF </strong>，3、<strong>footer </strong>，4、<strong>CRLF</strong> <strong>.</strong> 而每个chunk块由：chunk-size、chunk-ext（可选）、CRLF、chunk-data、CRLF组成。

4、消息长度的总结

其实，上面2中方法都可以归纳为是如何判断http消息的大小、消息的数量。RFC 2616 对 消息的长度总结如下：一个消息的transfer-length（传输长度）是指消息中的message-body（消息体）的长度。当应用了 transfer-coding（传输编码），每个消息中的message-body（消息体）的长度（transfer-length）由以下几种情况 决定（优先级由高到低）：
任何不含有消息体的消息（如1XXX、204、304等响应消息和任何头(HEAD，首部)请求的响应消息），总是由一个空行（CLRF）结束。
如果出现了Transfer-Encoding头字段 并且值为非“identity”，那么transfer-length由“chunked” 传输编码定义，除非消息由于关闭连接而终止。
如果出现了Content-Length头字段，它的值表示entity-length（实体长度）和transfer-length（传输长 度）。如果这两个长度的大小不一样（i.e.设置了Transfer-Encoding头字段），那么将不能发送Content-Length头字段。并 且如果同时收到了Transfer-Encoding字段和Content-Length头字段，那么必须忽略Content-Length字段。
如果消息使用媒体类型“multipart/byteranges”，并且transfer-length 没有另外指定，那么这种自定界（self-delimiting）媒体类型定义transfer-length 。除非发送者知道接收者能够解析该类型，否则不能使用该类型。
由服务器关闭连接确定消息长度。（注意：关闭连接不能用于确定请求消息的结束，因为服务器不能再发响应消息给客户端了。）
为了兼容HTTP/1.0应用程序，HTTP/1.1的请求消息体中必须包含一个合法的Content-Length头字段，除非知道服务器兼容 HTTP/1.1。一个请求包含消息体，并且Content-Length字段没有给定，如果不能判断消息的长度，服务器应该用用400 (bad request) 来响应；或者服务器坚持希望收到一个合法的Content-Length字段，用 411 (length required)来响应。

所有HTTP/1.1的接收者应用程序必须接受“chunked” transfer-coding (传输编码)，因此当不能事先知道消息的长度，允许使用这种机制来传输消息。消息不应该够同时包含 Content-Length头字段和non-identity transfer-coding。如果一个消息同时包含non-identity transfer-coding和Content-Length ，必须忽略Content-Length 。

5、HTTP头字段总结

最后我总结下HTTP协议的头部字段。
1、 Accept：告诉WEB服务器自己接受什么介质类型，/ 表示任何类型，type/* 表示该类型下的所有子类型，type/sub-type。
2、 Accept-Charset： 浏览器申明自己接收的字符集 Accept-Encoding： 浏览器申明自己接收的编码方法，通常指定压缩方法，是否支持压缩，支持什么压缩方法（gzip，deflate） Accept-Language：浏览器申明自己接收的语言 语言跟字符集的区别：中文是语言，中文有多种字符集，比如big5，gb2312，gbk等等。
3、 Accept-Ranges：WEB服务器表明自己是否接受获取其某个实体的一部分（比如文件的一部分）的请求。bytes：表示接受，none：表示不接受。
4、 Age：当代理服务器用自己缓存的实体去响应请求时，用该头部表明该实体从产生到现在经过多长时间了。
5、 Authorization：当客户端接收到来自WEB服务器的 WWW-Authenticate 响应时，用该头部来回应自己的身份验证信息给WEB服务器。
6、 Cache-Control：请求：no-cache（不要缓存的实体，要求现在从WEB服务器去取） max-age：（只接受 Age 值小于 max-age 值，并且没有过期的对象） max-stale：（可以接受过去的对象，但是过期时间必须小于 max-stale 值） min-fresh：（接受其新鲜生命期大于其当前 Age 跟 min-fresh 值之和的缓存对象） 响应：public(可以用 Cached 内容回应任何用户) private（只能用缓存内容回应先前请求该内容的那个用户） no-cache（可以缓存，但是只有在跟WEB服务器验证了其有效后，才能返回给客户端） max-age：（本响应包含的对象的过期时间） ALL: no-store（不允许缓存）
7、 Connection：请求：close（告诉WEB服务器或者代理服务器，在完成本次请求的响应后，断开连接，不要等待本次连接的后续请求了）。 keepalive（告诉WEB服务器或者代理服务器，在完成本次请求的响应后，保持连接，等待本次连接的后续请求）。 响应：close（连接已经关闭）。 keepalive（连接保持着，在等待本次连接的后续请求）。 Keep-Alive：如果浏览器请求保持连接，则该头部表明希望 WEB 服务器保持连接多长时间（秒）。例如：Keep-Alive：300
8、 Content-Encoding：WEB服务器表明自己使用了什么压缩方法（gzip，deflate）压缩响应中的对象。例如：Content-Encoding：gzip
9、Content-Language：WEB 服务器告诉浏览器自己响应的对象的语言。
10、Content-Length： WEB 服务器告诉浏览器自己响应的对象的长度。例如：Content-Length: 26012
11、Content-Range： WEB 服务器表明该响应包含的部分对象为整个对象的哪个部分。例如：Content-Range: bytes 21010-47021/47022
12、Content-Type： WEB 服务器告诉浏览器自己响应的对象的类型。例如：Content-Type：application/xml
13、ETag：就是一个对象（比如URL）的标志值，就一个对象而言，比如一个 html 文件，如果被修改了，其 Etag 也会别修改，所以ETag 的作用跟 Last-Modified 的作用差不多，主要供 WEB 服务器判断一个对象是否改变了。比如前一次请求某个 html 文件时，获得了其 ETag，当这次又请求这个文件时，浏览器就会把先前获得的 ETag 值发送给WEB 服务器，然后 WEB 服务器会把这个 ETag 跟该文件的当前 ETag 进行对比，然后就知道这个文件有没有改变了。
14、 Expired：WEB服务器表明该实体将在什么时候过期，对于过期了的对象，只有在跟WEB服务器验证了其有效性后，才能用来响应客户请求。是 HTTP/1.0 的头部。例如：Expires：Sat, 23 May 2009 10:02:12 GMT
15、 Host：客户端指定自己想访问的WEB服务器的域名/IP 地址和端口号。例如：Host：rss.sina.com.cn
16、 If-Match：如果对象的 ETag 没有改变，其实也就意味著对象没有改变，才执行请求的动作。
17、 If-None-Match：如果对象的 ETag 改变了，其实也就意味著对象也改变了，才执行请求的动作。
18、 If-Modified-Since：如果请求的对象在该头部指定的时间之后修改了，才执行请求的动作（比如返回对象），否则返回代码304，告诉浏览器 该对象没有修改。例如：If-Modified-Since：Thu, 10 Apr 2008 09:14:42 GMT
19、 If-Unmodified-Since：如果请求的对象在该头部指定的时间之后没修改过，才执行请求的动作（比如返回对象）。
20、 If-Range：浏览器告诉 WEB 服务器，如果我请求的对象没有改变，就把我缺少的部分给我，如果对象改变了，就把整个对象给我。浏览器通过发送请求对象的 ETag 或者 自己所知道的最后修改时间给 WEB 服务器，让其判断对象是否改变了。总是跟 Range 头部一起使用。
21、 Last-Modified：WEB 服务器认为对象的最后修改时间，比如文件的最后修改时间，动态页面的最后产生时间等等。例如：Last-Modified：Tue, 06 May 2008 02:42:43 GMT
22、 Location：WEB 服务器告诉浏览器，试图访问的对象已经被移到别的位置了，到该头部指定的位置去取。例如：Location：http://i0.sinaimg.cn/dy/deco/2008/0528/sinahome_0803_ws_005_text_0.gif</a>
23、 Pramga：主要使用 Pramga: no-cache，相当于 Cache-Control： no-cache。例如：Pragma：no-cache
24、 Proxy-Authenticate： 代理服务器响应浏览器，要求其提供代理身份验证信息。Proxy-Authorization：浏览器响应代理服务器的身份验证请求，提供自己的身份信息。
25、 Range：浏览器（比如 Flashget 多线程下载时）告诉 WEB 服务器自己想取对象的哪部分。例如：Range: bytes=1173546-
26、 Referer：浏览器向 WEB 服务器表明自己是从哪个 网页/URL 获得/点击 当前请求中的网址/URL。例如：Referer：http://www.sina.com/</a>
27、 Server: WEB 服务器表明自己是什么软件及版本等信息。例如：Server：Apache/2.0.61 (Unix)
28、 User-Agent: 浏览器表明自己的身份（是哪种浏览器）。例如：User-Agent：Mozilla/5.0 (Windows; U; Windows NT 5.1; zh-CN; rv:1.8.1.14) Gecko/20080404 Firefox/2、0、0、14
29、 Transfer-Encoding: WEB 服务器表明自己对本响应消息体（不是消息体里面的对象）作了怎样的编码，比如是否分块（chunked）。例如：Transfer-Encoding: chunked
30、 Vary: WEB服务器用该头部的内容告诉 Cache 服务器，在什么条件下才能用本响应所返回的对象响应后续的请求。假如源WEB服务器在接到第一个请求消息时，其响应消息的头部为：Content- Encoding: gzip; Vary: Content-Encoding那么 Cache 服务器会分析后续请求消息的头部，检查其 Accept-Encoding，是否跟先前响应的 Vary 头部值一致，即是否使用相同的内容编码方法，这样就可以防止 Cache 服务器用自己 Cache 里面压缩后的实体响应给不具备解压能力的浏览器。例如：Vary：Accept-Encoding
31、 Via： 列出从客户端到 OCS 或者相反方向的响应经过了哪些代理服务器，他们用什么协议（和版本）发送的请求。当客户端请求到达第一个代理服务器时，该服务器会在自己发出的请求里面添 加 Via 头部，并填上自己的相关信息，当下一个代理服务器收到第一个代理服务器的请求时，会在自己发出的请求里面复制前一个代理服务器的请求的Via 头部，并把自己的相关信息加到后面，以此类推，当 OCS 收到最后一个代理服务器的请求时，检查 Via 头部，就知道该请求所经过的路由。例如：Via：1.0 236.D0707195.sina.com.cn:80 (squid/2.6.STABLE13)
=============================================================================== HTTP 请求消息头部实例： Host：rss.sina.com.cn User-Agent：Mozilla/5、0 (Windows; U; Windows NT 5、1; zh-CN; rv:1、8、1、14) Gecko/20080404 Firefox/2、0、0、14 Accept：text/xml,application/xml,application/xhtml+xml,text/html;q=0、9,text/plain;q=0、8,image/png,/;q=0、5 Accept-Language：zh-cn,zh;q=0、5 Accept-Encoding：gzip,deflate Accept-Charset：gb2312,utf-8;q=0、7,*;q=0、7 Keep-Alive：300 Connection：keep-alive Cookie：userId=C5bYpXrimdmsiQmsBPnE1Vn8ZQmdWSm3WRlEB3vRwTnRtW &lt;-- Cookie If-Modified-Since：Sun, 01 Jun 2008 12:05:30 GMT Cache-Control：max-age=0 HTTP 响应消息头部实例： Status：OK - 200 -- 响应状态码，表示 web 服务器处理的结果。 Date：Sun, 01 Jun 2008 12:35:47 GMT Server：Apache/2.0.61 (Unix) Last-Modified：Sun, 01 Jun 2008 12:35:30 GMT Accept-Ranges：bytes Content-Length：18616 Cache-Control：max-age=120 Expires：Sun, 01 Jun 2008 12:37:47 GMT Content-Type：application/xml Age：2 X-Cache：HIT from 236-41.D07071951.sina.com.cn -- 反向代理服务器使用的 HTTP 头部 Via：1.0 236-41.D07071951.sina.com.cn:80 (squid/2.6.STABLE13) Connection：close

这两天重点补充了一下前端的缓存机制。缓存是个好东西啊，可以减少网络IO消耗，提高访问速度，对前端性能优化有着显著的效果。

       依在下愚见，前端的存储分为两个方向，本地缓存以及浏览器缓存。今天主要介绍一下浏览器缓存的原理。


浏览器缓存是依赖于http协议的，所以一下就成之谓HTTP缓存～

HTTP缓存也可以分为两种，强缓存和协商缓存，优先级较高的是强缓存，在命中强缓存失败的情况下，才会走协商缓存

强缓存的实现原理：
        强缓存是利用http头中的Expires和Cache-Control两个字段来控制的。当请求“再次”发起，浏览器来检测expire和cache-control来判断目标资源是否符合强缓存，若符合则直接从缓存中获取资源，不会再与服务端发生通信。


来看一张配图吧，我们把目光关注到Expires和Cache-Control上，这两位爷台想表达什么意思呢？


        我看过的资料中把Expires翻译成一个时间戳，其实在下看来，从字面意思理解，它就是个过期时间嘛，就像是我们超市里买的营养快线一样，都有保质期的嘛汪。所以当浏览器再次，一定是再次向服务器请求资源的时候，浏览器就会先对比本地时间和expires的保质期，如果本地时间小于expires的保质期，那么就直接去缓存中取这个资源喵。

       “那要是我把本地时间改了的话，是不是就无法达到预期的效果啦？”

       “不错！问得好，紫薇”


        在考虑到expires这个过期时间的方案不靠谱之后，HTTP1.1新增了Cache-Control来接expires的班儿。Cache-Control可以视作是expires的完全替代方案，那么它是怎么给expires干掉的呢，灯光师，来个特写来。


        各位官爷请看，在Cache-Control中，我们通过max-age来控制资源的有效期。这次我们用的不是一个时间点，而是一个时间长度，31536000代表31536000秒，在31536000秒之内不管访问几次，都是走我浏览器缓存，如果cache-control和expires同时存在，那么优先考虑cache-control。耶！

        关于Cache-Control的另一个参数public，是在告诉浏览器是否可以被代理服务器缓存，如果我们只想浏览器缓存，默认设置private便可，Cache-Control还有一个为代理浏览器服务的参数叫做s-maxage，如果两者同时出现且s-maxage未过期，则向代理服务器请求其缓存内容。在大型项目中，架构会依赖各种代理服务器，所以我们不得不去考虑代理服务器的缓存问题（这一块儿只是在下还未参透，就不在各位看官面前妄言了）。

协商缓存的实现原理：
        前文书说到，强缓存会比协商缓存的优先级要高，这是因为协商缓存也会向服务器发送http请求的！

        协商缓存是依赖Last-Modified和Etag，他们会埋伏在响应头中，来告诉浏览器这个服务器中的资源有没有修改过呀，如果有的话从服务器中拉取最新资源，如果没有的话继续使用浏览器中缓存汪。


        Last-Modified也是个时间戳，它会在我们首次请求的时候随着Response Headers返回，告诉浏览器我们最后一次修改时间是18年六月，随后我们每次请求都会带上一个叫If-Modified-Since的时间戳字段，它的值正是上一次请求时Last-Modified的值。


        此时服务器接收到了这个时间戳，会比对该时间戳和服务器上资源的最后修改时间是否一致，如果一致则返回304状态码告诉浏览器就用缓存就好啦~

        但是！！！
        1，如果我们编辑了文件，但是并没有对文件内容做修改，Last-Modified时间戳也会变，不该重新请求的时候也去重新请求了。

        2，当我们修改文件的手速巨快，由于If-Modified-Since只能检查到以秒为最小计量单位的时间差，所以它是感知不到这个改动的，该重新请求的时候也不会重新请求了。

        所以，etag作为Last-Modified的补充就出现了。

        Etag和Last-Modified类似，当我们首次请求时，会在响应头里获取到一个最初的标识字符串，这个是服务器算出的hash值，比last-modified更准确：


        等到用户下一次请求的时候，请求头里会带上一个名为If-None-Match的值给服务器来比对，是的，这个值就是etag中取出来的


        服务器在比对了这串hash值之后就可以知道文件有没有修改过，从而是返回304还是获取最新资源了。

        但是etag的hash值在生成过程中也是需要服务器额外付出开销的，这也是他的弊端喵。

        到此为止http的缓存已经全部讲完了，在下再为各位模拟一下协商缓存中浏览器与服务器的对话哦。

强缓存：
        场景：我访问了一下自己的网站，刷出了自己的帅照~

        浏览器：“老服，我要一张 Yabble最帅.jpg 图片，我展示用”

        服务器：“知道了，给你这个 Yabble最帅.jpg 你先用着，这货也不爱打扮，3个月之内用它都可以(cache-control: 大概3个月的秒数)。”

        场景：我立刻又访问了一次~

        浏览器：“老服，这糙货又访问了，快给我图片...诶，还没过期呀，那直接用旧的就得了呗”

协商缓存：
        场景：我访问了一下自己的网站，刷出了自己的帅照~

        浏览器：“老服，我要一张 Yabble最帅.jpg 图片，我展示”（嘿嘿，原谅我偷懒了）

        服务器：“知道了，给你这个 Yabble最帅.jpg 你先用着，这货最近还敷上面膜了，不知道哪儿天就变小鲜肉了，给你个暗号（etag）记着点，等他哪天捯饬的不一样了我告诉你~”

        场景：我又访问了一次

        浏览器：“老服，他又访问了，快给我这糙货的图片，这是暗号（If-none-match）”

        服务器：“来，暗号给我看看，没变样嘛，还用那张旧的吧”

        场景：我面膜美颜加刮胡子之后上传替换了一张自拍照，访问网站ing~

        浏览器：“老服，他又访问了，快快，这是暗号（If-none-match）”     

        服务器：“哦哦，暗号给我来看，哟呵，和现在的照片不一样了嘿，给你新的。刷刷刷~
