---
title: scala_partion_function
layout: post
category: spark
author: 夏泽民
---
<!-- more -->
从使用case语句构造匿名函数谈起

在Scala里，我们可以使用case语句来创建一个匿名函数（函数字面量），这有别于一般的匿名函数创建方法。来看个例子：

scala> List(1,2,3) map {case i:Int=>i+1}
res1: List[Int] = List(2, 3, 4)
这很有趣，case i:Int=>i+1构建的匿名函数等同于(i:Int)=>i+1，也就是下面这个样子：

scala> List(1,2,3) map {(i:Int)=>i+1}
res2: List[Int] = List(2, 3, 4)
《Scala In Programming》一书对独立的case语句作为匿名函数（函数字面量）有权威的解释：

Essentially, a case sequence is a function literal, only more general. Instead of having a single entry point and list of parameters, a case sequence has multiple entry points, each with their own list of parameters. Each case is an entry point to the function, and the parameters are specified with the pattern. 
一个case语句就是一个独立的匿名函数，如果有一组case语句的话，从效果上看，构建出的这个匿名函数会有多种不同的参数列表，每一个case对应一种参数列表，参数是case后面的变量声明，其值是通过模式匹配赋予的。

使用case语句构造匿名函数的“额外”好处

使用case语句构造匿名函数是有“额外”好处的，这个“好处”在下面这个例子中得到了充分的体现：

List(1, 3, 5, "seven") map { case i: Int => i + 1 } // won't work
// scala.MatchError: seven (of class java.lang.String)
List(1, 3, 5, "seven") collect { case i: Int => i + 1 }
// verify
assert(List(2, 4, 6) == (List(1, 3, 5, "seven") collect { case i: Int => i + 1 }))
在这个例子中：传递给map的case语句构建的是一个普通的匿名函数，在把这个函数适用于”seven”元素时发生了类型匹配错误。而对于collect,它声明接受的是一个偏函数：PartialFunction，传递的case语句能良好的工作说明这个case语句被编译器自动编译成了一个PartialFunction！这就是case语句“额外”的好处：case语句（组合）除了可以被编译为匿名函数（类型是FunctionX，在Scala里，所有的函数字面量都是一个对象，这个对象的类型是FunctionX），还可以非常方便的编译为一个偏函数PartialFunction！（注意：PartialFunction同时是Function1的子类）编译器会根据调用处的函数类型声明自动帮我们判定如何编译这个case语句（组合）。

上面我们直接抛出了偏函数的概念，这会让人头晕，我们可以只从collect这个示例的效果上去理解偏函数：它只对会作用于指定类型的参数或指定范围值的参数实施计算，超出它的界定范围之外的参数类型和值它会忽略（未必会忽略，这取决于你打算怎样处理）。就像上面例子中一样，case i: Int => i + 1只声明了对Int参数的处理，在遇到”seven”元素时，不在偏函数的适用范围内，所以这个元素被忽略了。

正式认识偏函数Partial Function

如同在一开始的例子中那样，我们手动实现了一个与case i:Int=>i+1等价的那个匿名函数(i:Int)=>i+1,那么在上面的collect方法中使用到的case i: Int => i + 1它的等价函数是什么呢？显然，不可能是(i:Int)=>i+1了，因为我们已经解释了，collect接受的参数类型是PartialFunction[Any,Int],而不是(Int)=>Int。 那个case语句对应的偏函数具体是什么样的呢？来看：

scala> val inc = new PartialFunction[Any, Int] {
     | def apply(any: Any) = any.asInstanceOf[Int]+1
     | def isDefinedAt(any: Any) = if (any.isInstanceOf[Int]) true else false
     | }
inc: PartialFunction[Any,Int] = <function1>

scala> List(1, 3, 5, "seven") collect inc
res4: List[Int] = List(2, 4, 6)
PartialFunction特质规定了两个要实现的方法：apply和isDefinedAt，isDefinedAt用来告知调用方这个偏函数接受参数的范围，可以是类型也可以是值，在我们这个例子中我们要求这个inc函数只处理Int型的数据。apply方法用来描述对已接受的值如何处理，在我们这个例子中，我们只是简单的把值+1，注意，非Int型的值已被isDefinedAt方法过滤掉了，所以不用担心类型转换的问题。

