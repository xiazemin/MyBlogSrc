---
title: GPBMetadata  DescriptorPool
layout: post
category: php
author: 夏泽民
---
DescriptorPool
任何时候想要查询一个Descriptor , 都是去DescriptorPool里面查询。 
DescriptorPool 实现了这样的机制 ：

缓存所有查询的文件的Descriptor 。
查找Descriptor的时候，如果自身缓存查到就直接返回结果， 
否则去自带的DescriptorDatabase中查FileDescriptorProto， 
查到就转化成Descriptor， 返回结果并且缓存.

Descriptor – 用来描述 消息
FieldDescriptor – 用来描述 字段
OneofDescriptor – 用来描述 联合体
EnumDescriptor – 用来描述 枚举
EnumValueDescriptor – 用来描述 枚举值
ServiceDescriptor – 用来描述 服务器
MethodDescriptor – 用来描述 服务器方法
FileDescriptor – 用来描述 文件

https://blog.csdn.net/boshuzhang/article/details/66969353
<!-- more -->
descriptor pool

一个 protobuf 的 message 的 reflection API 是通过 descriptor 建立起来的，通常我们可以通过一个 message 的 static method：descriptor 拿到对应的 descriptor 继而进行一些工作。那么如果我们没有编译时候的类型信息，我们希望 runtime 通过一个 type erased 的类型获得相关的 descriptor 的话，可以依赖 descriptor pool，事实上 proto2::DescriptorPool 提供的 API 可以做非常多的事情，我们这里仅仅讨论一个对应的实现 DescriptorPool::generated_pool() 所返回的 pool 可以用来查询当前 binary 里面所链接的所有 message，这意味着我们可以通过一个 protobuf message 的 canonical name（字符串）获取对应的 Descriptor。

很显然如果我们需要将 message 类型本身通过 serialization/deserialization 的 boundary，字符串也要比 Descriptor* 容易很多。

options

protobuf 提供的一种 annotation 称为 option，它出现在 protobuf 的定义中，可以为几乎所有的东西提供 annotation：message 类型、field 甚至 enum value 等等。我们通过每个类型对应的 descriptor API 可以获得这些 annotation，这样我们就可以根据这些额外的 annotation 进行一些改变。

一个简单的例子：比如一个 query language 在处理某个 protobuf 的数据时，protobuf 本身其实作为一个 schema 可以将类型信息用在 query parsing 和 type verification 上，但是比如整数的 formatting 等问题并不能完全被 protobuf 本身涵盖，这更接近 query 这个 domain 希望如何 interpret 整数：uint64 是不是其实是一个 IP 地址，或者一个 timestamp 还是就是最好 hex 输出的 hash value

因此每个 domain 很可能希望有一套自己的 annotation，这样尽管生成数据的人尽管不关心 query，但是通过简单的 annotation 就可以在 query 自己数据的时候获得需要的格式了。

protobuf 的 option 也是通过 extension 来实现的。我们首先需要把 annotation 本身用一个 message type 来描述出来，然后在合适的 scope 中 extend protobuf 的 options，将自定义的 annotation type 以某种 extension 加入（如 message 是 proto2.MessageOptions）。这样我们就可以使用 option 关键字来 annotate 一些地方（field 和 value 的 annotation 语法稍微有点不同）。

除了以上这种用法，我们有时候还希望为某些类型关联一些 metadata，一种方法是通过 traits 在 C++ 代码里面提供，另一种就是当这些类型是 protobuf type 的时候我们也可以通过 options 来实现。很显然前者更为 generic，能够提供关联的方法，后者一般用相同的类型来存储 metadata，可以用 reflection/descriptor 的 API 来操作。

http://pages.cs.wisc.edu/~starr/bots/Undermind-src/html/classgoogle_1_1protobuf_1_1DescriptorPool.html
一、前提环境
    系统信息:Linux raspberrypi 4.9.24-v7+ #993 SMP Wed Apr 26 18:01:23 BST 2017 armv7l GNU/Linux

    系统已经预先安装php， 我是源码编译安装的php， 版本为5.5.38

    特别注意: 源码编译php的时候， 在配置中请加上 --enable-bcmath，protobuf php需要bcadd()函数支持。

   

