---
title: file_get_contents
layout: post
category: php
author: 夏泽民
---
1: 用file_get_contents 以get方式获取内容

<?php 
$url='http://www.baidu.com/'; 
$html=file_get_contents($url); 
//print_r($http_response_header); 
ec($html); 
printhr(); 
printarr($http_response_header); 
printhr(); 
?> 

示例代码2: 用fopen打开url, 以get方式获取内容

<? 
$fp=fopen($url,'r'); 
printarr(stream_get_meta_data($fp)); 
printhr(); 
while(!feof($fp)){ 
$result.=fgets($fp,1024); 
} 
echo"url body: $result"; 
printhr(); 
fclose($fp); 
?> 

示例代码3：用file_get_contents函数,以post方式获取url

<?php 
$data=array('foo'=>'bar'); 
$data=http_build_query($data); 

$opts=array( 
'http'=>array( 
'method'=>'POST', 
'header'=>"Content-type: application/x-www-form-urlencoded\r\n". 
"Content-Length: ".strlen($data)."\r\n", 
'content'=>$data 
), 
); 
$context=stream_context_create($opts); 
$html=file_get_contents('http://localhost/e/admin/test.html',false,$context); 
echo$html; 
?> 

示例代码4：用fsockopen函数打开url，以get方式获取完整的数据，包括header和body

<? 
functionget_url($url,$cookie=false){ 
$url=parse_url($url); 
$query=$url[path]."?".$url[query]; 
ec("Query:".$query); 
$fp=fsockopen($url[host],$url[port]?$url[port]:80,$errno,$errstr,30); 
if(!$fp){ 
returnfalse; 
}else{ 
$request="GET$queryHTTP/1.1\r\n"; 
$request.="Host:$url[host]\r\n"; 
$request.="Connection: Close\r\n"; 
if($cookie)$request.="Cookie: $cookie\n"; 
$request.="\r\n"; 
fwrite($fp,$request); 
while(!@feof($fp)){ 
$result.=@fgets($fp,1024); 
} 
fclose($fp); 
return$result; 
} 
} 
//获取url的html部分，去掉header 
functionGetUrlHTML($url,$cookie=false){ 
$rowdata=get_url($url,$cookie); 
if($rowdata) 
{ 
$body=stristr($rowdata,"\r\n\r\n"); 
$body=substr($body,4,strlen($body)); 
return$body; 
} 
returnfalse; 
} 
?> 

示例代码5：用fsockopen函数打开url，以POST方式获取完整的数据，包括header和body

<? 
functionHTTP_Post($URL,$data,$cookie,$referrer=""){ 
// parsing the given URL 
$URL_Info=parse_url($URL); 

// Building referrer 
if($referrer=="")// if not given use this script. as referrer 
$referrer="111"; 

// making string from $data 
foreach($dataas$key=>$value) 
$values[]="$key=".urlencode($value); 
$data_string=implode("&",$values); 

// Find out which port is needed - if not given use standard (=80) 
if(!isset($URL_Info["port"])) 
$URL_Info["port"]=80; 

// building POST-request: 
$request.="POST ".$URL_Info["path"]." HTTP/1.1\n"; 
$request.="Host: ".$URL_Info["host"]."\n"; 
$request.="Referer:$referer\n"; 
$request.="Content-type: application/x-www-form-urlencoded\n"; 
$request.="Content-length: ".strlen($data_string)."\n"; 
$request.="Connection: close\n"; 
$request.="Cookie: $cookie\n"; 
$request.="\n"; 
$request.=$data_string."\n"; 

$fp=fsockopen($URL_Info["host"],$URL_Info["port"]); 
fputs($fp,$request); 
while(!feof($fp)){ 
$result.=fgets($fp,1024); 
} 
fclose($fp); 
return$result; 
} 
printhr(); 
?> 

示例代码6:使用curl库，使用curl库之前，你可能需要查看一下php.ini，查看是否已经打开了curl扩展

<? 
$ch = curl_init(); 
$timeout = 5; 
curl_setopt ($ch, CURLOPT_URL, 'http://www.baidu.com/'); 
curl_setopt ($ch, CURLOPT_RETURNTRANSFER, 1); 
curl_setopt ($ch, CURLOPT_CONNECTTIMEOUT, $timeout); 
$file_contents = curl_exec($ch); 
curl_close($ch); 
echo $file_contents; 
?> 

关于curl库：
curl官方网站http://curl.haxx.se/
curl 是使用URL语法的传送文件工具，支持FTP、FTPS、HTTP HTPPS SCP SFTP TFTP TELNET DICT FILE和LDAP。curl 支持SSL证书、HTTP POST、HTTP PUT 、FTP 上传，kerberos、基于HTT格式的上传、代理、cookie、用户＋口令证明、文件传送恢复、http代理通道和大量其他有用的技巧

<? 
functionprintarr(array$arr) 
{ 
echo"<br> Row field count: ".count($arr)."<br>"; 
foreach($arras$key=>$value) 
{ 
echo"$key=$value <br>"; 
} 
} 
?> 

PHP抓取远程网站数据的代码
现在可能还有很多程序爱好者都会遇到同样的疑问,就是要如何像搜索引擎那样去抓取别人网站的HTML代码,然后把代码收集整理成为自己有用的数据!今天就等我介绍一些简单例子吧.

Ⅰ.抓取远程网页标题的例子:
以下是代码片段：

<?php 
/* 
+------------------------------------------------------------- 
+抓取网页标题的代码,直接拷贝本代码片段,另存为.php文件执行即可. 
+------------------------------------------------------------- 
*/ 

error_reporting(7); 
$file = fopen ("http://www.jb51.net/", "r"); 
if (!$file) { 
echo "<font color=red>Unable to open remote file.</font>\n"; 
exit; 
} 
while (!feof ($file)) { 
$line = fgets ($file, 1024); 
if (eregi ("<title>(.*)</title>", $line, $out)) { 
$title = $out[1]; 
echo "".$title.""; 
break; 
} 
} 
fclose($file); 

//End 
?> 

Ⅱ.抓取远程网页HTML代码的例子:

<? php 
/* 
+---------------- 
+DNSing Sprider 
+---------------- 
*/ 

