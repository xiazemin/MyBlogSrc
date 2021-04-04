---
title: mysql Cancellations
layout: post
category: storage
author: 夏泽民
---
https://www.alexedwards.net/blog/how-to-manage-database-timeouts-and-cancellations-in-go
https://medium.com/@rocketlaunchr.cloud/canceling-mysql-in-go-827ed8f83b30


<!-- more -->
https://golangrepo.com/repo/go-sql-driver-mysql-go-database-drivers

https://blog.zhenlanghuo.top/2019/07/21/Golang%20databasesql%E4%B8%8Ego-sql-drivermysql%20%E6%BA%90%E7%A0%81%E9%98%85%E8%AF%BB%E7%AC%94%E8%AE%B0%20--%20go-sql-drivermysql%E7%AF%87/

https://github.com/ngrok/sqlmw
https://github.com/gchaincl/sqlhooks

https://github.com/shogo82148/go-sql-proxy

https://github.com/shogo82148/go-sql-proxy/security

sql 的query 后不能直接取消，需要等到rows.Scan后，因为共用了数据结构，否则会出现 scan到的数据为空的情况
