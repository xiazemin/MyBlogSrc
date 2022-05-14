---
title: netpoller
layout: post
category: golang
author: 夏泽民
---
Go 基于 I/O multiplexing 和 goroutine scheduler 构建了一个简洁而高性能的原生网络模型(基于 Go 的 I/O 多路复用 netpoller )，提供了 goroutine-per-connection 这样简单的网络编程模式。在这种模式下，开发者使用的是同步的模式去编写异步的逻辑，极大地降低了开发者编写网络应用时的心智负担，且借助于 Go runtime scheduler 对 goroutines 的高效调度，这个原生网络模型不论从适用性还是性能上都足以满足绝大部分的应用场景。

然而，在工程性上能做到如此高的普适性和兼容性，最终暴露给开发者提供接口/模式如此简洁，其底层必然是基于非常复杂的封装，做了很多取舍，也有可能放弃了一些追求极致性能的设计和理念。事实上 Go netpoller 底层就是基于 epoll/kqueue/iocp 这些 I/O 多路复用技术来做封装的，最终暴露出 goroutine-per-connection 这样的极简的开发模式给使用者。

Go netpoller 在不同的操作系统，其底层使用的 I/O 多路复用技术也不一样，可以从 Go 源码目录结构和对应代码文件了解 Go 在不同平台下的网络 I/O 模式的实现。比如，在 Linux 系统下基于 epoll，freeBSD 系统下基于 kqueue，以及 Windows 系统下基于 iocp。

本文将基于 Linux 平台来解析 Go netpoller 之 I/O 多路复用的底层是如何基于 epoll 封装实现的，从源码层层推进，全面而深度地解析 Go netpoller 的设计理念和实现原理，以及 Go 是如何利用 netpoller 来构建它的原生网络模型的。主要涉及到的一些概念：I/O 模型、用户/内核空间、epoll、Linux 源码、goroutine scheduler 等等，我会尽量简单地讲解，如果有对相关概念不熟悉的同学，还是希望能提前熟悉一下。

用户空间与内核空间
现代操作系统都是采用虚拟存储器，那么对 32 位操作系统而言，它的寻址空间（虚拟存储空间）为 4G（2 的 32 次方）。操作系统的核心是内核，独立于普通的应用程序，可以访问受保护的内存空间，也有访问底层硬件设备的所有权限。为了保证用户进程不能直接操作内核（kernel），保证内核的安全，操心系统将虚拟空间划分为两部分，一部分为内核空间，一部分为用户空间。针对 Linux 操作系统而言，将最高的 1G 字节（从虚拟地址 0xC0000000 到 0xFFFFFFFF），供内核使用，称为内核空间，而将较低的 3G 字节（从虚拟地址 0x00000000 到 0xBFFFFFFF），供各个进程使用，称为用户空间。



现代的网络服务的主流已经完成从 CPU 密集型到 IO 密集型的转变，所以服务端程序对 I/O 的处理必不可少，而一旦操作 I/O 则必定要在用户态和内核态之间来回切换。

I/O 模型
在神作《UNIX 网络编程》里，总结归纳了 5 种 I/O 模型，包括同步和异步 I/O：

阻塞 I/O (Blocking I/O)
非阻塞 I/O (Nonblocking I/O)
I/O 多路复用 (I/O multiplexing)
信号驱动 I/O (Signal driven I/O)
异步 I/O (Asynchronous I/O)
操作系统上的 I/O 是用户空间和内核空间的数据交互，因此 I/O 操作通常包含以下两个步骤：

等待网络数据到达网卡(读就绪)/等待网卡可写(写就绪) –> 读取/写入到内核缓冲区
从内核缓冲区复制数据 –> 用户空间(读)/从用户空间复制数据 -> 内核缓冲区(写)
而判定一个 I/O 模型是同步还是异步，主要看第二步：数据在用户和内核空间之间复制的时候是不是会阻塞当前进程，如果会，则是同步 I/O，否则，就是异步 I/O。基于这个原则，这 5 种 I/O 模型中只有一种异步 I/O 模型：Asynchronous I/O，其余都是同步 I/O 模型。

这 5 种 I/O 模型的对比如下：



Non-blocking I/O
什么叫非阻塞 I/O，顾名思义就是：所有 I/O 操作都是立刻返回而不会阻塞当前用户进程。I/O 多路复用通常情况下需要和非阻塞 I/O 搭配使用，否则可能会产生意想不到的问题。比如，epoll 的 ET(边缘触发) 模式下，如果不使用非阻塞 I/O，有极大的概率会导致阻塞 event-loop 线程，从而降低吞吐量，甚至导致 bug。

Linux 下，我们可以通过 fcntl 系统调用来设置 O_NONBLOCK 标志位，从而把 socket 设置成 Non-blocking。当对一个 Non-blocking socket 执行读操作时，流程是这个样子：


当用户进程发出 read 操作时，如果 kernel 中的数据还没有准备好，那么它并不会 block 用户进程，而是立刻返回一个 EAGAIN error。从用户进程角度讲 ，它发起一个 read 操作后，并不需要等待，而是马上就得到了一个结果。用户进程判断结果是一个 error 时，它就知道数据还没有准备好，于是它可以再次发送 read 操作。一旦 kernel 中的数据准备好了，并且又再次收到了用户进程的 system call，那么它马上就将数据拷贝到了用户内存，然后返回。

所以，Non-blocking I/O 的特点是用户进程需要不断的主动询问 kernel 数据好了没有。下一节我们要讲的 I/O 多路复用需要和 Non-blocking I/O 配合才能发挥出最大的威力！

I/O 多路复用
所谓 I/O 多路复用指的就是 select/poll/epoll 这一系列的多路选择器：支持单一线程同时监听多个文件描述符（I/O 事件），阻塞等待，并在其中某个文件描述符可读写时收到通知。 I/O 复用其实复用的不是 I/O 连接，而是复用线程，让一个 thread of control 能够处理多个连接（I/O 事件）。

select & poll
 1#include <sys/select.h>
 2
 3/* According to earlier standards */
 4#include <sys/time.h>
 5#include <sys/types.h>
 6#include <unistd.h>
 7
 8int select(int nfds, fd_set *readfds, fd_set *writefds, fd_set *exceptfds, struct timeval *timeout);
 9
10// 和 select 紧密结合的四个宏：
11void FD_CLR(int fd, fd_set *set);
12int FD_ISSET(int fd, fd_set *set);
13void FD_SET(int fd, fd_set *set);
14void FD_ZERO(fd_set *set);
select 是 epoll 之前 Linux 使用的 I/O 事件驱动技术。

理解 select 的关键在于理解 fd_set，为说明方便，取 fd_set 长度为 1 字节，fd_set 中的每一 bit 可以对应一个文件描述符 fd，则 1 字节长的 fd_set 最大可以对应 8 个 fd。select 的调用过程如下：

执行 FD_ZERO(&set), 则 set 用位表示是 0000,0000
若 fd＝5, 执行 FD_SET(fd, &set); 后 set 变为 0001,0000(第 5 位置为 1)
再加入 fd＝2, fd=1，则 set 变为 0001,0011
执行 select(6, &set, 0, 0, 0) 阻塞等待
若 fd=1, fd=2 上都发生可读事件，则 select 返回，此时 set 变为 0000,0011 (注意：没有事件发生的 fd=5 被清空)
基于上面的调用过程，可以得出 select 的特点：

可监控的文件描述符个数取决于 sizeof(fd_set) 的值。假设服务器上 sizeof(fd_set)＝512，每 bit 表示一个文件描述符，则服务器上支持的最大文件描述符是 512*8=4096。fd_set 的大小调整可参考 【原创】技术系列之 网络模型（二） 中的模型 2，可以有效突破 select 可监控的文件描述符上限
将 fd 加入 select 监控集的同时，还要再使用一个数据结构 array 保存放到 select 监控集中的 fd，一是用于在 select 返回后，array 作为源数据和 fd_set 进行 FD_ISSET 判断。二是 select 返回后会把以前加入的但并无事件发生的 fd 清空，则每次开始 select 前都要重新从 array 取得 fd 逐一加入（FD_ZERO 最先），扫描 array 的同时取得 fd 最大值 maxfd，用于 select 的第一个参数
可见 select 模型必须在 select 前循环 array（加 fd，取 maxfd），select 返回后循环 array（FD_ISSET 判断是否有事件发生）
所以，select 有如下的缺点：

最大并发数限制：使用 32 个整数的 32 位，即 32*32=1024 来标识 fd，虽然可修改，但是有以下第 2, 3 点的瓶颈
每次调用 select，都需要把 fd 集合从用户态拷贝到内核态，这个开销在 fd 很多时会很大
性能衰减严重：每次 kernel 都需要线性扫描整个 fd_set，所以随着监控的描述符 fd 数量增长，其 I/O 性能会线性下降
poll 的实现和 select 非常相似，只是描述 fd 集合的方式不同，poll 使用 pollfd 结构而不是 select 的 fd_set 结构，poll 解决了最大文件描述符数量限制的问题，但是同样需要从用户态拷贝所有的 fd 到内核态，也需要线性遍历所有的 fd 集合，所以它和 select 只是实现细节上的区分，并没有本质上的区别。

epoll
epoll 是 Linux kernel 2.6 之后引入的新 I/O 事件驱动技术，I/O 多路复用的核心设计是 1 个线程处理所有连接的 等待消息准备好 I/O 事件，这一点上 epoll 和 select&poll 是大同小异的。但 select&poll 错误预估了一件事，当数十万并发连接存在时，可能每一毫秒只有数百个活跃的连接，同时其余数十万连接在这一毫秒是非活跃的。select&poll 的使用方法是这样的： 返回的活跃连接 == select(全部待监控的连接) 。

什么时候会调用 select&poll 呢？在你认为需要找出有报文到达的活跃连接时，就应该调用。所以，select&poll 在高并发时是会被频繁调用的。这样，这个频繁调用的方法就很有必要看看它是否有效率，因为，它的轻微效率损失都会被 高频 二字所放大。它有效率损失吗？显而易见，全部待监控连接是数以十万计的，返回的只是数百个活跃连接，这本身就是无效率的表现。被放大后就会发现，处理并发上万个连接时，select&poll 就完全力不从心了。这个时候就该 epoll 上场了，epoll 通过一些新的设计和优化，基本上解决了 select&poll 的问题。

epoll 的 API 非常简洁，涉及到的只有 3 个系统调用：

1#include <sys/epoll.h>  
2int epoll_create(int size); // int epoll_create1(int flags);
3int epoll_ctl(int epfd, int op, int fd, struct epoll_event *event);
4int epoll_wait(int epfd, struct epoll_event *events, int maxevents, int timeout);
其中，epoll_create 创建一个 epoll 实例并返回 epollfd；epoll_ctl 注册 file descriptor 等待的 I/O 事件(比如 EPOLLIN、EPOLLOUT 等) 到 epoll 实例上；epoll_wait 则是阻塞监听 epoll 实例上所有的 file descriptor 的 I/O 事件，它接收一个用户空间上的一块内存地址 (events 数组)，kernel 会在有 I/O 事件发生的时候把文件描述符列表复制到这块内存地址上，然后 epoll_wait 解除阻塞并返回，最后用户空间上的程序就可以对相应的 fd 进行读写了：

1#include <unistd.h>
2ssize_t read(int fd, void *buf, size_t count);
3ssize_t write(int fd, const void *buf, size_t count);
epoll 的工作原理如下：



与 select&poll 相比，epoll 分清了高频调用和低频调用。例如，epoll_ctl 相对来说就是非频繁调用的，而 epoll_wait 则是会被高频调用的。所以 epoll 利用 epoll_ctl 来插入或者删除一个 fd，实现用户态到内核态的数据拷贝，这确保了每一个 fd 在其生命周期只需要被拷贝一次，而不是每次调用 epoll_wait 的时候都拷贝一次。 epoll_wait 则被设计成几乎没有入参的调用，相比 select&poll 需要把全部监听的 fd 集合从用户态拷贝至内核态的做法，epoll 的效率就高出了一大截。

在实现上 epoll 采用红黑树来存储所有监听的 fd，而红黑树本身插入和删除性能比较稳定，时间复杂度 O(logN)。通过 epoll_ctl 函数添加进来的 fd 都会被放在红黑树的某个节点内，所以，重复添加是没有用的。当把 fd 添加进来的时候时候会完成关键的一步：该 fd 会与相应的设备（网卡）驱动程序建立回调关系，也就是在内核中断处理程序为它注册一个回调函数，在 fd 相应的事件触发（中断）之后（设备就绪了），内核就会调用这个回调函数，该回调函数在内核中被称为： ep_poll_callback ，这个回调函数其实就是把这个 fd 添加到 rdllist 这个双向链表（就绪链表）中。epoll_wait 实际上就是去检查 rdllist 双向链表中是否有就绪的 fd，当 rdllist 为空（无就绪 fd）时挂起当前进程，直到 rdllist 非空时进程才被唤醒并返回。

