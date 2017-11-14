---
title: spark toDF 失败原因总结
layout: post
category: spark
author: 夏泽民
---
<!-- more -->
错误提示
value toDF is not a member of org.apache.spark.rdd.RDD[(org.apache.spark.ml.

解决办法
 val conf = new SparkConf().setAppName("SimpleParamsExample1")
    val sc = new SparkContext(conf)
  
  val sqlContext= new org.apache.spark.sql.SQLContext(sc)
  import sqlContext.implicits._ 
  
  
错误: 找不到或无法加载主类 example.Statistics
译器顺序：右键项目-properties-scala Compiler -Build manager ：
 set the compile order to JavaThenScala instead of Mixed
 
 
 右键项目-properties-scala Compiler -Standard 
 选择安装的scala 版本
