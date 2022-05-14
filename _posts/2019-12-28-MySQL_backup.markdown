---
title: MySQL_backup
layout: post
category: storage
author: 夏泽民
---
Mysql最常用的三种备份工具分别是mysqldump、Xtrabackup（innobackupex工具）、lvm-snapshot快照。


https://github.com/zonezoen/MySQL_backup
全量备份
/usr/bin/mysqldump -uroot -p123456  --lock-all-tables --flush-logs test > /home/backup.sql

，其功能是将 test 数据库全量备份。其中 MySQL 用户名为：root 密码为：123456备份的文件路径为：/home （当然这个路径也是可以按照个人意愿修改的。）备份的文件名为：backup.sql 参数 —flush-logs：使用一个新的日志文件来记录接下来的日志参数 —lock-all-tables：锁定所有数据库
<!-- more -->
恢复全量备份
mysql -h localhost -uroot -p123456 < bakdup.sql
或者

mysql> source /path/backup/bakdup.sql

定时备份
输入如下命令，进入 crontab 定时任务编辑界面：
crontab -e
添加如下命令，其意思为：每分钟执行一次备份脚本：
* * * * * sh /usr/your/path/mysqlBackup.sh
复制代码每五分钟执行 ：
*/5 * * * * sh /usr/your/path/mysqlBackup.sh
每小时执行：
0 * * * * sh /usr/your/path/mysqlBackup.sh
复制代码每天执行：
0 0 * * * sh /usr/your/path/mysqlBackup.sh
每周执行：
0 0 * * 0 sh /usr/your/path/mysqlBackup.sh
复制代码每月执行：
0 0 1 * * sh /usr/your/path/mysqlBackup.sh
每年执行:
0 0 1 1 * sh /usr/your/path/mysqlBackup.sh

增量备份
首先在进行增量备份之前需要查看一下配置文件，查看 log_bin 是否开启，因为要做增量备份首先要开启 log_bin 。首先，进入到 myslq 命令行，输入如下命令：
show variables like '%log_bin%';
复制代码
如下命令所示，则为未开启
mysql> show variables like '%log_bin%';
+---------------------------------+-------+
| Variable_name                   | Value |
+---------------------------------+-------+
| log_bin                         | OFF   |
| log_bin_basename                |       |
| log_bin_index                   |       |
| log_bin_trust_function_creators | OFF   |
| log_bin_use_v1_row_events       | OFF   |
| sql_log_bin                     | ON    |
+---------------------------------+-------+
复制代码
修改 MySQL 配置项到如下代码段：vim /etc/mysql/mysql.conf.d/mysqld.cnf
# Copyright (c) 2014, 2016, Oracle and/or its affiliates. All rights reserved.
#
# This program is free software; you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation; version 2 of the License.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
#
# You should have received a copy of the GNU General Public License
# along with this program; if not, write to the Free Software
# Foundation, Inc., 51 Franklin St, Fifth Floor, Boston, MA  02110-1301 USA

#
# The MySQL  Server configuration file.
#
# For explanations see
# http://dev.mysql.com/doc/mysql/en/server-system-variables.html

[mysqld]
pid-file	= /var/run/mysqld/mysqld.pid
socket		= /var/run/mysqld/mysqld.sock
datadir		= /var/lib/mysql
#log-error	= /var/log/mysql/error.log
# By default we only accept connections from localhost
#bind-address	= 127.0.0.1
# Disabling symbolic-links is recommended to prevent assorted security risks
symbolic-links=0

#binlog setting，开启增量备份的关键
log-bin=/var/lib/mysql/mysql-bin
server-id=123454
复制代码
修改之后，重启 mysql 服务，输入：
show variables like '%log_bin%';
复制代码
状态如下：
mysql> show variables like '%log_bin%';
+---------------------------------+--------------------------------+
| Variable_name                   | Value                          |
+---------------------------------+--------------------------------+
| log_bin                         | ON                             |
| log_bin_basename                | /var/lib/mysql/mysql-bin       |
| log_bin_index                   | /var/lib/mysql/mysql-bin.index |
| log_bin_trust_function_creators | OFF                            |
| log_bin_use_v1_row_events       | OFF                            |
| sql_log_bin                     | ON                             |
+---------------------------------+--------------------------------+
复制代码
好了，做好了充足的准备，那我们就开始学习增量备份了。
查看当前使用的 mysql_bin.000*** 日志文件，
show master status;
复制代码
状态如下：
mysql> show master status;
+------------------+----------+--------------+------------------+-------------------+
| File             | Position | Binlog_Do_DB | Binlog_Ignore_DB | Executed_Gtid_Set |
+------------------+----------+--------------+------------------+-------------------+
| mysql-bin.000015 |      610 |              |                  |                   |
+------------------+----------+--------------+------------------+-------------------+
复制代码
当前正在记录日志的文件名为 mysql-bin.000015 。
当前数据库中有如下数据：
mysql> select * from users;
+-------+------+----+
| name  | sex  | id |
+-------+------+----+
| zone  | 0    |  1 |
| zone1 | 1    |  2 |
| zone2 | 0    |  3 |
+-------+------+----+
复制代码
我们插入一条数据：
insert into `zone`.`users` ( `name`, `sex`, `id`) values ( 'zone3', '0', '4');
复制代码
查看效果：
mysql> select * from users;
+-------+------+----+
| name  | sex  | id |
+-------+------+----+
| zone  | 0    |  1 |
| zone1 | 1    |  2 |
| zone2 | 0    |  3 |
| zone3 | 0    |  4 |
+-------+------+----+
复制代码
我们执行如下命令，使用新的日志文件：
mysqladmin -uroot -123456 flush-logs
复制代码
日志文件从 mysql-bin.000015 变为  mysql-bin.000016，而 mysql-bin.000015 则记录着刚刚 insert 命令的日志。上句代码的效果如下：
mysql> show master status;
+------------------+----------+--------------+------------------+-------------------+
| File             | Position | Binlog_Do_DB | Binlog_Ignore_DB | Executed_Gtid_Set |
+------------------+----------+--------------+------------------+-------------------+
| mysql-bin.000016 |      154 |              |                  |                   |
+------------------+----------+--------------+------------------+-------------------+
复制代码
那么到现在为止，其实已经完成了增量备份了。
恢复增量备份
那么现在将刚刚插入的数据删除，效果如下：
delete from `zone`.`users` where `id`='4' 

mysql> select * from users;
+-------+------+----+
| name  | sex  | id |
+-------+------+----+
| zone  | 0    |  1 |
| zone1 | 1    |  2 |
| zone2 | 0    |  3 |
+-------+------+----+
复制代码
那么现在就是重点时间了，从 mysql-bin.000015 中恢复数据：
mysqlbinlog /var/lib/mysql/mysql-bin.000015 | mysql -uroot -p123456 zone;
复制代码
上一句代码指定了，需要恢复的 mysql_bin 文件，指定了用户名：root 、密码：123456 、数据库名：zone。效果如下：
mysql> select * from users;
+-------+------+----+
| name  | sex  | id |
+-------+------+----+
| zone  | 0    |  1 |
| zone1 | 1    |  2 |
| zone2 | 0    |  3 |
| zone3 | 0    |  4 |
+-------+------+----+
复制代码
OK，整一个增量备份的操作流程都在这里了，那么我们如何将它写成脚本文件呢，代码如下：
#!/bin/bash
#在使用之前，请提前创建以下各个目录
backupDir=/usr/local/work/backup/daily
#增量备份时复制mysql-bin.00000*的目标目录，提前手动创建这个目录
mysqlDir=/var/lib/mysql
#mysql的数据目录
logFile=/usr/local/work/backup/bak.log
BinFile=/var/lib/mysql/mysql-bin.index
#mysql的index文件路径，放在数据目录下的

mysqladmin -uroot -p123456 flush-logs
#这个是用于产生新的mysql-bin.00000*文件
# wc -l 统计行数
# awk 简单来说awk就是把文件逐行的读入，以空格为默认分隔符将每行切片，切开的部分再进行各种分析处理。
Counter=`wc -l $BinFile |awk '{print $1}'`
NextNum=0
#这个for循环用于比对$Counter,$NextNum这两个值来确定文件是不是存在或最新的
for file in `cat $BinFile`
do
    base=`basename $file`
    echo $base
    #basename用于截取mysql-bin.00000*文件名，去掉./mysql-bin.000005前面的./
    NextNum=`expr $NextNum + 1`
    if [ $NextNum -eq $Counter ]
    then
        echo $base skip! >> $logFile
    else
        dest=$backupDir/$base
        if(test -e $dest)
        #test -e用于检测目标文件是否存在，存在就写exist!到$logFile去
        then
            echo $base exist! >> $logFile
        else
            cp $mysqlDir/$base $backupDir
            echo $base copying >> $logFile
         fi
     fi
done
echo `date +"%Y年%m月%d日 %H:%M:%S"` $Next Bakup succ! >> $logFile

#NODE_ENV=$backUpFolder@$backUpFileName /root/node/v8.11.3/bin/node /usr/local/upload.js

一、binlog二进制日志通常作为备份的重要资源，所以再说备份方案之前先总结一下binlog日志~~
1.binlog日志内容
1）引起mysql服务器改变的任何操作。
2）复制功能依赖于此日志。
3）slave服务器通过复制master服务器的二进制日志完成主从复制，在执行之前保存于中继日志(relay log)中。
4）slave服务器通常可以关闭二进制日志以提升性能。

2.binlog日志文件的文件表现形式
1）默认在安装目录下，存在mysql-bin.00001, mysql-bin.00002的二进制文件（binlog日志文件名依据my.cnf配置中的log-bin参数后面的设置为准）
2）还有mysql-bin.index用来记录被mysql管理的二进制文件列表
3）如果需要删除二进制日志时，切勿直接删除二进制文件，这样会使得mysql管理混乱。

3.binlog日志文件查看相关mysql命令
1）SHOW MASTER STATUS ; 查看正在使用的二进制文件
MariaDB [(none)]> SHOW MASTER STATUS ;
+------------------+----------+--------------+------------------+
| File | Position | Binlog_Do_DB | Binlog_Ignore_DB |
+------------------+----------+--------------+------------------+
| mysql-bin.000003 | 245 | | |
+------------------+----------+--------------+------------------+
2）FLUSH LOGS; 手动滚动二进制日志
MariaDB [(none)]> FLUSH LOGS;
MariaDB [(none)]> SHOW MASTER STATUS ;
+------------------+----------+--------------+------------------+
| File | Position | Binlog_Do_DB | Binlog_Ignore_DB |
+------------------+----------+--------------+------------------+
| mysql-bin.000004 | 245 | | |
+------------------+----------+--------------+------------------+
滚动以后，mysql重新创建一个新的日志mysql-bin.000004
3）SHOW BINARY LOGS 显示使用过的二进制日志文件
MariaDB [(none)]> SHOW BINARY LOGS ;
+------------------+-----------+
| Log_name | File_size |
+------------------+-----------+
| mysql-bin.000001 | 30373 |
| mysql-bin.000002 | 1038814 |
| mysql-bin.000003 | 288 |
| mysql-bin.000004 | 245 |
4）SHOW BINLOG EVENTS 以表的形式查看二进制文件
命令格式：SHOW BINLOG EVENTS [IN 'log_name'] [FROM pos] [LIMIT [offset,] row_count]
MariaDB [(none)]> SHOW BINLOG EVENTS IN 'mysql-bin.000001' \G;
*************************** 99. row ***************************
Log_name: mysql-bin.000001
Pos: 30225
Event_type: Query
Server_id: 1
End_log_pos: 30354
Info: use `mysql`; DROP TEMPORARY TABLE `tmp_proxies_priv` /* generated by server */

4.MySQL二进制文件读取工具mysqlbinlog
命令格式：mysqlbinlog [参数] log-files
有以下四种参数选择：
--start-datetime
--stop-datetime
--start-position
--stop-position
[root@test-huanqiu ~]# mysqlbinlog --start-position 30225 --stop-position 30254 mysql-bin.000001
截取一下结果：
# at 30225
#151130 12:43:35 server id 1 end_log_pos 30354 Querythread_id=1exec_time=0error_code=0
use `mysql`/*!*/;
SET TIMESTAMP=1448858615/*!*/;
SET @@session.pseudo_thread_id=1/*!*/

根据以上截取结果第二行，进行解释二进制日志内容
1）时间点： 151130 12:43:35
2）服务器ID： server id 1
服务器ID主要用于标记日志产生的服务器，主要用于双主模型中，互为主从，确保二进制文件不会被相互循环复制
3）记录类型： Query
4) 线程号： thread_id = 1
5) 语句的时间戳和写入二进制日志文件的时间差； exec_time=0
6) 事件内容
7）事件位置 #at 30225
8) 错误代码 error_code=0
9) 事件结束位置 end_log_pos也就是下一事件开始的位置

5.二进制日志格式
由bin_log_format={statement|row|mixed}定义
1）statement: 基于语句，记录生成数据的语句
缺点在于如果当时插入信息为函数生成，有可能不同时间点执行结果不一样，
例如： INSERT INTO t1 VALUE (CURRENT_DATE());
2）row: 基于行数据
缺点在于，有时候数据量会过大
3）mixed： 混合模式，又mysql自行决定何时使用statement, 何时使用row 模式

6.二进制相关参数总结
1）log_bin = {ON|OFF}
还可以是个文件路径，自定义binlog日志文件名，使用“log_bin=“或“log-bin=“都可以，主要用于控制全局binlog的存放位置和是否开启binlog日志功能。
比如：log_bin=mysql-bin 或者 log-bin=mysql-bin，这样binlog日志默认会和mysql数据放在同一目录下。
2） log_bin_trust_function_creators
是否记录在
3） sql_log_bin = {ON|OFF}
会话级别是否关闭binlog， 如果关闭当前会话内的操作将不会记录
4） sync_binlog 是否马上同步事务类操作到二进制日志中
5） binlog_format = {statement|row|mixed} 二进制日志的格式，上面单独提到了
6） max_binlog_cache_size =
二进制日志缓冲空间大小，仅用于缓冲事务类的语句；
7） max_binlog_stmt_cache_size =
语句缓冲，非事务类和事务类共用的空间大小
8） max_binlog_size =
二进制日志文件上限，超过上限后则滚动
9) 删除二进制日志
命令格式：PURGE { BINARY | MASTER } LOGS { TO 'log_name' | BEFORE datetime_expr }
MariaDB> PURGE BINARY LOGS TO 'mysql-bin.010';
MariaDB> PURGE BINARY LOGS BEFORE '2016-11-02 22:46:26';
建议：切勿将二进制日志与数据文件放在一同设备；可以将binlog日志实时备份到远程设备上，以防出现机器故障进行数据恢复；

二、接下来说下binlog二进制日志备份和恢复
1.为什么做备份：
（1）灾难恢复
（2）审计，数据库在过去某一个时间点是什么样的
（3）测试

2.备份的目的：
（1）用于恢复数据
（2）备份结束后，需要周期性的做恢复测试

3.备份类型：
（1）根据备份时，服务器是否在线
1）冷备(cold backup)： 服务器离线，读写操作都不能进行
2）温备份： 全局施加共享锁，只能读不能写
3）热备(hot backup)：数据库在线，读写照样进行
（2）根据备份时的数据集分类
1）完全备份(full backup)
2）部分备份(partial backup)
（3）根据备份时的接口
1）物理备份（physical backup）：直接复制数据文件 ，打包归档
特点：
不需要额外工具，直接归档命令即可，但是跨平台能力比较差；如果数据量超过几十个G，则适用于物理备份
2）逻辑备份(logical backup)： 把数据抽取出来保存在sql脚本中
特点：
可以使用文本编辑器编辑；导入方便，直接读取sql语句即可；逻辑备份恢复时间慢，占据空间大；无法保证浮点数的精度；恢复完数据库后需要重建索引。
（4）根据备份整个数据还是变化数据
1） 全量备份 full backup
2） 增量备份 incremental backup
在不同时间点起始备份一段数据，比较节约空间；针对的是上一次备份后有变化的数据，备份数据少，备份快，恢复慢
3） 差异备份 differential backup
备份从每个时间点到上一次全部备份之间的数据，随着时间增多二增多；比较容易恢复；对于很大的数据库，可以考虑主从模型，备份从服务器的内容。针对的是上一次全量备份后有变化的数据，备份数据多，备份慢，恢复快。
（5）备份策略，需要考虑因素如下
备份方式
备份实践
备份成本
锁时间
时长
性能开销
恢复成本
恢复时长
所能够容忍丢失的数据量
（6）备份内容
1）数据库中的数据
2）配置文件
3）mysql中的代码： 存储过程，存储函数，触发器
4）OS 相关的配置文件，chrontab 中的备份策略脚本
5）如果是主从复制的场景中： 跟复制相关的信息
6）二进制日志文件需要定期备份，一旦发现二进制文件出现问题，需马上对数据进行完全备份

(7)Mysql最常用的三种备份工具：
1）mysqldump：
通常为小数据情况下的备份
innodb： 热备，温备
MyISAM, Aria: 温备
单线程备份恢复比较慢
2）Xtrabackup（通常用innobackupex工具）：
备份mysql大数据
InnoDB热备，增量备份；
MyISAM温备，不支持增量，只有完全备份
属于物理备份，速度快；
3）lvm-snapshot：
接近于热备的工具：因为要先请求全局锁，而后创建快照，并在创建快照完成后释放全局锁；
使用cp、tar等工具进行物理备份；
备份和恢复速度较快；
很难实现增量备份，并且请求全局需要等待一段时间，在繁忙的服务器上尤其如此；


除此之外，还有其他的几个备份工具：
-->mysqldumper: 多线程的mysqldump
-->SELECT clause INTO OUTFILE '/path/to/somefile' LOAD DATA INFILE '/path/from/somefile'
部分备份工具， 不会备份关系定义，仅备份表中的数据；
逻辑备份工具，快于mysqldump，因为不备份表格式信息
-->mysqlhotcopy: 接近冷备，基本没用

 

mysqldump工具基本使用
1. mysqldump [OPTIONS] database [tables…]
还原时库必须存在，不存在需要手动创建
    --all-databases:　备份所有库
    --databases db1 db2 ...: 备份指定的多个库，如果使用此命令，恢复时将不用手动创建库。或者是-B db1 db2 db3 ....
    --lock-all-tables：请求锁定所有表之后再备份，对MyISAM、InnoDB、Aria做温备
    --lock-table: 对正在备份的表加锁，但是不建议使用，如果其它表被修改，则备份后表与表之间将不同步
    --single-transaction: 能够对InnoDB存储引擎实现热备；
启动一个很大的大事物，基于MOCC可以保证在事物内的表版本一致
自动加锁不需要，再加--lock-table, 可以实现热备
备份代码：
   --events: 备份事件调度器代码
   --routines: 备份存储过程和存储函数
   --triggers：备份触发器
备份时滚动日志：
   --flush-logs: 备份前、请求到锁之后滚动日志；
