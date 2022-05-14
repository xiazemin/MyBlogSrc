---
title: prometheus
layout: post
category: linux
author: 夏泽民
---
https://github.com/prometheus/prometheus
Prometheus 是由 SoundCloud 开源监控告警解决方案。

prometheus
prometheus存储的是时序数据，即按相同时序(相同名称和标签)，以时间维度存储连续的数据的集合。

时序(time series)是由名字(Metric)以及一组key/value标签定义的，具有相同的名字以及标签属于相同时序。
<!-- more -->
PromQL
PromQL (Prometheus Query Language) 是 Prometheus 自己开发的数据查询 DSL 语言。

查询结果类型：

瞬时数据 (Instant vector): 包含一组时序，每个时序只有一个点，例如：http_requests_total

区间数据 (Range vector): 包含一组时序，每个时序有多个点，例如：http_requests_total[5m]

纯量数据 (Scalar): 纯量只有一个数字，没有时序，例如：count(http_requests_total)

https://www.jianshu.com/p/93c840025f01

https://baijiahao.baidu.com/s?id=1643112707938591118&wfr=spider&for=pc

高维度数据模型

自定义查询语言

可视化数据展示

高效的存储策略

易于运维

提供各种客户端开发库

警告和报警

数据导出

https://www.kubernetes.org.cn/tags/prometheus

http://www.eryajf.net/2468.html
