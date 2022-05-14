---
title: scala tuple
layout: post
category: spark
author: 夏泽民
---
<!-- more -->
<div>
Scala元组将固定数量的项目组合在一起，以便它们可以作为一个整体传递。 与数组或列表不同，元组可以容纳不同类型的对象，但它们也是不可变的。
val t = (1, "hello", Console)
val t = new Tuple3(1, "hello", Console)
元组的实际类型取决于它包含的数量和元素以及这些元素的类型。 因此，(99，"Luftballons")的类型是Tuple2 [Int，String]

元组是类型Tuple1，Tuple2，Tuple3等等。目前在Scala中只能有22个上限，如果您需要更多个元素，那么可以使用集合而不是元组。 对于每个TupleN类型，其中上限为1 <= N <= 22，Scala定义了许多元素访问方法。给定以下定义 -
val t = (4,3,2,1)
要访问元组t的元素，可以使用t._1方法访问第一个元素，t._2方法访问第二个元素

scala>  val t2=("test",1)
t2: (String, Int) = (test,1)

scala> t2.getClass
res0: Class[_ <: (String, Int)] = class scala.Tuple2

scala> val t1=("test")
t1: String = test

scala> t1.getClass
res3: Class[_ <: String] = class java.lang.String

scala> val t1=(Tuple1)("test")
t1: (String,) = (test,)

scala> t1.getClass
res8: Class[_ <: (String,)] = class scala.Tuple1

tuple1.apply 的作用：将任何类型的元素装箱为tuple1的对象，可以toDF（）转换了
{% highlight scala %}
 val data = Seq(
      Vectors.sparse(4, Seq((0, 1.0), (3, -2.0))),
      Vectors.dense(4.0, 5.0, 0.0, 3.0),
      Vectors.dense(6.0, 7.0, 0.0, 8.0),
      Vectors.sparse(4, Seq((0, 9.0), (3, 1.0)))
    )
    data.foreach(println(_))
    data.map(Tuple1.apply).foreach(println(_))
    val df = data.map(Tuple1.apply).toDF("features")
    df.show(false)
{% endhighlight %}  
装箱前
(4,[0,3],[1.0,-2.0])
[4.0,5.0,0.0,3.0]
[6.0,7.0,0.0,8.0]
(4,[0,3],[9.0,1.0])
装箱后
((4,[0,3],[1.0,-2.0]))
([4.0,5.0,0.0,3.0])
([6.0,7.0,0.0,8.0])
((4,[0,3],[9.0,1.0]))
df显示
+--------------------+
|features            |
+--------------------+
|(4,[0,3],[1.0,-2.0])|
|[4.0,5.0,0.0,3.0]   |
|[6.0,7.0,0.0,8.0]   |
|(4,[0,3],[9.0,1.0]) |
+--------------------+
</div>
