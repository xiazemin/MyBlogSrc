---
title: innobackupex
layout: post
category: storage
author: 夏泽民
---
xtrabackup :
这个备份工具是挺好的，但是有缺陷，只可以备份innodb；但是我们也需要备份myisam，然后就出来了一个工具：innobackupex
<!-- more -->
一、innobackupex 备份:
1.1 查看数据目录：

[[email protected]03 ~]# ls /data/mysql/
auto.cnf  db1  ibdata1  ib_logfile0  ib_logfile1  mysql  mysql2  performance_schema  test  zhdy01  zhdy02  zhdy03  zhdy-03.err  zhdy-03.pid
其实我们完全可以使用mysqldump去备份myisam和innodb，但是速度有点慢，如果遇到大的数据库很浪费时间！

下面我们就对如上的一些数据进行备份：

1.2 安装percona-xtrabackup工具：

[[email protected]03 ~]# rpm -ivh http://www.percona.com/downloads/percona-release/redhat/0.1-3/percona-release-0.1-3.noarch.rpm

[[email protected]03 ~]# yum list | grep percona

[[email protected]03 ~]# yum install percona-xtrabackup -y
1.3 创建一个备份的用户：

先使用root账户登录;

mysql> GRANT RELOAD, LOCK TABLES, REPLICATION CLIENT ON *.* TO 'bakuser'@'localhost' identified by 'zhangduanya';
Query OK, 0 rows affected (0.00 sec)

创建一个bakuser，且授予RELOAD,LOCK TABLES,REPLICATION CLIENT权限。

mysql> flush privileges;
Query OK, 0 rows affected (0.00 sec)

刷新一下权限；
1.4 备份：

首先创建一个备份的目录；
[root@zhdy-03 ~]# mkdir -p /data/backup/

[root@zhdy-03 ~]# innobackupex --defaults-file=/etc/my.cnf --user=bakuser --password='zhangduanya' -S /tmp/mysql.sock  /data/backup

[root@zhdy-03 backup]# du -sh *
92M	2017-08-23_21-23-46
如果在备份的时候有任何的错误，它会自动的停止，并输出错误的信息！

1.5 备份对比：

[[email protected]03 backup]# ls /data/mysql/
auto.cnf  db1  ibdata1  ib_logfile0  ib_logfile1  mysql  mysql2  performance_schema  test  zhdy01  zhdy02  zhdy03  zhdy-03.err  zhdy-03.pid

[[email protected]03 backup]# ls /data/backup/2017-08-23_21-23-46/
backup-my.cnf  ibdata1  mysql2              test                    xtrabackup_info     zhdy01  zhdy03
db1            mysql    performance_schema  xtrabackup_checkpoints  xtrabackup_logfile  zhdy02
其实备份的文件+目录几乎是一样的，但是是不可以直接恢复使用的！

二、innobackupex 恢复：
2.1 模拟数据库被删除：

先停掉数据库；
[[email protected] backup]# /etc/init.d/mysqld stop 
Shutting down MySQL.. SUCCESS! 

[[email protected]03 backup]# mv /data/mysql /data/mysql.bak
[[email protected]03 backup]# ls /data/
backup  mysql.bak  wwwroot
[[email protected]03 backup]# mkdir -p /data/mysql
[[email protected]03 backup]# chown -R mysql.mysql /data/mysql
2.2 恢复数据：

[[email protected]03 backup]# innobackupex --use-memory=512M --apply-log 2017-08-23_21-23-46/

-use-memory=512M：意思是恢复数据指定使用的内存为512M；（因为这是虚拟机，所以只是测试，要是线上的服务器64G我们可以使用32G来恢复数据，这样速度会更快些）
--apply-log:指定需要恢复的日志文件
如上我们只是初始化了一下；

2.3 现在进行恢复：

[[email protected]03 backup]# innobackupex --defaults-file=/etc/my.cnf --copy-back ./2017-08-23_21-23-46/

即可恢复咱们的备份到mysql目录。
再次检查：

[[email protected]03 data]# du -sh *
190M	backup
188M	mysql
188M	mysql.bak
132M	wwwroot
三、innobackupex增量备份
3.1 先全量
innobackupex --defaults-file=/etc/my.cnf --user=bakuser --password=zhangduanya /data/backup

xtrabackup: Transaction log of lsn (3037472) to (3037472) was copied.
170824 22:57:58 completed OK!

-----------
等待出现如上消息，意味着我们已经对全量备份完毕。
查看全量备份：

[[email protected]03 ~]# ls /data/backup/
2017-08-24_22-57-43

[[email protected]03 ~]# du -sh /data/backup/
92M	/data/backup/
3.2 创建增量备份

再开始之前，我们先模拟增加一个库，这个就是我们所谓的增加的数据！

[[email protected]03 ~]# mysql -uroot -pzhangduanya -e "create database db123"      //创建一个db123库

[[email protected]03 ~]# mysql -uroot -pzhangduanya db123 < /tmp/mysqlbak.sql   //把之前备份的数据恢复得到db123库
第一次增量备份：

[[email protected] ~]# innobackupex --user=bakuser --password='zhangduanya' --incremental /data/backup --incremental-basedir /data/backup/2017-08-24_22-57-43

[[email protected] backup]# du -sh *
92M	2017-08-24_22-57-43
16M	2017-08-24_23-10-21

---------------------------
第一次增量备份的数据只有16M；
3.3 模拟执行第二次增量备份

[[email protected] backup]# mysql -uroot -pzhangduanya -e "create database lalala"

[[email protected]03 backup]# mysql -uroot -pzhangduanya lalala < /tmp/mysqlbak.sql 

