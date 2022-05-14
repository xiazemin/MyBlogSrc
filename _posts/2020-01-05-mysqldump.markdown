---
title: mysqldump Mysql 大量数据快速导出
layout: post
category: storage
author: 夏泽民
---
mysqldump -u root -p -q -e -t  webgps4 dn_location2 > dn_location2.sql 
mysqldump -u root -p -q -e -t --single-transaction  webgps4 dn_location2 > dn_location2.sql
source dn_location2.sql
以上是导入导出数据的语句，该方法15分钟导出1.6亿条记录，导出的文件中平均7070条记录拼成一个insert语句，通过source进行批量插入，导入1.6亿条数据耗时将近5小时。平均速度：3200W条/h。后来尝试加上--single-transaction参数，结果影响不大。另外，若在导出时增加-w参数，表示对导出数据进行筛选，那么导入导出的速度基本不变，筛选出的数据量越大，时间越慢而已。对于其中的参数这里进行说明：
<!-- more -->
–quick，-q
该选项在导出大表时很有用，它强制 mysqldump 从服务器查询取得记录直接输出而不是取得所有记录后将它们缓存到内存中。


--extended-insert, -e
使用具有多个VALUES列的INSERT语法。这样使导出文件更小，并加速导入时的速度。默认为打开状态，使用--skip-extended-insert取消选项。



--single-transaction

该选项在导出数据之前提交一个BEGIN SQL语句，BEGIN 不会阻塞任何应用程序且能保证导出时数据库的一致性状态。它只适用于多版本存储引擎，仅InnoDB。本选项和--lock-tables 选项是互斥的，因为LOCK TABLES 会使任何挂起的事务隐含提交。要想导出大表的话，应结合使用--quick 选项。在本例子中没有起到加快速度的作用
mysqldump -uroot -p --host=localhost --all-databases --single-transaction



-t 仅导出表数据，不导出表结构


一般的数据备份用 ：mysql路径+bin/mysqldump -u 用户名 -p 数据库名 > 导出的文件名 

数据还原是：到mysql命令行下面，用：source   文件名;的方法。

但是这种方法对大数据量的表进行操作就非常慢。因为他不仅导出了数据还导出了表结构。

在针对大数据量的表时，我们可以用infile和 outfile来操作。

outfile导出数据库数据的用法：

可以看到6百多万数据35秒就搞定了

在infile导入数据的时候，我们还可以做一些优化。我们可以用 

alter table table_name disable keys   关闭普通索引。等数据导入玩，再用：

alter table table_name enable keys    来开启普通索引。这样就不会边导入数据，边整理索引的二叉树儿影响导数据的效率。

如果可以保证 数据的正确性，我们可以将表的唯一索引也关闭，之后再开启,不是每条数据就算是唯一的他都要去检测一遍。命令：

 set unique_checks=0; #关闭唯一校验

 set unique_checks=1;#开启唯一校验

如果是InnoDB存储引擎，我们还可以set auto commit=0;关闭自动提交，来提高效率。InnoDB是按主键的顺序保存的，我们将其主键顺序排列也可以提高效率。



下面我们对myisam引擎的表做个测试，我们先不关索引，导入数据(用了近4分钟)：





