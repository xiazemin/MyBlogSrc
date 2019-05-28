---
title: SO_REUSEPORT 惊群
layout: post
category: linux
author: 夏泽民
---
单个进程监听多个端口
单个进程创建多个 socket 绑定不同的端口，TCP, UDP 都行

多个进程监听同一个端口(multiple processes listen on same port)
通过 fork 创建子进程的方式可以实现

从Linux 3.9内核版本之后Linux网络协议栈开始支持SO_REUSEPORT套接字选项，这个新的选项允许一个主机上的多个套接字绑定到同一个端口上，它的目的是提高运行在多核CPU上的多线程网络服务应用的处理性能。

他的使用也非常简单，如果多个进程或者线程都设置了下面这个选项，则他们可以同时绑定到同一个端口上：

int sfd = socket(domain, socktype, 0);

int optval = 1;
setsockopt(sfd, SOL_SOCKET, SO_REUSEPORT, &optval, sizeof(optval));

bind(sfd, (struct sockaddr *) &addr, addrlen);
只要第一个进程在绑定端口时设置了这个选项，则其他进程也可以通过设置这个选项来绑定到同一个端口上。 要求第一个进程必须设置SO_REUSEPORT这个选项的原因是防止端口劫持–一些流氓进程通过绑定正在被使用的端口上，来获取其进程接收到的连接请求和数据。为了防止其他不必要的进程通过SO_REUSEPORT选项劫持端口，所有之后绑定这个端口的进程都需要设置和第一个进程相同的user ID。
<!-- more -->
TCP和UDP都可以使用SO_REUSEPORT选项。对于TCP它允许多个套接字监听同一个端口号，这样每个线程都可以调用accept()来处理连接， 避免了传统多线程服务中通常使用一个单一进程处理连接请求，而这个单一进程很可能会成为整个系统的瓶颈。 传统多线程服务中的另一种处理方法是多个线程或者进程对同一个套接字循环调用accept()函数处理连接请求，形式如下：

while (1) {
    new_fd = accept(...);
    process_connection(new_fd);
}
这种处理方式也会有一个问题：多个线程之间不能均衡的处理请求，有些线程处理了大量请求，有些线程处理了少量请求，这种不均衡会降低多核CPU的利用率。 而SO_REUSEPORT会更加均衡的分发请求到不同线程或者进程上。

SO_REUSEPORT选项分发数据包的方法是计算对端IP、端口加上本地IP、端口这四个值的哈希值，通过这个哈希值将数据包分发到不同进程上。 这样就可以保证同一个连接的数据包都被分发到同一个进程中去处理。

SO_REUSEPORT套接字选项在内核中的实现
这里只看UDP协议的实现， 当设置了SO_REUSEPORT套接字选项之后，绑定在同一个端口号的套接字在内核中会形成一个数组，保存在sock_reuseport结构体中， 在调用bind()函数时，会调用到/net/core/sock_reuseport.c文件中的reuseport_add_sock函数，此函数用来将当前套接字添加到数组中。

/* /include/net/sock_reuseport.h */
struct sock_reuseport {
	struct rcu_head		rcu;

	u16			max_socks;	/* length of socks */
	u16			num_socks;	/* elements in socks */
	struct bpf_prog __rcu	*prog;		/* optional BPF sock selector */
	struct sock		*socks[0];	/* 绑定在同一个端口号的套接字指针数组 */
};
在这篇文章中Linux协议栈–UDP协议的发送和接收我们说过当UDP数据到达IP层之后，会调用__udp4_lib_rcv函数将数据包存放到UDP的数据接收缓冲区中。在存放之前会调用__udp4_lib_lookup_skb函数找到这个数据包对应的sock。最终会调用__udp4_lib_lookup函数进行实际的查找工作：

