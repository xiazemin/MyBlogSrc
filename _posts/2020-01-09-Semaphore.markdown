---
title: Mutex/Semaphore/Spinlock
layout: post
category: linux
author: 夏泽民
---
Mutex是一把钥匙，一个人拿了就可进入一个房间，出来的时候把钥匙交给队列的第一个。一般的用法是用于串行化对critical section代码的访问，保证这段代码不会被并行的运行。


 

Semaphore是一件可以容纳N人的房间，如果人不满就可以进去，如果人满了，就要等待有人出来。对于N=1的情况，称为binary semaphore。一般的用法是，用于限制对于某一资源的同时访问。
<!-- more -->
Binary semaphore与Mutex的差异：

在有的系统中Binary semaphore与Mutex是没有差异的。在有的系统上，主要的差异是mutex一定要由获得锁的进程来释放。而semaphore可以由其它进程释放（这时的semaphore实际就是个原子的变量，大家可以加或减），因此semaphore可以用于进程间同步。Semaphore的同步功能是所有系统都支持的，而Mutex能否由其他进程释放则未定，因此建议mutex只用于保护critical section。而semaphore则用于保护某变量，或者同步。

 

另一个概念是spin lock，这是一个内核态概念。spin lock与semaphore的主要区别是spin lock是busy waiting，而semaphore是sleep。对于可以sleep的进程来说，busy waiting当然没有意义。对于单CPU的系统，busy waiting当然更没意义（没有CPU可以释放锁）。因此，只有多CPU的内核态非进程空间，才会用到spin lock。Linux kernel的spin lock在非SMP的情况下，只是关irq，没有别的操作，用于确保该段程序的运行不会被打断。其实也就是类似mutex的作用，串行化对 critical section的访问。但是mutex不能保护中断的打断，也不能在中断处理程序中被调用。而spin lock也一般没有必要用于可以sleep的进程空间。
---------------------------------------------------------------------------------------------

内核同步措施

    为了避免并发，防止竞争。内核提供了一组同步方法来提供对共享数据的保护。我们的重点不是介绍这些方法的详细用法，而是强调为什么使用这些方法和它们之间的差别。
    Linux 使用的同步机制可以说从2.0到2.6以来不断发展完善。从最初的原子操作，到后来的信号量，从大内核锁到今天的自旋锁。这些同步机制的发展伴随 Linux从单处理器到对称多处理器的过度；伴随着从非抢占内核到抢占内核的过度。锁机制越来越有效，也越来越复杂。
    目前来说内核中原子操作多用来做计数使用，其它情况最常用的是两种锁以及它们的变种:一个是自旋锁，另一个是信号量。我们下面就来着重介绍一下这两种锁机制。


自旋锁

    自旋锁是专为防止多处理器并发而引入的一种锁，它在内核中大量应用于中断处理等部分(对于单处理器来说，防止中断处理中的并发可简单采用关闭中断的方式，不需要自旋锁)。
    自旋锁最多只能被一个内核任务持有，如果一个内核任务试图请求一个已被争用(已经被持有)的自旋锁，那么这个任务就会一直进行忙循环——旋转——等待锁重新可用。要是锁未被争用，请求它的内核任务便能立刻得到它并且继续进行。自旋锁可以在任何时刻防止多于一个的内核任务同时进入临界区，因此这种锁可有效地避免多处理器上并发运行的内核任务竞争共享资源。
    事实上，自旋锁的初衷就是：在短期间内进行轻量级的锁定。一个被争用的自旋锁使得请求它的线程在等待锁重新可用的期间进行自旋(特别浪费处理器时间)，所以自旋锁不应该被持有时间过长。如果需要长时间锁定的话, 最好使用信号量。
自旋锁的基本形式如下：
    spin_lock(&mr_lock);
    //临界区
    spin_unlock(&mr_lock);

    因为自旋锁在同一时刻只能被最多一个内核任务持有，所以一个时刻只有一个线程允许存在于临界区中。这点很好地满足了对称多处理机器需要的锁定服务。在单处理器上，自旋锁仅仅当作一个设置内核抢占的开关。如果内核抢占也不存在，那么自旋锁会在编译时被完全剔除出内核。
    简单的说，自旋锁在内核中主要用来防止多处理器中并发访问临界区，防止内核抢占造成的竞争。另外自旋锁不允许任务睡眠(持有自旋锁的任务睡眠会造成自死锁——因为睡眠有可能造成持有锁的内核任务被重新调度，而再次申请自己已持有的锁)，它能够在中断上下文中使用。
    死锁：假设有一个或多个内核任务和一个或多个资源，每个内核都在等待其中的一个资源，但所有的资源都已经被占用了。这便会发生所有内核任务都在相互等待，但它们永远不会释放已经占有的资源，于是任何内核任务都无法获得所需要的资源，无法继续运行，这便意味着死锁发生了。自死琐是说自己占有了某个资源，然后自己又申请自己已占有的资源，显然不可能再获得该资源，因此就自缚手脚了。


信号量
    Linux中的信号量是一种睡眠锁。如果有一个任务试图获得一个已被持有的信号量时，信号量会将其推入等待队列，然后让其睡眠。这时处理器获得自由去执行其它代码。当持有信号量的进程将信号量释放后，在等待队列中的一个任务将被唤醒，从而便可以获得这个信号量。
    信号量的睡眠特性，使得信号量适用于锁会被长时间持有的情况；只能在进程上下文中使用，因为中断上下文中是不能被调度的；另外当代码持有信号量时，不可以再持有自旋锁。

信号量基本使用形式为：
static DECLARE_MUTEX(mr_sem);//声明互斥信号量
if(down_interruptible(&mr_sem))
    //可被中断的睡眠，当信号来到，睡眠的任务被唤醒
    //临界区
up(&mr_sem);


信号量和自旋锁区别
    虽然听起来两者之间的使用条件复杂，其实在实际使用中信号量和自旋锁并不易混淆。注意以下原则:
    如果代码需要睡眠——这往往是发生在和用户空间同步时——使用信号量是唯一的选择。由于不受睡眠的限制，使用信号量通常来说更加简单一些。如果需要在自旋锁和信号量中作选择，应该取决于锁被持有的时间长短。理想情况是所有的锁都应该尽可能短的被持有，但是如果锁的持有时间较长的话，使用信号量是更好的选择。另外，信号量不同于自旋锁，它不会关闭内核抢占，所以持有信号量的代码可以被抢占。这意味者信号量不会对影响调度反应时间带来负面影响。


自旋锁对信号量

需求                     建议的加锁方法

低开销加锁               优先使用自旋锁
短期锁定                 优先使用自旋锁
长期加锁                 优先使用信号量
中断上下文中加锁          使用自旋锁
持有锁是需要睡眠、调度     使用信号量

---------------------------------------------------------------------------------------------

 