$fp = fsockopen("www.dnsing.com", 80, $errno, $errstr, 30); 
if (!$fp) { 
echo "$errstr ($errno)<br/>\n"; 
} else { 
$out = "GET / HTTP/1.1\r\n"; 
$out .= "Host:www.dnsing.com\r\n"; 
$out .= "Connection: Close \r\n\r\n"; 
fputs($fp, $out); 
while (!feof($fp)) { 
echo fgets($fp, 128); 
} 
fclose($fp); 
} 
//End 
?> 

以上两个代码片段都直接Copy回去运行就知道效果了,上面的例子只是抓取网页数据的雏形,要使其更适合自己的使用,情况有各异.所以,在此各位程序爱好者自己好好研究一下吧.

===============================

稍微有点意义的函数是：get_content_by_socket(), get_url(), get_content_url(), get_content_object 几个函数，也许能够给你点什么想法。

<?php 

//获取所有内容url保存到文件 
function get_index($save_file, $prefix="index_"){ 
$count = 68; 
$i = 1; 
if (file_exists($save_file)) @unlink($save_file); 
$fp = fopen($save_file, "a+") or die("Open ". $save_file ." failed"); 
while($i<$count){ 
$url = $prefix . $i .".htm"; 
echo "Get ". $url ."..."; 
$url_str = get_content_url(get_url($url)); 
echo " OK\n"; 
fwrite($fp, $url_str); 
++$i; 
} 
fclose($fp); 
} 

//获取目标多媒体对象 
function get_object($url_file, $save_file, $split="|--:**:--|"){ 
if (!file_exists($url_file)) die($url_file ." not exist"); 
$file_arr = file($url_file); 
if (!is_array($file_arr) || empty($file_arr)) die($url_file ." not content"); 
$url_arr = array_unique($file_arr); 
if (file_exists($save_file)) @unlink($save_file); 
$fp = fopen($save_file, "a+") or die("Open save file ". $save_file ." failed"); 
foreach($url_arr as $url){ 
if (empty($url)) continue; 
echo "Get ". $url ."..."; 
$html_str = get_url($url); 
echo $html_str; 
echo $url; 
exit; 
$obj_str = get_content_object($html_str); 
echo " OK\n"; 
fwrite($fp, $obj_str); 
} 
fclose($fp); 
} 

//遍历目录获取文件内容 
function get_dir($save_file, $dir){ 
$dp = opendir($dir); 
if (file_exists($save_file)) @unlink($save_file); 
$fp = fopen($save_file, "a+") or die("Open save file ". $save_file ." failed"); 
while(($file = readdir($dp)) != false){ 
if ($file!="." && $file!=".."){ 
echo "Read file ". $file ."..."; 
$file_content = file_get_contents($dir . $file); 
$obj_str = get_content_object($file_content); 
echo " OK\n"; 
fwrite($fp, $obj_str); 
} 
} 
fclose($fp); 
} 


//获取指定url内容 
function get_url($url){ 
$reg = '/^http:\/\/[^\/].+$/'; 
if (!preg_match($reg, $url)) die($url ." invalid"); 
$fp = fopen($url, "r") or die("Open url: ". $url ." failed."); 
while($fc = fread($fp, 8192)){ 
$content .= $fc; 
} 
fclose($fp); 
if (empty($content)){ 
die("Get url: ". $url ." content failed."); 
} 
return $content; 
} 

//使用socket获取指定网页 
function get_content_by_socket($url, $host){ 
$fp = fsockopen($host, 80) or die("Open ". $url ." failed"); 
$header = "GET /".$url ." HTTP/1.1\r\n"; 
$header .= "Accept: */*\r\n"; 
$header .= "Accept-Language: zh-cn\r\n"; 
$header .= "Accept-Encoding: gzip, deflate\r\n"; 
$header .= "User-Agent: Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1; SV1; Maxthon; InfoPath.1; .NET CLR 2.0.50727)\r\n"; 
$header .= "Host: ". $host ."\r\n"; 
$header .= "Connection: Keep-Alive\r\n"; 
//$header .= "Cookie: cnzz02=2; rtime=1; ltime=1148456424859; cnzz_eid=56601755-\r\n\r\n"; 
$header .= "Connection: Close\r\n\r\n"; 

fwrite($fp, $header); 
while (!feof($fp)) { 
$contents .= fgets($fp, 8192); 
} 
fclose($fp); 
return $contents; 
} 


//获取指定内容里的url 
function get_content_url($host_url, $file_contents){ 

//$reg = '/^(#|javascript.*?|ftp:\/\/.+|http:\/\/.+|.*?href.*?|play.*?|index.*?|.*?asp)+$/i'; 
//$reg = '/^(down.*?\.html|\d+_\d+\.htm.*?)$/i'; 
$rex = "/([hH][rR][eE][Ff])\s*=\s*['\"]*([^>'\"\s]+)[\"'>]*\s*/i"; 
$reg = '/^(down.*?\.html)$/i'; 
preg_match_all ($rex, $file_contents, $r); 
$result = ""; //array(); 
foreach($r as $c){ 
if (is_array($c)){ 
foreach($c as $d){ 
if (preg_match($reg, $d)){ $result .= $host_url . $d."\n"; } 
} 
} 
} 
return $result; 
} 

//获取指定内容中的多媒体文件 
function get_content_object($str, $split="|--:**:--|"){ 
$regx = "/href\s*=\s*['\"]*([^>'\"\s]+)[\"'>]*\s*(<b>.*?<\/b>)/i"; 
preg_match_all($regx, $str, $result); 

if (count($result) == 3){ 
$result[2] = str_replace("<b>多媒体： ", "", $result[2]); 
$result[2] = str_replace("</b>", "", $result[2]); 
$result = $result[1][0] . $split .$result[2][0] . "\n"; 
} 
return $result; 
} 

?> 
同一域名对应多个IP时，PHP获取远程网页内容的函数

fgc就是简单的读取过来，把一切操作封装了
fopen也进行了一些封装，但是需要你循环读取得到所有数据。
fsockopen这是直板板的socket操作。
如果仅仅是读取一个html页面，fgc更好。
如果公司是通过防火墙上网，一 般的file_get_content函数就不行了。当然，通过一些socket操作，直接向proxy写http请求也是可以的，但是比较麻烦。
如果你能确认文件很小，可以任选以上两种方式fopen ,join(’’,file($file));。比如，你只操作小于1k的文件，那最好还是用file_get_contents吧。

