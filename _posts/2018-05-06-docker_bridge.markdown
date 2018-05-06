---
title: docker_bridge 容器桥接到物理网络
layout: post
category: docker
author: 夏泽民
---
在5台物理主机上的虚拟机中都装了Docker，每台中都有3个容器，现在要解决容器跨主机通信，五种方案：
一、利用OpenVSwitch
二、利用Weave
三、Docker在1.9之后支持的Overlay network（这个好像是官方的做法）
四、将多个物理机的容器组到一个物理网络来
        1.创建自己的网桥br0
        2.将docker默认网桥绑定到br0
五、修改主机docker默认的虚拟网段，然后在各自主机上分别把对方的docker网段加入到路由表中，配合iptables即可实现docker容器跨主机通信
<!-- more -->
跨主机容器通信大致有以下几种方式：
1. NAT方式
NAT就是传统的docker网络，利用宿主机的IP和Iptables来达到容器，主机之间的通信。容器对外IP都是宿主机的IP。只要宿主机之间三层可达，容器之间就可以通信。
2. Tunnel（overlay）方式
VPN，ipip，VXLAN等都是tunnel技术，简单讲就是在容器的数据包间封装一层或多层其他的数据协议头，达到连通的效果。
Docker的Libnetwork就支持vxlan的overlay方式，weave也支持UDP和Vxlan的overlay模式，flannel，calico等都支持overlay模式。
•      Overlay一般通过vSwitch来实现，比如OVS
•      一般需要一个全局的KV
store（sdn controller、etcd、consul）来保存控制信息。
3. Routing方式
路由方案主要是通过路由设置的方式让容器对容器，容器对宿主机之间相通信。例如：calico的BGP路由方案
4. “大二层”方式
大二层就是将物理网卡和容器网络桥接到同一个linux bridge。容器网络可以直接接入该linux bridge，也可以将容器网络接入容器网桥docker0，再把docker0桥接到linux bridge上。目的是使得容器网络和宿主机网络在同一个二层网络。
docker 默认的桥接网卡是docker0,它只会在本机桥接所有的容器网卡，举例来说容器的虚拟网卡在主机上看一般叫做veth***  而docker只是把所有这些网卡桥接在一起,这样就可以把这个网络看成是一个私有的网络，通过nat 连接外网，如果要让外网连接到容器中，就需要做端口映射，即-p参数
 当rocker启动时，会在主机上创建一个docker0的虚拟网卡。他随机挑选RFC1918私有网络中的一段地址给docker0。比如172.17.42.1/16,16位掩码的网段可以拥有65534个地址可以使用，这对主机和容器来说应该足够了。docker0 不是普通的网卡，他是桥接到其他网卡的虚拟网卡，容器使用它来和主机相互通信。当创建一个docker容器的时候，它就创建了一个对接口，当数据包发送到一个接口时，另外一个接口也可以收到相同的数据包，它们是绑在一起一对孪生接口。这对接口在容器中的那个的名字是eth0，主机上的接口会指定一个唯一的名字，比如vethAQI2QT这样的名字，这种接口名字不再主机的命名空间中。所有的veth*的接口都会桥接到docker0，这样docker就创建了在主机和所有容器之间一个虚拟共享网络。
