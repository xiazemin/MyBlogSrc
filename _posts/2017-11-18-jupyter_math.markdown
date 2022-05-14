---
title: jupyter 数学公式
layout: post
category: spark
author: 夏泽民
---
<!-- more -->
<div class="container">

<div class="row">
1、数学公式的前后要加上 $ 或 \( 和 \)，比如：$f(x) = 3x + 7$ 和 f(x)=3x+7 效果是一样的；如果用 \[ 和 \]，或者使用 $$ 和 $$，则该公式独占一行；如果用 \begin{equation} 和 \end{equation}，则公式除了独占一行还会自动被添加序号， 如何公式不想编号则使用 \begin{equation*} 和\end{equation*}.
2、字符
除了# $ % & ~ _ ^ \ { }普通字符在数学公式中含义一样，若要在数学环境中表示这些符号# $ % & _ { }，需要分别表示为\# \$ \% \& \_ \{ \}，即在个字符前加上\。

3、上标和下标
用 ^ 来表示上标，用 _ 来表示下标，看一简单例子：
LaTeX可以通过这符号 $^$ 和 $_$ 来设置上标和下标。使用可以参见：技巧十。
用 ^ 来表示上标，用 _ 来表示下标，如果上标的内容多于一个字符，注意用 { } 把上标括起来，上下标是可以嵌套的，下面是一些简单例子：
$\sum_{i=1}^n a_i=0$
$f(x)=x^{x^x}$
4、希腊字母
5、数学函数
例如sin x， 输入应该为\sin x
6、在公式中插入文本可以通过 \mbox{text} 在公式中添加text，比如：
   \documentclass{article}
	\usepackage{CJK}
	\begin{CJK*}{GBK}{song} 
	\begin{document} 
	$$\mbox{对任意的$x>0$}, \mbox{有 }f(x)>0. $$ 
	\end{CJK*}
	\end{document}
7、分数及开方
\frac{numerator}{denominator} \sqrt{expression_r_r_r}表示开平方，
\sqrt[n]{expression_r_r_r} 表示开 n 次方.
8、省略号（3个点）
\ldots 表示跟文本底线对齐的省略号；\cdots 表示跟文本中线对齐的省略号，
9、括号和分隔符
() 和 [ ] 和 ｜ 对应于自己；
{} 对应于 \{ \}；
|| 对应于 \|。
当要显示大号的括号或分隔符时，要对应用 \left 和 \right
10、多行的数学公式
其中&是对其点，表示在此对齐。
*使latex不自动显示序号，如果想让latex自动标上序号，则把*去掉
11、矩阵
12、导数、极限、求和、积分(Derivatives, Limits, Sums and Integrals)
$\frac{du}{dt}   $
$  \frac{d^2 u}{dx^2}$
$\lim_{x \to +\infty}, \inf_{x > s}$
$\frac{1}{\lim_{u \rightarrow \infty}}, \frac{1}{\lim\limits_{u \rightarrow \infty}} or
\frac{1}{ \displaystyle \lim_{u \rightarrow \infty}}$
</div>
</div>
