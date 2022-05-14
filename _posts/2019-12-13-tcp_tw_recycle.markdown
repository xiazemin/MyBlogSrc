---
title: tcp_tw_recycle
layout: post
category: linux
author: 夏泽民
---
Linux 只能收到 SYN 包 不能回包
问题
如果用户发现云主机不能登录，例如无法远程 22 端口或其他端口，但是更换网络环境正常，服务端抓包发现客户端发包只有 SYN，没有回包，可以执行 netstat -s |grep rejec 查看下是否是 tcp_timestamps 的问题
[root@hfgo2 ~]# netstat -s |grep rejec
    8316 passive connections rejected because of time stamp
    780 packets rejects in established connections because of timestamp
如果出现很多数据包的 timestamp 被拒绝，则检查下内核参数 tcp_tw_recycle 是否开启，如果开启，将其关闭即可。
[root@hfgo2 ~]# cat  /proc/sys/net/ipv4/tcp_tw_recycle
原因
这个主意是和内核的 2 个参数相关

net.ipv4.tcp_timestamps
tcp_timestamps 会记录数据包的发送时间。基本的步骤如下：

发送方在发送数据时，将一个timestamp(表示发送时间)放在包里面
接收方在收到数据包后，在对应的ACK包中将收到的timestamp返回给发送方(echo back)
发送发收到ACK包后，用当前时刻now - ACK包中的timestamp就能得到准确的RTT
timestamps一个双向的选项，当一方不开启时，两方都将停用timestamps。比如client端发送的SYN包中带有timestamp选项，但server端并没有开启该选项。则回复的 SYN-ACK 将不带 timestamp 选项，同时 client 后续回复的ACK也不会带有 timestamp 选项。如果client发送的SYN包中就不带 timestamp，双向都将停用 timestamp。

tcp数据包中 timestamps 的 value 是系统开机时间到现在时间的（毫秒级）时间戳。

如果用户是在一个NAT 环境下，或者出口IP 为1个，如果同时要多个用户连接云服务器，则可能会出现这种问题，根据上述SYN包处理规则，在tcp_tw_recycle和tcp_timestamps同时开启的条件下，timestamp大的主机访问serverN成功，而timestmap小的主机访问失败。

net.ipv4.tcp_tw_recycle
TCP 规范中规定的处于 TIME_WAIT 的 TCP 连接必须等待 2msl 时间。但在linux中，如果开启了 tcp_tw_recycle，TIME_WAIT 的 TCP 连接就不会等待 2msl 时间（而是rto或者60s），从而达到快速重用（回收）处于 TIME_WAIT 状态的 tcp 连接的目的。这就可能导致连接收到之前连接的数据。为此，linux 在打开 tcp_tw_recycle 的情况下，会记录下 TIME_WAIT 连接的对端（peer）信息，包括IP地址、时间戳等。这样，当内核收到同一个 IP 的 SYN 包时，就会去比较时间戳，检查 SYN 包的时间戳是否滞后，如果滞后，就将其丢掉（认为是旧连接的数据）。这在绝大部分情况下是没有问题的，但是对于我们实际的 client-server 的服务，访问我们服务的用户一般都位于 NAT之后，如果NAT之后有多个用户访问同一个服务，就有可能因为时间戳滞后的连接被丢掉。

解决办法
 # echo "0" > /proc/sys/net/ipv4/tcp_tw_recycle
<!-- more -->
TIME_WAIT 是 TCP 协议栈中比较特殊的状态，其主要目的是保证不同的链接不会相互干扰，但是对于一些高性能的场景，就可能由于较多的 TIME_WAIT 状态最终导致链接不可用。

TIME_WAIT 是 TCP 协议栈中比较特殊的状态，其主要目的是保证不同的链接不会相互干扰，但是对于一些高性能的场景，就可能由于较多的 TIME_WAIT 状态最终导致链接不可用。

如下简单介绍如何充分利用该状态。

简介
如下是 TCP 的状态转换图，主动关闭链接的一端会进入 TIME_WAIT 状态。

TCP/IP Finite State Machine FSM

可以通过如下命令统计当前不同状态的链接数。

