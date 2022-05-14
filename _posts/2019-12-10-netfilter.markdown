---
title: netfilter iptables
layout: post
category: linux
author: 夏泽民
---
https://www.netfilter.org/
https://arthurchiao.github.io/blog/deep-dive-into-iptables-and-netfilter-arch-zh/
https://www.digitalocean.com/community/tutorials/a-deep-dive-into-iptables-and-netfilter-architecture
<!-- more -->
	<img src="{{site.url}}{{site.baseurl}}/img/Netfilter-packet-flow.svg"/>
	
	https://arthurchiao.github.io/blog/nat-zh/
	
	防火墙是保护服务器和基础设施安全的重要工具。在 Linux 生态系统中，iptables 是使 用很广泛的防火墙工具之一，它基于内核的包过滤框架（packet filtering framework） netfilter。如果管理员或用户不了解这些系统的架构，那可能就无法创建出可靠的防火 墙策略，一方面是因为 iptables 的语法颇有挑战性，另外一方面是 netfilter 框架内部 相互交织而变得错综复杂。

本文将带领读者深入理解 iptables 框架，让那些需要创建防火墙策略的用户对它有一个 更全面的认识。我们会讨论 iptables 是如何与 netfilter 交互的，几个组件是如何组织 成一个全面的过滤和矫正系统（a comprehensive filtering and mangling system）的。

1 IPTables 和 Netfilter 是什么？
Linux 上最常用的防火墙工具是 iptables。iptables 与协议栈内有包过滤功能的 hook 交 互来完成工作。这些内核 hook 构成了 netfilter 框架。

每个进入网络系统的包（接收或发送）在经过协议栈时都会触发这些 hook，程序 可以通过注册 hook 函数的方式在一些关键路径上处理网络流量。iptables 相关的内核模 块在这些 hook 点注册了处理函数，因此可以通过配置 iptables 规则来使得网络流量符合 防火墙规则。

2. Netfilter Hooks
netfilter 提供了 5 个 hook 点。包经过协议栈时会触发内核模块注册在这里的处理函数 。触发哪个 hook 取决于包的方向（是发送还是接收）、包的目的地址、以及包在上一个 hook 点是被丢弃还是拒绝等等。

下面几个 hook 是内核协议栈中已经定义好的：

NF_IP_PRE_ROUTING: 接收到的包进入协议栈后立即触发此 hook，在进行任何路由判断 （将包发往哪里）之前
NF_IP_LOCAL_IN: 接收到的包经过路由判断，如果目的是本机，将触发此 hook
NF_IP_FORWARD: 接收到的包经过路由判断，如果目的是其他机器，将触发此 hook
NF_IP_LOCAL_OUT: 本机产生的准备发送的包，在进入协议栈后立即触发此 hook
NF_IP_POST_ROUTING: 本机产生的准备发送的包或者转发的包，在经过路由判断之后， 将触发此 hook
注册处理函数时必须提供优先级，以便 hook 触发时能按照 优先级高低调用处理函数。这使得多个模块（或者同一内核模块的多个实例）可以在同一 hook 点注册，并且有确定的处理顺序。内核模块会依次被调用，每次返回一个结果给 netfilter 框架，提示该对这个包做什么操作。

3 IPTables 表和链（Tables and Chains）
iptables 使用 table 来组织规则，根据用来做什么类型的判断（the type of decisions they are used to make）标准，将规则分为不同 table。例如，如果规则是处理 网络地址转换的，那会放到 nat table；如果是判断是否允许包继续向前，那可能会放到 filter table。

在每个 table 内部，规则被进一步组织成 chain，内置的 chain 是由内置的 hook 触发 的。chain 基本上能决定（basically determin）规则何时被匹配。

下面可以看出，内置的 chain 名字和 netfilter hook 名字是一一对应的：

PREROUTING: 由 NF_IP_PRE_ROUTING hook 触发
INPUT: 由 NF_IP_LOCAL_IN hook 触发
FORWARD: 由 NF_IP_FORWARD hook 触发
OUTPUT: 由 NF_IP_LOCAL_OUT hook 触发
POSTROUTING: 由 NF_IP_POST_ROUTING hook 触发
chain 使管理员可以控制在包的传输路径上哪个点（where in a packet’s delivery path）应用策略。因为每个 table 有多个 chain，因此一个 table 可以在处理过程中的多 个地方施加影响。特定类型的规则只在协议栈的特定点有意义，因此并不是每个 table 都 会在内核的每个 hook 注册 chain。

