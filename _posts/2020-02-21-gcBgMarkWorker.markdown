---
title: gcBgMarkWorker
layout: post
category: golang
author: 夏泽民
---
https://gocn.vip/topics/9822
一. Go GC 要点
先来回顾一下 GC 的几个重要的阶段:

Mark Prepare - STW
做标记阶段的准备工作，需要停止所有正在运行的 goroutine(即 STW)，标记根对象，启用内存屏障，内存屏障有点像内存读写钩子，它用于在后续并发标记的过程中，维护三色标记的完备性 (三色不变性)，这个过程通常很快，大概在 10-30 微秒。

Marking - Concurrent
标记阶段会将大概 25%(gcBackgroundUtilization) 的 P 用于标记对象，逐个扫描所有 G 的堆栈，执行三色标记，在这个过程中，所有新分配的对象都是黑色，被扫描的 G 会被暂停，扫描完成后恢复，这部分工作叫后台标记 (gcBgMarkWorker)。这会降低系统大概 25% 的吞吐量，比如MAXPROCS=6，那么 GC P 期望使用率为6*0.25=1.5，这 150%P 会通过专职 (Dedicated)/兼职 (Fractional)/懒散 (Idle) 三种工作模式的 Worker 共同来完成。


这还没完，为了保证在 Marking 过程中，其它 G 分配堆内存太快，导致 Mark 跟不上 Allocate 的速度，还需要其它 G 配合做一部分标记的工作，这部分工作叫辅助标记 (mutator assists)。在 Marking 期间，每次 G 分配内存都会更新它的” 负债指数”(gcAssistBytes)，分配得越快，gcAssistBytes 越大，这个指数乘以全局的” 负载汇率”(assistWorkPerByte)，就得到这个 G 需要帮忙 Marking 的内存大小 (这个计算过程叫revise)，也就是它在本次分配的 mutator assists 工作量 (gcAssistAlloc)。
<!-- more -->
Mark Termination - STW
标记阶段的最后工作是 Mark Termination，关闭内存屏障，停止后台标记以及辅助标记，做一些清理工作，整个过程也需要 STW，大概需要 60-90 微秒。在此之后，所有的 P 都能继续为应用程序 G 服务了。

Sweeping - Concurrent
在标记工作完成之后，剩下的就是清理过程了，清理过程的本质是将没有被使用的内存块整理回收给上一个内存管理层级 (mcache -> mcentral -> mheap -> OS)，清理回收的开销被平摊到应用程序的每次内存分配操作中，直到所有内存都 Sweeping 完成。当然每个层级不会全部将待清理内存都归还给上一级，避免下次分配再申请的开销，比如 Go1.12 对 mheap 归还 OS 内存做了优化，使用NADV_FREE延迟归还内存。

STW
在Go 调度模型中我们已经提到，Go 没有真正的实时抢占机制，而是一套协作式抢占 (cooperative preemption)，即给 G(groutine) 打个标记，等待 G 在调用函数时检查这个标记，以此作为一个安全的抢占点 (GC safe-point)。但如果其它 P 上的 G 都停了，某个 G 还在执行如下代码:

func add(numbers []int) int {
     var v int
     for _, n := range numbers {
         v += n
     }
     return v
}
add 函数的运行时间取决于切片的长度，并且在函数内部是没有调用其它函数的，也就是没有抢占点。就会导致整个运行时都在等待这个 G 调用函数 (以实现抢占，开始处理 GC)，其它 P 也被挂起。这就是 Go GC 最大的诟病: GC STW 时间会受到 G 调用函数的时机的影响并被延长，甚至如果某个 G 在执行无法抢占的死循环 (即循环内部没有发生函数调用的死循环)，那么整个 Go 的 runtime 都会挂起，CPU 100%，节点无法响应任何消息，连正常停服都做不到。pprof 这类调试工具也用不了，只能通过 gdb，delve 等外部调试工具来找到死循环的 goroutine 正在执行的堆栈。如此后果比没有被 defer 的 panic 更严重，因为那个时候的节点内部状态是无法预期的。

因此有 Gopher 开始倡议 Go 使用非协作式抢占 (non-cooperative preemption)，通过堆栈和寄存器来保存抢占上下文，避免对抢占不友好的函数导致 GC STW 延长 (毕竟第三方库代码的质量也是参差不齐的)。相关的 Issue 在这里。好消息是，Go1.14(目前还是 Beta1 版本，还未正式发布) 已经支持异步抢占，也就是说:

