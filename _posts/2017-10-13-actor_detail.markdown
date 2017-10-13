---
title: Actor系统的实体
layout: post
category: spark
author: 夏泽民
---
<!-- more -->
在Actor系统中，actor之间具有树形的监管结构，并且actor可以跨多个网络节点进行透明通信。 
对于一个Actor而言，其源码中存在Actor，ActorContext，ActorRef等多个概念，它们都是为了描述Actor对象而进行的不同层面的抽象。 
我们先给出一个官方的示例图，再对各个概念进行解释。

   <img src="{{site.url}}{{site.baseurl}}/img/ActorPath.png"/>

    上图很清晰的展示了一个actor在源码层面的不同抽象，和不同actor之间的父子关系： 
Actor类的一个成员context是ActorContext类型，ActorContext存储了Actor类的上下文，包括self、sender。 
ActorContext还混入了ActorRefFactory特质，其中实现了actorOf方法用来创建子actor。 
这是Actor中context的源码：

{% highlight scala linenos %}
trait Actor {
  /**
   * Stores the context for this actor, including self, and sender.
   * It is implicit to support operations such as `forward`.
   *
   * WARNING: Only valid within the Actor itself, so do not close over it and
   * publish it to other threads!
   *
   * [[akka.actor.ActorContext]] is the Scala API. `getContext` returns a
   * [[akka.actor.UntypedActorContext]], which is the Java API of the actor
   * context.
   */
  implicit val context: ActorContext = {
    val contextStack = ActorCell.contextStack.get
    if ((contextStack.isEmpty) || (contextStack.head eq null))
      throw ActorInitializationException(
        s"You cannot create an instance of [${getClass.getName}] explicitly using the constructor (new). " +
          "You have to use one of the 'actorOf' factory methods to create a new actor. See the documentation.")
    val c = contextStack.head
    ActorCell.contextStack.set(null :: contextStack)
    c
  }
  {% endhighlight %}
ActorCell的self成员是ActorRef类型，ActorRef是一个actor的不可变，可序列化的句柄（handle），它可能不在本地或同一个ActorSystem中，它是实现网络空间位置透明性的关键设计。 
这是ActorContext中self的源码：

trait ActorContext extends ActorRefFactory {

