---
title: epoll_event
layout: post
category: linux
author: 夏泽民
---
结构体epoll_event被用于注册所感兴趣的事件和回传所发生待处理的事件，定义如下：
    typedef union epoll_data {
        void *ptr;
         int fd;
         __uint32_t u32;
         __uint64_t u64;
     } epoll_data_t;//保存触发事件的某个文件描述符相关的数据
     struct epoll_event {
         __uint32_t events;      /* epoll event */
         epoll_data_t data;      /* User data variable */
};

其中events表示感兴趣的事件和被触发的事件，可能的取值为：
EPOLLIN：表示对应的文件描述符可以读；
EPOLLOUT：表示对应的文件描述符可以写；
EPOLLPRI：表示对应的文件描述符有紧急的数可读；

EPOLLERR：表示对应的文件描述符发生错误；
EPOLLHUP：表示对应的文件描述符被挂断；
EPOLLET：    ET的epoll工作模式；
<!-- more -->
1、epoll_create函数
   函数声明：int epoll_create(int size)

  功能：该函数生成一个epoll专用的文件描述符，其中的参数是指定生成描述符的最大范围；


2、epoll_ctl函数
   函数声明：int epoll_ctl(int epfd, int op, int fd, struct epoll_event *event)

   功能：用于控制某个文件描述符上的事件，可以注册事件，修改事件，删除事件。

   @epfd：由epoll_create生成的epoll专用的文件描述符；

    @op：要进行的操作，EPOLL_CTL_ADD注册、EPOLL_CTL_MOD修改、EPOLL_CTL_DEL删除；

    @fd：关联的文件描述符；

   @event：指向epoll_event的指针；

  成功：0；失败：-1


3、epoll_wait函数
  函数声明:int epoll_wait(int epfd,struct epoll_event * events,int maxevents,int timeout)

  功能：该函数用于轮询I/O事件的发生；

   @epfd：由epoll_create生成的epoll专用的文件描述符；

   @epoll_event：用于回传代处理事件的数组；

   @maxevents：每次能处理的事件数；

   @timeout：等待I/O事件发生的超时值；

  成功：返回发生的事件数；失败：-1

Epoll的ET模式与LT模式
ET（Edge Triggered）与LT（Level Triggered）的主要区别可以从下面的例子看出
eg：
1． 标示管道读者的文件句柄注册到epoll中；
2． 管道写者向管道中写入2KB的数据；
3． 调用epoll_wait可以获得管道读者为已就绪的文件句柄；
4． 管道读者读取1KB的数据
5． 一次epoll_wait调用完成
如果是ET模式，管道中剩余的1KB被挂起，再次调用epoll_wait，得不到管道读者的文件句柄，除非有新的数据写入管道。如果是LT模式，只要管道中有数据可读，每次调用epoll_wait都会触发。

另一点区别就是设为ET模式的文件句柄必须是非阻塞的。
三、 Epoll的实现
Epoll 的源文件在/usr/src/linux/fs/eventpoll.c，在module_init时注册一个文件系统 eventpoll_fs_type，对该文件系统提供两种操作poll和release，所以epoll_create返回的文件句柄可以被poll、 select或者被其它epoll epoll_wait。对epoll的操作主要通过三个系统调用实现：
1． sys_epoll_create
2． sys_epoll_ctl
3． sys_epoll_wait
下面结合源码讲述这三个系统调用。
1.1 long sys_epoll_create (int size)
该系统调用主要分配文件句柄、inode以及file结构。在linux-2.4.32内核中，使用hash保存所有注册到该epoll的文件句柄，在该系统调用中根据size大小分配hash的大小。具体为不小于size，但小于2size的2的某次方。最小为2的9次方（512），最大为2的17次方（128 x 1024）。在linux-2.6.10内核中，使用红黑树保存所有注册到该epoll的文件句柄，size参数未使用。
1.2 long sys_epoll_ctl(int epfd, int op, int fd, struct epoll_event event)
1． 注册句柄 op = EPOLL_CTL_ADD
注册过程主要包括：
A．将fd插入到hash（或rbtree）中，如果原来已经存在返回-EEXIST，
B．给fd注册一个回调函数，该函数会在fd有事件时调用，在该函数中将fd加入到epoll的就绪队列中。
C．检查fd当前是否已经有期望的事件产生。如果有，将其加入到epoll的就绪队列中，唤醒epoll_wait。

2． 修改事件 op = EPOLL_CTL_MOD
修改事件只是将新的事件替换旧的事件，然后检查fd是否有期望的事件。如果有，将其加入到epoll的就绪队列中，唤醒epoll_wait。

3． 删除句柄 op = EPOLL_CTL_DEL
将fd从hash（rbtree）中清除。
1.3 long sys_epoll_wait(int epfd, struct epoll_event events, int maxevents,int timeout)
如果epoll的就绪队列为空，并且timeout非0，挂起当前进程，引起CPU调度。
如果epoll的就绪队列不空，遍历就绪队列。对队列中的每一个节点，获取该文件已触发的事件，判断其中是否有我们期待的事件，如果有，将其对应的epoll_event结构copy到用户events。

revents = epi->file->f_op->poll(epi->file, NULL);
epi->revents = revents & epi->event.events;
if (epi->revents) {
……
copy_to_user;
……
}
需要注意的是，在LT模式下，把符合条件的事件copy到用户空间后，还会把对应的文件重新挂接到就绪队列。所以在LT模式下，如果一次epoll_wait某个socket没有read/write完所有数据，下次epoll_wait还会返回该socket句柄。


