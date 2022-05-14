---
title: mcall systemstack
layout: post
category: golang
author: 夏泽民
---
https://studygolang.com/articles/28553?fr=sidebar

不同硬件平台使用的汇编文件不同，本文分析的函数mcall, systemstack, asmcgocall是基于asm_arm64.s汇编文件。
不用操作系统平台使用的系统调用不同，本文分析的函数syscall是基于asm_linux_arm64.s汇编文件。

CPU的上下文
这些函数的本质都是为了切换goroutine，goroutine切换时需要切换CPU执行的上下文，主要有2个寄存器的值SP（当前线程使用的栈的栈顶地址），PC（下一个要执行的指令的地址）。

mcall函数
mcall函数的定义如下，mcall传入的是函数指针，传入函数的类型如下，只有一个参数goroutine的指针，无返回值。

func mcall(fn func(*g)
mcall函数的作用是在系统栈中执行调度代码，并且调度代码不会返回，将在运行过程中又一次执行mcall。mcall的流程是保存当前的g的上下文，切换到g0的上下文，传入函数参数，跳转到函数代码执行。

// void mcall(fn func(*g))
// Switch to m->g0's stack, call fn(g).
// Fn must never return. It should gogo(&g->sched)
// to keep running g.
TEXT runtime·mcall(SB), NOSPLIT|NOFRAME, $0-8
    // Save caller state in g->sched
    //此时线程当前的sp pc bp等上下文都存在寄存器中 需要将寄存器的值写回g 下面就是写回g的过程
    MOVD    RSP, R0  // R0 = RSP
    MOVD    R0, (g_sched+gobuf_sp)(g)  // g_sp = RO 保存sp寄存器的值
    MOVD    R29, (g_sched+gobuf_bp)(g) // g_bp = R29 (R29保存bp值)
    MOVD    LR, (g_sched+gobuf_pc)(g)  // g_pc = LR (LR保存pc值)
    MOVD    $0, (g_sched+gobuf_lr)(g)  // g_lr = 0
    MOVD    g, (g_sched+gobuf_g)(g)    // ???

    // Switch to m->g0 & its stack, call fn.
    // 将当前的g切为g0
    MOVD    g, R3  // R3 = g (g表示当前调用mcall时的goutine)
    MOVD    g_m(g), R8 // R8 = g.m (R8表示g绑定的m 即当前的m)
    MOVD    m_g0(R8), g // g = m.g0 (将当前g切换为g0)
    BL  runtime·save_g(SB) // ???
    CMP g, R3 // g == g0  R3 == 调用mcall的g 必不相等
    BNE 2(PC) // 如果不想等则正常执行
    B   runtime·badmcall(SB) // 相等则说明有bug 调用badmcall
    // fn是要调用的函数 写入寄存器
    MOVD    fn+0(FP), R26           // context R26存的是fn的pc
    MOVD    0(R26), R4          // code pointer R4也是fn的pc值
    MOVD    (g_sched+gobuf_sp)(g), R0  // g0的 sp值赋给寄存器
    MOVD    R0, RSP // sp = m->g0->sched.sp
    MOVD    (g_sched+gobuf_bp)(g), R29 // g0的bp值赋给对应的寄存器
    MOVD    R3, -8(RSP) // R3在之前被赋值为调用mcall的g 现在写入g0的栈中 作为fn的函数参数
    MOVD    $0, -16(RSP) // 此处的空值不太理解 只有一个参数且无返回值 为何要在栈中预留8字节
    SUB $16, RSP // 对栈进行偏移16byte（上面g $0 各占8byte）
    BL  (R4) // R4此时是fn的pc值 跳到该 PC执行fn
    B   runtime·badmcall2(SB) // 该函数永远不会返回 因此这一步理论上永远执行不到
常见的调用mcall执行的函数有：

mcall(gosched_m)
mcall(park_m)
mcall(goexit0)
mcall(exitsyscall0)
mcall(preemptPark)
mcall(gopreempt_m)
<!-- more -->
systemstack函数
systemstack函数的定义如下，传入的函数无参数，无返回值。

func systemstack(fn func())
systemstack函数的作用是在系统栈中执行只能由g0(或gsignal?)执行的调度代码，和mcall不同的是，在执行完调度代码后会切回到现在正在执行的代码。
该部分的源码注释有只有个大概的流程的理解，许多细节推敲不出来。主要流程是先判断当前运行的g是否为g0或者gsignal，如果是则直接运行，不是则先切换到g0，执行完函数后切换为g返回调用处。

// systemstack_switch is a dummy routine that systemstack leaves at the bottom
// of the G stack. We need to distinguish the routine that
// lives at the bottom of the G stack from the one that lives
// at the top of the system stack because the one at the top of
// the system stack terminates the stack walk (see topofstack()).
TEXT runtime·systemstack_switch(SB), NOSPLIT, $0-0
    UNDEF
    BL  (LR)    // make sure this function is not leaf
    RET

// func systemstack(fn func())
TEXT runtime·systemstack(SB), NOSPLIT, $0-8
    MOVD    fn+0(FP), R3    // R3 = fn
    MOVD    R3, R26     // context R26 = R3 = fn
    MOVD    g_m(g), R4  // R4 = m

    MOVD    m_gsignal(R4), R5   // R5 = m.gsignal
    CMP g, R5  // m.gsignal是有权限执行fn的g
    BEQ noswitch // 如果相等说明已经是m.gsignale了 则不需要切换

    MOVD    m_g0(R4), R5    // R5 = g0
    CMP g, R5  // 如果当前的g已经是g0 则说明不用切换
    BEQ noswitch

    MOVD    m_curg(R4), R6 // R6 = m.curg
    CMP g, R6 // m.curg == g
    BEQ switch

    // Bad: g is not gsignal, not g0, not curg. What is it?
    // Hide call from linker nosplit analysis.
    MOVD    $runtime·badsystemstack(SB), R3
    BL  (R3)
    B   runtime·abort(SB)

switch:
    // save our state in g->sched. Pretend to
    // be systemstack_switch if the G stack is scanned.
    MOVD    $runtime·systemstack_switch(SB), R6
    ADD $8, R6  // get past prologue
    // 以下是常规的保存当前g的上下文
    MOVD    R6, (g_sched+gobuf_pc)(g)
    MOVD    RSP, R0
    MOVD    R0, (g_sched+gobuf_sp)(g)
    MOVD    R29, (g_sched+gobuf_bp)(g)
    MOVD    $0, (g_sched+gobuf_lr)(g)
    MOVD    g, (g_sched+gobuf_g)(g)

    // switch to g0
    MOVD    R5, g  // g = R5 = g0
    BL  runtime·save_g(SB)
    MOVD    (g_sched+gobuf_sp)(g), R3 // R3 = sp
    // make it look like mstart called systemstack on g0, to stop traceback
    SUB $16, R3  // sp地址 内存对齐
    AND $~15, R3
    MOVD    $runtime·mstart(SB), R4
    MOVD    R4, 0(R3)
    MOVD    R3, RSP
    MOVD    (g_sched+gobuf_bp)(g), R29 // R29 = g0.gobuf.bp

    // call target function
    MOVD    0(R26), R3  // code pointer
    BL  (R3)

    // switch back to g
    MOVD    g_m(g), R3
    MOVD    m_curg(R3), g
    BL  runtime·save_g(SB)
    MOVD    (g_sched+gobuf_sp)(g), R0
    MOVD    R0, RSP
    MOVD    (g_sched+gobuf_bp)(g), R29
    MOVD    $0, (g_sched+gobuf_sp)(g)
    MOVD    $0, (g_sched+gobuf_bp)(g)
    RET

noswitch:
    // already on m stack, just call directly
    // Using a tail call here cleans up tracebacks since we won't stop
    // at an intermediate systemstack.
    MOVD    0(R26), R3  // code pointer  R3 = R26 = fn
    MOVD.P  16(RSP), R30    // restore LR  R30 = RSP + 16(systemstack调用完成后下条指令的PC值？)
    SUB $8, RSP, R29    // restore FP  R29 = RSP - 8 表示栈的
    B   (R3)
asmcgocall函数
asmcgocall函数定义如下，传入的参数有2个为函数指针和参数指针，返回参数为int32。

func asmcgocall(fn, arg unsafe.Pointer) int32
asmcgocall函数的作用是执行cgo代码，该部分代码只能在g0(或gsignal, osthread)的栈执行，因此流程是先判断当前的栈是否要切换，如果无需切换则直接执行nosave然后返回，否则先保存当前g的上下文，然后切换到g0，执行完cgo代码后切回g，然后返回。

// func asmcgocall(fn, arg unsafe.Pointer) int32
// Call fn(arg) on the scheduler stack,
// aligned appropriately for the gcc ABI.
// See cgocall.go for more details.
TEXT ·asmcgocall(SB),NOSPLIT,$0-20
    MOVD    fn+0(FP), R1  // R1 = fn
    MOVD    arg+8(FP), R0  // R2 = arg

    MOVD    RSP, R2     // save original stack pointer
    CBZ g, nosave  // 如果g为nil 则跳转到 nosave。 g == nil是否说明当前是osthread？
    MOVD    g, R4  // R4 = g

    // Figure out if we need to switch to m->g0 stack.
    // We get called to create new OS threads too, and those
    // come in on the m->g0 stack already.
    MOVD    g_m(g), R8 // R8 = g.m
    MOVD    m_gsignal(R8), R3 // R3 = g.m.gsignal
    CMP R3, g  // 如果g == g.m.signal jump nosave
    BEQ nosave
    MOVD    m_g0(R8), R3 // 如果g== m.g0 jump nosave
    CMP R3, g
    BEQ nosave

    // Switch to system stack.
    // save g的上下文
    MOVD    R0, R9  // gosave<> and save_g might clobber R0
    BL  gosave<>(SB)
    MOVD    R3, g
    BL  runtime·save_g(SB)
    MOVD    (g_sched+gobuf_sp)(g), R0
    MOVD    R0, RSP
    MOVD    (g_sched+gobuf_bp)(g), R29
    MOVD    R9, R0

    // Now on a scheduling stack (a pthread-created stack).
    // Save room for two of our pointers /*, plus 32 bytes of callee
    // save area that lives on the caller stack. */
    MOVD    RSP, R13
    SUB $16, R13
    MOVD    R13, RSP  // RSP = RSP - 16
    MOVD    R4, 0(RSP)  // save old g on stack  RSP.0 = R4 = oldg
    MOVD    (g_stack+stack_hi)(R4), R4 // R4 = old.g.stack.hi
    SUB R2, R4  // R4 = oldg.stack.hi - old_RSP
    MOVD    R4, 8(RSP)  // save depth in old g stack (can't just save SP, as stack might be copied during a callback)
    BL  (R1) // R1 = fn
    MOVD    R0, R9 // R9 = R0 = errno?

    // Restore g, stack pointer. R0 is errno, so don't touch it
    MOVD    0(RSP), g  // g = RSP.0 = oldg
    BL  runtime·save_g(SB)
    MOVD    (g_stack+stack_hi)(g), R5 // R5 = g.stack.hi
    MOVD    8(RSP), R6 // R6 = RSP + 8 = oldg.stack.hi - old_RSP
    SUB R6, R5 // R5 = R5 - R6 = old_RSP
    MOVD    R9, R0 // R0 = R9 = errno
    MOVD    R5, RSP // RSP = R5 = old_RSP

    MOVW    R0, ret+16(FP) // ret = R0 = errno
    RET

nosave:
    // Running on a system stack, perhaps even without a g.
    // Having no g can happen during thread creation or thread teardown
    // (see needm/dropm on Solaris, for example).
    // This code is like the above sequence but without saving/restoring g
    // and without worrying about the stack moving out from under us
    // (because we're on a system stack, not a goroutine stack).
    // The above code could be used directly if already on a system stack,
    // but then the only path through this code would be a rare case on Solaris.
    // Using this code for all "already on system stack" calls exercises it more,
    // which should help keep it correct.
    MOVD    RSP, R13 
    SUB $16, R13  
    MOVD    R13, RSP // RSP = RSP - 16
    MOVD    $0, R4 // R4 = 0
    MOVD    R4, 0(RSP)  // Where above code stores g, in case someone looks during debugging.
    MOVD    R2, 8(RSP)  // Save original stack pointer.  RSP + 8 = old_R2
    BL  (R1)
    // Restore stack pointer.
    MOVD    8(RSP), R2  // R2 = RSP + 8 = old_R2
    MOVD    R2, RSP // RSP = old_R2 = old_RSP
    MOVD    R0, ret+16(FP) // ret = R0 = errno
    RET
syscall函数
Syscall函数的定义如下，传入4个参数，返回3个参数。

func syscall(fn, a1, a2, a3 uintptr) (r1, r2 uintptr, err Errno)
syscall函数的作用是传入系统调用的地址和参数，执行完成后返回。流程主要是系统调用前执行entersyscall，设置g p的状态，然后入参，执行后，写返回值然后执行exitsyscall设置g p的状态。
entersyscall和exitsyscall在g的调用中细讲。

// func Syscall(trap int64, a1, a2, a3 uintptr) (r1, r2, err uintptr);
// Trap # in AX, args in DI SI DX R10 R8 R9, return in AX DX
// Note that this differs from "standard" ABI convention, which
// would pass 4th arg in CX, not R10.

// 4个入参：PC param1 param2 param3
TEXT ·Syscall(SB),NOSPLIT,$0-56
    // 调用entersyscall 判断是执行条件是否满足 记录调度信息 切换g p的状态
    CALL    runtime·entersyscall(SB)
    // 将参数存入寄存器中
    MOVQ    a1+8(FP), DI
    MOVQ    a2+16(FP), SI
    MOVQ    a3+24(FP), DX
    MOVQ    trap+0(FP), AX  // syscall entry
    SYSCALL
    CMPQ    AX, $0xfffffffffffff001
    JLS ok
    // 执行失败时 写返回值
    MOVQ    $-1, r1+32(FP)
    MOVQ    $0, r2+40(FP)
    NEGQ    AX
    MOVQ    AX, err+48(FP)
    // 调用exitsyscall 记录调度信息
    CALL    runtime·exitsyscall(SB)
    RET
ok:
    // 执行成功时 写返回值
    MOVQ    AX, r1+32(FP)
    MOVQ    DX, r2+40(FP)
    MOVQ    $0, err+48(FP)
    CALL    runtime·exitsyscall(SB)
    RET
除了Syscal还有Syscall6（除fn还有6个参数）对应有6个参数的系统调用。实现大同小异，这里不分析。

总结与思考
1.汇编函数的作用。为什么golang一定要引入汇编函数呢？因为CPU执行时的上下文是寄存器，只有汇编语言才能操作寄存器。
2.CPU的上下文和g.sched(gobuf)结构体中的字段一一对应，只有10个以内的字段，因此切换上下文效率非常的高。
3.除了golang，其它在用的语言是否要有类似的汇编来实现语言和操作系统之间的交互？

最后
除了mcall函数，其它函数在具体执行细节上理解不够深，后面加强汇编相关的知识后再把这个坑填上。

https://studygolang.com/articles/28553?fr=sidebar


golang里面用户态进行系统调用时候的一些原理，主要关注点将会放在system call与scheduler之间的关联。

https://blog.csdn.net/u010853261/article/details/88312904

1.入口
系统调用的入口根据不同系统有不同实现，对于AMD64, Linux环境是：syscall/asm_linux_amd64.s

函数声明如下：

func Syscall(trap, a1, a2, a3 uintptr) (r1, r2 uintptr, err syscall.Errno)

func Syscall6(trap, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2 uintptr, err syscall.Errno)

func RawSyscall(trap, a1, a2, a3 uintptr) (r1, r2 uintptr, err syscall.Errno)

func RawSyscall6(trap, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2 uintptr, err syscall.Errno)

这些函数的实现都是汇编，按照 linux 的 syscall 调用规范，我们只要在汇编中把参数依次传入寄存器，并调用 SYSCALL 指令即可进入内核处理逻辑，系统调用执行完毕之后，返回值放在 RAX 中:

Syscall 和 Syscall6 的区别只有传入参数不一样, 具体源码与实现请看golang的开源源码。

这里只列出Syscall和RawSyscall的源码：

//Syscall
TEXT ·Syscall(SB),NOSPLIT,$0-56
	CALL	runtime·entersyscall(SB)
	MOVQ	a1+8(FP), DI
	MOVQ	a2+16(FP), SI
	MOVQ	a3+24(FP), DX
	MOVQ	$0, R10
	MOVQ	$0, R8
	MOVQ	$0, R9
	MOVQ	trap+0(FP), AX	// syscall entry
	SYSCALL
	CMPQ	AX, $0xfffffffffffff001
	JLS	ok
	MOVQ	$-1, r1+32(FP)
	MOVQ	$0, r2+40(FP)
	NEGQ	AX
	MOVQ	AX, err+48(FP)
	CALL	runtime·exitsyscall(SB)
	RET
ok:
	MOVQ	AX, r1+32(FP)
	MOVQ	DX, r2+40(FP)
	MOVQ	$0, err+48(FP)
	CALL	runtime·exitsyscall(SB)
	RET

//RawSyscall
TEXT ·RawSyscall(SB),NOSPLIT,$0-56
	MOVQ	a1+8(FP), DI
	MOVQ	a2+16(FP), SI
	MOVQ	a3+24(FP), DX
	MOVQ	$0, R10
	MOVQ	$0, R8
	MOVQ	$0, R9
	MOVQ	trap+0(FP), AX	// syscall entry
	SYSCALL
	CMPQ	AX, $0xfffffffffffff001
	JLS	ok1
	MOVQ	$-1, r1+32(FP)
	MOVQ	$0, r2+40(FP)
	NEGQ	AX
	MOVQ	AX, err+48(FP)
	RET
ok1:
	MOVQ	AX, r1+32(FP)
	MOVQ	DX, r2+40(FP)
	MOVQ	$0, err+48(FP)
	RET

Syscall和RawSyscall的实现比较典型，可以看到这两个实现最主要的区别在于：
Syscall在进入系统调用的时候，调用了runtime·entersyscall(SB)函数，在结束系统调用的时候调用了runtime·exitsyscall(SB)。做到进入和退出syscall的时候通知runtime。

这两个函数runtime·entersyscall和runtime·exitsyscall的实现在proc.go文件里面。其实在runtime·entersyscall函数里面，通知系统调用时候，是会将g的M的P解绑，P可以去继续获取M执行其余的g，这样提升效率。

所以如果用户代码使用了 RawSyscall 来做一些阻塞的系统调用，是有可能阻塞其它的 g 的。RawSyscall 只是为了在执行那些一定不会阻塞的系统调用时，能节省两次对 runtime 的函数调用消耗。

runtime·entersyscall和runtime·exitsyscall这两个函数也是与scheduler交互的地方，后面会对源码进行分析。

2.系统调用管理
系统调用的定义文件: /syscall/syscall_linux.go

可以把系统调用分为三类:

阻塞系统调用
非阻塞系统调用非阻塞系统调用
wrapped 系统调用
阻塞系统调用会定义成下面这样的形式:

//sys   Madvise(b []byte, advice int) (err error)
1
非阻塞系统调用:

//sysnb    EpollCreate(size int) (fd int, err error)
1
然后，根据这些注释，mksyscall.pl 脚本会生成对应的平台的具体实现。mksyscall.pl 是一段 perl 脚本，感兴趣的同学可以自行查看，这里就不再赘述了。

看看阻塞和非阻塞的系统调用的生成结果:

func Madvise(b []byte, advice int) (err error) {
    var _p0 unsafe.Pointer
    if len(b) > 0 {
        _p0 = unsafe.Pointer(&b[0])
    } else {
        _p0 = unsafe.Pointer(&_zero)
    }
    _, _, e1 := Syscall(SYS_MADVISE, uintptr(_p0), uintptr(len(b)), uintptr(advice))
    if e1 != 0 {
        err = errnoErr(e1)
    }
    return
}

func EpollCreate(size int) (fd int, err error) {
    r0, _, e1 := RawSyscall(SYS_EPOLL_CREATE, uintptr(size), 0, 0)
    fd = int(r0)
    if e1 != 0 {
        err = errnoErr(e1)
    }
    return
}

标记为 sys(阻塞)的系统调用使用的是 Syscall 或者 Syscall6，标记为 sysnb(非阻塞) 的系统调用使用的是 RawSyscall 或 RawSyscall6。

wrapped 的系统调用是怎么一回事呢？

func Rename(oldpath string, newpath string) (err error) {
    return Renameat(_AT_FDCWD, oldpath, _AT_FDCWD, newpath)
}

可能是觉得系统调用的名字不太好，或者参数太多，我们就简单包装一下。没啥特别的。

3.runtime 中的 SYSCALL
除了上面提到的阻塞非阻塞和 wrapped syscall，runtime 中还定义了一些 low-level 的 syscall，这些是不暴露给用户程序的。

提供给用户的 syscall 库，在使用时，会使 goroutine 和 p 分别进入 Gsyscall 和 Psyscall 状态。但 runtime 自己封装的这些 syscall 无论是否阻塞，都不会调用 entersyscall 和 exitsyscall。 虽说是 “low-level” 的 syscall， 不过和暴露给用户的 syscall 本质是一样的。这些代码在 runtime/sys_linux_amd64.s 中，举个具体的例子:

TEXT runtime·write(SB),NOSPLIT,$0-28
    MOVQ    fd+0(FP), DI
    MOVQ    p+8(FP), SI
    MOVL    n+16(FP), DX
    MOVL    $SYS_write, AX
    SYSCALL
    CMPQ    AX, $0xfffffffffffff001
    JLS    2(PC)
    MOVL    $-1, AX
    MOVL    AX, ret+24(FP)
    RET

TEXT runtime·read(SB),NOSPLIT,$0-28
    MOVL    fd+0(FP), DI
    MOVQ    p+8(FP), SI
    MOVL    n+16(FP), DX
    MOVL    $SYS_read, AX
    SYSCALL
    CMPQ    AX, $0xfffffffffffff001
    JLS    2(PC)
    MOVL    $-1, AX
    MOVL    AX, ret+24(FP)
    RET

这些 syscall 理论上都是不会在执行期间被调度器剥离掉 p 的，所以执行成功之后 goroutine 会继续执行，而不像用户的 goroutine 一样，若被剥离 p 会进入等待队列。

4.用户代码的系统调用和调度交互
既然要和调度交互，那就要友好地通知我要 syscall 了: entersyscall，我完事了: exitsyscall。

所以这里的交互指的是用户代码使用 syscall 库时和调度器的交互。runtime 里的 syscall 不走这套流程。

entersyscall和exitsyscall的pipeline
                      +-----------------------------------------------------+
                      |runtime.entersyscall()                               |
                      |1) save() goroutine Save on site                     |
