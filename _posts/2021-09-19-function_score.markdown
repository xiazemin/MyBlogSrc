---
title: es function_score
layout: post
category: storage
author: 夏泽民
---
为了使推荐更接近用户偏好，搜索时使用了function_score功能对文档进行了重新打分，改变排序规则。

es内置了几种预先定义好了的函数：

      1、weight：对每份文档适用一个简单的提升，且该提升不会被归约：当weight为2时，结果为2 * _score。


      2、field_value_factor：使用文档中某个字段的值来改变_score，比如将受欢迎程度或者投票数量考虑在内。

      3、random_score：使用一致性随机分值计算来对每个用户采用不同的结果排序方式，对相同用户仍然使用相同的排序方式。

      4、Decay Functions：衰减函数，衰减函数是利用从给定的原点到某个用户数字类型字段的值的距离的衰减进行打分的。这类似于一个范围查询，而且边缘是光滑的。

      es内部支持的衰减函数有gauss（高斯）、exp（指数）、linear（线性）

      5、 script_score：使用自定义的脚本来完全控制分值计算逻辑。
<!-- more -->
https://www.cnblogs.com/a-du/p/10755787.html

https://blog.csdn.net/lijingjingchn/article/details/106405577

https://blog.csdn.net/weixin_40341116/article/details/80913045

https://www.elastic.co/guide/cn/elasticsearch/guide/current/function-score-query.html

https://www.cnblogs.com/atomicbomb/p/9105102.html


