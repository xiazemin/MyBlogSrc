---
title: Conductor
layout: post
category: web
author: 夏泽民
---
Netflix Conductor框架是典型的服务编排框架，通过Conductor还可以实现工作流和分布式调度，性能非常卓越。

流程
流程引擎默认是用DSL来编写流程定义文件，这是一种JSON格式的文件，我们的工作流案例就是以这个定义文件为驱动的，但是很可惜目前Conductor只支持手写定义，无法通过界面生成，这块就需要后面通过改造Conductor来增加相应功能。
任务
这里面包括的主要是和任务相关的功能，通过这个功能可以进行简单工作流的实现，还可以进行并行计算。
历史
如果想要查看之前进行过的（完成，失败等终态）历史任务，通过这个功能就可以实现。
监控
当工作流任务流程非常冗长的时候，我们对每个节点的任务运行情况并不了解，这时候就需要有一个任务监控功能及时知道流程的状态方便我们做出相应决策。同时还有一个重要功能是任务调度，通过这个功能可以实现类似xxl-job的功能，满足分布式定时调度的需求。
客户端和通信
这二个功能本是一体的，既然Conductor是分布式的任务流程那么核心原理就是通过Server+Worker的方式，利用核心状态机发消息的方式来驱动客户端的任务执行，而Worker的实现是跨语言的，可以用JAVA、Python、go等语言实现，而Worker需要长轮询Server端的状态来判断当然是否有自己的任务来执行。
管理后台
通过管理后台可以查看任务和工作流的元数据定义，工作流的执行状态等。
<!-- more -->
https://www.jianshu.com/p/4eae1af8afa8
https://github.com/Netflix/conductor


https://netflix.github.io/conductor/intro/

https://netflix.github.io/conductor/intro/#installing-and-running

https://www.jianshu.com/p/c0611dada7d6

git clone https://github.com/Netflix/conductor.git

cd server进入server目录，再执行../gradlew server，第一次启动可能会很慢，如果失败就重试几次。

https://www.jianshu.com/p/75b4ac6deb50
http://wiki.hualala.com/pages/viewpage.action?pageId=28446373

Conductor 是 Netflix 受需要运行全球流媒体业务流程的启发，构建的基于云的微服务编排引擎。

Conductor 管理工作流，可以暂停和重新启动进程，并使用基于 JSON DSL 的蓝图来定义执行流。 它还具有可视化流程流的用户界面，并可扩展到数百万个并发运行的流程流。

https://www.oschina.net/p/netflix-conductor?hmsr=aladdin1e1

https://www.cnblogs.com/rongfengliang/p/11362106.html

https://blog.csdn.net/u011868076/article/details/74231528

http://blog.sina.com.cn/s/blog_493a84550102yvg8.html
http://www.uml.org.cn/zjjs/201702082.asp?artid=18972

http://dockone.io/article/2298
https://gitbook.cn/books/5db93922b3d7b75070d1c6a9/index.html
我们将 Conductor 作为“编排引擎”构建，以满足以下需求，在应用程序中消除了模板，同时提供反应流：

使用基于 JSON DSL 的 Blueprint 定义执行流程
跟踪并管理工作流
能够暂停，恢复和重新启动流程
用户界面可视化处理流程
能够在需要时同步处理所有任务
能够扩展为百万级并发运行的流程
由客户端提取出来的的队列服务支撑
能够通过 HTTP 或其他方式操作，例如 GRPC
Conductor 旨在满足上述需求，现在已在 Netflix 使用了将近一年。迄今为止，它调度超过 260 万个工作流，包括从简单的线性工作流到运行多天的非常复杂的动态工作流。

https://my.oschina.net/u/347227/blog/1499635
https://www.oschina.net/p/netflix-conductor