struct sock *__udp4_lib_lookup(struct net *net, __be32 saddr,
		__be16 sport, __be32 daddr, __be16 dport, int dif,
		int sdif, struct udp_table *udptable, struct sk_buff *skb)
{
    ...
begin:
	result = NULL;
	badness = 0;
    /* 遍历链表 */
	sk_for_each_rcu(sk, &hslot->head) {
        /* 根据五元组等信息来进行匹配 */
		score = compute_score(sk, net, saddr, sport,
				      daddr, hnum, dif, sdif, exact_dif);
		if (score > badness) {
            /* 匹配到之后，判断是否设置了 SO_REUSEPORT 选项 */
			if (sk->sk_reuseport) {
                /* 根据源端口、IP和接收端口、IP这四个值计算一个哈希值 */
				hash = udp_ehashfn(net, daddr, hnum,
						   saddr, sport);
                /* 根据这个哈希值，将数据包分发到对应的sock上 */
				result = reuseport_select_sock(sk, hash, skb,
							sizeof(struct udphdr));
				if (result)
					return result;
			}
			result = sk;
			badness = score;
		}
	}
	return result;
}
找到对应的sock之后，调用udp_queue_rcv_skb函数将数据包存放到此套接字的缓冲区中，之后调用sk->sk_data_ready(sk)函数指针，此函数指针在创建套接字的时候初始化为sock_def_readable函数。这个函数会将对应的进程唤醒，来接收数据包。

static void sock_def_readable(struct sock *sk)
{
	struct socket_wq *wq;

	rcu_read_lock();
	wq = rcu_dereference(sk->sk_wq);
	if (skwq_has_sleeper(wq))
		wake_up_interruptible_sync_poll(&wq->wait, EPOLLIN | EPOLLPRI |
						EPOLLRDNORM | EPOLLRDBAND);
	sk_wake_async(sk, SOCK_WAKE_WAITD, POLL_IN);
	rcu_read_unlock();
}

惊群现象就是当多个进程或线程在同时阻塞等待同一个事件时，如果这个事件发生，会唤醒所有的进程，但最终只可能有一个进程/线程对该事件进行处理，其他进程/线程会在失败后重新休眠，这种性能浪费就是惊群。

主进程执行 socket()+bind()+listen() 后，fork() 多个子进程，每个子进程都通过 accept() 循环处理这个 socket；此时，每个进程都阻塞在 accpet() 调用上，当一个新连接到来时，所有的进程都会被唤醒，但其中只有一个进程会 accept() 成功，其余皆失败，重新休眠。这就是 accept 惊群。

如果只用一个进程去 accept 新连接，并通过消息队列等同步方式使其他子进程处理这些新建的连接，那么将会造成效率低下；因为这个进程只能用来 accept 连接，该进程可能会造成瓶颈。

epoll()
另外还有一个是关于 epoll_wait() 的，目前来仍然存在惊群现象。

主进程仍执行 socket()+bind()+listen() 后，将该 socket 加入到 epoll 中，然后 fork 出多个子进程，每个进程都阻塞在 epoll_wait() 上，如果有事件到来，则判断该事件是否是该 socket 上的事件，如果是，说明有新的连接到来了，则进行 accept 操作。
#include <netdb.h>
#include <stdio.h>
#include <fcntl.h>
#include <string.h>
#include <unistd.h>
#include <stdlib.h>
#include <sys/wait.h>
#include <sys/epoll.h>
#include <sys/types.h>
#include <sys/socket.h>

#define PROCESS_NUM 10
#define MAXEVENTS 64