user code             |2) casgstatus(_g_, _Grunning, _Gsyscall)             |
 syscall   ---------->|3) atomic.Store(&_g_.m.p.ptr().status, _Psyscall)    |
                      |                                                     |
                      |a) the M is blocking;                                |
                      |b) the status of P is _Psyscall, So the P can be     |
                      |schedule to execute other goroutine                  |
                      +-----------------------------------------------------+
                                                                             
                                                                             
                                                                             
                                                                             
                             +--------------------------+                    
     user code               |runtime.exitsyscall()     |                    
 syscall finished ---------->|1) disable preemption     |                    
                             |                          |                    
                             +--+--------------------+--+                    
                                |                    |                                       
                                |                    |                           
                                |                    |                       
                                v                    v                       
                        try to re-acquire   try to get any other             
                            the last P             idle P                    
                                |                    |                       
                                |                    |                       
                                |                    |                       
                      success---+--------------------+-----+                 
                          |                                |                 
                          |                              fail                
                          |                                |                 
                          v                                |                 
               +---------------------+          +----------+------>------+   
               |there is a P to run G|          |not get P               |   
               |runtime.exexute(G)   |          |1.put G into global tail|   
               |schedule loop        |          |2.idel this M           |   
               |                     |          |                        |   
               +---------------------+          +------------------------+   