相比于 select&poll 调用时会将全部监听的 fd 从用户态空间拷贝至内核态空间并线性扫描一遍找出就绪的 fd 再返回到用户态，epoll_wait 则是直接返回已就绪 fd，因此 epoll 的 I/O 性能不会像 select&poll 那样随着监听的 fd 数量增加而出现线性衰减，是一个非常高效的 I/O 事件驱动技术。

由于使用 epoll 的 I/O 多路复用需要用户进程自己负责 I/O 读写，从用户进程的角度看，读写过程是阻塞的，所以 select&poll&epoll 本质上都是同步 I/O 模型，而像 Windows 的 IOCP 这一类的异步 I/O，只需要在调用 WSARecv 或 WSASend 方法读写数据的时候把用户空间的内存 buffer 提交给 kernel，kernel 负责数据在用户空间和内核空间拷贝，完成之后就会通知用户进程，整个过程不需要用户进程参与，所以是真正的异步 I/O。

延伸
另外，我看到有些文章说 epoll 之所以性能高是因为利用了 Linux 的 mmap 内存映射让内核和用户进程共享了一片物理内存，用来存放就绪 fd 列表和它们的数据 buffer，所以用户进程在 epoll_wait 返回之后用户进程就可以直接从共享内存那里读取/写入数据了，这让我很疑惑，因为首先看 epoll_wait 的函数声明：

1int epoll_wait(int epfd, struct epoll_event *events, int maxevents, int timeout);
第二个参数：就绪事件列表，是需要在用户空间分配内存然后再传给 epoll_wait 的，如果内核会用 mmap 设置共享内存，直接传递一个指针进去就行了，根本不需要在用户态分配内存，多此一举。其次，内核和用户进程通过 mmap 共享内存是一件极度危险的事情，内核无法确定这块共享内存什么时候会被回收，而且这样也会赋予用户进程直接操作内核数据的权限和入口，非常容易出现大的系统漏洞，因此一般极少会这么做。所以我很怀疑 epoll 是不是真的在 Linux kernel 里用了 mmap，我就去看了下最新版本（5.3.9）的 Linux kernel 源码：

 1/*
 2 * Implement the event wait interface for the eventpoll file. It is the kernel
 3 * part of the user space epoll_wait(2).
 4 */
 5static int do_epoll_wait(int epfd, struct epoll_event __user *events,
 6			 int maxevents, int timeout)
 7{
 8	...
 9  
10	/* Time to fish for events ... */
11	error = ep_poll(ep, events, maxevents, timeout);
12}
13
14// 如果 epoll_wait 入参时设定 timeout == 0, 那么直接通过 ep_events_available 判断当前是否有用户感兴趣的事件发生，如果有则通过 ep_send_events 进行处理
15// 如果设置 timeout > 0，并且当前没有用户关注的事件发生，则进行休眠，并添加到 ep->wq 等待队列的头部；对等待事件描述符设置 WQ_FLAG_EXCLUSIVE 标志
16// ep_poll 被事件唤醒后会重新检查是否有关注事件，如果对应的事件已经被抢走，那么 ep_poll 会继续休眠等待
17static int ep_poll(struct eventpoll *ep, struct epoll_event __user *events, int maxevents, long timeout)
18{
19	...
20  
21	send_events:
22	/*
23	 * Try to transfer events to user space. In case we get 0 events and
24	 * there's still timeout left over, we go trying again in search of
25	 * more luck.
26	 */
27  
28	// 如果一切正常, 有 event 发生, 就开始准备数据 copy 给用户空间了
29	// 如果有就绪的事件发生，那么就调用 ep_send_events 将就绪的事件 copy 到用户态内存中，
30	// 然后返回到用户态，否则判断是否超时，如果没有超时就继续等待就绪事件发生，如果超时就返回用户态。
31	// 从 ep_poll 函数的实现可以看到，如果有就绪事件发生，则调用 ep_send_events 函数做进一步处理
32	if (!res && eavail &&
33			!(res = ep_send_events(ep, events, maxevents)) && !timed_out)
34		goto fetch_events;
35  
36	...
37}
38
39// ep_send_events 函数是用来向用户空间拷贝就绪 fd 列表的，它将用户传入的就绪 fd 列表内存简单封装到
40// ep_send_events_data 结构中，然后调用 ep_scan_ready_list 将就绪队列中的事件写入用户空间的内存；
41// 用户进程就可以访问到这些数据进行处理
42static int ep_send_events(struct eventpoll *ep,
43				struct epoll_event __user *events, int maxevents)
44{
45	struct ep_send_events_data esed;
46
47	esed.maxevents = maxevents;
48	esed.events = events;
49	// 调用 ep_scan_ready_list 函数检查 epoll 实例 eventpoll 中的 rdllist 就绪链表，
50	// 并注册一个回调函数 ep_send_events_proc，如果有就绪 fd，则调用 ep_send_events_proc 进行处理
51	ep_scan_ready_list(ep, ep_send_events_proc, &esed, 0, false);
52	return esed.res;
53}
54
55// 调用 ep_scan_ready_list 的时候会传递指向 ep_send_events_proc 函数的函数指针作为回调函数，
56// 一旦有就绪 fd，就会调用 ep_send_events_proc 函数
57static __poll_t ep_send_events_proc(struct eventpoll *ep, struct list_head *head, void *priv)
58{
59	...
60  
61	/*
62	 * If the event mask intersect the caller-requested one,
63	 * deliver the event to userspace. Again, ep_scan_ready_list()
64	 * is holding ep->mtx, so no operations coming from userspace
65	 * can change the item.
66	 */
67	revents = ep_item_poll(epi, &pt, 1);
68	// 如果 revents 为 0，说明没有就绪的事件，跳过，否则就将就绪事件拷贝到用户态内存中
69	if (!revents)
70		continue;
71	// 将当前就绪的事件和用户进程传入的数据都通过 __put_user 拷贝回用户空间,
72	// 也就是调用 epoll_wait 之时用户进程传入的 fd 列表的内存
73	if (__put_user(revents, &uevent->events) || __put_user(epi->event.data, &uevent->data)) {
74		list_add(&epi->rdllink, head);
75		ep_pm_stay_awake(epi);
76		if (!esed->res)
77			esed->res = -EFAULT;
78		return 0;
79	}
80  
81	...
82}
从 do_epoll_wait 开始层层跳转，我们可以很清楚地看到最后内核是通过 __put_user 函数把就绪 fd 列表和事件返回到用户空间，而 __put_user 正是内核用来拷贝数据到用户空间的标准函数。此外，我并没有在 Linux kernel 的源码中和 epoll 相关的代码里找到 mmap 系统调用做内存映射的逻辑，所以基本可以得出结论：epoll 在 Linux kernel 里并没有使用 mmap 来做用户空间和内核空间的内存共享，所以那些说 epoll 使用了 mmap 的文章都是误解。

Go netpoller 核心
Go netpoller 基本原理

Go netpoller 通过在底层对 epoll/kqueue/iocp 的封装，从而实现了使用同步编程模式达到异步执行的效果。总结来说，所有的网络操作都以网络描述符 netFD 为中心实现。netFD 与底层 PollDesc 结构绑定，当在一个 netFD 上读写遇到 EAGAIN 错误时，就将当前 goroutine 存储到这个 netFD 对应的 PollDesc 中，同时调用 gopark 把当前 goroutine 给 park 住，直到这个 netFD 上再次发生读写事件，才将此 goroutine 给 ready 激活重新运行。显然，在底层通知 goroutine 再次发生读写等事件的方式就是 epoll/kqueue/iocp 等事件驱动机制。

总所周知，Go 是一门跨平台的编程语言，而不同平台针对特定的功能有不用的实现，这当然也包括了 I/O 多路复用技术，比如 Linux 里的 I/O 多路复用有 select、poll 和 epoll，而 freeBSD 或者 MacOS 里则是 kqueue，而 Windows 里则是基于异步 I/O 实现的 iocp，等等；因此，Go 为了实现底层 I/O 多路复用的跨平台，分别基于上述的这些不同平台的系统调用实现了多版本的 netpollers，具体的源码路径如下：

src/runtime/netpoll_epoll.go
src/runtime/netpoll_kqueue.go
src/runtime/netpoll_solaris.go
src/runtime/netpoll_windows.go
src/runtime/netpoll_aix.go
src/runtime/netpoll_fake.go
本文的解析基于 epoll 版本，如果读者对其他平台的 netpoller 底层实现感兴趣，可以在阅读完本文后自行翻阅其他 netpoller 源码，所有实现版本的机制和原理基本类似，所以了解了 epoll 版本的实现后再去学习其他版本实现应该没什么障碍。

接下来让我们通过分析最新的 Go 源码（v1.15.3），全面剖析一下整个 Go netpoller 的运行机制和流程。

数据结构
netFD
net.Listen("tcp", ":8888") 方法返回了一个 *TCPListener，它是一个实现了 net.Listener 接口的 struct，而通过 listener.Accept() 接收的新连接 *TCPConn 则是一个实现了 net.Conn 接口的 struct，它内嵌了 net.conn struct。仔细阅读上面的源码可以发现，不管是 Listener 的 Accept 还是 Conn 的 Read/Write 方法，都是基于一个 netFD 的数据结构的操作， netFD 是一个网络描述符，类似于 Linux 的文件描述符的概念，netFD 中包含一个 poll.FD 数据结构，而 poll.FD 中包含两个重要的数据结构 Sysfd 和 pollDesc，前者是真正的系统文件描述符，后者对是底层事件驱动的封装，所有的读写超时等操作都是通过调用后者的对应方法实现的。

netFD 和 poll.FD 的源码：

 1// Network file descriptor.
 2type netFD struct {
 3	pfd poll.FD
 4
 5	// immutable until Close
 6	family      int
 7	sotype      int
 8	isConnected bool // handshake completed or use of association with peer
 9	net         string
10	laddr       Addr
11	raddr       Addr
12}
13
14// FD is a file descriptor. The net and os packages use this type as a
15// field of a larger type representing a network connection or OS file.
16type FD struct {
17	// Lock sysfd and serialize access to Read and Write methods.
18	fdmu fdMutex
19
20	// System file descriptor. Immutable until Close.
21	Sysfd int
22
23	// I/O poller.
24	pd pollDesc
25
26	// Writev cache.
27	iovecs *[]syscall.Iovec
28
29	// Semaphore signaled when file is closed.
30	csema uint32
31
32	// Non-zero if this file has been set to blocking mode.
33	isBlocking uint32
34
35	// Whether this is a streaming descriptor, as opposed to a
36	// packet-based descriptor like a UDP socket. Immutable.
37	IsStream bool
38
39	// Whether a zero byte read indicates EOF. This is false for a
40	// message based socket connection.
41	ZeroReadIsEOF bool
42
43	// Whether this is a file rather than a network socket.
44	isFile bool
45}
pollDesc
前面提到了 pollDesc 是底层事件驱动的封装，netFD 通过它来完成各种 I/O 相关的操作，它的定义如下：