[[email protected]03 backup]# ls /data/mysql/
auto.cnf  db123    ib_logfile0  lalala  mysql2              test    zhdy02  zhdy-03.err
db1       ibdata1  ib_logfile1  mysql   performance_schema  zhdy01  zhdy03  zhdy-03.pid

找到我们创建的两个db123，和lalala

[[email protected]03 backup]# innobackupex --user=bakuser --password='zhangduanya' --incremental /data/backup --incremental-basedir /data/backup/2017-08-24_23-10-21

[[email protected]03 backup]# du -sh *
92M	2017-08-24_22-57-43
16M	2017-08-24_23-10-21
17M	2017-08-24_23-41-33

-----------------------------
这里有个注意点，也是困扰很多人的一个关键操作，我们再次做增量备份的时候要基于刚刚已经做了的基础上面再次增量，也就是2017-08-24_23-10-21。也即是说这样这次的增量里面才会有刚刚咱们添加的db123库的信息。
四、增量备份的恢复
4.1 为了还原真实性，我模拟删除数据库，并且停掉mysql，利用咱们已经备份的数据去恢复它。

[[email protected] backup]# /etc/init.d/mysqld stop
Shutting down MySQL.. SUCCESS! 

[[email protected]03 data]# ls
backup  mysql  wwwroot

[[email protected]03 data]# mv /data/mysql/ /data/mysqlbak

[[email protected]03 data]# ls /data/
backup  mysqlbak  wwwroot

[[email protected]03 data]# mkdir /data/mysql

[[email protected]03 data]# ls 
backup  mysql  mysqlbak  wwwroot
为了不容易混淆，我先把backup目录中的这些备份展示出来：

[[email protected]03 backup]# ls
2017-08-24_22-57-43  2017-08-24_23-10-21  2017-08-24_23-41-33
4.2 先初始化全量备份：

[[email protected]03 data]# innobackupex --apply-log --redo-only /data/backup/2017-08-24_22-57-43

----
/data/backup/2017-08-24_22-57-43：此为咱们第一次全量备份的数据。
4.3 初始化整合第一次的增量：

[[email protected]03 backup]# innobackupex --apply-log --redo-only /data/backup/2017-08-24_22-57-43 --incremental-dir=/data/backup/2017-08-24_23-10-21
4.4 初始化整合第二次的增量：

[[email protected]03 backup]# innobackupex --apply-log --redo-only /data/backup/2017-08-24_22-57-43 --incremental-dir=/data/backup/2017-08-24_23-41-33
4.5 再次把整合好的增量再次初始化一下：

[[email protected]03 backup]# innobackupex --apply-log  /data/backup/2017-08-24_22-57-43
4.6 最后一步恢复：

但是在最后一步出错了：

[[email protected]03 backup]# innobackupex --copy-back  /data/backup/2017-08-24_22-57-43/
170825 00:12:44 innobackupex: Starting the copy-back operation

IMPORTANT: Please check that the copy-back run completes successfully.
           At the end of a successful copy-back run innobackupex
           prints "completed OK!".

innobackupex version 2.3.6 based on MySQL server 5.6.24 Linux (x86_64) (revision id: )
Error: datadir must be specified.
其原因是，我没有定义/etc/my.cnf中的datadir

[mysqld]

datadir = /data/mysql
这样就可以了！

4.7 然后开始恢复！

[[email protected]03 ~]# innobackupex --copy-back  /data/backup/2017-08-24_22-57-43/
再次检查数据：

[root@zhdy-03 ~]# ls /data/mysql
db1  db123  ibdata1  ib_logfile0  ib_logfile1  lalala  mysql  mysql2  performance_schema  test  xtrabackup_info  zhdy01  zhdy02  zhdy03
刚刚创建的db123和lalala也已经全部恢复！

4.8 启动mysql报错：

[[email protected] ~]# /etc/init.d/mysqld start
Starting MySQL.Logging to '/data/mysql/zhdy-03.err'.
. ERROR! The server quit without updating PID file (/data/mysql/zhdy-03.pid).
还记得刚刚咱们模拟删除库，自己创建的/data/mysqlm吗？

[[email protected]03 ~]# ls -ld /data/mysql
drwxr-xr-x. 2 root root 6 Aug 24 23:50 /data/mysql
[[email protected]03 ~]# chown -R mysql.mysql /data/mysql

-----
所属者和所属组都属于root，修改为mysql即可。
再次启动：

[[email protected] ~]# /etc/init.d/mysqld start
Starting MySQL. SUCCESS!

 Xtrabackup是一个对InnoDB做数据备份的工具，支持在线热备份（备份时不影响数据读写），是商业备份工具InnoDB Hotbackup的一个很好的替代品。它能对InnoDB和XtraDB存储引擎的数据库非阻塞地备份（对于MyISAM的备份同样需要加表锁）。XtraBackup支持所有的Percona Server、MySQL、MariaDB和Drizzle。几年前使用过，但现在忘记的差不多了，所以就重新拾起看看。

xtrabackup有两个主要的工具：xtrabackup、innobackupex

(1).xtrabackup只能备份InnoDB和XtraDB 两种数据表

(2).innobackupex则封装了xtrabackup,同时可以备份MyISAM数据表
Innobackupex完整备份后生成了几个重要的文件：
xtrabackup_binlog_info：记录当前最新的LOG Position

xtrabackup_binlog_pos_innodb:innodb log postion

xtrabackup_checkpoints: 存放备份的起始位置beginlsn和结束位置endlsn,增量备份需要这个lsn[增量备份可以在这里面看from和to两个值的变化]
Xtrabackup特点：
复制代码
(1)备份过程快速、可靠

(2)备份过程不会打断正在执行的事务

(3)能够基于压缩等功能节约磁盘空间和流量

(4)自动实现备份检验

(5)还原速度快