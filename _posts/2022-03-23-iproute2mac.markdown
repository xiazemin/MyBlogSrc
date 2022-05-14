---
title: iproute2mac IProute2
layout: post
category: linux
author: 夏泽民
---
OS X是BSD系的内核，高级路由用pf（旧版的可能有ipfw可以用）

brew install iproute2mac

https://github.com/brona/iproute2mac
<!-- more -->
https://www.zhihu.com/question/22761563


官方仓库：https://github.com/shemminger/iproute2

https://blog.csdn.net/m0_58562922/article/details/123549691
https://wiki.linuxfoundation.org/networking/iproute2

$ git clone git://git.kernel.org/pub/scm/network/iproute2/iproute2.git
or

 $ git clone https://github.com/shemminger/iproute2.git
https://github.com/shemminger/iproute2

The Docker Engine uses Linux-specific kernel features, so to run it on OS X we need to use a lightweight virtual machine (vm). You use the OS X Docker client to control the virtualized Docker Engine to build, run, and manage Docker containers.

https://serverfault.com/questions/603763/is-it-possible-to-do-linux-network-namespaces-netns-on-macos

“Object "netns" is unknown, try "ip help".\n'”报错
By default, CentOS 6.4 does not support network namespaces. If one wants totest the new virtualization platforms (Docker, OpenStack, & co…) on aCentOS server, all features won’t be available.
For OpenStack for example, Neutron won’t work as expected, since it actuallyneeds network namespace to create networks,

https://blog.csdn.net/weixin_33978016/article/details/91833057

https://serverfault.com/questions/603763/is-it-possible-to-do-linux-network-namespaces-netns-on-macos

https://apple.stackexchange.com/questions/429079/does-macos-have-network-namespaces-like-linux


https://www.unix.com/man-page/osx/8/ip-netns/

ip [ OPTIONS ] netns  { COMMAND | help }

       ip netns [ list ]

       ip netns add NETNSNAME

       ip [-all] netns del [ NETNSNAME ]

       ip netns set NETNSNAME NETNSID

       ip netns identify [ PID ]

       ip netns pids NETNSNAME

       ip [-all] netns exec [ NETNSNAME ] command...

       ip netns monitor

       ip netns list-id
  
  https://github.com/vishvananda/netns
  



