---
title: php怎么判断函数，类，类方法是不是存在
layout: post
category: php
author: 夏泽民
---
<!-- more -->
{% highlight php linenos %}
（1）php判断系统函数或自己写的函数是否存在
bool function_exists ( string $function_name ) 判断函数是否已经定义，例如：
if(function_exists('curl_init')){
    curl_init();
}else{
    echo 'not function curl_init';
}
（2）php判断类是否存在
bool class_exists ( string $class_name [, bool $autoload = true ] ) 检查一个类是否已经定义，一定以返回true，否则返回false，例如：
if(class_exists('MySQL')){
    $myclass=new MySQL();
}
（3）php判断类里面的某个方法是否已经定义
bool method_exists ( mixed $object , string $method_name ) 检查类的方法是否存在，例如：
$directory=new Directory;
if(!method_exists($directory,'read')){
    echo '未定义read方法！';
}
{% endhighlight %}