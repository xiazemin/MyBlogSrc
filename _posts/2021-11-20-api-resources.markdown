---
title: api-resources
layout: post
category: k8s
author: 夏泽民
---
Binding: 已弃用。用于记录一个object和另一个object的绑定关系。实际上主要用于将pod和node关系，所以在1.7版本后已经改为在pods.bindings中记录了。
ComponentStatus: 是一个全局的list（即不受命名空间影响），记录了k8s中所有的组件的的相关信息，比如创建时间，现在状态等。
Configmap: 是一种用于记录pod本身或其内部配置信息的API资源，可以认为是通过API形式存储的配置文件。
Endpoints: 用于记录每个service的pod的真实物理ip和port的对应关系，包括service是TCP还是UDP等。
Event: 用于记录集群中的事件，可以认为类似于日志里的一条记录。
LimitRange: 用于记录各个命名空间中的pod或container对每种资源的使用限制，一般被包含在pod的定义中。
Namespace: 是一个全局的list，保存集群中所有的命名空间。
Node: 是一个全局的list，详细记录了每个节点的name, labels, PodCIDR, host IP, hostname, 总资源（cpu，内存），可分配资源，各心跳状态（网络，内存，硬盘，PID数量，kubelet等），kubelet的物理port，各k8s组件image信息，node环境信息（os, CRI version, kubeProxy version, kubelet version等）。
PersistentVolumeClaim: 记录用户对持久化存储的要求。
PersistentVolume: 是一个全局的object，记录了所有的持久化存储设备的信息（类似于node）
Pod: 是对于使用k8s的开发者而言最重要的资源，其中包含ownerReference (Node, Demonset等），containers相关信息（image，启动命令，probe，资源信息，存储信息，结束时行，是否接受service注入环境变量为等），网络设置（dns设置，port设置等），集群调度相关信息（优先级，tolerations，affinity，重启规则等），pod状态（hostIP，podIP，启动时间等）
PodTemplate: 一般是被包含在其它资源中的一部分，比如Jobs, DaemonSets, Replication Controllers。其初始化刚被创建的pod的k8s相关的信息，一般是label等。
Replication Controller: 是系统内建的最常用的controller，用来保证Pod的实际运行数量满足定义。如果不足则负责创建，如果过多则通知一些pod terminate。
ResourceQuota: 用于记录和限制某个namespace的中的总的资源消耗，一般用于多用户下利用namespace对资源进行限制。
Secrets: 实际上将文件内容通过base64编码后存在etcd中。在Pod中container启动时可以将secretes作为文件挂载在某一路径下，如此避免重要信息存储在image中。
ServiceAccout: 用于授权集群内的pod访问apiServer。
Service: 非常重要且常见的资源，用于对外提供统一的Service IP和port，将流量负载均衡调整至集群中多个pod。重要的配置有：cluster IP，port，selector（选择转发流量的目的pod），sessionAffinity等。这里提供的负载均衡是L3 TCP的。
MutatingWebhookConfiguration: 不明（内部object）
ValidatingWebhookConfiguration: 不明（内部object）
CustomerResourceDefinitions: 自定义资源也是非常重要的一种资源，是各种k8s插件能够存在的基础。比如当要实现Clico之类的自定义插件时，首先需要考虑的就是apiServer如何能够处理相关的请求信息。自定义资源的定义便是apiServer处理资源的依据。这个话题比较复杂，在这里不详细讨论。
APIService: 定义API服务的资源。一个API请求有两种形式，/apis/GROUP/VERSION/*这种不被包含在namespace中的（即全局的）和/apis/GROUP/VERSION/namespaces/NAMESPACE/*这种被包含在namespace中的。当一个请求到达apiServer后，必然需要有相应的代码去处理它。每一对GROUP和VERSION确定一种API，响应每一种API请求的代码被抽象为一种服务（service）。想象一下自定义资源的相关API请求到达apiServer后如何被处理呢？相关的service也是自定义的并且运行在master中，k8s正是根据APIService来正确地将请求与正确的service关联。在这里可以定义service名称，安全设置，优先级等。
ControllerRevision: 是一个beta功能，用于Controller保存自己的历史状态便于更新和回滚。
Daemenset: 常见的Pod set种类，用于控制每种pod状态（数量，计算资源使用，probe等）在定义的范围内，且在每node上最多有一个。
Replicaset: 常见的Pod set种类但现在基本上不直接使用，用于控制每种pod的状态（数量，计算资源使用，probe等）在定义的范围内。一个Replicasets中的各个pod都应是等同的、可互换的，即对外表现完全相同。就好比所有的氢原子（1质子0中子）都是不可区分的。
Deployment: 最常见的Pod set种类，可以拥有Replicasets和Pod。用于控制拥有的资源的状态（数量，计算资源使用，probe等）在定义的范围内。
StatefulSet: 常见的Pod set种类。和Deployment的区别之处是它控制的pod不是可互换的而是在整个生命周期有不变的标签。这样，每个pod可以有自己的DNS名，存储等。即使pod被删除后这些信息也会被恢复。
TokenReview: 不明，似乎和apiServer的Webhook的token授权相关。
LocalSubjectAccessReview: 不明（内部object），和一个命名空间对用户/组的授权检查相关。
SelfSubjectAccessReview: 不明（内部object），和当前用户检查自己是否有权限对一个命名空间进行操作相关。
SelfSubjectRulesReivew: 不明（内部object），含有当前用户在一个命名空间内能进行的操作的列表。和apiServer的授权安全模式相关
SubjectAccessReviews: 不明（内部object），和用户/组的授权检查相关，并不限定于某个命名空间。
HorizontalPodAutoScaler: 控制Pod set（比如Deployment）的pod数量的资源。可以根据pod的CPU、内存、自定义数据动态调节pod数量。在这里可以找到相关的例子。
CronJob: 定时运行Job pod的资源。
Job: 常见的Pod set种类，会创建一定数量的pod，仅当特定数量的pod成功结束后这个Job才算成功结束。创建的pod在结束后不会被重启。
CertificateSigningRequests: 可以认为是一个接口，便于Pod等资源申请获得一个X.509证书。这个证书应该被controller approve或者被手动approve，之后被合适的对象签名。具体可以参考这里。
Lease: 是一个在1.13版本中加入的资源类型，用于Node向master通知自己的心跳信息。之前的版本中kebulet是通过更新NodeStatus通知master心跳，后来发现NodeStatus太大了而心跳信息更新频繁，导致master压力较大，于是增加了Lease这种资源。
EndpointSlice: 是含有一个service的Endpoint的部分信息的资源。原因和Lease类似，对于含有较多信息的service（比如有很多pod分布在多个node上）一个endpoint object可能会比较大而且被频繁访问，所以这种情况下会有多个endpointSlice被创建减轻master的压力。
Event: 描述一个集群内的事件的资源，含有message，event，reason，报告来源等详细的信息。
Ingresse (APIGroup=extensions): 将被deprecated。
Ingresse (APIGroup=http://networking.k8s.io): 可以简单理解为是定义loadbalancer的资源。其中含有一系列规则，定义了不同url的对应后端，SSL termination等。为什么这个新的API会取代前面那个APIGroup=extensions的Ingress API呢？我查了很多地方没有找到具体的文字解释，但是可以推测是Ingress正式成为k8s的网络模块的一部分，对应的server（代码）从extensions迁移到http://networking.k8s.io。
NetworkPolicy: 定义了那些网络流量可以去哪些pod的资源。一个NetworkPolicy可以指定一组pods，定义只有满足了特定条件（比如源/目的IP，port，pod名等）的网络流量可以被相应的pod收发。
RuntimeClass: 这是2019年讨论加入的新API资源。文档说明其目的是将容器运行时（Container Runtime）环境的属性暴露给k8s的控制层，便于在一个集群或节点中支持多种容器运行时环境。这样便于未来创建更具有兼容性的k8s集群。
PodDisruptionBudget: 这一个API资源使用户可以对一组pod定义“k8s可以容忍的实际running状态的pod数量与预期的差距”。考虑这样一个情景：一集群中某个service后一共有5个相同pod处理其流量，要求至少有一半的pod是可用的，但其中3个pod由于调度运行在node A上。如果出现node A突然故障等情况导致服务不可用，暂时没有好的办法处理这种不可避免地意外情况（或者需要让调度算法知道这些pod应该被尽量均匀分布在个节点上，但目前k8s没有功能强制这种调度）。但除此之外还有很多可以避免的意外情况，比如在集群维护或者其它事件的处理过程中，集群管理员可能drain node A，导致三个pod同时被结束从而影响业务。针对这种可避免的意外，若一组pod的数量因为可避免的k8s操作将会低于可容忍程度（在PodDisruptionBudget中定义），那么这个命令会被阻止并返回失败。
PodSecurityPolicy: 定义了一个pod在集群中被创建/运行/更新时需要满足的条件。
ClusterRole: 定义了集群中policy rule的一些常见集合，比如system-node等，用于控制账户权限。
ClusterRoleBinding: 定义了某个账户/组对ClusterRole的引用，用于赋权。
Roles: 和前面ClusterRole类似，但是顾名思义ClusterRole是和集群账户相关，Role则被用于其它的账户（比如controller使用的service account）
RoleBindings: 定义了某个账户/组对Role的引用，用于赋权。
PriorityClass: 定义了pod优先级名称和对应值的映射。比如system-cluster-critical对应的优先级为2000000000。值越大代表优先级越高，那么当集群资源不足等情况发生必须终止一些pod时，优先级小的pod会先被终止。为什么不直接用数值代表优先级呢？因为这样子很容易出现确定随意性。比如开发人员A开发了一个非常重要的pod，于是在代码中将其优先级的值设置为9999。但是集群集群管理员B可能认为9999是一个小数字，他创建的随便一个pod的优先级都是999999+。于是需要PriorityClass来进行优先级的统一管理和比较。
CSIDriver: 定义了集群中容器存储驱动的API资源。CSI代表的是Container Storage Interface，即容器存储接口。k8s应该可以利用各种各样的存储服务，各家云厂商的活开源的。k8s如何知道怎么去用这些存储服务呢？那么就是通过这个CSIDriver资源找到相应的驱动。
CSINode: 前面CSIDriver产生的节点相关的信息便存在CSINode中。
StorageClass: 定义了可以存在的存储类型的API资源。
Volumeattachments: 定义了对一个node分配/回收存储空间的请求的API资源。
NetworkSets: 接下来的都是Calico自定义API资源，就不一一分析了，都与网络协议/安全/管理相关。
NetworkPolicies: Calico自定义API资源
IPPools: Calico自定义API资源
IPAMHandles: Calico自定义API资源
IPAMConfigs: Calico自定义API资源
IPAMBlocks: Calico自定义API资源
HostEndpoints: Calico自定义API资源
GlobalNetworkSets: Calico自定义API资源
GlobalNetworkPolicies: Calico自定义API资源
FelixConfiguration: Calico自定义API资源
ClusterInformation: Calico自定义API资源
BlockAffinity: Calico自定义API资源
BGPPeer: Calico自定义API资源
BGPConfiguration: Calico自定义API资源
<!-- more -->
https://zhuanlan.zhihu.com/p/115903242

https://www.jianshu.com/p/57ee9f8129e7

资源对象
资源对象通常有3个组成部分：

metadata：这是关于资源的元数据，比如它的名称、类型、api版本、注释和标签。
spec：这是由用户定义的希望系统最终达到的状态，比如启动多少个replica、cpu和内存的限制等等。
status：系统的当前状态，由服务器去更新。
资源的操作
资源通常有这几种操作：创建（Create），更新（Update），读取（Read），删除（Delete）。其中更新又分为替换（Replace）和打补丁（Patch），区别是替换是把整个spec替换掉，会有乐观锁保证读写安全；打补丁则是修改某些指定的字段，没有锁，最后一次写会成功。

有部分资源还会支持下列操作：

Rollback: 将PodTemplate回滚到以前的版本。
Read / Write Scale: 读取或更新资源的副本数量。
Read / Write Status: 读取或更新资源对象的状态。


