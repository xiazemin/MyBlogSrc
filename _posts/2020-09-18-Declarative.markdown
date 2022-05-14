---
title: Declarative
layout: post
category: k8s
author: 夏泽民
---
Declarative（声明式设计）指的是这么一种软件设计理念和做法：我们向一个工具描述我们想要让一个事物达到的目标状态，由这个工具自己内部去figure out如何令这个事物达到目标状态。

和Declarative（声明式设计）相对的是Imperative或Procedural（过程式设计）。两者的区别是：在Declarative中，我们描述的是目标状态（Goal State），而在Imperative模式中，我们描述的是一系列的动作。这一系列的动作如果被正确的顺利执行，最终结果是这个事物达到了我们期望的目标状态的。

声明式（Declarative）的编程方式一直都会被工程师们拿来与命令式（Imperative）进行对比，这两者是完全不同的编程方法。我们最常接触的其实是命令式编程，它要求我们描述为了达到某一个效果或者目标所需要完成的指令，常见的编程语言 Go、Ruby、C++ 其实都为开发者了命令式的编程方法，

声明式和命令式是两种截然不同的编程方式:

在命令式 API 中，我们可以直接发出服务器要执行的命令，例如： “运行容器”、“停止容器”等；
在声明式 API 中，我们声明系统要执行的操作，系统将不断向该状态驱动。
<!-- more -->
SQL 其实就是一种常见的声明式『编程语言』，它能够让开发者自己去指定想要的数据是什么。或者说，告诉数据库想要的结果是什么，数据库会帮我们设计获取这个结果集的执行路径，并返回结果集。众所周知，使用 SQL 语言获取数据，要比自行编写处理过程去获取数据容易的多。

SELECT * FROM posts WHERE user_id = 1 AND title LIKE 'hello%';
我们来看看相同设计的 YAML，利用它，我们可以告诉 Kubernetes 最终想要的是什么，然后 Kubernetes 会完成目标。

例如，在 Kubernetes 中，我们可以直接使用 YAML 文件定义服务的拓扑结构和状态：

apiVersion: v1
kind: Pod
metadata:
  name: rss-site
  labels:
    app: web
spec:
  containers:
    - name: front-end
      image: nginx
      ports:
        - containerPort: 80
    - name: rss-reader
      image: nickchase/rss-php-nginx:v1
      ports:
        - containerPort: 88


https://skyao.io/learning-cloudnative/declarative/

https://cloud.tencent.com/developer/article/1080886

https://zhuanlan.zhihu.com/p/89446640

https://www.cnblogs.com/yuxiaoba/p/9803284.html


首先，什么是 Kubernetes 的 API ？
Kubernetes 有很多能力，这些能力都是通过各种 API 对象来提供。也就是说，API 对象正是我们使用 Kubernetes 的接口，我们正是通过操作这些提供的 API 对象来使用 Kubernetes 能力的。
可以看一下整个 Kubernetes 里的所有 API 对象，实际上就可以用如下的树形结构表示出来：
API 对象树形结构图
图 1 Kubernetes 里的所有 API 对象树形结构图
注：Kubernetes 项目中，一个 API 对象在 Etcd 里的完整资源路径，是由：Group（API 组）、Version（API 版本）和 Resource（API 资源类型）三个部分组成的。
有了这些 API ，就可以向这些 API 发送请求来操作 Kubernetes 了。

那么，声明式 API 的 “声明式” 是什么意思？
对于我们使用 Kubernetes API 对象的方式，一般会编写对应 API 对象的 YAML 文件交给 Kubernetes（而不是使用一些命令来直接操作 API）。所谓“声明式”，指的就是我只需要提交一个定义好的 API 对象来“声明”（这个 YAML 文件其实就是一种“声明”），表示所期望的最终状态是什么样子就可以了。而如果提交的是一个个命令，去指导怎么一步一步达到期望状态，这就是“命令式”了。
“命令式 API”接收的请求只能一个一个实现，否则会有产生冲突的可能；“声明式API”一次能处理多个写操作，并且具备 Merge 能力。
二、Kubernetes 编程范式
首先借用磊哥的一句话：从“使用 Kubernetes 部署代码”，到“使用 Kubernetes 编写代码”的蜕变过程，正是你从一个 Kubernetes 用户，到 Kubernetes 玩家的晋级之路。
“使用 Kubernetes 编写代码”，就要遵循 Kubernetes 所提供的编程范式来进行开发。这个编程范式有2个过程：
（1）为 Kubernetes 添加自定义 API 对象。
（2）为自定义的 API 对象添加控制逻辑。

2.1 为 Kubernetes 添加自定义 API 对象
添加一个 Kubernetes 风格的 API 资源类型得益于一个全新的 API 插件机制：CRD。CRD 的全称是 Custom Resource Definition。它指的就是，允许用户在 Kubernetes 中添加一个跟 Pod、Node 类似的、新的 API 资源类型，即：自定义 API 资源。从而就可以达到通过标准的 kubectl create 和 get 操作，来管理自定义 API 对象的目的。
CRD（Custom Resource Definition）就是我们来提供一个 Definition，然后让 Kubernetes 认识其中定义的 CR（Custom Resource）。举个例子：

