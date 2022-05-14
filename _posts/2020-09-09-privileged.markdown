---
title: privileged 在docker内部运行docker
layout: post
category: docker
author: 夏泽民
---
https://github.com/lizrice/containers-from-scratch/blob/master/main.go

直接dssh 运行
[root@25e1d9737615 containers-from-scratch]$ /home/go/bin/go run main.go run /bin/bash
Running [/bin/bash]
panic: fork/exec /proc/self/exe: operation not permitted

goroutine 1 [running]:
main.must(0x52f220, 0xc42000e480)
	/home/xiaoju/containers-from-scratch/main.go:74 +0x4a
main.run()

原因：用自己的容器 docker exec -it --privileged xxx bash

公司的容器 用的是docker的默认root权限，阉割版的。

https://github.com/google/gvisor/issues/144
<!-- more -->
Docker exec 命令用来在宿主机直接在容器内部运行一个命令，不需要进入到容器。帮助信息如下：

[root@localhost ~]# docker exec --help
 
Usage:  docker exec [OPTIONS] CONTAINER COMMAND [ARG...]
 
Run a command in a running container
 
Options:
  -d, --detach               Detached mode: run command in the background
      --detach-keys string   Override the key sequence for detaching a container
  -e, --env list             Set environment variables
  -i, --interactive          Keep STDIN open even if not attached
      --privileged           Give extended privileges to the command
  -t, --tty                  Allocate a pseudo-TTY
  -u, --user string          Username or UID (format: <name|uid>[:<group|gid>])
  -w, --workdir string       Working directory inside the container
其中：

--detach-keys string 覆盖用于分离容器的键序列

-e, --env list 设置环境变量

-i, --interactive 保持STDIN打开，即使没有连接

--privileged 为该命令授予扩展特权

-t, --tty 分配一个pseudo-TTY

-u, --user string 用户名或UID (格式: <name|uid>[:<group|gid>])

-w, --workdir string 容器内的工作目录


https://github.com/kubernetes-sigs/kind/issues/1458

// 特权模式创建容器
docker run -t -i --privileged -v /usr/java:/mnt --name ContainerName ImageId /usr/sbin/init
// 采用host网络模式创建容器
docker run -t -i -d --privileged -v /usr/java:/mnt --net=host --name Container ImageId /usr/sbin/init


privileged参数
 
$ docker help run 
...
--privileged=false         Give extended privileges to this container
...
大约在0.6版，privileged被引入docker。
使用该参数，container内的root拥有真正的root权限。
否则，container内的root只是外部的一个普通用户权限。
privileged启动的容器，可以看到很多host上的设备，并且可以执行mount。
甚至允许你在docker容器中启动docker容器。

未设置privileged启动的容器：

[root@localhost ~]# docker run -t -i centos:latest bash
[root@65acccbba42f /]# ls /dev
console  fd  full  fuse  kcore  null  ptmx  pts  random  shm  stderr  stdin  stdout  tty  urandom  zero
[root@65acccbba42f /]# mkdir /home/test/
[root@65acccbba42f /]# mkdir /home/test2/
[root@65acccbba42f /]# mount -o bind /home/test  /home/test2
mount: permission denied
设置privileged启动的容器：
 
[root@localhost ~]# docker run -t -i --privileged centos:latest bash
[root@c39330902b45 /]# ls /dev/
autofs           dm-1  hidraw0       loop1               null    ptp3    sg0  shm       tty10  tty19  tty27  tty35  tty43  tty51  tty6   ttyS1    usbmon3  vcs5   vfio
bsg              dm-2  hidraw1       loop2               nvram   pts     sg1  snapshot  tty11  tty2   tty28  tty36  tty44  tty52  tty60  ttyS2    usbmon4  vcs6   vga_arbiter
btrfs-control    dm-3  hpet          loop3               oldmem  random  sg2  snd       tty12  tty20  tty29  tty37  tty45  tty53  tty61  ttyS3    usbmon5  vcsa   vhost-net
bus              dm-4  input         mapper              port    raw     sg3  stderr    tty13  tty21  tty3   tty38  tty46  tty54  tty62  uhid     usbmon6  vcsa1  watchdog
console          dm-5  kcore         mcelog              ppp     rtc0    sg4  stdin     tty14  tty22  tty30  tty39  tty47  tty55  tty63  uinput   vcs      vcsa2  watchdog0
cpu              dm-6  kmsg          mem                 ptmx    sda     sg5  stdout    tty15  tty23  tty31  tty4   tty48  tty56  tty7   urandom  vcs1     vcsa3  zero
cpu_dma_latency  fd    kvm           net                 ptp0    sda1    sg6  tty       tty16  tty24  tty32  tty40  tty49  tty57  tty8   usbmon0  vcs2     vcsa4
crash            full  loop-control  network_latency     ptp1    sda2    sg7  tty0      tty17  tty25  tty33  tty41  tty5   tty58  tty9   usbmon1  vcs3     vcsa5
dm-0             fuse  loop0         network_throughput  ptp2    sda3    sg8  tty1      tty18  tty26  tty34  tty42  tty50  tty59  ttyS0  usbmon2  vcs4     vcsa6
[root@c39330902b45 /]# mkdir /home/test/
[root@c39330902b45 /]# mkdir /home/test2/