如果确定文件很大，或者不能确定文件的大小，那就最好使用文件流了。fopen一个1K的文件和fopen一个1G的文件没什么明显的区别。内容长，就可以花更长的时间去读，而不是让脚本死掉。

----------------------------------------------------
http://www.phpcake.cn/archives/tag/fsockopen
PHP获取远程网页内容有多种方式，例如用自带的file_get_contents、fopen等函数。

<?php 

echo file_get_contents("http://img.jb51.net/abc.php"); 
?> 
1
2
3
4
但是，在DNS轮询等负载均衡中，同一域名，可能对应多台服务器，多个IP。假设img.jb51.net被DNS解析到 72.249.146.213、72.249.146.214、72.249.146.215三个IP，用户每次访问img.jb51.net，系统会根据负载均衡的相应算法访问其中的一台服务器。
　　上周做一个视频项目时，就碰到这样一类需求：需要依次访问每台服务器上的一个PHP接口程序（假设为abc.php），查询这台服务器的传输状态。

这时就不能直接用file_get_contents访问http://img.jb51.net/abc.php了，因为它可能一直重复访问某一台服务器。

而采用依次访问http://72.249.146.213/abc.php、http://72.249.146.214/abc.php、http://72.249.146.215/abc.php的方法，在这三台服务器上的Web Server配有多个虚拟主机时，也是不行的。

通过设置本地hosts也不行，因为hosts不能设置多个IP对应同一个域名。

那就只有通过PHP和HTTP协议来实现：访问abc.php时，在header头中加上img.jb51.net域名。于是，我写了下面这个PHP函数：

<?php 

/************************ 
* 函数用途：同一域名对应多个IP时，获取指定服务器的远程网页内容 
* 创建时间：2008-12-09 
* 创建人：张宴（img.jb51.net） 
* 参数说明： 
* $ip 服务器的IP地址 
* $host 服务器的host名称 
* $url 服务器的URL地址（不含域名） 
* 返回值： 
* 获取到的远程网页内容 
* false 访问远程网页失败 
************************/ 
function HttpVisit($ip, $host, $url) 
{ 
$errstr = ''; 
$errno = ''; 
$fp = fsockopen ($ip, 80, $errno, $errstr, 90); 
if (!$fp) 
{ 
return false; 
} 
else 
{ 
$out = "GET {$url} HTTP/1.1\r\n"; 
$out .= "Host:{$host}\r\n"; 
$out .= "Connection: close\r\n\r\n"; 
fputs ($fp, $out); 

while($line = fread($fp, 4096)){ 
$response .= $line; 
} 
fclose( $fp ); 

//去掉Header头信息 
$pos = strpos($response, "\r\n\r\n"); 
$response = substr($response, $pos + 4); 

return $response; 
} 
} 

//调用方法： 
$server_info1 = HttpVisit("72.249.146.213", "img.jb51.net", "/abc.php"); 
$server_info2 = HttpVisit("72.249.146.214", "img.jb51.net", "/abc.php"); 
$server_info3 = HttpVisit("72.249.146.215", "img.jb51.net", "/abc.php"); 
?> 
<!-- more -->
https://www.php.net/manual/en/function.file-get-contents.php

早在2010年时候遇到过这样的事情，因为file_get_contents函数造成服务器挂掉的情况，现在觉得很有必要总结下。

公司里有经常有这样的业务，需要调用第三方公司提供的HTTP接口，在把接口提供的信息显示到网页上，代码是这样写的: file_get_contents("http://example.com/") 。
有一天突然接到运维同事的报告，说是服务器挂了，查出原因说是因为file_get_contents函数造成的，那么为什么一个函数会把服务器给搞挂掉呢？
经过详细的查询发现第三方公司提供接口已经坏掉了，就是因为接口坏掉了，才导致服务器挂掉。
问题分析如下：
    我们代码是“file_get_contents("http://example.com/") “获取一个 URL 的返回内容，如果第三方公司提供的URL响应速度慢，或者出现问题，我们服务器的PHP程序将会一直执行去获得这个URL，我 们知道，在 php.ini 中，有一个参数 max_execution_time 可以设置 PHP 脚本的最大执行时间，但是，在 php-cgi(php-fpm) 中，该参数不会起效。真正能够控制 PHP 脚本最大执行时间的是 php-fpm.conf 配置文件中的以下参数： <value name="request_terminate_timeout">0s</value>  　默认值为 0 秒，也就是说，PHP 脚本会一直执行下去，当请求越来越多的情况下会导致php-cgi 进程都卡在 file_get_contents() 函数时，这台 Nginx+PHP 的 WebServer 已经无法再处理新的 PHP 请求了，Nginx 将给用户返回“502 Bad Gateway”。CPU的利用率达到100% ，时间一长服务器就会挂掉。
问题的解决：
     已经找到问题，那么我们该怎么解决呢？
     当时想到的解决问题的办法就是设置PHP的超时时间，用set_time_limit; 设置超时时间，这样就不会卡住了。代码上线后发现服务器还是会挂掉，好像根本不起作用。后来查了资料才知道，set_time_limit设置的是PHP程序的超时时间，而不是file_get_contents函数读取URL的超时时间。set_time_limit和修改php.ini文件里max_execution_time  效果是一样的。
要设置file_get_contents函数的超时时间，可以用resource $context的timeout参数，代码如下：
复制代码
复制代码
1 $opts = array(
2 　　'http'=>array(
3 　　　　'method'=>"GET",
4　　　　 'timeout'=>10,
5　　 )
6 );
7 $context = stream_context_create($opts);
8 $html =file_get_contents('http://www.example.com', false, $context);
9 echo $html;
复制代码
复制代码
代码中的timeout就是file_get_contents读取url的超时时间。

另外还有一个说法也可以改变读取url的超时时间，就是修改php.ini中的default_socket_timeout的值，或者 ini_set('default_socket_timeout',    10);  但是我没有测试过不知道行不行。
有了解决方法之后，服务器就不会挂掉了。
在解决的过程中我还发现起到关键作用的是stream_context_create方法，里面method 可以是GET，那么能否可以POST呢？还有没有其他的参数？
还有一个为老同事告诉我们还有一个比file_get_contents更好的办法，就是用CURL。
请看下面两篇。

目录(?)[+]

