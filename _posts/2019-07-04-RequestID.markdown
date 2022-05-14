---
title: RequestID fastcgi
layout: post
category: php
author: 夏泽民
---
为什么需要在消息头发送 RequestID 这个标识？
如果是每个连接仅处理一个请求，发送 RequestID 则略显多余。

但是我们的 Web 服务器和 FastCGI 进程之间的连接可能处理多个请求，即一个连接可以处理多个请求。所以才需要采用数据包协议而不是直接使用单个数据流的原因：以实现「多路复用」。

因此，由于每个数据包都包含唯一的 RequestID，所以 Web 服务器才能在一个连接上发送任意数量的请求，并且 FastCGI 进程也能够从一个连接上接收到任意数量的请求数据包。

另外我们还需要明确一点就是 Web 服务器 与 FastCGI 进程间通信是 无序的。即使我们在交互过程中看起来一个请求是有序的，但是我们的 Web 服务器也有可能在同一时间发出几十个 BEGIN_REQUEST 类型的数据包，以此类推

为什么不是5元组？因为一般服务器和cgi的服务器都是固定的，甚至一对一
https://fast-cgi.github.io/original/
<!-- more -->
CGI(Common Gateway Interface:通用网关接口)是Web 服务器运行时外部程序的规范,按CGI 编写的程序可以扩展服务器功能。CGI 应用程序能与浏览器进行交互,还可通过数据库API 与数据库服务器等外部数据源进行通信,从数据库服务器中获取数据。--百度百科

CGI协议同 HTTP 协议一样是一个「应用层」协议，它的 功能 是为了解决 Web 服务器与 PHP 应用（或其他 Web 应用）之间的通信问题。

既然它是一个「协议」，换言之它与语言无关，即只要是实现类 CGI 协议的应用就能够实现相互的通信。

深入CGI协议
我们已经知道了 CGI 协议是为了完成 Web 服务器和应用之间进行数据通信这个问题。那么，这一节我们就来看看究竟它们之间是如何进行通信的。

简单来讲 CGI 协议它描述了 Web 服务器和应用程序之间进行数据传输的格式，并且只要我们的编程语言支持标准输入（STDIN）、标准输出（STDOUT）以及环境变量等处理，你就可以使用它来编写一个 CGI 程序。

CGI的运行原理
当用户访问我们的 Web 应用时，会发起一个 HTTP 请求。最终 Web 服务器接收到这个请求。

Web 服务器创建一个新的 CGI 进程。在这个进程中，将 HTTP 请求数据已一定格式解析出来，并通过标准输入和环境变量传入到 URL 指定的 CGI 程序（PHP 应用 $_SERVER）。

Web 应用程序处理完成后将返回数据写入到标准输出中，Web 服务器进程则从标准输出流中读取到响应，并采用 HTTP 协议返回给用户响应。

一句话就是 Web 服务器中的 CGI 进程将接收到的 HTTP 请求数据读取到环境变量中，通过标准输入转发给 PHP 的 CGI 程序；当 PHP 程序处理完成后，Web 服务器中的 CGI 进程从标准输出中读取返回数据，并转换回 HTTP 响应消息格式，最终将页面呈献给用户。然后 Web 服务器关闭掉这个 CGI 进程。

可以说 CGI 协议特别擅长处理 Web 服务器和 Web 应用的通信问题。然而，它有一个严重缺陷，对于每个请求都需要重新 fork 出一个 CGI 进程，处理完成后立即关闭。

CGI协议的缺陷
每次处理用户请求，都需要重新 fork CGI 子进程、销毁 CGI 子进程。

一系列的 I/O 开销降低了网络的吞吐量，造成了资源的浪费，在大并发时会产生严重的性能问题。

深入FastCGI协议
从功能上来讲，CGI 协议已经完全能够解决 Web 服务器与 Web 应用之间的数据通信问题。但是由于每个请求都需要重新 fork 出 CGI 子进程导致性能堪忧，所以基于 CGI 协议的基础上做了改进便有了 FastCGI 协议，它是一种常驻型的 CGI 协议。

