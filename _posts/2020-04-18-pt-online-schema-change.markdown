---
title: pt-online-schema-change
layout: post
category: storage
author: 夏泽民
---
在线数据库的维护中，总会涉及到研发修改表结构的情况，修改一些小表影响很小，而修改大表时，往往影响业务的正常运转，如表数据量超过500W，1000W，甚至过亿时

在线修改大表的可能影响
在线修改大表的表结构执行时间往往不可预估，一般时间较长
由于修改表结构是表级锁，因此在修改表结构时，影响表写入操作
如果长时间的修改表结构，中途修改失败，由于修改表结构是一个事务，因此失败后会还原表结构，在这个过程中表都是锁着不可写入
修改大表结构容易导致数据库CPU、IO等性能消耗，使MySQL服务器性能降低
在线修改大表结构容易导致主从延时，从而影响业务读取
pt-online-schema-change介绍
pt-online-schema-change是percona公司开发的一个工具，在percona-toolkit包里面可以找到这个功能，它可以在线修改表结构

原理:

首先它会新建一张一模一样的表，表名一般是_new后缀
然后在这个新表执行更改字段操作
然后在原表上加三个触发器，DELETE/UPDATE/INSERT，将原表中要执行的语句也在新表中执行
最后将原表的数据拷贝到新表中，然后替换掉原表
使用pt-online-schema-change执行SQL的日志
SQL语句:
ALTER TABLE tmp_task_user ADD support tinyint(1) unsigned NOT NULL DEFAULT '1';

sh pt.sh tmp_task_user "ADD COLUMN support tinyint(1) unsigned NOT NULL DEFAULT '1'"

tmp_task_user
ADD COLUMN support tinyint(1) unsigned NOT NULL DEFAULT '1'
No slaves found.  See --recursion-method if host h=127.0.0.1,P=3306 has slaves.
Not checking slave lag because no slaves were found and --check-slave-lag was not specified.
Operation, tries, wait:
  analyze_table, 10, 1
  copy_rows, 10, 0.25
  create_triggers, 10, 1
  drop_triggers, 10, 1
  swap_tables, 10, 1
  update_foreign_keys, 10, 1
Altering `test_db`.`tmp_task_user`...
Creating new table...
Created new table test_db._tmp_task_user_new OK.
Altering new table...
Altered `test_db`.`_tmp_task_user_new` OK.
2018-05-14T18:14:21 Creating triggers...
2018-05-14T18:14:21 Created triggers OK.
2018-05-14T18:14:21 Copying approximately 6 rows...
2018-05-14T18:14:21 Copied rows OK.
2018-05-14T18:14:21 Analyzing new table...
2018-05-14T18:14:21 Swapping tables...
2018-05-14T18:14:21 Swapped original and new tables OK.
2018-05-14T18:14:21 Dropping old table...
2018-05-14T18:14:21 Dropped old table `test_db`.`_tmp_task_user_old` OK.
2018-05-14T18:14:21 Dropping triggers...
2018-05-14T18:14:21 Dropped triggers OK.
Successfully altered `test_db`.`tmp_task_user`.
好处:

降低主从延时的风险
可以限速、限资源，避免操作时MySQL负载过高
建议:

在业务低峰期做，将影响降到最低
<!-- more -->
pt-online-schema-change安装
1.去官网下载对应的版本，官网下载地址:https://www.percona.com/downl...

2.下载解压之后就可以看到pt-online-schema-change

clipboard.png

3.该工具需要一些依赖包，直接执行不成功时一般会有提示，这里可以提前yum安装

yum install perl-DBI
yum install perl-DBD-MySQL
yum install perl-Time-HiRes
yum install perl-IO-Socket-SSL
pt-online-schema-change使用
1.参数
./bin/pt-online-schema-change --help 可以查看参数的使用，我们只是要修改个表结构，只需要知道几个简单的参数就可以了

--user=        连接mysql的用户名
--password=    连接mysql的密码
--host=        连接mysql的地址
P=3306         连接mysql的端口号
D=             连接mysql的库名
t=             连接mysql的表名
--alter        修改表结构的语句
--execute      执行修改表结构
--charset=utf8 使用utf8编码，避免中文乱码
--no-version-check  不检查版本，在阿里云服务器中一般加入此参数，否则会报错
2.为避免每次都要输入一堆参数，写个脚本复用一下，pt.sh

