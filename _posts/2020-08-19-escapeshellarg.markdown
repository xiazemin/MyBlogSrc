---
title: escapeshellarg
layout: post
category: php
author: 夏泽民
---
(PHP 4 >= 4.0.3, PHP 5, PHP 7)

escapeshellarg — 把字符串转码为可以在 shell 命令里使用的参数

说明 ¶
escapeshellarg ( string $arg ) : string
escapeshellarg() 将给字符串增加一个单引号并且能引用或者转码任何已经存在的单引号，这样以确保能够直接将一个字符串传入 shell 函数，并且还是确保安全的。对于用户输入的部分参数就应该使用这个函数。shell 函数包含 exec(), system() 执行运算符 。

https://www.php.net/manual/zh/function.escapeshellarg.php
<!-- more -->
我们详细分析一下：

传入的参数是：172.17.0.2' -v -d a=1
经过escapeshellarg处理后变成了'172.17.0.2'\'' -v -d a=1'，即先对单引号转义，再用单引号将左右两部分括起来从而起到连接的作用。
经过escapeshellcmd处理后变成'172.17.0.2'\\'' -v -d a=1\'，这是因为escapeshellcmd对\以及最后那个不配对儿的引号进行了转义：http://php.net/manual/zh/function.escapeshellcmd.php
最后执行的命令是curl '172.17.0.2'\\'' -v -d a=1\'，由于中间的\\被解释为\而不再是转义字符，所以后面的'没有被转义，与再后面的'配对儿成了一个空白连接符。所以可以简化为curl 172.17.0.2\ -v -d a=1'，即向172.17.0.2\发起请求，POST 数据为a=1'。
回到mail中，我们的 payload 最终在执行时变成了'-fa'\\''\( -OQueueDirectory=/tmp -X/var/www/html/test.php \)@a.com\'，分割后就是-fa\(、-OQueueDirectory=/tmp、-X/var/www/html/test.php、)@a.com'，最终的参数就是这样被注入的。

谁的锅？

仔细想想其实这可以算是escapeshellarg和escapeshellcmd的设计问题，因为先转义参数再转义命令是很正常的想法，但是它们在配合时并没有考虑到单引号带来的隐患。

在 PHPMailer 的这次补丁中，作者使用escapeshellarg意在防止参数注入，但是却意外的为新漏洞打了助攻，想想也是很有趣的 xD。

攻击面

如果应用使用escapeshellarg -> escapeshellcmd这样的流程来处理输入是存在隐患的，mail就是个很好的例子，因为它函数内部使用了escapeshellcmd，如果开发人员仅用escapeshellarg来处理输入再传给mail那这层防御几乎是可以忽略的。

如果可以注入参数，那利用就是各种各样的了，例如 PHPMailer 和 RoundCube 中的mail和 Naigos Core 中的 curl都是很好的参数注入的例子。

有一点需要注意的是，由于注入的命令中会带有中间的\和最后的'，有可能会影响到命令的执行结果，还要结合具体情况再做分析。

https://paper.seebug.org/164/

0x01 知识铺垫
我们先来看一下 escapeshellarg 函数的用法吧：

escapeshellarg — 把字符串转码为可以在 shell 命令里使用的参数

功能 ：escapeshellarg() 将给字符串增加一个单引号并且能引用或者转码任何已经存在的单引号，这样以确保能够直接将一个字符串传入 shell 函数，shell 函数包含 exec(), system() 执行运算符(反引号)

定义 ：string escapeshellarg ( string $arg )

具体功能作用如下：

1

经过 escapeshellarg 函数处理过的参数被拼凑成 shell 命令，并且被双引号包裹这样就会造成漏洞，这主要在于bash中双引号和单引号解析变量是有区别的。

在解析单引号的时候 , 被单引号包裹的内容中如果有变量 , 这个变量名是不会被解析成值的，但是双引号不同 , bash 会将变量名解析成变量的值再使用。

2

