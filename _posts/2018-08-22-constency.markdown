---
title: 主从一致性架构优化4种方法
layout: post
category: storage
author: 夏泽民
---
<!-- more -->
需求缘起
大部分互联网的业务都是“读多写少”的场景，数据库层面，读性能往往成为瓶颈。业界通常采用“一主多从，读写分离，冗余多个读库”的数据库架构来提升数据库的读性能。
这种架构的一个潜在缺点是，业务方有可能读取到并不是最新的旧数据：
（1）系统先对DB-master进行了一个写操作，写主库
（2）很短的时间内并发进行了一个读操作，读从库，此时主从同步没有完成，故读取到了一个旧数据
（3）主从同步完成
有没有办法解决或者缓解这类“由于主从延时导致读取到旧数据”的问题呢，这是本文要集中讨论的问题。

方案一（半同步复制）
不一致是因为写完成后，主从同步有一个时间差，假设是500ms，这个时间差有读请求落到从库上产生的。有没有办法做到，等主从同步完成之后，主库上的写请求再返回呢？答案是肯定的，就是大家常说的“半同步复制”semi-sync：
（1）系统先对DB-master进行了一个写操作，写主库
（2）等主从同步完成，写主库的请求才返回
（3）读从库，读到最新的数据（如果读请求先完成，写请求后完成，读取到的是“当时”最新的数据）
方案优点：利用数据库原生功能，比较简单
方案缺点：主库的写请求时延会增长，吞吐量会降低
方案二（强制读主库）
如果不使用“增加从库”的方式来增加提升系统的读性能，完全可以读写都落到主库，这样就不会出现不一致了：


方案优点：“一致性”上不需要进行系统改造
方案缺点：只能通过cache来提升系统的读性能，这里要进行系统改造
方案三（数据库中间件）
如果有了数据库中间件，所有的数据库请求都走中间件，这个主从不一致的问题可以这么解决：
（1）所有的读写都走数据库中间件，通常情况下，写请求路由到主库，读请求路由到从库
（2）记录所有路由到写库的key，在经验主从同步时间窗口内（假设是500ms），如果有读请求访问中间件，此时有可能从库还是旧数据，就把这个key上的读请求路由到主库
（3）经验主从同步时间过完后，对应key的读请求继续路由到从库
方案优点：能保证绝对一致
方案缺点：数据库中间件的成本比较高

方案四（缓存记录写key法）
既然数据库中间件的成本比较高，有没有更低成本的方案来记录某一个库的某一个key上发生了写请求呢？很容易想到使用缓存，当写请求发生的时候：
（1）将某个库上的某个key要发生写操作，记录在cache里，并设置“经验主从同步时间”的cache超时时间，例如500ms
（2）修改数据库

而读请求发生的时候：
（1）先到cache里查看，对应库的对应key有没有相关数据
（2）如果cache hit，有相关数据，说明这个key上刚发生过写操作，此时需要将请求路由到主库读最新的数据
（3）如果cache miss，说明这个key上近期没有发生过写操作，此时将请求路由到从库，继续读写分离
方案优点：相对数据库中间件，成本较低
方案缺点：为了保证“一致性”，引入了一个cache组件，并且读写数据库时都多了一步cache操作
总结
为了解决主从数据库读取旧数据的问题，常用的方案有四种：
（1）半同步复制
（2）强制读主
（3）数据库中间件
（4）缓存记录写key

监控主从同步延迟，同步延迟的检查工作主要从下面两方面着手：
1.一般的做法就是根据Seconds_Behind_Master的值来判断slave的延迟状态。
可以通过监控show slave status\G命令输出的Seconds_Behind_Master参数的值来判断，是否有发生主从延时。

需要监控下面三个参数：
   1）Slave_IO_Running：该参数可作为io_thread的监控项，Yes表示io_thread的和主库连接正常并能实施复制工作，No则说明与主库通讯异常，多数情况是由主从间网络引起的问题；
   2）Slave_SQL_Running：该参数代表sql_thread是否正常，YES表示正常，NO表示执行失败，具体就是语句是否执行通过，常会遇到主键重复或是某个表不存在。
   3）Seconds_Behind_Master：是通过比较sql_thread执行的event的timestamp和io_thread复制好的event的timestamp(简写为ts)进行比较，而得到的这么一个差值；
        NULL—表示io_thread或是sql_thread有任何一个发生故障，也就是该线程的Running状态是No，而非Yes。
        0 — 该值为零，是我们极为渴望看到的情况，表示主从复制良好，可以认为lag不存在。
        正值 — 表示主从已经出现延时，数字越大表示从库落后主库越多。
        负值 — 几乎很少见，我只是听一些资深的DBA说见过，其实，这是一个BUG值，该参数是不支持负值的，也就是不应该出现。
