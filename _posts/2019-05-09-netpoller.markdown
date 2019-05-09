---
title: netpoller
layout: post
category: golang
author: 夏泽民
---
Go中网络交互采用多路复用的技术，具体到各个平台，即Kqueue、Epoll、Select、Poll等，下面以Linux下的Epoll实现为例进行分析。

netpoller代码分析
所谓的netpoller，其实是Golang实利用了OS提供的非阻塞IO访问模式，并配合epll/kqueue等IO事件监控机制；为了弥合OS的异步机制与Golang接口的差异，而在runtime上做的一层封装。以此来实现网络IO优化。

实际的实现(epoll/kqueue)必须定义以下函数：

func netpollinit() // 初始化轮询器
func netpollopen(fd uintptr, pd *pollDesc) int32 // 为fd和pd启动边缘触发通知
<!-- more -->
当一个goroutine进行io阻塞时，会去被放到等待队列。这里面就关键的就是建立起文件描述符和goroutine之间的关联。 pollDesc结构体就是完成这个任务的。代码参见src/runtime/netpoll.go。

type pollDesc struct { // Poller对象
    link *pollDesc // 链表
    lock mutex // 保护下面字段
    fd uintptr // fd是底层网络io文件描述符，整个生命期内，不能改变值
    closing bool
    seq uintptr // protect from stale(过时) timers and ready notifications
    rg uintptr // reader goroutine addr
    rt timer
    rd int64
    wg uintptr // writer goroutine addr
    wt timer
    wd int64
    user int32 // user-set cookie用户自定义数据
}

type pollCache struct { // 全局Poller链表
    lock mutex // 保护Poller链表
    first *pollDesc
}

func poll_runtime_pollServerInit() // 调用netpollinit()
func poll_runtime_pollOpen() // 调用netpollopen()
func poll_runtime_pollClose() // 调用netpollclose()
func poll_runtime_pollReset(pd, mode) // 先check(netpollcheckerr(pd, mode))是否有err发生，没有的话重置pd对应字段
func poll_runtime_pollWait(pd, mode) // 先chekerr，再调用netpollblock(pd, mode, false)
func poll_runtime_pollWaitCanceled(pd, mode) // windows下专用
func poll_runtime_pollSetDeadline(pd, deadline, mode)
//1. 重置定时器，并seq++
//2. 设置超时函数netpollDeadline(或者netpollReadDeadline、netpollWriteDeadline)
//3. 如果已经过期，调用netpollunblock和netpollgoready
func poll_runtime_pollUnblock(pd) // netpollUnblock、netpollgoready

/*------------------部分实现------------------*/
func netpollcheckerr(pd, mode) // 检查是否超时或正在关闭
func netpollblockcommit(gp *g, gpp unsafe.Pointer)
func netpollready(gpp *guintptr, pd, mode) // 调用netpollunblock，更新g的schedlink
func netpollgoready(gp *g, traceskip) // 更新统计数据，调用goready --- 通知调度器协程g从parked变为ready
func netpollblock(pd, mode, waitio) // Set rg/wg = pdWait，调用gopark挂起pd对应的g。
func netpollunblock(pd, mode, ioready)
func netpoll(Write/Read)Deadline(arg, seq)
pollCache是pollDesc链表入口，加锁保护链表安全。

pollDesc中，rg、wg有些特殊，它可能有如下3种状态:

pdReady == 1:  网络io就绪通知，goroutine消费完后应置为nil
pdWait == 2: goroutine等待被挂起，后续可能有3种情况:
goroutine被调度器挂起，置为goroutine地址
收到io通知，置为pdReady
超时或者被关闭，置为nil
Goroutine地址: 被挂起的goroutine的地址，当io就绪时、或者超时、被关闭时，此goroutine将被唤醒，同时将状态改为pdReady或者nil。
另外，由于wg、rg是goroutine的地址，因此当GC发生后，如果goroutine被回收(在heap区)，代码就崩溃了(指针无效)。所以，进行网络IO的goroutine不能在heap区分配内存。

