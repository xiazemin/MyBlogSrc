---
title: UNIX Domain Socket
layout: post
category: linux
author: 夏泽民
---
UNIX Domain Socket, 简称UDS, 
UDS的优势:

UDS传输不需要经过网络协议栈,不需要打包拆包等操作,只是数据的拷贝过程
UDS分为SOCK_STREAM(流套接字)和SOCK_DGRAM(数据包套接字),由于是在本机通过内核通信,不会丢包也不会出现发送包的次序和接收包的次序不一致的问题
流程介绍
如果熟悉Socket的话,UDS也是同样的方式, 区别如下:

UDS不需要IP和Port, 而是通过一个文件名来表示
domain 为 AF_UNIX
UDS中使用sockaddr_un表示
struct sockaddr_un {
    sa_family_t sun_family; /* AF_UNIX */
    char sun_path[UNIX_PATH_MAX];   /* pathname */
};
SOCK_STREAM(流) : 提供有序，可靠的双向连接字节流。 可以支持带外数据传输机制, 
无论多大的数据都不会截断 
SOCK_DGRAM(数据报)：支持数据报（固定最大长度的无连接，不可靠的消息),数据报超过最大长度,会被截断.
<!-- more -->
一、 概述

UNIX Domain Socket是在socket架构上发展起来的用于同一台主机的进程间通讯（IPC），它不需要经过网络协议栈，不需要打包拆包、计算校验和、维护序号和应答等，只是将应用层数据从一个进程拷贝到另一个进程。UNIX Domain Socket有SOCK_DGRAM或SOCK_STREAM两种工作模式，类似于UDP和TCP，但是面向消息的UNIX Domain Socket也是可靠的，消息既不会丢失也不会顺序错乱。

UNIX Domain Socket可用于两个没有亲缘关系的进程，是全双工的，是目前使用最广泛的IPC机制，比如X Window服务器和GUI程序之间就是通过UNIX Domain Socket通讯的。

二、工作流程
UNIX Domain socket与网络socket类似，可以与网络socket对比应用。
上述二者编程的不同如下：
address family为AF_UNIX
因为应用于IPC，所以UNIXDomain socket不需要IP和端口，取而代之的是文件路径来表示“网络地址”。这点体现在下面两个方面。
地址格式不同，UNIXDomain socket用结构体sockaddr_un表示，是一个socket类型的文件在文件系统中的路径，这个socket文件由bind()调用创建，如果调用bind()时该文件已存在，则bind()错误返回。
UNIX Domain Socket客户端一般要显式调用bind函数，而不象网络socket一样依赖系统自动分配的地址。客户端bind的socket文件名可以包含客户端的pid，这样服务器就可以区分不同的客户端。
UNIX Domain socket的工作流程简述如下（与网络socket相同）。
服务器端：创建socket—绑定文件（端口）—监听—接受客户端连接—接收/发送数据—…—关闭
客户端：创建socket—绑定文件（端口）—连接—发送/接收数据—…—关闭

1. unix域的数据报服务是否可靠

        man unix 手册即可看到，unix domain socket 的数据报既不会丢失也不会乱序 （据我所知，在Linux下的确是这样）。不过最新版本的内核，仍然又提供了一个保证次序的类型 “ kernel 2.6.4 SOCK_SEQPACKET ”。
2. STREAM 和 DGRAM 的主要区别

        既然数据报不丢失也可靠，那不是和 STREAM 很类似么？我理解也确实是这样，而且我觉得 DGRAM 相对还要好一些，因为发送的数据可以带边界。二者另外的区别在于收发时的数据量不一样，基于 STREAM 的套接字，send 可以传入超过 SO_SNDBUF 长的数据，recv 时同 TCP 类似会存在数据粘连。

        采用阻塞方式使用API，在unix domain socket 下调用 sendto 时，如果缓冲队列已满，会阻塞。而UDP因为不是可靠的，无法感知对端的情况，即使对端没有及时收取数据，基本上sendto都能立即返回成功（如果发端疯狂sendto就另当别论，因为过快地调用sendto在慢速网络的环境下，可能撑爆套接字的缓冲区，导致sendto阻塞）。