1type pollDesc struct {
2	runtimeCtx uintptr
3}
这里的 struct 只包含了一个指针，而通过 pollDesc 的 init 方法，我们可以找到它具体的定义是在 runtime.pollDesc 这里：

 1func (pd *pollDesc) init(fd *FD) error {
 2	serverInit.Do(runtime_pollServerInit)
 3	ctx, errno := runtime_pollOpen(uintptr(fd.Sysfd))
 4	if errno != 0 {
 5		if ctx != 0 {
 6			runtime_pollUnblock(ctx)
 7			runtime_pollClose(ctx)
 8		}
 9		return syscall.Errno(errno)
10	}
11	pd.runtimeCtx = ctx
12	return nil
13}
14
15// Network poller descriptor.
16//
17// No heap pointers.
18//
19//go:notinheap
20type pollDesc struct {
21	link *pollDesc // in pollcache, protected by pollcache.lock
22
23	// The lock protects pollOpen, pollSetDeadline, pollUnblock and deadlineimpl operations.
24	// This fully covers seq, rt and wt variables. fd is constant throughout the PollDesc lifetime.
25	// pollReset, pollWait, pollWaitCanceled and runtime·netpollready (IO readiness notification)
26	// proceed w/o taking the lock. So closing, everr, rg, rd, wg and wd are manipulated
27	// in a lock-free way by all operations.
28	// NOTE(dvyukov): the following code uses uintptr to store *g (rg/wg),
29	// that will blow up when GC starts moving objects.
30	lock    mutex // protects the following fields
31	fd      uintptr
32	closing bool
33	everr   bool    // marks event scanning error happened
34	user    uint32  // user settable cookie
35	rseq    uintptr // protects from stale read timers
36	rg      uintptr // pdReady, pdWait, G waiting for read or nil
37	rt      timer   // read deadline timer (set if rt.f != nil)
38	rd      int64   // read deadline
39	wseq    uintptr // protects from stale write timers
40	wg      uintptr // pdReady, pdWait, G waiting for write or nil
41	wt      timer   // write deadline timer
42	wd      int64   // write deadline
43}
这里重点关注里面的 rg 和 wg，这里两个 uintptr "万能指针"类型，取值分别可能是 pdReady、pdWait、等待 file descriptor 就绪的 goroutine 也就是 g 数据结构以及 nil，它们是实现唤醒 goroutine 的关键。

runtime.pollDesc 包含自身类型的一个指针，用来保存下一个 runtime.pollDesc 的地址，以此来实现链表，可以减少数据结构的大小，所有的 runtime.pollDesc 保存在 runtime.pollCache 结构中，定义如下：

1type pollCache struct {
2	lock  mutex
3	first *pollDesc
4	// PollDesc objects must be type-stable,
5	// because we can get ready notification from epoll/kqueue
6	// after the descriptor is closed/reused.
7	// Stale notifications are detected using seq variable,
8	// seq is incremented when deadlines are changed or descriptor is reused.
9}
因为 runtime.pollCache 是一个在 runtime 包里的全局变量，因此需要用一个互斥锁来避免 data race 问题，从它的名字也能看出这是一个用于缓存的数据结构，也就是用来提高性能的，具体如何实现呢？

 1const pollBlockSize = 4 * 1024
 2
 3func (c *pollCache) alloc() *pollDesc {
 4	lock(&c.lock)
 5	if c.first == nil {
 6		const pdSize = unsafe.Sizeof(pollDesc{})
 7		n := pollBlockSize / pdSize
 8		if n == 0 {
 9			n = 1
10		}
11		// Must be in non-GC memory because can be referenced
12		// only from epoll/kqueue internals.
13		mem := persistentalloc(n*pdSize, 0, &memstats.other_sys)
14		for i := uintptr(0); i < n; i++ {
15			pd := (*pollDesc)(add(mem, i*pdSize))
16			pd.link = c.first
17			c.first = pd
18		}
19	}
20	pd := c.first
21	c.first = pd.link
22	lockInit(&pd.lock, lockRankPollDesc)
23	unlock(&c.lock)
24	return pd
25}
Go runtime 会在调用 poll_runtime_pollOpen 往 epoll 实例注册 fd 之时首次调用 runtime.pollCache.alloc方法时批量初始化大小 4KB 的 runtime.pollDesc 结构体的链表，初始化过程中会调用 runtime.persistentalloc 来为这些数据结构分配不会被 GC 回收的内存，确保这些数据结构只能被 epoll和kqueue 在内核空间去引用。

再往后每次调用这个方法则会先判断链表头是否已经分配过值了，若是，则直接返回表头这个 pollDesc，这种批量初始化数据进行缓存而后每次都直接从缓存取数据的方式是一种很常见的性能优化手段，在这里这种方式可以有效地提升 netpoller 的吞吐量。

Go runtime 会在关闭 pollDesc 之时调用 runtime.pollCache.free 释放内存：

1func (c *pollCache) free(pd *pollDesc) {
2	lock(&c.lock)
3	pd.link = c.first
4	c.first = pd
5	unlock(&c.lock)
6}
实现原理
使用 Go 编写一个典型的 TCP echo server:

 1package main
 2
 3import (
 4	"log"
 5	"net"
 6)
 7
 8func main() {
 9	listen, err := net.Listen("tcp", ":8888")
10	if err != nil {
11		log.Println("listen error: ", err)
12		return
13	}
14
15	for {
16		conn, err := listen.Accept()
17		if err != nil {
18			log.Println("accept error: ", err)
19			break
20		}
21
22		// start a new goroutine to handle the new connection.
23		go HandleConn(conn)
24	}
25}
26
27func HandleConn(conn net.Conn) {
28	defer conn.Close()
29	packet := make([]byte, 1024)
30	for {
31		// block here if socket is not available for reading data.
32		n, err := conn.Read(packet)
33		if err != nil {
34			log.Println("read socket error: ", err)
35			return
36		}
37
38		// same as above, block here if socket is not available for writing.
39		_, _ = conn.Write(packet[:n])
40	}
41}
上面是一个基于 Go 原生网络模型（基于 netpoller）编写的一个 TCP server，模式是 goroutine-per-connection ，在这种模式下，开发者使用的是同步的模式去编写异步的逻辑而且对于开发者来说 I/O 是否阻塞是无感知的，也就是说开发者无需考虑 goroutines 甚至更底层的线程、进程的调度和上下文切换。而 Go netpoller 最底层的事件驱动技术肯定是基于 epoll/kqueue/iocp 这一类的 I/O 事件驱动技术，只不过是把这些调度和上下文切换的工作转移到了 runtime 的 Go scheduler，让它来负责调度 goroutines，从而极大地降低了程序员的心智负担！

Go 的这种同步模式的网络服务器的基本架构通常如下：



上面的示例代码中相关的在源码里的几个数据结构和方法：

 1// TCPListener is a TCP network listener. Clients should typically
 2// use variables of type Listener instead of assuming TCP.
 3type TCPListener struct {
 4	fd *netFD
 5	lc ListenConfig
 6}
 7
 8// Accept implements the Accept method in the Listener interface; it
 9// waits for the next call and returns a generic Conn.
10func (l *TCPListener) Accept() (Conn, error) {
11	if !l.ok() {
12		return nil, syscall.EINVAL
13	}
14	c, err := l.accept()
15	if err != nil {
16		return nil, &OpError{Op: "accept", Net: l.fd.net, Source: nil, Addr: l.fd.laddr, Err: err}
17	}
18	return c, nil
19}
20
21func (ln *TCPListener) accept() (*TCPConn, error) {
22	fd, err := ln.fd.accept()
23	if err != nil {
24		return nil, err
25	}
26	tc := newTCPConn(fd)
27	if ln.lc.KeepAlive >= 0 {
28		setKeepAlive(fd, true)
29		ka := ln.lc.KeepAlive
30		if ln.lc.KeepAlive == 0 {
31			ka = defaultTCPKeepAlive
32		}
33		setKeepAlivePeriod(fd, ka)
34	}
35	return tc, nil
36}
37
38// TCPConn is an implementation of the Conn interface for TCP network
39// connections.
40type TCPConn struct {
41	conn
42}
43
44// Conn
45type conn struct {
46	fd *netFD
47}
48
49type conn struct {
50	fd *netFD
51}
52
53func (c *conn) ok() bool { return c != nil && c.fd != nil }
54
55// Implementation of the Conn interface.
56
57// Read implements the Conn Read method.
58func (c *conn) Read(b []byte) (int, error) {
59	if !c.ok() {
60		return 0, syscall.EINVAL
61	}
62	n, err := c.fd.Read(b)
63	if err != nil && err != io.EOF {
64		err = &OpError{Op: "read", Net: c.fd.net, Source: c.fd.laddr, Addr: c.fd.raddr, Err: err}
65	}
66	return n, err
67}
68
69// Write implements the Conn Write method.
70func (c *conn) Write(b []byte) (int, error) {
71	if !c.ok() {
72		return 0, syscall.EINVAL
73	}
74	n, err := c.fd.Write(b)
75	if err != nil {
76		err = &OpError{Op: "write", Net: c.fd.net, Source: c.fd.laddr, Addr: c.fd.raddr, Err: err}
77	}
78	return n, err
79}
net.Listen
调用 net.Listen 之后，底层会通过 Linux 的系统调用 socket 方法创建一个 fd 分配给 listener，并用以来初始化 listener 的 netFD ，接着调用 netFD 的 listenStream 方法完成对 socket 的 bind&listen 操作以及对 netFD 的初始化（主要是对 netFD 里的 pollDesc 的初始化），调用链是 runtime.runtime_pollServerInit --> runtime.poll_runtime_pollServerInit --> runtime.netpollGenericInit，主要做的事情是：

调用 epollcreate1 创建一个 epoll 实例 epfd，作为整个 runtime 的唯一 event-loop 使用；
调用 runtime.nonblockingPipe 创建一个用于和 epoll 实例通信的管道，这里为什么不用更新且更轻量的 eventfd 呢？我个人猜测是为了兼容更多以及更老的系统版本；
将 netpollBreakRd 通知信号量封装成 epollevent 事件结构体注册进 epoll 实例。
相关源码如下：

  1// 调用 linux 系统调用 socket 创建 listener fd 并设置为为阻塞 I/O
  2s, err := socketFunc(family, sotype|syscall.SOCK_NONBLOCK|syscall.SOCK_CLOEXEC, proto)
  3// On Linux the SOCK_NONBLOCK and SOCK_CLOEXEC flags were
  4// introduced in 2.6.27 kernel and on FreeBSD both flags were
  5// introduced in 10 kernel. If we get an EINVAL error on Linux
  6// or EPROTONOSUPPORT error on FreeBSD, fall back to using
  7// socket without them.
  8
  9socketFunc        func(int, int, int) (int, error)  = syscall.Socket
 10
 11// 用上面创建的 listener fd 初始化 listener netFD
 12if fd, err = newFD(s, family, sotype, net); err != nil {
 13	poll.CloseFunc(s)
 14	return nil, err
 15}
 16
 17// 对 listener fd 进行 bind&listen 操作，并且调用 init 方法完成初始化
 18func (fd *netFD) listenStream(laddr sockaddr, backlog int, ctrlFn func(string, string, syscall.RawConn) error) error {
 19	...
 20  
 21	// 完成绑定操作
 22	if err = syscall.Bind(fd.pfd.Sysfd, lsa); err != nil {
 23		return os.NewSyscallError("bind", err)
 24	}
 25  
 26	// 完成监听操作
 27	if err = listenFunc(fd.pfd.Sysfd, backlog); err != nil {
 28		return os.NewSyscallError("listen", err)
 29	}
 30  
 31	// 调用 init，内部会调用 poll.FD.Init，最后调用 pollDesc.init
 32	if err = fd.init(); err != nil {
 33		return err
 34	}
 35	lsa, _ = syscall.Getsockname(fd.pfd.Sysfd)
 36	fd.setAddr(fd.addrFunc()(lsa), nil)
 37	return nil
 38}
 39
 40// 使用 sync.Once 来确保一个 listener 只持有一个 epoll 实例
 41var serverInit sync.Once
 42
 43// netFD.init 会调用 poll.FD.Init 并最终调用到 pollDesc.init，
 44// 它会创建 epoll 实例并把 listener fd 加入监听队列
 45func (pd *pollDesc) init(fd *FD) error {
 46	// runtime_pollServerInit 通过 `go:linkname` 链接到具体的实现函数 poll_runtime_pollServerInit，
 47	// 接着再调用 netpollGenericInit，然后会根据不同的系统平台去调用特定的 netpollinit 来创建 epoll 实例
 48	serverInit.Do(runtime_pollServerInit)
 49  
 50	// runtime_pollOpen 内部调用了 netpollopen 来将 listener fd 注册到 
 51	// epoll 实例中，另外，它会初始化一个 pollDesc 并返回
 52	ctx, errno := runtime_pollOpen(uintptr(fd.Sysfd))
 53	if errno != 0 {
 54		if ctx != 0 {
 55			runtime_pollUnblock(ctx)
 56			runtime_pollClose(ctx)
 57		}
 58		return syscall.Errno(errno)
 59	}
 60	// 把真正初始化完成的 pollDesc 实例赋值给当前的 pollDesc 代表自身的指针，
 61	// 后续使用直接通过该指针操作
 62	pd.runtimeCtx = ctx
 63	return nil
 64}
 65
 66var (
 67	// 全局唯一的 epoll fd，只在 listener fd 初始化之时被指定一次
 68	epfd int32 = -1 // epoll descriptor
 69)
 70
 71// netpollinit 会创建一个 epoll 实例，然后把 epoll fd 赋值给 epfd，
 72// 后续 listener 以及它 accept 的所有 sockets 有关 epoll 的操作都是基于这个全局的 epfd
 73func netpollinit() {
 74	epfd = epollcreate1(_EPOLL_CLOEXEC)
 75	if epfd < 0 {
 76		epfd = epollcreate(1024)
 77		if epfd < 0 {
 78			println("runtime: epollcreate failed with", -epfd)
 79			throw("runtime: netpollinit failed")
 80		}
 81		closeonexec(epfd)
 82	}
 83	r, w, errno := nonblockingPipe()
 84	if errno != 0 {
 85		println("runtime: pipe failed with", -errno)
 86		throw("runtime: pipe failed")
 87	}
 88	ev := epollevent{
 89		events: _EPOLLIN,
 90	}
 91	*(**uintptr)(unsafe.Pointer(&ev.data)) = &netpollBreakRd
 92	errno = epollctl(epfd, _EPOLL_CTL_ADD, r, &ev)
 93	if errno != 0 {
 94		println("runtime: epollctl failed with", -errno)
 95		throw("runtime: epollctl failed")
 96	}
 97	netpollBreakRd = uintptr(r)
 98	netpollBreakWr = uintptr(w)
 99}
