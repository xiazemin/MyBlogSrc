---
title: percona-toolkit  Maatkit
layout: post
category: storage
author: 夏泽民
---
在mysql工作中接触最多的就是mysql replication，mysql在复制方面还是会有一些常规问题，比如主库宕机或者从库宕机有可能会导致复制中断，通常需要进行人为修复，或者很多时候需要把一个从库提升为主库，但对从库和主库的数据一致性不能保证一样。这种情况下就需要使用percona-toolkit工具的pt-table-checksum组件来检查主从数据的一致性；如果发现不一致的数据，可以通过pt-table-sync修复；还可以通过pt-heartbeat监控主从复制延迟。当然如果数据量小，slave只是当做一个备份使用，那么出现数据不一致完全可以重做，或者通过其他方法解决。如果数据量非常大，重做就是非常蛋碎的一件事情了。比如说，线上数据库做了主从同步环境，数据库在进行了迁移后，需要对mysql迁移（Replication）后的数据一致性进行校验，但又不能对生产环境使用造成影响，pt-table-checksum成为了绝佳也是唯一的检查工具。

percona-toolkit介绍
percona-toolkit是一组高级命令行工具的集合，用来执行各种通过手工执行非常复杂和麻烦的mysql和系统任务，这些任务包括：
   1）检查master和slave数据的一致性
   2）有效地对记录进行归档
   3）查找重复的索引
   4）对服务器信息进行汇总
   5）分析来自日志和tcpdump的查询
   6）当系统出问题的时候收集重要的系统信息
percona-toolkit源自Maatkit和Aspersa工具，这两个工具是管理mysql的最有名的工具。不过，现在Maatkit工具已经不维护了，所以以后推荐还是使用percona-toolkit工具！
这些工具主要包括开发、性能、配置、监控、复制、系统、实用六大类，作为一个优秀的DBA，里面有的工具非常有用，如果能掌握并加以灵活应用，将能极大的提高工作效率。

percona-toolkit工具中最主要的三个组件分别是：
   1）pt-table-checksum 负责监测mysql主从数据一致性
   2）pt-table-sync 负责当主从数据不一致时修复数据，让它们保存数据的一致性
   3）pt-heartbeat 负责监控mysql主从同步延迟
下面就对这三个组件的使用做一记录，当然percona-toolkit工具也有很多其他组件，后面会一一说明。

