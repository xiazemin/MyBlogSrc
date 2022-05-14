---
title: mysql no rows in result set
layout: post
category: storage
author: 夏泽民
---
 // 查询一条记录时, 不能使用类似if err := db.QueryRow().Scan(&...); err != nil {}的处理方式
    // 因为查询单条数据时, 可能返回var ErrNoRows = errors.New("sql: no rows in result set")该种错误信息
    // 而这属于正常错误
    
    所有查询出来的字段都不允许有NULL, 避免该方式最好的办法就是建表字段时, 不要设置类似DEFAULT NULL属性
    // 还有一些无法避免的情况, 比如下面这个查询
    // 该种查询, 如果不存在, 返回值为NULL, 而非0, 针对该种简单的查询, 直接使用HAVING子句即可
    // 具体的查询, 需要在编码的过程中自行处理
    var age int32
    err = db.QueryRow(`
        SELECT
            SUM(age) age
        FROM user
        WHERE id = ?
        HAVING age <> NULL
    `, 10).Scan(&age)
    switch {
    case err == sql.ErrNoRows:
    case err != nil:
        fmt.Println(err)
    }
    fmt.Println(age)
}
<!-- more -->
https://studygolang.com/articles/9957

ifnull(sum(size),0)
0 

select sum(ifnull(size,0)) from xx
null

原因是，如果没有记录，sum里面的函数不会执行，直接返回null


注意用ifnull(max(size),0)，而不是max（ifnull(size,0)）
	前者能走索引

ifnull(max(id),0) 
max(ifnull(id,0)) 
sum(ifnull(id,0))
