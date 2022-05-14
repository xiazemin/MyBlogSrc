---
title: innodb-buffer-pool-size
layout: post
category: mysql
author: 夏泽民
---
用于缓存索引和数据的内存大小，这个当然是越多越好， 数据读写在内存中非常快， 减少了对磁盘的读写。 当数据提交或满足检查点条件后才一次性将内存数据刷新到磁盘中。然而内存还有操作系统或数据库其他进程使用， 根据经验，推荐设置innodb-buffer-pool-size为服务器总可用内存的75%。 若设置不当， 内存使用可能浪费或者使用过多。 对于繁忙的服务器， buffer pool 将划分为多个实例以提高系统并发性， 减少线程间读写缓存的争用。buffer pool 的大小首先受 innodb_buffer_pool_instances 影响， 当然影响较小。

Innodb_buffer_pool_pages_data
Innodb buffer pool缓存池中包含数据的页的数目，包括脏页。单位是page。
eg、show global status like 'Innodb_buffer_pool_pages_data';

  
Innodb_buffer_pool_pages_total
innodb buffer pool的页总数目。单位是page。
eg：show global status like 'Innodb_buffer_pool_pages_total';

show global status like 'Innodb_page_size';

查看@@innodb_buffer_pool_size大小，单位字节
SELECT @@innodb_buffer_pool_size/1024/1024/1024; #字节转为G

在线调整InnoDB缓冲池大小，如果不设置，默认为128M
set global innodb_buffer_pool_size = 4227858432; ##单位字节

计算Innodb_buffer_pool_pages_data/Innodb_buffer_pool_pages_total*100%
当结果 > 95% 则增加 innodb_buffer_pool_size， 建议使用物理内存的 75%
当结果 < 95% 则减少 innodb_buffer_pool_size， 
建议设置大小为： Innodb_buffer_pool_pages_data * Innodb_page_size * 1.05 / (1024*1024*1024)
<!-- more -->
https://www.cnblogs.com/yaochunhui/p/14809865.html

innodb_buffer_pool_size参数

对于参数 innodb_buffer_pool_size的参数到底需要设置为多少比较合适？我们这里展开讨论一下。

概念理解

它的是一个内存区域，用来缓存 InnoDB存储引擎的表中的数据和索引数据。以便提高对 InnoDB存储引擎表中数据的查询访问速度。

在查看 innodb_buffer_pool_size这个参数值的时候，是以 字节(byte)为单位进行展现的。它的默认值为 134217728字节，也就是 128MB(134217728/1204/1024=128MB，1MB=1024KB,1KB=1024Byte)。

在MySQL 5.7.5版本后，可以在不重启MySQL数据库服务的情况下，动态修改这个参数的值，但是在此版本之前的老版本修改后需要重启数据库服务才可以生效。

相关连的参数

innodb_buffer_pool_chunk_size

innodb_buffer_pool_chunk_size默认值为 134217728字节，即 128MB。它可以按照 1MB的单位进行增加或减小。

可以简单的把它理解成是 innodb_buffer_pool_size增加或缩小最小单位。 innodb_buffer_pool_size是有一个或多个 innodb_buffer_pool_chunk_size组成的。

如果修改了 innodb_buffer_pool_chunk_size的值将会导致 innodb_buffer_pool_size的值改变。在修改该参数的时候，需要计算好最后的 innodb_buffer_pool_size是否符合服务器的硬件配置。

参数 innodb_buffer_pool_chunk_size的修改，需要重启MySQL数据库服务，不能再数据库服务运行的过程中修改。如果修改MySQL的配置文件 my.cnf之后，需要重启MySQL服务。

innodb_buffer_pool_instances

innodb_buffer_pool_instances的默认值为1，最大可以设置为64。

当 innodb_buffer_pool_instances不为1的时候，表示需要启用多个缓冲池实例，即把整个 innodb_buffer_pool_size在逻辑上划分为多个缓存池，多实例可以提高并发性，可以减少不同线程读写缓存页面时的争用。

参数 innodb_buffer_pool_instances的修改，需要重启MySQL数据库服务，不能再数据库服务运行的过程中修改。

关联参数之间的关系

innodb_buffer_pool_size这个参数的设置，离不开另外两个参数的影响，它们分别是： innodb_buffer_pool_chunk_size和 innodb_buffer_pool_instances。

它们三者存在这样的数学公式如下：

innodb_buffer_pool_size = innodb_buffer_pool_instances * innodb_buffer_pool_chunk_size * x

注意：