如上可知， 即使参数用了 escapeshellarg 函数过滤单引号，但参数在拼接命令的时候用了双引号的话还是会导致命令执行的漏洞。

紧接着看看 escapeshellcmd 函数的作用吧。

escapeshellcmd — shell 元字符转义

功能：escapeshellcmd() 对字符串中可能会欺骗 shell 命令执行任意命令的字符进行转义。 此函数保证用户输入的数据在传送到 exec() 或 system() 函数，或者 执行操作符 之前进行转义。

反斜线（\）会在以下字符之前插入： &#;`|\?~<>^()[]{}$*, \x0A 和 \xFF*。 *’ 和 “ 仅在不配对儿的时候被转义。 在 Windows 平台上，所有这些字符以及 % 和 ! 字符都会被空格代替。

定义 ：string escapeshellcmd ( string $command)

具体作用举个例子吧：

24

那么 escapeshellcmd 和 escapeshellarg 一起使用的时候会发生什么问题呢，我们继续看看，这两个函数都会对单引号进行处理，但还是有区别的，区别如下:

25

对于单个单引号, escapeshellarg 函数转义后,还会在左右各加一个单引号,但 escapeshellcmd 函数是直接加一个转义符，对于成对的单引号, escapeshellcmd 函数默认不转义,但 escapeshellarg 函数转义:

escapeshellcmd() 和 escapeshellarg() 一起出现会有什么问题呢，我们举个简单例子如下：

23

详细分析一下这个过程：

传入的参数是

1
127.0.0.1' -v -d a=1
由于escapeshellarg先对单引号转义，再用单引号将左右两部分括起来从而起到连接的作用。所以处理之后的效果如下：

1
'127.0.0.1'\'' -v -d a=1'
经过escapeshellcmd针对第二步处理之后的参数中的\以及a=1'中的单引号进行处理转义之后的效果如下所示：

1
'127.0.0.1'\\'' -v -d a=1\'
由于第三步处理之后的payload中的\\被解释成了\而不再是转义字符，所以单引号配对连接之后将payload分割为三个部分，具体如下所示：

18

所以这个payload可以简化为curl 127.0.0.1\ -v -d a=1'，即向127.0.0.1\发起请求，POST 数据为a=1'。

但是如果是先用 escapeshellcmd 函数过滤,再用的 escapeshellarg 函数过滤,则没有这个问题。

0x02 漏洞分析
这里实例分析选择 PHPMailer 命令执行漏洞 （ CVE-2016-10045 和 CVE-2016-10033 ）。项目代码可以通过以下方式下载：

1
2
3
git clone https://github.com/PHPMailer/PHPMailer
cd PHPMailer
git checkout -b CVE-2016-10033 v5.2.17
一. 前情提要
php 的内置函数 mail 可能会造成的命令执行漏洞，在进行漏洞分析之前补一点基础知识，我们先看看 php 自带的 mail 函数如何使用吧。

1
2
3
4
5
6
7
bool mail (
    string $to ,
    string $subject ,
    string $message [,
    string $additional_headers [,
    string $additional_parameters ]]
)
其参数含义分别表示如下：

to，指定邮件接收者，即接收人
subject，邮件的标题
message，邮件的正文内容
additional_headers，指定邮件发送时其他的额外头部，如发送者From，抄送CC，隐藏抄送BCC
additional_parameters，指定传递给发送程序sendmail的额外参数。
在Linux系统上， php 的 mail 函数在底层中已经写好了，默认调用 Linux 的 sendmail 程序发送邮件。而在额外参数中， sendmail 支持主要选项有以下三种：

-O option = value
QueueDirectory = queuedir 选择队列消息
-X logfile
这个参数可以指定一个目录来记录发送邮件时的详细日志情况，我们正式利用这个参数来达到我们的目的。
-f from email
这个参数可以让我们指定我们发送邮件的邮箱地址。
举个简单例子方便理解:

3

上面这个样例代码会在 /var/www/html/rce.php 中写入如下数据：

