---
title: goroutine 启动运行详细过程
layout: post
category: golang
author: 夏泽民
---
https://juejin.im/post/5ec72e0951882542f346e672
它跟线程有啥区别？原理是啥？
都说他好，他好在哪里？
使用上面有啥注意的？
<!-- more -->
package main

import (
	"fmt"
)

func worker(stop chan bool) {
	
	for i:=0;i<10;i++ {
		fmt.Println("干活....")
	}
	stop <- true
}

func main() {
	stop := make(chan bool)
	go worker(stop)
	<- stop
}
复制代码我们在main中新起了一个goroutine来干活。后台实现是runtime.newproc调用,函数体如下
// 使用siz字节参数创建一个运行fn的新g。 
// 将其放在g等待运行的队列中。 编译器将go语句转换为对此的调用。 
// 无法拆分堆栈，因为它假定参数在＆fn;之后顺序可用。 
// 如果发生堆栈拆分，则不会复制它们。
//go:nosplit
func newproc(siz int32, fn *funcval) {
    // 从 fn 的地址增加一个指针的长度，从而获取第一参数地址
	argp := add(unsafe.Pointer(&fn), sys.PtrSize)
	// 获取当前的运行的g
	gp := getg()
	// getcallerpc返回其调用方的程序计数器（PC）。用于存放下一条指令所在单元的地址的地方。
	pc := getcallerpc()
	// systemstack在系统堆栈上运行
	// 如果从每个OS线程（g0）堆栈调用systemstack
	// ，或者从信号处理（gsignal）堆栈调用systemstack ，
	// systemstack直接调用fn并返回。
	// 否则，从普通goroutine的有限堆栈中调用systemstack。
	// 在这种情况下，系统堆栈切换到每个OS线程堆栈，调用fn，然后切回。
	// 通常使用func字面量作为参数，以便与调用系统堆栈周围的代码共享输入和输出
	systemstack(func() {
	    // 原型：func newproc1(fn *funcval, argp *uint8, narg int32, callergp *g, callerpc uintptr) 
	    // 创建一个新的g，运行fn，其中narg个字节的参数从argp开始。
	    // callerpc是创建它的go语句的地址。新g放入g等待运行的队列中。
		newproc1(fn, (*uint8)(argp), siz, gp, pc)
	})
}
复制代码newproc1是重头戏，也比较复杂，可能目前还不能看的很明白，但是，大致先了解一下：
func newproc1(fn *funcval, argp *uint8, narg int32, callergp *g, callerpc uintptr) {
	_g_ := getg()

	if fn == nil {
		_g_.m.throwing = -1 // do not dump full stacks
		throw("go of nil func value")
	}
	_g_.m.locks++ // disable preemption because it can be holding p in a local var
	siz := narg
	siz = (siz + 7) &^ 7

	// We could allocate a larger initial stack if necessary.
	// Not worth it: this is almost always an error.
	// 4*sizeof(uintreg): extra space added below
	// sizeof(uintreg): caller's LR (arm) or return address (x86, in gostartcall).
	if siz >= _StackMin-4*sys.RegSize-sys.RegSize {
		throw("newproc: function arguments too large for new goroutine")
	}

	_p_ := _g_.m.p.ptr() 
	newg := gfget(_p_) // 根据 p 获得一个新的 g
	// 初始化阶段，gfget 是不可能找到 g 的
	// 也可能运行中本来就已经耗尽了
	if newg == nil {
		newg = malg(_StackMin) // 创建一个拥有 _StackMin 大小的栈的 g
		casgstatus(newg, _Gidle, _Gdead) // 将新创建的 g 从 _Gidle 更新为 _Gdead 状态
		allgadd(newg) // 将 Gdead 状态的 g 添加到 allg，这样 GC 不会扫描未初始化的栈
	}
	if newg.stack.hi == 0 {
		throw("newproc1: newg missing stack")
	}

	if readgstatus(newg) != _Gdead {
		throw("newproc1: new g is not Gdead")
	}

	totalSize := 4*sys.RegSize + uintptr(siz) + sys.MinFrameSize // extra space in case of reads slightly beyond frame
	totalSize += -totalSize & (sys.SpAlign - 1)                  // align to spAlign
	sp := newg.stack.hi - totalSize
	spArg := sp
	if usesLR {
		// caller's LR
		*(*uintptr)(unsafe.Pointer(sp)) = 0
		prepGoExitFrame(sp)
		spArg += sys.MinFrameSize
	}
	if narg > 0 {
		memmove(unsafe.Pointer(spArg), unsafe.Pointer(argp), uintptr(narg))
		// This is a stack-to-stack copy. If write barriers
		// are enabled and the source stack is grey (the
		// destination is always black), then perform a
		// barrier copy. We do this *after* the memmove
		// because the destination stack may have garbage on
		// it.
		if writeBarrier.needed && !_g_.m.curg.gcscandone {
			f := findfunc(fn.fn)
			stkmap := (*stackmap)(funcdata(f, _FUNCDATA_ArgsPointerMaps))
			if stkmap.nbit > 0 {
				// We're in the prologue, so it's always stack map index 0.
				bv := stackmapdata(stkmap, 0)
				bulkBarrierBitmap(spArg, spArg, uintptr(bv.n)*sys.PtrSize, 0, bv.bytedata)
			}
		}
	}

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
	if isSystemGoroutine(newg, false) {
		atomic.Xadd(&sched.ngsys, +1)
	}
	newg.gcscanvalid = false
	casgstatus(newg, _Gdead, _Grunnable)

	if _p_.goidcache == _p_.goidcacheend {
		// Sched.goidgen is the last allocated id,
		// this batch must be [sched.goidgen+1, sched.goidgen+GoidCacheBatch].
		// At startup sched.goidgen=0, so main goroutine receives goid=1.
		_p_.goidcache = atomic.Xadd64(&sched.goidgen, _GoidCacheBatch)
		_p_.goidcache -= _GoidCacheBatch - 1
		_p_.goidcacheend = _p_.goidcache + _GoidCacheBatch
	}
	newg.goid = int64(_p_.goidcache)
	_p_.goidcache++
	if raceenabled {
		newg.racectx = racegostart(callerpc)
	}
	if trace.enabled {
		traceGoCreate(newg, newg.startpc)
	}
	runqput(_p_, newg, true)

	if atomic.Load(&sched.npidle) != 0 && atomic.Load(&sched.nmspinning) == 0 && mainStarted {
		wakep()
	}
	_g_.m.locks--
	if _g_.m.locks == 0 && _g_.preempt { // restore the preemption request in case we've cleared it in newstack
		_g_.stackguard0 = stackPreempt
	}
}
复制代码也就是说，刚开始的时候，p上并没有可以使用的g,所以创建了一个具有很少栈容量的g.
// 分配一个新的g，其堆栈足以容纳stacksize字节。
func malg(stacksize int32) *g {
	newg := new(g)
	if stacksize >= 0 {
		stacksize = round2(_StackSystem + stacksize)
		systemstack(func() {
			newg.stack = stackalloc(uint32(stacksize))
		})
		newg.stackguard0 = newg.stack.lo + _StackGuard
		newg.stackguard1 = ^uintptr(0)
	}
	return newg
}
复制代码分配的g是一个结构体指针,如果stacksize大于零，还将分配stack堆栈，该结构体具体内容如下：
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
	m              *m      // current m; offset known to arm liblink
	sched          gobuf
	syscallsp      uintptr        // if status==Gsyscall, syscallsp = sched.sp to use during gc
	syscallpc      uintptr        // if status==Gsyscall, syscallpc = sched.pc to use during gc
	stktopsp       uintptr        // expected sp at top of stack, to check in traceback
	param          unsafe.Pointer // passed parameter on wakeup
	atomicstatus   uint32
	stackLock      uint32 // sigprof/scang lock; TODO: fold in to atomicstatus
	goid           int64
	schedlink      guintptr
	waitsince      int64      // approx time when the g become blocked
	waitreason     waitReason // if status==Gwaiting
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
	lockedm        muintptr
	sig            uint32
	writebuf       []byte
	sigcode0       uintptr
	sigcode1       uintptr
	sigpc          uintptr
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
复制代码东西太多，目前能看懂的就是new过g后，分配了一个round2(_StackSystem + stacksize)个字节的stack.
newg.stackguard0 = newg.stack.lo + _StackGuard
newg.stackguard1 = ^uintptr(0)
复制代码然后将新生成的g的状态由_Gidle变成_Gdead。将 Gdead 状态的 g 添加到 allg切片中。
var (
	allgs    []*g
	allglock mutex
)

