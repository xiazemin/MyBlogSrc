---
title: prestissimo Composer 加速
layout: post
category: php
author: 夏泽民
---
https://github.com/hirak/prestissimo
https://zhuanlan.zhihu.com/p/64419387
它的作用就是在安装拓展的时候提升 Composer 的安装速度，其原理是使用多进程下载的方式来解决。我体验了一下，真的是爽到不行。它的安装也是非常简单。

       composer global require hirak/prestissimo
如果安装的时候报错了，是一个404错误。那是因为 Composer 源没有，更换一个 Composer 源即可。该源是 Composer 官网提供的。

composer config -g repo.packagist composer https://packagist.laravel-china.org
然后再执行上面的命令，你就可以成功安装该插件。从此以后， Composer 的安装速度就起飞啦
<!-- more -->
要求
composer >=1.0.0 (includes dev-master)
PHP >=5.3, (suggest >=5.5, because curl_share_init)
ext-curl
安装
$ composer global require hirak/prestissimo
卸载
$ composer global remove hirak/prestissimo
基准测试效果
288s -> 26s

$ composer create-project laravel/laravel laravel1 --no-progress --profile --prefer-dist

为什么 php composer速度这么慢？
因为composer实现了file_get_contents()。没有TCP优化，没有Keep-Alive，没有多路复用...

我创建了一个并行下载软件包来下载composer插件。
 https://packagist.org/packages/hirak/prestissimo
 
 https://stackoverflow.com/questions/28436237/why-is-php-composer-so-slow
 
 
 php的file_get_contents为什么很慢？
 
 <?php
function sendGetByCurl($url, $time)
{
    $ch = curl_init();
    curl_setopt($ch, CURLOPT_URL, $url);
    curl_setopt($ch, CURLOPT_RETURNTRANSFER, 1);
    curl_setopt($ch, CURLOPT_TIMEOUT, $time);
    curl_setopt($ch, CURLOPT_CONNECTTIMEOUT, $time);
    $return = curl_exec($ch);
    curl_close($ch);
    return $return;
}

$url = 'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTIrQKRNquic8GwsU951TC7PDCFzIew3RFwTOFoNx8u1fln0FOzHv04YBoEqXPTHHfyU0Xa1qoFULCw/132';

$start1 = microtime(true);
$data1 = file_get_contents($url);
file_put_contents('1.jpg', $data1);
$end1 = microtime(true);
$span1 = $end1 - $start1;
echo $span1 . PHP_EOL;

$start2 = microtime(true);
$data2 = sendGetByCurl($url, 3);
file_put_contents('2.jpg', $data2);
$end2 = microtime(true);
$span2 = $end2 - $start2;
echo $span2 . PHP_EOL;

exit;
 
 这个确实有意思。
这是其他人提到的类似bug问题：

PHP file_get_contents very slow when using full url
这是找到的比较靠谱的说法：

PHP:file_get_contents获取微信头像缓慢问题定位
综上，应该是5.6.14之前确实存在此问题，但是可以通过手动增加'header'=>'Connection: close\r\n'来解决，5.6.14之后就不需要了，但是至于为什么php7还出现这个问题，那应该是微信的锅吧，没有正确的响应请求。

这是file_get_contents的一个bug，在最新php版本中已经修复。
https://stackoverflow.com/questions/3629504/php-file-get-contents-very-slow-when-using-full-url

https://github.com/php/php-src/commit/4b1dff6f438f84f7694df701b68744edbdd86153

By analyzing it with Wireshark, the issue (in my case and probably yours too) was that the remote web server DIDN'T CLOSE THE TCP CONNECTION UNTIL 15 SECONDS (i.e. "keep-alive").

Indeed, file_get_contents doesn't send a "connection" HTTP header, so the remote web server considers by default that's it's a keep-alive connection and doesn't close the TCP stream until 15 seconds (It might not be a standard value - depends on the server conf).

A normal browser would consider the page is fully loaded if the HTTP payload length reaches the length specified in the response Content-Length HTTP header. File_get_contents doesn't do this and that's a shame.

SOLUTION

SO, if you want to know the solution, here it is:

$context = stream_context_create(array('http' => array('header'=>'Connection: close\r\n')));
file_get_contents("http://www.something.com/somepage.html",false,$context);
The thing is just to tell the remote web server to close the connection when the download is complete, as file_get_contents isn't intelligent enough to do it by itself using the response Content-Length HTTP header.


