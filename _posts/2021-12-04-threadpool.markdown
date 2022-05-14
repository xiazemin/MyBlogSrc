---
title: threadpool
layout: post
category: mysql
author: 夏泽民
---
By default, for every client connection the MySQL server spawns a separate thread which will process all statements for this connection. This is the ‘one-thread-per-connection’ model. 

There are several implementations of the Thread Pool model:
– Commercial: Oracle/MySQL provides thread_pool plugin as part of the Enterprise subscription
– Open Source: The MariaDB thread_pool developed by Vladislav Vaintroub

https://www.percona.com/blog/2014/01/23/percona-server-improve-scalability-percona-thread-pool/
https://www.percona.com/blog/2019/02/25/mysql-challenge-100k-connections/

Step 1. 10,000 connections
This one is very easy, as there is not much to do to achieve this. We can do this with only one client. But you may face the following error on the client side:

FATAL: error 2004: Can't create TCP/IP socket (24)

This is caused by the open file limit, which is also a limit of TCP/IP sockets. This can be fixed by setting   ulimit -n 100000  on the client.

Step 2. 25,000 connections
With 25,000 connections, we hit an error on MySQL side:

Can't create a new thread (errno 11); if you are not out of available memory, you can consult the manual for a possible OS-dependent bug

Step 3. 50,000 connections
This is where we encountered the biggest challenge. At first, trying to get 50,000 connections in sysbench we hit the following error:

FATAL: error 2003: Can't connect to MySQL server on '139.178.82.47' (99)

Error (99) is cryptic and it means: Cannot assign requested address.

It comes from the limit of ports an application can open. By default on my system it is

cat /proc/sys/net/ipv4/ip_local_port_range : 32768   60999

This says there are only 28,231 available ports — 60999 minus 32768 — or the limit of TCP connections you can establish from or to the given IP address.

You can extend this using a wider range, on both the client and the server:

echo 4000 65000 > /proc/sys/net/ipv4/ip_local_port_range

Step 4. 100,000 connections
There is nothing eventful to achieve75k and 100k connections. We just spin up an additional server and start sysbench. For 100,000 connections we need four servers for sysbench, each shows:

Shell
[ 101s] threads: 25000, tps: 0.00, reads: 8033.83, writes: 0.00, response time: 3320.21ms (95%), errors: 0.00, reconnects:  0.00
[ 102s] threads: 25000, tps: 0.00, reads: 8065.02, writes: 0.00, response time: 3405.77ms (95%), errors: 0.00, reconnects:  0.00
1
2
[ 101s] threads: 25000, tps: 0.00, reads: 8033.83, writes: 0.00, response time: 3320.21ms (95%), errors: 0.00, reconnects:  0.00
[ 102s] threads: 25000, tps: 0.00, reads: 8065.02, writes: 0.00, response time: 3405.77ms (95%), errors: 0.00, reconnects:  0.00
So we have the same throughput (8065*4=32260 tps in total) with 3405ms 95% response time.
<!-- more -->
MySQL线程用完，新连接无法连接解决
如果服务器有大量的闲置连接，这样就会白白的浪费内存，且如果一直在累加而不断开的话，就会达到连接上限，报"too many connections”的错误。可通过命令"show process list”查看，若发现后台有大量的sleep线程，此时就需要调整上述参数了。

show variables like '%max_connections%';  --查询当前连接
可以在/etc/my.cnf里面设置数据库的最大连接数
[mysqld]
max_connections = 1000
mysql> show status like 'Threads%';
+-------------------+-------+
| Variable_name     | Value |
+-------------------+-------+
| Threads_cached    | 58    |
| Threads_connected | 57    |   ###这个数值指的是打开的连接数
| Threads_created   | 3676  |
| Threads_running   | 4     |   ###这个数值指的是激活的连接数，这个数值一般远低于connected数值

https://blog.51cto.com/renzhiyuan/1861841

https://www.kancloud.cn/thinkphp/mysql-faq/47447

MySQL Threads Running
Threads_running	MySQL
0 - 10	Normal：几乎所有硬件都没问题
10 - 30	Busy：大多数硬件通常都可以，因为服务器多核
30 - 50	High：很少有工作负载需要运行这么多线程。它可以短期爆发（<5min），但如果持续时间较长，则响应时间很可能很慢
50 - 100	Overloaded：某些硬件可以处理此问题，但是不能期望在此范围内成功运行。对于我们的本地部署硬件而言，此范围内的瞬时突发（<5s）通常是可以的。
> 100	Failing：在极少数情况下，MySQL可以运行大于100个线程，但在此范围内可能会失败
    建议指导值：

Threads_running < 50
1:1000 Threads_running：QPS

让我们换个角度来阐明一个重要的问题：MySQL线程是一个数据库连接。Threads_running是活动查询的数据库连接数。请记住，每个应用程序实例都有其自己的数据库连接池，这一点很重要。因此，最大可能的连接（线程）为：

    连接池大小为100是合理的，但是如果将应用程序部署到5个应用程序实例，则可能有500个数据库连接。通常是这样：应用程序通常具有数百个空闲数据库连接，这是连接池的用途。（对连接的线程，MySQL也有一个尺度）直到同时运行过多的连接（线程），这才成为问题。
    一个应用程序不止一次被扩展（即部署更多的应用程序实例）以处理更多的请求，但这样做会使数据库过载，运行的线程太多。没有快速或简单的解决方案来解决这种类型的数据库性能限制。原因很简单：如果您希望MySQL在相同的时间（每秒）内执行更多的工作（查询），则每个工作必须花费较少的时间，否则无法进行计算。如果每个查询花费100毫秒，则执行速度不可能超过10QPS。
    
    https://blog.csdn.net/zsx0728/article/details/114536258
    