----- 建议使用后者，要快很多
# netstat -ant | awk '/^tcp/ {++state[$NF]} END {for(key in state) print key,state[key]}' | sort -rnk2
# ss -ant | awk '!/^State/ {++state[$1]} END {for(key in state) print key,state[key]}' | sort -rnk2
如上，TIME_WAIT 是在发起主动关闭一方，完成四次挥手后 TCP 状态转换为 TIME_WAIT 状态，并且该状态会保持两个 MSL 。

在 Linux 里一个 MSL 为 30s，不可配置，在 include/net/tcp.h 中定义如下。

#define TCP_TIMEWAIT_LEN (60*HZ) /* how long to wait to destroy TIME-WAIT
                                  * state, about 60 seconds */
原因
之所以采用两个 MSL 主要是为了可靠安全的关闭 TCP 连接。

报文延迟
最常见的是为了防止前一个链接的延迟报文被下一个链接接收，当然这两个链接的四元组 (source address, source port, destination address, destination port) 是相同的，而且序列号 sequence 也是在特定范围内的。

虽然满足上述场景的概率很小，但是对于 fast connections 以及滑动窗口很大时，会增加出现这种场景的概率。关于 TIME_WAIT 详细讨论可以参考 rfc1337 。

下图是一个延迟延迟报文，在原链接已经正常关闭，并使用相同的端口建立了新链接，那么上个链接发送的报文可能混入新的链接中。
<img src="{{site.url}}{{site.baseurl}}/img/tcpip_timewait_duplicate_segment.png"/>
确保远端以关闭
当发送的最后一个 ACK 丢失后，远端处于 LAST_ACK 状态，在此状态时，如果没有收到最后的 ACK 报文，那么就会重发 FIN 。

如下图，当 FIN 丢失，那么被动关闭方会处于 LAST_ACK 状态，那么尝试重新建立链接时，会直接发送一个 RST 关闭链接，影响新链接创建。
<img src="{{site.url}}{{site.baseurl}}/img/tcpip_timewait_last_ack.png"/>
注意，处于 LAST_ACK 状态的链接，如果没有收到最后一个 ACK 报文，那么就会一致重发 FIN 报文，直到如下条件：

由于超时自动关闭链接；
收到了 ACK 报文然后关闭链接；
收到了 RST 报文并关闭链接。
简言之，通过 2MSL 等待时间，保证前一个链接的报文已经在网络上消失，保证双方都已经正常关闭链接。

如果有较多的链接处于 TIME_WAIT 状态 (可以通过 ss -tan state time-wait | wc -l 查看)，那么一般会有如下几个方面的影响：

占用文件描述符 (fd) ，会导致无法创建相同类型的链接；
内核中 socket 结构体占用的内存资源；
额外的 CPU 消耗。

文件描述符
一个链接在 TIME_WAIT 保存一分钟，那么相同四元组 (saddr, sport, daddr, dport) 的链接就无法创建；实际上，如果从内核角度看，实际上根据配置项，还可能包含了 sequence 以及 timestamp 。

如果服务器是挂载在一个 L7 Load-Balancer 之后的，那么源地址是相同的，而 Linux 中的 Port 分配范围是通过 net.ipv4.ip_local_port_range 参数配置的，默认是 3W 左右的可用端口，那么平均是每秒 500 个链接。

在客户端，那么就会报 EADDRNOTAVAIL (Cannot assign requested address) 错误，而服务端的排查相比来说要复杂的多，如下可以查看当前与客户端的链接数。

----- sport就是Local Address列对应的端口
# ss -tan 'sport = :80' | awk '/^TIME-WAIT/ {print $(NF)" "$(NF-1)}' | \
    sed 's/:[^ ]*//g' | sort | uniq -c | sort -rn
针对这一场景，那么就是如何增加四元组，按照难度依次列举如下：

客户端添加更多的端口范围，设置 net.ipv4.ip_local_port_range；
服务端增加监听端口，例如额外增加 81、82、83 …；
注意，在客户端，不同版本的内核其行为也略有区别，老版本内核会查找空闲的本地二元组 (source address, source port)，此时增加服务端的 IP 以及 port 时不会增大链接数；而 Linux 3.2 之后，对于不同的四元组那么可以复用本地的二元组。

最后就是调整 net.ipv4.tcp_tw_reuse 和 net.ipv4.tcp_tw_recycle 参数，不过在一些场景下可能会有问题，所以尽量不要使用，下面再详细介绍。

