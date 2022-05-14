---
title: route print netstat nr
layout: post
category: linux
author: 夏泽民
---
route print命令可以查看路由表，在dos下面输入route print 就可以了，如何读懂路由表
Active Routes:   
Network Destination        Netmask          Gateway       Interface  Metric   
          0.0.0.0          0.0.0.0    202.256.257.1  202.256.257.258      1   
        127.0.0.0        255.0.0.0        127.0.0.1       127.0.0.1       1   
    202.256.257.0    255.255.255.0  202.256.257.258  202.256.257.258      1   
  202.256.257.258  255.255.255.255        127.0.0.1       127.0.0.1       1   
  202.256.257.255  255.255.255.255  202.256.257.258  202.256.257.258      1   
        224.0.0.0        224.0.0.0  202.256.257.258  202.256.257.258      1   
  255.255.255.255  255.255.255.255  202.256.257.258  202.256.257.258      1   
Default Gateway:     202.256.257.1   

首先是最上方给出了接口列表，一个本地循环，一个网卡接口，网卡结构给出了网卡的mac地址。

再说说每一列的内容吧，从做到右依次是：Network Destination（目的地址），Netmask（掩码），Gateway（网关），Interface（接口），Metric（也不知道是什么，姑且认为是一个度量值或是管理距离）。 

下面说说每一行内容代表的内容，首先是  
Network Destination        Netmask          Gateway       Interface  Metric  
          0.0.0.0          0.0.0.0    202.256.257.1  202.256.257.258      1  
这表示发向任意网段的数据通过本机接口202.256.257.258被送往一个默认的网关：202.256.257.1，它的管理距离是1，这里对管理距离说说，管理距离指的是在路径选择的过程中信息的可信度，管理距离越小的，可信度越高。 

再看看第二行  
Network Destination        Netmask          Gateway       Interface  Metric  
        127.0.0.0        255.0.0.0        127.0.0.1       127.0.0.1       1  
A类地址中127.0.0.0留住本地调试使用，所以路由表中所以发向127.0.0.0网络的数据通过本地回环127.0.0.1发送给指定的网关：127.0.0.1，也就是从自己的回环接口发到自己的回环接口，这将不会占用局域网带宽。 

第三行  
Network Destination        Netmask          Gateway       Interface  Metric  
    202.256.257.0    255.255.255.0  202.256.257.258  202.256.257.258      1  
这里的目的网络与本机处于一个局域网，所以发向网络202.256.257.0（也就是发向局域网的数据）使用本机:202.256.257.258作为网关，这便不再需要路由器路由或不需要交换机交换，增加了传输效率。 

 

第四行  
Network Destination        Netmask          Gateway       Interface  Metric  
  202.256.257.258  255.255.255.255        127.0.0.1       127.0.0.1       1  
表示从自己的主机发送到自己主机的数据包，如果使用的是自己主机的IP地址，跟使用回环地址效果相同，通过同样的途径被路由，也就是如果我有自己的站点，我要浏览自己的站点，在IE地质栏里面输入localhost与202.256.257.258是一样的，尽管localhost被解析为 127.0.0.1。 

 

第五行  
Network Destination        Netmask          Gateway       Interface  Metric  
  202.256.257.255  255.255.255.255  202.256.257.258  202.256.257.258      1  
这里的目的地址是一个局域广播地址，系统对这样的数据包的处理方法是把本机202.256.257.258作为网关，发送局域广播帧，这个帧将被路由器过滤。 

第六行  
Network Destination        Netmask          Gateway       Interface  Metric  
        224.0.0.0        224.0.0.0  202.256.257.258  202.256.257.258      1  
这里的目的地址是一个组播（muticast）网络，组播指的是数据包同时发向几个指定的IP地址，其他的地址不会受到影响。系统的处理依然是适用本机作为网关，进行路由。这里有一点要说明的组播可被路由器转发，如果路由器不支持组播，则采用广播方式转发。 

最后一行  
Network Destination        Netmask          Gateway       Interface  Metric  
  255.255.255.255  255.255.255.255  202.256.257.258  202.256.257.258      1  
目的地址是一个广域广播，同样适用本机为网关，广播广播帧，这样的包到达路由器之后被转发还是丢弃根据路由器的配置决定。 

还有个半行没有解释  
Default Gateway:     202.256.257.1 

这是一个缺省的网关，要是发送的数据的目的地址根前面例举的都不匹配的时候，就将数据发送到这个缺省网关，由其决定路由
https://blog.csdn.net/jinrich/article/details/5146307
<!-- more -->
    [ Route就是用来显示、人工添加和修改路由表项目的。]大多数主机一般都是驻留在只连接一台路由器的网段上。由于只有一台路由器，因此不存在使用哪一台路由器将数据包发表到远程计算机上去的问题，该路由器的IP地址可作为该网段上所有计算机的缺省网关来输入。 但是，当网络上拥有两个或多个路由器时，你就不一定想只依赖缺省网关了。实际上你可能想让你的某些远程IP地址通过某个特定的路由器来传递，而其他的远程IP则通过另一个路由器来传递。
https://blog.csdn.net/weixin_44893633/article/details/103683634

http://blog.sina.com.cn/s/blog_493cafbb0101hml2.html

https://www.cnblogs.com/youxin233/archive/2010/02/09/1666192.html

https://superuser.com/questions/521791/what-is-an-interface-address


mac上面查看路由表
本来想使用linux上面的命令route -n查看mac上面的路由表的，结果显示mac上面的route命令不是这样玩的。
netstat -nr
解决
Mac上面需要使用netstat的命令来解决route不能查看路由表的问题。
https://blog.csdn.net/weixin_33958585/article/details/91773534

查看所有路由表
netstat -rn
http://witmax.cn/mac-route.html
http://blog.joylau.cn/2018/12/14/MacOS-Route/

http://ciscoelearning.blogspot.com/2009/03/router-interface.html

https://blog.csdn.net/chengqiuming/article/details/80489180

Neutron的Router模型中，蕴含着三种路由：直连路由、默认静态路由和静态路由。前两种路由不需要显示地增加路由表项，也不会体现在路由表（routers）中，当增加一个Port时（add_router_interface），Neutron会自动增加一个直连路由；当增加一个外部网关信息时（external_gateway_info）,Neutron会增加一个默认静态路由。
路由表中的路由也是静态路由，它与默认路由一样，都是通往外部网络。外部网络，指的是neutron管理范围之外的网络，不过，静态路由中的外部网络，一般指的是私网，而默认静态路由中的外部网络，一般指的是公网。所有，在外部网关信息中，有一个开关enable_snat，当它为true时，需要启动SNAT。
外部网关信息（external_gateway_info）中，只有外部网关IP（蕴含在subnet_id所对应的subnet的gateway_ip字段中），而无目的网段，所以被称为默认静态路由。同时，也只有创建外部网关信息时，Neutron才会自动在对应的Router上创建一个Port，其余场景都需要人主动创建Port作关联。
Float IP首先是一个SNAT/DNAT转换规则：floating_ip_address(外网/公网IP)与fixed_ip_address（内网/私网）互相转换。然后，从实现角度来讲，它才是绑定到一个Router（router_id）的端口（Port）上，以让报文在进出这个端口时，Router能对其做SNAT/DNAT转换。

https://blog.csdn.net/chengqiuming/article/details/80500300

