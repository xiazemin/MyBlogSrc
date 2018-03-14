---
title: 从Lambda演算到组合子演算
layout: post
category: lang
author: 夏泽民
---
<!-- more -->

让我们来看看三个简单的组合子：

S：S是一个函数应用组合子： S = lambda x y z . (x z (y z))
K：K生成一个返回特定常数值的函数： K = lambda x . (lambda y . x)。 （即扔掉第二个参数，返回第一个参数）
I：恒等函数： I = lambda x . x
乍一看，这是一个很奇怪的组合。S的应用机制尤为奇怪 —— 它并不是接受两个参数x和y，并应用x到y，它除了x和y外还用到了第三个值z，先将x应用到z上，再将y应用到z上，最后用前者的结果应用到了后者的结果上。

这是有道理的。以下各行各做了一步规约：

S K K x = 
(K x) (K x) = 
x 
噗！ 我们根本用不着I。我们仅用S和K就创建了I的等价。但是，这仅仅是个开始：事实上，我们可以只用S和K组合子，甚至一个变量都不用，创建任意lambda演算表达式的等价。

例如，Y组合子可以写成：

Y = S S K (S (K (S S (S (S S K)))) K) 
在我们继续深入之前，有一个重要的事情要指出。我在上面说的是，使用S K K，我们创建了I的等价，然而它并没有规约为lambda x . x。

到目前为止，我们说在Lambda演算中，“x = y”，当且仅当x和y相同，或通过Alpha转化后相同。（这样lambda x y . x + y等于lambda a b . a + b ，但不等于lambda x y . y + x ）这就是所谓的内涵等价(intensional equivalence) 。 然而，另一种相等也非常有用，这就是所谓的外延等价（extensional equivalence）或外延相等（extensional equality）。外延相等时，表达式X等于一个表达式Y，当且仅当X等同Y（模Alpha），或者 for all a . X a = Y a。

从现在起，我们使用「=」表示外延相等。我们可以将任何 Lambda表达式转换为外延相等的组合子形式。我们定义一个从Lambda形式到组合子形式的变换函数C：

C{x} = x
C{E1 E2} = C{E1} C{E2}
C{lambda x . E} = K C{E}，如果x在E中非自由
C{lambda x . x} = I
C{lambda x . E1 E2} = (S C{lambda x . E1} C {lambda x . E2})
C{lambda x . (lambda y . E)} = C {lambda x . C {lambda y . E}}，如果x在E中是自由变量
让我们演进一下 C{lambda x y . y x} ：

柯里化函数： C{lambda x . (lambda y . y x)}
根据规则6： C{lambda x . C{lambda y . y x}}
根据规则5： C{lambda x . S C{lambda y . y} C{lambda y . x}}
根据规则4： C{lambda x . S I C{lambda y . x}}
根据规则3： C{lambda x . S I (K C{x})}
通过规则1： C{lambda x . S I (K x)}
根据规则5： S C{lambda x . S I} C{lambda x . (K x)}
根据规则3： S (K (S I)) C{lambda x . K x}
根据规则5： S (K (S I)) (S C{lambda x . K} C{lambda x . x})
通过规则1： S (K (S I)) (S C{lambda x . K} I)
根据规则3： S (K (S I)) (S (K K) I)
现在，让我们尝试使用“x”和“y”作为参数传递给该组合子表达式，并规约：

S (K (S I)) (S (K K) I) x y
让我们创建一些别名，以方便阅读：A = (K (S I)), B = (S (K K) I)，所以我们的表达式现在成了：S A B x y
展开S: (A x (B x)) y
让我们去掉别名B：(A x ((S (K K) I) x)) y
现在让我们去掉S：(A x ((K K) x (I x))) y
以及I：(A x ((K K) x x)) y
规约(K K) x ：(A x (K x)) y
展开别名A： ((K (S I)) x (K x)) y
规约(K (S I)) x ，得到： ((S I) (K x)) y
规约S：I y (K x) y
规约I：y (K x) y
最后规约(K x) y，剩下：y x