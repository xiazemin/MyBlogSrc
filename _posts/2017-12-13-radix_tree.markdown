---
title: radix tree
layout: post
category: linux
author: 夏泽民
---

基数树

对于长整型数据的映射，如何解决Hash冲突和Hash表大小的设计是一个很头疼的问题。
radix树就是针对这种稀疏的长整型数据查找，能快速且节省空间地完成映射。借助于Radix树，我们可以实现对于长整型数据类型的路由。利用radix树可以根据一个长整型（比如一个长ID）快速查找到其对应的对象指针。这比用hash映射来的简单，也更节省空间，使用hash映射hash函数难以设计，不恰当的hash函数可能增大冲突，或浪费空间。

radix tree是一种多叉搜索树，树的叶子结点是实际的数据条目。每个结点有一个固定的、2^n指针指向子结点（每个指针称为槽slot，n为划分的基的大小）

Radix树在Linux中的应用：
Linux基数树（radix tree）是将long整数键值与指针相关联的机制，它存储有效率，并且可快速查询，用于整数值与指针的映射（如：IDR机制）、内存管理等。
IDR（ID Radix）机制是将对象的身份鉴别号整数值ID与对象指针建立关联表，完成从ID与指针之间的相互转换。IDR机制使用radix树状结构作为由id进行索引获取指针的稀疏数组，通过使用位图可以快速分配新的ID，IDR机制避免了使用固定尺寸的数组存放指针。IDR机制的API函数在lib/idr.c中实现。

Linux radix树最广泛的用途是用于内存管理，结构address_space通过radix树跟踪绑定到地址映射上的核心页，该radix树允许内存管理代码快速查找标识为dirty或writeback的页。其使用的是数据类型unsigned long的固定长度输入的版本。每级代表了输入空间固定位数。Linux radix树的API函数在lib/radix-tree.c中实现。（把页指针和描述页状态的结构映射起来，使能快速查询一个页的信息。）

Linux内核利用radix树在文件内偏移快速定位文件缓存页。 
Linux(2.6.7) 内核中的分叉为 64(2^6)，树高为 6(64位系统)或者 11(32位系统)，用来快速定位 32 位或者 64 位偏移，radix tree 中的每一个叶子节点指向文件内相应偏移所对应的Cache项。

【radix树为稀疏树提供了有效的存储，代替固定尺寸数组提供了键值到指针的快速查找。】 

radix树概述

radix树是通用的字典类型数据结构，radix树又称为PAT位树（Patricia Trie or crit bit tree）。Linux内核使用了数据类型unsigned long的固定长度输入的版本。每级代表了输入空间固定位数。
radix tree是一种多叉搜索树，树的叶子结点是实际的数据条目。每个结点有一个固定的、2^n指针指向子结点（每个指针称为槽slot），并有一个指针指向父结点。
Linux内核利用radix树在文件内偏移快速定位文件缓存页，图4是一个radix树样例，该radix树的分叉为4(22)，树高为4，树的每个叶子结点用来快速定位8位文件内偏移，可以定位4x4x4x4=256页，如：图中虚线对应的两个叶子结点的路径组成值0x00000010和0x11111010，指向文件内相应偏移所对应的缓存页
<!-- more -->
<div class="container">
	<div class="row">
	<img src="{{site.url}}{{site.baseurl}}/img/radix_tree.jpg"/>
	</div>
</div>
Linux radix树每个结点有64个slot，与数据类型long的位数相同，图1显示了一个有3级结点的radix树，每个数据条目（item）可用3个6位的键值（key）进行索引，键值从左到右分别代表第1~3层结点位置。没有孩子的结点在图中不出现。因此，radix树为稀疏树提供了有效的存储，代替固定尺寸数组提供了键值到指针的快速查找。 
radix树slot数

Linux内核根用户配置将树的slot数定义为4或6，即每个结点有16或64个slot，如图2所示，当树高为1时，64个slot对应64个页，当树高为2时，对应64*64个页。 