func allgadd(gp *g) {
	if readgstatus(gp) == _Gidle {
		throw("allgadd: bad status Gidle")
	}

	lock(&allglock)
	allgs = append(allgs, gp)
	allglen = uintptr(len(allgs))
	unlock(&allglock)
}
复制代码之后对g相关的sched字段进行初始化赋值,该字段类型是个结构体，
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
	sp   uintptr
	pc   uintptr
	g    guintptr
	ctxt unsafe.Pointer
	ret  sys.Uintreg
	lr   uintptr
	bp   uintptr // for GOEXPERIMENT=framepointer
}
复制代码该字段的功能，目前我们不得而知，先看
	memclrNoHeapPointers(unsafe.Pointer(&newg.sched), unsafe.Sizeof(newg.sched))
	newg.sched.sp = sp
	newg.stktopsp = sp
	newg.sched.pc = funcPC(goexit) + sys.PCQuantum // +PCQuantum so that previous instruction is in same function
	newg.sched.g = guintptr(unsafe.Pointer(newg))
	gostartcallfn(&newg.sched, fn)
复制代码调整Gobuf就像执行对fn的调用一样，然后立即执行gosave.
func gostartcallfn(gobuf *gobuf, fv *funcval) {
	var fn unsafe.Pointer
	if fv != nil {
		fn = unsafe.Pointer(fv.fn)
	} else {
		fn = unsafe.Pointer(funcPC(nilfunc))
	}
	gostartcall(gobuf, fn, unsafe.Pointer(fv))
}
复制代码之后有一个将当前g的状态调整的动作
casgstatus(newg, _Gdead, _Grunnable)
复制代码可运行状态的g会被放入到本地的可运行队列中，
runqput(_p_, newg, true)
复制代码该函数体如下：
// runqput尝试将g放置在本地可运行队列中。 
// 如果next为false，则runqput将g添加到可运行队列的尾部。
// 如果next为true，则runqput将g放在_p_.runnext插槽中。 
// 如果运行队列已满，则runnext将g放入全局队列。 
// 仅由所有者P执行。
func runqput(_p_ *p, gp *g, next bool) {
	if randomizeScheduler && next && fastrand()%2 == 0 {
		next = false
	}

	if next {
	retryNext:
		oldnext := _p_.runnext
		if !_p_.runnext.cas(oldnext, guintptr(unsafe.Pointer(gp))) {
			goto retryNext
		}
		if oldnext == 0 {
			return
		}
		// Kick the old runnext out to the regular run queue.
		gp = oldnext.ptr()
	}

retry:
	h := atomic.LoadAcq(&_p_.runqhead) // load-acquire, synchronize with consumers
	t := _p_.runqtail
	if t-h < uint32(len(_p_.runq)) {
		_p_.runq[t%uint32(len(_p_.runq))].set(gp)
		atomic.StoreRel(&_p_.runqtail, t+1) // store-release, makes the item available for consumption
		return
	}
	if runqputslow(_p_, gp, h, t) {
		return
	}
	// the queue is not full, now the put above must succeed
	goto retry
}
复制代码以上，关于g的内容，我们有了一个大致的了解，当我们将创建的g放到本地队列时，提到了一个结构体p,这个东西是什么呢？下面是他的结构体
type p struct {
	lock mutex

	id          int32
	status      uint32 // one of pidle/prunning/...
	link        puintptr
	schedtick   uint32     // incremented on every scheduler call
	syscalltick uint32     // incremented on every system call
	sysmontick  sysmontick // last tick observed by sysmon
	m           muintptr   // back-link to associated m (nil if idle)
	mcache      *mcache
	racectx     uintptr

	deferpool    [5][]*_defer // pool of available defer structs of different sizes (see panic.go)
	deferpoolbuf [5][32]*_defer

	// Cache of goroutine ids, amortizes accesses to runtime·sched.goidgen.
	goidcache    uint64
	goidcacheend uint64

	// Queue of runnable goroutines. Accessed without lock.
	// 可运行goroutines的队列，访问无需锁，这个就是我们上述创建的g存放的位置
	runqhead uint32
	runqtail uint32
	runq     [256]guintptr
	// runnext（如果不是nil）是当前G准备好的可运行G，
	// 如果正在运行的G的时间片中还有剩余时间，则应下一个运行，而不是从runq中获取G。
	// 它将继承当前时间片中剩余的时间。
	// 如果将一组goroutine锁定为通信等待模式，
	// 则此调度会将其设置为一个单元，
	// 并消除（可能很大的）调度延迟，
	// 否则该延迟可能是由于将就绪的goroutine添加到运行队列的末尾而引起的。
	runnext guintptr

	// Available G's (status == Gdead)
	gFree struct {
		gList
		n int32
	}

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

	pad cpu.CacheLinePad
}
复制代码在newproc函数中，从当前g获取p结构时，通过的是g的m字段，该字段是个什么呢？是个m结构体指针，m的结构体原型为：
type m struct {
	g0      *g     // 用于执行调度指令的 goroutine
	morebuf gobuf  // gobuf arg to morestack
	divmod  uint32 // div/mod denominator for arm - known to liblink

	// Fields not known to debuggers.
	procid        uint64       // for debuggers, but offset not hard-coded
	gsignal       *g           // 处理 signal 的 g
	goSigStack    gsignalStack // Go-allocated signal handling stack
	sigmask       sigset       // storage for saved signal mask
	tls           [6]uintptr   // 线程本地存储
	mstartfn      func()
	curg          *g       // 当前运行的G
	caughtsig     guintptr // goroutine running during fatal signal
	p             puintptr // 执行 go 代码时持有的 p (如果没有执行则为 nil)
	nextp         puintptr
	oldp          puintptr // the p that was attached before executing a syscall
	id            int64
	mallocing     int32
	throwing      int32
	preemptoff    string // if != "", keep curg running on this m
	locks         int32
	dying         int32
	profilehz     int32
	spinning      bool // m 当前没有运行 work 且正处于寻找 work 的活跃状态
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
	alllink       *m // on allm
	schedlink     muintptr
	mcache        *mcache
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
复制代码看了上面的结构体感觉很空洞，都是些什么呢？就知道newproc时，创建的G,放到了关联的P的本地可运行队列中，要明白这些东西是什么，就要从他们是如何产生的说起？
➜  goroutinetest gdb main
GNU gdb (GDB) 8.3
Copyright (C) 2019 Free Software Foundation, Inc.
License GPLv3+: GNU GPL version 3 or later <http://gnu.org/licenses/gpl.html>
This is free software: you are free to change and redistribute it.
There is NO WARRANTY, to the extent permitted by law.
Type "show copying" and "show warranty" for details.
This GDB was configured as "x86_64-apple-darwin16.7.0".
Type "show configuration" for configuration details.
For bug reporting instructions, please see:
<http://www.gnu.org/software/gdb/bugs/>.
Find the GDB manual and other documentation resources online at:
    <http://www.gnu.org/software/gdb/documentation/>.

For help, type "help".
Type "apropos word" to search for commands related to "word"...
Reading symbols from main...
(No debugging symbols found in main)
Loading Go Runtime support.
(gdb) info files
Symbols from "/Users/zhaojunwei/workspace/src/just.for.test/goroutinetest/main".
Local exec file:
	`/Users/zhaojunwei/workspace/src/just.for.test/goroutinetest/main', file type mach-o-x86-64.
	Entry point: 0x1052770
	0x0000000001001000 - 0x0000000001093194 is .text
	0x00000000010931a0 - 0x00000000010e1ace is __TEXT.__rodata
	0x00000000010e1ae0 - 0x00000000010e1be2 is __TEXT.__symbol_stub1
	0x00000000010e1c00 - 0x00000000010e2864 is __TEXT.__typelink
	0x00000000010e2868 - 0x00000000010e28d0 is __TEXT.__itablink
	0x00000000010e28d0 - 0x00000000010e28d0 is __TEXT.__gosymtab
	0x00000000010e28e0 - 0x000000000115c108 is __TEXT.__gopclntab
	0x000000000115d000 - 0x000000000115d158 is __DATA.__nl_symbol_ptr
	0x000000000115d160 - 0x0000000001169c9c is __DATA.__noptrdata
	0x0000000001169ca0 - 0x0000000001170610 is .data
	0x0000000001170620 - 0x000000000118be50 is .bss
	0x000000000118be60 - 0x000000000118e418 is __DATA.__noptrbss
(gdb)
(gdb) b *0x1052770
Breakpoint 1 at 0x1052770
(gdb) info br
Num     Type           Disp Enb Address            What
1       breakpoint     keep y   0x0000000001052770 <_rt0_amd64_darwin>
(gdb)
复制代码查看一下_rt0_amd64_darwin是什么？
#include "textflag.h"

TEXT _rt0_amd64_darwin(SB),NOSPLIT,$-8
	JMP	_rt0_amd64(SB)

// When linking with -shared, this symbol is called when the shared library
// is loaded.
TEXT _rt0_amd64_darwin_lib(SB),NOSPLIT,$0
	JMP	_rt0_amd64_lib(SB)

复制代码_rt0_amd64是使用内部链接时大多数amd64系统的通用启动代码。 这是内核中普通-buildmode = exe程序的程序入口点。 堆栈保存参数数量和C风格的argv。
TEXT _rt0_amd64(SB),NOSPLIT,$-8
	MOVQ	0(SP), DI	// argc
	LEAQ	8(SP), SI	// argv
	JMP	runtime·rt0_go(SB)
复制代码最终调用的是runtime.rt0_go方法
TEXT runtime·rt0_go(SB),NOSPLIT,$0
	// SP = stack; R0 = argc; R1 = argv

	SUB	$32, RSP
	MOVW	R0, 8(RSP) // argc
	MOVD	R1, 16(RSP) // argv

	// create istack out of the given (operating system) stack.
	// _cgo_init may update stackguard.
	MOVD	$runtime·g0(SB), g
	MOVD	RSP, R7
	MOVD	$(-64*1024)(R7), R0
	MOVD	R0, g_stackguard0(g)
	MOVD	R0, g_stackguard1(g)
	MOVD	R0, (g_stack+stack_lo)(g)
	MOVD	R7, (g_stack+stack_hi)(g)

	// if there is a _cgo_init, call it using the gcc ABI.
	MOVD	_cgo_init(SB), R12
	CMP	$0, R12
	BEQ	nocgo

	MRS_TPIDR_R0			// load TLS base pointer
	MOVD	R0, R3			// arg 3: TLS base pointer
#ifdef TLSG_IS_VARIABLE
	MOVD	$runtime·tls_g(SB), R2 	// arg 2: &tls_g
#else
	MOVD	$0, R2		        // arg 2: not used when using platform's TLS
#endif
	MOVD	$setg_gcc<>(SB), R1	// arg 1: setg
	MOVD	g, R0			// arg 0: G
	SUB	$16, RSP		// reserve 16 bytes for sp-8 where fp may be saved.
	BL	(R12)
	ADD	$16, RSP

nocgo:
	BL	runtime·save_g(SB)
	// update stackguard after _cgo_init
	MOVD	(g_stack+stack_lo)(g), R0
	ADD	$const__StackGuard, R0
	MOVD	R0, g_stackguard0(g)
	MOVD	R0, g_stackguard1(g)

	// set the per-goroutine and per-mach "registers"
	MOVD	$runtime·m0(SB), R0

	// save m->g0 = g0
	MOVD	g, m_g0(R0)
	// save m0 to g0->m
	MOVD	R0, g_m(g)

	BL	runtime·check(SB)

	MOVW	8(RSP), R0	// copy argc
	MOVW	R0, -8(RSP)
	MOVD	16(RSP), R0		// copy argv
	MOVD	R0, 0(RSP)
	BL	runtime·args(SB)
	BL	runtime·osinit(SB)
	BL	runtime·schedinit(SB)

	// create a new goroutine to start program
	MOVD	$runtime·mainPC(SB), R0		// entry
	MOVD	RSP, R7
	MOVD.W	$0, -8(R7)
	MOVD.W	R0, -8(R7)
	MOVD.W	$0, -8(R7)
	MOVD.W	$0, -8(R7)
	MOVD	R7, RSP
	BL	runtime·newproc(SB)
	ADD	$32, RSP

	// start this M
	BL	runtime·mstart(SB)

	MOVD	$0, R0
	MOVD	R0, (R0)	// boom
	UNDEF
复制代码首先进行g0和m0的初始化，之后进行本地线程存储的检测设置。之后尽心调度器的初始化，并创建一个新的goroutine运行程序，最后开启我们的M.
// The bootstrap sequence is:
//
//	call osinit
//	call schedinit
//	make & queue new G
//	call runtime·mstart
//
// The new G calls runtime·main.
func schedinit() {
	// raceinit must be the first call to race detector.
	// In particular, it must be done before mallocinit below calls racemapshadow.
	_g_ := getg()
	if raceenabled {
		_g_.racectx, raceprocctx0 = raceinit()
	}
	// 设置最多启动10000个操作系统线程，也是最多10000个M
	sched.maxmcount = 10000

	tracebackinit()
	moduledataverify()
	stackinit()
	mallocinit()
	mcommoninit(_g_.m) // 初始化m0，因为从前面的代码我们知道g0->m = &m0
	cpuinit()       // must run before alginit
	alginit()       // maps must not be used before this call
	modulesinit()   // provides activeModules
	typelinksinit() // uses maps, activeModules
	itabsinit()     // uses activeModules

	msigsave(_g_.m)
	initSigmask = _g_.m.sigmask

	goargs()
	goenvs()
	parsedebugvars()
	gcinit()

	sched.lastpoll = uint64(nanotime())
	// 系统中有多少核，就创建和初始化多少个p结构体对象
	procs := ncpu
	if n, ok := atoi32(gogetenv("GOMAXPROCS")); ok && n > 0 {
	    // 如果环境变量指定了GOMAXPROCS，则创建指定数量的p
		procs = n
	}
	// 创建和初始化全局变量allp
	if procresize(procs) != nil {
		throw("unknown runnable goroutine during bootstrap")
	}

	// For cgocheck > 1, we turn on the write barrier at all times
	// and check all pointer writes. We can't do this until after
	// procresize because the write barrier needs a P.
	if debug.cgocheck > 1 {
		writeBarrier.cgo = true
		writeBarrier.enabled = true
		for _, p := range allp {
			p.wbBuf.reset()
		}
	}

	if buildVersion == "" {
		// Condition should never trigger. This code just serves
		// to ensure runtime·buildVersion is kept in the resulting binary.
		buildVersion = "unknown"
	}
}
复制代码我们来关注一下m0是如何初始化的
func mcommoninit(mp *m) {
	_g_ := getg()

	// g0 stack won't make sense for user (and is not necessary unwindable).
	if _g_ != _g_.m.g0 {
		callers(1, mp.createstack[:])
	}

	lock(&sched.lock)
	if sched.mnext+1 < sched.mnext {
		throw("runtime: thread ID overflow")
	}
	// m0分配的id,schedt结构体的mnext字段标识下一个可用的thread id.
	mp.id = sched.mnext
	sched.mnext++
	
	checkmcount()

	mp.fastrand[0] = 1597334677 * uint32(mp.id)
	mp.fastrand[1] = uint32(cputicks())
	if mp.fastrand[0]|mp.fastrand[1] == 0 {
		mp.fastrand[1] = 1
	}

	mpreinit(mp)
	if mp.gsignal != nil {
		mp.gsignal.stackguard1 = mp.gsignal.stack.lo + _StackGuard
	}

	// Add to allm so garbage collector doesn't free g->m
	// when it is just in a register or thread-local storage.
	// allm挂到这里，防止被垃圾回收
	mp.alllink = allm

	// NumCgoCall() iterates over allm w/o schedlock,
	// so we need to publish it safely.
	atomicstorep(unsafe.Pointer(&allm), unsafe.Pointer(mp))
	unlock(&sched.lock)

	// Allocate memory to hold a cgo traceback if the cgo call crashes.
	if iscgo || GOOS == "solaris" || GOOS == "windows" {
		mp.cgoCallers = new(cgoCallers)
	}
}
复制代码调度器初始化最后一部分工作就是p的初始化

初始化调度后，开启新的goroutine运行我们的主程序，然后调用runtime.mstart开启M.
func mstart() {
	_g_ := getg()

    // 通过检查 g 执行占的边界来确定是否为系统栈
	osStack := _g_.stack.lo == 0
	if osStack {
		// Initialize stack bounds from system stack.
		// Cgo may have left stack size in stack.hi.
		// minit may update the stack bounds.
		size := _g_.stack.hi
		if size == 0 {
			size = 8192 * sys.StackGuardMultiplier
		}
		_g_.stack.hi = uintptr(noescape(unsafe.Pointer(&size)))
		_g_.stack.lo = _g_.stack.hi - size + 1024
	}
	// Initialize stack guards so that we can start calling
	// both Go and C functions with stack growth prologues.
	_g_.stackguard0 = _g_.stack.lo + _StackGuard
	_g_.stackguard1 = _g_.stackguard0
	// 启动m
	mstart1()

	// Exit this thread.
	if GOOS == "windows" || GOOS == "solaris" || GOOS == "plan9" || GOOS == "darwin" || GOOS == "aix" {
		// Window, Solaris, Darwin, AIX and Plan 9 always system-allocate
		// the stack, but put it in _g_.stack before mstart,
		// so the logic above hasn't set osStack yet.
		osStack = true
	}
	mexit(osStack)
}


func mstart1() {
	_g_ := getg()

	if _g_ != _g_.m.g0 {
		throw("bad runtime·mstart")
	}

	// Record the caller for use as the top of stack in mcall and
	// for terminating the thread.
	// We're never coming back to mstart1 after we call schedule,
	// so other calls can reuse the current frame.
	save(getcallerpc(), getcallersp())
	asminit()
	minit()

	// Install signal handlers; after minit so that minit can
	// prepare the thread to be able to handle the signals.
	if _g_.m == &m0 {
		mstartm0()
	}

	if fn := _g_.m.mstartfn; fn != nil {
		fn()
	}
    // 如果当前 m 并非 m0，则要求绑定 p
	if _g_.m != &m0 {
		acquirep(_g_.m.nextp.ptr())
		_g_.m.nextp = 0
	}
	schedule()
}
复制代码在mstart1中，调用了schedule函数：一轮调度程序：找到一个可运行的goroutine并执行它。永不return.
func schedule() {
	_g_ := getg()

	if _g_.m.locks != 0 {
		throw("schedule: holding locks")
	}

	if _g_.m.lockedg != 0 {
		stoplockedm()
		execute(_g_.m.lockedg.ptr(), false) // Never returns.
	}

	// We should not schedule away from a g that is executing a cgo call,
	// since the cgo call is using the m's g0 stack.
	if _g_.m.incgo {
		throw("schedule: in cgo")
	}

top:
	if sched.gcwaiting != 0 {
		gcstopm()
		goto top
	}
	if _g_.m.p.ptr().runSafePointFn != 0 {
		runSafePointFn()
	}

	var gp *g
	var inheritTime bool
	if trace.enabled || trace.shutdown {
		gp = traceReader()
		if gp != nil {
			casgstatus(gp, _Gwaiting, _Grunnable)
			traceGoUnpark(gp, 0)
		}
	}
	if gp == nil && gcBlackenEnabled != 0 {
		gp = gcController.findRunnableGCWorker(_g_.m.p.ptr())
	}
	if gp == nil {
		// // 说明不在 GC
		//
		// 每调度 61 次，就检查一次全局队列，保证公平性
		// 否则两个 goroutine 可以通过互相 respawn 一直占领本地的 runqueue
		if _g_.m.p.ptr().schedtick%61 == 0 && sched.runqsize > 0 {
			lock(&sched.lock)
			// 从全局队列中偷 g
			gp = globrunqget(_g_.m.p.ptr(), 1)
			unlock(&sched.lock)
		}
	}
	if gp == nil {
		gp, inheritTime = runqget(_g_.m.p.ptr())
		if gp != nil && _g_.m.spinning {
			throw("schedule: spinning with local work")
		}
	}
	if gp == nil {
		gp, inheritTime = findrunnable() // 如果偷都偷不到，则休眠，在此阻塞
	}

	// 该线程将运行goroutine，并且不再自旋，
	// 因此，如果将其标记为正在自旋，则需要立即将其重置并可能启动新自旋的M。
	if _g_.m.spinning {
		resetspinning()
	}

	if sched.disable.user && !schedEnabled(gp) {
		// Scheduling of this goroutine is disabled. Put it on
		// the list of pending runnable goroutines for when we
		// re-enable user scheduling and look again.
		lock(&sched.lock)
		if schedEnabled(gp) {
			// Something re-enabled scheduling while we
			// were acquiring the lock.
			unlock(&sched.lock)
		} else {
			sched.disable.runnable.pushBack(gp)
			sched.disable.n++
			unlock(&sched.lock)
			goto top
		}
	}

	if gp.lockedm != 0 {
		// Hands off own p to the locked m,
		// then blocks waiting for a new p.
		startlockedm(gp)
		goto top
	}

    // 开始执行
	execute(gp, inheritTime)
}

复制代码如果m处在自旋的状态，那么将调用resetspinning方法，
func resetspinning() {
	_g_ := getg()
	if !_g_.m.spinning {
		throw("resetspinning: not a spinning m")
	}
	_g_.m.spinning = false
	nmspinning := atomic.Xadd(&sched.nmspinning, -1)
	if int32(nmspinning) < 0 {
		throw("findrunnable: negative nmspinning")
	}
	// M的唤醒策略故意有些保守，因此请检查是否需要在此处唤醒另一个P。
	// 有关详细信息，请参见文件顶部的“工作线程park/unpark”注释。
	if nmspinning == 0 && atomic.Load(&sched.npidle) > 0 {
		wakep()
	}
}
复制代码wakep()尝试再添加一个P以执行G。 当G变为可运行时调用（newproc，就绪）.该函数会调用startm(nil, true).startm函数调度一些M以运行p（必要时创建M）。 如果p == nil，则尝试获取一个空闲P，如果没有空闲P则不执行任何操作。 可以与m.p == nil一起运行，因此不允许写障碍。 如果设置了旋转，则调用者已增加nmspinning，并且startm将减少nmspinning或在新启动的M中设置m.spinning。