---
title: mysql buffer pool
layout: post
category: storage
author: 夏泽民
---
应用系统分层架构，为了加速数据访问，会把最常访问的数据，放在缓存(cache)里，避免每次都去访问数据库。

 

操作系统，会有缓冲池(buffer pool)机制，避免每次访问磁盘，以加速数据的访问。

 

MySQL作为一个存储系统，同样具有缓冲池(buffer pool)机制，以避免每次查询数据都进行磁盘IO。

InnoDB的缓冲池缓存什么？有什么用？

缓存表数据与索引数据，把磁盘上的数据加载到缓冲池，避免每次访问都进行磁盘IO，起到加速访问的作用。

 

速度快，那为啥不把所有数据都放到缓冲池里？

凡事都具备两面性，抛开数据易失性不说，访问快速的反面是存储容量小：

（1）缓存访问快，但容量小，数据库存储了200G数据，缓存容量可能只有64G；

（2）内存访问快，但容量小，买一台笔记本磁盘有2T，内存可能只有16G；

因此，只能把“最热”的数据放到“最近”的地方，以“最大限度”的降低磁盘访问。

 

如何管理与淘汰缓冲池，使得性能最大化呢？

 

在介绍具体细节之前，先介绍下“预读”的概念。

 

什么是预读？

磁盘读写，并不是按需读取，而是按页读取，一次至少读一页数据（一般是4K），如果未来要读取的数据就在页中，就能够省去后续的磁盘IO，提高效率。

 

预读为什么有效？

数据访问，通常都遵循“集中读写”的原则，使用一些数据，大概率会使用附近的数据，这就是所谓的“局部性原理”，它表明提前加载是有效的，确实能够减少磁盘IO。

 

按页(4K)读取，和InnoDB的缓冲池设计有啥关系？

（1）磁盘访问按页读取能够提高性能，所以缓冲池一般也是按页缓存数据；

（2）预读机制启示了我们，能把一些“可能要访问”的页提前加入缓冲池，避免未来的磁盘IO操作；

 

InnoDB是以什么算法，来管理这些缓冲页呢？

最容易想到的，就是LRU(Least recently used)。

画外音：memcache，OS都会用LRU来进行页置换管理，但MySQL的玩法并不一样。

 

传统的LRU是如何进行缓冲页管理？

 

最常见的玩法是，把入缓冲池的页放到LRU的头部，作为最近访问的元素，从而最晚被淘汰。这里又分两种情况：

（1）页已经在缓冲池里，那就只做“移至”LRU头部的动作，而没有页被淘汰；

（2）页不在缓冲池里，除了做“放入”LRU头部的动作，还要做“淘汰”LRU尾部页的动作；

传统的LRU缓冲池算法十分直观，OS，memcache等很多软件都在用，MySQL为啥这么矫情，不能直接用呢？

这里有两个问题：

（1）预读失效；

（2）缓冲池污染；

 

什么是预读失效？

由于预读(Read-Ahead)，提前把页放入了缓冲池，但最终MySQL并没有从页中读取数据，称为预读失效。

 

如何对预读失效进行优化？

要优化预读失效，思路是：

（1）让预读失败的页，停留在缓冲池LRU里的时间尽可能短；

（2）让真正被读取的页，才挪到缓冲池LRU的头部；

以保证，真正被读取的热数据留在缓冲池里的时间尽可能长。

 

具体方法是：

（1）将LRU分为两个部分：

新生代(new sublist)

老生代(old sublist)

（2）新老生代收尾相连，即：新生代的尾(tail)连接着老生代的头(head)；

（3）新页（例如被预读的页）加入缓冲池时，只加入到老生代头部：

如果数据真正被读取（预读成功），才会加入到新生代的头部

如果数据没有被读取，则会比新生代里的“热数据页”更早被淘汰出缓冲池
<!-- more -->
https://blog.csdn.net/wuhenyouyuyouyu/article/details/93377605

https://www.cnblogs.com/wxlevel/p/12995324.html

buffer pool是什么?
是一块内存区域，当数据库操作数据的时候，把硬盘上的数据加载到buffer pool，不直接和硬盘打交道，操作的是buffer pool里面的数据
数据库的增删改查都是在buffer pool上进行，和undo log/redo log/redo log buffer/binlog一起使用，后续会把数据刷到硬盘上
默认大小 128M
数据页
磁盘文件被分成很多数据页，一个数据页里面有很多行数据
一个数据页默认大小 16K
更新一行数据，实际上是把行数据所在的 数据页 整个加载到buffer pool中

https://www.cnblogs.com/wasitututu/p/13612605.html

buffer pool的配置
innodb_buffer_pool_size：缓存区域的大小。
innodb_buffer_pool_chunk_size：当增加或减少innodb_buffer_pool_size时，操作以块（chunk）形式执行。块大小由innodb_buffer_pool_chunk_size配置选项定义，默认值128M。
innodb_buffer_pool_instances：当buffer pool比较大的时候（超过1G），innodb会把buffer pool划分成几个instances，这样可以提高读写操作的并发，减少竞争。读写page都使用hash函数分配给一个instances。
当增加或者减少buffer pool大小的时候，实际上是操作的chunk。buffer pool的大小必须是innodb_buffer_pool_chunk_sizeinnodb_buffer_pool_instances，如果配置的innodb_buffer_pool_size不是innodb_buffer_pool_chunk_sizeinnodb_buffer_pool_instances的倍数，buffer pool的大小会自动调整为innodb_buffer_pool_chunk_size*innodb_buffer_pool_instances的倍数，自动调整的值不少于指定的值。
如果指定的buffer大小是9G，instances的个数是16，chunk默认的大小是128M，那么buffer会自动调整为10G。具体的配置可以参考mysql官网的介绍mysql reference

https://blog.csdn.net/qq_27347991/article/details/81052728
https://www.jianshu.com/p/f9ab1cb24230