临界区（Critical Section）

    保证在某一时刻只有一个线程能访问数据的简便办法。在任意时刻只允许一个线程对共享资源进行访问。如果有多个线程试图同时访问临界区，那么在有一个线程进入后其他所有试图访问此临界区的线程将被挂起，并一直持续到进入临界区的线程离开。临界区在被释放后，其他线程可以继续抢占，并以此达到用原子方式操作共享资源的目的。

 

    在使用临界区时，一般不允许其运行时间过长，只要进入临界区的线程还没有离开，其他所有试图进入此临界区的线程都会被挂起而进入到等待状态，并会在一定程度上影响。程序的运行性能。尤其需要注意的是不要将等待用户输入或是其他一些外界干预的操作包含到临界区。如果进入了临界区却一直没有释放，同样也会引起其他线程的长时间等待。虽然临界区同步速度很快，但却只能用来同步本进程内的线程，而不可用来同步多个进程中的线程。


互斥量（Mutex）

   互斥（Mutex）是一种用途非常广泛的内核对象。能够保证多个线程对同一共享资源的互斥访问。同临界区有些类似，只有拥有互斥对象的线程才具有访问资源的权限，由于互斥对象只有一个，因此就决定了任何情况下此共享资源都不会同时被多个线程所访问。当前占据资源的线程在任务处理完后应将拥有的互斥对象交出，以便其他线程在获得后得以访问资源。与其他几种内核对象不同，互斥对象在操作系统中拥有特殊代码，并由操作系统来管理，操作系统甚至还允许其进行一些其他内核对象所不能进行的非常规操作。互斥量跟临界区很相似，只有拥有互斥对象的线程才具有访问资源的权限，由于互斥对象只有一个，因此就决定了任何情况下此共享资源都不会同时被多个线程所访问。当前占据资源的线程在任务处理完后应将拥有的互斥对象交出，以便其他线程在获得后得以访问资源。互斥量比临界区复杂。因为使用互斥不仅仅能够在同一应用程序不同线程中实现资源的安全共享，而且可以在不同应用程序的线程之间实现对资源的安全共享。

   　
信号量（Semaphores）

    信号量对象对线程的同步方式与前面几种方法不同，信号允许多个线程同时使用共享资源，这与操作系统中的PV操作相同。它指出了同时访问共享资源的线程最大数目。它允许多个线程在同一时刻访问同一资源，但是需要限制在同一时刻访问此资源的最大线程数目。在用CreateSemaphore（）创建信号量时即要同时指出允许的最大资源计数和当前可用资源计数。一般是将当前可用资源计数设置为最大资源计数，每增加一个线程对共享资源的访问，当前可用资源计数就会减1，只要当前可用资源计数是大于0的，就可以发出信号量信号。但是当前可用计数减小到0时则说明当前占用资源的线程数已经达到了所允许的最大数目，不能在允许其他线程的进入，此时的信号量信号将无法发出。线程在处理完共享资源后，应在离开的同时通过ReleaseSemaphore（）函数将当前可用资源计数加1。在任何时候当前可用资源计数决不可能大于最大资源计数。信号量是通过计数来对线程访问资源进行控制的，而实际上信号量确实也被称作Dijkstra计数器。


    PV操作及信号量的概念都是由荷兰科学家E.W.Dijkstra提出的。信号量S是一个整数，S大于等于零时代表可供并发进程使用的资源实体数，但S小于零时则表示正在等待使用共享资源的进程数。

    P操作申请资源：
    （1）S减1；
    （2）若S减1后仍大于等于零，则进程继续执行；
    （3）若S减1后小于零，则该进程被阻塞后进入与该信号相对应的队列中，然后转入进程调度。
　　
    V操作释放资源：
    （1）S加1；
    （2）若相加结果大于零，则进程继续执行；
    （3）若相加结果小于等于零，则从该信号的等待队列中唤醒一个等待进程，然后再返回原进程继续执行或转入进程调度。


　　信号量的使用特点使其更适用于对Socket（套接字）程序中线程的同步。例如，网络上的HTTP服务器要对同一时间内访问同一页面的用户数加以限制，这时可以为没一个用户对服务器的页面请求设置一个线程，而页面则是待保护的共享资源，通过使用信号量对线程的同步作用可以确保在任一时刻无论有多少用户对某一页面进行访问，只有不大于设定的最大用户数目的线程能够进行访问，而其他的访问企图则被挂起，只有在有用户退出对此页面的访问后才有可能进入。


总结：

    1．互斥量与临界区的作用非常相似，但互斥量是可以命名的，也就是说它可以跨越进程使用。所以创建互斥量需要的资源更多，所以如果只为了在进程内部是用的话使用临界区会带来速度上的优势并能够减少资源占用量。因为互斥量是跨进程的互斥量一旦被创建，就可以通过名字打开它。

    2．互斥量（Mutex），信号灯（Semaphore）都可以被跨越进程使用来进行同步数据操作，而其他的对象与数据同步操作无关，但对于进程和线程来讲，如果进程和线程在运行状态则为无信号状态，在退出后为有信号状态。

    3．通过互斥量可以指定资源被独占的方式使用，但如果有下面一种情况通过互斥量就无法处理，比如现在一位用户购买了一份三个并发访问许可的数据库系统，可以根据用户购买的访问许可数量来决定有多少个线程/进程能同时进行数据库操作，这时候如果利用互斥量就没有办法完成这个要求，信号灯对象可以说是一种资源计数器
    
   一、 以2.6.38以前的内核为例， 讲spinlock、 mutex 以及 semaphore
1. spinlock更原始，效率高，但讲究更多，不能随便用。
2. 个人觉得初级阶段不要去深挖mutex 以及 semaphore的不同，用法类似。在内核代码里面搜索，感觉 DEFINE_MUTEX + mutex_lock_xx + mutex_unlock 用的更多。
3. 在内核里面这三个符号发挥的作用就是： 自旋锁与互斥体。
semaphore:内核中的信号量通常用作mutex互斥体（信号量初值初始化为1，即binary semaphore的方式，就达到了互斥的效果）。
mutex：顾名思义， 互斥体。
所以在内核里面，mutex_lock()和down()的使用情景基本上相同。

//spinlock.h
/******
*API
*spin_lock
*spin_lock_bh
*spin_lock_irq
*spin_trylock
*spin_trylock_bh
*spin_trylock_irq
*spin_unlock
*spin_unlock_bh
*spin_unlock_irq
*spin_unlock_irqrestore
*spin_unlock_wait
******/

//semaphore.h
用 DECLARE_MUTEX 定义了一个count==1 的信号量（binary semaphore)。

#define DECLARE_MUTEX(name)    \
    struct semaphore name = __SEMAPHORE_INITIALIZER(name, 1)
    
struct semaphore {
    spinlock_t        lock;
    unsigned int        count;
    struct list_head    wait_list;
};    

