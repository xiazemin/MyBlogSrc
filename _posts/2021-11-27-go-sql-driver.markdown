---
title: go-sql-driver
layout: post
category: mysql
author: 夏泽民
---
golang go-sql-driver 数据库报错 bad connection

原因
这是因为Mysql服务器主动关闭了Mysql链接。
在项目中使用了一个mysql链接，同时使用了事务，处理多个表操作。处理时间长。
导致空闲链接超时，Mysql关闭了链接。而客户端保持了已经关闭的链接。

具体原因是：
beego没有调用db.SetConnMaxLifetime 这个方法，导致客户端保持了已经关闭的链接。

https://wwblog.csdn.net/whatday/article/details/103952962

go driver: bad connection
driver: bad connection

原因：rows没Close
<!-- more -->

Mysql 查看连接数,状态 最大并发数
-- show variables like '%max_connections%'; 查看最大连接数
set global max_connections=1000 重新设置
mysql> show status like 'Threads%';
+-------------------+-------+
| Variable_name     | Value |
+-------------------+-------+
| Threads_cached    | 58    |
| Threads_connected | 57    |   ###这个数值指的是打开的连接数
| Threads_created   | 3676  |
| Threads_running   | 4     |   ###这个数值指的是激活的连接数，这个数值一般远低于connected数值
+-------------------+-------+
 
Threads_connected 跟show processlist结果相同，表示当前连接数。准确的来说，Threads_running是代表当前并发数

命令： show processlist; 
如果是root帐号，你能看到所有用户的当前连接。如果是其它普通帐号，只能看到自己占用的连接。 
show processlist;只列出前100条，如果想全列出请使用show full processlist; 
MySQL> show processlist;

https://blog.csdn.net/wsf568582678/article/details/53636747

https://blog.csdn.net/weixin_39862697/article/details/113316610



