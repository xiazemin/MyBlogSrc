---
title: μJavaActors
layout: post
category: algorithm
author: 夏泽民
---
即使 Java 6 和 Java 7 中引入并发性更新，Java 语言仍然无法让并行编程变得特别容易。Java 线程、synchronized 代码块、wait/notify 和 java.util.concurrent 包都拥有自己的位置，但面对多核系统的容量压力，Java 开发人员正在依靠其他语言中开创的技术。actor 模型就是这样一项技术，它已在 Erlang、Groovy 和 Scala 中实现。
<!-- more -->
μJavaActors 库 是一个紧凑的库，用于在 Java 平台上实现基于 actor 的系统（μ 表示希腊字母 Mμ，意指 “微型”）。
通过实现一种消息传递 模式，使并行处理更容易编码。在此模式中，系统中的每个 actor 都可接收消息；执行该消息所表示的操作；然后将消息发送给其他 actor（包括它们自己）以执行复杂的操作序列。actor 之间的所有消息是异步的，这意味着发送者会在收到任何回复之前继续进行处理。因此，一个 actor 可能终生都陷入接收和处理消息的无限循环中。

当使用多个 actor 时，独立的活动可轻松分配到多个可并行执行消息的线程上（进而分配在多个处理器上）。一般而言，每个 actor 都在一个独立线程上处理消息。一些 actor 系统静态地向 actor 分配线程；而其他系统（比如本文中介绍的系统）则会动态地分配它们。

μJavaActors 简介
μJavaActors 是 actor 系统的一个简单的 Java 实现。只有 1,200 行代码，μJavaActors 虽然很小，但很强大。在下面的练习中，您将学习如何使用 μJavaActors 动态地创建和管理 actor，将消息传送给它们。
μJavaActors 围绕 3 个核心界面而构建：

消息 是在 actor 之间发送的消息。Message 是 3 个（可选的）值和一些行为的容器：
source 是发送 actor。
subject 是定义消息含义的字符串（也称为命令）。
data 是消息的任何参数数据；通常是一个映射、列表或数组。参数可以是要处理和/或其他 actor 要与之交互的数据。
subjectMatches() 检查消息主题是否与字符串或正则表达式匹配。
μJavaActors 包的默认消息类是 DefaultMessage。
ActorManager 是一个 actor 管理器。它负责向 actor 分配线程（进而分配处理器）来处理消息。ActorManager 拥有以下关键行为或特征：
createActor() 创建一个 actor 并将它与此管理器相关联。
startActor() 启动一个 actor。
detachActor() 停止一个 actor 并将它与此管理器断开。
send()/broadcast() 将一条消息发送给一个 actor、一组 actor、一个类别中的任何 actor 或所有 actor。
在大部分程序中，只有一个 ActorManager，但如果您希望管理多个线程和/或 actor 池，也可以有多个 ActorManager。此接口的默认实现是 DefaultActorManager。
Actor 是一个执行单元，一次处理一条消息。Actor 具有以下关键行为或特征：
每个 actor 有一个 name，该名称在每个 ActorManager 中必须是惟一的。
每个 actor 属于一个 category；类别是一种向一组 actor 中的一个成员发送消息的方式。一个 actor 一次只能属于一个类别。
只要 ActorManager 可以提供一个执行 actor 的线程，系统就会调用 receive()。为了保持最高效率，actor 应该迅速处理消息，而不要进入漫长的等待状态（比如等待人为输入）。
willReceive() 允许 actor 过滤潜在的消息主题。
peek() 允许该 actor 和其他 actor 查看是否存在挂起的消息（或许是为了选择主题）。
remove() 允许该 actor 和其他 actor 删除或取消任何尚未处理的消息。
getMessageCount() 允许该 actor 和其他 actor 获取挂起的消息数量。
getMaxMessageCount() 允许 actor 限制支持的挂起消息数量；此方法可用于预防不受控制地发送。
大部分程序都有许多 actor，这些 actor 常常具有不同的类型。actor 可在程序启动时创建或在程序执行时创建（和销毁）。
actor 包 包含一个名为 AbstractActor 的抽象类，actor 实现基于该类。

