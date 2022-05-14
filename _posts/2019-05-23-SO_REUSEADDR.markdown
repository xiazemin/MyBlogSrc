---
title: SO_REUSEADDR SO_REUSEPORT
layout: post
category: linux
author: 夏泽民
---
每次kill掉该服务器进程并重新启动的时候，都会出现bind错误：error:98，Address already in use。然而再kill掉该进程，再次重新启动的时候，就bind成功了。
端口正在被使用，并处于TCP中的TIME_WAIT状态。再过两分钟，我再执行命令netstat -an|grep 9877


SO_REUSEADDR和SO_REUSEPORT主要是影响socket绑定ip和port的成功与否。
绑定规则
规则1：socket可以指定绑定到一个特定的ip和port，例如绑定到192.168.0.11:9000上；
规则2：同时也支持通配绑定方式，即绑定到本地"any address"（例如一个socket绑定为 0.0.0.0:21，那么它同时绑定了所有的本地地址）；
规则3：默认情况下，任意两个socket都无法绑定到相同的源IP地址和源端口。
<!-- more -->
SO_REUSEADDR

SO_REUSEADDR的作用主要包括两点
1、改变了通配绑定时处理源地址冲突的处理方式，其具体的表现方式为：未设置SO_REUSEADDR时，socketA先绑定到0.0.0.0:21，后socketB绑定192.168.0.1:21将失败，不符合规则3。但在设置SO_REUSEADDR后socketB将绑定成功。并且这个设置对于socketA（通配绑定）和socketB（特定绑定）的绑定是顺序无关的。

对于linux，要使这个设置达到预期效果，对于绑定的顺序的有要求的，即在设置了SO_REUSEADDR，须先进行特定绑定，后进行通配绑定，后者才能成功；如果先进行通配绑定，后面的绑定（端口相同情况下）地址只要和通配绑定中的一个相同都将失败。

2、改变了系统对处于TIME_WAIT状态的socket的看待方式，要理解这个句话，首先先简单介绍以下什么是处于TIME_WAIT状态的socket？

socket通常都有发送缓冲区，当调用send()函数成功后，只是将数据放到了缓冲区，并不意味着所有数据真正被发送出去。对于TCP socket，在加入缓冲区和真正被发送之间的时延会相当长。这就导致当close一个TCP socket的时候，可能在发送缓冲区中保存着等待发送的数据。为了确保TCP的可靠传输，TCP的实现是close一个TCP socket时，如果它仍然有数据等待发送，那么该socket会进入TIME_WAIT状态。这种状态将持续到数据被全部发送或者发生超时（这个超时时间通常被称为Linger Time，大多数系统默认为2分钟）
在未设置SO_REUSEADDR时，内核将一个处于TIME_WAIT状态的socketA仍然看成是一个绑定了指定ip和port的有效socket，因此此时如果另外一个socketB试图绑定相同的ip和port都将失败（不满足规则3），直到socketA被真正释放后，才能够绑定成功。如果socketB设置SO_REUSEADDR（仅仅只需要socketB进行设置），这种情况下socketB的绑定调用将成功返回，但真正生效需要在socketA被真正释放后。（这个地方的理解可能有点问题，待后续验证一下）。总结一下：内核在处理一个设置了SO_REUSEADDR的socket绑定时，如果其绑定的ip和port和一个处于TIME_WAIT状态的socket冲突时，内核将忽略这种冲突，即改变了系统对处于TIME_WAIT状态的socket的看待方式。
SO_REUSEPORT

SO_REUSEPORT作用就比较明显直观，即打破了上面的规则3
1、允许将多个socket绑定到相同的地址和端口，前提每个socket绑定前都需设置SO_REUSEPORT。如果第一个绑定的socket未设置SO_REUSEPORT，那么其他的socket无论有没有设置SO_REUSEPORT都无法绑定到该地址和端口直到第一个socket释放了绑定。

2、attention：SO_REUSEPORT并不表示SO_REUSEADDR，即不具备上述SO_REUSEADDR的第二点作用（对TIME_WAIT状态的socket处理方式）。因此当有个socketA未设置SO_REUSEPORT绑定后处在TIME_WAIT状态时，如果socketB仅设置了SO_REUSEPORT在绑定和socketA相同的ip和端口时将会失败。解决方案
（1）、socketB设置SO_REUSEADDR 或者socketB即设置SO_REUSEADDR也设置SO_REUSEPORT
（2）、两个socket上都设置SO_REUSEPORT
Linux 内核3.9加入了SO_REUSEPORT。除上述功能外其额外实现了
1、为了阻止port 劫持Port hijacking，限制所有使用相同ip和port的socket都必须拥有相同的有效用户id(effective user ID)。
2、linux内核在处理SO_REUSEPORT socket的集合时，进行了简单的负载均衡操作，即对于UDP socket，内核尝试平均的转发数据报，对于TCP监听socket，内核尝试将新的客户连接请求(由accept返回)平均的交给共享同一地址和端口的socket(监听socket)。

TCP/UDP是由以下五元组唯一地识别的： 
{<protocol>, <src addr>, <src port>, <dest addr>, <dest port>}

这些数值组成的任何独特的组合可以唯一地确一个连接。那么，对于任意连接，这五个值都不能完全相同。否则的话操作系统就无法区别这些连接了。

一个socket的协议是在用socket()初始化的时候就设置好的。源地址（source address）和源端口（source port）在调用bind()的时候设置。目的地址（destination address）和目的端口（destination port）在调用connect()的时候设置。其中UDP是无连接的，UDP socket可以在未与目的端口连接的情况下使用。但UDP也可以在某些情况下先与目的地址和端口建立连接后使用。在使用无连接UDP发送数据的情况下，如果没有显式地调用bind()，草错系统会在第一次发送数据时自动将UDP socket与本机的地址和某个端口绑定（否则的话程序无法接受任何远程主机回复的数据）。同样的，一个没有绑定地址的TCP socket也会在建立连接时被自动绑定一个本机地址和端口。

如果我们手动绑定一个端口，我们可以将socket绑定至端口0，绑定至端口0的意思是让系统自己决定使用哪个端口（一般是从一组操作系统特定的提前决定的端口数范围中），所以也就是任何端口的意思。同样的，我们也可以使用一个通配符来让系统决定绑定哪个源地址（ipv4通配符为0.0.0.0，ipv6通配符为::）。而与端口不同的是，一个socket可以被绑定到主机上所有接口所对应的地址中的任意一个。基于连接在本socket的目的地址和路由表中对应的信息，操作系统将会选择合适的地址来绑定这个socket，并用这个地址来取代之前的通配符IP地址。

在默认情况下，任意两个socket不能被绑定在同一个源地址和源端口组合上。比如说我们将socketA绑定在A:X地址，将socketB绑定在B:Y地址，其中A和B是IP地址，X和Y是端口。那么在A==B的情况下X!=Y必须满足，在X==Y的情况下A!=B必须满足。需要注意的是，如果某一个socket被绑定在通配符IP地址下，那么事实上本机所有IP都会被系统认为与其绑定了。例如一个socket绑定了0.0.0.0:21，在这种情况下，任何其他socket不论选择哪一个具体的IP地址，其都不能再绑定在21端口下。因为通配符IP0.0.0.0与所有本地IP都冲突。


