---
title: Semisynchronous mysql半同步复制
layout: post
category: storage
author: 夏泽民
---
从MySQL5.5开始，MySQL以插件的形式支持半同步复制。如何理解半同步呢？首先我们来看看异步，全同步的概念

 

异步复制（Asynchronous replication）

MySQL默认的复制即是异步的，主库在执行完客户端提交的事务后会立即将结果返给给客户端，并不关心从库是否已经接收并处理，这样就会有一个问题，主如果crash掉了，此时主上已经提交的事务可能并没有传到从上，如果此时，强行将从提升为主，可能导致新主上的数据不完整。

 

全同步复制（Fully synchronous replication）

指当主库执行完一个事务，所有的从库都执行了该事务才返回给客户端。因为需要等待所有从库执行完该事务才能返回，所以全同步复制的性能必然会收到严重的影响。

 

半同步复制（Semisynchronous replication）

介于异步复制和全同步复制之间，主库在执行完客户端提交的事务后不是立刻返回给客户端，而是等待至少一个从库接收到并写到relay log中才返回给客户端。相对于异步复制，半同步复制提高了数据的安全性，同时它也造成了一定程度的延迟，这个延迟最少是一个TCP/IP往返的时间。所以，半同步复制最好在低延时的网络中使用。
<!-- more -->
半同步复制的潜在问题

客户端事务在存储引擎层提交后，在得到从库确认的过程中，主库宕机了，此时，可能的情况有两种

 

事务还没发送到从库上

此时，客户端会收到事务提交失败的信息，客户端会重新提交该事务到新的主上，当宕机的主库重新启动后，以从库的身份重新加入到该主从结构中，会发现，该事务在从库中被提交了两次，一次是之前作为主的时候，一次是被新主同步过来的。

 

事务已经发送到从库上

此时，从库已经收到并应用了该事务，但是客户端仍然会收到事务提交失败的信息，重新提交该事务到新的主上。

 

无数据丢失的半同步复制

针对上述潜在问题，MySQL 5.7引入了一种新的半同步方案：Loss-Less半同步复制。

 

针对上面这个图，“Waiting Slave dump”被调整到“Storage Commit”之前。

 

当然，之前的半同步方案同样支持，MySQL 5.7.2引入了一个新的参数进行控制-rpl_semi_sync_master_wait_point

rpl_semi_sync_master_wait_point有两种取值

 

AFTER_SYNC

这个即新的半同步方案，Waiting Slave dump在Storage Commit之前。

 

AFTER_COMMIT

老的半同步方案，如图所示。

 

半同步复制的安装部署

要想使用半同步复制，必须满足以下几个条件：

1. MySQL 5.5及以上版本

2. 变量have_dynamic_loading为YES

3. 异步复制已经存在

 

首先加载插件

因用户需执行INSTALL PLUGIN, SET GLOBAL, STOP SLAVE和START SLAVE操作，所以用户需有SUPER权限。

主：

mysql> INSTALL PLUGIN rpl_semi_sync_master SONAME 'semisync_master.so';

从：

mysql> INSTALL PLUGIN rpl_semi_sync_slave SONAME 'semisync_slave.so';

 

查看插件是否加载成功

有两种方式

1. 

mysql> show plugins;

rpl_semi_sync_master       | ACTIVE   | REPLICATION        | semisync_master.so | GPL  
2. 

mysql> SELECT PLUGIN_NAME, PLUGIN_STATUS FROM INFORMATION_SCHEMA.PLUGINS  WHERE PLUGIN_NAME LIKE '%semi%';

+----------------------+---------------+
| PLUGIN_NAME          | PLUGIN_STATUS |
+----------------------+---------------+
| rpl_semi_sync_master | ACTIVE        |
+----------------------+---------------+
1 row in set (0.00 sec)
 

启动半同步复制

在安装完插件后，半同步复制默认是关闭的，这时需设置参数来开启半同步

主：

mysql> SET GLOBAL rpl_semi_sync_master_enabled = 1;

从：

mysql> SET GLOBAL rpl_semi_sync_slave_enabled = 1;

 

以上的启动方式是在命令行操作，也可写在配置文件中。

主：

plugin-load=rpl_semi_sync_master=semisync_master.so
rpl_semi_sync_master_enabled=1
从：

plugin-load=rpl_semi_sync_slave=semisync_slave.so
rpl_semi_sync_slave_enabled=1
在有的高可用架构下，master和slave需同时启动，以便在切换后能继续使用半同步复制

