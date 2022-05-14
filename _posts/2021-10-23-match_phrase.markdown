---
title: match_phrase和term区别
layout: post
category: elasticsearch
author: 夏泽民
---
term是将传入的文本原封不动地（不分词）拿去查询。
match会对输入进行分词处理后再去查询，部分命中的结果也会按照评分由高到低显示出来。
match_phrase是按短语查询，只有存在这个短语的文档才会被显示出来。

也就是说，term和match_phrase都可以用于精确匹配，而match用于模糊匹配。

之前我以为match_phrase不会被分词，看来理解错了，其官方解释如下：

Like the match query, the match_phrase query first analyzes the query string to produce a list of terms. It then searches for all the terms, but keeps only documents that contain all of the search terms, in the same positions relative to each other.

总结下，这段话的3个要点：

match_phrase还是分词后去搜的
目标文档需要包含分词后的所有词
目标文档还要保持这些词的相对顺序和文档中的一致
只有当这三个条件满足，才会命中文档！
<!-- more -->
https://blog.csdn.net/timothytt/article/details/86775114

https://www.jianshu.com/p/74653f4b476c

term:词项，即索引的最小单元，文本搜索时最小的匹配单元。

match查询语句，match和term查询的最大区别在于，term查询会将查询词当为词项，并在倒排索引中进行全匹配。match查询会先进行分词处理，再将解析后的词项去查询，"minimum_should_match",可以控制match的查询词中最小应该匹配的比例。

match_phrase，句子查询，和match的区别，phrase是句子，句子内部要保持信息一致，所以match_phrase查询将全匹配句子所有文字，并且保证文字之间的相对位置。es提供了slop等查询控制，给用户去调整文字间相对位置的距离。slop:1 以为着 查询词“帅哥”，可以匹配到“帅*哥”，中间可以有一个文本的距离。

boost，控制单个查询语句在整体查询语句中的权重。

bool逻辑查询，should，must，should_not, must_not,可以和match、term查询进行嵌套。"minimum_should_match"，在这里也可以控制should的处理个数。（ 可以组合match bool，match_phrase来保证文本的相对位置，以及允许少匹配文本个数。）

aggs，聚合查询。强大的聚合查询，根据用户设置的桶处理条件，可以进行桶內数据的sum，min，max，统计。terms桶，会根据terms处理字段，统计桶內同一文体聚类数量。聚合查询支持，嵌套桶，时间范围的桶等。

https://zhuanlan.zhihu.com/p/29350723

https://segmentfault.com/a/1190000021520174

match_all 查询
match_all 查询简单的匹配所有文档。在没有指定查询方式时，它是默认的查询：

https://www.elastic.co/guide/cn/elasticsearch/guide/current/_most_important_queries.html