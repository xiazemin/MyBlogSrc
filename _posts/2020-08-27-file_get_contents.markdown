---
title: file_get_contents
layout: post
category: php
author: 夏泽民
---
PHP中file() 函数和file_get_contents() 函数的作用都是将整个文件读入某个介质，其主要区别就在于这个介质的不同。file() 函数是将文件读入一个数组中，而file_get_contents()是将文件读入一个字符串中。

file() 函数是把整个文件读入一个数组中，然后将文件作为一个数组返回。数组中的每个单元都是文件中相应的一行，包括换行符在内。如果失败，则返回 false。

file_get_contents() 函数是把整个文件读入一个字符串中。和 file() 一样，不同的是file_get_contents() 把文件读入一个字符串。file_get_contents() 函数是用于将文件的内容读入到一个字符串中的首选方法。如果操作系统支持，还会使用内存映射技术来增强性能。
<!-- more -->
我们来看一个例子吧，下面是test.txt的内容如下：

springload
news
eightnight
file_get_contents的操作方法如下：

$content = file_get_contents('test.txt');
$temp =str_replace(chr(13),'|',$content);
$arr =explode('|',$temp);
print_r($arr);
file的操作方法如下：

print_r(file('test.txt'));
效果与上面完全相同的，结果如下：

Array
(
    [0] => springload
    [1] => news
    [2] => eightnight
)
file()与file_get_contents()得到文件的内容方式是直接得到。这个与fgets不一样，fgets得到的文件内容方式是由fopen打开的资源流。

在可以用file_get_contents替代file、fopen、feof、fgets等系列方法的情况下，尽量用 file_get_contents，因为他的效率高得多!但是要注意file_get_contents在打开一个url文件时候的php版本问题。


{% raw %}
原有服务器的PHP环境为5.5，云服务的PHP环境为5.6。当时，抓取远程内容的函数用的是：get_file_content()，迁移之后，发现PDF文件打不开，经过调试，原来PHP5.5时，抓取URL远程内容时，不会自动gzip压缩内容，而PHP5.6时，抓取URL远程内容时，会自动gzip压缩，恰恰 get_file_content()，不能自动通过gizp压缩的方式抓取，故导致无发抓取远程内容。

解决方案：使用curl方式远程抓取，代码如下

/**
 *param  String $url url地址
 *param  Array  $param 请求参数 
 *param  Boolen $ispost 请求方式：true-POST请求|false-GET请求
 *param  Boolen $gzip 是否gizp压缩方式请求：true-是|false-否
 */
function get_curl($url, $params = '', $ispost = false,$gzip=false)
	{
		$ch = curl_init();
		curl_setopt($ch, CURLOPT_USERAGENT, 'Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.22 (KHTML, like Gecko) Chrome/25.0.1364.172 Safari/537.22');
		curl_setopt($ch, CURLOPT_CONNECTTIMEOUT, 30);
		curl_setopt($ch, CURLOPT_TIMEOUT, 30);
		curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
		if ($ispost) {
			if (is_array($params)) $params = http_build_query($params);
			curl_setopt($ch, CURLOPT_POST, true);
			curl_setopt($ch, CURLOPT_POSTFIELDS, $params);
			curl_setopt($ch, CURLOPT_URL, $url);
		} else {
			if ($params) {
				if (is_array($params)) $params = http_build_query($params);
				curl_setopt($ch, CURLOPT_URL, $url . '?' . $params);
			} else {
				curl_setopt($ch, CURLOPT_URL, $url);
			}
		}
		
		if($gzip){
			curl_setopt($ch, CURLOPT_ENCODING, "gzip");
		}
		
		$response = curl_exec($ch);
		if ($response === FALSE) {
			return false;
		}
		curl_close($ch);
 
		return $response;
}
{% endraw %}

php中 curl， fsockopen ，file_get_contents 三个函数 都可以实现采集模拟发言 。三者有什么区别，或者讲究么

赵永斌:
有些时候用file_get_contents()调用外部文件,容易超时报错。换成curl后就可以.具体原因不清楚
curl 效率比file_get_contents()和fsockopen()高一些,原因是CURL会自动对DNS信息进行缓存(亮点啊有我待亲测)

范佳鹏:
file_get_contents curl fsockopen
在当前所请求环境下选择性操作，没有一概而论：
具我们公司开发KBI应用来看：
刚开始采用：file_get_contents
后来采用：fsockopen
最后到至今采用：curl

（远程）我个人理解到的表述如下（不对请指出，不到位请补充）
file_get_contents 需要php.ini里开启allow_url_fopen,请求http时，使用的是http_fopen_wrapper，不会keeplive.curl是可以的。
file_get_contents()单个执行效率高，返回没有头的信息。
这个是读取一般文件的时候并没有什么问题，但是在读取远程问题的时候就会出现问题。
如果是要打一个持续连接，多次请求多个页面。那么file_get_contents和fopen就会出问题。
取得的内容也可能会不对。所以做一些类似采集工作的时候，肯定就有问题了。
sock较底层，配置麻烦，不易操作。 返回完整信息。

潘少宁-腾讯：
file_get_contents 虽然可以获得某URL的内容，但不能post get啊。
curl 则可以post和get啊。还可以获得head信息
而socket则更底层。可以设置基于UDP或是TCP协议去交互
file_get_contents 和 curl 能干的，socket都能干。
socket能干的，curl 就不一定能干了
file_get_contents 更多的时候 只是去拉取数据。效率比较高 也比较简单。
赵的情况这个我也遇到过，我通过CURL设置host 就OK了。 这和网络环境有关系