内核一共只有 5 个 netfilter hook，因此不同 table 的 chain 最终都是注册到这几个点 。例如，有三个 table 有 PRETOUTING chain。当这些 chain 注册到对应的 NF_IP_PRE_ROUTING hook 点时，它们需要指定优先级，应该依次调用哪个 table 的 PRETOUTING chain，优先级从高到低。我们一会就会看到 chain 的优先级问题。

4. table 种类
先来看看 iptables 提供的 table 类型。这些 table 是按规则类型区分的。

4.1 Filter Table
filter table 是最常用的 table 之一，用于判断是否允许一个包通过。

在防火墙领域，这通常称作“过滤”包（”filtering” packets）。这个 table 提供了防火墙 的一些常见功能。

4.2 NAT Table
nat table 用于实现网络地址转换规则。

当包进入协议栈的时候，这些规则决定是否以及如何修改包的源/目的地址，以改变包被 路由时的行为。nat table 通常用于将包路由到无法直接访问的网络。

4.3 Mangle Table
mangle （修正）table 用于修改包的 IP 头。

例如，可以修改包的 TTL，增加或减少包可以经过的跳数。

这个 table 还可以对包打只在内核内有效的“标记”（internal kernel “mark”），后 续的 table 或工具处理的时候可以用到这些标记。标记不会修改包本身，只是在包的内核 表示上做标记。

4.4 Raw Table
iptables 防火墙是有状态的：对每个包进行判断的时候是依赖已经判断过的包。

建立在 netfilter 之上的连接跟踪（connection tracking）特性使得 iptables 将包 看作已有的连接或会话的一部分，而不是一个由独立、不相关的包组成的流。连接跟踪逻 辑在包到达网络接口之后很快就应用了。

raw table 定义的功能非常有限，其唯一目的就是提供一个让包绕过连接跟踪的框架。

4.5 Security Table
security table 的作用是给包打上 SELinux 标记，以此影响 SELinux 或其他可以解读 SELinux 安全上下文的系统处理包的行为。这些标记可以基于单个包，也可以基于连接。

5 每种 table 实现的 chain
前面已经分别讨论了 table 和 chain，接下来看每个 table 里各有哪些 chain。另外，我 们还将讨论注册到同一 hook 的不同 chain 的优先级问题。例如，如果三个 table 都有 PRETOUTING chain，那应该按照什么顺序调用它们呢？

下面的表格展示了 table 和 chain 的关系。横向是 table， 纵向是 chain，Y 表示 这个 table 里面有这个 chain。例如，第二行表示 raw table 有 PRETOUTING 和 OUTPUT 两 个 chain。具体到每列，从上倒下的顺序就是 netfilter hook 触发的时候，（对应 table 的）chain 被调用的顺序。

有几点需要说明一下。在下面的图中，nat table 被细分成了 DNAT （修改目的地址） 和 SNAT（修改源地址），以更方便地展示他们的优先级。另外，我们添加了路由决策点 和连接跟踪点，以使得整个过程更完整全面：

