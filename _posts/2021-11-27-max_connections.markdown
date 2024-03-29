---
title: max_connections profiling和explain
layout: post
category: mysql
author: 夏泽民
---
-- 数据库最大连接数

show variables like '%max_connections%';

-- 已使用连接数
show global status like 'Max_used_connections';

-- 连接线程数
show status like 'Threads%';

-- 连接详细信息
show FULL PROCESSLIST;
 

二、修改连接数

1、临时修改，重启就失效

SET GLOBAL max_connections=2000;

2、永久修改

可以在/etc/my.cnf里面设置数据库的最大连接数

[mysqld]

set-variable=max_connections=500  #最大连接数
https://blog.csdn.net/zsd_31/article/details/90812355
<!-- more -->
在我们的系统出现性能问题时，往往避不开调查各种类型 Lock Wait，如Row Lock Wait、Page Lock Wait、Page IO Latch Wait等。从中找出可能的异常等待，为性能优化做一定的参考 

https://www.cnblogs.com/SameZhao/p/4719146.html

Linux中我们常用mpstat、vmstat、iostat、sar和top来查看系统的性能状态。

§ mpstat： mpstat是Multiprocessor Statistics的缩写，是实时系统监控工具。其报告为CPU的一些统计信息，这些信息存放在/proc/stat文件中。在多CPUs系统里，其不但能查看所有CPU的平均状况信息，而且能够查看特定CPU的信息。mpstat最大的特点是可以查看多核心cpu中每个计算核心的统计数据，而类似工具vmstat只能查看系统整体cpu情况。

§ vmstat：vmstat命令是最常见的Linux/Unix监控工具，可以展现给定时间间隔的服务器的状态值，包括服务器的CPU使用率，内存使用，虚拟内存交换情况，IO读写情况。这个命令是我查看Linux/Unix最喜爱的命令，一个是Linux/Unix都支持，二是相比top，我可以看到整个机器的CPU、内存、IO的使用情况，而不是单单看到各个进程的CPU使用率和内存使用率(使用场景不一样)。

§ iostat: 主要用于监控系统设备的IO负载情况，iostat首次运行时显示自系统启动开始的各项统计信息，之后运行iostat将显示自上次运行该命令以后的统计信息。用户可以通过指定统计的次数和时间来获得所需的统计信息。

§ sar： sar（System Activity Reporter系统活动情况报告）是目前 Linux 上最为全面的系统性能分析工具之一，可以从多方面对系统的活动进行报告，包括：文件的读写情况、系统调用的使用情况、磁盘I/O、CPU效率、内存使用状况、进程活动及IPC有关的活动等。

§ top：top命令是Linux下常用的性能分析工具，能够实时显示系统中各个进程的资源占用状况，类似于Windows的任务管理器。top显示系统当前的进程和其他状况，是一个动态显示过程,即可以通过用户按键来不断刷新当前状态.如果在前台执行该命令，它将独占前台，直到用户终止该程序为止。比较准确的说，top命令提供了实时的对系统处理器的状态监视。它将显示系统中CPU最“敏感”的任务列表。该命令可以按CPU使用。内存使用和执行时间对任务进行排序；而且该命令的很多特性都可以通过交互式命令或者在个人定制文件中进行设定。

https://www.jianshu.com/p/3c79039e82aa

MySQL5.0.37版本以上支持了Profiling – 官方手册。此工具可用来查询 SQL 会执行多少时间，System lock和Table lock 花多少时间等等，对定位一条语句的 I/O消耗和CPU消耗 非常重要。
查看profiling；
　　select @@profiling;
启动profiling:
set @@profiling=1 
关闭profiling ：
set @@profiling=0;

　　sql语句；
1.查看profile记录
show profiles;

Duration:我需要时间；
query：执行的sql语句；
2.查看详情：
show profile for query 2；
 
3.查看cup和io情况
show profile cpu,block io for query 2;

1.id：一组数字，操作顺序，如果id相同，则执行顺序由上至下，如果是子查询，id的序号递增，值越大优先级越高，越先被执行；

2.select_type:表示每个字句的类型，简单还是复杂，取值如下；

　　a>simple :简单查询，无子查询或union等；

　　b>primary:查询中若包含复杂的子部分，最外层则被标记为primary；