int main (int argc, char *argv[])
{
    int sfd, efd;
    int flags;
    int n, i, k;
    struct epoll_event event;
    struct epoll_event *events;
    struct sockaddr_in serveraddr;

    sfd = socket(PF_INET, SOCK_STREAM, 0);
    serveraddr.sin_family = AF_INET;
    serveraddr.sin_addr.s_addr = htonl(INADDR_ANY);
    serveraddr.sin_port = htons(atoi("1234"));
    bind(sfd, (struct sockaddr*)&serveraddr, sizeof(serveraddr));

    flags = fcntl (sfd, F_GETFL, 0);
    flags |= O_NONBLOCK;
    fcntl (sfd, F_SETFL, flags);

    if (listen(sfd, SOMAXCONN) < 0) {
        perror ("listen");
        exit(EXIT_SUCCESS);
    }

    if ((efd = epoll_create(MAXEVENTS)) < 0) {
        perror("epoll_create");
        exit(EXIT_SUCCESS);
    }

    event.data.fd = sfd;
    event.events = EPOLLIN; // | EPOLLET;
    if (epoll_ctl(efd, EPOLL_CTL_ADD, sfd, &event) < 0) {
        perror("epoll_ctl");
        exit(EXIT_SUCCESS);
    }

    /* Buffer where events are returned */
    events = (struct epoll_event*)calloc(MAXEVENTS, sizeof event);

    for(k = 0; k < PROCESS_NUM; k++) {
        if (fork() == 0) { /* children process */
            while (1) {    /* The event loop */
                n = epoll_wait(efd, events, MAXEVENTS, -1);
                printf("process #%d return from epoll_wait!\n", getpid());
                sleep(2);  /* sleep here is very important!*/
                for (i = 0; i < n; i++) {
                    if ((events[i].events & EPOLLERR) ||
                        (events[i].events & EPOLLHUP) ||
                        (!(events[i].events & EPOLLIN))) {
                        /* An error has occured on this fd, or the socket is not
                         * ready for reading (why were we notified then?)
                         */
                        fprintf (stderr, "epoll error\n");
                        close (events[i].data.fd);
                        continue;
                    } else if (sfd == events[i].data.fd) {
                        /* We have a notification on the listening socket, which
                         * means one or more incoming connections.
                         */
                        struct sockaddr in_addr;
                        socklen_t in_len;
                        int infd;
                        //char hbuf[NI_MAXHOST], sbuf[NI_MAXSERV];

                        in_len = sizeof in_addr;

                        infd = accept(sfd, &in_addr, &in_len);
                        if (infd == -1) {
                            printf("process %d accept failed!\n", getpid());
                            break;
                        }
                        printf("process %d accept successed!\n", getpid());

                        /* Make the incoming socket non-blocking and add it to the
                        list of fds to monitor. */
                        close(infd);
                    }
                }
            }
        }
    }
    int status;
    wait(&status);
    free (events);
    close (sfd);
    return EXIT_SUCCESS;
}
注意：上述的处理中添加了 sleep() 函数，实际上，如果很快处理完了这个 accept() 请求，那么其余进程可能还没有来得及被唤醒，内核队列上已经没有这个事件，无需唤醒其他进程。

那么，为什么只解决了 accept() 的惊群问题，而没有解决 epoll() 的？

当接收到一个报文后，显然只能由一个进程处理 (accept)；而 epoll() 却不同，因为内核不知道对应的触发事件具体由哪些进程处理，那么只能是唤醒所有的进程，然后由不同的进程进行处理。

nginx 的每个 worker 进程都会在函数 ngx_process_events_and_timers() 中处理不同的事件，然后通过 ngx_process_events() 封装了不同的事件处理机制，在 Linux 上默认采用 epoll_wait()。

主要在 ngx_process_events_and_timers() 函数中解决惊群现象。

