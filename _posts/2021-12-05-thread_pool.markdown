---
title: MySQL 线程池总结
layout: post
category: mysql
author: 夏泽民
---
在 MySQL 5.6出现以前，MySQL 处理连接的方式是 One-Connection-Per-Thread,即对于每一个数据库连接，MySQL-Server都会创建一个独立的线程服务，请求结束后，销毁线程。再来一个连接请求，则再创建一个连接，结束后再进行销毁。这种方式在高并发情况下，会导致线程的频繁创建和释放。当然，通过 thread-cache，我们可以将线程缓存起来，以供下次使用，避免频繁创建和释放的问题，但是无法解决高连接数的问题。One-Connection-Per-Thread 方式随着连接数暴增，导致需要创建同样多的服务线程，高并发线程意味着高的内存消耗，更多的上下文切换(cpu cache命中率降低)以及更多的资源竞争，导致服务出现抖动。相对于 One-Thread-Per-Connection 方式，一个线程对应一个连接，Thread-Pool 实现方式中，线程处理的最小单位是statement(语句)，一个线程可以处理多个连接的请求。这样，在保证充分利用硬件资源情况下(合理设置线程池大小)，可以避免瞬间连接数暴增导致的服务器抖动。
<!-- more -->
调度方式实现
MySQL-Server 同时支持3种连接管理方式，包括No-Threads，One-Thread-Per-Connection 和 Pool-Threads。

No-Threads 表示处理连接使用主线程处理，不额外创建线程，这种方式主要用于调试；
One-Thread-Per-Connection 是线程池出现以前最常用的方式，为每一个连接创建一个线程服务；
Pool-Threads 则是本文所讨论的线程池方式。Mysql-Server通过一组函数指针来同时支持3种连接管理方式，对于特定的方式，将函数指针设置成特定的回调函数，连接管理方式通过thread_handling参数控制，代码如下：
if (thread_handling <= SCHEDULER_ONE_THREAD_PER_CONNECTION)   
   one_thread_per_connection_scheduler(thread_scheduler,&max_connections, &connection_count);
else if (thread_handling == SCHEDULER_NO_THREADS)
     one_thread_scheduler(thread_scheduler);
else                                 
    pool_of_threads_scheduler(thread_scheduler, &max_connections,&connection_count); 
连接管理流程
通过poll监听mysql端口的连接请求 收到连接后，调用accept接口，创建通信socket 初始化thd实例，vio对象等 根据thread_handling方式设置，初始化thd实例的scheduler函数指针 调用scheduler特定的add_connection函数新建连接 下面代码展示了scheduler_functions模板和线程池对模板回调函数的实现，这个是多种连接管理的核心。

struct scheduler_functions                        
{  
uint   max_threads;
uint   *connection_count;                          
ulong *max_connections;                          
bool (*init)(void);                              
bool (*init_new_connection_thread)(void);       
void (*add_connection)(THD *thd);
void (*thd_wait_begin)(THD *thd, int wait_type); 
void (*thd_wait_end)(THD *thd);                  
void (*post_kill_notification)(THD *thd);        
bool (*end_thread)(THD *thd, bool cache_thread);
void (*end)(void);
};
static scheduler_functions tp_scheduler_functions=
{ 
  0, // max_threads
  NULL,
  NULL, 
  tp_init, // init
  NULL, // init_new_connection_thread
  tp_add_connection, // add_connection
  tp_wait_begin, // thd_wait_begin            
  tp_wait_end, // thd_wait_end
  tp_post_kill_notification,  // post_kill_notification 
  NULL,   // end_thread
  tp_end  // end
};
线程池的相关参数
thread_handling: 表示线程池模型。thread_pool_size:表示线程池的group个数，一般设置为当前CPU核心数目。理想情况下，一个group一个活跃的工作线程，达到充分利用CPU的目的。thread_pool_stall_limit:用于timer线程定期检查group是否“停滞”，参数表示检测的间隔。thread_pool_idle_timeout:当一个worker空闲一段时间后会自动退出，保证线程池中的工作线程在满足请求的情况下，保持比较低的水平。thread_pool_oversubscribe:该参数用于控制CPU核心上“超频”的线程数。这个参数设置值不含listen线程计数。threadpool_high_prio_mode:表示优先队列的模式。线程池实现

