---
title: tcp_tw_reuse
layout: post
category: linux
author: 夏泽民
---
linux TIME_WAIT 相关参数:

net.ipv4.tcp_tw_reuse = 0    表示开启重用。允许将TIME-WAIT sockets重新用于新的TCP连接，默认为0，表示关闭
net.ipv4.tcp_tw_recycle = 0  表示开启TCP连接中TIME-WAIT sockets的快速回收，默认为0，表示关闭
net.ipv4.tcp_fin_timeout = 60  表示如果套接字由本端要求关闭，这个参数决定了它保持在FIN-WAIT-2状态的时间
注意：

- 不像Windows 可以修改注册表修改2MSL 的值，linux 需要修改内核宏定义重新编译，tcp_fin_timeout 不是2MSL 而是Fin-WAIT-2状态超时时间.

- tcp_tw_reuse 和 SO_REUSEADDR 是两个完全不同的东西

   SO_REUSEADDR 允许同时绑定 127.0.0.1 和 0.0.0.0 同一个端口； SO_RESUSEPORT linux 3.7才支持，用于绑定相同ip:port，像nginx 那样 fork方式也能实现

 

1. tw_reuse，tw_recycle 必须在客户端和服务端 timestamps 开启时才管用（默认打开）

2. tw_reuse 只对客户端起作用，开启后客户端在1s内回收

3. tw_recycle 对客户端和服务器同时起作用，开启后在 3.5*RTO 内回收，RTO 200ms~ 120s 具体时间视网络状况。

　　内网状况比tw_reuse 稍快，公网尤其移动网络大多要比tw_reuse 慢，优点就是能够回收服务端的TIME_WAIT数量

 

对于客户端
1. 作为客户端因为有端口65535问题，TIME_OUT过多直接影响处理能力，打开tw_reuse 即可解决，不建议同时打开tw_recycle，帮助不大；

2. tw_reuse 帮助客户端1s完成连接回收，基本可实现单机6w/s短连接请求，需要再高就增加IP数量；

3. 如果内网压测场景，且客户端不需要接收连接，同时 tw_recycle 会有一点点好处；

4. 业务上也可以设计由服务端主动关闭连接。

 

对于服务端
1. 打开tw_reuse无效

2. 线上环境 tw_recycle 不建议打开

   服务器处于NAT 负载后，或者客户端处于NAT后（基本公司家庭网络基本都走NAT）；

　公网服务打开就可能造成部分连接失败，内网的话到时可以视情况打开；

   像我所在公司对外服务都放在负载后面，负载会把 timestamp 都给清空，就算你打开也不起作用。

3. 服务器TIME_WAIT 高怎么办

   不像客户端有端口限制，处理大量TIME_WAIT Linux已经优化很好了，每个处于TIME_WAIT 状态下连接内存消耗很少，

而且也能通过tcp_max_tw_buckets = 262144 配置最大上限，现代机器一般也不缺这点内存。
<!-- more -->
slabtop
251461  95%    0.25K  17482       15     69928K tw_sock_TCP
 ss -s
Total: 259 (kernel 494)
TCP:   262419 (estab 113, closed 262143, orphaned 156, synrecv 0, timewait 262143/0), ports 80

Transport Total     IP        IPv6
*         494       -         -        
RAW       1         1         0        
UDP       0         0         0        
TCP       276       276       0        
INET      277       277       0        
FRAG      0         0         0

原理分析
 1. MSL 由来

　　发起连接关闭方回复最后一个fin 的ack，为避免对方ack 收不到、重发的或还在中间路由上的fin 把新连接给丢掉了，等个2MSL（linux 默认2min）。

　　也就是连接有谁关闭的那一方有time_wait问题，被关那方无此问题。

2. reuse、recycle

     通过timestamp的递增性来区分是否新连接，新连接的timestamp更大，那么保证小的timestamp的 fin 不会fin掉新连接，不用等2MSL。

3. reuse

     通过timestamp 递增性，客户端、服务器能够处理outofbind fin包

4. recycle

    对于服务端，同一个src ip，可能会是NAT后很多机器，这些机器timestamp递增性无可保证，服务器会拒绝非递增请求连接。
    
在执行sysctl -p操作时突然报错如下：sysctl: cannot stat /proc/sys/net/ipv4/tcp_tw_recycle: No such file or directory复制代码2、问题原因Linux 从4.12内核版本开始移除了 tcp_tw_recycle 配置。参考：[1]tcp:remove tcp_tw_recycle 4396e460移除sysctl.conf中关于net.ipv4.tcp_tw_recycle的配置内容，再次尝试sysctl -p就不再提示报错了。3、深入解析tcp_tw_recycle通常会和tcp_tw_reuse参数一起使用，用于解决服务器TIME_WAIT状态连接过多的问题。

