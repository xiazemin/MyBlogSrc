---
title: php流Streams、包装器wrapper 
layout: post
category: php
author: 夏泽民
---
流Streams这个概念是在php4.3引进的，是对流式数据的抽象，用于统一数据操作，比如文件数据、网络数据、压缩数据等，以使可以共享同一套函数，
php的文件系统函数就是这样的共享，比如file_get_contents()函数即可打开本地文件也可以访问url就是这一体现。简单点讲，流就是表现出流式数据行为的资源对象。
以线性方式进行读写，并可以在流里面任意位置进行搜索。
流有点类似数据库抽象层，在数据库抽象层方面，不管使用何种数据库，在抽象层之上都使用相同的方式操作数据，
而流是对数据的抽象，它不管是本地文件还是远程文件还是压缩文件等等，只要来的是流式数据，那么操作方式就是一样的
有了流这个概念就引申出了包装器wrapper这个概念，每个流都对应一种包装器，
流是从统一操作这个角度产生的一个概念，而包装器呢是从理解流数据内容出发产生的一个概念，也就是这个统一的操作方式怎么操作或配置不同的内容；
这些内容都是以流的方式呈现，但内容规则是不一样的，比如http协议传来的数据是流的方式，但只有http包装器才理解http协议传来的数据的意思，
可以这么理解，流就是一根流水的管子，只不过它流出的是数据，包装器就是套在流这根管子外层的一个解释者，它理解流出的数据的意思，并能操作它
官方手册说：“一个包装器是告诉流怎么处理特殊协议或编码的附加代码”明白这句话的意思了吗？

包装器可以嵌套，一个流外面包裹了一个包装器后，还可以在外层继续包裹包装器，这个时候里层的包装器相对于外层的包装器充当流的角色

在php自身底层实现的c语言开发文档有这样的解释：
流API操作一对不同级别：在基本级别，api定义了php_stream对象表示流式数据源，在稍微高一点的级别，api定义了php_stream_wrapper对象
它包裹低一级别的php_stream对象，以提供取回URL的内容和元数据、添加上下文参数的能力，调整包装器行为；

每一种流打开后都可以应用任意数量的过滤器在上面，流数据会经过过滤器的处理，笔者认为过滤器这个词用得有点不准确，有些误导人
从字面意思看好像是去掉一些数据的感觉，应该称为数据调整器，因为它既可去掉一些数据，也可以添加，还可以修改，但历史原因约定俗成，
也就称为过滤器了，大家心里明白就好。

我们经常看到下面的词，来解释下他们的区别：
资源和数据：资源是比较宏观的说法，通常包含数据，而数据是比较具象的说法，在开发程序的时候经常说是数据，而在软件规划时说是资源，他们是近义词，就像软件设计和程序开发的区别一样。
上下文和参数：上下文是比较宏观的说法，经常用在沟通上面，具体点讲就是一次沟通本身的参数，而参数这个说法往往用在比较具体的事情上面，比如说函数

上面解释了概念性的东西，下面来看看具体内容：
php支持的协议和包装器请看这里：http://php.net/manual/zh/wrappers.php：
默认的支持了一些协议和包装器，请用stream_get_wrappers()函数查看.也可以自定义一个包装器，用stream_wrapper_register()注册
尽管RFC 3986里面可以使用:做分割符，但php只允许://，所以url请使用"scheme://target"这样的格式
    file:// — 访问本地文件系统，在用文件系统函数时默认就使用该包装器
    http:// — 访问 HTTP(s) 网址
    ftp:// — 访问 FTP(s) URLs
    php:// — 访问各个输入/输出流（I/O streams）
    zlib:// — 压缩流
    data:// — 数据（RFC 2397）
    glob:// — 查找匹配的文件路径模式
    phar:// — PHP 归档
    ssh2:// — Secure Shell 2
    rar:// — RAR
    ogg:// — 音频流
    expect:// — 处理交互式的流
如何实现一个自定义的包装器：

在用fopen、fwrite、fread、fgets、feof、rewind、file_put_contents、file_get_contents等等文件系统函数操作流时，数据是先传给定义的包装器类对象，包装器再去操作流。
如何实现一个自定义的流包装器呢？php提供了一个类原型，只是原型而已，不是接口也不是类，不能用于继承：
 streamWrapper {
/* 属性 */
public resource $context ;
/* 方法 */
__construct ( void )
__destruct ( void )
public bool dir_closedir ( void )
public bool dir_opendir ( string $path , int $options )
public string dir_readdir ( void )
public bool dir_rewinddir ( void )
public bool mkdir ( string $path , int $mode , int $options )
public bool rename ( string $path_from , string $path_to )
public bool rmdir ( string $path , int $options )
public resource stream_cast ( int $cast_as )
public void stream_close ( void )
public bool stream_eof ( void )
public bool stream_flush ( void )
public bool stream_lock ( int $operation )
public bool stream_metadata ( string $path , int $option , mixed $value )
public bool stream_open ( string $path , string $mode , int $options , string &$opened_path )
public string stream_read ( int $count )
public bool stream_seek ( int $offset , int $whence = SEEK_SET )
public bool stream_set_option ( int $option , int $arg1 , int $arg2 )
public array stream_stat ( void )
public int stream_tell ( void )
public bool stream_truncate ( int $new_size )
public int stream_write ( string $data )
public bool unlink ( string $path )
public array url_stat ( string $path , int $flags )
}