percona-toolkit工具安装（建议主库和从库服务器上都安装）
软件下载并在主库服务器上安装 [百度云盘下载地址：https://pan.baidu.com/s/1bp1OOgf   （提取密码：y462）]
[root@master-server src]# wget https://www.percona.com/downloads/percona-toolkit/2.2.7/RPM/percona-toolkit-2.2.7-1.noarch.rpm
[root@master-server src]# rpm -ivh percona-toolkit-2.2.7-1.noarch.rpm     //安装后，percona-toolkit工具的各个组件命令就有有了（输入ht-，按TAB键就会显示）

安装该工具依赖的软件包
[root@master-server src]# yum install perl-IO-Socket-SSL perl-DBD-MySQL perl-Time-HiRes perl perl-DBI -y

一、pt-table-checksum使用梳理
pt-table-checksum 是 Percona-Toolkit的组件之一，用于检测MySQL主、从库的数据是否一致。其原理是在主库执行基于statement的sql语句来生成主库数据块的checksum，把相同的sql语句传递到从库执行，并在从库上计算相同数据块的checksum，最后，比较主从库上相同数据块的checksum值，由此判断主从数据是否一致。检测过程根据唯一索引将表按row切分为块（chunk），以为单位计算，可以避免锁表。检测时会自动判断复制延迟、 master的负载， 超过阀值后会自动将检测暂停，减小对线上服务的影响。
pt-table-checksum 默认情况下可以应对绝大部分场景，官方说，即使上千个库、上万亿的行，它依然可以很好的工作，这源自于设计很简单，一次检查一个表，不需要太多的内存和多余的操作；必要时，pt-table-checksum 会根据服务器负载动态改变 chunk 大小，减少从库的延迟。

为了减少对数据库的干预，pt-table-checksum还会自动侦测并连接到从库，当然如果失败，可以指定--recursion-method选项来告诉从库在哪里。它的易用性还体现在，复制若有延迟，在从库 checksum 会暂停直到赶上主库的计算时间点（也通过选项--设定一个可容忍的延迟最大值，超过这个值也认为不一致）。

为了保证主数据库服务的安全，该工具实现了许多保护措施：
    1）自动设置 innodb_lock_wait_timeout 为1s，避免引起
    2）默认当数据库有25个以上的并发查询时，pt-table-checksum会暂停。可以设置 --max-load 选项来设置这个阀值
    3）当用 Ctrl+C 停止任务后，工具会正常的完成当前 chunk 检测，下次使用 --resume 选项启动可以恢复继续下一个 chunk

pt-table-checksum [OPTIONS] [DSN]
pt-table-checksum：在主（master）上通过执行校验的查询对复制的一致性进行检查，对比主从的校验值，从而产生结果。DSN指向的是主的地址，该工具的退出状态不为零，如果发现有任何差别，或者如果出现任何警告或错误。注意：第一次运行的时候需要加上--create-replicate-table参数，生成checksums表！！如果不加这个参数，那么就需要在对应库下手工添加这张表了,表结构SQL如下：
CREATE TABLE checksums (
   db             char(64)     NOT NULL,
   tbl            char(64)     NOT NULL,
   chunk          int          NOT NULL,
   chunk_time     float            NULL,
   chunk_index    varchar(200)     NULL,
   lower_boundary text             NULL,
   upper_boundary text             NULL,
   this_crc       char(40)     NOT NULL,
   this_cnt       int          NOT NULL,
   master_crc     char(40)         NULL,
   master_cnt     int              NULL,
   ts             timestamp    NOT NULL,
   PRIMARY KEY (db, tbl, chunk),
   INDEX ts_db_tbl (ts, db, tbl)
) ENGINE=InnoDB;
常用参数解释：
--nocheck-replication-filters ：不检查复制过滤器，建议启用。后面可以用--databases来指定需要检查的数据库。
--no-check-binlog-format : 不检查复制的binlog模式，要是binlog模式是ROW，则会报错。
--replicate-check-only :只显示不同步的信息。
--replicate= ：把checksum的信息写入到指定表中，建议直接写到被检查的数据库当中。
--databases= ：指定需要被检查的数据库，多个则用逗号隔开。
--tables= ：指定需要被检查的表，多个用逗号隔开
h= ：Master的地址
u= ：用户名
p=：密码
P= ：端口

最重要的一点就是：
要在主库上授权，能让主库ip访问。这一点不能忘记！（实验证明从库上可以不授权，但最好还是从库也授权）
注意：
1）根据测试，需要一个即能登录主库，也能登录从库的账号；
2）只能指定一个host，必须为主库的IP；
3）在检查时会向表加S锁；
4）运行之前需要从库的同步IO和SQL进程是YES状态。

例如：（本文例子中：192.168.1.101是主库ip，192.168.1.102是从库ip）
在主库执行授权（一定要对主库ip授权，授权的用户名和密码可以自行定义，不过要保证这个权限能同时登陆主库和从库）
mysql> GRANT SELECT, PROCESS, SUPER, REPLICATION SLAVE,CREATE,DELETE,INSERT,UPDATE ON *.* TO 'root'@'192.168.1.101' identified by '123456';
mysql> flush privileges;

在从库上执行授权
mysql> GRANT SELECT, PROCESS, SUPER, REPLICATION SLAVE ON *.* TO 'root'@'192.168.1.101' IDENTIFIED BY '123456';
mysql> flush privileges;

如下，在主库上执行的一个检查主从数据一致性的命令（别忘了第一次运行的时候需要添加--create-replicate-table参数，后续再运行时就不需要加了）：
下面命令中的192.168.1.101是主库ip
检查的是huanqiu库下的haha表的数据（当然，命令中也可以不跟表，直接检查某整个库的数据；如下去掉--tables=haha表，直接检查huanqiu库的数据）
[root@master-server ~]# pt-table-checksum --nocheck-replication-filters --no-check-binlog-format --replicate=huanqiu.checksums --create-replicate-table --databases=huanqiu --tables=haha h=192.168.1.101,u=root,p=123456,P=3306
Diffs cannot be detected because no slaves were found.  Please read the --recursion-method documentation for information.
            TS ERRORS  DIFFS     ROWS  CHUNKS SKIPPED    TIME TABLE
01-08T04:04:54      0      0        4       1       0   0.009 huanqiu.haha
上面有报错：
Diffs cannot be detected because no slaves were found. Please read the --recursion-method documentation for information
上面的提示信息很清楚，因为找不到从，所以执行失败，提示用参数--recursion-method 可以指定模式解决。
其实是因为从库的slave关闭了。
在主库上执行：
mysql> show processlist;
+----+------+-----------+------+---------+------+-------+------------------+
| Id | User | Host      | db   | Command | Time | State | Info             |
+----+------+-----------+------+---------+------+-------+------------------+
| 10 | root | localhost | NULL | Query   |    0 | init  | show processlist |
+----+------+-----------+------+---------+------+-------+------------------+
发现没有slave在运行。

在从库上开启slave
mysql> start slave;
mysql> show slave status\G;

再在主库上执行：
mysql> show processlist;
+----+-------+---------------------+------+-------------+------+-----------------------------------------------------------------------+------------------+
| Id | User  | Host                | db   | Command     | Time | State                                                                 | Info             |
+----+-------+---------------------+------+-------------+------+-----------------------------------------------------------------------+------------------+
| 10 | root  | localhost           | NULL | Query       |    0 | init                                                                  | show processlist |
| 18 | slave | 192.168.1.102:37115 | NULL | Binlog Dump |    5 | Master has sent all binlog to slave; waiting for binlog to be updated | NULL             |
+----+-------+---------------------+------+-------------+------+-----------------------------------------------------------------------+------------------+
发现已有slave在运行。

再次执行检查命令：
[root@master-server ~]# pt-table-checksum --nocheck-replication-filters --no-check-binlog-format --replicate=huanqiu.checksums --databases=huanqiu --tables=haha h=192.168.1.101,u=root,p=123456,P=3306
            TS ERRORS  DIFFS     ROWS  CHUNKS SKIPPED    TIME TABLE
01-08T04:11:03      0      0        4       1       0   1.422 huanqiu.haha
解释：
TS ：完成检查的时间。
ERRORS ：检查时候发生错误和警告的数量。
DIFFS ：0表示一致，1表示不一致。当指定--no-replicate-check时，会一直为0，当指定--replicate-check-only会显示不同的信息。
ROWS ：表的行数。
CHUNKS ：被划分到表中的块的数目。
SKIPPED ：由于错误或警告或过大，则跳过块的数目。
TIME ：执行的时间。
TABLE ：被检查的表名。

二、pt-table-sync用法梳理
如果通过pt-table-checksum 检查找到了不一致的数据表，那么如何同步数据呢？即如何修复MySQL主从不一致的数据，让他们保持一致性呢？
这时候可以利用另外一个工具pt-table-sync。
使用方法：
pt-table-sync: 高效的同步MySQL表之间的数据，他可以做单向和双向同步的表数据。他可以同步单个表，也可以同步整个库。它不同步表结构、索引、或任何其他模式对象。所以在修复一致性之前需要保证他们表存在。

假如上面检查数据时发现主从不一致
[root@master-server ~]# pt-table-checksum --nocheck-replication-filters --no-check-binlog-format --replicate=huanqiu.checksums --databases=huanqiu --tables=haha h=192.168.1.101,u=root,p=123456,P=3306
            TS ERRORS  DIFFS     ROWS  CHUNKS SKIPPED    TIME TABLE
01-08T04:18:07      0      1        4       1       0   0.843 huanqiu.haha
现在需要DIFFS为1可知主从数据不一致，需要修复！修复命令如下：
先master的ip，用户，密码，然后是slave的ip，用户，密码
[root@master-server ~]# pt-table-sync --replicate=huanqiu.checksums h=192.168.1.101,u=root,p=123456 h=192.168.1.102,u=root,p=123456 --print
REPLACE INTO `huanqiu`.`haha`(`id`, `name`) VALUES ('1', 'wangshibo') /*percona-toolkit src_db:huanqiu src_tbl:haha src_dsn:h=192.168.1.101,p=...,u=root dst_db:huanqiu dst_tbl:haha dst_dsn:h=192.168.1.102,p=...,u=root lock:1 transaction:1 changing_src:huanqiu.checksums replicate:huanqiu.checksums bidirectional:0 pid:23676 user:root host:master-server*/;
REPLACE INTO `huanqiu`.`haha`(`id`, `name`) VALUES ('2', 'wangshikui') /*percona-toolkit src_db:huanqiu src_tbl:haha src_dsn:h=192.168.1.101,p=...,u=root dst_db:huanqiu dst_tbl:haha dst_dsn:h=192.168.1.102,p=...,u=root lock:1 transaction:1 changing_src:huanqiu.checksums replicate:huanqiu.checksums bidirectional:0 pid:23676 user:root host:master-server*/;
REPLACE INTO `huanqiu`.`haha`(`id`, `name`) VALUES ('3', 'limeng') /*percona-toolkit src_db:huanqiu src_tbl:haha src_dsn:h=192.168.1.101,p=...,u=root dst_db:huanqiu dst_tbl:haha dst_dsn:h=192.168.1.102,p=...,u=root lock:1 transaction:1 changing_src:huanqiu.checksums replicate:huanqiu.checksums bidirectional:0 pid:23676 user:root host:master-server*/;
REPLACE INTO `huanqiu`.`haha`(`id`, `name`) VALUES ('4', 'wanghi') /*percona-toolkit src_db:huanqiu src_tbl:haha src_dsn:h=192.168.1.101,p=...,u=root dst_db:huanqiu dst_tbl:haha dst_dsn:h=192.168.1.102,p=...,u=root lock:1 transaction:1 changing_src:huanqiu.checksums replicate:huanqiu.checksums bidirectional:0 pid:23676 user:root host:master-server*/;
参数解释：
--replicate= ：指定通过pt-table-checksum得到的表，这2个工具差不多都会一直用。
--databases= : 指定执行同步的数据库。
--tables= ：指定执行同步的表，多个用逗号隔开。
--sync-to-master ：指定一个DSN，即从的IP，他会通过show processlist或show slave status 去自动的找主。
h= ：服务器地址，命令里有2个ip，第一次出现的是Master的地址，第2次是Slave的地址。
u= ：帐号。
p= ：密码。
--print ：打印，但不执行命令。
--execute ：执行命令。

上面命令介绍完了，接下来开始执行修复：
通过（--print）打印出来了修复数据的sql语句，可以手动的在slave从库上执行，让他们数据保持一致性，这样比较麻烦！
可以直接在master主库上执行修复操作，通过--execute参数，如下：
[root@master-server ~]# pt-table-sync --replicate=huanqiu.checksums h=192.168.1.101,u=root,p=123456 h=192.168.1.102,u=root,p=123456 --execute

如上修复后，再次检查，发现主从库数据已经一致了！
[root@master-server ~]# pt-table-checksum --nocheck-replication-filters --no-check-binlog-format --replicate=huanqiu.checksums --databases=huanqiu --tables=haha h=192.168.1.101,u=root,p=123456,P=3306
            TS ERRORS  DIFFS     ROWS  CHUNKS SKIPPED    TIME TABLE
01-08T04:36:43      0      0        4       1       0   0.040 huanqiu.haha
-----------------------------------------------------------------------------------------------------------------------
建议:
修复数据的时候，最好还是用--print打印出来的好，这样就可以知道那些数据有问题，可以人为的干预下。
不然直接执行了，出现问题之后更不好处理。总之还是在处理之前做好数据的备份工作。

注意：要是表中没有唯一索引或则主键则会报错：
Can't make changes on the master because no unique index exists at /usr/local/bin/pt-table-sync line 10591.
-----------------------------------------------------------------------------------------------------------------------
为了确保主从数据的一致性，可以编写监控脚本，定时检查。当检查到主从数据不一致时，强制修复数据。
[root@master-server ~]# cat /root/pt_huanqiu.sh
#!/bin/bash
NUM=$(/usr/bin/pt-table-checksum --nocheck-replication-filters --no-check-binlog-format --replicate=huanqiu.checksums --databases=huanqiu  h=192.168.1.101,u=root,p=123456,P=3306|awk -F" " '{print $3}'|sed -n '2p')
if [ $NUM -eq 1 ];then
  /usr/bin/pt-table-sync --replicate=huanqiu.checksums h=192.168.1.101,u=root,p=123456 h=192.168.1.102,u=root,p=123456 --print
  /usr/bin/pt-table-sync --replicate=huanqiu.checksums h=192.168.1.101,u=root,p=123456 h=192.168.1.102,u=root,p=123456 --execute
else
  echo "data is ok"
fi
[root@master-server ~]# cat /root/pt_huanpc.sh 
#!/bin/bash
NUM=$(/usr/bin/pt-table-checksum --nocheck-replication-filters --no-check-binlog-format --replicate=huanpc.checksums --databases=huanpc  h=192.168.1.101,u=root,p=123456,P=3306|awk -F" " '{print $3}'|sed -n '2p')
if [ $NUM -eq 1 ];then
  /usr/bin/pt-table-sync --replicate=huanpc.checksums h=192.168.1.101,u=root,p=123456 h=192.168.1.102,u=root,p=123456 --print
  /usr/bin/pt-table-sync --replicate=huanpc.checksums h=192.168.1.101,u=root,p=123456 h=192.168.1.102,u=root,p=123456 --execute
else
  echo "data is ok"
fi
[root@master-server ~]# crontab -l
#检查主从huanqiu库数据一致性
* * * * * /bin/bash -x /root/pt_huanqiu.sh > /dev/null 2>&1
* * * * * sleep 10;/bin/bash -x /root/pt_huanqiu.sh > /dev/null 2>&1
* * * * * sleep 20;/bin/bash -x /root/pt_huanqiu.sh > /dev/null 2>&1
* * * * * sleep 30;/bin/bash -x /root/pt_huanqiu.sh > /dev/null 2>&1
* * * * * sleep 40;/bin/bash -x /root/pt_huanqiu.sh > /dev/null 2>&1
* * * * * sleep 50;/bin/bash -x /root/pt_huanqiu.sh > /dev/null 2>&1

#检查主从huanpc库数据一致性
* * * * * /bin/bash -x /root/root/pt_huanpc.sh > /dev/null 2>&1
* * * * * sleep 10;/bin/bash -x /root/pt_huanpc.sh > /dev/null 2>&1
* * * * * sleep 20;/bin/bash -x /root/pt_huanpc.sh > /dev/null 2>&1
* * * * * sleep 30;/bin/bash -x /root/pt_huanpc.sh > /dev/null 2>&1
* * * * * sleep 40;/bin/bash -x /root/pt_huanpc.sh > /dev/null 2>&1
* * * * * sleep 50;/bin/bash -x /root/pt_huanpc.sh > /dev/null 2>&1

-----------------------------------------------------------------------------------------------------------------------
最后总结：
pt-table-checksum和pt-table-sync工具很给力，工作中常常在使用。注意使用该工具需要授权，一般SELECT, PROCESS, SUPER, REPLICATION SLAVE等权限就已经足够了。

-----------------------------------------------------------------------------------------------------------------------
另外说一个问题：
在上面的操作中，在主库里添加pt-table-checksum检查的权限（从库可以不授权）后，进行数据一致性检查操作，会在操作的库（实例中是huanqiu、huanpc）下产生一个checksums表！
这张checksums表是pt-table-checksum检查过程中产生的。这张表一旦产生了，默认是删除不了的，并且这张表所在的库也默认删除不了，删除后过一会儿就又会出来。
也就是说，checksums表一旦产生，不仅这张表默认删除不了，连同它所在的库，要是想删除它们，只能如上操作先撤销权限。

三、pt-heartbeat监控mysql主从复制延迟梳理

对于MySQL数据库主从复制延迟的监控，可以借助percona的有力武器pt-heartbeat来实现。
pt-heartbeat的工作原理通过使用时间戳方式在主库上更新特定表，然后在从库上读取被更新的时间戳然后与本地系统时间对比来得出其延迟。具体流程：
   1）在主上创建一张heartbeat表，按照一定的时间频率更新该表的字段（把时间更新进去）。监控操作运行后，heartbeat表能促使主从同步！
   2）连接到从库上检查复制的时间记录，和从库的当前系统时间进行比较，得出时间的差异。

使用方法（主从和从库上都可以执行监控操作）：
pt-heartbeat [OPTIONS] [DSN] --update|--monitor|--check|--stop
注意：需要指定的参数至少有 --stop，--update，--monitor，--check。
其中--update，--monitor和--check是互斥的，--daemonize和--check也是互斥。
--ask-pass     隐式输入MySQL密码
--charset     字符集设置
--check      检查从的延迟，检查一次就退出，除非指定了--recurse会递归的检查所有的从服务器。
--check-read-only    如果从服务器开启了只读模式，该工具会跳过任何插入。
--create-table    在主上创建心跳监控的表，如果该表不存在，可以自己手动建立，建议存储引擎改成memory。通过更新该表知道主从延迟的差距。
CREATE TABLE heartbeat (
  ts                    varchar(26) NOT NULL,
  server_id             int unsigned NOT NULL PRIMARY KEY,
  file                  varchar(255) DEFAULT NULL,
  position              bigint unsigned DEFAULT NULL,
  relay_master_log_file varchar(255) DEFAULT NULL,
  exec_master_log_pos   bigint unsigned DEFAULT NULL
);
heratbeat   表一直在更改ts和position,而ts是我们检查复制延迟的关键。
--daemonize   执行时，放入到后台执行
--user=-u，   连接数据库的帐号
--database=-D，    连接数据库的名称
--host=-h，     连接的数据库地址
--password=-p，     连接数据库的密码
--port=-P，     连接数据库的端口
--socket=-S，    连接数据库的套接字文件
--file 【--file=output.txt】   打印--monitor最新的记录到指定的文件，很好的防止满屏幕都是数据的烦恼。
--frames 【--frames=1m,2m,3m】  在--monitor里输出的[]里的记录段，默认是1m,5m,15m。可以指定1个，如：--frames=1s，多个用逗号隔开。可用单位有秒（s）、分钟（m）、小时（h）、天（d）。
--interval   检查、更新的间隔时间。默认是见是1s。最小的单位是0.01s，最大精度为小数点后两位，因此0.015将调整至0.02。
--log    开启daemonized模式的所有日志将会被打印到制定的文件中。
--monitor    持续监控从的延迟情况。通过--interval指定的间隔时间，打印出从的延迟信息，通过--file则可以把这些信息打印到指定的文件。
--master-server-id    指定主的server_id，若没有指定则该工具会连到主上查找其server_id。
--print-master-server-id    在--monitor和--check 模式下，指定该参数则打印出主的server_id。
--recurse    多级复制的检查深度。模式M-S-S...不是最后的一个从都需要开启log_slave_updates，这样才能检查到。
--recursion-method     指定复制检查的方式,默认为processlist,hosts。
--update    更新主上的心跳表。
--replace     使用--replace代替--update模式更新心跳表里的时间字段，这样的好处是不用管表里是否有行。
--stop    停止运行该工具（--daemonize），在/tmp/目录下创建一个“pt-heartbeat-sentinel” 文件。后面想重新开启则需要把该临时文件删除，才能开启（--daemonize）。
--table   指定心跳表名，默认heartbeat。
实例说明：
master：192.168.1.101
slave：192.168.1.102
同步的库：huanqiu、huanpc
主从库都能使用root账号、密码123456登录

先操作针对huanqiu库的检查，其他同步的库的检查操作类似！
mysql> use huanqiu;                   
Database changed
 
mysql> CREATE TABLE heartbeat (            //主库上的对应库下创建heartbeat表，一般创建后从库会同步这张表（不同步的话，就在从库那边手动也手动创建）
    ->   ts                    varchar(26) NOT NULL,
    ->   server_id             int unsigned NOT NULL PRIMARY KEY,
    ->   file                  varchar(255) DEFAULT NULL,
    ->   position              bigint unsigned DEFAULT NULL,
    ->   relay_master_log_file varchar(255) DEFAULT NULL,
    ->   exec_master_log_pos   bigint unsigned DEFAULT NULL
    -> );
Query OK, 0 rows affected (0.02 sec)
更新主库上的heartbeat,--interval=1表示1秒钟更新一次（注意这个启动操作要在主库服务器上执行）
[root@master-server ~]# pt-heartbeat --user=root --ask-pass --host=192.168.1.101 --create-table -D huanqiu --interval=1 --update --replace --daemonize
Enter password: 
[root@master-server ~]# 
[root@master-server ~]# ps -ef|grep pt-heartbeat
root 15152 1 0 19:49 ? 00:00:00 perl /usr/bin/pt-heartbeat --user=root --ask-pass --host=192.168.1.101 --create-table -D huanqiu --interval=1 --update --replace --daemonize
root 15154 14170 0 19:49 pts/3 00:00:00 grep pt-heartbeat

在主库运行监测同步延迟：
<!-- more -->
如何关闭上面在主库上执行的heartbeat更新进程呢？
方法一：可以用参数--stop去关闭
[root@master-server ~]# ps -ef|grep heartbeat
root 15152 1 0 19:49 ? 00:00:02 perl /usr/bin/pt-heartbeat --user=root --ask-pass --host=192.168.1.101 --create-table -D huanqiu --interval=1 --update --replace --daemonize
root 15310 1 0 19:59 ? 00:00:01 perl /usr/bin/pt-heartbeat -D huanqiu --table=heartbeat --monitor --host=192.168.1.102 --user=root --password=123456 --log=/opt/master-slave.txt --daemonize
root 15555 31932 0 20:13 pts/2 00:00:00 grep heartbeat
[root@master-server ~]# pt-heartbeat --stop
Successfully created file /tmp/pt-heartbeat-sentinel
[root@master-server ~]# ps -ef|grep heartbeat
root 15558 31932 0 20:14 pts/2 00:00:00 grep heartbeat
[root@master-server ~]#

这样就把在主上开启的进程杀掉了。
但是后续要继续开启后台进行的话，记住一定要先把/tmp/pt-heartbeat-sentinel 文件删除，否则启动不了

方法二：直接kill掉进程pid（推荐这种方法）
[root@master-server ~]# ps -ef|grep heartbeat
root 15152 1 0 19:49 ? 00:00:02 perl /usr/bin/pt-heartbeat --user=root --ask-pass --host=192.168.1.101 --create-table -D huanqiu --interval=1 --update --replace --daemonize
root 15310 1 0 19:59 ? 00:00:01 perl /usr/bin/pt-heartbeat -D huanqiu --table=heartbeat --monitor --host=192.168.1.102 --user=root --password=123456 --log=/opt/master-slave.txt --daemonize
root 15555 31932 0 20:13 pts/2 00:00:00 grep heartbeat
[root@master-server ~]# kill -9 15152
[root@master-server ~]# ps -ef|grep heartbeat
root 15558 31932 0 20:14 pts/2 00:00:00 grep heartbeat

最后总结：
通过pt-heartbeart工具可以很好的弥补默认主从延迟的问题，但需要搞清楚该工具的原理。
默认的Seconds_Behind_Master值是通过将服务器当前的时间戳与二进制日志中的事件时间戳相对比得到的，所以只有在执行事件时才能报告延时。备库复制线程没有运行，也会报延迟null。
还有一种情况：大事务，一个事务更新数据长达一个小时，最后提交。这条更新将比它实际发生时间要晚一个小时才记录到二进制日志中。当备库执行这条语句时，会临时地报告备库延迟为一个小时，执行完后又很快变成0。

---------------------------------------percona-toolkit其他组件命令用法---------------------------------- 

下面这些工具最好不要直接在线上使用，应该作为上线辅助或故障后离线分析的工具，也可以做性能测试的时候配合着使用。

1）pt-online-schema-change
功能介绍：
功能为:在alter操作更改表结构的时候不用锁定表，也就是说执行alter的时候不会阻塞写和读取操作，注意执行这个工具的时候必须做好备份，操作之前最好要充分了解它的原理。
工作原理是:创建一个和你要执行alter操作的表一样的空表结构，执行表结构修改，然后从原表中copy原始数据到表结构修改后的表，当数据copy完成以后就会将原表移走，用新表代替原表，默认动作是将原表drop掉。在copy数据的过程中，任何在原表的更新操作都会更新到新表，因为这个工具在会在原表上创建触发器，触发器会将在原表上更新的内容更新到新表。如果表中已经定义了触发器这个工具就不能工作了。

用法介绍：
pt-online-schema-change [OPTIONS] DSN
options可以自行查看help（或加--help查看有哪些选项），DNS为你要操作的数据库和表。
有两个参数需要注意一下：
--dry-run 这个参数不建立触发器，不拷贝数据，也不会替换原表。只是创建和更改新表。
--execute 这个参数的作用和前面工作原理的介绍的一样，会建立触发器，来保证最新变更的数据会影响至新表。注意：如果不加这个参数，这个工具会在执行一些检查后退出。这一举措是为了让使用这充分了解了这个工具的原理。

3）pt-slave-find
功能介绍：
查找和打印mysql所有从服务器复制层级关系
用法介绍：
pt-slave-find [OPTION...] MASTER-HOST
原理：连接mysql主服务器并查找其所有的从，然后打印出所有从服务器的层级关系。
4）pt-show-grants
功能介绍：
规范化和打印mysql权限，让你在复制、比较mysql权限以及进行版本控制的时候更有效率！
用法介绍：
pt-show-grants [OPTION...] [DSN]
选项自行用help查看，DSN选项也请查看help，选项区分大小写。
使用示例：
查看指定mysql的所有用户权限：
5）pt-upgrade
功能介绍：
这个工具用来检查在新版本中运行的SQL是否与老版本一样，返回相同的结果，最好的应用场景就是数据迁移的时候。这在升级服务器的时候非常有用，可以先安装并导数据到新的服务器上，然后使用这个工具跑一下sql看看有什么不同，可以找出不同版本之间的差异。
用法介绍：
pt-upgrade [OPTION...] DSN [DSN...] [FILE]
比较文件中每一个查询语句在每台服务器上执行的结果（主要是针对不同版本的执行结果）。（--help查看选项）
使用示例：
查看某个sql文件在两个服务器的运行结果范例：
6）pt-index-usage
功能介绍：
这个工具主要是用来分析慢查询的索引使用情况。从log文件中读取插叙语句，并用explain分析他们是如何利用索引。完成分析之后会生成一份关于索引没有被查询使用过的报告。
用法介绍：
pt-index-usage [OPTION...] [FILE...]
可以直接从慢查询中获取sql，FILE文件中的sql格式必须和慢查询中个是一致，如果不是一直需要用pt-query-digest转换一下。也可以不生成报告直接保存到数据库中，具体的见后面的示例
注意：使用这个工具需要MySQL必须要有密码，另外运行时可能报找不到/var/lib/mysql/mysql.sock的错，简单的从mysql启动后的sock文件做一个软链接即可。
重点要说明的是pt-index-usage只能分析慢查询日志，所以如果想全面分析所有查询的索引使用情况就得将slow_launch_time设置为0，因此请谨慎使用该工具，线上使用的话最好在凌晨进行分析，尤其分析大量日志的时候是很耗CPU的。
整体来说这个工具是不推荐使用的，要想实现类似的分析可以考虑一些其他第三方的工具，比如：mysqlidxchx, userstat和check-unused-keys。网上比较推荐的是userstat，一个Google贡献的patch。
使用示例：
从满查询中的sql查看索引使用情况范例：
[root@master-server ~]# pt-index-usage --host=localhost --user=root --password=123456 /data/mysql/data/mysql-slow.log
将分析结果保存到数据库范例：
[root@master-server ~]# pt-index-usage --host=localhost --user=root --password=123456 /data/mysql/data/mysql-slow.log  --no-report --create-save-results-database
7）pt-visual-explain
功能介绍：
格式化explain出来的执行计划按照tree方式输出，方便阅读。
用法介绍：
pt-visual-explain [OPTION...] [FILE...]
8）pt-config-diff
功能介绍：
比较mysql配置文件和服务器参数
用法介绍：
pt-config-diff [OPTION...] CONFIG CONFIG [CONFIG...]
CONFIG可以是文件也可以是数据源名称，最少必须指定两个配置文件源，就像unix下面的diff命令一样，如果配置完全一样就不会输出任何东西。
9）pt-mysql-summary
功能介绍：
精细地对mysql的配置和sataus信息进行汇总，汇总后你直接看一眼就能看明白。
工作原理：连接mysql后查询出status和配置信息保存到临时目录中，然后用awk和其他的脚本工具进行格式化。OPTIONS可以查阅官网的相关页面。
用法介绍：
pt-mysql-summary [OPTIONS] [-- MYSQL OPTIONS]
10）pt-deadlock-logger
功能介绍：
提取和记录mysql死锁的相关信息
用法介绍：
pt-deadlock-logger [OPTION...] SOURCE_DSN
收集和保存mysql上最近的死锁信息，可以直接打印死锁信息和存储死锁信息到数据库中，死锁信息包括发生死锁的服务器、最近发生死锁的时间、死锁线程id、死锁的事务id、发生死锁时事务执行了多长时间等等非常多的信息。
使用示例：
查看本地mysql的死锁信息
[root@master-server ~]# pt-deadlock-logger  --user=root --password=123456 h=localhost D=test,t=deadlocks
server ts thread txn_id txn_time user hostname ip db tbl idx lock_type lock_mode wait_hold victim query
localhost 2017-01-11T11:00:33 188 0 0 root  192.168.1.101 huanpc checksums PRIMARY RECORD X w 1 REPLACE INTO `huanpc`.`checksums` (db, tbl, chunk, chunk_index, lower_boundary, upper_boundary, this_cnt, this_crc) SELECT 'huanpc', 'heihei', '1', NULL, NULL, NULL, COUNT(*) AS cnt, COALESCE(LOWER(CONV(BIT_XOR(CAST(CRC32(CONCAT_WS('#', `member`, `city`)) AS UNSIGNED)), 10, 16)), 0) AS crc FROM `huanpc`.`heihei` /*checksum table*/
localhost 2017-01-11T11:00:33 198 0 0 root  192.168.1.101 huanpc checksums PRIMARY RECORD X w 0 REPLACE INTO `huanpc`.`checksums` (db, tbl, chunk, chunk_index, lower_boundary, upper_boundary, this_cnt, this_crc) SELECT 'huanpc', 'heihei', '1', NULL, NULL, NULL, COUNT(*) AS cnt, COALESCE(LOWER(CONV(BIT_XOR(CAST(CRC32(CONCAT_WS('#', `member`, `city`)) AS UNSIGNED)), 10, 16)), 0) AS crc FROM `huanpc`.`heihei` /*checksum table*/
11）pt-mext
功能介绍：
并行查看SHOW GLOBAL STATUS的多个样本的信息。
用法介绍：
pt-mext [OPTIONS] -- COMMAND
原理：pt-mext执行你指定的COMMAND，并每次读取一行结果，把空行分割的内容保存到一个一个的临时文件中，最后结合这些临时文件并行查看结果。
使用示例：
每隔10s执行一次SHOW GLOBAL STATUS，并将结果合并到一起查看

1
[root@master-server ~]# pt-mext  -- mysqladmin ext -uroot -p123456  -i10 -c3
12）pt-query-digest
功能介绍：
分析查询执行日志，并产生一个查询报告，为MySQL、PostgreSQL、 memcached过滤、重放或者转换语句。
pt-query-digest可以从普通MySQL日志，慢查询日志以及二进制日志中分析查询，甚至可以从SHOW PROCESSLIST和MySQL协议的tcpdump中进行分析，如果没有指定文件，它从标准输入流（STDIN）中读取数据。
用法介绍：
pt-query-digest [OPTION...] [FILE]
解析和分析mysql日志文件
整个输出分为三大部分：
1）整体概要（Overall）
这个部分是一个大致的概要信息(类似loadrunner给出的概要信息)，通过它可以对当前MySQL的查询性能做一个初步的评估，比如各个指标的最大值(max)，平均值(min)，95%分布值，中位数(median)，标准偏差(stddev)。
这些指标有查询的执行时间（Exec time），锁占用的时间（Lock time），MySQL执行器需要检查的行数（Rows examine），最后返回给客户端的行数（Rows sent），查询的大小。
 
2）查询的汇总信息（Profile）
这个部分对所有“重要”的查询(通常是比较慢的查询)做了个一览表。
每个查询都有一个Query ID，这个ID通过Hash计算出来的。pt-query-digest是根据这个所谓的Fingerprint来group by的。
Rank整个分析中该“语句”的排名，一般也就是性能最常的。
Response time  “语句”的响应时间以及整体占比情况。
Calls 该“语句”的执行次数。
R/Call 每次执行的平均响应时间。
V/M 响应时间的差异平均对比率。
在尾部有一行输出，显示了其他2个占比较低而不值得单独显示的查询的统计数据。
 
3）详细信息
这个部分会列出Profile表中每个查询的详细信息：
包括Overall中有的信息、查询响应时间的分布情况以及该查询”入榜”的理由。
pt-query-digest还有很多复杂的操作，这里就不一一介绍了。比如：从PROCESSLIST中查询某个MySQL中最慢的查询：
从tcpdump中分析：
[root@master-server ~]# tcpdump -s 65535 -x -nn -q -tttt -i any -c 1000 port 3306 > mysql.tcp.txt
tcpdump: verbose output suppressed, use -v or -vv for full protocol decode
listening on any, link-type LINUX_SLL (Linux cooked), capture size 65535 bytes
 
然后打开另一个终端窗口：
[root@master-server ~]# pt-query-digest --type tcpdump mysql.tcp.txt
Pipeline process 3 (TcpdumpParser) caused an error: substr outside of string at /usr/bin/pt-query-digest line 3628, <> chunk 93.
Will retry pipeline process 2 (TcpdumpParser) 100 more times.
 
# 320ms user time, 20ms system time, 24.93M rss, 204.84M vsz
# Current date: Mon Jan 16 13:24:50 2017
# Hostname: master-server
# Files: mysql.tcp.txt
# Overall: 31 total, 4 unique, 4.43 QPS, 0.00x concurrency _______________
# Time range: 2017-01-16 13:24:43.000380 to 13:24:50.001205
# Attribute          total     min     max     avg     95%  stddev  median
# ============     ======= ======= ======= ======= ======= ======= =======
# Exec time           30ms    79us     5ms   967us     4ms     1ms   159us
# Rows affecte          14       0       2    0.45    1.96    0.82       0
# Query size         1.85k      17     200   61.16  192.76   72.25   17.65
.........
13）pt-slave-delay
功能介绍：
设置从服务器落后于主服务器指定时间。
用法介绍：
pt-slave-delay [OPTION...] SLAVE-HOST [MASTER-HOST]
原理：通过启动和停止复制sql线程来设置从落后于主指定时间。默认是基于从上relay日志的二进制日志的位置来判断，因此不需要连接到主服务器，如果IO进程不落后主服务器太多的话，这个检查方式工作很好，如果网络通畅的话，一般IO线程落后主通常都是毫秒级别。一般是通过--delay and --delay"+"--interval来控制。--interval是指定检查是否启动或者停止从上sql线程的频繁度，默认的是1分钟检查一次。
使用示例：
范例1：使从落后主1分钟，并每隔1分钟检测一次，运行10分钟
[root@master-server ~]# pt-slave-delay --user=root --password=123456 --delay 1m --run-time 10m --host=192.168.1.102
2017-01-16T13:32:31 slave running 0 seconds behind
2017-01-16T13:32:31 STOP SLAVE until 2017-01-16T13:33:31 at master position mysql-bin.000005/102554361
范例2：使从落后主1分钟，并每隔15秒钟检测一次，运行10分钟：
[root@master-server ~]# pt-slave-delay --user=root --password=123456 --delay 1m --interval 15s --run-time 10m --host=192.168.1.102
2017-01-16T13:38:22 slave running 0 seconds behind
2017-01-16T13:38:22 STOP SLAVE until 2017-01-16T13:39:22 at master position mysql-bin.000005/102689359
14）pt-slave-restart
功能介绍：
监视mysql复制错误，并尝试重启mysql复制当复制停止的时候
用法介绍：
pt-slave-restart [OPTION...] [DSN]
监视一个或者多个mysql复制错误，当从停止的时候尝试重新启动复制。你可以指定跳过的错误并运行从到指定的日志位置。
使用示例：
范例1：监视192.168.1.101的从，跳过1个错误
[root@master-server ~]# pt-slave-restart --user=root --password=123456 --host=192.168.1.101 --skip-count=1
范例2：监视192.168.1.101的从，跳过错误代码为1062的错误。
[root@master-server ~]# pt-slave-restart --user=root --password=123456 --host=192.168.1.101 --error-numbers=1062
15）pt-diskstats
功能介绍：
是一个对GUN/LINUX的交互式监控工具
用法介绍：
pt-diskstats [OPTION...] [FILES]
为GUN/LINUX打印磁盘io统计信息，和iostat有点像，但是这个工具是交互式并且比iostat更详细。可以分析从远程机器收集的数据。
使用示例：
范例1：查看本机所有的磁盘的状态情况：
[root@master-server ~]# pt-diskstats
范例2：只查看本机sdc1磁盘的状态情况：
[root@master-server ~]# pt-diskstats  --devices-regex vdc1
  #ts device    rd_s rd_avkb rd_mb_s rd_mrg rd_cnc   rd_rt    wr_s wr_avkb wr_mb_s wr_mrg wr_cnc   wr_rt busy in_prg    io_s  qtime stime
  0.9 vdc1       0.0     0.0     0.0     0%    0.0     0.0     5.9     4.0     0.0     0%    0.0     1.0   0%      0     5.9    0.6   0.4
  1.0 vdc1       0.0     0.0     0.0     0%    0.0     0.0     2.0     6.0     0.0    33%    0.0     0.7   0%      0     2.0    0.0   0.7
16）pt-summary
功能介绍：
友好地收集和显示系统信息概况，此工具并不是一个调优或者诊断工具，这个工具会产生一个很容易进行比较和发送邮件的报告。
用法介绍：
pt-summary
原理：此工具会运行和多命令去收集系统状态和配置信息，先保存到临时目录的文件中去，然后运行一些unix命令对这些结果做格式化，最好是用root用户或者有权限的用户运行此命令。
使用示例：
查看本地系统信息概况
[root@master-server ~]# pt-summary
17）pt-stalk
功能介绍：
出现问题的时候收集mysql的用于诊断的数据
用法介绍：
pt-stalk [OPTIONS] [-- MYSQL OPTIONS]
pt-stalk等待触发条件触发，然后收集数据帮助错误诊断，它被设计成使用root权限运行的守护进程，因此你可以诊断那些你不能直接观察的间歇性问题。默认的诊断触发条件为SHOW GLOBAL STATUS。也可以指定processlist为诊断触发条件 ，使用--function参数指定。
使用示例：
范例1：指定诊断触发条件为status，同时运行语句超过20的时候触发，收集的数据存放在目标目录/tmp/test下：
[root@master-server ~]# pt-stalk  --function status --variable Threads_running --threshold 20 --dest /tmp/test  -- -uroot -p123456  -h192.168.1.101
范例2：指定诊断触发条件为processlist，超过20个状态为statistics触发，收集的数据存放在/tmp/test目录下：
[root@master-server ~]# pt-stalk  --function processlist --variable State --match statistics --threshold 20 --dest /tmp/test -- -uroot -p123456  -h192.168.1.101
.......
2017_01_15_17_31_49-hostname
2017_01_15_17_31_49-innodbstatus1
2017_01_15_17_31_49-innodbstatus2
2017_01_15_17_31_49-interrupts
2017_01_15_17_31_49-log_error
2017_01_15_17_31_49-lsof
2017_01_15_17_31_49-meminfo
18）pt-archiver
功能介绍：
将mysql数据库中表的记录归档到另外一个表或者文件
用法介绍：
pt-archiver [OPTION...] --source DSN --where WHERE
这个工具只是归档旧的数据，不会对线上数据的OLTP查询造成太大影响，你可以将数据插入另外一台服务器的其他表中，也可以写入到一个文件中，方便使用source命令导入数据。另外你还可以用它来执行delete操作。特别注意：这个工具默认的会删除源中的数据！！
19）pt-find
功能介绍：
查找mysql表并执行指定的命令，和gnu的find命令类似。
用法介绍：
pt-find [OPTION...] [DATABASE...]
默认动作是打印数据库名和表名
20）pt-kill
功能介绍：
Kill掉符合指定条件mysql语句
用法介绍：
pt-kill [OPTIONS]
加入没有指定文件的话pt-kill连接到mysql并通过SHOW PROCESSLIST找到指定的语句，反之pt-kill从包含SHOW PROCESSLIST结果的文件中读取mysql语句

Maatkit是一个开源的工具包，为mySQL日常管理提供了帮助，它包含很多工具，这里主要说下面两个：
1）mk-table-checksum                     用来检测master和slave上的表结构和数据是否一致的；
2）mk-table-sync                             在主从数据不一致时，用来修复数据的；先主后从有效保证表一致的工具，不必重载从表而能够保证一致。 
上面两个perl脚本在运行时都会锁表，表的大小取决于执行的快慢，勿在高峰期间运行，可选择凌晨。
-----------------------------------------------------------------------------------------------------

下面记录了这一操作过程：
基本信息：
master：192.168.1.101
slave：192.168.1.102
版本：mysql5.6
同步的库：huanqiu、huanpc

Maatkit安装过程：（主库和从库服务器上都可以安装，因为数据一致性检查操作在主库或从库机器上都可以运行；建议主从机器上都安装）
1）安装该工具依赖的软件包
[root@master-server src]# yum install perl-IO-Socket-SSL perl-DBD-MySQL perl-Time-HiRes perl perl-DBI -y

2）maatkit下载安装    [需要FQ到官网下载：https://code.google.com/archive/p/maatkit/downloads]
百度云盘下载地址：https://pan.baidu.com/s/1c1AufW8    （提取密码：vbi1）
[root@master-server ~]# tar -zvxf maatkit-7540.tar.gz  && cd maatkit-7540 
[root@master-server maatkit-7540]# perl Makefile.PL
使用mk-table-checksum检查主从数据一致性
1）在主库服务器上运行数据一致性检查操作（也可以在从库服务器上进行数据一致性检查操作，命令跟下面一样）
[root@master-server ~]# mk-table-checksum h=192.168.1.101,u=data_check,p=check@123,P=3306 h=192.168.1.102,u=data_check,p=check@123,P=3306
使用mk-table-sync修复主从不同步的数据
顾名思义，mk-table-sync用来修复多个实例之间数据的不一致。它可以让主从的数据修复到最终一致，也可以使通过应用双写或多写的多个不相关的数据库实例修复到一致。同时它还内部集成了pt-table-checksum的校验功能，可以一边校验一边修复，也可以基于pt-table-checksum的计算结果来进行修复。
mk-table-sync工作原理
1）单行数据checksum值的计算
计算逻辑与pt-table-checksum一样，也是先检查表结构，并获取每一列的数据类型，把所有数据类型都转化为字符串，然后用concat_ws()函数进行连接，由此计算出该行的checksum值。checksum默认采用crc32计算。

2）数据块checksum值的计算
同pt-table-checksum工具一样，pt-table-sync会智能分析表上的索引，然后把表的数据split成若干个chunk，计算的时候以chunk为单位。可以理解为把chunk内所有行的数据拼接起来，再计算crc32的值，即得到该chunk的checksum值。

3）坏块检测和修复
前面两步，pt-table-sync与pt-table-checksum的算法和原理一样。再往下，就开始有所不同：
pt-table-checksum只是校验，所以它把checksum结果存储到统计表，然后把执行过的sql语句记录到binlog中，任务就算完成。语句级的复制把计算逻辑传递到从库，并在从库执行相同的计算。pt-table-checksum的算法本身并不在意从库的延迟，延迟多少都一样计算(有同事对此不理解，可以参考我的前一篇文章)，不会影响计算结果的正确性(但是我们还是会检测延迟，因为延迟太多会影响业务，所以总是要加上—max-lag来限流)。 
pt-table-sync则不同。它首先要完成chunk的checksum值的计算，一旦发现主从上同样的chunk的checksum值不同，就深入到该chunk内部，逐行比较并修复有问题的行。其计算逻辑描述如下(以修复主从结构的数据不一致为例，业务双写的情况修复起来更复杂—因为涉及到冲突解决和基准选择的问题，限于篇幅，这里不介绍)：
    1）对每一个从库，每一个表，循环进行如下校验和修复过程。
    2）对每一个chunk，在校验时加上for update锁。一旦获得锁，就记录下当前主库的show master status值。
    3）在从库上执行select master_pos_wait()函数，等待从库sql线程执行到show master status得到的位置。以此保证，主从上关于这个chunk的内容均不再改变。
    4）对这个chunk执行checksum，然后与主库的checksum进行比较。
    5）如果checksum相同，说明主从数据一致，就继续下一个chunk。
    6）如果checksum不同，说明该chunk有不一致。深入chunk内部，逐行计算checksum并比较(单行checksum比较过程与chunk的比较过程一样，单行实际是chunk的size为1的特例)。
    7）如果发现某行不一致，则标记下来。继续检测剩余行，直到这个chunk结束。
    8）对找到的主从不一致的行，采用replace into语句，在主库执行一遍以生成该行全量的binlog，并同步到从库，这会以主库数据为基准来修复从库；对于主库有的行而从库没有的行，采用replace在主库上插入(必须不能是insert)；对于从库有而主库没有的行，通过在主库执行delete来删除(pt-table-sync强烈建议所有的数据修复都只在主库进行，而不建议直接修改从库数据；但是也有特例，后面会讲到)。
    9）直到修复该chunk所有不一致的行。继续检查和修复下一个chunk。
   10）直到这个从库上所有的表修复结束。开始修复下一个从库。

重要选项
--print                    显示同步需要执行的语句
--execute               执行数据同步
--charset=utf8        设置字符集，避免从库乱码。

实例说明：
mk-table-sync的工作方式是：先一行一行检查主从库的表是否一样，如果哪里不一样，就执行删除，更新，插入等操作，使其达到一致。
通过上面mk-table-checksum的检查结果可以看出，同步的两个库huanqiu和huanpc的数据并不一致，这时就可以使用mk-table-sync进行数据修复了。

数据修复命令如下：（如果mysql端口是默认的3306，则下面命令中的P=3306可以省略）
由于上面在mk-table-checksum检查时用的data_check只有select权限，权限太小，不能用于mk-table-sync修复数据只用。
所以还需要在主库和从库数据库里创建用于mk-table-sync修复数据之用的账号权限