　　c>subquery：在select或where中若包含子查询，则该子查询被标记为subquery；

　　d>derived：from中包含子查询，被标记为derived；

　　e>union：若select出现在union之后，则被标记为union；

　　f>union result：从union表中获取结果的select将被标记为union result；

3.table 查询的数据库表名称

4.type 联合查询使用的类型

　　all :全表扫描

　　index：全表扫描,只是扫描表的时候按照索引次序 进行而不是行。主要优点就是避免了排序, 但是开销仍然非常大。

range:索引范围扫描

　　ref:非唯一性索引扫描，交返回匹配单独值的所有行，常见于使用非唯一性索引或唯一性索引的非唯一前缀进行的查找。

　　eq_ref:唯一性索引扫描

　　const、system：当mysql对查询的某部分进行优化，并转换为一个常量时。如将主键置于where列表中，mysql就能将该查询转换为一个常量。system是const的特例，当查询的表只有一行的情况下，即可使用system。
　　
　　https://zhuanlan.zhihu.com/p/114440123
　　
　　开启慢查询日志
在 MySQL中，提供了慢查询查询日志，基于性能方面的考虑，该配置默认为OFF(关闭) 状态。那么如何开启慢日志查询呢？其步骤如下：

在 MySQL 中，慢查询日志默认为OFF状态，通过如下命令进行查看：
mysql> show variables like "slow_query_log";
通过如下命令进行设置为 ON 状态：
set global slow_query_log = "ON";

mysql> set global long_query_time = 5;
当设置值小于0时，默认为 0。

https://zhuanlan.zhihu.com/p/112307303

比较的五款常用工具
 
mysqldumpslow, mysqlsla, myprofi, mysql-explain-slow-log, mysqllogfilter

mysqldumpslow, mysql官方提供的慢查询日志分析工具. 输出图表如下


主要功能是, 统计不同慢sql的
出现次数(Count), 
执行最长时间(Time), 
累计总耗费时间(Time), 
等待锁的时间(Lock), 
发送给客户端的行总数(Rows), 
扫描的行总数(Rows), 
用户以及sql语句本身(抽象了一下格式, 比如 limit 1, 20 用 limit N,N 表示).
讲一下有用的参数： 
-s 排序选项：c 查询次数 r 返回记录行数 t 查询时间 
-t 只显示top n条查询 
mysqldumpslow -s r -t 10 slow.log 

 

mysqlsla, hackmysql.com推出的一款日志分析工具(该网站还维护了 mysqlreport, mysqlidxchk 等比较实用的mysql工具)


整体来说, 功能非常强大. 数据报表,非常有利于分析慢查询的原因, 包括执行频率, 数据量, 查询消耗等.
 

 
格式说明如下:
总查询次数 (queries total), 去重后的sql数量 (unique)
输出报表的内容排序(sorted by)
最重大的慢sql统计信息, 包括 平均执行时间, 等待锁时间, 结果行的总数, 扫描的行总数.

 
Count, sql的执行次数及占总的slow log数量的百分比.
Time, 执行时间, 包括总时间, 平均时间, 最小, 最大时间, 时间占到总慢sql时间的百分比.
95% of Time, 去除最快和最慢的sql, 覆盖率占95%的sql的执行时间.
Lock Time, 等待锁的时间.
95% of Lock , 95%的慢sql等待锁时间.
Rows sent, 结果行统计数量, 包括平均, 最小, 最大数量.
Rows examined, 扫描的行数量.
Database, 属于哪个数据库
Users, 哪个用户,IP, 占到所有用户执行的sql百分比

 
Query abstract, 抽象后的sql语句
Query sample, sql语句

https://developer.aliyun.com/article/435883

MySQL日志文件系统的组成
   a、错误日志：记录启动、运行或停止mysqld时出现的问题。
   b、通用日志：记录建立的客户端连接和执行的语句。
   c、更新日志：记录更改数据的语句。该日志在MySQL 5.1中已不再使用。
   d、二进制日志：记录所有更改数据的语句。还用于复制。
   e、慢查询日志：记录所有执行时间超过long_query_time秒的所有查询或不使用索引的查询。
   f、Innodb日志：innodb redo log

https://blog.csdn.net/enweitech/article/details/80239189


