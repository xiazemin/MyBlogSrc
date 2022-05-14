---
title: query_string
layout: post
category: elasticsearch
author: 夏泽民
---
text字段和keyword字段的区别
说明text类型的字段会被分词，查询的时候如果用拆开查可以查询的到，但是要是直接全部查，就是查询不到。
所以将字段设置成keyword的时候查询的时候已有的值不会被分词。
注意“1, 2”会被拆分成[1, 2]，但是"1,2"是不拆分的，少了个空格。

match和term的区别
1.term

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

query_string查询key类型的字段，试过了，无法查询。
query_string查询text类型的字段。
和match_phrase区别的是，不需要连续，顺序还可以调换。
<!-- more -->
https://www.cnblogs.com/chenmz1995/p/10199147.html