entersyscall
直接看源码：

// 用户代码使用 syscall 库时和调度器的交互；
// runtime本身的syscall不走这一套流程
// Standard syscall entry used by the go syscall library and normal cgo calls.
//go:nosplit
func entersyscall() {
	reentersyscall(getcallerpc(), getcallersp())
}

//go:nosplit
func reentersyscall(pc, sp uintptr) {
	_g_ := getg()

	// 需要禁止 g 的抢占
	_g_.m.locks++

	// entersyscall 中不能调用任何会导致栈增长/分裂的函数
	// (See details in comment above.)
	// Catch calls that might, by replacing the stack guard with something that
	// will trip any stack check and leaving a flag to tell newstack to die.
	_g_.stackguard0 = stackPreempt
	_g_.throwsplit = true

	// Leave SP around for GC and traceback.
	//保存现场，在 syscall 之后会依据这些数据恢复现场
	save(pc, sp)
	_g_.syscallsp = sp
	_g_.syscallpc = pc
	casgstatus(_g_, _Grunning, _Gsyscall)
	if _g_.syscallsp < _g_.stack.lo || _g_.stack.hi < _g_.syscallsp {
		systemstack(func() {
			print("entersyscall inconsistent ", hex(_g_.syscallsp), " [", hex(_g_.stack.lo), ",", hex(_g_.stack.hi), "]\n")
			throw("entersyscall")
		})
	}

	if trace.enabled {
		systemstack(traceGoSysCall)
		// systemstack itself clobbers g.sched.{pc,sp} and we might
		// need them later when the G is genuinely blocked in a
		// syscall
		save(pc, sp)
	}

	if atomic.Load(&sched.sysmonwait) != 0 {
		systemstack(entersyscall_sysmon)
		save(pc, sp)
	}

	if _g_.m.p.ptr().runSafePointFn != 0 {
		// runSafePointFn may stack split if run on this stack
		systemstack(runSafePointFn)
		save(pc, sp)
	}

	_g_.m.syscalltick = _g_.m.p.ptr().syscalltick
	_g_.sysblocktraced = true
	_g_.m.mcache = nil
	// 解绑P与M的关系
	_g_.m.p.ptr().m = 0
	atomic.Store(&_g_.m.p.ptr().status, _Psyscall)
	if sched.gcwaiting != 0 {
		systemstack(entersyscall_gcwait)
		save(pc, sp)
	}

	_g_.m.locks--
}

