---
title: tcp_retries
layout: post
category: linux
author: 夏泽民
---
https://pracucci.com/linux-tcp-rto-min-max-and-tcp-retries2.html
TCP retransmits an unacknowledged packet up to tcp_retries2 sysctl setting times (defaults to 15) using an exponential backoff timeout for which each retransmission timeout is between TCP_RTO_MIN (200 ms) and TCP_RTO_MAX (120 seconds). Once the 15th retry expires (by default), the TCP stack will notify the layers above (ie. app) of a broken connection.

The value of TCP_RTO_MIN and TCP_RTO_MAX is hardcoded in the Linux kernel and defined by the following constants:

#define TCP_RTO_MAX ((unsigned)(120*HZ))
#define TCP_RTO_MIN ((unsigned)(HZ/5))
<!-- more -->
http://perthcharles.github.io/2015/09/06/wiki-rtt-estimator/

http://perthcharles.github.io/2015/09/07/wiki-tcp-retries/

Linux中确实定义了两个参数来限定超时重传的次数的，以下是源码中Documentation/networking/ip-sysctl.txt文档中的描述

tcp_retries1 - INTEGER
    This value influences the time, after which TCP decides, that
    something is wrong due to unacknowledged RTO retransmissions,
    and reports this suspicion to the network layer.
    See tcp_retries2 for more details.

    RFC 1122 recommends at least 3 retransmissions, which is the
    default.

tcp_retries2 - INTEGER
    This value influences the timeout of an alive TCP connection,
    when RTO retransmissions remain unacknowledged.
    Given a value of N, a hypothetical TCP connection following
    exponential backoff with an initial RTO of TCP_RTO_MIN would
    retransmit N times before killing the connection at the (N+1)th RTO.

    The default value of 15 yields a hypothetical timeout of 924.6
    seconds and is a lower bound for the effective timeout.
    TCP will effectively time out at the first RTO which exceeds the
    hypothetical timeout.

    RFC 1122 recommends at least 100 seconds for the timeout,
    which corresponds to a value of at least 8.
就是这样一段话，可能由于过于概括，会令人产生很多疑问，甚至产生一些误解。
比如常见的问题有：
a. 超过tcp_retries1这个阈值后，到底是report了怎样一种suspicion呢？
b. tcp_retries1和tcp_retries2的数字是表示RTO重传的次数上限，对吗？
c. 文档中提到，924.6s is a lower bound for the effective timeout。
这里的effective timeout是指什么？
为什么是lower bound，tcp_retries2不应该是限制重传次数的upper bound吗？

// RTO timer的处理函数是tcp_retransmit_timer()，与tcp_retries1相关的代码调用关系如下  
tcp_retransmit_timer()
    => tcp_write_timeout()  // 判断是否重传了足够的久
        => retransmit_timed_out(sk, sysctl_tcp_retries1, 0, 0)  // 判断是否超过了阈值

// tcp_write_timeout()的具体相关内容  
...
if ((1 << sk->sk_state) & (TCPF_SYN_SENT | TCPF_SYN_RECV)) {
    // 如果超时发生在三次握手期间，此时有专门的tcp_syn_retries来负责限定重传次数
    ...
} else {    // 如果超时发生在数据发送期间
    // 这个函数负责判断重传是否超过阈值，返回真表示超过。后续会详细分析这个函数  
    if (retransmits_timed_out(sk, sysctl_tcp_retries1, 0, 0)) { 
        /* Black hole detection */
        tcp_mtu_probing(icsk, sk);  // 如果开启tcp_mtu_probing(默认关闭)了，则执行PMTU

        dst_negative_advice(sk);    // 更新路由缓存
    }
    ...
}
从以上的代码可以看到，一旦重传超过阈值tcp_retries1，主要的动作就是更新路由缓存。
用以避免由于路由选路变化带来的问题。

重传超过tcp_retries2会怎样
会直接放弃重传，关闭TCP流

// 依然还是在tcp_write_timeout()中，retry_until一般是tcp_retries2
...
if (retransmits_timed_out(sk, retry_until, syn_set ? 0 : icsk->icsk_user_timeout, syn_set)) {
    /* Has it gone just too far? */
    tcp_write_err(sk);      // 调用tcp_done关闭TCP流
    return 1;
}
retries限制的重传次数吗
咋一看文档，很容易想到retries的数字就是限定的重传的次数，甚至源码中对于retries常量注释中都写着”This is how many retries it does…”