https://www.php.net/manual/zh/function.file-get-contents.php

一句话总结：使用file_get_contents()进行分段读取，file_get_contents()函数可以分段读取
 

1、读取大文件是，file_get_contents()函数为什么会发生错误？
发生内存溢出而打开错误

当我们遇到文本文件体积很大时，比如超过几十M甚至几百M几G的大文件，用记事本或者其它编辑器打开往往不能成功，因为他们都需要把文件内容全部放到内存里面，这时就会发生内存溢出而打开错误

因为file_get_contents()只能读取长度为 maxlen 的内容

 

2、file_get_contents()函数的机制是什么？
file_get_contents() 把文件读入一个字符串。将在参数 offset 所指定的位置开始读取长度为 maxlen 的内容。如果失败，file_get_contents() 将返回 FALSE。

 

3、file_get_contents()函数如何读取大文件？
使用file_get_contents()进行分段读取

$str = $content=file_get_contents("2.sql",FALSE,NULL,1024*1024,1024);
echo $str;

 

$u ='www.bKjia.c0m'; //此文件为100GB

$a =file_get_contents( $u,100,1000 );

读取成功了

 

 

4、fread()函数如何分段读取？
其实就是判断文件是否结束，没结束就一直读

$fp=fopen('2.sql','r');
while (!feof($fp)){
$str.=fread($fp, filesize ($filename)/10);//每次读出文件10分之1

 

5、如何设置file_get_contents函数的超时时间？
<?php
//设置超时参数
$opts=array(
        "http"=>array(
                "method"=>"GET",
                "timeout"=>3
                ),
        );
////创建数据流上下文
$context = stream_context_create($opts);

//$url请求的地址，例如：

$result =file_get_contents($url, false, $context);

// 打印结果
print_r($result);

?>
 

 

 

回到顶部
二、php 使用file_get_contents读取大文件的方法
当我们遇到文本文件体积很大时，比如超过几十M甚至几百M几G的大文件，用记事本或者其它编辑器打开往往不能成功，因为他们都需要把文件内容全部放到内存里面，这时就会发生内存溢出而打开错误，遇到这种情况我们可以使用PHP的文件读取函数file_get_contents()进行分段读取。

函数说明
string file_get_contents ( string $filename [, bool $use_include_path [, resource $context [, int $offset [, int $maxlen ]]]] )
和 file() 一样，只除了 file_get_contents() 把文件读入一个字符串。将在参数 offset 所指定的位置开始读取长度为 maxlen 的内容。如果失败，file_get_contents() 将返回 FALSE。

file_get_contents() 函数是用来将文件的内容读入到一个字符串中的首选方法。如果操作系统支持还会使用内存映射技术来增强性能。

应用：

代码如下:

$str = $content=file_get_contents("2.sql",FALSE,NULL,1024*1024,1024);
echo $str;
如果针对较小文件只是希望分段读取并以此读完可以使用fread()函数

代码如下:

$fp=fopen('2.sql','r');
while (!feof($fp)){
$str.=fread($fp, filesize ($filename)/10);//每次读出文件10分之1
//进行处理
}
echo $str;

以上就是如何使用file_get_contents函数读取大文件的方法，超级简单吧，需要的小伙伴直接搬走！

 

 

参考：php 使用file_get_contents读取大文件的方法_php技巧_脚本之家
https://www.jb51.net/article/57380.htm

 

回到顶部
三、解决php中file_get_contents 读取大文件返回false问题
file_get_contents文件是用来读写文件的，但我发现用file_get_contents 读取大文件出错提示Note: string can be as large as 2GB了，这个就是不能超过2G了，有没有办法解决呢，下面我来一起来看。

如果我读取一个 www.bKjia.c0m文件

 代码如下
复制代码

$u ='www.bKjia.c0m'; //此文件为100GB

$a =file_get_contents( $u );

运行提示

Note: string can be as large as 2GB

不能大于2GB了，我们去官方看此函数参考

string file_get_contents ( string $filename [, bool $use_include_path = false [, resource $context [, int $offset = -1 [, int $maxlen ]]]] )

发现有个

 file_get_contents() 把文件读入一个字符串。将在参数 offset 所指定的位置开始读取长度为 maxlen 的内容。如果失败， file_get_contents() 将返回 FALSE。

原来如此，这样我们对程序进行修改即可

 代码如下
复制代码

$u ='www.bKjia.c0m'; //此文件为100GB

$a =file_get_contents( $u,100,1000 );

读取成功了

总结

file_get_contents如果正常返回，会把文件内容储存到某个字符串中，所以它不应该返回超过2G长度的字符串。

如果文件内容超过2G，不加offset和maxlen调用file_get_contents的话，肯定会返回false，

 

 

参考：解决php中file_get_contents 读取大文件返回false问题 | E网新时代
http://www.jxtobo.com/56854.html

 
 
回到顶部
四、PHP中file_get_contents($url)的超时处理
PHP中file_get_contents函数的作用是获取一个 URL 的返回内容。如果是url响应速度慢，或者网络等因素，会造成等待时间较长的情况。只需设置一下file_get_contents函数的超时时间即可解决。示例代码如下：

<?php
//设置超时参数
$opts=array(
        "http"=>array(
                "method"=>"GET",
                "timeout"=>3
                ),
        );
////创建数据流上下文
$context = stream_context_create($opts);

//$url请求的地址，例如：

$result =file_get_contents($url, false, $context);

// 打印结果
print_r($result);

?>


