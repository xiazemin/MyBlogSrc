---
title: mysql_real_escape_string
layout: post
category: php
author: 夏泽民
---
mysql_real_escape_string — 转义 SQL 语句中使用的字符串中的特殊字符，并考虑到连接的当前字符集

Warning
本扩展自 PHP 5.5.0 起已废弃，并在自 PHP 7.0.0 开始被移除。应使用 MySQLi 或 PDO_MySQL 扩展来替换之。参见 MySQL：选择 API 指南以及相关 FAQ 来获取更多信息。用以替代本函数的有：

mysqli_real_escape_string()
PDO::quote()
说明 ¶
mysql_real_escape_string ( string $unescaped_string [, resource $link_identifier = NULL ] ) : string
本函数将 unescaped_string 中的特殊字符转义，并计及连接的当前字符集，因此可以安全用于 mysql_query()。

mysql_real_escape_string() 调用mysql库的函数 mysql_real_escape_string, 在以下字符前添加反斜杠: \x00, \n, \r, \, ', " 和 \x1a.

为了安全起见，在像MySQL传送查询前，必须调用这个函数（除了少数例外情况）。

Caution
安全提示: 默认字符集
The character set must be set either at the server level, or with the API function mysql_set_charset() for it to affect mysql_real_escape_string(). See the concepts section on character sets for more information.

参数 ¶
unescaped_string
The string that is to be escaped.

link_identifier
MySQL 连接。如不指定连接标识，则使用由 mysql_connect() 最近打开的连接。如果没有找到该连接，会尝试不带参数调用 mysql_connect() 来创建。如没有找到连接或无法建立连接，则会生成 E_WARNING 级别的错误。

https://www.php.net/manual/zh/function.mysql-real-escape-string.php
<!-- more -->
https://www.jb51.net/w3school/php/func_mysql_real_escape_string.htm
The Attack
那么，让我们以显示攻击开始…

mysql_query('SET NAMES gbk');
$var = mysql_real_escape_string("\xbf\x27 OR 1=1 /*");
mysql_query("SELECT * FROM test WHERE name = '$var' LIMIT 1");
在某些情况下，这将返回超过1行。我们来分析一下这里发生了什么：

选择字符集

mysql_query（'SET NAMES gbk'）;
为了使这种攻击能够正常工作，我们需要服务器在连接上预期的编码，如ASCII码0x27 和 >有一些字符的最后一个字节是ASCII \ie 0x5c。事实证明，默认情况下，MySQL 5.6支持5种这样的编码：big5，cp932，gb2312，gbk代码>和<代码> sjis 。我们在这里选择 gbk `。

现在，在这里注意SET NAMES的使用非常重要。这将设置字符集 在服务器上 。如果我们使用了C API函数mysql_set_charset（）的调用，那么我们会很好（自2006年以来在MySQL上发布）。但更多的是为什么在一分钟…

有效负载
我们要用于这个注入的有效负载从字节序列0xbf27开始。在gbk中，这是一个无效的多字节字符;在latin1中，它是字符串Â''。请注意，在latin1 和 gbk中，0x27是一个字面的'字符。

我们选择了这个有效载荷，因为如果我们在它上面调用addslashes（），我们会插入一个ASCII \ 即 0x5c
在‘字符之前。因此，我们将在 gbk 中使用 0xbf5c27 结尾，它是一个两字符序列： 0xbf5c 后跟 0x27 。或者换句话说，一个 _有效的_ 字符后跟一个非转义的‘。但是我们不使用 addslashes（）`。所以下一步…

