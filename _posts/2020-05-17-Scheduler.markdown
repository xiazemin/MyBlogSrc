---
title: Scheduler 原理解析
layout: post
category: golang
author: 夏泽民
---
Golang Scheduler原理解析
Section1 Scheduler原理
1.基础知识
2.调度模型
3.调度核心问题
Section2 主要模型的源码分析
2.1 实体M
2.2 实体P(processor)
2.3 实体G(goroutine)
Section3 主要调度流程的源码分析
3.1 预备知识
3.1.1 golang的函数调用规范
3.1.2 TLS(thread local storage)
3.1.3 栈扩张
3.1.4 写屏障(write barrier)
3.1.5 m0和g0
3.1.6 go中线程的种类
3.2 main线程启动执行
3.3 新建goroutine过程
3.4 循环调度schedule过程
3.5 抢占式调度实现(sysmon线程)
Section4：scheduler与memory allocation、channel、garbage collection关联部分
本文主要分析Golang底层对于协程的调度原理，本文与Golang的memory allocation、garbage collection这两个主题是紧密相关的，本文scheduler作为系列的第一篇文章。
文章大体上的思路是这样的：
section1：主要图示和文字介绍scheduler的原理；
section2：主要模型的角度介绍scheduler原理；
section3：从主要调度流程介绍scheduler原理；
section4：分析scheduler与memory allocation、channel、garbage collection关联部分
基于源码 Go SDK 1.11

Section1 Scheduler原理
1.基础知识
Golang支持语言级别的并发，并发的最小逻辑单位叫做goroutine，goroutine就是Go为了实现并发提供的用户态线程，这种用户态线程是运行在内核态线程(OS线程)之上。当我们创建了大量的goroutine并且同时运行在一个或则多个内核态线程上时(内核线程与goroutine是m:n的对应关系)，就需要一个调度器来维护管理这些goroutine，确保所有的goroutine都有相对公平的机会使用CPU。

这里再次强调一次，goroutine与内核OS线程的映射关系是M:N，这样多个goroutine就可以在多个内核线程上面运行。goroutine的切换大部分场景下都没有走OS线程的切换所带来的开销，这样整体运行效率相比OS线程的调度会高很多，但是这样带来的问题就是goroutine调度模型的复杂。
https://blog.csdn.net/u010853261/article/details/84790392
<!-- more -->
2.调度模型
Golang的调度模型主要有几个主要的实体：G、M、P、schedt。

G：代表一个goroutine实体，它有自己的栈内存，instruction pointer和一些相关信息(比如等待的channel等等)，是用于调度器调度的实体。
M：代表一个真正的内核OS线程，和POSIX里的thread差不多，属于真正执行指令的人。
P：代表M调度的上下文，可以把它看做一个局部的调度器，调度协程go代码在一个内核线程上跑。P是实现协程与内核线程的N:M映射关系的关键。P的上限是通过系统变量runtime.GOMAXPROCS (numLogicalProcessors)来控制的。golang启动时更新这个值，一般不建议修改这个值。P的数量也代表了golang代码执行的并发度，即有多少goroutine可以并行的运行。
schedt：runtime全局调度时使用的数据结构，这个实体其实只是一个壳，里面主要有M的全局idle队列，P的全局idle队列，一个全局的就绪的G队列以及一个runtime全局调度器级别的锁。当对M或P等做一些非局部调度器的操作时，一般需要先锁住全局调度器。
为了解释清楚这几个实体之间的关系，我们先抽象G、M、P、schedt的关系，主要的workflow如下图所示：


从上图我们可以分析出几个结论：

我们通过 go func()来创建一个goroutine；
有两个存储goroutine的队列，一个是局部调度器P的local queue、一个是全局调度器数据模型schedt的global queue。新创建的goroutine会先保存在local queue，如果local queue已经满了就会保存在全局的global queue；
goroutine只能运行在M中，一个M必须持有一个P，M与P是1：1的关系。M会从P的local queue弹出一个Runable状态的goroutine来执行，如果P的local queue为空，就会执行work stealing；
一个M调度goroutine执行的过程是一个loop；
当M执行某一个goroutine时候如果发生了syscall或则其余阻塞操作，M会阻塞，如果当前有一些G在执行，runtime会把这个线程M从P中摘除(detach)，然后再创建一个新的操作系统的线程(如果有空闲的线程可用就复用空闲线程)来服务于这个P；
当M系统调用结束时候，这个goroutine会尝试获取一个空闲的P执行，并放入到这个P的local queue。如果获取不到P，那么这个线程M会park它自己(休眠)， 加入到空闲线程中，然后这个goroutine会被放入schedt的global queue。
Go运行时会在下面的goroutine被阻塞的情况下运行另外一个goroutine：

syscall、
network input、
channel operations、
primitives in the sync package。
3.调度核心问题
前面已经大致介绍了scheduler的一些核心调度原理，介绍的都是比较抽象的内容。听下来还有几个疑问需要分析，主要通过后面的源码来进行细致分析。

Question1：如果一个goroutine一直占有CPU又不会有阻塞或则主动让出CPU的调度，scheduler怎么做抢占式调度让出CPU？
Answer1：有一个sysmon线程做抢占式调度，当一个goroutine占用CPU超过10ms之后，调度器会根据实际情况提供不保证的协程切换机制，具体细节见后面源码分析。

Question2：我们知道P的上限是系统启动时候设定的并且一般不会更改，那么内核线程M上限是多少？遇到需要新的M时候是选取IDEL的M还是创建新的M，整体策略是怎样的？
Answer2：在golang系统启动时候会设置内核线程上限是10000，这里先解释一下，M的数量与P的数量和G的数量没有直接关系，实际情况要看调度器的执行情况。
至于具体的策略见后面的源码分析。

Question3：P、M、G的状态机运转，主要是协程对象G。
Answer3：状态机见后面的源码分析

Question4：每一个协程goroutine的栈空间是保存在哪里的？P、M、G分别维护的与内存有关的数据有哪些？
Answer4：golang scheduler use thrid-level cache，each goroutine stack space is applied in heap。

Question5：当syscall、网络IO、channel时，如果这些阻塞返回了，对应的G会被保存在哪个地方？global Queue或则是local queue? 为什么？系统初始化时候的G会被保存在哪里？为什么？
Answer5：唤醒的M首先会尝试获取一个空闲P，然后将G放到P的local queue，如果获取失败，就放进 global queue，然后M自己放进the M idle 列表。

Notice：scheduler的源码主要在两个文件：

runtime/runtime2.go 主要是实体G、M、P、schedt的数据模型
runtime/proc.go 主要是调度实现的逻辑部分。
Section2 主要模型的源码分析
这一部分主要结合源码进行分析，主要分为两个部分，一部分介绍主要模型G、M、P、schedt的职责、维护的数据域以及它们的联系。

2.1 实体M
实体M在模型上等同于系统内核OS线程，M运行的go代码类型有两种：

goroutine代码, M运行go代码需要一个实体P进行局部调度；
原生代码, 例如阻塞的syscall, M运行原生代码不需要P
M会从runqueue(local or global)中抽取G并运行，如果G运行完毕或则G进入了睡眠态，就会从runqueue中取出下一个runnable状态的G运行, 循环调度。

G有时会执行一些阻塞调用(syscall)，这时M会释放持有的P并进入阻塞态，其他的M会取得这个idel状态的P并继续运行队列中的G。Golang需要保证有足够的M可以运行G, 不让CPU闲着, 也需要保证M的数量不能过多。通常创建一个M的原因是由于没有足够的M来关联P并运行其中可运行的G。而且运行时系统执行系统监控的时候，或者GC的时候也会创建M。