内存
假设每秒要处理 1W 的链接，那么在 Linux 中处于 TIME_WAIT 状态的链接数就有 60W ，这样的话是否会消耗很多的内存资源？

首先，应用端的 TIME_WAIT 状态链接已经关闭，所以不会消耗资源；主要资源的消耗在内核中。内核中保存了三个不同的结构体：

HASH TABLE
connection hash table 用于快速定位一个现有的链接，例如接收到一个报文快速定位到链接结构体。

$ dmesg | grep "TCP established hash table"
[    0.292951] TCP established hash table entries: 65536 (order: 7, 524288 bytes)
可以在内核启动时通过 thash_entries 参数设置，系统启动时会调用 alloc_large_system_hash() 函数初始化，并打印上述的启动日志信息。

在内核中，处于 TIME_WAIT 状态的链接使用的是 struct tcp_timewait_sock 结构体，而其它状态则使用的是 struct tcp_sock 。

struct inet_timewait_sock {
    struct sock_common  __tw_common;

    int                     tw_timeout;
    volatile unsigned char  tw_substate;
    unsigned char           tw_rcv_wscale;
    __be16 tw_sport;
    unsigned int tw_ipv6only     : 1,
                 tw_transparent  : 1,
                 tw_pad          : 6,
                 tw_tos          : 8,
                 tw_ipv6_offset  : 16;
    unsigned long            tw_ttd;
    struct inet_bind_bucket *tw_tb;
    struct hlist_node        tw_death_node;
};

struct tcp_timewait_sock {
    struct inet_timewait_sock tw_sk;
    u32    tw_rcv_nxt;
    u32    tw_snd_nxt;
    u32    tw_rcv_wnd;
    u32    tw_ts_offset;
    u32    tw_ts_recent;
    long   tw_ts_recent_stamp;
};
TIME_WAIT
用于判断处于 TIME_WAIT 状态的 socket 还有多长时间过期，与上述的 hash 表使用相同的结构体，也就是对应了 struct inet_timewait_sock 结构体中的 struct hlist_node tw_death_node 。

HASH TABLE
bind hash table 用于保存本地已经保存的端口以及相关参数，用于绑定 listen 端口以及查找可用端口。

$ dmesg | grep "TCP bind hash table"
[    0.293146] TCP bind hash table entries: 65536 (order: 8, 1048576 bytes)
这里每个元素使用 struct inet_bind_socket 结构体，大小与 connection 相同。

重点关注两个 hash 表使用的内存即可。

# slabtop -o | grep -E '(^  OBJS|tw_sock_TCP|tcp_bind_bucket)'
  OBJS ACTIVE  USE OBJ SIZE  SLABS OBJ/SLAB CACHE SIZE NAME                   
 50955  49725  97%    0.25K   3397       15     13588K tw_sock_TCP            
 44840  36556  81%    0.06K    760       59      3040K tcp_bind_bucket
内核参数调优
相关参数可以查看内核文档 ip-sysctl.txt ，仅摘取与此相关的参数介绍：

tcp_tw_reuse - BOOLEAN
    Allow to reuse TIME-WAIT sockets for new connections when it is
    safe from protocol viewpoint. Default value is 0.
    It should not be changed without advice/request of technical experts.

tcp_tw_recycle - BOOLEAN
    Enable fast recycling TIME-WAIT sockets. Default value is 0.
    It should not be changed without advice/request of technical experts.
注意，在使用上述参数时，需要先开启 tcp_timestamps 参数，否则这里的配置时无效的，也就是需要开启如下参数。

----- 查看当前状态
# sysctl net.ipv4.tcp_tw_reuse
# sysctl net.ipv4.tcp_tw_recycle
# sysctl net.ipv4.tcp_timestamps

----- 将对应的参数开启
# sysctl net.ipv4.tcp_tw_reuse=1
# sysctl net.ipv4.tcp_tw_recycle=1
# sysctl net.ipv4.tcp_timestamps=1

----- 写入到配置文件，持久化，并使配置立即生效
# cat << EOF >> /etc/sysctl.conf
net.ipv4.tcp_tw_reuse = 1
net.ipv4.tcp_tw_recycle = 1
net.ipv4.tcp_timestamps = 1
EOF
# /sbin/sysctl -p
tcp_tw_reuse
如上所述，TIME_WAIT 状态主要是为了防止延迟报文影响新链接，不过在有些场景下可以确保不会出现这种情况。