如果在企业内部应用，或则做多个物理主机的集群，可能需要将多个物理主机的容器组到一个物理网络中来，那么就需要将这个网桥桥接到我们指定的网卡上。下面以ubuntu为例创建多个主机的容器联网。
第一步：创建自己的网桥
编辑/etc/network/interface文件
auro br0
iface bro inet static
address
netmask
gateway
bridge_ports em1
bridge_stp off
dns-nameservers 8.8.8.8 
将docker的默认网桥绑定到这个新建的br0上面，这样就将这台机器上容器绑定到em1这个网卡所对应的物理网络上了。ubuntu修改/etc/default/docker文件 添加最后一行内容
DOCKER_OPTS="-b=br0"
这改变默认的docker网卡绑定，你也可以创建多个网桥绑定到多个物理网卡上，在启动docker的时候 使用-b参数 绑定到多个不同的物理网络上。
重启docker服务后，再进入容器可以看到它已经绑定到你的物理网络上了，这样就直接把容器暴露到你的物理网络上了，多台物理主机的容器也可以相互联网了。需要注意的是，这样就需要自己来保证容器的网络安全了。
Docker发布了1.9版本，实现了原本实验性的Networking组件的支持。用户可以在Swarm中使用它或者将其作为Compose 工具。创建虚拟网络并将其连接到容器上，可实现多个主机上容器相互通信，并且实现不同的应用程序或者应用程序不同部分能够相互隔离。互联子系统设计为可插拔式，兼容VXLAN或者IPVLAN等技术。
Docker 容器实现跨主机通讯主要通过几种方式：自带overlay network插件，第三方插件如weave、ovs等，docker swarm（虽然也是通过key value service进行调用）。在这里主要介绍直接使用自带插件，以及通过docker swarm的两种实现方式。
直接使用自带插件实现容器跨主机访问
测试环境为ubuntu14.04。根据docker建议要求，将内核升级到3.19.如果内核版本过低，将会出现overlay network创建失败或加入失败等一系列问题。
内核升级
 #apt-get install linux-generic-lts-vivid
升级完成，重启主机。
 # uname -a
Linux ukub09 3.19.0-33-generic #38~14.04.1-Ubuntu SMP Fri Nov 6 18:17:28 UTC 2015 x86_64 x86_64 x86_64 GNU/Linux
启动key value service
可以选择使用zookeeper、etcd、consul等多种组件。本例使用额外一台虚拟机，启用和官方文档相一致的consul镜像。
docker $(docker-machine config consul2) run -d \
-p "8501:8500" \
-h "consul" \
progrium/consul -server -bootstrap
此时可以通过ip:8501对consul服务进行访问。
将key value service信息加入docker daemon
修改希望加入docker overlay network的所有主机上docker daemon服务。
#vi /etc/dafault/docker
DOCKER_OPTS='
-H tcp://0.0.0.0:2376
-H unix:///var/run/docker.sock
--cluster-store=consul://172.25.2.43:8501
--cluster-advertise=eth0:2376'
通过修改以下两个参数实现
--cluster-store= 参数指向docker daemon所使用key value service的地址
--cluster-advertise= 参数决定了所使用网卡以及docker daemon端口信息
重启docker服务。
#service docker restart
docker 网络的创建与操作
此时可以通过docker network create创建-d overlay属性的驱动。
#docker network create -d overlay over
列出当前网络
root@ukub10:~# docker network ls
NETWORK ID          NAME                DRIVER
b95efb2ab985        over                overlay             
0a4a123fc278        bridge              bridge              
2da9ecadf108        none                null                
5ca3275ec2a9        host                host        
在另一台已加入key value service指向的主机上同样可以看到这个网络
root@ukub09:~# docker network ls
NETWORK ID          NAME                DRIVER
b95efb2ab985        over                overlay             
b5f9713c9f2b        bridge              bridge              
16e1cdf30f97        none                null                
cf257529463c        host                host
docker网络的连接方式
可以在两台主机上分别创建容器，互相进行ping或ssh访问。官方文档例子创建一个nginx容器，并在另一个容器内抓取其状态进行验证。这里不再详述。
创建容器接入指定网络
#docker run -itd --name=testcontainer --net=over  nginx
将容器从指定网络中退出
#docker network disconnect over testcontainer 
重新连入指定网络
#docker network connect over testcontainer
连接与断开将会是实时的。在连接之后，即可通过容器名访问本网络中所有容器。
使用Swarm创建docker cluster实现跨主机管理和网络访问
根据官方文档，需要搭建起基于swarm的cluster，来实现Docker overlay network。为部署整套Docker cluster，需要安装docker machine/compose/swarm等组件。
使用第三方网络插件实现跨主机容器通讯也变得更容易了。
安装docker machine
docker machine用来进行docker的推送安装，可以在许多driver上安装部署docker：
amazonec2 azure digitalocean exoscale generic google none
openstack rackspace softlayer virtualbox vmwarevcloudair
vmwarevsphere
首先需要安装好docker engine，之后下载并将docker machine的库文件移动到bin下。
#curl -L https://github.com/docker/machine/releases/download/v0.5.0/docker-machine_linux-amd64.zip >machine.zip && \
unzip machine.zip && \
rm machine.zip && \
mv docker-machine* /usr/local/bin
验证安装
# docker-machine -v
docker-machine version 0.5.0 (04cfa58)
docker-machine ls
NAME   ACTIVE   DRIVER   STATE   URL   SWARM
官方文档的使用virtualbox进行各种角色主机的启用。本篇文档则在已有的虚拟机环境中进行。
使用docker-machine进行推送安装
首先需要拥有一台linux系统，并使docker-machine主机可以无密码访问到这台系统。
#ssh-copy-id -i /root/.ssh/id_rsa.pub root@172.25.2.43
之后即可通过docker-machine进行安装。
安装通过internet下载相关包，并通过docker-machine进行设置安装的。必须保证网络畅通，否则可能安装失败。
#docker-machine create -d generic --generic-ip-address=172.25.2.43 --generic-ssh-user=root test5
安装后，可以查看所有被docker-machine安装控制后的主机状态
# docker-machine ls
NAME        ACTIVE   DRIVER    STATE     URL                      SWARM
c0-master   -        generic   Running   tcp://172.25.2.44:2376   c0-master (master)
c0-n1       -        generic   Running   tcp://172.25.2.45:2376   c0-master
c0-n2       -        generic   Running   tcp://172.25.2.46:2376   c0-master
dm01        -        generic   Running   tcp://172.25.2.43:2376   
local       -        generic   Running   tcp://172.25.2.34:2376
可以查看相应主机的环境变量
#docker-machine env local 
把当前操作环境变更为指定主机
#eval "$(docker-machine env local)"
在docker-machine主机上操纵host
 #docker $(docker-machine config local) run -tid ubuntu /bin/bash