#define TCP_RETR1       3   /*
                             * This is how many retries it does before it
                             * tries to figure out if the gateway is
                             * down. Minimal RFC value is 3; it corresponds
                             * to ~3sec-8min depending on RTO.
                             */

#define TCP_RETR2       15  /*
                             * This should take at least
                             * 90 minutes to time out.
                             * RFC1122 says that the limit is 100 sec.
                             * 15 is ~13-30min depending on RTO.
                             */
那就就来看看retransmits_timed_out的具体实现，看看到底是不是限制的重传次数

/* This function calculates a "timeout" which is equivalent to the timeout of a
 * TCP connection after "boundary" unsuccessful, exponentially backed-off
 * retransmissions with an initial RTO of TCP_RTO_MIN or TCP_TIMEOUT_INIT if
 * syn_set flag is set.
 */
static bool retransmits_timed_out(struct sock *sk,
                              unsigned int boundary,
                              unsigned int timeout,
                              bool syn_set)
{
    unsigned int linear_backoff_thresh, start_ts;
    // 如果是在三次握手阶段，syn_set为真
    unsigned int rto_base = syn_set ? TCP_TIMEOUT_INIT : TCP_RTO_MIN;

    if (!inet_csk(sk)->icsk_retransmits)
            return false;

    // retrans_stamp记录的是数据包第一次发送的时间，在tcp_retransmit_skb()中设置
    if (unlikely(!tcp_sk(sk)->retrans_stamp))
            start_ts = TCP_SKB_CB(tcp_write_queue_head(sk))->when;
    else
            start_ts = tcp_sk(sk)->retrans_stamp;

    // 如果用户态未指定timeout，则算一个出来
    if (likely(timeout == 0)) {
            /* 下面的计算过程，其实就是算一下如果以rto_base为第一次重传间隔，
             * 重传boundary次需要多长时间
             */
            linear_backoff_thresh = ilog2(TCP_RTO_MAX/rto_base);

            if (boundary <= linear_backoff_thresh)
                    timeout = ((2 << boundary) - 1) * rto_base;
            else
                    timeout = ((2 << linear_backoff_thresh) - 1) * rto_base +
                            (boundary - linear_backoff_thresh) * TCP_RTO_MAX;
    }
    // 如果数据包第一次发送的时间距离现在的时间间隔，超过了timeout值，则认为重传超于阈值了
    return (tcp_time_stamp - start_ts) >= timeout;
}
从以上的代码分析可以看到，真正起到限制重传次数的并不是真正的重传次数。
而是以tcp_retries1或tcp_retries2为boundary，以rto_base(如TCP_RTO_MIN 200ms)为初始RTO，计算得到一个timeout值出来。如果重传间隔超过这个timeout，则认为超过了阈值。
上面这段话太绕了，下面举两个个例子来说明

以判断是否放弃TCP流为例，如果tcp_retries2=15，那么计算得到的timeout=924600ms。

1. 如果RTT比较小，那么RTO初始值就约等于下限200ms
   由于timeout总时长是924600ms，表现出来的现象刚好就是重传了15次，超过了timeout值，从而放弃TCP流

2. 如果RTT较大，比如RTO初始值计算得到的是1000ms
   那么根本不需要重传15次，重传总间隔就会超过924600ms。
   比如我测试的一个RTT=400ms的情况，当tcp_retries2=10时，仅重传了3次就放弃了TCP流
另外几个小问题
理解了Linux决定重传次数的真实机制，就不难回答一下几个问题了

>> effective timeout指的是什么？  
<< 就是retransmits_timed_out计算得到的timeout值

>> 924.6s是怎么算出来的？
<< 924.6s = (( 2 << 9) -1) * 200ms + (15 - 9) * 120s

>> 为什么924.6s是lower bound？
<< 重传总间隔必须大于timeout值，即 (tcp_time_stamp - start_ts) >= timeout

