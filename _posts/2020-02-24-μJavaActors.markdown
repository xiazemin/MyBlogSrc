---
title: μJavaActors
layout: post
category: algorithm
author: 夏泽民
---
https://www.ibm.com/developerworks/cn/java/j-javaactors/

Kilim	http://www.malhar.net/sriram/kilim/	一个支持基于轻型线程的多生成者、单使用者邮箱模型的 Java 库。
Akka	http://akka.io/	尝试使用函数语言模拟 actor 的模式匹配，一般使用 instanceof 类型检查（但 μJavaActors 一般使用字符串同等性或正则表达式匹配）。
GPars	http://gpars.codehaus.org/Actor	Groovy Actor 库。
<!-- more -->
即使 Java 6 和 Java 7 中引入并发性更新，Java 语言仍然无法让并行编程变得特别容易。Java 线程、synchronized 代码块、wait/notify 和 java.util.concurrent 包都拥有自己的位置，但面对多核系统的容量压力，Java 开发人员正在依靠其他语言中开创的技术。actor 模型就是这样一项技术，它已在 Erlang、Groovy 和 Scala 中实现。本文为那些希望体验 actor 但又要继续编写 Java 代码的开发人员带来了 μJavaActors 库。

用于 JVM 的另外 3 个 actor 库
请参阅 “表 1：对比 JVM actor 库”，快速了解 3 个用于 JVM 的流行的 actor 库与 μJavaActors 的对比特征。

μJavaActors 库 是一个紧凑的库，用于在 Java 平台上实现基于 actor 的系统（μ 表示希腊字母 Mμ，意指 “微型”）。在本文中，我使用 μJavaActors 探讨 actor 在 Producer/Consumer 和 Map/Reduce 等常见设计模式中的工作原理。

您随时可以 下载 μJavaActors 库的源代码。

Java 平台上的 actor 并发性
这个名称有何含义？具有任何其他名称的 actor 也适用！
基于 actor 的系统 通过实现一种消息传递 模式，使并行处理更容易编码。在此模式中，系统中的每个 actor 都可接收消息；执行该消息所表示的操作；然后将消息发送给其他 actor（包括它们自己）以执行复杂的操作序列。actor 之间的所有消息是异步的，这意味着发送者会在收到任何回复之前继续进行处理。因此，一个 actor 可能终生都陷入接收和处理消息的无限循环中。

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
大部分程序都有许多 actor，这些 actor 常常具有不同的类型。actor 可在程序启动时创建或在程序执行时创建（和销毁）。本文中的 actor 包 包含一个名为 AbstractActor 的抽象类，actor 实现基于该类。
图 1 显示了 actor 之间的关系。每个 actor 可向其他 actor 发送消息。这些消息保存在一个消息队列（也称为邮箱；从概念上讲，每个 actor 有一个队列，当 ActorManager 看到某个线程可用于处理消息时，就会从队列中删除该消息，并将它传送给在线程下运行的 actor，以便处理该消息。

图 1. actor 之间的关系
actor 通过线程向其他 actor 发送消息
μJavaActors 的并行执行功能
现在您已可开始使用 μJavaActors 实现并行执行了。首先要创建一组 actor。这些是简单的 actor，因为它们所做的只是延迟少量时间并将消息发送给其他 actor。这样做的效果是创建一个消息风暴，您首先会看到如何创建 actor，然后会看到如何逐步分派它们来处理消息。

有两种消息类型：

initialization (init) 会导致 actor 初始化。仅需为每个 actor 发送一次这种类型的消息。
repeat 会导致 actor 发送 N-1 条消息，其中 N 是一个传入的消息参数。
清单 1 中的 TestActor 实现从 AbstractActor 继承的抽象方法。activate 和 deactivate 方法向 actor 通知它的寿命信息；此示例中不会执行任何其他操作。runBody 方法是在收到任何消息之前、首次创建 actor 的时候调用的。它通常用于将第一批消息引导至 actor。testMessage 方法在 actor 即将收到消息时调用；这里 actor 可拒绝或接受消息。在本例中，actor 使用继承的 testMessage 方法测试消息接受情况；因此接受了所有消息。