关键接口
tp_add_connection[处理新连接]

创建一个connection对象
根据thread_id%group_count确定connection分配到哪个group
将connection放进对应group的队列
如果当前活跃线程数为0，则创建一个工作线程
worker_main[工作线程]

调用get_event获取请求
如果存在请求，则调用handle_event进行处理
否则，表示队列中已经没有请求，退出结束。
get_event[获取请求]

获取一个连接请求
如果存在，则立即返回，结束
若此时group内没有listener，则线程转换为listener线程，阻塞等待
若存在listener，则将线程加入等待队列头部
线程休眠指定的时间(thread_pool_idle_timeout)
如果依然没有被唤醒，是超时，则线程结束，结束退出
否则，表示队列里有连接请求到来，跳转1
备注：获取连接请求前，会判断当前的活跃线程数是否超过thread_pool_oversubscribe+1，若超过，则将线程进入休眠状态。

handle_event[处理请求]

判断连接是否进行登录验证，若没有，则进行登录验证
关联thd实例信息
获取网络数据包，分析请求
调用do_command函数循环处理请求
获取thd实例的套接字句柄，判断句柄是否在epoll的监听列表中
若没有，调用epoll_ctl进行关联
结束
listener[监听线程]

调用epoll_wait进行对group关联的套接字监听，阻塞等待
若请求到来，从阻塞中恢复
根据连接的优先级别，确定是放入普通队列还是优先队列
判断队列中任务是否为空
若队列为空，则listener转换为worker线程
若group内没有活跃线程，则唤醒一个线程
备注：这里epoll_wait监听group内所有连接的套接字，然后将监听到的连接请求push到队列，worker线程从队列中获取任务，然后执行。

timer_thread[监控线程]

若没有listener线程，并且最近没有io_event事件
则创建一个唤醒或创建一个工作线程
若group最近一段时间没有处理请求，并且队列里面有请求，则
表示group已经stall，则唤醒或创建线程
检查是否有连接超时
备注：timer线程通过调用check_stall判断group是否处于stall状态，通过调用timeout_check检查客户端连接是否超时。

tp_wait_begin[进入等待状态流程]

active_thread_count减1，waiting_thread_count加1

设置connection->waiting= true

若活跃线程数为0，并且任务队列不为空，或者没有监听线程，则

唤醒或创建一个线程

tp_wait_end[结束等待状态流程]

设置connection的waiting状态为false
active_thread_count加1，waiting_thread_count减1
备注：

waiting_threads这个list里面的线程是空闲线程，并非等待线程，所谓空闲线程是随时可以处理任务的线程，而等待线程则是因为等待锁，或等待io操作等无法处理任务的线程。

tp_wait_begin和tp_wait_end的主要作用是由于汇报状态，即使更新active_thread_count和waiting_thread_count的信息。

tp_init/tp_end

分别调用thread_group_init和thread_group_close来初始化和销毁线程池

线程池与连接池
连接池通常实现在 Client 端，是指应用(客户端)创建预先创建一定的连接，利用这些连接服务于客户端所有的DB请求。如果某一个时刻，空闲的连接数小于DB的请求数，则需要将请求排队，等待空闲连接处理。通过连接池可以复用连接，避免连接的频繁创建和释放，从而减少请求的平均响应时间，并且在请求繁忙时，通过请求排队，可以缓冲应用对DB的冲击。

线程池实现在server端，通过创建一定数量的线程服务DB请求，相对于 one-conection-per-thread 的一个线程服务一个连接的方式，线程池服务的最小单位是语句，即一个线程可以对应多个活跃的连接。通过线程池，可以将 server 端的服务线程数控制在一定的范围，减少了系统资源的竞争和线程上下文切换带来的消耗，同时也避免出现高连接数导致的高并发问题。

