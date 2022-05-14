---
title: 系统调度
layout: post
category: golang
author: 夏泽民
---
程序其实就是一堆需要按照顺序一个接一个执行的机器指令。为此，操作系统使用了一个线程的概念。线程的工作就是按顺序执行分配给它的指令集，直到没有指令可以执行了为止。
你运行的每一个程序都会创建一个进程，并且每一个进程都会有一个初始线程。线程拥有创建更多线程的能力。这些不同的线程都是独立运行的，调度策略都是在线程这一级别上的，而不是进程级别（或者说调度的最小单元是线程而不是进程）。线程是可以并发执行的（轮流使用同一个核），或并行（每个线程互不干扰的同时在不同的核上跑）。线程还维护这他们自己状态，好保证安全、隔离、独立的执行自己的指令。
系统调度器负责保证当有线程可以执行时，CPU 是不能处于空闲状态的。它还必须创建一个所有线程同时都在运行的假象。在创造这个假象的过程中，调度器需要优先运行优先级更高的线程。但是低优先级的线程又不能被饿死（就是一直不被运行）。调度器还需要通过快速、明智的决策尽可能的最小化调度延迟。
https://www.ardanlabs.com/blog/2018/08/scheduling-in-go-part1.html
https://www.ardanlabs.com/blog/2018/08/scheduling-in-go-part2.html
<!-- more -->
执行指令
Program counter(PC)，有时候也被叫做指令指针(instruction pointer, 简称IP)，线程用它来跟踪下一个要执行的指令。在大多数处理器中，PC 指向的是下一个指令，而不是当前指令。
如果你曾经看过 Go 程序的 stack trace，你可能注意到了每行的最后都有一个 16 进制数字。比如 +0x39和0x72。
goroutine 1 [running]:
   main.example(0xc000042748, 0x2, 0x4, 0x106abae, 0x5, 0xa)
       stack_trace/example1/example1.go:13 +0x39                 <- LOOK HERE
   main.main()
       stack_trace/example1/example1.go:8 +0x72                  <- LOOK HERE


这些数字表示的是 PC 值距离所在函数的顶部的偏移量。0x39 这个 PC 偏移量代表了线程要执行的在 example 函数中的下一个指令（如果程序没有崩溃）。0x72 代表的是程序返回到 main 后，要执行的下一个指令。更重要的是，在这个指针之前的那个指令，表示的是正在执行的指令。

源代码。
func main() {
    example(make([]string, 2, 4), "hello", 10)
}

