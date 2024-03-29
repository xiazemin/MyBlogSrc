---
title: Normalization
layout: post
category: elasticsearch
author: 夏泽民
---
自定义文本分析
文本分析由分析器执行，分析器是一组控制整个过程的规则。Elasticsearch 包含一个默认分析器standard analyzer，适用于大多数开箱即用的用例。

自定义分析器可让您控制分析过程的每个步骤，包括：

分词前对文本的更改
文本如何转换为Token
在索引或搜索前，对Token所做的正规化处理
二、相关概念
分析器
无论内置的还是自定义的分析器（analyzer）都由三部分组成，即

字符过滤器（character filters）：接收字符流，可以添加、移除或改变字符来转换流
例如：去除字符里的HTML元素
分词器（tokenizers）：接收处理后的字符流，分解为单独的 token（通常是单个单词）。
例如：whitespace分词器会在遇到空格时对文本进行拆分；
分词器还负责记录每个term的顺序或位置以及该term所代表的原始单词的开始和结束字符偏移量；
一个分析器有且只有一个分词器。
token过滤器（token filters）：token过滤器在遇到token流时，可以添加、删除或改变Token
例如，lowercase令牌过滤器将所有令牌转换为小写；stop令牌过滤器从令牌流中删除常用词（停用词）the；synonym令牌过滤器将同义词引入令牌流。
Token过滤器不允许更改每个Token的位置或字符偏移量
分析器可能有零个或多个 token过滤器，按序生效
内置分析器将这些构建块预先打包成适用于不同语言和文本类型的分析器。Elasticsearch 还公开了各个构建块，以便将它们组合起来定义新的自定义分析器。
<!-- more -->
https://www.codetd.com/en/article/13381953

Elasticsearch这种全文搜索引擎，会用某种算法对建立的文档进行分析，从文档中提取出有效信息（Token）

对于es来说，有内置的分析器（Analyzer）和分词器（Tokenizer）

1：分析器
ES内置分析器

standard	分析器划分文本是通过词语来界定的，由Unicode文本分割算法定义。它删除大多数标点符号，将词语转换为小写(就是按照空格进行分词)
simple	分析器每当遇到不是字母的字符时，将文本分割为词语。它将所有词语转换为小写。
keyword	可以接受任何给定的文本，并输出与单个词语相同的文本
pattern	分析器使用正则表达式将文本拆分为词语，它支持小写和停止字
language	语言分析器
whitespace	（空白）分析器每当遇到任何空白字符时，都将文本划分为词语。它不会将词语转换为小写
custom	自定义分析器

https://www.cnblogs.com/niutao/p/10909147.html
https://chowdera.com/2021/12/202112131732102543.html
https://www.elastic.co/guide/cn/elasticsearch/guide/current/standard-tokenizer.html
2：更新分析器
1：要先关闭索引
​
2：添加分析器
​
3：打开索引

3：分词器
Es中也支持非常多的分词器

Standard	默认的分词器根据 Unicode 文本分割算法，以单词边界分割文本。它删除大多数标点符号。<br/>它是大多数语言的最佳选择
Letter	遇到非字母时分割文本
Lowercase	类似 letter ，遇到非字母时分割文本，同时会将所有分割后的词元转为小写
Whitespace	遇到空白字符时分割位文本

https://www.shuzhiduo.com/A/n2d9yv00dD/

https://unicode.org/reports/tr29/
https://doc.codingdict.com/elasticsearch/347/
