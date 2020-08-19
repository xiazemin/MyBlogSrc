---
title: 用phar-composer来构建基于composer的phar包
layout: post
category: php
author: 夏泽民
---
phar.readonly=off
本文成功执行的前提依然是php.ini中的phar.readonly=off。

如果您自己的项目，是基于composer的。并且需要打包成phar文件。那么本条内容，就是您所需要的。我们需要在当前项目根目录下面，执行如下语句。
phar-composer build
但是这条语句顺利生成一个phar文件的前提是：该目录使用composer进行管理。也就是说：根目录下面应该有个composer.json。这个文件，我们可以通过下面命令生成。
composer init
这样的话，我们就符合这个phar-compser工具的执行条件了。事实上，代码还是我们的代码，仅仅多了个没有什么用途的composer.json文件而已。这个就免得我们大家再写个构建脚本了。还是挺容易的，推荐使用。如果您还是需要依赖包的话，您还可能需要下列命令。
composer install

通过git地址打包phar
还是以上述需求为例，我们打包一下aliyun的oss库。

Bash
phar-composer build https://github.com/aliyun/aliyun-oss-php-sdk.git

通过命令行交互式打包
Bash
phar-composer
输入要打包的项目关键词，选择要打包的序号。
<!-- more -->
Git Clone代码到本地:复制
git clone http://www.github.com/clue/phar-composer

https://www.kutu66.com/GitHub/article_102438

https://github.com/clue/phar-composer

https://github.com/clue/phar-composer#phar-composer


{% raw %}
humbug/box 是一款快速的、零配置的 PHAR 打包工具。

还记得前些天的《SMProxy, 让你的数据库操作快三倍！》吗，该项目的 PHAR 便是使用 Box 打包完成的。

该项目是 box-project/box2 的 Fork 分支，原项目已经不再维护。新项目的作者呼吁我们支持该 Fork。

Box 的可配置项有很多，为了能够快速帮助大家了解用法，接下来我将使用 SMProxy 的 box.json 作为例子给大家做一个简单的介绍。

推荐一篇预备知识，可以帮你简单了解 PHAR 的部分用途：使用 phar 上线你的代码包。

首先，正如 Box 作者的描述：

Fast, zero config application bundler with PHARs.

我们默认无需任何配置，在你的 PHP 应用的根目录执行：

composer require humbug/box
vendor/bin/box compile 
即可生成一个基本的 PHAR 包文件。

Box 的配置文件为应用根目录的 box.json，例如 SMProxy 项目的该文件内容为：

{
    "main": "bin/SMProxy",
    "output": "SMProxy.phar",
    "directories": [
        "bin",
        "src"
    ],
    "finder": [
        {
            "notName": "/LICENSE|.*\\.md|.*\\.dist|composer\\.json|composer\\.lock/",
            "exclude": [
                "doc",
                "docs",
                "test",
                "test_old",
                "tests",
                "Tests",
                "vendor-bin"
            ],
            "in": "vendor"
        },
        {
            "name": "composer.json",
            "in": "."
        }
    ],
    "compression": "NONE",
    "compactors": [
        "KevinGH\\Box\\Compactor\\Json",
        "KevinGH\\Box\\Compactor\\Php"
    ],
    "git": "phar-version"
}
main 用于设定应用的入口文件，也就是打包 PHAR 后，直接运行该 PHAR 包所执行的代码，你可以在某种意义上理解为 index.php。
output 用于设定 PHAR 的输出文件，可以包含目录，相对路径或绝对路径。
directories 用于指定打包的 PHP 源码目录。
finder 配置相对比较复杂，底层是使用 Symfony/Finder 实现，与 PHP-CS-Fixer 的 Finder 规则类似。在以上例子中，包含两个 Finder；第一个定义在 vendor 文件夹内，排除指定名称的文件和目录；第二个表示包含应用根目录的 composer.json。
compression 用于设定 PHAR 文件打包时使用的压缩算法。可选值有：GZ（最常用） / BZ2 / NONE（默认）。但有一点需要注意：使用 GZ 要求运行 PHAR 的 PHP 环境已启用 Gzip 扩展，否则会造成报错。
compactors 用于设定压缩器，但此处的压缩器不同于上文所介绍的 compression；一个压缩器类实例可压缩特定文件类型，降低文件大小，例如以下 Box 自带的压缩器：
KevinGH\Box\Compactor\Json：压缩 JSON 文件，去除空格和缩进等。
KevinGH\Box\Compactor\Php：压缩 PHP 文件，去除注释和 PHPDoc 等。
KevinGH\Box\Compactor\PhpScoper：使用 humbug/php-scoper 隔离代码。
git 用于设定一个「占位符」，打包时将会扫描文件内是否含有此处定义的占位符，若存在将会替换为使用 Git 最新 Tag 和 Commit 生成的版本号（例如 2.0.0 或 2.0.0@e558e33）。你可以参考 这里 的代码来更加深入地理解该用法。
{% endraw %}


