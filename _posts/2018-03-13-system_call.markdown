---
title: 系统调用
layout: post
category: linux
author: 夏泽民
---
<!-- more -->
为什么需要系统调用
   现代的操作系统通常都具有多任务处理的功能，通常靠进程来实现。由于操作系统快速的在每个进程间切换执行，所以一切看起来就会像是同时的。同时这也带来了很多安全问题，例如，一个进程可以轻易的修改进程的内存空间中的数据来使另一个进程异常或达到一些目的，因此操作系统必须保证每一个进程都能安全的执行。这一问题的解决方法是在处理器中加入基址寄存器和界限寄存器。这两个寄存器中的内容用硬件限制了对储存器的存取指令所访问的储存器的地址。这样就可以在系统切换进程时写入这两个寄存器的内容到该进程被分配的地址范围，从而避免恶意软件。
   为了防止用户程序修改基址寄存器和界限寄存器中的内容来达到访问其他内存空间的目的，这两个寄存器必须通过一些特殊的指令来访问。通常，处理器设有两种模式：“用户模式”与“内核模式”，通过一个标签位来鉴别当前正处于什么模式。一些诸如修改基址寄存器内容的指令只有在内核模式中可以执行，而处于用户模式的时候硬件会直接跳过这个指令并继续执行下一个。
   当操作系统接收到系统调用请求后，会让处理器进入内核模式，从而执行诸如I/O操作，修改基址寄存器内容等指令，而当处理完系统调用内容后，操作系统会让处理器返回用户模式，来执行用户代码。

系统中的程序类型及状态
   操作系统中的状态分为管态（核心态）和目态（用户态）。特权指令：一类只能在核心态下运行而不能在用户态下运行的特殊指令。不同的操作系统特权指令会有所差异，但是一般来说主要是和硬件相关的一些指令。访管指令：本身是一条特殊的指令，但不是特权指令。（trap指令）。基本功能：“自愿进管”，能引起访管异常。用户程序只在用户态下运行，有时需要访问系统核心功能，这时通过系统调用接口使用系统调用。
   系统功能调用：就是用户在程序中使用“访管指令”调用由操作系统提供的子功能集合。其中每一个系统子功能称为一个系统调用命令，也叫广义指令。
   
   系统调用本质上是一种过程调用，但它是一种特殊的过程调用，与一般用户程序中的过程调用有明显的区别 。
   系统调用的调用过程和被调用过程运行在不同的状态，而普通的过程调用一般运行在相同的状态。
   系统调用必须通过软中断机制首先进入系统核心，然后才能转向相应的命令处理程序。普通过程调用可以直接由调用过程转向被调用过程
   用户态和内核态的区分：

       现代计算机机中都有几种不同的指令级别，在高执行级别下，代码可以执行特权指令，访问任意的物理地址，这种CPU执行级别就对应着内核态，而在相应的低级别执行状态下，代码的掌控范围会受到限制，只能在对应级别允许的范围内活动。举例：Intrel x86 CPU有四种不同的执行级别0-3，Linux只使用了其中的0级和3级来分别表示内核态和用户态。操作系统让系统本身更为稳定的方式，这样程序员自己写的用户态代码很难把整个系统都给搞崩溃，内核的代码经过仔细的分析有专业的人员写的代码会更加健壮一些，整个程序会更加稳定一些，注意：这里所说的地址空间是逻辑地址而不是物理地址。

     用户态和内核态的很显著的区分就是：CS和EIP， CS寄存器的最低两位表明了当前代码的特权级别；CPU每条指令的读取都是通过CS:EIP这两个寄存器：其中CS是代码段选择寄存器，EIP是偏移量寄存器，上述判断由硬件完成。一般来说在Linux中，地址空间是一个显著的标志：0xc0000000以上的地址空间只能在内核态下访问，0xc00000000-0xbfffffff的地址空间在两种状态下都可以访问。

