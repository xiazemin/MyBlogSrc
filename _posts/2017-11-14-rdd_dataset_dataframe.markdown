---
title: RDD/Dataset/DataFrame互转
layout: post
category: spark
author: 夏泽民
---
<!-- more -->
1.RDD -> Dataset 
val ds = rdd.toDS()

2.RDD -> DataFrame 
val df = spark.read.json(rdd)

3.Dataset -> RDD
val rdd = ds.rdd

4.Dataset -> DataFrame
val df = ds.toDF()

5.DataFrame -> RDD
val rdd = df.toJSON.rdd

6.DataFrame -> Dataset
val ds = df.toJSON
