---
title: auto_increment
layout: post
category: mysql
author: 夏泽民
---
auto_increment的基本特性
MySQL的中AUTO_INCREMENT类型的属性用于为一个表中记录自动生成ID功能，可在一定程度上代替Oracle，PostgreSQL等数据库中的sequence。

在数据库应用，我们经常要用到唯一编号，以标识记录。在MySQL中可通过数据列的AUTO_INCREMENT属性来自动生成。

可在建表时可用“AUTO_INCREMENT=n”选项来指定一个自增的初始值。
可用alter table table_name AUTO_INCREMENT=n命令来重设自增的起始值。
当插入记录时，如果为AUTO_INCREMENT数据列明确指定了一个数值，则会出现两种情况，
情况一，如果插入的值与已有的编号重复，则会出现出错信息，因为AUTO_INCREMENT数据列的值必须是唯一的；
情况二，如果插入的值大于已编号的值，则会把该插入到数据列中，并使在下一个编号将从这个新值开始递增。也就是说，可以跳过一些编号。
如果自增序列的最大值被删除了，则在插入新记录时，该值被重用。
如果用UPDATE命令更新自增列，如果列值与已有的值重复，则会出错。如果大于已有值，则下一个编号从该值开始递增。

在使用AUTO_INCREMENT时，应注意以下几点：
AUTO_INCREMENT是数据列的一种属性，只适用于整数类型数据列。
设置AUTO_INCREMENT属性的数据列应该是一个正数序列，所以应该把该数据列声明为UNSIGNED，这样序列的编号个可增加一倍。
AUTO_INCREMENT数据列必须有唯一索引，以避免序号重复(即是主键或者主键的一部分)。
AUTO_INCREMENT数据列必须具备NOT NULL属性。
AUTO_INCREMENT数据列序号的最大值受该列的数据类型约束，如TINYINT数据列的最大编号是127,如加上UNSIGNED，则最大为255。一旦达到上限，AUTO_INCREMENT就会失效。
当进行全表删除时，MySQL AUTO_INCREMENT会从1重新开始编号。全表删除的意思是发出以下两条语句时：
[php] view plain copy

delete from table_name;  
或者  
truncate table table_name   
delete from table_name;
或者
truncate table table_name 
1
2
这是因为进行全表操作时，MySQL(和PHP搭配之最佳组合)实际是做了这样的优化操作：先把数据表里的所有数据和索引删除，然后重建数据表。
如果想删除所有的数据行又想保留序列编号信息，可这样用一个带where的delete命令以抑制MySQL(和PHP搭配之最佳组合)的优化：
[php] view plain copy

delete from table_name where 1;   
delete from table_name where 1; 

可用last_insert_id（）获取刚刚自增过的值。
关于mysql auto_increment所带来的锁表操作
在mysql5.1.22之前，mysql的“INSERT-like”语句（包INSERT, INSERT…SELECT, REPLACE,REPLACE…SELECT, and LOAD DATA)会在执行整个语句的过程中使用一个AUTO-INC锁将表锁住，直到整个语句结束（而不是事务结束）。
因此在使用INSERT…SELECT、INSERT…values(…),values(…)时，LOAD DATA等耗费时间较长的操作时，会将整个表锁住，而阻塞其他的“INSERT-like”、Update等语句，推荐使用程序将这些语句分成多条语句，一一插入，减少单一时间的锁表时间。
mysql5.1.22之后mysql进行了改进，引入了参数 innodb_autoinc_lock_mode，通过这个参数控制mysql的锁表逻辑。
在介绍这个之前先引入几个术语，方便说明 innodb_autoinc_lock_mode。
1.“INSERT-like”：
INSERT, INSERT … SELECT, REPLACE, REPLACE … SELECT, and LOAD DATA, INSERT … VALUES(),VALUES()
2.“Simple inserts”：
就是通过分析insert语句可以确定插入数量的insert语句, INSERT, INSERT … VALUES(),VALUES()
3.“Bulk inserts”：
就是通过分析insert语句不能确定插入数量的insert语句, INSERT … SELECT, REPLACE … SELECT, LOAD DATA
4.“Mixed-mode inserts”：
不确定是否需要分配auto_increment id，一般是下面两种情况
INSERT INTO t1 (c1,c2) VALUES (1,’a’), (NULL,’b’), (5,’c’), (NULL,’d’);
INSERT … ON DUPLICATE KEY UPDATE