M结构体定义在runtime2.go如下：

type m struct {
	/*
        1.  所有调用栈的Goroutine,这是一个比较特殊的Goroutine。
        2.  普通的Goroutine栈是在Heap分配的可增长的stack,而g0的stack是M对应的线程栈。
        3.  所有与调度相关的代码,都会先切换到g0的栈再执行。
    */
	g0      *g     // goroutine with scheduling stack
	morebuf gobuf  // gobuf arg to morestack
	divmod  uint32 // div/mod denominator for arm - known to liblink

	// Fields not known to debuggers.
	procid        uint64       // for debuggers, but offset not hard-coded
	gsignal       *g           // signal-handling g
	goSigStack    gsignalStack // Go-allocated signal handling stack
	sigmask       sigset       // storage for saved signal mask
	tls           [6]uintptr   // thread-local storage (for x86 extern register)
	// 表示M的起始函数。其实就是我们 go 语句携带的那个函数。
	mstartfn      func()
	// M中当前运行的goroutine
	curg          *g       // current running goroutine
	caughtsig     guintptr // goroutine running during fatal signal
	// 与M绑定的P，如果为nil表示空闲
	p             puintptr // attached p for executing go code (nil if not executing go code)
	// 用于暂存于当前M有潜在关联的P。 （预联）当M重新启动时，即用预联的这个P做关联啦
	nextp         puintptr
	id            int64
	mallocing     int32
	throwing      int32
	// 当前m是否关闭抢占式调度
	preemptoff    string // if != "", keep curg running on this m
	/**  */
	locks         int32
	dying         int32
	profilehz     int32
	// 不为0表示此m在做帮忙gc。helpgc等于n只是一个编号
	helpgc        int32
	// 自旋状态，表示当前M是否正在自旋寻找G。在寻找过程中M处于自旋状态。
	spinning      bool // m is out of work and is actively looking for work
	blocked       bool // m is blocked on a note
	inwb          bool // m is executing a write barrier
	newSigstack   bool // minit on C thread called sigaltstack
	printlock     int8
	incgo         bool   // m is executing a cgo call
	freeWait      uint32 // if == 0, safe to free g0 and delete m (atomic)
	fastrand      [2]uint32
	needextram    bool
	traceback     uint8
	ncgocall      uint64      // number of cgo calls in total
	ncgo          int32       // number of cgo calls currently in progress
	cgoCallersUse uint32      // if non-zero, cgoCallers in use temporarily
	cgoCallers    *cgoCallers // cgo traceback if crashing in cgo call
	park          note
	//这个域用于链接allm
	alllink       *m // on allm
	schedlink     muintptr
	mcache        *mcache
	// 表示与当前M锁定的那个G。运行时系统会把 一个M 和一个G锁定，一旦锁定就只能双方相互作用，不接受第三者。
	lockedg       guintptr
	createstack   [32]uintptr    // stack that created this thread.
	lockedExt     uint32         // tracking for external LockOSThread
	lockedInt     uint32         // tracking for internal lockOSThread
	nextwaitm     muintptr       // next m waiting for lock
	waitunlockf   unsafe.Pointer // todo go func(*g, unsafe.pointer) bool
	waitlock      unsafe.Pointer
	waittraceev   byte
	waittraceskip int
	startingtrace bool
	syscalltick   uint32
	thread        uintptr // thread handle
	freelink      *m      // on sched.freem

	// these are here because they are too large to be on the stack
	// of low-level NOSPLIT functions.
	libcall   libcall
	libcallpc uintptr // for cpu profiler
	libcallsp uintptr
	libcallg  guintptr
	syscall   libcall // stores syscall parameters on windows

	vdsoSP uintptr // SP for traceback while in VDSO call (0 if not in call)
	vdsoPC uintptr // PC for traceback while in VDSO call

	mOS
}
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
24
25
26
27
28
29
30
31
32
33
34
35
36
37
38
39
40
41
42
43
44
45
46
47
48
49
50
51
52
53
54
55
56
57
58
59
60
61
62
63
64
65
66
67
68
69
70
71
72
73
74
75
76
77
78
79
80
81
82
83
84
上面字段很多，核心的主要是以下几个字段：

g0      *g     // goroutine with scheduling stack，也是运行局部调度器的g
mstartfn      func()
curg          *g       // current running goroutine
p             puintptr // attached p for executing go code (nil if not executing go code)
nextp         puintptr
helpgc        int32
spinning      bool // m is out of work and is actively looking for work
alllink       *m // on allm
lockedg       guintptr
1
2
3
4
5
6
7
8
9
这些字段主要功能如下：

g0: Golang runtime系统在线程创建的时候创建的，g0的栈使用的是内核线程的栈，主要用于局部调度器执行调度逻辑时使用的栈，也就是执行调度逻辑时的线程栈。
mstartfn：表示M的起始函数。其实就是我们 go 关键字后面携带的那个函数。
curg：存放当前正在运行的G的指针。
p：指向当前与M关联的那个P。
nextp：用于暂存于当前M有潜在关联的P。 （预联）当M重新启动时，即用预联的这个P做关联啦
spinning：自旋状态标志位，表示当前M是否正在寻找G。
alllink：连接到所有的m链表的一个指针。
lockedg：表示与当前M锁定的那个G。运行时系统会把 一个M 和一个G锁定，一旦锁定就只能双方相互作用，不接受第三者。
M的状态机比较简单，因为M是golang对内核OS线程的更上一层抽象，所以M也没有专门字段来维护状态，简单来说有一下几种状态：

自旋中(spinning): M正在从运行队列获取G, 这时候M会拥有一个P；
执行go代码中: M正在执行go代码, 这时候M会拥有一个P；
执行原生代码中: M正在执行原生代码或者阻塞的syscall, 这时M并不拥有P；
休眠中: M发现无待运行的G时会进入休眠，并添加到空闲M链表中, 这时M并不拥有P。
上面的几种状态中，spinning这个状态非常重要，是否需要唤醒或者创建新的M取决于当前自旋中的M的数量。

M在被创建之初会被加入到全局的M列表 【runtime.allm】。 接着，M的起始函数（mstartfn）和准备关联的P都会被设置。最后，runtime会为M专门创建一个新的内核线程并与之关联。这时候这个新的M就为执行G做好了准备。其中起始函数（mstartfn）仅当runtime要用此M执行系统监控或者垃圾回收等任务的时候才会被设置。【runtime.allm】的作用是runtime在需要的时候会通过它获取到所有的M的信息，同时防止M被gc。

在新的M被创建后会做一些初始化工作。其中包括了对自身所持的栈空间以及信号的初始化。在上述初始化完成后 mstartfn 函数就会被执行 (如果存在的话)。【注意】：如果mstartfn 代表的是系统监控任务的话，那么该M会一直在执行mstartfn 而不会有后续的流程。否则 mstartfn 执行完后，当前M将会与那个准备与之关联的P完成关联。至此，一个并发执行环境才真正完成。之后就是M开始寻找可运行的G并运行之。

runtime管辖的M会在GC任务执行的时候被停止，这时候系统会对M的属性做某些必要的重置并把M放置入全局调度器的空闲M列表。【很重要】因为调度器在需要一个未被使用的M时，运行时系统会先去这个空闲列表获取M。(只有都没有的时候才会创建M)

M本身是无状态的。M是否是空闲态仅根据它是否存在于调度器的空闲M列表 【runtime.sched.midle】 中来判定（注意：空闲列表不是那个全局列表）。

单个Go程序所使用的M的最大数量是可以被设置的。在我们使用命令运行Go程序时候，有一个引导程序先会被启动的。在这个引导程序中会为Go程序的运行建立必要的环境。引导程序对M的数量进行初始化设置，默认最大值是10000【一个Go程序最多可以使用10000个M，即：理想状态下，可以同时有1W个内核线程被同时运行】。可以使用 runtime/debug.SetMaxThreads() 函数设置。

