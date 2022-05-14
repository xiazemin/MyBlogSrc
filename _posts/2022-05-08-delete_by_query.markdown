---
title: delete_by_query
layout: post
category: elasticsearch
author: 夏泽民
---
curl -X POST "localhost:9200/twitter/_delete_by_query" -H 'Content-Type: application/json' -d'
{
  "query": { 
    "match": {
      "name": "测试删除"
    }
  }
}
'
<!-- more -->
查询必须是有效的键值对，query是键，这和Search API是同样的方式。在search api中q参数和上面效果是一样的。

返回数据格式，告诉你用时和删除多少数据等

在执行_delete_by_query期间，为了删除匹配到的所有文档，多个搜索请求是按顺序执行的。每次找到一批文档时，将会执行相应的批处理请求来删除找到的全部文档。如果搜索或者批处理请求被拒绝，_delete_by_query根据默认策略对被拒绝的请求进行重试（最多10次）。达到最大重试次数后，会造成_delete_by_query请求中止，并且会在failures字段中响应 所有的故障。已经删除的仍会执行。换句话说，该过程没有回滚，只有中断。
在第一个请求失败引起中断，失败的批处理请求的所有故障信息都会记录在failures元素中；并返回回去。因此，会有不少失败的请求。
如果你想计算有多少个版本冲突，而不是中止，可以在URL中设置为conflicts=proceed或者在请求体中设置"conflicts": "proceed"。

curl -X POST "localhost:9200/twitter/_doc/_delete_by_query?conflicts=proceed" -H 'Content-Type: application/json' -d'
{
  "query": {
    "match_all": {}
  }
}
'

https://www.cnblogs.com/wangmo/p/11008928.html
