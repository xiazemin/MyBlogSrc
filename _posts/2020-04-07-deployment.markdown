---
title: kubernetes 控制器 Deployment 核心机制
layout: post
category: k8s
author: 夏泽民
---
https://banzaicloud.com/blog/generating-go-code/
https://gocn.vip/topics/10174
Deployment 本质上其实只是一种部署策略，在了解其实现之前，先简单介绍一下部署系统里面常见的概念，Deployment 里面的各种参数和设计其实也都是围绕着这些展开的

1.1 ReplicaSet
Deployment 本身并不直接操作 Pod，每当其更新的时候通过构建 ReplicaSet 来进行版本更新，在更新的过程中通过 scale up(新的 RS) 和 scale down(旧的 RS) 来完成

1.2 部署状态
在 k8s 的官方文档中主要是介绍了 Deployment 的三种状态, 对应的 Condition 分别为 Available、Progressing、ReplicaFailure 三种状态, 并且每个状态下面又会很有导致对应状态切换的不同的 Reson，Reson 可能是运维过程中最需要关注的点
<!-- more -->
1.3 部署策略
部署策略是 Deployment 控制 ReplicaSet 更新的策略，通过对新旧 ReplicaSet 的扩缩容，再满足部署策略的情况下，将系统更新至最新的目标状态，Deployment 本身并没有太多可选的策略，默认只有两种 Recreate 和 RollingUpdate
image.png

在一些大版本产品更新的时候，新旧版本的数据库模型都不一致的情况下，通常会选择停服操作，此时可以选择 Recreate 即将所有老的副本都干掉，然后重新创建一批。当然默认情况下大部分业务还是 RollingUpdate 即滚动更新即可

1.4 可用性与高低水位
部署过程中最常被提到的可能就是可用性问题了，即在更新的过程中 (RollingUpdate 策略下) 需要保证系统中可用的 Pod 在一个指定的水位，保证对应服务的可用性
image.png

高低水位 (deployment 并没有这个词) 其实就是对应的上面的可用性来说的，Deployment 通过一些参数让我们可以自由控制在滚动更新的过程中，我们可以创建的 Pod 的最多数量 (高水位) 和可以删除的最多的 Pod(低水位)， 从而达到可用性保护的目标

部署的概念就介绍到这里， 接下来就一起看看 Deployment 中这些关键机制的具体实现

2. 核心实现
Deployment 的实现上相对复杂一点，但是从场景上又可以简单的分为：删除、暂停、回滚、扩缩容、更新几个大的场景

2.1 暂停部署
暂停部署是用于中断 Deployment 更新流程的一种方式，但由于 k8s 中是基于事件驱动的最终一致性的系统，这里的中断仅仅意味着 Deployment 层不会进行的进行后续的副本变更，而底层的 replicaSet 此时如果还没有达到目标的副本，则就需要继续更新， 同时在暂停的过程中如果发现并没有尝试进行回滚到指定版本的操作，这时候还会进行一些副本的清理工作，即只保留最近的指定数量的历史副本

2.2 回滚控制
回滚控制里面的信息跟其他参数有些不同，其主要是通过在 Annotations 中存储的 DeprecatedRollbackTo 来进行指定版本的回滚
image.png
回滚的实现本质上就是从指定的 Revisions 中获取对应的 replicaset 的 Pod 模板，去覆盖当前的 Deployment 的 Pod 模板，并且更新 Deployment 即可, 那如果对应的版本不存在怎么办，如果是这种情况，其实就需要你自己去寻找历史版本了，并且 k8s 会给新添加一个 RollbackRevisionNotFound 类型的事件提示你版本不存在 

2.3 扩缩容机制
扩缩容机制主要是指的 Deployment 的 scale 操作，在进行 Deployment 更新之前，会首先检查对应 Deployment 的副本的期望是否得到满足，只有期望的副本数得到满足,才会进行更新操作，所以在 k8s 中如果之前进行了扩缩容操作，则在该操作完成之前，是不会进行模板更新的

2.4 Recreate 策略
image.png

Recreate 部署策略在实现上通过两种机制保证之前的 Pod 一定被删除：所有活跃副本都为 0 和所有 Pod 都处于 (PodFailed 和 PodSucceeded) 两种状态下，然后才会创建新的副本,如果对应的副本完全就绪，还会进行清理历史副本

2.5 RollingUpdate 策略
RollingUpdate 策略可能是最复杂的部分之一了，里面很有多的参数控制，都作用于该策略，来一起看下