一、innodb_autoinc_lock_mode = 0 (“traditional” lock mod，传统模式)。
这种方式就和mysql5.1.22以前一样，为了向后兼容而保留了这种模式，如同前面介绍的一样，这种方式的特点就是“表级锁定”，并发性较差。
二、innodb_autoinc_lock_mode = 1 (“consecutive” lock mode，连续模式)。
这种方式是新版本中的默认方式，推荐使用，并发性相对较高，特点是“consecutive”，即保证同一条insert语句中新插入的auto_increment id都是连续的。
这种模式下：
“Simple inserts”：直接通过分析语句，获得要插入的数量，然后一次性分配足够的auto_increment id，只会将整个分配的过程锁住。
“Bulk inserts”：因为不能确定插入的数量，因此使用和以前的模式相同的表级锁定。
“Mixed-mode inserts”：直接分析语句，获得最坏情况下需要插入的数量，然后一次性分配足够的auto_increment id，只会将整个分配的过程锁住。
需要注意的是，这种方式下，会分配过多的id，而导致“浪费”。
比如INSERT INTO t1 (c1,c2) VALUES (1,’a’), (NULL,’b’), (5,’c’), (NULL,’d’);会一次性的分配5个id，而不管用户是否指定了部分id；
INSERT … ON DUPLICATE KEY UPDATE一次性分配，而不管将来插入过程中是否会因为duplicate key而仅仅执行update操作。
注意：当master mysql版本<5.1.22，slave mysql版本>=5.1.22时，slave需要将innodb_autoinc_lock_mode设置为0，因为默认的innodb_autoinc_lock_mode为1，对于INSERT … ON DUPLICATE KEY UPDATE和INSERT INTO t1 (c1,c2) VALUES (1,’a’), (NULL,’b’), (5,’c’), (NULL,’d’);的执行结果不同，现实环境一般会使用INSERT
… ON DUPLICATE KEY UPDATE。
三、innodb_autoinc_lock_mode = 2 (“interleaved” lock mode，交叉模式)。
这种模式是来一个分配一个，而不会锁表，只会锁住分配id的过程，和innodb_autoinc_lock_mode = 1的区别在于，不会预分配多个，这种方式并发性最高。
但是在replication中当binlog_format为statement-based时（简称SBR statement-based replication）存在问题，因为是来一个分配一个，这样当并发执行时，“Bulk inserts”在分配的时会同时向其他的INSERT分配，会出现主从不一致（从库执行结果和主库执行结果不一样），因为binlog只会记录开始的insert id。
测试SBR，执行begin;insert values(),();insert values(),();commit;会在binlog中每条insert values(),();前增加SET INSERT_ID=18/!/;。
但是row-based replication RBR时不会存在问题。
另外RBR的主要缺点是日志数量在包括语句中包含大量的update delete（update多条语句，delete多条语句）时，日志会比SBR大很多；假如实际语句中这样语句不是很多的时候（现实中存在很多这样的情况），推荐使用RBR配合innodb_autoinc_lock_mode，不过话说回来，现实生产中“Bulk inserts”本来就很少，因此innodb_autoinc_lock_mode = 1应该是够用了
<!-- more -->
https://blog.csdn.net/xiong9999/article/details/53688277
第 4 篇提过自增主键，让主键索引尽量递增顺序插入，避免页分裂，索引更紧凑。

不能保证连续递增。什么情况出现 “空洞”？主键冲突、回滚、批量申请。

id 自增主键、c 唯一索引


一、自增值保存在哪儿？
insert into t values(null, 1, 1); show create table

图 1 自动生成的 AUTO_INCREMENT 值
AUTO_INCREMENT=2，下一次插入，生成 id=2。

表的结构定义存放在后缀名为.frm 的文件中，不会保存自增值。

MyISAM 保存数据文件；InnoDB 自增值保存内存里（MySQL 8.0记录redo log 中），重启恢复为重启前值

MySQL 5.7 前，没持久化。最大值 max(id)=10，AUTO_INCREMENT=11。删除 id=10 的行，AUTO_INCREMENT 还是 11。重启AUTO_INCREMENT 0。重启修改AUTO_INCREMENT。

二、自增值修改机制
id 被定义为 AUTO_INCREMENT，

1. 插入时 id 字段指定 0、null 或未指定值，AUTO_INCREMENT 填自增字段；

2.  id 指定具体值，用指定值。

插入 X，自增值 Y。

1.   X（2）<Y（3），自增值不变；

2.  如果 （2）X≥Y（1），前自增值改为新自增值X

新自增值生成算法是：自增的初始值auto_increment_offset(默认1)开始， 步长auto_increment_increment (默认1)为，持续叠加，第一个大于 X 值，作为新的自增值。

ps：用的不全是默认值。如双 M 主备结构里双写，auto_increment_increment=2，让自增 id 都是奇数，另一都是偶数，避免主键冲突。

插入的值 >= 当前自增值，新自增值就是“准备插入的值 +1”；否则不变。

两个参数都设置为 1 的时候，自增主键 id 却不能保证是连续的，这是什么原因呢？