1
2
3
4
5
6
7
17220 <<< To: Alice@example.com
 17220 <<< Subject: Hello Alice!
 17220 <<< X-PHP-Originating-Script: 0:test.php
 17220 <<< CC: somebodyelse@example.com
 17220 <<<
 17220 <<< <?php phpinfo(); ?>
 17220 <<< [EOF]
有的时候我们会使用以下代码来保证输入参数的可靠性。

1
filter_var($email, FILTER_VALIDATE_EMAIL)
这串代码的主要作用是确保在第5个参数中只使用有效的电子邮件地址 mail() 。我们先简单了解一下 filter_var()这个函数吧。

filter_var ：使用特定的过滤器过滤一个变量

1
mixed filter_var ( mixed $variable [, int $filter = FILTER_DEFAULT [, mixed $options ]] )
功能 ：这里主要是根据第二个参数filter过滤一些想要过滤的东西。

关于 filter_var() 中 FILTER_VALIDATE_EMAIL 这个选项作用，我们可以看看这个帖子 PHP FILTER_VALIDATE_EMAIL 。这里面有个结论引起了我的注意： none of the special characters in this local part are allowed outside quotation marks ，表示所有的特殊符号必须放在双引号中。 filter_var() 问题在于，我们能够在双引号中嵌套转义空格仍然能够通过检测。同时由于底层正则表达式的原因，我们通过重叠单引号和双引号，欺骗 filter_val() 使其认为我们仍然在双引号中，我们就可以绕过检测。下面举个简单的例子，方便理解：

4

当然由于引入的特殊符号，虽然绕过了 filter_var() 针对邮箱的检测，但是由于PHP的 mail() 函数在底层实现中，调用了 escapeshellcmd() 函数对用户输入的邮箱地址进行检测，具体代码在https://github.com/php/php-src/blob/PHP-5.6.29/ext/standard/mail.c ，其中第167-177行如下：

14

所以导致即使存在特殊符号，也会被 escapeshellcmd() 函数处理转义，这样可能没办法达到命令执行的作用了。

那我们前面说过了PHP的 mail() 函数在底层调用了 escapeshellcmd() 函数对用户输入的邮箱地址进行处理，即使我们使用带有特殊字符的payload绕过了 filter_var() 的检测，但是还是会被 escapeshellcmd() 处理，这时候如果按照先用 escapeshellarg() 函数过滤,再用的 escapeshellcmd() 函数过滤的顺序，则可能会发生参数逃逸的问题。

二. 原理介绍
1. CVE-2016-10033
在github上直接diff一下，对比一下不同版本的 class.phpmailer.php 文件，差异如下：

5

这里在 sendmailSend 函数中加了 validateAddress 函数，来针对发送的数据进行判断，判断邮箱地址合法性。另外针对传入的数据，调用了 escapeshellarg 函数来转义特殊符号，防止注入参数。然而这样做，会引入一个新的问题。因为同时使用 escapeshellarg 函数和 escapeshellcmd() 函数，会导致单引号逃逸。由于程序没有对传命令参数的地方进行转义，所以我们可以结合 mail 函数的第五个参数 -X 写入 webshell 。

下面详细看一下代码，漏洞具体位置在 class.phpmailer.php 中，我们截取部分相关代码，下图第12行 ：

7

这里没有 $params 变量进行严格过滤，只是简单地判断是否为 null ，所以可以直接传入命令。继续往下看，我们发现在上图第12行，当 safe_mode 模式处于关闭状态， mail() 函数才会传入 $params 变量。

进一步跟一下这个 $params 参数，看看是怎么来的，这个参数的位置在 class.phpmailer.php 中，我们截取部分相关代码，具体看下图 第11行 ：

8

很明显 $params 是从 $this->Sender 传进来的，继续跟进一下 $this->Sender ，这个函数位置在 class.phpmailer.php 中，截取部分相关代码，具体在 第11行 ：

9

这里在 setFrom 函数中将 $address 经过某些处理之后赋值给 $this->Sender 。

