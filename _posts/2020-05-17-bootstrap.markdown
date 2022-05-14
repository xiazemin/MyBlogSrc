---
title: bootstrap
layout: post
category: golang
author: 夏泽民
---
https://blog.csdn.net/u010853261/article/details/84790392
https://blog.csdn.net/u010853261/article/details/84901386


<!-- more -->
1. 环境
要分析runtime相关内部机制，首先从系统启动开始。首先准备分析环境：golang、OS、gdb


2. 引导程序宏观流程
在go代码里面，用户逻辑从main.main()开始，那么runtime如何启动？怎么初始化？初始化做了哪些工作呢？

这里我们从函数运行的起点开始分析。我们先编写一个最简单的go代码：

package main

import "fmt"

func main() {
	fmt.Println("hello,golang")
}
1
2
3
4
5
6
7
这是最简单的一个源码版本了，然后生成可执行文件，使用GDB动态查看即可。

$go build -o test2
$gdb test2
1
2
先用go build编译可执行文件，然后用GDB命令gdb test2进入调试分析源码界面。

使用gdb命令可以获取到系统的Entry Point然后找到函数入口。发现是在

(gdb) info symbol 0x1051fd0
_rt0_amd64_darwin in section .text
1
2
在rt0_darwin_amd64.s里面：

TEXT _rt0_amd64_darwin(SB),NOSPLIT,$-8
	JMP	_rt0_amd64(SB)
1
2
然后在：asm_amd64.s里面发现了_rt0_amd64代码段，如下：

TEXT _rt0_amd64(SB),NOSPLIT,$-8
	MOVQ	0(SP), DI	// argc
	LEAQ	8(SP), SI	// argv
	JMP	runtime·rt0_go(SB)
1
2
3
4
最后调用了runtime.runtime·rt0_go，系统初始化主要逻辑也是在这个地方(已删减不重要逻辑)；
具体的源码见 asm_amd64.s

TEXT runtime·rt0_go(SB),NOSPLIT,$0
	......
	CALL	runtime·args(SB)
	CALL	runtime·osinit(SB)
	CALL	runtime·schedinit(SB)

	// create a new goroutine to start program
	MOVQ	$runtime·mainPC(SB), AX		// entry
	PUSHQ	AX
	PUSHQ	$0			// arg size
	CALL	runtime·newproc(SB)
	POPQ	AX
	POPQ	AX

	// start this M
	CALL	runtime·mstart(SB)

	CALL	runtime·abort(SB)	// mstart should never return
	RET

	// Prevent dead-code elimination of debugCallV1, which is
	// intended to be called by debuggers.
	MOVQ	$runtime·debugCallV1(SB), AX
	RET

在完成命令行初始化、OS初始化、调度器初始化之后。使用newproc创建一个goroutine放入待运行队列。然后mstart让主线程进入任务调度模式，从队列提出main goroutine并执行。

bootstrap overview


3. 初始化流程
初始化的流程里面，命令行初始化和OS初始化与调度器的机制关系不大，这里就不做主要讲解，这里主要关心调度器的初始化，初始化函数schedinit()在runtime/proc.go 这个源码里面。

3.1 schedule init
初始化的代码在proc.go里面的schedule_init函数，源码如下(只保留核心流程代码)：

// The bootstrap sequence is:
//
//	call osinit
//	call schedinit
//	make & queue new G
//	call runtime·mstart
//
// The new G calls runtime·main.
func schedinit() {
	// M最大数量限制
	sched.maxmcount = 10000
	
	// 内存相关初始化
	stackinit()
	mallocinit()
	// m相关初始化
	mcommoninit(_g_.m)
	// 命令参数和环境初始化
	goargs()
	goenvs()
	// GC 初始化
	gcinit()
	// 设置GOMAXPROCS
	procs := ncpu
	if n, ok := atoi32(gogetenv("GOMAXPROCS")); ok && n > 0 {
		procs = n
	}
	if procresize(procs) != nil {
		throw("unknown runnable goroutine during bootstrap")
	}
}