2.2 实体P(processor)
P是一个抽象的概念，并不代表一个具体的实体，抽象地表示M运行G所需要的资源。P并不代表CPU核心数，而是表示执行go代码的并发度。有一点需要注意的是，执行原生代码的时候并不受P数量的限制。

同一时间只有一个线程(M)可以拥有P， 局部调度器P维护的数据都是锁自由(lock free)的, 读写这些数据的效率会非常的高。

P是使G能够在M中运行的关键。Go的runtime适当地让P与不同的M建立或者断开联系，以使得P中的那些可运行的G能够在需要的时候及时获得运行时机。

P结构体定义在runtime2.go如下：

type p struct {
	lock mutex

	id          int32
	// 当前p的状态
	status      uint32 // one of pidle/prunning/...
	// 链接
	link        puintptr
	schedtick   uint32     // incremented on every scheduler call
	syscalltick uint32     // incremented on every system call
	sysmontick  sysmontick // last tick observed by sysmon
	// p反向链接到关联的m(空闲时为nil)
	m           muintptr   // back-link to associated m (nil if idle)
	mcache      *mcache
	racectx     uintptr

	deferpool    [5][]*_defer // pool of available defer structs of different sizes (see panic.go)
	deferpoolbuf [5][32]*_defer

	// Cache of goroutine ids, amortizes accesses to runtime·sched.goidgen.
	goidcache    uint64
	goidcacheend uint64

	// Queue of runnable goroutines. Accessed without lock.
	runqhead uint32
	runqtail uint32
	runq     [256]guintptr
	// runnext, if non-nil, is a runnable G that was ready'd by
	// the current G and should be run next instead of what's in
	// runq if there's time remaining in the running G's time
	// slice. It will inherit the time left in the current time
	// slice. If a set of goroutines is locked in a
	// communicate-and-wait pattern, this schedules that set as a
	// unit and eliminates the (potentially large) scheduling
	// latency that otherwise arises from adding the ready'd
	// goroutines to the end of the run queue.
	runnext guintptr

	// Available G's (status == Gdead)
	gfree    *g
	gfreecnt int32

	sudogcache []*sudog
	sudogbuf   [128]*sudog

	tracebuf traceBufPtr

	// traceSweep indicates the sweep events should be traced.
	// This is used to defer the sweep start event until a span
	// has actually been swept.
	traceSweep bool
	// traceSwept and traceReclaimed track the number of bytes
	// swept and reclaimed by sweeping in the current sweep loop.
	traceSwept, traceReclaimed uintptr

	palloc persistentAlloc // per-P to avoid mutex

	// Per-P GC state
	gcAssistTime         int64 // Nanoseconds in assistAlloc
	gcFractionalMarkTime int64 // Nanoseconds in fractional mark worker
	gcBgMarkWorker       guintptr
	gcMarkWorkerMode     gcMarkWorkerMode

	// gcMarkWorkerStartTime is the nanotime() at which this mark
	// worker started.
	gcMarkWorkerStartTime int64

	// gcw is this P's GC work buffer cache. The work buffer is
	// filled by write barriers, drained by mutator assists, and
	// disposed on certain GC state transitions.
	gcw gcWork

	// wbBuf is this P's GC write barrier buffer.
	//
	// TODO: Consider caching this in the running G.
	wbBuf wbBuf

	runSafePointFn uint32 // if 1, run sched.safePointFn at next safe point

	pad [sys.CacheLineSize]byte
}
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
24
25
26
27
28
29
30
31
32
33
34
35
36
37
38
39
40
41
42
43
44
45
46
47
48
49
50
51
52
53
54
55
56
57
58
59
60
61
62
63
64
65
66
67
68
69
70
71
72
73
74
75
76
77
78
79
80
81
这些字段里面比较核心的字段如下：

lock mutex
status      uint32 // one of pidle/prunning/...
link        puintptr
m           muintptr   // back-link to associated m (nil if idle)
// local runnable queue. Accessed without lock. implement through array
runqhead uint32
runqtail uint32
runq     [256]guintptr
runnext guintptr
// Available G's (status == Gdead)
gfree    *g
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
通过runtime.GOMAXPROCS函数我们可以改变单个Go程序可以拥有P的最大数量，如果不做设置会有一个默认值。

每一个P都必须关联一个M才能使其中的G得以运行。

【注意】：runtime会将M与关联的P分离开来。但是如果该P的runqueue中还有未运行的G，那么runtime就会找到一个空的M（在调度器的空闲队列中的M） 或者创建一个空的M，并与该P关联起来（为了运行G而做准备）。

runtime.GOMAXPROCS只能够设置P的数量，并不会影响到M（内核线程）数量，所以runtime.GOMAXPROCS不是控制线程数，只能影响局部调度器P的数量。

在runtime初始化时会确认P的最大数量，之后会根据这个最大值初始化全局P列表【runtime.allp】。类似全局M列表，【runtime.allp】包含了runtime创建的所有P。随后，runtime会把调度器的可运行G队列【runtime.schedt.runq】中的所有G均匀的放入全局的P列表中的各个P的可执行G队列 local queue中。到这里为止，runtime需要用到的所有P都准备就绪了。

类似M的空闲列表，调度器也存在一个空闲P的列表【runtime.shcedt.pidle】，当一个P不再与任何M关联的时候，runtime会把该P放入这个列表，而一个空闲的P关联了某个M之后会被从【runtime.shcedt.pidle】中取出来。【注意：一个P加入了空闲列表，其G的可运行local queue也不一定为空】。

和M不同，P是有状态机的（五种）：

Pidel：当前P未和任何M关联
Prunning：当前P已经和某个M关联，M在执行某个G
Psyscall：当前P中的被运行的那个G正在进行系统调用
Pgcstop：runtime正在进行GC（runtime会在gc时试图把全局P列表中的P都处于此种状态）
Pdead：当前P已经不再被使用（在调用runtime.GOMAXPROCS减少P的数量时，多余的P就处于此状态）
在对P初始化的时候就是Pgcstop的状态，但是这个状态保持时间很短，在初始化并填充P中的G队列之后，runtime会将其状态置为Pidle并放入调度器的空闲P列表【runtime.schedt.pidle】中，其中的P会由调度器根据实际情况进行取用。具体的状态机流转图如下图所示：



从上图我们可以看到，除了Pdead状态以外的其余状态，在runtime进行GC的时候，P都会被指定成Pgcstop。在GC结束后状态不会回复到GC前的状态，而是都统一直接转到了Pidle 【这意味着，他们都需要被重新调度】。

【注意】除了Pgcstop 状态的P，其他状态的P都会在调用runtime.GOMAXPROCS 函数减少P数目时，被认为是多余的P而状态转为Pdead，这时候其带的可运行G的队列中的G都会被转移到调度器的可运行G队列中，它的自由G队列 【gfree】也是一样被移到调度器的自由列表【runtime.sched.gfree】中。

【注意】每个P中都有一个可运行G队列及自由G队列。自由G队列包含了很多已经完成的G，随着被运行完成的G的积攒到一定程度后，runtime会把其中的部分G转移到全局调度器的自由G队列 【runtime.sched.gfree】中。

【注意】当我们每次用 go关键字启用一个G的时候，首先都是尽可能复用已经执行完的G。具体过程如下：运行时系统都会先从P的自由G队列获取一个G来封装我们提供的函数 (go 关键字后面的函数) ，如果发现P中的自由G过少时，会从调度器的自由G队列中移一些G过来，只有连调度器的自由G列表都弹尽粮绝的时候，才会去创建新的G。