// 简单起见，没用channel协同
func main() {
  go func() {
    for {
    }
  }()

  time.Sleep(time.Millisecond)
  runtime.GC()
  println("OK")
}
这段代码在 Go1.14 中终于能输出OK了。这个提了近五年的 Issue: runtime: tight loops should be preemptible #10958前几天终于关闭了。不得不说，这是 Go Runtime 的一大进步，它不止避免了单个 goroutine 死循环导致整个 runtime 卡死的问题，更重要的是，它为 STW 提供了最坏预期，避免了 GC STW 造成了性能抖动隐患。

二. Go GC 度量
1. go tool prof
Go 基础性能分析工具，pprof 的用法和启动方式参考go pprof 性能分析，其中的 heap 即为内存分配分析，go tool 默认是查看正在使用的内存 (inuse_heap)，如果要看其它数据，使用go tool pprof --alloc_space|inuse_objects|alloc_objects。

需要注意的是，go pprof 本质是数据采样分析，其中的值并不是精确值，适用于性能热点优化，而非真实数据统计。

2. go tool trace
go tool trace 可以将 GC 统计信息以可视化的方式展现出来。要使用 go tool trace，可以通过以下方式生成采样数据:

API: trace.Start
go test: go test -trace=trace.out pkg
net/http/pprof: curl http://127.0.0.1:6060/debug/pprof/trace?seconds=20
得到采样数据后，之后即可以通过 go tool trace trace.out 启动一个 HTTP Server，在浏览器中查看可视化 trace 数据:



里面提供了各种 trace 和 prof 的可视化入口，点击第一个 View trace 可以看到追踪总览:



包含的信息量比较广，横轴为时间线，各行为各种维度的度量，通过 A/D 左右移动，W/S 放大放小。以下是各行的意义:

Goroutines: 包含 GCWaiting，Runnable，Running 三种状态的 Goroutine 数量统计
Heap: 包含当前堆使用量 (Allocated) 和下次 GC 阈值 (NextGC) 统计
Threads: 包含正在运行和正在执行系统调用的 Threads 数量
GC: 哪个时间段在执行 GC
ProcN: 各个 P 上面的 goroutine 调度情况
除了View trace之外，trace 目录的第二个Goroutine analysis也比较有用，它能够直观统计 Goroutine 的数量和执行状态:





通过它可以对各个 goroutine 进行健康诊断，各种 network,syscall 的采样数据下载下来之后可以直接通过go tool pprof分析，因此，实际上 pprof 和 trace 两套工具是相辅相成的。

3. GC Trace
GC Trace 是 Golang 提供的非侵入式查看 GC 信息的方案，用法很简单，设置GCDEBUG=gctrace=1环境变量即可:

GODEBUG=gctrace=1 bin/game
gc 1 @0.039s 3%: 0.027+4.5+0.015 ms clock, 0.11+2.3/4.0/5.5+0.063 ms cpu, 4->4->2 MB, 5 MB goal, 4 P
gc 2 @0.147s 1%: 0.007+1.2+0.008 ms clock, 0.029+0.15/1.1/2.0+0.035 ms cpu, 5->5->3 MB, 6 MB goal, 4 P
gc 3 @0.295s 0%: 0.010+2.3+0.013 ms clock, 0.040+0.14/2.1/4.3+0.053 ms cpu, 7->7->4 MB, 8 MB goal, 4 P
下面是各项指标的解释:

gc 1 @0.039s 3%: 0.027+4.5+0.015 ms clock, 0.11+2.3/4.0/5.5+0.063 ms cpu, 4->4->2 MB, 5 MB goal, 4 P

// 通用参数
gc 2: 程序运行后的第2次GC
@0.147s: 到目前为止程序运行的时间
3%: 到目前为止程序花在GC上的CPU%

// Wall-Clock 流逝的系统时钟
0.027ms+4.5ms+0.015 ms   : 分别是 STW Mark Prepare，Concurrent Marking，STW Mark Termination 的时钟时间

// CPU Time 消耗的CPU时间
0.11+2.3/4.0/5.5+0.063 ms : 以+分隔的阶段同上，不过将Concurrent Marking细分为Mutator Assists Time, Background GC Time(包括Dedicated和Fractional Worker), Idle GC Time三种。其中0.11=0.027*4，0.063=0.015*4。