上述公式中的 x的值要满足 x>=1。
如果设置的 innodb_buffer_pool_size的值不是 innodb_buffer_pool_chunk_size * innodb_buffer_pool_instances两者的乘积或两者乘积的整数倍(即 x倍)，那么 innodb_buffer_pool_size会自动调整为另外一个值，使得这个值是上述两者的乘积或者乘积的整数倍，所以当我们设置 innodb_buffer_pool_size的值之后，最后的结果并不一定会是我们所设置的结果。

从上图可以看出三个参数之间大概的关系

chunk的数目最多到1000。
instance的数目最多到64。
一个 instance里面可以有一个或多个 chunk，到底有几个 chunk，取决于 chunk的大小是多少。在 buffer pool大小一定的情况下， instance数目一定的情况下，每一个 instance下面放多少给 chunk，取决于每一个 chunk的大小。 例如： buffer pool大小设置为 4GB， instance的值设置为 8，那么 chunk的值可以设置为默认的 128MB，也可以设置为 512MB，也可以设置为 256MB。 因为：
4GB=8*128MB*4，此时 chunk的值为 128MB,每个 instance里有 4个 chunk
4GB=8*512MB*1，此时 chunk的值为 512MB,每个 instance里有 1个 chunk
4GB=8*256MB*2，此时 chunk的值为 256MB,每个 instance里有 2个 chunk
其中的4、1、2就是每一个instance中chunk的数目。

设置innodb_buffer_pool_size时的注意事项

除了上面的数学公式之外，还需要满足以下几个约束：

innodb_buffer_pool_instances的取值范围为[1,64]。

为了避免潜在的性能问题，chunk的数目不要超过1000这个阈值。也就是说： innodb_buffer_pool_size / innodb_buffer_pool_chunk_size的值需要在[1,1000)之间。
当 innodb_buffer_pool_instances的值大于 1的时候，也就把说当把总的 innodb_buffer_pool_size的大小按照 innodb_buffer_pool_instances的数量划分为多个逻辑上的 buffer_pool的时候，每个逻辑上的 buffer_pool的大小要大于1GB。换算成公式就是说： innodb_buffer_pool_size / innodb_buffer_pool_instances > 1GB
innodb_buffer_pool_size的取值范围建议设置为物理服务器内存的50%~75%，要留一些内存供服务器操作系统本身使用。如果在数据库服务器上面除了部署MySQL数据库服务，还有部署其他应用和服务，则可以适当再减少对应的百分比。原则上是去除其他应用服务，剩余的可用内存的50%~70%作为 innodb_buffer_pool_size的配置值。
占50%~75%的内存是否真合适

这里，我们再仔细想一想，前面我们提到的 innodb_buffer_pool_size的值配置成服务器可用内存的75%是否对所有的服务都适用？

比如你有一个 256GB内存的服务器来作为MySQL数据库服务器，并且这个服务器上面没有其他的应用，只是用来部署MySQL数据库服务。那么此时我们把 innodb_buffer_pool_size的大小设置为 192GB(256GB * 0.75 = 192GB)，那么此时剩余的内存为 64GB(256GB - 192GB = 64GB)。

此时我们的服务器操作系统在运行的时候，是需要使用 64GB的内存吗？显然是不需要那么多内存给服务器操作系统使用的。 8GB的内存够操作系统使用了吧！

那还有 56GB(64GB - 8GB = 56GB)的内存就是浪费掉了。所以这个75%的比例是不太合适的。应该是留足操作系统使用的内存，其余的就都给 innodb_buffer_pool_size来使用就可以了。

我们不管它目前是不是需要那么多内存，既然我们已经觉得把一个256GB的服务器单独给MySQL来使用了，那么它现在用不了那么多的 innodb_buffer_pool_size，不代表1年以后仍然用不了那么多的 innodb_buffer_pool_size。所以我们在配置部署的时候，就把分配好的资源都配置好，即便现在不用，以后用到也不用再次排查调优配置了。

所以我们应该把服务器划分为不同的层级，然后再不同的层级采用不同的分配方式。