#define __SEMAPHORE_INITIALIZER(name, n)                \
{                                    \
    .lock        = __SPIN_LOCK_UNLOCKED((name).lock),        \
    .count        = n,                        \
    .wait_list    = LIST_HEAD_INIT((name).wait_list),        \
}
    
#define init_MUTEX(sem)        sema_init(sem, 1)
#define init_MUTEX_LOCKED(sem)    sema_init(sem, 0)

/*****
*API：
*#define init_MUTEX(sem)        sema_init(sem, 1)
*#define init_MUTEX_LOCKED(sem)    sema_init(sem, 0)
*extern void down(struct semaphore *sem);
*extern int __must_check down_interruptible(struct semaphore *sem);
*extern int __must_check down_killable(struct semaphore *sem);
*extern int __must_check down_trylock(struct semaphore *sem);
*extern int __must_check down_timeout(struct semaphore *sem, long jiffies);
*extern void up(struct semaphore *sem);
****/

//mutex.h
#define DEFINE_MUTEX(mutexname) \
    struct mutex mutexname = __MUTEX_INITIALIZER(mutexname)
    
    #define __MUTEX_INITIALIZER(lockname) \
        { .count = ATOMIC_INIT(1) \
        , .wait_lock = __SPIN_LOCK_UNLOCKED(lockname.wait_lock) \
        , .wait_list = LIST_HEAD_INIT(lockname.wait_list) \
        __DEBUG_MUTEX_INITIALIZER(lockname) \
        __DEP_MAP_MUTEX_INITIALIZER(lockname) }

/*
 * Simple, straightforward mutexes with strict semantics:
 *
 * - only one task can hold the mutex at a time
 * - only the owner can unlock the mutex
 * - multiple unlocks are not permitted
 * - recursive locking is not permitted
 * - a mutex object must be initialized via the API
 * - a mutex object must not be initialized via memset or copying
 * - task may not exit with mutex held
 * - memory areas where held locks reside must not be freed
 * - held mutexes must not be reinitialized
 * - mutexes may not be used in hardware or software interrupt
 *   contexts such as tasklets and timers
 *
 * These semantics are fully enforced when DEBUG_MUTEXES is
 * enabled. Furthermore, besides enforcing the above rules, the mutex
 * debugging code also implements a number of additional features
 * that make lock debugging easier and faster:
 *
 * - uses symbolic names of mutexes, whenever they are printed in debug output
 * - point-of-acquire tracking, symbolic lookup of function names
 * - list of all locks held in the system, printout of them
 * - owner tracking
 * - detects self-recursing locks and prints out all relevant info
 * - detects multi-task circular deadlocks and prints out all affected
 *   locks and tasks (and only those tasks)
 */
struct mutex {
    /* 1: unlocked, 0: locked, negative: locked, possible waiters */
    atomic_t        count;
    spinlock_t        wait_lock;
    struct list_head    wait_list;
#if defined(CONFIG_DEBUG_MUTEXES) || defined(CONFIG_SMP)
    struct thread_info    *owner;
#endif
#ifdef CONFIG_DEBUG_MUTEXES
    const char         *name;
    void            *magic;
#endif
#ifdef CONFIG_DEBUG_LOCK_ALLOC
    struct lockdep_map    dep_map;
#endif
};

/*******
*API:
*extern void mutex_lock(struct mutex *lock);
*extern int __must_check mutex_lock_interruptible(struct mutex *lock);
*extern int __must_check mutex_lock_killable(struct mutex *lock);
*extern void mutex_unlock(struct mutex *lock);
********/

复制代码
EG1-1:
    spinlock_t rtc_lock;
    spin_lock_init(&rtc_lock);//每个驱动都会事先初始化，只需要这一次初始化
    
    spin_lock_irq(&rtc_lock);
    //临界区
    spin_unlock_irq(&rtc_lock);

EG1-2:
    unsigned long flags;
    static spinlock_t i2o_drivers_lock;
    spin_lock_init(&i2o_drivers_lock);//每个驱动都会事先初始化，只需要这一次初始化

    spin_lock_irqsave(&i2o_drivers_lock, flags);
    //临界区
    spin_unlock_irqrestore(&i2o_drivers_lock, flags);

EG2:
    static DECLARE_MUTEX(start_stop_sem);
    down(&start_stop_sem);
    //临界区
    up(&start_stop_sem);

EG3:
    static DEFINE_MUTEX(adutux_mutex);
    mutex_lock_interruptible(&adutux_mutex);
    //临界区
    mutex_unlock(&adutux_mutex);
复制代码
 



二、 2.6.38以后DECLARE_MUTEX替换成DEFINE_SEMAPHORE（命名改变）, DEFINE_MUTEX用法不变

复制代码
    static DEFINE_SEMAPHORE(msm_fb_pan_sem);// DECLARE_MUTEX
    down(&adb_probe_mutex);
    //临界区
    up(&adb_probe_mutex);
    
    static DEFINE_SEMAPHORE(bnx2x_prev_sem);
    down_interruptible(&bnx2x_prev_sem);
    //临界区
    up(&bnx2x_prev_sem); 
复制代码
Linux 2.6.36以后file_operations和DECLARE_MUTEX 的变化  http://blog.csdn.net/heanyu/article/details/6757917

在include/linux/semaphore.h 中将#define DECLARE_MUTEX(name)   改成了 #define DEFINE_SEMAPHORE(name)   【命名】

 　　　　 
三、自旋锁与信号量
1. 自旋锁
简单的说，自旋锁在内核中主要用来防止多处理器中并发访问临界区，防止内核抢占造成的竞争。【适用于多处理器】【自旋锁会影响内核调度】
另外自旋锁不允许任务睡眠(持有自旋锁的任务睡眠会造成自死锁——因为睡眠有可能造成持有锁的内核任务被重新调度，而再次申请自己已持有的锁)，它能够在中断上下文中使用。【不允许任务睡眠】

锁定一个自旋锁的函数有四个：
void spin_lock(spinlock_t *lock); //最基本得自旋锁函数，它不失效本地中断。
void spin_lock_irqsave(spinlock_t *lock, unsigned long flags);//在获得自旋锁之前禁用硬中断（只在本地处理器上），而先前的中断状态保存在flags中
void spin_lockirq(spinlock_t *lock);//在获得自旋锁之前禁用硬中断（只在本地处理器上），不保存中断状态
void spin_lock_bh(spinlock_t *lock);//在获得锁前禁用软中断，保持硬中断打开状态

2. 信号量
内核中的信号量通常用作mutex互斥体（信号量初值初始化为1就达到了互斥的效果）。

