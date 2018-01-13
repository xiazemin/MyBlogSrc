---
title: spark on hive
layout: post
category: spark
author: 夏泽民
---
$cp hive/hive/conf/hive-site.xml spark/spark/conf/
<!-- more -->
启动hive
启动spark

import org.apache.spark.sql.SparkSession
val sparkSession = SparkSession.builder.master("local").enableHiveSupport().getOrCreate();

Caused by: ERROR XJ040: Failed to start database 'metastore_db' with class loader org.apache.spark.sql.hive.client.IsolatedClientLoader$$anon$1@294b045b, see the next exception for details.
	at org.apache.derby.iapi.error.StandardException.newException(Unknown Source)
	at org.apache.derby.impl.jdbc.SQLExceptionFactory.wrapArgsForTransportAcrossDRDA(Unknown Source)
	... 150 more
Caused by: ERROR XSDB6: Another instance of Derby may have already booted the database /Users/didi/metastore_db.

在使用Hive on Spark模式操作hive里面的数据时，报以上错误，原因是因为HIVE采用了derby这个内嵌数据库作为数据库，它不支持多用户同时访问,解决办法就是把derby数据库换成mysql数据库即可

解决方式：

不使用默认的内嵌数据库derby，采用mysql作为统计的存储信息。

修改相关配置信息（hive-site.xml）：

<property>

       <name>hive.stats.dbclass</name>

       <value>jdbc:mysql</value>

</property>

<property>

       <name>hive.stats.jdbcdriver</name>

       <value>com.mysql.jdbc.Driver</value>

</property>

<property>

       <name>hive.stats.dbconnectionstring</name>

       <value>jdbc:mysql://localhost:3306/TempStatsStore</value>

</property>

修改完成保存。
另外后面还有一个步骤就是要在mysql里创建TempStatsStore这个数据库（mysql里不会自动创建该库，在derby里会自动创建）

方式二
mv  metastore_db metastore_db_bak

scala> val sparkSession = SparkSession.builder.master("local").enableHiveSupport().getOrCreate();
18/01/13 15:55:44 WARN spark.SparkContext: Using an existing SparkContext; some configuration may not take effect.
18/01/13 15:55:45 WARN conf.HiveConf: HiveConf of name hive.conf.hidden.list does not exist
18/01/13 15:55:47 WARN conf.HiveConf: HiveConf of name hive.conf.hidden.list does not exist
18/01/13 15:55:49 WARN metastore.ObjectStore: Version information not found in metastore. hive.metastore.schema.verification is not enabled so recording the schema version 1.2.0
18/01/13 15:55:50 WARN metastore.ObjectStore: Failed to get database default, returning NoSuchObjectException
18/01/13 15:55:51 WARN metastore.ObjectStore: Failed to get database global_temp, returning NoSuchObjectException
18/01/13 15:55:51 WARN conf.HiveConf: HiveConf of name hive.conf.hidden.list does not exist
sparkSession: org.apache.spark.sql.SparkSession = org.apache.spark.sql.SparkSession@14c16388


beeline>  !connect jdbc:hive2://localhost:10000

