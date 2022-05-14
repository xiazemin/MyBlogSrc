---
title: prototype
layout: post
category: node
author: 夏泽民
---
在任意对象和Object.prototype之间,
存在着一条以非标准属性__proto__进行连接的链,
我们将这条链称为原型链,
在默认情况下,一个任意的对象的原型是Object.prototype

作用：

表明一个对象的所属类
在“对象.成员”访问失败时，会沿着原型链继续访问

我们每创建一个函数，都有一个prototype(原型)属性，
这个属性是一个指针，指向一个对象，
这个对象包含特定类型所有实例共享的属性和方法

原型对象的内容，包含着两个重要的成员：

constructor(构造函数)
__proto __(原型链)

js原型和原型链 js中this、call、apply、bind简介
https://blog.csdn.net/qingbaicai/article/details/106791954
<!-- more -->
https://www.cnblogs.com/codderYouzg/p/12873535.html

JavaScript严格模式(use strict)
严格模式通过在脚本或函数的头部添加 "use strict"; 表达式来声明。

严格模式下不能使用没有定义的变量，如果在严格模式下是用了未定义的变量，控制台就会报错。


使用严格模式的优点：

消除代码运行的一些不安全之处，保证代码运行的安全；
消除Javascript语法的一些不合理、不严谨之处，减少一些怪异行为;
提高编译器效率，增加运行速度；
为未来新版本的Javascript做好铺垫。
"严格模式"体现了Javascript更合理、更安全、更严谨的发展方向，包括IE 10在内的主流浏览器，都已经支持它，许多大项目已经开始全面拥抱它。

另一方面，同样的代码，在"严格模式"中，可能会有不一样的运行结果；一些在"正常模式"下可以运行的语句，在"严格模式"下将不能运行。掌握这些内容，有助于更细致深入地理解Javascript，让你变成一个更好的程序员工程师。

https://www.shuzhiduo.com/A/amd02vq6Jg/
