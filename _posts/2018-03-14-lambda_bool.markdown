---
title: lambda 演算中的布尔值和选择
layout: post
category: lang
author: 夏泽民
---
<!-- more -->
我们在lambda演算中引入了数字，只差两件事情就可以表达任意计算了：一个是如何表达选择（分支），另一个是如何表示重复。在这篇文章中，我将讨论布尔值和选择，下一篇将介绍重复和递归。

我们希望能够写出形如 if / then / else语句的表达式，就像我们在大多数编程语言做的那样。继像丘奇数那样将数字表示为函数之后，我们也将true和false值表示为对其参数执行一个if-then-else操作的函数：

let TRUE = lambda x y . x 
let FALSE = lambda x y . y 
于是，现在我们可以写一个“if”函数，它的第一个参数是一个条件表达式，第二个参数是如果条件为真时才进行运算的表达式，第三个参数则如果条件为假时要进行的运算。

let IfThenElse = lambda cond true_expr false_expr . cond true_expr false_expr 
此外我们还需要定义常用的逻辑运算：

let BoolAnd = lambda x y . x y FALSE 
let BoolOr = lambda x y. x TRUE y 
let BoolNot = lambda x . x FALSE TRUE 
现在，就让我们过一遍这些定义。让我们先看看BoolAnd：

BoolAnd TRUE FALSE，展开TRUE和FALSE定义：BoolAnd (lambda x y . x) (lambda x y . y)
alpha变换true和false：BoolAnd (lambda xt yt . xt) (lambda xf yf . yf)
现在，展开BoolAnd：(lambda x y. x y FALSE) (lambda xt yt . xt) (lambda xf yf . yf)
beta规约：(lambda xt yt.xt) (lambda xf yf. yf) FALSE
再次beta规约：(lambda xf xf . yf)
于是我们得到结果：BoolAnd TRUE FALSE = FALSE。再让我们来看看BoolAnd FALSE TRUE：

BoolAnd (lambda x y . y) (lambda x y .x)
alpha变换：BoolAnd (lambda xf yf . yf) (lambda xt yt . xt)
展开BoolAnd： (lambda x y .x y FALSE) (lambda xf yf . yf) (lambda xt yt . xt)
beta规约：(lambda xf yf . yf) (lambda xt yt . xt) FALSE
再beta规约：FALSE
所以，BoolAnd FALSE TRUE = FALSE

最后让我们来算算，BoolAnd TRUE TRUE：

展开两个TRUE： BoolAnd (lambda x y . x) (lambda x y . x)
alpha变换： BoolAnd (lambda xa ya . xa) (lambda xb yb . xb)
展开BoolAnd： (lambda x y . x y FALSE) (lambda xa ya . xa) (lambda xb yb . xb)
beta规约： (lambda xa ya . xa) (lambda xb yb . xb) FALSE
beta规约： (lambda xb yb .xb)
所以，BoolAnd TRUE TRUE = TRUE
