---
title: boot2docker
layout: post
category: web
author: 夏泽民
---
因为有了 Docker 这个东西。它依赖于 LXC(Linux Container)，能从网络上获得配置好的 Linux 镜像，非常容易在隔离的系统中运行自己的应用。也因为它的底层核心是个 LXC，所以在 Mac OS X 下需要在 VirtualBox 中跑一个精小的 LXC(这里是一个 Tiny Core Linux，完全在内存中运行，个头只约 24MB，启动时间小于 5 秒的 boot2docker) 虚拟机，构建在 VirtualBox 中。以后的通信过程就是 docker--> boot2docker --> container，端口或磁盘映射也是遵照这一关系。
<!-- more -->
Docker 安装过程

1. 安装 VirtualBox, 不多讲, 因要在它当中创建一个 boot2docker-vm 虚拟机

2. 安装 boot2docker

brew install boot2docker

你也可以手工安装

curl https://raw.github.com/steeve/boot2docker/master/boot2docker > boot2docker; chmod +x boot2docker; sudo mv boot2docker /usr/local/bin

3. 安装 Docker

brew install docker

也可手工安装

curl -o docker http://get.docker.io/builds/Darwin/x86_64/docker-latest; chmod +x docker; sudo cp docker /usr/local/bin

4. 配置 Docker 客户端

export DOCKER_HOST=tcp://127.0.0.1:4243

把它写到 ~/.bash_profile 中，如果你是用的 bash 的话。我工作在 fish 下，所以在 ~/.config/fish/config.fish 中加了 set -x DOCKER_HOST tcp://127.0.0.1:4243

5. boot2docker 初始化与启动

boot2docker init

完成后就能在 VirtualBox 中看到一个叫做 boot2docker-vm的虚拟机，以后只需用 boot2docker 命令来控制这个虚拟机的行为，启动，停止等。

boot2docker up

启动，boot2docker-vm虚拟机，我们能在 VirtualBox 中看到该虚拟机变成 Running 状态

直接执行 boot2docker 可以看到可用的参数

Usage /usr/local/bin/boot2docker {init|start|up|save|pause|stop|restart|status|info|delete|ssh|download}

6. 启动 Docker 守护进程

sudo docker -d

这时可执行

boot2docker ssh，输入密码  tcuser 进到该虚拟机的控制台下，如果要用户名的话请输入docker

Mac 启动了 4243 端口，在 boot2docker 虚拟机中也有 4243 端口，并在 /var/run/docker.sock 上监听。借此回顾下 docker 的通信过程，dock 命令是与 Docker daemon 在  Mac 上开启的  4243 端口通信，该端口映射到 boot2docker 的  4243 端口上，进而通过 /var/run/docker.sock 与其中的容器进行通信。

所以在执行  docker version 时如果没有启动 Docker daemon 会提示

2014/05/16 06:52:48 Cannot connect to the Docker daemon. Is 'docker -d' running on this host?

如果没有启动 boot2docker 会得到提示

Get http:///var/run/docker.sock/v1.11/version: dial unix /var/run/docker.sock: no such file or directory

boot2docker, docker 都准备就绪了, 现在开始进入 dock  的操作了

端口映射
比如我们现在要做的映射关系是 Mac OS X(50080) --> boot2docker(40080) --> container(80)：
可以有两种办法

1)boot2docker ssh -L 50080:localhost:40080  #这条命令可以在  boot2docker-vm  运行时执行，建立多个不同的映射就是执行多次

docker run -i -t -p 40080:80 learn/tutorial
root@c79b5070a972:/# apachectl start

然后在 Mac 的浏览器中打开 http://localhost:50080

2)VBoxManage modifyvm "boot2docker-vm" --natpf1 "tcp-port_50080:80,tcp,,50080,,40080"

docker run -i -t -p 40080:80 learn/tutorial
root@c79b5070a972:/# apachectl start

这是直接修改了  boot2docker-vm 的配置，可以在 VirtualBox 中看到这条配置，配置 nat 命令见 http://www.virtualbox.org/manual/ch06.html#natforward. 也能建立许多的端口映射

boot2docker ssh -L 50080:localhost:40080  #这条命令可以在  boot2docker-vm  运行时执行，建立多个不同的映射就是执行多次，映射本机的50080到vm的40080
docker run -i -t -p 40080:80 learn/tutorial bash # 映射vm的40080到learn/tutorial容器的80端口。

docker容器中安装ssh服务，并启动。commit后退出。

boot2docker ssh -L 50080:localhost:40080 -L 50443:localhost:40443 -L50022:localhost:40022 ，映射多个ip到localhost，并启动vm。

docker run -i -t -p 40080:80 -p 40443:443 -p 40022:22 daimin/test bash ，映射端口方式启动docker容器，这个时候40022映射到了容器的22端口，也就是ssh端口，然后开启ssh服务。

在控制台输入 ssh -p 50022 root@localhost ，发现提示输入用户密码，然后成功登录。


在linux下我们可以在docker中新建容器，然后通过端口转发直接访问到容器。但是在mac下中间又通过了一层虚拟机，所以端口转发就需要在多做一点。

1 把本地端口和虚拟机对应起来：可以通过命令来做：
＃VBoxManage modifyvm "boot2docker-vm" --natpf1 "containerssh,tcp,,2222,,2222"
2，也可以通过手动添加：
选中虚拟机，设置，网络，网络地址转发（NAT），端口转发 



在 Mac 下使用 Docker 除了可用 boot2docker 作为 LXC，还有个替代品 VAGRANT 。
