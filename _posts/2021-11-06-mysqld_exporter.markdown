---
title: prometheus监控插件mysqld_exporter redis_exporter
layout: post
category: mysql
author: 夏泽民
---
https://github.com/prometheus/mysqld_exporter
https://www.cnblogs.com/xiaobaozi-95/p/11453734.html
https://blog.csdn.net/allway2/article/details/106986309/

https://www.cnblogs.com/klvchen/p/10062754.html
https://www.cnblogs.com/xiangsikai/p/11289675.html

redis 监控插件
https://github.com/oliver006/redis_exporter
<!-- more -->
Prometheus 已经成为云原生应用监控行业的标准，在很多流行的监控系统中都已经实现了 Prometheus的监控接口，例如 etcd、Kubernetes、CoreDNS等，它们可以直接被Prometheus监控，但大多数监控对象都没办法直接提供监控接口，主要原因有：
（1）很多系统在Prometheus诞生前的很多年就已发布，例如MySQL、Redis等；
（2）它们本身不支持 HTTP 接口，例如对于硬件性能指标，操作系统并没有原生的HTTP接口可以获取；
（3）考虑到安全性、稳定性及代码耦合等因素的影响，软件作者并不愿意将监控代码加入现有代码中。
这些都导致无法通过一个规范解决所有监控问题。在此背景之下，Exporter 应运而生。Exporter 是一个采集监控数据并通过 Prometheus 监控规范对外提供数据的组件。除了官方实现的Exporter如Node Exporter、HAProxy Exporter、MySQLserver Exporter，还有很多第三方实现如Redis Exporter和RabbitMQ Exporter等


Exporter获取监控数据的方式
Exporter 主要通过被监控对象提供的监控相关的接口获取监控数据，主要有如下几种方式：
（1）HTTP/HTTPS方式。例如 RabbitMQ exporter通过 RabbitMQ的 HTTPS接口获取监控数据。
（2）TCP方式。例如Redis exporter通过Redis提供的系统监控相关命令获取监控指标，MySQL server exporter通过MySQL开放的监控相关的表获取监控指标。
（3）本地文件方式。例如Node exporter通过读取proc文件系统下的文件，计算得出整个操作系统的状态。
（4）标准协议方式。

https://www.jianshu.com/p/3029e58f7141

MySQL监控指标及采集方法
一、用户、连接类
1、查看每个客户端IP过来的连接消耗资源情况。

select * from sys.host_summary;



2、查看每个用户消耗资源情况

select * from sys.user_summary;



3、查看当前连接情况（有多少连接就应该有多少行）

select host,current_connections,statements from sys.host_summary;



4、查看当前正在执行的SQL

和执行show full processlist的结果差不多

select conn_id,pid,user,db,command,current_statement,last_statement,time,lock_latency from sys.session



二、SQL 和io类
1、查看发生IO请求前5名的文件。

select * from sys.io_global_by_file_by_bytes order by total limit 5;



三、buffer pool 、内存
1、查看总共分配了多少内存

select * from sys.memory_global_total;select * from sys.memory_global_by_current_bytes;



2、每个库（database）占用多少buffer pool

select * from sys.innodb_buffer_stats_by_schema order by allocated desc;



pages是指在buffer pool中的page数量；pages_old指在LUR 列表中处于后37%位置的page。

当出现buffer page不够用时，就会征用这些page所占的空间。37%是默认位置，具体可以自定义。

3、统计每张表具体在InnoDB中具体的情况，比如占多少页？

注意和前面的pages的总数都是相等的，也可以借用sum（pages）运算验证一下。

select * from sys.innodb_buffer_stats_by_table;



4、查询每个连接分配了多少内存

利用session表和memory_by_thread_by_current_bytes分配表进行关联查询。

SELECT b.USER, current_count_used, current_allocated, current_avg_alloc, current_max_alloc, total_allocated, current_statement FROM sys.memory_by_thread_by_current_bytes a, sys.SESSION b WHERE a.thread_id = b.thd_id;



四、字段、索引、锁
1、查看表自增字段最大值和当前值，有时候做数据增长的监控，可以作为参考。

select * from sys.schema_auto_increment_columns;