#!/bin/bash
table=$1
alter_conment=$2

cnn_host='127.0.0.1'
cnn_user='user'
cnn_pwd='password'
cnn_db='database_name'

echo "$table"
echo "$alter_conment"
/root/percona-toolkit-2.2.19/bin/pt-online-schema-change --charset=utf8 --no-version-check --user=${cnn_user} --password=${cnn_pwd} --host=${cnn_host}  P=3306,D=${cnn_db},t=$table --alter 
"${alter_conment}" --execute
3.添加表字段
如添加表字段SQL语句为:
ALTER TABLE tb_test ADD COLUMN column1 tinyint(4) DEFAULT NULL;
那么使用pt-online-schema-change则可以这样写
sh pt.sh tb_test "ADD COLUMN column1 tinyint(4) DEFAULT NULL"

4.修改表字段
SQL语句：
ALTER TABLE tb_test MODIFY COLUMN num int(11) unsigned NOT NULL DEFAULT '0';

pt-online-schema-change工具:
sh pt.sh tb_test "MODIFY COLUMN num int(11) unsigned NOT NULL DEFAULT '0'"

5.修改表字段名
SQL语句:
ALTER TABLE tb_test CHANGE COLUMN age adress varchar(30);

pt-online-schema-change工具:
sh pt.sh tb_test "CHANGE COLUMN age address varchar(30)"

6.添加索引
SQL语句:
ALTER TABLE tb_test ADD INDEX idx_address(address);
pt-online-schema-change工具:
sh pt.sh tb_test "ADD INDEX idx_address(address)"

其他
pt-online-schema-change工具还有很多其他的参数，可以有很多限制，比如限制CPU、线程数量、从库状态等等，不过我做过一个超过6000W表的结构修改，发现几乎不影响性能，很稳定很流畅的就修改了表结构，所以，对以上常规参数的使用基本能满足业务
一定要在业务低峰期做，这样才能确保万无一失

MySQL ddl 的问题现状

在 运维mysql数据库时，我们总会对数据表进行ddl 变更，修改添加字段或者索引，对于mysql 而已，ddl 显然是一个令所有MySQL dba 诟病的一个功能，因为在MySQL中在对表进行ddl时，会锁表，当表比较小比如小于1w上时，对前端影响较小，当时遇到千万级别的表 就会影响前端应用对表的写操作。

 

工作原理：

1、如果存在外键，根据alter-foreign-keys-method参数的值，检测外键相关的表，做相应设置的处理。没有使用 --alter-foreign-keys-method=rebuild_constraints 指定特定的值，该工具不予执行
2、创建一个和源表表结构一样的临时表(_tablename_new)，执行alter修改临时表表结构。
3、在原表上创建3个于inser delete update对应的触发器. （用于copy 数据的过程中，在原表的更新操作 更新到新表.）

4、从原表拷贝数据到临时表，拷贝过程中在原表进行的写操作都会更新到新建的临时表
5、修改外键相关的子表，根据修改后的数据，修改外键关联的子表。
6、rename源数据表为old表，把新表rename为源表名，并将old表删除。
7、删除触发器。

 

执行条件：

1．操作的表必须有主键或则唯一索引，否则报如下错误。

Cannot chunk the original table `test`.`t_driver_allowance_test`: There is no good index and the table is oversized. at /usr/local/bin/pt-online-schema-change line 5486.

2 . 该表不能定义触发器，否则报如下错误。

The table `taotao`.`tttt` has triggers.  This tool needs to create its own triggers, so the table cannot already have triggers.

 

用法介绍：

pt-online-schema-change --help  查看参数选项

这里有几个参数需要介绍一下：

--dry-run  这个参数不建立触发器，不拷贝数据，也不会替换原表。只是创建和更改新表。

--execute  这个参数的作用和前面工作原理的介绍的一样，会建立触发器，来保证最新变更的数据会影响至新表。注意：如果不加这个参数，这个工具会在执行一些检查后退出。

--critical-load  每次chunk操作前后，会根据show global status统计指定的状态量的变化，默认是统计Thread_running。目的是为了安全，防止原始表上的触发器引起负载过高。这也是为了防止在线DDL对线上的影响。超过设置的阀值，就会终止操作，在线DDL就会中断。提示的异常如上报错信息。

