---
title: register_globals
layout: post
category: php
author: 夏泽民
---
可能 PHP 中最具争议的变化就是从 PHP » 4.2.0 版开始配置文件中 PHP 指令 register_globals 的默认值从 on 改为 off 了。对此选项的依赖是如此普遍以至于很多人根本不知道它的存在而以为 PHP 本来就是这么工作的。本节会解释用这个指令如何写出不安全的代码，但要知道这个指令本身没有不安全的地方，误用才会。

当 register_globals 打开以后，各种变量都被注入代码，例如来自 HTML 表单的请求变量。再加上 PHP 在使用变量之前是无需进行初始化的，这就使得更容易写出不安全的代码。这是个很艰难的抉择，但 PHP 社区还是决定默认关闭此选项。当打开时，人们使用变量时确实不知道变量是哪里来的，只能想当然。

https://www.php.net/manual/zh/security.globals.php
<!-- more -->
https://stackoverflow.com/questions/3593210/what-are-register-globals-in-php

https://sites.google.com/site/phnessu4/register_globals%E4%BD%BF%E7%94%A8%E8%AF%A6%E8%A7%A3

register_globals使用详解
register_globals是php.ini里的一个配置，这个配置影响到php如何接收传递过来的参数，如果你的问题是：为什么我的表单无法传递数据？为什么我的程序无法得到传递过来的变量？等等，那么你需要仔细的阅读以下的内容。

register_globals的值可以设置为：On或者Off，我们举一段代码来分别描述它们的不同。
代码:

<form name="frmTest" id="frmTest" action="URL">
<input type="text" name="user_name" id="user_name">
<input type="password" name="user_pass" id="user_pass">
<input type="submit" value="login">
</form>



当register_globals=Off的时候，下一个程序接收的时候应该用$_GET['user_name']和$_GET['user_pass']来接受传递过来的值。（注：当<form>的method属性为post的时候应该用$_POST['user_name']和$_POST['user_pass']）

当register_globals=On的时候，下一个程序可以直接使用$user_name和$user_pass来接受值。

顾名思义，register_globals的意思就是注册为全局变量，所以当On的时候，传递过来的值会被直接的注册为全局变量直接使用，而Off的时候，我们需要到特定的数组里去得到它。所以，碰到上边那些无法得到值的问题的朋友应该首先检查一下你的register_globals的设置和你获取值的方法是否匹配。（查看可以用phpinfo()函数或者直接查看php.ini）

那我们为什么要使用Off呢？原因有2：
1、php以后的新版本默认都用Off，虽然你可以设置它为On，但是当你无法控制服务器的时候，你的代码的兼容性就成为一个大问题，所以，你最好从现在就开始用Off的风格开始编程
2、这里有两篇文章介绍为什么要Off而不用On
http://www.linuxforum.net/forum/gshowflat.php?Cat=&Board=php3&Number=292803&page=0&view=collapsed&sb=5&o=all&fpart=
http://www.php.net/manual/en/security.registerglobals.php

现在还有一个问题就是，以前用On风格写的大量脚本怎么办？
如果你以前的脚本规划得好，有个公共包含文件，比如config.inc.php一类的文件，在这个文件里加上以下的代码来模拟一下（这个代码不保证100%可以解决你的问题，因为我没有大量测试，但是我觉得效果不错）。另外，这个帖子里的解决方法也可以参考一下（http://www.chinaunix.net/forum/viewtopic.php?t=159284）。
代码:

<?php
if ( !ini_get('register_globals') )
{
extract($_POST);
extract($_GET);
extract($_SERVER);
extract($_FILES);
extract($_ENV);
extract($_COOKIE);

if ( isset($_SESSION) )
{
extract($_SESSION);
}
}
?>

register_globals = Off的情况不仅仅影响到如何获取从<form>、url传递过来的数据，也影响到session、cookie，对应的，得到session、cookie的方式应该为：$_SESSION[]、$_COOKIE。同时对于session的处理也有一些改变，比如，session_register()没有必要而且失效，具体的变化，请查看php manual里的Session handling functions

$_REQUEST中间的内容实际上还是来源于$_GET $_POST $_COOKIE，缺点是无法判断变量到底来自于get post 还是cookie，对要求比较严格的场合不适用。 php manual 写到:

Variables provided to the scrīpt via the GET, POST, and COOKIE input mechanisms, and which therefore cannot be trusted. The presence and order of variable inclusion in this array is defined according to the PHP variables_order configuration directive.


