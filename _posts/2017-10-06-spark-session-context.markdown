---
title: spark-session-context
layout: post
category: spark
author: 夏泽民
---
<!-- more -->
初始化Spark
一个Spark程序首先必须要做的是创建一个SparkContext对象，这个对象告诉Spark如何访问一个集群。为了创建一个SparkContext，你首先需要构建一个包含了关于你的应用的信息的SparkConf对象。

只有一个SparkContext可以激活一个JVM。你必须在创建一个新的SparkContext之前stop()掉这个激活的SparkContext。

val conf = new SparkConf().setAppName(appName).setMaster(master)
new SparkContext(conf)
参数appName是在集群UI上展示的你的应用的名字。master是一个Spark，Mesos或者YARN集群的URL，或者是表示以本地模式运行的一个特定的字符串”local”。在实践中，当运行在集群上时，你是不会想在程序中硬编码master的，而是通过spark-submit启动应用并指定master的。然而，对于本地测试和单元测试，你可以通过传递”local”来在进程内运行Spark。

使用shell
在Spark shell中，一个特殊的SparkContext已经为你创建好了，变量名是sc。如果再创建你自己的SparkContext就不起作用了。你可以使用参数--master来设置context连接到哪个master，还可以通过给参数--jars传递逗号分隔的列表来给classpath增加JARs。你还可以通过给参数--packages提供逗号分隔的Maven坐标来给你的shell会话增加依赖（例如Spark包）。任何可能存在依赖的附加库（例如Sonatype）都可以被传递给参数--repositories

　在Spark的早期版本，sparkContext是进入Spark的切入点。我们都知道RDD是Spark中重要的API，然而它的创建和操作得使用sparkContext提供的API；对于RDD之外的其他东西，我们需要使用其他的Context。比如对于流处理来说，我们得使用StreamingContext；对于SQL得使用sqlContext；而对于hive得使用HiveContext。然而DataSet和Dataframe提供的API逐渐称为新的标准API，我们需要一个切入点来构建它们，所以在 Spark 2.0中我们引入了一个新的切入点(entry point)：SparkSession

　　SparkSession实质上是SQLContext和HiveContext的组合（未来可能还会加上StreamingContext），所以在SQLContext和HiveContext上可用的API在SparkSession上同样是可以使用的。SparkSession内部封装了sparkContext，所以计算实际上是由sparkContext完成的。

创建SparkSession

　　SparkSession的设计遵循了工厂设计模式（factory design pattern），下面代码片段介绍如何创建SparkSession
[python] view plain copy
val sparkSession = SparkSession.builder.  
      master("local")  
      .appName("spark session example")  
      .getOrCreate()  
上面代码类似于创建一个SparkContext，master设置为local，然后创建了一个SQLContext封装它。如果你想创建hiveContext，可以使用下面的方法来创建SparkSession，以使得它支持Hive：
[python] view plain copy
val sparkSession = SparkSession.builder.  
      master("local")  
      .appName("spark session example")  
      .enableHiveSupport()  
      .getOrCreate()  
enableHiveSupport 函数的调用使得SparkSession支持hive，类似于HiveContext。

spark2.0 主要变化
1 更容易的SQL和Streamlined APIs

   Spark 2.0主要聚焦于两个方面：（1）、对标准的SQL支持（2）、统一DataFrame和Dataset API。

　　在SQL方面，Spark 2.0已经显著地扩大了它的SQL功能，比如引进了一个新的ANSI SQL解析器和对子查询的支持。现在Spark 2.0已经可以运行TPC-DS所有的99个查询，这99个查询需要SQL 2003的许多特性。因为SQL是Spark应用程序的主要接口之一，Spark 2.0 SQL的扩展大幅减少了应用程序往Spark迁移的代价。

　　在编程API方面，我们对API进行了精简。

　　1、统一Scala和Java中DataFrames和Datasets的API：从Spark 2.0开始，DataFrame仅仅是Dataset的一个别名。有类型的方法(typed methods)（比如：map, filter, groupByKey）和无类型的方法(untyped methods)(比如：select, groupBy)目前在Dataset类上可用。同样，新的Dataset接口也在Structured Streaming中使用。因为编译时类型安全(compile-time type-safety)在Python和R中并不是语言特性，所以Dataset的概念并不在这些语言中提供相应的API。而DataFrame仍然作为这些语言的主要编程抽象。

　　2、SparkSession：一个新的切入点，用于替代旧的SQLContext和HiveContext。对于那些使用DataFrame API的用户，一个常见的困惑就是我们正在使用哪个context？现在我们可以使用SparkSession了，其涵括了SQLContext和HiveContext，仅仅提供一个切入点。需要注意的是为了向后兼容，旧的SQLContext和HiveContext目前仍然可以使用。

　　3、简单以及性能更好的Accumulator API：Spark 2.0中设计出一种新的Accumulator API，它拥有更加简洁的类型层次，而且支持基本类型。为了向后兼容，旧的Accumulator API仍然可以使用。

　　4、基于DataFrame的Machine Learning API可以作为主要的ML API了：在Spark 2.0中， spark.ml包以其pipeline API将会作为主要的机器学习API了，而之前的spark.mllib仍然会保存，将来的开发会聚集在基于DataFrame的API上。

　　5、Machine learning pipeline持久化：现在用户可以保存和加载Spark支持所有语言的Machine learning pipeline和models。

　　6、R的分布式算法：在R语言中添加支持了Generalized Linear Models (GLM), Naive Bayes, Survival Regression, and K-Means。
　　
2 更快：Spark作为编译器

Spark 2.0中附带了第二代Tungsten engine，这一代引擎是建立在现代编译器和MPP数据库的想法上，并且把它们应用于数据的处理过程中。主要想法是通过在运行期间优化那些拖慢整个查询的代码到一个单独的函数中，消除虚拟函数的调用以及利用CPU寄存器来存放那些中间数据。我们把这些技术称为"整段代码生成"(whole-stage code generation)。

3 更加智能：Structured Streaming