100
101// netpollopen 会被 runtime_pollOpen 调用，注册 fd 到 epoll 实例，
102// 注意这里使用的是 epoll 的 ET 模式，同时会利用万能指针把 pollDesc 保存到 epollevent 的一个 8 位的字节数组 data 里
103func netpollopen(fd uintptr, pd *pollDesc) int32 {
104	var ev epollevent
105	ev.events = _EPOLLIN | _EPOLLOUT | _EPOLLRDHUP | _EPOLLET
106	*(**pollDesc)(unsafe.Pointer(&ev.data)) = pd
107	return -epollctl(epfd, _EPOLL_CTL_ADD, int32(fd), &ev)
108}
我们前面提到的 epoll 的三个基本调用，Go 在源码里实现了对那三个调用的封装：

1#include <sys/epoll.h>  
2int epoll_create(int size);  
3int epoll_ctl(int epfd, int op, int fd, struct epoll_event *event);  
4int epoll_wait(int epfd, struct epoll_event * events, int maxevents, int timeout);
5
6// Go 对上面三个调用的封装
7func netpollinit()
8func netpollopen(fd uintptr, pd *pollDesc) int32
9func netpoll(block bool) gList
netFD 就是通过这三个封装来对 epoll 进行创建实例、注册 fd 和等待事件操作的。

Listener.Accept()
netpoll accept socket 的工作流程如下：

服务端的 netFD 在 listen 时会创建 epoll 的实例，并将 listenerFD 加入 epoll 的事件队列
netFD 在 accept 时将返回的 connFD 也加入 epoll 的事件队列
netFD 在读写时出现 syscall.EAGAIN 错误，通过 pollDesc 的 waitRead 方法将当前的 goroutine park 住，直到 ready，从 pollDesc 的 waitRead 中返回
Listener.Accept() 接收来自客户端的新连接，具体还是调用 netFD.accept 方法来完成这个功能：

 1// Accept implements the Accept method in the Listener interface; it
 2// waits for the next call and returns a generic Conn.
 3func (l *TCPListener) Accept() (Conn, error) {
 4	if !l.ok() {
 5		return nil, syscall.EINVAL
 6	}
 7	c, err := l.accept()
 8	if err != nil {
 9		return nil, &OpError{Op: "accept", Net: l.fd.net, Source: nil, Addr: l.fd.laddr, Err: err}
10	}
11	return c, nil
12}
13
14func (ln *TCPListener) accept() (*TCPConn, error) {
15	fd, err := ln.fd.accept()
16	if err != nil {
17		return nil, err
18	}
19	tc := newTCPConn(fd)
20	if ln.lc.KeepAlive >= 0 {
21		setKeepAlive(fd, true)
22		ka := ln.lc.KeepAlive
23		if ln.lc.KeepAlive == 0 {
24			ka = defaultTCPKeepAlive
25		}
26		setKeepAlivePeriod(fd, ka)
27	}
28	return tc, nil
29}
30
31func (fd *netFD) accept() (netfd *netFD, err error) {
32	// 调用 poll.FD 的 Accept 方法接受新的 socket 连接，返回 socket 的 fd
33	d, rsa, errcall, err := fd.pfd.Accept()
34	if err != nil {
35		if errcall != "" {
36			err = wrapSyscallError(errcall, err)
37		}
38		return nil, err
39	}
40	// 以 socket fd 构造一个新的 netFD，代表这个新的 socket
41	if netfd, err = newFD(d, fd.family, fd.sotype, fd.net); err != nil {
42		poll.CloseFunc(d)
43		return nil, err
44	}
45	// 调用 netFD 的 init 方法完成初始化
46	if err = netfd.init(); err != nil {
47		fd.Close()
48		return nil, err
49	}
50	lsa, _ := syscall.Getsockname(netfd.pfd.Sysfd)
51	netfd.setAddr(netfd.addrFunc()(lsa), netfd.addrFunc()(rsa))
52	return netfd, nil
53}
netFD.accept 方法里会再调用 poll.FD.Accept ，最后会使用 Linux 的系统调用 accept 来完成新连接的接收，并且会把 accept 的 socket 设置成非阻塞 I/O 模式：

 1// Accept wraps the accept network call.
 2func (fd *FD) Accept() (int, syscall.Sockaddr, string, error) {
 3	if err := fd.readLock(); err != nil {
 4		return -1, nil, "", err
 5	}
 6	defer fd.readUnlock()
 7
 8	if err := fd.pd.prepareRead(fd.isFile); err != nil {
 9		return -1, nil, "", err
10	}
11	for {
12		// 使用 linux 系统调用 accept 接收新连接，创建对应的 socket
13		s, rsa, errcall, err := accept(fd.Sysfd)
14		// 因为 listener fd 在创建的时候已经设置成非阻塞的了，
15		// 所以 accept 方法会直接返回，不管有没有新连接到来；如果 err == nil 则表示正常建立新连接，直接返回
16		if err == nil {
17			return s, rsa, "", err
18		}
19		// 如果 err != nil，则判断 err == syscall.EAGAIN，符合条件则进入 pollDesc.waitRead 方法
20		switch err {
21		case syscall.EAGAIN:
22			if fd.pd.pollable() {
23				// 如果当前没有发生期待的 I/O 事件，那么 waitRead 会通过 park goroutine 让逻辑 block 在这里
24				if err = fd.pd.waitRead(fd.isFile); err == nil {
25					continue
26				}
27			}
28		case syscall.ECONNABORTED:
29			// This means that a socket on the listen
30			// queue was closed before we Accept()ed it;
31			// it's a silly error, so try again.
32			continue
33		}
34		return -1, nil, errcall, err
35	}
36}
37
38// 使用 linux 的 accept 系统调用接收新连接并把这个 socket fd 设置成非阻塞 I/O
39ns, sa, err := Accept4Func(s, syscall.SOCK_NONBLOCK|syscall.SOCK_CLOEXEC)
40// On Linux the accept4 system call was introduced in 2.6.28
41// kernel and on FreeBSD it was introduced in 10 kernel. If we
42// get an ENOSYS error on both Linux and FreeBSD, or EINVAL
43// error on Linux, fall back to using accept.
44
45// Accept4Func is used to hook the accept4 call.
46var Accept4Func func(int, int) (int, syscall.Sockaddr, error) = syscall.Accept4
pollDesc.waitRead 方法主要负责检测当前这个 pollDesc 的上层 netFD 对应的 fd 是否有『期待的』I/O 事件发生，如果有就直接返回，否则就 park 住当前的 goroutine 并持续等待直至对应的 fd 上发生可读/可写或者其他『期待的』I/O 事件为止，然后它就会返回到外层的 for 循环，让 goroutine 继续执行逻辑。

poll.FD.Accept() 返回之后，会构造一个对应这个新 socket 的 netFD，然后调用 init() 方法完成初始化，这个 init 过程和前面 net.Listen() 是一样的，调用链：netFD.init() --> poll.FD.Init() --> poll.pollDesc.init()，最终又会走到这里：

 1var serverInit sync.Once
 2
 3func (pd *pollDesc) init(fd *FD) error {
 4	serverInit.Do(runtime_pollServerInit)
 5	ctx, errno := runtime_pollOpen(uintptr(fd.Sysfd))
 6	if errno != 0 {
 7		if ctx != 0 {
 8			runtime_pollUnblock(ctx)
 9			runtime_pollClose(ctx)
10		}
11		return syscall.Errno(errno)
12	}
13	pd.runtimeCtx = ctx
14	return nil
15}
然后把这个 socket fd 注册到 listener 的 epoll 实例的事件队列中去，等待 I/O 事件。

Conn.Read/Conn.Write
我们先来看看 Conn.Read 方法是如何实现的，原理其实和 Listener.Accept 是一样的，具体调用链还是首先调用 conn 的 netFD.Read ，然后内部再调用 poll.FD.Read ，最后使用 Linux 的系统调用 read: syscall.Read 完成数据读取：

 1// Implementation of the Conn interface.
 2
 3// Read implements the Conn Read method.
 4func (c *conn) Read(b []byte) (int, error) {
 5	if !c.ok() {
 6		return 0, syscall.EINVAL
 7	}
 8	n, err := c.fd.Read(b)
 9	if err != nil && err != io.EOF {
10		err = &OpError{Op: "read", Net: c.fd.net, Source: c.fd.laddr, Addr: c.fd.raddr, Err: err}
11	}
12	return n, err
13}
14
15func (fd *netFD) Read(p []byte) (n int, err error) {
16	n, err = fd.pfd.Read(p)
17	runtime.KeepAlive(fd)
18	return n, wrapSyscallError("read", err)
19}
20
21// Read implements io.Reader.
22func (fd *FD) Read(p []byte) (int, error) {
23	if err := fd.readLock(); err != nil {
24		return 0, err
25	}
26	defer fd.readUnlock()
27	if len(p) == 0 {
28		// If the caller wanted a zero byte read, return immediately
29		// without trying (but after acquiring the readLock).
30		// Otherwise syscall.Read returns 0, nil which looks like
31		// io.EOF.
32		// TODO(bradfitz): make it wait for readability? (Issue 15735)
33		return 0, nil
34	}
35	if err := fd.pd.prepareRead(fd.isFile); err != nil {
36		return 0, err
37	}
38	if fd.IsStream && len(p) > maxRW {
39		p = p[:maxRW]
40	}
41	for {
42		// 尝试从该 socket 读取数据，因为 socket 在被 listener accept 的时候设置成
43		// 了非阻塞 I/O，所以这里同样也是直接返回，不管有没有可读的数据
44		n, err := syscall.Read(fd.Sysfd, p)
45		if err != nil {
46			n = 0
47			// err == syscall.EAGAIN 表示当前没有期待的 I/O 事件发生，也就是 socket 不可读
48			if err == syscall.EAGAIN && fd.pd.pollable() {
49				// 如果当前没有发生期待的 I/O 事件，那么 waitRead
50				// 会通过 park goroutine 让逻辑 block 在这里
51				if err = fd.pd.waitRead(fd.isFile); err == nil {
52					continue
53				}
54			}
55
56			// On MacOS we can see EINTR here if the user
57			// pressed ^Z.  See issue #22838.
58			if runtime.GOOS == "darwin" && err == syscall.EINTR {
59				continue
60			}
61		}
62		err = fd.eofError(n, err)
63		return n, err
64	}
65}
conn.Write 和 conn.Read 的原理是一致的，它也是通过类似 pollDesc.waitRead 的 pollDesc.waitWrite 来 park 住 goroutine 直至期待的 I/O 事件发生才返回恢复执行。