主要流程如下：

设置_g_.m.locks++，禁止g被强占
设置_g_.stackguard0 = stackPreempt，禁止调用任何会导致栈增长/分裂的函数
保存现场，在 syscall 之后会依据这些数据恢复现场
更新G的状态为_Gsyscall
释放局部调度器P：解绑P与M的关系；
更新P状态为_Psyscall
g.m.locks–解除禁止强占。
可以看到，进入 syscall 的 G 是铁定不会被抢占的。

此外进入系统调用的goroutine会阻塞，导致内核M会阻塞。此时P会被剥离掉，所以P可以继续去获取其余的空闲M执行其余的goroutine。

exitsyscall
直接看源码

// g 已经退出了 syscall
// 需要准备让 g 在 cpu 上重新运行
// 不能有 write barrier，因为 P 可能已经被偷走了
//go:nosplit
//go:nowritebarrierrec
func exitsyscall() {
	_g_ := getg()
	// 禁止强占
	_g_.m.locks++ // see comment in entersyscall
	if getcallersp() > _g_.syscallsp {
		throw("exitsyscall: syscall frame is no longer valid")
	}

	_g_.waitsince = 0
	oldp := _g_.m.p.ptr()
	if exitsyscallfast() {
		if _g_.m.mcache == nil {
			throw("lost mcache")
		}
		if trace.enabled {
			if oldp != _g_.m.p.ptr() || _g_.m.syscalltick != _g_.m.p.ptr().syscalltick {
				systemstack(traceGoStart)
			}
		}
		// There's a cpu for us, so we can run.
		_g_.m.p.ptr().syscalltick++
		// We need to cas the status and scan before resuming...
		casgstatus(_g_, _Gsyscall, _Grunning)

		// Garbage collector isn't running (since we are),
		// so okay to clear syscallsp.
		_g_.syscallsp = 0
		_g_.m.locks--
		if _g_.preempt {
			// restore the preemption request in case we've cleared it in newstack
			_g_.stackguard0 = stackPreempt
		} else {
			// otherwise restore the real _StackGuard, we've spoiled it in entersyscall/entersyscallblock
			_g_.stackguard0 = _g_.stack.lo + _StackGuard
		}
		_g_.throwsplit = false
		return
	}

	_g_.sysexitticks = 0
	if trace.enabled {
		// Wait till traceGoSysBlock event is emitted.
		// This ensures consistency of the trace (the goroutine is started after it is blocked).
		for oldp != nil && oldp.syscalltick == _g_.m.syscalltick {
			osyield()
		}
		// We can't trace syscall exit right now because we don't have a P.
		// Tracing code can invoke write barriers that cannot run without a P.
		// So instead we remember the syscall exit time and emit the event
		// in execute when we have a P.
		_g_.sysexitticks = cputicks()
	}

	_g_.m.locks--

	// Call the scheduler.
	mcall(exitsyscall0)

	if _g_.m.mcache == nil {
		throw("lost mcache")
	}

	// Scheduler returned, so we're allowed to run now.
	// Delete the syscallsp information that we left for
	// the garbage collector during the system call.
	// Must wait until now because until gosched returns
	// we don't know for sure that the garbage collector
	// is not running.
	_g_.syscallsp = 0
	_g_.m.p.ptr().syscalltick++
	_g_.throwsplit = false
}