--max-load 选项定义一个阀值，在每次chunk操作后，查看show global status状态值是否高于指定的阀值。该参数接受一个mysql status状态变量以及一个阀值，如果没有给定阀值，则定义一个阀值为为高于当前值的20%。注意这个参数不会像--critical-load终止操作，而只是暂停操作。当status值低于阀值时，则继续往下操作。是暂停还是终止操作这是--max-load和--critical-load的差别。

--charset=utf8连接到MySQL后运行SET NAMES UTF8

--check-replication-filters 检查复制中是否设置了过滤条件，如果设置了，程序将退出

--nocheck-replication-filters 不检查复制中是否设置了过滤条件

--set-vars 设置mysql的变量值

--check-slave-lag 检查主从延迟

 

例子：

添加字段
pt-online-schema-change -u root  -p 123456  --alter='add column vid int ' --execute D=taotao,t=tttt --max-load=Threads_connected:650 --critical-load=Threads_running=550 --charset=utf8  --nocheck-replication-filters --set-vars innodb_lock_wait_timeout=30000

删除字段
pt-online-schema-change -u root  -p 123456  --alter='drop column vid  '  --execute D=taotao,t=tttt --max-load=Threads_connected:650 --critical-load=Threads_running=550 --charset=utf8  --nocheck-replication-filters --set-vars innodb_lock_wait_timeout=30000

修改字段
pt-online-schema-change -u root  -p 123456  --alter='modify  column sid bigint(25) '  --execute D=taotao,t=tttt --max-load=Threads_connected:650 --critical-load=Threads_running=550 --charset=utf8  --nocheck-replication-filters --set-vars innodb_lock_wait_timeout=30000

添加索引
pt-online-schema-change -u root  -p 123456  --alter=' add key indx_sid(sid) '  --execute D=taotao,t=tttt --max-load=Threads_connected:650 --critical-load=Threads_running=550 --charset=utf8  --nocheck-replication-filters --set-vars innodb_lock_wait_timeout=30000

删除索引
pt-online-schema-change -u root  -p 123456  --alter=' drop  key indx_sid '  --execute D=taotao,t=tttt --max-load=Threads_connected:650 --critical-load=Threads_running=550 --charset=utf8  --nocheck-replication-filters --set-vars innodb_lock_wait_timeout=30000

 

 

pt-online-schema-change --user=dba_user --password=msds007 -S /app/mysqldata/3306/mysql.sock  --charset=utf8 --no-check-replication-filters --alter "modify HomeworkID bigint(20) AUTO_INCREMENT" --no-drop-old-table D=test,t=wx_edu_homework --alter-foreign-keys-method=rebuild_constraints --print --execute

 

考虑从库延迟情况 ，意味这要注意这几个选项的设置

--max-lag
--check-interval
--recursion-method
--check-slave-lag
 

pt-online-schema-change --user=dba_user --password=msds007 -S /app/mysqldata/3306/mysql.sock  --charset=utf8 --no-check-replication-filters --alter "modify HomeworkID bigint(20) AUTO_INCREMENT" --no-drop-old-table D=test,t=wx_edu_homework --alter-foreign-keys-method=rebuild_constraints --print --execute --max-lag=1s --check-interval=10s --check-slave-lag=h=192.168.1.121,u=root,p=msds007,P=3306

 

pt-online-schema-change --user=dba_user --password=msds007 -S /app/mysqldata/3306/mysql.sock  --charset=utf8 --no-check-replication-filters --alter "modify HomeworkID bigint(20) AUTO_INCREMENT" --no-drop-old-table D=test,t=wx_edu_homework --alter-foreign-keys-method=rebuild_constraints --print --execute --max-lag=1s --check-interval=10s --recursion-method=processlist

 

流程：

1.判断各种参数

2.根据原表"t",创建一个名称为"_t_new"的新表

3.执行ALTER TABLE语句修改新表"_t_new"

4.创建3个触发器,名称格式为pt_osc_库名_表名_操作类型,比如

CREATE TRIGGER `pt_osc_dba_t_del` AFTER DELETE ON `dba`.`t` FOR EACH ROW DELETE IGNORE FROM `dba`.`_t_new` WHERE `dba`.`_t_new`.`id` <=> OLD.`id`

