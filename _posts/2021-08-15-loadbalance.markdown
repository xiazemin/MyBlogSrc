---
title: k8s loadbalance
layout: post
category: k8s
author: 夏泽民
---
K8S不支持长连接的负载均衡，所以负载可能不是很均衡。如果你在使用HTTP/2，gRPC, RSockets, AMQP 或者任何长连接场景，你需要考虑客户端负载均衡。

Kubernetes Services不存在，没有进程监听服务的IP地址和端口。

您可以通过访问Kubernetes集群中的任何节点并执行netstat -ntlp来检查是否存在这种情况。

甚至在任何地方都找不到IP地址,Services的IP地址由控制器管理器中的控制平面分配，并存储在数据库etcd中。然后，另一个组件将使用相同的IP地址：kube-proxy。

Kube-proxy读取所有Services的IP地址列表，并在每个节点中写入一组iptables规则。这些规则的意思是：“如果看到此Services IP地址，则改写请求并选择Pod之一作为目的地”。Services IP地址仅用作占位符-这就是为什么没有进程监听IP地址或端口的原因。


<!-- more -->
https://blog.icorer.com/index.php/archives/507/

https://blog.51cto.com/nxlhero/2713860

内核模式的转发就是iptables和ipvs。

3.1.3.2 iptables和ipvs
这俩其实都是依赖的一个共同的Linux内核模块：Netfilter。Netfilter是Linux 2.4.x引入的一个子系统，它作为一个通用的、抽象的框架，提供一整套的hook函数的管理机制，使得诸如数据包过滤、网络地址转换(NAT)和基于协议类型的连接跟踪成为了可能。
Netfilter的架构就是在整个网络流程的若干位置放置了一些检测点（HOOK），而在每个检测点上登记了一些处理函数进行处理。

https://kubernetes.io/zh/docs/concepts/services-networking/dns-pod-service/