1.问题描述
前几周在做微信需求开发的时候一个功能需要拉取微信用户头像，使用了file_get_contents。但是发现拉取非常缓慢，网上查询资料说使用curl即可解决，试了一下确实如此。

但是为何造成这种差异，网上资料解释也五花八门，什么HTTP头不一样、DNS缓存造成的........之类。

2.抓包
为了更深入了解这种差异的原因，我特意编译了一个带debug符号的php和libcurl方便必要的时候进行源码级别调试。这里我先用Wireshark抓包，看file_get_contents和curl在TCP流程上是否存在差异。

file_get_contents复现测试代码

$url = "http://thirdwx.qlogo.cn/mmopen/g3MonUZtNHkdmzicIlibx6iaFqAc56vxLSUfpb6n5WKSYVY0ChQKkiaJSgQ1dZuTOgvLLrhJbERQQ4eMsv84eavHiaiceqxibJxCfHe/0";
$data = file_get_contents($url);
curl测试代码


$url = "http://thirdwx.qlogo.cn/mmopen/g3MonUZtNHkdmzicIlibx6iaFqAc56vxLSUfpb6n5WKSYVY0ChQKkiaJSgQ1dZuTOgvLLrhJbERQQ4eMsv84eavHiaiceqxibJxCfHe/0";
$ch = curl_init($url);
curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
$data = curl_exec($ch);


TCP的握手、传输、关闭流程这里就不概述，可以参考其它文章了解。从抓包结果可以看出，握手、传输流程差异不大，但是到了关闭就有了很大差异。file_get_contents最后由微信服务端关闭(80 -> 本地端口)，并且在等待对方关闭后回复的FIN ACK包耗费了大量时间，可以从Time那一栏看出。而curl最后由自己主动关闭而非等待微信服务端关闭(本地端口 -> 80)，所以可以看到整个流程耗时非常短，这也是造成了使用file_get_contents和curl分别拉取微信头像造成这么大耗时差异的原因。

从file_get_contents抓包结果里面看到，微信服务端隔了一段时间才关闭并回复FIN ACK，那么是不是由于file_get_content的HTTP头带的有Connect:keep-alive，而造成对方维持了一段时间长连接，然后我从抓包结果里面筛出file_get_content的HTTP请求头内容。

可以看出，这里的设置并不是Connect:keep-alive，所以未要求服务端维持长连接，但是微信服务端在返回完所有数据后过了段时间才调用close并返回FIN ACK包进入关闭状态，整体表现非常像设置了Connect:keep-alive，但因为这是微信的服务端我不能进一步追踪，所以只能猜测可能是微信自家的服务端对HTTP协议支持不完全造成的或者这是有什么其它妙用。

上面的线索断了后，大体知道了问题所在，因为file_get_contents在等待对方服务端调用close回复的FIN ACK上耗费了大量时间，所以造成拉取缓慢。但curl因为是主动调用close所以直接进入了关闭状态。

3.file_get_contents调试
但是，为什么curl就能正常处理这种情况？file_get_contents就会发生这种情况。所以我后面从源代码入手跟踪双方在数据传输及连接关闭的流程上有什么差异。

file_get_contents的实现在php源代码目录ext/standard/file.c的521行，主要流程如下图


file_get_content主要流程
php_stream_open_wrapper_ex找到对应的协议实现，进行设置并打开。我们这里是http协议，内部协议相关实现会找到域名对应IP地址、创建socket连接、创建HTTP请求头并通过socket发送等常规操作。

而我们比较关心的传输部分则是在php_stream_copy_to_mem(PHP源码目录/main/streams/stream.c 1393行)这个调用里，而我们比较关心的传输部分的核心逻辑如下图


从服务端读取数据流程
可以看到整体逻辑是先分配一个php的string类型当缓冲区，不断调用php_stream_read直到没有数据为止。当缓冲区大小不够时，会扩容缓冲区，最后返回到php应用层就是php经常用的字符串了。

而php_stream_read封装的最终调用核心逻辑如下图(PHP源码目录/main/streams/xp_socket.c 153行)。


可以看到在调用socket的recv之前，先调用了php_sock_stream_wait_for_data(PHP源代码/main/streams/xp_socket.c 121行)等待数据，而调试跟踪过程中也发现是会在这里阻塞一段时间，然后这个函数的最终调用是poll。最后看调用poll时监控了哪些事件。



PHP_POLLREADABLE的定义如下