中断处理是从用户态进入到内核态的主要的方式：

      也可能是用户态程序执行的过程中调用了一个系统调用陷入了内核态当中，这个叫做trap,系统调用只是一种特殊的中断。
      寄存器上下文：
            ——从用户态切换到内核态的时候
                  必须保存用户态的寄存器上下文
                  要保存哪些？
                  保存在哪里？
      中断/int指令会在堆栈上保存一些寄存器的值
            ——如：用户态栈顶地址、当时的状态字、当时的cs:eip的值
      中断发生的后的第一件事就是保护现场，保护现场就是进入中断的程序保存需要用到的寄存器数据，恢复现场就是退出中断程序，恢复、保存寄存器的数据。
         Libc库定义的一些API引用了封装例成，唯一目的就是发布系统调用：1.一般每个系统调用对应一个封装例程；2.库函数再用这些封装例程定义出给用户的API（把系统调用封装成很多歌方便程序员使用的函数，不是每个API都对应一个特定的系统调用）
     API可能直接提供用户态的服务 如：一些数学函数 1.一个单独的API可能调用几个系统调用2.不同的API可能调用了同一个系统调用返回：大部分封装例程返回一个整数，其值的含义依赖于相应的系统调用-1在多数情况下表示内核不能满足进程的请求，Libc中定义的errno变量包含特定的出错码
     我们要扒开系统调用的三层皮，我们讲这三层皮分别是：xyz、system_call和sys_xyz
     第一个就是API、第二个就是中断向量对应的这些也就是中断服务程序，中断向量对用的系统调用它有很多种不同的服务程序，比如sys_xyz,这就是三层皮。
      我们仔细看一下系统调用的服务历程：中断向量0x80与system_call绑定起来：
      当用户态进程调用一个系统调用时，CPU切换到内核态并开始执行第一个内核函数
      1.在Linux中是通过执行ini $0x80来执行系统调用的，这条汇编指令产生向量为128的编程异常
      2. Intel Pentium ll中引进了sysenter指令（快速系统调用）
系统调用号将xyz和sys_xyz关联起来了：
     传参：
      1.内核实现了很多不同的系统调用
      2.进程必须指明需要哪些系统调用，这需要传递一个系统调用号的参数，使用eax寄存器
     系统调用也需要输入输出参数，例如： 
     1.实际的值 2.用户态进程地址空间的变量的地址 3.甚至是包含指向用户态函数的指针的数据结构的地址
 system_call是linux中所有系统调用的入口点，每个系统调用至少有一个参数，即由eax传递的系统调用号
     2.一个应用程序调用fork(0封装例程，那么在执行int $0x80之前就把eax寄存器的值置为2（即_NR_fork)
     3.这个寄存器的设置是libc库中封装例程进行的，因此用户一般不关心系统调用号
     4.进入sys_call之后，立即将eax的值压入内核堆栈
寄存器传递参数有如下限制：
    1.每个参数的长度不能超过寄存器的长度，即32位
     2.在系统调用号eax之外，参数的个数不能超过6个（ebx,ecx,edx,esi,edi,ebp)
     超过6个怎么办？做一个把某个寄存器作为指针，指向一块内存，这样进入内核态之后可以访问所有内存空间，这就是系统调用的参数传递方式。
     
＃为什么系统调用开销大
平时说的系统调用开销大，主要是相对于函数调用来说的。对于一个函数调用，汇编层面上就是一个CALL或者JMP，这种指令在硬件层面上虽然首次是会打乱流水线的，但如果是十分有规律的情况下，大多数CPU都能很好的处理。对于一个CALL指令来说，CPU层面上做的事情是（来自intel手册，near call）：When executing a near call, the processor does the following (see Figure 6-2):1. Pushes the current value of the EIP register on the stack.2. Loads the offset of the called procedure in the EIP register.3. Begins execution of the called procedure.When executing a near return, the processor performs these actions:1. Pops the top-of-stack value (the return instruction pointer) into the EIP register.2. If the RET instruction has an optional n argument, increments the stack pointer by the number of bytes specified with the n operand to release parameters from the stack.3. Resumes execution of the calling procedure.其实就是存入EIP，载入新的EIP，执行；对于系统调用来说，麻烦就大了，过去Linux采用的是INT 80H中断的方式处理系统调用，一个带有栈切换的中断的流程如下：If a stack switch does occur, the processor does the following:1. Temporarily saves (internally) the current contents of the SS, ESP, EFLAGS, CS, and EIP registers.2. Loads the segment selector and stack pointer for the new stack (that is, the stack for the privilege level being called) from the TSS into the SS and ESP registers and switches to the new stack.3. Pushes the temporarily saved SS, ESP, EFLAGS, CS, and EIP values for the interrupted procedure’s stack onto the new stack.4. Pushes an error code on the new stack (if appropriate).5. Loads the segment selector for the new code segment and the new instruction pointer (from the interrupt gate or trap gate) into the CS and EIP registers, respectively.6. If the call is through an interrupt gate, clears the IF flag in the EFLAGS register.7. Begins execution of the handler procedure at the new privilege level.A return from an interrupt or exception handler is initiated with the IRET instruction. The IRET instruction is similar to the far RET instruction, except that it also restores the contents of the EFLAGS register for the interrupted procedure.When executing a return from an interrupt or exception handler from the same privilege level as the interrupted procedure, the processor performs these actions:1. Restores the CS and EIP registers to their values prior to the interrupt or exception.2. Restores the EFLAGS register.3. Increments the stack pointer appropriately.4. Resumes execution of the interrupted procedure.简单点说，就是保存的东西多了，CPU要处理的事情也多了。系统调用指令本身的开销就比一般的CALL和JMP要多一些，是因为同时又要进行一些额外的检查（权限、有效性等）。同时，因为可能涉及到任务栈的切换，会导致部分cache失效，这在CPU性能上的损失是很大的，甚至会导致TLB失效的情况（老版本Linux有此问题）。并且，由于系统调用属于“大事”，一般操作系统都会在系统调用内部进行多次检查，这又导致了一部分软件层面上的开销。所以，系统调用开销大基本上可以总结为：1. CPU要做的事情太多；2. 软件要做的事情也太多；因为INT指令开销太大，所以intel后来推出了SYSCALL/SYSENTER/SYSEXIT指令，这些指令不再查IDT表项了，直接从寄存器里取值，同时这些指令不保存堆栈和返回地址，同时不搞那么多权限检查了（因为这些指令必然是R3-R0之间切换的），所以CPU的开销会小的多。而且由于没有堆栈切换，实际上对流水线基本上没有多少破坏，所以CPU的性能会提升很多。但整体来说，由于系统调用CPU要做的事情还是多于一般的CALL/JMP指令，所以系统调用的开销肯定要比一般的函数调用开销要大，并且大的多。大概只能这么泛泛的说说，建议看Intel的手册，这里头的坑实在太多。
   这是操作系统为程序员做的统一接口，为了有好的程序员用户体验，做了很多业务。而普通的move, load指令，就很单一，执行后马上执行下一条指令不用考虑下面这些乱七八糟的，所以很快。系统调用即软中断（x86为int 0x80指令），用户态调这个会陷入内核，操作系统要有很多的现场环境保存，各种寄存器值，压栈操作…… 待内核任务执行完后，要切换到用户态，这时还会伴随一次任务调度的抢占，运气差了，cpu做其他进程的事了，运气好也要有恢复原先的现场的操作，而且context_swicth都会涉及到加锁，解锁，cpu资源抢占(多核)，涉及到数据地址的，内核要把数据复制用户态，因为用户态内核态不同地址空间（除非做了内存映射）