调度器的初始化主要相关内容上面注释都描述的比较清楚了，主要是初始化栈空间分配器, GC, 按cpu核心数量或GOMAXPROCS的值生成P等。

特别要主要的是生成P的操作在procresize函数里面，如果在runtime需要更改P的数量也是调用这个函数。这个函数主要逻辑如下源码：

// Returns list of Ps with local work, they need to be scheduled by the caller.
func procresize(nprocs int32) *p {
	// Grow allp if necessary.
	if nprocs > int32(len(allp)) {
		// Synchronize with retake, which could be running
		// concurrently since it doesn't run on a P.
		lock(&allpLock)
		if nprocs <= int32(cap(allp)) {
			allp = allp[:nprocs]
		} else {
			nallp := make([]*p, nprocs)
			// Copy everything up to allp's cap so we
			// never lose old allocated Ps.
			copy(nallp, allp[:cap(allp)])
			allp = nallp
		}
		unlock(&allpLock)
	}

	// initialize new P's
	for i := int32(0); i < nprocs; i++ {
		pp := allp[i]
		if pp == nil {
			pp = new(p)
			pp.id = i
			pp.status = _Pgcstop
			pp.sudogcache = pp.sudogbuf[:0]
			for i := range pp.deferpool {
				pp.deferpool[i] = pp.deferpoolbuf[i][:0]
			}
			pp.wbBuf.reset()
			atomicstorep(unsafe.Pointer(&allp[i]), unsafe.Pointer(pp))
		}
		if pp.mcache == nil {
			if old == 0 && i == 0 {
				if getg().m.mcache == nil {
					throw("missing mcache?")
				}
				pp.mcache = getg().m.mcache // bootstrap
			} else {
				pp.mcache = allocmcache()
			}
		}
	}

	// free unused P's
	for i := nprocs; i < old; i++ {
		p := allp[i]
		// move all runnable goroutines to the global queue
		for p.runqhead != p.runqtail {
			// pop from tail of local queue
			p.runqtail--
			gp := p.runq[p.runqtail%uint32(len(p.runq))].ptr()
			// push onto head of global queue
			globrunqputhead(gp)
		}
		if p.runnext != 0 {
			globrunqputhead(p.runnext.ptr())
			p.runnext = 0
		}
		// if there's a background worker, make it runnable and put
		// it on the global queue so it can clean itself up
		if gp := p.gcBgMarkWorker.ptr(); gp != nil {
			casgstatus(gp, _Gwaiting, _Grunnable)
			if trace.enabled {
				traceGoUnpark(gp, 0)
			}
			globrunqput(gp)
			// This assignment doesn't race because the
			// world is stopped.
			p.gcBgMarkWorker.set(nil)
		}
		// Flush p's write barrier buffer.
		if gcphase != _GCoff {
			wbBufFlush1(p)
			p.gcw.dispose()
		}
		for i := range p.sudogbuf {
			p.sudogbuf[i] = nil
		}
		p.sudogcache = p.sudogbuf[:0]
		for i := range p.deferpool {
			for j := range p.deferpoolbuf[i] {
				p.deferpoolbuf[i][j] = nil
			}
			p.deferpool[i] = p.deferpoolbuf[i][:0]
		}
		freemcache(p.mcache)
		p.mcache = nil
		gfpurge(p)
		traceProcFree(p)
		if raceenabled {
			raceprocdestroy(p.racectx)
			p.racectx = 0
		}
		p.gcAssistTime = 0
		p.status = _Pdead
		// can't free P itself because it can be referenced by an M in syscall
	}

	// Trim allp.
	if int32(len(allp)) != nprocs {
		lock(&allpLock)
		allp = allp[:nprocs]
		unlock(&allpLock)
	}

	_g_ := getg()
	if _g_.m.p != 0 && _g_.m.p.ptr().id < nprocs {
		// continue to use the current P
		_g_.m.p.ptr().status = _Prunning
	} else {
		// release the current P and acquire allp[0]
		if _g_.m.p != 0 {
			_g_.m.p.ptr().m = 0
		}
		_g_.m.p = 0
		_g_.m.mcache = nil
		p := allp[0]
		p.m = 0
		p.status = _Pidle
		acquirep(p)
	}
	var runnablePs *p
	for i := nprocs - 1; i >= 0; i-- {
		p := allp[i]
		if _g_.m.p.ptr() == p {
			continue
		}
		p.status = _Pidle
		if runqempty(p) {
			pidleput(p)
		} else {
			p.m.set(mget())
			p.link.set(runnablePs)
			runnablePs = p
		}
	}
	stealOrder.reset(uint32(nprocs))
	var int32p *int32 = &gomaxprocs // make compiler check that gomaxprocs is an int32
	atomic.Store((*uint32)(unsafe.Pointer(int32p)), uint32(nprocs))
	return runnablePs
}