0: jdbc:hive2://localhost:10000> show tables;
Error: Error while processing statement: FAILED: Execution Error, return code 1 from org.apache.hadoop.hive.ql.exec.DDLTask. MetaException(message:For direct MetaStore DB connections, we don't support retries at the client level.) (state=08S01,code=1)


scala> sparkSession.sql("CREATE TABLE IF NOT EXISTS src (key INT, value STRING) ROW FORMAT DELIMITED FIELDS TERMINATED BY '\t' ")
18/01/13 16:04:23 WARN metastore.HiveMetaStore: Location: hdfs://localhost:8020/user/hive/warehouse/src specified for non-external table:src
res8: org.apache.spark.sql.DataFrame = []


scala> sparkSession.sql("show tables").collect().foreach(println)
[default,src,false]

scala> sparkSession.sql("show databases").collect().foreach(println)
[default]

scala> sparkSession.sql("show create  table src").collect().foreach(println)
[CREATE TABLE `src`(`key` int, `value` string)
ROW FORMAT SERDE 'org.apache.hadoop.hive.serde2.lazy.LazySimpleSerDe'
WITH SERDEPROPERTIES (
  'field.delim' = '	',
  'serialization.format' = '	'
)
STORED AS
  INPUTFORMAT 'org.apache.hadoop.mapred.TextInputFormat'
  OUTPUTFORMAT 'org.apache.hadoop.hive.ql.io.HiveIgnoreKeyTextOutputFormat'
TBLPROPERTIES (
  'transient_lastDdlTime' = '1515830663'
)
]


   hive不支持用insert语句一条一条的进行插入操作，也不支持update操作。数据是以load的方式加载到建立好的表中。数据一旦导入就不可以修改。

DML包括：INSERT插入、UPDATE更新、DELETE删除


scala> sparkSession.sql("insert into table src key=1234,value='abc'").collect().foreach(println)
org.apache.spark.sql.catalyst.parser.ParseException:
no viable alternative at input 'key'(line 1, pos 22)

== SQL ==
insert into table src key=1234,value='abc'
----------------------^^^

 既然Hive没有行级别的数据插入、更新和删除操作，那么往表中装载数据的唯一途径就是使用一种”大量“的数据装载操作。我们以如下格式文件演示五种数据导入Hive方式    
 
 Tom         24    NanJing   Nanjing University  
Jack        29    NanJing   Southeast China University  
Mary Kake   21    SuZhou    Suzhou University  
John Doe    24    YangZhou  YangZhou University  
Bill King   23    XuZhou    Xuzhou Normal University 


数据格式以\t分隔，分别表示：姓名、年龄、地址、学校

一、从本地文件系统中导入数据
 (1) 创建test1测试表
scala> sparkSession.sql("CREATE TABLE test1(name STRING,age INT, address STRING,school STRING)   ROW FORMAT DELIMITED FIELDS TERMINATED BY '\t'  STORED AS TEXTFILE").collect().foreach(println)
18/01/13 16:20:50 WARN metastore.HiveMetaStore: Location: hdfs://localhost:8020/user/hive/warehouse/test1 specified for non-external table:test1
(2) 从本地加载数据 
scala> sparkSession.sql(" LOAD DATA LOCAL INPATH '/Users/didi/hive/testHive.txt' INTO TABLE test1").collect().foreach(println) 
(3) 查看导入结果
scala> sparkSession.sql("select * from test1").collect().foreach(println)
[Tom         24    NanJing   Nanjing University  ,null,null,null]
[Jack        29    NanJing   Southeast China University  ,null,null,null]
[Mary Kake   21    SuZhou    Suzhou University  ,null,null,null]
[John Doe    24    YangZhou  YangZhou University  ,null,null,null]
[Bill King   23    XuZhou    Xuzhou Normal University,null,null,null] 
        注意：此处使用的是LOCAL，表示从本地文件系统中加载数据到Hive中，同时没有OVERWRITE关键字，仅仅会把新增的文件增加到目标文件夹而不会删除之前的数据。如果使用OVERWRITE关键字，那么目标文件夹中之前的数据将会被先删除掉。
二、从HDFS文件系统加载数据到Hive
(1) 清空之前创建的表中数据
insert overwrite table test1  select * from test1 where 1=0;  //清空表，一般不推荐这样操作  
 (2) 从HDFS加载数据
hive> LOAD DATA INPATH "/input/test1.txt"  
    > OVERWRITE INTO TABLE test1;  
Loading data to table hive.test1  
rmr: DEPRECATED: Please use 'rm -r' instead.  
Deleted hdfs://secondmgt:8020/hive/warehouse/hive.db/test1  
Table hive.test1 stats: [numFiles=1, numRows=0, totalSize=201, rawDataSize=0]  
OK  
Time taken: 0.355 seconds  
 (3) 查询结果
hive> select * from test1;  
OK  
Tom     24.0    NanJing Nanjing University  
Jack    29.0    NanJing Southeast China University  
Mary Kake       21.0    SuZhou  Suzhou University  
John Doe        24.0    YangZhou        YangZhou University  
Bill King       23.0    XuZhou  Xuzhou Normal University  
Time taken: 0.054 seconds, Fetched: 5 row(s)  
        注意：此处没有LOCAL关键字，表示分布式文件系统中的路径，这就是和第一种方法的主要区别，同时由日志可以发现，因为此处加了OVERWRITE关键字，执行了Deleted操作，即先删除之前存储的数据，然后再执行加载操作。
       同时，INPATH子句中使用的文件路径还有一个限制，那就是这个路径下不可以包含任何文件夹。

三、通过查询语句向表中插入数据
(1) 创建test4测试表
scala> sparkSession.sql(" CREATE TABLE test4(name STRING,age FLOAT,address STRING,school STRING)   ROW FORMAT DELIMITED FIELDS TERMINATED BY '\t' STORED AS TEXTFILE").collect().foreach(println)
18/01/13 16:27:53 WARN metastore.HiveMetaStore: Location: hdfs://localhost:8020/user/hive/warehouse/test4 specified for non-external table:test4
 创建表过程基本和前面一样，此处不细讲
(2) 从查询结果中导入数据
scala> sparkSession.sql("INSERT INTO TABLE test4 SELECT * FROM test1").collect().foreach(println) 
        注意：新建表的字段数，一定要和后面SELECT中查询的字段数一样，且要注意数据类型。如test4包含四个字段：name、age、address和school，则SELECT查询出的结果也应该对应这四个字段。
(3) 查看导入结果
scala> sparkSession.sql("select * from test4").collect().foreach(println)
[Tom         24    NanJing   Nanjing University  ,null,null,null]
[Jack        29    NanJing   Southeast China University  ,null,null,null]
[Mary Kake   21    SuZhou    Suzhou University  ,null,null,null]
[John Doe    24    YangZhou  YangZhou University  ,null,null,null]
[Bill King   23    XuZhou    Xuzhou Normal University,null,null,null]
四、分区插入
        分区插入有两种，一种是静态分区，另一种是动态分区。如果混合使用静态分区和动态分区，则静态分区必须出现在动态分区之前。现分别介绍这两种分区插入
(1) 静态分区插入

①创建分区表
hive> CREATE TABLE test2(name STRING,address STRING,school STRING)  
    > PARTITIONED BY(age float)  
    > ROW FORMAT DELIMITED FIELDS TERMINATED BY '\t'  
    > STORED AS TEXTFILE ;  
OK  
Time taken: 0.144 seconds  
       此处创建了一个test2的分区表，以年龄分区

②从查询结果中导入数据      
hive> INSERT INTO  TABLE test2 PARTITION (age='24') SELECT * FROM test1;  
FAILED: SemanticException [Error 10044]: Line 1:19 Cannot insert into target table because column number/types are different ''24'': Table insclause-0 has 3 columns, but query has 4 columns.  
 此处报了一个错误。是因为test2中是以age分区的，有三个字段，SELECT * 语句中包含有四个字段，所以出错。正确如下：
[html] view plain copy
hive> INSERT INTO  TABLE test2 PARTITION (age='24') SELECT name,address,school FROM test1;  
 ③ 查看插入结果
hive> select * from test2;  
OK  
Tom     NanJing Nanjing University      24.0  
Jack    NanJing Southeast China University      24.0  
Mary Kake       SuZhou  Suzhou University       24.0  
John Doe        YangZhou        YangZhou University     24.0  
Bill King       XuZhou  Xuzhou Normal University        24.0  
Time taken: 0.079 seconds, Fetched: 5 row(s)  
 由查询结果可知，每条记录的年龄均为24，插入成功。
(2) 动态分区插入
静态分区需要创建非常多的分区，那么用户就需要写非常多的SQL！Hive提供了一个动态分区功能，其可以基于查询参数推断出需要创建的分区名称。

① 创建分区表，此过程和静态分区创建表一样，此处省略

② 参数设置
hive> set hive.exec.dynamic.partition=true;  
hive> set hive.exec.dynamic.partition.mode=nonstrict;  
 注意：动态分区默认情况下是没有开启的。开启后，默认是以”严格“模式执行的，在这种模式下要求至少有一列分区字段是静态的。这有助于阻止因设计错误导致查询产生大量的分区。但是此处我们不需要静态分区字段，估将其设为nonstrict。
③ 数据动态插入
hive> insert into table test2 partition (age) select name,address,school,age from test1;  
Total jobs = 1  
Launching Job 1 out of 1  
Number of reduce tasks not specified. Estimated from input data size: 1  
In order to change the average load for a reducer (in bytes):  
  set hive.exec.reducers.bytes.per.reducer=<number>  
In order to limit the maximum number of reducers:  
  set hive.exec.reducers.max=<number>  
In order to set a constant number of reducers:  
  set mapreduce.job.reduces=<number>  
Starting Job = job_1419317102229_0029, Tracking URL = http://secondmgt:8088/proxy/application_1419317102229_0029/  
Kill Command = /home/hadoopUser/cloud/hadoop/programs/hadoop-2.2.0/bin/hadoop job  -kill job_1419317102229_0029  
Hadoop job information for Stage-1: number of mappers: 1; number of reducers: 1  
2014-12-28 20:45:07,996 Stage-1 map = 0%,  reduce = 0%  
2014-12-28 20:45:21,488 Stage-1 map = 100%,  reduce = 0%, Cumulative CPU 2.67 sec  
2014-12-28 20:45:32,926 Stage-1 map = 100%,  reduce = 100%, Cumulative CPU 7.32 sec  
MapReduce Total cumulative CPU time: 7 seconds 320 msec  
Ended Job = job_1419317102229_0029  
Loading data to table hive.test2 partition (age=null)  
        Loading partition {age=29.0}  
        Loading partition {age=23.0}  
        Loading partition {age=21.0}  
        Loading partition {age=24.0}  
Partition hive.test2{age=21.0} stats: [numFiles=1, numRows=1, totalSize=35, rawDataSize=34]  
Partition hive.test2{age=23.0} stats: [numFiles=1, numRows=1, totalSize=42, rawDataSize=41]  
Partition hive.test2{age=24.0} stats: [numFiles=1, numRows=2, totalSize=69, rawDataSize=67]  
Partition hive.test2{age=29.0} stats: [numFiles=1, numRows=1, totalSize=40, rawDataSize=39]  
MapReduce Jobs Launched:  
Job 0: Map: 1  Reduce: 1   Cumulative CPU: 7.32 sec   HDFS Read: 415 HDFS Write: 375 SUCCESS  
Total MapReduce CPU Time Spent: 7 seconds 320 msec  
OK  
Time taken: 41.846 seconds  
 注意：查询语句select查询出来的age字段必须放在最后，和分区字段对应，不然结果会出错
④ 查看插入结果
[html] view plain copy
hive> select * from test2;  
OK  
Mary Kake       SuZhou  Suzhou University       21.0  
Bill King       XuZhou  Xuzhou Normal University        23.0  
John Doe        YangZhou        YangZhou University     24.0  
Tom     NanJing Nanjing University      24.0  
Jack    NanJing Southeast China University      29.0  
五、单个查询语句中创建表并加载数据
         在实际情况中，表的输出结果可能太多，不适于显示在控制台上，这时候，将Hive的查询输出结果直接存在一个新的表中是非常方便的，我们称这种情况为CTAS（create table .. as select）
        (1) 创建表
hive> CREATE TABLE test3  
    > AS  
    > SELECT name,age FROM test1;  
Total jobs = 3  
Launching Job 1 out of 3  
Number of reduce tasks is set to 0 since there's no reduce operator  
Starting Job = job_1419317102229_0030, Tracking URL = http://secondmgt:8088/proxy/application_1419317102229_0030/  
Kill Command = /home/hadoopUser/cloud/hadoop/programs/hadoop-2.2.0/bin/hadoop job  -kill job_1419317102229_0030  
Hadoop job information for Stage-1: number of mappers: 1; number of reducers: 0  
2014-12-28 20:59:59,375 Stage-1 map = 0%,  reduce = 0%  
2014-12-28 21:00:10,795 Stage-1 map = 100%,  reduce = 0%, Cumulative CPU 2.68 sec  
MapReduce Total cumulative CPU time: 2 seconds 680 msec  
Ended Job = job_1419317102229_0030  
Stage-4 is selected by condition resolver.  
Stage-3 is filtered out by condition resolver.  
Stage-5 is filtered out by condition resolver.  
Moving data to: hdfs://secondmgt:8020/hive/scratchdir/hive_2014-12-28_20-59-45_494_6763514583931347886-1/-ext-10001  
Moving data to: hdfs://secondmgt:8020/hive/warehouse/hive.db/test3  
Table hive.test3 stats: [numFiles=1, numRows=5, totalSize=63, rawDataSize=58]  
MapReduce Jobs Launched:  
Job 0: Map: 1   Cumulative CPU: 2.68 sec   HDFS Read: 415 HDFS Write: 129 SUCCESS  
Total MapReduce CPU Time Spent: 2 seconds 680 msec  
OK  
Time taken: 26.583 seconds  
 (2) 查看插入结果
hive> select * from test3;  
OK  
Tom     24.0  
Jack    29.0  
Mary Kake       21.0  
John Doe        24.0  
Bill King       23.0  
Time taken: 0.045 seconds, Fetched: 5 row(s)  
   


Exception in thread "main" java.lang.NoClassDefFoundError: scala/Product$class
	at org.apache.spark.internal.config.ConfigBuilder.<init>(ConfigBuilder.scala:176)
	at org.apache.spark.sql.internal.SQLConf$.buildConf(SQLConf.scala:58)
	at org.apache.spark.sql.internal.SQLConf$.<init>(SQLConf.scala:67)
	at org.apache.spark.sql.internal.SQLConf$.<clinit>(SQLConf.scala)
	at org.apache.spark.sql.internal.StaticSQLConf$.<init>(StaticSQLConf.scala:31)
	at org.apache.spark.sql.internal.StaticSQLConf$.<clinit>(StaticSQLConf.scala)
	at org.apache.spark.sql.SparkSession$Builder.enableHiveSupport(SparkSession.scala:843)
	at main.scala.hiveConnection$.main(hiveConnection.scala:6)
	at main.scala.hiveConnection.main(hiveConnection.scala)
	

在使用Log4j时若提示如下信息：
log4j:WARN No appenders could be found for logger (org.apache.ibatis.logging.LogFactory).  
log4j:WARN Please initialize the log4j system properly. 
则，解决办法为：在项目的src下面新建file名为log4j.properties文件，内容如下: 

# Configure logging for testing: optionally with log file
#可以设置级别：debug>info>error
#debug:可以显式debug,info,error
#info:可以显式info,error
#error:可以显式error

log4j.rootLogger=debug,appender1
#log4j.rootLogger=info,appender1
#log4j.rootLogger=error,appender1

#输出到控制台
log4j.appender.appender1=org.apache.log4j.ConsoleAppender
#样式为TTCCLayout
log4j.appender.appender1.layout=org.apache.log4j.TTCCLayout

然后，存盘退出。再次运行程序就会显示Log信息了。


通过配置文件可知，我们需要配置3个方面的内容：
1、根目录（级别和目的地）；
2、目的地（控制台、文件等等）；
3、输出样式。



或者，使用下面的内容也可以：
 # Configure logging for testing: optionally with log file
log4j.rootLogger=WARN, stdout
 # log4j.rootLogger=WARN, stdout, logfile
log4j.appender.stdout=org.apache.log4j.ConsoleAppender
log4j.appender.stdout.layout=org.apache.log4j.PatternLayout
log4j.appender.stdout.layout.ConversionPattern=%d %p [%c] - %m%n
log4j.appender.logfile=org.apache.log4j.FileAppender
log4j.appender.logfile.File=target/spring.log
log4j.appender.logfile.layout=org.apache.log4j.PatternLayout
log4j.appender.logfile.layout.ConversionPattern=%d %p [%c] - %m%n