pollDesc.waitRead/pollDesc.waitWrite
pollDesc.waitRead 内部调用了 poll.runtime_pollWait --> runtime.poll_runtime_pollWait 来达成无 I/O 事件时 park 住 goroutine 的目的：

  1//go:linkname poll_runtime_pollWait internal/poll.runtime_pollWait
  2func poll_runtime_pollWait(pd *pollDesc, mode int) int {
  3	err := netpollcheckerr(pd, int32(mode))
  4	if err != pollNoError {
  5		return err
  6	}
  7	// As for now only Solaris, illumos, and AIX use level-triggered IO.
  8	if GOOS == "solaris" || GOOS == "illumos" || GOOS == "aix" {
  9		netpollarm(pd, mode)
 10	}
 11	// 进入 netpollblock 并且判断是否有期待的 I/O 事件发生，
 12	// 这里的 for 循环是为了一直等到 io ready
 13	for !netpollblock(pd, int32(mode), false) {
 14		err = netpollcheckerr(pd, int32(mode))
 15		if err != 0 {
 16			return err
 17		}
 18		// Can happen if timeout has fired and unblocked us,
 19		// but before we had a chance to run, timeout has been reset.
 20		// Pretend it has not happened and retry.
 21	}
 22	return 0
 23}
 24
 25// returns true if IO is ready, or false if timedout or closed
 26// waitio - wait only for completed IO, ignore errors
 27func netpollblock(pd *pollDesc, mode int32, waitio bool) bool {
 28	// gpp 保存的是 goroutine 的数据结构 g，这里会根据 mode 的值决定是 rg 还是 wg，
 29	// 前面提到过，rg 和 wg 是用来保存等待 I/O 就绪的 gorouine 的，后面调用 gopark 之后，
 30	// 会把当前的 goroutine 的抽象数据结构 g 存入 gpp 这个指针，也就是 rg 或者 wg
 31	gpp := &pd.rg
 32	if mode == 'w' {
 33		gpp = &pd.wg
 34	}
 35
 36	// set the gpp semaphore to WAIT
 37	// 这个 for 循环是为了等待 io ready 或者 io wait
 38	for {
 39		old := *gpp
 40		// gpp == pdReady 表示此时已有期待的 I/O 事件发生，
 41		// 可以直接返回 unblock 当前 goroutine 并执行响应的 I/O 操作
 42		if old == pdReady {
 43			*gpp = 0
 44			return true
 45		}
 46		if old != 0 {
 47			throw("runtime: double wait")
 48		}
 49		// 如果没有期待的 I/O 事件发生，则通过原子操作把 gpp 的值置为 pdWait 并退出 for 循环
 50		if atomic.Casuintptr(gpp, 0, pdWait) {
 51			break
 52		}
 53	}
 54
 55	// need to recheck error states after setting gpp to WAIT
 56	// this is necessary because runtime_pollUnblock/runtime_pollSetDeadline/deadlineimpl
 57	// do the opposite: store to closing/rd/wd, membarrier, load of rg/wg
 58
 59	// waitio 此时是 false，netpollcheckerr 方法会检查当前 pollDesc 对应的 fd 是否是正常的，
 60	// 通常来说  netpollcheckerr(pd, mode) == 0 是成立的，所以这里会执行 gopark
 61	// 把当前 goroutine 给 park 住，直至对应的 fd 上发生可读/可写或者其他『期待的』I/O 事件为止，
 62	// 然后 unpark 返回，在 gopark 内部会把当前 goroutine 的抽象数据结构 g 存入
 63	// gpp(pollDesc.rg/pollDesc.wg) 指针里，以便在后面的 netpoll 函数取出 pollDesc 之后，
 64	// 把 g 添加到链表里返回，接着重新调度 goroutine
 65	if waitio || netpollcheckerr(pd, mode) == 0 {
 66		// 注册 netpollblockcommit 回调给 gopark，在 gopark 内部会执行它，保存当前 goroutine 到 gpp
 67		gopark(netpollblockcommit, unsafe.Pointer(gpp), waitReasonIOWait, traceEvGoBlockNet, 5)
 68	}
 69	// be careful to not lose concurrent READY notification
 70	old := atomic.Xchguintptr(gpp, 0)
 71	if old > pdWait {
 72		throw("runtime: corrupted polldesc")
 73	}
 74	return old == pdReady
 75}
 76
 77// gopark 会停住当前的 goroutine 并且调用传递进来的回调函数 unlockf，从上面的源码我们可以知道这个函数是
 78// netpollblockcommit
 79func gopark(unlockf func(*g, unsafe.Pointer) bool, lock unsafe.Pointer, reason waitReason, traceEv byte, traceskip int) {
 80	if reason != waitReasonSleep {
 81		checkTimeouts() // timeouts may expire while two goroutines keep the scheduler busy
 82	}
 83	mp := acquirem()
 84	gp := mp.curg
 85	status := readgstatus(gp)
 86	if status != _Grunning && status != _Gscanrunning {
 87		throw("gopark: bad g status")
 88	}
 89	mp.waitlock = lock
 90	mp.waitunlockf = unlockf
 91	gp.waitreason = reason
 92	mp.waittraceev = traceEv
 93	mp.waittraceskip = traceskip
 94	releasem(mp)
 95	// can't do anything that might move the G between Ms here.
 96	// gopark 最终会调用 park_m，在这个函数内部会调用 unlockf，也就是 netpollblockcommit，
 97	// 然后会把当前的 goroutine，也就是 g 数据结构保存到 pollDesc 的 rg 或者 wg 指针里
 98	mcall(park_m)
 99}
100
101// park continuation on g0.
102func park_m(gp *g) {
103	_g_ := getg()
104
105	if trace.enabled {
106		traceGoPark(_g_.m.waittraceev, _g_.m.waittraceskip)
107	}
108
109	casgstatus(gp, _Grunning, _Gwaiting)
110	dropg()
111
112	if fn := _g_.m.waitunlockf; fn != nil {
113		// 调用 netpollblockcommit，把当前的 goroutine，
114		// 也就是 g 数据结构保存到 pollDesc 的 rg 或者 wg 指针里
115		ok := fn(gp, _g_.m.waitlock)
116		_g_.m.waitunlockf = nil
117		_g_.m.waitlock = nil
118		if !ok {
119			if trace.enabled {
120				traceGoUnpark(gp, 2)
121			}
122			casgstatus(gp, _Gwaiting, _Grunnable)
123			execute(gp, true) // Schedule it back, never returns.
124		}
125	}
126	schedule()
127}
128
129// netpollblockcommit 在 gopark 函数里被调用
130func netpollblockcommit(gp *g, gpp unsafe.Pointer) bool {
131	// 通过原子操作把当前 goroutine 抽象的数据结构 g，也就是这里的参数 gp 存入 gpp 指针，
132	// 此时 gpp 的值是 pollDesc 的 rg 或者 wg 指针
133	r := atomic.Casuintptr((*uintptr)(gpp), pdWait, uintptr(unsafe.Pointer(gp)))
134	if r {
135		// Bump the count of goroutines waiting for the poller.
136		// The scheduler uses this to decide whether to block
137		// waiting for the poller if there is nothing else to do.
138		atomic.Xadd(&netpollWaiters, 1)
139	}
140	return r
141}
pollDesc.waitWrite 的内部实现原理和 pollDesc.waitRead 是一样的，都是基于 poll.runtime_pollWait --> runtime.poll_runtime_pollWait，这里就不再赘述。

netpoll
前面已经从源码的层面分析完了 netpoll 是如何通过 park goroutine 从而达到阻塞 Accept/Read/Write 的效果，而通过调用 gopark，goroutine 会被放置在某个等待队列中，这里是放到了 epoll 的 "interest list" 里，底层数据结构是由红黑树实现的 eventpoll.rbr，此时 G 的状态由 _Grunning为_Gwaitting ，因此 G 必须被手动唤醒(通过 goready )，否则会丢失任务，应用层阻塞通常使用这种方式。

所以我们现在可以来从整体的层面来概括 Go 的网络业务 goroutine 是如何被规划调度的了：



首先，client 连接 server 的时候，listener 通过 accept 调用接收新 connection，每一个新 connection 都启动一个 goroutine 处理，accept 调用会把该 connection 的 fd 连带所在的 goroutine 上下文信息封装注册到 epoll 的监听列表里去，当 goroutine 调用 conn.Read 或者 conn.Write 等需要阻塞等待的函数时，会被 gopark 给封存起来并使之休眠，让 P 去执行本地调度队列里的下一个可执行的 goroutine，往后 Go scheduler 会在循环调度的 runtime.schedule() 函数以及 sysmon 监控线程中调用 runtime.netpoll 以获取可运行的 goroutine 列表并通过调用 injectglist 把剩下的 g 放入全局调度队列或者当前 P 本地调度队列去重新执行。

那么当 I/O 事件发生之后，netpoller 是通过什么方式唤醒那些在 I/O wait 的 goroutine 的？答案是通过 runtime.netpoll。

runtime.netpoll 的核心逻辑是：

根据调用方的入参 delay，设置对应的调用 epollwait 的 timeout 值；
调用 epollwait 等待发生了可读/可写事件的 fd；
循环 epollwait 返回的事件列表，处理对应的事件类型， 组装可运行的 goroutine 链表并返回。
  1// netpoll checks for ready network connections.
  2// Returns list of goroutines that become runnable.
  3// delay < 0: blocks indefinitely
  4// delay == 0: does not block, just polls
  5// delay > 0: block for up to that many nanoseconds
  6func netpoll(delay int64) gList {
  7	if epfd == -1 {
  8		return gList{}
  9	}
 10
 11	// 根据特定的规则把 delay 值转换为 epollwait 的 timeout 值
 12	var waitms int32
 13	if delay < 0 {
 14		waitms = -1
 15	} else if delay == 0 {
 16		waitms = 0
 17	} else if delay < 1e6 {
 18		waitms = 1
 19	} else if delay < 1e15 {
 20		waitms = int32(delay / 1e6)
 21	} else {
 22		// An arbitrary cap on how long to wait for a timer.
 23		// 1e9 ms == ~11.5 days.
 24		waitms = 1e9
 25	}
 26	var events [128]epollevent
 27retry:
 28	// 超时等待就绪的 fd 读写事件
 29	n := epollwait(epfd, &events[0], int32(len(events)), waitms)
 30	if n < 0 {
 31		if n != -_EINTR {
 32			println("runtime: epollwait on fd", epfd, "failed with", -n)
 33			throw("runtime: netpoll failed")
 34		}
 35		// If a timed sleep was interrupted, just return to
 36		// recalculate how long we should sleep now.
 37		if waitms > 0 {
 38			return gList{}
 39		}
 40		goto retry
 41	}
 42
 43	// toRun 是一个 g 的链表，存储要恢复的 goroutines，最后返回给调用方
 44	var toRun gList
 45	for i := int32(0); i < n; i++ {
 46		ev := &events[i]
 47		if ev.events == 0 {
 48			continue
 49		}
 50
 51		// Go scheduler 在调用 findrunnable() 寻找 goroutine 去执行的时候，
 52		// 在调用 netpoll 之时会检查当前是否有其他线程同步阻塞在 netpoll，
 53		// 若是，则调用 netpollBreak 来唤醒那个线程，避免它长时间阻塞
 54		if *(**uintptr)(unsafe.Pointer(&ev.data)) == &netpollBreakRd {
 55			if ev.events != _EPOLLIN {
 56				println("runtime: netpoll: break fd ready for", ev.events)
 57				throw("runtime: netpoll: break fd ready for something unexpected")
 58			}
 59			if delay != 0 {
 60				// netpollBreak could be picked up by a
 61				// nonblocking poll. Only read the byte
 62				// if blocking.
 63				var tmp [16]byte
 64				read(int32(netpollBreakRd), noescape(unsafe.Pointer(&tmp[0])), int32(len(tmp)))
 65				atomic.Store(&netpollWakeSig, 0)
 66			}
 67			continue
 68		}
 69
 70		// 判断发生的事件类型，读类型或者写类型等，然后给 mode 复制相应的值，
 71		// mode 用来决定从 pollDesc 里的 rg 还是 wg 里取出 goroutine
 72		var mode int32
 73		if ev.events&(_EPOLLIN|_EPOLLRDHUP|_EPOLLHUP|_EPOLLERR) != 0 {
 74			mode += 'r'
 75		}
 76		if ev.events&(_EPOLLOUT|_EPOLLHUP|_EPOLLERR) != 0 {
 77			mode += 'w'
 78		}
 79		if mode != 0 {
 80			// 取出保存在 epollevent 里的 pollDesc
 81			pd := *(**pollDesc)(unsafe.Pointer(&ev.data))
 82			pd.everr = false
 83			if ev.events == _EPOLLERR {
 84				pd.everr = true
 85			}
 86			// 调用 netpollready，传入就绪 fd 的 pollDesc，
 87			// 把 fd 对应的 goroutine 添加到链表 toRun 中
 88			netpollready(&toRun, pd, mode)
 89		}
 90	}
 91	return toRun
 92}
 93
 94// netpollready 调用 netpollunblock 返回就绪 fd 对应的 goroutine 的抽象数据结构 g
 95func netpollready(toRun *gList, pd *pollDesc, mode int32) {
 96	var rg, wg *g
 97	if mode == 'r' || mode == 'r'+'w' {
 98		rg = netpollunblock(pd, 'r', true)
 99	}
100	if mode == 'w' || mode == 'r'+'w' {
101		wg = netpollunblock(pd, 'w', true)
102	}
103	if rg != nil {
104		toRun.push(rg)
105	}
106	if wg != nil {
107		toRun.push(wg)
108	}
109}
110
111// netpollunblock 会依据传入的 mode 决定从 pollDesc 的 rg 或者 wg 取出当时 gopark 之时存入的
112// goroutine 抽象数据结构 g 并返回
113func netpollunblock(pd *pollDesc, mode int32, ioready bool) *g {
114	// mode == 'r' 代表当时 gopark 是为了等待读事件，而 mode == 'w' 则代表是等待写事件
115	gpp := &pd.rg
116	if mode == 'w' {
117		gpp = &pd.wg
118	}
119
120	for {
121		// 取出 gpp 存储的 g
122		old := *gpp
123		if old == pdReady {
124			return nil
125		}
126		if old == 0 && !ioready {
127			// Only set READY for ioready. runtime_pollWait
128			// will check for timeout/cancel before waiting.
129			return nil
130		}
131		var new uintptr
132		if ioready {
133			new = pdReady
134		}
135		// 重置 pollDesc 的 rg 或者 wg
136		if atomic.Casuintptr(gpp, old, new) {
137			// 如果该 goroutine 还是必须等待，则返回 nil
138			if old == pdWait {
139				old = 0
140			}
141			// 通过万能指针还原成 g 并返回
142			return (*g)(unsafe.Pointer(old))
143		}
144	}
145}
146
147// netpollBreak 往通信管道里写入信号去唤醒 epollwait
148func netpollBreak() {
149	// 通过 CAS 避免重复的唤醒信号被写入管道，
150	// 从而减少系统调用并节省一些系统资源
151	if atomic.Cas(&netpollWakeSig, 0, 1) {
152		for {
153			var b byte
154			n := write(netpollBreakWr, unsafe.Pointer(&b), 1)
155			if n == 1 {
156				break
157			}
158			if n == -_EINTR {
159				continue
160			}
161			if n == -_EAGAIN {
162				return
163			}
164			println("runtime: netpollBreak write failed with", -n)
165			throw("runtime: netpollBreak write failed")
166		}
167	}
168}
Go 在多种场景下都可能会调用 netpoll 检查文件描述符状态，netpoll 里会调用 epoll_wait 从 epoll 的 eventpoll.rdllist 就绪双向链表返回，从而得到 I/O 就绪的 socket fd 列表，并根据取出最初调用 epoll_ctl 时保存的上下文信息，恢复 g。所以执行完netpoll 之后，会返回一个就绪 fd 列表对应的 goroutine 链表，接下来将就绪的 goroutine 通过调用 injectglist 加入到全局调度队列或者 P 的本地调度队列中，启动 M 绑定 P 去执行。