2.5.1 缩容过多的新版本
首先在更新的时候要做一致性检测，如果发现新版本的 ReplicaSet 比当前的 deployment 设定的副本数目多，则首先干掉这部分 Pod, 同时会根据当前的 Deployment 的副本数来设定当前的期望副本数 DesiredReplicasAnnotation, 并且根据 maxSurge 来计算当前最大的副本数量 MaxReplicasAnnotation, 同时在这个同步 ReplicaSet 的 minReadySeconds

2.5.2 扩容新副本
如果说新副本的数量不足，则就需要根据当前的 maxSurege 来设定，同时会再次计算当前的 RS 的所有 Pod,如果发现 Pod 数量过多即超过 Deployment 的 Replicas+maxSurge，则也不会进行操作

// Find the total number of pods
currentPodCount := GetReplicaCountForReplicaSets(allRSs)
// 最大pod数量
maxTotalPods := *(deployment.Spec.Replicas) + int32(maxSurge)
// 当前pod数量》总的运行的pod数量
if currentPodCount >= maxTotalPods {
    // Cannot scale up.
    return *(newRS.Spec.Replicas), nil
}
否则则就会进行计算允许 scale up 的数量

scaleUpCount := maxTotalPods - currentPodCount
scaleUpCount = int32(integer.IntMin(int(scaleUpCount), int(*(deployment.Spec.Replicas)-*(newRS.Spec.Replicas))))
2.5.3 缩容旧的副本
缩容计数器的算法计算主要是根据 Deployment 的 Replcas 和 maxUnavailable(通过 surge 和 maxUnavailable) 共同计算而来，最终的公式其实如下，有了缩容的数量，就可以更新旧 ReplicaSet 的数量了

minAvailable := *(deployment.Spec.Replicas) - maxUnavailable
// 新副本不可用数量
newRSUnavailablePodCount := *(newRS.Spec.Replicas) - newRS.Status.AvailableReplicas
// 最大缩容大小=所有pod统计-最小不可用-新副本不可用副本
maxScaledDown := allPodsCount - minAvailable - newRSUnavailablePodCount
整理如下

最小可用副本   = Deployment的副本-最大不可用副本
新副本不可用统计= 新副本数量-可用副本数量
最大缩容数量   = 全部副本Pod计数-最小可用副本-新副本不可用统计
至此我们知道了 Deployment 扩缩容的核心的副本计算实现，也知道了扩缩容的流程，那还缺什么呢？答案是状态

2.5.4 Available 状态
Deployment 的状态主要是由新旧副本以及当前集群中的 Pod 决定的，其计算公式如下， 则认为当前可用否则即为不可用

前可用的副本数量计数>=Deployment的副本数量-最大不可用副本计数
if availableReplicas >= *(deployment.Spec.Replicas)-deploymentutil.MaxUnavailable(*deployment) {
    minAvailability := deploymentutil.NewDeploymentCondition(apps.DeploymentAvailable, v1.ConditionTrue, deploymentutil.MinimumReplicasAvailable, "Deployment has minimum availability.")
    deploymentutil.SetDeploymentCondition(&status, *minAvailability)
} else {
    noMinAvailability := deploymentutil.NewDeploymentCondition(apps.DeploymentAvailable, v1.ConditionFalse, deploymentutil.MinimumReplicasUnavailable, "Deployment does not have minimum availability.")
    deploymentutil.SetDeploymentCondition(&status, *noMinAvailability)
}
2.5.5 Processing 状态
首先并不是所有的 Deployment 都有该状态，只有设置了 progressDeadlineSeconds 参数的才会有该状态，其主要实在 Deployment 未完成的时候，进行一些状态决策，从而避免一个 Deployment 无期限的运行，其关键状态有两个即运行中与超时决策， 其流程实现上分为两步 1）首先如果判断是正在运行中，就更新 LastTransitionTime

condition := util.NewDeploymentCondition(apps.DeploymentProgressing, v1.ConditionTrue, util.ReplicaSetUpdatedReason, msg)

condition.LastTransitionTime = currentCond.LastTransitionTime
2）超时检测则及时通过记录之前的转换时间，然后决策是否超时

from := condition.LastUpdateTime
now := nowFn()
delta := time.Duration(*deployment.Spec.ProgressDeadlineSeconds) * time.Second
timedOut := from.Add(delta).Before(now)

2.5.6 ReplicaFailure 状态
 该状态相对简单，检测当前的所有的副本，如果发现有副本失败，就取最新的一条失败的信息来填充 Condition

3. 小结
image.png

更新机制的核心实现可能就这些，代码实现上还是相对复杂的，主要是集中在为了保证伸缩和更新时为了保证可用性而做了大量的计算，还有很多的边界条件的处理
