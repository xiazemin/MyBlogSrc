---
title: expose_php
layout: post
category: php
author: 夏泽民
---
https://manual.phpdoc.org/
expose_php boolean
决定是否暴露 PHP 被安装在服务器上（例如在 Web 服务器的信息头中加上其签名：X-Powered-By: PHP/5.3.7)。 The PHP logo guids are also exposed, thus appending them to the URL of a PHP enabled site will display the appropriate logo (e.g., » https://www.php.net/?=PHPE9568F34-D428-11d2-A769-00AA001ACF42). This also affects the output of phpinfo(), as when disabled, the PHP logo and credits information will not be displayed.
https://www.php.net/manual/zh/ini.core.php
<!-- more -->
一、如何触发PHP彩蛋？
我们只要在运行PHP的服务器上，在域名后面输入下面的字符参数，就能返回一些意想不到的信息。当然有些服务器是把菜单屏蔽了的。彩蛋只有这4个，PHP是开放源代码的，所以不必担心还有其他。

?=PHPB8B5F2A0-3C92-11d3-A3A9-4C7B08C10000 (PHP信息列表)
?=PHPE9568F34-D428-11d2-A769-00AA001ACF42 (PHP的LOGO)
?=PHPE9568F35-D428-11d2-A769-00AA001ACF42 (Zend LOGO)
?=PHPE9568F36-D428-11d2-A769-00AA001ACF42 (PHP LOGO 蓝色大象)


https://zhang.ge/4983.html