//exitsyscallfast
//go:nosplit
func exitsyscallfast() bool {
	_g_ := getg()

	// Freezetheworld sets stopwait but does not retake P's.
	if sched.stopwait == freezeStopWait {
		_g_.m.mcache = nil
		_g_.m.p = 0
		return false
	}

	// Try to re-acquire the last P.
	if _g_.m.p != 0 && _g_.m.p.ptr().status == _Psyscall && atomic.Cas(&_g_.m.p.ptr().status, _Psyscall, _Prunning) {
		// There's a cpu for us, so we can run.
		exitsyscallfast_reacquired()
		return true
	}

	// Try to get any other idle P.
	oldp := _g_.m.p.ptr()
	_g_.m.mcache = nil
	_g_.m.p = 0
	if sched.pidle != 0 {
		var ok bool
		systemstack(func() {
			ok = exitsyscallfast_pidle()
			if ok && trace.enabled {
				if oldp != nil {
					// Wait till traceGoSysBlock event is emitted.
					// This ensures consistency of the trace (the goroutine is started after it is blocked).
					for oldp.syscalltick == _g_.m.syscalltick {
						osyield()
					}
				}
				traceGoSysExit(0)
			}
		})
		if ok {
			return true
		}
	}
	return false
}

// exitsyscall0
// exitsyscall slow path on g0.
// Failed to acquire P, enqueue gp as runnable.
//
//go:nowritebarrierrec
func exitsyscall0(gp *g) {
	_g_ := getg()

	casgstatus(gp, _Gsyscall, _Grunnable)
	dropg()
	lock(&sched.lock)
	_p_ := pidleget()
	if _p_ == nil {
		globrunqput(gp)
	} else if atomic.Load(&sched.sysmonwait) != 0 {
		atomic.Store(&sched.sysmonwait, 0)
		notewakeup(&sched.sysmonnote)
	}
	unlock(&sched.lock)
	if _p_ != nil {
		acquirep(_p_)
		execute(gp, false) // Never returns.
	}
	if _g_.m.lockedg != 0 {
		// Wait until another thread schedules gp and so m again.
		stoplockedm()
		execute(gp, false) // Never returns.
	}
	stopm()
	schedule() // Never returns.
}

