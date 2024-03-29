---
title: SearchType
layout: post
category: elasticSearch
author: 夏泽民
---
ES的搜索分scatter/gather两个步骤：

scatter阶段:客户端向10个分片发起搜索请求；

gather阶段:10个分片完成搜索，符合条件的结果返回；

客户端，将返回的结果进行重新拍下和排名，最后返回给用户。
<!-- more -->
搜索面临的两个问题：

问题1：假如集群一个节点有10个分片，词语“土豆”在分片的相关性从分片0到分片9依次降低（即分片0存储词语“土豆”相关性最大，分片9相关性最小），如果搜索词语“土豆”需要10个分片的计算；

问题2：如果针对特定的分片进行搜索，因为词语“土豆”在每个分片的相关性不一致，可能返回的结果也存在偏差。

针对以上问题，ES给允许设置search_type来解决上述问题

SearchType共四种类型：

1、query and fetch

向索引的所有分片（shard）都发出查询请求，各分片返回的时候把元素文档（document）和计算后的排名信息一起返回。这种搜索方式是最快的。因为相比下面的几种搜索方式，这种查询方法只需要去shard查询一次。但是各个shard返回的结果的数量之和可能是用户要求的size的n倍。

2、query then fetch（默认的搜索方式）

如果你搜索时，没有指定搜索方式，就是使用的这种搜索方式。这种搜索方式，大概分两个步骤，第一步，先向所有的shard发出请求，各分片只返回排序和排名相关的信息（注意，不包括文档document)，然后按照各分片返回的分数进行重新排序和排名，取前size个文档。然后进行第二步，去相关的shard取document。这种方式返回的document与用户要求的size是相等的。

3、DFS query and fetch

这种方式比第一种方式多了一个初始化散发(initial scatter)计算全局词频（term frequencies）步骤，有这一步，据说可以更精确控制搜索打分和排名。先对所有分片发送请求， 把所有分片中的词频和文档频率等打分依据全部汇总到一块， 再执行后面的操作。优点很明显，数据量是准确并且排名也准确，但性能是最差的。

4、DFS query then fetch

比第2种方式多了一个初始化散发(initial scatter)计算全局词频（term frequencies）步骤，过程与上一种类似，优点是排名准确，但返回的数据量不准确，可能返回(N*分片数量)的数据

https://blog.csdn.net/HuoqilinHeiqiji/article/details/103460430

DSF是什么缩写？初始化散发是一个什么样的过程？

从es的官方网站我们可以指定，初始化散发其实就是在进行真正的查询之前，先把各个分片的词频率和文档频率收集一下，然后进行词搜索的时候，各分片依据全局的词频率和文档频率进行搜索和排名。显然如果使用DFS_QUERY_THEN_FETCH这种查询方式，效率是最低的，因为一个搜索，可能要请求3次分片。但，使用DFS方法，搜索精度应该是最高的。

https://www.cnblogs.com/ningskyer/articles/5984346.html

https://blog.csdn.net/caipeichao2/article/details/46418413/