TIME_WAIT永远是出现在主动发送断开连接请求的一方(下文中我们称之为客户)，划重点：这一点面试的时候经常会被问到。客户在收到服务器端发送的FIN(表示"我们也要断开连接了")后发送ACK报文，并且进入TIME_WAIT状态，等待2MSL(MaximumSegmentLifetime 最大报文生存时间)。对于Linux，字段为TCP_TIMEWAIT_LEN硬编码为30秒，对于windows为2分钟(可自行调整)。为什么客户端不直接进入CLOSED状态，而是要在TIME_WAIT等待那么久呢，基于如下考虑：1.确保远程端处于关闭状态。也就是说需要确保客户端发出的最后一个ACK报文能够到达服务器。由于网络不可靠，有可能最后一个ACK报文丢失，如果服务器没有收到客户端的ACK，则会重新发送FIN报文，客户端就可以在2MSL时间段内收到这个这个重发的报文，并且重发ACK报文。但如果客户端跳过TIME_WAIT阶段进入了CLOSED，服务端始终无法得到响应，就会处于LAST-ACK状态，此时假如客户端发起了一个新连接，则会以失败告终。

2.防止上一次连接中的包，迷路后重新出现，影响新连接(经过2MSL,上一次连接中所有的重复包都会消失)，这一点和为啥要执行三次握手而不是两次的原因是一样的。

查看方式有两种：（1）ss -tan state time-wait|wc -l（2）netstat -n | awk '/^tcp/ {++S[$NF]} END {for(a in S) print a, S[a]}'3.2、TIME_WAIT的危害对于一个处理大量连接的处理器TIME_WAIT是有危害的，表现如下：1.占用连接资源TIME_WAIT占用的1分钟时间内，相同四元组(源地址，源端口，目标地址，目标端口)的连接无法创建，通常一个ip可以开启的端口为net.ipv4.ip_local_port_range指定的32768-61000，如果TIME_WAIT状态过多，会导致无法创建新连接。2.占用内存资源这个占用资源并不是很多，可以不用担心。3.3、TIME_WAIT的解决可以考虑如下方式：1.修改为长连接，代价较大，长连接对服务器性能有影响。2.增加可用端口范围(修改net.ipv4.ip_local_port_range); 增加服务端口，比如采用80，81等多个端口提供服务; 增加客户端ip(适用于负载均衡，比如nginx，采用多个ip连接后端服务器); 增加服务端ip; 这些方式治标不治本，只能缓解问题。3.将net.ipv4.tcp_max_tw_buckets设置为很小的值(默认是18000). 当TIME_WAIT连接数量达到给定的值时，所有的TIME_WAIT连接会被立刻清除，并打印警告信息。但这种粗暴的清理掉所有的连接，意味着有些连接并没有成功等待2MSL，就会造成通讯异常。4.修改TCP_TIMEWAIT_LEN值，减少等待时间，但这个需要修改内核并重新编译。5.打开tcp_tw_recycle和tcp_timestamps选项。6.打开tcp_tw_reuse和tcp_timestamps选项。3.4、net.ipv4.tcp_tw_{reuse,recycle}需要明确两个点：解决方式已经给出，那我们需要了解一下net.ipv4.tcp_tw_reuse和net.ipv4.tcp_tw_recycle有啥区别1.两个选项都需要打开对TCP时间戳的支持，即net.ipv4.tcp_timestamps=1(默认即为1)。RFC 1323中实现了TCP拓展规范，以便保证网络繁忙的情况下的高可用。并定义了一个新的TCP选项-两个四字节的timestamp字段，第一个是TCP发送方的当前时钟时间戳，第二个是从远程主机接收到的最新时间戳。2.两个选项默认都是关闭状态，即等于0。3.4.1 - net.ipv4.tcp_tw_reuse：更安全的设置将处于TIME_WAIT状态的socket用于新的TCP连接，影响连出的连接。[2]kernel sysctl 官方指南中是这么写的：Allow to reuse TIME-WAIT sockets for new connections when it is safe from protocol viewpoint. Default value is 0.It should not be changed without advice/request of technical experts.协议安全主要指的是两点：1.只适用于客户端(连接发起方)net/ipv4/inet_hashtables.cstatic int __inet_check_established(struct inet_timewait_death_row *death_row,
                    struct sock *sk, __u16 lport,
                    struct inet_timewait_sock **twp)
{
    /* ……省略…… */
    sk_nulls_for_each(sk2, node, &head->chain) {
            if (sk2->sk_hash != hash)
                        continue;
                        
            if (likely(INET_MATCH(sk2, net, acookie,
                    saddr, daddr, ports, dif))) {
                        if (sk2->sk_state == TCP_TIME_WAIT) {
                            tw = inet_twsk(sk2);
                            if (twsk_unique(sk, sk2, twp))
                                break;
            }
            goto not_unique;
        }
    }
    /* ……省略…… */
}复制代码2.TIME_WAIT创建时间超过1秒才可以被复用net/ipv4/tcp_ipv4.cint tcp_twsk_unique(struct sock *sk, struct sock *sktw, void *twp)
{
    /* ……省略…… */
    if (tcptw->tw_ts_recent_stamp &&
        (!twp || (sock_net(sk)->ipv4.sysctl_tcp_tw_reuse &&
         get_seconds() - tcptw->tw_ts_recent_stamp > 1))) {
         /* ……省略…… */
         return 1;
    }
    return 0;
}复制代码满足以上两个条件才会被认为是"safe from protocol viewpoint"的状况。启用net.ipv4.tcp_tw_reuse后，如果新的时间戳比之前存储的时间戳更大，那么Linux将会从TIME-WAIT状态的存活连接中选取一个，重新分配给新的连接出去的的TCP连接，这种情况下，TIME-WAIT的连接相当于只需要1秒就可以被复用了。重新回顾为什么要引入TIME-WAIT：第一个作用就是避免新连接接收到重复的数据包，由于使用了时间戳，重复的数据包会因为时间戳过期被丢弃。第二个作用是确保远端不是处于LAST-ACK状态，如果ACK包丢失，远端没有成功获取到最后一个ACK包，则会重发FIN包。直到：1.放弃(连接断开)2.收到ACK包3.收到RST包如果FIN包被及时接收到，并且本地端仍然是TIME-WAIT状态，那ACK包会被发送，此时就是正常的四次挥手流程。如果TIME-WAIT的条目已经被新连接所复用，则新连接的SYN包会被忽略掉，并且会收到FIN包的重传，本地会回复一个RST包(因为此时本地连接为SYN-SENT状态)，这会让远程端跳出LAST-ACK状态，最初的SYN包也会在1秒后重新发送，然后完成连接的建立，整个过程不会中断，只是有轻微的延迟。

