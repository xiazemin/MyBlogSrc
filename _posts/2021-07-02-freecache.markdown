---
title: freecache 无gc的go cache
layout: post
category: golang
author: 夏泽民
---
https://github.com/coocood/freecache
特性
能存储数亿个条目

零 GC 开销

高并发线程安全访问

纯 Go 实现

过期支持

类似 LRU 算法

严格限制内存使用

附带一个小服务器，支持带有管道功能的基本 Redis 命令
<!-- more -->
Set 性能比内置 map 快约 2 倍；Get 性能比内置 map 慢约 1/2 倍。由于它是基于单线程做的基准测试，因此在多线程环境中，FreeCache 应该比单锁保护的内置 map 快许多倍

它怎么做到零 GC 开销的？
FreeCache 通过减少指针数量避免了 GC 开销。无论其中存储了多少条目，都只有 512 个指针。通过 key 的哈希值将数据集分割为 256 个段。每个段只有两个指针，一个是存储键和值的环形缓冲区，另一个是用于查找条目的索引切片。每个段都有自己的锁 

https://freecache将缓存划分为256个segment，对于一个key的操作，freecache通过hash方法（xxhash）计算得到一个64位的hashValue，并取低8位作为segId，定位到具体的segment，并对segment加锁。由于只对segment加锁，不同segment之间可以并发进行key操作，所以freecache支持高并发线程安全访问。

segment底层实际上是由两个切片组成的复杂数据结构，其中一个切片用来实现环形缓冲区RingBuf，存储了所有的entry （entry=24 byte header + key + value）。另一个切片则是用于查找entry的索引切片slotData，slotData被逻辑上切分为256个slot，每个slot上的entry索引都是按照hash16有序排列的。可以看出，不管freecache缓存了多少数据，底层永远都只会有512个指针，所以freecache的对GC开销几乎为零。

对于一个key的set操作，首先判断key是否存在，不存在的情况处理比较简单，直接追加到环尾；如果存在的话，则看一下原来为entry预留的value容量是否充足，充足的话，直接覆盖，否则删掉原来的entry，并将新的entry追加到环尾，新的entry会给value预留多一点空间。

2.2.2 set操作为什么高效
采用二分查找，极大的减少查找entry索引的时间开销。slot切片上的entry索引是根据hash16值有序排列的，对于有序集合，可以采用二分查找算法进行搜索，假设缓存了n个key，那么查找entry索引的时间复杂度为log2(n * 2^-16) = log2(n) - 16。
对于key不存在的情况下（找不到entry索引）。
如果Ringbuf容量充足，则直接将entry追加到环尾，时间复杂度为O(1)。
如果RingBuf不充足，需要将一些key移除掉，情况会复杂点，后面会单独讲解这块逻辑，不过freecache通过一定的措施，保证了移除数据的时间复杂度为O(1)，所以对于RingBuf不充足时，entry追加操作的时间复杂度也是O(1)。
对于已经存在的key（找到entry索引）。
如果原来给entry的value预留的容量充足的话，则直接更新原来entry的头部和value，时间复杂度为O(1)。
如果原来给entry的value预留的容量不足的话，freecache为了避免移动底层数组数据，不直接对原来的entry进行扩容，而是将原来的entry标记为删除（懒删除），然后在环形缓冲区RingBuf的环尾追加新的entry，时间复杂度为O(1)。

freecache追加新entry时候，如果RingBuf的可用容量不足时，会从环头开始，通过近乎LRU的置换算法，将旧数据删掉，空出足够的空间出来，

所以atomic.LoadInt64(&seg.totalTime)/atomic.LoadInt64(&seg.totalCount)表示RingBuf中的entry最近一次访问时间戳的平均值，当一个entry的accessTime小于等于这个平均值，则认为这个entry是可以被置换掉的。这里我简单的总结一下freecache的entry置换算法：

