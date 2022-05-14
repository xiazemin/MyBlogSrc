---
title: php json_encode 的坑
layout: post
category: php
author: 夏泽民
---
成功则返回 JSON 编码的 string 或者在失败时返回 FALSE 。
<?php
$a="在水";
$b=substr($a,0,1);
var_dump($b);
//string(1) "�"
var_dump(json_encode($b));
//bool(false)
var_dump(json_encode(false));
//string(5) "false"
var_dump(json_encode($b,JSON_UNESCAPED_SLASHES | JSON_UNESCAPED_UNICODE));
//bool(false)
<!-- more -->
json_encode的字符串里面包含无法解析的编码，比如URLdecode的转义不完整，比如转义出来的是你好⊙◆？带有乱码，解决办法去除字符串中的乱码或者用正则preg_match_all()把需要的字符串匹配出来，重新进行json_encode

php判断字符串包含乱码的话就不显示，提示字符串中含有乱码，
可以运用黑魔法之json_encode( $string) === 'null'来判断。如果字符串中含有乱码，json_encode该字符串就会返回null。

可以用正则匹配，但是你需要知道乱码大概包括的符号有哪些。
出现PHP substr中文乱码的情况，可能会导致程序无法正常运行。解决办法主要有两种：

一、使用mbstring扩展库的mb_substr()截取就不会出现乱码了。

可以用mb_substr()/mb_strcut()这个函数，mb_substr()/mb_strcut()的用法与substr()相似，只是在mb_substr()/mb_strcut最后要加入多一个参数，以设定字符串的编码，但是一般的服务器都没打开php_mbstring.dll，需要在php.ini在把php_mbstring.dll打开。
<?php
  echo mb_substr("php中文字符encode",0,4,"utf-8");
?>
如果未指定最后一个编码参数，会是三个字节为一个中文，这就是utf-8编码的特点，若加上utf-8字符集说明，所以，是以一个字为单位来截取的。