file_get_contents函数
一般的也就是使用file_get_contents($url)，但是关于这个函数还有很多没有注意到的地方。

先看关于手册：

file_get_contents(path,include_path,context,start,max_length)
参数

描述

path

必需。规定要读取的文件。

include_path

可选。如果也想在 include_path 中搜寻文件的话，可以将该参数设为 “1″。

context

可选。规定文件句柄的环境。

context 是一套可以修改流的行为的选项。若使用 null，则忽略。

start

可选。规定在文件中开始读取的位置。该参数是 PHP 5.1 新加的。

max_length

可选。规定读取的字节数。该参数是 PHP 5.1 新加的。

——————-可以选择读取文件位置和长度这个选项不错。但是关于context的选项是做什么用的呢？

强大的context——stream_context_create
context 就是文本流的意思。而在php中创建文本流的函数是：stream_context_create

参看官方手册：http://php.net/manual/en/function.stream-context-create.php

stream_context_create是用来创建打开文件的上下文件选项的，比如用POST访问，使用代理，发送header等。看到没有之前用 curl实现的所谓代理，post，header方法都可以使用file_get_contents+stream_context_create来实 现。

之前在《PHP批量采集下载美女图片》中抱怨file_get_contents采集图片时候经常会遇到慢资源造成cpu负载过高，不能设置超时时间，最后使用curl来实现，其实file_get_contents也可以设置超时时间。

file_get_contents超时设置
1
$opts  = array('http'=>array('timeout'=>10));
2
$context  = stream_context_create($opts);
3
echo  file_get_contents($url,false,$context);
这样就可以实现设置10s的超时时间

更强大的file_get_contents

file_get_contents实现post
参看官方手册的例子

1
$opts  = array('http'  =>
2
  array(
3
    'method'   => 'POST',
4
    'header'   => "Content-Type: text/xmlrn".
5
      "Authorization: Basic ".base64_encode("$https_user:$https_password")."rn",
6
    'content'  => $body,
7
    'timeout'  => 60
8
  )
9
);
10
 
11
$context   = stream_context_create($opts);
12
$url  = 'https://'.$https_server;
13
$result  = file_get_contents($url, false, $context, -1, 40000);
还可以实现get请求，header代理等等功能，理论上curl可以实现的功能file_get_contents都可以实现，但是关于 stream_context_create的解释网络上资源不是很多，也注定在采集程序方面curl的应用更广，另外curl是一种通信模式，不是单纯 的php-curl。

之前写过关于解决gzip乱码的问题《:file_get_contents获取gzip网页乱码》

更多高级使用方法参看官方手册的实例：http://php.net/manual/en/function.stream-context-create.php，http://php.net/manual/en/function.file-get-contents.php

file_put_contents函数
语法：

file_put_contents(file,data,mode,context)
参数

描述

file

必需。规定要写入数据的文件。如果文件不存在，则创建一个新文件。

data

可选。规定要写入文件的数据。可以是字符串、数组或数据流。

mode

可选。规定如何打开/写入文件。可能的值：

FILE_USE_INCLUDE_PATH
FILE_APPEND
LOCK_EX
context

可选。规定文件句柄的环境。

context 是一套可以修改流的行为的选项。若使用 null，则忽略。

注意事项：
file_put_contents等于依次调用 fopen()，fwrite() 以及 fclose() 功能一样，但是效率要更高。
data不仅仅是字符串，也包括数组格式和文本流，当是数组格式的时候（只能是一维数组，不能是多维数组），需要把数组分割implode(”, $array)， 其实还是转换为字符串，如果不分割的话文本存储的内容就是$array[0]$array[1]$array[2]$array[3]这种，不利于读取。 文本流这个就更好理解了，例如存储file_get_contents(‘aa.jpg’)这一张图片的二进制流也是可以存储的。
模式：FILE_APPEND 是追加模式，默认的写入方式是覆盖之前的内容，但是使用FILE_APPEND 模式后就可以不覆盖之前的内容了。LOCK_EX是文本锁，防止并行写入冲突。
context 和上面的file_get_contents一样，可以增加文本流选项，官方的一个例子
1
<?php
2
 /* set the FTP hostname */
3
 $user  = "test";
4
 $pass  = "myFTP";
5
 $host  = "example.com";
6
 $file  = "test.txt";
7
 $hostname  = $user  . ":"  . $pass  . "@"  . $host  . "/"  . $file; 
8
 
9
 /* the file content */
10
 $content  = "this is just a test."; 
11
 
12
 /* create a stream context telling PHP to overwrite the file */
13
 $options  = array('ftp'  => array('overwrite'  => true));
14
 $stream  = stream_context_create($options); 
15
 
16
 /* and finally, put the contents */
17
 file_put_contents($hostname, $content, 0, $stream);
18
?>
PHP中使用CURL实现GET和POST请求
AloneMonkey 2014年7月3日 11
一、什么是CURL？

cURL 是一个利用URL语法规定来传输文件和数据的工具，支持很多协议，如HTTP、FTP、TELNET等。最爽的是，PHP也支持 cURL 库。使用PHP的cURL库可以简单和有效地去抓网页。你只需要运行一个脚本，然后分析一下你所抓取的网页，然后就可以以程序的方式得到你想要的数据了。 无论是你想从从一个链接上取部分数据，或是取一个XML文件并把其导入数据库，那怕就是简单的获取网页内容，cURL 是一个功能强大的PHP库。

 

二、CURL函数库。

curl_close — 关闭一个curl会话
curl_copy_handle — 拷贝一个curl连接资源的所有内容和参数
curl_errno — 返回一个包含当前会话错误信息的数字编号
curl_error — 返回一个包含当前会话错误信息的字符串
curl_exec — 执行一个curl会话
curl_getinfo — 获取一个curl连接资源句柄的信息
curl_init — 初始化一个curl会话
curl_multi_add_handle — 向curl批处理会话中添加单独的curl句柄资源
curl_multi_close — 关闭一个批处理句柄资源
curl_multi_exec — 解析一个curl批处理句柄
curl_multi_getcontent — 返回获取的输出的文本流
curl_multi_info_read — 获取当前解析的curl的相关传输信息
curl_multi_init — 初始化一个curl批处理句柄资源
curl_multi_remove_handle — 移除curl批处理句柄资源中的某个句柄资源
curl_multi_select — Get all the sockets associated with the cURL extension, which can then be “selected”
curl_setopt_array — 以数组的形式为一个curl设置会话参数
curl_setopt — 为一个curl设置会话参数
curl_version — 获取curl相关的版本信息