3. SO_SNDBUF 和 SO_REVBUF

        对于 unix domain socket，设置 SO_SNDBUF 会影响 sendto 最大的报文长度，但是任何针对 SO_RCVBUF 的设置都是无效的 。实际上 unix domain socket 的数据报还是得将数据放入内核所申请的内存块里面，再由另一个进程通过 recvfrom 从内核读取，因此具体可以发送的数据报长度受限于内核的 slab 策略 。在 linux 平台下，早先版本（如 2.6.2）可发送最大数据报长度约为 128 k ，新版本的内核支持更大的长度。
4. 使用 DGRAM 时，缓冲队列的长度

        有几个因素会影响缓冲队列的长度，一个是上面提到的 slab 策略，另一个则是系统的内核参数 /proc/sys/net/unix/max_dgram_qlen。缓冲队列长度是这二者共同决定的。

        如 max_dgram_qlen 默认为 10，在数据报较小时（如1k），先挂起接收数据的进程后，仍可以 sendto 10 次并顺利返回；

        但是如果数据报较大（如120k）时，就要看 slab “size-131072” 的 limit 了。
5. 使用 unix domain socket 进行进程间通信 vs 其他方式
        · 需要先确定操作系统类型，以及其所对应的最大 DGRAM 长度，如果有需要传送超过该长度的数据报，建议拆分成几个发送，接收后组装即可（不会乱序，个人觉得这样做比用 STREAM 再切包方便得多）

        · 同管道相比，unix 域的数据报不但可以维持数据的边界，还不会碰到在写入管道时的原子性问题。

        · 同共享内存相比，不能独立于进程缓存大量数据，但是却避免了同步互斥的考量。

        · 同普通 socket 相比，开销相对较小（不用计算报头），DGRAM 的报文长度可以大于 64k，不过不能像普通 socket 那样将进程切换到不同机器 。
6. 其他
        其实在本机 IPC 时，同普通 socket 的 UDP 相比，unix domain socket 的数据报只不过是在收发时分别少计算了一下校验和而已，本机的 UDP 会走 lo 接口，不会进行 IP 分片，也不会真正跑到网卡的链路层上去（不会占用网卡硬件） 。也就是说，在本机上使用普通的 socket UDP，只是多耗了一些 CPU（之所以说一些，是因为校验和的计算很简单），此外本机的 UDP 也可以保证数据不丢失、不乱序 。

Unix domain socket
    Unix domain socket 或者 IPC socket是一种终端，可以使同一台操作系统上的两个或多个进程进行数据通信。与管道相比，Unix domain sockets 既可以使用字节流，又可以使用数据队列，而管道通信则只能使用字节流。Unix domain sockets的接口和Internet socket很像，但它不使用网络底层协议来通信。Unix domain socket 的功能是POSIX操作系统里的一种组件。

   Unix domain sockets 使用系统文件的地址来作为自己的身份。它可以被系统进程引用。所以两个进程可以同时打开一个Unix domain sockets来进行通信。不过这种通信方式是发生在系统内核里而不会在网络里传播。

   UNIX Domain Socket是在socket架构上发展起来的用于同一台主机的进程间通讯（IPC），它不需要经过网络协议栈，不需要打包拆包、计算校验和、维护序号和应答等，只是将应用层数据从一个进程拷贝到另一个进程。UNIX Domain Socket有SOCK_DGRAM或SOCK_STREAM两种工作模式，类似于UDP和TCP，但是面向消息的UNIX Domain Socket也是可靠的，消息既不会丢失也不会顺序错乱。UNIX Domain Socket可用于两个没有亲缘关系的进程，是全双工的，是目前使用最广泛的IPC机制，比如X Window服务器和GUI程序之间就是通过UNIX Domain Socket通讯的。

     因为应用于IPC，所以UNIXDomain socket不需要IP和端口，取而代之的是文件路径来表示“网络地址”，这点体现在下面两个方面：


   1. 地址格式不同，UNIXDomain socket用结构体sockaddr_un表示，是一个socket类型的文件在文件系统中的路径，这个socket文件由bind()调用创建，如果调用bind()时该文件已存在，则bind()错误返回。                                                                                                  
  2. UNIX Domain Socket客户端一般要显式调用bind函数，而不象网络socket一样依赖系统自动分配的地址。客户端bind的socket文件名可以包含客户端的pid，这样服务器就可以区分不同的客户端。

