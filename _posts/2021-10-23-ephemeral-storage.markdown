---
title: ephemeral-storage
layout: post
category: k8s
author: 夏泽民
---
Kubernetes在1.8的版本中引入了一种类似于CPU，内存的新的资源模式：ephemeral-storage，并且在1.10的版本在kubelet中默认打开了这个特性。

Alpha release target (x.y): 1.7/1.8
Beta release target (x.y): 1.10
Stable release target (x.y): 1.11
ephemeral-storage是为了管理和调度Kubernetes中运行的应用的短暂存储。

在每个Kubernetes的节点上，kubelet的根目录(默认是/var/lib/kubelet)和日志目录(/var/log)保存在节点的主分区上，这个分区同时也会被Pod的EmptyDir类型的volume、容器日志、镜像的层、容器的可写层所占用。ephemeral-storage便是对这块主分区进行管理，通过应用定义的需求(requests)和约束(limits)来调度和管理节点上的应用对主分区的消耗。
<!-- more -->
ephemeral-storage的eviction逻辑
在节点上的kubelet启动的时候，kubelet会统计当前节点的主分区的可分配的磁盘资源，或者你可以覆盖节点上kubelet的配置来自定义可分配的资源。在创建Pod时会根据存储需求调度到满足存储的节点，在Pod使用超过限制的存储时会对其做驱逐的处理来保证不会耗尽节点上的磁盘空间。

如果运行时指定了别的独立的分区，比如修改了docker的镜像层和容器可写层的存储位置(默认是/var/lib/docker)所在的分区，将不再将其计入ephemeral-storage的消耗。

EmptyDir 的使用量超过了他的 SizeLimit，那么这个 pod 将会被驱逐
Container 的使用量（log，如果没有 overlay 分区，则包括 imagefs）超过了他的 limit，则这个 pod 会被驱逐
Pod 对本地临时存储总的使用量（所有 emptydir 和 container）超过了 pod 中所有container 的 limit 之和，则 pod 被驱逐

https://www.sunshanpeng.com/2018/11/28/%E4%BD%BF%E7%94%A8ephemeral-storage%E7%AE%A1%E7%90%86%E5%AE%B9%E5%99%A8%E7%9A%84%E4%B8%B4%E6%97%B6%E5%AD%98%E5%82%A8/

每次探测都将获得以下三种结果之一：

Success（成功）：容器通过了诊断。
Failure（失败）：容器未通过诊断。
Unknown（未知）：诊断失败，因此不会采取任何行动。

https://kubernetes.io/zh/docs/concepts/workloads/pods/pod-lifecycle/