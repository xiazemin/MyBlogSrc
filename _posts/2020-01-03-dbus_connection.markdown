---
title: Failed to get D-Bus connection
layout: post
category: linux
author: 夏泽民
---
I'm trying to install ambari 2.6 on a docker centos7 image but in the the ambari setup step and exactly while intializing the postgresql db I receive this error:

Failed to get D-Bus connection: Operation not permitted

I've got this error every time I try to run a serice on my docker image.

I tried every solution in the net but nothing worked yet.
<!-- more -->
Use this command

docker run -d -it --privileged ContainerId /usr/sbin/init

And access root in container

systemctl start httpd.service

在centos7的容器里面出现了一个BUG，就是serveice启动服务的时候出现报错，不能用service启动服务。
[root@e13c3d3802d0 /]# service httpd start
Redirecting to /bin/systemctl start  httpd.service
Failed to get D-Bus connection: Operation not permitted

首先恭喜你使用centos7镜像，然后就是不幸告诉你这个问题是个BUG 将在centos7.2解决。
  Currently, systemd in CentOS 7 has been removed and replaced with a fakesystemd package for dependency resolution. This is due to systemd requiring the CAP_SYS_ADMIN capability, as well as being able to read the host's cgroups. If you wish to replace the fakesystemd package and use systemd normally, please follow the steps below.
  
我查了好多地方都说是个BUG不能解决
有的说创建的时候加上 --privileged选项

我试了这些然而并没有任何的卵用

最后实在是没办法就 rpm -ql 软件包 查看安装的时候有哪些命令在PATH下，用这些命令去启动，这个是一种解决的方法
例如apache的启动就是用命令 httpd


这几天研究了个解决的办法比较靠谱，亲身实测好使：
systemctl start http.service
Failed to get D-Bus connection: No connection to service manager.

   这个的原因是因为dbus-daemon没能启动。其实systemctl并不是不可以使用。将你的CMD或者entrypoint设置为/usr/sbin/init即可。会自动将dbus等服务启动起来。
   然后就可以使用systemctl了。命令如下：
   docker run --privileged  -ti -e "container=docker"  -v /sys/fs/cgroup:/sys/fs/cgroup  centos  /usr/sbin/init
   
   因为centos7镜像的原因，默认是无法开启systemctl的，需要特权启动，并挂载cgroup，启动点设置为init，启动命令如下：

docker run -d -it –network=mynet –ip 172.18.0.8 –restart=always -h test –name test -v /sys/fs/cgroup:/sys/fs/cgroup:ro –privileged=true centos /usr/sbin/init

docker centos 使用 systemctl Failed to get D-Bus connection: Operation not permitted
转载倾-尽 发布于2017-12-22 15:35:01 阅读数 11766  收藏
展开
我们知道，Docker运行一个容器起来的时候，只是为你提供特定的文件系统层和进程隔离，它给你一个VM的感觉却并不是VM，所以你可能偶尔会想要像在物理机那样使用systemctl start|status|stop来管理服务进程，然后你通常会看到
Failed to get D-Bus connection: Operation not permitted
这个错误。
原因很简单：

你需要启动systemd进程
你需要特权
所以你如果想要一个可以使用Systemd的容器，你可以尝试这样启动容器:

1.解决办法一：给权限
cat /etc/redhat-release 
//CentOS Linux release 7.2.1511 (Core) 
docker run -tdi --privileged centos init

在容器中，你可以使用systemd管理服务进程了:

yum install -y vsftpd
systemctl start vsftpd
systemctl status vsftpd

2.解决办法二：init.d
除此之外你也可以通过其他启动软件的方式: init.d来达到 systemctl 的效果：
比如要启动 mysql 服务：/etc/init.d/mysql start


Docker的设计理念是在容器里面不运行后台服务，容器本身就是宿主机上的一个独立的主进程，也可以间接的理解为就是容器里运行服务的应用进程。一个容器的生命周期是围绕这个主进程存在的，所以正确的使用容器方法是将里面的服务运行在前台。

再说到systemd，这个套件已经成为主流Linux发行版（比如CentOS7、Ubuntu14+）默认的服务管理，取代了传统的SystemV风格服务管理。systemd维护系统服务程序，它需要特权去会访问Linux内核。而容器并不是一个完整的操作系统，只有一个文件系统，而且默认启动只是普通用户这样的权限访问Linux内核，也就是没有特权，所以自然就用不了！

因此，请遵守容器设计原则，一个容器里运行一个前台服务！



我就想这样运行，难道解决不了吗？

答：可以，以特权模式运行容器。



创建容器：

# docker run -d -name centos7 --privileged=true centos:7 /usr/sbin/init

进入容器：

# docker exec -it centos7 /bin/bash

这样可以使用systemctl启动服务了。