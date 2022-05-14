---
title: spark_memory
layout: post
category: spark
author: 夏泽民
---

在 spark-env.sh 中添加：

export JAVA_HOME=/home/tom/jdk1.8.0_73/
export SCALA_HOME=/home/tom//scala-2.10.6
export SPARK_MASTER_IP=localhost
export SPARK_WORKER_MEMORY=4G

内存设置太少会导致运行失败
可以代码中设置内存
//set new runtime options
spark.conf.set("spark.sql.shuffle.partitions", 1)
spark.conf.set("spark.executor.memory", "512m")

Spark虽然是in memory的运算平台，但从官方资料看，似乎本身对内存的要求并不是特别苛刻。官方网站只是要求内存在8GB之上即可（Impala要求机器配置在128GB）Spark建议需要提供至少75%的内存空间分配给Spark，至于其余的内存空间，则分配给操作系统与buffer cache。Spark对内存的消耗主要分为三部分（即取决于你的应用程序的需求）：
      数据集中对象的大小；
      访问这些对象的内存消耗；
      垃圾回收GC的消耗
      
      
spark executor都是装载在container里运行，container默认的内存是1G（参数yarn.scheduler.minimum-allocation-mb定义），executor分配的内存是executor-memory，所以向YARN申请的内存是（executor-memory + 1）* num-executors。 Executor 内存的大小，和性能本身当然并没有直接的关系，但是几乎所有运行时性能相关的内容都或多或少间接和内存大小相关。这个参数最终会被设置到Executor的JVM的heap尺寸上。如果Executor的数量和内存大小受机器物理配置影响相对固定，那么你就需要合理规划每个分区任务的数据规模，例如采用更多的分区，用增加任务数量（进而需要更多的批次来运算所有的任务）的方式来减小每个任务所需处理的数据大小。