主要逻辑如下：

如果新的GOMAXPROCS的值大于已有的全局变量里面P的数量，那么需要扩容
加锁allpLock然后重新生成新的GOMAXPROCS数量的P数组，并将原来已有的P数组copy过来，重新构造P数组之后释放allpLock锁
初始化并new新的P(GOMAXPROCS数量大于old P数量)
如果GOMAXPROCS小于old P数量，那么旧的P需要free掉，free时候涉及到一些P持有相关资源的释放和转移
move all runnable goroutines to the global queue(contain local queue and p.runnext)
if there’s a background worker, make it runnable and put it on the global queue so it can clean itself up
release related cache
mark P state Pdead
重新赋值runtime.allp
如果当前g的P仍然有效，那么继续使用当前P，如果当前g的P已经被release，那么选择runtime.allp[0]作为当前g的新P。
更新runtime.allp的状态机为Pidle
当我们更新GOMAXPROCS数量时候，调度器会被加上全局锁，也会出发Stop the World。所以除了系统启动时候调用外，不建议在其余地方触发这个函数，十分影响性能。

3.2 runtime.main
完成核心初始化之后创建并执行main goroutine，其实也就是runtime.main函数。

有一个概念需要区分清楚，前面描述的初始化属于内核层面的初始化，还有一种初始化属于用户层面逻辑初始化，包括runtime包、标准库、用户、第三方包初始化函数。逻辑层面的初始化可能关系到很多的同步关系，代码依赖等等。

函数的核心代码如下：proc.go#main

// The main goroutine.
func main() {
	// Max stack size is 1 GB on 64-bit, 250 MB on 32-bit.
	if sys.PtrSize == 8 {
		maxstacksize = 1000000000
	} else {
		maxstacksize = 250000000
	}

	// Allow newproc to start new Ms.
	mainStarted = true
	
	// 启动后台监控线程
	if GOARCH != "wasm" { // no threads on wasm yet, so no sysmon
		systemstack(func() {
			newm(sysmon, nil)
		})
	}
	// runtime包里面的初始化函数
	runtime_init() // must be before defer
	
	// Record when the world started. Must be after runtime_init
	// because nanotime on some platforms depends on startNano.
	runtimeInitTime = nanotime()
	//启动GC
	gcenable()
	//执行用户、标准库，三方包里面的初始化函数
	fn := main_init // make an indirect call, as the linker doesn't know the address of the main package when laying down the runtime
	fn()
	// 如果是库方式就不执行用户入口函数
	if isarchive || islibrary {
		// A program compiled with -buildmode=c-archive or c-shared
		// has a main, but it is not executed.
		return
	}
	//执行用户入口函数
	fn = main_main // make an indirect call, as the linker doesn't know the address of the main package when laying down the runtime
	fn()
	// 退出进程
	exit(0)
	for {
		var x *int32
		*x = 0
	}
}

有几个需要注意的点：
1）goroutine的栈最大值是1G(64位系统)
2）会启动一个后台监控线程，这个监控线程也是完成抢占式调度的地方
3）目标是runtime.init、main.init、以及main.main这个入口函数

runtime.init函数和main.init由编译器编译生成，负责调用所有的初始化函数。

runtime.init函数仅仅负责runtime包，
main.init函数负责用户、标准库和第三方所有init函数，
所有初始化都在main goroutine中完成
全部初始化完成再执行main.main函数。