三、相关API

int socket(int domain, int type, int protocol)

domain:说明我们网络程序所在的主机采用的通讯协族(AF_UNIX和AF_INET等). AF_UNIX只能够用于单一的Unix系统进程间通信,而AF_INET是针对               Internet的,因而可以允许在远程主机之间通信
type:我们网络程序所采用的通讯协议(SOCK_STREAM,SOCK_DGRAM等) SOCK_STREAM表明我们用的是TCP协议,这样会提供按顺序的,可靠,双向,面向连接的       比特流. SOCK_DGRAM 表明我们用的是UDP协议,这样只会提供定长的,不可靠,无连接的通信.
protocol:由于我们指定了type,所以这个地方我们一般只要用0来代替就可以了
socket为网络通讯做基本的准备.成功时返回文件描述符,失败时返回-1,看errno可知道出错的详细情况


int bind( int sockfd , const struct sockaddr * my_addr, socklen_t addrlen)

sockfd:是由socket调用返回的文件描述符.
addrlen:是sockaddr结构的长度.
my_addr:是一个指向sockaddr的指针. 在中有 sockaddr的定义
struct sockaddr{
unisgned short as_family;
char sa_data[14];
};

int listen(int sockfd,int backlog)


sockfd:是bind后的文件描述符.
backlog:设置请求排队的最大长度.当有多个客户端程序和服务端相连时, 使用这个表示可以介绍的排队长度. listen函数将bind的文件描述符变为监听套接字.返回的情况和bind一样.

int accept(int sockfd, struct sockaddr *addr,int *addrlen)
sockfd:是listen后的文件描述符.
addr,addrlen是用来给客户端的程序填写的,服务器端只要传递指针就可以了. bind,listen和accept是服务器端用的函数,accept调用时,服务器端的程序会一直阻塞到有一个客户程序发出了连接. accept成功时返回最后的服务器端的文件描述符,这个时候服务器端可以向该描述符写信息了. 失败时返回-1

过程：有人从很远的地方通过一个你在侦听(listen())的端口连接(connect())到你的机器。它的连接将加入到等待接受(accept())的队列中。你调用accept()告诉它你有空闲的连接。它将返回一个新的套接字文件描述符！这样你就有两个套接字了，原来的一个还在侦听你的那个端口，新的在准备发送(send())和接收(recv())数据。这就是Linux Accept函数的过程！

Ps：Linux Accept函数注意事项，在系统调用send()和recv()中你应该使用新的套接字描述符new_fd。如果你只想让一个连接进来，那么你可以使用close()去关闭原来的文件描述符sockfd来避免同一个端口更多的连接。

int connect(int sockfd, struct sockaddr * serv_addr,int addrlen)
sockfd:socket返回的文件描述符.
serv_addr:储存了服务器端的连接信息.其中sin_add是服务端的地址
addrlen:serv_addr的长度
connect函数是客户端用来同服务端连接的.成功时返回0,sockfd是同服务端通讯的文件描述符失败时返回-1

ssize_t recv(int sockfd, void *buff, size_t nbytes, int flags);

ssize_t send(int sockfd, const void *buff, size_t nbytes, int flags)

recv 和send的前3个参数等同于read和write。