方恢复备份时间点以后的内容
复制时的同步位置标记：主从架构中的，主服务器数据。效果相当于标记一个时间点。
   --master-data=[0|1|2]
   0: 不记录
   1：记录为CHANGE MASTER语句
   2：记录为注释的CHANGE MASTER语句

2. 使用mysqldump备份大体过程：
1) 请求锁：--lock-all-tables或使用–singe-transaction进行innodb热备；
2) 滚动日志：--flush-logs
3) 选定要备份的库：--databases
4) 记录二进制日志文件及位置：--master-data=
FLUSH TABLES5 WITH READ LOCK;

3. 恢复：
恢复过程无需写到二进制日志中
建议：关闭二进制日志，关闭其它用户连接；

4. 备份策略：基于mysqldump
备份：mysqldump+二进制日志文件；（“mysqldump >”）
周日做一次完全备份：备份的同时滚动日志
周一至周六：备份二进制日志；
恢复：（“mysql < ”）或在mysql数据库中直接执行“source sql备份文件;”进行恢复。如果sql执行语句比较多，可以将sql语句放在一个文件内，将文件名命名为.sql结尾，然后在mysql数据库中使用"source 文件.sql;"命令进行执行即可！
完全备份+各二进制日志文件中至此刻的事件

5. 实例说明：
参考：Mysql备份系列（2）--mysqldump备份（全量+增量）方案操作记录

 

lvm-snapshot：基于LVM快照的备份
1.关于快照：
1）事务日志跟数据文件必须在同一个卷上；
2）刚刚创立的快照卷，里面没有任何数据，所有数据均来源于原卷
3）一旦原卷数据发生修改，修改的数据将复制到快照卷中，此时访问数据一部分来自于快照卷，一部分来自于原卷
4）当快照使用过程中，如果修改的数据量大于快照卷容量，则会导致快照卷崩溃。
5）快照卷本身不是备份，只是提供一个时间一致性的访问目录。

2.基于快照备份几乎为热备：
1）创建快照卷之前，要请求MySQL的全局锁；在快照创建完成之后释放锁；
2）如果是Inoodb引擎， 当flush tables 后会有一部分保存在事务日志中，却不在文件中。 因此恢复时候，需要事务日志和数据文件
但释放锁以后，事务日志的内容会同步数据文件中，因此备份内容并不绝对是锁释放时刻的内容，由于有些为完成的事务已经完成，但在备份数据中因为没完成而回滚。 因此需要借助二进制日志往后走一段

3.基于快照备份注意事项：
1）事务日志跟数据文件必须在同一个卷上；
2）创建快照卷之前，要请求MySQL的全局锁；在快照创建完成之后释放锁；
3）请求全局锁完成之后，做一次日志滚动；做二进制日志文件及位置标记(手动进行)；

4.为什么基于MySQL快照的备份很好？
原因如下几点：
1）几乎是热备 在大多数情况下，可以在应用程序仍在运行的时候执行备份。无需关机，只需设置为只读或者类似只读的限制。
2）支持所有基于本地磁盘的存储引擎 它支持MyISAM, Innodb, BDB，还支持 Solid, PrimeXT 和 Falcon。
3）快速备份 只需拷贝二进制格式的文件，在速度方面无以匹敌。
4）低开销 只是文件拷贝，因此对服务器的开销很细微。
5）容易保持完整性 想要压缩备份文件吗？把它们备份到磁带上，FTP或者网络备份软件 -- 十分简单，因为只需要拷贝文件即可。
6）快速恢复 恢复的时间和标准的MySQL崩溃恢复或数据拷贝回去那么快，甚至可能更快，将来会更快。
7）免费 无需额外的商业软件，只需Innodb热备工具来执行备份。

快照备份mysql的缺点：
1）需要兼容快照 -- 这是明显的。
2）需要超级用户(root) 在某些组织，DBA和系统管理员来自不同部门不同的人，因此权限各不一样。
3）停工时间无法预计，这个方法通常指热备，但是谁也无法预料到底是不是热备 -- FLUSH TABLES WITH READ LOCK 可能会需要执行很长时间才能完成。
4）多卷上的数据问题 如果你把日志放在独立的设备上或者你的数据库分布在多个卷上，这就比较麻烦了，因为无法得到全部数据库的一致性快照。不过有些系统可能能自动做到多卷快照。

5.备份与恢复的大体步骤
备份：
1）请求全局锁，并滚动日志
mysql> FLUSH TABLES WITH READ LOCK;
mysql> FLUSH LOGS;
2）做二进制日志文件及位置标记(手动进行)；
[root@test-huanqiu ~]# mysql -e 'show master status' > /path/to/orignal_volume
3）创建快照卷
[root@test-huanqiu ~]# lvcreate -L -s -n -p r /path/to/some_lv
4）释放全局锁
5）挂载快照卷并备份
6）备份完成之后，删除快照卷

恢复：
1）二进制日志保存好；
提取备份之后的所有事件至某sql脚本中；
2）还原数据，修改权限及属主属组等，并启动mysql
3）做即时点还原
4）生产环境下， 一次大型恢复后，需要马上进行一次完全备份。

备份与恢复实例说明：
环境， 实现创建了一个test_vg卷组，里面有个mylv1用来装mysql数据，挂载到/data/mysqldata

备份实例：
1. 创建备份专用的用户，授予权限FLUSH LOGS 和 LOCK TABLES
MariaDB > GRANT RELOAD,LOCK TABLES,SUPER ON *.* TO 'lvm'@'192.168.1.%' IDENTIFIED BY 'lvm';
MariaDB > FLUSH PRIVILEGES;

2. 记录备份点
[root@test-huanqiu ~]# mysql -ulvm -h192.168.1.10 -plvm -e 'SHOW MASTER STATUS' > /tmp/backup_point.txt

3. 创建快照卷并挂载快照卷
[root@test-huanqiu ~]# lvcreate -L 1G -s -n lvmbackup -p r /dev/test_vg/mylv1
[root@test-huanqiu ~]# mount -t ext4 /dev/test_vg/lvmbackup /mnt/

4. 释放锁
[root@test-huanqiu ~]# mysql -ulvm -h192.168.98.10 -plvm -e 'UNLOCK TABLES'
做一些模拟写入工作
MariaDB [test]> create database testdb2

5. 复制文件
[root@test-huanqiu ~]# cp /data/mysqldata /tmp/backup_mysqldata -r

6. 备份完成卸载，删除快照卷
[root@test-huanqiu ~]# umount /mnt
[root@test-huanqiu ~]# lvmremove /dev/test_vg/lvmbackup

还原实例：
假如整个mysql服务器崩溃，并且目录全部被删除

1. 数据文件复制回源目录
[root@test-huanqiu ~]# cp -r /tmp/backup_mysqldata/* /data/mysqldata/
MariaDB [test]> show databases ;
+--------------------+
| Database |
+--------------------+
| information_schema |
| hellodb |
| mysql |
| mysqldata |
| openstack |
| performance_schema |
| test |
+--------------------+
此时还没有testdb2， 因为这个是备份之后创建的，因此需要通过之前记录的二进制日志

2. 查看之前记录的记录点。向后还原
[root@test-huanqiu ~]# cat /tmp/backup_point.txt
FilePositionBinlog_Do_DBBinlog_Ignore_DB
mysql-bin.000001245
[root@test-huanqiu ~]# mysqlbinlog /data/binlog/mysql-bin.000001 --start-position 245 > tmp.sql
MariaDB [test]> source /data/mysqldata/tmp.sql
MariaDB [test]> show databases ;
+--------------------+
| Database |
+--------------------+
| information_schema |
| hellodb |
| mysql |
| mysqldata |
| openstack |
| performance_schema |
| test |
| testdb2 |
+--------------------+
8 rows in set (0.00 sec)
testdb2 已经被还原回来。