通过ssh连接到相应主机
 # docker-machine ssh dm01
Docker Compose安装
docker compose用于同时启动多个容器，构成同一组应用环境，在这里用于支持swarm的安装。
curl -L https://github.com/docker/comp ... ose-X 10X- uname -m > /usr/local/bin/docker-compose
chmod +x /usr/local/bin/docker-compose
 Docker Swarm集群安装
启用一台主机安装key value service，用于节点、主机发现等工作。
本例使用了consul的镜像，也可以使用如etcd、zookeeper等组件。
在dm01上启动一个consul服务：
docker $(docker-machine config dm01) run -d \
-p "8500:8500" \
-h "consul" \
progrium/consul -server -bootstrap
创建swarm集群可以通过key value服务，也可以使用swarm cluster ID字符串创建。
创建swarm manager
docker-machine create \
-d generic --generic-ip-address=172.25.2.44 --generic-ssh-user=root\
--swarm \
--swarm-master \
--swarm-discovery="consul://$(docker-machine ip dm01):8500" \
--engine-opt="cluster-store=consul://$(docker-machine ip dm01):8500" \
--engine-opt="cluster-advertise=eth0:2376" \
c0-master
创建swarm node
docker-machine create \
-d generic --generic-ip-address=172.25.2.45 --generic-ssh-user=root \
--swarm \
--swarm-discovery="consul://$(docker-machine ip dm01):8500" \
--engine-opt="cluster-store=consul://$(docker-machine ip dm01):8500" \
--engine-opt="cluster-advertise=eth0:2376" \
c0-n1
docker-machine create \
-d generic --generic-ip-address=172.25.2.46 --generic-ssh-user=root \
--swarm \
--swarm-discovery="consul://$(docker-machine ip dm01):8500" \
--engine-opt="cluster-store=consul://$(docker-machine ip dm01):8500" \
--engine-opt="cluster-advertise=eth0:2376" \
c0-n2
Overlay Network创建
设置环境变量
 #eval "$(docker-machine env --swarm c0-master)"