void ngx_process_events_and_timers(ngx_cycle_t *cycle)
{
    ... ...
    // 是否通过对accept加锁来解决惊群问题，需要工作线程数>1且配置文件打开accetp_mutex
    if (ngx_use_accept_mutex) {
        // 超过配置文件中最大连接数的7/8时，该值大于0，此时满负荷不会再处理新连接，简单负载均衡
        if (ngx_accept_disabled > 0) {
            ngx_accept_disabled--;
        } else {
            // 多个worker仅有一个可以得到这把锁。获取锁不会阻塞过程，而是立刻返回，获取成功的话
            // ngx_accept_mutex_held被置为1。拿到锁意味着监听句柄被放到本进程的epoll中了，如果
            // 没有拿到锁，则监听句柄会被从epoll中取出。
            if (ngx_trylock_accept_mutex(cycle) == NGX_ERROR) {
                return;
            }
            if (ngx_accept_mutex_held) {
                // 此时意味着ngx_process_events()函数中，任何事件都将延后处理，会把accept事件放到
                // ngx_posted_accept_events链表中，epollin|epollout事件都放到ngx_posted_events链表中
                flags |= NGX_POST_EVENTS;
            } else {
                // 拿不到锁，也就不会处理监听的句柄，这个timer实际是传给epoll_wait的超时时间，修改
                // 为最大ngx_accept_mutex_delay意味着epoll_wait更短的超时返回，以免新连接长时间没有得到处理
                if (timer == NGX_TIMER_INFINITE || timer > ngx_accept_mutex_delay) {
                    timer = ngx_accept_mutex_delay;
                }
            }
        }
    }
    ... ...
    (void) ngx_process_events(cycle, timer, flags);   // 实际调用ngx_epoll_process_events函数开始处理
    ... ...
    if (ngx_posted_accept_events) { //如果ngx_posted_accept_events链表有数据，就开始accept建立新连接
        ngx_event_process_posted(cycle, &ngx_posted_accept_events);
    }

    if (ngx_accept_mutex_held) { //释放锁后再处理下面的EPOLLIN EPOLLOUT请求
        ngx_shmtx_unlock(&ngx_accept_mutex);
    }

    if (delta) {
        ngx_event_expire_timers();
    }

    ngx_log_debug1(NGX_LOG_DEBUG_EVENT, cycle->log, 0, "posted events %p", ngx_posted_events);
	// 然后再处理正常的数据读写请求。因为这些请求耗时久，所以在ngx_process_events里NGX_POST_EVENTS标
    // 志将事件都放入ngx_posted_events链表中，延迟到锁释放了再处理。
}
关于 ngx_use_accept_mutex、ngx_accept_disabled 的修改可以直接 grep 查看。

SO_REUSEPORT
Linux 内核的 3.9 版本带来了 SO_REUSEPORT 特性，该特性支持多个进程或者线程绑定到同一端口，提高服务器程序的性能，允许多个套接字 bind() 以及 listen() 同一个 TCP 或者 UDP 端口，并且在内核层面实现负载均衡。

在未开启 SO_REUSEPORT 时，由一个监听 socket 将新接收的链接请求交给各个 worker 处理。

在使用 SO_REUSEPORT 后，多个进程可以同时监听同一个 IP:Port ，然后由内核决定将新链接发送给那个进程，显然会降低各个 worker 接收新链接时锁竞争。

其实在Linux2.6版本以后，内核内核已经解决了accept()函数的“惊群”问题，大概的处理方式就是，当内核接收到一个客户连接后，只会唤醒等待队列上的第一个进程或线程。所以，如果服务器采用accept阻塞调用方式，在最新的Linux系统上，已经没有“惊群”的问题了。

但是，对于实际工程中常见的服务器程序，大都使用select、poll或epoll机制，此时，服务器不是阻塞在accept，而是阻塞在select、poll或epoll_wait，这种情况下的“惊群”仍然需要考虑。

在早期的Linux版本中，内核对于阻塞在epoll_wait的进程，也是采用全部唤醒的机制，所以存在和accept相似的“惊群”问题。新版本的的解决方案也是只会唤醒等待队列上的第一个进程或线程，所以，新版本Linux 部分的解决了epoll的“惊群”问题。所谓部分的解决，意思就是：对于部分特殊场景，使用epoll机制，已经不存在“惊群”的问题了，但是对于大多数场景，epoll机制仍然存在“惊群”。

epoll存在惊群的场景如下：在worker保持工作的状态下，都会被唤醒
Nginx中使用mutex互斥锁解决这个问题，具体措施有使用全局互斥锁，每个子进程在epoll_wait()之前先去申请锁，申请到则继续处理，获取不到则等待，并设置了一个负载均衡的算法（当某一个子进程的任务量达到总设置量的7/8时，则不会再尝试去申请锁）来均衡各个进程的任务量。

通过五元组（协议、源地址、源端口、目的地址、目的端口）可以 唯一定位 一个连接。通过 socket、bind、connect 三个系统调用可以为一个 socket 分配五元组。