plugin-load = "rpl_semi_sync_master=semisync_master.so;rpl_semi_sync_slave=semisync_slave.so"
rpl-semi-sync-master-enabled = 1
rpl-semi-sync-slave-enabled = 1
 

重启从上的IO线程

mysql> STOP SLAVE IO_THREAD;

mysql> START SLAVE IO_THREAD;

如果没有重启，则默认还是异步复制，重启后，slave会在master上注册为半同步复制的slave角色。

这时候，主的error.log中会打印如下信息：

2016-08-05T10:03:40.104327Z 5 [Note] While initializing dump thread for slave with UUID <ce9aaf22-5af6-11e6-850b-000c2988bad2>, found a zombie dump thread with the same UUID. Master is killing the zombie dump thread(4).
2016-08-05T10:03:40.111175Z 4 [Note] Stop asynchronous binlog_dump to slave (server_id: 2)
2016-08-05T10:03:40.119037Z 5 [Note] Start binlog_dump to master_thread_id(5) slave_server(2), pos(mysql-bin.000003, 621)
2016-08-05T10:03:40.119099Z 5 [Note] Start semi-sync binlog_dump to slave (server_id: 2), pos(mysql-bin.000003, 621)
 

查看半同步是否在运行

主：

mysql> show status like 'Rpl_semi_sync_master_status';

+-----------------------------+-------+
| Variable_name               | Value |
+-----------------------------+-------+
| Rpl_semi_sync_master_status | ON    |
+-----------------------------+-------+
1 row in set (0.00 sec)
从：

mysql> show status like 'Rpl_semi_sync_slave_status';

+----------------------------+-------+
| Variable_name              | Value |
+----------------------------+-------+
| Rpl_semi_sync_slave_status | ON    |
+----------------------------+-------+
1 row in set (0.20 sec)
这两个变量常用来监控主从是否运行在半同步复制模式下。

至此，MySQL半同步复制搭建完毕~

 

事实上，半同步复制并不是严格意义上的半同步复制

当半同步复制发生超时时（由rpl_semi_sync_master_timeout参数控制，单位是毫秒，默认为10000，即10s），会暂时关闭半同步复制，转而使用异步复制。当master dump线程发送完一个事务的所有事件之后，如果在rpl_semi_sync_master_timeout内，收到了从库的响应，则主从又重新恢复为半同步复制。

一、异步复制（Asynchronous replication）
1、逻辑上

MySQL默认的复制即是异步的，主库在执行完客户端提交的事务后会立即将结果返给给客户端，并不关心从库是否已经接收并处理，这样就会有一个问题，主如果crash掉了，此时主上已经提交的事务可能并没有传到从库上，如果此时，强行将从提升为主，可能导致新主上的数据不完整。

2、技术上

主库将事务 Binlog 事件写入到 Binlog 文件中，此时主库只会通知一下 Dump 线程发送这些新的 Binlog，然后主库就会继续处理提交操作，而此时不会保证这些 Binlog 传到任何一个从库节点上。

二、全同步复制（Fully synchronous replication）

1、逻辑上

指当主库执行完一个事务，所有的从库都执行了该事务才返回给客户端。因为需要等待所有从库执行完该事务才能返回，所以全同步复制的性能必然会收到严重的影响。

2、技术上

当主库提交事务之后，所有的从库节点必须收到、APPLY并且提交这些事务，然后主库线程才能继续做后续操作。但缺点是，主库完成一个事务的时间会被拉长，性能降低。

三、半同步复制（Semisynchronous replication）

1、逻辑上

是介于全同步复制与全异步复制之间的一种，主库只需要等待至少一个从库节点收到并且 Flush Binlog 到 Relay Log 文件即可，主库不需要等待所有从库给主库反馈。同时，这里只是一个收到的反馈，而不是已经完全完成并且提交的反馈，如此，节省了很多时间。

2、技术上

介于异步复制和全同步复制之间，主库在执行完客户端提交的事务后不是立刻返回给客户端，而是等待至少一个从库接收到并写到relay log中才返回给客户端。相对于异步复制，半同步复制提高了数据的安全性，同时它也造成了一定程度的延迟，这个延迟最少是一个TCP/IP往返的时间。所以，半同步复制最好在低延时的网络中使用。

四、选型及设置说明

如何设置到相应的同步方式上呢？

mysql主从模式默认是异步复制的，而MySQL Cluster是同步复制的，只要设置为相应的模式即是在使用相应的同步策略。

从MySQL5.5开始，MySQL以插件的形式支持半同步复制。其实说明半同步复制是更好的方式，兼顾了同步和性能的问题。
