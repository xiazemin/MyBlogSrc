---
title: Scheduler
layout: post
category: golang
author: 夏泽民
---
https://chanjarster.github.io/post/go/scheduling-in-go-part-1/
https://chanjarster.github.io/post/go/scheduling-in-go-part-2/
https://chanjarster.github.io/post/go/scheduling-in-go-part-3/
几个数字
operation	cost
1纳秒	可以执行12条指令
OS上下文切换	~1000到~1500 nanosecond，相当于~12k到~18k条指令。
Go程上下文切换	~200 nanoseconds，相当于~2.4k instructions条指令。
访问主内存	~100到~300 clock cycles
访问CPU cache	~3到~40 clock cycles（根据不同的cache类型）
操作系统线程调度器
你的程序实际上就是一系列需要执行的指令，而这些指令是跑线程里的。

线程可以并发运行：每个线程轮流占用一个core；也可以并行运行：每个线程跑在不同core上。

操作系统线程调度器负责保证充分利用core来执行线程。
<!-- more -->
程序指令是如何执行的
程序计数器（program counter，PC），有时也称为指令指针（instruction pointer，IP），用来告诉线程下一个要执行的指令（注意不是当前正在执行的指令）的位置。它是一种寄存器（register）。

每次执行指令的时候都会更新PC，因此程序才能够顺序执行。

线程状态
Waiting：等待中。原因：等待硬件（比如磁盘、网络）、正在系统调用（syscall）、阻塞在同步上（atomic、mutex）
Runnable：可以运行，正在等待调度。越多线程等待调度，大家就等的越久，且分配到的时间就越少。
Executing：正在某个core上运行。
任务类型
CPU绑定：这种任务永远不会让线程进入Waiting状态，比如计算Pi。
IO绑定：这种任务会让线程进入Waiting状态。
上下文切换
Linux、Mac和Windows使用的是抢占式调度器，所以：

你无法预测调度器什么时候会运行哪个线程。线程优先级混合事件（比如接收网络数据），也使得预测调度器行为变得不可能。
如果你要有确定的行为，那么就应该对线程做同步和编排（synchronization and orchestration）。否则你观察到现在是这个样子的，无法保证下次还是这个样子的。
在一个core上切换线程的物理行为称为上下文切换（context switching）。调度器把一个线程从core上换下来，然后把另一个线程换上去。换上去的线程状态从Runnable->Executing，换下来的线程的状态从Executing->Runnable（如果依然可以运行），或者Executing->Waiting（因为等待所以被换下来）。

上下文切换的代价比较高，大概在~1000到~1500 nanosecond之间，考虑到core大致每纳秒可以执行12条指令，那么就相当于浪费了~12k到~18k的指令。

如果是IO绑定任务，那么上下文切换能够有效利用CPU，因为A线程进入Waiting那么B线程就可以顶上使用CPU。

如果是CPU绑定任务，那么上下文切换会造成性能损失，因为把CPU能力白白浪费在上下文切换上了（浪费了~12k到~18k的指令）。

少即是多
越少的线程带来越少的调度开销，每个线程能分配到的时间就越多，那么就能完成越多的工作。

Cache line
访问主内存（main memory）的数据的延迟大概在~100到~300 clock cycles。

访问cache的数据延迟大概在 ~3到~40 clock cycles（根据不同的cache类型）。


CPU会把数据从主内存中copy到cache中，以cache line为单位，每条cache line为64 bytes。所以多线程修改内存会造成性能损失。

多个并行运行的线程访问同一个数据或者相邻的数据，那么它们可能就会访问同一条cache line。任何线程跑在任何core上都有一份自己的cache line copy。所以就有了False Sharing问题：


只要一个线程操作了自己core上的某个cache line，那么这个cache line在其他core就会变脏（cache coherency），当一个线程访问一个脏cache line的时候，就要访问一下main memory（~100到~300 clock cycles）。当单处理器core变多的时候，以及当有多个处理器（处理器间通信）的时候，这个开销就变得很大了。

