---
title: timestamp
layout: post
category: storage
author: 夏泽民
---
其实，表达方式还是蛮多的，汇总如下：

CURRENT_TIMESTAMP

CURRENT_TIMESTAMP()

NOW()

LOCALTIME

LOCALTIME()

LOCALTIMESTAMP

LOCALTIMESTAMP()

TIMESTAMP和DATETIME的相同点：

1> 两者都可用来表示YYYY-MM-DD HH:MM:SS[.fraction]类型的日期。

 

TIMESTAMP和DATETIME的不同点：

1> 两者的存储方式不一样

对于TIMESTAMP，它把客户端插入的时间从当前时区转化为UTC（世界标准时间）进行存储。查询时，将其又转化为客户端当前时区进行返回。

而对于DATETIME，不做任何改变，基本上是原样输入和输出。

两者所能存储的时间范围不一样

timestamp所能存储的时间范围为：'1970-01-01 00:00:01.000000' 到 '2038-01-19 03:14:07.999999'。

datetime所能存储的时间范围为：'1000-01-01 00:00:00.000000' 到 '9999-12-31 23:59:59.999999'。
<!-- more -->
https://www.cnblogs.com/mxwz/p/7520309.html

https://www.cnblogs.com/liuxs13/p/9760812.html

1.CURRENT_TIMESTAMP 

当要向数据库执行insert操作时，如果有个timestamp字段属性设为 

CURRENT_TIMESTAMP，则无论这个字段有木有set值都插入当前系统时间 

2.ON UPDATE CURRENT_TIMESTAMP

当执行update操作是，并且字段有ON UPDATE CURRENT_TIMESTAMP属性。则字段无论值有没有变化，他的值也会跟着更新为当前UPDATE操作时的时间。

 https://www.cnblogs.com/banye/p/7066021.html
 
 