<强> mysql_real_escape_string（）
对mysql_real_escape_string（）的C API调用与>
addslashes（）的区别在于它知道连接字符集。所以它可以对服务器期望的字符集进行正确的转义。但是，到目前为止，客户端认为我们仍然在使用latin1 来进行连接，因为我们从来没有告诉过它。我们告诉服务器_我们使用的是 gbk ，但是 _客户端_ 仍然认为它是 latin1 `。

因此，对mysql_real_escape_string（）的调用会插入反斜杠，并且在我们的“转义”内容中有一个免费的`字符。实际上，如果我们要查看 gbk 字符集中的 $ var `，我们会看到：

ç¸-'OR 1 = 1 / * 
这是[正是攻击所需要的。

查询
这部分只是一个形式，但是这里是渲染的查询：

SELECT FROM test WHERE name =’ç¸-‘OR 1 = 1 / ‘LIMIT 1

恭喜你，你用mysql_real_escape_string（）成功地攻击了一个程序…

The Bad
情况变得更糟PDO默认使用MySQL模拟_准备好的语句。这意味着在客户端，基本上通过mysql_real_escape_string（）（在C库中）执行sprintf，这意味着以下操作将导致注入成功：

$pdo->query('SET NAMES gbk');
$stmt = $pdo->prepare('SELECT * FROM test WHERE name = ? LIMIT 1');
$stmt->execute(array("\xbf\x27 OR 1=1 /*"));
现在，值得注意的是，您可以通过禁用模拟的准备语句来防止这种情况：

$pdo->setAttribute(PDO::ATTR_EMULATE_PREPARES, false);
这通常会导致准备好的语句（即数据从查询中以单独的数据包发送）。但是，请注意，PDO将默默[后备**
a>来模拟MySQL不能本地准备的语句：它可以是列出，但请注意选择合适的服务器版本）。

The Ugly
我刚开始说如果我们用mysql_set_charset（'gbk'）代替SET NAMES gbk，我们可以阻止所有这些。如果您从2006年开始使用MySQL，那就是真的。

如果您使用的是较早版本的MySQL版本，则可以使用错误 mysql_real_escape_string（）表示无效的多字节字符（如我们的有效载荷中的那些字符）被视为单个字节用于转义，即使客户端已经正确地被通知了连接编码这次攻击仍然会成功。这个错误在MySQL中被修复了[4.1.20 < a>， 5.0.22和 5.1.11

但最糟糕的是PDO在5.3.6之前没有公开mysql_set_charset（）的C API，所以在之前的版本中 **不能 < /
strong>防止每一个可能的命令这个攻击！ 它现在已经公开为 DSN参数。

The Saving Grace
正如我们在一开始所说的，要使这个攻击行得通，数据库连接必须使用易受攻击的字符集进行编码。 utf8mb4不是 很容易
，但是可以支持每一个 Unicode字符：所以你可以选择使用它，但是它只在MySQL 5.5.3以后才可用。另一种方法是 utf8，也是不容易的，并且可以支持整个Unicode
基本多语言平面#Basic_Multilingual_Plane)。

或者，您可以启用 NO_BACKSLASH_ESCAPES SQL模式（其中包括）改变mysql_real_escape_string（）的操作。启用此模式后，0x27将被替换为0x2727而不是0x5c27，因此转义过程 不能 创建有效字符在任何以前不存在的易受攻击的编码中（即0xbf27仍然是0xbf27等） -
所以服务器仍然会拒绝该字符串为无效。不过，请参阅 @
eggyal的回答，了解使用此SQL模式可能导致的另一个漏洞。

Safe Examples
下面的例子是安全的：

mysql_query('SET NAMES utf8');
$var = mysql_real_escape_string("\xbf\x27 OR 1=1 /*");
mysql_query("SELECT * FROM test WHERE name = '$var' LIMIT 1");
因为服务器期望utf8 …

mysql_set_charset('gbk');
$var = mysql_real_escape_string("\xbf\x27 OR 1=1 /*");
mysql_query("SELECT * FROM test WHERE name = '$var' LIMIT 1");
因为我们已经正确设置了字符集，所以客户端和服务器匹配。

$pdo->setAttribute(PDO::ATTR_EMULATE_PREPARES, false);
$pdo->query('SET NAMES gbk');
$stmt = $pdo->prepare('SELECT * FROM test WHERE name = ? LIMIT 1');
$stmt->execute(array("\xbf\x27 OR 1=1 /*"));
因为我们已经关闭了模拟的准备好的语句。

$pdo = new PDO('mysql:host=localhost;dbname=testdb;charset=gbk', $user, $password);
$stmt = $pdo->prepare('SELECT * FROM test WHERE name = ? LIMIT 1');
$stmt->execute(array("\xbf\x27 OR 1=1 /*"));
因为我们已经正确设置了字符集。

$mysqli->query('SET NAMES gbk');
$stmt = $mysqli->prepare('SELECT * FROM test WHERE name = ? LIMIT 1');
$param = "\xbf\x27 OR 1=1 /*";
$stmt->bind_param('s', $param);
$stmt->execute();
因为MySQLi一直都在准备好声明。

Wrapping Up
如果你：

使用MySQL的现代版本（5.1后，所有5.5,5.6等）** mysql_set_charset（） / ` $ mysqli-＆gt; set_charset / code> / PDO的DSN字符集参数（PHP≥5.3.6）
**或

请勿使用易受攻击的字符集进行连接编码（您只能使用utf8 / latin1 / ascii / etc）
你100％安全。

否则，即使您正在使用mysql_real_escape_string（） ，您也是脆弱的 …

http://www.tracholar.top/2018/02/18/sql-injection-that-gets-around-mysql-real-escape-string/

尽量避免使用system，exec，popen,passthru等。如果必须使用这些函数，需要用escapeshellarg过滤参数。

https://cloud.tencent.com/developer/article/1088430

https://www.cnblogs.com/wangtanzhi/p/12238343.html
https://zhuanlan.zhihu.com/p/110351539