CREATE TRIGGER `pt_osc_dba_t_upd` AFTER UPDATE ON `dba`.`t` FOR EACH ROW REPLACE INTO `dba`.`_t_new` (`id`, `a`, `b`, `c1`) VALUES (NEW.`id`, NEW.`a`, NEW.`b`, NEW.`c1`)

CREATE TRIGGER `pt_osc_dba_t_ins` AFTER INSERT ON `dba`.`t` FOR EACH ROW REPLACE INTO `dba`.`_t_new` (`id`, `a`, `b`, `c1`) VALUES (NEW.`id`, NEW.`a`, NEW.`b`, NEW.`c1`)

5.开始复制数据,比如

INSERT LOW_PRIORITY IGNORE INTO `dba`.`_t_new` (`id`, `a`, `b`, `c1`) SELECT `id`, `a`, `b`, `c1` FROM `dba`.`t` LOCK IN SHARE MODE /*pt-online-schema-change 28014 copy table*/

注意：对原表加共享锁，会阻塞所有排他锁

6.复制完成后,交互原表和新表,执行RENAME命令,如 RENAME TABLE t to _t_old, _t_new to t;

7.删除老表,_t_old

8.删除触发器

9.修改完成

 

注意：如果异常终止程序，触发器不会自动删除，如果要删除新表，那么要先删除触发器，否则向老表插入数据会因为找不到新表而报错

 

 

注意事项：

pt-online-schema-change 在线DDL工具，虽然说不会锁表，但是对性能还是有一定的影响，执行过程中对全表做一次select。这个过程会将buffer_cache中活跃数据全部交换一遍，这就导致活跃数据的请求都要从磁盘获取，导致慢SQL增多，file_reads增大。所以对于大表应在业务低峰期执行该操作
执行 RENAME 时，你不能有任何锁定的表或活动的事务。你同样也必须有对原初表的 ALTER 和 DROP 权限，以及对新表的 CREATE 和 INSERT 权限。当业务量较大时，修改操作会等待没有数据修改后，执行最后的rename操作。因此，在修改表结构时，应该尽量选择在业务相对空闲时，至少修改表上的数据操作较低时，执行较为妥当。如果在多表更名中，MySQL 遭遇到任何错误，它将对所有被更名的表进行倒退更名，将每件事物退回到最初状态。
对表的慢查询操作，慢查询还未结束执行osc操作，会报错，超时错误，在创建触发器的时候退出。
对于主从复制架构。 考虑到主从的一致性，应该在主库上执行pt-online-schema-change操作。
ps：如果在误在从库上执行了pt-online-schema-change操作，未执行完成不要取消，等到执行完成了，在修改成原来的状态。

如果在误在从库上执行了pt-online-schema-change操作，未执行完成取消的话，删除有 pt-online-schema-change在从库上创建的临时表和触发器即可。

####################################################################

在原表上建立三个触发器，如下：
（1）CREATETRIGGER mk_osc_del AFTER DELETE ON $table ” “FOR EACH ROW ”
“DELETE IGNORE FROM $new_table “”WHERE$new_table.$chunk_column = OLD.$chunk_column”;
（2）CREATETRIGGER mk_osc_ins AFTER INSERT ON $table ” “FOR EACH ROW ”
“REPLACEINTO $new_table ($columns) ” “VALUES($new_values)”;
（3）CREATETRIGGER mk_osc_upd AFTER UPDATE ON $table ” “FOR EACH ROW ”
“REPLACE INTO $new_table ($columns) “”VALUES ($new_values)”;