最理性的情况下，即消息不过期、没有消息被标记删除、key被set进去之后，就没有再被访问过，在这种情况下，确实可以完全满足LRU算法，不过这种情况是不会发生的。
freecache选择将accessTime小于等于平均accessTime的entry进行置换，从大局来看，确实是将最近较少使用的缓存置换出去，从某种程度来将，是一种近LRU的置换算法。
freecache为什么不完全实现LRU置换算法呢？如果采用hash表+数组来实现LRU算法，维护hash表所带来的空间开销先不说，找出来的entry在环中的位置还是随机的，这种随机置换会产生空间碎片，如果要解决碎片问题性能将会大打折扣。如果不采用hash表来实现，则需要遍历所有entry索引，而且同样也会产生空间碎片。
在特殊情况下，环头的数据都比较新时，会导致一直找不到合适的entry进行置换，空出足够的空间，为了不影响set操作的性能，当连续5次出现环头entry不符合置换条件时，第6次置换如果entry还是不满足置换条件，也会被强制置换出去。

2.3 过期与删除实现
2.3.1 key过期
对于过期的数据，freecache会让它继续存储在RingBuf中，RingBuf从一开始初始化之后，就固定不变了，是否删掉数据，对RingBuf的实际占用空间不会产生影响。
当get到一个过期缓存时，freecache会删掉缓存的entry索引（但是不会将缓存从RingBuf中移除），然后对外报ErrNotFound错误。
当RingBuf的容量不足时，会从环头开始遍历，如果key已经过期，这时才会将它删除掉。
如果一个key已经过期时，在它被freecache删除之前，如果又重新set进来（过期不会主动删除entry索引，理论上有被重新set的可能），过期的entry容量充足的情况下，则会重新复用这个entry。
freecache这种过期机制，一方面减少了维护过期数据的工作，另一方面，freecache底层存储是采用数组来实现，要求缓存数据必须连续，缓存过期的剔除会带来空间碎片，挪动数组来维持缓存数据的连续性不是一个很好的选择。
2.3.2 key删除
freecache有一下两种情况会进行删除key操作：
外部主动调用del接口删除key。
set缓存时，发现key已经存在，但是为entry预留的cap不足时，会选择将旧的数据删掉，然后再环尾追加新的数据。
freecache的删除机制也是懒删除，删除缓存时，只会删掉entry索引，但是缓存还是会继续保留在RingBuf中，只是被标记为删除，等到RingBuf容量不足需要置换缓存时，才会对标记为删除的缓存数据做最后的删除工作。
freecache删除一个key，需要搜索entry索引和标记缓存数据，搜索entry索引的时间复杂度前面已经分析过了，为O(log2(n) - 16)，而标记缓存数据的时间复杂度为O(1)，所以删除操作性能上还是挺不错的。
2.4 entry索引
2.4.1 前提说明
  256个slot底层其实是共用同一个entry索引切片seg.slotsData，下面的所有图文描述的数组下标值，都是站在整个entry索引切片seg.slotsData看的，描述的结果可能会和freecache源码计算得到的结果不一致，不过不影响我们理解entry索引相关操作的原理。如果要和代码实际计算的值对应上，在entry索引切片没有扩容之前，可以减掉1024就是代码里slot的下标值；在扩容之后，减掉2048就是代码里slot的下标位置。

3.1 freecache的不足
需要一次性申请所有缓存空间。用于实现segment的RingBuf切片，从缓存被创建之后，其容量就是固定不变的，申请的内存也会一直被占用着，空间换时间，确实避免不了。
freecache的entry置换算法不是完全LRU，而且在某些情况下，可能会把最近经常被访问的缓存置换出去。
entry索引切片slotsData无法一次性申请足够的容量，当slotsData容量不足时，会进行空间容量x2的扩容，这种自动扩容机制，会带来一定的性能开销。
由于entry过期时，不会主动清理缓存数据，这些过期缓存的entry索引还会继续保存slot切片中，这种机制会加快entry索引切片提前进行扩容，而实际上除掉这些过期缓存的entry索引，entry索引切片的容量可能还是完全充足的。
为了保证LRU置换能够正常进行，freecache要求entry的大小不能超过缓存大小的1/1024，而且这个限制还不给动态修改，具体可以参考github上的issues。
3.2 使用freecache的注意事项
缓存的数据如果可以的话，大小尽量均匀一点，可以减少RingBuf容量不足时的置换工作开销。
缓存的数据不易过大，这样子才能缓存更多的key，提高缓存命中率。

https://blog.csdn.net/chizhenlian/article/details/108435024