具体实例说明，参考：Mysql备份系列（4）--lvm-snapshot备份mysql数据(全量+增量）操作记录

 

使用Xtrabackup进行MySQL备份：

参考：Mysql备份系列（3）--innobackupex备份mysql大数据(全量+增量）操作记录

 

--------------------------------------------------------------------------------------
关于备份和恢复的几点经验之谈

备份注意：
1. 将数据和备份放在不同的磁盘设备上；异机或异地备份存储较为理想；
2. 备份的数据应该周期性地进行还原测试；
3. 每次灾难恢复后都应该立即做一次完全备份；
4. 针对不同规模或级别的数据量，要定制好备份策略；
5. 二进制日志应该跟数据文件在不同磁盘上，并周期性地备份好二进制日志文件；

从备份中恢复应该遵循步骤：
1. 停止MySQL服务器；
2. 记录服务器的配置和文件权限；
3. 将数据从备份移到MySQL数据目录；其执行方式依赖于工具；
4. 改变配置和文件权限；
5. 以限制访问模式重启服务器；mysqld的--skip-networking选项可跳过网络功能；
方法：编辑my.cnf配置文件，添加如下项：
skip-networking
socket=/tmp/mysql-recovery.sock
6. 载入逻辑备份（如果有）；而后检查和重放二进制日志；
7. 检查已经还原的数据；
8. 重新以完全访问模式重启服务器；
注释前面第5步中在my.cnf中添加的选项，并重启；

在日常运维工作中，对mysql数据库的备份是万分重要的，以防在数据库表丢失或损坏情况出现，可以及时恢复数据。

线上数据库备份场景：
每周日执行一次全量备份，然后每天下午1点执行MySQLdump增量备份.

下面对这种备份方案详细说明下：
1.MySQLdump增量备份配置
执行增量备份的前提条件是MySQL打开binlog日志功能，在my.cnf中加入
log-bin=/opt/Data/MySQL-bin
“log-bin=”后的字符串为日志记载目录，一般建议放在不同于MySQL数据目录的磁盘上。

1
2
3
4
5
6
7
8
9
-----------------------------------------------------------------------------------
mysqldump >       导出数据
mysql <           导入数据  （或者使用source命令导入数据，导入前要先切换到对应库下）
 
注意一个细节：
若是mysqldump导出一个库的数据，导出文件为a.sql，然后mysql导入这个数据到新的空库下。
如果新库名和老库名不一致，那么需要将a.sql文件里的老库名改为新库名，
这样才能顺利使用mysql命令导入数据（如果使用source命令导入就不需要修改a.sql文件了）。
-----------------------------------------------------------------------------------
2.MySQLdump增量备份
假定星期日下午1点执行全量备份，适用于MyISAM存储引擎。
[root@test-huanqiu ~]# MySQLdump --lock-all-tables --flush-logs --master-data=2 -u root -p test > backup_sunday_1_PM.sql

对于InnoDB将--lock-all-tables替换为--single-transaction
--flush-logs为结束当前日志，生成新日志文件；
--master-data=2 选项将会在输出SQL中记录下完全备份后新日志文件的名称，

用于日后恢复时参考，例如输出的备份SQL文件中含有：
CHANGE MASTER TO MASTER_LOG_FILE=’MySQL-bin.000002′, MASTER_LOG_POS=106;

3.MySQLdump增量备份其他说明：
如果MySQLdump加上–delete-master-logs 则清除以前的日志，以释放空间。但是如果服务器配置为镜像的复制主服务器，用MySQLdump –delete-master-logs删掉MySQL二进制日志很危险，因为从服务器可能还没有完全处理该二进制日志的内容。在这种情况下，使用 PURGE MASTER LOGS更为安全。

每日定时使用 MySQLadmin flush-logs来创建新日志，并结束前一日志写入过程。并把前一日志备份，例如上例中开始保存数据目录下的日志文件 MySQL-bin.000002 , ...

1.恢复完全备份
mysql -u root -p < backup_sunday_1_PM.sql

2.恢复增量备份
mysqlbinlog MySQL-bin.000002 … | MySQL -u root -p注意此次恢复过程亦会写入日志文件，如果数据量很大，建议先关闭日志功能

--compatible=name
它告诉 MySQLdump，导出的数据将和哪种数据库或哪个旧版本的 MySQL 服务器相兼容。值可以为 ansi、MySQL323、MySQL40、postgresql、oracle、mssql、db2、maxdb、no_key_options、no_tables_options、no_field_options 等，要使用几个值，用逗号将它们隔开。当然了，它并不保证能完全兼容，而是尽量兼容。

--complete-insert，-c
导出的数据采用包含字段名的完整 INSERT 方式，也就是把所有的值都写在一行。这么做能提高插入效率，但是可能会受到 max_allowed_packet 参数的影响而导致插入失败。因此，需要谨慎使用该参数，至少我不推荐。

--default-character-set=charset
指定导出数据时采用何种字符集，如果数据表不是采用默认的 latin1 字符集的话，那么导出时必须指定该选项，否则再次导入数据后将产生乱码问题。

--disable-keys
告诉 MySQLdump 在 INSERT 语句的开头和结尾增加 /*!40000 ALTER TABLE table DISABLE KEYS */; 和 /*!40000 ALTER TABLE table ENABLE KEYS */; 语句，这能大大提高插入语句的速度，因为它是在插入完所有数据后才重建索引的。该选项只适合 MyISAM 表。

--extended-insert = true|false
默认情况下，MySQLdump 开启 --complete-insert 模式，因此不想用它的的话，就使用本选项，设定它的值为 false 即可。

--hex-blob
使用十六进制格式导出二进制字符串字段。如果有二进制数据就必须使用本选项。影响到的字段类型有 BINARY、VARBINARY、BLOB。

--lock-all-tables，-x
在开始导出之前，提交请求锁定所有数据库中的所有表，以保证数据的一致性。这是一个全局读锁，并且自动关闭 --single-transaction 和 --lock-tables 选项。

--lock-tables
它和 --lock-all-tables 类似，不过是锁定当前导出的数据表，而不是一下子锁定全部库下的表。本选项只适用于 MyISAM 表，如果是 Innodb 表可以用 --single-transaction 选项。

--no-create-info，-t
只导出数据，而不添加 CREATE TABLE 语句。

--no-data，-d
不导出任何数据，只导出数据库表结构。
mysqldump --no-data --databases mydatabase1 mydatabase2 mydatabase3 > test.dump
将只备份表结构。--databases指示主机上要备份的数据库。

--opt
这只是一个快捷选项，等同于同时添加 --add-drop-tables --add-locking --create-option --disable-keys --extended-insert --lock-tables --quick --set-charset 选项。本选项能让 MySQLdump 很快的导出数据，并且导出的数据能很快导回。该选项默认开启，但可以用 --skip-opt 禁用。注意，如果运行 MySQLdump 没有指定 --quick 或 --opt 选项，则会将整个结果集放在内存中。如果导出大数据库的话可能会出现问题。

--quick，-q
该选项在导出大表时很有用，它强制 MySQLdump 从服务器查询取得记录直接输出而不是取得所有记录后将它们缓存到内存中。

--routines，-R
导出存储过程以及自定义函数。

--single-transaction
该选项在导出数据之前提交一个 BEGIN SQL语句，BEGIN 不会阻塞任何应用程序且能保证导出时数据库的一致性状态。它只适用于事务表，例如 InnoDB 和 BDB。本选项和 --lock-tables 选项是互斥的，因为 LOCK TABLES 会使任何挂起的事务隐含提交。要想导出大表的话，应结合使用 --quick 选项。

--triggers
同时导出触发器。该选项默认启用，用 --skip-triggers 禁用它。

跨主机备份
使用下面的命令可以将host1上的sourceDb复制到host2的targetDb，前提是host2主机上已经创建targetDb数据库：
-C 指示主机间的数据传输使用数据压缩
mysqldump --host=host1 --opt sourceDb| mysql --host=host2 -C targetDb

结合Linux的cron命令实现定时备份
比如需要在每天凌晨1:30备份某个主机上的所有数据库并压缩dump文件为gz格式
30 1 * * * mysqldump -u root -pPASSWORD --all-databases | gzip > /mnt/disk2/database_`date '+%m-%d-%Y'`.sql.gz

一个完整的Shell脚本备份MySQL数据库示例。比如备份数据库opspc
[root@test-huanqiu ~]# vim /root/backup.sh
#!bin/bash
echo "Begin backup mysql database"
mysqldump -u root -ppassword opspc > /home/backup/mysqlbackup-`date +%Y-%m-%d`.sql
echo "Your database backup successfully completed"

[root@test-huanqiu ~]# crontab -e
30 1 * * * /bin/bash -x /root/backup.sh > /dev/null 2>&1

mysqldump全量备份+mysqlbinlog二进制日志增量备份
1）从mysqldump备份文件恢复数据会丢失掉从备份点开始的更新数据，所以还需要结合mysqlbinlog二进制日志增量备份。
首先确保已开启binlog日志功能。在my.cnf中包含下面的配置以启用二进制日志：
[mysqld]
log-bin=mysql-bin

2）mysqldump命令必须带上--flush-logs选项以生成新的二进制日志文件：
mysqldump --single-transaction --flush-logs --master-data=2 > backup.sql
其中参数--master-data=[0|1|2]
0: 不记录
1：记录为CHANGE MASTER语句
2：记录为注释的CHANGE MASTER语句

mysqldump全量+增量备份方案的具体操作可参考下面两篇文档：
数据库误删除后的数据恢复操作说明
解说mysql之binlog日志以及利用binlog日志恢复数据

--------------------------------------------------------------------------
下面分享一下自己用过的mysqldump全量和增量备份脚本

应用场景：
1）增量备份在周一到周六凌晨3点，会复制mysql-bin.00000*到指定目录；
2）全量备份则使用mysqldump将所有的数据库导出，每周日凌晨3点执行，并会删除上周留下的mysq-bin.00000*，然后对mysql的备份操作会保留在bak.log文件中。

脚本实现：
1）全量备份脚本（假设mysql登录密码为123456；注意脚本中的命令路径）：
[root@test-huanqiu ~]# vim /root/Mysql-FullyBak.sh
#!/bin/bash
# Program
# use mysqldump to Fully backup mysql data per week!
# History
# Path
BakDir=/home/mysql/backup
LogFile=/home/mysql/backup/bak.log
Date=`date +%Y%m%d`
Begin=`date +"%Y年%m月%d日 %H:%M:%S"`
cd $BakDir
DumpFile=$Date.sql
GZDumpFile=$Date.sql.tgz
/usr/local/mysql/bin/mysqldump -uroot -p123456 --quick --events --all-databases --flush-logs --delete-master-logs --single-transaction > $DumpFile
/bin/tar -zvcf $GZDumpFile $DumpFile
/bin/rm $DumpFile
Last=`date +"%Y年%m月%d日 %H:%M:%S"`
echo 开始:$Begin 结束:$Last $GZDumpFile succ >> $LogFile
cd $BakDir/daily
/bin/rm -f *

2）增量备份脚本（脚本中mysql的数据存放路径是/home/mysql/data，具体根据自己的实际情况进行调整）
[root@test-huanqiu ~]# vim /root/Mysql-DailyBak.sh
#!/bin/bash
# Program
# use cp to backup mysql data everyday!
# History
# Path
BakDir=/home/mysql/backup/daily                     //增量备份时复制mysql-bin.00000*的目标目录，提前手动创建这个目录
BinDir=/home/mysql/data                                   //mysql的数据目录
LogFile=/home/mysql/backup/bak.log
BinFile=/home/mysql/data/mysql-bin.index           //mysql的index文件路径，放在数据目录下的
/usr/local/mysql/bin/mysqladmin -uroot -p123456 flush-logs
#这个是用于产生新的mysql-bin.00000*文件
Counter=`wc -l $BinFile |awk '{print $1}'`
NextNum=0
#这个for循环用于比对$Counter,$NextNum这两个值来确定文件是不是存在或最新的
for file in `cat $BinFile`
do
    base=`basename $file`
    #basename用于截取mysql-bin.00000*文件名，去掉./mysql-bin.000005前面的./
    NextNum=`expr $NextNum + 1`
    if [ $NextNum -eq $Counter ]
    then
        echo $base skip! >> $LogFile
    else
        dest=$BakDir/$base
        if(test -e $dest)
        #test -e用于检测目标文件是否存在，存在就写exist!到$LogFile去
        then
            echo $base exist! >> $LogFile
        else
            cp $BinDir/$base $BakDir
            echo $base copying >> $LogFile
         fi
     fi
done
echo `date +"%Y年%m月%d日 %H:%M:%S"` $Next Bakup succ! >> $LogFile

3）设置crontab任务，执行备份脚本。先执行的是增量备份脚本，然后执行的是全量备份脚本：
[root@test-huanqiu ~]# crontab -e
#每个星期日凌晨3:00执行完全备份脚本
0 3 * * 0 /bin/bash -x /root/Mysql-FullyBak.sh >/dev/null 2>&1
#周一到周六凌晨3:00做增量备份
0 3 * * 1-6 /bin/bash -x /root/Mysql-DailyBak.sh >/dev/null 2>&1

4）手动执行上面两个脚本，测试下备份效果
[root@test-huanqiu backup]# pwd
/home/mysql/backup
[root@test-huanqiu backup]# mkdir daily
[root@test-huanqiu backup]# ll
total 4
drwxr-xr-x. 2 root root 4096 Nov 29 11:29 daily
[root@test-huanqiu backup]# ll daily/
total 0

先执行增量备份脚本
[root@test-huanqiu backup]# sh /root/Mysql-DailyBak.sh
[root@test-huanqiu backup]# ll
total 8
-rw-r--r--. 1 root root 121 Nov 29 11:29 bak.log
drwxr-xr-x. 2 root root 4096 Nov 29 11:29 daily
[root@test-huanqiu backup]# ll daily/
total 8
-rw-r-----. 1 root root 152 Nov 29 11:29 mysql-binlog.000030
-rw-r-----. 1 root root 152 Nov 29 11:29 mysql-binlog.000031
[root@test-huanqiu backup]# cat bak.log
mysql-binlog.000030 copying
mysql-binlog.000031 copying
mysql-binlog.000032 skip!
2016年11月29日 11:29:32 Bakup succ!

然后执行全量备份脚本
[root@test-huanqiu backup]# sh /root/Mysql-FullyBak.sh
20161129.sql
[root@test-huanqiu backup]# ll
total 152
-rw-r--r--. 1 root root 145742 Nov 29 11:30 20161129.sql.tgz
-rw-r--r--. 1 root root 211 Nov 29 11:30 bak.log
drwxr-xr-x. 2 root root 4096 Nov 29 11:30 daily
[root@test-huanqiu backup]# ll daily/
total 0
[root@test-huanqiu backup]# cat bak.log
mysql-binlog.000030 copying
mysql-binlog.000031 copying
mysql-binlog.000032 skip!
2016年11月29日 11:29:32 Bakup succ!
开始:2016年11月29日 11:30:38 结束:2016年11月29日 11:30:38 20161129.sql.tgz succ

在日常的linux运维工作中，大数据量备份与还原，始终是个难点。关于mysql的备份和恢复，比较传统的是用mysqldump工具，今天这里推荐另一个备份工具innobackupex。innobackupex和mysqldump都可以对mysql进行热备份的，mysqldump对mysql的innodb的备份可以使用single-transaction参数来开启一个事务，利用innodb的mvcc来不进行锁表进行热备份，mysqldump备份是逻辑备份，备份出来的文件是sql语句，所以备份和恢复的时候很慢，但是备份和恢复时候很清楚。当MYSQL数据超过10G时，用mysqldump来导出备份就比较慢了，此种情况下用innobackupex这个工具就比mysqldump要快很多。利用它对mysql做全量和增量备份，仅仅依据本人实战操作做一记录，如有误述，敬请指出~

一、innobackupex的介绍

Xtrabackup是由percona开发的一个开源软件，是使用perl语言完成的脚本工具，能够非常快速地备份与恢复mysql数据库，且支持在线热备份（备份时不影响数据读写），此工具调用xtrabackup和tar4ibd工具，实现很多对性能要求并不高的任务和备份逻辑，可以说它是innodb热备工具ibbackup的一个开源替代品。
Xtrabackup中包含两个工具：
1）xtrabackup ：只能用于热备份innodb,xtradb两种数据引擎表的工具，不能备份其他表。
2）innobackupex：是一个对xtrabackup封装的perl脚本，提供了用于myisam(会锁表)和innodb引擎，及混合使用引擎备份的能力。主要是为了方便同时备份InnoDB和MyISAM引擎的表，但在处理myisam时需要加一个读锁。并且加入了一些使用的选项。如slave-info可以记录备份恢 复后，作为slave需要的一些信息，根据这些信息，可以很方便的利用备份来重做slave。
innobackupex比xtarbackup有更强的功能，它整合了xtrabackup和其他的一些功能，它不但可以全量备份/恢复，还可以基于时间的增量备份与恢复。innobackupex同时支持innodb,myisam。

Xtrabackup可以做什么
1）在线(热)备份整个库的InnoDB, XtraDB表
2）在xtrabackup的上一次整库备份基础上做增量备份（innodb only）
3）以流的形式产生备份，可以直接保存到远程机器上（本机硬盘空间不足时很有用）

MySQL数据库本身提供的工具并不支持真正的增量备份，二进制日志恢复是point-in-time(时间点)的恢复而不是增量备份。Xtrabackup工具支持对InnoDB存储引擎的增量备份，工作原理如下：
1）首先完成一个完全备份，并记录下此时检查点的LSN(Log Sequence Number)。
2）在进程增量备份时，比较表空间中每个页的LSN是否大于上次备份时的LSN，如果是，则备份该页，同时记录当前检查点的LSN。首先，在logfile中找到并记录最后一个checkpoint(“last checkpoint LSN”)，然后开始从LSN的位置开始拷贝InnoDB的logfile到xtrabackup_logfile；接着，开始拷贝全部的数据文件.ibd；在拷贝全部数据文件结束之后，才停止拷贝logfile。因为logfile里面记录全部的数据修改情况，所以，即时在备份过程中数据文件被修改过了，恢复时仍然能够通过解析xtrabackup_logfile保持数据的一致。

innobackupex备份mysql数据的流程
innobackupex首先调用xtrabackup来备份innodb数据文件，当xtrabackup完成后，innobackupex就查看文件xtrabackup_suspended ；然后执行“FLUSH TABLES WITH READ LOCK”来备份其他的文件。

innobackupex恢复mysql数据的流程
innobackupex首先读取my.cnf，查看变量(datadir,innodb_data_home_dir,innodb_data_file_path,innodb_log_group_home_dir)对应的目录是存在，确定相关目录存在后，然后先copy myisam表和索引，然后在copy innodb的表、索引和日志。

------------------------------------------------------------------------------------------------------------------------------------------
下面详细说下innobackupex备份和恢复的工作原理：

（1）备份的工作原理
       如果在程序启动阶段未指定模式，innobackupex将会默认以备份模式启动。
       默认情况下，此脚本以--suspend-at-end选项启动xtrabackup，然后xtrabackup程序开始拷贝InnoDB数据文件。当xtrabackup程序执行结束，innobackupex将会发现xtrabackup创建了xtrabackup_suspended_2文件，然后执行FLUSH TABLES WITH READ LOCK，此语句对所有的数据库表加读锁。然后开始拷贝其他类型的文件。
       如果--ibbackup未指定，innobackupex将会自行尝试确定使用的xtrabackup的binary。其确定binary的逻辑如下：首先判断备份目录中xtrabackup_binary文件是否存在，如果存在，此脚本将会依据此文件确定使用的xtrabackup binary。否则，脚本将会尝试连接database server，通过server版本确定binary。如果连接无法建立，xtrabackup将会失败，需要自行指定binary文件。
       在binary被确定后，将会检查到数据库server的连接是否可以建立。其执行逻辑是：建立连接、执行query、关闭连接。若一切正常，xtrabackup将以子进程的方式启动。
       FLUSH TABLES WITH READ LOCK是为了备份MyISAM和其他非InnoDB类型的表，此语句在xtrabackup已经备份InnoDB数据和日志文件后执行。在这之后，将会备份 .frm, .MRG, .MYD, .MYI, .TRG, .TRN, .ARM, .ARZ, .CSM, .CSV, .par, and .opt 类型的文件。
       当所有上述文件备份完成后，innobackupex脚本将会恢复xtrabackup的执行，等待其备份上述逻辑执行过程中生成的事务日志文件。接下来，表被解锁，slave被启动，到server的连接被关闭。接下来，脚本会删掉xtrabackup_suspended_2文件，允许xtrabackup进程退出。

（2）恢复的工作原理
       为了恢复一个备份，innobackupex需要以--copy-back选项启动。
       innobackupex将会首先通过my.cnf文件读取如下变量：datadir, innodb_data_home_dir, innodb_data_file_path, innodb_log_group_home_dir，并确定这些目录存在。
       接下来，此脚本将会首先拷贝MyISAM表、索引文件、其他类型的文件（如：.frm, .MRG, .MYD, .MYI, .TRG, .TRN, .ARM, .ARZ, .CSM, .CSV, par and .opt files），接下来拷贝InnoDB表数据文件，最后拷贝日志文件。拷贝执行时将会保留文件属性，在使用备份文件启动MySQL前，可能需要更改文件的owener（如从拷贝文件的user更改到mysql用户）。
---------------------------------------------------------------------------------------------------------------------------------------------

二、innobackupex针对mysql数据库的备份环境部署
1）源码安装Xtrabackup，将源码包下载到/usr/local/src下
源码包下载
[root@test-huanqiu ~]# cd /usr/local/src
先安装依赖包
[root@test-huanqiu src]# yum -y install cmake gcc gcc-c++ libaio libaio-devel automake autoconf bzr bison libtool  zlib-devel libgcrypt-devel  libcurl-devel  crypt*  libgcrypt* python-sphinx openssl   imake libxml2-devel expat-devel   ncurses5-devel ncurses-devle   vim-common  libgpg-error-devel   libidn-devel perl-DBI  perl-DBD-MySQL perl-Time-HiRes perl-IO-Socket-SSL 
[root@test-huanqiu src]# wget http://www.percona.com/downloads/XtraBackup/XtraBackup-2.1.9/source/percona-xtrabackup-2.1.9.tar.gz
[root@test-huanqiu src]# tar -zvxf percona-xtrabackup-2.1.9.tar.gz
[root@test-huanqiu src]# cd percona-xtrabackup-2.1.9
[root@test-huanqiu percona-xtrabackup-2.1.9]# ./utils/build.sh                     //执行该安装脚本，会出现下面信息
Build an xtrabackup binary against the specified InnoDB flavor.

 

Usage: build.sh CODEBASE
where CODEBASE can be one of the following values or aliases:
innodb51 | plugin build against InnoDB plugin in MySQL 5.1
innodb55 | 5.5 build against InnoDB in MySQL 5.5
innodb56 | 5.6,xtradb56, build against InnoDB in MySQL 5.6
| mariadb100,galera56
xtradb51 | xtradb,mariadb51 build against Percona Server with XtraDB 5.1
| mariadb52,mariadb53
xtradb55 | galera55,mariadb55 build against Percona Server with XtraDB 5.5
根据上面提示和你使用的存储引擎及版本，选择相应的参数即可。因为我用的是MySQL 5.6版本，所以执行如下语句安装：
[root@test-huanqiu percona-xtrabackup-2.1.9]# ./utils/build.sh innodb56
以上语句执行成功后，表示安装完成。
最后，把生成的二进制文件拷贝到一个自定义目录下（本例中为/home/mysql/admin/bin/percona-xtrabackup-2.1.9），并把该目录放到环境变量PATH中。
[root@test-huanqiu  percona-xtrabackup-2.1.9]# mkdir -p /home/mysql/admin/bin/percona-xtrabackup-2.1.9/
[root@test-huanqiu  percona-xtrabackup-2.1.9]# cp ./innobackupex /home/mysql/admin/bin/percona-xtrabackup-2.1.9/
[root@test-huanqiu  percona-xtrabackup-2.1.9]# cp ./src/xtrabackup_56 ./src/xbstream  /home/mysql/admin/bin/percona-xtrabackup-2.1.9/
[root@test-huanqiu  percona-xtrabackup-2.1.9]# vim /etc/profile
.......
export PATH=$PATH:/home/mysql/admin/bin/percona-xtrabackup-2.1.9/
[root@test-huanqiu  percona-xtrabackup-2.1.9]# source /etc/profile

测试下innobackupex是否正常使用
[root@test-huanqiu percona-xtrabackup-2.1.9]# innobackupex --help
--------------------------------------------------------------------------------------------------------------------------------------------
可能报错1

Can't locate Time/HiRes.pm in @INC (@INC contains: /usr/local/lib64/perl5 /usr/local/share/perl5 /usr/lib64/perl5/vendor_perl /usr/share/perl5/vendor_perl /usr/lib64/perl5 /usr/share/perl5 .) at /home/mysql/admin/bin/percona-xtrabackup-2.1.9/innobackupex line 23.
BEGIN failed--compilation aborted at /home/mysql/admin/bin/percona-xtrabackup-2.1.9/innobackupex line 23.

解决方案：
.pm实际上是Perl的包，只需安装perl-Time-HiRes即可：

[root@test-huanqiu percona-xtrabackup-2.1.9]# yum install -y perl-Time-HiRes

可能报错2
Can't locate DBI.pm in @INC (@INC contains: /usr/lib64/perl5/site_perl/5.8.8/x86_64-linux-thread-multi /usr/lib/perl5/site_perl/5.8.8 /usr/lib/perl5/site_perl /usr/lib64/perl5/vendor_perl/5.8.8/x86_64-linux-thread-multi /usr/lib/perl5/vendor_perl/5.8.8 /usr/lib/perl5/vendor_perl /usr/lib64/perl5/5.8.8/x86_64-linux-thread-multi /usr/lib/perl5/5.8.8 .) at /usr/local/webserver/mysql5.1.57/bin/mysqlhotcopy line 25.
BEGIN failed--compilation aborted at /usr/local/webserver/mysql5.1.57/bin/mysqlhotcopy line 25.
报错原因：系统没有按安装DBI组件。
DBI(Database Interface)是perl连接数据库的接口。其是perl连接数据库的最优秀方法，他支持包括Orcal,Sybase,mysql,db2等绝大多数的数据库。

解决办法：
安装DBI组件（Can't locate DBI.pm in @INC-mysql接口）
或者单独装DBI、Data-ShowTable、DBD-mysql 三个组件
[root@test-huanqiu percona-xtrabackup-2.1.9]# yum -y install perl-DBD-MySQL

接着使用innobackupex命令测试是否正常
[root@test-huanqiu percona-xtrabackup-2.1.9]# innobackupex --help
Options:
--apply-log
Prepare a backup in BACKUP-DIR by applying the transaction log file
named "xtrabackup_logfile" located in the same directory. Also,
create new transaction logs. The InnoDB configuration is read from
the file "backup-my.cnf".

--compact
Create a compact backup with all secondary index pages omitted. This
option is passed directly to xtrabackup. See xtrabackup
documentation for details.

--compress
This option instructs xtrabackup to compress backup copies of InnoDB
data files. It is passed directly to the xtrabackup child process.
Try 'xtrabackup --help' for more details.
............
-------------------------------------------------------------------------------------------------------------------------------------------------------------

2）全量备份和恢复
---------------->全量备份操作<----------------
执行下面语句进行全备：
mysql的安装目录是/usr/local/mysql
mysql的配置文件路径/usr/local/mysql/my.cnf
mysql的密码是123456
全量备份后的数据存放目录是/backup/mysql/data
[root@test-huanqiu ~]# mkdir -p /backup/mysql/data
[root@test-huanqiu ~]# innobackupex --defaults-file=/usr/local/mysql/my.cnf --user=root --password=123456 /backup/mysql/data

InnoDB Backup Utility v1.5.1-xtrabackup; Copyright 2003, 2009 Innobase Oy
and Percona LLC and/or its affiliates 2009-2013. All Rights Reserved.
...................
161201 00:07:15 innobackupex: Connecting to MySQL server with DSN 'dbi:mysql:;mysql_read_default_file=/usr/local/mysql/my.cnf;mysql_read_default_group=xtrabackup' as 'root' (using password: YES).
161201 00:07:15 innobackupex: Connected to MySQL server
161201 00:07:15 innobackupex: Executing a version check against the server...
161201 00:07:15 innobackupex: Done.
..................
161201 00:07:19 innobackupex: Connection to database server closed
161201 00:07:19 innobackupex: completed OK!

出现上面的信息，表示备份已经ok。

上面执行的备份语句会将mysql数据文件（即由my.cnf里的变量datadir指定）拷贝至备份目录下（/backup/mysql/data）
注意：如果不指定--defaults-file，默认值为/etc/my.cnf。
备份成功后，将在备份目录下创建一个时间戳目录（本例创建的目录为/backup/mysql/data/2016-12-01_00-07-15），在该目录下存放备份文件。

查看备份数据：
[root@test-huanqiu ~]# ll /backup/mysql/data
total 4
drwxr-xr-x. 6 root root 4096 Dec 1 00:07 2016-12-01_00-07-15
[root@test-huanqiu ~]# ll /backup/mysql/data/2016-12-01_00-07-15/
total 12324
-rw-r--r--. 1 root root 357 Dec 1 00:07 backup-my.cnf
drwx------. 2 root root 4096 Dec 1 00:07 huanqiu
-rw-r-----. 1 root root 12582912 Dec 1 00:07 ibdata1
drwx------. 2 root root 4096 Dec 1 00:07 mysql
drwxr-xr-x. 2 root root 4096 Dec 1 00:07 performance_schema
drwxr-xr-x. 2 root root 4096 Dec 1 00:07 test
-rw-r--r--. 1 root root 13 Dec 1 00:07 xtrabackup_binary
-rw-r--r--. 1 root root 24 Dec 1 00:07 xtrabackup_binlog_info
-rw-r-----. 1 root root 89 Dec 1 00:07 xtrabackup_checkpoints
-rw-r-----. 1 root root 2560 Dec 1 00:07 xtrabackup_logfile

----------------------------------------------------------------------------------------------------------------------------------
可能报错1：
161130 05:56:48 innobackupex: Connecting to MySQL server with DSN 'dbi:mysql:;mysql_read_default_file=/usr/local/mysql/my.cnf;mysql_read_default_group=xtrabackup' as 'root' (using password: YES).
innobackupex: Error: Failed to connect to MySQL server as DBD::mysql module is not installed at /home/mysql/admin/bin/percona-xtrabackup-2.1.9/innobackupex line 2956.

解决办法：
[root@test-huanqiu ~]# yum -y install perl-DBD-MySQL.x86_64
......
Package perl-DBD-MySQL-4.013-3.el6.x86_64 already installed and latest version                //发现本机已经安装了

[root@test-huanqiu ~]# rpm -qa|grep perl-DBD-MySQL
perl-DBD-MySQL-4.013-3.el6.x86_64

发现本机已经安装了最新版的perl-DBD-MYSQL了，但是仍然报出上面的错误！！
莫慌~~继续下面的操作进行问题的解决

查看mysql.so依赖的lib库
[root@test-huanqiu ~]# ldd /usr/lib64/perl5/auto/DBD/mysql/mysql.so
linux-vdso.so.1 => (0x00007ffd291fc000)
libmysqlclient.so.16 => not found                                                   //这一项为通过检查，缺失libmysqlclient.so.16库导致
libz.so.1 => /lib64/libz.so.1 (0x00007f78ff9de000)
libcrypt.so.1 => /lib64/libcrypt.so.1 (0x00007f78ff7a7000)
libnsl.so.1 => /lib64/libnsl.so.1 (0x00007f78ff58e000)
libm.so.6 => /lib64/libm.so.6 (0x00007f78ff309000)
libssl.so.10 => /usr/lib64/libssl.so.10 (0x00007f78ff09d000)
libcrypto.so.10 => /usr/lib64/libcrypto.so.10 (0x00007f78fecb9000)
libc.so.6 => /lib64/libc.so.6 (0x00007f78fe924000)
libfreebl3.so => /lib64/libfreebl3.so (0x00007f78fe721000)
libgssapi_krb5.so.2 => /lib64/libgssapi_krb5.so.2 (0x00007f78fe4dd000)
libkrb5.so.3 => /lib64/libkrb5.so.3 (0x00007f78fe1f5000)
libcom_err.so.2 => /lib64/libcom_err.so.2 (0x00007f78fdff1000)
libk5crypto.so.3 => /lib64/libk5crypto.so.3 (0x00007f78fddc5000)
libdl.so.2 => /lib64/libdl.so.2 (0x00007f78fdbc0000)
/lib64/ld-linux-x86-64.so.2 (0x00007f78ffe1d000)
libkrb5support.so.0 => /lib64/libkrb5support.so.0 (0x00007f78fd9b5000)
libkeyutils.so.1 => /lib64/libkeyutils.so.1 (0x00007f78fd7b2000)
libresolv.so.2 => /lib64/libresolv.so.2 (0x00007f78fd597000)
libpthread.so.0 => /lib64/libpthread.so.0 (0x00007f78fd37a000)
libselinux.so.1 => /lib64/libselinux.so.1 (0x00007f78fd15a000)

以上结果说明缺少libmysqlclient.so.16这个二进制包，找个官方原版的mysql的libmysqlclient.so.16替换了即可！
[root@test-huanqiu~]# find / -name libmysqlclient.so.16                                   //查看本机并没有libmysqlclient.so.16库文件

查看mysql/lib下的libmysqlclinet.so库文件
[root@test-huanqiu~]# ll /usr/local/mysql/lib/
total 234596
-rw-r--r--. 1 mysql mysql 19520800 Nov 29 12:27 libmysqlclient.a
lrwxrwxrwx. 1 mysql mysql 16 Nov 29 12:34 libmysqlclient_r.a -> libmysqlclient.a
lrwxrwxrwx. 1 mysql mysql 17 Nov 29 12:34 libmysqlclient_r.so -> libmysqlclient.so
lrwxrwxrwx. 1 mysql mysql 20 Nov 29 12:34 libmysqlclient_r.so.18 -> libmysqlclient.so.18
lrwxrwxrwx. 1 mysql mysql 24 Nov 29 12:34 libmysqlclient_r.so.18.1.0 -> libmysqlclient.so.18.1.0
lrwxrwxrwx. 1 mysql mysql 20 Nov 29 12:34 libmysqlclient.so -> libmysqlclient.so.18
lrwxrwxrwx. 1 mysql mysql 24 Nov 29 12:34 libmysqlclient.so.18 -> libmysqlclient.so.18.1.0
-rwxr-xr-x. 1 mysql mysql 8858235 Nov 29 12:27 libmysqlclient.so.18.1.0
-rw-r--r--. 1 mysql mysql 211822074 Nov 29 12:34 libmysqld.a
-rw-r--r--. 1 mysql mysql 14270 Nov 29 12:27 libmysqlservices.a
drwxr-xr-x. 3 mysql mysql 4096 Nov 29 12:34 plugin

将mysql/lib/libmysqlclient.so.18.1.0库文件拷贝到/lib64下，拷贝后命名为libmysqlclient.so.16
[root@test-huanqiu~]# cp /usr/local/mysql/lib/libmysqlclient.so.18.1.0 /lib64/libmysqlclient.so.16

[root@test-huanqiu~]# cat /etc/ld.so.conf
include ld.so.conf.d/*.conf
/usr/local/mysql/lib/
/lib64/
[root@test-huanqiu~]# ldconfig

最后卸载perl-DBD-MySQL，并重新安装perl-DBD-MySQL
[root@test-huanqiu~]# rpm -qa|grep perl-DBD-MySQL
perl-DBD-MySQL-4.013-3.el6.x86_64
[root@test-huanqiu~]# rpm -e --nodeps perl-DBD-MySQL
[root@test-huanqiu~]# rpm -qa|grep perl-DBD-MySQL
[root@test-huanqiu~]# yum -y install perl-DBD-MySQL

待重新安装后，再次重新检查mysql.so依赖的lib库，发现已经都通过了
[root@test-huanqiu~]# ldd /usr/lib64/perl5/auto/DBD/mysql/mysql.so
linux-vdso.so.1 => (0x00007ffe3669b000)
libmysqlclient.so.16 => /usr/lib64/mysql/libmysqlclient.so.16 (0x00007f4af5c25000)
libz.so.1 => /lib64/libz.so.1 (0x00007f4af5a0f000)
libcrypt.so.1 => /lib64/libcrypt.so.1 (0x00007f4af57d7000)
libnsl.so.1 => /lib64/libnsl.so.1 (0x00007f4af55be000)
libm.so.6 => /lib64/libm.so.6 (0x00007f4af533a000)
libssl.so.10 => /usr/lib64/libssl.so.10 (0x00007f4af50cd000)
libcrypto.so.10 => /usr/lib64/libcrypto.so.10 (0x00007f4af4ce9000)
libc.so.6 => /lib64/libc.so.6 (0x00007f4af4955000)
libfreebl3.so => /lib64/libfreebl3.so (0x00007f4af4751000)
libgssapi_krb5.so.2 => /lib64/libgssapi_krb5.so.2 (0x00007f4af450d000)
libkrb5.so.3 => /lib64/libkrb5.so.3 (0x00007f4af4226000)
libcom_err.so.2 => /lib64/libcom_err.so.2 (0x00007f4af4021000)
libk5crypto.so.3 => /lib64/libk5crypto.so.3 (0x00007f4af3df5000)
libdl.so.2 => /lib64/libdl.so.2 (0x00007f4af3bf1000)
/lib64/ld-linux-x86-64.so.2 (0x00007f4af61d1000)
libkrb5support.so.0 => /lib64/libkrb5support.so.0 (0x00007f4af39e5000)
libkeyutils.so.1 => /lib64/libkeyutils.so.1 (0x00007f4af37e2000)
libresolv.so.2 => /lib64/libresolv.so.2 (0x00007f4af35c8000)
libpthread.so.0 => /lib64/libpthread.so.0 (0x00007f4af33aa000)
libselinux.so.1 => /lib64/libselinux.so.1 (0x00007f4af318b000)

可能报错2
sh: xtrabackup_56: command not found
innobackupex: Error: no 'mysqld' group in MySQL options at /home/mysql/admin/bin/percona-xtrabackup-2.1.9/innobackupex line 4350.

有可能是percona-xtrabackup编译安装后，在编译目录的src下存在xtrabackup_innodb56，只需要其更名为xtrabackup_56，然后拷贝到上面的/home/mysql/admin/bin/percona-xtrabackup-2.1.9/下即可！
----------------------------------------------------------------------------------------------------------------------------------

还可以在远程进行全量备份，命令如下：
[root@test-huanqiu ~]# innobackupex --defaults-file=/usr/local/mysql/my.cnf --user=root --password=123456 --host=127.0.0.1 --parallel=2 --throttle=200 /backup/mysql/data 2>/backup/mysql/data/bak.log 1>/backup/mysql/data/`data +%Y-%m-%d_%H-%M%S`
参数解释：
--user=root             备份操作用户名，一般都是root用户
--password=root123          数据库密码
--host=127.0.0.1            主机ip，本地可以不加（适用于远程备份）。注意要提前在mysql中授予连接的权限，最好备份前先测试用命令中的用户名、密码和host能否正常连接mysql。
--parallel=2 --throttle=200      并行个数，根据主机配置选择合适的，默认是1个，多个可以加快备份速度。
/backup/mysql/data            备份存放的目录
2>/backup/mysql/data/bak.log       备份日志，将备份过程中的输出信息重定向到bak.log

这种备份跟上面相比，备份成功后，不会自动在备份目录下创建一个时间戳目录，需要如上命令中自己定义。
[root@test-huanqiu ~]# cd /backup/mysql/data/
[root@test-huanqiu data]# ll
drwxr-xr-x. 6 root root 4096 Dec 1 03:18 2016-12-01_03-18-37
-rw-r--r--. 1 root root 5148 Dec 1 03:18 bak.log
[root@test-huanqiu data]# cat bak.log         //备份信息都记录在这个日志里，如果备份失败，可以到这里日志里查询

----------------------------------------------------------------------------------------------------------------------------

---------------->全量备份后的恢复操作<----------------

比如在上面进行全量备份后，由于误操作将数据库中的huanqiu库删除了。
[root@test-huanqiu ~]# mysql -p123456
.......
mysql> show databases;
+--------------------+
| Database |
+--------------------+
| information_schema |
| huanqiu |
| mysql |
| performance_schema |
| test |
+--------------------+
5 rows in set (0.00 sec)

mysql> use huanqiu;
Reading table information for completion of table and column names
You can turn off this feature to get a quicker startup with -A

Database changed
mysql> show tables;
+-------------------------+
| Tables_in_huanqiu |
+-------------------------+
| card_agent_file |
| product_sale_management |
+-------------------------+
2 rows in set (0.00 sec)

mysql> drop database huanqiu;
Query OK, 2 rows affected (0.12 sec)

mysql> show databases;
+--------------------+
| Database |
+--------------------+
| information_schema |
| mysql |
| performance_schema |
| test |
+--------------------+
4 rows in set (0.00 sec)

现在进行恢复数据操作
注意：恢复之前
1）要先关闭数据库
2）要删除数据文件和日志文件（也可以mv移到别的地方，只要确保清空mysql数据存放目录就行）
[root@test-huanqiu ~]# ps -ef|grep mysql
root 2442 21929 0 00:25 pts/2 00:00:00 grep mysql
root 28279 1 0 Nov29 ? 00:00:00 /bin/sh /usr/local/mysql//bin/mysqld_safe --datadir=/data/mysql/data --pid-file=/data/mysql/data/mysql.pid
mysql 29059 28279 0 Nov29 ? 00:09:07 /usr/local/mysql/bin/mysqld --basedir=/usr/local/mysql/ --datadir=/data/mysql/data --plugin-dir=/usr/local/mysql//lib/plugin --user=mysql --log-error=/data/mysql/data/mysql-error.log --pid-file=/data/mysql/data/mysql.pid --socket=/usr/local/mysql/var/mysql.sock --port=3306

由上面可查出mysql的数据和日志存放目录是/data/mysql/data
[root@test-huanqiu ~]# /etc/init.d/mysql stop
Shutting down MySQL.. SUCCESS!
[root@test-huanqiu ~]# rm -rf /data/mysql/data/*
[root@test-huanqiu ~]# ls /data/mysql/data
[root@test-huanqiu ~]#

查看备份数据
[root@[root@test-huanqiu ~]# ls /backup/mysql/data/
2016-12-01_00-07-15

恢复数据
[root@[root@test-huanqiu ~]# innobackupex --defaults-file=/usr/local/mysql/my.cnf --user=root --password=123456 --use-memory=4G --apply-log /backup/mysql/data/2016-12-01_00-07-15
[root@[root@test-huanqiu ~]# innobackupex --defaults-file=/usr/local/mysql/my.cnf --user=root --password=123456 --copy-back /backup/mysql/data/2016-12-01_00-07-15
........
innobackupex: Copying '/backup/mysql/data/2016-12-01_00-07-15/ib_logfile2' to '/data/mysql/data/ib_logfile2'
innobackupex: Copying '/backup/mysql/data/2016-12-01_00-07-15/ib_logfile0' to '/data/mysql/data/ib_logfile0'
innobackupex: Finished copying back files.

161201 00:31:33 innobackupex: completed OK!
出现上面的信息，说明数据恢复成功了！！

从上面的恢复操作可以看出，执行恢复分为两个步骤：
1）第一步恢复步骤是应用日志（apply-log），为了加快速度，一般建议设置--use-memory（如果系统内存充足，可以使用加大内存进行备份 ），这个步骤完成之后，目录/backup/mysql/data/2016-12-01_00-07-15下的备份文件已经准备就绪。
2）第二步恢复步骤是拷贝文件（copy-back），即把备份文件拷贝至原数据目录下。
恢复完成之后，一定要记得检查数据目录的所有者和权限是否正确。

[root@test-huanqiu ~]# ll /data/mysql/data/
total 110608
drwxr-xr-x. 2 root root 4096 Dec 1 00:31 huanqiu
-rw-r--r--. 1 root root 12582912 Dec 1 00:31 ibdata1
-rw-r--r--. 1 root root 33554432 Dec 1 00:31 ib_logfile0
-rw-r--r--. 1 root root 33554432 Dec 1 00:31 ib_logfile1
-rw-r--r--. 1 root root 33554432 Dec 1 00:31 ib_logfile2
drwxr-xr-x. 2 root root 4096 Dec 1 00:31 mysql
drwxr-xr-x. 2 root root 4096 Dec 1 00:31 performance_schema
drwxr-xr-x. 2 root root 4096 Dec 1 00:31 test
[root@test-huanqiu ~]# chown -R mysql.mysql /data/mysql/data/                                   //将数据目录的权限修改为mysql:mysql
[root@test-huanqiu ~]# ll /data/mysql/data/
total 110608
drwxr-xr-x. 2 mysql mysql 4096 Dec 1 00:31 huanqiu
-rw-r--r--. 1 mysql mysql 12582912 Dec 1 00:31 ibdata1
-rw-r--r--. 1 mysql mysql 33554432 Dec 1 00:31 ib_logfile0
-rw-r--r--. 1 mysql mysql 33554432 Dec 1 00:31 ib_logfile1
-rw-r--r--. 1 mysql mysql 33554432 Dec 1 00:31 ib_logfile2
drwxr-xr-x. 2 mysql mysql 4096 Dec 1 00:31 mysql
drwxr-xr-x. 2 mysql mysql 4096 Dec 1 00:31 performance_schema
drwxr-xr-x. 2 mysql mysql 4096 Dec 1 00:31 test
-------------------------------------------------------------------------------------------------------------------------------------------
可能报错：
sh: xtrabackup: command not found
innobackupex: Error: no 'mysqld' group in MySQL options at /home/mysql/admin/bin/percona-xtrabackup-2.1.9/innobackupex line 4350.

解决：将xtrabackup_56复制成xtrabackup即可
[root@test-huanqiu percona-xtrabackup-2.1.9]# ls
innobackupex xbstream xtrabackup_56
[root@test-huanqiu percona-xtrabackup-2.1.9]# cp xtrabackup_56 xtrabackup
[root@test-huanqiu percona-xtrabackup-2.1.9]# ls
innobackupex xbstream xtrabackup xtrabackup_56
-------------------------------------------------------------------------------------------------------------------------------------------

最后，启动mysql，查看数据是否恢复回来了
[root@test-huanqiu ~]# /etc/init.d/mysql start
Starting MySQL.. SUCCESS!
[root@test-huanqiu ~]# mysql -p123456
........
mysql> show databases;
+--------------------+
| Database |
+--------------------+
| information_schema |
| huanqiu |
| mysql |
| performance_schema |
| test |
+--------------------+
5 rows in set (0.00 sec)

mysql> use huanqiu;
Reading table information for completion of table and column names
You can turn off this feature to get a quicker startup with -A

Database changed
mysql> show tables;
+-------------------------+
| Tables_in_huanqiu |
+-------------------------+
| card_agent_file |
| product_sale_management |
+-------------------------+
2 rows in set (0.00 sec)

mysql>

3）增量备份和恢复
---------------->增量备份操作<----------------
特别注意：
innobackupex 增量备份仅针对InnoDB这类支持事务的引擎，对于MyISAM等引擎，则仍然是全备。

增量备份需要基于全量备份
先假设我们已经有了一个全量备份（如上面的/backup/mysql/data/2016-12-01_00-07-15），我们需要在该全量备份的基础上做第一次增量备份。
[root@test-huanqiu ~]# innobackupex --defaults-file=/usr/local/mysql/my.cnf --user=root --password=123456 --incremental-basedir=/backup/mysql/data/2016-12-01_00-07-15 --incremental /backup/mysql/data
其中：
--incremental-basedir     指向全量备份目录
--incremental       指向增量备份的目录
上面语句执行成功之后，会在--incremental执行的目录下创建一个时间戳子目录（本例中为：/backup/mysql/data/2016-12-01_01-12-22），在该目录下存放着增量备份的所有文件。
[root@test-huanqiu ~]# ll /backup/mysql/data/
total 8
drwxr-xr-x. 6 root root 4096 Dec 1 00:27 2016-12-01_00-07-15                //全量备份目录
drwxr-xr-x. 6 root root 4096 Dec 1 01:12 2016-12-01_01-12-22                //增量备份目录

在备份目录下，有一个文件xtrabackup_checkpoints记录着备份信息，其中可以查出
1）全量备份的信息如下：
[root@test-huanqiu 2016-12-01_00-07-15]# pwd
/backup/mysql/data/2016-12-01_00-07-15
[root@test-huanqiu 2016-12-01_00-07-15]# cat xtrabackup_checkpoints
backup_type = full-prepared
from_lsn = 0
to_lsn = 1631561
last_lsn = 1631561
compact = 0

2）基于以上全量备份的增量备份的信息如下：
[root@test-huanqiu 2016-12-01_01-12-22]# pwd
/backup/mysql/data/2016-12-01_01-12-22
[root@test-huanqiu 2016-12-01_01-12-22]# cat xtrabackup_checkpoints
backup_type = incremental
from_lsn = 1631561
to_lsn = 1631776
last_lsn = 1631776
compact = 0

从上面可以看出，增量备份的from_lsn正好等于全备的to_lsn。
那么，我们是否可以在增量备份的基础上再做增量备份呢？
答案是肯定的，只要把--incremental-basedir执行上一次增量备份的目录即可，如下所示：
[root@test-huanqiu ~]# innobackupex --defaults-file=/usr/local/mysql/my.cnf --user=root --password=123456 --incremental-basedir=/backup/mysql/data/2016-12-01_01-12-22 --incremental /backup/mysql/data
[root@test-huanqiu data]# ll
total 12
drwxr-xr-x. 6 root root 4096 Dec 1 00:27 2016-12-01_00-07-15            //全量备份目录
drwxr-xr-x. 6 root root 4096 Dec 1 01:12 2016-12-01_01-12-22            //增量备份目录1
drwxr-xr-x. 6 root root 4096 Dec 1 01:23 2016-12-01_01-23-23           //增量备份目录2

它的xtrabackup_checkpoints记录着备份信息如下：
[root@test-huanqiu 2016-12-01_01-23-23]# pwd
/backup/mysql/data/2016-12-01_01-23-23
[root@test-huanqiu 2016-12-01_01-23-23]# cat xtrabackup_checkpoints
backup_type = incremental
from_lsn = 1631776
to_lsn = 1638220
last_lsn = 1638220
compact = 0

可以看到，第二次增量备份的from_lsn是从上一次增量备份的to_lsn开始的。

---------------->增量备份后的恢复操作<----------------
增量备份的恢复要比全量备份复杂很多，增量备份与全量备份有着一些不同，尤其要注意的是：
1)需要在每个备份(包括完全和各个增量备份)上，将已经提交的事务进行“重放”。“重放”之后，所有的备份数据将合并到完全备份上。
2)基于所有的备份将未提交的事务进行“回滚”。于是，操作就变成了：不能回滚，因为有可能第一次备份时候没提交，在增量中已经成功提交

第一步是在所有备份目录下重做已提交的日志（注意备份目录路径要跟全路径）
1）innobackupex --apply-log --redo-only BASE-DIR
2）innobackupex --apply-log --redo-only BASE-DIR --incremental-dir=INCREMENTAL-DIR-1
3）innobackupex --apply-log BASE-DIR --incremental-dir=INCREMENTAL-DIR-2
其中：
BASE-DIR 是指全量备份的目录
INCREMENTAL-DIR-1 是指第一次增量备份的目录
INCREMENTAL-DIR-2 是指第二次增量备份的目录，以此类推。
这里要注意的是：
1）最后一步的增量备份并没有--redo-only选项！回滚进行崩溃恢复过程
2）可以使用--use_memory提高性能。
以上语句执行成功之后，最终数据在BASE-DIR（即全量目录）下，其实增量备份就是把增量目录下的数据，整合到全变量目录下，然后在进行，全数据量的还原。

第一步完成之后，我们开始下面关键的第二步，即拷贝文件，进行全部还原！注意：必须先停止mysql数据库，然后清空数据库目录(这里是指/data/mysql/data)下的文件。

4）innobackupex --copy-back BASE-DIR
同样地，拷贝结束之后，记得检查下数据目录(这里指/data/mysql/data)的权限是否正确(修改成mysql:mysql)，然后再重启mysql。

接下来进行案例说明：
假设我们已经有了一个全量备份2016-12-01_00-07-15
删除在上面测试创建的两个增量备份
[root@test-huanqiu ~]# cd /backup/mysql/data/
[root@test-huanqiu data]# ll
total 12
drwxr-xr-x. 6 root root 4096 Dec 1 00:27 2016-12-01_00-07-15
drwxr-xr-x. 6 root root 4096 Dec 1 01:12 2016-12-01_01-12-22
drwxr-xr-x. 6 root root 4096 Dec 1 01:23 2016-12-01_01-23-23
[root@test-huanqiu data]# rm -rf 2016-12-01_01-12-22/
[root@test-huanqiu data]# rm -rf 2016-12-01_01-23-23/
[root@test-huanqiu data]# ll
total 4
drwxr-xr-x. 6 root root 4096 Dec 1 00:27 2016-12-01_00-07-15

假设在全量备份后，mysql数据库中又有新数据写入
[root@test-huanqiu ~]# mysql -p123456
.........
mysql> create database ceshi;
Query OK, 1 row affected (0.00 sec)

mysql> use ceshi;
Database changed
mysql> create table test1(
-> id int3,
-> name varchar(20)
-> );
Query OK, 0 rows affected (0.07 sec)

mysql> insert into test1 values(1,"wangshibo");
Query OK, 1 row affected, 1 warning (0.03 sec)

mysql> select * from test1;
+------+-----------+
| id | name |
+------+-----------+
| 1 | wangshibo |
+------+-----------+
1 row in set (0.00 sec)

mysql> show databases;
+--------------------+
| Database |
+--------------------+
| information_schema |
| ceshi |
| huanqiu |
| mysql |
| performance_schema |
| test |
+--------------------+
6 rows in set (0.00 sec)

mysql>

然后进行一次增量备份：
[root@test-huanqiu ~]# innobackupex --defaults-file=/usr/local/mysql/my.cnf --user=root --password=123456 --incremental-basedir=/backup/mysql/data/2016-12-01_00-07-15 --incremental /backup/mysql/data
[root@test-huanqiu ~]# ll /backup/mysql/data/
total 8
drwxr-xr-x. 6 root root 4096 Dec 1 00:27 2016-12-01_00-07-15        //全量备份目录
drwxr-xr-x. 7 root root 4096 Dec 1 03:41 2016-12-01_03-41-41        //增量备份目录

接着再在mysql数据库中写入新数据
mysql> insert into test1 values(2,"guohuihui");
Query OK, 1 row affected, 1 warning (0.00 sec)

mysql> insert into test1 values(3,"wuxiang");
Query OK, 1 row affected, 1 warning (0.00 sec)

mysql> insert into test1 values(4,"liumengnan");
Query OK, 1 row affected, 1 warning (0.01 sec)

mysql> select * from test1;
+------+------------+
| id | name |
+------+------------+
| 1 | wangshibo |
| 2 | guohuihui |
| 3 | wuxiang |
| 4 | liumengnan |
+------+------------+
4 rows in set (0.00 sec)

接着在增量的基础上再进行一次增量备份
[root@test-huanqiu ~]# innobackupex --defaults-file=/usr/local/mysql/my.cnf --user=root --password=123456 --incremental-basedir=/backup/mysql/data/2016-12-01_03-41-41 --incremental /backup/mysql/data
[root@test-huanqiu ~]# ll /backup/mysql/data/
total 12
drwxr-xr-x. 6 root root 4096 Dec 1 00:27 2016-12-01_00-07-15       //全量备份目录
drwxr-xr-x. 7 root root 4096 Dec 1 02:24 2016-12-01_02-24-11       //增量备份目录1
drwxr-xr-x. 7 root root 4096 Dec 1 03:42 2016-12-01_03-42-43       //增量备份目录2

现在删除数据库huanqiu、ceshi
mysql> show databases;
+--------------------+
| Database |
+--------------------+
| information_schema |
| ceshi |
| huanqiu |
| mysql |
| performance_schema |
| test |
+--------------------+
6 rows in set (0.00 sec)

mysql> drop database huanqiu;
Query OK, 2 rows affected (0.02 sec)

mysql> drop database ceshi;
Query OK, 1 row affected (0.01 sec)

mysql> show databases;
+--------------------+
| Database |
+--------------------+
| information_schema |
| mysql |
| performance_schema |
| test |
+--------------------+
4 rows in set (0.00 sec)

mysql>

接下来就开始进行数据恢复操作：

先恢复应用日志（注意最后一个不需要加--redo-only参数）
[root@test-huanqiu ~]# innobackupex --defaults-file=/usr/local/mysql/my.cnf --user=root --password=123456  --apply-log --redo-only /backup/mysql/data/2016-12-01_00-07-15
[root@test-huanqiu ~]# innobackupex --defaults-file=/usr/local/mysql/my.cnf --user=root --password=123456  --apply-log --redo-only /backup/mysql/data/2016-12-01_00-07-15 --incremental-dir=/backup/mysql/data/2016-12-01_02-24-11
[root@test-huanqiu ~]# innobackupex --defaults-file=/usr/local/mysql/my.cnf --user=root --password=123456  --apply-log  /backup/mysql/data/2016-12-01_00-07-15 --incremental-dir=/backup/mysql/data/2016-12-01_03-42-43

到此，恢复数据工作还没有结束！还有最重要的一个环节，就是把增量目录下的数据整合到全量备份目录下，然后再进行一次全量还原。
停止mysql数据库，并清空数据目录
[root@test-huanqiu ~]# /etc/init.d/mysql stop
[root@test-huanqiu ~]# rm -rf /data/mysql/data/*

最后拷贝文件，并验证数据目录的权限
[root@test-huanqiu ~]# innobackupex --defaults-file=/usr/local/mysql/my.cnf --user=root --password=123456 --copy-back /backup/mysql/data/2016-12-01_00-07-15

[root@test-huanqiu ~]# chown -R mysql.mysql /data/mysql/data/*
[root@test-huanqiu ~]# /etc/init.d/mysql start

最后，检查下数据是否恢复
[root@test-huanqiu ~]# mysql -p123456
........
mysql> show databases;
+--------------------+
| Database |
+--------------------+
| information_schema |
| ceshi |
| huanqiu |
| mysql |
| performance_schema |
| test |
+--------------------+
6 rows in set (0.00 sec)
mysql> select * from ceshi.test1;
+------+------------+
| id | name |
+------+------------+
| 1 | wangshibo |
| 2 | guohuihui |
| 3 | wuxiang |
| 4 | liumengnan |
+------+------------+
4 rows in set (0.00 sec)

 

另外注意：
上面在做备份的时候，将备份目录和增量目录都放在了同一个目录路径下，其实推荐放在不同的路径下，方便管理！比如：
/backup/mysql/data/full 存放全量备份目录
/backup/mysql/data/daily1 存放第一次增量备份目录
/backup/mysql/data/daily2 存放第二次增量目录
以此类推

在恢复的时候，注意命令中的路径要跟对！

-----------------------------------------------------------------------------------------------------
innobackupex 常用参数说明
--defaults-file
同xtrabackup的--defaults-file参数

--apply-log
对xtrabackup的--prepare参数的封装

--copy-back
做数据恢复时将备份数据文件拷贝到MySQL服务器的datadir ；

--remote-host=HOSTNAME
通过ssh将备份数据存储到进程服务器上；

--stream=[tar]
备 份文件输出格式, tar时使用tar4ibd , 该文件可在XtarBackup binary文件中获得.如果备份时有指定--stream=tar, 则tar4ibd文件所处目录一定要在$PATH中(因为使用的是tar4ibd去压缩, 在XtraBackup的binary包中可获得该文件)。
在 使用参数stream=tar备份的时候，你的xtrabackup_logfile可能会临时放在/tmp目录下，如果你备份的时候并发写入较大的话 xtrabackup_logfile可能会很大(5G+)，很可能会撑满你的/tmp目录，可以通过参数--tmpdir指定目录来解决这个问题。

--tmpdir=DIRECTORY
当有指定--remote-host or --stream时, 事务日志临时存储的目录, 默认采用MySQL配置文件中所指定的临时目录tmpdir

--redo-only --apply-log组,
强制备份日志时只redo ,跳过rollback。这在做增量备份时非常必要。

--use-memory=#
该参数在prepare的时候使用，控制prepare时innodb实例使用的内存量

--throttle=IOS
同xtrabackup的--throttle参数

--sleep=是给ibbackup使用的，指定每备份1M数据，过程停止拷贝多少毫秒，也是为了在备份时尽量减小对正常业务的影响，具体可以查看ibbackup的手册 ；

--compress[=LEVEL]
对备份数据迚行压缩，仅支持ibbackup，xtrabackup还没有实现；

--include=REGEXP
对 xtrabackup参数--tables的封装，也支持ibbackup。备份包含的库表，例如：--include="test.*"，意思是要备份 test库中所有的表。如果需要全备份，则省略这个参数；如果需要备份test库下的2个表：test1和test2,则写 成：--include="test.test1|test.test2"。也可以使用通配符，如：--include="test.test*"。

--databases=LIST
列出需要备份的databases，如果没有指定该参数，所有包含MyISAM和InnoDB表的database都会被备份；

--uncompress
解压备份的数据文件，支持ibbackup，xtrabackup还没有实现该功能；

--slave-info,
备 份从库, 加上--slave-info备份目录下会多生成一个xtrabackup_slave_info 文件, 这里会保存主日志文件以及偏移, 文件内容类似于:CHANGE MASTER TO MASTER_LOG_FILE='', MASTER_LOG_POS=0

--socket=SOCKET
指定mysql.sock所在位置，以便备份进程登录mysql.

三、innobackupex全量、增量备份脚本

可以根据自己线上数据库情况，编写全量和增量备份脚本，然后结合crontab设置计划执行。
比如：每周日的1:00进行全量备份，每周1-6的1:00进行增量备份。
还可以在脚本里编写邮件通知信息（可以用mail或sendemail）

==================================================================

1
2
3
4
5
6
7
8
9
10
11
12
13
14
15
16
17
18
19
在使用xtrabackup对mysql执行备份操作的时候，出现下面的报错：
.....................
xtrabackup: innodb_log_file_size = 50331648
InnoDB: Error: log file ./ib_logfile0 is of different size 33554432 bytes
InnoDB: than specified in the .cnf file 50331648 bytes!
 
 
解决办法：
可以计算一下33554432的大小，33554432/1024/1024=32
查看my.cnf配置文件的innodb_log_file_size参数配置：
innodb_log_file_size = 32M
 
需要调整这个文件的大小
再计算一下50331648的大小，50331648/1024/1024=48
 
修改my.cnf配置文件的下面一行参数值：
innodb_log_file_size = 48M
 
然后重启mysql
============下面是曾经使用过的一个mysql通过innobackupex进行增量备份脚本===========

1
2
3
4
5
6
7
8
9
10
11
12
13
14
15
16
17
18
19
20
21
22
23
24
25
26
27
28
29
30
31
32
33
34
35
36
37
38
39
40
41
42
43
44
45
46
47
48
49
50
51
52
53
54
55
56
57
58
59
60
61
62
63
64
65
66
67
68
69
70
71
72
73
74
75
76
77
78
79
80
81
82
83
84
85
86
87
[root@mysql-node ~]# cat /data/backup/script/incremental-backup-mysql.sh
#!/bin/sh
#########################################################################
## Description: Mysql增量备份脚本
## File Name: incremental-backup-mysql.sh
## Author: wangshibo
## mail: wangshibo@kevin.com
## Created Time: 2018年1月11日 14:17:09
##########################################################################
today=`date +%Y%m%d`
datetime=`date +%Y%m%d-%H-%M-%S`
config=/etc/my.cnf
basePath=/data/backup
logfilePath=$basePath/logs
logfile=$logfilePath/incr_$datetime.log
USER=mybak
PASSWD=Mysql@!@#1988check
dataBases="activiti batchdb core scf_v2 midfax asset bc_asset"
 
pid=`ps -ef | grep -v "grep" |grep -i innobackupex|awk '{print $2}'|head -n 1`
if [ -z $pid ]
then
  echo " start incremental backup database " >> $logfile
  OneMonthAgo=`date -d "1 month ago"  +%Y%m%d`
  path=$basePath/incr_$datetime
  mkdir -p $path
  last_backup=`cat $logfilePath/last_backup_sucess.log| head -1`
  echo " last backup is ===> " $last_backup >> $logfile
sudo /data/backup/script/percona-xtrabackup-2.4.2-Linux-x86_64/bin/innobackupex  --defaults-file=$config  --user=$USER --password=$PASSWD --compress --compress-threads=2 --compress-chunk-size=64K --slave-info --safe-slave-backup  --host=localhost --incremental $path --incremental-basedir=$last_backup --databases="${dataBases}" --no-timestamp >> $logfile 2>&1
sudo chown app.app $path -R
  ret=`tail -n 2 $logfile |grep "completed OK"|wc -l`
  if [ "$ret" =  1 ] ; then
    echo 'delete expired backup ' $basePath/incr_$OneMonthAgo*  >> $logfile
    rm -rf $basePath/incr_$OneMonthAgo*
    rm -f $logfilePath/incr_$OneMonthAgo*.log
    echo $path > $logfilePath/last_backup_sucess.log
  else
    echo 'backup failure ,no delete expired backup'  >> $logfile
  fi
else
   echo "****** innobackupex in backup database  ****** "  >> $logfile
fi
 
 
增量备份文件放在了本机的/data/backup目录下，再编写一个rsync脚本同步到远程备份机上（192.168.10.130）：
[root@mysql-node ~]# cat /data/rsync.sh
#!/bin/bash
datetime=`date +%Y%m%d-%H-%M-%S`
logfile=/data/rsync.log
echo "$datetime Rsync backup mysql start "  >> $logfile
sudo rsync -e "ssh -p6666" -avpgolr /data/backup bigtree@192.168.10.130:/data/backup_data/bigtree/DB_bak/10.0.40.52/ >> $logfile 2>&1
 
ret=`tail -n 1 $logfile |grep "total size"|wc -l`
if [ "$ret" =  1 ] ; then
        echo "$datetime Rsync backup mysql finish " >> $logfile
else
        echo "$datetime Rsync backup failure ,pls sendmail"  >> $logfile
fi
 
 
结合crontab进行定时任务执行（每4个小时执行一次）
[root@mysql-node ~]# crontab -e
1 4,8,12,16,20,23 * * * /data/backup/script/incremental-backup-mysql.sh > /dev/null 2>&1
10 0,4,8,12,16,20,23 * * * /data/rsync.sh > /dev/null 2>&1
 
 
顺便看一下本机/data/backup目录下面的增量备份数据
[root@mysql-node ~]# ll /data/backup
总用量 786364
drwxr-xr-x  2 app app      4096 7月  31 01:01 2018-07-31
drwxr-xr-x  2 app app      4096 8月   1 01:01 2018-08-01
drwxr-xr-x 14 app app      4096 7月  31 00:01 full_20180731-00-01-01
drwxr-xr-x 14 app app      4096 8月   1 00:01 full_20180801-00-01-01
drwxr-xr-x  9 app app      4096 7月  31 04:01 incr_20180731-04-01-01
drwxr-xr-x  9 app app      4096 7月  31 08:01 incr_20180731-08-01-01
drwxr-xr-x  9 app app      4096 7月  31 12:01 incr_20180731-12-01-01
drwxr-xr-x  9 app app      4096 7月  31 16:01 incr_20180731-16-01-01
drwxr-xr-x  9 app app      4096 7月  31 20:01 incr_20180731-20-01-01
drwxr-xr-x  9 app app      4096 7月  31 23:01 incr_20180731-23-01-01
drwxr-xr-x  9 app app      4096 8月   1 04:01 incr_20180801-04-01-01
drwxr-xr-x  9 app app      4096 8月   1 08:01 incr_20180801-08-01-01
drwxr-xr-x  9 app app      4096 8月   1 12:01 incr_20180801-12-01-01
drwxr-xr-x  9 app app      4096 8月   1 16:01 incr_20180801-16-01-01
drwxr-xr-x  9 app app      4096 8月   1 20:01 incr_20180801-20-01-01
drwxr-xr-x  9 app app      4096 8月   1 23:01 incr_20180801-23-01-01
drwxrwxr-x  2 app app     20480 8月   9 08:01 logs
drwxrwxr-x  3 app app      4096 7月  12 17:43 script

lvm-snapshot：基于LVM快照的备份
1.关于快照：
1）事务日志跟数据文件必须在同一个卷上；
2）刚刚创立的快照卷，里面没有任何数据，所有数据均来源于原卷
3）一旦原卷数据发生修改，修改的数据将复制到快照卷中，此时访问数据一部分来自于快照卷，一部分来自于原卷
4）当快照使用过程中，如果修改的数据量大于快照卷容量，则会导致快照卷崩溃。
5）快照卷本身不是备份，只是提供一个时间一致性的访问目录。

2.基于快照备份几乎为热备：
1）创建快照卷之前，要请求MySQL的全局锁；在快照创建完成之后释放锁；
2）如果是Inoodb引擎， 当flush tables 后会有一部分保存在事务日志中，却不在文件中。 因此恢复时候，需要事务日志和数据文件
但释放锁以后，事务日志的内容会同步数据文件中，因此备份内容并不绝对是锁释放时刻的内容，由于有些为完成的事务已经完成，但在备份数据中因为没完成而回滚。 因此需要借助二进制日志往后走一段

3.基于快照备份注意事项：
1）事务日志跟数据文件必须在同一个卷上；
2）创建快照卷之前，要请求MySQL的全局锁；在快照创建完成之后释放锁；
3）请求全局锁完成之后，做一次日志滚动；做二进制日志文件及位置标记(手动进行)；

4.为什么基于MySQL快照的备份很好？
原因如下几点：
1）几乎是热备 在大多数情况下，可以在应用程序仍在运行的时候执行备份。无需关机，只需设置为只读或者类似只读的限制。
2）支持所有基于本地磁盘的存储引擎 它支持MyISAM, Innodb, BDB，还支持 Solid, PrimeXT 和 Falcon。
3）快速备份 只需拷贝二进制格式的文件，在速度方面无以匹敌。
4）低开销 只是文件拷贝，因此对服务器的开销很细微。
5）容易保持完整性 想要压缩备份文件吗？把它们备份到磁带上，FTP或者网络备份软件 -- 十分简单，因为只需要拷贝文件即可。
6）快速恢复 恢复的时间和标准的MySQL崩溃恢复或数据拷贝回去那么快，甚至可能更快，将来会更快。
7）免费 无需额外的商业软件，只需Innodb热备工具来执行备份。

快照备份mysql的缺点：
1）需要兼容快照 -- 这是明显的。
2）需要超级用户(root) 在某些组织，DBA和系统管理员来自不同部门不同的人，因此权限各不一样。
3）停工时间无法预计，这个方法通常指热备，但是谁也无法预料到底是不是热备 -- FLUSH TABLES WITH READ LOCK 可能会需要执行很长时间才能完成。
4）多卷上的数据问题 如果你把日志放在独立的设备上或者你的数据库分布在多个卷上，这就比较麻烦了，因为无法得到全部数据库的一致性快照。不过有些系统可能能自动做到多卷快照。

下面即是使用lvm-snapshot快照方式备份mysql的操作记录，仅依据本人实验中使用而述.

操作记录：
如下环境，本机是在openstack上开的云主机，在openstack上创建一个30G的云硬盘挂载到本机，然后制作lvm逻辑卷。

一、准备LVM卷，并将mysql数据恢复(或者说迁移）到LVM卷上：
1） 创建一个分区或保存到另一块硬盘上面
2） 创建PV、VG、LVM
3） 格式化 LV0
4） 挂载LV到临时目录
5） 确认服务处于stop状态
6） 将数据迁移到LV0
7） 重新挂载LV0到mysql数据库的主目录/data/mysql/data
8） 审核权限并启动服务
[root@test-huanqiu ~]# fdisk -l
.........
Disk /dev/vdc: 32.2 GB, 32212254720 bytes
16 heads, 63 sectors/track, 62415 cylinders
Units = cylinders of 1008 * 512 = 516096 bytes
Sector size (logical/physical): 512 bytes / 512 bytes
I/O size (minimum/optimal): 512 bytes / 512 bytes
Disk identifier: 0x00000000

[root@test-huanqiu ~]# fdisk /dev/vdc                            //依次输入p->n->p->1->回车->回车->w
.........
Command (m for help): p

Disk /dev/vdc: 32.2 GB, 32212254720 bytes
16 heads, 63 sectors/track, 62415 cylinders
Units = cylinders of 1008 * 512 = 516096 bytes
Sector size (logical/physical): 512 bytes / 512 bytes
I/O size (minimum/optimal): 512 bytes / 512 bytes
Disk identifier: 0x343250e4

Device Boot Start End Blocks Id System

Command (m for help): n
Command action
e extended
p primary partition (1-4)
p
Partition number (1-4): 1
First cylinder (1-62415, default 1):
Using default value 1
Last cylinder, +cylinders or +size{K,M,G} (1-62415, default 62415):
Using default value 62415

Command (m for help): w
The partition table has been altered!

Calling ioctl() to re-read partition table.
Syncing disks.

[root@test-huanqiu ~]# fdisk /dev/vdc

WARNING: DOS-compatible mode is deprecated. It's strongly recommended to
switch off the mode (command 'c') and change display units to
sectors (command 'u').

Command (m for help): p

Disk /dev/vdc: 32.2 GB, 32212254720 bytes
16 heads, 63 sectors/track, 62415 cylinders
Units = cylinders of 1008 * 512 = 516096 bytes
Sector size (logical/physical): 512 bytes / 512 bytes
I/O size (minimum/optimal): 512 bytes / 512 bytes
Disk identifier: 0x343250e4

Device Boot Start End Blocks Id System
/dev/vdc1 1 62415 31457128+ 5 Extended

Command (m for help):

[root@test-huanqiu ~]# pvcreate /dev/vdc1
Device /dev/vdc1 not found (or ignored by filtering).
[root@test-huanqiu ~]# vgcreate vg0 /dev/vdc1
Volume group "vg0" successfully created
[root@test-huanqiu ~]# lvcreate -L +3G -n lv0 vg0
Logical volume "lv0" created.
[root@test-huanqiu ~]# mkfs.ext4 /dev/vg0/lv0
[root@test-huanqiu ~]# mkdir /var/lv0/
[root@test-huanqiu ~]# mount /dev/vg0/lv0 /var/lv0/
[root@test-huanqiu ~]# df -h
Filesystem Size Used Avail Use% Mounted on
/dev/mapper/VolGroup00-LogVol00
8.1G 6.0G 1.7G 79% /
tmpfs 1.9G 0 1.9G 0% /dev/shm
/dev/vda1 190M 37M 143M 21% /boot
/dev/mapper/vg0-lv0 2.9G 4.5M 2.8G 1% /var/lv0

[root@test-huanqiu ~]# lvs
LV VG Attr LSize Pool Origin Data% Meta% Move Log Cpy%Sync Convert
LogVol00 VolGroup00 -wi-ao---- 8.28g
LogVol01 VolGroup00 -wi-ao---- 1.50g
lv0 vg0 -wi-a----- 3.00g

----------------------------------------------------------------------------------------------------
如果要想删除这个lvs，操作如下：
[root@test-huanqiu ~]# umount /data/mysql/data/            //先卸载掉这个lvs的挂载关系
[root@test-huanqiu ~]# lvremove /dev/vg0/lv0
[root@test-huanqiu ~]# vgremove vg0
[root@test-huanqiu ~]# pvremove /dev/vdc1
[root@test-huanqiu ~]# lvs
LV VG Attr LSize Pool Origin Data% Meta% Move Log Cpy%Sync Convert
LogVol00 VolGroup00 -wi-ao---- 8.28g
LogVol01 VolGroup00 -wi-ao---- 1.50g
----------------------------------------------------------------------------------------------------

mysql的数据目录是/data/mysql/data,密码是123456
[root@test-huanqiu ~]# ps -ef|grep mysql
mysql 2066 1286 0 07:33 ? 00:00:06 /usr/local/mysql/bin/mysqld --basedir=/usr/local/mysql/ --datadir=/data/mysql/data --plugin-dir=/usr/local/mysql//lib/plugin --user=mysql --log-error=/data/mysql/data/mysql-error.log --pid-file=/data/mysql/data/mysql.pid --socket=/usr/local/mysql/var/mysql.sock --port=3306
root 2523 2471 0 07:55 pts/1 00:00:00 grep mysql
[root@test-huanqiu ~]# /etc/init.d/mysql stop
Shutting down MySQL.... SUCCESS!

[root@test-huanqiu ~]# cd /data/mysql/data/
[root@test-huanqiu data]# tar -cf - . | tar xf - -C /var/lv0/

[root@test-huanqiu data]# umount /var/lv0/

[root@test-huanqiu data]# mount /dev/vg0/lv0 /data/mysql/data
[root@test-huanqiu data]# df -h
Filesystem Size Used Avail Use% Mounted on
/dev/mapper/VolGroup00-LogVol00
8.1G 6.0G 1.7G 79% /
tmpfs 1.9G 0 1.9G 0% /dev/shm
/dev/vda1 190M 37M 143M 21% /boot
/dev/mapper/vg0-lv0 2.9G 164M 2.6G 6% /data/mysql/data

删除挂载后产生的lost+found目录
[root@test-huanqiu data]# rm -rf lost+found

[root@test-huanqiu data]# ll -d /data/mysql/data
[root@test-huanqiu data]# ll -Z /data/mysql/data
[root@test-huanqiu data]# ll -Zd /data/mysql/data

需要注意的是:
当SElinux功能开启情况下，mysql数据库重启会失败，所以必须执行下面命令，恢复SElinux安全上下文.
[root@test-huanqiu data]# restorecon -R /data/mysql/data/
[root@test-huanqiu data]# /etc/init.d/mysql start
Starting MySQL... SUCCESS!

二、备份： (生产环境下一般都是整个数据库备份)
1）锁表
2）查看position号并记录，便于后期恢复
3）创建snapshot快照
4）解表
5）挂载snapshot
6）拷贝snapshot数据，进行备份。备份整个数据库之前，要关闭mysql服务（保护ibdata1文件）
7）移除快照

设置此变量为1，让每个事件尽可能同步到二进制日志文件里，以消耗IO来尽可能确保数据一致性。
mysql> SET GLOBAL sync_binlog=1;

查看二进制日志和position，以备后续进行binlog日志恢复增量数据（记住这个position节点记录，对后面的增量数据备份很重要） 
mysql> SHOW MASTER STATUS;
+------------------+----------+--------------+------------------+-------------------+
| File | Position | Binlog_Do_DB | Binlog_Ignore_DB | Executed_Gtid_Set |
+------------------+----------+--------------+------------------+-------------------+
| mysql-bin.000004 | 1434 | | | |
+------------------+----------+--------------+------------------+-------------------+
1 row in set (0.00 sec)

创建存放binlog日志的position节点记录的目录
所有的position节点记录都放在这同一个binlog.pos文件下（后面就使用>>符号追加到这个文件下）
[root@test-huanqiu ~]# mkdir /backup/mysql/binlog
[root@test-huanqiu ~]# mysql -p123456 -e "SHOW MASTER STATUS;" > /backup/mysql/binlog/binlog.pos
[root@test-huanqiu snap1]# cat /backup/mysql/binlog/binlog.pos
File Position Binlog_Do_DB Binlog_Ignore_DB Executed_Gtid_Set
mysql-bin.000004 1434

刷新日志，产生新的binlog日志，保证日志信息不会再写入到上面的mysql-bin.000004日志内。
mysql> FLUSH LOGS;

全局读锁，读锁请求到后不要关闭此mysql交互界面
mysql> FLUSH TABLES WITH READ LOCK;

在innodb表中，即使是请求到了读锁，但InnoDB在后台依然可能会有事务在进行读写操作，
可用"mysql> SHOW ENGINE INNODB STATUS;"查看后台进程的状态，等没有写请求后再做备份。

创建快照，以只读的方式（--permission r）创建一个3GB大小的快照卷snap1
-s：相当于--snapshot
[root@test-huanqiu ~]# mkdir /var/snap1
[root@test-huanqiu ~]# lvcreate -s -L 2G -n snap1 /dev/vg0/lv0 --permission r
Logical volume "snap1" created.

查看快照卷的详情（快照卷也是LV）：
[root@test-huanqiu ~]# lvdisplay

解除锁定
回到锁定表的mysql交互式界面，解锁：
mysql> UNLOCK TABLES;

此参数可以根据服务器磁盘IO的负载来调整
mysql> SET GLOBAL sync_binlog=0;

[root@test-huanqiu ~]# mount /dev/vg0/snap1 /var/snap1                //挂载快照卷
[root@test-huanqiu snap1]# df -h
Filesystem             Size   Used  Avail  Use%  Mounted on
/dev/mapper/VolGroup00-LogVol00
                            8.1G  5.8G  1.9G  76%    /
tmpfs                    1.9G  0       1.9G  0%      /dev/shm
/dev/vda1              190M 37M   143M 21%    /boot
/dev/mapper/vg0-lv0 2.9G 115M 2.7G 5% /data/mysql/data
/dev/mapper/vg0-snap1
                               2.9G 115M 2.7G 5% /var/snap1

[root@test-huanqiu ~]# cd /var/snap1/ && ll /var/snap1
[root@test-huanqiu snap1]# mkdir -p /backup/mysql/data/               //创建备份目录
total 0

对本机的数据库进行备份，备份整个数据库。

1
2
3
4
5
6
7
8
9
10
11
12
13
14
15
16
17
18
19
20
21
22
23
24
25
26
27
28
29
30
31
32
33
34
35
36
37
38
39
40
41
42
43
44
45
46
47
48
49
mysql> show databases;
+--------------------+
| Database           |
+--------------------+
| information_schema |
| mysql              |
| performance_schema |
| test               |
+--------------------+
4 rows in set (0.01 sec)
 
mysql> create database beijing;
Query OK, 1 row affected (0.00 sec)
 
mysql> use beijing;
Database changed
 
mysql> create table people(id int(5),name varchar(20));
Query OK, 0 rows affected (0.03 sec)
 
mysql> insert into people values("1","wangshibo");
Query OK, 1 row affected (0.00 sec)
 
mysql> insert into people values("2","guohuihui");
Query OK, 1 row affected (0.01 sec)
 
mysql> insert into people values("3","wuxiang");
Query OK, 1 row affected (0.01 sec)
 
mysql> select * from people;
+------+-----------+
| id   | name      |
+------+-----------+
|    1 | wangshibo |
|    2 | guohuihui |
|    3 | wuxiang   |
+------+-----------+
3 rows in set (0.00 sec)
mysql> show databases;
+--------------------+
| Database           |
+--------------------+
| information_schema |
| beijing            |
| mysql              |
| performance_schema |
| test               |
+--------------------+
5 rows in set (0.01 sec)
--------------------------------------------------------------------------------------------------------------------------
需要注意的是：
innodb表，一般会打开独立表空间模式(innodb_file_per_table)。
由于InnoDB默认会将所有的数据库InnoDB引擎的表数据存储在一个共享空间中：ibdata1文件。
增删数据库的时候，ibdata1文件不会自动收缩，这对单个或部分数据库的备份也将成为问题（如果不是整个数据库备份的情况下，ibdata1文件就不能备份，否则会影响全部数据库的数据）。
所以若是对单个数据库或部分数据库进行快照备份：
1）若是直接误删除mysql数据目录下备份库目录，可以直接将快照备份数据解压就能恢复
2）若是使用drop或delete误删除的数据，那么在使用快照备份数据恢复时，就会出问题！因为单库备份时ibdata1文件不能单独备份，恢复时会导致这个文件损坏！

所以正确的做法是：
要对整个数据库进行备份，并且一定要在mysql服务关闭的情况下（这样是为了保护ibdata1文件）。
因为mysql是采用缓冲方式来将数据写入到ibdata1文件中的，这正是fflush()函数存在的理由。当mysql在运行时，对ibdata1进行拷贝肯定会导致ibdata1文件中的数据出错，这样在数据恢复时，也就肯定会出现“ERROR 1146 (42S02): Table '****' doesn't exist“的报错！

在对启用innodb引擎的mysql数据库进行迁移的时候也是同理：
在对innodb数据库进行数据迁移的时候，即将msyql(innodb引擎)服务从一台服务器迁移到另一台服务器时，在对数据库目录进行整体拷贝的时候（当然就包括了对ibdata1文件拷贝），一定要在关闭对方mysql服务的情况下进行拷贝！

ibdata1用来储存文件的数据，而库名的文件夹里面的那些表文件只是结构而已，由于新版的mysql默认试innodb，所以ibdata1文件默认就存在了，少了这个文件有的数据表就会出错。要知道：数据库目录下的.frm文件是数据库中很多的表的结构描述文件；而ibdata1文件才是数据库的真实数据存放文件。

-------------------------------------------innodb_file_per_table参数说明------------------------------------------
线上环境的话，一般都建议打开这个独立表空间模式。
因为ibdata1文件会不断的增大，不会减少，无法向OS回收空间，容易导致线上出现过大的共享表空间文件，致使当前空间爆满。
并且ibdata1文件大到一定程序会影响insert、update的速度；并且
另外如果删表频繁的话，共享表空间产生的碎片会比较多。打开独立表空间，方便进行innodb表的碎片整理

使用MyISAM表引擎的数据库会分别创建三个文件：表结构、表索引、表数据空间。
可以将某个数据库目录直接迁移到其他数据库也可以正常工作。

然而当使用InnoDB的时候，一切都变了。
InnoDB默认会将所有的数据库InnoDB引擎的表数据存储在一个共享空间中：ibdata1文件。
增删数据库的时候，ibdata1文件不会自动收缩，单个数据库的备份也将成为问题。
通常只能将数据使用mysqldump 导出，然后再导入解决这个问题。

在MySQL的配置文件[mysqld]部分，增加innodb_file_per_table参数。
可以修改InnoDB为独立表空间模式，每个数据库的每个表都会生成一个数据空间。

它的优点：
1）每个表都有自已独立的表空间。
2）每个表的数据和索引都会存在自已的表空间中。
3）可以实现单表在不同的数据库中移动。
4）空间可以回收（除drop table操作处，表空不能自已回收）

Drop table操作自动回收表空间，如果对于统计分析或是日值表，删除大量数据后可以通过:alter table TableName engine=innodb;回缩不用的空间。
对于使innodb-plugin的Innodb使用turncate table也会使空间收缩。
对于使用独立表空间的表，不管怎么删除，表空间的碎片不会太严重的影响性能，而且还有机会处理。

它的缺点：
单表增加过大，如超过100个G。

结论：
共享表空间在Insert操作上少有优势。其它都没独立表空间表现好。当启用独立表空间时，请合理调整一下：innodb_open_files。
InnoDB Hot Backup（冷备）的表空间cp不会面对很多无用的copy了。而且利用innodb hot backup及表空间的管理命令可以实。

1）innodb_file_per_table设置.设置为1，表示打开了独立的表空间模式。 如果设置为0，表示关闭独立表空间模式，开启方法如下：
在my.cnf中[mysqld]下设置
innodb_file_per_table=1

2）查看是否开启：
mysql> show variables like "%per_table%";
+-----------------------+-------+
| Variable_name | Value |
+-----------------------+-------+
| innodb_file_per_table | ON |
+-----------------------+-------+
1 row in set (0.00 sec)

3）关闭独享表空间
innodb_file_per_table=0关闭独立的表空间
mysql> show variables like ‘%per_table%’;
-------------------------------------------innodb_file_per_table参数说明------------------------------------------
--------------------------------------------------------------------------------------------------------------------------

备份前，一定要关闭mysql数据库！因为里面会涉及到ibdata1文件备份，不关闭mysql的话，ibdata1文件备份后会损坏，从而导致恢复数据失败！
[root@test-huanqiu snap1]# /etc/init.d/mysql stop
Shutting down MySQL.... SUCCESS!
[root@test-huanqiu data]# lsof -i:3306
[root@test-huanqiu data]#

现在备份整个数据库
[root@test-huanqiu snap1]# tar -zvcf /backup/mysql/data/`date +%Y-%m-%d`dbbackup.tar.gz ./
[root@test-huanqiu snap1]# ll /backup/mysql/data/
total 384
-rw-r--r--. 1 root root 392328 Dec 5 22:15 2016-12-05dbbackup.tar.gz

释放快照卷，每次备份之后，应该删除快照，减少IO操作
先卸载，再删除
[root@test-huanqiu ~]# umount /var/snap1/
[root@test-huanqiu ~]# df -h                //确认上面的挂载关系已经没了
Filesystem Size Used Avail Use% Mounted on
/dev/mapper/VolGroup00-LogVol00
8.1G 5.8G 1.9G 76% /
tmpfs 1.9G 0 1.9G 0% /dev/shm
/dev/vda1 190M 37M 143M 21% /boot
/dev/mapper/vg0-lv0 2.9G 115M 2.7G 5% /data/mysql/data
[root@test-huanqiu ~]# lvremove /dev/vg0/snap1
Do you really want to remove active logical volume snap1? [y/n]: y
Logical volume "snap1" successfully removed


数据被快照备份后，可以启动数据库
[root@test-huanqiu ~]# /etc/init.d/mysql start
Starting MySQL.. SUCCESS!
[root@test-huanqiu ~]# lsof -i:3306
COMMAND PID USER FD TYPE DEVICE SIZE/OFF NODE NAME
mysqld 15943 mysql 16u IPv4 93348 0t0 TCP *:mysql (LISTEN)
[root@test-huanqiu ~]#

现在再进行新的数据写入：

1
2
3
4
5
6
7
8
9
10
11
12
13
14
15
16
17
18
19
20
21
22
23
24
25
26
27
28
29
30
31
32
33
34
35
36
37
38
39
40
41
42
43
44
45
46
47
48
49
mysql> use beijing;
Database changed
mysql> insert into people values("4","liumengnan");
Query OK, 1 row affected (0.02 sec)
 
mysql> insert into people values("5","zhangjuanjuan");
Query OK, 1 row affected (0.00 sec)
 
mysql> select * from people;
+------+---------------+
| id   | name          |
+------+---------------+
|    1 | wangshibo     |
|    2 | guohuihui     |
|    3 | wuxiang       |
|    4 | liumengnan    |
|    5 | zhangjuanjuan |
+------+---------------+
5 rows in set (0.00 sec)
 
mysql> create table heihei(name varchar(20),age varchar(20));
Query OK, 0 rows affected (0.02 sec)
 
mysql> insert into heihei values("jiujiujiu","nan");
Query OK, 1 row affected (0.00 sec)
 
mysql> select * from heihei;
+-----------+------+
| name      | age  |
+-----------+------+
| jiujiujiu | nan  |
+-----------+------+
1 row in set (0.00 sec)
 
mysql> create database shanghai;
Query OK, 1 row affected (0.01 sec)
 
mysql> show databases;
+--------------------+
| Database           |
+--------------------+
| information_schema |
| beijing            |
| mysql              |
| performance_schema |
| shanghai           |
| test               |
+--------------------+
6 rows in set (0.00 sec)
假设一不小心误操作删除beijing和shanghai库

1
2
3
4
5
6
7
8
9
10
11
12
13
14
15
16
mysql> drop database beijing;
Query OK, 2 rows affected (0.03 sec)
 
mysql> drop database shanghai;
Query OK, 0 rows affected (0.00 sec)
 
mysql> show databases;
+--------------------+
| Database           |
+--------------------+
| information_schema |
| mysql              |
| performance_schema |
| test               |
+--------------------+
4 rows in set (0.00 sec)
莫慌！接下来就说下数据恢复操作~~
三、恢复流程如下：
0）由于涉及到增量数据备份，所以提前将最近一次的binlog日志从mysql数据目录复制到别的路径下
1）在mysql数据库中执行flush logs命令，产生新的binlog日志，让日志信息写入到新的这个binlog日志中
1）关闭数据库，一定要关闭
2）删除数据目录下的文件
3）快照数据拷贝回来，position节点记录回放
4）增量数据就利用mysqlbinlog命令将上面提前拷贝的binlog日志文件导出为sql文件，并剔除其中的drop语句，然后进行恢复。
5）重启数据

先将最新一次的binlog日志备份到别处，用作增量数据备份。
比如mysql-bin.000006是最新一次的binlog日志
[root@test-huanqiu data]# cp mysql-bin.000006 /backup/mysql/data/

产生新的binlog日志，确保日志写入到这个新的binlog日志内，而不再写入到上面备份的binlog日志里。
mysql> flush logs;

[root@test-huanqiu data]# ll mysql-bin.000007
-rw-rw----. 1 mysql mysql 120 Dec 5 23:19 mysql-bin.000007

[root@test-huanqiu data]# /etc/init.d/mysql stop
Shutting down MySQL.... SUCCESS!
[root@test-huanqiu data]# lsof -i:3306
[root@test-huanqiu data]# pwd
/data/mysql/data
[root@test-huanqiu data]# rm -rf ./*
[root@test-huanqiu data]# tar -zvxf /backup/mysql/data/2016-12-05dbbackup.tar.gz ./

[root@test-huanqiu data]# /etc/init.d/mysql start
Starting MySQL SUCCESS!
[root@test-huanqiu data]# cat /backup/mysql/binlog/binlog.pos
File Position Binlog_Do_DB Binlog_Ignore_DB Executed_Gtid_Set
mysql-bin.000004 1434
[root@test-huanqiu data]# mysqlbinlog --start-position=1434 /data/mysql/data/mysql-bin.000004 | mysql -p123456

登陆数据库查看，发现这只是恢复到快照备份阶段的数据：

1
2
3
4
5
6
7
8
9
10
11
12
13
14
15
16
17
18
19
20
21
22
23
mysql> show databases;
+--------------------+
| Database           |
+--------------------+
| information_schema |
| beijing            |
| mysql              |
| performance_schema |
| test               |
+--------------------+
5 rows in set (0.00 sec)
 
mysql> select * from beijing.people;
+------+-----------+
| id   | name      |
+------+-----------+
|    1 | wangshibo |
|    2 | guohuihui |
|    3 | wuxiang   |
+------+-----------+
3 rows in set (0.00 sec)
 
mysql>
快照备份之后写入的数据要利用mysqlbinlog命令将上面拷贝的mysql-bin000006文件导出为sql文件，并剔除其中的drop语句，然后进行恢复。
[root@test-huanqiu ~]# cd /backup/mysql/data/
[root@test-huanqiu data]# ll
total 388
-rw-r--r--. 1 root root 392328 Dec 5 22:15 2016-12-05dbbackup.tar.gz
-rw-r-----. 1 root root 1274 Dec 5 23:19 mysql-bin.000006
[root@test-huanqiu data]# mysqlbinlog mysql-bin.000006 >000006bin.sql

剔除其中的drop语句
[root@test-huanqiu data]# vim 000006bin.sql          //手动删除sql语句中的drop语句

然后在mysql中使用source命令恢复数据
mysql> source /backup/mysql/data/000006bin.sql;

再次查看下，发现增量部分的数据也已经恢复回来了

1
2
3
4
5
6
7
8
9
10
11
12
13
14
15
16
17
18
19
20
21
22
23
24
25
26
27
28
29
30
31
32
33
34
35
36
37
38
39
40
41
42
43
mysql> show databases;
+--------------------+
| Database           |
+--------------------+
| information_schema |
| beijing            |
| mysql              |
| performance_schema |
| shanghai           |
| test               |
+--------------------+
6 rows in set (0.00 sec)
 
mysql> use beijing;
Database changed
mysql> show tables;
+-------------------+
| Tables_in_beijing |
+-------------------+
| heihei            |
| people            |
+-------------------+
2 rows in set (0.00 sec)
 
mysql> select * from people;
+------+---------------+
| id   | name          |
+------+---------------+
|    1 | wangshibo     |
|    2 | guohuihui     |
|    3 | wuxiang       |
|    4 | liumengnan    |
|    5 | zhangjuanjuan |
+------+---------------+
5 rows in set (0.00 sec)
 
mysql> select * from heihei;
+-----------+------+
| name      | age  |
+-----------+------+
| jiujiujiu | nan  |
+-----------+------+
1 row in set (0.00 sec)

lvm-snapshot：基于LVM快照的备份
1.关于快照：
1）事务日志跟数据文件必须在同一个卷上；
2）刚刚创立的快照卷，里面没有任何数据，所有数据均来源于原卷
3）一旦原卷数据发生修改，修改的数据将复制到快照卷中，此时访问数据一部分来自于快照卷，一部分来自于原卷
4）当快照使用过程中，如果修改的数据量大于快照卷容量，则会导致快照卷崩溃。 
5）快照卷本身不是备份，只是提供一个时间一致性的访问目录。

2.基于快照备份几乎为热备： 
1）创建快照卷之前，要请求MySQL的全局锁；在快照创建完成之后释放锁；
2）如果是Inoodb引擎， 当flush tables 后会有一部分保存在事务日志中，却不在文件中。 因此恢复时候，需要事务日志和数据文件
但释放锁以后，事务日志的内容会同步数据文件中，因此备份内容并不绝对是锁释放时刻的内容，由于有些为完成的事务已经完成，但在备份数据中因为没完成而回滚。 因此需要借助二进制日志往后走一段。

3.基于快照备份注意事项： 
1）事务日志跟数据文件必须在同一个卷上；
2）创建快照卷之前，要请求MySQL的全局锁；在快照创建完成之后释放锁；
3）请求全局锁完成之后，做一次日志滚动；做二进制日志文件及位置标记(手动进行)；

4.为什么基于MySQL快照的备份很好？
原因如下几点：
1）几乎是热备 在大多数情况下，可以在应用程序仍在运行的时候执行备份。无需关机，只需设置为只读或者类似只读的限制。
2）支持所有基于本地磁盘的存储引擎 它支持MyISAM, Innodb, BDB，还支持 Solid, PrimeXT 和 Falcon。
3）快速备份 只需拷贝二进制格式的文件，在速度方面无以匹敌。
4）低开销 只是文件拷贝，因此对服务器的开销很细微。
5）容易保持完整性 想要压缩备份文件吗？把它们备份到磁带上，FTP或者网络备份软件 -- 十分简单，因为只需要拷贝文件即可。
6）快速恢复 恢复的时间和标准的MySQL崩溃恢复或数据拷贝回去那么快，甚至可能更快，将来会更快。
7）免费 无需额外的商业软件，只需Innodb热备工具来执行备份。

快照备份mysql的缺点：
1）需要兼容快照 -- 这是明显的。
2）需要超级用户(root) 在某些组织，DBA和系统管理员来自不同部门不同的人，因此权限各不一样。
3）停工时间无法预计，这个方法通常指热备，但是谁也无法预料到底是不是热备 -- FLUSH TABLES WITH READ LOCK 可能会需要执行很长时间才能完成。
4）多卷上的数据问题 如果你把日志放在独立的设备上或者你的数据库分布在多个卷上，这就比较麻烦了，因为无法得到全部数据库的一致性快照。不过有些系统可能能自动做到多卷快照。

下面即是使用lvm-snapshot快照方式备份mysql的操作记录，仅依据本人实验中使用而述.

操作记录：
如下环境，本机是在openstack上开的云主机，在openstack上创建一个30G的云硬盘挂载到本机，然后制作lvm逻辑卷。

一、准备LVM卷，并将mysql数据恢复(或者说迁移）到LVM卷上：
1） 创建一个分区或保存到另一块硬盘上面
2） 创建PV、VG、LVM
3） 格式化 LV0
4） 挂载LV到临时目录
5） 确认服务处于stop状态
6） 将数据迁移到LV0
7） 重新挂载LV0到mysql数据库的主目录/data/mysql/data
8） 审核权限并启动服务

[root@test-huanqiu ~]# fdisk -l
.........
Disk /dev/vdc: 32.2 GB, 32212254720 bytes
16 heads, 63 sectors/track, 62415 cylinders
Units = cylinders of 1008 * 512 = 516096 bytes
Sector size (logical/physical): 512 bytes / 512 bytes
I/O size (minimum/optimal): 512 bytes / 512 bytes
Disk identifier: 0x00000000
[root@test-huanqiu ~]# fdisk /dev/vdc                            //依次输入p->n->p->1->回车->回车->w
.........
Command (m for help): p
Disk /dev/vdc: 32.2 GB, 32212254720 bytes
16 heads, 63 sectors/track, 62415 cylinders
Units = cylinders of 1008 * 512 = 516096 bytes
Sector size (logical/physical): 512 bytes / 512 bytes
I/O size (minimum/optimal): 512 bytes / 512 bytes
Disk identifier: 0x343250e4
Device Boot Start End Blocks Id System
Command (m for help): n
Command action
e extended
p primary partition (1-4)
p
Partition number (1-4): 1
First cylinder (1-62415, default 1): 
Using default value 1
Last cylinder, +cylinders or +size{K,M,G} (1-62415, default 62415): 
Using default value 62415
 
Command (m for help): w
The partition table has been altered!
 
Calling ioctl() to re-read partition table.
Syncing disks.
[root@test-huanqiu ~]# fdisk /dev/vdc
WARNING: DOS-compatible mode is deprecated. It's strongly recommended to
switch off the mode (command 'c') and change display units to
sectors (command 'u').
Command (m for help): p
Disk /dev/vdc: 32.2 GB, 32212254720 bytes
16 heads, 63 sectors/track, 62415 cylinders
Units = cylinders of 1008 * 512 = 516096 bytes
Sector size (logical/physical): 512 bytes / 512 bytes
I/O size (minimum/optimal): 512 bytes / 512 bytes
Disk identifier: 0x343250e4
Device Boot Start End Blocks Id System
/dev/vdc1 1 62415 31457128+ 5 Extended
Command (m for help):
[root@test-huanqiu ~]# pvcreate /dev/vdc1 
Device /dev/vdc1 not found (or ignored by filtering).
[root@test-huanqiu ~]# vgcreate vg0 /dev/vdc1
Volume group "vg0" successfully created
[root@test-huanqiu ~]# lvcreate -L +3G -n lv0 vg0
Logical volume "lv0" created.
[root@test-huanqiu ~]# mkfs.ext4 /dev/vg0/lv0 
[root@test-huanqiu ~]# mkdir /var/lv0/
[root@test-huanqiu ~]# mount /dev/vg0/lv0 /var/lv0/
[root@test-huanqiu ~]# df -h
Filesystem Size Used Avail Use% Mounted on
/dev/mapper/VolGroup00-LogVol00
8.1G 6.0G 1.7G 79% /
tmpfs 1.9G 0 1.9G 0% /dev/shm
/dev/vda1 190M 37M 143M 21% /boot
/dev/mapper/vg0-lv0 2.9G 4.5M 2.8G 1% /var/lv0
[root@test-huanqiu ~]# lvs
LV VG Attr LSize Pool Origin Data% Meta% Move Log Cpy%Sync Convert
LogVol00 VolGroup00 -wi-ao---- 8.28g 
LogVol01 VolGroup00 -wi-ao---- 1.50g 
lv0 vg0 -wi-a----- 3.00g
----------------------------------------------------------------------------------------------------
如果要想删除这个lvs，操作如下：
[root@test-huanqiu ~]# umount /data/mysql/data/            //先卸载掉这个lvs的挂载关系
[root@test-huanqiu ~]# lvremove /dev/vg0/lv0 
[root@test-huanqiu ~]# vgremove vg0
[root@test-huanqiu ~]# pvremove /dev/vdc1
[root@test-huanqiu ~]# lvs
LV VG Attr LSize Pool Origin Data% Meta% Move Log Cpy%Sync Convert
LogVol00 VolGroup00 -wi-ao---- 8.28g 
LogVol01 VolGroup00 -wi-ao---- 1.50g 
----------------------------------------------------------------------------------------------------
mysql的数据目录是/data/mysql/data,密码是123456
[root@test-huanqiu ~]# ps -ef|grep mysql
mysql 2066 1286 0 07:33 ? 00:00:06 /usr/local/mysql/bin/mysqld --basedir=/usr/local/mysql/ --datadir=/data/mysql/data --plugin-dir=/usr/local/mysql//lib/plugin --user=mysql --log-error=/data/mysql/data/mysql-error.log --pid-file=/data/mysql/data/mysql.pid --socket=/usr/local/mysql/var/mysql.sock --port=3306
root 2523 2471 0 07:55 pts/1 00:00:00 grep mysql
[root@test-huanqiu ~]# /etc/init.d/mysql stop
Shutting down MySQL.... SUCCESS!
 
[root@test-huanqiu ~]# cd /data/mysql/data/
[root@test-huanqiu data]# tar -cf - . | tar xf - -C /var/lv0/
 
[root@test-huanqiu data]# umount /var/lv0/
 
[root@test-huanqiu data]# mount /dev/vg0/lv0 /data/mysql/data
[root@test-huanqiu data]# df -h
Filesystem Size Used Avail Use% Mounted on
/dev/mapper/VolGroup00-LogVol00
8.1G 6.0G 1.7G 79% /
tmpfs 1.9G 0 1.9G 0% /dev/shm
/dev/vda1 190M 37M 143M 21% /boot
/dev/mapper/vg0-lv0 2.9G 164M 2.6G 6% /data/mysql/data
 
删除挂载后产生的lost+found目录
[root@test-huanqiu data]# rm -rf lost+found
 
[root@test-huanqiu data]# ll -d /data/mysql/data
[root@test-huanqiu data]# ll -Z /data/mysql/data
[root@test-huanqiu data]# ll -Zd /data/mysql/data
 
需要注意的是:
当SElinux功能开启情况下，mysql数据库重启会失败，所以必须执行下面命令，恢复SElinux安全上下文.
[root@test-huanqiu data]# restorecon -R /data/mysql/data/
[root@test-huanqiu data]# /etc/init.d/mysql start
Starting MySQL... SUCCESS!
二、备份： (生产环境下一般都是整个数据库备份)
1）锁表
2）查看position号并记录，便于后期恢复
3）创建snapshot快照
4）解表
5）挂载snapshot
6）拷贝snapshot数据，进行备份。备份整个数据库之前，要关闭mysql服务（保护ibdata1文件）
7）移除快照

设置此变量为1，让每个事件尽可能同步到二进制日志文件里，以消耗IO来尽可能确保数据一致性。 

mysql> SET GLOBAL sync_binlog=1;
查看二进制日志和position，以备后续进行binlog日志恢复增量数据（记住这个position节点记录，对后面的增量数据备份很重要） 

mysql> SHOW MASTER STATUS; 
+------------------+----------+--------------+------------------+-------------------+
| File | Position | Binlog_Do_DB | Binlog_Ignore_DB | Executed_Gtid_Set |
+------------------+----------+--------------+------------------+-------------------+
| mysql-bin.000004 | 1434 | | | |
+------------------+----------+--------------+------------------+-------------------+
1 row in set (0.00 sec)
创建存放binlog日志的position节点记录的目录
所有的position节点记录都放在这同一个binlog.pos文件下（后面就使用>>符号追加到这个文件下）

[root@test-huanqiu ~]# mkdir /backup/mysql/binlog
[root@test-huanqiu ~]# mysql -p123456 -e "SHOW MASTER STATUS;" > /backup/mysql/binlog/binlog.pos 
[root@test-huanqiu snap1]# cat /backup/mysql/binlog/binlog.pos
File Position Binlog_Do_DB Binlog_Ignore_DB Executed_Gtid_Set
mysql-bin.000004 1434
刷新日志，产生新的binlog日志，保证日志信息不会再写入到上面的mysql-bin.000004日志内。

mysql> FLUSH LOGS;
全局读锁，读锁请求到后不要关闭此mysql交互界面

mysql> FLUSH TABLES WITH READ LOCK;
在innodb表中，即使是请求到了读锁，但InnoDB在后台依然可能会有事务在进行读写操作，
可用"mysql> SHOW ENGINE INNODB STATUS;"查看后台进程的状态，等没有写请求后再做备份。

创建快照，以只读的方式（--permission r）创建一个3GB大小的快照卷snap1
-s：相当于--snapshot

[root@test-huanqiu ~]# mkdir /var/snap1
[root@test-huanqiu ~]# lvcreate -s -L 2G -n snap1 /dev/vg0/lv0 --permission r
Logical volume "snap1" created.
查看快照卷的详情（快照卷也是LV）：

[root@test-huanqiu ~]# lvdisplay
解除锁定
回到锁定表的mysql交互式界面，解锁：

mysql> UNLOCK TABLES;
此参数可以根据服务器磁盘IO的负载来调整

mysql> SET GLOBAL sync_binlog=0;
 
[root@test-huanqiu ~]# mount /dev/vg0/snap1 /var/snap1                //挂载快照卷
[root@test-huanqiu snap1]# df -h
Filesystem             Size   Used  Avail  Use%  Mounted on
/dev/mapper/VolGroup00-LogVol00
                            8.1G  5.8G  1.9G  76%    /
tmpfs                    1.9G  0       1.9G  0%      /dev/shm
/dev/vda1              190M 37M   143M 21%    /boot
/dev/mapper/vg0-lv0 2.9G 115M 2.7G 5% /data/mysql/data
/dev/mapper/vg0-snap1
                               2.9G 115M 2.7G 5% /var/snap1
 
[root@test-huanqiu ~]# cd /var/snap1/ && ll /var/snap1
[root@test-huanqiu snap1]# mkdir -p /backup/mysql/data/               //创建备份目录
total 0
对本机的数据库进行备份，备份整个数据库。

mysql> show databases;
+--------------------+
| Database           |
+--------------------+
| information_schema |
| mysql              |
| performance_schema |
| test               |
+--------------------+
4 rows in set (0.01 sec)
 
mysql> create database beijing;
Query OK, 1 row affected (0.00 sec)
 
mysql> use beijing;
Database changed
 
mysql> create table people(id int(5),name varchar(20));
Query OK, 0 rows affected (0.03 sec)
 
mysql> insert into people values("1","wangshibo");
Query OK, 1 row affected (0.00 sec)
 
mysql> insert into people values("2","guohuihui");
Query OK, 1 row affected (0.01 sec)
 
mysql> insert into people values("3","wuxiang");
Query OK, 1 row affected (0.01 sec)
 
mysql> select * from people;
+------+-----------+
| id   | name      |
+------+-----------+
|    1 | wangshibo |
|    2 | guohuihui |
|    3 | wuxiang   |
+------+-----------+
3 rows in set (0.00 sec)
mysql> show databases;
+--------------------+
| Database           |
+--------------------+
| information_schema |
| beijing            |
| mysql              |
| performance_schema |
| test               |
+--------------------+
5 rows in set (0.01 sec)
--------------------------------------------------------------------------------------------------------------------------
需要注意的是：
innodb表，一般会打开独立表空间模式(innodb_file_per_table)。
由于InnoDB默认会将所有的数据库InnoDB引擎的表数据存储在一个共享空间中：ibdata1文件。
增删数据库的时候，ibdata1文件不会自动收缩，这对单个或部分数据库的备份也将成为问题（如果不是整个数据库备份的情况下，ibdata1文件就不能备份，否则会影响全部数据库的数据）。
所以若是对单个数据库或部分数据库进行快照备份：
1）若是直接误删除mysql数据目录下备份库目录，可以直接将快照备份数据解压就能恢复
2）若是使用drop或delete误删除的数据，那么在使用快照备份数据恢复时，就会出问题！因为单库备份时ibdata1文件不能单独备份，恢复时会导致这个文件损坏！

所以正确的做法是：
要对整个数据库进行备份，并且一定要在mysql服务关闭的情况下（这样是为了保护ibdata1文件）。
因为mysql是采用缓冲方式来将数据写入到ibdata1文件中的，这正是fflush()函数存在的理由。当mysql在运行时，对ibdata1进行拷贝肯定会导致ibdata1文件中的数据出错，这样在数据恢复时，也就肯定会出现“ERROR 1146 (42S02): Table '****' doesn't exist“的报错！

在对启用innodb引擎的mysql数据库进行迁移的时候也是同理：
在对innodb数据库进行数据迁移的时候，即将msyql(innodb引擎)服务从一台服务器迁移到另一台服务器时，在对数据库目录进行整体拷贝的时候（当然就包括了对ibdata1文件拷贝），一定要在关闭对方mysql服务的情况下进行拷贝！

ibdata1用来储存文件的数据，而库名的文件夹里面的那些表文件只是结构而已，由于新版的mysql默认试innodb，所以ibdata1文件默认就存在了，少了这个文件有的数据表就会出错。要知道：数据库目录下的.frm文件是数据库中很多的表的结构描述文件；而ibdata1文件才是数据库的真实数据存放文件。

-------------------------------------------innodb_file_per_table参数说明------------------------------------------
线上环境的话，一般都建议打开这个独立表空间模式。
因为ibdata1文件会不断的增大，不会减少，无法向OS回收空间，容易导致线上出现过大的共享表空间文件，致使当前空间爆满。
并且ibdata1文件大到一定程序会影响insert、update的速度；并且
另外如果删表频繁的话，共享表空间产生的碎片会比较多。打开独立表空间，方便进行innodb表的碎片整理

使用MyISAM表引擎的数据库会分别创建三个文件：表结构、表索引、表数据空间。
可以将某个数据库目录直接迁移到其他数据库也可以正常工作。

然而当使用InnoDB的时候，一切都变了。
InnoDB默认会将所有的数据库InnoDB引擎的表数据存储在一个共享空间中：ibdata1文件。
增删数据库的时候，ibdata1文件不会自动收缩，单个数据库的备份也将成为问题。
通常只能将数据使用mysqldump 导出，然后再导入解决这个问题。

在MySQL的配置文件[mysqld]部分，增加innodb_file_per_table参数。
可以修改InnoDB为独立表空间模式，每个数据库的每个表都会生成一个数据空间。

它的优点：
1）每个表都有自已独立的表空间。
2）每个表的数据和索引都会存在自已的表空间中。
3）可以实现单表在不同的数据库中移动。
4）空间可以回收（除drop table操作处，表空不能自已回收）

Drop table操作自动回收表空间，如果对于统计分析或是日值表，删除大量数据后可以通过:alter table TableName engine=innodb;回缩不用的空间。
对于使innodb-plugin的Innodb使用turncate table也会使空间收缩。
对于使用独立表空间的表，不管怎么删除，表空间的碎片不会太严重的影响性能，而且还有机会处理。

它的缺点：
单表增加过大，如超过100个G。

结论：
共享表空间在Insert操作上少有优势。其它都没独立表空间表现好。当启用独立表空间时，请合理调整一下：innodb_open_files。
InnoDB Hot Backup（冷备）的表空间cp不会面对很多无用的copy了。而且利用innodb hot backup及表空间的管理命令可以实。

1）innodb_file_per_table设置.设置为1，表示打开了独立的表空间模式。 如果设置为0，表示关闭独立表空间模式，开启方法如下：
在my.cnf中[mysqld]下设置

innodb_file_per_table=1
2）查看是否开启：

mysql> show variables like "%per_table%";
+-----------------------+-------+
| Variable_name | Value |
+-----------------------+-------+
| innodb_file_per_table | ON |
+-----------------------+-------+
1 row in set (0.00 sec)
3）关闭独享表空间

innodb_file_per_table=0关闭独立的表空间
mysql> show variables like ‘%per_table%’;
-------------------------------------------innodb_file_per_table参数说明------------------------------------------
--------------------------------------------------------------------------------------------------------------------------

备份前，一定要关闭mysql数据库！因为里面会涉及到ibdata1文件备份，不关闭mysql的话，ibdata1文件备份后会损坏，从而导致恢复数据失败！

[root@test-huanqiu snap1]# /etc/init.d/mysql stop
Shutting down MySQL.... SUCCESS! 
[root@test-huanqiu data]# lsof -i:3306 
[root@test-huanqiu data]#
现在备份整个数据库

[root@test-huanqiu snap1]# tar -zvcf /backup/mysql/data/`date +%Y-%m-%d`dbbackup.tar.gz ./
[root@test-huanqiu snap1]# ll /backup/mysql/data/
total 384
-rw-r--r--. 1 root root 392328 Dec 5 22:15 2016-12-05dbbackup.tar.gz
释放快照卷，每次备份之后，应该删除快照，减少IO操作
先卸载，再删除

[root@test-huanqiu ~]# umount /var/snap1/ 
[root@test-huanqiu ~]# df -h                //确认上面的挂载关系已经没了
Filesystem Size Used Avail Use% Mounted on
/dev/mapper/VolGroup00-LogVol00
8.1G 5.8G 1.9G 76% /
tmpfs 1.9G 0 1.9G 0% /dev/shm
/dev/vda1 190M 37M 143M 21% /boot
/dev/mapper/vg0-lv0 2.9G 115M 2.7G 5% /data/mysql/data
[root@test-huanqiu ~]# lvremove /dev/vg0/snap1
Do you really want to remove active logical volume snap1? [y/n]: y
Logical volume "snap1" successfully removed

数据被快照备份后，可以启动数据库

[root@test-huanqiu ~]# /etc/init.d/mysql start
Starting MySQL.. SUCCESS! 
[root@test-huanqiu ~]# lsof -i:3306
COMMAND PID USER FD TYPE DEVICE SIZE/OFF NODE NAME
mysqld 15943 mysql 16u IPv4 93348 0t0 TCP *:mysql (LISTEN)
[root@test-huanqiu ~]#
现在再进行新的数据写入：

mysql> use beijing;
Database changed
mysql> insert into people values("4","liumengnan");
Query OK, 1 row affected (0.02 sec)
 
mysql> insert into people values("5","zhangjuanjuan");
Query OK, 1 row affected (0.00 sec)
 
mysql> select * from people;
+------+---------------+
| id   | name          |
+------+---------------+
|    1 | wangshibo     |
|    2 | guohuihui     |
|    3 | wuxiang       |
|    4 | liumengnan    |
|    5 | zhangjuanjuan |
+------+---------------+
5 rows in set (0.00 sec)
 
mysql> create table heihei(name varchar(20),age varchar(20));
Query OK, 0 rows affected (0.02 sec)
 
mysql> insert into heihei values("jiujiujiu","nan");
Query OK, 1 row affected (0.00 sec)
 
mysql> select * from heihei;
+-----------+------+
| name      | age  |
+-----------+------+
| jiujiujiu | nan  |
+-----------+------+
1 row in set (0.00 sec)
 
mysql> create database shanghai;
Query OK, 1 row affected (0.01 sec)
 
mysql> show databases;
+--------------------+
| Database           |
+--------------------+
| information_schema |
| beijing            |
| mysql              |
| performance_schema |
| shanghai           |
| test               |
+--------------------+
6 rows in set (0.00 sec)
假设一不小心误操作删除beijing和shanghai库

mysql> drop database beijing;
Query OK, 2 rows affected (0.03 sec)
 
mysql> drop database shanghai;
Query OK, 0 rows affected (0.00 sec)
 
mysql> show databases;
+--------------------+
| Database           |
+--------------------+
| information_schema |
| mysql              |
| performance_schema |
| test               |
+--------------------+
4 rows in set (0.00 sec)
莫慌！接下来就说下数据恢复操作~~
三、恢复流程如下：
0）由于涉及到增量数据备份，所以提前将最近一次的binlog日志从mysql数据目录复制到别的路径下
1）在mysql数据库中执行flush logs命令，产生新的binlog日志，让日志信息写入到新的这个binlog日志中
1）关闭数据库，一定要关闭
2）删除数据目录下的文件
3）快照数据拷贝回来，position节点记录回放
4）增量数据就利用mysqlbinlog命令将上面提前拷贝的binlog日志文件导出为sql文件，并剔除其中的drop语句，然后进行恢复。
5）重启数据

先将最新一次的binlog日志备份到别处，用作增量数据备份。
比如mysql-bin.000006是最新一次的binlog日志

[root@test-huanqiu data]# cp mysql-bin.000006 /backup/mysql/data/
产生新的binlog日志，确保日志写入到这个新的binlog日志内，而不再写入到上面备份的binlog日志里。

mysql> flush logs;
 
[root@test-huanqiu data]# ll mysql-bin.000007 
-rw-rw----. 1 mysql mysql 120 Dec 5 23:19 mysql-bin.000007
 
[root@test-huanqiu data]# /etc/init.d/mysql stop
Shutting down MySQL.... SUCCESS! 
[root@test-huanqiu data]# lsof -i:3306
[root@test-huanqiu data]# pwd
/data/mysql/data
[root@test-huanqiu data]# rm -rf ./*
[root@test-huanqiu data]# tar -zvxf /backup/mysql/data/2016-12-05dbbackup.tar.gz ./
 
[root@test-huanqiu data]# /etc/init.d/mysql start
Starting MySQL SUCCESS! 
[root@test-huanqiu data]# cat /backup/mysql/binlog/binlog.pos
File Position Binlog_Do_DB Binlog_Ignore_DB Executed_Gtid_Set
mysql-bin.000004 1434 
[root@test-huanqiu data]# mysqlbinlog --start-position=1434 /data/mysql/data/mysql-bin.000004 | mysql -p123456
登陆数据库查看，发现这只是恢复到快照备份阶段的数据：

mysql> show databases;
+--------------------+
| Database           |
+--------------------+
| information_schema |
| beijing            |
| mysql              |
| performance_schema |
| test               |
+--------------------+
5 rows in set (0.00 sec)
 
mysql> select * from beijing.people;
+------+-----------+
| id   | name      |
+------+-----------+
|    1 | wangshibo |
|    2 | guohuihui |
|    3 | wuxiang   |
+------+-----------+
3 rows in set (0.00 sec)
快照备份之后写入的数据要利用mysqlbinlog命令将上面拷贝的mysql-bin000006文件导出为sql文件，并剔除其中的drop语句，然后进行恢复。

[root@test-huanqiu ~]# cd /backup/mysql/data/
[root@test-huanqiu data]# ll
total 388
-rw-r--r--. 1 root root 392328 Dec 5 22:15 2016-12-05dbbackup.tar.gz
-rw-r-----. 1 root root 1274 Dec 5 23:19 mysql-bin.000006
[root@test-huanqiu data]# mysqlbinlog mysql-bin.000006 >000006bin.sql
剔除其中的drop语句

[root@test-huanqiu data]# vim 000006bin.sql          //手动删除sql语句中的drop语句
然后在mysql中使用source命令恢复数据

mysql> source /backup/mysql/data/000006bin.sql;
再次查看下，发现增量部分的数据也已经恢复回来了

mysql> show databases;
+--------------------+
| Database           |
+--------------------+
| information_schema |
| beijing            |
| mysql              |
| performance_schema |
| shanghai           |
| test               |
+--------------------+
6 rows in set (0.00 sec)
 
mysql> use beijing;
Database changed
mysql> show tables;
+-------------------+
| Tables_in_beijing |
+-------------------+
| heihei            |
| people            |
+-------------------+
2 rows in set (0.00 sec)
 
mysql> select * from people;
+------+---------------+
| id   | name          |
+------+---------------+
|    1 | wangshibo     |
|    2 | guohuihui     |
|    3 | wuxiang       |
|    4 | liumengnan    |
|    5 | zhangjuanjuan |
+------+---------------+
5 rows in set (0.00 sec)
 
mysql> select * from heihei;
+-----------+------+
| name      | age  |
+-----------+------+
| jiujiujiu | nan  |
+-----------+------+
1 row in set (0.00 sec)
-----------------------------------------------------------------------------------------------------------------
思路：
1）全库的快照备份只需要在开始时备份一份即可，这相当于全量备份。
2）后续只需要每天备份一次最新的binlog日志（备份后立即flush logs产生新的binlog日志），这相当于增量备份了。
3）利用快照备份恢复全量数据，利用备份的binlog日志进行增量数据恢复
4）crontab计划任务，每天定时备份最近一次的binlog日志即可。

mysql增量备份的方案，看到有2种：
1.也是使用mysqldump，然后将这一次的dump文件和前一次的文件用比对工具比对，生成patch补丁，将补丁打到上一次的dump文件中，生成一份新的dump文件。（stackoverflow上看到的，感觉思路很奇特，但是不知道这样会不会有问题）
2.增量拷贝 mysql 的数据文件，以及日志文件。典型的方式就是用官方mysqlbackup工具配合 --incremental 参数，但是 mysqlbackup 是mysql Enterprise版才有的工具(收费？)，现在项目安装的mysql版本貌似没有。还有v友中分享的各种增量备份脚本或工具也是基于这种方案？