二、安装protobuf
    1 先安装工具

    sudo apt-get install autoconf automake libtool curl make g++ unzip

    2 获取源码包

    https://github.com/google/protobuf           

    #这个编译会生成protoc可执行程序, 请参照https://github.com/google/protobuf/blob/master/src/README.md文档编译，在此就不再说明了

    https://github.com/allegro/php-protobuf    #这个编译php扩展， 即下文讲的源码

    3 将源码包上传到树莓派并解压到 /home/lighttpd/php-protobuf-master   (路径可以自己指定，后面的改成你自己的即可)
    4 进入源码路径

    cd  /home/lighttpd/php-protobuf-master

    5 将编译PHP生成的phpize可执行文件拷贝到当前路径, apt-get安装的也有这个文件，或者直接在终端检查有没有此命令

       若有就不需要执行此步骤

       如果是apt-get安装的php, 找不到phpize和php-config， 请先执行以下命令

       sudo apt-get install php5-dev

       然后查找php-config路径， 记下此路径，下一步备用

       sudo find / -name php-config

    6 进入源码路径php-protobuf-master, 执行以下命令
    ./phpize
    ./configure --with-php-config=/home/pi/webserver/php/build/x86/bin/php-config //此处路径为php安装路径下的php-config文件， 请改为自己的路径
    make
    make install
 
    7 看到以下内容即表示安装成功
     root@raspberrypi:/home/lighttpd/php-protobuf-master# make install
     Installing shared extensions:     /home/pi/webserver/php/build/x86/lib/php/extensions/no-debug-non-zts-20121212/
 
    8 查看版本:protoc --version

三，php加入protobuf扩展
    1 php开启模块支持， 打开php配置文件

    nano /home/pi/webserver/php/build/x86/lib/php.ini    //这是我的php.ini文件路径，请对照自己路径修改， nano和vim都是文本编辑器，也可以用vim
    加入此行，注意protobuf.so的路径即为前面第二节第7步安装成功后提示的路径
    extension=/home/pi/webserver/php/build/x86/lib/php/extensions/no-debug-non-zts-20121212/protobuf.so

    PS: 如果是apt-get安装的php， 

    sudo find / -name php.ini    #全局查找php.ini文件

    将找到的文件全部加上extension=.~~~~~/protobuf.so 这句话

    2 重启lighttpd服务器， 查看编写一个包含phpinfo()的脚本，浏览器访问， 检查有protobuf即OK

  