2.3 实体G(goroutine)
goroutine可以理解成被调度器管理的轻量级线程，goroutine使用go关键字创建。

goroutine的新建, 休眠, 恢复, 停止都受到go的runtime管理。
goroutine执行异步操作时会进入休眠状态, 待操作完成后再恢复, 无需占用系统线程。
goroutine新建或恢复时会添加到运行队列, 等待M取出并运行。

g和gobuf的结构定义和在runtime2.go如下：

type g struct {
	// Stack parameters.
	// stack describes the actual stack memory: [stack.lo, stack.hi).
	// stackguard0 is the stack pointer compared in the Go stack growth prologue.
	// It is stack.lo+StackGuard normally, but can be StackPreempt to trigger a preemption.
	// stackguard1 is the stack pointer compared in the C stack growth prologue.
	// It is stack.lo+StackGuard on g0 and gsignal stacks.
	// It is ~0 on other goroutine stacks, to trigger a call to morestackc (and crash).
	stack       stack   // offset known to runtime/cgo
	stackguard0 uintptr // offset known to liblink
	stackguard1 uintptr // offset known to liblink

	_panic         *_panic // innermost panic - offset known to liblink
	_defer         *_defer // innermost defer
	/**
	 *	有一个指针指向执行它的m，也即g隶属于m；
	 */
	m              *m      // current m; offset known to arm liblink
	// 进程切换时，利用sched域来保存上下文
	sched          gobuf
	syscallsp      uintptr        // if status==Gsyscall, syscallsp = sched.sp to use during gc
	syscallpc      uintptr        // if status==Gsyscall, syscallpc = sched.pc to use during gc
	stktopsp       uintptr        // expected sp at top of stack, to check in traceback
	param          unsafe.Pointer // passed parameter on wakeup
	// 状态Gidle,Grunnable,Grunning,Gsyscall,Gwaiting,Gdead
	atomicstatus   uint32
	stackLock      uint32 // sigprof/scang lock; TODO: fold in to atomicstatus
	goid           int64
	//????
	schedlink      guintptr
	waitsince      int64      // approx time when the g become blocked
	waitreason     waitReason // if status==Gwaiting
	// 抢占标志
	preempt        bool       // preemption signal, duplicates stackguard0 = stackpreempt
	paniconfault   bool       // panic (instead of crash) on unexpected fault address
	preemptscan    bool       // preempted g does scan for gc
	gcscandone     bool       // g has scanned stack; protected by _Gscan bit in status
	gcscanvalid    bool       // false at start of gc cycle, true if G has not run since last scan; TODO: remove?
	throwsplit     bool       // must not split stack
	raceignore     int8       // ignore race detection events
	sysblocktraced bool       // StartTrace has emitted EvGoInSyscall about this goroutine
	sysexitticks   int64      // cputicks when syscall has returned (for tracing)
	traceseq       uint64     // trace event sequencer
	tracelastp     puintptr   // last P emitted an event for this goroutine
	// G被锁定只能在这个m上运行
	lockedm        muintptr
	sig            uint32
	writebuf       []byte
	sigcode0       uintptr
	sigcode1       uintptr
	sigpc          uintptr
	// 创建这个goroutine的go表达式的pc
	gopc           uintptr         // pc of go statement that created this goroutine
	ancestors      *[]ancestorInfo // ancestor information goroutine(s) that created this goroutine (only used if debug.tracebackancestors)
	startpc        uintptr         // pc of goroutine function
	racectx        uintptr
	waiting        *sudog         // sudog structures this g is waiting on (that have a valid elem ptr); in lock order
	cgoCtxt        []uintptr      // cgo traceback context
	labels         unsafe.Pointer // profiler labels
	timer          *timer         // cached timer for time.Sleep
	selectDone     uint32         // are we participating in a select and did someone win the race?

	// Per-G GC state

	// gcAssistBytes is this G's GC assist credit in terms of
	// bytes allocated. If this is positive, then the G has credit
	// to allocate gcAssistBytes bytes without assisting. If this
	// is negative, then the G must correct this by performing
	// scan work. We track this in bytes to make it fast to update
	// and check for debt in the malloc hot path. The assist ratio
	// determines how this corresponds to scan work debt.
	gcAssistBytes int64
}

//用于保存G切换时上下文的缓存结构体
type gobuf struct {
	// The offsets of sp, pc, and g are known to (hard-coded in) libmach.
	//
	// ctxt is unusual with respect to GC: it may be a
	// heap-allocated funcval, so GC needs to track it, but it
	// needs to be set and cleared from assembly, where it's
	// difficult to have write barriers. However, ctxt is really a
	// saved, live register, and we only ever exchange it between
	// the real register and the gobuf. Hence, we treat it as a
	// root during stack scanning, which means assembly that saves
	// and restores it doesn't need write barriers. It's still
	// typed as a pointer so that any other writes from Go get
	// write barriers.
	sp   uintptr //当前的栈指针
	pc   uintptr //当前的计数器
	g    guintptr //g自身引用
	ctxt unsafe.Pointer
	ret  sys.Uintreg
	lr   uintptr
	bp   uintptr // for GOEXPERIMENT=framepointer
}
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
24
25
26
27
28
29
30
31
32
33
34
35
36
37
38
39
40
41
42
43
44
45
46
47
48
49
50
51
52
53
54
55
56
57
58
59
60
61
62
63
64
65
66
67
68
69
70
71
72
73
74
75
76
77
78
79
80
81
82
83
84
85
86
87
88
89
90
91
92
93
94
95
96
Go语言的编译器会把我们编写的goroutine编译为runtime的函数调用，并把go语句中的函数以及其参数传递给runtime的函数中。

runtime在接到这样一个调用后，会先检查一下go函数及其参数的合法性，紧接着会试图从局部调度器P的自由G队列中(或者全局调度器的自由G队列)中获取一个可用的自由G （P中有讲述了），如果没有则新创建一个G。类似M和P，G在运行时系统中也有全局的G列表【runtime.allg】，那些新建的G会先放到这个全局的G列表中，其列表的作用也是集中放置了当前运行时系统中给所有的G的指针。在用自由G封装go的函数时，运行时系统都会对这个G重新做一次初始化。

初始化：包含了被关联的go关键字后的函数及当前G的状态机G的ID等等。在G被初始化完成后就会被放置到当前本地的P的可运行队列中。只要时机成熟，调度器会立即尽心这个G的调度运行。

G的状态机会比较复杂一点，大致上和内核线程的状态机有一点类似，但是状态机流转有一些区别。G的各种状态如下：

Gidle：G被创建但还未完全被初始化。
Grunnable：当前G为可运行的，正在等待被运行。
Grunning：当前G正在被运行。
Gsyscall：当前G正在被系统调用
Gwaiting：当前G正在因某个原因而等待
Gdead：当前G完成了运行
初始化完的G是处于Grunnable的状态，一个G真正在M中运行时是处于Grunning的状态，G的状态机流转图如下图所示：


上图有一步是等待的事件到来，那么G在运行过程中，是否等待某个事件以及等待什么样的事件？完全由起封装的go关键字后的函数决定。（如：等待chan中的值、涉及网络I/O、time.Timer、time.Sleep等等事件）

G退出系统调用的过程非常复杂：runtime先会尝试获取空闲局部调度器P并直接运行当前G，如果没有就会把当前G转成Grunnable状态并放置入全局调度器的global queue。

最后，已经是Gdead状态的G是可以被重新初始化并使用的(从自由G队列取出来重新初始化使用)。而对比进入Pdead状态的P等待的命运只有被销毁。处于Gdead的G会被放置到本地P或者调度器的自由G列表中。

