---
title: spark_start问题原因及解决办法
layout: post
category: spark
author: 夏泽民
---

<console>:14: error: not found: value spark
       import spark.implicits._
              ^
<console>:14: error: not found: value spark
       import spark.sql
              ^
              
scala>  var rdd = sc.parallelize(1 to 10)
<console>:39: error: not found: value sc
        var rdd = sc.parallelize(1 to 10)
        
        
日志：
Caused by: java.lang.RuntimeException: java.net.ConnectException: Call From localhost/127.0.0.1 to localhost:8020 failed on connection exception: java.net.ConnectException: Connection refused; For more details see:  http://wiki.apache.org/hadoop/ConnectionRefused

hadoop没有启动
启动hadoop

Using Scala version 2.11.8 (Java HotSpot(TM) 64-Bit Server VM, Java 1.8.0_144)
Type in expressions to have them evaluated.
Type :help for more information.

scala> var rdd = sc.parallelize(1 to 10)
rdd: org.apache.spark.rdd.RDD[Int] = ParallelCollectionRDD[0] at parallelize at <console>:24