每个 actor 可向其他 actor 发送消息。这些消息保存在一个消息队列（也称为邮箱；从概念上讲，每个 actor 有一个队列，当 ActorManager 看到某个线程可用于处理消息时，就会从队列中删除该消息，并将它传送给在线程下运行的 actor，以便处理该消息。

首先要创建一组 actor。这些是简单的 actor，因为它们所做的只是延迟少量时间并将消息发送给其他 actor。这样做的效果是创建一个消息风暴，您首先会看到如何创建 actor，然后会看到如何逐步分派它们来处理消息。

有两种消息类型：
initialization (init) 会导致 actor 初始化。仅需为每个 actor 发送一次这种类型的消息。
repeat 会导致 actor 发送 N-1 条消息，其中 N 是一个传入的消息参数。


μJavaActors 库具有灵活的、动态的行为，为 Akka 等更加庞大的 actor 库提供了一个基于 Java 的替代方案。
μJavaActors 可跨一个执行线程池高效地分配 actor 消息处理工作。而且，可在用户界面中迅速确定是否需要更多线程。该界面还容易确定哪些 actor 渴求工作或者是否有一些 actor 负载过重。

DefaultActorManager（ActorManager 接口的默认实现）可保证没有 actor 会一次处理多条消息。因此这会减轻 actor 作者的负担，他们无需处理任何重新输入考虑因素。该实现还不需要 actor 同步，只要：(1) actor 仅使用私有（实例或方法本地的）数据，(2) 消息参数仅由消息发送者编写，以及 (3) 仅由消息接收者读取。

DefaultActorManager 的两个重要的设计参数是线程与 actor 的比率 以及要使用的线程总数。线程数量至少应该与计算机上的处理器一样多，除非一些线程为其他用途而保留。因为线程可能常常空闲（例如，当等待 I/O 时），所以正确的比率常常是线程是处理器的 2 倍或多倍。一般而言，应该有足够的 actor（其实是 actor 之间的消息比率）来保持线程池中大部分时间都很繁忙。（为了获得最佳的响应，应该有一些保留线程可用；通常平均 75% 到 80% 的活动比率最佳。）这意味着 actor 通常比线程更多，因为有时 actor 可能没有任何要处理的挂起消息。当然，您的情况可能有所不同。执行等待操作（比如等待一个人为响应）的 actor 将需要更多线程。（线程在等待时变为 actor 专用的，无法处理其他消息。）

DefaultActorManager 很好地利用了 Java 线程，因为在 actor 处理一条消息时，一个线程仅与一个特定的 actor 关联；否则，它可供其他 actor 自由使用。这允许一个固定大小的线程池为无限数量的 actor 提供服务。结果，需要为给定的工作负载创建的线程更少。这很重要，因为线程是重量级的对象，常常被主机操作系统限制于相对较少数量的实例。μJavaActors 库正是因为这一点而与为每个 actor 分配一个线程的 actor 系统区分开来；如果 actor 没有消息要处理，并且可能限制了可存在的 actor 实例数量，这么做可以让线程实际空闲下来。

在线程切换方面，μJavaActors 实现有很大不同。如果在消息处理完成时有一条新消息需要处理，则不会发生线程切换；而是会重复一个简单循环来处理该新消息。因此，如果等待的消息数量至少与线程一样多，则没有线程是空闲线程，因此不需要进行切换。如果存在足够的处理器（至少一个线程一个），则可以有效地将每个线程分配给一个处理器，而从不会发生线程切换。如果缓冲的消息不足，那么线程将会休眠，但这并不明显，因为只有在没有工作挂起时才会出现负载过重的现象。

1. Kilim	http://www.malhar.net/sriram/kilim/	一个支持基于轻型线程的多生成者、单使用者邮箱模型的 Java 库。	Kilim 需要字节代码调整。在 μJavaActors 中，每个 actor 也是其自身的邮箱，所以不需要独立的邮箱对象。
2. Akka	http://akka.io/	尝试使用函数语言模拟 actor 的模式匹配，一般使用 instanceof 类型检查（但 μJavaActors 一般使用字符串同等性或正则表达式匹配）。	Akka 功能更多（比如支持分布式 actor），因此比 μJavaActors 更大且有可能更复杂。
3. GPars	http://gpars.codehaus.org/Actor	Groovy Actor 库。	类似于 μJavaActors，但更适合 Groovy 开发人员。