我们可以看到这三个触发器分别对应于INSERT、UPDATE、DELETE三种操作：
（1）mk_osc_del，DELETE操作，我们注意到DELETEIGNORE，当新有数据时，我们才进行操作，也就是说，当在后续导入过程中，如果删除的这个数据还未导入到新表，那么我们可以不在新表执行操作，因为在以后的导入过程中，原表中改行数据已经被删除，已经没有数据，那么他也就不会导入到新表中；
（2）mk_osc_ins，INSERT操作，所有的INSERT INTO全部转换为REPLACEINTO，为了确保数据的一致性，当有新数据插入到原表时，如果触发器还未把原表数据未同步到新表，这条数据已经被导入到新表了，那么我们就可以利用replaceinto进行覆盖，这样数据也是一致的
（3）mk_osc_upd UPDATE操作，所有的UPDATE也转换为REPLACEINTO，因为当跟新的数据的行还未同步到新表时，新表是不存在这条记录的，那么我们就只能插入该条数据，如果已经同步到新表了，那么也可以进行覆盖插入，所有数据与原表也是一致的；
我们也能看出上述的精髓也就这这几条replaceinto操作，正是因为这几条replaceinto才能保证数据的一致性
4、拷贝原表数据到临时表中，在脚本中使用如下语句
INSERT IGNORE INTO $to_table ($columns) ” “SELECT $columns FROM $from_table “”WHERE ($chunks->[$chunkno])”，我们能看到他是通过一些查询（基本为主键、唯一键值）分批把数据导入到新的表中，在导入前，我们能通过参数–chunk-size对每次导入行数进行控制，已减少对原表的锁定时间，并且在导入时，我们能通过—sleep参数控制，在每个chunk导入后与下一次chunk导入开始前sleep一会，sleep时间越长，对于磁盘IO的冲击就越小
5、Rename 原表到old表中，在把临时表Rename为原表，
“RENAME TABLE `$db`.`$tmp_tbl`TO `$db`.`$tbl`”; 在rename过程，其实我们还是会导致写入读取堵塞的，所以从严格意思上说，我们的OSC也不是对线上环境没有一点影响，但由于rename操作只是一个修改名字的过程，也只会修改一些表的信息，基本是瞬间结束，故对线上影响不太大
6、清理以上过程中的不再使用的数据，如OLD表

https://www.percona.com/doc/percona-toolkit/2.1/pt-online-schema-change.html

https://www.percona.com/downloads/percona-toolkit/LATEST/

https://www.percona.com/doc/percona-toolkit/3.0/pt-online-schema-change.html

https://www.percona.com/downloads/percona-toolkit/LATEST/

在运维mysql数据库时，我们总会对数据表进行ddl 变更，修改添加字段或者索引，对于mysql 而已，ddl 显然是一个令所有MySQL dba 诟病的一个功能，因为在MySQL中在对表进行ddl时，会锁表，当表比较小比如小于1w上时，对前端影响较小，当时遇到千万级别的表 就会影响前端应用对表的写操作。

perconal 推出一个工具 pt-online-schema-change ，其特点是修改过程中不会造成读写阻塞。

原理：

1.建立一个与需要操作的表相同表结构的空表

2.给空表执行表结构修改

3.在原表上增加delete/update/insert的after trigger

4.copy数据到新表

5.将原表改名，并将新表改成原表名

6.删除原表

7.删除trigger

pt-osc限制条件：

1.表要有主键，否则会报错;

2.表不能有trigger;

使用方法：

1.下载

wget percona.com/get/percona-toolkit.tar.gz

2.安装

tar -zxvf percona-toolkit.tar.gz

cd percona-toolkit-3.0.4

perl Makefile.PL

(若执行Makefile出错 则需先执行yum install perl-ExtUtils-CBuilder perl-ExtUtils-MakeMaker)

make

make test

make install

3.使用方法

pt-online-schema-change --help 可查看参数帮助

若查看参数提示Can't locate Digest/MD5.pm in @INC错误 则需执行yum -y install perl-Digest-MD5安装相关组件

提示缺少perl-DBI模块，那么直接 yum install perl-DBI

场景1：增加列

pt-online-schema-change --host=192.168.0.0 -uroot -pyourpassword --alter "add column age int(11) default null" D=test,t='test_tb' --execute --print --statistics --no-check-alter

场景2：删除列

pt-online-schema-change --host=192.168.0.0 -uroot -pyourpassword --alter "drop column age" D=test,t='test_tb' --execute --print --statistics --no-check-alter

场景3：更改列

pt-online-schema-change --host=192.168.0.0 -uroot -pyourpassword --alter "CHANGE id id_num int(20)" D=test,t='test_tb' --execute --print --statistics --no-check-alter

场景4：创建索引

pt-online-schema-change --host=192.168.0.0 -uroot -pyourpassword --alter "add index indx_ukid(address_ukid)" D=test,t='address_tb' --execute --print --statistics --no-check-alter