>> 那RTO超时的间隔到底是不是源码注释的"15 is ~13-30min depending on RTO."呢？  
<< 显然不是! 虽然924.6s(15min)是一个lower bound，但是它同时也是一个upper bound!
   怎么理解？举例说明  
        1. 如果某个RTO值导致，在已经重传了14次后，总重传间隔开销是924s
        那么它还需要重传第15次，即使离924.6s只差0.6s。这就是发挥了lower bound的作用
        2. 如果某个RTO值导致，在重传了10次后，总重传间隔开销是924s
        重传第11次后，第12次超时触发时计算得到的总间隔变为1044s，超过924.6s
        那么此时会放弃第12次重传，这就是924.6s发挥了upper bound的作用
   总的来说，在Linux3.10中，如果tcp_retres2设置为15。总重传超时周期应该在如下范围内
        [924.6s, 1044.6s)
所以综合上述，Linux并不是直接拿tcp_retries1和tcp_retries2来限制重传次数的，而是用计算得到
的一个timeout值来判断是否要放弃重传的。真正的重传次数同时与RTT相关。

/proc/sys/net目录
　　所有的TCP/IP参数都位于/proc/sys/net目录下（请注意，对/proc/sys/net目录下内容的修改都是临时的，任何修改在系统重启后都会丢失），例如下面这些重要的参数：

参数（路径+文件）

描述

默认值

优化值

/proc/sys/net/core/rmem_default

默认的TCP数据接收窗口大小（字节）。

229376

256960

/proc/sys/net/core/rmem_max

最大的TCP数据接收窗口（字节）。

131071

513920

/proc/sys/net/core/wmem_default

默认的TCP数据发送窗口大小（字节）。

229376

256960

/proc/sys/net/core/wmem_max

最大的TCP数据发送窗口（字节）。

131071

513920

/proc/sys/net/core/netdev_max_backlog

在每个网络接口接收数据包的速率比内核处理这些包的速率快时，允许送到队列的数据包的最大数目。

1000

2000

/proc/sys/net/core/somaxconn

定义了系统中每一个端口最大的监听队列的长度，这是个全局的参数。

128

2048

/proc/sys/net/core/optmem_max

表示每个套接字所允许的最大缓冲区的大小。

20480

81920

/proc/sys/net/ipv4/tcp_mem

确 定TCP栈应该如何反映内存使用，每个值的单位都是内存页（通常是4KB）。第一个值是内存使用的下限；第二个值是内存压力模式开始对缓冲区使用应用压力 的上限；第三个值是内存使用的上限。在这个层次上可以将报文丢弃，从而减少对内存的使用。对于较大的BDP可以增大这些值（注意，其单位是内存页而不是字 节）。

94011  125351  188022

131072  262144  524288

/proc/sys/net/ipv4/tcp_rmem

为 自动调优定义socket使用的内存。第一个值是为socket接收缓冲区分配的最少字节数；第二个值是默认值（该值会被rmem_default覆 盖），缓冲区在系统负载不重的情况下可以增长到这个值；第三个值是接收缓冲区空间的最大字节数（该值会被rmem_max覆盖）。

4096  87380  4011232

8760  256960  4088000

/proc/sys/net/ipv4/tcp_wmem

为 自动调优定义socket使用的内存。第一个值是为socket发送缓冲区分配的最少字节数；第二个值是默认值（该值会被wmem_default覆 盖），缓冲区在系统负载不重的情况下可以增长到这个值；第三个值是发送缓冲区空间的最大字节数（该值会被wmem_max覆盖）。

4096  16384  4011232

8760  256960  4088000

/proc/sys/net/ipv4/tcp_keepalive_time

TCP发送keepalive探测消息的间隔时间（秒），用于确认TCP连接是否有效。

7200

1800

/proc/sys/net/ipv4/tcp_keepalive_intvl

探测消息未获得响应时，重发该消息的间隔时间（秒）。

75

30

/proc/sys/net/ipv4/tcp_keepalive_probes

在认定TCP连接失效之前，最多发送多少个keepalive探测消息。

9

3

/proc/sys/net/ipv4/tcp_sack

启用有选择的应答（1表示启用），通过有选择地应答乱序接收到的报文来提高性能，让发送者只发送丢失的报文段，（对于广域网通信来说）这个选项应该启用，但是会增加对CPU的占用。

1

1

/proc/sys/net/ipv4/tcp_fack

启用转发应答，可以进行有选择应答（SACK）从而减少拥塞情况的发生，这个选项也应该启用。

1

1

/proc/sys/net/ipv4/tcp_timestamps

