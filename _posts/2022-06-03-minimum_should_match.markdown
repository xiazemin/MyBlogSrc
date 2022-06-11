---
title: minimum_should_match
layout: post
category: elasticsearch
author: 夏泽民
---
对于minimum_should_match设置值:

1.minimum_should_match:"3"

无论可选子句的数量如何，都表示固定值.

2.minimum_should_match:"-2"

表示可选子句的总数减去此数字应该是必需的。

3.minimum_should_match:"75%"

表示最少匹配的子句个数,例如有五个可选子句,最少的匹配个数为5*75%=3.75.向下取整为3,这就表示五个子句最少要匹配其中三个才能查到.

4.minimum_should_match:"-25%"

和上面的类似,只是计算方式不一样,假如也是五个子句,5*25%=1.25,向下取整为1,5最少匹配个数为5-1=4.

5.minimum_should_match:"3<90%"

表示如果可选子句的数量等于（或小于）设置的值，则它们都是必需的，但如果它大于设置的值，则适用规范。在这个例子中：如果有1到3个子句，则它们都是必需的，但是对于4个或更多子句，只需要90％的匹配度.
<!-- more -->
https://blog.csdn.net/qq_22985751/article/details/90704189