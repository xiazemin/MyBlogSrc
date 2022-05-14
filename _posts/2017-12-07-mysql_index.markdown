---
title: mysql_index
layout: post
category: storage
author: 夏泽民
---
<!-- more -->
如果想知道MySQL数据库中每个表占用的空间、表记录的行数的话，可以打开MySQL的 information_schema 数据库。在该库中有一个 TABLES 表，这个表主要字段分别是：

TABLE_SCHEMA : 数据库名
TABLE_NAME：表名
ENGINE：所使用的存储引擎
TABLES_ROWS：记录数
DATA_LENGTH：数据大小
INDEX_LENGTH：索引大小

其他字段请参考MySQL的手册，我们只需要了解这几个就足够了。

1  首先查看某一实例下的所有占用磁盘空间（表数据+索引数据，得到的结果为B，这里做了数据处理转成M）：

select concat(round((sum(DATA_LENGTH)+sum(INDEX_LENGTH))/1024/1024,2),'M') from information_schema.tables where table_schema='实例名称';
 上面是查询所有的表计的累计量，下面是是查询单个表计的的SQL(按照实例名查询)：

select table_name,
DATA_LENGTH/1024/1024 as tablesData,
INDEX_LENGTH/1024/1024 as indexData 
from information_schema.tables
where table_schema='dsm'
ORDER BY  tablesData desc;
 
