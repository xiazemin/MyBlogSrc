---
title: uninterruptible D 状态
layout: post
category: linux
author: 夏泽民
---
D状态即无法中断的休眠进程，是由于在等待IO，比如磁盘IO，网络IO，其他外设IO，如果进程正在等待的IO在较长的时间内都没有响应，那么就被ps看到了，同时也就意味着很有可能有IO出了问题，可能是外设本身出了故障，也可能是比如挂载的远程文件系统已经不可访问等操作时出现的问题。

而task_uninterruptible状态存在的意义就在于，内核的某些处理流程是不能被打断的。如果响应异步信号，程序的执行流程中就会被插入一段用于处理异步信号的流程（这个插入的流程可能只存在于内核态，也可能延伸到用户态），于是原有的流程就被中断了。

在进程对某些硬件进行操作时（比如进程调用read系统调用对某个设备文件进行读操作，而read系统调用最终执行到对应设备驱动的代码，并与对应的物理设备进行交互），可能需要使用task_uninterruptible状态对进程进行保护，以避免进程与设备交互的过程被打断，造成设备陷入不可控的状态。这种情况下的task_uninterruptible状态总是非常短暂的，通过ps命令基本上不可能捕捉到。

我们通过vmstat 命令中procs下的b 可以来查看是否有处于uninterruptible 状态的进程。 该命令只能显示数量。

In computer operating systems terminology, a sleeping process can either be interruptible (woken via signals) or uninterruptible (woken explicitly). An uninterruptible sleep state is a sleep state that cannot handle a signal (such as waiting for disk or network IO (input/output)).

When the process is sleeping uninterruptibly, the signal will be noticed when the process returns from the system call or trap.

-- 这句是关键。 当处于uninterruptibly sleep 状态时，只有当进程从system 调用返回时，才通知signal。

A process which ends up in “D” state for any measurable length of time is trapped in the midst of a system call (usually an I/O operation on a device — thus the initial in the ps output).

Such a process cannot be killed — it would risk leaving the kernel in an inconsistent state, leading to a panic. In general you can consider this to be a bug in the device driver that the process is accessing.
引起D状态的根本原因是由于IO等待，若你对某个磁盘的IO操作特别频繁，就会造成后续的IO操作处于等待状态，即处于D状态。此时，你可以使用iostat命令查看当前操作的磁盘的IO是否达到瓶颈值，若达到可以从用户层面入手调整IO操作。
<!-- more -->
Linux进程状态：S (TASK_INTERRUPTIBLE)，可中断的睡眠状态。

处于这个状态的进程因为等待某某事件的发生（比如等待socket连接、等待信号量），而被挂起。这些进程的task_struct结构被放入对应事件的等待队列中。当这些事件发生时（由外部中断触发、或由其他进程触发），对应的等待队列中的一个或多个进程将被唤醒。

通过ps命令我们会看到，一般情况下，进程列表中的绝大多数进程都处于TASK_INTERRUPTIBLE状态（除非机器的负载很高）。毕竟CPU就这么一两个，进程动辄几十上百个，如果不是绝大多数进程都在睡眠，CPU又怎么响应得过来。

Linux进程状态：D (TASK_UNINTERRUPTIBLE)，不可中断的睡眠状态。

与TASK_INTERRUPTIBLE状态类似，进程处于睡眠状态，但是此刻进程是不可中断的。不可中断，指的并不是CPU不响应外部硬件的中断，而是指进程不响应异步信号。
绝大多数情况下，进程处在睡眠状态时，总是应该能够响应异步信号的。否则你将惊奇的发现，kill -9竟然杀不死一个正在睡眠的进程了！于是我们也很好理解，为什么ps命令看到的进程几乎不会出现TASK_UNINTERRUPTIBLE状态，而总是TASK_INTERRUPTIBLE状态。

而TASK_UNINTERRUPTIBLE状态存在的意义就在于，内核的某些处理流程是不能被打断的。如果响应异步信号，程序的执行流程中就会被插入一段用于处理异步信号的流程（这个插入的流程可能只存在于内核态，也可能延伸到用户态），于是原有的流程就被中断了。（参见《linux内核异步中断浅析》）
在进程对某些硬件进行操作时（比如进程调用read系统调用对某个设备文件进行读操作，而read系统调用最终执行到对应设备驱动的代码，并与对应的物理设备进行交互），可能需要使用TASK_UNINTERRUPTIBLE状态对进程进行保护，以避免进程与设备交互的过程被打断，造成设备陷入不可控的状态。这种情况下的TASK_UNINTERRUPTIBLE状态总是非常短暂的，通过ps命令基本上不可能捕捉到。

linux系统中也存在容易捕捉的TASK_UNINTERRUPTIBLE状态。执行vfork系统调用后，父进程将进入TASK_UNINTERRUPTIBLE状态，直到子进程调用exit或exec（参见《神奇的vfork》）。
通过下面的代码就能得到处于TASK_UNINTERRUPTIBLE状态的进程：

#include   
void main() {  
if (!vfork()) sleep(100);  
} 
编译运行，然后ps一下：

kouu@kouu-one:~/test$ ps -ax | grep a\.out  
4371 pts/0    D+     0:00 ./a.out  
4372 pts/0    S+     0:00 ./a.out  
4374 pts/1    S+     0:00 grep a.out 
然后我们可以试验一下TASK_UNINTERRUPTIBLE状态的威力。不管kill还是kill -9，这个TASK_UNINTERRUPTIBLE状态的父进程依然屹立不倒。

进程为什么会被置于uninterruptible sleep状态呢？处于uninterruptiblesleep状态的进程通常是在等待IO，比如磁盘IO，网络IO，其他外设IO，如果进程正在等待的IO在较长的时间内都没有响应，那么就很会不幸地被ps看到了，同时也就意味着很有可能有IO出了问题，可能是外设本身出了故障，也可能是比如挂载的远程文件系统已经不可访问了（由down掉的NFS服务器引起的D状态）。
正是因为得不到IO的相应，进程才进入了uninterruptible sleep状态，所以要想使进程从uninterruptiblesleep状态恢复，就得使进程等待的IO恢复，比如如果是因为从远程挂载的NFS卷不可访问导致进程进入uninterruptiblesleep状态的，那么可以通过恢复该NFS卷的连接来使进程的IO请求得到满足。
D状态，往往是由于 I/O 资源得不到满足，而引发等待，在内核源码 fs/proc/array.c 里，其文字定义为“ "D (disk sleep)", /* 2 */ ”（由此可知 D 原是Disk的打头字母），对应着 include/linux/sched.h 里的“ #define TASK_UNINTERRUPTIBLE 2 ”。举个例子，当 NFS 服务端关闭之时，若未事先 umount 相关目录，在 NFS 客户端执行 df 就会挂住整个登录会话，按 Ctrl+C 、Ctrl+Z 都无济于事。断开连接再登录，执行 ps axf 则看到刚才的 df 进程状态位已变成了 D ，kill -9 无法杀灭。正确的处理方式，是马上恢复 NFS 服务端，再度提供服务，刚才挂起的 df 进程发现了其苦苦等待的资源，便完成任务，自动消亡。若 NFS 服务端无法恢复服务，在 reboot 之前也应将 /etc/mtab 里的相关 NFS mount 项删除，以免 reboot 过程例行调用 netfs stop 时再次发生等待资源，导致系统重启过程挂起。