epoll原理简述：epoll = 一颗红黑树 + 一张准备就绪句柄链表 + 少量的内核cache
select/poll 每次调用时都要传递你所要监控的所有socket给select/poll系统调用，这意味着需要将用户态的socket列表copy到内核态，如果以万计的句柄会导致每次都要copy几十几百KB的内存到内核态，非常低效。
调用epoll_create后，内核就已经在内核态开始准备帮你存储要监控的句柄了，每次调用epoll_ctl只是在往内核的数据结构里塞入新的socket句柄。
epoll向内核注册了一个文件系统，用于存储上述的被监控socket。当调用epoll_create时，就会在这个虚拟的epoll文件系统里创建一个file结点。当然这个file不是普通文件，它只服务于epoll。
epoll在被内核初始化时（操作系统启动），同时会开辟出epoll自己的内核高速cache区，用于安置每一个我们想监控的socket，这些socket会以红黑树的形式保存在内核cache里，以支持快速的查找、插入、删除。
这个内核高速cache区，就是建立连续的物理内存页，然后在之上建立slab层，简单的说，就是物理上分配好你想要的size的内存对象，每次使用时都是使用空闲的已分配好的对象。
调用epoll_create时，内核除了帮我们在epoll文件系统里建了个file结点，在内核cache里建了个红黑树用于存储以后epoll_ctl传来的socket外，还会再建立一个list链表，用于存储准备就绪的事件，当epoll_wait调用时，仅仅观察这个list链表里有没有数据即可。有数据就返回，没有数据就sleep，等到timeout时间到后即使链表没数据也返回。
通常情况下即使我们要监控百万计的句柄，大多一次也只返回很少量的准备就绪句柄而已，所以，epoll_wait仅需要从内核态copy少量的句柄到用户态而已
这个准备就绪list链表是怎么维护的呢？当我们执行epoll_ctl时，除了把socket放到epoll文件系统里file对象对应的红黑树上之外，还会给内核中断处理程序注册一个回调函数，告诉内核，如果这个句柄的中断到了，就把它放到准备就绪list链表里。
所以，当一个socket上有数据到了，内核在把网卡上的数据copy到内核中后就来把socket插入到准备就绪链表里了。

【转】epoll的两种工作模式：
Epoll的2种工作方式-水平触发（LT）和边缘触发（ET）

假如有这样一个例子：

我们已经把一个用来从管道中读取数据的文件句柄(RFD)添加到epoll描述符
这个时候从管道的另一端被写入了2KB的数据
调用epoll_wait(2)，并且它会返回RFD，说明它已经准备好读取操作
然后我们读取了1KB的数据
调用epoll_wait(2)……
Edge Triggered 工作模式：
如果我们在第1步将RFD添加到epoll描述符的时候使用了EPOLLET标志，那么在第5步调用epoll_wait(2)之后将有可能会挂起，因为剩余的数据还存在于文件的输入缓冲区内，而且数据发出端还在等待一个针对已经发出数据的反馈信息。只有在监视的文件句柄上发生了某个事件的时候 ET 工作模式才会汇报事件。因此在第5步的时候，调用者可能会放弃等待仍在存在于文件输入缓冲区内的剩余数据。在上面的例子中，会有一个事件产生在RFD句柄上，因为在第2步执行了一个写操作，然后，事件将会在第3步被销毁。因为第4步的读取操作没有读空文件输入缓冲区内的数据，因此我们在第5步调用 epoll_wait(2)完成后，是否挂起是不确定的。epoll工作在ET模式的时候，必须使用非阻塞套接口，以避免由于一个文件句柄的阻塞读/阻塞写操作把处理多个文件描述符的任务饿死。最好以下面的方式调用ET模式的epoll接口，在后面会介绍避免可能的缺陷。

基于非阻塞文件句柄
只有当read(2)或者write(2)返回EAGAIN时才需要挂起，等待。但这并不是说每次read()时都需要循环读，直到读到产生一个EAGAIN才认为此次事件处理完成，当read()返回的读到的数据长度小于请求的数据长度时，就可以确定此时缓冲中已没有数据了，也就可以认为此事读事件已处理完成。
Level Triggered 工作模式：
相反的，以LT方式调用epoll接口的时候，它就相当于一个速度比较快的poll(2)，并且无论后面的数据是否被使用，因此他们具有同样的职能。因为即使使用ET模式的epoll，在收到多个chunk的数据的时候仍然会产生多个事件。调用者可以设定EPOLLONESHOT标志，在 epoll_wait(2)收到事件后epoll会与事件关联的文件句柄从epoll描述符中禁止掉。因此当EPOLLONESHOT设定后，使用带有

EPOLL_CTL_MOD标志的epoll_ctl(2)处理文件句柄就成为调用者必须作的事情。

struct epoll_event
typedef union epoll_data {
    void *ptr;
    int fd;
    __uint32_t u32;
    __uint64_t u64;
} epoll_data_t;
 //感兴趣的事件和被触发的事件
struct epoll_event {
    __uint32_t events; /* Epoll events */
    epoll_data_t data; /* User data variable */
};
注意epoll_data是个union而不是struct，两者的区别: http://www.gonglin91.com/cpp-struct-union/

所以当我们在epoll中为一个文件注册一个事件时，如果不需要什么额外的信息，只需要简单在epoll_data.fd = 当前文件的fd

（尽管这并不是必须的，在网上找到的几乎所有的博客都是这样写，呵呵）

但是当我们需要一些额外的数据的时候，就要用上void* ptr;

我们可以自定义一种类型，例如my_data，然后让ptr指向它，以后取出来的时候就能用了；

最后，记住union的特性，它只会保存最后一个被赋值的成员，所以不要data.fd,data.ptr都赋值；通用的做法就是给ptr赋值，fd只是有些时候为了演示或怎么样罢了

所以个人的想法是，其实epoll_data完全可以就是一个void*，不过可能为了一些简单的场景，才设计成union，包含了fd，u32，u64