如果代码需要睡眠——这往往是发生在和用户空间同步时——使用信号量是唯一的选择。由于不受睡眠的限制，使用信号量通常来说更加简单一些。【信号量使用简单】
如果需要在自旋锁和信号量中作选择，应该取决于锁被持有的时间长短。理想情况是所有的锁都应该尽可能短的被持有，但是如果锁的持有时间较长的话，使用信号量是更好的选择。【如果锁占用的时间较长，信号量更好】
另外，信号量不同于自旋锁，它不会关闭内核抢占，所以持有信号量的代码可以被抢占。这意味者信号量不会对影响调度反应时间带来负面影响。【信号量不会影响内核调度】

3. 使用情景对比
=============================================
需求                     建议的加锁方法 
低开销加锁               优先使用自旋锁
短期锁定                 优先使用自旋锁
长期加锁                 优先使用信号量
中断上下文中加锁         使用自旋锁
持有锁是需要睡眠、调度    使用信号量

1 spin_lock 

      自旋锁的实现是为了保护一段短小的临界区操作代码，保证这个临界区的操作是原子的，从而避免并发的竞争冒险。在Linux内核中，自旋锁通常用于包含内核数据结构的操作，你可以看到在许多内核数据结构中都嵌入有spinlock，这些大部分就是用于保证它自身被操作的原子性，在操作这样的结构体时都经历这样的过程：上锁-操作-解锁。

      如果内核控制路径发现自旋锁“开着”（可以获取），就获取锁并继续自己的执行。相反，如果内核控制路径发现锁由运行在另一个CPU上的内核控制路径“锁着”，就在原地“旋转”，反复执行一条紧凑的循环检测指令，直到锁被释放。 自旋锁是循环检测“忙等”，即等待时内核无事可做（除了浪费时间），进程在CPU上保持运行，所以它保护的临界区必须小，且操作过程必须短。不过，自旋锁通常非常方便，因为很多内核资源只锁1毫秒的时间片段，所以等待自旋锁的释放不会消耗太多CPU的时间。

      自旋锁需要阻止在代码运行过程中出现的任何并发干扰。这些“干扰”包括： 中断，包括硬件中断和软件中断 （仅在中断代码可能访问临界区时需要）,内核抢占（仅存在于可抢占内核中）,其他处理器对同一临界区的访问 （仅SMP系统），spin lock需要在芯片底层实现物理上的内存地址独占访问，并且在实现上使用特殊的汇编指令访问。请看参考资料中对于自旋锁的实现分析。以arm为例，从存在SMP的ARM构架指令集开始（V6、V7），采用LDREX和STREX指令实现真正的自旋等待。

    深入分析_linux_spinlock_实现机制    Linux内核同步 - spin_lock

2  semaphore 

      由于信号量只能进行两种操作等待和发送信号，即P(sv)和V(sv),他们的行为是这样的：P(sv)：如果sv的值大于零，就给它减1；如果它的值为零，就挂起该进程的执行V(sv)：如果有其他进程因等待sv而被挂起，就让它恢复运行，如果没有进程因等待sv而挂起，就给它加1.
    举个例子，就是 两个进程共享信号量sv，一旦其中一个进程执行了P(sv)操作，它将得到信号量，并可以进入临界区，使sv减1。而第二个进程将被阻止进入临界区，因为 当它试图执行P(sv)时，sv为0，它会被挂起以等待第一个进程离开临界区域并执行V(sv)释放信号量，这时第二个进程就可以恢复执行。

      信号量其相应的接口也有两种(posix信号量)(systemv信号量 )这两者之间略有区别， systemv实际上是一个信号量组,posix信号量)，两者有很大的区别。

    1 system V的信号量是信号量集，可以包括多个信号灯（有个数组），每个操作可以同时操作多个信号灯posix是单个信号灯，POSIX有名信号灯支持进程间通信，无名信号灯放在共享内存中时可以用于进程间通信。
     2 POSIX信号量在有些平台并没有被实现，比如：SUSE8，而SYSTEM V大多数LINUX/UNIX都已经实现。两者都可以用于进程和线程间通信。但一般来说，system v信号量用于 进程间同步、有名信号灯既可用于线程间的同步，又可以用于进程间的同步、posix无名用于同一个进程的不同线程间，如果无名信号量要用于进程间同步，信号量要放在共享内存中。
      3 POSIX有两种类型的信号量，有名信号量和无名信号量。有名信号量像system v信号量一样由一个名字标识。
      4 POSIX通过sem_open单一的调用就完成了信号量的创建、初始化和权限的设置，而system v要两步。也就是说posix 信号是多线程，多进程安全的，而system v不是，可能会出现问题。
      5 system V信号量通过一个int类型的值来标识自己（类似于调用open()返回的fd），而sem_open函数返回sem_t类型（长整形）作为posix信号量的标识值。
      6 对于System V信号量你可以控制每次自增或是自减的信号量计数，而在Posix里面，信号量计数每次只能自增或是自减1。
      7 Posix无名信号量提供一种非常驻的信号量机制。
      8 相关进程: 如果进程是从一已经存在的进程创建，并最终操作这个创建进程的资源，那么这些进程被称为相关的。 

      基于system V是一个信号量集并且不单单只能自增自减1，我们就可以用两个信号量一个读信号量和写信号量组成一个信号量集实现一个基于进程的读写锁，很简‘单，读的时候把读信号量加一(p操作)，写的时候把写信号量加一，但是读之前要等写信号量为减为0(v操作)，写之前要等到读信号量和写信号量都减为0。

 linux 2.6 内核结构

struct semaphore {
     spinlock_t lock;
     unsigned int count;
     struct list_head wait_list;
};
    semaphore工作原理及其使用案例     进程间通信之-信号量semaphore--linux内核剖析

3 mutex

      mutex是用作互斥的，而semaphore是用作同步的。如果不去深究linux内核mutex和信号量的底层实现，在其功能上可以理解为mutex是一个受限制的信号量，也就是说，mutex一定是为0或者1，而semaphore可以是任意的数，所以如果使用mutex，那第一个进入临界区的进程一定可以执行，而其他的进程必须等待。而semaphore则不一定，如果一开始初始化为0，则所有进程都必须等待。互斥量的加锁和解锁必须由同一线程分别对应使用，信号量可以由一个线程释放，另一个线程得到。

 linux 2.6  内核结构