可以用docker info看到现在的集群状况
 #docker network create -d overlay myStack1
root@local:~# docker network ls
NETWORK ID          NAME                DRIVER
685ed9e9f701        myStack1            overlay             
4c88e7c52a7c        c0-n1/bridge        bridge              
7fd49a4f6a64        c0-n1/none          null                
8b2372139902        c0-n2/host          host                
d998a91dcfea        c0-n2/bridge        bridge              
afd7a1190fb4        c0-n2/none          null                
0f2b5dae67f7        c0-n1/host          host                
a45b8049f8e6        c0-master/bridge    bridge              
d13eed83f250        c0-master/none      null                
364a592ae9ae        c0-master/host      host  
测试
在c0-n1上启动一个nginx，并加入myStack1网络
 #docker run -itd --name=web --net=myStack1 --env="constraint:node==c0-n1" nginx
在c0-n2上启动一个shell，同样加入myStack1网络
 #docker run -ti --name=webtest --net=myStack1 --env="constraint:node==c0-n2" ubuntu /bin/bash
在shell中将可以通过容器名或ip获取到web容器的状态
{{{#apt-get install wget
wget -O- http://web }}}
docker在官方文档上主推以swarm的方式创建cluster，相对比较复杂，但提供了整体Cluster的解决方案，可以将容器进行整体管理、推送、使用。

如何使不同主机上的docker容器互相通信
docker启动时，会在宿主主机上创建一个名为docker0的虚拟网络接口，默认选择172.17.42.1/16，一个16位的子网掩码给容器提供了65534个IP地址。docker0只是一个在绑定到这上面的其他网卡间自动转发数据包的虚拟以太网桥，它可以使容器和主机相互通信，容器与容器间通信。
问题是，如何让位于不同主机上的docker容器可以通信？
最简单的思路，修改一台主机docker默认的虚拟网段，然后在各自主机上分别把对方的docker网段加入到路由表中，即可实现docker容器夸主机通信。
现有两台虚拟机
v1：192.168.124.51
v2：192.168.124.52
更改虚拟机docker0网段，修改为
v1：172.17.1.1/24
v2：172.17.2.1/24
 #v1
sudo ifconfig docker0 172.17.1.1 netmask 255.255.255.0
sudo service docker restart
 #v2
sudo ifconfig docker0 172.17.2.1 netmask 255.255.255.0
sudo service docker restart
然后在v1，v2上把对方的docker0网段加入到自己的路由表中
 #v1
sudo route add -net 172.17.2.0 netmask 255.255.255.0 gw 192.168.124.52
sudo iptables -t nat -F POSTROUTING
sudo iptables -t nat -A POSTROUTING -s 172.17.1.0/24 ! -d 172.17.0.0/16 -j MASQUERADE
 #v2
sudo route add -net 172.17.1.0  netmask 255.255.255.0  gw 192.168.124.51
sudo iptables -t nat -F POSTROUTING
sudo iptables -t nat -A POSTROUTING -s 172.17.2.0/24 ! -d 172.17.0.0/16 -j MASQUERADE
测试，v1，v2创建容器test1，test2
＃v1
docker run --rm --name test1 -i -t base:latest bin/bash
docker inspect --format '{{.NetworkSettings.IPAddress}}' test1
 #172.17.1.1
v2
docker run --rm --name test2 -i -t base:latest bin/bash
docker inspect --format '{{.NetworkSettings.IPAddress}}' test2
 #172.17.2.1
主机上可以ping通对方容器ip，至此也就ok了。