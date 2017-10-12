---
title: java_scala
layout: post
category: spark
author: 夏泽民
---
<!-- more -->

Java似乎可以无缝操纵Scala语言中定义的类，在trait那一节中我们提到，如果trait中全部是抽象成员，则它与java中的interface是等同的，这时候java可以把它当作接口来使用，但如果trait中定义了具体成员，则它有着自己的内部实现，此时在java中使用的时候需要作相应的调整。


Scala可以直接调用Java实现的任何类，只要符合scala语法就可以，不过某些方法在JAVA类中不存在，此时只要引入scala.collection.JavaConversions._包就可以了，它会我们自动地进行隐式转换，从而可以使用scala中的一些非常方便的高阶函数，如foreach方法,还可以显式地进行转换


Java中的泛型可以直接转换成Scala中的泛型，在前面的课程中我们已经有所涉及，例如Java中的Comparator<T> 可以直接转换成 Scala中的Comparator[T] 使用方法完全一样，不同的只是语法上的。

Scala中的异常处理是通过模式匹配来实现的