至此，G、M、P的初步描述已经完毕，下面我们来看一看一些核心的队列：

中文名	源码名称	作用域	简要说明
全局M列表	runtime.allm	运行时系统	存放所有M
全局P列表	runtime.allp	运行时系统	存放所有P
全局G列表	runtime.allg	运行时系统	存放所有G
调度器中的空闲M列表	runtime.schedt.midle	调度器	存放空闲M，链表结构
调度器中的空闲P列表	runtime.schedt.pidle	调度器	存放空闲P，链表结构
调度器中的可运行G队列	runtime.schedt.runq	调度器	存放可运行G，链表结构
调度器中的自由G列表	runtime.schedt.gfree	调度器	存放自由G， 链表结构
P中的可运行G队列	runq	本地P	存放当前P中的可运行G，环形队列，数组实现
P中的自由G列表	gfree	本地P	存放当前P中的自由G，链表结构
三个全局的列表主要为了统计runtime的所有G、M、P。我们主要关心剩下的这些容器，尤其是和G相关的四个。

在runtime创建的G都会被保存在全局的G列表中，值得注意的是：

从Gsyscall转出来的G，如果不能马上获取空闲的P执行，就会被放置到全局调度器的可运行队列中(global queue)。
被runtime初始化的G会被放置到本地P的可运行队列中(local queue)
从Gwaiting转出来的G，除了因网络IO陷入等待的G之外，都会被防止到本地P可运行的G队列中。
转成Gdead状态的G会先被放置在本地P的自由G列表。
调度器中的与G、M、P相关的列表其实只是起了一个暂存的作用。
一句话概括三者关系：

G需要绑定在M上才能运行
M需要绑定P才能运行
这三者之间的实体关系是：


内核调度实体(Kernel Scheduling Entry)与三者的关系是：



可知：一个G的执行需要M和P的支持。一个M在于一个P关联之后就形成一个有效的G运行环境 【内核线程 + 上下文环境】。每个P都含有一个 可运行G的队列【runq】。队列中的G会被一次传递给本地P关联的M并且获得运行时机。

M 与 KSE 的关系是绝对的一对一，一个M仅能代表一个内核线程。在一个M的生命周期内，仅会和一个内核KSE产生关联。M与P以及P与G之间的关联时多变的，总是会随着调度器的实际调度策略而变化。

这里我们再回顾下G、M、P里面核心成员

G里面的核心成员

stack ：当前g使用的栈空间, 有lo和hi两个成员
stackguard0 ：检查栈空间是否足够的值, 低于这个值会扩张栈, 0是go代码使用的
stackguard1 ：检查栈空间是否足够的值, 低于这个值会扩张栈, 1是原生代码使用的
m ：当前g对应的m
sched ：g的调度数据, 当g中断时会保存当前的pc和rsp等值到这里, 恢复运行时会使用这里的值
atomicstatus: g的当前状态
schedlink: 下一个g, 当g在链表结构中会使用
preempt: g是否被抢占中
lockedm: g是否要求要回到这个M执行, 有的时候g中断了恢复会要求使用原来的M执行
M里面的核心成员

g0: 用于调度的特殊g, 调度和执行系统调用时会切换到这个g
curg: 当前运行的g
p: 当前拥有的P
nextp: 唤醒M时, M会拥有这个P
park: M休眠时使用的信号量, 唤醒M时会通过它唤醒
schedlink: 下一个m, 当m在链表结构中会使用
mcache: 分配内存时使用的本地分配器, 和p.mcache一样(拥有P时会复制过来)
lockedg: lockedm的对应值
P里面的核心成员

status: p的当前状态
link: 下一个p, 当p在链表结构中会使用
m: 拥有这个P的M
mcache: 分配内存时使用的本地分配器
runqhead: 本地运行队列的出队序号
runqtail: 本地运行队列的入队序号
runq: 本地运行队列的数组, 可以保存256个G
gfree: G的自由列表, 保存变为_Gdead后可以复用的G实例
gcBgMarkWorker: 后台GC的worker函数, 如果它存在M会优先执行它
gcw: GC的本地工作队列, 详细将在下一篇(GC篇)分析
调度器除了设计上面的三个结构体，还有一个全局调度器数据结构schedt:

type schedt struct {
	// accessed atomically. keep at top to ensure alignment on 32-bit systems.
	// // 下面两个变量需以原子访问访问。保持在 struct 顶部，确保其在 32 位系统上可以对齐
	goidgen  uint64
	lastpoll uint64

	lock mutex

	// When increasing nmidle, nmidlelocked, nmsys, or nmfreed, be
	// sure to call checkdead().
	//=====与m数量相关的变量================================================
	// 空闲m列表指针
	midle        muintptr // idle m's waiting for work
	// 空闲m的数量
	nmidle       int32    // number of idle m's waiting for work
	// 被锁住的m空闲数量
	nmidlelocked int32    // number of locked m's waiting for work
	// 已经创建的m的数目和下一个m ID
	mnext        int64    // number of m's that have been created and next M ID
	// 允许创建的m的最大数量
	maxmcount    int32    // maximum number of m's allowed (or die)
	// 不计入死锁的m的数量
	nmsys        int32    // number of system m's not counted for deadlock
	// 释放m的累计数量
	nmfreed      int64    // cumulative number of freed m's

	//系统的goroutine的数量
	ngsys uint32 // number of system goroutines; updated atomically

	//=====与p数量相关的变量================================================
	// 空闲的p列表
	pidle      puintptr // idle p's
	// 空闲p的数量
	npidle     uint32
	//
	nmspinning uint32 // See "Worker thread parking/unparking" comment in proc.go.

	// Global runnable queue.
	// 全局runable g链表的head地址
	runqhead guintptr
	// 全局runable g链表的tail地址
	runqtail guintptr
	// 全局runable g链表的大小
	runqsize int32

	// Global cache of dead G's.
	gflock       mutex
	gfreeStack   *g
	gfreeNoStack *g
	ngfree       int32

	// Central cache of sudog structs.
	sudoglock  mutex
	sudogcache *sudog

	// Central pool of available defer structs of different sizes.
	deferlock mutex
	deferpool [5]*_defer

	// freem is the list of m's waiting to be freed when their
	// m.exited is set. Linked through m.freelink.
	freem *m

	gcwaiting  uint32 // gc is waiting to run
	stopwait   int32
	stopnote   note
	sysmonwait uint32
	sysmonnote note

	// safepointFn should be called on each P at the next GC
	// safepoint if p.runSafePointFn is set.
	safePointFn   func(*p)
	safePointWait int32
	safePointNote note

	profilehz int32 // cpu profiling rate

	procresizetime int64 // nanotime() of last change to gomaxprocs
	totaltime      int64 // ∫gomaxprocs dt up to procresizetime
}
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
24
25
26
27
28
29
30
31
32
33
34
35
36
37
38
39
40
41
42
43
44
45
46
47
48
49
50
51
52
53
54
55
56
57
58
59
60
61
62
63
64
65
66
67
68
69
70
71
72
73
74
75
76
77
78
79
80
全局调度器，全局只有一个schedt类型的实例。

sudoG 结构体：

// sudog 代表在等待列表里的 g，比如向 channel 发送/接收内容时
// 之所以需要 sudog 是因为 g 和同步对象之间的关系是多对多的
// 一个 g 可能会在多个等待队列中，所以一个 g 可能被打包为多个 sudog
// 多个 g 也可以等待在同一个同步对象上
// 因此对于一个同步对象就会有很多 sudog 了
// sudog 是从一个特殊的池中进行分配的。用 acquireSudog 和 releaseSudog 来分配和释放 sudog
 
