---
title: stream_set_timeout
layout: post
category: php
author: 夏泽民
---
PHP函数stream_set_timeout（Stream Functions）作用于读取流时的时间控制。fsockopen函数的timeout只管创建连接时的超时，对于连接后读取流时的超时，则需要用到 stream_set_timeout函数。由于国内的网络环境不是很稳定，尤其是连接国外的时候，不想程序出现Fatal error: Maximum execution time of 30 seconds exceeded in …的错误，该函数尤其有用。stream_set_timeout需配合stream_get_meta_data使用，如果没有timeout， stream_get_meta_data返回数组中time_out为空，反之为1，可根据此判断是否超时。另外由于PHP默认的Maximum execution time为30秒，这是一次执行周期的时间，为了不出现上述的Fatal error，还需要设置一个总的读取流的时间
<!-- more -->
$server="www.yahoo.com";  
$port = 80;  
  
$data="GET / HTTP/1.0rn";  
$data.="Connection: Closern";  
$data.="User-Agent: Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1)rnrn";  
  
$start_time = time();  
$fp=fsockopen($server, $port, $errno, $errstr, 5);  
if (!$fp) {  
die("Connect Timeout.n");  
} else {  
stream_set_blocking($fp, True);  
stream_set_timeout($fp, 3);  
  
fputs($fp, "$data");  
while (!feof($fp)) {  
$text .= fread($fp, 2000);  
  
$diff = time() - $start_time;  
if ($diff > 24) {  
die("Timeout!n");  
}  
  
$status = stream_get_meta_data($fp);  
if ($status[’timed_out’]) {  
die("Stream Timeout!n");  
}  
}  
}  
  
fclose($fp);

接口超时502 bad gate way
超时限制失效
页面调用接口的时候，响应10s。通过使用xhprof监控，发现是调用接口超时。但是实际上接口超时的时候最多只允许2s，通过stream_set_timeout进行了读流超时限制。

查询手册
手册中有一些例子，其中有一段是这样说的：

If you are using fsockopen() to create a connection, first going to write into the stream and then waiting for the reply (e.g. simulating HTTP request with some extra headers), then stream_set_timeout() must be set only after the write - if it is before write, it has no effect on the read timeout :-(
但是实际上，按照这种写法进行了逻辑修改，依然是无效的。

https://www.jb51.net/shouce/php5/zh/function.stream-set-timeout.html

PHP函数stream_set_timeout（Stream Functions）作用于读取流时的时间控制。fsockopen函数的timeout只管创建连接时的超时，对于连接后读取流时的超时，则需要用到stream_set_timeout函数。由于国内的网络环境不是很稳定，尤其是连接国外的时候，不想程序出现Fatal error: Maximum execution time of 30 seconds exceeded in …的错误，该函数尤其有用。stream_set_timeout需配合stream_get_meta_data使用，如果没有timeout， stream_get_meta_data返回数组中time_out为空，反之为1，可根据此判断是否超时。另外由于PHP默认的Maximum execution time为30秒，这是一次执行周期的时间，为了不出现上述的Fatal error，还需要设置一个总的读取流的时间

stream_set_timeout
(PHP 4 >= 4.3.0, PHP 5)

stream_set_timeout — Set timeout period on a stream

说明
bool stream_set_timeout ( resource $stream , int $seconds [, int $microseconds = 0 ] )
Sets the timeout value on stream, expressed in the sum of seconds and microseconds.

When the stream times out, the 'timed_out' key of the array returned by stream_get_meta_data() is set to TRUE, although no error/warning is generated.

参数
stream
The target stream.

seconds
The seconds part of the timeout to be set.

microseconds
The microseconds part of the timeout to be set.

返回值
成功时返回 TRUE， 或者在失败时返回 FALSE.

更新日志
版本	说明
4.3.0	As of PHP 4.3, this function can (potentially) work on any kind of stream. In PHP 4.3, socket based streams are still the only kind supported in the PHP core, although streams from other extensions may support this function.
范例

Example #1 stream_set_timeout() example

<?php
$fp = fsockopen("www.example.com", 80);
if (!$fp) {
    echo "Unable to open ";
} else {

    fwrite($fp, "GET / HTTP/1.0 ");
    stream_set_timeout($fp, 2);
    $res = fread($fp, 2000);

    $info = stream_get_meta_data($fp);
    fclose($fp);

    if ($info['timed_out']) {
        echo 'Connection timed out!';
    } else {
        echo $res;
    }

}
?>
注释
Note:

This function doesn't work with advanced operations like stream_socket_recvfrom(), use stream_select() with timeout parameter instead.

This function was previously called as set_socket_timeout() and later socket_set_timeout() but this usage is deprecated.

参见
fsockopen() - Open Internet or Unix domain socket connection
fopen() - 打开文件或者 URL


stream_set_timeout和stream_socket_sendto猛一看好像没啥问题，从函数名看上去也没啥问题，但是通过php手册可以看到stream_set_timeout是不能控制stream_socket_sendto的

其实我们知道fgets就可以读取流的数据，并且stream_set_timeout是可以控制fread的，所以我们只需要简单的把

stream_socket_sendto替换成fwrite就可以了，stream_socket_recvfrom替换为fread即可