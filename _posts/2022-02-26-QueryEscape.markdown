---
title: QueryEscape
layout: post
category: golang
author: 夏泽民
---
scheme:[//[user:password@]host[:port]]path[?query][#fragment]
Go provides the following two functions to encode or escape a string so that it can be safely placed inside a URL:

QueryEscape(): Encode a string to be safely placed inside a URL query string.
PathEscape(): Encode a string to be safely placed inside a URL path segment.
<!-- more -->

https://www.callicoder.com/golang-url-encoding-decoding/
将中文转换为unicode码，使用golang中的strconv包中的QuoteToASCII直接进行转换，将unicode码转换为中文就比较麻烦一点，先对unicode编码按\u进行分割，然后使用strconv.ParseInt，将16进制数字转换Int64，在使用fmt.Sprintf将数字转换为字符，最后将其连接在一起，这样就变成了中文字符串了。

https://www.cnblogs.com/borey/p/5622812.html

网址URL中特殊字符转义编码

字符    -    URL编码值

空格    -    %20
"          -    %22
#         -    %23
%        -    %25
&         -    %26
(          -    %28
)          -    %29
+         -    %2B
,          -    %2C
/          -    %2F
:          -    %3A
;          -    %3B
<         -    %3C
=         -    %3D
>         -    %3E
?         -    %3F
@       -    %40
\          -    %5C
|          -    %7C 

https://blog.csdn.net/p312011150/article/details/78928003

您可以为此创建包装器shell脚本，并使用颜色转义序列为其着色。这是Linux上的一个简单示例（我不确定它在Windows上会有什么样子，但我想有办法...... :)）

go test -v . | sed ''/PASS/s//$(printf "\033[32mPASS\033[0m")/'' | sed ''/FAIL/s//$(printf "\033[31mFAIL\033[0m")/''

https://www.thinbug.com/q/27242652
https://stackoverflow.com/questions/55802157/parse-error-near-while-setting-heroku-config-vars
