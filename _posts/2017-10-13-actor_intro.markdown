---
title: Actor模型原理
layout: post
category: spark
author: 夏泽民
---
<!-- more -->
1.Actor模型
在使用Java进行并发编程时需要特别的关注锁和内存原子性等一系列线程问题，而Actor模型内部的状态由它自己维护即它内部数据只能由它自己修改(通过消息传递来进行状态修改)，所以使用Actors模型进行并发编程可以很好地避免这些问题，Actor由状态(state)、行为(Behavior)和邮箱(mailBox)三部分组成

状态(state)：Actor中的状态指的是Actor对象的变量信息，状态由Actor自己管理，避免了并发环境下的锁和内存原子性等问题
行为(Behavior)：行为指定的是Actor中计算逻辑，通过Actor接收到消息来改变Actor的状态
邮箱(mailBox)：邮箱是Actor和Actor之间的通信桥梁，邮箱内部通过FIFO消息队列来存储发送方Actor消息，接受方Actor从邮箱队列中获取消息
Actor的基础就是消息传递

2.使用Actor模型的好处：
事件模型驱动--Actor之间的通信是异步的，即使Actor在发送消息后也无需阻塞或者等待就能够处理其他事情
强隔离性--Actor中的方法不能由外部直接调用，所有的一切都通过消息传递进行的，从而避免了Actor之间的数据共享，想要
观察到另一个Actor的状态变化只能通过消息传递进行询问
位置透明--无论Actor地址是在本地还是在远程机上对于代码来说都是一样的
轻量性--Actor是非常轻量的计算单机，单个Actor仅占400多字节，只需少量内存就能达到高并发
3.Actor模型原理

创建ActorSystem
ActorSystem作为顶级Actor，可以创建和停止Actors,甚至可关闭整个Actor环境，
此外Actors是按层次划分的，ActorSystem就好比Java中的Object对象，Scala中的Any，
是所有Actors的根，当你通过ActorSystem的actof方法创建Actor时，实际就是在ActorSystem
下创建了一个子Actor。
可通过以下代码来初始化ActorSystem

val system = ActorSystem("UniversityMessageSystem")

通过ActorSystem创建TeacherActor的代理(ActorRef)
看看TeacherActor的代理的创建代码

val teacherActorRef:ActorRef = system.actorOf(Props[TeacherActor])

ActorSystem通过actorOf创建Actor，但其并不返回TeacherActor而是返
回一个类型为ActorRef的东西。
ActorRef作为Actor的代理，使得客户端并不直接与Actor对话，这种Actor
模型也是为了避免TeacherActor的自定义/私有方法或变量被直接访问，所
以你最好将消息发送给ActorRef，由它去传递给目标Actor
发送QuoteRequest消息到代理中
你只需通过!方法将QuoteReques消息发送给ActorRef(注意：ActorRef也有个tell方法,其作用就委托回调给!)

techerActorRef!QuoteRequest
等价于teacherActorRef.tell(QuoteRequest, teacherActorRef)

MailBox

每个Actor都有一个MailBox,同样，Teacher也有个MailBox，其会检查MailBox并处理消息。
MailBox内部采用的是FIFO队列来存储消息，有一点不同的是，现实中我们的最新邮件
会在邮箱的最前面。
Dispatcher

Dispatcher从ActorRef中获取消息并传递给MailBox,Dispatcher封装了一个线程池，之后在
线程池中执行MailBox。

protected[akka] override def registerForExecution(mbox: Mailbox, ...): Boolean = {
  ...
 try {
 executorService execute mbox
 ...
}

为什么能执行MailBox?
看看MailBox的实现,没错，其实现了Runnable接口

private[akka] abstract class Mailbox(val messageQueue: MessageQueue) extends SystemMessageQueue with Runnable

