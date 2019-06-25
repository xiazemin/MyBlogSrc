---
title: docker for mac
layout: post
category: web
author: 夏泽民
---
$ screen ~/Library/Containers/com.docker.docker/Data/com.docker.driver.amd64-linux/tty
<!-- more -->
用 Docker for Mac 已经很久了，用它跑本地开发环境可以说是非常方便，自打有了它，就再也没打开过 VMware Fusion。前段时间 Docker for Mac 又加入了 Kubernetes 支持，能在本地自动起一个单节点 Kubernetes 集群（就当它是集群），这就省去了手工搭建的步骤，简单好用。可以参考官方文档了解一下如何开启 Docker for Mac 的 Kubernetes。

但是 Docker for Mac 自诞生以来就一直有一个问题，那就是在宿主机上看不到 docker0，无法访问容器所在的网络，也就是说不能 ping 通 Docker 给 Container 所分配的 IP 地址。关于这个问题，官方文档中有描述，Known limitations, use cases, and workarounds。对于 docker run 启动的 Container 来说，通常会通过 -p 参数映射相应的服务端口，一般不会遇到要直接访问容器 IP 的情况。但对于 Kubernetes 来讲，很多时候都想要直接访问 Pod IP 或者 Service IP 进行调试，在 Mac 上却没办法实现。

Docker for Mac 基本原理
要解决这个问题，得先搞清楚 Docker for Mac的原理。我们都知道 Docker 是利用 Linux 的 Namespace 和 Cgroups 来实现资源的隔离和限制，容器共享宿主机内核，所以 Mac 本身没法运行 Docker 容器。不过不支持不要紧，我们可以跑虚拟机，最早还没有 Docker for Mac 的时候，就是通过 docker-machine 在 Virtual Box 或者 VMWare 直接起一个 Linux 的虚拟机，然后在主机上用 Docker Client 操作虚拟机里的 Docker Server。

Docker for Mac Architecture

Docker for Mac 也是在本地跑了一个虚拟机来运行 Docker，不过 Hypervisor 采用的是 xhyve，而 xhyve 又基于 Mac 自带的虚拟化方案 Hypervisor.framework，虚拟机里运行的发行版是 Docker 自己打包的 LinuxKit，之前用的发行版好像是 Alpine Linux。总而言之就是 Docker for Mac 跑的这个虚拟机非常轻量级，性能也会更好。

可以打开虚拟机的 tty 看看：
$ screen ~/Library/Containers/com.docker.docker/Data/com.docker.driver.amd64-linux/tty
Linux linuxkit-025000000001 4.9.87-linuxkit-aufs #1 SMP Fri Mar 16 18:16:33 UTC 2018 x86_64 Linux
linuxkit-025000000001:~#
问题解决思路
很简单，在 Docker for Mac 的虚拟机里跑一个 OpenVPN Server，然后从本地连过去。鉴于 Docker for Mac 在重启的时候不会保留虚拟机里的改动，所以这个 OpenVPN Server 必须要跑在容器里，并且网络模式需要设置为 host，这样才可以访问到所有的 Docker 网络。

流程如下：
Mac <-> Tunnelblick <-> socat/service <-> OpenVPN Server <-> Containers
Development Toolkit for Kubernetes on Docker for Mac
整理了些配置和脚本，在 Docker for Mac 开了 Kubernetes 功能后会比较有用，其中包括上面提到的 OpenVPN 服务。

GitHub 链接 docker-for-mac-kubernetes-devkit，使用方法在 README.md 里也有详细介绍。

访问 Pod/Docker 网络
提供了两种方式：

如果开了 Kubernetes，可以用 helm 安装。
如果没开，可以用 docker-compose 起 OpenVPN 服务。
准备工作
首先安装 Mac 的 OpenVPN 的客户端 Tunnelblick。

然后将代码 clone 下来，并进入 docker-for-mac-openvpn 目录。
$ git clone git@github.com:pengsrc/docker-for-mac-kubernetes-devkit.git
$ cd docker-for-mac-kubernetes-devkit/docker-for-mac-openvpn
用 Kubernetes 运行
装 helm。
在 local/values.yaml 创建一个配置文件，用于指定一些本地目录。Docker for Mac 的 File Sharing 配置里必须包含这些目录，否则 OpenVPN 是无法启动的。
dirPaths:
  # 项目目录
  data: /tmp/docker-for-mac-kubernetes-devkit/docker-for-mac-openvpn
  # 生成的 OpenVPN Client 的配置文件会放到这里
  local: /tmp/docker-for-mac-kubernetes-devkit/docker-for-mac-openvpn/local
  # 生成的 OpenVPN Server 的配置文件会放到这里
  configs: /tmp/docker-for-mac-kubernetes-devkit/docker-for-mac-openvpn/local/configs
运行 OpenVPN Server。注：-n 参数用于指定要安装到的 namespace。
1
$ helm install -n docker-for-mac -f local/values.yaml .
用 Docker Compose 运行
用 docker-compose 起服务更简单些。
$ # 启动
$ docker-compose up -d
$ # 看看日志
$ docker-compose logs -f
配置客户端
起 Server 的时候会生成 OpenVPN Client 的配置文件，放在 ./local/docker-for-mac.ovpn。最后，在配置文件末尾加上要访问的网段。

例如：
route 172.16.0.0 255.255.0.0
route 10.96.0.1 255.240.0.0
测试
跑个 Nginx 看下效果。
$ # 启动 Nginx
$ docker run --rm -it nginx

$ # 找到容器 IP
$ docker inspect `docker ps | grep nginx | awk '{print $1}'` | grep '"IPAddress"'
"IPAddress": "172.16.0.11",

$ # 访问下
$ curl http://172.16.0.11
<!DOCTYPE html>
<html>
<head>
<title>Welcome to nginx!</title>
...
Nginx Ingress Controller
这一段算是跑题了，docker-for-mac-kubernetes-devkit 还提供了在 Docker for Mac 上部署 Nginx Ingress Controller 的 Kubernetes Objects，批量 Apply 一下就可以使用。
$ kubectl apply -f ingress-nginx/namespaces
namespace "ingress-nginx" created

$ kubectl apply -f ingress-nginx/configmaps
configmap "nginx-configuration" created
configmap "tcp-services" created
configmap "udp-services" created

$ kubectl apply -f ingress-nginx/deployments
deployment "default-http-backend" created
deployment "nginx-ingress-controller" created

$ kubectl apply -f ingress-nginx/services
service "default-http-backend" created
service "nginx-ssl" created
service "nginx" created

$ curl http://127.0.0.1
default backend - 404

