---
title: tokenlizer
layout: post
category: elasticsearch
author: 夏泽民
---
Elasticsearch 的内置分析器都是全局可用的，不需要提前配置，它们也可以在字段映射中直接指定在某字段

配置语言分析器
语言分析器都不需要任何配置，开箱即用， 它们中的大多数都允许你控制它们的各方面行为，具体来说：

词干提取排除
想象下某个场景，用户们想要搜索 World Health Organization 的结果， 但是却被替换为搜索 organ health 的结果。有这个困惑是因为 organ 和 organization 有相同的词根： organ 。 通常这不是什么问题，但是在一些特殊的文档中就会导致有歧义的结果，所以我们希望防止单词 organization 和 organizations 被缩减为词干。

自定义停用词
英语中默认的停用词列表如下：

a, an, and, are, as, at, be, but, by, for, if, in, into, is, it,
no, not, of, on, or, such, that, the, their, then, there, these,
they, this, to, was, will, with
关于单词 no 和 not 有点特别，这俩词会反转跟在它们后面的词汇的含义。或许我们应该认为这两个词很重要，不应该把他们看成停用词

https://www.cntofu.com/book/40/language-intro.html
<!-- more -->
https://www.elastic.co/guide/en/elasticsearch/reference/current/analysis-lang-analyzer.html

https://www.elastic.co/guide/en/elasticsearch/reference/current/analysis-lang-analyzer.html#arabic-analyzer

https://github.com/msarhan/lucene-arabic-analyzer

https://xiaoxiami.gitbook.io/elasticsearch/ji-chu/33-analysisfen-679029/333analyzersfen-xi-566829/language-analyzersff08-yu-yan-fen-xi-qi-ff09

系统默认分词器：
1、standard 分词器
中文分词器：
1、ik-analyzer
采用了多子处理器分析模式，支持：英文字母、数字、中文词汇等分词处理，兼容韩文、日文字符

优化的词典存储，更小的内存占用。支持用户词典扩展定义

https://segmentfault.com/a/1190000012553894
Analysis 和 Analyzer
Analysis： 文本分析是把全文本转换一系列单词(term/token)的过程，也叫分词。Analysis是通过Analyzer来实现的。

当一个文档被索引时，每个Field都可能会创建一个倒排索引（Mapping可以设置不索引该Field）。

倒排索引的过程就是将文档通过Analyzer分成一个一个的Term,每一个Term都指向包含这个Term的文档集合。

当查询query时，Elasticsearch会根据搜索类型决定是否对query进行analyze，然后和倒排索引中的term进行相关性查询，匹配相应的文档。

2 、Analyzer组成
分析器（analyzer）都由三种构件块组成的：character filters ， tokenizers ， token filters。

https://www.cnblogs.com/qdhxhz/p/11585639.html

https://github.com/ARBML/tkseem

https://www.elastic.co/guide/en/elasticsearch/reference/8.0/specify-analyzer.html#specify-index-field-analyzer