  def self: ActorRef
ActorRef的path成员是ActorPath类型，ActorPath是actor树结构中唯一的地址，它定义了根actor到子actor的顺序。 
这是ActorRef中path的源码：

abstract class ActorRef extends java.lang.Comparable[ActorRef] with Serializable {
  /**
   * Returns the path for this actor (from this actor up to the root actor).
   */
  def path: ActorPath
Actor引用

Actor引用是ActorRef的子类，它的最重要功能是支持向它所代表的actor发送消息。每个actor通过self来访问它的标准（本地）引用，在发送给其它actor的消息中也缺省包含这个引用。反过来，在消息处理过程中，actor可以通过sender来访问到当前消息的发送者的引用。

不同类型的Actor引用

根据actor系统的配置，支持几种不同的actor引用：

纯本地引用被配置成不支持网络功能的，这些actor引用发送的消息不能通过一个网络发送到另一个远程的JVM。
支持远程调用的本地引用使用在支持同一个jvm中actor引用之间的网络功能的actor系统中。为了在发送到其它网络节点后被识别，这些引用包含了协议和远程地址信息。
本地actor引用有一个子类是用在路由（比如，混入了Router trait的actor）。它的逻辑结构与之前的本地引用是一样的，但是向它们发送的消息会被直接重定向到它的子actor。
远程actor引用代表可以通过远程通讯访问的actor，i.e. 从别的jvm向他们发送消息时，Akka会透明地对消息进行序列化。
有几种特殊的actor引用类型，在实际用途中比较类似本地actor引用： 
PromiseActorRef表示一个Promise，作用是从一个actor返回的响应来完成，它是由akka.pattern.ask调用来创建的
DeadLetterActorRef是死信服务的缺省实现，所有接收方被关闭或不存在的消息都在此被重新路由。
EmptyLocalActorRef是查找一个不存在的本地actor路径时返回的：它相当于DeadLetterActorRef，但是它保有其路径因此可以在网络上发送，以及与其它相同路径的存活的actor引用进行比较，其中一些存活的actor引用可能在该actor消失之前得到了。
然后有一些内部实现，你可能永远不会用上： 
有一个actor引用并不表示任何actor，只是作为根actor的伪监管者存在，我们称它为“时空气泡穿梭者”。
在actor创建设施启动之前运行的第一个日志服务是一个伪actor引用，它接收日志事件并直接显示到标准输出上；它就是Logging.StandardOutLogger。
获得Actor引用

创建Actor

一个actor系统通常是在根actor上使用ActorSystem.actorOf创建actor，然后使用ActorContext.actorOf从创建出的actor中生出actor树来启动的。这些方法返回指向新创建的actor的引用。每个actor都拥有到它的父亲，它自己和它的子actor的引用。这些引用可以与消息一直发送给别的actor，以便接收方直接回复。

具体路径查找

另一种查找actor引用的途径是使用ActorSystem.actorSelection方法，也可以使用ActorContext.actorSelection来在actor之中查询。它会返回一个（未验证的）本地、远程或集群actor引用。向这个引用发送消息或试图观察它的存活状态会在actor系统树中从根开始一层一层从父向子actor发送消息，直到消息到达目标或是出现某种失败，i.e.路径中的某一个actor名字不存在（在实际中这个过程会使用缓存来优化，但相较使用物理actor路径来说仍然增加了开销，因为物理路径能够从actor的响应消息中的发送方引用中获得），这个消息传递过程由Akka自动完成的，对客户端代码不可见。 
使用相对路径向兄弟actor发送消息：

context.actorSelection("../brother") ! msg
1
也可以用绝对路径：

context.actorSelection("/user/serviceA") ! msg
1
查询逻辑Actor层次结构

由于actor系统是一个类似文件系统的树形结构，对actor的匹配与unix shell中支持的一样：你可以将路径（中的一部分）用通配符(«*» 和«?»)替换来组成对0个或多个实际actor的匹配。由于匹配的结果不是一个单一的actor引用，它拥有一个不同的类型ActorSelection，这个类型不完全支持ActorRef的所有操作。同样，路径选择也可以用ActorSystem.actorSelection或ActorContext.actorSelection两种方式来获得，并且支持发送消息。 
下面是将msg发送给包括当前actor在内的所有兄弟actor：

context.actorSelection("../*") ! msg
1
与远程部署之间的互操作

当一个actor创建一个子actor，actor系统的部署者会决定新的actor是在同一个jvm中或是在其它的节点上。如果是在其他节点创建actor，actor的创建会通过网络连接来到另一个jvm中进行，结果是新的actor会进入另一个actor系统。 远程系统会将新的actor放在一个专为这种场景所保留的特殊路径下。新的actor的监管者会是一个远程actor引用（代表会触发创建动作的actor）。这时，context.parent（监管者引用）和context.path.parent（actor路径上的父actor）表示的actor是不同的。但是在其监管者中查找这个actor的名称能够在远程节点上找到它，保持其逻辑结构，e.g.当向另外一个未确定(unresolved)的actor引用发送消息时。 


因为设计分布式执行会带来一些限制，最明显的一点就是所有通过电缆发送的消息都必须可序列化。虽然有一点不太明显的就是包括闭包在内的远程角色工厂，用来在远程节点创建角色（即Props内部）。 
另一个结论是，要意识到所有交互都是完全异步的，它意味着在一个计算机网络中一条消息需要几分钟才能到达接收者那里（基于配置），而且可能比在单JVM中有更高丢失率，后者丢失率接近于0（还没有确凿的证据）。

Akka使用的特殊路径

在路径树的根上是根监管者，所有的的actor都可以从通过它找到。在第二个层次上是以下这些：

"/user"是所有由用户创建的顶级actor的监管者，用ActorSystem.actorOf创建的actor在其下一个层次 are found at the next level。
"/system" 是所有由系统创建的顶级actor（如日志监听器或由配置指定在actor系统启动时自动部署的actor）的监管者。
"/deadLetters" 是死信actor，所有发往已经终止或不存在的actor的消息会被送到这里。
"/temp"是所有系统创建的短时actor(i.e.那些用在ActorRef.ask的实现中的actor)的监管者。
"/remote" 是一个人造的路径，用来存放所有其监管者是远程actor引用的actor。
附录-Actor模型概述：

Actor模型为编写并发和分布式系统提供了一种更高的抽象级别。它将开发人员从显式地处理锁和线程管理的工作中解脱出来，使编写并发和并行系统更加容易。Actor模型是在1973年Carl Hewitt的论文中提的，但是被Erlang语言采用后才变得流行起来，一个成功案例是爱立信使用Erlang非常成功地创建了高并发的可靠的电信系统。

Actor的树形结构

像一个商业组织一样，actor自然会形成树形结构。程序中负责某一个功能的actor可能需要把它的任务分拆成更小的、更易管理的部分。为此它启动子Actor并监管它们。要知道每个actor有且仅有一个监管者，就是创建它的那个actor。 


Actor系统的精髓在于任务被分拆开来并进行委托，直到任务小到可以被完整地进行处理。 这样做不仅使任务本身被清晰地划分出结构，而且最终的actor也能按照它们“应该处理的消息类型”，“如何完成正常流程的处理”以及“失败流程应如何处理”来进行解析。如果一个actor对某种状况无法进行处理，它会发送相应的失败消息给它的监管者请求帮助。这样的递归结构使得失败能够在正确的层次进行处理。

可以将这与分层的设计方法进行比较。分层的设计方法最终很容易形成防御性编程，以防止任何失败被泄露出来。把问题交由正确的人处理会是比将所有的事情“藏在深处”更好的解决方案。

现在，设计这种系统的难度在于如何决定谁应该监管什么。这当然没有一个唯一的最佳方案，但是有一些可能会有帮助的原则：

如果一个actort管理另一个actor所做的工作，如分配一个子任务，那么父actor应该监督子actor，原因是父actor知道可能会出现哪些失败情况，知道如何处理它们。
如果一个actor携带着重要数据（i.e. 它的状态要尽可能地不被丢失），这个actor应该将任何可能的危险子任务分配给它所监管的子actor，并酌情处理子任务的失败。视请求的性质，可能最好是为每一个请求创建一个子actor，这样能简化收集回应时的状态管理。这在Erlang中被称为“Error Kernel Pattern”。
如果actor A需要依赖actor B才能完成它的任务，A应该观测B的存活状态并对收到B的终止提醒消息进行响应。这与监管机制不同，因为观测方对监管机制没有影响，需要指出的是，仅仅是功能上的依赖并不足以用来决定是否在树形监管体系中添加子actor。
Actor实体

一个Actor是一个容器，它包含了 状态，行为，一个邮箱，子Actor和一个监管策略。所有这些包含在一个Actor引用里。

状态

Actor对象通常包含一些变量来反映actor所处的可能状态。这可能是一个明确的状态机，或是一个计数器，一组监听器，待处理的请求，等等。这些数据使得actor有价值，并且必须将这些数据保护起来不被其它的actor所破坏。

好消息是在概念上每个Akka actor都有它自己的轻量线程，这个线程是完全与系统其它部分隔离的。这意味着你不需要使用锁来进行资源同步，可以完全不必担心并发性地来编写你的actor代码。

在幕后，Akka会在一组线程上运行一组Actor，通常是很多actor共享一个线程，对某一个actor的调用可能会在不同的线程上进行处理。Akka保证这个实现细节不影响处理actor状态的单线程性。
由于内部状态对于actor的操作是至关重要的，所以状态不一致是致命的。当actor失败并由其监管者重新启动，状态会进行重新创建，就象第一次创建这个actor一样。这是为了实现系统的“自愈合”。

行为

每次当一个消息被处理时，消息会与actor的当前的行为进行匹配。行为是一个函数，它定义了处理当前消息所要采取的动作，例如如果客户已经授权过了，那么就对请求进行处理，否则拒绝请求。

邮箱

Actor的用途是处理消息，这些消息是从其它的actor（或者从actor系统外部）发送过来的。连接发送者与接收者的纽带是actor的邮箱：每个actor有且仅有一个邮箱，所有的发来的消息都在邮箱里排队。排队按照发送操作的时间顺序来进行，这意味着从不同的actor发来的消息在运行时没有一个固定的顺序，这是由于actor分布在不同的线程中。从另一个角度讲，从同一个actor发送多个消息到相同的actor，则消息会按发送的顺序排队。

可以有不同的邮箱实现供选择，缺省的是FIFO：actor处理消息的顺序与消息入队列的顺序一致。这通常是一个好的选择，但是应用可能需要对某些消息进行优先处理。在这种情况下，可以使用优先邮箱来根据消息优先级将消息放在某个指定的位置，甚至可能是队列头，而不是队列末尾。如果使用这样的队列，消息的处理顺序是由队列的算法决定的，而不是FIFO。

Akka与其它actor模型实现的一个重要差别在于当前的行为必须处理下一个从队列中取出的消息，Akka不会去扫描邮箱来找到下一个匹配的消息。无法处理某个消息通常是作为失败情况进行处理，除非actor覆盖了这个行为。

子Actor

每个actor都是一个潜在的监管者：如果它创建了子actor来委托处理子任务，它会自动地监管它们。子actor列表维护在actor的上下文中，actor可以访问它。对列表的更改是通过context.actorOf(...)创建或者context.stop(child)停止子actor来实现，并且这些更改会立刻生效。实际的创建和停止操作在幕后以异步的方式完成，这样它们就不会“阻塞”其监管者。

监督策略

Actor的最后一部分是它用来处理其子actor错误状况的机制。错误处理是由Akka透明地进行处理的。由于策略是actor系统组织结构的基础，所以一旦actor被创建了它就不能被修改。

考虑对每个actor只有唯一的策略，这意味着如果一个actor的子actor们应用了不同的策略，这些子actor应该按照相同的策略来进行分组，生成中间的监管者，又一次倾向于根据任务到子任务的划分来组织actor系统的结构。