pt-online-schema-change 是什么？
pt-online-schema-change是Percona工具包的一员，用于修改表而不会造成读锁或者写锁；

pt-online-schema-change详细描述：
pt-online-schema-change工作在一个副本上，所以原表没有被lock，client可以继续进行read和write操作；pt-online-schema-change会创建一个空表，然后根据需要对空表进行modify，之后将原表数据复制到新表，复制完成后就用RENAME将新表替换原表，完成后 drop 原表；

pt-online-schema-change修改以块为单位，修改中根据情况自行调整块大小（类似checksum） 

为了安全，如果没执行--execute，则不会发生修改；工具会有多种措施，防止不必要的负载，如：

1， 多数情况下，如果没有PRIMARY KEY或者UNIQUE KEY，则拒绝操作（--alter）

2， 如果检测到replication filters，则拒绝操作（--[no]check-replication-filters）

3， 如果有较大的延时，则暂停操作（"--max-lag"）

4， 如果有较大的负载则暂停或放弃操作（"--max-load" and "--critical-load）

5， 工具会设置自己会话的锁超时时间"innodb_lock_wait_timeout=1"lock_wait_timeout=60，所以会更少 地影响其他事物（修改时间"--set-vars"）

6， 默认如果有外键约束则拒绝修改（忽略约束"--alter-foreign-keys-method"）

7， 不支持Percona XtraDB Cluster（Percona高可用解决方案）

工具会忽略mysql 5.7+版本的"GENERATED"列，因为它根据表达式计算列值

输出信息：

    复制表的时候会显示进度， 如果要显示额外信息，则--print, 如果指定--statistics 最后会显示计数信息

选项OPTIONS
--alter

     不需要加alter table关键字，如pt-online-schema-change  --alter "ADD COLUMN c2 INT" D=apple,t=a 

     限制： 

     1， 必须要有PRIMARY KEY或者UNIQUE KEY，因为它会在程序进行的时候创建DELETE 触发器，来保证新表跟着老表一起更新；

     2， RENAME不会被采用去RENAME TABLE

     3， 如果增加列的时候指定NOT NULL，并且未提供默认值，则失败；

     4， 如果要删除外键，需要指定 "_constraint_name"而不是 "constraint_name"

     5， 如果要将MYISAM表修改为Innodb会报错

--[no]analyze-before-swap  当新表与原表交换的前执行analyze table；

--ask-pass 连接的时候会要求提供密码

--charset 默认字符集

--[no]check-alter   对于危险的操作会产生告警，如删除主键或者重命名列；

--check-interval   当主从延时为--max-lag 秒的时候，检查sleep的时间；例如设置为10的时候，主从延时产生，则暂停10秒，再检查延时，如果依然有有延时，停止10s再检查； 默认1s

--[no]check-plan   是否在执行前用EXPLAIN检查语句

--[no]check-replication-filters  是否检查过滤器，如果有任何过滤器的时候，操作都会被放弃；（如：binlog_ignore_db  replicate_do_db）

--check-slave-lag  如果slave的时间戳低于"--max-lag"的时候则暂停操作

--chunk-index  指定chunk的索引，工具默认会选择最合适的索引，当然你也可以手动指定

--chunk-index-columns 有复合索引的时候，指定索引列

--chunk-size   chunk的行数，默认1000

 --chunk-size-limit   默认4，chunk超过应有值的4倍大小则跳过。如果没有唯一索引，则chunk 大小是不精确的，工具会用EXPLAIN评估大小，如果超过需要chunk大小的n倍，则跳过该chunk；

--chunk-time 默认0.5， 工具会根据每次复制数据花费的时间自动调整chunk大小，尽可能使每次时间都相同；=0则不调节

--config  读取配置文件，多个文件逗号分隔；如果指定则必须是第一个参数

--critical-load  默认Threads_running=50；  每次chunk执行后会自动用SHOW GLOBAL STATUS检查负载情况，如果超过阈值则放弃；

--database -D 指定库

--default-engine 该选项将导致新表使用系统默认引擎；   而不是原有表一致的引擎

--data-dir  v 5.6+ 在不同的分区上创建新表

--remove-data-dir 如果原表已通过--data-dir 创建，该参数会删除它并在默认datadir下创建新表