3.4.2 - net.ipv4.tcp_tw_recycle：更激进的设置启用TIME_WAIT 状态的sockets的快速回收，影响所有连入和连出的连接[3]kernel sysctl 官方指南是这么写的：Enable fast recycling TIME-WAIT sockets. Default value is 0. It should not be changed without advice/request of technical experts.这次表述的更加模糊，继续翻看源码：net/ipv4/tcp_input.cint tcp_conn_request(struct request_sock_ops *rsk_ops,
            const struct tcp_request_sock_ops *af_ops,
            struct sock *sk, struct sk_buff *skb)
{
 /* ……省略…… */
 if (!want_cookie && !isn) {
     /* ……省略…… */
     if (net->ipv4.tcp_death_row.sysctl_tw_recycle) {
         bool strict;

dst = af_ops->route_req(sk, &fl, req, &strict);
 
if (dst && strict &&
              !tcp_peer_is_proven(req, dst, true,
                      tmp_opt.saw_tstamp)) {
              NET_INC_STATS(sock_net(sk), LINUX_MIB_PAWSPASSIVEREJECTED);
              goto drop_and_release;
       }
     }
     /* ……省略…… */
     isn = af_ops->init_seq(skb, &tcp_rsk(req)->ts_off);
   }
/* ……省略…… */

drop_and_release:
            dst_release(dst);
       drop_and_free:
            reqsk_free(req);
       drop:
            tcp_listendrop(sk);
            return 0;
}复制代码简单来说就是，Linux会丢弃所有来自远端的timestramp时间戳小于上次记录的时间戳(由同一个远端发送的)的任何数据包。也就是说要使用该选项，则必须保证数据包的时间戳是单调递增的。问题在于，此处的时间戳并不是我们通常意义上面的绝对时间，而是一个相对时间。很多情况下，我们是没法保证时间戳单调递增的，比如使用了nat，lvs等情况。而这也是很多优化文章中并没有提及的一点，大部分文章都是简单的推荐将net.ipv4.tcp_tw_recycle设置为1，却忽略了该选项的局限性，最终造成严重的后果(比如我们之前就遇到过部署在nat后端的业务网站有的用户访问没有问题，但有的用户就是打不开网页)。3.5、被抛弃的tcp_tw_recycle如果说之前内核中tcp_tw_recycle仅仅不适用于nat和lvs环境，那么从4.10内核开始，官方修改了时间戳的生成机制。参考：[4]tcp: randomize tcp timestamp offsets for each connection 95a22ca在这种情况下，无论任何时候，tcp_tw_recycle都不应该开启。故被抛弃也是理所应当的了。4、总结tcp_tw_recycle 选项在4.10内核之前还只是不适用于NAT/LB的情况(其他情况下，我们也非常不推荐开启该选项)，但4.10内核后彻底没有了用武之地，并且在4.12内核中被移除.tcp_tw_reuse 选项仍然可用。在服务器上面，启用该选项对于连入的TCP连接来说不起作用，但是对于客户端(比如服务器上面某个服务以客户端形式运行，比如nginx反向代理)等是一个可以考虑的方案。修改TCP_TIMEWAIT_LEN是非常不建议的行为。

https://www.cnxct.com/coping-with-the-tcp-time_wait-state-on-busy-linux-servers-in-chinese-and-dont-enable-tcp_tw_recycle/