TCP时间戳（会在TCP包头增加12个字节），以一种比重发超时更精确的方法（参考RFC 1323）来启用对RTT 的计算，为实现更好的性能应该启用这个选项。

1

1

/proc/sys/net/ipv4/tcp_window_scaling

启用RFC 1323定义的window scaling，要支持超过64KB的TCP窗口，必须启用该值（1表示启用），TCP窗口最大至1GB，TCP连接双方都启用时才生效。

1

1

/proc/sys/net/ipv4/tcp_syncookies

表示是否打开TCP同步标签（syncookie），内核必须打开了CONFIG_SYN_COOKIES项进行编译，同步标签可以防止一个套接字在有过多试图连接到达时引起过载。

1

1

/proc/sys/net/ipv4/tcp_tw_reuse

表示是否允许将处于TIME-WAIT状态的socket（TIME-WAIT的端口）用于新的TCP连接 。

0

1

/proc/sys/net/ipv4/tcp_tw_recycle

能够更快地回收TIME-WAIT套接字。

0

1

/proc/sys/net/ipv4/tcp_fin_timeout

对于本端断开的socket连接，TCP保持在FIN-WAIT-2状态的时间（秒）。对方可能会断开连接或一直不结束连接或不可预料的进程死亡。

60

30

/proc/sys/net/ipv4/ip_local_port_range

表示TCP/UDP协议允许使用的本地端口号

32768  61000

1024  65000

/proc/sys/net/ipv4/tcp_max_syn_backlog

对于还未获得对方确认的连接请求，可保存在队列中的最大数目。如果服务器经常出现过载，可以尝试增加这个数字。

2048

2048

/proc/sys/net/ipv4/tcp_low_latency

允许TCP/IP栈适应在高吞吐量情况下低延时的情况，这个选项应该禁用。

0

 

/proc/sys/net/ipv4/tcp_westwood

启用发送者端的拥塞控制算法，它可以维护对吞吐量的评估，并试图对带宽的整体利用情况进行优化，对于WAN 通信来说应该启用这个选项。

0

 

/proc/sys/net/ipv4/tcp_bic

为快速长距离网络启用Binary Increase Congestion，这样可以更好地利用以GB速度进行操作的链接，对于WAN通信应该启用这个选项。

1

 

 

/etc/sysctl.conf文件

　　/etc /sysctl.conf是一个允许你改变正在运行中的Linux系统的接口。它包含一些TCP/IP堆栈和虚拟内存系统的高级选项，可用来控制 Linux网络配置，由于/proc/sys/net目录内容的临时性，建议把TCPIP参数的修改添加到/etc/sysctl.conf文件, 然后保存文件，使用命令“/sbin/sysctl –p”使之立即生效。具体修改方案参照上文：

net.core.rmem_default = 256960

net.core.rmem_max = 513920

net.core.wmem_default = 256960

net.core.wmem_max = 513920

net.core.netdev_max_backlog = 2000

net.core.somaxconn = 2048

net.core.optmem_max = 81920

net.ipv4.tcp_mem = 131072  262144  524288

net.ipv4.tcp_rmem = 8760  256960  4088000

net.ipv4.tcp_wmem = 8760  256960  4088000

net.ipv4.tcp_keepalive_time = 1800

net.ipv4.tcp_keepalive_intvl = 30

net.ipv4.tcp_keepalive_probes = 3

net.ipv4.tcp_sack = 1

net.ipv4.tcp_fack = 1

net.ipv4.tcp_timestamps = 1

net.ipv4.tcp_window_scaling = 1

net.ipv4.tcp_syncookies = 1

net.ipv4.tcp_tw_reuse = 1

net.ipv4.tcp_tw_recycle = 1

net.ipv4.tcp_fin_timeout = 30

net.ipv4.ip_local_port_range = 1024  65000

net.ipv4.tcp_max_syn_backlog = 2048

 

Doc2：

   可调优的内核变量存在两种主要接口：sysctl命令和/proc文件系统，proc中与进程无关的所有信息都被移植到sysfs中。IPV4协议栈的 sysctl参数主要是sysctl.net.core、sysctl.net.ipv4，对应的/proc文件系统是/proc/sys/net /ipv4和/proc/sys/net/core。只有内核在编译时包含了特定的属性，该参数才会出现在内核中。

    对于内核参数应该谨慎调节，这些参数通常会影响到系统的整体性能。内核在启动时会根据系统的资源情况来初始化特定的变量，这种初始化的调节一般会满足通常的性能需求。

    应用程序通过socket系统调用和远程主机进行通讯，每一个socket都有一个读写缓冲区。读缓冲区保存了远程主机发送过来的数据，如果缓冲区已满， 则数据会被丢弃，写缓冲期保存了要发送到远程主机的数据，如果写缓冲区已慢，则系统的应用程序在写入数据时会阻塞。可知，缓冲区是有大小的。