主要的pipeline如下：

设置 g.m.locks++ 禁止强占
调用 exitsyscallfast() 快速退出系统调用
2.1. Try to re-acquire the last P，如果成功就直接接return;
2.2. Try to get any other idle P from allIdleP list;
2.3. 没有获取到空闲的P
如果快速获取到了P：
3.1. 更新G 的状态是_Grunning
3.2. 与G绑定的M会在退出系统调用之后继续执行
没有获取到空闲的P：
4.1. 调用mcall()函数切换到g0的栈空间；
4.2. 调用exitsyscall0函数：
4.2.1. 更新G 的状态是_Grunning
4.2.2. 调用dropg()：解除当前g与M的绑定关系；
4.2.3. 调用globrunqput将G插入global queue的队尾，
4.2.4. 调用stopm()释放M，将M加入全局的idel M列表，这个调用会阻塞，知道获取到可用的P。
4.2.5. 如果4.2.4中阻塞结束，M获取到了可用的P，会调用schedule()函数，执行一次新的调度。
需要注意的是，调用 exitsyscall0 时，会切换到 g0 栈。

entersyscallblock
用户代码进行系统调用时候，知道自己会 block，直接就把 p 交出来了。

代码实现和 entersyscall 一样，就是会直接把 P 给交出去，因为知道自己是会阻塞的。

