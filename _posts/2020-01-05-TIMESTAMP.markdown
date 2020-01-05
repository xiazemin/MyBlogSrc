---
title: TIMESTAMP
layout: post
category: storage
author: 夏泽民
---
明明数据有更新，update_time字段却还停留在创建数据的时候。

按常理来说这个字段应该是自动更新的才对。

查了一下表结构，`update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP

发现update_time字段的类型是datetime

由此牵扯出两个问题，（1）timestamp与datetime的区别；（2）CURRENT_TIMESTAMP为什么能用于datetime类型

timestamp与datetime的区别
a）DATETIME的默认值为null；TIMESTAMP的字段默认不为空（not null）,默认值为当前时间（CURRENT_TIMESTAMP），如果不做特殊处理，并且update语句中没有指定该列的更新值，则默认更新为当前时间。
这个区别就解释了为什么平时我们都不用可以管这个字段就能自动更新了，因为多数时候用的是timestamp；而此处用的是datetime，不会有自动更新当前时间的机制，所以需要在上层手动更新该字段

b）DATETIME使用8字节的存储空间，TIMESTAMP的存储空间为4字节。因此，TIMESTAMP比DATETIME的空间利用率更高。

这个区别解释了为啥timestamp类型用的多

c）两者的存储方式不一样 ，对于TIMESTAMP，它把客户端插入的时间从当前时区转化为UTC（世界标准时间）进行存储。查询时，将其又转化为客户端当前时区进行返回。而对于DATETIME，不做任何改变，基本上是原样输入和输出。

d）两者所能存储的时间范围不一样 

timestamp所能存储的时间范围为：’1970-01-01 00:00:01.000000’ 到 ‘2038-01-19 03:14:07.999999’；

datetime所能存储的时间范围为：’1000-01-01 00:00:00.000000’ 到 ‘9999-12-31 23:59:59.999999’。

CURRENT_TIMESTAMP为什么能用于datetime类型
在mysql 5.6之前的版本，CURRENT_TIMESTAMP只能用于timestamp类型，
5.6版本之后，CURRENT_TIMESTAMP也能用于datetime类型了
select version()查了一下数据库发现确实版本是5.6.29
<!-- more -->
一个完整的日期格式如下：YYYY-MM-DD HH:MM:SS[.fraction]，它可分为两部分：date部分和time部分，其中，date部分对应格式中的“YYYY-MM-DD”，time部分对应格式中的“HH:MM:SS[.fraction]”。对于date字段来说，它只支持date部分，如果插入了time部分的内容，它会丢弃掉该部分的内容，并提示一个warning。

timestamp和datetime的相同点：

（1） 两者都可用来表示YYYY-MM-DD HH:MM:SS[.fraction]类型的日期。

timestamp和datetime的不同点：

（1）两者的存储方式不一样

对于TIMESTAMP，它把客户端插入的时间从当前时区转化为UTC（世界标准时间）进行存储。查询时，将其又转化为客户端当前时区进行返回。

而对于DATETIME，不做任何改变，基本上是原样输入和输出。

（2）两者所能存储的时间范围不一样

timestamp所能存储的时间范围为：'1970-01-01 00:00:01.000000' 到 '2038-01-19 03:14:07.999999'。

datetime所能存储的时间范围为：'1000-01-01 00:00:00.000000' 到 '9999-12-31 23:59:59.999999'。

MySQL 中常用的两种时间储存类型分别是datetime和 timestamp。如何在它们之间选择是建表时必要的考虑。下面就谈谈他们的区别和怎么选择。

1 区别
1.1 占用空间
类型	占据字节	表示形式
datetime	8 字节	yyyy-mm-dd hh:mm:ss
timestamp	4 字节	yyyy-mm-dd hh:mm:ss
1.2 表示范围
类型	表示范围
datetime	'1000-01-01 00:00:00.000000' to '9999-12-31 23:59:59.999999'
timestamp	'1970-01-01 00:00:01.000000' to '2038-01-19 03:14:07.999999'
timestamp翻译为汉语即"时间戳"，它是当前时间到 Unix元年(1970 年 1 月 1 日 0 时 0 分 0 秒)的秒数。对于某些时间的计算，如果是以 datetime 的形式会比较困难，假如我是 1994-1-20 06:06:06 出生，现在的时间是 2016-10-1 20:04:50 ，那么要计算我活了多少秒钟用 datetime 还需要函数进行转换，但是 timestamp 直接相减就行。

1.3 时区
timestamp 只占 4 个字节，而且是以utc的格式储存， 它会自动检索当前时区并进行转换。

datetime以 8 个字节储存，不会进行时区的检索.

也就是说，对于timestamp来说，如果储存时的时区和检索时的时区不一样，那么拿出来的数据也不一样。对于datetime来说，存什么拿到的就是什么。

还有一个区别就是如果存进去的是NULL，timestamp会自动储存当前时间，而 datetime会储存 NULL。

2 测试
我们新建一个表

image

插入数据

image

查看数据，可以看到存进去的是NULL，timestamp会自动储存当前时间，而 datetime会储存NULL

image

把时区修改为东 9 区，再查看数据，会会发现 timestamp 比 datetime 多一小时

image

如果插入的是无效的呢？假如插入的是时间戳

image

结果是0000-00-00 00:00:00，根据官方的解释是插入的是无效的话会转为 0000-00-00 00:00:00，而时间戳并不是MySQL有效的时间格式。

那么什么形式的可以插入呢，下面列举三种

//下面都是 MySQL 允许的形式，MySQL 会自动处理
2016-10-01 20:48:59
2016#10#01 20/48/59
20161001204859
3 选择
如果在时间上要超过Linux时间的，或者服务器时区不一样的就建议选择 datetime。

如果是想要使用自动插入时间或者自动更新时间功能的，可以使用timestamp。

如果只是想表示年、日期、时间的还可以使用 year、 date、 time，它们分别占据 1、3、3 字节，而datetime就是它们的集合。


1）timestamp：4个字节，（北京时间：2038年1月19日中午11:14:07）之后无法正常工作

2）datetime：8个字节

当涉及到日期计算、应用需要跨多个时区（国际业务）等，使用时间戳。

timestamp 在不同时区下能确保时间的精确性。

总体来说，存储时间优先使用时间戳较好



使用DATETIME(3),TIMESTAMP(6)能达到存储毫秒微秒的效果，而且还会自动四舍五入。



应用中如何录入毫秒微秒呢？使用now(3)，now(6)

mysql> select now(3),now(6);
+-------------------------+----------------------------+
| now(3)                  | now(6)                     |
+-------------------------+----------------------------+
| 2017-07-18 09:24:14.969 | 2017-07-18 09:24:14.969081 |
+-------------------------+----------------------------+
1 row in set (0.00 sec)


应用中如何查询当前的毫秒数，微秒数呢？

mysql> SELECT SUBSTR(MICROSECOND(now(3)),1,3) '毫 秒',MICROSECOND(NOW(3)) '微 秒';       
+--------+--------+
| 毫秒   | 微秒   |
+--------+--------+
| 332    | 332000 |
+--------+--------+


一、MySQL 获得毫秒、微秒及对毫秒、微秒的处理

MySQL 较新的版本中（MySQL 6.0.5），也还没有产生微秒的函数，now() 只能精确到秒。 MySQL 中也没有存储带有毫秒、微秒的日期时间类型。

但，奇怪的是 MySQL 已经有抽取（extract）微秒的函数。例如：

select microsecond('12:00:00.123456');                          -- 123456
select microsecond('1997-12-31 23:59:59.000010');               -- 10
select extract(microsecond from '12:00:00.123456');             -- 123456
select extract(microsecond from '1997-12-31 23:59:59.000010');  -- 10
select date_format('1997-12-31 23:59:59.000010', '%f');         -- 000010
尽管如此，想在 MySQL 获得毫秒、微秒还是要在应用层程序中想办法。假如在应用程序中获得包含微秒的时间：1997-12-31 23:59:59.000010，在 MySQL 存放时，可以设计两个字段：c1 datetime, c2 mediumint，分别存放日期和微秒。为什么不采用 char 来存储呢？用 char 类型需要 26 bytes，而 datetime + mediumint 只有 11（8+3） 字节。