type sudog struct {
	// The following fields are protected by the hchan.lock of the
	// channel this sudog is blocking on. shrinkstack depends on
	// this for sudogs involved in channel ops.
 
	g          *g
	selectdone *uint32 // CAS to 1 to win select race (may point to stack)
	next       *sudog
	prev       *sudog
	elem       unsafe.Pointer // data element (may point to stack)
 
	// The following fields are never accessed concurrently.
	// For channels, waitlink is only accessed by g.
	// For semaphores, all fields (including the ones above)
	// are only accessed when holding a semaRoot lock.
 
	acquiretime int64
	releasetime int64
	ticket      uint32
	parent      *sudog // semaRoot binary tree
	waitlink    *sudog // g.waiting list or semaRoot
	waittail    *sudog // semaRoot
	c           *hchan // channel
}
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
24
25
26
27
28
29
30
31
Section3 主要调度流程的源码分析
下面主要分析调度器调度流程的源码，主要分为几个部分：
1）预备知识点
2）main程序启动初始化过程（单独一篇文章写这个）
3）新建 goroutine的过程
4）循环调度的过程
5）抢占式调度的实现
6）初始化P过程
7）初始化M过程
8）初始化G过程

3.1 预备知识
在学习源码之前，需要了解一些关于Golang的一些规范和预备知识。

3.1.1 golang的函数调用规范
在golang里面调用函数，必须要关注的就是参数的传入和返回值。Golang有自己的一套函数调用规范，这个规范定义，所有参数都通过栈传递，返回值也是通过栈传递。

比如，对于函数：

type MyStruct struct { 
	X int
	P *int 
}
func someFunc(x int, s MyStruct) (int, MyStruct) { ... }
1
2
3
4
5
调用函数时候，栈的内容如下：


可以看到，参数和返回值都是从低位到高位排列，go函数可以有多个返回值的原因也在此，因为返回值都通过栈传递了。

需要注意这里的"返回地址"是x86和x64上的, arm的返回地址会通过LR寄存器保存, 内容会和这里的稍微不一样.

另外注意的是go和c不一样, 传递 struct 时整个struct 的内容都会复制到栈上, 如果构造体很大将会影响性能。

3.1.2 TLS(thread local storage)
TLS全称是Thread Local Storage，代表每个线程中的本地数据。写入TLS中的数据不会干扰到其余线程中的值。

Go的协程实现非常依赖于TLS机制，会用于获取系统线程中当前的G和G所属于的M实例。

Go操作TLS会使用系统原生的接口，以Linux X64为例，
go在新建M时候会调用arch_prctl 这个syscall来设置FS寄存器的值为M.tls的地址，
运行中每个M的FS寄存器都会指向它们对应的M实例的tls，linux内核调度线程时FS寄存器会跟着线程一起切换，
这样go代码只需要访问FS寄存机就可以获取到线程本地的数据。

3.1.3 栈扩张
go的协程设计是stackful coroutine，每一个goroutine都需要有自己的栈空间，
栈空间的内容再goroutine休眠时候需要保留的，等到重新调度时候恢复(这个时候整个调用树是完整的)。
这样就会引出一个问题，如果系统存在大量的goroutine，给每一个goroutine都预先分配一个足够的栈空间那么go就会使用过多的内存。

为了避免内存使用过多问题，go在一开始时候，会默认只为goroutine分配一个很小的栈空间，它的大小在1.92版本中是2k。
当函数发现栈空间不足时，会申请一块新的栈空间并把原来的栈复制过去。

g实例里面的g.stack、g.stackguard0两个变量来描述goroutine实例的栈。

3.1.4 写屏障(write barrier)
go支持并行GC的，GC的扫描阶段和go代码可以同时运行。这样带来的问题是，GC扫描的过程中go代码的执行可能改变了对象依赖树。

比如：开始扫描时候发现根对象A和B，B拥有C的指针，GC先扫描A，然后B把C的指针交给A，GC再扫描B，这时C就不会被扫描到。
为了避免这个问题，go在GC扫描标记阶段会启用写屏障（Write Barrier）

启用了Write barrier之后，当B把C指针交给A时，GC会认为在这一轮扫描中C的指针是存活的，即使A 可能在稍后丢掉C，那么C在下一轮GC中再回收。

Write barrier只针对指针启用，而且只在GC的标记阶段启用，平时会直接把值写入到目标地址。

3.1.5 m0和g0
go中有特殊的M和G，它们分别是m0和g0。

m0是启动程序后的主线程，这个M对应的实例会在全局变量runtime.m0中，不需要在heap上分配，
m0负责执行初始化操作和启动第一个g， 在之后m0就和其他的m一样了。

g0是仅用于负责调度的G，g0不指向任何可执行的函数, 每个m都会有一个自己的g0。

在调度或系统调用时会使用g0的栈空间, 全局变量的g0是m0的g0。

3.1.6 go中线程的种类
在 runtime 中有三种线程：

一种是主线程，
一种是用来跑 sysmon 的线程，
一种是普通的用户线程。
主线程在 runtime 由对应的全局变量: runtime.m0 来表示。用户线程就是普通的线程了，和 p 绑定，执行 g 中的任务。虽然说是有三种，实际上前两种线程整个 runtime 就只有一个实例。用户线程才会有很多实例。

主线程中用来跑 runtime.main，流程线性执行，没有跳转。

3.2 main线程启动执行
main线程的启动是伴随着go的main goroutine一起启动的，具体的启动流程可看另外一篇博文：
Golang-bootstrap分析 里面关于scheduler.main函数的分析。

3.3 新建goroutine过程
前面已经讲过了，当我们用 go func() 创建一个写的goroutine时候，compiler会编译成对runtime.newproc()的调用。堆栈的结构如下：


runtime.newproc()源码如下：

// Create a new g running fn with siz bytes of arguments.
// Put it on the queue of g's waiting to run.
// The compiler turns a go statement into a call to this.
// Cannot split the stack because it assumes that the arguments
// are available sequentially after &fn; they would not be
// copied if a stack split occurred.
//go:nosplit
// 根据 参数 fn 和 siz 创建一个 g
// 并把它放置入 自由g队列中等待唤醒
func newproc(siz int32, fn *funcval) {
	//获取栈上的参数的指针地址
	argp := add(unsafe.Pointer(&fn), sys.PtrSize)
	gp := getg()
	// 获取pc指针地址
	pc := getcallerpc()
	// 用g0的栈创建G对象
	systemstack(func() {
		newproc1(fn, (*uint8)(argp), siz, gp, pc)
	})
}
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
newproc只做了三件事：

计算参数的地址 argp
获取调用端的地址(返回地址) pc
使用systemstack调用 newproc1 函数，也就是用g0的栈创建g对象
systemstack 会切换当前的 g 到 g0, 并且使用g0的栈空间, 然后调用传入的函数, 再切换回原来的g和原来的栈空间。
切换到g0后会假装返回地址是mstart, 这样traceback的时候可以在mstart停止。
这里传给systemstack的是一个闭包, 调用时会把闭包的地址放到寄存器rdx, 具体可以参考上面对闭包的分析。

下面主要看 newproc1 函数主要做的事情：