逻辑组件
P：Logical Processor，你有多少个虚拟core就有多少个P。之所以说虚拟core是因为如果处理器支持支持一个core有多个硬件线程（Hyper-Threading，超线程），那么每个硬件线程就算作一个虚拟core。runtime.NumCPU()能够得到虚拟core的数量。

M：操作系统线程。这个线程依然由操作系统调度。每个P被分配一个M。

G：Go程。Go程实际上是协程（Coroutine）。和操作系统线程有点像，只是操作系统线程上下文切换在core上，而Go程上下文切换在M上。Go程的上下文切换发生在用户空间，开销更低。

运行队列（队列中的G都是runnable的）：

LRQ（Local Run Queue）。每个P会给一个LRQ，P负责把LRQ中的G上下文切换到M上。
GRQ（Global Run Queue），GRQ放还未分配到P的G。
这是一张全景图：


协作式调度器
和操作系统的抢占式调度器不同，Go采用的是协作式调度器。Go调度器是Go运行时的一部分，而Go运行时在内置在你的程序里。所以Go调度器运行在用户空间。

Go调度器运行在用户空间，那么就需要定义明确的发生在safepoint的用户空间事件来做调度决策。不过程序员不需要太多关心这个，同时也无法控制调度行为，所以Go调度器虽然是协作式的但看起来像是抢占式的。

Go程状态
Waiting：Go程停止了，且在等待什么事情发生。比如等待操作系统（syscall）、同步调用（atomic和mutex操作）
Runnable：Go程想要M的时间来执行指令。越多的Go程想要时间，就以为着等待越长的时间，每个Go程能分到的时间就越少。
Executing：Go程正在M上执行指令。
上下文切换
Go调度程序需要定义明确的用户空间事件，这些事件发生在代码中的安全点，以便进行上下文切换。安全点体现在函数调用中。所以函数调用很重要。在Go 1.11之前，如果程序在跑一个很长的循环且循环里没有函数调用，那么就会导致调度器和GC被推迟。

4类事件允许调度器做调度决策，注意调度不一定会发生，而是给了调度器一个机会而已：

使用go
垃圾收集
系统调用
同步和编排（Synchronization and Orchestration）
使用go

go创建了一个新的Go程，自然调度器有机会做一个调度决策。

垃圾收集

垃圾收集跑在自己的Go程里，需要征用M来运行，因此调度器也需要做决策

系统调用

系统调用会导致Go程阻塞M。调度器有些时候会把这个G从M换下（上下文切换），然后把新的G换上M。也有可能创建一个新的M，用来执行P的LRQ中的G。

同步和编排

如果atomic、mutex、channel操作阻塞了一个G，调度器会把一个新的G去运行，等到它又能运行了（从阻塞中解除），那么再把它放到队列中，然后最终跑在到M上。

异步系统调用
比如MacOS中的kqueue、Linux中的epoll、Windows的iocp都是异步网络库。G做这些异步系统调用并不会阻塞M，那么就意味着M可以用来执行LRQ中的其他M。下面是图解：

G1准备做网络调用：


G1移到了Net Poller，然后M可以跑G2


G1就绪了，就回到LRQ中，等待被调度，整个过程不需要新的M：


同步系统调用
文件IO不是异步的，所以G会把M给阻塞，那么Go调度器会这么做：

G1调用了阻塞系统调用：


M1连带G1从P脱离（此时M1因为阻塞被操作系统上下文切换下去了），创建新的M2给P，把G2调度到M2上：


而后G1从阻塞中恢复，追加到LRQ中等待下次调度，M1则保留下来等待以后使用：


Work Stealing
虽然名字叫做工作偷窃，但实际上是好事。简单来说就是当P1没有G时，把P2的LRQ中的G“偷”过来执行，借此来提高M的利用率。

