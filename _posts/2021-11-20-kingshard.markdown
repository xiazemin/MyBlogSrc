---
title: kingshard shardingsphere
layout: post
category: mysql
author: 夏泽民
---
	kingshard是一个由Go开发高性能MySQL Proxy项目，kingshard在满足基本的读写分离的功能上，致力于简化MySQL分库分表操作；能够让DBA通过kingshard轻松平滑地实现MySQL数据库扩容。 kingshard的性能是直连MySQL性能的80%以上。
<!-- more -->
1. 基础功能
支持SQL读写分离。
支持透明的MySQL连接池，不必每次新建连接。
支持平滑上线DB或下线DB，前端应用无感知。
支持多个slave，slave之间通过权值进行负载均衡。
支持强制读主库。
支持主流语言（java,php,python,C/C++,Go)SDK的mysql的prepare特性。
支持到后端DB的最大连接数限制。
支持SQL日志及慢日志输出。
支持SQL黑名单机制。
支持客户端IP访问白名单机制，只有白名单中的IP才能访问kingshard（支持IP 段）。
支持字符集设置。
支持last_insert_id功能。
支持热加载配置文件，动态修改kingshard配置项（具体参考管理端命令）。
支持以Web API调用的方式管理kingshard。
支持多用户模式，不同用户之间的表是权限隔离的，互不感知。
2. sharding功能
支持按整数的hash和range分表方式。
支持按年、月、日维度的时间分表方式。
支持跨节点分表，子表可以分布在不同的节点。
支持跨节点的count,sum,max和min等聚合函数。
支持单个分表的join操作，即支持分表和另一张不分表的join操作。
支持跨节点的order by,group by,limit等操作。
支持将sql发送到特定的节点执行。
支持在单个节点上执行事务，不支持跨多节点的分布式事务。
支持非事务方式更新（insert,delete,update,replace）多个node上的子表。

https://github.com/flike/kingshard/blob/master/README_ZH.md


https://github.com/flike/kingshard

https://github.com/go-pg/sharding

https://github.com/gaoxianglong/shark/wiki

通过kingshard可以非常方便地动态迁移子表，从而保证MySQL节点的不至于负载压力太大。大致步骤如下所述：

通过自动数据迁移工具开始数据迁移。
数据差异小于某一临界值，阻塞老子表写操作（read-only）
等待新子表数据同步完毕
更改kingshard配置文件中的对应子表的路由规则。
删除老节点上的子表。

https://segmentfault.com/a/1190000003001545
https://github.com/apache/shardingsphere

mysql-proxy是官方提供的mysql中间件产品能够实现负载平衡，读写分离，failover等，但其不支持大数据量的分库分表且性能较差。下面介绍几款能代替其的mysql开源中间件产品，Atlas，cobar，tddl，让咱们看看它们各自有些什么优势和新特性吧。前端

Atlasmysql

Atlas是由 Qihoo 360, Web平台部基础架构团队开发维护的一个基于MySQL协议的数据中间层项目。它是在mysql-proxy 0.8.2版本的基础上，对其进行了优化，增长了一些新的功能特性。360内部使用Atlas运行的mysql业务，天天承载的读写请求数达几十亿条。
Altas架构：
Atlas是一个位于应用程序与MySQL之间，它实现了MySQL的客户端与服务端协议，做为服务端与应用程序通信，同时做为客户端与MySQL通信。它对应用程序屏蔽了DB的细节，同时为了下降MySQL负担，它还维护了链接池。

https://www.shangmayuan.com/a/edb384be1fd94fafbe73f5d0.html

https://www.jianshu.com/p/92a565a8eb37

https://www.cnblogs.com/firstdream/p/6768177.html

https://blog.csdn.net/gxl1989225/article/details/84761289
