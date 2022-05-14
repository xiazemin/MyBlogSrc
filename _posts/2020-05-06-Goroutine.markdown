---
title: Goroutine的创建与销毁
layout: post
category: golang
author: 夏泽民
---
https://www.jianshu.com/p/181dc7845bb8

<!-- more -->
下面的分析均基于Golang1.14版本。
go func(){} 只是一个语法糖，在编译时会替换为newproc函数。

一、创建---newproc
阅读建议：g的创建涉及的数据结构主要有g, p的结构体和全局数据结构allgs, sched，阅读时对照这些数据结构阅读源码。
newproc的调用过程：在newproc中切换到g0栈执行newproc1.

//go:nosplit
func newproc(siz int32, fn *funcval) { // fn表示要运行的函数 siz表示要运行的函数的总参数的大小
    argp := add(unsafe.Pointer(&fn), sys.PtrSize) // argp表示fn参数存放的指针 该指针紧跟在fn后面
    gp := getg()                                  // gp 当前正在运行的g 非g0
    pc := getcallerpc()                           // 获取调用当前函数的函数的pc值 即当前函数返回后的下一个指令
    systemstack(func() {                          // 切换到g0栈
        newproc1(fn, argp, siz, gp, pc)
    })
}

func newproc1(fn *funcval, argp unsafe.Pointer, narg int32, callergp *g, callerpc uintptr) {
    _g_ := getg() // 当前正在运行的g 即g0

    if fn == nil {
        _g_.m.throwing = -1 // do not dump full stacks
        throw("go of nil func value")
    }
    acquirem() // disable preemption because it can be holding p in a local var  为当前m加锁 避免preempt
    siz := narg
    siz = (siz + 7) &^ 7 // narg 表示总参数大小 siz为narg内存对齐后的值

    // We could allocate a larger initial stack if necessary.
    // Not worth it: this is almost always an error.
    // 4*sizeof(uintreg): extra space added below
    // sizeof(uintreg): caller's LR (arm) or return address (x86, in gostartcall).
    if siz >= _StackMin-4*sys.RegSize-sys.RegSize { // 判断参数长度是否溢出 g的出事大小为_StackMin并且初始化时还需要部分内存存储其他参数
        throw("newproc: function arguments too large for new goroutine")
    }

    _p_ := _g_.m.p.ptr()
    // 获取新的g 先尝试从当前p空闲列表获取 如果获取不到则创建 具体的不深究
    newg := gfget(_p_)
    if newg == nil {
        newg = malg(_StackMin) // 从空闲列表获取g失败 尝试创建新的g 传入的参数_StackMin表示g的栈大小
        casgstatus(newg, _Gidle, _Gdead)
        // 将g放入 allgs切片（allgs管理所有的g）
        allgadd(newg) // publishes with a g->status of Gdead so GC scanner doesn't look at uninitialized stack.
    }
    if newg.stack.hi == 0 { // 如果栈的高地址为0 说明g的栈内存分配失败 抛出异常
        throw("newproc1: newg missing stack")
    }

    if readgstatus(newg) != _Gdead { // 判断g的状态
        throw("newproc1: new g is not Gdead")
    }

    // 计算fn参数所需的空间大小 额外申请了一些空间 具体作用未知
    totalSize := 4*sys.RegSize + uintptr(siz) + sys.MinFrameSize // extra space in case of reads slightly beyond frame
    // 对齐的原理 以8位对齐举例 sys.SpAlign - 1表示对齐位后面的bit均为1 即0111 减去 得到的数&totalSize 表示将totalSize中小于SpAlign的部分减掉
    totalSize += -totalSize & (sys.SpAlign - 1) // align to spAlign
    sp := newg.stack.hi - totalSize // 根据栈的起始地址和传入参数大小计算g栈的sp的值 注意栈增长是高地址向低地址增长
    spArg := sp // 此时 newg.stack.hi -> sp(spArg)
    if usesLR { // 这部分代码暂时看不懂 略过
        // caller's LR
        *(*uintptr)(unsafe.Pointer(sp)) = 0
        prepGoExitFrame(sp)
        spArg += sys.MinFrameSize
    }
    if narg > 0 {
        // 将传入的参数拷贝到g的栈中
        memmove(unsafe.Pointer(spArg), argp, uintptr(narg)) //将参数从argp拷贝到 spArg --> spArg + narg
        // This is a stack-to-stack copy. If write barriers
        // are enabled and the source stack is grey (the
        // destination is always black), then perform a
        // barrier copy. We do this *after* the memmove
        // because the destination stack may have garbage on
        // it.
        // 如果正在GC copy stack会触发写屏障 具体的操作在GC中分析
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

    // 初始化g的调度的上下文信息sched
    memclrNoHeapPointers(unsafe.Pointer(&newg.sched), unsafe.Sizeof(newg.sched))
    newg.sched.sp = sp
    newg.stktopsp = sp
    // 设置g的 pc为goexit函数+1 +1的原因参考goexit实现
    newg.sched.pc = funcPC(goexit) + sys.PCQuantum // +PCQuantum so that previous instruction is in same function
    newg.sched.g = guintptr(unsafe.Pointer(newg))
    // 相当于执行一次gosave 相当于是有个虚拟的函数f f先执行fn再执行goexit 在执行fn前保存上下文
    gostartcallfn(&newg.sched, fn)
    // 记录调用newproc的 pc值 父g 和初始函数
    newg.gopc = callerpc
    newg.ancestors = saveAncestors(callergp)
    newg.startpc = fn.fn
    if _g_.m.curg != nil {
        newg.labels = _g_.m.curg.labels // 从父g中继承labels labels作用?
    }

    // 根据函数名称是否以runtime开头判断是否为系统函数
    if isSystemGoroutine(newg, false) {
        atomic.Xadd(&sched.ngsys, +1)
    }
    casgstatus(newg, _Gdead, _Grunnable)

    // 如果当前p的goid缓存用完 则再分配
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
    runqput(_p_, newg, true) // 将g放入p的空闲列表中 true表示放入p中的next中 会在下一次调度中被执行

    // 如果当前有p空闲 并且当前没有正在自旋的m(执行findrunnable的m)  且maingoroutine已经初始化完成
    if atomic.Load(&sched.npidle) != 0 && atomic.Load(&sched.nmspinning) == 0 && mainStarted {
        wakep()
    }
    releasem(_g_.m) // 释放m上的锁
}
g的获取途径有2种，一种是从g的空闲列表中获取（gfget函数），一种是重新分配（malg函数）。

// Get from gfree list.
// If local list is empty, grab a batch from global list.
func gfget(_p_ *p) *g {
retry:
    // 如果P缓存的g为空 且全局的g空闲列表不为空 则尝试在全局的空闲列表中获取g
    if _p_.gFree.empty() && (!sched.gFree.stack.empty() || !sched.gFree.noStack.empty()) {
        lock(&sched.gFree.lock)
        // Move a batch of free Gs to the P.
        for _p_.gFree.n < 32 {
            // Prefer Gs with stacks.
            gp := sched.gFree.stack.pop()
            if gp == nil {
                gp = sched.gFree.noStack.pop()
                if gp == nil {
                    break
                }
            }
            sched.gFree.n--
            _p_.gFree.push(gp)
            _p_.gFree.n++
        }
        unlock(&sched.gFree.lock)
        goto retry
    }
    gp := _p_.gFree.pop()
    if gp == nil { // 如果获取g失败 则返回nil 在外面调用malg函数分配g
        return nil
    }
    _p_.gFree.n--
    // 如果分配的g中栈为空 则为其分配栈 (旧的g中的栈可能未释放)
    if gp.stack.lo == 0 {
        // Stack was deallocated in gfput. Allocate a new one.
        systemstack(func() {
            gp.stack = stackalloc(_FixedStack)
        })
        gp.stackguard0 = gp.stack.lo + _StackGuard
    } else {
        if raceenabled {
            racemalloc(unsafe.Pointer(gp.stack.lo), gp.stack.hi-gp.stack.lo)
        }
        if msanenabled {
            msanmalloc(unsafe.Pointer(gp.stack.lo), gp.stack.hi-gp.stack.lo)
        }
    }
    return gp
}

// Allocate a new g, with a stack big enough for stacksize bytes.
func malg(stacksize int32) *g {
    newg := new(g)
    if stacksize >= 0 {
        // round2把数值向上调整为2的幂次
        stacksize = round2(_StackSystem + stacksize)
        systemstack(func() {
            // 主要是栈大小的调整和栈内存的实际分配
            newg.stack = stackalloc(uint32(stacksize))
        })
        newg.stackguard0 = newg.stack.lo + _StackGuard
        newg.stackguard1 = ^uintptr(0)
        // Clear the bottom word of the stack. We record g
        // there on gsignal stack during VDSO on ARM and ARM64.
        *(*uintptr)(unsafe.Pointer(newg.stack.lo)) = 0
    }
    return newg
}
将当前的pc设置为goexit函数的pc值后，需要调用gostartcallfn保存一次当前的上下文。

// adjust Gobuf as if it executed a call to fn
// and then did an immediate gosave.
func gostartcallfn(gobuf *gobuf, fv *funcval) {
    var fn unsafe.Pointer
    // 获取函数指针
    if fv != nil {
        fn = unsafe.Pointer(fv.fn)
    } else {
        fn = unsafe.Pointer(funcPC(nilfunc))
    }
    gostartcall(gobuf, fn, unsafe.Pointer(fv))
}
// adjust Gobuf as if it executed a call to fn with context ctxt
// and then did an immediate gosave.
func gostartcall(buf *gobuf, fn, ctxt unsafe.Pointer) {
    sp := buf.sp
    if sys.RegSize > sys.PtrSize {  // 这一段if没看懂
        sp -= sys.PtrSize
        *(*uintptr)(unsafe.Pointer(sp)) = 0
    }
    sp -= sys.PtrSize // 存入goexit指令后 偏移sp
    *(*uintptr)(unsafe.Pointer(sp)) = buf.pc  // 在这里将当前pc存入栈中
    buf.sp = sp
    buf.pc = uintptr(fn) // 将pc设置为fn函数的入口
    buf.ctxt = ctxt // ctxt表示当前正在执行的 funcval
}
初始化后的g处于runnable状态，通过runqput放入到队列中运行。

// runqput tries to put g on the local runnable queue.
// If next is false, runqput adds g to the tail of the runnable queue.
// If next is true, runqput puts g in the _p_.runnext slot.
// If the run queue is full, runnext puts g on the global queue.
// Executed only by the owner P.
func runqput(_p_ *p, gp *g, next bool) {
    // 如果是 randomizeScheduler 状态且50%概率随机 即使next为true也设置为false
    if randomizeScheduler && next && fastrand()%2 == 0 {
        next = false
    }

    if next {
    retryNext:
        // 将p的next设置为新生成的g 如果当前的next不为空 则尝试放入p中缓存的队列
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
    // 如果p中可运行g的缓存队列未满 则放入缓存队列中
    if t-h < uint32(len(_p_.runq)) {
        _p_.runq[t%uint32(len(_p_.runq))].set(gp)
        atomic.StoreRel(&_p_.runqtail, t+1) // store-release, makes the item available for consumption
        return
    }
    // 放入全局的可运行的g的缓存队列中
    if runqputslow(_p_, gp, h, t) {
        return
    }
    // the queue is not full, now the put above must succeed
    goto retry
}
二、销毁---goexit
// The top-most function running on a goroutine
// returns to goexit+PCQuantum.
TEXT runtime·goexit(SB),NOSPLIT|NOFRAME|TOPFRAME,$0-0
    MOVD    R0, R0  // NOP 空操作符 占一个机器字节
    BL  runtime·goexit1(SB) // does not return
在newproc1中设置返回的pc值时，newg.sched.pc = funcPC(goexit) + sys.PCQuantum。该pc值取的是goexit+一个指令大小，应当是BL runtime.goexit1(SB)指令，说明实际执行goexit1。

func goexit1() {
    if raceenabled {
        racegoend()
    }
    if trace.enabled {
        traceGoEnd()
    }

    // 使用mcall调用goexit0 切换当前的栈为g0栈
    mcall(goexit0)
}

// goexit continuation on g0.
func goexit0(gp *g) {
    _g_ := getg()

    // g 状态转换
    casgstatus(gp, _Grunning, _Gdead)
    if isSystemGoroutine(gp, false) {
        atomic.Xadd(&sched.ngsys, -1)
    }
    // 将引用的参数设置为0
    gp.m = nil
    locked := gp.lockedm != 0
    gp.lockedm = 0
    _g_.m.lockedg = 0
    gp.preemptStop = false
    gp.paniconfault = false
    gp._defer = nil // should be true already but just in case.
    gp._panic = nil // non-nil for Goexit during panic. points at stack-allocated data.
    gp.writebuf = nil
    gp.waitreason = 0
    gp.param = nil
    gp.labels = nil
    gp.timer = nil

    // 如果正在GC且和GC相关的数据不为0
    if gcBlackenEnabled != 0 && gp.gcAssistBytes > 0 {
        // Flush assist credit to the global pool. This gives
        // better information to pacing if the application is
        // rapidly creating an exiting goroutines.
        scanCredit := int64(gcController.assistWorkPerByte * float64(gp.gcAssistBytes))
        atomic.Xaddint64(&gcController.bgScanCredit, scanCredit)
        gp.gcAssistBytes = 0
    }

    // 解除g->m  m->g的引用
    dropg()

    // 一些特殊的状态判断
    if GOARCH == "wasm" { // no threads yet on wasm
        gfput(_g_.m.p.ptr(), gp)
        schedule() // never returns
    }

    if _g_.m.lockedInt != 0 {
        print("invalid m->lockedInt = ", _g_.m.lockedInt, "\n")
        throw("internal lockOSThread error")
    }
    // 释放g 将g放回空闲列表并且释放g->stack指向的内存
    gfput(_g_.m.p.ptr(), gp)
    if locked {
        // The goroutine may have locked this thread because
        // it put it in an unusual kernel state. Kill it
        // rather than returning it to the thread pool.

        // Return to mstart, which will release the P and exit
        // the thread.
        if GOOS != "plan9" { // See golang.org/issue/22227.
            gogo(&_g_.m.g0.sched)
        } else {
            // Clear lockedExt on plan9 since we may end up re-using
            // this thread.
            _g_.m.lockedExt = 0
        }
    }
    // 调度寻找下一个可执行的g
    schedule()
}
其中通过gfput将g放入空闲列表并且尝试释放g使用的栈。

// Put on gfree list.
// If local list is too long, transfer a batch to the global list.
func gfput(_p_ *p, gp *g) {
    if readgstatus(gp) != _Gdead {
        throw("gfput: bad status (not Gdead)")
    }

    stksize := gp.stack.hi - gp.stack.lo

    // 由注释可猜测 非标准的栈大小则释放 栈在分配时malg传入大小为StackMin(进一步计算处理过)
    //gfget时栈大小为_FixedStack 非标准的应该是malg分配的栈或者扩容后的栈
    if stksize != _FixedStack {
        // non-standard stack size - free it.
        stackfree(gp.stack)
        gp.stack.lo = 0
        gp.stack.hi = 0
        gp.stackguard0 = 0
    }

    // 将g放入p的gFree链表中 如果p的链表数据达到64 则释放32个到全局的sched空闲列表中
    _p_.gFree.push(gp)
    _p_.gFree.n++
    if _p_.gFree.n >= 64 {
        lock(&sched.gFree.lock)
        for _p_.gFree.n >= 32 {
            _p_.gFree.n--
            gp = _p_.gFree.pop()
            if gp.stack.lo == 0 {
                sched.gFree.noStack.push(gp)
            } else {
                sched.gFree.stack.push(gp)
            }
            sched.gFree.n++
        }
        unlock(&sched.gFree.lock)
    }
}
三、总结
1.TLS(thread local storage)的思想，从g的分配可看出，p充当了m的数据的一级缓存角色，因为p中的数据同一时刻只会被和其绑定的m访问，因此可以无锁使用。
2.goexit在初始化时写入栈中的做法值得细细品味。