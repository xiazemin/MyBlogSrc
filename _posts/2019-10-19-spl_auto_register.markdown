---
title: spl_auto_register
layout: post
category: lang
author: 夏泽民
---
__autoload的作用：当我们在一个页面使用其他文件的类方法时候，经常使用的是require ,require_once ,include,include_once,

但是，如果有100个界面时，我们还都要一个个的require_once进来吗？

__autoload的作用就出来，当我们调用一个从未定义的类时，就会加载__autoload方法，你可以使用__autoload方法加载文件

比如.

auto.php
<?php
function __autoload($className){
    $className = $className.'.php';
    if(is_file($className)){
        require_once "$className";
    }
}
那 spl_auto_register()有什么作用呢?

他可以加载自己写的函数来覆盖__autoload()函数

auto_register.php

复制代码
<?php
function load($className){
    $fileName = $className.'.php';
    if(is_file($fileName)){
        require_once "$fileName";
    }
}
spl_autoload_register('load');
<!-- more -->
解释：

参数$classname为实例化的类名称，此例中就是Test
通过__autoload()找不到的类就会自定进行加载，省去了大量的require重复操作
__autoload()只是在出错失败前多了一次加载所需类的机会，如果__autoload()里还是没找到所需的类，依旧会报错
__autoload()只能定义一次
注意：PHP7.2之后已经弃用了此魔术方法

Warning This feature has been DEPRECATED as of PHP 7.2.0. Relying on this feature is highly discouraged.
spl_autoload_register()
此函数的功能跟__autoload()差不多，都是当实例化当类不存在的时候调用此方法

用法跟__autoload差不多，当类不存在时触发spl_autoload_register()方法，此方法绑定了myLoad方法，执行myLoad方法，相当于__autoload改为我们自己定义的方法了，官方推荐用spl_autoload_register()替代__autoload()

也可以调用静态方法
<?php
class LoadTest
{
    public static function myLoad($classname)
    {
        $classpath='./'.$classname.'.php';
        if (file_exists($classpath)) {
            require_once $classpath;
        }
    }
}
spl_autoload_register(array('LoadTest','myLoad'));

一：什么是自动加载#
我们在new出一个class的时候，不需要手动去require或include来导入这个class文件，而是程序自动帮你导入这个文件
不需要手动的require那么多class文件了

 #
二：怎么样才能自动加载呢#
PHP提供了2种方法，一个是魔术方法 __autoload($classname)，另外一个是函数 spl_autoload_register()

 

三：__autoload 自动加载#
3.1 原理#
当我们new一个classname的时候，如果php找不到这个类，就会去调用 __autoload($classname)，new的这个classname就是这个函数的参数
所以我们就能根据这个classname去require对应路径的类文件，从而实现自动加载

3.2 使用#
student.php

复制代码
<?php
class student {
      function __construct() {
            echo "i am a student";
      }
}
?>
复制代码
 

index.php

复制代码
<?php
$stu = new student();

function __autoload($classname) {
     require $classname.'.php';
}
?>
复制代码
 

 

四：spl_autoload_register 自动加载#
4.1 为什么又出现了个spl_autoload_register 呢#
因为一个项目中只能有一个__autoload，项目小，文件少，一个__autoload 足够用了， 但是随着需求的增加，项目文件变的越变越多，我们需要不同的自动加载来加载不同路径的文件，这时候只有一个 __autoload 就不够用了，如果写2个__autoload，就会报错，所以 spl_autoload_register 函数应运而生，这个函数比 __autoload更好用，更方便

4.2 spl_autoload_register 函数说明#
当我们new一个classname的时候，php找不到classname，php就会去调用spl_autoload_register 注册的函数，这个函数通过参数传递进去

函数原型：

bool spl_autoload_register ([ callable autoload_function[,bool throw = true [, bool $prepend = false ]]] )
autoload_function:
欲注册的自动装载函数。如果没有提供任何参数，则自动注册 autoload 的默认实现函数spl_autoload()。


throw:
此参数设置了 autoload_function 无法成功注册时， spl_autoload_register()是否抛出异常。


prepend:
如果是 true，spl_autoload_register() 会添加函数到队列之首，而不是队列尾部。

4.3 几种参数形式的调用#
复制代码
sql_autoload_resister('load_func'); //函数名
sql_autoload_resister(array('class_object', 'load_func')); //类和静态方法
sql_autoload_resister('class_object::load_func'); //类和方法的静态调用

//php 5.3之后，也可以像这样支持匿名函数了。
spl_autoload_register(function($className){
    if (is_file('./lib/' . $className . '.php')) {
       require './lib/' . $className . '.php';
   }
});
复制代码
 

函数加载 spl_load_func.php

复制代码
<?php
function load_func($classname) {
    require $classname.'.php';
}

spl_autoload_register('load_func');

$stu = new student();
?>
复制代码
 

类加载 spl_load_class.php
类加载的方式必须是static静态方法

复制代码
<?php
class load_class {
    public static function load($classname) {
         require $classname.'.php';
  }
}
// 2种方法调用
spl_autoload_register(array('load_class', 'load'));
spl_autoload_register('load_class::load');

$stu = new student();  // php会自动找到student类并加载
?>
