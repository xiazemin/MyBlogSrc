---
title: binlog同步kafka方案
layout: post
category: storage
author: 夏泽民
---
Canel，Databus，Puma等，这些都是需要部署server和client的。其中server端是由这些工具实现，配置了就可以读binlog，而client端是需要我们动手编写程序的，远没有达到我即插即用的期望和懒人的标准。
　　再来看看flume，只需要写一个配置文件，就可以完成数据同步的操作。官网：http://flume.apache.org/FlumeUserGuide.html#flume-sources。它的数据源默认是没有读取binlog日志实现的，也没有读数据库表的官方实现，只能用开源的自定义source：https://github.com/keedio/flume-ng-sql-source
同步的格式
　　原作者的插件flume-ng-sql-source只支持csv的格式，如果开始同步之后，数据库表需要增减字段，则会给开发者造成很大的困扰。所以我添加了一个分支版本，用来将数据以JSON的格式，同步到kafka，字段语义更加清晰。
　　sql-json插件包下载地址：https://github.com/yucy/flume-ng-sql-source-json/releases/download/1.0/flume-ng-sql-source-json-1.0.jar
　　将此jar包下载之后，和相应的数据库驱动包，一起放到flume的lib目录之下即可。
处理机制
flume-ng-sql-source在【status.file.name】文件中记录读取数据库表的偏移量，进程重启后，可以接着上次的进度，继续增量读表。
启动说明
说明：启动命令里的【YYYYMM=201711】，会传入到flume.properties里面，替换${YYYYMM}
[test@localhost ~]$ YYYYMM=201711 bin/flume-ng agent -c conf -f conf/flume.properties -n sync &
 -c：表示配置文件的目录，在此我们配置了flume-env.sh，也在conf目录下；
 -f：指定配置文件，这个配置文件必须在全局选项的--conf参数定义的目录下，就是说这个配置文件要在前面配置的conf目录下面；
 -n：表示要启动的agent的名称，也就是我们flume.properties配置文件里面，配置项的前缀，这里我们配的前缀是【sync】；
flume的配置说明
flume-env.sh
 # 配置JVM堆内存和java运行参数，配置-DpropertiesImplementation参数是为了在flume.properties配置文件中使用环境变量
export JAVA_OPTS="-Xms512m -Xmx512m -Dcom.sun.management.jmxremote -DpropertiesImplementation=org.apache.flume.node.EnvVarResolverProperties"
 关于propertiesImplementation参数的官方说明：http://flume.apache.org/FlumeUserGuide.html#using-environment-variables-in-configuration-files
flume.properties
<!-- more -->
canal是阿里巴巴旗下的一款开源项目，纯Java开发。基于数据库增量日志解析，提供增量数据订阅&消费，目前主要支持了MySQL（也支持mariaDB）。

起源：早期，阿里巴巴B2B公司因为存在杭州和美国双机房部署，存在跨机房同步的业务需求。不过早期的数据库同步业务，主要是基于trigger的方式获取增量变更，不过从2010年开始，阿里系公司开始逐步的尝试基于数据库的日志解析，获取增量变更进行同步，由此衍生出了增量订阅&消费的业务，从此开启了一段新纪元。

基于日志增量订阅&消费支持的业务：

数据库镜像
数据库实时备份
多级索引 (卖家和买家各自分库索引)
search build
业务cache刷新
价格变化等重要业务消息
原理相对比较简单：
1、canal模拟mysql slave的交互协议，伪装自己为mysql slave，向mysql master发送dump协议
2、mysql master收到dump请求，开始推送binary log给slave(也就是canal)
3、canal解析binary log对象(原始为byte流)
架构设计 
个人理解，数据增量订阅与消费应当有如下几个点：
1、增量订阅和消费模块应当包括binlog日志抓取，binlog日志解析，事件分发过滤（EventSink），存储（EventStore）等主要模块。
2、如果需要确保HA可以采用Zookeeper保存各个子模块的状态，让整个增量订阅和消费模块实现无状态化，当然作为consumer(客户端)的状态也可以保存在zk之中。
3、整体上通过一个Manager System进行集中管理，分配资源。
源码以及项目介绍： https://github.com/alibaba/canal 
canal消费端项目开源:Otter(分布式数据库同步系统)，地址：https://github.com/alibaba/otter

mysql主备复制分成三步：

master将改变记录到二进制日志(binary log)中（这些记录叫做二进制日志事件，binary log events，可以通过show binlog events进行查看）；
slave将master的binary log events拷贝到它的中继日志(relay log)；
slave重做中继日志中的事件，将改变反映它自己的数据。
原理相对比较简单：