lock锁对象保护了pollOpen, pollSetDeadline, pollUnblock和deadlineimpl操作。而这些操作又完全包含了对seq, rt, tw变量。fd在PollDesc整个生命过程中都是一个常量。处理pollReset, pollWait, pollWaitCanceled和runtime.netpollready(IO就绪通知)不需要用到锁，所以closing, rg, rd, wg和wd的所有操作都是一个无锁的操作。

多路复用三部曲
初始化PollServer
初始化在下面注册fd监听时顺便处理了，调用runtime_pollServerInit()，并使用sync.Once()机制保证只会被初始化一次。全局使用同一个EpollServer(同一个Epfd)。

func poll_runtime_ServerInit() {
    netpollinit() // 具现化到Linux下，调用epoll_create
    ...
}
注册监听fd
所有Unix文件在初始化时，如果支持Poll，都会加入到PollServer的监听中。源码下搜索runtime_pollOpen，即见分晓。

/*****************internal/poll/fd_unix.go*******************/
type FD struct {
    // Lock sysfd and serialize access to Read and Write methods.
    fdmu fdMutex
    // System file descriptor. Immutable until Close.
    Sysfd int
    // I/O poller.
    pd pollDesc
    ...
}
func(fd *FD) Init(net string, pollable bool) error {
    ...
    err := fd.pd.init(fd) // 初始化pd
    ...
}
...
/*****************internal/poll/fd_poll_runtime.go*****************/
type pollDesc struct {
    runtimeCtx uintptr
}
func (pd *pollDesc) init(fd *FD) error {
    serverInit.Do(runtime_pollServerInit) // 初始化PollServer(sync.Once)
    ctx, errno := runtime_pollOpen(uintptr(fd.Sysfd))
    ...
    runtimeCtx = ctx
    return nil
}
...
/*****************runtime/netpoll.go*****************/
func poll_runtime_pollOpen(fd uintptr) (*epDesc, int32) {
    ...
    errno := netpollopen(fd, pd) // 具现化到Linux下，调用epoll_ctl
    ...
}
取消fd的监听与此流程类似，最终调用epoll_ctl.