2、MySQL索引使用情况统计

select * from sys.schema_index_statistics order by rows_selected desc;



3、MySQL中有哪些冗余索引和无用索引

若库中展示没有冗余索引，则没有数据；当有联合索引idx_abc(a,b,c)和idx_a(a)，那么idx_a就算冗余索引了。

select * from sys.schema_redundant_indexes;



4、查看INNODB 锁信息

在未来的版本将被移除，可以采用其他方式

select * from sys.innodb_lock_waits



5、查看库级别的锁信息，这个需要先打开MDL锁的监控：

--打开MDL锁监控update performance_schema.setup_instruments set enabled='YES',TIMED='YES' where name='wait/lock/metadata/sql/mdl';select * from sys.schema_table_lock_waits;



五、线程类
1、MySQL内部有多个线程在运行，线程类型及数量

select user,count(*) from sys.`processlist` group by user;



六、主键自增
查看MySQL自增id的使用情况

SELECT table_schema, table_name, ENGINE, Auto_increment FROM information_schema.TABLES WHERE TABLE_SCHEMA NOT IN ( "INFORMATION_SCHEMA", "PERFORMANCE_SCHEMA", "MYSQL", "SYS" )



背景：线上生产环境MySQL的架构是一主双从，为了更好了解MySQL集群运行状况，我们需要对以下指标进行监控！

一、对数据库服务可用性进行监控

思路：

1.1 通过测试账号ping命令返回的信息判断数据库可以通过网络连接

[root@host-39-108-217-12 scripts]# /usr/bin/mysqladmin -uroot -p123456 ping

mysqld is alive



1.2 确认数据库是否可读写

a.检查数据库的read_only参数是否为off

[root@host-47-106-141-17 scripts]# mysql -uroot -p123456 -P3306 -e "show global variables like 'read_only'" | grep read_only

read_only OFF



b.执行简单的数据库查询，如：select @@version;

[root@host-47-106-141-17 scripts]# mysql -uroot -p123456 -P3306 -e "select @@version" | grep MariaDB

5.5.56-MariaDB



二、对数据库性能进行监控

2.1 监控数据库连接数可用性

a.数据库最大连接数

[root@host-47-106-141-17 scripts]# mysql -uroot -p123456 -e "show variables like 'max_connections'"

+-----------------+-------+

| Variable_name | Value |

+-----------------+-------+

| max_connections | 151 |

+-----------------+-------+



b.数据库当前打开的连接数

[root@host-47-106-141-17 scripts]# mysqladmin -uroot -p123456 extended-status | grep -w "Threads_connected"

| Threads_connected | 1 |



注：如何计算当前打开的连接数占用最大连接数的比例呢？

result=Threads_connected/max_connections,在做监控报警或可视化监控时能够很好的根据这个比例及时调整最大连接数。



2.2 数据库性能监控

a.QPS:每秒的查询数

QPS计算方法

Questions = SHOW GLOBAL STATUS LIKE 'Questions';

Uptime = SHOW GLOBAL STATUS LIKE 'Uptime';

QPS=Questions/Uptime



b.TPS:每秒的事物量（commit与rollback的之和）

TPS计算方法

Com_commit = SHOW GLOBAL STATUS LIKE 'Com_commit';

Com_rollback = SHOW GLOBAL STATUS LIKE 'Com_rollback';

Uptime = SHOW GLOBAL STATUS LIKE 'Uptime';

TPS=(Com_commit + Com_rollback)/Uptime



2.3 数据库并发请求数量

MariaDB [(none)]> SHOW GLOBAL STATUS LIKE 'Threads_running';

+-----------------+-------+

| Variable_name | Value |

+-----------------+-------+

| Threads_running | 3 |

+-----------------+-------+

1 row in set (0.00 sec)

注：并发请求数量通常会远小于同一时间内连接到数据库的连接数数量。



2.4 监控innodb阻塞情况

a. innodb



三、对主从复制进行监控

3.1 主从复制链路状态的监控

3.2 主从复制延迟时间的监控

3.3 定期确认主从复制的数据是否一致

https://www.cnblogs.com/wwcom123/p/10759494.html