连接池和线程池相辅相成，通过连接池可以减少连接的创建和释放，提高请求的平均响应时间，并能很好地控制一个应用的DB连接数，但无法控制整个应用集群的连接数规模，从而导致高连接数，通过线程池则可以很好地应对高连接数，保证server端能提供稳定的服务。

线程池优化
调度死锁解决
引入线程池解决了多线程高并发的问题，但也带来一个隐患。假设，A，B两个事务被分配到不同的group中执行，A事务已经开始，并且持有锁，但由于A所在的group比较繁忙，导致A执行一条语句后，不能立即获得调度执行；而B事务依赖A事务释放锁资源，虽然B事务可以被调度起来，但由于无法获得锁资源，导致仍然需要等待，这就是所谓的调度死锁。由于一个group会同时处理多个连接，但多个连接不是对等的。比如，有的连接是第一次发送请求；而有的连接对应的事务已经开启，并且持有了部分锁资源。为了减少锁资源争用，后者显然应该比前者优先处理，以达到尽早释放锁资源的目的。因此在group里面，可以添加一个优先级队列，将已经持有锁的连接，或者已经开启的事务的连接发起的请求放入优先队列，工作线程首先从优先队列获取任务执行。

大查询处理
假设一种场景，某个group里面的连接都是大查询，那么group里面的工作线程数很快就会达到thread_pool_oversubscribe参数设置值，对于后续的连接请求，则会响应不及时(没有更多的连接来处理)，这时候group就发生了stall。

通过前面分析知道，timer线程会定期检查这种情况，并创建一个新的worker线程来处理请求。如果长查询来源于业务请求，则此时所有group都面临这种问题，此时主机可能会由于负载过大，导致hang住的情况。这种情况线程池本身无能为力，因为源头可能是烂SQL并发，或者SQL没有走对执行计划导致，通过其他方法，比如SQL高低水位限流或者SQL过滤手段可以应急处理。

但是，还有另外一种情况，就是dump任务。很多下游依赖于数据库的原始数据，通常通过dump命令将数据拉到下游，而这种dump任务通常都是耗时比较长，所以也可以认为是大查询。如果dump任务集中在一个group内，并导致其他正常业务请求无法立即响应，这个是不能容忍的，因为此时数据库并没有压力，只是因为采用了线程池策略，才导致了请求响应不及时，为了解决这个问题，我们将group中处理dump任务的线程不计入thread_pool_oversubscribe累计值，避免上述问题。

https://mp.weixin.qq.com/s?__biz=MzI4NjExMDA4NQ==&mid=2648456189&idx=1&sn=bd263ff9262bea2bdb585afce9a62d90&scene=58&subscene=0

1、减少线程重复创建与销毁部分的开销，提高性能



线程池技术通过预先创建一定数量的线程，在监听到有新的请求时，线程池直接从现有的线程中分配一个线程来提供服务，服务结束后这个线程不会直接销毁，而是又去处理其他的请求。这样就避免了线程和内存对象频繁创建和销毁，减少了上下文切换，提高了资源利用率，从而在一定程度上提高了系统的性能和稳定性。



2、对系统起到保护作用



线程池技术限制了并发线程数，相当于限制了MySQL的runing线程数，无论系统目前有多少连接或者请求，超过最大设置的线程数的都需要排队，让系统保持高性能水平，从而防止DB出现雪崩，对底层DB起到保护作用。



可能有人会问，使用连接池能否也达到类似的效果？



也许有的DBA会把线程池和连接池混淆，但其实两者是有很大区别的：连接池一般在客户端设置，而线程池是在DB服务器上配置；另外连接池可以起到避免了连接频繁创建和销毁，但是无法控制MySQL活动线程数的目标，在高并发场景下，无法起到保护DB的作用。比较好的方式是将连接池和线程池结合起来使用。

