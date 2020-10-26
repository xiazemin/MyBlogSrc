---
title: type_hinting
layout: post
category: php
author: 夏泽民
---
<!-- more -->
从PHP5开始，我们可以使用类型提示来指定定义函数时，函数接收的参数类型。如果在定义函数时，指定了参数的类型，那么当我们调用函数时，如果实参的类型与指定的类型不符，那么PHP会产生一个致命级别的错误(Catchable fatal error)。

类名称和数组

在定义函数时，PHP只支持两种类型声明：类名称和数组。类名称表名该参数接收的实参为对应类实例化的对象，数组表明接收的实参为数组类型。下面是一个例子：
复制代码 代码如下:

function demo(array $options){
  var_dump($options);
}

在定义demo()函数的时候，指定了函数接收的参数类型为数组。如果我们调用函数时，传入的参数不是数组类型，例如像下面这样的调用：
复制代码 代码如下:

$options='options';
demo($options);

那么将产生以下错误：
复制代码 代码如下:

Catchable fatal error: Argument 1 passed to demo() must be of the type array, string given,
可以使用null作为默认参数

注意

有一点需要特别注意的是，PHP只支持两种类型的类型声明,其他任何标量类型的声明都是不支持的，比如下下面的代码都将产生错误:
复制代码 代码如下:

function demo(string $str){
}
$str="hello";
demo($str)
当我们运行上面的代码时，string会被当做类名称对待，因此会报下面的错误:
Catchable fatal error: Argument 1 passed to demo() must be an instance of string, string given,

总结

类型声明也是PHP面向对象的一个进步吧，尤其是在捕获某种指定类型的异常时非常有用。
使用类型声明，也可以增加代码的可读性。
但是，由于PHP是弱类型的语言，使用类型声明又于PHP设计的初衷相悖。