curl_init()函数的作用初始化一个curl会话，curl_init()函数唯一的一个参数是可选的，表示一个url地址。
curl_exec()函数的作用是执行一个curl会话，唯一的参数是curl_init()函数返回的句柄。
curl_close()函数的作用是关闭一个curl会话，唯一的参数是curl_init()函数返回的句柄。

 

三、PHP建立CURL请求的基本步骤

①：初始化

curl_init()

②：设置属性

curl_setopt().有一长串cURL参数可供设置，它们能指定URL请求的各个细节。

③：执行并获取结果

curl_exec()

④：释放句柄

curl_close()

 

四、CURL实现GET和POST

①：GET方式实现

1
2
3
4
5
6
7
8
9
10
11
12
13
14
15
16
<?php
    //初始化
    $curl = curl_init();
    //设置抓取的url
    curl_setopt($curl, CURLOPT_URL, 'http://www.baidu.com');
    //设置头文件的信息作为数据流输出
    curl_setopt($curl, CURLOPT_HEADER, 1);
    //设置获取的信息以文件流的形式返回，而不是直接输出。
    curl_setopt($curl, CURLOPT_RETURNTRANSFER, 1);
    //执行命令
    $data = curl_exec($curl);
    //关闭URL请求
    curl_close($curl);
    //显示获得的数据
    print_r($data);
?>
运行结果：

image

②：POST方式实现

1
2
3
4
5
6
7
8
9
10
11
12
13
14
15
16
17
18
19
20
21
22
23
24
<?php
    //初始化
    $curl = curl_init();
    //设置抓取的url
    curl_setopt($curl, CURLOPT_URL, 'http://www.baidu.com');
    //设置头文件的信息作为数据流输出
    curl_setopt($curl, CURLOPT_HEADER, 1);
    //设置获取的信息以文件流的形式返回，而不是直接输出。
    curl_setopt($curl, CURLOPT_RETURNTRANSFER, 1);
    //设置post方式提交
    curl_setopt($curl, CURLOPT_POST, 1);
    //设置post数据
    $post_data = array(
        "username" => "coder",
        "password" => "12345"
        );
    curl_setopt($curl, CURLOPT_POSTFIELDS, $post_data);
    //执行命令
    $data = curl_exec($curl);
    //关闭URL请求
    curl_close($curl);
    //显示获得的数据
    print_r($data);
?>
③：如果获得的数据时json格式的，使用json_decode函数解释成数组。

$output_array = json_decode($output,true);

如果使用json_decode($output)解析的话，将会得到object类型的数据。

 

五、我自己封装的一个函数

1
2
3
4
5
6
7
8
9
10
11
12
13
14
15
16
17
18
19
20
21
22
23
24
25
26
27
28
29
30
31
32
33
 //参数1：访问的URL，参数2：post数据(不填则为GET)，参数3：提交的$cookies,参数4：是否返回$cookies
 function curl_request($url,$post='',$cookie='', $returnCookie=0){
        $curl = curl_init();
        curl_setopt($curl, CURLOPT_URL, $url);
        curl_setopt($curl, CURLOPT_USERAGENT, 'Mozilla/5.0 (compatible; MSIE 10.0; Windows NT 6.1; Trident/6.0)');
        curl_setopt($curl, CURLOPT_FOLLOWLOCATION, 1);
        curl_setopt($curl, CURLOPT_AUTOREFERER, 1);
        curl_setopt($curl, CURLOPT_REFERER, "http://XXX");
        if($post) {
            curl_setopt($curl, CURLOPT_POST, 1);
            curl_setopt($curl, CURLOPT_POSTFIELDS, http_build_query($post));
        }
        if($cookie) {
            curl_setopt($curl, CURLOPT_COOKIE, $cookie);
        }
        curl_setopt($curl, CURLOPT_HEADER, $returnCookie);
        curl_setopt($curl, CURLOPT_TIMEOUT, 10);
        curl_setopt($curl, CURLOPT_RETURNTRANSFER, 1);
        $data = curl_exec($curl);
        if (curl_errno($curl)) {
            return curl_error($curl);
        }
        curl_close($curl);
        if($returnCookie){
            list($header, $body) = explode("\r\n\r\n", $data, 2);
            preg_match_all("/Set\-Cookie:([^;]*);/", $header, $matches);
            $info['cookie']  = substr($matches[1][0], 1);
            $info['content'] = $body;             return $info;         }else{             return $data;         } }
 

附可选参数说明：

