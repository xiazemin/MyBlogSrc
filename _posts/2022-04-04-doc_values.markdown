---
title: doc_values
layout: post
category: elasticsearch
author: 夏泽民
---
在es5.6版本的时候，使用聚合语句查询es的时候，会出现异常，关键信息如下：
can't load fielddata on because fielddata is unsupported on fields of type xxx, use doc value instead...
主要是因为，es 索引的mapping里面的doc_values属性的值设置的是false。
这个属性是干啥的呢？快速了解一下。
字段的 doc_values 属性有两个值， true、false。默认为 true ，即开启。
当 doc_values 为 fasle 时，无法基于该字段排序、聚合、在脚本中访问字段值。
当 doc_values 为 true 时，ES 会增加一个相应的正排索引，这增加的磁盘占用，也会导致索引数据速度慢一些。
Doc Values 默认对所有字段启用，除了 analyzed strings 可分词的字段外。也就是说所有的数字、地理坐标、日期、IP 和不可分词的（not_analyzed）字符类型都会默认开启。
为了加快排序、聚合操作，在建立倒排索引的时候，额外增加一个列式存储映射，是一个空间换时间的做法。默认是开启的，对于确定不需要聚合或者排序的字段可以关闭。
在ES保持文档,构建倒排索引的同时doc_values就被生成了, doc_values数据太大时, 它存储在电脑磁盘上.
doc_values是列式存储结构, 它擅长做聚合和排序
对于非分词字段, doc_values默认值是true(开启的), 如果确定某字段不参与聚合和排序,可以把该字段的doc_values设为false
例如SessionID, 它是keyword类型, 对它聚合或排序毫无意义, 需要把doc_values设为false, 节约磁盘空间
分词字段不能用doc_values
通过设置 doc_values: false ，这个字段将不能被用于聚合、排序以及脚本操作
那，解决报错的方法就是修改这个索引的mapping，把这个属性=false的给删掉。或者给改成true
<!-- more -->
https://blog.csdn.net/qq_27093465/article/details/114639346


当被分词的field数据被加载，会一直贮存在内存，直到被驱逐或节点崩溃。

Fielddata is loaded lazily. If you never aggregate on an analyzed string, you’ll never load fielddata into memory. Furthermore, fielddata is loaded on a per-field basis, meaning only actively used fields will incur the “fielddata tax”.

FiledData是懒加载的，只有当你使用被分词字段查询或聚合时才会被载入内存。fileddata是基于field的，意味着只有常用的fields才会引起 “field tax”。

https://blog.csdn.net/u012246455/article/details/78805874

Fielddata针对text字段在默认时是禁用的
Fielddata会占用大量堆空间，尤其是在加载大量的文本字段时。 一旦将字段数据加载到堆中，它在该段的生命周期内将一直保留在那里。 同样，加载字段数据是一个昂贵的过程，可能导致用户遇到延迟的情况。 这就是默认情况下禁用字段数据的原因。

https://www.cnblogs.com/sanduzxcvbnm/p/12092298.html

1.fielddata含义
当ES进行排序（sort），统计（aggs）时，ES把涉及到的字段数据全部读取到内存（JVM Heap）中进行操作。相当于进行了数据缓存，提升查询效率。
所以fielddata是延迟加载的，在加载的时候是这个字段所有的字段都要加载。
ES中利用fielddata这个正排索引，即从文档到item，来加快统计排序等操作，fielddata实际存储方式为列式存储。

https://www.jianshu.com/p/49260d54beaf?clicktime=1574594663