struct mutex {
   atomic_t  count; //引用计数器,1: 所可以利用,小于等于0：该锁已被获取，需要等待
   spinlock_t  wait_lock;//自旋锁类型，保证多cpu下，对等待队列访问是安全的。
   spinlock_t  wait_lock; //等待队列，如果该锁被获取，任务将挂在此队列上，等待调度。
   struct list_head wait_list;
};
  linux 2.6 互斥锁的实现-源码分析

       信号量mutex是sleep-waiting。 就是说当没有获得mutex时，会有上下文切换，将自己、加到忙等待队列中，直到另外一个线程释放mutex并唤醒它，而这时CPU是空闲的，可以调度别的任务处理。

        而自旋锁spin lock是busy-waiting。就是说当没有可用的锁时，就一直忙等待并不停的进行锁请求，直到得到这个锁为止。这个过程中cpu始终处于忙状态，不能做别的任务。
        Mutex适合对锁操作非常频繁的场景，并且具有更好的适应性。尽管相比spin lock它会花费更多的开销（主要是上下文切换），但是它能适合实际开发中复杂的应用场景，在保证一定性能的前提下提供更大的灵活度。spin lock的lock/unlock性能更好(花费更少的cpu指令)，但是它只适应用于临界区运行时间很短的场景。

        而在实际软件开发中，除非程序员对自己的程序的锁操作行为非常的了解，否则使用spin lock不是一个好主意(通常一个多线程程序中对锁的操作有数以万次，如果失败的锁操作(contended lock requests)过多的话就会浪费很多的时间进行空等待)。更保险的方法或许是先（保守的）使用 Mutex，然后如果对性能还有进一步的需求，可以尝试使用spin lock进行调优。毕竟我们的程序不像Linux kernel那样对性能需求那么高(Linux Kernel最常用的锁操作是spin lock和rw lock)。

应用场景

低开销加锁               优先使用自旋锁
短期锁定                 优先使用自旋锁
长期加锁                 优先使用信号量
中断上下文中加锁          使用自旋锁
持有锁是需要睡眠、调度     使用信号量

最后说到futex

   当然上面提到的内核数据结构如果我们不是开发linux内核，我们在实际开发中不可能直接操作这些数据结构，我们只能调用glibc接口。glibc是GNU发布的libc库，即c运行库。glibc是linux系统中最底层的api，几乎其它任何运行库都会依赖于glibc。glibc除了封装linux操作系统所提供的系统服务外，它本身也提供了许多其它一些必要功能服务的实现。由于 glibc 囊括了几乎所有的 UNIX通行的标准，可以想见其内容包罗万象。不管任何语言实现的任何各种各样的锁都是基于glibc和汇编来实现的，因为几乎所有的现代编程语言都是用c或者c++来写的，当你去看glibc的源码时你会看到pthread和信号量等接口的实现都出现了futex这个东西。这个东西有什么用呢？

     Futex是一种用户态和内核态混合的同步机制。首先，同步的进程间通过mmap共享一段内存，futex变量就位于这段共享的内存中且操作是原子的，当进程尝试进入互斥区或者退出互斥区的时候，先去查看共享内存中的futex变量，如果没有竞争发生，则只修改futex,而不用再执行系统调用了。当通过访问futex变量告诉进程有竞争发生，则还是得执行系统调用去完成相应的处理(wait 或者 wake up)。简单的说，futex就是通过在用户态的检查，（motivation）如果了解到没有竞争就不用陷入内核了，大大提高了low-contention时候的效率。
     为什么要有futex， 经研究发现，很多同步是无竞争的，即某个进程进入互斥区，到再从某个互斥区出来这段时间，常常是没有进程也要进这个互斥区或者请求同一同步变量的。但是在这种情况下，这个进程也要陷入内核去看看有没有人和它竞争，退出的时侯还要陷入内核去看看有没有进程等待在同一同步变量上。这些不必要的系统调用(或者说内核陷入)造成了大量的性能开销。为了解决这个问题，Futex就应运而生。前面的概念已经说了，futex是一种用户态和内核态混合同步机制，为什么会是用户态+内核态，听起来有点复杂，由于我们应用程序很多场景下多线程都是非竞争的，也就是说多任务在同一时刻同时操作临界区的概率是比较小的，大多数情况是没有竞争的，在早期内核同步互斥操作必须要进入内核态，由内核来提供同步机制，这就导致在非竞争的情况下，互斥操作扔要通过系统调用进入内核态。

    其实就是一句话，通过futex我们可以避免无谓的互斥操作，大大增加同步效率。

为什么需要内核锁?
多核处理器下，会存在多个进程处于内核态的情况，而在内核态下，进程是可以访问所有内核数据的，因此要对共享数据进行保护，即互斥处理

有哪些内核锁机制?
(1)原子操作
atomic_t数据类型，atomic_inc(atomic_t *v)将v加1
原子操作比普通操作效率要低，因此必要时才使用，且不能与普通操作混合使用
如果是单核处理器，则原子操作与普通操作相同
(2)自旋锁
spinlock_t数据类型，spin_lock(&lock)和spin_unlock(&lock)是加锁和解锁
等待解锁的进程将反复检查锁是否释放，而不会进入睡眠状态(忙等待)，所以常用于短期保护某段代码
同时，持有自旋锁的进程也不允许睡眠，不然会造成死锁——因为睡眠可能造成持有锁的进程被重新调度，而再次申请自己已持有的锁
如果是单核处理器，则自旋锁定义为空操作，因为简单的关闭中断即可实现互斥
(3)信号量与互斥量
struct semaphore数据类型，down(struct semaphore * sem)和up(struct semaphore * sem)是占用和释放
struct mutex数据类型，mutex_lock(struct mutex *lock)和mutex_unlock(struct mutex *lock)是加锁和解锁
竞争信号量与互斥量时需要进行进程睡眠和唤醒，代价较高，所以不适于短期代码保护，适用于保护较长的临界区

互斥量与信号量的区别?(转载但找不到原文出处)
(1)互斥量用于线程的互斥，信号线用于线程的同步
这是互斥量和信号量的根本区别，也就是互斥和同步之间的区别
互斥：是指某一资源同时只允许一个访问者对其进行访问，具有唯一性和排它性。但互斥无法限制访问者对资源的访问顺序，即访问是无序的
同步：是指在互斥的基础上（大多数情况），通过其它机制实现访问者对资源的有序访问。在大多数情况下，同步已经实现了互斥，特别是所有写入资源的情况必定是互斥的。少数情况是指可以允许多个访问者同时访问资源
(2)互斥量值只能为0/1，信号量值可以为非负整数
也就是说，一个互斥量只能用于一个资源的互斥访问，它不能实现多个资源的多线程互斥问题。信号量可以实现多个同类资源的多线程互斥和同步。当信号量为单值信号量是，也可以完成一个资源的互斥访问
(3)互斥量的加锁和解锁必须由同一线程分别对应使用，信号量可以由一个线程释放，另一个线程得到

mutex，一句话：保护共享资源。
典型的例子就是买票：
票是共享资源，现在有两个线程同时过来买票。如果你不用mutex在线程里把票锁住，那么就可能出现“把同一张票卖给两个不同的人（线程）”的情况。

一般人不明白semaphore和mutex的区别，根本原因是不知道semaphore的用途。
semaphore的用途，一句话：调度线程。
有的人用semaphore也可以把上面例子中的票“保护"起来以防止共享资源冲突，必须承认这是可行的，但是semaphore不是让你用来做这个的；如果你要做这件事，请用mutex。