看下图中P1和P2都有3个G等待调度，GRQ中有一个G


这个时候P1先把自己的G都处理完了：


P1会“偷”P2 LRQ中一半的G，偷窃算法如下，简单来说就是先偷P2的G，如果没有再从GRQ中取：

runtime.schedule() {
    // only 1/61 of the time, check the global runnable queue for a G.
    // if not found, check the local queue.
    // if not found,
    //     try to steal from other Ps.
    //     if not, check the global runnable queue.
    //     if not found, poll network.
}

当P2把G都做完了，然后P1没有G在LRQ中时：


根据前面讲的算法，P2会拿GRQ中的G来运行：


实际的例子
下面拿一个实际的例子来告诉你Go调度器是如何比你直接用OS线程做更多工作的。

协作式OS线程程序
有两个OS线程，T1和T2，它们之间的交互式这样的：

T2等待消息，T1发送消息，T1等待消息
T2接收消息，T2发送消息，T2等待消息
T1接收消息。。。
。。。
T1一开始在C1上，T2处于等待状态：


T1发送消息给T2，进入等待，从C1脱离；T2收到消息后调度到C2上：


T2发送消息给T1，进入等待，从C2脱离；T1收到消息后调度到C3上：


所以你可以看到T1和T2频繁发生OS上下文切换，而这个代价是很高的（见文头表格）。同时每次切换到不同core上，导致cache miss，所以还存在访问主内存的开销。

协作式Go程程序
下面来看看Go调度器怎么做的：

G1一开始在M1上，而M1和C1绑定，G2处于等待状态：


G1发消息给G2，进入等待，从M1脱离；G2收到消息被调度到M1：


G2发消息给G1，进入等待，从M1脱离；G1收到消息被调度到M1：


所以Go程调度的优势：

OS线程始终保持运行，没有进入waiting
Go程上下文切换不是发生在OS层面，代价相对低， ~200 nanoseconds 或 ~2.4k instructions。
始终都是在同一个core上，优化了cache miss的问题，这个对于NUMA架构特别友好。


不是所有问题都可以concurrency
不论你遇到什么种类的问题，你应该先求出一个正确的sequential解，然后再看这个问题是否有可能作出concurrency解。

什么是concurrency
Concurrency意味着乱序执行，拿一组原本顺序执行的指令，把它们乱序执行依然能够得到相同的结果。对于你来说就要去权衡concurrency之后得到的性能好处和其带来的复杂度。而且有些问题乱序执行压根就没道理，只能顺序执行。

并行和并发的区别在于，并行是指在不同的OS线程上，OS线程在不同的core上同时执行不相干的指令。


上图中，P1和P2有自己的OS线程，OS线程在不同的core上，因此G1和G2是并行的。

但是在P1和P2自己看来，它有3个G要执行，而这三个G共享同一个OS线程/core，而且执行顺序是不定的，它们是并发执行的。

工作负载
前面讲过，有两种类型的工作负载：

CPU-Bound：永远不会使得Go程处于waiting状态，永远处于runnable/executing状态的纯计算型任务。
IO-Bound：天然的会使得Go程进入waiting状态。比如访问网络资源、syscall、访问文件。同时也把同步事件归到此类（atomic、mutex）
对于CPU-Bound任务来说，你需要利用并行。如果Go程数量多于OS线程/core数量，那么就会使得Go程被上下文切换，从而带来性能损失。

对于IO-Bound任务来说，你可以不需要利用并行，一个OS线程/core可以很轻松的处理这种天然就会进出waiting状态的任务。Go程数量大于OS线程/core数量可以大大提高OS线程/core的利用率，提高任务的处理速度。Go程的上下文切换不会造成性能损失，因为你的任务自己就会停止。

那么使用多个Go程能带来多大好处，以及多少个Go程能带来最大效果，那么就需要benchmark才能知道。