--defaults-file   从文件中读取选项

 --[no]drop-new-table 如果复制原表失败则删除新表； 也可以no-xxx来保留新表

 --[no]drop-old-table  rename新表后drop旧表，可以no-xxx来保留旧表

--[no]drop-triggers  删除旧表的触发器；"--no-drop-triggers" forces "--no-drop-old-table".

--dry-run 创建并修改新表，但不创建触发器，也不复制表，或者替换原表，与--execute互斥

--execute 执行操作 与 --dry-run互斥

--[no]check-unique-key-change  如果--alter尝试增加唯一索引的的话，则工具不会运行； 因为增加唯一索引如果列有重复值，会发生丢失数据的情况；因为INSERT的时候默认采用"INSERT IGNORE"

--force  强制运行，可能打破外键约束

--help 帮助
--host | -h  指定host

--max-lag 默认1s， 如果主从延时的时间超过这个值，则复制会暂停"--check-interval"秒时间；然后再检查，直到主从延时小于该值；如果指定了 "--check-slave-lag"，则只会检查指定slave延时，而不是检查所有slave;如果有任何SLAVE stop了，那么工具会一直等待下去；每次停止的时候都会打印报告

--max-load 默认 : Threads_running=25 

        每次复制chunk后检查负载，如果超过该值则暂停；类似--critical-load

--preserve-triggers     保留原表的触发器，不删除

--new-table-name   指定新表的名字，默认是 tablename_new

--null-to-not-null  修改允许null值为not null

--password 指定密码

--pause-file 改参数指定的文件存在的时，操作将会暂停

--pid 新建pid文件

--port 指定连接端口

--print  将会显示工具执行的命令

--progress 显示进度，默认：  time,30

       两部分组成，第一部分percentage,time或者iterations；  第二部分多久更新一次

--quiet    -q

  不打印消息到屏幕； 禁用--progress,  ERROR和告警还是会打印

--recurse   

    当发现有slave的时候，指定递归的层数 ， int值

--recursion-method     查找Slave的递归方法

                METHOD       USES
             ===========  ==================
             processlist  SHOW PROCESSLIST
             hosts        SHOW SLAVE HOSTS
             dsn=DSN      DSNs from a table
             none         Do not find slaves  

--skip-check-slave-lag  检查SLAVE的时候，指定该SLAVE跳过；

--slave-user    设置连接slave的用户；可以指定用更少权限的用户连接SLAVE，但是所有slave都必须有该账户

--slave-password    --slave-user 的密码

--set-vars  默认：

              wait_timeout=10000
              innodb_lock_wait_timeout=1
              lock_wait_timeout=60

--sleep  复制每个chunk的时候暂停多少s，默认0

--socket 指定连接socket

--statistics   打印计数器的数据，当有告警的时候这个参数比较有用

-[no]swap-tables  指定是否替换旧表

--tries  重试次数；

--user 连接用户名

--version  版本

--[no]version-check 检查Percona Toolkit的最新版本

[root@localhost tmp]# pt-online-schema-change  --alter "ADD COLUMN c3 INT" D=apple,t=a -u root -h localhost -p password --execute
No slaves found.  See --recursion-method if host localhost.localdomain has slaves.
Not checking slave lag because no slaves were found and --check-slave-lag was not specified.
 
# A software update is available:
Operation, tries, wait:
  analyze_table, 10, 1
  copy_rows, 10, 0.25
  create_triggers, 10, 1
  drop_triggers, 10, 1
  swap_tables, 10, 1
  update_foreign_keys, 10, 1
Altering `apple`.`a`...
Creating new table...
Created new table apple._a_new OK.
Altering new table...
Altered `apple`.`_a_new` OK.
2018-09-05T12:57:33 Creating triggers...
2018-09-05T12:57:33 Created triggers OK.
2018-09-05T12:57:33 Copying approximately 2 rows...
2018-09-05T12:57:33 Copied rows OK.
2018-09-05T12:57:33 Analyzing new table...
2018-09-05T12:57:33 Swapping tables...
2018-09-05T12:57:34 Swapped original and new tables OK.
2018-09-05T12:57:34 Dropping old table...
2018-09-05T12:57:34 Dropped old table `apple`.`_a_old` OK.
2018-09-05T12:57:34 Dropping triggers...
2018-09-05T12:57:34 Dropped triggers OK.