第一类：
对于下面的这些option的可选参数，value应该被设置一个bool类型的值：
选项
可选value值
备注
CURLOPT_AUTOREFERER
当根据Location:重定向时，自动设置header中的Referer:信息。
CURLOPT_BINARYTRANSFER
在启用CURLOPT_RETURNTRANSFER的时候，返回原生的（Raw）输出。
CURLOPT_COOKIESESSION
启用时curl会仅仅传递一个session cookie，忽略其他的cookie，默认状况下cURL会将所有的cookie返回给服务端。session cookie是指那些用来判断服务器端的session是否有效而存在的cookie。
CURLOPT_CRLF
启用时将Unix的换行符转换成回车换行符。
CURLOPT_DNS_USE_GLOBAL_CACHE
启用时会启用一个全局的DNS缓存，此项为线程安全的，并且默认启用。
CURLOPT_FAILONERROR
显示HTTP状态码，默认行为是忽略编号小于等于400的HTTP信息。
CURLOPT_FILETIME
启用时会尝试修改远程文档中的信息。结果信息会通过 curl_getinfo()函数的CURLINFO_FILETIME选项返回。curl_getinfo().
CURLOPT_FOLLOWLOCATION
启用时会将服务器服务器返回的”Location: “放在header中递归的返回给服务器，使用CURLOPT_MAXREDIRS可以限定递归返回的数量。
CURLOPT_FORBID_REUSE
在完成交互以后强迫断开连接，不能重用。
CURLOPT_FRESH_CONNECT
强制获取一个新的连接，替代缓存中的连接。
CURLOPT_FTP_USE_EPRT
启用时当FTP下载时，使用EPRT (或 LPRT)命令。设置为FALSE时禁用EPRT和LPRT，使用PORT命令 only.
CURLOPT_FTP_USE_EPSV
启用时，在FTP传输过程中回复到PASV模式前首先尝试EPSV命令。设置为FALSE时禁用EPSV命令。
CURLOPT_FTPAPPEND
启用时追加写入文件而不是覆盖它。
CURLOPT_FTPASCII
CURLOPT_TRANSFERTEXT的别名。
CURLOPT_FTPLISTONLY
启用时只列出FTP目录的名字。
CURLOPT_HEADER
启用时会将头文件的信息作为数据流输出。
CURLINFO_HEADER_OUT
启用时追踪句柄的请求字符串。
从 PHP 5.1.3 开始可用。CURLINFO_前缀是故意的(intentional)。
CURLOPT_HTTPGET
启用时会设置HTTP的method为GET，因为GET是默认是，所以只在被修改的情况下使用。
CURLOPT_HTTPPROXYTUNNEL
启用时会通过HTTP代理来传输。
CURLOPT_MUTE
启用时将cURL函数中所有修改过的参数恢复默认值。
CURLOPT_NETRC
在连接建立以后，访问~/.netrc文件获取用户名和密码信息连接远程站点。
CURLOPT_NOBODY
启用时将不对HTML中的BODY部分进行输出。
CURLOPT_NOPROGRESS
启用时关闭curl传输的进度条，此项的默认设置为启用。
Note:
PHP自动地设置这个选项为TRUE，这个选项仅仅应当在以调试为目的时被改变。
CURLOPT_NOSIGNAL
启用时忽略所有的curl传递给php进行的信号。在SAPI多线程传输时此项被默认启用。
cURL 7.10时被加入。
CURLOPT_POST
启用时会发送一个常规的POST请求，类型为：application/x-www-form-urlencoded，就像表单提交的一样。
CURLOPT_PUT
启用时允许HTTP发送文件，必须同时设置CURLOPT_INFILE和CURLOPT_INFILESIZE。
CURLOPT_RETURNTRANSFER
将 curl_exec()获取的信息以文件流的形式返回，而不是直接输出。
CURLOPT_SSL_VERIFYPEER
禁 用后cURL将终止从服务端进行验证。使用CURLOPT_CAINFO选项设置证书使用CURLOPT_CAPATH选项设置证书目录 如果CURLOPT_SSL_VERIFYPEER(默认值为2)被启用，CURLOPT_SSL_VERIFYHOST需要被设置成TRUE否则设置为 FALSE。
自cURL 7.10开始默认为TRUE。从cURL 7.10开始默认绑定安装。
CURLOPT_TRANSFERTEXT
启用后对FTP传输使用ASCII模式。对于LDAP，它检索纯文本信息而非HTML。在Windows系统上，系统不会把STDOUT设置成binary模式。
CURLOPT_UNRESTRICTED_AUTH
在使用CURLOPT_FOLLOWLOCATION产生的header中的多个locations中持续追加用户名和密码信息，即使域名已发生改变。
CURLOPT_UPLOAD
启用后允许文件上传。
CURLOPT_VERBOSE
启用时会汇报所有的信息，存放在STDERR或指定的CURLOPT_STDERR中。