清单 1. TestActor
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
12
13
14
15
16
17
18
19
20
21
22
23
class TestActor extends AbstractActor {
 
  @Override
  public void activate() {
    super.activate();
  }
 
  @Override
  public void deactivate() {
    super.deactivate();
  }
 
  @Override
  protected void runBody() {
    sleeper(1);  // delay up to 1 second
    DefaultMessage dm = new DefaultMessage("init", 8);
    getManager().send(dm, null, this);
  }
 
  @Override
  protected Message testMessage() {
    return super.testMessage();
  }
loopBody 方法（如清单 2 中所示）在 actor 收到一条消息时调用。在通过较短延迟来模拟某种一般性处理之后，才开始处理该消息。如果消息为 “repeat”，那么 actor 基于 count 参数开始发送另外 N-1 条消息。这些消息通过调用 actor 管理器的 send 方法发送给一个随机 actor。

清单 2. loopBody()
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
12
13
@Override
protected void loopBody(Message m) {
  sleeper(1);
  String subject = m.getSubject();
  if ("repeat".equals(subject)) {
    int count = (Integer) m.getData();
    if (count > 0) {
      DefaultMessage dm = new DefaultMessage("repeat", count - 1);
      String toName = "actor" + rand.nextInt(TEST_ACTOR_COUNT);
      Actor to = testActors.get(toName);
      getManager().send(dm, this, to);
    }
  }
如果消息为 “init”，那么 actor 通过向随机选择的 actor 或一个属于 common 类别的 actor 发送两组消息，启动 repeat 消息队列。一些消息可立即处理（实际上在 actor 准备接收它们且有一个线程可用时即可处理）；其他消息则必须等待几秒才能运行。这种延迟的消息处理对本示例不是很重要，但它可用于实现对长期运行的流程（比如等待用户输入或等待对网络请求的响应到达）的轮询。

清单 3. 一个初始化序列
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
12
13
14
else if ("init".equals(subject)) {
  int count = (Integer) m.getData();
  count = rand.nextInt(count) + 1;
  for (int i = 0; i < count; i++) {
    DefaultMessage dm = new DefaultMessage("repeat", count);
    String toName = "actor" + rand.nextInt(TEST_ACTOR_COUNT);
    Actor to = testActors.get(toName);
    getManager().send(dm, this, to);
     
    dm = new DefaultMessage("repeat", count);
    dm.setDelayUntil(new Date().getTime() + (rand.nextInt(5) + 1) * 1000);
    getManager().send(dm, this, "common");
  }
}
否则，表明消息不适合并会报告一个错误：

1
2
3
4
5
6
    else {
      System.out.printf("TestActor:%s loopBody unknown subject: %s%n", 
        getName(), subject);
    }
  }
}
主要程序包含清单 4 中的代码，它在 common 类别中创建了 2 个 actor，在 default 类别中创建了 5 个 actor，然后启动它们。然后 main 至多会等待 120 秒（sleeper 等待它的参数值的时间约为 1000ms），定期显示进度消息。

清单 4. createActor、startActor
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
12
13
14
15
16
17
18
19
20
21
22
23
DefaultActorManager am = DefaultActorManager.getDefaultInstance();
:
Map<String, Actor> testActors = new HashMap<String, Actor>();
for (int i = 0; i < 2; i++) {
    Actor a = am.createActor(TestActor.class, "common" + i);
    a.setCategory("common");
    testActors.put(a.getName(), a);
}
for (int i = 0; i < 5; i++) {
    Actor a = am.createActor(TestActor.class, "actor" + i);
    testActors.put(a.getName(), a);
}
for (String key : testActors.keySet()) {
   am.startActor(testActors.get(key));
}    
for (int i = 120; i > 0; i--) {
    if (i < 10 || i % 10 == 0) {
        System.out.printf("main waiting: %d...%n", i);
    }
    sleeper(1);
}
:
am.terminateAndWait();

