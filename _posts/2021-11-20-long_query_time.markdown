---
title: long_query_time
layout: post
category: mysql
author: 夏泽民
---
MySQL的慢查询日志是MySQL提供的一种日志记录，它用来记录在MySQL中响应时间超过阀值的语句，具体指运行时间超过long_query_time值的SQL，则会被记录到慢查询日志中。long_query_time的默认值为10，意思是运行10S以上的语句。默认情况下，Mysql数据库并不启动慢查询日志，需要我们手动来设置这个参数，当然，如果不是调优需要的话，一般不建议启动该参数，因为开启慢查询日志会或多或少带来一定的性能影响。慢查询日志支持将日志记录写入文件，也支持将日志记录写入数据库表。

MySQL慢查询定义
分析MySQL语句查询性能的方法除了使用 EXPLAIN 输出执行计划，还可以让MySQL记录下查询超过指定时间的语句，我们将超过指定时间的SQL语句查询称为“慢查询”。

MySQL慢查询的体现
慢查询主要体现在慢上，通常意义上来讲，只要返回时间大于 >1 sec上的查询都可以称为慢查询。慢查询会导致CPU，内存消耗过高。数据库服务器压力陡然过大，那么大部分情况来讲，肯定是由某些慢查询导致的。
<!-- more -->
mysql慢查询开启
1.查看当前慢查询设置情况

#查看慢查询时间，默认10s，建议降到1s或以下
mysql> show variables like "long_query_time";

2.如何开启慢查询功能
方法一：在服务器上找到mysql的配置文件my.cnf , 然后再mysqld模块里追加一下内容，这样的好处是会一直生效，不好就是需要重启mysql进程。

vim my.cnf
[mysqld]
slow_query_log = ON
#定义慢查询日志的路径
slow_query_log_file = /tmp/slow_querys.log
#定义查过多少秒的查询算是慢查询，我这里定义的是1秒，5.6之后允许设置少于1秒，例如0.1秒
long_query_time = 1
#用来设置是否记录没有使用索引的查询到慢查询记录,默认关闭,看需求开启,会产生很多日志,可动态修改
#log-queries-not-using-indexes
管理指令也会被记录到慢查询。比如OPTIMEZE TABLE, ALTER TABLE,默认关闭,看需求开启,会产生很多日志,可动态修改
#log-slow-admin-statements

然后重启mysql服务器即可，这是通过一下命令看一下慢查询日志的情况：

tail -f /tmp/slow_querys.log

方法二：通过修改mysql的全局变量来处理，这样做的好处是，不用重启mysql服务器，登陆到mysql上执行一下sql脚本即可，不过重启后就失效了。

#开启慢查询功能，1是开启，0是关闭
mysql> set global slow_query_log=1;
#定义查过多少秒的查询算是慢查询，我这里定义的是1秒，5.6之后允许设置少于1秒，例如0.1秒
mysql> set global long_query_time=1;
#定义慢查询日志的路径
mysql> set global slow_query_log_file='/tmp/slow_querys.log';
#关闭功能：set global slow_query_log=0;
然后通过一下命令查看是否成功
mysql> show variables like 'long%';
mysql> show variables like 'slow%';
#设置慢查询记录到表中
#set global log_output='TABLE';

特别要注意的是long_query_time的设置，5.6之后支持设置低于0.1秒，所以记录的详细程度，就看你自己的需求，数据库容量比较大的，超过0.1秒还是比较多，所以就变得有点不合理了。

MYSQL慢查询日志的记录定义
直接查看mysql的慢查询日志分析，比如我们可以tail -f slow_query.log查看里面的内容

https://www.cnblogs.com/yangfei123/p/12651543.html

https://www.cnblogs.com/kerrycode/p/5593204.html

1.错误日志作用
错误日志记录了mysql启动和停止时。以及server执行过程中发生不论什么严重性错误的相关信息。当数据库出现不论什么故障导致无法启动时候。比方mysql启动异常。我们可首先检查此日志。在mysql中，错误日志日志（还有其它日志），不只能够存储在文件里。当然还能够存储到数据的表中。至于实现方式。笔者也正在研究中···
2.错误日志控制与使用
1.配置
通过log-error=[file-name]来配置（在mysql的配置文件里），假设没有指定file_name，mysqld使用错误日志名为host_name.err(host_name为主机名)。并默认在參数datadir（保存数据的文件夹）指定的文件夹中写入日志文件。
比方我本地使用的是WampServer集成环境
当中log-error=D:/wamp/logs/mysql.log

2.查看错误日志
错误日志的格式：时间 [错误级别] 错误信息
假设你感觉通过mysql配置文件来定位错误日志所在位置比較麻烦。你全然能够通过再client通过命令来查看错误日志所在位置
使用命令式：show variables like 'log_error';

二进制日志
1.作用二进制日志（又叫binlog日志）记录了全部的DDL（数据定义语言）语句和DML（数据操作语言）语句。可是不包含数据查询语句。语句是以“事件”的形式保存的，它描写叙述数据更改的过程。该日志的两个主要功能是：数据的恢复与数据的复制。
数据的恢复：MySQL本身具备数据备份和恢复功能。
比方，我们每天午夜12:00进行数据的备份。

假设某天，下午13:00,数据库出现问题。导致数据库内容丢失。

我们能够通过二进制日志解决问题。解决思路是。能够先将前一天午夜12:00的数据备份文件恢复到数据库，然后再使用二进制日志回复从前一天午夜12:00到当天13:00对数据库的操作。

数据复制：MySQL支持主从server间的数据复制功能，并通过该功能实现数据库的冗余机制以保证数据库的可用性和提高数据库德性能。MySQL正是通过二进制日志实现数据的传递。
主server上的二进制日志内容会被发送到各个从server上。并在每一个从server上运行，从而保证了主从server之间数据的一致性。

2.二进制日志控制与使用
1.开启
在默认情况下，mySQL不会记录二进制日志。如何才干开启MySQL的二进制日志记录功能呢？
我们能够通过MySQL的配置文件来控制MySQL启动二进制日志记录功能。通过改动參数log-bin=[base_name]来启动MySQL二进制日志。
mySQL会将改动的数据库内容的语句记录到以 base_name-bin.0000x为名的日志文件里。当中bin代表binary。后缀00000x代表二进制日志文件的顺序，每次启动Mysql，日志文件顺序会自己主动加1.假设base_name未定义。MySQL将使用pid-file參数设置的值作为二进制日志文件的基础名字。

比方我将log-bin文件名称定为mybinlog。那么将会在D:/wamp/bin/mysql/mysql5.6.17/data文件夹下，生成mybinlog.00000x的二进制日志文件。

https://www.cnblogs.com/zhchoutai/p/8520295.html


