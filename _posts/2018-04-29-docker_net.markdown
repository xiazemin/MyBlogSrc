---
title: docker_net
layout: post
category: docker
author: 夏泽民
---
由于Boot2Docker的存在，造成了三层Mac->VirtualBox->Docker网络,由VirtualBox到Docker的映射可以通过run容器的时候指定-p参数实现，而从宿主机到VirtualBox端口映射需要通过下述方法实现:
查询虚拟机网络： 
VBoxManagelistvms查询虚拟机网络状态，默认虚拟机名为′default′ VBoxManage showvminfo “default” | grep NIC 
2.关闭运行中的虚拟机 
由于Boot2Docker会自动运行VirtualBox中的虚拟机，所以在设置网络映射时必须先关闭运行中的虚拟机。否则，将出现The machine ‘boot2docker’ is already locked for a session (or being unlocked)的错误提示
$ VBoxManage controlvm "default" poweroff
修改虚拟机与Mac系统的网络映射 
根据实际需要进行网络映射，其中
rulename: 自定义规则名称
hostip: Mac访问地址，可不填
hostport: Mac映射端口
guestip: 虚拟机访问地址，可不填
guestport: 虚拟机映射端口
$ VBoxManage modifyvm “dufault” –natpf1 “,

Docker 在一个包装中联合了以上功能,并称之为容器格式。
Libcontainer
默认的容器格式被称为 libcontainer 。
Docker 也支持使用 LXC 的传统 Linux 容器。在将来, Docker 可能会支持其他 的容器格式,比如结合 BSD jails或者 Solaris Zones。
执行驱动程序是一种特殊容器格式的实现,用来运行 docker 容器。在最新的版 本中, libcontainer 有以下特性:
∙是运行 docker 容器的默认执行驱动程序。
∙和 LXC 同时装载。
∙使用没有任何其他依赖关系的 Go 语言设计的库,来直接访问内核容器的 API 。 ▪在 Docker 0.9中, LXC 现在可以选择关闭。
▪注意:LXC 在将来会继续被支持。
▪如果想要重新使用 LXC 驱动,只需输入指令 docker -d – e lxc, 然后重启Docker 。
▪目前的 Docker 涵盖的功能有:命名空间使用, cgroups 管理, capabilities 权限 集,进程运行的环境变量配置以及网络接口防火墙设置 —— 所有功能是固定可 预测的,不依赖 LXC 或者其它任何用户区软件包。
▪只需提供一个根文件系统,和 libcontainer 对容器的操作配置,它会帮你完成剩 下的事情。
▪支持新建容器或者添加到现有的容器。
▪事实上,对 libcontainer 最迫切的需求是稳定,开发团队也将其设为了默认。
用户命名空间
Docker 不是虚拟化,相反的,它是一个支持命名空间抽象的内核 , 提供了独立工作空间 (或容器 ) 。当你运行一个容器的时候, Docker 为容器新建了一系列的 namespace 。
▪一些 Docker 使用的 linux 命名空间:
∙pid namespace
▪用作区分进程(PID: Process ID)。
▪容器中运行的进程就如同在普通的 Linux 系统运行一样,尽管它们和其他进程 共享一个底层内核。
∙net namespace
▪用作管理网络接口。
▪DNAT 允许你单独配置主机中每个用户的的网络,并且有一个方便的接口传输 它们之间的数据。
▪当然,你也可以通过使用网桥用物理接口替换它。
∙ipc namespace
▪用作管理对 IPC (IPC: InterProcess Communication)资源的访问。
∙mnt namespace
▪用作管理 mount-points (MNT: Mount)。
∙uts namespace
▪用作区分内核和版本标识符 (UTS: Unix Timesharing System)。
Cgroups
Linux 上的 Docker 使用了被称为 cgroups 的技术。因为每个虚拟机都是一个进 程,所有普通 Linux 的资源管理应用可以被应用到虚拟机。此外,资源分配和 调度只有一个等级,因为一个容器化的 Linux 系统只有一个内核并且这个内核 对容器完全可见。
∙总之, cgroups 可以让 Docker :
∙实现组进程并且管理它们的资源总消耗。
∙分享可用的硬件资源到容器。
∙限制容器的内存和 CPU 使用。
▪可以通过更改相应的 cgroup 来调整容器的大小。
▪通过检查 Linux 中的 /sys/fs/cgroup对照组来获取容器中的资源使用信息。 ∙提供了一种可靠的结束容器内所有进程的方法。
Capabilities
Linux 使用的是 “POSIX capabilities” 。这些权限是所有强大的 root 权限分割而成的一系 列权限。在 Linux manpages上可以找到所有可用权限的清单。 Docker 丢弃了除了所 需权限外的所有权限,使用了白名单而不是黑名单。
一般服务器(裸机或者虚拟机)需要以 root 权限运行一系列进程。包括:
∙SSH
∙cron
∙syslogd
∙硬件管理工具 (比如负载模块 )
∙网络配置工具 (比如处理 DHCP, WPA, or VPNs)等。
每个容器都是不同的,因为几乎所有这些任务都由围绕容器的基础设施进行处理。默 认的, Docker 启用一个严格限制权限的容器。大多数案例中,容器不需要真正的 root 权限。举个例子,进程(比如说网络服务)只需要绑定一个小于 1024的端口而不需要 root 权限:他们可以被授予 CAP_NET_BIND_SERVICE来代替。因此,容器可以被降 权运行:意味着容器中的 root 权限比真正的 root 权限拥有更少的特权。
Capabilities 只是现代 Linux 内核提供的众多安全功能中的一个。为了加固一个 Docker 主机,你可以使用现有知名的系统:
∙TOMOYO
∙AppArmor
∙SELinux
∙GRSEC, etc.
如果你的发行版本附带了 Docker 容器的安全模块,你现在就可以使用它们。 比如,装载了 AppArmor 模板的 Docker 和 Red Hat自带 SELinux 策略的 Docker 
一、编辑网卡配置文件，桥接网卡
生成br0网卡配置文件
cd /etc/sysconfig/network-scripts/
cp ifcfg-eth0  ifcfg-eth0.bak
cp ifcfg-eth0  ifcfg-br0
eth0网卡配置 (不需要UUID，存在则两个网卡都有IP，网络不通)
 # cat ifcfg-eth0
