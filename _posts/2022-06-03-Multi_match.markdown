---
title: Multi_match
layout: post
category: elasticsearch
author: 夏泽民
---
多字段检索，是组合查询的另一种形态，考试的时候如果考察多字段检索，并不一定必须使用multi_match，使用bool query

GET <index>/_search
{
  "query": {
    "multi_match": {
      "query": "<query keyword>",
      "type": "<multi_match_type>",
      "fields": [
        "<field_a>",
        "<field_b>"
      ]
    }
  }
}
<!-- more -->
3.1 best_fields：
3.1.1 概念：
侧重于字段维度，单个字段的得分权重大，对于同一个query，单个field匹配更多的term，则优先排序。

注意，best_fields是multi_match中type的默认值

在best_fields策略中给其他剩余字段设置的权重值，取值范围 [0,1]，其中 0 代表使用 dis_max 最佳匹配语句的普通逻辑，1表示所有匹配语句同等重要。最佳的精确值需要根据数据与查询调试得出，但是合理值应该与零接近（处于 0.1 - 0.4 之间），这样就不会颠覆 dis_max （Disjunction Max Query）最佳匹配性质的根本。

3.2 most_fields：
3.2.1 概念
侧重于查询维度，单个查询条件的得分权重大，如果一次请求中，对于同一个doc，匹配到某个term的field越多，则越优先排序


3.3 cross_fields:
注意：理解cross_fields的概念之前，需要对ES的评分规则有基本的了解，戳：评分，学习ES评分的基本原理

3.3.1 概念
将任何与任一查询匹配的文档作为结果返回，但只将最佳匹配的评分作为查询的评分结果返回

评分基本规则：

词频（TF term frequency ）：关键词在每个doc中出现的次数，词频越高，评分越高
反词频（ IDF inverse doc
frequency）：关键词在整个索引中出现的次数，反词频越高，评分越低
每个doc的长度，越长相关度评分越低


# 吴 必须包含在 name.姓 或者 name.名 里
# 并且
# 磊 必须包含在 name.姓 或者 name.名 里
GET teacher/_search
{
  "query": {
    "multi_match" : {
      "query":      "吴磊",
      "type":       "cross_fields",
      "fields":     [ "name.姓", "name.名" ],
      "operator":   "and"
    }
  }
}


https://blog.csdn.net/wlei0618/article/details/120451249