具体调用 netpoll 的地方，首先在 Go runtime scheduler 循环调度 goroutines 之时就有可能会调用 netpoll 获取到已就绪的 fd 对应的 goroutine 来调度执行。

首先 Go scheduler 的核心方法 runtime.schedule() 里会调用一个叫 runtime.findrunable() 的方法获取可运行的 goroutine 来执行，而在 runtime.findrunable() 方法里就调用了 runtime.netpoll 获取已就绪的 fd 列表对应的 goroutine 列表：

 1// One round of scheduler: find a runnable goroutine and execute it.
 2// Never returns.
 3func schedule() {
 4	...
 5  
 6  if gp == nil {
 7		gp, inheritTime = findrunnable() // blocks until work is available
 8	}
 9  
10	...
11}
12
13// Finds a runnable goroutine to execute.
14// Tries to steal from other P's, get g from global queue, poll network.
15func findrunnable() (gp *g, inheritTime bool) {
16  ...
17  
18  // Poll network.
19	if netpollinited() && (atomic.Load(&netpollWaiters) > 0 || pollUntil != 0) && atomic.Xchg64(&sched.lastpoll, 0) != 0 {
20		atomic.Store64(&sched.pollUntil, uint64(pollUntil))
21		if _g_.m.p != 0 {
22			throw("findrunnable: netpoll with p")
23		}
24		if _g_.m.spinning {
25			throw("findrunnable: netpoll with spinning")
26		}
27		if faketime != 0 {
28			// When using fake time, just poll.
29			delta = 0
30		}
31		list := netpoll(delta) // 同步阻塞调用 netpoll，直至有可用的 goroutine
32		atomic.Store64(&sched.pollUntil, 0)
33		atomic.Store64(&sched.lastpoll, uint64(nanotime()))
34		if faketime != 0 && list.empty() {
35			// Using fake time and nothing is ready; stop M.
36			// When all M's stop, checkdead will call timejump.
37			stopm()
38			goto top
39		}
40		lock(&sched.lock)
41		_p_ = pidleget() // 查找是否有空闲的 P 可以来就绪的 goroutine
42		unlock(&sched.lock)
43		if _p_ == nil {
44			injectglist(&list) // 如果当前没有空闲的 P，则把就绪的 goroutine 放入全局调度队列等待被执行
45		} else {
46			// 如果当前有空闲的 P，则 pop 出一个 g，返回给调度器去执行，
47			// 并通过调用 injectglist 把剩下的 g 放入全局调度队列或者当前 P 本地调度队列
48			acquirep(_p_)
49			if !list.empty() {
50				gp := list.pop()
51				injectglist(&list)
52				casgstatus(gp, _Gwaiting, _Grunnable)
53				if trace.enabled {
54					traceGoUnpark(gp, 0)
55				}
56				return gp, false
57			}
58			if wasSpinning {
59				_g_.m.spinning = true
60				atomic.Xadd(&sched.nmspinning, 1)
61			}
62			goto top
63		}
64	} else if pollUntil != 0 && netpollinited() {
65		pollerPollUntil := int64(atomic.Load64(&sched.pollUntil))
66		if pollerPollUntil == 0 || pollerPollUntil > pollUntil {
67			netpollBreak()
68		}
69	}
70	stopm()
71	goto top
72}
另外， sysmon 监控线程会在循环过程中检查距离上一次 runtime.netpoll 被调用是否超过了 10ms，若是则会去调用它拿到可运行的 goroutine 列表并通过调用 injectglist 把 g 列表放入全局调度队列或者当前 P 本地调度队列等待被执行：

 1// Always runs without a P, so write barriers are not allowed.
 2//
 3//go:nowritebarrierrec
 4func sysmon() {
 5		...
 6  
 7		// poll network if not polled for more than 10ms
 8		lastpoll := int64(atomic.Load64(&sched.lastpoll))
 9		if netpollinited() && lastpoll != 0 && lastpoll+10*1000*1000 < now {
10			atomic.Cas64(&sched.lastpoll, uint64(lastpoll), uint64(now))
11			list := netpoll(0) // non-blocking - returns list of goroutines
12			if !list.empty() {
13				// Need to decrement number of idle locked M's
14				// (pretending that one more is running) before injectglist.
15				// Otherwise it can lead to the following situation:
16				// injectglist grabs all P's but before it starts M's to run the P's,
17				// another M returns from syscall, finishes running its G,
18				// observes that there is no work to do and no other running M's
19				// and reports deadlock.
20				incidlelocked(-1)
21				injectglist(&list)
22				incidlelocked(1)
23			}
24		}
25  
26  ...
27}
Go runtime 在程序启动的时候会创建一个独立的 M 作为监控线程，叫 sysmon ，这个线程为系统级的 daemon 线程，无需 P 即可运行， sysmon 每 20us~10ms 运行一次。 sysmon 中以轮询的方式执行以下操作（如上面的代码所示）：

以非阻塞的方式调用 runtime.netpoll ，从中找出能从网络 I/O 中唤醒的 g 列表，并通过调用 injectglist 把 g 列表放入全局调度队列或者当前 P 本地调度队列等待被执行，调度触发时，有可能从这个全局 runnable 调度队列获取 g。然后再循环调用 startm ，直到所有 P 都不处于 _Pidle 状态。
调用 retake ，抢占长时间处于 _Psyscall 状态的 P。
综上，Go 借助于 epoll/kqueue/iocp 和 runtime scheduler 等的帮助，设计出了自己的 I/O 多路复用 netpoller，成功地让 Listener.Accept / conn.Read / conn.Write 等方法从开发者的角度看来是同步模式。

Go netpoller 的价值
通过前面对源码的分析，我们现在知道 Go netpoller 依托于 runtime scheduler，为开发者提供了一种强大的同步网络编程模式；然而，Go netpoller 存在的意义却远不止于此，Go netpoller I/O 多路复用搭配 Non-blocking I/O 而打造出来的这个原生网络模型，它最大的价值是把网络 I/O 的控制权牢牢掌握在 Go 自己的 runtime 里，关于这一点我们需要从 Go 的 runtime scheduler 说起，Go 的 G-P-M 调度模型如下：



G 在运行过程中如果被阻塞在某个 system call 操作上，那么不光 G 会阻塞，执行该 G 的 M 也会解绑 P(实质是被 sysmon 抢走了)，与 G 一起进入 sleep 状态。如果此时有 idle 的 M，则 P 与其绑定继续执行其他 G；如果没有 idle M，但仍然有其他 G 要去执行，那么就会创建一个新的 M。当阻塞在 system call 上的 G 完成 syscall 调用后，G 会去尝试获取一个可用的 P，如果没有可用的 P，那么 G 会被标记为 _Grunnable 并把它放入全局的 runqueue 中等待调度，之前的那个 sleep 的 M 将再次进入 sleep。

现在清楚为什么 netpoll 为什么一定要使用非阻塞 I/O 了吧？就是为了避免让操作网络 I/O 的 goroutine 陷入到系统调用从而进入内核态，因为一旦进入内核态，整个程序的控制权就会发生转移(到内核)，不再属于用户进程了，那么也就无法借助于 Go 强大的 runtime scheduler 来调度业务程序的并发了；而有了 netpoll 之后，借助于非阻塞 I/O ，G 就再也不会因为系统调用的读写而 (长时间) 陷入内核态，当 G 被阻塞在某个 network I/O 操作上时，实际上它不是因为陷入内核态被阻塞住了，而是被 Go runtime 调用 gopark 给 park 住了，此时 G 会被放置到某个 wait queue 中，而 M 会尝试运行下一个 _Grunnable 的 G，如果此时没有 _Grunnable 的 G 供 M 运行，那么 M 将解绑 P，并进入 sleep 状态。当 I/O available，在 epoll 的 eventpoll.rdr 中等待的 G 会被放到 eventpoll.rdllist 链表里并通过 netpoll 中的 epoll_wait 系统调用返回放置到全局调度队列或者 P 的本地调度队列，标记为 _Grunnable ，等待 P 绑定 M 恢复执行。

Goroutine 的调度
这一小节主要是讲处理网络 I/O 的 goroutines 阻塞之后，Go scheduler 具体是如何像前面几个章节所说的那样，避免让操作网络 I/O 的 goroutine 陷入到系统调用从而进入内核态的，而是封存 goroutine 然后让出 CPU 的使用权从而令 P 可以去调度本地调度队列里的下一个 goroutine 的。

温馨提示：这一小节属于延伸阅读，涉及到的知识点更偏系统底层，需要有一定的汇编语言基础才能通读，另外，这一节对 Go scheduler 的讲解仅仅涉及核心的一部分，不会把整个调度器都讲一遍（事实上如果真要解析 Go scheduler 的话恐怕重开一篇几万字的文章才能基本讲清楚。。。），所以也要求读者对 Go 的并发调度器有足够的了解，因此这一节可能会稍显深奥。当然这一节也可选择不读，因为通过前面的整个解析，我相信读者应该已经能够基本掌握 Go netpoller 处理网络 I/O 的核心细节了，以及能从宏观层面了解 netpoller 对业务 goroutines 的基本调度了。而这一节主要是通过对 goroutines 调度细节的剖析，能够加深读者对整个 Go netpoller 的彻底理解，接上前面几个章节，形成一个完整的闭环。如果对调度的底层细节没兴趣的话这也可以直接跳过这一节，对理解 Go netpoller 的基本原理影响不大，不过还是建议有条件的读者可以看看。

从源码可知，Go scheduler 的调度 goroutine 过程中所调用的核心函数链如下：

1runtime.schedule --> runtime.execute --> runtime.gogo --> goroutine code --> runtime.goexit --> runtime.goexit1 --> runtime.mcall --> runtime.goexit0 --> runtime.schedule
Go scheduler 会不断循环调用 runtime.schedule() 去调度 goroutines，而每个 goroutine 执行完成并退出之后，会再次调用 runtime.schedule()，使得调度器回到调度循环去执行其他的 goroutine，不断循环，永不停歇。

