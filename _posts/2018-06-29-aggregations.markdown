---
title: aggregations
layout: post
category: elasticsearch
author: 夏泽民
---
框架集合由搜索查询选择的所有数据。框架中包含许多构建块，有助于构建复杂的数据描述或摘要。聚合的基本结构如下所示 -

"aggregations" : {
   "<aggregation_name>" : {
      "<aggregation_type>" : {
         <aggregation_body>
      }

      [,"meta" : { [<meta_data_body>] } ]?
      [,"aggregations" : { [<sub_aggregation>]+ } ]?
   }
}
JSON
有以下不同类型的聚合，每个都有自己的目的 -

指标聚合
这些聚合有助于从聚合文档的字段值计算矩阵，并且某些值可以从脚本生成。
数字矩阵或者是平均聚合的单值，或者是像stats一样的多值。

平均聚合
此聚合用于获取聚合文档中存在的任何数字字段的平均值。 例如，

POST http://localhost:9200/schools/_search
请求正文

{
   "aggs":{
      "avg_fees":{"avg":{"field":"fees"}}
   }
}
基数聚合
此聚合给出特定字段的不同值的计数。 例如，

POST http://localhost:9200/schools*/_search
请求正文

{
   "aggs":{
      "distinct_name_count":{"cardinality":{"field":"name"}}
   }
}
扩展统计聚合
此聚合生成聚合文档中特定数字字段的所有统计信息。 例如，

POST http://localhost:9200/schools/school/_search
请求正文

{
   "aggs" : {
      "fees_stats" : { "extended_stats" : { "field" : "fees" } }
   }
}
最大聚合
此聚合查找聚合文档中特定数字字段的最大值。 例如，
POST http://localhost:9200/schools*/_search
请求正文

{
   "aggs" : {
      "max_fees" : { "max" : { "field" : "fees" } }
   }
}

最小聚合
此聚合查找聚合文档中特定数字字段的最小值。 例如，

POST http://localhost:9200/schools*/_search
请求正文

{
   "aggs" : {
      "min_fees" : { "min" : { "field" : "fees" } }
   }
}
总和聚合
此聚合计算聚合文档中特定数字字段的总和。 例如，

POST http://localhost:9200/schools*/_search
请求正文

{
   "aggs" : {
      "total_fees" : { "sum" : { "field" : "fees" } }
   }
}

<!-- more -->
在特殊情况下使用的一些其他度量聚合，例如地理边界聚集和用于地理位置的地理中心聚集。

桶聚合
这些聚合包含用于具有标准的不同类型的桶聚合，该标准确定文档是否属于某一个桶。桶聚合已经在下面描述 -

子聚集

此存储桶聚合会生成映射到父存储桶的文档集合。类型参数用于定义父索引。 例如，我们有一个品牌及其不同的模型，然后模型类型将有以下_parent字段 -

{
   "model" : {
      "_parent" : {
         "type" : "brand"
      }
   }
}
JSON
还有许多其他特殊的桶聚合，这在许多其他情况下是有用的，它们分别是 -

日期直方图汇总/聚合
日期范围汇总/聚合
过滤聚合
过滤器聚合
地理距离聚合
GeoHash网格聚合
全局汇总
直方图聚合
IPv4范围聚合
失踪聚合
嵌套聚合
范围聚合
反向嵌套聚合
采样器聚合
重要条款聚合
术语聚合
聚合元数据
可以通过使用元标记在请求时添加关于聚合的一些数据，并可以获得响应。 例如，

POST http://localhost:9200/school*/report/_search
请求正文

{
   "aggs" : {
      "min_fees" : { "avg" : { "field" : "fees" } ,
         "meta" :{
            "dsc" :"Lowest Fees"
         }
      }
   }
}
JSON
响应

{
   "aggregations":{"min_fees":{"meta":{"dsc":"Lowest Fees"}, "value":2180.0}}
}
