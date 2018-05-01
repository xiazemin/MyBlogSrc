---
title: linux namespace
layout: post
category: linux
author: 夏泽民
---
Linux Namespace是Linux提供的一种内核级别环境隔离的方法。不知道你是否还记得很早以前的Unix有一个叫chroot的系统调用（通过修改根目录把用户jail到一个特定目录下），chroot提供了一种简单的隔离模式：chroot内部的文件系统无法访问外部的内容。Linux Namespace在此基础上，提供了对UTS、IPC、mount、PID、network、User等的隔离机制。
举个例子，我们都知道，Linux下的超级父亲进程的PID是1，所以，同chroot一样，如果我们可以把用户的进程空间jail到某个进程分支下，并像chroot那样让其下面的进程 看到的那个超级父进程的PID为1，于是就可以达到资源隔离的效果了（不同的PID namespace中的进程无法看到彼此）
Linux Namespace 有如下种类，官方文档在这里《Namespace in Operation》
分类	系统调用参数	相关内核版本
Mount namespaces	CLONE_NEWNS	Linux 2.4.19
UTS namespaces	CLONE_NEWUTS	Linux 2.6.19
IPC namespaces	CLONE_NEWIPC	Linux 2.6.19
PID namespaces	CLONE_NEWPID	Linux 2.6.24
Network namespaces	CLONE_NEWNET	始于Linux 2.6.24 完成于 Linux 2.6.29
User namespaces	CLONE_NEWUSER	始于 Linux 2.6.23 完成于 Linux 3.8)
主要是三个系统调用
clone() – 实现线程的系统调用，用来创建一个新的进程，并可以通过设计上述参数达到隔离。
unshare() – 使某进程脱离某个namespace
setns() – 把某进程加入到某个namespace
User Namespace主要是用了CLONE_NEWUSER的参数。使用了这个参数后，内部看到的UID和GID已经与外部不同了，默认显示为65534。那是因为容器找不到其真正的UID所以，设置上了最大的UID（其设置定义在/proc/sys/kernel/overflowuid）。
要把容器中的uid和真实系统的uid给映射在一起，需要修改 /proc/<pid>/uid_map 和 /proc/<pid>/gid_map 这两个文件。这两个文件的格式为：
ID-inside-ns ID-outside-ns length
其中：
第一个字段ID-inside-ns表示在容器显示的UID或GID，
第二个字段ID-outside-ns表示容器外映射的真实的UID或GID。
第三个字段表示映射的范围，一般填1，表示一一对应。
比如，把真实的uid=1000映射成容器内的uid=0
$ cat /proc/2465/uid_map
         0       1000          1
再比如下面的示例：表示把namespace内部从0开始的uid映射到外部从0开始的uid，其最大范围是无符号32位整形
$ cat /proc/$$/uid_map
         0          0          4294967295
另外，需要注意的是：
写这两个文件的进程需要这个namespace中的CAP_SETUID (CAP_SETGID)权限（可参看Capabilities）
写入的进程必须是此user namespace的父或子的user namespace进程。
另外需要满如下条件之一：1）父进程将effective uid/gid映射到子进程的user namespace中，2）父进程如果有CAP_SETUID/CAP_SETGID权限，那么它将可以映射到父进程中的任一uid/gid。
<!-- more -->
当你启动一个Docker容器后，你可以使用ip link show或ip addr show来查看当前宿主机的网络情况（我们可以看到有一个docker0，还有一个veth22a38e6的虚拟网卡——给容器用的） 
 # 把容器里的 veth-ns1改名为 eth0 （容器外会冲突，容器内就不会了）
ip netns exec ns1  ip link set dev veth-ns1 name eth0 
 # 为容器中的网卡分配一个IP地址，并激活它
ip netns exec ns1 ifconfig eth0 192.168.10.11/24 up
 # 上面我们把veth-ns1这个网卡按到了容器中，然后我们要把lxcbr0.1添加上网桥上
brctl addif lxcbr0 lxcbr0.1
 # 为容器增加一个路由规则，让容器可以访问外面的网络
ip netns exec ns1     ip route add default via 192.168.10.1
 # 在/etc/netns下创建network namespce名称为ns1的目录，
 # 然后为这个namespace设置resolv.conf，这样，容器内就可以访问域名了
mkdir -p /etc/netns/ns1
echo "nameserver 8.8.8.8" > /etc/netns/ns1/resolv.conf
上面基本上就是docker网络的原理了，只不过，
Docker的resolv.conf没有用这样的方式，而是用了Mount Namesapce的那种方式
另外，docker是用进程的PID来做Network Namespace的名称的。