这里详细看看 $address 变量是如何被处理的，主要处理函数在 class.phpmailer.php 文件中，我们截取了部分相关代码，在 第三行 使用了 validateAddress 来处理 $address 变量。

10

所以跟进一下 validateAddress 函数，这个函数位置在 class.phpmailer.php 文件中，截取部分相关代码，我们看看具体做了哪些操作。

11

分析一下这段代码，大概意思就是针对环境进行了判断，就是说如果没有 prce 并且 php 版本 <5.2.0 ，则 $patternselect = ‘noregex’ 。

接着往下看，在 class.phpmailer.php 文件中，有部分关于 $patternselect 的 swich 操作，我只选择了我们需要的那个，跟踪到下面的 noregex 。

12

这里简单的只是根据 @ 符号来处理字符，所以这里的payload很简单。

1
a( -OQueueDirectory=/tmp -X/var/www/html/x.php )@a.com
然后通过 linux 自身的 sendmail 写入log的方式，把log写入了web目录下，成功写入了一个webshell。

2. CVE-2016-10045
diff一下5.2.20和5.2.18发现针对 escapeshellcmd 和 escapeshellarg 做了改动。

13

这里其实有个很奇妙的漏洞，针对入使用 escapeshellarg 处理，最新版本中使用之前的 payload 攻击是失败的，例如：

1
a( -OQueueDirectory=/tmp -X/var/www/html/x.php )@a.com
在最新版中可以使用这个 payload 可以攻击成功：

1
a'( -OQueueDirectory=/tmp -X/var/www/html/x.php )@a.com
这里抛出一个疑问，为什么单引号绕过 escapeshellarg() 限制呢。这时候我们就需要看看， php 的 mail 函数在底层是怎么处理的了，具体代码在https://github.com/php/php-src/blob/PHP-5.6.29/ext/standard/mail.c ，其中第167-177行如下：

14

我们可以针对这个漏洞的 payload 进行测试如下:

16

我们的 payload 最终在执行时变成了'-fa'\\''\( -OQueueDirectory=/tmp -X/var/www/html/test.php \)@a.com\'，分割后就是-fa\(、-OQueueDirectory=/tmp、-X/var/www/html/test.php、)@a.com'，最终的参数就是这样被注入的。

从上面我们可以看到，这里我们想要的东西其实已经逃逸出来了，所以是 escapeshellarg 和 escapeshellcmd函数一起处理参数时会导致的问题。

0x03 后话
因为phpmailer的这个洞，其实就是参数注入的问题。而参数注入漏洞是指，在执行命令的时候，用户控制了命令中的某个参数，并通过一些危险的参数功能，达成攻击的目的，有的时候开发者理解了命令注入，但是参数注入所造成的问题，可能由于自身认知有限，没办法很好的处理。

应用ph师傅博客中的例子：

最典型是案例是Wordpress PwnScriptum漏洞，PHP mail函数的第五个参数，允许直接注入参数，用户通过注入-X参数，导致写入任意文件，最终getshell。

另一个典型的例子是php-cgi CVE-2012-1823 ，在cgi模式中，用户传入的querystring将作为cgi的参数传给php-cgi命令。而php-cgi命令可以用-d参数指定配置项，我们通过指定auto_prepend_file=php://input，最终导致任意代码执行。

客户端上也出现过类似的漏洞，比如Electron CVE-2018-1000006，我们通过注入参数--gpu-launcher=cmd.exe /c start calc，来让electron内置的chromium执行任意命令。electron的最早给出的缓解措施也是在拼接点前面加上“—”。

http://www.lmxspace.com/2018/07/16/%E8%B0%88%E8%B0%88escapeshellarg%E5%8F%82%E6%95%B0%E7%BB%95%E8%BF%87%E5%92%8C%E6%B3%A8%E5%85%A5%E7%9A%84%E9%97%AE%E9%A2%98/

https://www.anquanke.com/post/id/107336

https://www.leavesongs.com/PENETRATION/escapeshellarg-and-parameter-injection.html