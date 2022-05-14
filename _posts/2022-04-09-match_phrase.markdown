---
title: match_phrase
layout: post
category: elasticsearch
author: 夏泽民
---
es数据库text类型和keyword类型数据中match、match_phrase、query_string、term之间区别

将字段设置成keyword的时候查询的时候已有的值不会被分词
1）term查询keyword字段。

 term不会分词。而keyword字段也不分词。需要完全匹配才可。

2）term查询text字段。

因为text字段会分词，而term不分词，所以term查询的条件必须是text字段分词后的某一个。

1）match查询keyword字段

match会被分词，而keyword不会被分词，match的需要跟keyword的完全匹配可以。

2）match查询text字段

match分词，text也分词，只要match的分词结果和text的分词结果有相同的就匹配。

1）match_phrase匹配keyword字段。

这个同上必须跟keywork一致才可以。

2）match_phrase匹配text字段。

match_phrase是分词的，text也是分词的。match_phrase的分词结果必须在text字段分词中都包含，而且顺序必须相同，而且必须都是连续的。
<!-- more -->
1）query_string查询key类型的字段，试过了，无法查询。

2）query_string查询text类型的字段。

和match_phrase区别的是，不需要连续，顺序还可以调换。

https://blog.csdn.net/java173842219/article/details/117745993

https://zhuanlan.zhihu.com/p/142641300

https://mp.weixin.qq.com/s?__biz=MzIxMjE3NjYwOQ==&mid=2247483805&idx=1&sn=d890e383339c2fe5bf265b2d53a8a270&chksm=974b5a13a03cd30558de75c891fc9d471f3477981f49b53c8bb26a84256d9e986e97f5d33808&scene=21#wechat_redirect

https://mp.weixin.qq.com/s?__biz=MzIxMjE3NjYwOQ==&mid=2247483699&idx=1&sn=36735645bc6983d9f0229ece49d4ac7b&chksm=974b5abda03cd3abbe925c4e77cb899f12b9ca1eedab420e64336a358ddb46c0f253a11b4409&scene=21#wechat_redirect

https://mp.weixin.qq.com/s?__biz=MzIxMjE3NjYwOQ==&mid=2247483734&idx=1&sn=dac2e9f092303b57314f8744a82fb9ff&chksm=974b5ad8a03cd3ce7be4ff7b4e942cf57645c773b54c855fcd3777755d95f5fb50dad0aab013&scene=21#wechat_redirect

https://mp.weixin.qq.com/s?__biz=MzIxMjE3NjYwOQ==&mid=2247483805&idx=1&sn=d890e383339c2fe5bf265b2d53a8a270&chksm=974b5a13a03cd30558de75c891fc9d471f3477981f49b53c8bb26a84256d9e986e97f5d33808&scene=21#wechat_redirect

https://mp.weixin.qq.com/s?__biz=MzIxMjE3NjYwOQ==&mid=2247483825&idx=1&sn=0b294bd614ed4504577b5d155706fd5f&chksm=974b5a3fa03cd329c9393f184299b430728118a7005e236cce049e1af9d31f8c75037b80a1f9&scene=21#wechat_redirect
可以使用_catAPI进行集群健康检查：

API格式：

GET /_cat/health?v
复制


https://www.qikegu.com/docs/3053

