---
title: ClickHouse
layout: post
category: elasticsearch
author: 夏泽民
---
随着日志量不断增加，一些问题逐渐暴露出来：

一方面 ES 服务器越来越多，投入的成本越来越高。

另一方面用户的满意度不高，日志写入延迟、查询慢甚至查不出来的问题一直困扰着用户。



而从运维人员的角度看，ES 的运维成本较高，运维的压力越来越大。


为什么选择 ClickHouse



ClickHouse 是一款高性能列式分布式数据库管理系统，我们对 ClickHouse 进行了测试，发现有下列优势：

ClickHouse 写入吞吐量大，单服务器日志写入量在 50MB 到 200MB/s，每秒写入超过 60w 记录数，是 ES 的 5 倍以上。


②在 ES 中比较常见的写 Rejected 导致数据丢失、写入延迟等问题，在 ClickHouse 中不容易发生。


③查询速度快，官方宣称数据在 pagecache 中，单服务器查询速率大约在 2-30GB/s；没在 pagecache 的情况下，查询速度取决于磁盘的读取速率和数据的压缩率。经测试 ClickHouse 的查询速度比 ES 快 5-30 倍以上。


ClickHouse 比 ES 服务器成本更低：

一方面 ClickHouse 的数据压缩比比 ES 高，相同数据占用的磁盘空间只有 ES 的 1/3 到 1/30，节省了磁盘空间的同时，也能有效的减少磁盘 IO，这也是 ClickHouse 查询效率更高的原因之一。

另一方面 ClickHouse 比 ES 占用更少的内存，消耗更少的 CPU 资源。我们预估用 ClickHouse 处理日志可以将服务器成本降低一半。


④相比 ES，ClickHouse 稳定性更高，运维成本更低。



⑤ES 中不同的 Group 负载不均衡，有的 Group 负载高，会导致写 Rejected 等问题，需要人工迁移索引；在 ClickHouse 中通过集群和 Shard 策略，采用轮询写的方法，可以让数据比较均衡的分布到所有节点。


⑥ES 中一个大查询可能导致 OOM 的问题；ClickHouse 通过预设的查询限制，会查询失败，不影响整体的稳定性。



⑦ES 需要进行冷热数据分离，每天 200T 的数据搬迁，稍有不慎就会导致搬迁过程发生问题，一旦搬迁失败，热节点可能很快就会被撑爆，导致一大堆人工维护恢复的工作。


⑧ClickHouse 按天分 Partition，一般不需要考虑冷热分离，特殊场景用户确实需要冷热分离的，数据量也会小很多，ClickHouse 自带的冷热分离机制就可以很好的解决。


⑨ClickHouse 采用 SQL 语法，比 ES 的 DSL 更加简单，学习成本更低。


结合携程的日志分析场景，日志进入 ES 前已经格式化成 JSON，同一类日志有统一的 Schema，符合 ClickHouse Table 的模式。


日志查询的时候，一般按照某一维度统计数量、总量、均值等，符合 ClickHouse 面向列式存储的使用场景。


偶尔有少量的场景需要对字符串进行模糊查询，也是先经过一些条件过滤掉大量数据后，再对少量数据进行模糊匹配，ClickHouse 也能很好的胜任。


另外我们发现 90% 以上的日志没有使用 ES 的全文索引特性，因此我们决定尝试用 ClickHouse 来处理日志。
<!-- more -->
https://mp.weixin.qq.com/s/l4RgNQPxvdNIqx52LEgBnQ

https://artifacthub.io/packages/search?repo=elastic