在网上、包括stackoverflow等著名论坛上，有一个流传很广的厕所例子：
mutex是一个厕所一把钥匙，谁抢上钥匙谁用厕所，谁没抢上谁就等着；semaphore是多个同样厕所多把同样的钥匙 ---- 只要你能拿到一把钥匙，你就可以随便找一个空着的厕所进去。
事实上，这个例子对初学者、特别是刚刚学过mutex的初学者来说非常糟糕 ----- 我第一次读到这个例子的第一反应是：semaphore是线程池？？？所以，请务必忘记这个例子。
另外，有人也会说：mutex就是semaphore的value等于1的情况。
这句话不能说不对，但是对于初学者来说，请先把这句话视为错误；等你将来彻底融会贯通这部分知识了，你才能真正理解上面这句话到底是什么意思。总之请务必记住：mutex干的活儿和semaphore干的活儿不要混起来。



在这里，我模拟一个最典型的使用semaphore的场景：
a源自一个线程，b源自另一个线程，计算c = a + b也是一个线程。（即一共三个线程）

显然，第三个线程必须等第一、二个线程执行完毕它才能执行。
在这个时候，我们就需要调度线程了：让第一、二个线程执行完毕后，再执行第三个线程。
此时，就需要用semaphore了。
int a, b, c;
void geta()
{
    a = calculatea();
    semaphore_increase();
}
 
void getb()
{
    b = calculateb();
    semaphore_increase();
}
 
 
void getc()
{
    semaphore_decrease();
    semaphore_decrease();
    c = a + b;
}
 
t1 = thread_create(geta);
t2 = thread_create(getb);
t3 = thread_create(getc);
thread_join(t3);
这就是semaphore最典型的用法。
说白了，调度线程，就是：一些线程生产（increase）同时另一些线程消费（decrease），semaphore可以让生产和消费保持合乎逻辑的执行顺序。
而线程池是程序员根据具体的硬件水平和不同的设计需求、为了达到最佳的运行效果而避免反复新建和释放线程同时对同一时刻启动的线程数量的限制，这完全是两码事。
比如如果你要计算z = a + b +...+ x + y ...的结果，同时每个加数都是一个线程，那么计算z的线程和每个加数的线程之间的逻辑顺序是通过semaphore来调度的；而至于你运行该程序的时候到底要允许最多同时启动几个线程，则是用线程池来实现的。

请回头看那个让大家忘记的厕所例子。我之所以让大家忘记这个例子，是因为如果你从这个角度去学习semaphore的话，一定会和mutex混为一谈。semaphore的本质就是调度线程 ---- 在充分理解了这个概念后，我们再看这个例子。

semaphore是通过一个值来实现线程的调度的，因此借助这种机制，我们也可以实现对线程数量的限制。例子我就不写了，如果你看懂了上面的c = a + b的例子，相信你可以轻松写出来用semaphore限制线程数量的例子。而当我们把线程数量限制为1时，你会发现：共享资源受到了保护 ------ 任意时刻只有一个线程在运行，因此共享资源当然等效于受到了保护。但是我要再提醒一下，如果你要对共享资源进行保护，请用mutex；到底应该用条件锁还是用semaphore，请务必想清楚。通过semaphore来实现对共享资源的保护的确可行但是是对semaphore的一种错用。

只要你能搞清楚锁、条件锁和semaphore为什么而生、或者说它们是面对什么样的设计需求、为了解决什么样类型的问题才出现的，你自然就不会把他们混淆起来。

semaphore和条件锁的区别：
条件锁，本质上还是锁，它的用途，还是围绕“共享资源”的。条件锁最典型的用途就是：防止不停地循环去判断一个共享资源是否满足某个条件。
比如还是买票的例子：
我们除了买票的线程外，现在再加一个线程：如果票数等于零，那么就要挂出“票已售完”的牌子。这种情况下如果没有条件锁，我们就不得不在“挂牌子”这个线程里不断地lock和unlock而在大多数情况下票数总是不等于零，这样的结果就是：占用了很多CPU资源但是大多数时候什么都没做。
另外，假如我们还有一个线程，是在票数等于零时向上级部门申请新的票。同理，问题和上面的一样。而如果有了条件锁，我们就可以避免这种问题，而且还可以一次性地通知所有被条件锁锁住的线程。
这里有个问题，是关于条件锁的：pthread_cond_wait 为什么需要传递 mutex 参数？
不清楚条件锁的朋友可以看一下。
总之请记住：条件锁，是为了避免绝大多数情况下都是lock ---> 判断条件 ----> unlock的这种很占资源但又不干什么事情的线程。它和semaphore的用途是不同的。


简而言之，锁是服务于共享资源的；而semaphore是服务于多个线程间的执行的逻辑顺序的。




[转]mutex和spin lock的区别
mutex和spin lock的区别和应用(sleep-waiting和busy-waiting的区别)2011-10-19 11:43
信号量mutex是sleep-waiting。 就是说当没有获得mutex时，会有上下文切换，将自己、加到忙等待队列中，直到另外一个线程释放mutex并唤醒它，而这时CPU是空闲的，可以调度别的任务处理。

而自旋锁spin lock是busy-waiting。就是说当没有可用的锁时，就一直忙等待并不停的进行锁请求，直到得到这个锁为止。这个过程中cpu始终处于忙状态，不能做别的任务。

例如在一个双核的机器上有两个线程(线程A和线程B)，它们分别运行在Core0 和Core1上。 用spin-lock，coer0上的线程就会始终占用CPU。
另外一个值得注意的细节是spin lock耗费了更多的user time。这就是因为两个线程分别运行在两个核上，大部分时间只有一个线程能拿到锁，所以另一个线程就一直在它运行的core上进行忙等待，CPU占用率一直是100%；而mutex则不同，当对锁的请求失败后上下文切换就会发生，这样就能空出一个核来进行别的运算任务了。（其实这种上下文切换对已经拿着锁的那个线程性能也是有影响的，因为当该线程释放该锁时它需要通知操作系统去唤醒那些被阻塞的线程，这也是额外的开销）

总结
（1）Mutex适合对锁操作非常频繁的场景，并且具有更好的适应性。尽管相比spin lock它会花费更多的开销（主要是上下文切换），但是它能适合实际开发中复杂的应用场景，在保证一定性能的前提下提供更大的灵活度。

（2）spin lock的lock/unlock性能更好(花费更少的cpu指令)，但是它只适应用于临界区运行时间很短的场景。而在实际软件开发中，除非程序员对自己的程序的锁操作行为非常的了解，否则使用spin lock不是一个好主意(通常一个多线程程序中对锁的操作有数以万次，如果失败的锁操作(contended lock requests)过多的话就会浪费很多的时间进行空等待)。