第二类：
对于下面的这些option的可选参数，value应该被设置一个integer类型的值：
选项
可选value值
备注
CURLOPT_BUFFERSIZE
每次获取的数据中读入缓存的大小，但是不保证这个值每次都会被填满。
在cURL 7.10中被加入。
CURLOPT_CLOSEPOLICY
不是CURLCLOSEPOLICY_LEAST_RECENTLY_USED就是CURLCLOSEPOLICY_OLDEST，还存在另外三个CURLCLOSEPOLICY_，但是cURL暂时还不支持。
CURLOPT_CONNECTTIMEOUT
在发起连接前等待的时间，如果设置为0，则无限等待。
CURLOPT_CONNECTTIMEOUT_MS
尝试连接等待的时间，以毫秒为单位。如果设置为0，则无限等待。
在cURL 7.16.2中被加入。从PHP 5.2.3开始可用。
CURLOPT_DNS_CACHE_TIMEOUT
设置在内存中保存DNS信息的时间，默认为120秒。
CURLOPT_FTPSSLAUTH
FTP验证方式：CURLFTPAUTH_SSL (首先尝试SSL)，CURLFTPAUTH_TLS (首先尝试TLS)或CURLFTPAUTH_DEFAULT (让cURL自动决定)。
在cURL 7.12.2中被加入。
CURLOPT_HTTP_VERSION
CURL_HTTP_VERSION_NONE (默认值，让cURL自己判断使用哪个版本)，CURL_HTTP_VERSION_1_0 (强制使用 HTTP/1.0)或CURL_HTTP_VERSION_1_1 (强制使用 HTTP/1.1)。
CURLOPT_HTTPAUTH
使用的HTTP验证方法，可选的值有：CURLAUTH_BASIC、CURLAUTH_DIGEST、CURLAUTH_GSSNEGOTIATE、CURLAUTH_NTLM、CURLAUTH_ANY和CURLAUTH_ANYSAFE。
可以使用|位域(或)操作符分隔多个值，cURL让服务器选择一个支持最好的值。
CURLAUTH_ANY等价于CURLAUTH_BASIC | CURLAUTH_DIGEST | CURLAUTH_GSSNEGOTIATE | CURLAUTH_NTLM.
CURLAUTH_ANYSAFE等价于CURLAUTH_DIGEST | CURLAUTH_GSSNEGOTIATE | CURLAUTH_NTLM.
CURLOPT_INFILESIZE
设定上传文件的大小限制，字节(byte)为单位。
CURLOPT_LOW_SPEED_LIMIT
当传输速度小于CURLOPT_LOW_SPEED_LIMIT时(bytes/sec)，PHP会根据CURLOPT_LOW_SPEED_TIME来判断是否因太慢而取消传输。
CURLOPT_LOW_SPEED_TIME
当传输速度小于CURLOPT_LOW_SPEED_LIMIT时(bytes/sec)，PHP会根据CURLOPT_LOW_SPEED_TIME来判断是否因太慢而取消传输。
CURLOPT_MAXCONNECTS
允许的最大连接数量，超过是会通过CURLOPT_CLOSEPOLICY决定应该停止哪些连接。
CURLOPT_MAXREDIRS
指定最多的HTTP重定向的数量，这个选项是和CURLOPT_FOLLOWLOCATION一起使用的。
CURLOPT_PORT
用来指定连接端口。（可选项）
CURLOPT_PROTOCOLS
CURLPROTO_* 的位域指。如果被启用，位域值会限定libcurl在传输过程中有哪些可使用的协议。这将允许你在编译libcurl时支持众多协议，但是限制只是用它们 中被允许使用的一个子集。默认libcurl将会使用全部它支持的协议。参见CURLOPT_REDIR_PROTOCOLS.
可用的协议选项 为：CURLPROTO_HTTP、CURLPROTO_HTTPS、CURLPROTO_FTP、CURLPROTO_FTPS、 CURLPROTO_SCP、CURLPROTO_SFTP、CURLPROTO_TELNET、CURLPROTO_LDAP、 CURLPROTO_LDAPS、CURLPROTO_DICT、CURLPROTO_FILE、CURLPROTO_TFTP、 CURLPROTO_ALL
在cURL 7.19.4中被加入。
CURLOPT_PROXYAUTH
HTTP代理连接的验证方式。使用在CURLOPT_HTTPAUTH中的位域标志来设置相应选项。对于代理验证只有CURLAUTH_BASIC和CURLAUTH_NTLM当前被支持。
在cURL 7.10.7中被加入。
CURLOPT_PROXYPORT
代理服务器的端口。端口也可以在CURLOPT_PROXY中进行设置。
CURLOPT_PROXYTYPE
不是CURLPROXY_HTTP (默认值) 就是CURLPROXY_SOCKS5。
在cURL 7.10中被加入。
CURLOPT_REDIR_PROTOCOLS
CURLPROTO_* 中的位域值。如果被启用，位域值将会限制传输线程在CURLOPT_FOLLOWLOCATION开启时跟随某个重定向时可使用的协议。这将使你对重定向 时限制传输线程使用被允许的协议子集默认libcurl将会允许除FILE和SCP之外的全部协议。这个和7.19.4预发布版本种无条件地跟随所有支持 的协议有一些不同。关于协议常量，请参照CURLOPT_PROTOCOLS。
在cURL 7.19.4中被加入。
CURLOPT_RESUME_FROM
在恢复传输时传递一个字节偏移量（用来断点续传）。
CURLOPT_SSL_VERIFYHOST
1 检查服务器SSL证书中是否存在一个公用名(common name)。译者注：公用名(Common Name)一般来讲就是填写你将要申请SSL证书的域名 (domain)或子域名(sub domain)。2 检查公用名是否存在，并且是否与提供的主机名匹配。
CURLOPT_SSLVERSION
使用的SSL版本(2 或 3)。默认情况下PHP会自己检测这个值，尽管有些情况下需要手动地进行设置。
CURLOPT_TIMECONDITION
如 果在CURLOPT_TIMEVALUE指定的某个时间以后被编辑过，则使用CURL_TIMECOND_IFMODSINCE返回页面，如果没有被修改 过，并且CURLOPT_HEADER为true，则返回一个”304 Not Modified”的header，        CURLOPT_HEADER为false，则使用CURL_TIMECOND_IFUNMODSINCE，默认值为 CURL_TIMECOND_IFUNMODSINCE。
CURLOPT_TIMEOUT
设置cURL允许执行的最长秒数。
CURLOPT_TIMEOUT_MS
设置cURL允许执行的最长毫秒数。
在cURL 7.16.2中被加入。从PHP 5.2.3起可使用。
CURLOPT_TIMEVALUE
设置一个CURLOPT_TIMECONDITION使用的时间戳，在默认状态下使用的是CURL_TIMECOND_IFMODSINCE。

第三类：
对于下面的这些option的可选参数，value应该被设置一个string类型的值：
选项
可选value值
备注
CURLOPT_CAINFO
一个保存着1个或多个用来让服务端验证的证书的文件名。这个参数仅仅在和CURLOPT_SSL_VERIFYPEER一起使用时才有意义。 .
CURLOPT_CAPATH
一个保存着多个CA证书的目录。这个选项是和CURLOPT_SSL_VERIFYPEER一起使用的。
CURLOPT_COOKIE
设定HTTP请求中”Cookie: “部分的内容。多个cookie用分号分隔，分号后带一个空格(例如， “fruit=apple; colour=red”)。
CURLOPT_COOKIEFILE
包含cookie数据的文件名，cookie文件的格式可以是Netscape格式，或者只是纯HTTP头部信息存入文件。
CURLOPT_COOKIEJAR
连接结束后保存cookie信息的文件。
CURLOPT_CUSTOMREQUEST
使 用一个自定义的请求信息来代替”GET”或”HEAD”作为HTTP请求。这对于执行”DELETE” 或者其他更隐蔽的HTTP请求。有效值如”GET”，”POST”，”CONNECT”等等。也就是说，不要在这里输入整个HTTP请求。例如输 入”GET /index.html HTTP/1.0\r\n\r\n”是不正确的。
Note:
在确定服务器支持这个自定义请求的方法前不要使用。
CURLOPT_EGDSOCKET
类似CURLOPT_RANDOM_FILE，除了一个Entropy Gathering Daemon套接字。
CURLOPT_ENCODING
HTTP请求头中”Accept-Encoding: “的值。支持的编码有”identity”，”deflate”和”gzip”。如果为空字符串””，请求头会发送所有支持的编码类型。
在cURL 7.10中被加入。
CURLOPT_FTPPORT
这个值将被用来获取供FTP”POST”指令所需要的IP地址。”POST”指令告诉远程服务器连接到我们指定的IP地址。这个字符串可以是纯文本的IP地址、主机名、一个网络接口名（UNIX下）或者只是一个’-’来使用默认的IP地址。
CURLOPT_INTERFACE
网络发送接口名，可以是一个接口名、IP地址或者是一个主机名。
CURLOPT_KRB4LEVEL
KRB4 (Kerberos 4) 安全级别。下面的任何值都是有效的(从低到高的顺序)：”clear”、”safe”、”confidential”、”private”.。如果字符串 和这些都不匹配，将使用”private”。这个选项设置为NULL时将禁用KRB4 安全认证。目前KRB4 安全认证只能用于FTP传输。
CURLOPT_POSTFIELDS
全 部数据使用HTTP协议中的”POST”操作来发送。要发送文件，在文件名前面加上@前缀并使用完整路径。这个参数可以通过urlencoded后的字符 串类似’para1=val1¶2=val2&…’或使用一个以字段名为键值，字段数据为值的数组。如果value是一个数组，Content- Type头将会被设置成multipart/form-data。
CURLOPT_PROXY
HTTP代理通道。
CURLOPT_PROXYUSERPWD
一个用来连接到代理的”[username]:[password]“格式的字符串。
CURLOPT_RANDOM_FILE
一个被用来生成SSL随机数种子的文件名。
CURLOPT_RANGE
以”X-Y”的形式，其中X和Y都是可选项获取数据的范围，以字节计。HTTP传输线程也支持几个这样的重复项中间用逗号分隔如”X-Y,N-M”。
CURLOPT_REFERER
在HTTP请求头中”Referer: “的内容。
CURLOPT_SSL_CIPHER_LIST
一个SSL的加密算法列表。例如RC4-SHA和TLSv1都是可用的加密列表。
CURLOPT_SSLCERT
一个包含PEM格式证书的文件名。
CURLOPT_SSLCERTPASSWD
使用CURLOPT_SSLCERT证书需要的密码。
CURLOPT_SSLCERTTYPE
证书的类型。支持的格式有”PEM” (默认值), “DER”和”ENG”。
在cURL 7.9.3中被加入。
CURLOPT_SSLENGINE
用来在CURLOPT_SSLKEY中指定的SSL私钥的加密引擎变量。
CURLOPT_SSLENGINE_DEFAULT
用来做非对称加密操作的变量。
CURLOPT_SSLKEY
包含SSL私钥的文件名。
CURLOPT_SSLKEYPASSWD
在CURLOPT_SSLKEY中指定了的SSL私钥的密码。
Note:
由于这个选项包含了敏感的密码信息，记得保证这个PHP脚本的安全。
CURLOPT_SSLKEYTYPE
CURLOPT_SSLKEY中规定的私钥的加密类型，支持的密钥类型为”PEM”(默认值)、”DER”和”ENG”。
CURLOPT_URL
需要获取的URL地址，也可以在 curl_init()函数中设置。
CURLOPT_USERAGENT
在HTTP请求中包含一个”User-Agent: “头的字符串。
CURLOPT_USERPWD
传递一个连接中需要的用户名和密码，格式为：”[username]:[password]“。

