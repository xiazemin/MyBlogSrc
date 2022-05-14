---
title: elasticsearch 文件的存储
layout: post
category: elasticsearch
author: 夏泽民
---
在传统的数据库里面，对数据关系描述无外乎三种，一对一，一对多和多对多的关系，如果有关联关系的数据，通常我们在建表的时候会添加主外键来建立数据联系，然后在查询或者统计时候通过join来还原或者补全数据，最终得到我们需要的结果数据，那么转化到ElasticSearch里面，如何或者怎样来处理这些带有关系的数据。
ElasticSearch是一个NoSQL类型的数据库，本身是弱化了对关系的处理，因为像lucene，es，solr这样的全文检索框架对性能要求都是比较高的，一旦出现join这样的操作，性能会非常差，所以在使用搜索框架时，我们应该避免把搜索引擎当做关系型数据库用。

当然，现实数据肯定是有关系的，那么在es里面是如何处理和管理这些带有关系的数据呢？

大家都知道，es天生对json数据支持的非常完美，只要是标准的json结构的数据，无论多么复杂，无论是嵌套多少层，都能存储到es里面，进而能够查询和分析，检索。在这种机制上，es处理和管理关系主要有三种方式：

一，使用objcet和array[object]的字段类型自动存储多层结构的json数据
这是es默认的机制，也就是我们并没有设置任何mapping，直接向es服务端插入一条复杂的json数据，也能成功插入，并能支持检索，（能这样操作是因为es默认用的是动态mapping，只要插入的是标准的json结构就会自动转换，当然我们也能控制mapping类型，es里面有动态mapping和静态maping，静态mapping还分严格类型，弱类型，一般类型
{
  "name" : "Zach",
  "car" : [
    {
      "make" : "Saturn",
      "model" : "SL"
    },
    {
      "make" : "Subaru",
      "model" : "Imprezza"
    }
  ]
最终转化成的存储结构是下面这样的：
{
  "name" : "Zach",
  "car.make" : ["Saturn", "Subaru"]
  "car.model" : ["SL", "Imprezza"]
}
因为es的底层lucene是天生支持多值域的存储，所以在上面看起来像数组的结构，其实在es里面存储的就是这个字段多值域。

然后检索的时候.符号就能检索相对应的内容。这样的一条数据，其实已经包含了数据和关系，看起来像一对多的关系，一个人拥有多辆汽车。但实际上并不能算严格意义上的关系，因为lucene底层是扁平化存储的，这样以来多个汽车的数据实际都是存到一起的混杂的，你没办法单独获取到这个人某一辆汽车的数据，因为整条数据都是一个整体，无论什么操作整条数据都会返回。

二，使用nested[object]类型，存储拥有多级关系的数据
在方案一里面，我们指出了array存储的数组对象，并不是严格意义的关系，因为第二层的数据是没有分离的，如果想要分离，就必须使用nested类型来显式定义数据结构。只有这样，第二层的多个汽车数据才是独立的互不影响，也就是说可以单独获取或查询某一辆汽车的数据
在方案1里面，最终到es里面会存储一条数据，在第二种类型里面，而如果声明了car类型是nested，那么最终存储到es的数量会显示3，这里解释一下3是怎么来的 = 1个root文档+2个汽车文档，nested声明类型，每一个实例都是一个新的document，所以在查询的时候才能够独立进行查询，并且性能还不错，因为es底层会把整条数据存在同一个shard的lucene的sengment里面，缺点是更新的代价比较大，每一个子文档的更新都要重建整个结构体的索引，所以nested适合不经常update的嵌套多级关系的场景。

nested类型的数据，需要用其指定的查询和聚合方法才能生效，普通的es查询只能查询1级也就是root级的属性，嵌套的属性是不能查的，如果想要查，必须用嵌套查询或者聚合才行。

嵌套应用有两种模式：

第一种：嵌套查询

每个查询都是单个文档内生效，包括排序,

第二种：嵌套聚合或者过滤

对同一层级的所有文档都是全局生效，包括过滤排序

三，parent/children 父子关系
parent/children 模式与nested非常类似，但是应用场景侧重点有所不同。

在使用parent/children管理关联关系时，es会在每个shard的内存中维护一张关系表，在检索时，通过has_parent和has_child过滤器来得到关联的数据，这种模式下父文档与子文档也是独立的，查询性能会比nested模式稍低，因为父文档和子文档在插入的时候会通过route使得他们都分布在同一个shard里面，但并不保证在同一个lucene的sengment索引段里面，所以检索性能稍低，除此之外，每次检索es都需要从内存的关系表里面得到数据关联的信息，也需要花费一定的时间，相比nested的优势在于，父文档或者子文档的更新，并不影响其他的文档，所以对于更新频繁的多级关系，使用parent/children模式，最为合适不过。
插入数据时，需要先插入父文档：
然后插入子文档时，需要加上路由字段：
最终，父文档zach就关联上了两个子文档，在查询时候可以通过parent/children特定查询来获取数据。

总结：

方法一：

（1）简单，快速，性能较高

（2）对维护一对一的关系比较擅长

（3）不需要特殊的查询

方法二：

（1）由于底层存储在同一个lucene的sengment里，所以读取和查询性能对比方法三更快

（2）更新单个子文档，会重建整个数据结构，所以不适合更新频繁的嵌套场景

（3）可以维护一对多和多对多的存储关系

方法三：

（1）多个关系数据，存储完全独立，但是存在同一个shard里面，所以读取和查询性能比方法二稍低

（2）需要额外的内存，维护管理关系列表

（3）更新文档不影响其他的子文档，所以适合更新频繁的场景

（4）排序和评分操作比较麻烦，需要额外的脚本函数支持

原始数据都存在主分片中，副分片只是主分片的一个副本，便于节点管理数据，规避宕机等丢失数据的风险
Elasticsearch 如何知道一个文档应该存放到哪个分片中呢？当我们创建文档时，它如何决定这个文档应当被存储在分片 1 还是分片 2 中呢？ 
首先这肯定不会是随机的，否则将来要获取文档的时候我们就不知道从何处寻找了。实际上，这个过程是根据下面这个公式决定的： 
shard = hash(routing) % number_of_primary_shards 
routing 是一个可变值，默认是文档的 _id ，也可以设置成一个自定义的值。 routing 通过 hash 函数生成一个数字，然后这个数字再除以 number_of_primary_shards （主分片的数量）后得到 余数 。这个分布在 0 到 number_of_primary_shards-1 之间的余数，就是我们所寻求的文档所在分片的位置。

这就解释了为什么我们要在创建索引的时候就确定好主分片的数量 并且永远不会改变这个数量：因为如果数量变化了，那么所有之前路由的值都会无效，文档也再也找不到了。

所有的文档 API（ get 、 index 、 delete 、 bulk 、 update 以及 mget ）都接受一个叫做 routing 的路由参数 ，通过这个参数我们可以自定义文档到分片的映射。一个自定义的路由参数可以用来确保所有相关的文档——例如所有属于同一个用户的文档——都被存储到同一个分片中。我们也会在扩容设计这一章中详细讨论为什么会有这样一种需求。

主分片和副本分片如何交互
为了说明目的, 我们假设有一个集群由三个节点组成。 它包含一个叫 blogs 的索引，有两个主分片，每个主分片有两个副本分片。相同分片的副本不会放在同一节点
在主副分片和任何副本分片上面 成功新建，索引和删除文档所需要的步骤顺序：

客户端向 Node 1 发送新建、索引或者删除请求。
节点使用文档的 _id 确定文档属于分片 0。请求会被转发到Node 3，因为分片 0 的主分片目前被分配在Node 3 上。
Node 3 在主分片上面执行请求。如果成功了，它将请求并行转发到 Node 1 和 Node 2 的副本分片上。一旦所有的副本分片都报告成功, Node 3 将向协调节点报告成功，协调节点向客户端报告成功。
在客户端收到成功响应时，文档变更已经在主分片和所有副本分片执行完成，变更是安全的。 
有一些可选的请求参数允许您影响这个过程，可能以数据安全为代价提升性能。这些选项很少使用，因为Elasticsearch已经很快，但是为了完整起见，在这里阐述如下： 
一致性 
默认情况下，主分片 需要 规定数量(quorum),或大多数的分片 (其中分片副本可以是主分片或者副本分片)在写入操作时可用。这是为了防止将数据写入到网络分区的‘`背面’’。规定的数量定义公式如下： 
int( (primary + number_of_replicas) / 2 ) + 1

允许的 一致性 值是 一个 （只是主分片）或者 所有（主分片和副本分片）, 或者默认的规定数量或者大多数的副本分片。

注意 number_of_replicas 是在索引中的设置指定的分片数，不是当前处理活动状态的副本分片数。如果你指定索引应该有三个副本分片，那规定数量计算公式是： 
int( (primary + 3 replicas) / 2 ) + 1 = 3

但是如果只启动两个节点，则活动分片副本无法满足规定数量，并且您将无法索引和删除任何文档。 
超时 如果没有足够的副本分片会发生什么？ Elasticsearch会等待，希望更多的分片出现。默认情况下，它最多等待1分钟。 如果你需要，你可以使用 timeout 参数 使它更早终止： 100 100毫秒，30s 是30秒。

Note 
新索引默认有 1 个副本分片，这意味着为满足 规定数量 应该 需要两个活动的分片副本。 但是，这些默认的设置会阻止我们在单一节点上做任何事情。为了避免这个问题，要求只有当 number_of_replicas 大于1的时候，规定数量才会执行。

总结
在增删改的时候，也就是说如果像主例子那样，我们新建了一个索引，规定这个索引的主分片是2个，每个主分片有2个副分片，那么它能保证数据一致性的要求就是它有3 =（2+2）/2 + 1 个活动的副分片，如果我们这个时候只启动了两个节点，也就是这样子的运行： 
Node1:R0 P1 
Node2:R1 P0 
此时只有2个副分片，这个时候无法保证在操作的时候保证主分片的值和所有副分片一致。 
即最好启动3个节点。
从主分片或者副本分片检索文档的步骤顺序：

1、客户端向 Node 1 发送获取请求。

2、节点使用文档的 _id 来确定文档属于分片 0 。分片 0 的副本分片存在于所有的三个节点上。 在这种情况下，它将请求转发到 Node 2 。（为了读取请求，协调节点在每次请求的时候将选择不同的副本分片来达到负载均衡；通过轮询所有的副本分片。 
）

3、Node 2 将文档返回给 Node 1 ，然后将文档返回给客户端。

在文档被检索时，已经被索引的文档可能已经存在于主分片上但是还没有复制到副本分片。 在这种情况下，副本分片可能会报告文档不存在，但是主分片可能成功返回文档。 一旦索引请求成功返回给用户，文档在主分片和副本分片都是可用的。

总结：
在查询时，也就是说不像增删改操作那样必须到主分片执行，可以轮询访问所有的包含文档的主副分片，如果副分片此时不存在，也会再去访问主分片返回文档。此时如果返回了文档，此文档在副分片也是可查的了。

局部更新文档
部分更新一个文档的步骤：

客户端向 Node 1 发送更新请求。
它将请求转发到主分片所在的 Node 3 。
Node 3 从主分片检索文档，修改 _source 字段中的 JSON ，并且尝试重新索引主分片的文档。 如果文档已经被另一个进程修改，它会重试步骤 3 ，超过 retry_on_conflict 次后放弃。
如果 Node 3 成功地更新文档，它将新版本的文档并行转发到 Node 1 和 Node 2 上的副本分片，重新建立索引。 一旦所有副本分片都返回成功， Node 3 向协调节点也返回成功，协调节点向客户端返回成功。
update API 还接受在 新建、索引和删除文档 章节中介绍的 routing 、 replication 、 consistency 和 timeout 参数。

基于文档的复制

当主分片把更改转发到副本分片时， 它不会转发更新请求。 相反，它转发完整文档的新版本。请记住，这些更改将会异步转发到副本分片，并且不能保证它们以发送它们相同的顺序到达。 如果Elasticsearch仅转发更改请求，则可能以错误的顺序应用更改，导致得到损坏的文档。

多文档模式
mget 和 bulk API 的 模式类似于单文档模式。区别在于协调节点知道每个文档存在于哪个分片中。 它将整个多文档请求分解成 每个分片 的多文档请求，并且将这些请求并行转发到每个参与节点
使用单个 mget 请求取回多个文档所需的步骤顺序： 
1. 客户端向 Node 1 发送 mget 请求。 
2. Node 1 为每个分片构建多文档获取请求，然后并行转发这些请求到托管在每个所需的主分片或者副本分片的节点上。一旦收到所有答复， Node 1 构建响应并将其返回给客户端。 
可以对 docs 数组中每个文档设置 routing 参数。
bulk API 按如下步骤顺序执行：

客户端向 Node 1 发送 bulk 请求。
Node 1 为每个节点创建一个批量请求，并将这些请求并行转发到每个包含主分片的节点主机。
主分片一个接一个按顺序执行每个操作。当每个操作成功时，主分片并行转发新文档（或删除）到副本分片，然后执行下一个操作。 一旦所有的副本分片报告所有操作成功，该节点将向协调节点报告成功，协调节点将这些响应收集整理并返回给客户端。
bulk API 还可以在整个批量请求的最顶层使用 consistency 参数，以及在每个请求中的元数据中使用 routing 参数。

为什么是有趣的格式？
当我们早些时候在代价较小的批量操作章节了解批量请求时， 您可能会问自己， “为什么 bulk API 需要有换行符的有趣格式，而不是发送包装在 JSON 数组中的请求，例如 mget API？” 。 
为了回答这一点，我们需要解释一点背景：在批量请求中引用的每个文档可能属于不同的主分片， 每个文档可能被分配给集群中的任何节点。这意味着批量请求 bulk 中的每个 操作 都需要被转发到正确节点上的正确分片。 
如果单个请求被包装在 JSON 数组中，那就意味着我们需要执行以下操作： 
• 将 JSON 解析为数组（包括文档数据，可以非常大） 
• 查看每个请求以确定应该去哪个分片 
• 为每个分片创建一个请求数组 
• 将这些数组序列化为内部传输格式 
• 将请求发送到每个分片 
这是可行的，但需要大量的 RAM 来存储原本相同的数据的副本，并将创建更多的数据结构，Java虚拟机（JVM）将不得不花费时间进行垃圾回收。 
相反，Elasticsearch可以直接读取被网络缓冲区接收的原始数据。 它使用换行符字符来识别和解析小的 action/metadata 行来决定哪个分片应该处理每个请求。 
这些原始请求会被直接转发到正确的分片。没有冗余的数据复制，没有浪费的数据结构。整个请求尽可能在最小的内存中处理

<!-- more -->
elasticsearch副本提供了高可靠性；它可以保证节点丢失而不会中断服务，但是副本不能做到容灾备份，所以需要把elasticsearch的数据被分到hdfs中。
操作步骤
安装repository-hdfs

 进入ES的目录，执行命令：bin/elasticsearch-plugin install repository-hdfs
 如果需要移除插件，执行命令:bin/elasticsearch-plugin remove repository-hdfs
建立仓库命令

curl -H "Content-Type: application/json" -XPUT 'http://192.168.2.227:9200/_snapshot/backup' -d '{"type":"hdfs", "settings":{ "path":"/elasticsearch/respositories/my_hdfs_repository", "uri":"hdfs://192.168.2.202:9000" }}’
备注:一个仓库可以包含多个快照

查看仓库

  curl -XGET  'http://192.168.2.227:9200/_snapshot/backup?pretty'
快照特定的索引

  curl -H "Content-Type: application/json" -XPUT 'http://192.168.2.227:9200/_snapshot/backup/snapshot_1' -d '{"indices":"blog"}'
快照多个索引

  curl -H "Content-Type: application/json" -XPUT 'http://192.168.2.227:9200/_snapshot/backup/snapshot_1' -d '{"indices":"blog1,blog2"}'
经过这一步操作可以查看到hdfs里多了备份文件

查看特定快照信息

  curl -XGET 'http://192.168.2.227:9200/_snapshot/backup/snapshot_1?pretty'
INITIALIZING：分片在检查集群状态看看自己是否可以被快照。这个一般是非常快的。
STARTED：数据正在被传输到仓库。
FINALIZING：数据传输完成；分片现在在发送快照元数据。
DONE：快照完成！
FAILED：快照处理的时候碰到了错误，这个分片/索引/快照不可能完成了。检查你的日志获取更多信息。

恢复特定索引

  curl -XPOST 'http://192.168.2.227:9200/_snapshot/backup/snapshot_1/_restore?pretty'
删除快照

  curl -XDELETE 'http://192.168.2.227:9200/_snapshot/backup/snapshot_1’
监控快照

  curl -XGET 'http://192.168.2.227:9200/_snapshot/backup/snapshot_1/_status'
执行脚本

 #!/bin/bash
 current_time=$(date +%Y%m%d%H%M%S)
 command_prefix="http://192.168.2.227:9200/_snapshot/backup/all_"
 command=$command_prefix$current_time
 curl -H "Content-Type: application/json" -XPUT $command -d '{"indices":"index*,logstash*,nginx*,magicianlog*,invokelog*,outside*"}'
crontab定时任务，每天备份一次

0 0 * * * /root/shell/snapshot_all_hdfs.sh >>/root/shell/logs/snapshot_all_day.log 2>&1
常见问题处理
    failed to create snapshot","caused_by":{"type":"access_control_exception","reason":"Permission denied: user=elasticsearch, access=WRITE, inode=\"/elasticsearch/respositories/my_hdfs_repository
修改hdfs-site.xml，把dfs.permissions改成false
提前在hdfs中新建备份的文件目录， /elasticsearch/respositories/my_hdfs_repository
并修改hdfs文件权限，bin/hdfs dfs -chmod 777

插件git地址：https://github.com/elastic/ela … -hdfs 
下载地址：https://download.elastic.co/elasticsearch/elasticsearch-repository-hdfs/elasticsearch-repository-hdfs-2.2.0-hadoop2.zip 
在线安装 
进入ES的目录，执行命令：bin/elasticsearch-plugin install repository-hdfs 
在线安装 
现将下载好的zip包，放在指定目录，如/home/hadoop/elk/es-reporitory.zip，然后执行命令：bin/plugin install file:///home/hadoop/elk/es-reporitory.zip

 elasticsearch-hadoop是一个深度集成Hadoop和ElasticSearch的项目，也是ES官方来维护的一个子项目，通过实现Hadoop和ES之间的输入输出，可以在Hadoop里面对ES集群的数据进行读取和写入，充分发挥Map-Reduce并行处理的优势，为Hadoop数据带来实时搜索的可能。 
项目网址：http://www.elasticsearch.org/overview/hadoop/

hadoop-ES读写最主要的就是ESInputFormat、ESOutputFormat的参数配置（Configuration）。

另外 其它数据源操作（Mysql等）也是类似，找到对应的InputFormat，OutputFormat配置上环境参数。