上面这个例子写起来真得非常笨拙，和前面的case语句方式比起来真是差太多了。这个例子从反面展示了：通过case语句组合去是实现一个偏函数是多么简洁。实际上case语句组合与偏函数的用意是高度贴合的，所以使用case语句组合是最简单明智的选择，同样是上面的inc函数，换成case去写如下：

scala> def inc: PartialFunction[Any, Int] =
     | { case i: Int => i + 1 }
inc: PartialFunction[Any,Int]

scala> List(1, 3, 5, "seven") collect inc
res5: List[Int] = List(2, 4, 6)

当然，如果偏函数的逻辑非常复杂，可能通过定义一个专门的类并继承PartialFunction是更好选择。

Case语句是如何被编译成偏函数的

关于这个问题在《Programming In Scala》中有较为详细的解释。对于这样一个使用case写在的偏函数：

val second: PartialFunction[List[Int],Int] = {
    case x :: y :: _ => y
}
In fact, such an expression gets ranslated by the Scala compiler to a partial function by translating the patterns twice—once for the implementation of the real function, and once to test whether the function is defined or not. For instance, the function literal { case x :: y :: _ => y }above gets translated to the following partialfunction value:

new PartialFunction[List[Int], Int] {
    def apply(xs: List[Int]) = xs match {
        case x :: y :: _ => y
    }
    def isDefinedAt(xs: List[Int]) = xs match {
        case x :: y :: _ => true
        case _ => false
    }
}
为什么偏函数需要抽象成一个专门的Trait

首先，在Scala里，一切皆对象，函数字面量（匿名函数）也不例外！这也是为什么我们可以把函数字面量赋给一个变量的原因, 是对象就有对应的类型，那么一个函数字面量的真实类型是什么呢？看下面这个例子：

scala> var inc = (x: Int) => x + 1
inc: Int => Int = <function1>

scala> inc.isInstanceOf[Function1[Int,Int]]
res0: Boolean = true
在Scala的scala包里，有一系列Function trait，它们实际上就是函数字面量作为“对象”存在时对应的类型。Function类型有多个版本，Function0表示无参数函数，Function1表示只有一个参数的函数，以此类推。至此我们解释的是一个普遍性问题：是函数就是对象，是对象就有类型。那么，接下来我们看一下偏函数又应该是什么样的一种“类型”？

从语义上讲，偏函数区别于普通函数的唯一特征就是：偏函数会自主地告诉调用方它的处理参数的范围，范围既可是值也可以是类型。针对这样的场景，我们需要给函数安插一种明确的“标识”，告诉编译器：这个函数具有这种特征。所以特质PartialFunction就被创建出来用于“标记”这类函数的，这个特质最主要的方法就是isDefinedAt！同时你也记得PartialFunction还是Function1的子类，所以它也要有apply方法，这是非常自然的，偏函数本身首先是一个函数嘛。

从另一个角度思考，偏函数的逻辑是可以通过普通函数去实现的，只是偏函数是更为优雅的一种方式，同时偏函数特质PartialFunction的存在对调用方和实现方都是一种语义更加丰富的约定，比如collect方法声明使用一个偏函数就暗含着它不太可能对每一个元素进行操作，它的返回结果仅仅是针对偏函数“感兴趣”的元素计算出来的

为什么偏函数只能有一个参数？

为什么只有针对单一参数的偏函数，而不是像Function特质那样，拥有多个版本的PartialFunction呢？在刚刚接触偏函数时，这也让我感到费解，但看透了偏函数的实质之后就会觉得很合理了。我们说所谓的偏函数本质上是由多个case语句组成的针对每一种可能的参数分别进行处理的一种“结构较为特殊”的函数，那特殊在什么地方呢？对，就是case语句，前面我们提到，case语句声明的变量就是偏函数的参数，既然case语句只能声明一个变量，那么偏函数受限于此，也只能有一个参数！说到底，类型PartialFunction无非是为由一组case语句描述的函数字面量提供一个类型描述而已，case语句只接受一个参数，则偏函数的类型声明自然就只有一个参数。

但是，上这并不会对编程造成什么阻碍，如果你想给一个偏函数传递多个参数，完全可以把这些参数封装成一个Tuple传递过去