甚至可以为正在运行的docker容器增加一个新的网卡：
ip link add peerA type veth peer name peerB 
brctl addif docker0 peerA 
ip link set peerA up 
ip link set peerB netns ${container-pid} 
ip netns exec ${container-pid} ip link set dev peerB name eth1 
ip netns exec ${container-pid} ip link set eth1 up ; 
ip netns exec ${container-pid} ip addr add ${ROUTEABLE_IP} dev eth1 ;
上面的示例是我们为正在运行的docker容器，增加一个eth1的网卡，并给了一个静态的可被外部访问到的IP地址。
这个需要把外部的“物理网卡”配置成混杂模式，这样这个eth1网卡就会向外通过ARP协议发送自己的Mac地址，然后外部的交换机就会把到这个IP地址的包转到“物理网卡”上，因为是混杂模式，所以eth1就能收到相关的数据，一看，是自己的，那么就收到。这样，Docker容器的网络就和外部通了。
当然，无论是Docker的NAT方式，还是混杂模式都会有性能上的问题，NAT不用说了，存在一个转发的开销，混杂模式呢，网卡上收到的负载都会完全交给所有的虚拟网卡上，于是就算一个网卡上没有数据，但也会被其它网卡上的数据所影响。
这两种方式都不够完美，我们知道，真正解决这种网络问题需要使用VLAN技术，于是Google的同学们为Linux内核实现了一个IPVLAN的驱动，这基本上就是为Docker量身定制的。
Namespace文件
上面就是目前Linux Namespace的玩法。 现在，我来看一下其它的相关东西。
让我们运行一下上篇中的那个pid.mnt的程序（也就是PID Namespace中那个mount proc的程序），然后不要退出。
$ sudo ./pid.mnt 
[sudo] password for hchen: 
Parent [ 4599] - start a container!
Container [    1] - inside the container!
我们到另一个shell中查看一下父子进程的PID：
 $ pstree -p 4599
pid.mnt(4599)───bash(4600)
我们可以到proc下（/proc//ns）查看进程的各个namespace的id（内核版本需要3.8以上）。
下面是父进程的：
 $ sudo ls -l /proc/4599/ns
total 0
lrwxrwxrwx 1 root root 0  4月  7 22:01 ipc -> ipc:[4026531839]
lrwxrwxrwx 1 root root 0  4月  7 22:01 mnt -> mnt:[4026531840]
lrwxrwxrwx 1 root root 0  4月  7 22:01 net -> net:[4026531956]
lrwxrwxrwx 1 root root 0  4月  7 22:01 pid -> pid:[4026531836]
lrwxrwxrwx 1 root root 0  4月  7 22:01 user -> user:[4026531837]
lrwxrwxrwx 1 root root 0  4月  7 22:01 uts -> uts:[4026531838]
下面是子进程的：
 $ sudo ls -l /proc/4600/ns
total 0
lrwxrwxrwx 1 root root 0  4月  7 22:01 ipc -> ipc:[4026531839]
lrwxrwxrwx 1 root root 0  4月  7 22:01 mnt -> mnt:[4026532520]
lrwxrwxrwx 1 root root 0  4月  7 22:01 net -> net:[4026531956]
lrwxrwxrwx 1 root root 0  4月  7 22:01 pid -> pid:[4026532522]
lrwxrwxrwx 1 root root 0  4月  7 22:01 user -> user:[4026531837]
lrwxrwxrwx 1 root root 0  4月  7 22:01 uts -> uts:[4026532521]
我们可以看到，其中的ipc，net，user是同一个ID，而mnt,pid,uts都是不一样的。如果两个进程指向的namespace编号相同，就说明他们在同一个namespace下，否则则在不同namespace里面。
这些文件还有另一个作用，那就是，一旦这些文件被打开，只要其fd被占用着，那么就算PID所属的所有进程都已经结束，创建的namespace也会一直存在。比如：我们可以通过：mount –bind /proc/4600/ns/uts ~/uts 来hold这个namespace。
另外，我们在上篇中讲过一个setns的系统调用，其函数声明如下：
int setns(int fd, int nstype);
其中第一个参数就是一个fd，也就是一个open()系统调用打开了上述文件后返回的fd，比如：
fd = open("/proc/4600/ns/nts", O_RDONLY);  // 获取namespace文件描述符
setns(fd, 0); // 加入新的namespace