当我们使用 go 关键字启动一个新 goroutine 时，最终会调用 runtime.newproc --> runtime.newproc1，来得到 g，runtime.newproc1 会先从 P 的 gfree 缓存链表中查找可用的 g，若缓存未生效，则会新创建 g 给当前的业务函数，最后这个 g 会被传给 runtime.gogo 去真正执行。

这里首先需要了解一个 gobuf 的结构体，它用来保存 goroutine 的调度信息，是 runtime.gogo 的入参：

 1// gobuf 存储 goroutine 调度上下文信息的结构体
 2type gobuf struct {
 3	// The offsets of sp, pc, and g are known to (hard-coded in) libmach.
 4	//
 5	// ctxt is unusual with respect to GC: it may be a
 6	// heap-allocated funcval, so GC needs to track it, but it
 7	// needs to be set and cleared from assembly, where it's
 8	// difficult to have write barriers. However, ctxt is really a
 9	// saved, live register, and we only ever exchange it between
10	// the real register and the gobuf. Hence, we treat it as a
11	// root during stack scanning, which means assembly that saves
12	// and restores it doesn't need write barriers. It's still
13	// typed as a pointer so that any other writes from Go get
14	// write barriers.
15	sp   uintptr  // Stack Pointer 栈指针
16	pc   uintptr  // Program Counter 程序计数器
17	g    guintptr // 持有当前 gobuf 的 goroutine
18	ctxt unsafe.Pointer
19	ret  sys.Uintreg
20	lr   uintptr
21	bp   uintptr // for GOEXPERIMENT=framepointer
22}
执行 runtime.execute()，进而调用 runtime.gogo：

 1func execute(gp *g, inheritTime bool) {
 2	_g_ := getg()
 3
 4	// Assign gp.m before entering _Grunning so running Gs have an
 5	// M.
 6	_g_.m.curg = gp
 7	gp.m = _g_.m
 8	casgstatus(gp, _Grunnable, _Grunning)
 9	gp.waitsince = 0
10	gp.preempt = false
11	gp.stackguard0 = gp.stack.lo + _StackGuard
12	if !inheritTime {
13		_g_.m.p.ptr().schedtick++
14	}
15
16	// Check whether the profiler needs to be turned on or off.
17	hz := sched.profilehz
18	if _g_.m.profilehz != hz {
19		setThreadCPUProfiler(hz)
20	}
21
22	if trace.enabled {
23		// GoSysExit has to happen when we have a P, but before GoStart.
24		// So we emit it here.
25		if gp.syscallsp != 0 && gp.sysblocktraced {
26			traceGoSysExit(gp.sysexitticks)
27		}
28		traceGoStart()
29	}
30	// gp.sched 就是 gobuf
31	gogo(&gp.sched)
32}
这里还需要了解一个概念：g0，Go G-P-M 调度模型中，g 代表 goroutine，而实际上一共有三种 g：

执行用户代码的 g；
执行调度器代码的 g，也即是 g0；
执行 runtime.main 初始化工作的 main goroutine；
第一种 g 就是使用 go 关键字启动的 goroutine，也是我们接触最多的一类 g；第三种 g 是调度器启动之后用来执行的一系列初始化工作的，包括但不限于启动 sysmon 监控线程、内存初始化和启动 GC 等等工作；第二种 g 叫 g0，用来执行调度器代码，g0 在底层和其他 g 是一样的数据结构，但是性质上有很大的区别，首先 g0 的栈大小是固定的，比如在 Linux 或者其他 Unix-like 的系统上一般是固定 8MB，不能动态伸缩，而普通的 g 初始栈大小是 2KB，可按需扩展，g0 其实就是线程栈，我们知道每个线程被创建出来之时都需要操作系统为之分配一个初始固定的线程栈，就是前面说的 8MB 大小的栈，g0 栈就代表了这个线程栈，因此每一个 m 都需要绑定一个 g0 来执行调度器代码，然后跳转到执行用户代码的地方。

runtime.gogo 是真正去执行 goroutine 代码的函数，这个函数由汇编实现，为什么需要用汇编？因为 gogo 的工作是完成线程 M 上的堆栈切换：从系统堆栈 g0 切换成 goroutine gp，也就是 CPU 使用权和堆栈的切换，这种切换本质上是对 CPU 的 PC、SP 等寄存器和堆栈指针的更新，而这一类精度的底层操作别说是 Go，就算是最贴近底层的 C 也无法做到，这种程度的操作已超出所有高级语言的范畴，因此只能借助于汇编来实现。

runtime.gogo 在不同的 CPU 架构平台上的实现各不相同，但是核心原理殊途同归，我们这里选用 amd64 架构的汇编实现来分析，我会在关键的地方加上解释：

 1// func gogo(buf *gobuf)
 2// restore state from Gobuf; longjmp
 3TEXT runtime·gogo(SB), NOSPLIT, $16-8
 4	// 将第一个 FP 伪寄存器所指向的 gobuf 的第一个参数存入 BX 寄存器, 
 5	// gobuf 的一个参数即是 SP 指针
 6	MOVQ	buf+0(FP), BX
 7	MOVQ	gobuf_g(BX), DX  // 将 gp.sched.g 保存到 DX 寄存器
 8	MOVQ	0(DX), CX		// make sure g != nil
 9
10	// 将 tls (thread local storage) 保存到 CX 寄存器，然后把 gp.sched.g 放到 tls[0]，
11	// 这样以后调用 getg() 之时就可以通过 TLS 直接获取到当前 goroutine 的 g 结构体实例，
12	// 进而可以得到 g 所在的 m 和 p，TLS 里一开始存储的是系统堆栈 g0 的地址
13	get_tls(CX)
14	MOVQ	DX, g(CX)
15
16	// 下面的指令则是对函数栈的 BP/SP 寄存器(指针)的存取，
17	// 最后进入到指定的代码区域，执行函数栈帧
18	MOVQ	gobuf_sp(BX), SP	// restore SP
19	MOVQ	gobuf_ret(BX), AX
20	MOVQ	gobuf_ctxt(BX), DX
21	MOVQ	gobuf_bp(BX), BP
22
23	// 这里是在清空 gp.sched，因为前面已经把 gobuf 里的字段值都存入了寄存器，
24	// 所以 gp.sched 就可以提前清空了，不需要等到后面 GC 来回收，减轻 GC 的负担
25	MOVQ	$0, gobuf_sp(BX)	// clear to help garbage collector
26	MOVQ	$0, gobuf_ret(BX)
27	MOVQ	$0, gobuf_ctxt(BX)
28	MOVQ	$0, gobuf_bp(BX)
29
30	// 把 gp.sched.pc 值放入 BX 寄存器
31	// PC 指针指向 gogo 退出时需要执行的函数地址
32	MOVQ	gobuf_pc(BX), BX
33	// 用 BX 寄存器里的值去修改 CPU 的 IP 寄存器，
34	// 这样就可以根据 CS:IP 寄存器的段地址+偏移量跳转到 BX 寄存器里的地址，也就是 gp.sched.pc
35	JMP	BX
runtime.gogo 函数接收 gp.sched 这个 gobuf 结构体实例，其中保存了函数栈寄存器 SP/PC/BP，如果熟悉操作系统原理的话可以知道这些寄存器是 CPU 进行函数调用和返回时切换对应的函数栈帧所需的寄存器，而 goroutine 的执行和函数调用的原理是一致的，也是 CPU 寄存器的切换过程，所以这里的几个寄存器当前存的就是 G 的函数执行栈，当 goroutine 在处理网络 I/O 之时，如果恰好处于 I/O 就绪的状态的话，则正常完成 runtime.gogo，并在最后跳转到特定的地址，那么这个地址是哪里呢？

我们知道 CPU 执行函数的时候需要知道函数在内存里的代码段地址和偏移量，然后才能去取来函数栈执行，而典型的提供代码段地址和偏移量的寄存器就是 CS 和 IP 寄存器，而 JMP BX 指令则是用 BX 寄存器去更新 IP 寄存器，而 BX 寄存器里的值是 gp.sched.pc，那么这个 PC 指针究竟是指向哪里呢？让我们来看另一处源码。

众所周知，启动一个新的 goroutine 是通过 go 关键字来完成的，而 go compiler 会在编译期间利用 cmd/compile/internal/gc.state.stmt 和 cmd/compile/internal/gc.state.call 这两个函数将 go 关键字翻译成 runtime.newproc 函数调用，而 runtime.newproc 接收了函数指针和其大小之后，会获取 goroutine 和调用处的程序计数器，接着再调用 runtime.newproc1：

 1// Create a new g in state _Grunnable, starting at fn, with narg bytes
 2// of arguments starting at argp. callerpc is the address of the go
 3// statement that created this. The caller is responsible for adding
 4// the new g to the scheduler.
 5//
 6// This must run on the system stack because it's the continuation of
 7// newproc, which cannot split the stack.
 8//
 9//go:systemstack
10func newproc1(fn *funcval, argp unsafe.Pointer, narg int32, callergp *g, callerpc uintptr) *g {
11  ...
12  
13  memclrNoHeapPointers(unsafe.Pointer(&newg.sched), unsafe.Sizeof(newg.sched))
14	newg.sched.sp = sp
15	newg.stktopsp = sp
16	// 把 goexit 函数地址存入 gobuf 的 PC 指针里
17	newg.sched.pc = funcPC(goexit) + sys.PCQuantum // +PCQuantum so that previous instruction is in same function
18	newg.sched.g = guintptr(unsafe.Pointer(newg))
19	gostartcallfn(&newg.sched, fn)
20	newg.gopc = callerpc
21	newg.ancestors = saveAncestors(callergp)
22	newg.startpc = fn.fn
23	if _g_.m.curg != nil {
24		newg.labels = _g_.m.curg.labels
25	}
26	if isSystemGoroutine(newg, false) {
27		atomic.Xadd(&sched.ngsys, +1)
28	}
29	casgstatus(newg, _Gdead, _Grunnable)
30  
31  ...
32}
这里可以看到，newg.sched.pc 被设置了 runtime.goexit 的函数地址，newg 就是后面 runtime.gogo 执行的 goroutine，因此 runtime.gogo 最后的汇编指令 JMP BX是跳转到了 runtime.goexit，让我们来继续看看这个函数做了什么：

 1// The top-most function running on a goroutine
 2// returns to goexit+PCQuantum. Defined as ABIInternal
 3// so as to make it identifiable to traceback (this
 4// function it used as a sentinel; traceback wants to
 5// see the func PC, not a wrapper PC).
 6TEXT runtime·goexit<ABIInternal>(SB),NOSPLIT,$0-0
 7	BYTE	$0x90	// NOP
 8	CALL	runtime·goexit1(SB)	// does not return
 9	// traceback from goexit1 must hit code range of goexit
10	BYTE	$0x90	// NOP
这个函数也是汇编实现的，但是非常简单，就是直接调用 runtime·goexit1：

 1// Finishes execution of the current goroutine.
 2func goexit1() {
 3	if raceenabled {
 4		racegoend()
 5	}
 6	if trace.enabled {
 7		traceGoEnd()
 8	}
 9	mcall(goexit0)
10}
调用 runtime.mcall函数：

 1// func mcall(fn func(*g))
 2// Switch to m->g0's stack, call fn(g).
 3// Fn must never return. It should gogo(&g->sched)
 4// to keep running g.
 5
 6// 切换回 g0 的系统堆栈，执行 fn(g)
 7TEXT runtime·mcall(SB), NOSPLIT, $0-8
 8	// 取入参 funcval 对象的指针存入 DI 寄存器，此时 fn.fn 是 goexit0 的地址
 9	MOVQ	fn+0(FP), DI
