---
title: MySQL的表类型的（存储引擎）
layout: post
category: storage
author: 夏泽民
---
查看数据库支持的存储引擎： \G或者分号表示命令结束
MySQL>show engines \G
or show variables like 'have%';
mysql -h localhost -u root -p123
show databases;
use database;
show tables;


创建指定存储引擎的表：
create table tableA(
i bigint(20) not null auto_increment,
primary key(i)
)engine=MyISAM default charset=gbk;


改变表的存储引擎：
alter
table tableA engine=innodb;

> show engines \G
ERROR 2006 (HY000): MySQL server has gone away
No connection. Trying to reconnect...
Connection id:    1
Current database: carpool

*************************** 1. row ***************************
      Engine: InnoDB
     Support: DEFAULT
     Comment: Supports transactions, row-level locking, and foreign keys
Transactions: YES
          XA: YES
  Savepoints: YES
*************************** 2. row ***************************
      Engine: MRG_MYISAM
     Support: YES
     Comment: Collection of identical MyISAM tables
Transactions: NO
          XA: NO
  Savepoints: NO
*************************** 3. row ***************************
      Engine: MEMORY
     Support: YES
     Comment: Hash based, stored in memory, useful for temporary tables
Transactions: NO
          XA: NO
  Savepoints: NO
*************************** 4. row ***************************
      Engine: BLACKHOLE
     Support: YES
     Comment: /dev/null storage engine (anything you write to it disappears)
Transactions: NO
          XA: NO
  Savepoints: NO
*************************** 5. row ***************************
      Engine: MyISAM
     Support: YES
     Comment: MyISAM storage engine
Transactions: NO
          XA: NO
  Savepoints: NO
*************************** 6. row ***************************
      Engine: CSV
     Support: YES
     Comment: CSV storage engine
Transactions: NO
          XA: NO
  Savepoints: NO
*************************** 7. row ***************************
      Engine: ARCHIVE
     Support: YES
     Comment: Archive storage engine
Transactions: NO
          XA: NO
  Savepoints: NO
*************************** 8. row ***************************
      Engine: PERFORMANCE_SCHEMA
     Support: YES
     Comment: Performance Schema
Transactions: NO
          XA: NO
  Savepoints: NO
*************************** 9. row ***************************
      Engine: FEDERATED
     Support: NO
     Comment: Federated MySQL storage engine
Transactions: NULL
          XA: NULL
  Savepoints: NULL
  
  各种存储引擎的特性：
### 一
MyISAM：不支持事务和外键；
访问速度快，对事务完整性没有要求的或者以select insert为主的应用适合；
每个MyISAM 在磁盘上存储三个文件，文件名都和表名相同，但扩展名分别是：
.frm 存储表定义
myd 存储数据
myi 存储索引

数据文件和存储文件可以放在不同目录，平均分布io,获得更快的速度；
表容易被损坏，check table, repair table;
MyISAM的表还支持三种不同的存储格式：
静态表（固定长度）：存储迅速，容易缓存，出现故障容易恢复；占用空间比动态表多；空格补充，返回时去掉空格；
create databases test_mysql;
create table Myisam_char (name char(10))engine=myisam;
insert into Myisam_char 
values('abcde'),('abcde '),(' abcde'),(' abcde '); //5 5 7 7

select name,length(name) 
from Myisam_char;
插入进去前面的空格保留，后面的空格都被去掉了。

动态表包含变长字段，记录不是固定长度的，这样存储：占用的空间少，频繁地删除记录会产生碎片，定期执行优化改善性能（optimize table 或者 myisamchk-r ），出故障恢复困难；

压缩表：由myisampack工具创建，占据非常小的磁盘空间，每条记录被单独压缩，只有非常小访问开支；

InnoDB:
具有提交，回滚，崩溃恢复能力的事务安全。相遇于MyISAM引擎，写处理速度略差，占用更多存储空间保留数据和索引。

2.1自动增长列：插入的是0或者null，则实际插入的是增长后的值；
create table autoincre_demo ( i smallint not null auto_increment, name varchar(10), primary key(i) ) engine=innodb;

insert into autoincre_demo values(1,'1'),(0,'2'),(3,'3');
select * from autoincre_demo;

select last_insert_id(); //查询当前线程最后插入记录使用的值；如果是多条，返回第一条；
改变初始值，初始值默认为1，存在内存中，重启数据库需要重新设定；
alert table tableA auto_increment=n;

InnoDB表：自动增长列必须是索引，如果是组合索引，必须是组合索引的第一列；

