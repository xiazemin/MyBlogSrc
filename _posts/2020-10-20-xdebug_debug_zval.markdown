---
title: xdebug_debug_zval
layout: post
category: php
author: 夏泽民
---
众所周知php的变量都是通过zend引擎来处理的 而zval结构体则是管理我们变量的一个容器
而 xdebug_debug_zval 函数则是我们调试 zval结构体的一个好工具

<?php

//php zval变量容器

$a = range(0, 3);

xdebug_debug_zval('a');
上面的代码 在浏览器中 会输出 以下结果

a: 
(refcount=1, is_ref=0),
array (size=4)
  0 => (refcount=1, is_ref=0),int 0
  1 => (refcount=1, is_ref=0),int 1
  2 => (refcount=1, is_ref=0),int 2
  3 => (refcount=1, is_ref=0),int 3

refcount代表的是 有多少个变量指向这一个内存空间 is_ref代表的是是不是引用

https://www.php.net/manual/zh/features.gc.refcounting-basics.php
<!-- more -->
概述

    在5.2及更早版本的PHP中，没有专门的垃圾回收器GC（Garbage Collection），引擎在判断一个变量空间是否能够被释放的时候是依据这个变量的zval的refcount的值，如果refcount为0，那么变量的空间可以被释放，否则就不释放，这是一种非常简单的GC实现。然而在这种简单的GC实现方案中，出现了意想不到的变量内存泄漏情况（Bug:http://bugs.php.net/bug.php?id=33595），引擎将无法回收这些内存，于是在PHP5.3中出现了新的GC，新的GC有专门的机制负责清理垃圾数据，防止内存泄漏。本文将详细的阐述PHP5.3中新的GC运行机制。

    目前很少有详细的资料介绍新的GC，本文将是目前国内最为详细的从源码角度介绍PHP5.3中GC原理的文章。其中关于垃圾产生以及算法简介部分由笔者根据手册翻译而来，当然其中融入了本人的一些看法。手册中相关内容：Garbage Collection

    在介绍这个新的GC之前，读者必须先了解PHP中变量的内部存储相关知识，请先阅读 变量的内部存储：引用和计数 

 

什么算垃圾

    首先我们需要定义一下“垃圾”的概念，新的GC负责清理的垃圾是指变量的容器zval还存在，但是又没有任何变量名指向此zval。因此GC判断是否为垃圾的一个重要标准是有没有变量名指向变量容器zval。

    假设我们有一段PHP代码，使用了一个临时变量$tmp存储了一个字符串，在处理完字符串之后，就不需要这个$tmp变量了，$tmp变量对于我们来说可以算是一个“垃圾”了，但是对于GC来说，$tmp其实并不是一个垃圾，$tmp变量对我们没有意义，但是这个变量实际还存在，$tmp符号依然指向它所对应的zval，GC会认为PHP代码中可能还会使用到此变量，所以不会将其定义为垃圾。

    那么如果我们在PHP代码中使用完$tmp后，调用unset删除这个变量，那么$tmp是不是就成为一个垃圾了呢。很可惜，GC仍然不认为$tmp是一个垃圾，因为$tmp在unset之后，refcount减少1变成了0(这里假设没有别的变量和$tmp指向相同的zval),这个时候GC会直接将$tmp对应的zval的内存空间释放，$tmp和其对应的zval就根本不存在了。此时的$tmp也不是新的GC所要对付的那种“垃圾”。那么新的GC究竟要对付什么样的垃圾呢，下面我们将生产一个这样的垃圾。  

 

顽固垃圾的产生过程

    如果读者已经阅读了变量内部存储相关的内容，想必对refcount和isref这些变量内部的信息有了一定的了解。这里我们将结合手册中的一个例子来介绍垃圾的产生过程：

 

<?php

$a = "new string";

?>

在这么简单的一个代码中，$a变量内部存储信息为

a: (refcount=1, is_ref=0)='new string'

 

当把$a赋值给另外一个变量的时候，$a对应的zval的refcount会加1

<?php

$a = "new string";

$b = $a;

?>
此时$a和$b变量对应的内部存储信息为

a,b: (refcount=2, is_ref=0)='new string'

当我们用unset删除$b变量的时候，$b对应的zval的refcount会减少1

<?php

$a = "new string"; //a: (refcount=1, is_ref=0)='new string'

$b = $a;                 //a,b: (refcount=2, is_ref=0)='new string'

unset($b);              //a: (refcount=1, is_ref=0)='new string'

?>

 

对于普通的变量来说，这一切似乎很正常，但是在复合类型变量（数组和对象）中，会发生比较有意思的事情：

<?php

$a = array('meaning' => 'life', 'number' => 42);

?>

a的内部存储信息为:

a: (refcount=1, is_ref=0)=array (
   'meaning' => (refcount=1, is_ref=0)='life',
   'number' => (refcount=1, is_ref=0)=42
)

数组变量本身($a)在引擎内部实际上是一个哈希表，这张表中有两个zval项 meaning和number，

所以实际上那一行代码中一共生成了3个zval,这3个zval都遵循变量的引用和计数原则

https://blog.csdn.net/guoyuqi0554/article/details/8075435

http://yangxikun.github.io/php/2013/08/24/php-garbage-collection-mechanism.html

https://www.iminho.me/wiki/blog-18.html

https://www.php.net/manual/zh/features.gc.collecting-cycles.php

https://www.laruence.com/2011/03/29/1949.html

https://blog.text.wiki/2016/06/16/session-and-gc.html


由于对象进行了分代处理，因此垃圾回收区域、时间也不一样。GC有两种类型：Scavenge GC和Full GC。

1 Scavenge GC
  一般情况下，当新对象生成，并且在Eden申请空间失败时，就会触发Scavenge GC，对Eden区域进行GC，清除非存活对象，并且把尚且存活的对象移动到Survivor区。然后整理Survivor的两个区。这种方式的GC是对年轻代的Eden区进行，不会影响到年老代。因为大部分对象都是从Eden区开始的，同时Eden区不会分配的很大，所以Eden区的GC会频繁进行。因而，一般在这里需要使用速度快、效率高的算法，使Eden去能尽快空闲出来。

2 Full GC
  对整个堆进行整理，包括Young、Tenured和Perm。Full GC因为需要对整个堆进行回收，所以比Scavenge GC要慢，因此应该尽可能减少Full GC的次数。在对JVM调优的过程中，很大一部分工作就是对于Full GC的调节。有如下原因可能导致Full GC：

a) 年老代（Tenured）被写满；

b) 持久代（Perm）被写满；

c) System.gc()被显示调用；

d) 上一次GC之后Heap的各域分配策略动态变化；


正常的变量在生命周期完成之后的回收。

  这种情况也就是说当zend_value中refCount==0的时候，这时候属于正常的内存回收。
  
  垃圾回收
所谓垃圾： 就是指通过循环引用（自己引用自己，目前只在array,object类型中有出现）的形式而导致refcount永远不为0。这种情况下，如果不处理，但是这些内存无法释放，到时内存泄露
https://segmentfault.com/a/1190000013971525