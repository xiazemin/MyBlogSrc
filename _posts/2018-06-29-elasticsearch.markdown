---
title: elasticsearch
layout: post
category: elasticsearch
author: 夏泽民
---
插入一条数据
curl -XPUT 'localhost:9200/shakespeare?pretty' -H 'Content-Type: application/json' -d'
 {
  "mappings": {
   "doc": {
    "properties": {
     "speaker": {"type": "keyword"},
     "play_name": {"type": "keyword"},
     "line_id": {"type": "integer"},
     "speech_number": {"type": "integer"}
    }
   }
  }
 }
'
插入文件
curl -H 'Content-Type: application/x-ndjson' -XPOST 'localhost:9200/shakespeare/doc/_bulk?pretty' --data-binary @shakespeare_6.0.json

查看插入是否成功
curl -XGET 'localhost:9200/_cat/indices?v&pretty'
health status index                         uuid                   pri rep docs.count docs.deleted store.size pri.store.size
green  open   .monitoring-es-6-2018.06.29   PUGr8fC1QbGnlxDoVeAxkQ   1   0       7906           45      4.8mb          4.8mb
green  open   .monitoring-alerts-6          R_QVVTIqToymAC8qwjCT9w   1   0          6            0     35.6kb         35.6kb
green  open   .watcher-history-7-2018.06.29 zmHqjXsPR1q0Y-DkNNcO2Q   1   0        944            0      1.5mb          1.5mb
yellow open   shakespeare                   di23ZjwlSa-gIHVsgBOktw   5   1          0            0      1.1kb          1.1kb
green  open   .watches                      h3Krrmk4T82wqZ57UTJnyg   1   0          6            0     58.5kb         58.5kb
green  open   .triggered_watches            btt3Gxx7TGOndfCFu2Fg6A   1   0          0            0     33.1kb         33.1kb
进入kibana  http://localhost:5601/ 
点击：Set up index patterns

No default index pattern. You must select or create one to continue.

勾选Include system indices
Index pattern
输入：＊
next

Time Filter field name
I don‘t want the Time Filter

custom index id
shakespeare

create index patton
	<img src="{{site.url}}{{site.baseurl}}/img/pattern.png"/>
<!-- more -->

