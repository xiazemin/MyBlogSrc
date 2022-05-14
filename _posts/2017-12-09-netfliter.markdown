---
title: netfliter
layout: post
category: linux
author: 夏泽民
---
<!-- more -->
Netfilter是Linux 2.4.x引入的一个子系统，它作为一个通用的、抽象的框架，提供一整套的hook函数的管理机制，使得诸如数据包过滤、网络地址转换(NAT)和基于协议类型的连接跟踪成为了可能。
netfilter的架构就是在整个网络流程的若干位置放置了一些检测点HOOK），而在每个检测点上登记了一些处理函数进行处理。

netfilter是由Rusty Russell提出的Linux 2.4内核防火墙框架，该框架既简洁又灵活，可实现安全策略应用中的许多功能，如数据包过滤、数据包处理、地址伪装、透明代理、动态网络地址转换（Network Address Translation，NAT），以及基于用户及媒体访问控制（Media Access Control，MAC）地址的过滤和基于状态的过滤、包速率限制等。
框架
netfilter提供了一个抽象、通用化的框架[1]，作为中间件，为每种网络协议（IPv4、IPv6等）定义一套钩子函数。Ipv4定义了5个钩子函数，这些钩子函数在数据报流过协议栈的5个关键点被调用，也就是说，IPv4协议栈上定义了5个“允许垂钓点”。在每一个“垂钓点”，都可以让netfilter放置一个“鱼钩”，把经过的网络包（Packet）钓上来，与相应的规则链进行比较，并根据审查的结果，决定包的下一步命运，即是被原封不动地放回IPv4协议栈，继续向上层递交；还是经过一些修改，再放回网络；或者干脆丢弃掉。
Ipv4中的一个数据包通过netfilter系统的过程如图1所示。
图1 Netfilter的功能框架
关键技术
netfilter主要采用连线跟踪（Connection Tracking）、包过滤（Packet Filtering）、地址转换、包处理（Packet Mangling)4种关键技术。
⒈2.1 连线跟踪
连线跟踪是包过滤、地址转换的基础，它作为一个独立的模块运行。采用连线跟踪技术在协议栈低层截取数据包，将当前数据包及其状态信息与历史数据包及其状态信息进行比较，从而得到当前数据包的控制信息，根据这些信息决定对网络数据包的操作，达到保护网络的目的。
当下层网络接收到初始化连接同步（Synchronize，SYN）包，将被netfilter规则库检查。该数据包将在规则链中依次序进行比较。如果该包应被丢弃，发送一个复位（Reset，RST）包到远端主机，否则连接接收。这次连接的信息将被保存在连线跟踪信息表中，并表明该数据包所应有的状态。这个连线跟踪信息表位于内核模式下，其后的网络包就将与此连线跟踪信息表中的内容进行比较，根据信息表中的信息来决定该数据包的操作。因为数据包首先是与连线跟踪信息表进行比较，只有SYN包才与规则库进行比较，数据包与连线跟踪信息表的比较都是在内核模式下进行的，所以速度很快。
⒈2.2 包过滤
包过滤检查通过的每个数据包的头部，然后决定如何处置它们，可以选择丢弃，让包通过，或者更复杂的操作。
⒈2.3 地址转换
网络地址转换 分为源NAT（Source NAT，SNAT）和目的NAT(Destination NAT,DNAT)2种不同的类型。SNAT是指修改数据包的源地址（改变连接的源IP）。SNAT会在数据包送出之前的最后一刻做好转换工作。地址伪装（Masquerading）是SNAT的一种特殊形式。DNAT 是指修改数据包的目标地址（改变连接的目的IP）。DNAT 总是在数据包进入以后立即完成转换。端口转发、负载均衡和透明代理都属于DNAT。
⒈2.4 包处理
利用包处理可以设置或改变数据包的服务类型（Type of Service,TOS）字段；改变包的生存期（Time to Live,TTL）字段；在包中设置标志值，利用该标志值可以进行带宽限制和分类查询。



<div class="container">
	<div class="row">
	<img src="{{site.url}}{{site.baseurl}}/img/netfliter.png"/>
	</div>
	<div class="row">
	</div>
</div>
Netfilter 是Linux内核中进行数据包过滤、连接跟踪、地址转换等的主要实现框架。当我们希望过滤特定的数据包或者需要修改数据包的内容再发送出去，这些动作主要都在netfilter中完成。
iptables工具就是用户空间和内核的Netfilter模块通信的手段，iptables命令提供很多选项来实现过滤数据包的各种操作，所以，我们在定义数据包过滤规则时，并不需要去直接修改内核中的netfilter模块，后面会讲到iptables命令如何作用于内核中的netfilter。
Netfilter的实质就是定义一系列的hook点（挂钩），每个hook点上可以挂载多个hook函数，hook函数中就实现了我们要对数据包的内容做怎样的修改、以及要将数据包放行还是过滤掉。数据包进入netfilter框架后，实际上就是依次经过所有hook函数的处理，数据包的命运就掌握在这些hook函数的手里。
本文基于内核版本2.6.31。
所有的hook点都放在一个全局的二维数组，每个hook点上的hook函数按照优先级顺序注册到一个链表中，注册的接口为nf_register_hook()。这个二维数组的定义如下：
struct list_head nf_hooks[NFPROTO_NUMPROTO][NF_MAX_HOOKS]__read_mostly;
其中NFPROTO_NUMPROTO 为netfilter支持的协议类型：
enum {
   NFPROTO_UNSPEC=  0,
   NFPROTO_IPV4   =  2,
   NFPROTO_ARP    =  3,
   NFPROTO_BRIDGE=  7,
   NFPROTO_IPV6   = 10,
   NFPROTO_DECNET= 12,
   NFPROTO_NUMPROTO,
};