对于内存较小的系统(<= 1GB) 对于内存少于1GB的系统，InnoDB缓冲池大小最好采用MySQL默认配置值128MB。
对于中等RAM (1GB - 32GB)的系统 考虑到系统的RAM大小为1GB - 32GB的情况，我们可以使用这个粗糙的启发式计算OS需要:256MB+256*log2(RAM大小GB) 这里的合理之处是，对于低RAM配置，我们从操作系统需要的256MB的基础值开始，并随着RAM数量的增加以对数的方式增加分配。这样，我们就可以得出一个确定性的公式来为我们的操作系统需求分配RAM。我们还将为MySQL的其他需求分配相同数量的内存。 例如，在一个有3GB内存的系统中，我们会合理分配660MB用于操作系统需求，另外660MB用于MySQL其他需求，这样InnoDB的缓冲池大小约为1.6GB。
对于拥有较大RAM (> 32GB)的系统 对于RAM大于32GB的系统，我们将返回到计算系统RAM大小的5%的OS需求，以及计算MySQL其他需求的相同数量。因此，对于一个RAM大小为192GB的系统，我们的方法将InnoDB缓冲池大小设为165GB左右，这也是使用的最佳值。
这里有一个简单他示例供大家参考：


buffer pool和服务器内存配比示例图
判断现在的innodb_buffer_pool_size是否合适

这个参数的默认值为128MB，我们该怎么判断一个数据库中的这个参数配置的值是否合适？我们是应该增加它的值？还是要减少它的值呢？

判断依据

在开始判断 innodb_buffer_pool_size设置的值是否合适之前，我们先来了解几个参数，因为在判断 innodb_buffer_pool_size的值是否合适的时候，需要用到这几个参数的值去做比对衡量。 所以我们需要先了解这几个参数的具体含义才可以进行

innodb_buffer_pool_reads： The number of logical reads that InnoDB could not satisfy from the buffer pool, and had to read directly from disk.翻译过来就是：缓存池中不能满足的逻辑读的次数，这些读需要从磁盘总直接读取。 备注：逻辑读是指从缓冲池中读，物理读是指从磁盘读。
innodb_buffer_pool_read_requests： The number of logical read requests.翻译过来就是：从buffer pool中逻辑读请求次数。
从查询SQL命中buffer pool缓存的数据的概率上来判断当前buffer pool是否满足需求，计算从buffer pool中读取数据的百分比如下：

percent = innodb_buffer_pool_read_requests / (innodb_buffer_pool_reads + innodb_buffer_pool_read_requests) * 100%

上述的 percent>=99%，则表示当前的buffer pool满足当前的需求。否则需要考虑增加 innodb_buffer_pool_size的值。

innodb_buffer_pool_pages_data： The number of pages in the InnoDB buffer pool containing data. The number includes both dirty and clean pages.翻译过来就是：innodb buffer pool中的数据页的个数，它的值包括存储数据的数据页和脏数据页。
innodb_buffer_pool_pages_total： The total size of the InnoDB buffer pool, in pages.翻译过来就是：innodb buffer pool中的总的数据页数目
innodb_page_size： InnoDB page size (default 16KB). Many values are counted in pages; the page size enables them to be easily converted to bytes.翻译过来就是：innodb中的单个数据页的大小。
根据上述的参数解释，我们可以大概得到如下的公式，这个公式可以用来验证我们MySQL中的 innodb_buffer_pool_pages_total的实际值是否匹配现在数据库中的两个参数的比值。

innodb_buffer_pool_pages_total = innodb_buffer_pool_size / innodb_page_size

从buffer pool中缓存的数据占用的数据页的百分比来判断buffer pool的大小是否合适。以上可以得到这样一个数学公式：

percent = Innodb_buffer_pool_pages_data / Innodb_buffer_pool_pages_total * 100%

上述的 percent>=95%则表示我们的 innodb_buffer_pool_size的值是可以满足当前的需求的。 否则可以考虑增加 innodb_buffer_pool_size的大小。 具体设置大小的值，可以参考下面的计算出来的值。

innodb_buffer_pool_size = Innodb_buffer_pool_pages_data * Innodb_page_size * 1.05 / (1024*1024*1024)

以上就是对MySQL中innodb_buffer_pool_size参数的理解，希望能帮助到你。

关注我，获取更多的分享。