// Create a new g running fn with narg bytes of arguments starting
// at argp. callerpc is the address of the go statement that created
// this. The new g is put on the queue of g's waiting to run.
// 根据函数参数和函数地址，创建一个新的G，然后将这个G加入队列等待运行
func newproc1(fn *funcval, argp *uint8, narg int32, callergp *g, callerpc uintptr) {
	// get g0
	_g_ := getg()
	// 设置g0对应的m的locks++, 禁止抢占
	_g_.m.locks++ // disable preemption because it can be holding p in a local var
	// get the p that m has
	_p_ := _g_.m.p.ptr()
	// new a g
	newg := gfget(_p_)
	if newg == nil {
		newg = malg(_StackMin)
		casgstatus(newg, _Gidle, _Gdead)
		allgadd(newg) // publishes with a g->status of Gdead so GC scanner doesn't look at uninitialized stack.
	}
	totalSize := 4*sys.RegSize + uintptr(siz) + sys.MinFrameSize // extra space in case of reads slightly beyond frame
	totalSize += -totalSize & (sys.SpAlign - 1)                  // align to spAlign
	sp := newg.stack.hi - totalSize
	spArg := sp
	
	//  初始化 g，g 的 gobuf 现场，g 的 m 的 curg
	// 以及各种寄存器
	memclrNoHeapPointers(unsafe.Pointer(&newg.sched), unsafe.Sizeof(newg.sched))
	newg.sched.sp = sp
	newg.stktopsp = sp
	newg.sched.pc = funcPC(goexit) + sys.PCQuantum // +PCQuantum so that previous instruction is in same function
	newg.sched.g = guintptr(unsafe.Pointer(newg))
	gostartcallfn(&newg.sched, fn)
	newg.gopc = callerpc
	newg.ancestors = saveAncestors(callergp)
	newg.startpc = fn.fn
	if _g_.m.curg != nil {
		newg.labels = _g_.m.curg.labels
	}
	if isSystemGoroutine(newg) {
		atomic.Xadd(&sched.ngsys, +1)
	}
	newg.gcscanvalid = false
	casgstatus(newg, _Gdead, _Grunnable)

	newg.goid = int64(_p_.goidcache)
	_p_.goidcache++
	//
	runqput(_p_, newg, true)

	if atomic.Load(&sched.npidle) != 0 && atomic.Load(&sched.nmspinning) == 0 && mainStarted {
		wakep()
	}
	_g_.m.locks--
	if _g_.m.locks == 0 && _g_.preempt { // restore the preemption request in case we've cleared it in newstack
		_g_.stackguard0 = stackPreempt
	}
}
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
24
25
26
27
28
29
30
31
32
33
34
35
36
37
38
39
40
41
42
43
44
45
46
47
48
49
50
51
52
53
54
55
56
处理流程如下：

调用getg(汇编实现)获取当前的g, 会编译为读取FS寄存器(TLS), 这里会获取到g0
设置g对应的m的locks++, 禁止抢占
获取m拥有的p
新建一个g策略：
调用 gfget函数，这里是复用优先策略
首先从p的gfree获取回收的g，如果p.gfree链表为空，就从全局调度器sched里面的gfree链表里面steal 32个free的g给p.gfree。
将p.gfree链表的head元素获取返回。
如果获取不到freeg时调用malg()函数新建一个g, 初始的栈空间大小是2K。
把参数复制到g的栈上
把返回地址复制到g的栈上, 这里的返回地址是goexit, 表示调用完目标函数后会调用goexit
设置g的调度数据(sched)
设置sched.sp等于参数+返回地址后的rsp地址
设置sched.pc等于目标函数的地址, 查看gostartcallfn和gostartcall
设置sched.g等于g
设置g的状态为待运行(_Grunnable)
调用runqput函数把g放到运行队列
首先随机把g放到p.runnext, 如果放到runnext则入队原来在runnext的g；
然后尝试把g放到P的local queue；
如果local queue（256 capacity）满了则调用runqputslow函数把g放到"全局运行队列"（操作全局 sched 时，需要获取全局 sched.lock 锁，全局锁争抢的开销较大，所以才称之为 slow
runqputslow会把本地运行队列中一半的g放到全局运行队列, 这样下次就可以快速使用local queue.
如果当前有空闲的P，但是没有自旋的M(nmspinning等于0)，并且主函数已执行，则唤醒或新建一个M来调度一个P执行
这一步非常重要, 用于保证当前有足够的M运行G, 具体请查看上面的"空闲M链表"
唤醒或新建一个M会通过调用wakep函数
首先交换nmspinning到1, 成功再继续, 多个线程同时执行wakep函数只有一个会继续
调用startm函数
调用pidleget从"空闲P链表"获取一个空闲的P
调用mget从"空闲M链表"获取一个空闲的M
如果没有空闲的M, 则调用newm新建一个M
newm会新建一个m的实例, m的实例包含一个g0, 然后调用newosprocclone一个系统线程
newosproc会调用syscall clone创建一个新的线程
线程创建后会设置TLS, 设置TLS中当前的g为g0, 然后执行mstart
调用notewakeup(&mp.park)唤醒线程
创建goroutine的流程就这么多了, 接下来看看M是如何调度的.

3.4 循环调度schedule过程
从前面描述的调度器的工作流可知，scheduler是一个循环的过程。

M启动时会调用mstart函数，m0在初始化后调用，其他的m在线程启动后调用。proc.go源码里面这个函数的注释就是：

//Called to start an M. begin scheduling
1
mstart函数的处理如下:

首先调用getg函数获取当前的g, 这里会获取到g0
如果g0未分配栈则从当前的栈空间(系统栈空间)上分配, 也就是说g0会使用系统栈空间
调用mstart1函数
调用gosave函数保存当前的状态到g0的调度数据中, 以后每次调度都会从这个栈地址开始；
调用asminit函数, 不做任何事情；
调用minit函数, 设置当前线程可以接收的信号(signal)；
调用schedule函数。
调用schedule函数后就进入了调度循环，整个流程可以简单总结为：

schedule函数获取g => [必要时休眠] => [唤醒后继续获取] => execute函数执行g => 执行后返回到goexit => 重新执行schedule函数
schedule函数处理流程如下：

获取当前调度的g，也就是g0，g0在执行调度逻辑；
如果当前GC需要停止整个世界（STW), 则调用gcstopm休眠当前的M；
如果M拥有的P中指定了需要在安全点运行的函数(P.runSafePointFn), 则运行它；
快速获取待运行的G, 以下处理如果有一个获取成功后面就不会继续获取：
如果当前GC正在标记阶段, 则查找有没有待运行的GC Worker, GC Worker也是一个G；
为了公平起见, 每61次调度从全局运行队列获取一次G, (一直从本地获取可能导致全局运行队列中的G不被运行)；
从P的本地运行队列中获取G, 调用runqget函数。
快速获取失败时, 调用findrunnable函数获取待运行的G, 会阻塞到获取成功为止：
如果当前GC需要停止整个世界（STW), 则调用stopm休眠当前的M；
如果M拥有的P中指定了需要在安全点运行的函数(P.runSafePointFn), 则运行它；
如果有析构器待运行则使用"运行析构器的G"；
从P的本地运行队列中获取G, 调用runqget函数，如果获取到就返回；
从全局运行队列获取G, 调用globrunqget函数, 需要上锁，获取到就返回。；
从网络事件反应器获取G, 函数netpoll会获取哪些fd可读可写或已关闭, 然后返回等待fd相关事件的G；
如果从local 和 global 都获取不到G, 则执行Work Stealing：
调用runqsteal尝试从其他P的本地运行队列盗取一半的G。
如果还是获取不到G, 就需要休眠M了, 接下来是休眠的步骤：
再次检查当前GC是否在标记阶段, 在则查找有没有待运行的GC Worker, GC Worker也是一个G；
再次检查如果当前GC需要停止整个世界, 或者P指定了需要再安全点运行的函数, 则跳到findrunnable的顶部重试；
再次检查全局运行队列中是否有G, 有则获取并返回；
释放M拥有的P, P会变为空闲(_Pidle)状态；
把P添加到"空闲P链表"中；
让M离开自旋状态, 这里的处理非常重要, 参考上面的"空闲M链表"；
首先减少表示当前自旋中的M的数量的全局变量nmspinning；
再次检查所有P的本地运行队列, 如果不为空则让M重新进入自旋状态, 并跳findrunnable的顶部重试；
再次检查有没有待运行的GC Worker, 有则让M重新进入自旋状态, 并跳到findrunnable的顶部重试
再次检查网络事件反应器是否有待运行的G, 这里对netpoll的调用会阻塞, 直到某个fd收到了事件；
如果最终还是获取不到G, 调用stopm休眠当前的M；
唤醒后跳到findrunnable的顶部重试。
成功获取到一个待运行的G；
让M离开自旋状态, 调用resetspinning, 这里的处理和上面的不一样：
如果当前有空闲的P, 但是无自旋的M(nmspinning等于0), 则唤醒或新建一个M；
上面离开自旋状态是为了休眠M, 所以会再次检查所有队列然后休眠；
这里离开自选状态是为了执行G, 所以会检查是否有空闲的P, 有则表示可以再开新的M执行G。
如果G要求回到指定的M(例如上面的runtime.main)：
调用startlockedm函数把G和P交给该M, 自己进入休眠；
从休眠唤醒后跳到schedule的顶部重试
调用execute函数在当前M上执行G。
execute函数执行gp的处理如下:

调用getg获取当前的g(g0)；
把G(gp)的状态由待运行(_Grunnable)改为运行中(_Grunning)；
设置G的stackguard, 栈空间不足时可以扩张；
增加P中记录的调度次数(对应上面的每61次优先获取一次全局运行队列)；
设置g.m.curg = g；
设置gp.m = g0.m；
调用gogo函数(在M上执行gp，通过汇编代码实现的，在asm_amd64.s里面有gogo函数在64位平台上实现)：
这个函数会根据g.sched中保存的状态恢复各个寄存器的值并继续运行g
首先针对g.sched.ctxt调用写屏障(GC标记指针存活), ctxt中一般会保存指向[函数+参数]的指针
设置TLS中的g为g.sched.g, 也就是g自身
设置rsp寄存器为g.sched.rsp
设置rax寄存器为g.sched.ret
设置rdx寄存器为g.sched.ctxt (上下文)
设置rbp寄存器为g.sched.rbp
清空sched中保存的信息
跳转到g.sched.pc
因为前面创建goroutine的newproc1函数把返回地址设为了goexit, 函数运行完毕返回时将会调用goexit函数。
自此一次调度过程就结束了。
g.sched.pc在G首次运行时会指向目标函数的第一条机器指令,
如果G被抢占或者等待资源而进入休眠, 在休眠前会保存状态到g.sched,
g.sched.pc会变为唤醒后需要继续执行的地址, "保存状态"的实现将在下面讲解.

目标函数执行完毕后会调用goexit函数，goexit函数会调用goexit1函数，goexit1函数会通过mcall调用goexit0函数。
mcall这个函数就是用于实现"保存状态"的, 处理如下:

设置g.sched.pc等于当前的返回地址
设置g.sched.sp等于寄存器rsp的值
设置g.sched.g等于当前的g
设置g.sched.bp等于寄存器rbp的值
切换TLS中当前的g等于m.g0
设置寄存器rsp等于g0.sched.sp, 使用g0的栈空间
设置第一个参数为原来的g
设置rdx寄存器为指向函数地址的指针(上下文)
调用指定的函数, 不会返回。
mcall这个函数保存当前的运行状态到g.sched, 然后切换到g0和g0的栈空间, 再调用指定的函数。
回到g0的栈空间这个步骤非常重要, 因为这个时候g已经中断, 继续使用g的栈空间且其他M唤醒了这个g将会产生灾难性的后果。
G在中断或者结束后都会通过mcall回到g0的栈空间继续调度, 从goexit调用的mcall的保存状态其实是多余的, 因为G已经结束了。

goexit1函数会通过mcall调用goexit0函数, goexit0函数调用时已经回到了g0的栈空间, 处理如下:

把G的状态由运行中(_Grunning)改为已中止(_Gdead)
清空G的成员
调用dropg函数解除M和G之间的关联
调用gfput函数把G放到P的自由列表中, 下次创建G时可以复用
调用schedule函数继续调度
G结束后回到schedule函数, 这样就结束了一个调度循环。
不仅只有G结束会重新开始调度, G被抢占或者等待资源也会重新进行调度, 下面继续来看这两种情况。

3.5 抢占式调度实现(sysmon线程)
Golang-bootstrap分析 这篇文章里面我提到了runtime.main会创建一个额外的M运行sysmon函数，抢占式调度就是在sysmon中实现的。

sysmon函数(4249行) 会进入一个无限循环，第一轮回休眠20us，之后每次休眠时间倍增，最终每一轮都会休眠10ms。

sysmon中有netpool(获取fd事件)，retake(抢占)，forcegc(按时间强制执行gc),scavenge heap(释放自由列表中多余的项减少内存占用)等处理。

retake函数负责处理抢占，流程是:

枚举所有的P：
如果P在系统调用中(_Psyscall), 且经过了一次sysmon循环(20us~10ms), 则抢占这个P
调用handoffp解除M和P之间的关联
如果P在运行中(_Prunning), 且经过了一次sysmon循环并且G运行时间超过forcePreemptNS(10ms), 则抢占这个P
调用preemptone函数
设置g.preempt = true
设置g.stackguard0 = stackPreempt
为什么设置了stackguard就可以实现抢占?
因为这个值用于检查当前栈空间是否足够, go函数的开头会比对这个值判断是否需要扩张栈.
stackPreempt是一个特殊的常量, 它的值会比任何的栈地址都要大, 检查时一定会触发栈扩张.

栈扩张调用的是morestack_noctxt函数, morestack_noctxt函数清空rdx寄存器并调用morestack函数。
morestack函数会保存G的状态到g.sched, 切换到g0和g0的栈空间, 然后调用newstack函数。
newstack函数判断g.stackguard0等于stackPreempt, 就知道这是抢占触发的, 这时会再检查一遍是否要抢占：

如果M被锁定(函数的本地变量中有P), 则跳过这一次的抢占并调用gogo函数继续运行G
如果M正在分配内存, 则跳过这一次的抢占并调用gogo函数继续运行G
如果M设置了当前不能抢占, 则跳过这一次的抢占并调用gogo函数继续运行G
如果M的状态不是运行中, 则跳过这一次的抢占并调用gogo函数继续运行G
即使这一次抢占失败, 因为g.preempt等于true, runtime中的一些代码会重新设置stackPreempt以重试下一次的抢占。
如果判断可以抢占, 则继续判断是否GC引起的, 如果是则对G的栈空间执行标记处理(扫描根对象)然后继续运行。
如果不是GC引起的则调用gopreempt_m函数完成抢占。

gopreempt_m函数会调用goschedImpl函数, goschedImpl函数的流程是:

把G的状态由运行中(_Grunnable)改为待运行(_Grunnable)
调用dropg函数解除M和G之间的关联
调用globrunqput把G放到全局运行队列
调用schedule函数继续调度
为全局运行队列的优先度比较低, 各个M会经过一段时间再去重新获取这个G执行，

抢占机制保证了不会有一个G长时间的运行导致其他G无法运行的情况发生。

Section4：scheduler与memory allocation、channel、garbage collection关联部分
这部分需要完成其余三个主题的学习了之后才能解释~

这篇文章写的真是艰难，参考了很多文章博客，对于scheduler的设计有了基本的了解，学完之后发现其实与channel原理、memory allocation和GC其实关联挺深的，需要详细学习其余部分之后对本文进行再一次修改。