定期Poll
结合上述实现，必然有处逻辑定期执行epoll_wait来检测fd状态。在代码中搜索下netpoll(，即可发现是在sysmon、startTheWorldWithSema、pollWork、findrunnable中调用的，以sysmon为例:

// runtime/proc.go
...
lastpoll := int64(atomic.Load64(&sched.lastpoll))
now := nanotime()
// 如果10ms内没有poll过，则poll。(1ms=1000000ns)
if lastpoll != 0 && lastpoll+10*1000*1000 < now {
    atomic.Cas64(&sched.lastpoll, uint64(lastpoll), uint64(now))
    gp := netpoll(false) // netpoll在Linux具现为epoll_wait
    if gp != nil {
       injectglist(gp) //把g放到sched中去执行，底层仍然是调用的之前在goroutine里面提到的startm函数。
	}
}
...
以读等待挂起为例
加入监听
golang中客户端与服务端进行通讯时，常用如下方法:

conn, err := net.Dial("tcp", "localhost:1208")
...
从net.Dial看进去，最终会调用net/net_posix.go中的socket函数，大致流程如下:

func socket(...) ... {
	/*
	1. 调用sysSocket创建原生socket
	2. 调用同名包下netFd()，初始化网络文件描述符netFd
	3. 调用fd.dial()，其中最终有调用runtime_pollOpen()加入监听列表
	*/
}
至此，文件描述符已经加入pollServer监听列表。

读等待
主要是挂起goroutine，并建立gorotine和fd之间的关联。

当从netFd读取数据时，调用system call，循环从fd.sysfd读取数据：

func (fd *FD) Read(p []byte) (int, error) {
    if err := fd.pd.prepareRead(fd.isFile); err != nil {
        return 0, err
    }
    if fd.IsStream && len(p) > maxRW {
        p = p[:maxRW]
    }
    for {
        n, err := syscall.Read(fd.Sysfd, p)
        if err != nil {
            n = 0
            if err == syscall.EAGAIN && fd.pd.pollable() {
                if err = fd.pd.waitRead(fd.isFile); err == nil {
                    continue
                }
            }
        }
        err = fd.eofError(n, err)
        return n, err
    }
}
读取的时候只处理EAGAIN类型的错误，其他错误一律返回给调用者，因为对于非阻塞的网络连接的文件描述符，如果错误是EAGAIN，说明Socket的缓冲区为空，未读取到任何数据，则调用fd.pd.WaitRead：

func (pd *pollDesc) waitRead(isFile bool) error {
    return pd.wait('r', isFile)
}

func (pd *pollDesc) wait(mode int, isFile bool) error {
    if pd.runtimeCtx == 0 {
        return errors.New("waiting for unsupported file type")
    }
    res := runtime_pollWait(pd.runtimeCtx, mode)
    return convertErr(res, isFile)
}
res是runtime_pollWait函数返回的结果，由conevertErr函数包装后返回：

func convertErr(res int, isFile bool) error {
    switch res {
    case 0:
        return nil
    case 1:
        return errClosing(isFile)
    case 2:
        return ErrTimeout
    }
    println("unreachable: ", res)
    panic("unreachable")
}
其中0表示io已经准备好了，1表示链接意见关闭，2表示io超时。再来看看pollWait的实现：

func poll_runtime_pollWait(pd *pollDesc, mode int) int {
    err := netpollcheckerr(pd, int32(mode))
    if err != 0 {
        return err
    }
    for !netpollblock(pd, int32(mode), false) {
        err = netpollcheckerr(pd, int32(mode))
        if err != 0 {
            return err
        }
    }
    return 0
}
调用netpollblock来判断IO是否准备好了：

func netpollblock(pd *pollDesc, mode int32, waitio bool) bool {
	gpp := &pd.rg
	if mode == 'w' {
		gpp = &pd.wg
	}
    for {
        old := *gpp
        if old == pdReady {
            *gpp = 0
            return true
        }
        if old != 0 {
            throw("runtime: double wait")
        }
        if atomic.Casuintptr(gpp, 0, pdWait) {
            break
        }
    }
    if waitio || netpollcheckerr(pd, mode) == 0 {
	    gopark(netpollblockcommit, unsafe.Pointer(gpp), "IO wait", traceEvGoBlockNet, 5)
    }
    old := atomic.Xchguintptr(gpp, 0)
    if old > pdWait {
        throw("runtime: corrupted polldesc")
    }
    return old == pdReady
}
返回true说明IO已经准备好，返回false说明IO操作已经超时或者已经关闭。否则当waitio为false, 且io不出现错误或者超时才会挂起当前goroutine。最后的gopark函数，就是将当前的goroutine(调用者)设置为waiting状态。

就绪唤醒
如上所述，go在多种场景下都会调用netpoll检查文件描述符状态。寻找到IO就绪的socket文件描述符，并找到这些socket文件描述符对应的轮询器中附带的信息，根据这些信息将之前等待这些socket文件描述符就绪的goroutine状态修改为Grunnable。执行完netpoll之后，会找到一个就绪的goroutine列表，接下来将就绪的goroutine加入到调度队列中，等待调度运行。

总结
总的来说，netpoller的最终的效果就是用户层阻塞，底层非阻塞。当goroutine读或写阻塞时会被放到等待队列，这个goroutine失去了运行权，但并不是真正的整个系统“阻塞”于系统调用。而通过后台的poller不停地poll，所有的文件描述符都被添加到了这个poller中的，当某个时刻一个文件描述符准备好了，poller就会唤醒之前因它而阻塞的goroutine，于是goroutine重新运行起来。

和使用Unix系统中的select或是poll方法不同地是，Golang的netpoller查询的是能被调度的goroutine而不是那些函数指针、包含了各种状态变量的struct等，这样你就不用管理这些状态，也不用重新检查函数指针等，这些都是你在传统Unix网络I/O需要操心的问题。

netpoll只是一种框架和一些接口，只有依赖这个框架和接口实现的netpoll实例，netpoll才能发挥它的功能。类似于kernel中的vfs，vfs本身并不会去做具体的文件操作，只是为不同的文件系统提供了一个框架。netpoll不依赖于网络协议栈，因此在内核网络及I/O子系统尚未可用时，也可以发送或接收数据包。当然netpoll能够处理的数据包类型也很有限，只有UDP和ARP数据包，并且只能是以太网报文。注意这里对UDP数据包的处理并不像四层的UDP协议那样复杂，并且netpoll可以发挥作用要依赖网络设备的支持。
  1、netpoll结构和netpoll_info结构
  netpoll结构用来描述接收和发送数据包的必要信息，每一个依赖netpoll的模块在使用这个框架前都必须实现并注册netpoll实例。
  netpoll结构定义如下：

[cpp] view plaincopy
struct netpoll {  
    struct net_device *dev;  
    char dev_name[IFNAMSIZ];  
   
    const char *name;  
   
    void (*rx_hook)(struct netpoll *, int, char *, int);  
   
    __be32 local_ip, remote_ip;  
    u16 local_port, remote_port;  
   
    u8 remote_mac[ETH_ALEN];  
};  
  dev成员存储的是绑定的网络设备实例，netpoll实例只能通过特定的网络设备接收和发送数据包。该设备在注册netpoll实例时设置。
  dev_name存储的是网络设备名，通过它调用dev_get_by_name()获取指定的网络设备实例，并保存在dev中。
  name是netpoll实例的名称。
  netpoll实例有两种：能接收数据包和不能接收数据包。如果要接收数据包的话，必须实现rx_hook接口。如果不接收数据包的话，则不用。
  local_ip和remote_ip分别存储的是远端和本地的IP，由netpoll实例指定。
  local_port和remote_port分别存储的是远端和本地的port。
  remote_mac存储的MAC地址。netpoll只支持以太网数据包，所以这里的MAC地址是以太网MAC地址。
  当支持netpoll时，网络设备的net_device实例必须实现npinfo成员，即网络设备的netpoll_info信息块，描述结构为netpoll_info，定义如下：
[cpp] view plaincopy
struct netpoll_info {  
   
    atomic_t refcnt;  
   
    int rx_flags;  
   
    spinlock_t rx_lock;  
   
    struct netpoll *rx_np; /* netpoll that registered an rx_hook */  
   
    struct sk_buff_head arp_tx; /* list of arp requests to reply to */  
   
    struct sk_buff_head txq;  
    struct delayed_work tx_work;  
};  
  refcnt是引用计数。每个netpoll_info实例被多个netpoll实例引用，每次引用时都对该成员加1.
  rx_flags是标识接收的特性，可取的值为NETPOLL_RX_ENABLED和NETPOLL_RX_DROP（尚未使用）。如果所属的netpoll实例允许接收数据包，则会设置为NETPOLL_RX_ENABLED，否则为0.
  rx_lock用来保证同一时刻只有一个CPU在进行相关的netpoll的输入操作。除此之外，在清理netpoll实例操作与netpoll的输入操作互斥，参见netpoll_cleanup().
  如果注册的netpoll实例可以接收数据包，则将实例存储在rx_np成员中，不过该成员在发送数据包时也会使用，参见arp_reply().
  arp_tx存储的是接收到的ARP报文。这里存储的ARP报文是在service_arp_queue()中处理的，而调用该函数的是netpoll_poll()，后面再讨论netpoll_poll()函数。
  如果netpoll没有能成功发送数据包或者设备繁忙，则将待输出报文缓存到txq队列中，重新调度tx_work工作队列，等待再次尝试发送。
  2、netpoll的输入
  netpoll_rx()函数是netpoll接收数据包的入口函数，在netif_rx()和netif_receive_skb()中都会调用到。如果该函数返回0，则表示当前数据包不是netpoll想要的，继续传递到上层协议栈继续处理；如果返回1，则表示由netpoll来处理，不再向上层传递。
  如果是ARP包，是否接收还要看静态变量trapped。trapped默认状态下是0，只有在poll_one_api()中调用网络设备的poll接口接收数据包前，才会加1，接收完后又会减1，如下所示：
[cpp] view plaincopy
static int poll_one_napi(struct netpoll_info *npinfo,  
             struct napi_struct *napi, int budget)  
{  
    int work;  
   
    /* net_rx_action's ->poll() invocations and our's are 
     * synchronized by this test which is only made while 
     * holding the napi->poll_lock. 
     */  
    if (!test_bit(NAPI_STATE_SCHED, &napi->state))  
        return budget;  
   
    npinfo->rx_flags |= NETPOLL_RX_DROP;  
    atomic_inc(&trapped);  
    set_bit(NAPI_STATE_NPSVC, &napi->state);  
   
    work = napi->poll(napi, budget);  
    trace_napi_poll(napi);  
   
    clear_bit(NAPI_STATE_NPSVC, &napi->state);  
    atomic_dec(&trapped);  
    npinfo->rx_flags &= ~NETPOLL_RX_DROP;  
   
    return budget - work;  
}  
  poll_one_api()是由poll_napi()调用的，如果当前CPU和接收数据包的CPU不是一个CPU，并且此时网卡被放置到轮询列表，即设置了NAPI_STATE_SCHED，才会去执行接收操作。所以netpoll在调度接收网卡的数据包过程中会trap数据包（trapped不为0），这种情况下ARP包会被接收。如果trapped为0，即不trap数据包，并且是ARP数据包，则会传递到上层协议栈。不过，在__netpoll_rx()中返回之前，trapped此时不为0，会丢弃ARP包。
  如果是UDP数据包，则主要是检查校验和和IP地址、端口号等信息，确定是否是netpoll想要的数据包，如果不是，则根据trapped决定是丢弃数据包还是传递到上层协议栈。如果是netpoll实例感兴趣的报文，则会调用其注册的rx_hook接口来接收，然后释放掉SKB包（注意，是在调用rx_hook之后立即释放）。
  还有一点需要注意，如果netpoll在调度接收网卡的数据包，即trapped不为0，这个过程中会直接释放掉所有不是netpoll想要的数据包（只是netpoll实例绑定的网卡上的数据包）。个人理解是，netpoll只有在发送数据包没有成功或者分配skb失败（尝试10次）时（都是在向外输出数据包的时候）才会调度网卡接收数据包，如果出现这种情况，则说明网卡非常繁忙，并且很多数据包没来得及处理，此时丢掉数据包也是合理的。
3、netpoll的输出
  如果接收到ARP报文，会调用arp_reply()来发送ARP响应；如果是UDP报文，则由netpoll实例处理，发送数据包调用的是netpoll_send_udp()。不过这两个接口最终都是在封装好要发送的数据包后，交给netpoll_send_skb()来发送，如下所示：
[cpp] view plaincopy
static void netpoll_send_skb(struct netpoll *np, struct sk_buff *skb)  
{  
    ......  
   
    /* don't get messages out of order, and no recursion */  
    if (skb_queue_len(&npinfo->txq) == 0 && !netpoll_owner_active(dev)) {  
        struct netdev_queue *txq;  
        unsigned long flags;  
   
        txq = netdev_get_tx_queue(dev, skb_get_queue_mapping(skb));  
   
        local_irq_save(flags);  
        /* try until next clock tick */  
        for (tries = jiffies_to_usecs(1)/USEC_PER_POLL;  
             tries > 0; --tries) {  
            if (__netif_tx_trylock(txq)) {  
                if (!netif_tx_queue_stopped(txq)) {  
                    status = ops->ndo_start_xmit(skb, dev);  
                    if (status == NETDEV_TX_OK)  
                        txq_trans_update(txq);  
                }  
                __netif_tx_unlock(txq);  
   
                if (status == NETDEV_TX_OK)  
                    break;  
   
            }  
   
            /* tickle device maybe there is some cleanup */  
            netpoll_poll(np);  
   
            udelay(USEC_PER_POLL);  
        }  
   
        WARN_ONCE(!irqs_disabled(),  
            "netpoll_send_skb(): %s enabled interrupts in poll (%pF)\n",  
            dev->name, ops->ndo_start_xmit);  
   
        local_irq_restore(flags);  
    }  
   
    if (status != NETDEV_TX_OK) {  
        skb_queue_tail(&npinfo->txq, skb);  
        schedule_delayed_work(&npinfo->tx_work,0);  
    }  
}  
  从上面的代码我们可以看到，只有在（skb_queue_len(&npinfo->txq) == 0 && !netpoll_owner_active(dev)）为真时才会尝试，发送数据包，否则直接缓存到txq队列中。
  如果npinfo->txq队列不为空，说明tx_work工作队列已经被调度执行，此时直接将数据包缓存到txq队列中，通过tx_work工作队列来输出。注意，这里调用schedule_delayed_work()的时候，延迟时间设置的是0，所以如果重新调度的话，tx_work工作队列会立即开始执行。
  如果npinfo->txq队列为空，是否将数据包直接缓存到txq队列，取决于netpoll_owner_active()的返回值。netpoll_owner_active()源码如下：
[cpp] view plaincopy
static int netpoll_owner_active(struct net_device *dev)  
{  
    struct napi_struct *napi;  
   
    list_for_each_entry(napi, &dev->napi_list, dev_list) {  
        if (napi->poll_owner == smp_processor_id())  
            return 1;  
    }  
    return 0;  
}  
  poll_owner是在接收数据包的软中断处理函数net_rx_action()中设置的，保存的是当前处理软中断的CPU的ID。如果netpoll实例绑定的网卡没有在接收数据包，也就是网卡没有放到设备轮询列表上，此时会直接返回0.如果此时网卡被放到轮询列表上，但是接收数据包的CPU不是当前的CPU，也会返回0。如果此时绑定的网卡正在接收数据包，并且是当前CPU，才会返回1，这时netpoll在发送SKB包时，会直接将数据包放到txq队列中，等待tx_work工作队列发送。  
  如果不是上述情况，netpoll_send_skb()会立即调用网络设备的ndo_start_xmit接口发送数据包。如果发送失败，则会尝试多次，直到下一次时钟节拍。如果仍然没有发送成功，则会将数据包缓存到txq队列中。
  在尝试重新发送的过程中，netpoll对调用netpoll_poll()接口来模拟网络设备接收到数据包的中断，然后借助其他CPU来接收数据包，源码如下：
[cpp] view plaincopy
void netpoll_poll(struct netpoll *np)  
{  
    struct net_device *dev = np->dev;  
    const struct net_device_ops *ops;  
   
    if (!dev || !netif_running(dev))  
        return;  
   
    ops = dev->netdev_ops;  
    if (!ops->ndo_poll_controller)  
        return;  
   
    /* Process pending work on NIC */  
    ops->ndo_poll_controller(dev);  
   
    poll_napi(dev);  
   
    /* 
      * 处理arp_tx队列中的ARP报文 
      */  
    service_arp_queue(dev->npinfo);  
   
    zap_completion_queue();  
}  
  模拟中断的接口是ndo_poll_controller，如果网卡不支持，则直接返回。模拟中断后，网卡设备会被放到轮询列表上，在poll_api()中会检查接收数据包的CPU和当前CPU是否是同一个CPU，如果不是，则会调用poll_one_napi()去使用网络设备的poll接口来接收数据包，否则直接返回，避免在UP上出现递归的情况。如果可以接收数据包，则trapped会加1，此时netpoll会trap数据包，该网卡上不是netpoll想要的数据包都会被直接丢掉，也只有在这段时间netpoll才可以接收ARP报文。所以我们看到，处理netpoll接收到的ARP包的接口，只在netpoll_poll()中调用，也只有在此时才有必要去处理接收到的ARP包。
  综上所述，netpoll_poll()会加速网卡对数据包的处理，这样下次发送数据包时就更容易成功。
4、netpoll应用
  netconsole是依赖netpoll实现的，可以将本机的dmesg系统信息，通过网络的方式输出到另一台主机上。这样就可以实现远程监控某些主机的dmesg信息，给开发人员调试内核提供了非常方便的途径。netconolse的使用方法参见内核文档netconsole.txt，里面介绍的非常详细 