在2.4.4版内核中，狭义上的系统调用共有221个，你可以在<内核源码目录>/include/asm-i386/unistd.h中找到它们的原本，也可以通过命令"man 2 syscalls"察看它们的目录（man pages的版本一般比较老，可能有很多最新的调用都没有包含在内）。

系统调用是怎么工作的？
一般的，进程是不能访问内核的。它不能访问内核所占内存空间也不能调用内核函数。CPU硬件决定了这些（这就是为什么它被称作"保护模式"）。系统调用是这些规则的一个例外。其原理是进程先用适当的值填充寄存器，然后调用一个特殊的指令，这个指令会跳到一个事先定义的内核中的一个位置（当然，这个位置是用户进程可读但是不可写的）。在Intel CPU中，这个由中断0x80实现。硬件知道一旦你跳到这个位置，你就不是在限制模式下运行的用户，而是作为操作系统的内核--所以你就可以为所欲为。

进程可以跳转到的内核位置叫做sysem_call。这个过程检查系统调用号，这个号码告诉内核进程请求哪种服务。然后，它查看系统调用表(sys_call_table)找到所调用的内核函数入口地址。接着，就调用函数，等返回后，做一些系统检查，最后返回到进程（或到其他进程，如果这个进程时间用尽）。如果你希望读这段代码，它在<内核源码目录>/kernel/entry.S，Entry(system_call)的下一行。

因为内核实现了许多不同的系统调用，为了区别他们，进程必须传递一个系统调用号的参数来识别所需的系统调用。EAX寄存器是负责传递系统调用号的。

系统调用处理程序执行下列操作：

（1）在内核栈保存大多数寄存器的内容（这个操作对所有的系统调用都是通用的，并用汇编语言编写）。

（2）调用系统调用服务例程的相应的C函数处理系统调用。

（3）通过syscall_exit_work()函数从系统调用返回（这个函数用汇编语言编写）。



1.初始化系统调用

内核初始化期间调用trap_init()函数建立IDT表（中断描述符表）中128号向量对应的表项，语句如下：

set_system_gate(SYSCALL_VECTOR, &system_call);

其中SYSCALL_VECTOR是一个宏定义，其值为0x80，该调用把下列值装入这个门描述符的相应域。

（1）段选择子：因为系统调用处理程序属于内核代码，填写内核代码段__KERNEL_CS的段选择子。

（2）偏移量：指向system_call()系统调用处理程序。

（3）类型：置为15.表示这个异常是一个陷阱门，相应的处理程序不禁止可屏蔽中断。

（4）DPL（描述符权级）:置为3。这就允许用户态进程调用这个异常处理函数。



2.system_call()函数