在这个原型里面定义的方法，根据自己需要去定义，并不要求全部实现，这就是为什么不定义成接口的原因，因为有些实现根本用不着某些方法，
这带来很多灵活性，比如包装器是不支持删除目录rmdir功能的，那么就不需要实现streamWrapper::rmdir
由于未实现它，如果用户在包装器上调用rmdir将有错误抛出，要自定义这个错误那么也可以实现它并在其内部抛出错误

streamWrapper也不是一个预定义类，测试class_exists("streamWrapper")就知道，它只是一个指导开发者的原型

官方手册提供了一个例子：http://php.net/manual/zh/stream.streamwrapper.example-1.php

本博客提供一个从drupal8系统中抽取修改过的包装器例子，请看drupal8源码分析关于流那一部分
流系列函数，官方手册：http://php.net/manual/zh/ref.stream.php
常用的函数如下：
stream_bucket_append函数：为队列添加数据　
stream_bucket_make_writeable函数：从操作的队列中返回一个数据对象
stream_bucket_new函数：为当前队列创建一个新的数据
stream_bucket_prepend函数：预备数据到队列　
stream_context_create函数：创建数据流上下文
stream_context_get_default函数：获取默认的数据流上下文
stream_context_get_options函数：获取数据流的设置
stream_context_set_option函数：对数据流、数据包或者上下文进行设置
stream_context_set_params函数：为数据流、数据包或者上下文设置参数
stream_copy_to_stream函数：在数据流之间进行复制操作
stream_filter_append函数：为数据流添加过滤器
stream_filter_prepend函数：为数据流预备添加过滤器
stream_filter_register函数：注册一个数据流的过滤器并作为PHP类执行
stream_filter_remove函数：从一个数据流中移除过滤器
stream_get_contents函数：读取数据流中的剩余数据到字符串
stream_get_filters函数：返回已经注册的数据流过滤器列表
stream_get_line函数：按照给定的定界符从数据流资源中获取行
stream_get_meta_data函数：从封装协议文件指针中获取报头/元数据
stream_get_transports函数：返回注册的Socket传输列表
stream_get_wrappers函数：返回注册的数据流列表
stream_register_wrapper函数：注册一个用PHP类实现的URL封装协议
stream_select函数：接收数据流数组并等待它们状态的改变
stream_set_blocking函数：将一个数据流设置为堵塞或者非堵塞状态
stream_set_timeout函数：对数据流进行超时设置
stream_set_write_buffer函数：为数据流设置缓冲区
stream_socket_accept函数：接受由函数stream_ socket_server()创建的Socket连接
stream_socket_client函数：打开网络或者UNIX主机的Socket连接
stream_socket_enable_crypto函数：为一个已经连接的Socket打开或者关闭数据加密
stream_socket_get_name函数：获取本地或者网络Socket的名称
stream_socket_pair函数：创建两个无区别的Socket数据流连接
stream_socket_recvfrom函数：从Socket获取数据，不管其连接与否
stream_socket_sendto函数：向Socket发送数据，不管其连接与否
stream_socket_server函数：创建一个网络或者UNIX Socket服务端
stream_wrapper_restore函数：恢复一个事先注销的数据包
stream_wrapper_unregister函数：注销一个URL地址包

一个过滤器的列子及解释：

相关链接：

用户过滤器基类：http://php.net/manual/zh/class.php-user-filter.php

过滤器注册：http://php.net/manual/zh/function.stream-filter-register.php

<?php
 
/* 定义一个过滤器 */
class strtoupper_filter extends php_user_filter {
  function filter($in, $out, &$consumed, $closing)
  {
    while ($bucket = stream_bucket_make_writeable($in)) { //从流里面取出一段数据
      $bucket->data = strtoupper($bucket->data);
      $consumed += $bucket->datalen;
      stream_bucket_append($out, $bucket); //将修改后的数据送到输出的地方
    }
    return PSFS_PASS_ON;
  }
}
 
/* 注册过滤器到php */
stream_filter_register("strtoupper", "strtoupper_filter")
    or die("Failed to register filter");
 
$fp = fopen("foo-bar.txt", "w");
 
/* 应用过滤器到一个流 */
stream_filter_append($fp, "strtoupper");
 
fwrite($fp, "Line1\n");
fwrite($fp, "Word - 2\n");
fwrite($fp, "Easy As 123\n");
 
fclose($fp);
 
//读取并显示内容 将全部变为大写
readfile("foo-bar.txt");
 
?>

<!-- more -->
