---
title: mysql parseTime
layout: post
category: storage
author: 夏泽民
---
golang：unsupported Scan, storing driver.Value type []uint8 into type *time.Time

解决：
在open连接后拼接参数：parseTime=true 即可
db, err := gorm.Open("mysql", "user:password@/dbname?charset=utf8&parseTime=True&loc=Local")
  defer db.Close()
  
  parseTime是查询结果是否自动解析为时间。
loc是MySQL的时区设置。
<!-- more -->
想要把 time.Time 直接存储入库，需要开启解析时间parseTime
db, err := sql.Open("mysql", "user:password@/dbname?charset=utf8mb4&parseTime=true")
golang 程序里 time.Time 为 2018-12-24 18:00:00 CST

转为 UTC 存储到 mysql 2018-12-24 10:00:00

golang 从 mysql 获取解析成 time.Time 为 2018-12-24 10:00:00 UTC

以上问题可以通过设置loc=Local解决
db, err := sql.Open("mysql", "user:password@/dbname?charset=utf8mb4&parseTime=true&loc=Local")
golang 程序里 time.Time 为 2018-12-24 18:00:00 CST

转为 UTC 存储到 mysql 2018-12-24 18:00:00

golang 从 mysql 获取解析成 time.Time 为 2018-12-24 18:00:00 CST

但要注意，这个方法不会修改连接的time_zone属性，而是在go-sql-driver程序里对 time.Time 做了时区转换，所以会遗留一个问题，如果在我们的程序里执行 SQL 语句，UNIX_TIMESTAMP(NOW()) - UNIX_TIMESTAMP(created_at)，会出现意想不到的问题，因为NOW()取到的仍然是 UTC 时间，而created_at为 GST 时间

这个问题可以通过同时指定loc=true&time_zone=*来解决
需要用url.QueryEscape对timezone进行编码，单引号不可省略
但这种方法我们需要手动配置时区名，这样并不方便，所以建议避免直接调用 MYSQL 的时间函数

timezone := "'Asia/Shanghai'"
db, err := sql.Open("mysql", "user:password@/dbname?charset=utf8mb4&parseTime=true&loc=Local&time_zone=" + url.QueryEscape(timezone))

https://studygolang.com/articles/17313?fr=sidebar