由于未学过汇编语言，可能理解的并不是很清晰。

system_call()函数实现了系统调用处理函数。它首先把系统调用号和这个异常处理程序可以用到的所有CPU寄存器保存到相应的栈中，当然，栈中还有CPU已自动保存的EFLAGS，CS，EPI，SS和ESP寄存器，也在DS和ES中装入内核数据段的段选择子。

然后对用户态进程传来的系统调用号进行有效性检查。如果这个号大于或等于NR_syscalls，系统调用处理程序终止。

如果系统调用号无效，跳转到syscall_babsys处执行，此时就把—ENOSYS值存放在栈中EAX寄存器（该寄存器即存放系统调用号也存放系统调用的返回值，前者为正数，后者为负数）所在的单元（从当前栈顶开始偏移为24的单元，即EXA寄存器所在的单元）。然后返回用户空间。当进程以这种方式恢复它在用户态的执行时，会在EXA中发现一个负数的返回码。

最后，根据EXA中所包含的系统调用号调用对应的服务例程，因为系统调用表种的每一表项占4个字节，因此首先要把EXA中的系统调用号乘以4再加上sys_call_table系统调用表的起始地址（相当于起始位置加偏移量从而找到所要找到的地址），然后从这个地址单元获取指向相应服务例程的指针，内核就找到了要调用的服务例程。

当服务例程执行结束时，system_call()从EAX获得他的返回值，并把这个返回值存放在栈中，让其位于用户态EXA寄存器曾存放的位置。然后执行syscall_exit代码段，终止系统调用处理程序的执行。



3.参数传递

与普通的函数相似，系统调用通常也需要输入/输出参数，这些参数可能是实际的值，也可能是函数的地址即用户态进程地址空间的变量。因为system_call()函数时LINUX中所有系统调用唯一的入口点，因此每隔系统调用至少由一个函数，即通过EXA寄存器传递来的系统调用号。例如，如果一个应用程序调用至少有一个参数，即通过EXA寄存器传递来的系统调用号。而EXA这个寄存器的设置是由libc种的封装例程进行的，所以程序员通常不用关心系统调用号。

普通函数的参数传递是通过把参数写入活动的程序栈（或用户态栈或内核态栈）。但是系统调用的参数通常是通过寄存器传递给系统调用处理程序的，然后在拷贝到内核态堆栈。（实际就是相应的服务例程中所需的参数）。

寄存器的使用时的系统调用处理程序的结构与其它的异常处理程序结构类似。

然而，为了用寄存器传递参数，必须满足以下两个条件。

（1）每隔参数的长度不得超过寄存器的长度，即32位。

（2）参数的个数不得超过6个（包括系统调用号），因为寄存器的数量是有限的。

第一条总能成立，因为根据POSIX的标准，不能存放在32位寄存器的长参数必须通过指定他们的地址来传递。

对于第二个条件，确实存在超过6个参数的系统调用。在这种情况下，用一个单独的寄存器指向进程地址空间中这些参数值所在的一个内存区间即可。当然，编程者并不用关心这个工作区。与任何C调用一样，当调用libc封装例程时，参数被自动的保存在栈中。封装例程将找到合适的方式把参数传递给内核。

system_call()使用SAVE_ALL宏把这些寄存器的值保存在内核堆栈态中。因此，当系统调用服务例程转换到内核堆栈态时，就会找到system_call()的返回地址，紧接着是存放在EAX中的参数（即系统调用的第一个参数），存放在其与寄存器的参数等。这种栈结构与普通函数调用的栈结构完全相同，因此，服务例程可以很容易的使用一般C语言构造的参数。

处理write()系统调用的sys_write()服务例程的声明如下：

int sys_write(unsigned int fd, const char *buf, unsigned int count);

C编辑器产生一个汇编语言函数，该函数可以在栈顶找到fd，buf和count参数，因为这些参数就位与返回地址的下面。

在少数情况下，系统调用不是用任何参数，但是相应的服务例程需要知道发出系统调用之前CPU寄存器中的内容。在这种情况下，一个类型位pt_regs的单独参数允许服务例程访问由SAVE_ALL宏保存在内核态堆栈中的值，例如系统调用fork()的服务例程sys_fork()：

int sys_fork(struct pt_regs regs);

服务例程的返回值必须写到EXA寄存器中，这是在执行return n 指令时由C编译程序自动完成时。



对于getpid（）函数的执行：

（1）程序调用libc库的封装函数getpid。该封装函数将系统调用号__NR_getpid压入EAX寄存器。

（2）在内核中首先执行system_call()函数，接着根据系统调用号在系统调用表中查到找到对应的系统调用服务例程sys_getpid。

（3）执行sys_getpid服务例程。

（4）执行完毕后，转入system_exit_work例程，从系统调用返回。

