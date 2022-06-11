---
title: too_many_clauses
layout: post
category: elasticsearch
author: 夏泽民
---
报错 too_many_clauses: maxClauseCount is set to 1024
<!-- more -->
报错原因是Search限制一个bool查询中最多只能有1024个值或子查询，当超过1024时，会抛出次异常。

详细说明
解决方案：

调整查询
当超过1024时可以将一个bool查询拆成两个子bool查询，使用must关键字，使得两个子bool查询是与的关系

增大限制
登录manager，添加参数
为Search添加自定义参数，indices.query.bool.max_clause_count

将它的值设置为10240，配置到elasticsearch.yml中，那么bool查询下的子段数量的限制将扩大到10240，基本可以满足正常使用。如果不够大，可以继续增大，但是一个bool查询中最好不要太多的子查询，会有一些性能损失，可以按照方案一做拆分。

https://blog.csdn.net/majixiang1996/article/details/105240614

https://github.com/alibaba/sentinel-golang
https://github.com/sony/gobreaker