本质上来将 FastCGI 和 CGI 协议几乎完全一样，它们都可以从 Web 服务器里接收到相同的数据，不同之处在于采取了不同的通信方式。

再来回顾一下 CGI 协议每次接收到 HTTP 请求时，都需要经历 fork 出 CGI 子进程、执行处理并销毁 CGI 子进程这一系列工作。

而 FastCGI 协议采用 进程间通信(IPC) 来处理用户的请求，下面我们就来看看它的运行原理。

FastCGI协议运行原理
FastCGI 进程管理器启动时会创建一个 主（Master） 进程和多个 CGI 解释器进程（Worker 进程），然后等待 Web 服务器的连接。

Web 服务器接收 HTTP 请求后，将 CGI 报文通过 套接字（UNIX 或 TCP Socket）进行通信，将环境变量和请求数据写入标准输入,转发到 CGI 解释器进程。

CGI 解释器进程完成处理后将标准输出和错误信息从同一连接返回给 Web 服务器。

CGI 解释器进程等待下一个 HTTP 请求的到来。

为什么是 FastCGI 而非 CGI 协议
如果仅仅因为工作模式的不同，似乎并没有什么大不了的。并没到非要选择 FastCGI 协议不可的地步。

然而，对于这个看似微小的差异，但意义非凡，最终的结果是实现出来的 Web 应用架构上的差异。

CGI 与 FastCGI 架构
在 CGI 协议中，Web 应用的生命周期完全依赖于 HTTP 请求的声明周期。

对每个接收到的 HTTP 请求，都需要重启一个 CGI 进程来进行处理，处理完成后必须关闭 CGI 进程，才能达到通知 Web 服务器本次 HTTP 请求处理完成的目的。

但是在 FastCGI 中完全不一样。

FastCGI 进程是常驻型的，一旦启动就可以处理所有的 HTTP 请求，而无需直接退出。

再看 FastCGI 协议
通过前面的讲解，我们相比已经可以很准确的说出来 FastCGI 是一种通信协议 这样的结论。现在，我们就将关注的焦点挪到协议本身，来看看这个协议的定义。

同 HTTP 协议一样，FastCGI 协议也是有消息头和消息体组成。

消息头信息
主要的消息头信息如下:

Version: 用于表示 FastCGI 协议版本号。

Type: 用于标识 FastCGI 消息的类型 - 用于指定处理这个消息的方法。

RequestID: 标识出当前所属的 FastCGI 请求。

Content Length: 数据包包体所占字节数。

消息类型定义
BEGIN_REQUEST: 从 Web 服务器发送到 Web 应用，表示开始处理新的请求。

ABORT_REQUEST: 从 Web 服务器发送到 Web 应用，表示中止一个处理中的请求。比如，用户在浏览器发起请求后按下浏览器上的「停止按钮」时，会触发这个消息。

END_REQUEST: 从 Web 应用发送给 Web 服务器，表示该请求处理完成。返回数据包里包含「返回的代码」，它决定请求是否成功处理。

PARAMS: 「流数据包」，从 Web 服务器发送到 Web 应用。此时可以发送多个数据包。发送结束标识为从 Web 服务器发出一个长度为 0 的空包。且 PARAMS 中的数据类型和 CGI 协议一致。即我们使用 $_SERVER 获取到的系统环境等。

STDIN: 「流数据包」，用于 Web 应用从标准输入中读取出用户提交的 POST 数据。

STDOUT: 「流数据报」，从 Web 应用写入到标准输出中，包含返回给用户的数据。

Web 服务器和 FastCGI 交互过程
Web 服务器接收用户请求，但最终处理请求由 Web 应用完成。此时，Web 服务器尝试通过套接字（UNIX 或 TCP 套接字，具体使用哪个由 Web 服务器配置决定）连接到 FastCGI 进程。

