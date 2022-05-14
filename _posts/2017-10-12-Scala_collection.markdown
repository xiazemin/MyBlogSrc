---
title: Scala_collection
layout: post
category: spark
author: 夏泽民
---
<!-- more -->
Scala 集合分为可变的和不可变的集合。

1	Scala List(列表)
List的特征是其元素以线性方式存储，集合中可以存放重复对象。
参考 API文档
2	Scala Set(集合)
Set是最简单的一种集合。集合中的对象不按特定的方式排序，并且没有重复对象。
参考 API文档
3	Scala Map(映射)
Map 是一种把键对象和值对象映射的集合，它的每一个元素都包含一对键对象和值对象。
参考 API文档
4	Scala 元组
元组是不同类型的值的集合
5	Scala Option
Option[T] 表示有可能包含值的容器，也可能不包含值。
6	Scala Iterator（迭代器）
迭代器不是一个容器，更确切的说是逐一访问容器内元素的方法。


　Scala中的三种集合类型包含:Array,List,Tuple．那么究竟这三种有哪些异同呢？说实话，我之前一直没弄明确，所以今天特意花了点时间学习了一下．

　　　　同样点:
　　　　　1.长度都是固定的，不可变长
　　　　　２.早期的Scala版本号,Array、List都不能混合类型，仅仅有Tuple能够,2.8版本号以后,3者的元素都能够混合不同的类型（转化为Any类型）

　　　　不同点:
　　　　　1.Array 中的元素值可变，List和Tuple中的元素值不可变
　　　　　２.Array通常是先确定长度，后赋值，而List和Tuple在声明的时候就须要赋值
　　　　　３.Array取单个元素的效率非常高。而List读取单个元素的效率是O(n)
　　　　　4.List和Array的声明不须要newkeyword。而Tuple声明无论有无new 都能够


          val arrayTest = Array(1,2,3,4)   //正确
          val arrayTest = Array(1,2,3,4)   //错误<span style="font-family: Arial, Helvetica, sans-serif;">  </span>
          val listTest = List(1,2,3,4)         //正确
          val listTest = new List(1,2,3,4)    //错误

          val tupleTest = Tuple(1,2,"aaa")        //正确
          val tupleTest = new Tuple(1,2,"aaa")    //正确
          val tupleTest = (1,2,"aaa")             //正确
　　
　　　　　5.当使用混合类型时，Array和List会将元素类型转化为Any类型,而Tuple则保留每个元素的初始类型

                    6.訪问方式不同。Array和List的下标从0開始，且使用小括号,而Tuple的下标从1開始，切使用点加下划线的方式訪问，如：arrayTest(0), listTest(0); Tuple訪问: tupleTest._1


 List 有个叫“ ::: ” 的方法实现叠加功能。你可以这么用：

val oneTwo = List(1, 2)
val threeFour = List(3, 4)

val oneTwoThreeFour = oneTwo ::: threeFour

 List 最常用的操作符是发音为“ cons ” 的‘ :: 。 Cons 把一个新元素组合到已有 List 的最前端，然后返回结果 List 。例如，若执行这个脚本：
val twoThree = list(2, 3)
val oneTwoThree = 1 :: twoThree

表达式“ 1 :: twoThree ” 中， :: 是它右操作数，列表 twoThree ，的方法。你或许会疑惑 :: 方法的关联性上有什么东西搞错了，不过这只是一个简单的需记住的规则：如果一个方法被用作操作符标注，如 a * b ，那么方法被左操作数调用，就像 a.*(b) 除非方法名以冒号结尾。这种情况下，方法被右操作数调用。因此， 1 :: twoThree 里， :: 方法被 twoThree 调用，传入 1 ，像这样： twoThree.::(1) 。