2、Thread Pool的组成



从架构图中可以看到Thread Pool由一个Timer线程和多个Thread Group组成，而每个Thread Group又由两个队列、一个listener线程和多个worker线程构成。下面分别来介绍各个部分的作用：



队列（高优先级队列和低优先级队列）



用来存放待执行的IO任务，分为高优先级队列和低优先级队列，高优先级队列的任务会优先被处理。



什么任务会放在高优先级队列呢？



事务中的语句会放到高优先级队列中，比如一个事务中有两个update的SQL，有1个已经执行，那么另外一个update的任务就会放在高优先级中。这里需要注意，如果是非事务引擎，或者开启了Autocommit的事务引擎，都会放到低优先级队列中。



还有一种情况会将任务放到高优先级队列中，如果语句在低优先级队列停留太久，该语句也会移到高优先级队列中，防止饿死。



listener线程



listener线程监听该线程group的语句，并确定当自己转变成worker线程，是立即执行对应的语句还是放到队列中，判断的标准是看队列中是否有待执行的语句。



如果队列中待执行的语句数量为0，而listener线程转换成worker线程，并立即执行对应的语句。如果队列中待执行的语句数量不为0，则认为任务比较多，将语句放入队列中，让其他的线程来处理。这里的机制是为了减少线程的创建，因为一般SQL执行都非常快。



worker线程



worker线程是真正干活的线程。



Timer线程



Timer线程是用来周期性检查group是否处于处于阻塞状态，当出现阻塞的时候，会通过唤醒线程或者新建线程来解决。



具体的检测方法为：通过queue_event_count的值和IO任务队列是否为空来判断线程组是否为阻塞状态。



每次worker线程检查队列中任务的时候，queue_event_count会+1，每次Timer检查完group是否阻塞的时候会将queue_event_count清0，如果检查的时候任务队列不为空，而queue_event_count为0，则说明任务队列没有被正常处理，此时该group出现了阻塞，Timer线程会唤醒worker线程或者新建一个wokrer线程来处理队列中的任务，防止group长时间被阻塞。



3、Thread Pool的是如何运作的？



下面描述极简的Thread Pool运作，只是简单描述，省略了大量的复杂逻辑，请不要挑刺~



Step1：请求连接到MySQL，根据threadid%thread_pool_size确定落在哪个group；



Step2：group中的listener线程监听到所在的group有新的请求以后，检查队列中是否有请求还未处理。如果没有，则自己转换为worker线程立即处理该请求，如果队列中还有未处理的请求，则将对应请求放到队列中，让其他的线程处理；



Step3：group中的thread线程检查队列的请求，如果队列中有请求，则进行处理，如果没有请求，则休眠，一直没有被唤醒，超过thread_pool_idle_timeout后就自动退出。线程结束。当然，获取请求之前会先检查group中的running线程数是否超过thread_pool_oversubscribe+1，如果超过也会休眠；



Step4：timer线程定期检查各个group是否有阻塞，如果有，就对wokrer线程进行唤醒或者创建一个新的worker线程。



4、Thread Pool的分配机制



线程池会根据参数thread_pool_size的大小分成若干的group，每个group各自维护客户端发起的连接，当客户端发起连接到MySQL的时候，MySQL会跟进连接的线程id（thread_id）对thread_pool_size进行取模，从而落到对应的group。



thread_pool_oversubscribe参数控制每个group的最大并发线程数，每个group的最大并发线程数为thread_pool_oversubscribe+1个。若对应的group达到了最大的并发线程数，则对应的连接就需要等待。这个分配机制在某个group中有多个慢SQL的场景下会导致普通的SQL运行时间很长，这个问题会在后面做详细描述。

https://mp.weixin.qq.com/s?__biz=MzkwOTIxNDQ3OA==&mid=2247533328&idx=1&sn=8c96aa41bd0b31f70fb43146ca714a34&scene=58&subscene=0