RFC 1323 提供了一个 TCP 扩展项，其中定义了两个四字节的时间戳选项，第一个是发送时的发送端时间戳，第二个是从远端接收报文的最近时间戳。

TCP Timestamps Option (TSopt):
    +-------+-------+---------------------+---------------------+
    |Kind=8 |  10   |   TS Value (TSval)  |TS Echo Reply (TSecr)|
    +-------+-------+---------------------+---------------------+
        1       1              4                     4
在开启了 net.ipv4.tcp_tw_reuse 参数后，对于 outgoing 的链接，Linux 会复用处于 TIME_WAIT 状态的资源，当然时间戳要大于前一个链接的最近时间戳，这也就意味着一般一秒之后就可以重新使用该链接。

对于上面说的第一种异常场景，也就是延迟报文被新链接接收，这一场景通过 timestamp 就可以判断上个链接的报文已经过期，从而会直接丢弃。

对于第二种场景，也即防止远端处于 LAST_ACK 状态，此时当复用处于 TIME_WAIT 状态的链接时，第一个发送的 SYN 报文，由于对方校验时间戳不通过，会直接丢弃，并重发一个 FIN 报文，由于此时处于 SYN_SENT 状态，则会发送一个 RST 报文关闭对端的链接。

接下来，就可以正常创建新链接，只是时间会略有增加。

tcpip timewait last ack reuse

tcp_tw_recycle
同样需要开启 timestamp 参数，不过会影响到 incoming 和 outcoming 。

正常来说，处于 TIME_WAIT 状态的超时时间为 TCP_TIMEWAIT_LEN ，也即 60s；当开启了该参数后，该 socket 的释放时间与 RTO (通过RTT计算) 相关。

关于当前链接的 RTO 可以通过 ss --info sport = :2112 dport = :4057 命令查看。

这一过程的处理可以查看源码中的 tcp_time_wait() 函数，内容摘取如下。

if (tcp_death_row.sysctl_tw_recycle && tp->rx_opt.ts_recent_stamp)
    recycle_ok = icsk->icsk_af_ops->remember_stamp(sk);
......

if (timeo < rto)
    timeo = rto;

if (recycle_ok) {
    tw->tw_timeout = rto;
} else {
    tw->tw_timeout = TCP_TIMEWAIT_LEN;
    if (state == TCP_TIME_WAIT)
        timeo = TCP_TIMEWAIT_LEN;
}

inet_twsk_schedule(tw, &tcp_death_row, timeo,
           TCP_TIMEWAIT_LEN);
在内核源码可以发现 tcp_tw_recycle 和 tcp_timestamps 都开启的条件下，如果 60s 内同一源 IP 主机的 socket connect 请求中的 timestamp 必须是递增的。

可以从 TCP 的三次握手时的 SYN 包的处理函数中，也就是 tcp_conn_request() 函数中。

bool tcp_peer_is_proven(struct request_sock *req, struct dst_entry *dst,
            bool paws_check, bool timestamps)
{
    struct tcp_metrics_block *tm;
    bool ret;

    if (!dst)
        return false;

    rcu_read_lock();
    tm = __tcp_get_metrics_req(req, dst);
    if (paws_check) {
        if (tm &&
            // 判断表示该源ip的上次tcp通讯发生在60s内
            (u32)get_seconds() - tm->tcpm_ts_stamp < TCP_PAWS_MSL &&
            // 该条件判断表示该源ip的上次tcp通讯的timestamp 大于本次tcp
            ((s32)(tm->tcpm_ts - req->ts_recent) > TCP_PAWS_WINDOW ||
             !timestamps))
            ret = false;
        else
            ret = true;
    } else {
        if (tm && tcp_metric_get(tm, TCP_METRIC_RTT) && tm->tcpm_ts_stamp)
            ret = true;
        else
            ret = false;
    }
    rcu_read_unlock();

    return ret;
}

int tcp_conn_request(struct request_sock_ops *rsk_ops,
             const struct tcp_request_sock_ops *af_ops,
             struct sock *sk, struct sk_buff *skb)
{
    ... ...
        if (tcp_death_row.sysctl_tw_recycle) { // 本机系统开启tcp_tw_recycle选项
            bool strict;

            dst = af_ops->route_req(sk, &fl, req, &strict);

            if (dst && strict &&
                !tcp_peer_is_proven(req, dst, true,
                        tmp_opt.saw_tstamp)) { // 该socket支持tcp_timestamp
                NET_INC_STATS_BH(sock_net(sk), LINUX_MIB_PAWSPASSIVEREJECTED);
                goto drop_and_release;
            }
        }
    ... ...
}
上述种的状态统计可以查看 /proc/net/netstat 文件中的 PAWSPassive 。