三、自增值的修改时机
3.1主键冲突
已经有了 (1,1,1)   insert into t  values(null, 1, 1);

1.  写入一行(0,1,1);

2.  没有指定自增 id 的值，t 当前的自增值 2；

3.  改成(2,1,1)；自增值改成 3

4.  已经存在 c=1 的记录，所以报 Duplicate key error，语句返回。

图 2 insert(null, 1,1) 唯一键冲突
没有插入成功，自增值不再改回去。不连续。

图 3 一个自增主键 id 不连续的复现步骤
3.2回滚也会产生类似现象

自增值不能回退：提升性能

两个并行事务，加锁顺序申请。

1.   A 申请到了 id=2， B  id=3， t 的自增值是 4，

2.   B 提交， A 出现唯一键冲突。

3.  如允许A 自增 id 回退， t 改回 2，问题：id=3 再申请到 id=3（已有）“主键冲突”

解决主键冲突方法（ 都导致性能问题，放弃）：
1. 申请前判断，存在跳过。成本高。因为，本来申请 id 是一个很快的操作，现在还要再去主键索引树上判断 id 是否存在。

2. 锁范围扩大，事务完提交，下一个再申请自增 id。粒度太大，并发能力下降。

四、自增锁的优化
4.1 innodb_autoinc_lock_mode控制锁粒度
0：5.0 策略，语句结束释放；

1 (默认)：普通 insert ，自增锁申请后释放；

                insert …select 批量插入，语句结束释放；

2 ：申请后释放锁

ps：5.0 版本，自增锁，语句级别。申请表自增锁，结束释放，影响并发度。

4.1 为什么默认、insert … select用语句级锁？默认值不是 2？
数据的一致性

图 4 批量插入数据的自增锁
t1 插入4 行，创建相同结构表 t2，同时向 t2 插入。

如果 B 申请自增值后马上释放自增锁，情况：

B 先插入了两个记录，(1,1,1)、(2,2,2)；

A 来申请自增 id 得到 id=3，插入了（3,5,5)；

B 插入两条记录 (4,3,3)、 (5,4,4)。

B 本身就没要求 t2 跟A 相同。如果binlog_format=statement，binlog 里id 连续。数据不一致。

问题原因：B 的 insert 语句，生成id 不连续。statement 格式的 binlog 串行执行，执行不出来。

解决两种思路：
1.  批量插入数据语句，固定生成连续 id 值。语句结束释放

2.  binlog如实记录进来，备库执行，不依赖于自增主键生成。innodb_autoinc_lock_mode = 2，binlog_format = row。生产上，尤其insert …select 批量插入数据时，提升并发性，不会数据一致性。

批量插入数据包含： insert …select、replace … select 和 load data 

普通 insert 多个 value 情况，innodb_autoinc_lock_mode 设置 1，精确计算出要多少个 id 的，一次性申请，释放。

4.2 批量申请自增 id 的策略（ 不连续第三种原因 ）：
1.  第一次申请自增 id，会分配 1 个；

2.  第二次申请自增 id，会分配 2 个；

3. 第三次申请自增 id，会分配 4 个；依此类推，上一次的两倍


实际上t2 中插入 4 行，分三次，1，第二次id=2 和 id=3， 第三次id=4 到 id=7。

 id=5 到 id=7 浪费掉

再执行 insert into t2  values(null, 5,5)，实际上插入的数据就是（8,5,5)。

小结
自增值存储。

MyISAM 里，被写数据文件上。 InnoDB 中，记录内存的。重启前后不变。

自增值改变时机，回滚不能回收自增 id。

innodb_autoinc_lock_mode，控制自增值申请锁范围。并发性能考虑，设置为 2，binlog_format =row。

思考题
最后例子，执行 insert into t2(c,d) select c,d from t; 隔离级别是可重复读（repeatable

read），binlog_format=statement。所有记录和间隙加锁。为什么这么做？

如果 insert …select 有其他线程操作原表，不会导致逻辑错误。如不加锁，就是快照读(执行期间，一致性视图是不会修改)。

不对t表所有记录和间隙加锁，，可重复读，其他提交t2看不到。但binlog=statement，备库或基于binlog恢复临时库t2看到，不一致。

评论1
自增id和写binlog是有先后顺序的。binlog=statement，A获取id=1，B获取id=2，B提交，写binlog，再A写binlog。

如果binlog重放，不会出现不一致，B的id为1，A的id为2的情况

因为binlog记录自增值语句前，前面多一句，指定“自增ID值多少”，对应主库自增值

http://events.jianshu.io/p/5b75339bf314
https://cdn.modb.pro/db/52819
https://www.jianshu.com/p/c570ed6c9811
https://www.jianshu.com/p/c570ed6c9811