socket缓冲区默认大小：
/proc/sys/net/core/rmem_default     对应net.core.rmem_default
/proc/sys/net/core/wmem_default     对应net.core.wmem_default
    上面是各种类型socket的默认读写缓冲区大小，然而对于特定类型的socket则可以设置独立的值覆盖默认值大小。例如tcp类型的socket就可以用/proc/sys/net/ipv4/tcp_rmem和tcp_wmem来覆盖。

socket缓冲区最大值：
/proc/sys/net/core/rmem_max        对应net.core.rmem_max
/proc/sys/net/core/wmem_max        对应net.core.wmem_max

/proc/sys/net/core/netdev_max_backlog    对应 net.core.netdev_max_backlog
    该参数定义了当接口收到包的速率大于内核处理包的速率时，设备的输入队列中的最大报文数。

/proc/sys/net/core/somaxconn        对应 net.core.somaxconn
    通过listen系统调用可以指定的最大accept队列backlog，当排队的请求连接大于该值时，后续进来的请求连接会被丢弃。

/proc/sys/net/core/optmem_max          对应 net.core.optmem_max
    每个socket的副缓冲区大小。

TCP/IPV4内核参数：
    在创建socket的时候会指定socke协议和地址类型。TCP socket缓冲区大小是他自己控制而不是由core内核缓冲区控制。
/proc/sys/net/ipv4/tcp_rmem     对应net.ipv4.tcp_rmem
/proc/sys/net/ipv4/tcp_wmem     对应net.ipv4.tcp_wmem
    以上是TCP socket的读写缓冲区的设置，每一项里面都有三个值，第一个值是缓冲区最小值，中间值是缓冲区的默认值，最后一个是缓冲区的最大值，虽然缓冲区的值不受core缓冲区的值的限制，但是缓冲区的最大值仍旧受限于core的最大值。

/proc/sys/net/ipv4/tcp_mem  
    该内核参数也是包括三个值，用来定义内存管理的范围，第一个值的意思是当page数低于该值时，TCP并不认为他为内存压力，第二个值是进入内存的压力区 域时所达到的页数，第三个值是所有TCP sockets所允许使用的最大page数，超过该值后，会丢弃后续报文。page是以页面为单位的，为系统中socket全局分配的内存容量。

socket的结构如下图：



/proc/sys/net/ipv4/tcp_window_scaling      对应net.ipv4.tcp_window_scaling
    管理TCP的窗口缩放特性，因为在tcp头部中声明接收缓冲区的长度为26位，因此窗口不能大于64K，如果大于64K，就要打开窗口缩放。

/proc/sys/net/ipv4/tcp_sack    对应net.ipv4.tcp_sack
    管理TCP的选择性应答，允许接收端向发送端传递关于字节流中丢失的序列号，减少了段丢失时需要重传的段数目，当段丢失频繁时，sack是很有益的。

/proc/sys/net/ipv4/tcp_dsack   对应net.ipv4.tcp_dsack
    是对sack的改进，能够检测不必要的重传。

/proc/sys/net/ipv4/tcp_fack    对应net.ipv4.tcp_fack
    对sack协议加以完善，改进tcp的拥塞控制机制。

TCP的连接管理：
/proc/sys/net/ipv4/tcp_max_syn_backlog    对应net.ipv4.tcp_max_syn_backlog
    每一个连接请求(SYN报文)都需要排队，直至本地服务器接收，该变量就是控制每个端口的 TCP SYN队列长度的。如果连接请求多余该值，则请求会被丢弃。

/proc/sys/net/ipv4/tcp_syn_retries    对应net.ipv4.tcp_syn_retries
    控制内核向某个输入的SYN/ACK段重新发送相应的次数，低值可以更好的检测到远程主机的连接失败。可以修改为3