FastCGI 进程查看接收到的连接。选择「接收」或「拒绝」连接。如果是「接收」连接，则从标准输入流中读取数据包。

如果 FastCGI 进程在指定时间内没有成功接收到连接，则该请求失败。否则，Web 服务器发送一个包含唯一的RequestID 的 BEGIN_REQUEST 类型消息给到 FastCGI 进程。后续所有数据包发送都包含这个 RequestID。 然后，Web 服务器发送任意数量的 PARAMS 类型消息到 FastCGI 进程。一旦发送完毕，Web 服务器通过发送一个空PARAMS 消息包，然后关闭这个流。 另外，如果用户发送了 POST 数据 Web 服务器会将其写入到 标准输入（STDIN） 发送给 FastCGI 进程。当所有 POST 数据发送完成，会发送一个空的 标准输入（STDIN） 来关闭这个流。

同时，FastCGI 进程接收到 BEGINREQUEST 类型数据包。它可以通过响应 ENDREQUEST 来拒绝这个请求。或者接收并处理这个请求。如果接收请求，FastCGI 进程会等待接收所有的 PARAMS 和 标准输入数据包。 然后，在处理请求并将返回结果写入 标准输出（STDOUT） 流。处理完成后，发送一个空的数据包到标准输出来关闭这个流，并且会发送一个 END_REQUEST 类型消息通知 Web 服务器，告知它是否发生错误异常。

为什么需要在消息头发送 RequestID 这个标识？
如果是每个连接仅处理一个请求，发送 RequestID 则略显多余。

但是我们的 Web 服务器和 FastCGI 进程之间的连接可能处理多个请求，即一个连接可以处理多个请求。所以才需要采用数据包协议而不是直接使用单个数据流的原因：以实现「多路复用」。

因此，由于每个数据包都包含唯一的 RequestID，所以 Web 服务器才能在一个连接上发送任意数量的请求，并且 FastCGI 进程也能够从一个连接上接收到任意数量的请求数据包。

另外我们还需要明确一点就是 Web 服务器 与 FastCGI 进程间通信是 无序的。即使我们在交互过程中看起来一个请求是有序的，但是我们的 Web 服务器也有可能在同一时间发出几十个 BEGIN_REQUEST 类型的数据包，以此类推。

PHP-FPM
PHP-FPM即PHP-FastCGI Process Manager.

PHP-FPM是FastCGI的实现，并提供了进程管理的功能。

进程包含 master 进程和 worker 进程两种进程。

master 进程只有一个，负责监听端口，接收来自 Web Server 的请求，而 worker 进程则一般有多个(具体数量根据实际需要配置)，每个进程内部都嵌入了一个 PHP 解释器，是 PHP 代码真正执行的地方。