canal模拟mysql slave的交互协议，伪装自己为mysql slave，向mysql master发送dump协议
mysql master收到dump请求，开始推送binary log给slave(也就是canal)
canal解析binary log对象(原始为byte流)
<img src="{{site.url}}{{site.baseurl}}/img/canel1.jpg"/>
说明：

server代表一个canal运行实例，对应于一个jvm
instance对应于一个数据队列  （1个server对应1..n个instance)
instance模块：

eventParser (数据源接入，模拟slave协议和master进行交互，协议解析)
eventSink (Parser和Store链接器，进行数据过滤，加工，分发的工作)
eventStore (数据存储)
metaManager (增量订阅&消费信息管理器)
知识科普
mysql的Binlay Log介绍

http://dev.mysql.com/doc/refman/5.5/en/binary-log.html
http://www.taobaodba.com/html/474_mysqls-binary-log_details.html
简单点说：

mysql的binlog是多文件存储，定位一个LogEvent需要通过binlog filename +  binlog position，进行定位
mysql的binlog数据格式，按照生成的方式，主要分为：statement-based、row-based、mixed。
mysql> show variables like 'binlog_format';  
+---------------+-------+  
| Variable_name | Value |  
+---------------+-------+  
| binlog_format | ROW   |  
+---------------+-------+  
1 row in set (0.00 sec)  
目前canal只能支持row模式的增量订阅(statement只有sql，没有数据，所以无法获取原始的变更日志)

EventParser设计
大致过程：
<img src="{{site.url}}{{site.baseurl}}/img/canel2.jpg"/>
整个parser过程大致可分为几步：

Connection获取上一次解析成功的位置  (如果第一次启动，则获取初始指定的位置或者是当前数据库的binlog位点)
Connection建立链接，发送BINLOG_DUMP指令
 // 0. write command number
 // 1. write 4 bytes bin-log position to start at
 // 2. write 2 bytes bin-log flags
 // 3. write 4 bytes server id of the slave
 // 4. write bin-log file name
Mysql开始推送Binaly Log
接收到的Binaly Log的通过Binlog parser进行协议解析，补充一些特定信息
// 补充字段名字，字段类型，主键信息，unsigned类型处理
传递给EventSink模块进行数据存储，是一个阻塞操作，直到存储成功
存储成功后，定时记录Binaly Log位置
mysql的Binlay Log网络协议：
<img src="{{site.url}}{{site.baseurl}}/img/canel3.png"/>
说明：
图中的协议4byte header，主要是描述整个binlog网络包的length
binlog event structure，详细信息请参考： http://dev.mysql.com/doc/internals/en/binary-log.html
EventSink设计
<img src="{{site.url}}{{site.baseurl}}/img/canel4.jpg"/>
明：

数据过滤：支持通配符的过滤模式，表名，字段内容等
数据路由/分发：解决1:n (1个parser对应多个store的模式)
数据归并：解决n:1 (多个parser对应1个store)
数据加工：在进入store之前进行额外的处理，比如join
数据1:n业务
  为了合理的利用数据库资源， 一般常见的业务都是按照schema进行隔离，然后在mysql上层或者dao这一层面上，进行一个数据源路由，屏蔽数据库物理位置对开发的影响，阿里系主要是通过cobar/tddl来解决数据源路由问题。
  所以，一般一个数据库实例上，会部署多个schema，每个schema会有由1个或者多个业务方关注
数据n:1业务
  同样，当一个业务的数据规模达到一定的量级后，必然会涉及到水平拆分和垂直拆分的问题，针对这些拆分的数据需要处理时，就需要链接多个store进行处理，消费的位点就会变成多份，而且数据消费的进度无法得到尽可能有序的保证。
  所以，在一定业务场景下，需要将拆分后的增量数据进行归并处理，比如按照时间戳/全局id进行排序归并.