10
11	get_tls(CX)
12	MOVQ	g(CX), AX	// save state in g->sched
13	MOVQ	0(SP), BX	// caller's PC
14	MOVQ	BX, (g_sched+gobuf_pc)(AX)
15	LEAQ	fn+0(FP), BX	// caller's SP
16	MOVQ	BX, (g_sched+gobuf_sp)(AX)
17	MOVQ	AX, (g_sched+gobuf_g)(AX)
18	MOVQ	BP, (g_sched+gobuf_bp)(AX)
19
20	// switch to m->g0 & its stack, call fn
21	MOVQ	g(CX), BX
22	MOVQ	g_m(BX), BX
23
24	// 把 g0 的栈指针存入 SI 寄存器，后面需要用到
25	MOVQ	m_g0(BX), SI
26	CMPQ	SI, AX	// if g == m->g0 call badmcall
27	JNE	3(PC)
28	MOVQ	$runtime·badmcall(SB), AX
29	JMP	AX
30
31	// 这两个指令是把 g0 地址存入到 TLS 里，
32	// 然后从 SI 寄存器取出 g0 的栈指针，
33	// 替换掉 SP 寄存器里存的当前 g 的栈指针
34	MOVQ	SI, g(CX)	// g = m->g0
35	MOVQ	(g_sched+gobuf_sp)(SI), SP	// sp = m->g0->sched.sp
36
37	PUSHQ	AX
38	MOVQ	DI, DX
39
40	// 入口处的第一个指令已经把 funcval 实例对象的指针存入了 DI 寄存器，
41	// 0(DI) 表示取出 DI 的第一个成员，即 goexit0 函数地址，再存入 DI
42	MOVQ	0(DI), DI
43	CALL	DI // 调用 DI 寄存器里的地址，即 goexit0
44	POPQ	AX
45	MOVQ	$runtime·badmcall2(SB), AX
46	JMP	AX
47	RET
可以看到 runtime.mcall 函数的主要逻辑是从当前 goroutine 切换回 g0 的系统堆栈，然后调用 fn(g)，此处的 g 即是当前运行的 goroutine，这个方法会保存当前运行的 G 的 PC/SP 到 g->sched 里，以便该 G 可以在以后被重新恢复执行，因为也涉及到寄存器和堆栈指针的操作，所以也需要使用汇编实现，该函数最后会在 g0 系统堆栈下执行 runtime.goexit0:

 1func goexit0(gp *g) {
 2	_g_ := getg()
 3
 4	casgstatus(gp, _Grunning, _Gdead)
 5	if isSystemGoroutine(gp, false) {
 6		atomic.Xadd(&sched.ngsys, -1)
 7	}
 8	gp.m = nil
 9	locked := gp.lockedm != 0
10	gp.lockedm = 0
11	_g_.m.lockedg = 0
12	gp.preemptStop = false
13	gp.paniconfault = false
14	gp._defer = nil // should be true already but just in case.
15	gp._panic = nil // non-nil for Goexit during panic. points at stack-allocated data.
16	gp.writebuf = nil
17	gp.waitreason = 0
18	gp.param = nil
19	gp.labels = nil
20	gp.timer = nil
21
22	if gcBlackenEnabled != 0 && gp.gcAssistBytes > 0 {
23		// Flush assist credit to the global pool. This gives
24		// better information to pacing if the application is
25		// rapidly creating an exiting goroutines.
26		scanCredit := int64(gcController.assistWorkPerByte * float64(gp.gcAssistBytes))
27		atomic.Xaddint64(&gcController.bgScanCredit, scanCredit)
28		gp.gcAssistBytes = 0
29	}
30
31	dropg()
32
33	if GOARCH == "wasm" { // no threads yet on wasm
34		gfput(_g_.m.p.ptr(), gp)
35		schedule() // never returns
36	}
37
38	if _g_.m.lockedInt != 0 {
39		print("invalid m->lockedInt = ", _g_.m.lockedInt, "\n")
40		throw("internal lockOSThread error")
41	}
42	gfput(_g_.m.p.ptr(), gp)
43	if locked {
44		// The goroutine may have locked this thread because
45		// it put it in an unusual kernel state. Kill it
46		// rather than returning it to the thread pool.
47
48		// Return to mstart, which will release the P and exit
49		// the thread.
50		if GOOS != "plan9" { // See golang.org/issue/22227.
51			gogo(&_g_.m.g0.sched)
52		} else {
53			// Clear lockedExt on plan9 since we may end up re-using
54			// this thread.
55			_g_.m.lockedExt = 0
56		}
57	}
58	schedule()
59}
runtime.goexit0 的主要工作是就是

利用 CAS 操作把 g 的状态从 _Grunning 更新为 _Gdead；
对 g 做一些清理操作，把一些字段值置空；
调用 runtime.dropg 解绑 g 和 m；
把 g 放入 p 存储 g 的 gfree 链表作为缓存，后续如果需要启动新的 goroutine 则可以直接从链表里取而不用重新初始化分配内存。
最后，调用 runtime.schedule() 再次进入调度循环去调度新的 goroutines，永不停歇。
另一方面，如果 goroutine 处于 I/O 不可用状态，我们前面已经分析过 netpoller 利用非阻塞 I/O + I/O 多路复用避免了陷入系统调用，所以此时会调用 runtime.gopark 并把 goroutine 暂时封存在用户态空间，并休眠当前的 goroutine，因此不会阻塞 runtime.gogo 的汇编执行，而是通过 runtime.mcall 调用 runtime.park_m：

 1func gopark(unlockf func(*g, unsafe.Pointer) bool, lock unsafe.Pointer, reason waitReason, traceEv byte, traceskip int) {
 2	if reason != waitReasonSleep {
 3		checkTimeouts() // timeouts may expire while two goroutines keep the scheduler busy
 4	}
 5	mp := acquirem()
 6	gp := mp.curg
 7	status := readgstatus(gp)
 8	if status != _Grunning && status != _Gscanrunning {
 9		throw("gopark: bad g status")
10	}
11	mp.waitlock = lock
12	mp.waitunlockf = unlockf
13	gp.waitreason = reason
14	mp.waittraceev = traceEv
15	mp.waittraceskip = traceskip
16	releasem(mp)
17	// can't do anything that might move the G between Ms here.
18	mcall(park_m)
19}
20
21func park_m(gp *g) {
22	_g_ := getg()
23
24	if trace.enabled {
25		traceGoPark(_g_.m.waittraceev, _g_.m.waittraceskip)
26	}
27
28	casgstatus(gp, _Grunning, _Gwaiting)
29	dropg()
30
31	if fn := _g_.m.waitunlockf; fn != nil {
32		ok := fn(gp, _g_.m.waitlock)
33		_g_.m.waitunlockf = nil
34		_g_.m.waitlock = nil
35		if !ok {
36			if trace.enabled {
37				traceGoUnpark(gp, 2)
38			}
39			casgstatus(gp, _Gwaiting, _Grunnable)
40			execute(gp, true) // Schedule it back, never returns.
41		}
42	}
43	schedule()
44}
runtime.mcall 方法我们在前面已经介绍过，它主要的工作就是是从当前 goroutine 切换回 g0 的系统堆栈，然后调用 fn(g)，而此时 runtime.mcall 调用执行的是 runtime.park_m，这个方法里会利用 CAS 把当前运行的 goroutine -- gp 的状态 从 _Grunning 切换到 _Gwaiting，表明该 goroutine 已进入到等待唤醒状态，此时封存和休眠 G 的操作就完成了，只需等待就绪之后被重新唤醒执行即可。最后调用 runtime.schedule() 再次进入调度循环，去执行下一个 goroutine，充分利用 CPU。

至此，我们完成了对 Go netpoller 原理剖析的整个闭环。

Go netpoller 的问题
Go netpoller 的设计不可谓不精巧、性能也不可谓不高，配合 goroutine 开发网络应用的时候就一个字：爽。因此 Go 的网络编程模式是及其简洁高效的，然而，没有任何一种设计和架构是完美的， goroutine-per-connection 这种模式虽然简单高效，但是在某些极端的场景下也会暴露出问题：goroutine 虽然非常轻量，它的自定义栈内存初始值仅为 2KB，后面按需扩容；海量连接的业务场景下， goroutine-per-connection ，此时 goroutine 数量以及消耗的资源就会呈线性趋势暴涨，虽然 Go scheduler 内部做了 g 的缓存链表，可以一定程度上缓解高频创建销毁 goroutine 的压力，但是对于瞬时性暴涨的长连接场景就无能为力了，大量的 goroutines 会被不断创建出来，从而对 Go runtime scheduler 造成极大的调度压力和侵占系统资源，然后资源被侵占又反过来影响 Go scheduler 的调度，进而导致性能下降。

Reactor 网络模型
目前 Linux 平台上主流的高性能网络库/框架中，大都采用 Reactor 模式，比如 netty、libevent、libev、ACE，POE(Perl)、Twisted(Python)等。

Reactor 模式本质上指的是使用 I/O 多路复用(I/O multiplexing) + 非阻塞 I/O(non-blocking I/O) 的模式。

通常设置一个主线程负责做 event-loop 事件循环和 I/O 读写，通过 select/poll/epoll_wait 等系统调用监听 I/O 事件，业务逻辑提交给其他工作线程去做。而所谓『非阻塞 I/O』的核心思想是指避免阻塞在 read() 或者 write() 或者其他的 I/O 系统调用上，这样可以最大限度的复用 event-loop 线程，让一个线程能服务于多个 sockets。在 Reactor 模式中，I/O 线程只能阻塞在 I/O multiplexing 函数上（select/poll/epoll_wait）。

Reactor 模式的基本工作流程如下：

Server 端完成在 bind&listen 之后，将 listenfd 注册到 epollfd 中，最后进入 event-loop 事件循环。循环过程中会调用 select/poll/epoll_wait 阻塞等待，若有在 listenfd 上的新连接事件则解除阻塞返回，并调用 socket.accept 接收新连接 connfd，并将 connfd 加入到 epollfd 的 I/O 复用（监听）队列。
当 connfd 上发生可读/可写事件也会解除 select/poll/epoll_wait 的阻塞等待，然后进行 I/O 读写操作，这里读写 I/O 都是非阻塞 I/O，这样才不会阻塞 event-loop 的下一个循环。然而，这样容易割裂业务逻辑，不易理解和维护。
调用 read 读取数据之后进行解码并放入队列中，等待工作线程处理。
工作线程处理完数据之后，返回到 event-loop 线程，由这个线程负责调用 write 把数据写回 client。
accept 连接以及 conn 上的读写操作若是在主线程完成，则要求是非阻塞 I/O，因为 Reactor 模式一条最重要的原则就是：I/O 操作不能阻塞 event-loop 事件循环。实际上 event loop 可能也可以是多线程的，只是一个线程里只有一个 select/poll/epoll_wait。

上面提到了 Go netpoller 在某些场景下可能因为创建太多的 goroutine 而过多地消耗系统资源，而在现实世界的网络业务中，服务器持有的海量连接中在极短的时间窗口内只有极少数是 active 而大多数则是 idle，就像这样（非真实数据，仅仅是为了比喻）：



那么为每一个连接指派一个 goroutine 就显得太过奢侈了，而 Reactor 模式这种利用 I/O 多路复用进而只需要使用少量线程即可管理海量连接的设计就可以在这样网络业务中大显身手了：

MultiReactors.png

在绝大部分应用场景下，我推荐大家还是遵循 Go 的 best practices，使用原生的 Go 网络库来构建自己的网络应用。然而，在某些极度追求性能、压榨系统资源以及技术栈必须是原生 Go （不考虑 C/C++ 写中间层而 Go 写业务层）的业务场景下，我们可以考虑自己构建 Reactor 网络模型。

gnet
gnet 是一个基于事件驱动的高性能和轻量级网络框架。它直接使用 epoll 和 kqueue 系统调用而非标准 Go 网络包：net 来构建网络应用，它的工作原理类似两个开源的网络库：netty 和 libuv，这也使得gnet 达到了一个远超 Go net 的性能表现。

gnet 设计开发的初衷不是为了取代 Go 的标准网络库：net，而是为了创造出一个类似于 Redis、Haproxy 能高效处理网络包的 Go 语言网络服务器框架。

gnet 的卖点在于它是一个高性能、轻量级、非阻塞的纯 Go 实现的传输层（TCP/UDP/Unix Domain Socket）网络框架，开发者可以使用 gnet 来实现自己的应用层网络协议(HTTP、RPC、Redis、WebSocket 等等)，从而构建出自己的应用层网络应用：比如在 gnet 上实现 HTTP 协议就可以创建出一个 HTTP 服务器 或者 Web 开发框架，实现 Redis 协议就可以创建出自己的 Redis 服务器等等。

gnet，在某些极端的网络业务场景，比如海量连接、高频短连接、网络小包等等场景，gnet 在性能和资源占用上都远超 Go 原生的 net 包（基于 netpoller）。

gnet 已经实现了 Multi-Reactors 和 Multi-Reactors + Goroutine Pool 两种网络模型，也得益于这些网络模型，使得 gnet 成为一个高性能和低损耗的 Go 网络框架
<!-- more -->

https://strikefreedom.top/go-netpoll-io-multiplexing-reactor

https://strikefreedom.top/reading-list