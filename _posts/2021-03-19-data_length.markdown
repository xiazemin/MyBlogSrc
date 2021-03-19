---
title: data length mysql
layout: post
category: storage
author: 夏泽民
---

https://blog.csdn.net/zgljl2012/article/details/78983706

https://zhuanlan.zhihu.com/p/58999756

select sum(DATA_LENGTH)+sum(INDEX_LENGTH)
from information_schema.tables  
where table_schema='数据库名';

在mysql中有一个information_schema数据库，这个数据库中装的是mysql的元数据，包括数据库信息、数据库中表的信息等。所以要想查询数据库占用磁盘的空间大小可以通过对information_schema数据库进行操作。
information_schema中的表主要有：
schemata表：这个表里面主要是存储在mysql中的所有的数据库的信息
tables表：这个表里存储了所有数据库中的表的信息，包括每个表有多少个列等信息。
columns表：这个表存储了所有表中的表字段信息。
statistics表：存储了表中索引的信息。
user_privileges表：存储了用户的权限信息。
schema_privileges表：存储了数据库权限。
table_privileges表：存储了表的权限。
column_privileges表：存储了列的权限信息。
character_sets表：存储了mysql可以用的字符集的信息。
collations表：提供各个字符集的对照信息。
collation_character_set_applicability表：相当于collations表和character_sets表的前两个字段的一个对比，记录了字符集之间的对照信息。
table_constraints表：这个表主要是用于记录表的描述存在约束的表和约束类型。
key_column_usage表：记录具有约束的列。
routines表：记录了存储过程和函数的信息，不包含自定义的过程或函数信息。
views表：记录了视图信息，需要有show view权限。
triggers表：存储了触发器的信息，需要有super权限

https://developer.aliyun.com/article/620517
<!-- more -->
https://www.jianshu.com/p/ea15158f39f7

https://zhuanlan.zhihu.com/p/88342863

https://segmentfault.com/a/1190000023781542

https://blog.csdn.net/qq_39455116/article/details/96480845


https://github.com/DATA-DOG/go-sqlmock

