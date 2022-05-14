---
title: spl_autoload_register
layout: post
category: php
author: 夏泽民
---
为什么要用spl_autoload_register其实我觉得这一段话基本可以解决所有的问题

尽管 __autoload() 函数也能自动加载类和接口，但更建议使用 spl_autoload_register() 函数。spl_autoload_register() 提供了一种更加灵活的方式来实现类的自动加载（同一个应用中，可以支持任意数量的加载器，比如第三方库中的）。因此，不再建议使用 __autoload() 函数，在以后的版本中它可能被弃用。

<!-- more -->

ClassA.php

<?php
class ClassA{ 
    public function __construct(){ 
        echo "ClassA load success!"; 
    } 
} 
ClassB.php

<?php
class ClassB extends ClassA { 
    public function __construct(){ 
        //parent::__construct(); 
        echo "ClassB load success!"; 
    } 
} 
index.php

<?php
function __autoload($classname){ 
    $classpath="./".$classname.'.php'; 
    if(file_exists($classpath)){ 
        require_once($classpath); 
    }else{ 
        echo "class file".$classpath.'not found!'; 
    } 
} 
 
$newobj = new ClassA(); 
$newobj = new ClassB(); 
?> 
这是个autoload的例子，三个文件在同一级目录下，大家看我慢慢改造他。

一、优点

1、真的很灵活

首先把__autoload注释掉，换成

spl_autoload_register();
依然是可以运行的。说明他自动封装了方法。

然后我们在方法下面加上

function my_autoload ($pClassName) {
    include(__DIR__ . "/classes/" . $pClassName . ".php");
}
spl_autoload_register("my_autoload");
并在当前目录下创建classes文件夹，把ClassA复制一份进去，改成C，你会发现C也是可以调用的。说明他支持多次注册，比原来更灵活了。

再有你可以更改文件的扩展名，我们把ClassA和ClassB两个类文件扩展名分别改为.abc，.ini

使用

spl_autoload_extensions('.abc,.ini')
 扩展名可以随便改，ClassC是因为你指定了，不能修改了。

再配合上spl_autoload_unregister基本无敌了。

2、更好的错误处理

<?php
spl_autoload_register(function ($name) {
    echo "Want to load $name.\n";
    throw new Exception("Unable to load $name.");
});
 
try {
    $obj = new NonLoadableClass();
} catch (Exception $e) {
    echo $e->getMessage(), "\n";
}
在 PHP 5.3 之前，__autoload 函数抛出的异常不能被 catch 语句块捕获并会导致一个致命错误（Fatal Error）。 自 PHP 5.3 起，能够 thrown 自定义的异常（Exception），随后自定义异常类即可使用。 __autoload 函数可以递归的自动加载自定义异常类。

这些例子都来自于文档。

二、参数

spl_autoload_register有三个参数

autoload_function
欲注册的自动装载函数。如果没有提供任何参数，则自动注册 autoload 的默认实现函数spl_autoload()。

throw
此参数设置了 autoload_function 无法成功注册时， spl_autoload_register()是否抛出异常。

prepend
如果是 true，spl_autoload_register() 会添加函数到队列之首，而不是队列尾部。

https://www.php.net/manual/zh/language.oop5.autoload.php


spl_autoload_register 与__autoload
spl_autoload_register 与__autoload都是尝试加载未定义的类。一般我们是实例化一个类比如new Foo();就是需要先include('Foo.php');使用spl_autoload_register 与__autoload就可以避免这种操作，减少include代码。

区别
spl_autoload_register允许存在多个自动加载器，而__autoload只允许一个加载器。
__autoload将在php7.2中放弃。

利用spl_autoload_register实现简单的命名空间自动加载案例
//Load.php
class Load{
    //声明命名空间与目录的映射
    public static $vendorMap = array(
        'app' => ROOT_DIR . DIRECTORY_SEPARATOR . 'app',
    );

    /**
     自动加载器
    */

    public static function autoLoad($class)
    {
        $file = self::findClassFile($class);
        if(file_exists($file))
        {
            include($file);
        }else
        {
            throw new Exception("Error $file Not Found", 1);
        }
    }

    /**
    * 解析文件路径
    */
    public static function findClassFile($class)
    {
        //class = app\controllers\Index
        //strpos 获取字符串在另一个字符串第一次出现的位置
        $vendor = substr($class,0,strpos($class, '\\')); //顶级命名空间
        $vendorDir = self::$vendorMap[$vendor]; // 文件基目录
        $filePath = substr($class, strlen($vendor)) . '.php'; // 文件相对路径
        return strtr($vendorDir . $filePath, '\\', DIRECTORY_SEPARATOR); // 文件标准路径
    }

}

//test.php
define('ROOT_DIR',__DIR__);
require './namespace/Load.php';
spl_autoload_register('Load::autoload'); // 注册自动加载

$controller = new app\controllers\Index();
echo $controller->IndexAction();

https://www.cnblogs.com/linqingvoe/p/10937709.html