NAT
不过上述的修改，在 NAT 场景下可能会引起更加复杂的问题，在 rfc1323 中，有如下描述。

An additional mechanism could be added to the TCP, a per-host
cache of the last timestamp received from any connection.
This value could then be used in the PAWS mechanism to reject
old duplicate segments from earlier incarnations of the
connection, if the timestamp clock can be guaranteed to have
ticked at least once since the old connection was open.  This
would require that the TIME-WAIT delay plus the RTT together
must be at least one tick of the sender's timestamp clock.
Such an extension is not part of the proposal of this RFC.
大概意思是说 TCP 可以缓存从每个主机收到报文的最新时间戳，后续请求的时间戳如果小于缓存时间戳，则认为该报文无效，数据包会被丢弃。

如上所述，这一行为在 Linux 中在 tcp_timestamps 和 tcp_tw_recycle 都开启时使用，而前者默认是开启的，所以当 tcp_tw_recycle 被开启后，实际上这种行为就被激活了，那么当客户端或服务端以 NAT 方式构建的时候就可能出现问题。

问题排查
针对上述讨论的部分异常，简单讨论下常见的场景。

EADDRNOTAVAIL
在客户端尝试建立链接时，如果报 EADDRNOTAVAIL (Cannot assign requested address) 错误，而且查看可能是已经达到了最大的端口可用数量。

----- 查看内核端口可用范围
$ sysctl net.ipv4.ip_local_port_range
net.ipv4.ip_local_port_range = 32768 60999
----- 上述范围是闭区间，实际可用的端口数量是
$ echo $((60999-32768+1))
28232
针对这一故障场景，可以通过如下步骤修复：

----- 增加本地可用端口数量，这里是临时修改
$ sysctl net.ipv4.ip_local_port_range="10240 61000"
----- 减少TIME_WAIT连接状态
$ sysctl net.ipv4.tcp_tw_reuse=1
NAT
当进行 SNAT 转换时，也即是在服务端看到的都是同一个 IP，那么对于服务端而言这些客户端实际上等同于一个，而客户端的时间戳会存在差异，那么服务端就会出现时间戳错乱的现象，进而直接丢弃时间戳小的数据包。

这类问题的现象是，客户端明发送的 SYN 报文，但服务端没有响应 ACK，可以通过下面命令来确认数据包不断被丢弃的现象：

$ netstat -s | grep timestamp
... packets rejects in established connections because of timestamp
所以，如果无法确定没有使用 NAT ，那么为了安全起见，需要禁止 tcp_tw_recycle 。

tcp_max_tw_buckets
该参数用于设置 TIME_WAIT 状态的 socket 数，超过该限制后就会删除掉，此时系统日志里会显示： TCP: time wait bucket table overflow 。

# sysctl net.ipv4.tcp_max_tw_buckets=100000
如果是 NAT 网络环境又存在大量访问，会产生各种连接不稳定断开的情况。

总结
处于 TIME_WAIT 状态的链接，最大的消耗不是内存或者 CPU，而是四元组对应的资源，所以最有效的方式是增加端口范围、增加地址等。

在 服务端，除非确定没有使用 NAT ，否则不要配置 tcp_tw_recycle 参数；另外，tcp_tw_reuse 参数对于 incoming 的报文无效。

在 客户端，可以直接开启 tcp_tw_reuse 参数，而 tcp_tw_recycle 参数的作用不是很大。

另外，在设计协议时，尽量不要让客户端先关闭链接，最好让服务端去处理这种场景。
https://vincent.bernat.ch/en/blog/2014-tcp-time-wait-state-linux
https://jin-yang.github.io/post/network-tcpip-timewait.html

最近发现一个PHP脚本时常出现连不上服务器的现象，调试了一下，发现是TIME_WAIT状态过多造成的，本文简要介绍一下解决问题的过程。


遇到这类问题，我习惯于先用strace命令跟踪了一下看看：