// 内存相关统计
4->4->2 MB: 分别是开始标记时，标记结束后的堆占用大小，以及标记结束后真正存活的(有效的)堆内存大小
5 MB goal: 下次GC Mark Termination后的目标堆占用大小，该值受GC Percentage影响，并且会影响mutator assist工作量(每次堆大小变更时都动态评估，如果快超出goal了，就需要其它goroutine帮忙干活了, https://github.com/golang/go/blob/dev.boringcrypto.go1.13/src/runtime/mgc.go#L484)

// Processors
4 P : P的数量，也就是GOMAXPROCS大小，可通过runtime.GoMaxProcs设置

// 其它
GC forced: 如果两分钟内没有执行GC，则会强制执行一次GC，此时会换行打印 GC forced
4. MemStats
runtime.MemStats记录了内存分配的一些统计信息，通过runtime.ReadMemStats(&ms)获取，它是runtime.mstats的对外版 (再次可见 Go 单一访问控制的弊端)，MemStats 字段比较多，其中比较重要的有:

// HeapSys 

// 以下内存大小字段如无特殊说明单位均为bytes
type MemStats struct {
    // 从开始运行到现在累计分配的堆内存数
    TotalAlloc uint64

    // 从OS申请的总内存数(包含堆、栈、内部数据结构等)
    Sys uint64

    // 累计分配的堆对象数量 (当前存活的堆对象数量=Mallocs-Frees)
    Mallocs uint64

    // 累计释放的堆对象数量
    Frees   uint64

    // 正在使用的堆内存数，包含可访问对象和暂未被GC回收的不可访问对象
    HeapAlloc uint64

    // 虚拟内存空间为堆保留的大小，包含还没被使用的(还没有映射物理内存，但这部分通常很小)
    // 以及已经将物理内存归还给OS的部分(即HeapReleased)
    // HeapSys = HeapInuse + HeapIdle
    HeapSys uint64

    // 至少包含一个对象的span字节数
    // Go GC是不会整理内存的
    // HeapInuse - HeapAlloc 是为特殊大小保留的内存，但是它们还没有被使用
    HeapInuse uint64

    // 未被使用的span中的字节数
    // 未被使用的span指没有包含任何对象的span，它们可以归还OS，也可以被重用，或者被用于栈内存
    // HeapIdle - HeadReleased 即为可以归还OS但还被保留的内存，这主要用于避免频繁向OS申请内存
    HeapIdle uint64

    // HeapIdle中已经归还给OS的内存量
    HeapReleased uint64

    // ....
}
程序可以通过定期调用runtime.ReadMemStatsAPI 来获取内存分配信息发往时序数据库进行监控。另外，该 API 是会 STW 的，但是很短，Google 内部也在用，用他们的话说:” STW 不可怕，长时间 STW 才可怕”，该 API 通常一分钟调用一次即可。

5. ReadGCStats
debug.ReadGCStats用于获取最近的 GC 统计信息，主要是 GC 造成的延迟信息:

// GCStats collect information about recent garbage collections.
type GCStats struct {
    LastGC         time.Time       // 最近一次GC耗费时间
    NumGC          int64           // 执行GC的次数
    PauseTotal     time.Duration   // 所有GC暂停时间总和
    Pause          []time.Duration // 每次GC的暂停时间，最近的排在前面
    ...
}
和 ReadMemStats 一样，ReadGCStats 也可以定时收集，发送给时序数据库做监控统计。

三. Go GC 调优
Go GC 相关的参数少得可怜，一如既往地精简:

1. debug.SetGCPercent
一个百分比数值，决定即本次 GC 后，下次触发 GC 的阈值，比如本次 GC Sweeping 完成后的内存占用为 200M，GC Percentage 为 100(默认值)，那么下次触发 GC 的内存阈值就是 400M。这个值通常不建议修改，因为优化 GC 开销的方法通常是避免不必要的分配或者内存复用，而非通过调整 GC Percent 延迟 GC 触发时机 (Go GC 本身也会根据当前分配速率来决定是否需要提前开启新一轮 GC)。另外，debug.SetGCPercent 传入<0 的值将关闭 GC。

2. runtime.GC
强制执行一次 GC，如果当前正在执行 GC，则帮助当前 GC 执行完成后，再执行一轮完整的 GC。该函数阻塞直到 GC 完成。

3. debug.FreeOSMemory
强制执行一次 GC，并且尽可能多地将不再使用的内存归还给 OS。

严格意义上说，以上几个 API 预期说调优，不如说是补救，它们都只是把 Go GC 本身就会做的事情提前或者延后了，通常是治标不治本的方法。真正的 GC 调优主要还是在应用层面。我在这篇文章聊了一些 Go 应用层面的内存优化。

以上主要从偏应用的角度介绍了 Golang GC 的几个重要阶段，STW，GC 度量/调试，以及相关 API 等。这些理论和方法能在在必要的时候派上用场，帮助更深入地了解应用程序并定位问题。

推荐文献:

Garbage Collection In Go
GC 20 问
A visual guide to Go Memory Allocator from scratch

https://www.ardanlabs.com/blog/2018/12/garbage-collection-in-go-part1-semantics.html

https://github.com/qcrao/Go-Questions/blob/master/GC/GC.md

https://blog.learngoprogramming.com/a-visual-guide-to-golang-memory-allocator-from-ground-up-e132258453ed

https://wudaijun.com/2020/01/go-gc-keypoint-and-monitor/
