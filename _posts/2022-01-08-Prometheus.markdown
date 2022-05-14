---
title: Prometheus
layout: post
category: k8s
author: 夏泽民
---
Histogram 由 <basename>_bucket{le="<upper inclusive bound>"}，
<basename>_bucket{le="+Inf"}, <basename>_sum，<basename>_count组成，
例如 Prometheus server中prometheus_local_storage_series_chunks_persisted, 
表示 Prometheus 中每个时序需要存储的 chunks 数量，我们可以用它计算待持久化的数据的分位数。


Summary 和 Histogram 类似，由 <basename>{quantile="<φ>"}，<basename>_sum，
<basename>_count 组成，
例如 Prometheus server 中 prometheus_target_interval_length_seconds。

Histogram vs Summary
都包含 <basename>_sum，<basename>_count
Histogram 需要通过 <basename>_bucket 计算 quantile, 
而 Summary 直接存储了 quantile 的值。
<!-- more -->
PromQL 查询结果主要有 3 种类型：
瞬时数据 (Instant vector): 包含一组时序，每个时序只有一个点，例如:
http_requests_total
区间数据 (Range vector): 包含一组时序，每个时序有多个点，例如：
http_requests_total[5m]
纯量数据 (Scalar): 纯量只有一个数字，没有时序，例如：
count(http_requests_total)

假设采样数据 metric 叫做 x, 如果 x 是 histogram 或 summary 类型必需满足以下条件：

采样数据的总和应表示为 x_sum。
采样数据的总量应表示为 x_count。
summary 类型的采样数据的 quantile 应表示为 x{quantile=”y”}。
histogram 类型的采样分区统计数据将表示为 x_bucket{le=”y”}。、
histogram 类型的采样必须包含 x_bucket{le=”+Inf”}, 它的值等于 x_count 的值。
summary 和 historam 中 quantile 和 le 必需按从小到大顺序排列。

https://blog.csdn.net/weixin_39843367/article/details/81777209