EventStore设计
1.  目前仅实现了Memory内存模式，后续计划增加本地file存储，mixed混合模式
2.  借鉴了Disruptor的RingBuffer的实现思路
RingBuffer设计
<img src="{{site.url}}{{site.baseurl}}/img/canel5.jpg"/>
定义了3个cursor
Put :  Sink模块进行数据存储的最后一次写入位置
Get :  数据订阅获取的最后一次提取位置
Ack :  数据消费成功的最后一次消费位置
借鉴Disruptor的RingBuffer的实现，将RingBuffer拉直来看：
<img src="{{site.url}}{{site.baseurl}}/img/canel6.jpg"/>
实现说明：
Put/Get/Ack cursor用于递增，采用long型存储
buffer的get操作，通过取余或者与操作。(与操作： cusor & (size - 1) , size需要为2的指数，效率比较高)
Instance设计
<img src="{{site.url}}{{site.baseurl}}/img/canel7.jpg"/>
instance代表了一个实际运行的数据队列，包括了EventPaser,EventSink,EventStore等组件。
抽象了CanalInstanceGenerator，主要是考虑配置的管理方式：
manager方式： 和你自己的内部web console/manager系统进行对接。(目前主要是公司内部使用)
spring方式：基于spring xml + properties进行定义，构建spring配置. 
Server设计
<img src="{{site.url}}{{site.baseurl}}/img/canel8.jpg"/>
server代表了一个canal的运行实例，为了方便组件化使用，特意抽象了Embeded(嵌入式) / Netty(网络访问)的两种实现
Embeded :  对latency和可用性都有比较高的要求，自己又能hold住分布式的相关技术(比如failover)
Netty :  基于netty封装了一层网络协议，由canal server保证其可用性，采用的pull模型，当然latency会稍微打点折扣，不过这个也视情况而定。(阿里系的notify和metaq，典型的push/pull模型，目前也逐步的在向pull模型靠拢，push在数据量大的时候会有一些问题) 
增量订阅/消费设计
<img src="{{site.url}}{{site.baseurl}}/img/canel9.jpg"/>
具体的协议格式，可参见：CanalProtocol.proto

get/ack/rollback协议介绍：

Message getWithoutAck(int batchSize)，允许指定batchSize，一次可以获取多条，每次返回的对象为Message，包含的内容为：
a. batch id 唯一标识
b. entries 具体的数据对象，对应的数据对象格式：EntryProtocol.proto
void rollback(long batchId)，顾命思议，回滚上次的get请求，重新获取数据。基于get获取的batchId进行提交，避免误操作
void ack(long batchId)，顾命思议，确认已经消费成功，通知server删除数据。基于get获取的batchId进行提交，避免误操作
canal的get/ack/rollback协议和常规的jms协议有所不同，允许get/ack异步处理，比如可以连续调用get多次，后续异步按顺序提交ack/rollback，项目中称之为流式api. 

流式api设计的好处：

get/ack异步化，减少因ack带来的网络延迟和操作成本 (99%的状态都是处于正常状态，异常的rollback属于个别情况，没必要为个别的case牺牲整个性能)
get获取数据后，业务消费存在瓶颈或者需要多进程/多线程消费时，可以不停的轮询get数据，不停的往后发送任务，提高并行化.  (作者在实际业务中的一个case：业务数据消费需要跨中美网络，所以一次操作基本在200ms以上，为了减少延迟，所以需要实施并行化) 
每次get操作都会在meta中产生一个mark，mark标记会递增，保证运行过程中mark的唯一性
每次的get操作，都会在上一次的mark操作记录的cursor继续往后取，如果mark不存在，则在last ack cursor继续往后取
进行ack时，需要按照mark的顺序进行数序ack，不能跳跃ack. ack会删除当前的mark标记，并将对应的mark位置更新为last ack cusor
一旦出现异常情况，客户端可发起rollback情况，重新置位：删除所有的mark, 清理get请求位置，下次请求会从last ack cursor继续往后取

说明：

可以提供数据库变更前和变更后的字段内容，针对binlog中没有的name,isKey等信息进行补全
可以提供ddl的变更语句

HA机制设计
canal的ha分为两部分，canal server和canal client分别有对应的ha实现

canal server:  为了减少对mysql dump的请求，不同server上的instance要求同一时间只能有一个处于running，其他的处于standby状态. 
canal client: 为了保证有序性，一份instance同一时间只能由一个canal client进行get/ack/rollback操作，否则客户端接收无法保证有序。
整个HA机制的控制主要是依赖了zookeeper的几个特性，watcher和EPHEMERAL节点(和session生命周期绑定)，可以看下我之前zookeeper的相关文章。
Canal Server: 
大致步骤：

canal server要启动某个canal instance时都先向zookeeper进行一次尝试启动判断  (实现：创建EPHEMERAL节点，谁创建成功就允许谁启动)
创建zookeeper节点成功后，对应的canal server就启动对应的canal instance，没有创建成功的canal instance就会处于standby状态
一旦zookeeper发现canal server A创建的节点消失后，立即通知其他的canal server再次进行步骤1的操作，重新选出一个canal server启动instance.
canal client每次进行connect时，会首先向zookeeper询问当前是谁启动了canal instance，然后和其建立链接，一旦链接不可用，会重新尝试connect.
Canal Client的方式和canal server方式类似，也是利用zokeeper的抢占EPHEMERAL节点的方式进行控制. 
项目的代码： https://github.com/alibabatech/canal
这里给出了如何快速启动Canal Server和Canal Client的例子
Quick Start
http://agapple.iteye.com/blog/1796070
https://github.com/alibabatech/canal/wiki/QuickStart
Client Example
http://agapple.iteye.com/blog/1796620
https://github.com/alibabatech/canal/wiki/ClientExample