#define PHP_POLLREADABLE    (POLLIN|POLLERR|POLLHUP)
其中POLLHUP是在关闭时触发的事件，到此file_get_contents的数据读取流程已经理顺了。

整体流程可以概括为不断调用recv获取数据，直到recv返回0(连接已经有序的关闭了)或者小于0为止，因为recv是非阻塞调用(传入了参数MSG_DONTWAIT)，所以在调用recv之前会调用poll并阻塞到有监控的事件发生的时候在返回。

因为实现是依靠不断调用recv，并靠它的返回值来判断是否读完了，所以在实际过程中，当不断调用poll + recv获取到所有http响应的数据后，因为TCP连接没有立即关闭，而且这个时候对方没有在发送数据，所以再次调用poll时会阻塞等待监听的事件发生，而微信服务端会隔一段时间在关闭并回复FIN ACK。所以poll在阻塞一段时间后，收到了这个回复的FIN ACK，再次调用recv，返回0，最后关闭这个连接，整个流程结束。所以file_get_contents的整个耗时都是被阻塞在等待这个对端关闭的FIN ACK回复上。

4.CURL调试
那么curl为什么没有这个问题？所以我马上又开始调试curl的这个流程，curl的主要流程处理是在CURL源码目录下/lib/multi.c 1288行的multi_runsingle方法，这是个长达800多行的if + swtich组合的判断逻辑（第一次调到这里简直懵逼了好吗！！），这个方法通过循环不断改变和处理连接的状态直到完成，涉及的状态如下图定义(CURL目录/lib/multihandle.h 36行)。


CURL连接状态

对应的英文注释应该能很好解释含义，这里我们只关注一个状态CURLM_STATE_PERFORM，这个是之前请求的状态已经处理完了，可以开始读数据了。

处理这个状态的逻辑在CURL源码目录下/lib/multi.c 1857行，需要注意的是，整个这个循环逻辑都要通过修改一个done变量来指示是否已经全部完成了，所以我们只要观察这个done变量什么时候会修改为true即可找到CURL对读完的处理是怎样判断的。


可以看到这个逻辑把done变量的内存地址传给了Curl_readwrite调用，那么可以肯定这个调用内部会修改这个变量的状态，然后跟踪到这个方法内部(CURL源码目录/lib/transfer.c 1238行)看这个方法在什么条件下会把这个done变量赋值为true。



这里可以看到当连接没有KEEP_RECV等标志时就判断为完成，KEEP_RECV标志代表是否还可以读取，那么我们找到这个标志什么时候被取消的，就知道CURL是如何判断读完了。

完成的读取流程也是由这个方法内部的1125行调用完成读取的。



这个readwrite_data方法实现是一个循环(CURL源代码目录/lib/transfer.c 482行)不断调用Curl_read（最终调用recv）获取数据然后解析。



在调试过程中，发现在726行的逻辑处理中判断了是否读完。



可以看到这个判断的条件是，如果k(struct SingleRequest)里的maxdownload不是-1，并且当前已读数量 + 前面调用Curl_read读到的数据大小如果大于=maxdownload，则在最后取消掉KEPP_RECV标识。

那么maxdownload又是在哪设置的？在随后的调试中发现在该方法内部的539行调用了解析HTTP头的方法。



随后调试到Curl_http_readwrite_headers方法实现(CURL目录/lib/http.c 3010行)的3580行。



我们可以看到，这个maxdownload(读取内容上限)是来自于HTTP头的Content-Length字段。

所以CURL之所以没有发生file_get_contents那样的情况，就是因为它读完Content-Length大小后就关闭连接了。

5.总结
通过抓包调试，我们知道了首先是微信服务器在返回完HTTP响应后并不会马上关闭，而且HTTP头的设置并不是Connect:keep-alive而是Connect:close，所以并不要求服务端维护一段时间长连接，因为file_get_contents的实现是通过不断循环调用socket的recv方法的返回值来判断是否读完所以导致了file_get_contents在微信服务端连接未关闭的时候会一直阻塞等待最后一个关闭回复的FIN ACK包，这就导致了file_get_contents获取微信头像会耗时长的原因。

而CURL则是优先按照HTTP响应头的Content-Length大小来读，并不像filet_get_contents是不断循环调用socket的recv，然后靠recv返回值来判断是否读完，CURL则是读完Content-Length个字节后马上主动关闭连接，所以就不存在等待对端连接关闭了