注意事(坑)项(点)
假如你打包的项目中，入口文件index.php 要引入(include or require)项目中的其他脚本，务必使用绝对路径，否则你打包成phar包之后，其他项目要引入这个phar就会路径出错!!,即如下:
<?php     //这是index.php 入口文件
  require __DIR__."/src/controller.php";  //要使用绝对路径
  require "./lib/tools.php";               //不要使用相对路径

https://dawnki.github.io/2017/07/04/Phar/
https://newsn.net/say/php-phar-create.html

https://packagist.org/

https://m.yisu.com/zixun/39302.html

一个php应用程序往往是由多个文件构成的，如果能把他们集中为一个文件来分发和运行是很方便的，这样的列子有很多，比如在window操作系统上面的安装程序、一个jquery库等等，为了做到这点php采用了phar文档文件格式，这个概念源自java的jar，但是在设计时主要针对 PHP 的 Web 环境，与 JAR 归档不同的是Phar 归档可由 PHP 本身处理，因此不需要使用额外的工具来创建或使用，使用php脚本就能创建或提取它。phar是一个合成词，由PHP 和 Archive构成，可以看出它是php归档文件的意思。

 

关于phar的官网文档请见http://php.net/manual/en/book.phar.php，本文档可以看做和官网文档互为补充

 

phar归档文件有三种格式：tar归档、zip归档、phar归档，前两种执行需要php安装Phar 扩展支持，用的也比较少，这里主要讲phar归档格式。

 

phar格式归档文件可以直接执行，它的产生依赖于Phar扩展，由自己编写的php脚本产生。

 

Phar 扩展对 PHP 来说并不是一个新鲜的概念，在php5.3已经内建于php中，它最初使用 PHP 编写并被命名为 PHP_Archive，然后在 2005 年被添加到 PEAR 库。由于在实际中，解决这一问题的纯 PHP 解决方案非常缓慢，因此 2007 年重新编写为纯 C 语言扩展，同时添加了使用 SPL 的 ArrayAccess 对象遍历 Phar 归档的支持。自那时起，人们做了大量工作来改善 Phar 归档的性能。

 

Phar 扩展依赖于php流包装器，关于此可访问笔者的另外一篇帖子：

http://blog.csdn.net/u011474028/article/details/52814049

 

很多php应用都是以phar格式分发并运行的，著名的有依赖管理：composer、单元测试：phpunit，下面我们来看一看如何创建、运行、提取还原。

 

phar文件的创建：

首先在php.ini中修改phar.readonly这个选项，去掉前面的分号，并改值为off，由于安全原因该选项默认是on，如果在php.ini中是禁用的（值为0或off），那么在用户脚本中可以开启或关闭，如果在php.ini中是开启的，那么用户脚本是无法关闭的，所以这里设置为off来展示示例。

我们来建立一个项目，在服务器根目录中建立项目文件夹为project，目录内的结构如下：

复制代码
file  
    -yunek.js  
    -yunke.css  
lib  
    -lib_a.php  
template  
    -msg.html  
index.php  
Lib.php 
复制代码
其中file文件夹有两个内容为空的js和css文件，仅仅演示phar可以包含多种文件格式

lib_a.php内容如下：

复制代码
<?php  
/** 
 * Created by yunke. 
 * User: yunke 
 * Date: 2017/2/10 
 * Time: 9:23 
 */  
function show(){  
    echo "l am show()";  
}
复制代码
msg.html内容如下：

复制代码
<!DOCTYPE html>  
<html lang="en">  
<head>  
    <meta charset="UTF-8">  
    <title>phar</title>  
</head>  
<body>  
<?=$str; ?>  
</body>  
</html>
复制代码
index.php内容如下：

复制代码
<?php  
/** 
 * Created by yunke. 
 * User: yunke 
 * Date: 2017/2/10 
 * Time: 9:17 
 */  
require "lib/lib_a.php";  
show();  
  
$str = isset($_GET["str"]) ? $_GET["str"] : "hello world";  
include "template/msg.html";  
复制代码
Lib.php内容如下：

复制代码
<?php  
/** 
 * Created by yunke. 
 * User: yunke 
 * Date: 2017/2/10 
 * Time: 9:20 
 */  