第四类
对于下面的这些option的可选参数，value应该被设置一个数组：
选项
可选value值
备注

CURLOPT_HTTP200ALIASES
200响应码数组，数组中的响应吗被认为是正确的响应，否则被认为是错误的。
在cURL 7.10.3中被加入。
CURLOPT_HTTPHEADER
一个用来设置HTTP头字段的数组。使用如下的形式的数组进行设置： array(‘Content-type: text/plain’, ‘Content-length: 100′)
CURLOPT_POSTQUOTE
在FTP请求执行完成后，在服务器上执行的一组FTP命令。
CURLOPT_QUOTE
一组先于FTP请求的在服务器上执行的FTP命令。

对于下面的这些option的可选参数，value应该被设置一个流资源 （例如使用 fopen()）：
选项
可选value值
CURLOPT_FILE
设置输出文件的位置，值是一个资源类型，默认为STDOUT (浏览器)。
CURLOPT_INFILE
在上传文件的时候需要读取的文件地址，值是一个资源类型。
CURLOPT_STDERR
设置一个错误输出地址，值是一个资源类型，取代默认的STDERR。
CURLOPT_WRITEHEADER
设置header部分内容的写入的文件地址，值是一个资源类型。
对于下面的这些option的可选参数，value应该被设置为一个回调函数名：
选项
可选value值
CURLOPT_HEADERFUNCTION
设置一个回调函数，这个函数有两个参数，第一个是cURL的资源句柄，第二个是输出的header数据。header数据的输出必须依赖这个函数，返回已写入的数据大小。
CURLOPT_PASSWDFUNCTION
设置一个回调函数，有三个参数，第一个是cURL的资源句柄，第二个是一个密码提示符，第三个参数是密码长度允许的最大值。返回密码的值。
CURLOPT_PROGRESSFUNCTION
设置一个回调函数，有三个参数，第一个是cURL的资源句柄，第二个是一个文件描述符资源，第三个是长度。返回包含的数据。

CURLOPT_READFUNCTION
拥有两个参数的回调函数，第一个是参数是会话句柄，第二是HTTP响应头信息的字符串。使用此函数，将自行处理返回的数据。返回值为数据大小，以字节计。返回0代表EOF信号。
CURLOPT_WRITEFUNCTION
拥有两个参数的回调函数，第一个是参数是会话句柄，第二是HTTP响应头信息的字符串。使用此回调函数，将自行处理响应头信息。响应头信息是整个字符串。设置返回值为精确的已写入字符串长度。发生错误时传输线程终止。

 

PHP中fopen,file_get_contents,curl函数的区别：

1.fopen /file_get_contents 每次请求都会重新做DNS查询，并不对 DNS信息进行缓存。但是CURL会自动对DNS信息进行缓存。对同一域名下的网页或者图片的请求只需要一次DNS查询。这大大减少了DNS查询的次数。所以CURL的性能比fopen /file_get_contents 好很多。

2.fopen /file_get_contents 在请求HTTP时，使用的是http_fopen_wrapper，不会keeplive。而curl却可以。这样在多次请求多个链接时，curl效率会好一些。

3.fopen / file_get_contents 函数会受到php.ini文件中allow_url_open选项配置的影响。如果该配置关闭了，则该函数也就失效了。而curl不受该配置的影响。

4.curl 可以模拟多种请求，例如：POST数据，表单提交等，用户可以按照自己的需求来定制请求。而fopen / file_get_contents只能使用get方式获取数据。
file_get_contents 获取远程文件时会把结果都存在一个字符串中 fiels函数则会储存成数组形式

因此，我还是比较倾向于使用curl来访问远程url。Php有curl模块扩展，功能很是强大。

这是别人做过的关于curl和file_get_contents的测试：

file_get_contents抓取google.com需用秒数：
2.31319094
2.30374217
2.21512604
3.30553889
2.30124092

curl使用的时间：
0.68719101
0.64675593
0.64326
0.81983113
0.63956594