PHP-FPM 是 FastCGI 进程管理器（PHP FastCGI Process Manager）(http://php.net/manual/zh/install.fpm.php)，用于替换 PHP 内核的 FastCGI 的大部分附加功能（或者说一种替代的 PHP FastCGI 实现），对于高负载网站是非常有用的。

PHP-FPM如何工作的？
PHP-FPM 进程管理器有两种进程组成，一个 Master 进程和多个 Worker 进程。Master 进程负责监听端口，接收来自 Web 服务器的请求，然后指派具体的 Worker 进程处理请求；worker 进程则一般有多个 (依据配置决定进程数)，每个进程内部都嵌入了一个 PHP 解释器，用来执行 PHP 代码。

Nginx 服务器如何与 FastCGI 协同工作
Nginx 服务器无法直接与 FastCGI 服务器进行通信，需要启用 ngx_http_fastcgi_module 模块进行代理配置，才能将请求发送给 FastCGI 服务。




传统CGI建立一个应用进程，用来响应一个请求，然后退出。FastCGI的初始状态比CGI/1.1的更精炼，因为FastCGI进程开始时并不会连接到任何东西，它没有传统的打开文件stdin,stdout和stderr,也不从环境变量中取得额外的信息。FastCGI的初始状态只是监听一个从服务器接受连接的socket。

在一个FastCGI进程在它监听的socket接受一个连接后，该进程执行一个简单的协议来接收和发送数据。协议有两个目的,首先，协议将多个独立的FastCGI请求复用到一个传输连接上，这可支持那些处理多个并发请求的应用程序使用事件驱动或者多线程编程技术。其次，在每个请求中，协议可双向提供多个独立的数据流。这样，stdout和stderr数据可通过一个从应用程序到服务器的单个传输连接来传送,而不是CGI/1.1那样要用不同的管道。

当应用开始执行时，web服务器只留下一个文件描述符FCGI_LISTENSOCK_FILENO, 这个描述符指向一个由web服务器创建的监听socket。

FCGI_LISTENSOCK_FILENO等于STDIN_FILENO(标准输入). 当应用开始执行时，标准描述符 STDOUT_FILENO和STDERR_FILENO被关闭。 应用程序判断自己是被CGI还是FastCGI调用的可靠的方法是，调用getpeername(FCGI_LISTENSOCK_FILENO)，如果返回-1 ，并且errno被设为ENOTCONN的话，就是使用FastCGI方式。

服务器的可靠连接，选择Unix流管道(AF_UNIX)还是TCP/IP(AF_INET),是隐含在FCGI_LISTENSOCK_FILENO Socket的内部状态中的。

一个FastCGI应用程序在一个由文件描述符FCGI_LISTENSOCK_FILENO定义的socket上调用accept() 来接受一个新的传输连接。如果accept()成功执行，而且也绑定了FCGI_WEB_SERVER_ADDRS环境变量，应用程序立即执行下面的处理
FCGI_WEB_SERVER_ADDRS:这个值是服务器 有效IP地址的列表

1)如果FCGI_WEB_SERVER_ADDRS被绑定，应用程序为列表中的成员 检查新连接的对方的的IP。如果检查失败（包括连接没有使用 TCP/IP传输的情况),应用程序将关闭这个连接。

2)FCGI_WEB_SERVER_ADDRS用逗号分开，一个合法的例子
FCGI_WEB_SERVER_ADDRS=199.170.183.28,199.170.183.71

应用程序执行一个使用简单协议从服务器送过来的请求。协议的细节依赖于应用程序的角色，但一般来讲，服务器首先发送一些参数和其他数据给应用程序，然后应用程序发送结果数据给服务器，最后应用程序发送给服务器一个指示，说请求已经完成。

在传输连接上的所有的数据都通过FastCGI records传送。FastCGI records完成两件事。
第一，记录将复用几个不独立的FastCGI请求到传输连接 上。复用支持应用程序处理多个使用事件驱动或者多线程编程技术的并发请求。
第二，记录在单个请求中提供多个不同方向的独立的数据流.以这种方式，stdout 和 stderr能通过单个传输连接，从应用程序传送到服务器，而不是要求不同的连接。

typedef struct {
            unsigned char version;
            unsigned char type;
            unsigned char requestIdB1;
            unsigned char requestIdB0;
            unsigned char contentLengthB1;
            unsigned char contentLengthB0;
            unsigned char paddingLength;
            unsigned char reserved;
            unsigned char contentData[contentLength];
            unsigned char paddingData[paddingLength];
        } FCGI_Record;
FastCGI记录 含有 一个固定长度的前缀,加上 可变长度的内容和填充字节
version:    标示FastCGI协议的版本。本规范的版本FCGI_VERSION_1
type:       标示FastCGI记录的类型。也就是 记录执行的一般功能.具体的功能后面有描述
requestId:  标示该记录属于哪一个FastCGI请求
contentLength: 后面的contentData元素的字节个数
paddingLength: 后面的paddingData元素的字节个数
contentData:   在0到65535之间，不同的记录类型有不同的意义
padddingData:  在0到255之间，填充字节。可以被忽略

很多情况下，FastCGI应用需要读取变长的值, 名-值对的格式为
名字长度+值长度+名字+值