flags参数值为0或：

 
flags	说明	recv	send
 MSG_DONTROUTE	绕过路由表查找 	 	  •
 MSG_DONTWAIT	仅本操作非阻塞 	  •    	  •
 MSG_OOB　　　　	发送或接收带外数据	  •	  •
 MSG_PEEK　　	窥看外来消息	  •	 
 MSG_WAITALL　　	等待所有数据 	  •	 
 1. send解析

 sockfd：指定发送端套接字描述符。

 buff：    存放要发送数据的缓冲区

 nbytes:  实际要改善的数据的字节数

 flags：   一般设置为0

 1) send先比较发送数据的长度nbytes和套接字sockfd的发送缓冲区的长度，如果nbytes > 套接字sockfd的发送缓冲区的长度, 该函数返回SOCKET_ERROR;

 2) 如果nbtyes <= 套接字sockfd的发送缓冲区的长度，那么send先检查协议是否正在发送sockfd的发送缓冲区中的数据，如果是就等待协议把数据发送完，如果协议还没有开始发送sockfd的发送缓冲区中的数据或者sockfd的发送缓冲区中没有数据，那么send就比较sockfd的发送缓冲区的剩余空间和nbytes

 3) 如果 nbytes > 套接字sockfd的发送缓冲区剩余空间的长度，send就一起等待协议把套接字sockfd的发送缓冲区中的数据发送完

 4) 如果 nbytes < 套接字sockfd的发送缓冲区剩余空间大小，send就仅仅把buf中的数据copy到剩余空间里(注意并不是send把套接字sockfd的发送缓冲区中的数据传到连接的另一端的，而是协议传送的，send仅仅是把buf中的数据copy到套接字sockfd的发送缓冲区的剩余空间里)。

 5) 如果send函数copy成功，就返回实际copy的字节数，如果send在copy数据时出现错误，那么send就返回SOCKET_ERROR; 如果在等待协议传送数据时网络断开，send函数也返回SOCKET_ERROR。

 6) send函数把buff中的数据成功copy到sockfd的改善缓冲区的剩余空间后它就返回了，但是此时这些数据并不一定马上被传到连接的另一端。如果协议在后续的传送过程中出现网络错误的话，那么下一个socket函数就会返回SOCKET_ERROR。（每一个除send的socket函数在执行的最开始总要先等待套接字的发送缓冲区中的数据被协议传递完毕才能继续，如果在等待时出现网络错误那么该socket函数就返回SOCKET_ERROR）

 7) 在unix系统下，如果send在等待协议传送数据时网络断开，调用send的进程会接收到一个SIGPIPE信号，进程对该信号的处理是进程终止。

2.recv函数

sockfd: 接收端套接字描述符

buff：   用来存放recv函数接收到的数据的缓冲区

nbytes: 指明buff的长度

flags:   一般置为0

 1) recv先等待s的发送缓冲区的数据被协议传送完毕，如果协议在传送sock的发送缓冲区中的数据时出现网络错误，那么recv函数返回SOCKET_ERROR

 2) 如果套接字sockfd的发送缓冲区中没有数据或者数据被协议成功发送完毕后，recv先检查套接字sockfd的接收缓冲区，如果sockfd的接收缓冲区中没有数据或者协议正在接收数据，那么recv就一起等待，直到把数据接收完毕。当协议把数据接收完毕，recv函数就把s的接收缓冲区中的数据copy到buff中（注意协议接收到的数据可能大于buff的长度，所以在这种情况下要调用几次recv函数才能把sockfd的接收缓冲区中的数据copy完。recv函数仅仅是copy数据，真正的接收数据是协议来完成的）

 3) recv函数返回其实际copy的字节数，如果recv在copy时出错，那么它返回SOCKET_ERROR。如果recv函数在等待协议接收数据时网络中断了，那么它返回0。

 4) 在unix系统下，如果recv函数在等待协议接收数据时网络断开了，那么调用 recv的进程会接收到一个SIGPIPE信号，进程对该信号的默认处理是进程终止。