四，在php中使用protobuf 序列化数据(这部分是坑最多的)
    1   安装composer， 后面要用到这个管理工具， 网上查找的安装方式几乎都不能用， 绝大多数是因为被墙的原因。不过你要是能安装上，可直接略过下面步骤。
         直接下载composer.phar，可以在CSDN下载， 或者中文官网也可以下载，但请下载较新版本的

         附CSDN连接:http://download.csdn.net/detail/webben/9787171   点击 CSDN下载

    2  下载完成后上传到树莓派任意路径

    3   执行以下命令，直接全局使用，注意执行之前给composer.phar附加执行权限

         cp composer.phar /usr/local/bin/composer
 

    4  执行以下命令 修改composer源到中国镜像，并关闭http安全验证，然后安装

        cd /home/lighttpd/php-protobuf-master  //首先进入到之前安装protobuf的路径, 否则最后执行composer install会有问题，请对照自己路径操作

        composer config -g repo.packagist composer https://packagist.phpcomposer.com  
        composer config secure-http false
        composer install
 
    5  更新 composer , 
        cd /home/lighttpd/php-protobuf-master/  

       编辑composer.json，在require字段处加入如下行"google/protobuf": "^3.2"     //3.2 是使用protoc --version命令看到的版本号,请对应操作,
       {
           "require": {
           "google/protobuf": "^3.2"
       }

       如果忽略上述步骤， 会造成php无法找到对应的google头文件

       新建一个目录， (又是一个坑)

       mkdir /home/lighttpd/php-protobuf-master/vendor/google/protobuf

       composer update
       composer install    //重新安装
  

    6  拷贝步骤5生成的文件夹vendor 到网站php脚本目录， 或者任意php脚本能访问到的目录都行
        cp -r  vendor  /var/www/htdocs/php/



    7 新建一个简单的proto文件, 路径可以任意指定， 后续的换成你自己的即可

       nano /var/www/htdocs/php/protobuf/test.proto

       写入以下内容,保存退出

       syntax = "proto3";
       package config;
       message VoiceConfig{
          int32 sample_rate = 1;
          int32 mic_num = 2;
          string voiceserver_address = 3;
       }

     8 执行下面命令，

       /home/pi/webserver/php/build/x86/bin/php /home/lighttpd/php-protobuf-master/protoc-gen-php.php test.proto

      其中/home/pi/webserver/php/build/x86/bin/php是php可执行文件路径，如果有全局php命令直接使用php即可

      /home/lighttpd/php-protobuf-master/protoc-gen-php.php是第一节安装protobuf路径下的文件，请换成你的路径

     9 新建PHP脚本test.php 内容如下

 <?php
    error_reporting(E_ERROR | E_WARNING | E_PARSE);  //禁止打印错误和警告信息
    require_once('vendor/autoload.php');   //步骤四.6生成的文件路径内有此文件
    require_once('protobuf/GPBMetadata/Config.php');  //步骤四.8 生成的文件路径内有此文件
 
    include 'protobuf/Config/VoiceConfig.php';   //步骤四.8 生成的文件路径内有此文件
 
    $foo = new \Config\VoiceConfig();
    $foo->setSampleRate(1600);
    $foo->setMicNum(8);
    $foo->setVoiceserverAddress("192.168.0.136");
 
    //to encode message
    $packed = $foo->serializeToString();
 
    //to decode message
    $param2 = new \Config\VoiceConfig();
    $param2->mergeFromString($packed);
 
 
    $jsonArr = array("getSampleRate"=> $param2->getSampleRate(),
        "setMicNum"=> $param2->getMicNum(),
        "setVoiceserverAddress"=> $param2->getVoiceserverAddress(),
        );
    echo json_encode($jsonArr);
?> 


10 通过浏览器访问test.php， 看到如下打印即OK

https://blog.csdn.net/lb1885727/article/details/74949883

https://my.oschina.net/hongjiang/blog/3111240
https://www.jianshu.com/p/1d550bb8509d

https://blog.csdn.net/panjican/article/details/97495326

https://learnku.com/articles/36277

http://www.suoniao.com/article/43760

https://segmentfault.com/a/1190000009389032

https://studygolang.com/articles/21667?fr=sidebar

https://www.bookstack.cn/read/hyperf/0fb5b9e27e11b5e6.md

https://www.yuanmas.com/info/4py2Bm8Kyb.html

http://www.suoniao.com/article/43760
https://blog.csdn.net/guyan0319/article/details/80613846?utm_source=blogxgwz8
https://segmentfault.com/a/1190000019234926

https://studygolang.com/articles/21709?fr=sidebar

https://www.jianshu.com/p/7392406e2450

https://segmentfault.com/a/1190000019688457?utm_source=tag-newest

https://www.jianshu.com/p/7b12fa3ca8e3

https://github.com/protocolbuffers/protobuf/issues/2971

https://github.com/protocolbuffers/protobuf/issues/3970
https://chromium.googlesource.com/chromium/src/third_party/+/master/protobuf/php/src/GPBMetadata/Google/Protobuf/Struct.php?autodive=0%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F


https://www.cnblogs.com/52fhy/p/11106670.html

BCMath 任意精度数学 ¶

https://www.php.net/manual/zh/book.bc.php

  PHP的linux版本需要手动安装BCMath扩展，在PHP的源码包中默认包含BCMath的安装文件，只需手动安装一次即可。
        编译安装
        1.进入PHP源码包目录下的ext/bcmath目录。
        2.执行phpize命令，phpize命令在PHP安装目录的bin目录下，如/usr/local/php-5.6.36/bin/phpize。
        3.执行./configure --with-php-config=/usr/local/php-5.6.36/bin/php-config。
        4.执行make && make install。
        5.将安装完成后得到bcmath.so文件拷贝到php.ini中extension_dir配置的目录中。
        6.在Dynamic Extensions配置块下添加一行extension=bcmath.so。
        7.重启php服务即可。

        在线安装
        1.安装BCMath，yum install php-bcmath。
        2.重启httpd，httpd -k restart。
        