MyISAM表：自动增长列可以是组合索引的其他列；
create table autoincre_demo ( d1 samllint not null auto_increment,
d2 smalliint not null, 
name varchar(10), index(d2,d1) )engine=MyISAM；
insert into autoincre_demo(d2,name) values(2,'2'),(3,'3'),(4,'4'),(2,'2'),(3,'3'),(4,'4');
select * from autoincre_demo;
2.2 外键约束：
mysql支持外键的存储引擎只有InnoDB；
  
key
 是数据库的物理结构，它包含两层意义，一是约束（偏重于约束和规范数据库的结构完整性），二是索引（辅助查询用的）
index是数据库的物理结构，它只是辅助查询的，它创建时会在另外的表空间（mysql中的innodb表空间）以一个类似目录的结构存储。索引要分类的话，分为前缀索引、全文本索引等；因此，索引只是索引，它不会去约束索引的字段的行为。

各种引擎对比：
MyISAM：读和插入为主，很少更新和删除；并对事务的完整性，并发性要求不高 ；
InnoDB：事务处理提交和回滚，支持外键。插入查询，更新删除。类似的计费系统财务系统；
MEMORY：快速定位记录；
MERGE：突破对单个MyISAM表大小的限制，使用VLＤＢ；

MyISAM
特性
不支持事务：MyISAM存储引擎不支持事务，所以对事务有要求的业务场景不能使用
表级锁定：其锁定机制是表级索引，这虽然可以让锁定的实现成本很小但是也同时大大降低了其并发性能
读写互相阻塞：不仅会在写入的时候阻塞读取，MyISAM还会在读取的时候阻塞写入，但读本身并不会阻塞另外的读
只会缓存索引：MyISAM可以通过key_buffer缓存以大大提高访问性能减少磁盘IO，但是这个缓存区只会缓存索引，而不会缓存数据
适用场景
不需要事务支持（不支持）
并发相对较低（锁定机制问题）
数据修改相对较少（阻塞问题）
以读为主
数据一致性要求不是非常高
最佳实践
尽量索引（缓存机制）
调整读写优先级，根据实际需求确保重要操作更优先
启用延迟插入改善大批量写入性能
尽量顺序操作让insert数据都写入到尾部，减少阻塞
分解大的操作，降低单个操作的阻塞时间
降低并发数，某些高并发场景通过应用来进行排队机制
对于相对静态的数据，充分利用Query Cache可以极大的提高访问效率
MyISAM的Count只有在全表扫描的时候特别高效，带有其他条件的count都需要进行实际的数据访问
InnoDB
特性
具有较好的事务支持：支持4个事务隔离级别，支持多版本读
行级锁定：通过索引实现，全表扫描仍然会是表锁，注意间隙锁的影响
读写阻塞与事务隔离级别相关
具有非常高效的缓存特性：能缓存索引，也能缓存数据
整个表和主键以Cluster方式存储，组成一颗平衡树
所有Secondary Index都会保存主键信息
适用场景
需要事务支持（具有较好的事务特性）
行级锁定对高并发有很好的适应能力，但需要确保查询是通过索引完成
数据更新较为频繁的场景
数据一致性要求较高
硬件设备内存较大，可以利用InnoDB较好的缓存能力来提高内存利用率，尽可能减少磁盘 IO
最佳实践
主键尽可能小，避免给Secondary index带来过大的空间负担
避免全表扫描，因为会使用表锁
尽可能缓存所有的索引和数据，提高响应速度
在大批量小插入的时候，尽量自己控制事务而不要使用autocommit自动提交
合理设置innodb_flush_log_at_trx_commit参数值，不要过度追求安全性
避免主键更新，因为这会带来大量的数据移动
NDBCluster
特性
分布式：分布式存储引擎，可以由多个NDBCluster存储引擎组成集群分别存放整体数据的一部分
支持事务：和Innodb一样，支持事务
可与mysqld不在一台主机：可以和mysqld分开存在于独立的主机上，然后通过网络和mysqld通信交互
内存需求量巨大：新版本索引以及被索引的数据必须存放在内存中，老版本所有数据和索引必须存在与内存中
适用场景
具有非常高的并发需求
对单个请求的响应并不是非常的critical
查询简单，过滤条件较为固定，每次请求数据量较少，又不希望自己进行水平Sharding
最佳实践
尽可能让查询简单，避免数据的跨节点传输
尽可能满足SQL节点的计算性能，大一点的集群SQL节点会明显多余Data节点
在各节点之间尽可能使用万兆网络环境互联，以减少数据在网络层传输过程中的延时
注：以上三个存储引擎是目前相对主流的存储引擎，还有其他类似如：Memory，Merge，CSV，Archive等存储引擎的使用场景都相对较少


<!-- more -->
<img src="{{site.url}}{{site.baseurl}}/img/jupyterSlider.png"/>
