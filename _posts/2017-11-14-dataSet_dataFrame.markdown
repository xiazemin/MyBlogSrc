---
title: dataSet和dataFrame的创建方法
layout: post
category: spark
author: 夏泽民
---
Spark创建DataFrame的三种方法
跟关系数据库的表(Table)一样，DataFrame是Spark中对带模式(schema)行列数据的抽象。DateFrame广泛应用于使用SQL处理大数据的各种场景。

通过导入(importing)Spark sql implicits, 就可以将本地序列(seq), 数组或者RDD转为DataFrame。只要这些数据的内容能指定数据类型即可。
本地seq + toDF创建DataFrame示例：

import sqlContext.implicits._
val df = Seq(
  (1, "First Value", java.sql.Date.valueOf("2010-01-01")),
  (2, "Second Value", java.sql.Date.valueOf("2010-02-01"))
).toDF("int_column", "string_column", "date_column")
注意：如果直接用toDF()而不指定列名字，那么默认列名为"_1", "_2", ...

通过case class + toDF创建DataFrame的示例

// sc is an existing SparkContext.
val sqlContext = new org.apache.spark.sql.SQLContext(sc)
// this is used to implicitly convert an RDD to a DataFrame.
import sqlContext.implicits._

// Define the schema using a case class.
// Note: Case classes in Scala 2.10 can support only up to 22 fields. To work around this limit,
// you can use custom classes that implement the Product interface.
case class Person(name: String, age: Int)

// Create an RDD of Person objects and register it as a table.
val people = sc.textFile("examples/src/main/resources/people.txt").map(_.split(",")).map(p => Person(p(0), p(1).trim.toInt)).toDF()
people.registerTempTable("people")

// 使用 sqlContext 执行 sql 语句.
val teenagers = sqlContext.sql("SELECT name FROM people WHERE age >= 13 AND age <= 19")

// 注：sql()函数的执行结果也是DataFrame，支持各种常用的RDD操作.
// The columns of a row in the result can be accessed by ordinal.
teenagers.map(t => "Name: " + t(0)).collect().foreach(println)

方法二，Spark中使用createDataFrame函数创建DataFrame

在SqlContext中使用createDataFrame也可以创建DataFrame。跟toDF一样，这里创建DataFrame的数据形态也可以是本地数组或者RDD。
通过row+schema创建示例

import org.apache.spark.sql.types._
val schema = StructType(List(
    StructField("integer_column", IntegerType, nullable = false),
    StructField("string_column", StringType, nullable = true),
    StructField("date_column", DateType, nullable = true)
))

val rdd = sc.parallelize(Seq(
  Row(1, "First Value", java.sql.Date.valueOf("2010-01-01")),
  Row(2, "Second Value", java.sql.Date.valueOf("2010-02-01"))
))
val df = sqlContext.createDataFrame(rdd, schema)
方法三，通过文件直接创建DataFrame

使用parquet文件创建

val df = sqlContext.read.parquet("hdfs:/path/to/file")
使用json文件创建

val df = spark.read.json("examples/src/main/resources/people.json")

// Displays the content of the DataFrame to stdout
df.show()
// +----+-------+
// | age|   name|
// +----+-------+
// |null|Michael|
// |  30|   Andy|
// |  19| Justin|
// +----+-------+
使用csv文件,spark2.0+之后的版本可用

//首先初始化一个SparkSession对象
val spark = org.apache.spark.sql.SparkSession.builder
        .master("local")
        .appName("Spark CSV Reader")
        .getOrCreate;

//然后使用SparkSessions对象加载CSV成为DataFrame
val df = spark.read
        .format("com.databricks.spark.csv")
        .option("header", "true") //reading the headers
        .option("mode", "DROPMALFORMED")
        .load("csv/file/path"); //.csv("csv/file/path") //spark 2.0 api

df.show()


DataSet数据集是一个强类型的域特定对象的集合，可以使用功能或关系操作并行转换.。每个数据集还有一个无类型的视图称为Dataframe，这是一个行(Row)的数据集。
在DataSet上的操作，分为transformations和actions。transformations会产生新的数据集（DataSet），而actions则是触发计算并产生结果。transformations包括：map, filter, select, and aggregate (`groupBy`). 等操作。而actions 包括： count, show 或把数据写入文件系统中。
DataSet是懒惰(‘lazy’)的，也就是说当一个action被调用时才会触发一个计算。在内部实现，数据集表示的是一个逻辑计划，它描述了生成数据所需的计算。当action被调用时，spark的查询优化器会优化这个逻辑计划，并生成一个物理计划，该物理计划可以通过并行和分布式的方式来执行。使用`explain`解释函数，来进行逻辑计划的探索和物理计划的优化。
为了有效地支持特定领域的对象，Encoder（编码器）是必需的。例如，给出一个Person的类，有两个字段：name(string)和age(int)，通过一个encoder来告诉spark在运行的时候产生代码把Person对象转换成一个二进制结构。这种二进制结构通常有更低的内存占用，以及优化的数据处理效率（例如在一个柱状格式）。若要了解数据的内部二进制表示，请使用schema(表结构)函数。

创建一个数据集有两种方式。
1, 第一种也是最常用的一种，就是利用SparkSession的read函数在存储系统中读取一个文件。代码如下：
val people = spark.read.parquet("...").as[Person]

2, 利用Datasets的转换(transformations)函数，例如下面的代码：利用filter来创建一个新的Dataset。


数据集(Dataset)的操作是无类型的，通过各种DSL(domain-specific-language)函数，这些函数是基于数据集Dataset , 类[[Column]],和 函数[[functions]]来定义的。这些操作非常类似基于R和Python抽象出来的data frame的操作。
从Dateset中选择一列
val ageCol = people("age") // in Scala
