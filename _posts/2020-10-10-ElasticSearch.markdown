---
title: ElasticSearch
layout: post
category: algorithm
author: 夏泽民
---
https://www.zhihu.com/question/41109030

1）数据量每天5G数据量，按存储1年的数据来考虑，大概是1.8T，es和hbase都能支持，并且两者都可以通过扩展集群来加大可存储的数据量。随着数据量的增加，es的读写性能会有所下降，从存储原始数据的角度来看，hbase要优于es2）数据更新es数据更新是对文档进行更新，需要先将es中的数据取出，设置更新字段后再写入es。hbase是列存储的，可以方便地更新任意字段的值。3）查询复杂度hbase支持简单的行、列或范围查询，若没有对查询字段做二级索引的话会引发扫全表操作，性能较差。而ES提供了丰富的查询语法，支持对多种类型的精确匹配、模糊匹配、范围查询、聚合等操作，ES对字段做了反向索引，即使在亿级数据量下还可以达到秒级的查询响应速度。4）字段扩展性hbase和es都对非结构化数据存储提供了良好的支持。es可以通过动态字段方便地对字段进行扩展，而hbase本身就是基于列存储的，可以很方便地添加qualifier来实现字段的扩展
<!-- more -->
https://zhuanlan.zhihu.com/p/88982181

https://www.jianshu.com/p/4e412f48e820

https://blog.csdn.net/sdksdk0/article/details/53966430

https://my.oschina.net/jiangbianwanghai/blog/779190

https://segmentfault.com/a/1190000018071516

https://www.xiaoheidiannao.com/212647.html

https://cloud.tencent.com/developer/article/1467461