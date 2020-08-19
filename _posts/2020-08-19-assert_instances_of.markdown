---
title: assert_instances_of instances_of
layout: post
category: php 
author: 夏泽民
---
$this->assertInstanceOf(User::class, $user);
$this->assertInstanceOf(get_class($expectedObject), $user);

https://stackoverflow.com/questions/16833923/phpunit-assertinstanceof-not-working
<!-- more -->
https://www.php.net/manual/zh/language.operators.type.php

作用：（1）判断一个对象是否是某个类的实例，（2）判断一个对象是否实现了某个接口。
第一种用法：

<?php
$obj = new A();
if ($obj instanceof A) {
   echo 'A';
}
?>


第二种用法：


<?php
interface ExampleInterface
{
     public function interfaceMethod();
 }

 class ExampleClass implements ExampleInterface
{
     public function interfaceMethod()
     {
         return 'Hello World!';
     }
 }

$exampleInstance = new ExampleClass();

 if($exampleInstance instanceof ExampleInterface){
     echo 'Yes, it is';
 }else{
     echo 'No, it is not';
} 
?>

instanceof 运算符 和 is_a() 方法都是判断：某对象是否属于该类 或 该类是此对象的父类 或 是否实现了某个接口
是的话返回 TRUE，不是的话返回 FALSE

区别：
instanceof 运算符是 PHP 5 引进的。在此之前用 is_a()，但是后来 is_a() 被废弃而用 instanceof 替代了。
注意：
PHP 5.3.0 起，又恢复使用 is_a() 了。

总结：
现在PHP的服务环境普遍都使用PHP5.0+了，所以尽量使用 instanceof 来代替 is_a()
综上，如果你不知道你的服务器环境，那么建议你使用instanceof，以免造成不必要的麻烦

bool is_a ( object $object , string $class_name )


https://www.jianshu.com/p/458bf12926ed