// 和 entersyscall 一样，就是会直接把 P 给交出去，因为知道自己是会阻塞的
//go:nosplit
func entersyscallblock(dummy int32) {
    _g_ := getg()

    _g_.m.locks++ // see comment in entersyscall
    _g_.throwsplit = true
    _g_.stackguard0 = stackPreempt // see comment in entersyscall
    _g_.m.syscalltick = _g_.m.p.ptr().syscalltick
    _g_.sysblocktraced = true
    _g_.m.p.ptr().syscalltick++

    // Leave SP around for GC and traceback.
    pc := getcallerpc()
    sp := getcallersp(unsafe.Pointer(&dummy))
    save(pc, sp)
    _g_.syscallsp = _g_.sched.sp
    _g_.syscallpc = _g_.sched.pc
    if _g_.syscallsp < _g_.stack.lo || _g_.stack.hi < _g_.syscallsp {
        sp1 := sp
        sp2 := _g_.sched.sp
        sp3 := _g_.syscallsp
        systemstack(func() {
            print("entersyscallblock inconsistent ", hex(sp1), " ", hex(sp2), " ", hex(sp3), " [", hex(_g_.stack.lo), ",", hex(_g_.stack.hi), "]\n")
            throw("entersyscallblock")
        })
    }
    casgstatus(_g_, _Grunning, _Gsyscall)
    if _g_.syscallsp < _g_.stack.lo || _g_.stack.hi < _g_.syscallsp {
        systemstack(func() {
            print("entersyscallblock inconsistent ", hex(sp), " ", hex(_g_.sched.sp), " ", hex(_g_.syscallsp), " [", hex(_g_.stack.lo), ",", hex(_g_.stack.hi), "]\n")
            throw("entersyscallblock")
        })
    }

    // 直接调用 entersyscallblock_handoff 把 p 交出来了
    systemstack(entersyscallblock_handoff)

    // Resave for traceback during blocked call.
    save(getcallerpc(), getcallersp(unsafe.Pointer(&dummy)))

    _g_.m.locks--
}

这个函数只有一个调用方 notesleepg，这里就不再赘述了。

func entersyscallblock_handoff() {
    handoffp(releasep())
}

5. 总结
提供给用户使用的系统调用，基本都会通知 runtime，以 entersyscall，exitsyscall 的形式来告诉 runtime，在这个 syscall 阻塞的时候，由 runtime 判断是否把 P 腾出来给其它的 M 用。解绑定指的是把 M 和 P 之间解绑，如果绑定被解除，在 syscall 返回时，这个 g 会被放入全局执行队列 global runq 中。

同时 runtime 又保留了自己的特权，在执行自己的逻辑的时候，我的 P 不会被调走，这样保证了在 Go 自己“底层”使用的这些 syscall 返回之后都能被立刻处理。

所以同样是 epollwait，runtime 用的是不能被别人打断的，你用的 syscall.EpollWait 那显然是没有这种特权的。
https://blog.csdn.net/u010853261/article/details/88312904