TYPE=Ethernet
DEVICE=eth0
ONBOOT=yes
BRIDGE=br0
br0网卡配置 (注意：清除copy的mac地址)
1. DHCP获取方式
 # cat ifcfg-br0
TYPE=Bridge
DEVICE=br0
ONBOOT=yes
BOOTPROTO=dhcp
2. 静态配置方式
TYPE=Bridge
DEVICE=br0
ONBOOT=yes
IPADDR=192.168.8.140
GATEWAY=192.168.8.2
redhat或centos 7 版本新增加了nmtui配置基本网络连接 
可以文本用户界面创建网桥，也可以建多网卡bond，比较方便。
二、pipework 指定物理网段容器IP地址
不需要(/etc/sysconfig/docker 里面把容器默认配置绑定网卡br0)，docker会从头分配ip，没用。 
每次重启容器，ip是会变的。
下载pipework
cd /usr/src
wget  -O pipework-master.zip https://codeload.github.com/jpetazzo/pipework/zip/master
 # 若没有unzip命令，安装 yum install -y unzip zip
unzip pipework-master.zip 
cp -p pipework-master/pipework /usr/local/bin/
pipework固定物理网段容器IP地址
pipework br0 test01 192.168.8.10/24@192.168.8.2
 #       网桥  容器名     IP           GW
重启容器后IP需要再次指定
Docker 四种网络模式：
bridge方式(默认)、none方式、host方式、container复用方式
1、bridge方式： –net=”bridge”
容器与Host网络是连通的： 
eth0实际上是veth pair的一端，另一端（vethb689485）连在docker0网桥上 
通过Iptables实现容器内访问外部网络
2、none方式： –net=”none” 
这样创建出来的容器完全没有网络，将网络创建的责任完全交给用户。可以实现更加灵活复杂的网络。 
另外这种容器可以可以通过link容器实现通信。
3、host方式： –net=”host” 
容器和主机公用网络资源，使用宿主机的IP和端口 
这种方式是不安全的。如果在隔离良好的环境中（比如租户的虚拟机中）使用这种方式，问题不大。
4、container复用方式： –net=”container:name or id” 
新创建的容器和已经存在的一个容器共享一个IP网络资源
<!-- more -->

