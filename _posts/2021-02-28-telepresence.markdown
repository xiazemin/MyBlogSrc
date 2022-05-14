---
title: telepresence 本地调试
layout: post
category: k8s
author: 夏泽民
---
https://github.com/telepresenceio/telepresence

https://kubernetes.io/zh/docs/tasks/debug-application-cluster/local-debugging/

Kubernetes 应用程序通常由多个独立的服务组成，每个服务都在自己的容器中运行。 在远端的 Kubernetes 集群上开发和调试这些服务可能很麻烦，需要 在运行的容器上打开 Shell， 然后在远端 Shell 中运行你所需的工具。

telepresence 是一种工具，用于在本地轻松开发和调试服务，同时将服务代理到远程 Kubernetes 集群。 使用 telepresence 可以为本地服务使用自定义工具（如调试器和 IDE）， 并提供对 Configmap、Secret 和远程集群上运行的服务的完全访问。
<!-- more -->
Kubernetes 集群安装完毕
配置好 kubectl 与集群交互
Telepresence 安装完毕
打开终端，不带参数运行 telepresence，以打开 telepresence Shell。 这个 Shell 在本地运行，使你可以完全访问本地文件系统。

telepresence Shell 的使用方式多种多样。 例如，在你的笔记本电脑上写一个 Shell 脚本，然后直接在 Shell 中实时运行它。 你也可以在远端 Shell 上执行此操作，但这样可能无法使用首选的代码编辑器，并且在容器终止时脚本将被删除。

开发和调试现有的服务
在 Kubernetes 上开发应用程序时，通常对单个服务进行编程或调试。 服务可能需要访问其他服务以进行测试和调试。 一种选择是使用连续部署流水线，但即使最快的部署流水线也会在程序或调试周期中引入延迟。

使用 --swap-deployment 选项将现有部署与 Telepresence 代理交换。 交换允许你在本地运行服务并能够连接到远端的 Kubernetes 集群。 远端集群中的服务现在就可以访问本地运行的实例。

要运行 telepresence 并带有 --swap-deployment 选项，请输入：

telepresence --swap-deployment $DEPLOYMENT_NAME

这里的 $DEPLOYMENT_NAME 是你现有的部署名称。

运行此命令将生成 Shell。在该 Shell 中，启动你的服务。 然后，你就可以在本地对源代码进行编辑、保存并能看到更改立即生效。 你还可以在调试器或任何其他本地开发工具中运行服务。

telepresence 和kubectl port-forward对比的好处是不用不用一个个配置，能够批量配置

https://cloud.google.com/community/tutorials/developing-services-with-k8s

Telepresence 是一个 CNCF 1 基金会下的项目。它的工作原理是在本地和 Kubernetes 集群中搭建一个透明的双向代理，这使得我们可以在本地用熟悉的 IDE 和调试工具来运行一个微服务，同时该服务还可以无缝的与 Kubernetes 集群中的其他服务进行交互，好像它就运行在这个集群中一样。

有了这些代理之后：

本地的服务就可以直接使用域名访问到远程集群中的其他服务
本地的服务直接访问到 Kubernetes 里的各种资源，包括环境变量、secrets、config map等
甚至集群中的服务还能直接访问到本地暴露出来的接口
安装
macOS:
brew cask install osxfuse  # required by sshfs to mount the pod's filesystem
brew install datawire/blackbird/telepresence
其他平台请参考：https://www.telepresence.io/reference/install

如果官方的安装包没有覆盖到你的平台，其实也可以从源代码安装，因为它本身就是用 Python3 写的，熟悉 Python 的朋友安装这个程序应该不难，我自己就在 CentOS 7 上安装成功了。

现在在本地用默认参数启动 Telepresence ，等它连接好集群：
$ telepresence

当运行 Telepresence 命令的时候，它创建了一个Deployment，这个Deployment的 Spec 是一个转发流量的代理容器，我们可以这样查看到它 kubectl get pod -l telepresence。
同时它还在本地创建了一个全局的 VPN，使得本地的所有程序都可以访问到集群中的服务。 Telepresence 其实还支持其他的网络代理模式（使用--method切换），vpn-tcp是默认的方式，其他的好像用处不大，inject-tcp甚至要在后续的版本中取消掉。
当本地的curl访问http://service-b:8000/时，对应的 DNS 查询和 HTTP 请求都被 VPN 路由到集群中刚刚创建的容器中处理。如果域名解析不了 (Could not resolve host)，可以试试加上 search 后缀：service-b.<NAMESPACE>.svc.cluster.local。

它的工作原理概述如下：

在集群中创建一个代理Deployment，并复制service-b的所有Label。
建立一个路由通道，将代理容器的所有流量转发到本地 8000 端口。
将service-b的 replicas 数设为0，这样 K8S Service 的 selector 就只能匹配到刚刚创建的代理容器上。
通过这样的方法，我们就有机会将集群中的请求转发到本地，然后在本地查看到具体的请求数据，调试逻辑，以及生成新的回复。

https://blog.betacat.io/post/develop-and-debug-k8s-microservices-locally-using-telepresence/

https://www.telepresence.io/



