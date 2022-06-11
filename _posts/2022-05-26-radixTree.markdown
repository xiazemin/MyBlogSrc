---
title: radixTree
layout: post
category: algorithm
author: 夏泽民
---
基数树是一种比较节省空间的树结构，下图展示了基数树的结构，其中key是树的构建方式，在这里，key是一个32位的整数，为了避免层数过深，所以使用两位代表子节点的索引，基数树就是依据二进制串来生成树结构。值value被储存在叶节点。
<!-- more -->
https://zhuanlan.zhihu.com/p/95814705
http://www.360doc.com/content/19/0305/18/496343_819431105.shtml

Linux基数树（radix tree）是将long整数键值与指针相关联的机制，它存储有效率。而且可高速查询，用于整数值与指针的映射（如：IDR机制）、内存管理等。
IDR（ID Radix）机制是将对象的身份鉴别号整数值ID与对象指针建立关联表。完毕从ID与指针之间的相互转换。
IDR机制使用radix树状结构作为由id进行索引获取指针的稀疏数组，通过使用位图能够高速分配新的ID，IDR机制避免了使用固定尺寸的数组存放指针。IDR机制的API函数在lib/idr.c中实现。


Linux radix树最广泛的用途是用于内存管理。结构address_space通过radix树跟踪绑定到地址映射上的核心页，该radix树同意内存管理代码高速查找标识为dirty或writeback的页。


https://www.cnblogs.com/wgwyanfs/p/6887889.html