function yunke()  
{  
    echo "l am yunke()";  
}
复制代码
项目文件准备好了，开始创建，现在在project文件夹同级目录建立一个yunkeBuild.php，用于产生phar格式文件，内容如下：

复制代码
<?php  
/** 
 * Created by yunke. 
 * User: yunke 
 * Date: 2017/2/10 
 * Time: 9:36 
 */  
  
//产生一个yunke.phar文件  
$phar = new Phar('yunke.phar', 0, 'yunke.phar');  
// 添加project里面的所有文件到yunke.phar归档文件  
$phar->buildFromDirectory(dirname(__FILE__) . '/project');  
//设置执行时的入口文件，第一个用于命令行，第二个用于浏览器访问，这里都设置为index.php  
$phar->setDefaultStub('index.php', 'index.php');  
复制代码
然后在浏览器中访问这个yunkeBuild.php文件，将产生一个yunke.phar文件，此时服务器根目录结构如下：

project  
yunkeBuild.php  
yunke.phar 
这就是产生一个phar归档文件最简单的过程了，更多内容请看官网，这里需要注意的是如果项目不具备单一执行入口则不宜使用phar归档文件

phar归档文件的使用：

我们在服务器根目录建立一个index.php文件来演示如何使用上面创建的phar文件，内容如下：

复制代码
<?php  
  
/** 
 * Created by yunke. 
 * User: yunke 
 * Date: 2017/2/8 
 * Time: 9:33 
 */  
  
require "yunke.phar";  
require "phar://yunke.phar/Lib.php";  
yunke(); 
复制代码
如果index.php文件中只有第一行，那么和不使用归档文件时，添加如下代码完全相同：

require "project/index.php";  
如果没有第二行，那么第三行的yunke()将提示未定义，所以可见require一个phar文件时并不是导入了里面所有的文件，而只是导入了入口执行文件而已，但在实际项目中往往在这个入口文件里导入其他需要使用的文件，在本例中入口执行文件为project/index.php

 

phar文件的提取还原：

我们有时候会好奇phar里面包含的文件源码，这个时候就需要将phar文件还原，如果只是看一看的话可以使用一些ide工具，比如phpstorm 10就能直接打开它，如果需要修改那么就需要提取操作了，为了演示，我们下载一个composer.phar放在服务器目录，在根目录建立一个get.php文件，内容如下：

复制代码
<?php  
/** 
 * Created by yunke. 
 * User: yunke 
 * Date: 2017/2/9 
 * Time: 19:02 
 */  
  
$phar = new Phar('composer.phar');  
$phar->extractTo('composer'); //提取一份原项目文件  
$phar->convertToData(Phar::ZIP); //另外再提取一份，和上行二选一即可  
复制代码
用浏览器访问这个文件，即可提取出来，以上列子展示了两种提取方式：第二行将建立一个composer目录，并将提取出来的内容放入，第三行将产生一个composer.zip文件，解压即可得到提取还原的项目文件。

 

补充：

1、在部署phar文件到生产服务器时需要调整服务器的配置，避免当访问时浏览器直接下载phar文件

2、可以为归档设置别名，别名保存在归档文件中永久保存，它可以用一个简短的名字引用归档，而不管归档文件在文件系统中存储在那里，设置别名：

$phar = new Phar('lib/yunke.phar', 0);  
$phar->setAlias ( "yun.phar");  
<?php  
require "lib/yunke.phar";  
require "phar://yun.phar/Lib.php";  //使用别名访问归档文件  
require "phar://lib/yunke.phar/Lib.php"; //当然仍然可以使用这样的方式去引用  
如果在制作phar文件时没有指定别名，也可以在存根文件里面使用Phar::mapPhar('yunke.phar');指定

3、归档文件中有一个存根文件，其实就是一段php执行代码，在制作归档时可以设置，直接执行归档文件时，其实就是执行它，所以它是启动文件；在脚本中包含归档文件时就像包含普通php文件一样包含它并运行，但直接以phar://的方式包含归档中某一个文件时不会执行存根代码， 往往在存根文件里面require包含要运行的其他文件，对存根文件的限制仅为以__HALT_COMPILER();结束，默认的存根设计是为在没有phar扩展时能够运行，它提取phar文件内容到一个临时目录再执行，不过从php5.3开始该扩展默认内置启用了

4、制作的phar文件不能被改动，因此配置文件之类的文件需要另外放置在归档文件外面

5、mapPhar函数：这个函数只应该在stub存根代码中调用，在没有设置归档别名的时候可以用来设置别名，打开一个引用映射到phar流

https://www.cnblogs.com/fps2tao/p/8717569.html