shell> strace php /path/to/file
EADDRNOTAVAIL (Cannot assign requested address)
从字面结果看似乎是网络资源相关问题。这里顺便介绍一点小技巧：在调试的时候一般是从后往前看strace命令的结果，这样更容易找到有价值的信息。

查看一下当前的网络连接情况，结果发现TIME_WAIT数非常大：

shell> netstat -ant | awk '
    {++s[$NF]} END {for(k in s) print k,s[k]}
'
补充一下，同netstat相比，ss要快很多：

shell> ss -ant | awk '
    {++s[$1]} END {for(k in s) print k,s[k]}
'
重复了几次测试，结果每次出问题的时候，TIME_WAIT都等于28233，这真是一个魔法数字！实际原因很简单，它取决于一个内核参数net.ipv4.ip_local_port_range：

shell> sysctl -a | grep port
net.ipv4.ip_local_port_range = 32768 61000
因为端口范围是一个闭区间，所以实际可用的端口数量是：

shell> echo $((61000-32768+1))
28233
问题分析到这里基本就清晰了，解决方向也明确了，内容所限，这里就不说如何优化程序代码了，只是从系统方面来阐述如何解决问题，无非就是以下两个方面：

首先是增加本地可用端口数量。这点可以用以下命令来实现：

shell> sysctl net.ipv4.ip_local_port_range="10240 61000"
其次是减少TIME_WAIT连接状态。网络上已经有不少相关的介绍，大多是建议：

shell> sysctl net.ipv4.tcp_tw_reuse=1
shell> sysctl net.ipv4.tcp_tw_recycle=1
注：通过sysctl命令修改内核参数，重启后会还原，要想持久化可以参考前面的方法。

这两个选项在降低TIME_WAIT数量方面可以说是立竿见影，不过如果你觉得问题已经完美搞定那就错了，实际上这样可能会引入一个更复杂的网络故障。

关于内核参数的详细介绍，可以参考官方文档。我们这里简要说明一下tcp_tw_recycle参数。它用来快速回收TIME_WAIT连接，不过如果在NAT环境下会引发问题。

RFC1323中有如下一段描述：

An additional mechanism could be added to the TCP, a per-host cache of the last timestamp received from any connection. This value could then be used in the PAWS mechanism to reject old duplicate segments from earlier incarnations of the connection, if the timestamp clock can be guaranteed to have ticked at least once since the old connection was open. This would require that the TIME-WAIT delay plus the RTT together must be at least one tick of the sender’s timestamp clock. Such an extension is not part of the proposal of this RFC.

大概意思是说TCP有一种行为，可以缓存每个主机最新的时间戳，后续请求中如果时间戳小于缓存的时间戳，即视为无效，相应的数据包会被丢弃。

Linux是否启用这种行为取决于tcp_timestamps和tcp_tw_recycle，因为tcp_timestamps缺省就是开启的，所以当tcp_tw_recycle被开启后，实际上这种行为就被激活了，当客户端或服务端以NAT方式构建的时候就可能出现问题，下面以客户端NAT为例来说明：

当多个客户端通过NAT方式联网并与服务端交互时，服务端看到的是同一个IP，也就是说对服务端而言这些客户端实际上等同于一个，可惜由于这些客户端的时间戳可能存在差异，于是乎从服务端的视角看，便可能出现时间戳错乱的现象，进而直接导致时间戳小的数据包被丢弃。如果发生了此类问题，具体的表现通常是是客户端明明发送的SYN，但服务端就是不响应ACK，我们可以通过下面命令来确认数据包不断被丢弃的现象：

shell> netstat -s | grep timestamp
... packets rejects in established connections because of timestamp
安全起见，通常要禁止tcp_tw_recycle。说到这里，大家可能会想到另一种解决方案：把tcp_timestamps设置为0，tcp_tw_recycle设置为1，这样不就可以鱼与熊掌兼得了么？可惜一旦关闭了tcp_timestamps，那么即便打开了tcp_tw_recycle，也没有效果。

好在我们还有另一个内核参数tcp_max_tw_buckets（一般缺省是180000）可用：

shell> sysctl net.ipv4.tcp_max_tw_buckets=100000
通过设置它，系统会将多余的TIME_WAIT删除掉，此时系统日志里可能会显示：「TCP: time wait bucket table overflow」，不过除非不得已，否则不要轻易使用。