（3）更保险的方法或许是先（保守的）使用 Mutex，然后如果对性能还有进一步的需求，可以尝试使用spin lock进行调优。毕竟我们的程序不像Linux kernel那样对性能需求那么高(Linux Kernel最常用的锁操作是spin lock和rw lock)。

 
总结：

互斥锁支持且只支持互斥访问资源，且只有加锁人才能解锁。即互斥量只支持互斥机制；

信号量支持互斥访问资源，加锁者和解锁者可以是不同人，因此可以实现生产者+消费者模式，可以使进程/线程间有序生产/消费资源。即信号量用于实现进程/线程间同步；

自旋锁与二者的区别就是不能阻塞，一直等待获得锁；



互斥：

是指某一资源同时只允许一个访问者对其进行访问，具有唯一性和排它性。但互斥无法限制访问者对资源的访问顺序，即访问是无序的。

某些资源只能被独占使用，或某些关键数据只能被一个进程修改，此时需要互斥机制。
设备有可能被多个用户使用，驱动程序就有可能被多个进程调用，因此如果有资源可能在多个进程间被共享，

比如设备的配置空间对应的内存就有可能被多个设备同时访问，例如nvme设备在被多个用户进程写入时，

io queue的head tail指针会随时更新，如果没有锁机制保证互斥访问，就可能会造成更新错误。

因此需要同步机制保证资源的一致性，避免同时访问造成冲突。

机制：
互斥锁Mutex lock/unlock机制，自旋锁spin lock/unlock机制，原子操作atomic_t机制。


同步：

是指在互斥的基础上（大多数情况），通过其它机制实现访问者对资源的有序访问。在大多数情况下，同步已经实现了互斥，特别是所有写入资源的情况必定是互斥的。少数情况是指可以允许多个访问者同时访问资源。

本质上就是进程间通信，A进程干完某个活之后告诉B进程可以开始干了，或者A进程释放了B进程需要的资源。
机制：
信号量semaphore donw/up机制，completion机制，wait event queue机制。

当资源需要独占使用且只有拥有者能释放时，使用互斥锁机制对资源进行保护和使用；

当资源数量大于1，且可以多个使用者同时使用，则可以使用信号量机制对资源进行分配和使用；

当资源是生产者+消费者模式时，适合使用信号量机制由不同的进程/线程生产资源和使用资源。

 
一、 以2.6.38以前的内核为例， 讲spinlock、 mutex 以及 semaphore
1. spinlock更原始，效率高，但讲究更多，不能随便用。
2. 个人觉得初级阶段不要去深挖mutex 以及 semaphore的不同，用法类似。在内核代码里面搜索，感觉 DEFINE_MUTEX + mutex_lock_xx + mutex_unlock 用的更多。
3. 在内核里面这三个符号发挥的作用就是： 自旋锁与互斥体。
semaphore:内核中的信号量通常用作mutex互斥体（信号量初值初始化为1，即binary semaphore的方式，就达到了互斥的效果）。
mutex：顾名思义， 互斥体。
所以在内核里面，mutex_lock()和down()的使用情景基本上相同。


 
//spinlock.h
/******
*API
*spin_lock
*spin_lock_bh
*spin_lock_irq
*spin_trylock
*spin_trylock_bh
*spin_trylock_irq
*spin_unlock
*spin_unlock_bh
*spin_unlock_irq
*spin_unlock_irqrestore
*spin_unlock_wait
******/


 
//semaphore.h
用 DECLARE_MUTEX 定义了一个count==1 的信号量（binary semaphore)。

#define DECLARE_MUTEX(name)    \
    struct semaphore name = __SEMAPHORE_INITIALIZER(name, 1)
    
struct semaphore {
    spinlock_t        lock;
    unsigned int        count;
    struct list_head    wait_list;
};

#define __SEMAPHORE_INITIALIZER(name, n)                \
{                                    \
    .lock        = __SPIN_LOCK_UNLOCKED((name).lock),        \
    .count        = n,                        \
    .wait_list    = LIST_HEAD_INIT((name).wait_list),        \
}
    
#define init_MUTEX(sem)        sema_init(sem, 1)
#define init_MUTEX_LOCKED(sem)    sema_init(sem, 0)


 
/*****
*API：
*#define init_MUTEX(sem)        sema_init(sem, 1)
*#define init_MUTEX_LOCKED(sem)    sema_init(sem, 0)
*extern void down(struct semaphore *sem);
*extern int __must_check down_interruptible(struct semaphore *sem);
*extern int __must_check down_killable(struct semaphore *sem);
*extern int __must_check down_trylock(struct semaphore *sem);
*extern int __must_check down_timeout(struct semaphore *sem, long jiffies);
*extern void up(struct semaphore *sem);
****/

//mutex.h
#define DEFINE_MUTEX(mutexname) \
    struct mutex mutexname = __MUTEX_INITIALIZER(mutexname)
    
    #define __MUTEX_INITIALIZER(lockname) \
        { .count = ATOMIC_INIT(1) \
        , .wait_lock = __SPIN_LOCK_UNLOCKED(lockname.wait_lock) \
        , .wait_list = LIST_HEAD_INIT(lockname.wait_list) \
        __DEBUG_MUTEX_INITIALIZER(lockname) \
        __DEP_MAP_MUTEX_INITIALIZER(lockname) }


 
/*
 * Simple, straightforward mutexes with strict semantics:
 *
 * - only one task can hold the mutex at a time
 * - only the owner can unlock the mutex
 * - multiple unlocks are not permitted
 * - recursive locking is not permitted
 * - a mutex object must be initialized via the API
 * - a mutex object must not be initialized via memset or copying
 * - task may not exit with mutex held
 * - memory areas where held locks reside must not be freed
 * - held mutexes must not be reinitialized
 * - mutexes may not be used in hardware or software interrupt
 *   contexts such as tasklets and timers
 *
 * These semantics are fully enforced when DEBUG_MUTEXES is
 * enabled. Furthermore, besides enforcing the above rules, the mutex
 * debugging code also implements a number of additional features
 * that make lock debugging easier and faster:
 *
 * - uses symbolic names of mutexes, whenever they are printed in debug output
 * - point-of-acquire tracking, symbolic lookup of function names
 * - list of all locks held in the system, printout of them
 * - owner tracking
 * - detects self-recursing locks and prints out all relevant info
 * - detects multi-task circular deadlocks and prints out all affected
 *   locks and tasks (and only those tasks)
 */
struct mutex {
    /* 1: unlocked, 0: locked, negative: locked, possible waiters */
    atomic_t        count;
    spinlock_t        wait_lock;
    struct list_head    wait_list;
#if defined(CONFIG_DEBUG_MUTEXES) || defined(CONFIG_SMP)
    struct thread_info    *owner;
#endif
#ifdef CONFIG_DEBUG_MUTEXES
    const char         *name;
    void            *magic;
#endif
#ifdef CONFIG_DEBUG_LOCK_ALLOC
    struct lockdep_map    dep_map;
#endif
};


 
/*******
*API:
*extern void mutex_lock(struct mutex *lock);
*extern int __must_check mutex_lock_interruptible(struct mutex *lock);
*extern int __must_check mutex_lock_killable(struct mutex *lock);
*extern void mutex_unlock(struct mutex *lock);
********/