func example(slice []string, str string, i int) {
   panic("Want stack trace")

16 进制数字 +0x39 代表了距离 example 函数第一条指令后面 57(0x39的10进制值) 字节的那个指令。下面我们通过对二进制文件执行 objdump，来看看这个 example 函数。从下面的汇编代码中找到第 12 条指令。注意上面代码中调用 panic 的那条指令。
$ go tool objdump -S -s "main.example" ./example1
TEXT main.example(SB) stack_trace/example1/example1.go
func example(slice []string, str string, i int) {
  0x104dfa0     65488b0c2530000000  MOVQ GS:0x30, CX
  0x104dfa9     483b6110        CMPQ 0x10(CX), SP
  0x104dfad     762c            JBE 0x104dfdb
  0x104dfaf     4883ec18        SUBQ $0x18, SP
  0x104dfb3     48896c2410      MOVQ BP, 0x10(SP)
  0x104dfb8     488d6c2410      LEAQ 0x10(SP), BP
    panic("Want stack trace")
  0x104dfbd     488d059ca20000  LEAQ runtime.types+41504(SB), AX
  0x104dfc4     48890424        MOVQ AX, 0(SP)
  0x104dfc8     488d05a1870200  LEAQ main.statictmp_0(SB), AX
  0x104dfcf     4889442408      MOVQ AX, 0x8(SP)
  0x104dfd4     e8c735fdff      CALL runtime.gopanic(SB)
  0x104dfd9     0f0b            UD2              <--- LOOK HERE PC(+0x39)

线程状态
另一个重要的概念就是线程状态，它决定了线程在调度器中的角色。一个线程有 3 中状态：阻塞态、就绪态和运行态。

实际情况中不止这 3 个，阻塞也分可中断和不可中断 2 种，此外还有僵尸态、初始化状态等。


阻塞态：表示线程已经停止，需要等待一些事情发生后才可继续。这有很多种原因，比如需要等待硬件（磁盘或网络），系统调用，或者互斥锁（atomic, mutexes）。这类情况导致的延迟，往往是性能不佳的根本原因。

译者注：如果你发现，机器的 CPU 利用率很低，同时程序的 QPS 还很低，待处理的请求还有很多堆在后面，速度就是上不去。那说明你的线程都处于阻塞态等待被唤醒，没有干活。

就绪态：这代表线程想要一个 CPU 核来执行被分配的机器指令。如果你有很多个线程需要 CPU，那么线程就不得不等待更长时间。此时，因为许多的线程都在争用 CPU，每个线程得到的运行时间也就缩短了。
运行态：这表示线程已经被分配了一个 CPU 核，正在执行它的指令。与应用相关的工作正在被完成。这是每个人都想要的状态。
负荷类型
线程的工作有 2 种类型的负荷。第一种叫做 CPU密集型，第二种叫 IO密集型。
CPU密集：处理这种工作的线程从来不会主动进入阻塞态。它一直都需要使用 CPU。这种工作通常都是数学计算。比如计算圆周率的第 n 位的工作就属于 CPU密集型的工作。
IO密集：这种工作会使得线程进入阻塞态。常见于通过网络请求资源，或者进行了系统调用。一个需要访问数据库的线程属于 IO密集的。互斥锁的使用也属于这种。
上下文切换
Linux，Mac 或者 Windows 系统上都拥有抢占式调度器。这表明了很重要的几点。
第一，线程何时被调度器选中，被分配 CPU 时间片是不可预测的。线程的优先级和事件同时都会对调度结果有影响，这导致你不可能确定调度去什么时候能调度你的线程。

译者注：你动一下鼠标，敲一下键盘，这些动作都会触发 CPU 的中断响应，也就是事件。这都会对调度器的结果产生影响的，所以说它是不可预测的。

第二，这表明了，你绝不能凭感觉来写代码，因为你的经验不能保证总是应验。人是很容易平感觉下定论的，同样的事情反复出现了上千次，就认为它是百分百的。如果你想要有百分百的确定性，必须在线程中使用互斥锁。
在同一个 CPU 核上交换线程的行为过程，称为上下文切换。上下文切换发生时，调度器把一个线程从核上拿下来，把另一个就绪态的线程放到核上。从就绪队列中选中的这个线程，就这样被置成了运行态。而被从核上拿来下的那个线程，被置为就绪态（如果它仍然可以被执行），或者进入阻塞态（如果它是因为执行了 IO 操作才被替换的）。
上下文切换是昂贵的，因为在一个核上交换线程需要时间片。上下文切换造成的计算损失受很多因素影响，一般是 50 到 100 纳秒左右。
假设一个 CPU 核平均每纳秒可执行 12 个机器指令，一次上下文切换要执行 600 到 1200 条指令，那么本质上，你的程序在上下文切换期间丢失了可以执行大量指令的机会。
如果你有一个 IO 密集的工作，那么上下文切换会有一定的优势。一旦一个线程进入了阻塞态，另一个就绪态的线程就可以立马执行。这使得 CPU 一直都在工作。这是调度中最重要的目标，就是如果有线程可执行（处于就绪态），就不能让 CPU 闲着。
如果你的程序是 CPU 密集的，那么上下文切换将会是性能的噩梦。因为线程总是有指令要执行，而上下文切换中断了这个过程。这与 IO 密集型的工作形成了鲜明的对比。

在早期，CPU 只有单核的。调度也就没那么复杂。因为你只有一个单核 CPU。在任意一个时间点，只有一个线程可以运行。有一种轮询调度的方法，它尝试对每个就绪态的线程都执行一段时间。使用调度周期，除以线程总数，就是每个线程应该执行的时间。
比如，如果你定义你的调度周期是 10 毫秒，现在有 2 个线程，那么在一个调度周期内，每个线程可以执行 5 毫秒。如果你有 5 个线程，那么每个线程可以执行 2 毫秒。但是，如果你有 1000 个线程呢？每个线程执行 10 微妙是没有意义的，因为你大部分时间都花在了上下文切换上。
这时你就需要限制最短的执行时间应该是多少。在上面的那个场景中，如果最短的执行时间是 2 毫秒，同时你有 100 个线程，那么调度周期就需要增加到 2000 毫秒（2秒）。如果你有 1000 个线程，调度周期就要变成 20 秒。这个简单的例子中，一个线程要等 20 秒才能执行一次。
要知道这我们只是举了最简单调度场景。实际上调度器在做调度策略时需要考虑很多事情。这是你应该会想到一个常见并发手段，就是线程池的使用。让线程的数量在控制之内。
所以游戏规则就是『少即是多』，越少的就绪态线程意味着越少的调度工作，每个线程就会得到更多的时间。越多的就绪态线程意味着每个线程会得到越少的时间，也就意味着同一时间你能完成的工作越少（其他的 CPU 时间都被操作系统拿去做调度用了）。
寻找平衡
你需要在 CPU 核数和线程数量上寻找一个平衡，来使你的应用能够拥有最高的吞吐。当需要维持这个平衡时，线程池是一个最好的解决方案

在 Go 语言中根本不需要线程池。我认为 Go 语言最优秀的一点就是，它使得并发编程更简单了。
在写 Go 之前，我使用 C++ 和 C# 在 NT 上开发。在那个操作系统上，主要是使用 IOCP 线程池来完成并发编程的。作为一个工程师，你需要指定需要多少线程池，每个线程池的最大线程数是多少，来保证在固定核上达到最高的性能。
当写 web 服务的时候，需要和数据库打交道，每核 3 个线程的配置，似乎总能在 NT 平台上德奥最高的吞吐量。换句话说，就是每核 3 个线程可以使上下文切换的代价最小，从而最大化线程的执行时间。
如果配置每核只用 2 个线程，它会花费更多时间把工作完成，因为 CPU 会经常处在空闲状态。如果我一个核创建 4 个线程，它也会花费更长时间，因为上下文切换的代价会升高。每核 3 个线程的平衡，总是能得到最好的结果，不知道什么原因，它就是个魔法数字。
那如果你的服务的即有 CPU 密集的工作也有 IO 密集的工作呢？这可能会产生不同类型的延迟。这种情况就不太可能找到一个魔法数字来适用于所有情况。当使用线程池来调整服务的性能时，找到一个正确的一致配置是很复杂的。
Cache Line
访问主内存中的数据是有很高延迟的。大约 100 ~ 300 个时钟周期。所以 CPU 往往都会有一个本地 cache，使得数据距离需要它的线程所在的核更近。访问 cache 中的数据是很快的，和访问寄存器差不多。今天，性能优化中很重要的一部分就是怎么才能让 CPU 更快的得到数据，来减少数据访问的延迟。写多线程应用时面对状态异变问题时，需要考虑 cache 系统的机制

cacheline1 线程私有
cacheline2 核私有
cacheline3 核共享

cache line 是 cache 与主内存交换数据的最小单位。一个 cache line 是一块 64 字节的内存，用以在主内存和 cache 系统之间交换数据。每个核都有一份它自己所需要数据的拷贝。这就是为什么在多线程应用中，内存异变是导致性能问题的噩梦。因为 CPU Core 上运行的线程变了，不同的线程需要访问的数据不同，cache 里的数据也就失效了。
当多线程并行时，如果他们访问同样的数据，或者相邻很近的数据。他们将会访问同一个 cache line 中的数据。运行这些线程的任何一个核，都会在自己的 cache 上对数据做一份拷贝。也就是说，每个 CPU 核的 cache 中都有同样的一份数据拷贝，这些拷贝对应于内存中的同一块地址。

如果一个核上的线程，对拷贝数据进行了修改，那么硬件会将其他所有核上的 cache line 拷贝都标记为『脏』。当其他核上的线程试图访问或修改这个数据时，需要重新从主内存上拷贝最新的数据到自己的 cache 中。
也许 2 核的 CPU，不会出大问题，但如果是 32 核的 CPU 并行的运行着 32 个线程呢？如果系统有 2 个 CPU ，每个 CPU 有 16 核呢？这更糟糕，因为增加了 CPU 之间的的通信延迟。这个应用将会出现内存颠簸现象，性能会急剧下降，然而你可能都不知道为什么。

调度决策场景
假如现在要求你在以上给出信息的基础上，设计一个系统调度器了。想象一下你需要解决的这种场景。记住，上面所描述的，只是做调度决策时，需要考虑的众多情况之一。
现在假设机器只有一个单核CPU。你启动了应用，线程被创建了并且在一个 CPU 核上运行。随着线程开始执行它的指令，cache line 也开始检索数据，因为指令需要数据。现在线程要创建一个新的线程做一些并发处理。现在的问题是。
线程一旦创建并进入了就绪态，调度器有以下几种选择：

直接进行上下文切换，把主线程从 CPU 上拿掉？
这是对性能有益的，因为新线程需要同样的数据，而这些数据之前已经存在与 cache  上了。但主线程就不得不把时间片分给子线程了。
让新线程等待主线程执行完它的所有时间片？
线程没执行，执行时不必从主内存中同步数据到 local cache 了。
让线程等待其他可用核？
这就意味着被选中的核，要拷贝一份 cache line 中的数据，这会导致一定的延迟。但是新线程会立刻开始执行，主线程也能继续完成自己的工作。

用哪种方式呢？这就是系统调度在做调度决策时需要考虑的一个有趣的问题。答案是，如果有空闲的核，那就直接用。我们的目标是，如果有工作要做，就决不让 CPU 闲着。

当你的 Go 程序启动之初，它会被分配一个逻辑处理器(P)，这是为这台机器定义的一个虚拟 CPU Core。如果你的 CPU 的每个核带有多个hardware thread（Hyper-Threading），每一个 hardware 都会对应 Go 语言中的一个虚拟 core。

可以看见我有一个 4 core 处理器。但这里没有给出的是一个物理 core 有多少个 hardware thread。Intel Core i7 处理器拥有 Hyper-Threading，表示一个 core 上可以同时跑 2 个线程。这就表示 Go 程序里有 8 个虚拟 core(P) 可以使用，来让系统线程并行。

package main

import (
    "fmt"
    "runtime"
)

func main() {

    // NumCPU returns the number of logical
    // CPUs usable by the current process.
    fmt.Println(runtime.NumCPU())
}

当我在自己的机器上运行这个程序，NumCPU() 函数的结果是 8。我机器上运行的任何一个 Go 程序均会有 8 个 P 可以使用。
每一个 P 会被分配一个系统线程(M)。这个 M 会被操作系统调度，操作系统仍然负责将线程(M)放到一个 CPU Core 上去执行。这意味着当我在我的机器上运行程序，我有 8 个线程可以使用去执行我的操作，每个线程都被绑定上了一个独立的 P。
每一个 Go 程序也被赋予了一个初始的 goroutine(G)，它是 Go 程序的执行路径。一个 Goroutine 本质上就是一个 Coroutine，只不过因为在 Go 语言里，就改了个名字。你可以认为 Goroutine 是应用级别的线程，它在很多方面跟系统的线程是相似的。就像系统线程不断的在一个 core 上做上下文切换一样，Goroutine 不断的在 M 上做上下文切换。
最后一个难题就是运行队列。在 Go 调度器中有 2 个不同的执行队列：全局队列（Global Run Queue, 简称 GRQ）和本地队列（Local Run Queue，简称 LRQ）。每一个 P 都会有一个 LRQ 来管理分配给 P 上的 Goroutine。这些 Goroutine 轮流被交付给 M 执行。GRQ 是用来保存还没有被分配到 P 的 Goroutine。会有一个逻辑将 Goroutine 从 GRQ 上移动到 LRQ 上

系统调度器的行为是抢占式的。本质上就意味着你不能够预测调度器将会做什么。系统内核决定了一切，而这一切都是不可确定的。运行在系统上的应用无法控制内核中的调度逻辑，除非使用互斥锁之类的操作。
Go 调度器是 Go 运行时的一部分，Go 运行时被编译到了你的程序里。这就表示 Go 调度器是运行在用户态的，在内核之上。当前版本的 Go 调度器实现并不是抢占式的，而是一个协同调度器。这就意味着调度器需要明确定义用户态事件来指定调度决策。
非抢占式调度器的精彩之处在于，它看上去是抢占式的。你不能预知 Go 调度器将会做什么。因为调度器的调度决策权并没有交给开发者，而是在运行时里。
Goroutine 状态
就像线程，Goroutine 也拥有同样的 3 个高级状态。这决定了他们在 Go 调度器中扮演的角色。一个 Goroutine 有 3 中状态：阻塞态，就绪态，运行态
阻塞态： 这表示 Goroutine 被暂停了，要等待一些事情发生了才能继续。有可能是因为要等待系统调用或者互斥调用（atomic 和 mutex 操作）。这些情况导致性能不佳的根因。
就绪态： 这代表 Goroutine 想要一个 M 来执行分配给它的指令。如果你有很多 Goroutine 都需要 M，那么 Goroutine 就需要等较长的时间。并且，每个 Goroutine 被分配的执行时间也就更短了。这种情况也会导致性能下降。
运行态：这表示 Goroutine 被交给了一个 M 执行，正在执行它的指令。这是每个人都希望的。
上下文切换
Go 调度器需要有明确定义的用户态事件和代码安全点，来实现切换操作。这些事件和安全点通过函数调用表现的。函数调用对 Go 调度器是至关重要的。今天（Go1.11或更低版本），如果你运行一个死循环，循环内不做任何函数调用，你将在进程调度和垃圾回收上出现延迟；因为它没有给调度器机会对它进行切换。函数调用在合理的范围内发生是至关重要的。

注意: 对于 1.12 版本有一个建议，在 Go 调度器中增加抢占式调度机制，来允许高速循环被抢占。

有 4 种事件会引起 Go 程序触发调度。这不意味着每次事件都会触发调度。Go 调度器会自己找合适的机会。

使用关键字 go

垃圾回收
系统调用
同步互斥操作，也就是 Lock()，Unlock() 等

使用 go 关键字
关键字 go 是用来创建 Goroutine 的，一旦一个新的 Goroutine 被创建了，它就会引起 Go 调度器进行调度。
垃圾回收
因为 GC 操作是使用自己的一组 Goroutine 来执行的，这些 Goroutine 需要一个 M 来运行。所以 GC 会导致调度混乱。但是，因为调度器是知道 Goroutine 要做什么的，所以它可以做出明智的决策。其中一个明智的决策是，在 GC 过程中，暂停那些需要访问堆空间的 Goroutine，运行那些不需要访问堆空间的。

统调用
如果一个 Goroutine 执行了系统调用，就会导致 Goroutine 把 M 给阻塞了（就是运行这个 Goroutine 的 M 进入了阻塞态），有时调度器有能力把 Goroutine 从 M 上拿走，把一个新的 Goroutine 放到这个 M 上；但有时候新的 M 需要被创建出来，来保证 P 队列中其他的 Goroutine 能被执行。这块内容后面会详细说明。
同步互斥
如果 atomic, mutex, 或者 channel 操作 导致了 Goroutine 阻塞。调度器可以切换一个新的 Goroutine 去执行。一旦 Goroutine 在此可以执行了（进入就绪态）。它会被重新放到队列中，最终被 M 执行。
异步系统调用
如果你的操作系统有能力异步处理系统调用，那么 network poller 可以更有效的来完成系统调用。这方面在 kqueue(MacOS)，epoll(Linux) 或 iocp(Windows) 中都有不同方式的实现。
现在我们用的大多数操作系统，在网络调用上都是可以被异步执行的。这就是 network poller 名字的由来，因为他的主要用途就是处理网络请求的。通过使用 network poller 完成网络系统调用，调度器可以避免 Goroutine 在执行系统调用时把 M 阻塞住。这使得 M 可以去运行 P 的 LRQ 中其他的 Goroutine，而不需要再创建一个新的 M。这就减少了 OS 层面上的负载。
理解它是如何运作的最好方式就是通过例子。

G1 在 M 上执行，同事其他 3 个 Goroutine 在 LRQ 中等待 M。现在 network poller 没有事情可做。
G1 要执行一个网络调用，所以 G1 被移动到了 network poller 上，等待完成网络系统调用。一旦 G1 被移到了 network poller 中，M 现在就可以去执行 LRQ 中其他的 Goroutine 了。在这里例子中，G2 被切换到了 M 上。
异步网络系统调用完成后，G1 又被放回到了 P 的 LRQ中。一旦 G1 可以被切换到 M 上，处理网络请求结果相关的 Go 代码又能被执行了。这里最大的优势在于，执行网络系统调用，不需要额外的 M。network poller 有一个系统线程处理，它可以高效的处理事件的轮询。

同步系统调用
如果 Goroutine 想要执行一个系统调用不能被异步执行时，会发生什么？这种情况，network poller 是用不了的。执行系统调用的 Goroutine 会导致 M 阻塞住。很不幸，但是没有其他办法可以阻止这种情况的发生。一个不能使系统调用异步执行的例子就是文件系统的调用。如果你用 CGO，可能还有其他调用 C 函数的场景导致 M 阻塞。

注意：Windows 系统有异步处理文件访问的系统调用。技术上讲，在 Windows 上运行时，是可以利用 network poller 的。

让我们来看一下同步的系统调用会发生什么。
 G1 将要执行一个同步的系统调用，这将会阻塞 M1。
 调度器有能力认出 G1 导致 M 阻塞了。这时，调度器会将 M 从 P 上分离出去，G1 依旧附在 M 上被一起分离了。然后调度器获取一个 M2 为 P 服务。此时，G2 会被选中切换到 M2 上执行。如果 M2 因为以前的切换操作已经存在了，那么这次转换就要比重新创建一个 M 要快。
由 G1 执行的阻塞式系统调用完成了。此时，G1 可以被放回到 LRQ 并等待被 P 在此调度。M1 会放在一边等待以后使用，以防止这种情况在此发生。如果空闲 M 很多，调度器会主动让其退出。
工作窃取
调度器的另一部分就是，它是一个工作窃取机制。这保证在一些场景下能保证高效的调度。
让我们来看一个例子。
我们有一个多线程的 Go 程序带有 2 个 P，每个 P 都有 4 个 Goroutine 要执行，还有一个 Goroutine 在 GRQ 中。如果其中一个 P 很快的把所有的 Goroutine 都执行完了，会发生什么呢？
P1 没有 Goroutine 可以执行了，但是仍然有 Goroutine 是处于就绪态的等待被执行，P2 的 LRQ 和 GRQ 中都有。这是 P1 就需要窃取工作了。工作窃取的逻辑如下。
runtime.schedule() {
    // only 1/61 of the time, check the global runnable queue for a G.
    // if not found, check the local queue.
    // if not found,
    //     try to steal from other Ps.
    //     if not, check the global runnable queue.
    //     if not found, poll network.
}

所以基于上述注释描述的逻辑，P1 需要检查 P2 的 LRQ 中的 Goroutine 列表，把其中一半的 Goroutine 拿到自己的队列中。

P2 里的一半的 Goroutine 被交给了 P1 执行。

如果 P2 执行完毕了所有的 Goroutine，同时 P1 的 LRQ 中也没有可执行的 Goroutine 了，怎么办呢？

P2 完成了它的所有工作，现在需要窃取一些。首先它会检查一下 P1 的 LRQ，但是没有 Goroutine 可偷。下一步它会检查 GRQ。会找到 G9.
P2 从 GRQ 上偷窃了 G9，并开始执行。这一切工作窃取的好处就在于，它使 M 保持繁忙而不是空闲。这方面还有一些其他的好处，JBD 在它的博客中解释的很好。

想象一个 C 语言写的多线程应用，程序的逻辑就是两个系统线程彼此互相传递消息。
有 2 个线程正在互相传递消息。线程1 被切换都了 Core 1 上，现在正在执行，把他的消息发送给了线程2。

注意：消息是怎么被传递的不重要。重要的是这些线程在编排过程中的状态。
一旦线程1完成了发送消息，它就需要等待返回结果。这就导致了线程1会被从 Core 1 上切换下来，并置为阻塞态。一旦线程2收到了消息通知，它就变成就绪态。现在操作系统执行上下文切换，把线程2放到一个 Core 上执行，这放生在 Core 2 上。下一步，线程 2 处理消息把新的消息返回给线程1。
线程再一次被切换了，因为线程2的消息被线程1收到了。现在线程2从运行态切换到了阻塞态，线程1从阻塞态切换到运行态，最后再进入运行态，是它能够处理消息并发送新消息。
所有的这些上下文切换和状态转换都需要时间执行，这影响了工作被完成的速度。每次上下文切换都会导致 50 纳秒的延迟，并且如果硬件每秒钟能执行 12 条指令，那么大约就有 600 个指令，执行切换的期间是卡在那里的。因为这些线程被绑定在不同的核上，因缓存无法被命中而导致的额外的开销的可能性也很大。
让我们用 Go 调度器调度 Goroutine 来完成通用的操作。
有 2 个 Goroutine 彼此直线通过互相消息来协同工作。G1 切换到 M1 上，M1 被调度到 Core 1 上，使得 G1 能执行工作，这个工作就是 G1 发送消息给 G2。
一旦 G1 结束了发送消息，它现在需要等待相应。这会导致 G1 从 M1 上切换下来，进入阻塞态。一旦 G2 收到了消息通知，它被置为就绪态。现在 Go 调度去可以执行一次切换，让 G2 在 M1 上执行，它将依旧在 Core 1 上执行。然后，G2 处理完消息后，发送新消息 G1。
再次发生了上下文切换，因为 G2 发送的消息被 G1 收到了。现在 G2 从运行态切换到了阻塞态，G1 从阻塞态切换成了就绪态，最终再进入运行态，使得它可以运行并返回新消息。
表面上看没有什么不同。上下文切换和状态的改变，无论是线程还是 Goroutine 都是一样的。使用线程和 Goroutine 的一个主要的区别乍看上去不怎么明显。
在使用 Goroutine 的情况下，同一个系统线程和 Core 应用于整个处理流程中。这就表示，透过操作系统来看，系统线程从来没有被进入过阻塞态。结果我们因上下文切换中而损失的所有时间片，在使用 Goroutine 时都没有丢失。
本质上，Go 是在 OS 层面上，将 IO/阻塞操作转换成了 CPU 操作。因为所有的上下文切换都发生在应用层，我们没有丢掉每次切换造成的 600 个指令损失。调度器同时也增加了 cache-line 的命中几率和 NUMA 。在 Go 中，事情变得更高效，因为 Go 调度器试图用更少的线程，每个线程做更多的事情，帮助我们减少系统和硬件层的调度该校。

Go 调度器的设计方面考虑到了操作系统和硬件的复杂情况。把系统层面的 IO/阻塞 操作转换成了 CPU密集 操作来最大化每个 CPU 的能力。这就是为什么你不需要超过虚拟核数的系统线程。你可以让每一个虚拟 Core 上都只跑一个线程来把所有事情做了，这是合理的。对于网络服务及其他不会阻塞系统线程的系统调用的服务来说，可以这样做。
作为一个开发者，你仍然需要理解你的应用要完成的哪一类型的工作。不要无节制的创建 Goroutine 以期望得到惊人的性能。少即是多，但是搞懂 Go 调度器的原理，你可以做出更好的工程决策。下一篇文章，我会探讨利用这些知识，来提升你的服务的性能，同时又能与代码的复杂度上保持一定的平衡。

什么是并发
并发意味着不按顺序执行。给定一组指令，可以按顺序执行，也可以找一种方式不按顺序执行，但仍能得到同样的结果。对于你眼前的问题，显然乱序执行能够增加价值。我所说的价值，意思是以复杂性为成本，而得到足够的性能增益。关键还是要看你的问题，它可能无法或甚至无法进行无序执行。
理解并发并不等于并行也是重要的。并行意味着 2 个或者 2 个以上的指令可以在同一时间一起执行。这和并发的概念不同。只有你的机器有至少 2 个 hardware threads 才能使你的程序达到并行效果。

你看到有两个逻辑处理器（P），每个都分配了一个独立的系统线程（M），线程被关联到了独立的 Core 上。你可以看到 2 个 Goroutine 正在并行执行，他们同时在不同的 Core 上执行各自的指令。对于每一个逻辑处理器，3 个 Groutine 正在轮流共享他们的系统线程。所有的 Goroutine 正在并发执行，没有按照特定的顺序执行他们的指令。
这里有个难点，就是有时候在没有并行能力机器上使用并发，实际上会降低你的性能。更有趣的是，有时候利用并行来达并发效果，并不会如你期望的那样得到更好的性能。
工作负荷
首先，你最好是先搞懂你的问题属于哪种负荷类型。是 CPU 密集型的还是 IO 密集型的。之前的文章有描述，这里不再重复。
CPU 密集型的工作你需要通过并行来达到并发。单个 hardware thread 处理多个 Goroutine 并不高效，因为这些 Goroutine 不会进入阻塞状态。使用比机器 hardware thread 数量更多的 Goroutine 会减慢工作，因为把 Goroutine 不断从系统线程上调来调去也是有成本。上下文切换会触发 Stop The World（简称STW）事件，因为在切换过程中你的工作不会被执行。
对于 IO 密集型的工作，你不需要并行来保证并发。单个 hardware thread 可以处理高效的处理多个 Goroutine，因为这些 Goroutine 会自动进出于阻塞态。拥有比 hardware thread 更多的 Goroutine 可以加速工作的完成，因为让 Goroutine 从系统线程上调来调去不会触发 STW 事件。你的工作会自动暂停，这允许其他不同的 Goroutine 使用同一个 hardware thread 来更高效的完成工作，而不是让它处于空闲状态。
你如何能知道，一个 hardware thread 上跑多少个 Goroutine 才能得到最大程度的吞吐量呢？太少的 Goroutine 导致你有太多的空闲资源，太多的 Goroutine 导致你有太多的延迟时间。这是您需要考虑的

unc add(numbers []int) int {
    var v int
    for _, n := range numbers {
        v += n
    }
    return v
}

问题：这个 add 函数可不可以乱序执行？答案可以的。一个数字集合可以被拆成更小的一些集合，这些小集合是可以并发处理的。一旦小集合求和完了，所有的结果再相加求和，能得到同样的答案。
但是，还有另外一个问题。应该把集合拆成多小，才能让速度最快呢？
为了回答这个问题你就必须知道 add 属于哪种工作。add 方法属于 CPU 密集型，因为算法就是不断执行数学运算，不会导致 Goroutine 进入阻塞态。这意味着每一个 hardware thread 跑一个 Goroutine 会让速度最快。

并发版本的 add。
44 func addConcurrent(goroutines int, numbers []int) int {
45     var v int64
46     totalNumbers := len(numbers)
47     lastGoroutine := goroutines - 1
48     stride := totalNumbers / goroutines
49
50     var wg sync.WaitGroup
51     wg.Add(goroutines)
52
53     for g := 0; g < goroutines; g++ {
54         go func(g int) {
55             start := g * stride
56             end := start + stride
57             if g == lastGoroutine {
58                 end = totalNumbers
59             }
60
61             var lv int
62             for _, n := range numbers[start:end] {
63                 lv += n
64             }
65
66             atomic.AddInt64(&v, int64(lv))
67             wg.Done()
68         }(g)
69     }
70
71     wg.Wait()
72
73     return int(v)
74 }

上面代码中，addConcurrent 函数就是并发版本的 add 方法。解释下这里面比较重要的几行代码。
48行：每个 Goroutine 获得一个独立的但是更小的数字集合进行相加。每个集合的大小是总集合的大小除以 Goroutine 的数量。
53行： 一堆 Goroutine 开始执行加法运算。
57-59行：最后一个 Goroutine 要把剩下的所有数字相加，这有可能使得它的集合要比其他集合大。
66行： 将小集合的求和结果，加在一起得到最终求和结果。
并发版本的实现，要比顺序版本的更复杂，但这到底值不值呢？最好的办法就是做压测。压测时我使用一千万数字的集合，并关闭掉 GC。
func BenchmarkSequential(b *testing.B) {
    for i := 0; i < b.N; i++ {
        add(numbers)
    }
}

func BenchmarkConcurrent(b *testing.B) {
    for i := 0; i < b.N; i++ {
        addConcurrent(runtime.NumCPU(), numbers)
    }
}


上面代码就是压测函数。这里我只分配一个 CPU Core 给程序。顺序版本的使用 1 个 Goroutine，并发版本的使用 runtie.NumCPU() 个（也就是 8 个）Goroutine。这种情况下并发版本因为无法使用并行来进行并发。

读文件
上面举了 2 个 CPU 密集型例子，但是 IO 密集型的呢？在 Goroutine 自然的不断频繁的进出阻塞态时，情况有什么不同？我们看一个 IO 密集的例子，就是读取一些文件然后执行查找操作。
第一个版本是顺序方式实现，函数名是 find
42 func find(topic string, docs []string) int {
43     var found int
44     for _, doc := range docs {
45         items, err := read(doc)
46         if err != nil {
47             continue
48         }
49         for _, item := range items {
50             if strings.Contains(item.Description, topic) {
51                 found++
52             }
53         }
54     }
55     return found
56 }


上面代码就是 find 方法的实现。功能就是从一组字符串中找到
下面是 read 函数的实现
33 func read(doc string) ([]item, error) {
34     time.Sleep(time.Millisecond) // Simulate blocking disk read.
35     var d document
36     if err := xml.Unmarshal([]byte(file), &d); err != nil {
37         return nil, err
38     }
39     return d.Channel.Items, nil
40 }


read 函数在使用了 time.Sleep 把自己挂起 1 毫秒。这个主要是为了模拟 IO 操作延迟。因为你实际去读文件，Goroutine 也是一样会被挂起一段时间。这 1 毫秒的延迟对于后面的压测结果至关重要。
下面是并发版本的实现。
58 func findConcurrent(goroutines int, topic string, docs []string) int {
59     var found int64
60
61     ch := make(chan string, len(docs))
62     for _, doc := range docs {
63         ch <- doc
64     }
65     close(ch)
66
67     var wg sync.WaitGroup
68     wg.Add(goroutines)
69
70     for g := 0; g < goroutines; g++ {
71         go func() {
72             var lFound int64
73             for doc := range ch {
74                 items, err := read(doc)
75                 if err != nil {
76                     continue
77                 }
78                 for _, item := range items {
79                     if strings.Contains(item.Description, topic) {
80                         lFound++
81                     }
82                 }
83             }
84             atomic.AddInt64(&found, lFound)
85             wg.Done()
86         }()
87     }
88
89     wg.Wait()
90
91     return int(found)
92 }


上面代码中 findConcurrent 函数就是 find 函数的并发实现版本。并发版本中使用适量的 Goroutine 完成不定数量的文档查询。代码太多，这里只解释几个重要的地方。
61-64行：创建了一个 Channel，用来发送文档。
65行： 关闭了 channel，这样当所有文档都处理完了后，所有的 Goroutine 就都自动退出了。
70行： Goroutine 被创建
73-83行：每个 Goroutine 从 Channel 中获取文档，把文档读到内存并查找 topic。当找到了，将本地变量 lFound 加一。
84行：每个 Goroutine 都将自己查找到的文档个数，加到全局变量found上。
下面我使用一千个文档进行压测，并关闭 GC。
func BenchmarkSequential(b *testing.B) {
    for i := 0; i < b.N; i++ {
        find("test", docs)
    }
}

func BenchmarkConcurrent(b *testing.B) {
    for i := 0; i < b.N; i++ {
        findConcurrent(runtime.NumCPU(), "test", docs)
    }
}


上面是压测代码。下面是仅使用 1 个 CPU Core 的压测结果。此时并发版本也只能使用一个 CPU Core 来执行 8 个 Goroutine。
10 Thousand Documents using 8 goroutines with 1 core
2.9 GHz Intel 4 Core i7
Concurrency WITHOUT Parallelism
-----------------------------------------------------------------------------
$ GOGC=off go test -cpu 1 -run none -bench . -benchtime 3s
goos: darwin
goarch: amd64
pkg: github.com/ardanlabs/gotraining/topics/go/testing/benchmarks/io-bound
BenchmarkSequential                3    1483458120 ns/op
BenchmarkConcurrent               20     188941855 ns/op : ~87% Faster
BenchmarkSequentialAgain           2    1502682536 ns/op
BenchmarkConcurrentAgain          20     184037843 ns/op : ~88% Faster


发现在单核的情况下。并发版本要比顺序版本快 87%~88%。这和我预期一致。因为所有的 Goroutine 共用一个 CPU Core，而 Goroutine 在调用 read 时候会自动进入阻塞态，这时会将 CPU Core 出让给其他 Goroutine 使用，这使得这个 CPU Core 更加繁忙。
下面看下使用多核进行压测。
10 Thousand Documents using 8 goroutines with 1 core
2.9 GHz Intel 4 Core i7
Concurrency WITH Parallelism
-----------------------------------------------------------------------------
$ GOGC=off go test -run none -bench . -benchtime 3s
goos: darwin
goarch: amd64
pkg: github.com/ardanlabs/gotraining/topics/go/testing/benchmarks/io-bound
BenchmarkSequential-8                  3    1490947198 ns/op
BenchmarkConcurrent-8                 20     187382200 ns/op : ~88% Faster
BenchmarkSequentialAgain-8             3    1416126029 ns/op
BenchmarkConcurrentAgain-8            20     185965460 ns/op : ~87% Faster


这个压测结果说明，更多的 CPU Core 并不会对性能有多大的提升。