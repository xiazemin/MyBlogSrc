---
title: dockershim
layout: post
category: docker
author: 夏泽民
---
当kubelet要创建一个容器时，需要以下几步：

Kubelet 通过 CRI 接口（gRPC）调用 dockershim，请求创建一个容器。CRI 即容器运行时接口（Container Runtime Interface），这一步中，Kubelet 可以视作一个简单的 CRI Client，而 dockershim 就是接收请求的 Server。目前 dockershim 的代码其实是内嵌在 Kubelet 中的，所以接收调用的凑巧就是 Kubelet 进程；
dockershim 收到请求后，转化成 Docker Daemon 能听懂的请求，发到 Docker Daemon 上请求创建一个容器。
Docker Daemon 早在 1.12 版本中就已经将针对容器的操作移到另一个守护进程——containerd 中了，因此 Docker Daemon 仍然不能帮我们创建容器，而是要请求 containerd 创建一个容器；
containerd 收到请求后，并不会自己直接去操作容器，而是创建一个叫做 containerd-shim 的进程，让 containerd-shim 去操作容器。这是因为容器进程需要一个父进程来做诸如收集状态，维持 stdin 等 fd 打开等工作。而假如这个父进程就是 containerd，那每次 containerd 挂掉或升级，整个宿主机上所有的容器都得退出了。而引入了 containerd-shim 就规避了这个问题（containerd 和 shim 并不是父子进程关系）；
我们知道创建容器需要做一些设置 namespaces 和 cgroups，挂载 root filesystem 等等操作，而这些事该怎么做已经有了公开的规范了，那就是 OCI（Open Container Initiative，开放容器标准）。它的一个参考实现叫做 runC。于是，containerd-shim 在这一步需要调用 runC 这个命令行工具，来启动容器；
runC 启动完容器后本身会直接退出，containerd-shim 则会成为容器进程的父进程，负责收集容器进程的状态，上报给 containerd，并在容器中 pid 为 1 的进程退出后接管容器中的子进程进行清理，确保不会出现僵尸进程。
<!-- more -->
https://blog.csdn.net/u011563903/article/details/90743853
https://www.cnblogs.com/charlieroro/articles/10998203.html
CRI 运行时有两个实现方案：

containerd
containerd 是 Docker 的一部分，提供的 CRI 都是由 Docker 提供的。
CRI-O：
CRI-O 在本质上属于纯 CRI 运行时，因此不包含除 CRI 之外的任何其他内容。
四、OCI 是啥？
当我们谈论「容器运行时」时，请注意我们到底是在谈论哪种类型的运行时，这里运行时分为两种：

CRI 运行时
OCI 运行时
OCI（Open Container Initiative），可以看做是「容器运行时」的一个标准，Ta 使用 Linux 内核系统调用（例如：cgroups 与命名空间）生成容器，按此标准实现的「容器运行时」有 runC 和 gVisor。
https://www.cnblogs.com/Gdavid/p/14818384.html
https://zhuanlan.zhihu.com/p/381235934
https://www.jianshu.com/p/75572613fd94
维护 dockershim 已经成为 Kubernetes 维护者肩头一个沉重的负担。 创建 CRI 标准就是为了减轻这个负担，同时也可以增加不同容器运行时之间平滑的互操作性。 但反观 Docker 却至今也没有实现 CRI，所以麻烦就来了。

Dockershim 向来都是一个临时解决方案（因此得名：shim）。 你可以进一步阅读 移除 Kubernetes 增强方案 Dockershim 以了解相关的社区讨论和计划。

此外，与 dockershim 不兼容的一些特性，例如：控制组（cgoups）v2 和用户名字空间（user namespace），已经在新的 CRI 运行时中被实现。 移除对 dockershim 的支持将加速这些领域的发展。


https://kubernetes.io/zh/blog/2020/12/02/dockershim-faq/