然后我们先把索引关闭试试(只用了一分钟多一点，快了不少啊

导出导出工作准备】

（1）导出前关闭日志，避免数据备份过程中频繁记录日志

（2）删除主键，关闭自动增长。在该表中主键其实作用不大，自动增长是需要的（mysql中自动增长的一列一定要为key，所以设置为主键），等待数据转移结束后重新设置回来

（3）删除表中索引。在插入数据时索引的存在会很大程度上影响速度，所以先关闭，转移后重新建立

（4）Mysql系统参数调优，如下：（具体含义后面给出）

innodb_data_file_path           = ibdata1:1G:autoextend  
innodb_file_per_table           = 1  
innodb_thread_concurrency       = 20  
innodb_flush_log_at_trx_commit  = 1  
innodb_log_file_size            = 256M  
innodb_log_files_in_group       = 3  
innodb_max_dirty_pages_pct      = 50  
innodb_lock_wait_timeout        = 120  
key_buffer_size=400M  
innodb_buffer_pool_size=4G  
innodb_additional_mem_pool_size=20M  
innodb_log_buffer_size=20M  
query_cache_size=40M  
read_buffer_size=4M  
read_rnd_buffer_size=8M  
tmp_table_size=16M  
max_allowed_packet = 32M 

操作方法及结果】
（1）create table t2 as select * from t1

[sql] view plain copy
CREATE TABLE dn_location3    
PARTITION BY RANGE (UNIX_TIMESTAMP(UPLOADTIME))  
 (   PARTITION p141109 VALUES LESS THAN (UNIX_TIMESTAMP('2014-11-09 00:00:00')),  
 PARTITION p141110 VALUES LESS THAN (UNIX_TIMESTAMP('2014-11-10 00:00:00')),     
PARTITION p141111 VALUES LESS THAN (UNIX_TIMESTAMP('2014-11-11 00:00:00')),     
PARTITION p141112 VALUES LESS THAN (UNIX_TIMESTAMP('2014-11-12 00:00:00'))   
)   
as select * from dn_location   
where uploadtime > '2014-08-04';  
  
create table t2 as select * from dn_location2;  
as创建出来的t2表（新表）缺少t1表（源表）的索引信息，只有表结构相同，没有索引。
此方法效率较高，在前面的实验环境下，42min内将一张表内4600W的数据转到一张新的表中，在create新表时我添加了分区的操作，因此新表成功创建为分区表，这样一步到位的既转移了数据又创建了分区表。此方法平均速度：6570W条/h ，至于该方法其他需要注意的地方，暂时没有去了解。



（2）使用MySQL的SELECT INTO OUTFILE 、Load data file

LOAD DATA INFILE语句从一个文本文件中以很高的速度读入一个表中。当用户一前一后地使用SELECT ... INTO OUTFILE 和LOAD DATA INFILE 将数据从一个数据库写到一个文件中，然后再从文件中将它读入数据库中时，两个命令的字段和行处理选项必须匹配。否则，LOAD DATA INFILE 将不能正确地解释文件内容。

假设用户使用SELECT ... INTO OUTFILE 以逗号分隔字段的方式将数据写入到一个文件中：

[sql] view plain copy
 在CODE上查看代码片派生到我的代码片
SELECT * INTO OUTFILE 'data.txt' FIELDS TERMINATED BY ',' FROM table2;  
为了将由逗号分隔的文件读回时，正确的语句应该是：

[sql] view plain copy
 在CODE上查看代码片派生到我的代码片
LOAD DATA INFILE 'data.txt' INTO TABLE table2 FIELDS TERMINATED BY ',';  
如果用户试图用下面所示的语句读取文件，它将不会工作，因为命令LOAD DATA INFILE 以定位符区分字段值：
[sql] view plain copy
 在CODE上查看代码片派生到我的代码片
LOAD DATA INFILE 'data.txt' INTO TABLE table2 FIELDS TERMINATED BY '\t';  
下面是我用来导入导出的命令：

[sql] view plain copy
 在CODE上查看代码片派生到我的代码片
select * into outfile 'ddd.txt' fields terminated by ',' from dn_location;  
load data infile 'ddd.txt' into table dn_location2  FIELDS TERMINATED BY ',';  
通过该方法导出的数据，是将各字段（只有数据，不导出表结构）数据存在一个文件中，中间以逗号分隔，因为文件中并不包含数据库名或者表名，因此需要在导入导出的时候些明确。该方法在18分钟内导出1.6亿条记录，46min内导入6472W条记录，平均速度：8442W条/h。mysql官方文档也说明了，该方法比一次性插入一条数据性能快20倍。



【额外测试1】在新的表结构中增加主键，并增加某一列自增，查看主键索引对插入效率的影响

【结论】导出效率没有变化，导入效率35min中导入4600W条记录，平均速度：7886W/h，考虑到测试次数很少，不能直接下结论，但至少明确该操作不会有明显的效率下降。

【测试语句】

[sql] view plain copy
SELECT MOTOR_ID,LAT,LON,UPLOADTIME,RECEIVETIME,STATE_ID,SYS_STATE_ID,SPEED,DIR,A,GPRS,DISTANCE,WEEKDAY,GPSLOCATE  INTO OUTFILE 'import2.txt' FROM dn_location3;  
LOAD DATA INFILE 'import2.txt' INTO TABLE dn_location_withkey(MOTOR_ID,LAT,LON,UPLOADTIME,RECEIVETIME,STATE_ID,SYS_STATE_ID,SPEED,DIR,A,GPRS,DISTANCE,WEEKDAY,GPSLOCATE);  
【额外测试2】在新建的表中对一个varchar类型字段增加索引，再往里导入数据，查看对插入效率的影响。

【结论】导入4600W条记录耗时47min，效率确实有所降低，比仅有主键索引的测试多了12分钟，从这里看插入效率排序： 没有任何索引 > 主键索引  >  主键索引+其他索引。

【额外测试3】在新建表中不加索引导入数据，完全导入后再建索引，查看建立索引时间

【结论】（1）表数据4600W，建立索引时间10min；表数据1.6亿条，建立索引时间41min，由此可见建立索引的时间与表的数据量有直接关系，其他影响因素比较少；（2）从此处看先插入数据再建索引与先建索引再批量插入数据时间上差距不大，前者稍快一些，开发中应根据实际情况选择。



（3）使用mysqldump ，source

[sql] view plain copy
 在CODE上查看代码片派生到我的代码片
mysqldump -u root -p -q -e -t  webgps4 dn_location2 > dn_location2.sql  
mysqldump -u root -p -q -e -t --single-transaction  webgps4 dn_location2 > dn_location2.sql  
source dn_location2.sql  
以上是导入导出数据的语句，该方法15分钟导出1.6亿条记录，导出的文件中平均7070条记录拼成一个insert语句，通过source进行批量插入，导入1.6亿条数据耗时将近5小时。平均速度：3200W条/h。后来尝试加上--single-transaction参数，结果影响不大。另外，若在导出时增加-w参数，表示对导出数据进行筛选，那么导入导出的速度基本不变，筛选出的数据量越大，时间越慢而已。对于其中的参数这里进行说明：
–quick，-q
该选项在导出大表时很有用，它强制 mysqldump 从服务器查询取得记录直接输出而不是取得所有记录后将它们缓存到内存中。


--extended-insert, -e
使用具有多个VALUES列的INSERT语法。这样使导出文件更小，并加速导入时的速度。默认为打开状态，使用--skip-extended-insert取消选项。



--single-transaction

该选项在导出数据之前提交一个BEGIN SQL语句，BEGIN 不会阻塞任何应用程序且能保证导出时数据库的一致性状态。它只适用于多版本存储引擎，仅InnoDB。本选项和--lock-tables 选项是互斥的，因为LOCK TABLES 会使任何挂起的事务隐含提交。要想导出大表的话，应结合使用--quick 选项。在本例子中没有起到加快速度的作用
mysqldump -uroot -p --host=localhost --all-databases --single-transaction



-t 仅导出表数据，不导出表结构



更多的mysqldump 参数说明请参考：http://blog.chinaunix.net/uid-26805356-id-4138986.html      

更多的mysql 参数调优说明参考：http://blog.csdn.net/yang1982_0907/article/details/20123055

                                                     http://blog.csdn.net/nightelve/article/details/17393631

extended-insert对mysqldump及导入性能的影响  http://blog.csdn.net/hw_libo/article/details/39583247

1. 对于Myisam类型的表，可以通过以下方式快速的导入大量的数据。

ALTER  TABLE  tblname  DISABLE  KEYS;

loading  the  data

 ALTER  TABLE  tblname  ENABLE  KEYS;

这两个命令用来打开或者关闭Myisam表非唯一索引的更新。在导入大量的数据到一 个非空的Myisam表时，通过设置这两个命令，可以提高导入的效率。对于导入大量 数据到一个空的Myisam表，默认就是先导入数据然后才创建索引的，所以不用进行 设置。

2. 而对于Innodb类型的表，这种方式并不能提高导入数据的效率。对于Innodb类型 的表，我们有以下几种方式可以提高导入的效率：

a. 因为Innodb类型的表是按照主键的顺序保存的，所以将导入的数据按照主键的顺 序排列，可以有效的提高导入数据的效率。如果Innodb表没有主键，那么系统会默认创建一个内部列作为主键，所以如果可以给表创建一个主键，将可以利用这个优势提高 导入数据的效率。

b. 在导入数据前执行SET  UNIQUE_CHECKS=0，关闭唯一性校验，在导入结束后执行SET  UNIQUE_CHECKS=1，恢复唯一性校验，可以提高导入的效率。
c. 如果应用使用自动提交的方式，建议在导入前执行SET  AUTOCOMMIT=0，关闭自动 提交，导入结束后再执行