https://baijiahao.baidu.com/s?id=1682511413024407482&wfr=spider&for=pc
innodb_buffer_pool_size
适用版本：8.0、5.7、5.6、5.5
默认值：134217728
修改完后是否需要重启：是
作用：InnoDB缓冲池大小， 这对Innodb表来说非常重要。Innodb相比MyISAM表对缓冲更为敏感。
修改建议：系统内存的3/4，如32G内存则设置24G。
innodb_log_file_size、innodb_log_files_in_group
通常，日志文件的总大小应足够大，以使服务器可以消除工作负载活动中的高峰和低谷，这通常意味着有足够的重做日志空间来处理一个小时以上的写活动。
值越大，缓冲池中需要的检查点刷新活动越少，从而节省了磁盘I / O。较大的日志文件也会使崩溃恢复变慢。
日志文件的总大小（innodb_log_file_size * innodb_log_files_in_group）不能超过略小于512GB的最大值。
可以在高峰期间采样1分钟产生日志量，计算1小时的日志量。设置为1~2小时的日志量即可。如：20M * 90 = 1800M；
innodb_log_files_in_group:2, innodb_log_file_size=900M 或 innodb_log_files_in_group:3, innodb_log_file_size=600M
innodb_flush_log_at_trx_commit
适用版本：8.0、5.7、5.6、5.5
默认值：1
修改完后是否需要重启：是
作用：控制提交操作的严格ACID遵从性与更高的性能之间的平衡，当重新安排并批量执行与提交相关的I/O操作时，可以实现更高的性能。可以通过更改默认值来获得更好的性能，但随后可能会在崩溃中丢失事务。
1: 事务提交时，把事务日志从缓存区写到日志文件中，并且立刻写入到磁盘上。
0: 日志每秒写入一次并刷新到磁盘。尚未刷新日志的事务可能会在崩溃中丢失。
2：事务提交时，把事务日志从缓存区写到日志文件中，但不一定立刻写入到磁盘上。日志文件会每秒写入到磁盘，如果写入前系统崩溃，就会导致最后1秒的日志丢失。
修改建议：能够容忍系统崩溃，丢失1秒的数据，并对性能要求极高的场景可以设置为2，否则建议保持默认1。
innodb_io_capacity
适用版本：8.0、5.7、5.6、5.5
默认值：200
修改完后是否需要重启：是
作用：参数定义了InnoDB后台任务每秒可用的I/O操作数（IOPS），例如用于从buffer pool中刷新脏页和从change buffer中合并数据。innodb后台进程最大的I/O性能指标，影响刷新赃页和插入缓冲的数量，在高转速磁盘下，尤其是现在SSD盘得到普及，可以根据需要适当提高该参数的值。
修改建议：建议设置为innodb_io_capacity_max的1/2(系统IOPS的40%~50%)，通常SSD设置为2000，**
innodb_io_capacity_max
适用版本：8.0、5.7、5.6、5.5
默认值：2000
作用：在压力下，控制当刷新脏数据时MySQL每秒执行的写IO量解释一下什么叫“在压力下”，MySQL中称为”紧急情况”，是当MySQL在后台刷新时，它需要刷新一些数据为了让新的写操作进来。
修改建议：建议设置为略低于系统IOPS，innodb_io_capacity的2倍，(系统IOPS的80%~90%)，通常SSD设置为4000，普通HHD设置200-400即可。**
sync_binlog
适用版本：8.0、5.7、5.6、5.5
默认值：1
修改完后是否需要重启：是
作用：控制MySQL服务器将二进制日志同步到磁盘的频率。
1: 事务提交后，将二进制日志文件写入磁盘并立即刷新，相当于同步写入磁盘，不经过系统缓存。
0: 禁止MySQL服务器将二进制日志同步到磁盘。
N：每写入1000次系统缓存就执行一次写入磁盘并刷新的操作，会有数据丢失的风险。
修改建议：能够容忍系统崩溃时丢失数据，并对性能要求极高的场景可以提高此参数如1000，否则建议保持默认1。
back_log
适用版本：8.0、5.7、5.6、5.5
默认值：1
修改完后是否需要重启：是
作用：MySQL每处理一个连接请求时都会创建一个新线程与之对应。在主线程创建新线程期间，如果前端应用有大量的短连接请求到达数据库，MySQL会限制这些新的连接进入请求队列，由参数back_log控制。如果等待的连接数量超过back_log的值，则将不会接受新的连接请求，所以如果需要MySQL能够处理大量的短连接，需要提高此参数的大小。
修改建议：3000
innodb_flush_method
适用版本：8.0、5.7、5.6、5.5
默认值：fsync(unix)
作用：定义用于将数据刷新到InnoDB数据文件和日志文件的方法，这可能会影响I/O吞吐量。
修改建议：当数据盘使用SSD时，不需要使用多级缓冲，设置成O_DIRECT；如果您使用或打算将其他存储设备用于redo log files和data files，并且您的数据文件驻留在具有非电池后备缓存的设备上，请改用O_DIRECT
max_connections、max_user_connections
适当加大相应配置
保证max_connections > max_user_connections 20这样即可
如: 5000/4800
interactive_timeout、wait_timeout 适当调整
interactive_timeout：交互式连接超时时间(mysql工具、mysqldump等)
wait_timeout：非交互式连接超时时间，默认的连接mysql api程序,jdbc连接数据库等
join_buffer_size、read_buffer_size、sort_buffer_size、 适当增加
join_buffer_size：用于普通索引扫描，范围索引扫描和不使用索引的联接的缓冲区的最小大小，从而执行全表扫描。通常，获得快速联接的最佳方法是添加索引。
sort_buffer_size: 每个必须执行排序的会话都会分配此大小的缓冲区。sort_buffer_size并非特定于任何存储引擎，而是以一般方式进行优化。
read_buffer_size：对MyISAM表进行顺序扫描的每个线程都会为其扫描的每个表分配此大小（以字节为单位）的缓冲区。如果进行多次顺序扫描，则可能需要增加此值，默认为131072，不使用MyISAM引擎可以不用改。
innodb_autoinc_lock_mode
适用版本：8.0、5.7、5.6、5.5
默认值：1
修改完后是否需要重启：是
作用：在MySQL 5.1.22后，InnoDB为了解决自增主键锁表的问题，引入了参数innodb_autoinc_lock_mode，用于控制自增主键的锁机制。该参数可以设置的值为0、1、2，RDS默认的参数值为1，表示InnoDB使用轻量级别的mutex锁来获取自增锁，替代最原始的表级锁。但是在load data（包括INSERT … SELECT和REPLACE … SELECT）场景下若使用自增表锁，则可能导致应用在并发导入数据时出现死锁。
现象：在load data（包括INSERT … SELECT和REPLACE … SELECT）场景下若使用自增表锁，在并发导入数据时出现如下死锁：
RECORD LOCKS space id xx page no xx n bits xx index PRIMARY of table xx.xx trx id xxx lock_mode X insert intention waiting. TABLE LOCK table xxx.xxx trx id xxxx lock mode AUTO-INC waiting；
修改建议：建议将该参数值改为2，表示所有情况插入都使用轻量级别的mutex锁（只针对row模式），这样就可以避免auto_inc的死锁，同时在INSERT … SELECT的场景下性能会有很大提升。
注：当该参数值为2时，binlog的格式需要被设置为row。
query_cache_size
适用版本：5.7、5.6、5.5
默认值：1048576
修改完后是否需要重启：否
作用：该参数用于控制MySQL query cache的内存大小。如果MySQL开启query cache，在执行每一个query的时候会先锁住query cache，然后判断是否存在于query cache中，如果存在则直接返回结果，如果不存在，则再进行引擎查询等操作。同时，insert、update和delete这样的操作都会将query cahce失效掉，这种失效还包括结构或者索引的任何变化。但是cache失效的维护代价较高，会给MySQL带来较大的压力。所以，当数据库不会频繁更新时，query cache是很有用的，但如果写入操作非常频繁并集中在某几张表上，那么query cache lock的锁机制就会造成很频繁的锁冲突，对于这一张表的写和读会互相等待query cache lock解锁，从而导致select的查询效率下降。
现象：数据库中有大量的连接状态为checking query cache for query、Waiting for query cache lock、storing result in query cache。
修改建议：MySQL默认是关闭query cache功能的，如果您的实例打开了query cache，当出现上述情况后可以关闭query cache。
net_write_timeout
适用版本：8.0、5.7、5.6、5.5
默认值：60
修改完后是否需要重启：否
作用：等待将一个block发送给客户端的超时时间。
现象：若参数设置过小，可能会导致客户端出现如下错误：
the last packet successfully received from the server was milliseconds ago或the last packet sent successfully to the server was milliseconds ago.
修改建议：一般在网络条件比较差时或者客户端处理每个block耗时较长时，由于net_write_timeout设置过小导致的连接中断很容易发生，建议增加该参数的大小。
tmp_table_size
适用版本：8.0、5.7、5.6、5.5
默认值：2097152
修改完后是否需要重启：否
作用：该参数用于决定内部内存临时表的最大值，每个线程都要分配，实际起限制作用的是tmp_table_size和max_heap_table_size的最小值。如果内存临时表超出了限制，MySQL就会自动地把它转化为基于磁盘的MyISAM表。优化查询语句的时候，要避免使用临时表，如果实在避免不了的话，要保证这些临时表是存在内存中的。
现象：如果复杂的SQL语句中包含了group by、distinct等不能通过索引进行优化而使用了临时表，则会导致SQL执行时间加长。
修改建议：如果应用中有很多group by、distinct等语句，同时数据库有足够的内存，可以增大tmp_table_size（max_heap_table_size）的值，以此来提升查询性能。

https://www.jianshu.com/p/4f59083a2a76