提交一个 CRD 的 YAML 文件来定义 CR
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: networks.samplecrd.k8s.io
spec:
  group: samplecrd.k8s.io
  version: v1
  names:
    kind: Network
    plural: networks
  scope: Namespaced
1
2
3
4
5
6
7
8
9
10
11
解释：这个 CRD 中，主要看 spec 字段，指定了“group: samplecrd.k8s.io”“version: v1”这样的 API 信息，也指定了这个 CR 的资源类型叫作 Network，复数（plural）是 networks。scope 是 Namespaced，代表定义的这个 Network 是一个属于 Namespace 的对象，类似于 Pod。

这样，就定义了 networks.samplecrd.k8s.io 这个 CR（Custom Resource），也就是/apis/samplecrd.k8s.io/v1/networks 这样一个自定义API对象了。

接下来就可以使用这个自定义 API 对象了
apiVersion: samplecrd.k8s.io/v1
kind: Network
metadata:
  name: example-network
spec:
  cidr: "192.168.0.0/16"
  gateway: "192.168.0.1"
1
2
3
4
5
6
7
解释： API 资源类型是 Network；API 组是samplecrd.k8s.io；API 版本是 v1。另外，可以看到这里还有“cidr”（网段），“gateway”（网关）这些关于对象描述的字段需要做些代码工作来让 Kubernetes“认识”了。

这样这个自定义的 API 对象有了，还需要为这个 API 对象添加“控制逻辑”，Kubernetes 才能根据 Network API 对象的“增、删、改”操作，在真实环境中做出相应的响应。也就是下一节内容。

2.2 为自定义的 API 对象添加控制逻辑
Kubernetes 所定义的 API 对象都是“声明式 API”，CRD 定义的 CR自然也不例外。“声明式 API”并不像“命令式 API”那样有着明显的执行逻辑，声明式 API 在于维护一个所声明的状态。这就使得基于声明式 API 的业务功能实现，往往需要通过控制器模式来“监视”API 对象的变化，然后以此来决定实际要执行的具体工作。

所以，我们需要开发一个“自定义控制器”来为 CRD 执行控制逻辑。

解释一下自定义控制器的工作原理，这张经典的图就可以了：
自定义控制器的工作流程示意图
控制器其实可以分为2部分来看：左边的 Informer，右面的”控制循环“，中间通过一个工作队列来进行协同。

Informer 主要用来与APIServer 进行数据同步。就是一个自带缓存和索引机制，可以触发 Handler 的客户端库。这个本地缓存在 Kubernetes 中一般被称为 Store，索引一般被称为 Index。Informer 使用了 Reflector 包，它是一个可以通过 ListAndWatch 机制获取并监视 API 对象变化的客户端封装。Reflector 和 Informer 之间，用到了一个“增量先进先出队列”进行协同。

”控制循环“部分是一个“无限循环”的任务，每一个循环周期执行的正是我们真正关心的业务逻辑，通过对比“期望状态”和“实际状态”的差异，不断完成调协（Reconcile）。

另外，这套流程不仅可以用在自定义 API 资源上，也完全可以用在 Kubernetes 原生的默认 API 对象上。在自定义控制器里面，可以通过对自定义 API 对象和默认 API 对象进行协同，从而实现更加复杂的编排功能。

一句话总结：CRD + 自定义控制器
Kubernetes 编程范式就是：CRD + 自定义控制器。
（CRD 用来创建自定义的 API 对象，自定义控制器来添加相应 API 的请求控制逻辑）

三、Operator 简单理解
Operator 是 Kubernetes 编程范式的聪明利用，为“有状态应用”设计 CRD 及其 Controller。在 Kubernetes 生态中，Operator 是一个相对更加灵活和编程友好的管理“有状态应用”的解决方案。

- 原理：
Operator 的工作原理实际上是利用了 Kubernetes 的自定义 API 资源（CRD），来描述我们想要部署的“有状态应用”；然后在自定义控制器里，根据自定义 API 对象的变化，来完成具体的部署和运维工作。

- 开发：
kubebuilder与operator-sdk ：是目前开发operator两种常用的SDK，或者叫 framework，不管是哪种方式只是规范了步骤和一些必要的元素，生成相应的模板。都是为了更快更好造出 CRD 及其 Controller ，实现 Operator 来更好管理 Kubernetes 有状态应用。


https://blog.csdn.net/KevinBetterQ/article/details/104012293

https://www.codenong.com/p12230853/

https://www.yunforum.net/group-topic-id-2780.html

http://penguincj.com/post/cloud/k8s/201904-k8s-declarative-api/
https://www.jianshu.com/p/46a7164d4598
https://aijishu.com/a/1060000000115599