-----------------------------------------------------------------------------------------------------------------------------
Seconds_Behind_Master的计算方式可能带来的问题：
relay-log和主库的bin-log里面的内容完全一样，在记录sql语句的同时会被记录上当时的ts，所以比较参考的值来自于binlog，其实主从没有必要与NTP进行同步，也就是说无需保证主从时钟的一致。其实比较动作真正是发生在io_thread与sql_thread之间，而io_thread才真正与主库有关联，于是，问题就出来了，当主库I/O负载很大或是网络阻塞，io_thread不能及时复制binlog（没有中断，也在复制），而sql_thread一直都能跟上io_thread的脚本，这时Seconds_Behind_Master的值是0，也就是我们认为的无延时，但是，实际上不是，你懂得。这也就是为什么大家要批判用这个参数来监控数据库是否发生延时不准的原因，但是这个值并不是总是不准，如果当io_thread与master网络很好的情况下，那么该值也是很有价值的。之前，提到Seconds_Behind_Master这个参数会有负值出现，我们已经知道该值是io_thread的最近跟新的ts与sql_thread执行到的ts差值，前者始终是大于后者的，唯一的肯能就是某个event的ts发生了错误，比之前的小了，那么当这种情况发生时，负值出现就成为可能。
-----------------------------------------------------------------------------------------------------------------------------

简单来说，就是监控slave同步状态中的：
1）Slave_IO_Running、Slave_SQL_Running状态值，如果都为YES，则表示主从同步；反之，主从不同步。
2）Seconds_Behind_Master的值，如果为0，则表示主从同步不延时，反之同步延时。

2.上面根据Seconds_Behind_Master的值来判断slave的延迟状态，这么做在大部分情况下尚可接受，但其实是并不够准确的。对于Slave延迟状态的监控，还应该做到下面的考虑：
首先，我们先看下slave的状态：
mysql> show slave status\G;

可以看到 Seconds_Behind_Master 的值是 3296，也就是slave至少延迟了 3296 秒。

我们再来看下slave上的2个REPLICATION进程状态：
mysql> show full processlist\G;
可以看到SQL线程一直在执行UPDATE操作，注意到 Time 的值是 3293，看起来像是这个UPDATE操作执行了3293秒，一个普通的SQL而已，肯定不至于需要这么久。
实际上，在REPLICATION进程中，Time 这列的值可能有几种情况：
   1）SQL线程当前执行的binlog（实际上是relay log）中的timestamp和IO线程最新的timestamp的差值，这就是通常大家认为的 Seconds_Behind_Master 值，并不是某个SQL的实际执行耗时；
   2）SQL线程当前如果没有活跃SQL在执行的话，Time值就是SQL线程的idle time；
而IO线程的Time值则是该线程自从启动以来的总时长（多少秒），如果系统时间在IO线程启动后发生修改的话，可能会导致该Time值异常，比如变成负数，或者非常大。
来看下面几个状态：
设置pager，只查看关注的几个status值
mysql> pager cat | egrep -i 'system user|Exec_Master_Log_Pos|Seconds_Behind_Master|Read_Master_Log_Pos';

这是没有活跃SQL的情况，Time值是idle time，并且 Seconds_Behind_Master 为 0
mysql> show processlist; show slave status\G;
检查到此，可以说下如何正确判断slave的延迟情况：
1）首先看 Relay_Master_Log_File 和 Master_Log_File 是否有差异；
2）如果Relay_Master_Log_File 和 Master_Log_File 是一样的话，再来看Exec_Master_Log_Pos 和 Read_Master_Log_Pos 的差异，对比SQL线程比IO线程慢了多少个binlog事件；
3）如果Relay_Master_Log_File 和 Master_Log_File 不一样，那说明延迟可能较大，需要从MASTER上取得binlog status，判断当前的binlog和MASTER上的差距；

因此，相对更加严谨的做法是：
在第三方监控节点上，对MASTER和slave同时发起SHOW BINARY LOGS和SHOW slave STATUS\G的请求，最后判断二者binlog的差异，以及 Exec_Master_Log_Pos 和Read_Master_Log_Pos 的差异。