/proc/sys/net/ipv4/tcp_retries1    对应net.ipv4.tcp_retries1
    该变量设置放弃回应一个tcp连接请求前，需要进行多少次重试。

/proc/sys/net/ipv4/tcp_retries2    对应net.ipv4.tcp_retries2
    控制内核向已经建立连接的远程主机重新发送数据的次数，低值可以更早的检测到与远程主机失效的连接，因此服务器可以更快的释放该连接，可以修改为5

TCP连接的保持：
/proc/sys/net/ipv4/tcp_keepalive_time        对应net.ipv4.tcp_keepalive_time
    如果在该参数指定的秒数内连接始终处于空闲状态，则内核向客户端发起对该主机的探测

/proc/sys/net/ipv4/tcp_keepalive_intvl    对应net.ipv4.tcp_keepalive_intvl
    该参数以秒为单位，规定内核向远程主机发送探测指针的时间间隔

/proc/sys/net/ipv4/tcp_keepalive_probes   对应net.ipv4.tcp_keepalive_probes
    该参数规定内核为了检测远程主机的存活而发送的探测指针的数量，如果探测指针的数量已经使用完毕仍旧没有得到客户端的响应，即断定客户端不可达，关闭与该客户端的连接，释放相关资源。

/proc/sys/net/ipv4/ip_local_port_range   对应net.ipv4.ip_local_port_range
    规定了tcp/udp可用的本地端口的范围。

TCP连接的回收：
/proc/sys/net/ipv4/tcp_max_tw_buckets     对应net.ipv4.tcp_max_tw_buckets
   该参数设置系统的TIME_WAIT的数量，如果超过默认值则会被立即清除。

/proc/sys/net/ipv4/tcp_tw_reuse           对应net.ipv4.tcp_tw_reuse
   该参数设置TIME_WAIT重用，可以让处于TIME_WAIT的连接用于新的tcp连接

/proc/sys/net/ipv4/tcp_tw_recycle         对应net.ipv4.tcp_tw_recycle
   该参数设置tcp连接中TIME_WAIT的快速回收。

/proc/sys/net/ipv4/tcp_fin_timeout       对应net.ipv4.tcp_fin_timeout
   设置TIME_WAIT2进入CLOSED的等待时间。

/proc/sys/net/ipv4/route/max_size
   内核所允许的最大路由数目。

/proc/sys/net/ipv4/ip_forward
   接口间转发报文

/proc/sys/net/ipv4/ip_default_ttl
   报文可以经过的最大跳数

虚拟内存参数：
/proc/sys/vm/


   在linux kernel 2.6.25之前通过ulimit -n(setrlimit(RLIMIT_NOFILE))设置每个进程的最大打开文件句柄数不能超过NR_OPEN(1024*1024),也就是 100多w(除非重新编译内核)，而在25之后，内核导出了一个sys接口可以修改这个最大值/proc/sys/fs/nr_open。shell里不 能直接更改，是因为登录的时候pam已经从limits.conf中设置了上限，ulimit命令只能在低于上限的范围内发挥了。

Linux中查看socket状态：
cat /proc/net/sockstat #（这个是ipv4的）

sockets: used 137
TCP: inuse 49 orphan 0 tw 3272 alloc 52 mem 46
UDP: inuse 1 mem 0
RAW: inuse 0
FRAG: inuse 0 memory 0
说明：
sockets: used：已使用的所有协议套接字总量
TCP: inuse：正在使用（正在侦听）的TCP套接字数量。其值≤ netstat –lnt | grep ^tcp | wc –l
TCP: orphan：无主（不属于任何进程）的TCP连接数（无用、待销毁的TCP socket数）
TCP: tw：等待关闭的TCP连接数。其值等于netstat –ant | grep TIME_WAIT | wc –l
TCP：alloc(allocated)：已分配（已建立、已申请到sk_buff）的TCP套接字数量。其值等于netstat –ant | grep ^tcp | wc –l
TCP：mem：套接字缓冲区使用量（单位不详。用scp实测，速度在4803.9kB/s时：其值=11，netstat –ant 中相应的22端口的Recv-Q＝0，Send-Q≈400）
UDP：inuse：正在使用的UDP套接字数量
RAW：
FRAG：使用的IP段数量

https://stackoverflow.com/questions/5907527/application-control-of-tcp-retransmission-on-linux

https://access.redhat.com/solutions/726753