EG1-:
    spinlock_t rtc_lock;
    spin_lock_init(&rtc_lock);//每个驱动都会事先初始化，只需要这一次初始化
 
    spin_lock_irq(&rtc_lock);
    //临界区
    spin_unlock_irq(&rtc_lock);
 
EG1-:
    unsigned long flags;
    static spinlock_t i2o_drivers_lock;
    spin_lock_init(&i2o_drivers_lock);//每个驱动都会事先初始化，只需要这一次初始化
 
    spin_lock_irqsave(&i2o_drivers_lock, flags);
    //临界区
    spin_unlock_irqrestore(&i2o_drivers_lock, flags);
 
EG2:
    static DECLARE_MUTEX(start_stop_sem);
    down(&start_stop_sem);
    //临界区
    up(&start_stop_sem);
 
EG3:
    static DEFINE_MUTEX(adutux_mutex);
    mutex_lock_interruptible(&adutux_mutex);
    //临界区
    mutex_unlock(&adutux_mutex);
二、 2.6.38以后DECLARE_MUTEX替换成DEFINE_SEMAPHORE（命名改变）, DEFINE_MUTEX用法不变

    static DEFINE_SEMAPHORE(msm_fb_pan_sem);// DECLARE_MUTEX
    down(&adb_probe_mutex);
    //临界区
    up(&adb_probe_mutex);
 
    static DEFINE_SEMAPHORE(bnx2x_prev_sem);
    down_interruptible(&bnx2x_prev_sem);
    //临界区
    up(&bnx2x_prev_sem); 
Linux 2.6.36以后file_operations和DECLARE_MUTEX 的变化http://blog.csdn.net/heanyu/article/details/6757917

在include/linux/semaphore.h 中将#define DECLARE_MUTEX(name)   改成了 #define DEFINE_SEMAPHORE(name)   【命名】


 
　　　　 
三、自旋锁与信号量
1. 自旋锁
简单的说，自旋锁在内核中主要用来防止多处理器中并发访问临界区，防止内核抢占造成的竞争。【适用于多处理器】【自旋锁会影响内核调度】
另外自旋锁不允许任务睡眠(持有自旋锁的任务睡眠会造成自死锁——因为睡眠有可能造成持有锁的内核任务被重新调度，而再次申请自己已持有的锁)，它能够在中断上下文中使用。【不允许任务睡眠】

锁定一个自旋锁的函数有四个：
void spin_lock(spinlock_t *lock); //最基本得自旋锁函数，它不失效本地中断。
void spin_lock_irqsave(spinlock_t *lock, unsigned long flags);//在获得自旋锁之前禁用硬中断（只在本地处理器上），而先前的中断状态保存在flags中
void spin_lockirq(spinlock_t *lock);//在获得自旋锁之前禁用硬中断（只在本地处理器上），不保存中断状态
void spin_lock_bh(spinlock_t *lock);//在获得锁前禁用软中断，保持硬中断打开状态

2. 信号量
内核中的信号量通常用作mutex互斥体（信号量初值初始化为1就达到了互斥的效果）。

如果代码需要睡眠——这往往是发生在和用户空间同步时——使用信号量是唯一的选择。由于不受睡眠的限制，使用信号量通常来说更加简单一些。【信号量使用简单】
如果需要在自旋锁和信号量中作选择，应该取决于锁被持有的时间长短。理想情况是所有的锁都应该尽可能短的被持有，但是如果锁的持有时间较长的话，使用信号量是更好的选择。【如果锁占用的时间较长，信号量更好】
另外，信号量不同于自旋锁，它不会关闭内核抢占，所以持有信号量的代码可以被抢占。这意味者信号量不会对影响调度反应时间带来负面影响。【信号量不会影响内核调度】

3. 使用情景对比
=============================================
需求                     建议的加锁方法 
低开销加锁               优先使用自旋锁
短期锁定                 优先使用自旋锁
长期加锁                 优先使用信号量
中断上下文中加锁         使用自旋锁
持有锁是需要睡眠、调度    使用信号量

信号量在内核中的定义如下：

struct semaphore {
	raw_spinlock_t		lock;///自旋锁
	unsigned int		count;///count=1时可进行互斥操作
	struct list_head	wait_list;
};
信号量的初始化：

sem_init(&sem,val);///var代表信号量的初始值

获取信号量：

down（&sem）;若此时信号量为0，则该进程会会处于睡眠状态，因此该函数不可用于中断上下文中。

接下来分析一下获取信号量的源码：

static noinline void __sched __down(struct semaphore *sem)
{
	__down_common(sem, TASK_UNINTERRUPTIBLE, MAX_SCHEDULE_TIMEOUT);
}
static inline int __sched __down_common(struct semaphore *sem, long state,
								long timeout)
{
	struct task_struct *task = current;
	struct semaphore_waiter waiter;

	list_add_tail(&waiter.list, &sem->wait_list);
	waiter.task = task;
	waiter.up = 0;
///死循环
	for (;;) {
	///如果当前进程被信号唤醒，则退出
		if (signal_pending_state(state, task))
			goto interrupted;
	///如果进程的等待时间超时，则退出
		if (timeout <= 0)
			goto timed_out;
		__set_task_state(task, state);
		raw_spin_unlock_irq(&sem->lock);
		///在等待队列中等待调度。
		timeout = schedule_timeout(timeout);
		raw_spin_lock_irq(&sem->lock);
		///如果调度是由信号量的释放而唤醒的，则返回0
		if (waiter.up)
			return 0;
	}

  ......
}
释放信号量

up(&sem);

互斥信号量：

struct mutex {
	/* 1: unlocked, 0: locked, negative: locked, possible waiters */
	atomic_t		count;
	spinlock_t		wait_lock;
	struct list_head	wait_list;
 ......
};
互斥信号量的初始化：

init_mutex(&sem);

同样作为同步操作，mutex、spinlock、semaphore有如下差异：

1、mutex的count初始化为1，而semaphore则初始化为0

2、mutex的使用者必须为同一线程，即必须成对使用，而semaphore可以由不同的线程执行P.V操作。

3、进程在获取不到信号量的时候执行的是sleep操作，而进程在获取不到自旋锁的时候执行的是忙等待操作。因此，不难看出，如果需要保护的临界区比较小，锁的持有时间比较短的情况下，通常使用spinlock。这样可以不需要对等待锁的进程执行睡眠/唤醒操作，大大节省了cpu时间。因此，spinlock通常作为多处理器之间的同步操作。