socket 中通过 SOCK_STREAM(TCP)、SOCK_DGRAM(UDP) 指定协议。
bind 来绑定源地址、源端口。bind 这步也可以省略，如果省略，内核会为该 socket 分配一个源地址、源端口。另外，如果 bind 的源地址为 0.0.0.0，内核也会从本地接口中分配一个作为源地址。如果 bind 的源端口是 0，内核也会自动分配源端口。
connect 来指定目的地址、目的端口。有同学会问，UDP 是无连接状态的，也能 connect 吗？当然可以的，详情请 man connect ，无非 UDP 在 connect 的时候没有三次握手嘛。如果 UDP 不通过 connect 来指定目的地址、目的端口，那发送数据包时，就必须使用 sendto 而不是 send ，在 sendto 的时候指定目的地址、目的端口。
为了能够通过五元组唯一定位一个连接，内核在给连接分配源地址、源端口的时候，必须使不同的连接具有不同的五元组。因此，在默认情况下，两个不同的 socket 不能 bind 到相同的源地址和源端口。

SO_REUSEADDR
在有些情况下，这个默认情况下的限制: 两个不同的 socket 不能 bind 到相同的源地址和源端口 ，带来很大的困扰。
TCP 的 TIME_WAIT 状态 时间过长，造成新 socket 无法复用这个端口，即使可以确定这个连接可以销毁。完全是 拉完屎还占着茅坑 。这个问题在重启 TCP Server 时更为严重。
0.0.0.0 和本地接口地址，也会被认为是相同的源地址，从而 bind 失败。
在 unix 系统中，SO_REUSEADDR 就是为了解决这些困扰而生。
SO_REUSEPORT
Linux 在内核 3.9 中添加了新的 socket option SO_REUSEPORT 。
如果在 bind 系统调用前，指定了 SO_REUSEPORT ，多个 socket 便可以 bind 到相同的源地址、源端口，比起 SO_REUSEADDR 更强大、更劲爆，有木有。不仅如此，还添加了权限保护，为了防止 端口劫持 ，在第一个 socket bind 成功后，后续的 socket bind 的用户必须或者是 root，或者跟第一个 socket 用户一致。

使用 SO_REUSEPORT 杜绝 accept 惊群
这里是本文的重点。以前的 TCP Server 开发中，为了充分利用多核的性能，所以在多进程中监听同一个端口。在没有 SO_REUSEPORT 的时代，可以通过 fork 来实现。父进程绑定一个端口监听 socket ，然后 fork 出多个子进程，子进程们开始循环 accept 这个 socket 。但是会带来一个问题：如果有新连接建立，哪个进程会被唤醒且能够成功 accept 呢？在 Linux 内核版本 2.6.18 以前，所有监听进程都会被唤醒，但是只有一个进程 accept 成功，其余失败。这种现象就是所谓的 惊群效应 。其实在 2.6.18 以后，这个问题得到修复，仅有一个进程被唤醒并 accept 成功。

但是，现在的 TCP Server，一般都是 多进程+多路IO复用(epoll) 的并发模型，比如我们常用的 nginx 。如果使用 epoll 去监听 accept socket fd 的读事件，当有新连接建立时，所有进程都会被触发。因为由于 fork 文件描述符继承的缘故，所有进程中的 accept socket fd 是相同的。惊群效应依然存在。nginx 也必然存在这个问题，nginx 为了解决问题，并且保证各个 worker 之前 accept 连接数的均衡，费了很大的力气。

有了 SO_REUSEPORT ，解决 多进程+多路IO复用(epoll) 并发模型 accept 惊群问题，就简单、高效很多。我们不需要通过 fork 的形式，让多进程监听同一个端口。只需要在各个进程中， 独自的 监听指定的端口，当然在监听前，我们需要为监听 socket 指定 SO_REUSEPORT ，否则会报错啦。由于没有采用 fork 的形式，各个进程中的 accept socket fd 不一样，加之有新连接建立时，内核只会唤醒一个进程来 accept，并且保证唤醒的 均衡性，因此使用 epoll 监听读事件，就不会触发所有啦。也有牛人为 nginx 提了 patch ，使用 SO_REUSEPORT 来杜绝 accept 惊群，并且还能够保证 worker 之间的均衡性哦。