Tables/Chains	PREROUTING	INPUT	FORWARD	OUTPUT	POSTROUTING
(路由判断)	 	 	 	Y	 
raw	Y	 	 	Y	 
(连接跟踪）	Y	 	 	Y	 
mangle	Y	Y	Y	Y	Y
nat (DNAT)	Y	 	 	Y	 
(路由判断)	Y	 	 	Y	 
filter	 	Y	Y	Y	 
security	 	Y	Y	Y	 
nat (SNAT)	 	Y	 	Y	Y
当一个包触发 netfilter hook 时，处理过程将沿着列从上向下执行。 触发哪个 hook （列）和包的方向（ingress/egress）、路由判断、过滤条件等相关。

特定事件会导致 table 的 chain 被跳过。例如，只有每个连接的第一个包会去匹配 NAT 规则，对这个包的动作会应用于此连接后面的所有包。到这个连接的应答包会被自动应用反 方向的 NAT 规则。

Chain 遍历优先级
假设服务器知道如何路由数据包，而且防火墙允许数据包传输，下面就是不同场景下包的游 走流程：

收到的、目的是本机的包：PRETOUTING -> INPUT
收到的、目的是其他主机的包：PRETOUTING -> FORWARD -> POSTROUTING
本地产生的包：OUTPUT -> POSTROUTING
综合前面讨论的 table 顺序问题，我们可以看到对于一个收到的、目的是本机的包： 首先依次经过 PRETOUTING chain 上面的 raw、mangle、nat table；然后依次经 过 INPUT chain 的 mangle、filter、security、nat table，然后才会到达本机 的某个 socket。

6 IPTables 规则
规则放置在特定 table 的特定 chain 里面。当 chain 被调用的时候，包会依次匹配 chain 里面的规则。每条规则都有一个匹配部分和一个动作部分。

6.1 匹配
规则的匹配部分指定了一些条件，包必须满足这些条件才会和相应的将要执行的动作（“ target”）进行关联。

匹配系统非常灵活，还可以通过 iptables extension 大大扩展其功能。规则可以匹配协 议类型、目的或源地址、目的或源端口、目的或源网段、接收或发送的接口（网卡）、协议 头、连接状态等等条件。这些综合起来，能够组合成非常复杂的规则来区分不同的网络流 量。

6.2 目标
包符合某种规则的条件而触发的动作（action）叫做目标（target）。目标分为两种类型：

终止目标（terminating targets）：这种 target 会终止 chain 的匹配，将控制权 转移回 netfilter hook。根据返回值的不同，hook 或者将包丢弃，或者允许包进行下一 阶段的处理
非终止目标（non-terminating targets）：非终止目标执行动作，然后继续 chain 的执行。虽然每个 chain 最终都会回到一个终止目标，但是在这之前，可以执行任意多 个非终止目标
每个规则可以跳转到哪个 target 依上下文而定，例如，table 和 chain 可能会设置 target 可用或不可用。规则里激活的 extensions 和匹配条件也影响 target 的可用性。

7 跳转到用户自定义 chain
这里要介绍一种特殊的非终止目标：跳转目标（jump target）。jump target 是跳转到其 他 chain 继续处理的动作。我们已经讨论了很多内置的 chain，它们和调用它们的 netfilter hook 紧密联系在一起。然而，iptables 也支持管理员创建他们自己的用于管理 目的的 chain。

向用户自定义 chain 添加规则和向内置的 chain 添加规则的方式是相同的。不同的地方 在于，用户定义的 chain 只能通过从另一个规则跳转（jump）到它，因为它们没有注册到 netfilter hook。

用户定义的 chain 可以看作是对调用它的 chain 的扩展。例如，用户定义的 chain 在结 束的时候，可以返回 netfilter hook，也可以继续跳转到其他自定义 chain。

这种设计使框架具有强大的分支功能，使得管理员可以组织更大更复杂的网络规则。

8 IPTables 和连接跟踪
在讨论 raw table 和 匹配连接状态的时候，我们介绍了构建在 netfilter 之上的连 接跟踪系统。连接跟踪系统使得 iptables 基于连接上下文而不是单个包来做出规则判 断，给 iptables 提供了有状态操作的功能。

连接跟踪在包进入协议栈之后很快（very soon）就开始工作了。在给包分配连接之前所做 的工作非常少，只有检查 raw table 和一些基本的完整性检查。

跟踪系统将包和已有的连接进行比较，如果包所属的连接已经存在就更新连接状态，否则就 创建一个新连接。如果 raw table 的某个 chain 对包标记为目标是 NOTRACK，那这 个包会跳过连接跟踪系统。

连接的状态
连接跟踪系统中的连接状态有：

NEW：如果到达的包关连不到任何已有的连接，但包是合法的，就为这个包创建一个新连接。对 面向连接的（connection-aware）的协议例如 TCP 以及非面向连接的（connectionless ）的协议例如 UDP 都适用

ESTABLISHED：当一个连接收到应答方向的合法包时，状态从 NEW 变成 ESTABLISHED。对 TCP 这个合法包其实就是 SYN/ACK 包；对 UDP 和 ICMP 是源和目 的 IP 与原包相反的包

RELATED：包不属于已有的连接，但是和已有的连接有一定关系。这可能是辅助连接（ helper connection），例如 FTP 数据传输连接，或者是其他协议试图建立连接时的 ICMP 应答包

INVALID：包不属于已有连接，并且因为某些原因不能用来创建一个新连接，例如无法 识别、无法路由等等

UNTRACKED：如果在 raw table 中标记为目标是 UNTRACKED，这个包将不会进入连 接跟踪系统

SNAT：包的源地址被 NAT 修改之后会进入的虚拟状态。连接跟踪系统据此在收到 反向包时对地址做反向转换

DNAT：包的目的地址被 NAT 修改之后会进入的虚拟状态。连接跟踪系统据此在收到 反向包时对地址做反向转换

这些状态可以定位到连接生命周期内部，管理员可以编写出更加细粒度、适用范围更大、更 安全的规则。

9 总结
netfilter 包过滤框架和 iptables 防火墙是 Linux 服务器上大部分防火墙解决方案的基 础。netfilter 的内核 hook 和协议栈足够紧密，提供了包经过系统时的强大控制功能。 iptables 防火墙基于这些功能提供了一个灵活的、可扩展的、将策略需求转化到内核的方 法。理解了这些不同部分是如何联系到一起的，就可以使用它们控制和保护你的的服务器环 境

netfilter 是 Linux 内置的一种防火墙机制，我们一般也称之为数据包过滤机制。iptables 则是一个命令行工具，用来配置 netfilter 防火墙。

netfilter 与 iptables 的关系
netfilter 指整个项目，其官网叫 netfilter.org。在这个项目里面，netfilter 特指内核中的 netfilter 框架，iptables 指运行在用户态的配置工具。

netfilter 在协议栈中添加了一些钩子，它允许内核模块通过这些钩子注册回调函数，这样经过钩子的所有数据包都会被注册在相应钩子上的函数所处理，包括修改数据包内容、给数据包打标记或者丢掉数据包等。netfilter 框架负责维护钩子上注册的处理函数或者模块，以及它们的优先级。netfilter 框架负责在需要的时候动态加载其它的内核模块，比如 ip_conntrack、nf_conntrack、NAT subsystem 等。

iptables 是运行在用户态的一个程序，通过 netlink 和内核的 netfilter 框架打交道，并负责往钩子上配置回调函数。

netfilter 防火墙原理
简单说 netfilter 机制就是对进出主机的数据包进行过滤。 我们可以通过 iptables 设置一些规则(rules)。所有进出主机的数据包都会按照一定的顺序匹配这些规则，如果匹配到某条规则，就执行这条规则对应的行为，比如抛弃该数据包或接受该数据包。下图展示了 netfilter 依据 iptables 规则对数据包过滤的大致过程：



对数据包进行过滤。检查通过则接受(ACCEPT)数据包进入主机获取资源，如果检查不通过，则予以丢弃(DROP)！如果所有的规则都没有匹配上，就通过默认的策略(Policy)决定数据包的去向。注意，上图中的规则是有顺序的！比如数据包与 rule1 指定的规则匹配，那么这个数据包就会执行 action1 指定的行为，而不会继续匹配后面的规则了。

下面我们看一个例子。假设我们的 Linux 主机提供了 web 服务，所以需要放行访问 80 端口的数据包。
但是你发现来自 13.76.1.65 的数据包总是恶意的尝试入侵我们的 web 服务器，所以需要丢弃来自 13.76.1.65 数据包。
我们的 web 服务器并不提供 web 服务之外的其它服务，所以直接丢弃所有的非 web 请求的数据包。
总结后就是我们需要下面三条规则：

rule1 丢弃来自 13.76.1.65 数据包
rule2 接受访问 web 服务的数据包
rule3 丢弃所有的数据包
如果我们不小心把上面的规则顺序写错了，比如写成了下面的样子：

rule1 接受访问 web 服务的数据包
rule2 丢弃来自 13.76.1.65 数据包
rule3 丢弃所有的数据包
这时来自 13.76.1.65 的数据包是可以访问 web 服务的，因为来自 13.76.1.65 的数据包是符合第一条规则的，所以会被接受，此时就不会再考虑第二条规则了。

iptables 中的 table 与 chain
iptables 用表(table)来分类管理它的规则(rule)，这也是 iptables 名称的由来。根据 rule 的作用分成了好几个表，比如用来过滤数据包的 rule 就会放到 filter 表中，用于处理地址转换的 rule 就会放到 nat 表中，其中 rule 就是应用在 netfilter 钩子上的函数，用来修改数据包的内容或过滤数据包。下面我们简单的介绍下最常用的 filter 表和 nat 表。
filter
从名字就可以看出，filter 表里面的规则主要用来过滤数据，用来控制让哪些数据可以通过，哪些数据不能通过，它是最常用的表。
nat
里面的规则都是用来处理网络地址转换的，控制要不要进行地址转换，以及怎样修改源地址或目的地址，从而影响数据包的路由，达到连通的目的，这是路由器必备的功能。

下图展示了 iptables 中常用的 tables 及其 rule chains：



从上图可以看出，filter 和 nat 表中默认都存在着数条 rule chain。也就是说表中的规则(rule)又被编入了不同的链(chain)，由 chain 来决定什么时候触发 chain 上的这些规则。
iptables 里面有 5 个内置的 chain：

PREROUTING：接收的数据包刚进来，还没有经过路由选择，即还不知道数据包是要发给本机还是其它机器。这时会触发该 chain 上的规则。
INPUT：已经经过路由选择，并且该数据包的目的 IP 是本机，进入本地数据包处理流程。此时会触发该 chain 上的规则。
FORWARD：已经经过路由选择，但该数据包的目的 IP 不是本机，而是其它机器，进入 forward 流程。此时会触发该 chain 上的规则。
OUTPUT：本地程序要发出去的数据包刚到 IP 层，还没进行路由选择。此时会触发该 chain 上的规则。
POSTROUTING：本地程序发出去的数据包，或者转发(forward)的数据包已经经过了路由选择，即将交由下层发送出去。此时会触发该 chain 上的规则。
我们可以通过下图来理解这五条默认的规则链：



从上图可知，不考虑特殊情况的话，一个数据包只会经过下面三个路径中的一个：
A，主机收到目的 IP 是本机的数据包
B，主机收到目的 IP 不是本机的数据包
C，本机发出去的数据包

路径 A，数据包进入 Linux 主机访问其资源，在路由判断后确定是向 Linux 主机请求数据的数据包，此时主要是通过 filter 表的 INPUT 链来进行控制。
路径 B，数据包经由 Linux 主机转发，没有使用主机资源，而是流向后端主机。在路由判断之前对数据包进行表头的修改后，发现数据包需要透过防火墙去后端，此时的数据包就会通过路径B。也就是说，该封包的目标并非我们的 Linux 主机。该场景下数据包主要经过的链是 filter 表的 FORWARD 以及 nat 表的 POSTROUTING 和 PREROUTING。
路径 C，数据包由 Linux 主机向外发送。比如响应客户端的请求要求，或者是 Linux 主机主动发送出去的数据包，都会通过路径 C。它会先进行路由判断，在确定了输出的路径后，再通过 filter 表的 OUTPUT 链来传送。当然，最终还是会经过 nat 表的 POSTROUTING 链。
由此我们可以总结出下面的规律。
filter 表主要跟进入 Linux 主机的数据包有关，其 chains 如下：

INPUT：主要与想要进入 Linux 主机的数据包有关
OUTPUT：主要与 Linux 主机所要发送出去的数据包有关
FORWARD：与 Linux 主机没有关系，它可以对数据包进行转发

nat(地址转换) 表主要在进行来源与目的之 IP 或 port 的转换，其 chains 如下：

PREROUTING：在进行路由判断之前所要执行的规则(DNAT)
POSTROUTING：在进行路由判断之后所要执行的规则(SNAT)
OUTPUT：与发送出去的数据包有关
iptables 中的规则(rules)
规则(rules)存放在特定表的特定 chain 上，每条 rule 包含下面两部分信息：

Matching
Matching 就是如何匹配一个数据包，匹配条件很多，比如协议类型、源/目的IP、源/目的端口、in/out接口、包头里面的数据以及连接状态等，这些条件可以任意组合从而实现复杂情况下的匹配。

Targets
Targets 就是找到匹配的数据包之后怎么办，常见的有下面几种：

DROP：直接将数据包丢弃，不再进行后续的处理
RETURN： 跳出当前 chain，该 chain 里后续的 rule 不再执行
QUEUE： 将数据包放入用户空间的队列，供用户空间的程序处理
ACCEPT： 同意数据包通过，继续执行后续的 rule
跳转到其它用户自定义的 chain 继续执行
比如下面的规则，只要是来自内网的(192.168.1.0/24)数据包都被接受：

$ sudo iptables -A INPUT -i eth0 -s 192.168.1.0/24 -j ACCEPT
用户自定义 chains
除了 iptables 预定义的 5 个 chain 之外，用户还可以在表中定义自己的 chain，用户自定义的 chain 中的规则和预定义 chain 里的规则没有区别，不过由于自定义的 chain 没有和 netfilter 里面的钩子进行绑定，所以它不会自动触发，只能从其它 chain 的规则中跳转过来。


netfilter和iptables是什么关系？常说的iptables里面的表(table)、链(chain)、规则(rule)都是什么东西？本篇将带着这些疑问介绍netfilter/iptables的结构和相关概念，帮助有需要的同学更好的理解netfilter/iptables，为进一步学习使用iptables做准备。

什么是netfilter和iptables
用通俗点的话来讲:

netfilter指整个项目，不然官网就不会叫www.netfilter.org了。

在这个项目里面，netfilter特指内核中的netfilter框架，iptables指用户空间的配置工具。

netfilter在协议栈中添加了5个钩子，允许内核模块在这些钩子的地方注册回调函数，这样经过钩子的所有数据包都会被注册在相应钩子上的函数所处理，包括修改数据包内容、给数据包打标记或者丢掉数据包等。

netfilter框架负责维护钩子上注册的处理函数或者模块，以及它们的优先级。

iptables是用户空间的一个程序，通过netlink和内核的netfilter框架打交道，负责往钩子上配置回调函数。

netfilter框架负责在需要的时候动态加载其它的内核模块，比如 ip_conntrack、nf_conntrack、NAT subsystem等。

在应用者的眼里，可能iptables代表了整个项目，代表了防火墙，但在开发者眼里，可能netfilter更能代表这个项目。

netfilter钩子（hooks）
在内核协议栈中，有5个跟netfilter有关的钩子，数据包经过每个钩子时，都会检查上面是否注册有函数，如果有的话，就会调用相应的函数处理该数据包，它们的位置见下图：

         |
         | Incoming
         ↓
+-------------------+
| NF_IP_PRE_ROUTING |
+-------------------+
         |
         |
         ↓
+------------------+
|                  |         +----------------+
| routing decision |-------->| NF_IP_LOCAL_IN |
|                  |         +----------------+
+------------------+                 |
         |                           |
         |                           ↓
         |                  +-----------------+
         |                  | local processes |
         |                  +-----------------+
         |                           |
         |                           |
         ↓                           ↓
 +---------------+          +-----------------+
 | NF_IP_FORWARD |          | NF_IP_LOCAL_OUT |
 +---------------+          +-----------------+
         |                           |
         |                           |
         ↓                           |
+------------------+                 |
|                  |                 |
| routing decision |<----------------+
|                  |
+------------------+
         |
         |
         ↓
+--------------------+
| NF_IP_POST_ROUTING |
+--------------------+
         |
         | Outgoing
         ↓
NF_IP_PRE_ROUTING: 接收的数据包刚进来，还没有经过路由选择，即还不知道数据包是要发给本机还是其它机器。

NF_IP_LOCAL_IN: 已经经过路由选择，并且该数据包的目的IP是本机，进入本地数据包处理流程。

NF_IP_FORWARD: 已经经过路由选择，但该数据包的目的IP不是本机，而是其它机器，进入forward流程。

NF_IP_LOCAL_OUT: 本地程序要发出去的数据包刚到IP层，还没进行路由选择。

NF_IP_POST_ROUTING: 本地程序发出去的数据包，或者转发（forward）的数据包已经经过了路由选择，即将交由下层发送出去。

关于这些钩子更具体的位置，请参考Linux网络数据包的接收过程和数据包的发送过程

从上面的流程中，我们还可以看出，不考虑特殊情况的话，一个数据包只会经过下面三个路径中的一个：

本机收到目的IP是本机的数据包: NF_IP_PRE_ROUTING -> NF_IP_LOCAL_IN

本机收到目的IP不是本机的数据包: NF_IP_PRE_ROUTING -> NF_IP_FORWARD -> NF_IP_POST_ROUTING

本机发出去的数据包: NF_IP_LOCAL_OUT -> NF_IP_POST_ROUTING

注意： netfilter所有的钩子（hooks）都是在内核协议栈的IP层，由于IPv4和IPv6用的是不同的IP层代码，所以iptables配置的rules只会影响IPv4的数据包，而IPv6相关的配置需要使用ip6tables。

iptables中的表（tables）
iptables用表（table）来分类管理它的规则（rule），根据rule的作用分成了好几个表，比如用来过滤数据包的rule就会放到filter表中，用于处理地址转换的rule就会放到nat表中，其中rule就是应用在netfilter钩子上的函数，用来修改数据包的内容或过滤数据包。目前iptables支持的表有下面这些：

Filter
从名字就可以看出，这个表里面的rule主要用来过滤数据，用来控制让哪些数据可以通过，哪些数据不能通过，它是最常用的表。

NAT
里面的rule都是用来处理网络地址转换的，控制要不要进行地址转换，以及怎样修改源地址或目的地址，从而影响数据包的路由，达到连通的目的，这是家用路由器必备的功能。

Mangle
里面的rule主要用来修改IP数据包头，比如修改TTL值，同时也用于给数据包添加一些标记，从而便于后续其它模块对数据包进行处理（这里的添加标记是指往内核skb结构中添加标记，而不是往真正的IP数据包上加东西）。

Raw
在netfilter里面有一个叫做connection tracking的功能（后面会介绍到），主要用来追踪所有的连接，而raw表里的rule的功能是给数据包打标记，从而控制哪些数据包不被connection tracking所追踪。

Security
里面的rule跟SELinux有关，主要是在数据包上设置一些SELinux的标记，便于跟SELinux相关的模块来处理该数据包。

chains
上面我们根据不同功能将rule放到了不同的表里面之后，这些rule会注册到哪些钩子上呢？于是iptables将表中的rule继续分类，让rule属于不同的链（chain），由chain来决定什么时候触发chain上的这些rule。

iptables里面有5个内置的chains，分别对应5个钩子：

PREROUTING: 数据包经过NF_IP_PRE_ROUTING时会触发该chain上的rule.

INPUT: 数据包经过NF_IP_LOCAL_IN时会触发该chain上的rule.

FORWARD: 数据包经过NF_IP_FORWARD时会触发该chain上的rule.

OUTPUT: 数据包经过NF_IP_LOCAL_OUT时会触发该chain上的rule.

POSTROUTING: 数据包经过NF_IP_POST_ROUTING时会触发该chain上的rule.

每个表里面都可以包含多个chains，但并不是每个表都能包含所有的chains，因为某些表在某些chain上没有意义或者有些多余，比如说raw表，它只有在connection tracking之前才有意义，所以它里面包含connection tracking之后的chain就没有意义。（connection tracking的位置会在后面介绍到）

多个表里面可以包含同样的chain，比如在filter和raw表里面，都有OUTPUT chain，那应该先执行哪个表的OUTPUT chain呢？这就涉及到后面会介绍的优先级的问题。

提示：可以通过命令iptables -L -t nat|grep policy|grep Chain查看到nat表所支持的chain，其它的表也可以用类似的方式查看到，比如修改nat为raw即可看到raw表所支持的chain。

每个表（table）都包含哪些chain，表之间的优先级是怎样的？
下图在上面那张图的基础上，详细的标识出了各个表的rule可以注册在哪个钩子上（即各个表里面支持哪些chain），以及它们的优先级。

图中每个钩子关联的表按照优先级高低，从上到下排列；

图中将nat分成了SNAT和DNAT，便于区分；

图中标出了connection tracking（可以简单的把connection tracking理解成一个不能配置chain和rule的表，它必须放在指定位置，只能enable和disable）。

                                    |
                                    | Incoming             ++---------------------++
                                    ↓                      || raw                 ||
                           +-------------------+           || connection tracking ||
                           | NF_IP_PRE_ROUTING |= = = = = =|| mangle              ||
                           +-------------------+           || nat (DNAT)          ||
                                    |                      ++---------------------++
                                    |
                                    ↓                                                ++------------++
                           +------------------+                                      || mangle     ||
                           |                  |         +----------------+           || filter     ||
                           | routing decision |-------->| NF_IP_LOCAL_IN |= = = = = =|| security   ||
                           |                  |         +----------------+           || nat (SNAT) ||
                           +------------------+                 |                    ++------------++
                                    |                           |
                                    |                           ↓
                                    |                  +-----------------+
                                    |                  | local processes |
                                    |                  +-----------------+
                                    |                           |
                                    |                           |                    ++---------------------++
 ++------------++                   ↓                           ↓                    || raw                 ||
 || mangle     ||           +---------------+          +-----------------+           || connection tracking ||
 || filter     ||= = = = = =| NF_IP_FORWARD |          | NF_IP_LOCAL_OUT |= = = = = =|| mangle              ||
 || security   ||           +---------------+          +-----------------+           || nat (DNAT)          ||
 ++------------++                   |                           |                    || filter              ||
                                    |                           |                    || security            ||
                                    ↓                           |                    ++---------------------++
                           +------------------+                 |
                           |                  |                 |
                           | routing decision |<----------------+
                           |                  |
                           +------------------+
                                    |
                                    |
                                    ↓
                           +--------------------+           ++------------++
                           | NF_IP_POST_ROUTING |= = = = = =|| mangle     ||
                           +--------------------+           || nat (SNAT) ||
                                    |                       ++------------++
                                    | Outgoing
                                    ↓
以NF_IP_PRE_ROUTING为例，数据包到了这个点之后，会先执行raw表中PREROUTING(chain)里的rule，然后执行connection tracking，接着再执行mangle表中PREROUTING(chain)里的rule，最后执行nat (DNAT)表中PREROUTING(chain)里的rule。

以filter表为例，它只能注册在NF_IP_LOCAL_IN、NF_IP_FORWARD和NF_IP_LOCAL_OUT上，所以它只支持INPUT、FORWARD和OUTPUT这三个chain。

以收到目的IP是本机的数据包为例，它的传输路径为：NF_IP_PRE_ROUTING -> NF_IP_LOCAL_IN，那么它首先要依次经过NF_IP_PRE_ROUTING上注册的raw、connection tracking 、mangle和nat (DNAT)，然后经过NF_IP_LOCAL_IN上注册的mangle、filter、security和nat (SNAT)。

iptables中的规则（Rules）
rule存放在特定表的特定chain上，每条rule包含下面两部分信息：

Matching
Matching就是如何匹配一个数据包，匹配条件很多，比如协议类型、源/目的IP、源/目的端口、in/out接口、包头里面的数据以及连接状态等，这些条件可以任意组合从而实现复杂情况下的匹配。详情请参考Iptables matches

Targets
Targets就是找到匹配的数据包之后怎么办，常见的有下面几种：

DROP：直接将数据包丢弃，不再进行后续的处理

RETURN： 跳出当前chain，该chain里后续的rule不再执行

QUEUE： 将数据包放入用户空间的队列，供用户空间的程序处理

ACCEPT： 同意数据包通过，继续执行后续的rule

跳转到其它用户自定义的chain继续执行

当然iptables包含的targets很多很多，但并不是每个表都支持所有的targets，
rule所支持的target由它所在的表和chain以及所开启的扩展功能来决定，具体每个表支持的targets请参考Iptables targets and jumps。

用户自定义Chains
除了iptables预定义的5个chain之外，用户还可以在表中定义自己的chain，用户自定义的chain中的rule和预定义chain里的rule没有区别，不过由于自定义的chain没有和netfilter里面的钩子进行绑定，所以它不会自动触发，只能从其它chain的rule中跳转过来。

连接追踪（Connection Tracking）
Connection Tracking发生在NF_IP_PRE_ROUTING和NF_IP_LOCAL_OUT这两个地方，一旦开启该功能，Connection Tracking模块将会追踪每个数据包（被raw表中的rule标记过的除外），维护所有的连接状态，然后这些状态可以供其它表中的rule引用，用户空间的程序也可以通过/proc/net/ip_conntrack来获取连接信息。下面是所有的连接状态：

这里的连接不仅仅是TCP的连接，两台设备的进程用UDP和ICMP（ping）通信也会被认为是一个连接

NEW: 当检测到一个不和任何现有连接关联的新包时，如果该包是一个合法的建立连接的数据包（比如TCP的sync包或者任意的UDP包），一个新的连接将会被保存，并且标记为状态NEW。

ESTABLISHED: 对于状态是NEW的连接，当检测到一个相反方向的包时，连接的状态将会由NEW变成ESTABLISHED，表示连接成功建立。对于TCP连接，意味着收到了一个SYN/ACK包， 对于UDP和ICMP，任何反方向的包都可以。

RELATED: 数据包不属于任何现有的连接，但它跟现有的状态为ESTABLISHED的连接有关系，对于这种数据包，将会创建一个新的连接，且状态被标记为RELATED。这种连接一般是辅助连接，比如FTP的数据传输连接（FTP有两个连接，另一个是控制连接），或者和某些连接有关的ICMP报文。

INVALID: 数据包不和任何现有连接关联，并且不是一个合法的建立连接的数据包，对于这种连接，将会被标记为INVALID，一般这种都是垃圾数据包，比如收到一个TCP的RST包，但实际上没有任何相关的TCP连接，或者别的地方误发过来的ICMP包。

UNTRACKED: 被raw表里面的rule标记为不需要tracking的数据包，这种连接将会标记成UNTRACKED。