[root@c39330902b45 /]# mount -o bind /home/test  /home/test2

--privilege是为了修改一些如/etc/sys里面的文件的时候才需要的，但正常运行一个容器，不建议开放这个权限。


privileged参数
$ docker help run 
...
--privileged=false         Give extended privileges to this container
...
大约在0.6版，privileged被引入docker。使用该参数，container内的root拥有真正的root权限。否则，container内的root只是外部的一个普通用户权限。privileged启动的容器，可以看到很多host上的设备，并且可以执行mount。甚至允许你在docker容器中启动docker容器。

未设置privileged启动的容器：

[root@xiexianbin_cn ~]# docker run -t -i centos:latest bash
[root@65acccbba42f /]# ls /dev
console  fd  full  fuse  kcore  null  ptmx  pts  random  shm  stderr  stdin  stdout  tty  urandom  zero
[root@65acccbba42f /]# mkdir /home/test/
[root@65acccbba42f /]# mkdir /home/test2/
[root@65acccbba42f /]# mount -o bind /home/test  /home/test2
mount: permission denied

设置privileged启动的容器：

[root@xiexianbin_cn ~]# docker run -t -i --privileged centos:latest bash
[root@c39330902b45 /]# ls /dev/
autofs           dm-1  hidraw0       loop1               null    ptp3    sg0  shm       tty10  tty19  tty27  tty35  tty43  tty51  tty6   ttyS1    usbmon3  vcs5   vfio
bsg              dm-2  hidraw1       loop2               nvram   pts     sg1  snapshot  tty11  tty2   tty28  tty36  tty44  tty52  tty60  ttyS2    usbmon4  vcs6   vga_arbiter
btrfs-control    dm-3  hpet          loop3               oldmem  random  sg2  snd       tty12  tty20  tty29  tty37  tty45  tty53  tty61  ttyS3    usbmon5  vcsa   vhost-net
bus              dm-4  input         mapper              port    raw     sg3  stderr    tty13  tty21  tty3   tty38  tty46  tty54  tty62  uhid     usbmon6  vcsa1  watchdog
console          dm-5  kcore         mcelog              ppp     rtc0    sg4  stdin     tty14  tty22  tty30  tty39  tty47  tty55  tty63  uinput   vcs      vcsa2  watchdog0
cpu              dm-6  kmsg          mem                 ptmx    sda     sg5  stdout    tty15  tty23  tty31  tty4   tty48  tty56  tty7   urandom  vcs1     vcsa3  zero
cpu_dma_latency  fd    kvm           net                 ptp0    sda1    sg6  tty       tty16  tty24  tty32  tty40  tty49  tty57  tty8   usbmon0  vcs2     vcsa4
crash            full  loop-control  network_latency     ptp1    sda2    sg7  tty0      tty17  tty25  tty33  tty41  tty5   tty58  tty9   usbmon1  vcs3     vcsa5
dm-0             fuse  loop0         network_throughput  ptp2    sda3    sg8  tty1      tty18  tty26  tty34  tty42  tty50  tty59  ttyS0  usbmon2  vcs4     vcsa6
[root@c39330902b45 /]# mkdir /home/test/
[root@c39330902b45 /]# mkdir /home/test2/
[root@c39330902b45 /]# mount -o bind /home/test  /home/test2

（十二）docker --privileged
1. privileged参数作用
--privileged                     Give extended privileges to this container
大约在0.6版，privileged被引入docker。
使用该参数，container内的root拥有真正的root权限。
否则，container内的root只是外部的一个普通用户权限。
privileged启动的容器，可以看到很多host上的设备，并且可以执行mount。
甚至允许你在docker容器中启动docker容器。

